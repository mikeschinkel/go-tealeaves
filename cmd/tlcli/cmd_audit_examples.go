package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// mdxBlock represents a ```go code block in an MDX file that has a verified comment.
type mdxBlock struct {
	mdxRelPath   string // e.g. "site/src/content/docs/components/confirm-dialog.mdx"
	line         int    // line number of the ```go fence
	examplePath  string // from verified comment, e.g. "site/examples/confirm-dialog/compile_test.go"
	exampleLine  int    // line number in the example file, e.g. 12
	contentHash  uint32 // CRC32 of the code block content
	contentLines []string
}

// sourceEntry represents one line:hash pair parsed from a // Source: comment.
type sourceEntry struct {
	line int
	hash uint32 // 0 if no hash present
}

// sourceRef represents a parsed // Source: comment from an example file.
type sourceRef struct {
	mdxRelPath string        // e.g. "site/src/content/docs/components/confirm-dialog.mdx"
	entries    []sourceEntry // line numbers with optional hashes
}

// sourceFixup describes one fix to apply to a Source comment entry.
type sourceFixup struct {
	examplePath string
	oldLine     int    // current line number in Source comment
	newLine     int    // correct line number (0 = keep oldLine)
	newHash     uint32 // hash to write (0 = keep existing)
}

func runAuditExamples(args []string) error {
	var fixFlag bool
	var remaining []string
	for _, a := range args {
		if a == "--fix" || a == "-fix" {
			fixFlag = true
		} else {
			remaining = append(remaining, a)
		}
	}

	srcDir := "."
	if len(remaining) >= 1 {
		srcDir = remaining[0]
	}

	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return fmt.Errorf("resolving source dir: %w", err)
	}

	// Phase 1: Parse MDX blocks with verified comments
	blocks := parseMDXBlocks(srcDir)

	// Phase 2: Parse Source comments from example files
	sources := parseSourceComments(srcDir)

	// Build set of example files referenced by MDX blocks
	referencedExamples := make(map[string]bool)

	// Categorize blocks
	type needsExample struct {
		mdxPath string
		line    int
		detail  string
	}
	type contentChanged struct {
		mdxPath     string
		line        int
		examplePath string
	}
	type lineStale struct {
		examplePath string
		oldLine     int
		newLine     int
		hasHash     bool // whether Source entry already had a hash
	}

	var (
		needsExampleList   []needsExample
		contentChangedList []contentChanged
		lineStaleList      []lineStale
		needsHashSeedList  []sourceFixup // exact line match but no hash
		matchedCount       int
	)

	for _, block := range blocks {
		referencedExamples[block.examplePath] = true

		// 1. File exists?
		absExample := filepath.Join(srcDir, block.examplePath)
		if _, statErr := os.Stat(absExample); os.IsNotExist(statErr) {
			needsExampleList = append(needsExampleList, needsExample{
				mdxPath: block.mdxRelPath,
				line:    block.line,
				detail:  fmt.Sprintf("referenced %s (not found)", block.examplePath),
			})
			continue
		}

		// 2. Bidirectional check — find Source comment for this example
		ref, hasSource := sources[block.examplePath]
		if !hasSource {
			needsExampleList = append(needsExampleList, needsExample{
				mdxPath: block.mdxRelPath,
				line:    block.line,
				detail:  fmt.Sprintf("no // Source: comment in %s", block.examplePath),
			})
			continue
		}

		// Check that the Source comment references the same MDX file
		if ref.mdxRelPath != block.mdxRelPath {
			needsExampleList = append(needsExampleList, needsExample{
				mdxPath: block.mdxRelPath,
				line:    block.line,
				detail:  fmt.Sprintf("Source comment references %s, not this file", ref.mdxRelPath),
			})
			continue
		}

		// 3. Find the matching entry by line number, hash, or nearby line
		var foundEntry *sourceEntry
		var matchKind string // "exact", "hash-shifted", "nearby"

		// Try exact line match
		for i := range ref.entries {
			if ref.entries[i].line == block.line {
				foundEntry = &ref.entries[i]
				matchKind = "exact"
				break
			}
		}

		// Try hash match (line shifted, hash preserved)
		if foundEntry == nil && block.contentHash != 0 {
			for i := range ref.entries {
				if ref.entries[i].hash != 0 && ref.entries[i].hash == block.contentHash {
					foundEntry = &ref.entries[i]
					matchKind = "hash-shifted"
					break
				}
			}
		}

		// Try nearby line match (±3) for entries without hashes
		if foundEntry == nil {
			for delta := 1; delta <= 3; delta++ {
				for _, tryLine := range []int{block.line - delta, block.line + delta} {
					for i := range ref.entries {
						if ref.entries[i].line == tryLine {
							foundEntry = &ref.entries[i]
							matchKind = "nearby"
							break
						}
					}
					if foundEntry != nil {
						break
					}
				}
				if foundEntry != nil {
					break
				}
			}
		}

		if foundEntry == nil {
			needsExampleList = append(needsExampleList, needsExample{
				mdxPath: block.mdxRelPath,
				line:    block.line,
				detail:  fmt.Sprintf("Source comment in %s doesn't reference line %d (or nearby)", block.examplePath, block.line),
			})
			continue
		}

		// 4. Categorize based on match kind
		switch matchKind {
		case "exact":
			if foundEntry.hash == 0 {
				// Line matches but no hash — needs seeding
				needsHashSeedList = append(needsHashSeedList, sourceFixup{
					examplePath: block.examplePath,
					oldLine:     foundEntry.line,
					newHash:     block.contentHash,
				})
			} else if foundEntry.hash != block.contentHash {
				// Line matches but hash differs — content changed
				contentChangedList = append(contentChangedList, contentChanged{
					mdxPath:     block.mdxRelPath,
					line:        block.line,
					examplePath: block.examplePath,
				})
				continue
			}
			matchedCount++

		case "hash-shifted":
			// Hash matches but line shifted
			lineStaleList = append(lineStaleList, lineStale{
				examplePath: block.examplePath,
				oldLine:     foundEntry.line,
				newLine:     block.line,
				hasHash:     true,
			})
			matchedCount++

		case "nearby":
			// Nearby line match — line shifted, no hash to confirm
			lineStaleList = append(lineStaleList, lineStale{
				examplePath: block.examplePath,
				oldLine:     foundEntry.line,
				newLine:     block.line,
				hasHash:     foundEntry.hash != 0,
			})
			matchedCount++
		}
	}

	// Find orphaned example files
	var orphanedList []string
	for exPath := range sources {
		if !referencedExamples[exPath] {
			orphanedList = append(orphanedList, exPath)
		}
	}
	sort.Strings(orphanedList)

	// Run go build on examples module
	examplesDir := filepath.Join(srcDir, "site", "examples")
	buildOutput, buildErr := runExamplesBuild(examplesDir)

	// --- Output ---
	fmt.Println("## Examples Audit")
	fmt.Println()

	fmt.Printf("### Needs Example Written (%d)\n", len(needsExampleList))
	if len(needsExampleList) == 0 {
		fmt.Println("None")
	} else {
		for _, e := range needsExampleList {
			if e.detail != "" {
				fmt.Printf("- %s:%d → %s\n", e.mdxPath, e.line, e.detail)
			} else {
				fmt.Printf("- %s:%d\n", e.mdxPath, e.line)
			}
		}
	}
	fmt.Println()

	fmt.Printf("### Content Changed (%d)\n", len(contentChangedList))
	if len(contentChangedList) == 0 {
		fmt.Println("None")
	} else {
		for _, c := range contentChangedList {
			fmt.Printf("- %s:%d → hash mismatch (MDX block was edited since last verification)\n", c.mdxPath, c.line)
		}
	}
	fmt.Println()

	fmt.Printf("### Orphaned Example File (%d)\n", len(orphanedList))
	if len(orphanedList) == 0 {
		fmt.Println("None")
	} else {
		for _, o := range orphanedList {
			fmt.Printf("- %s → no MDX block references this file (delete it)\n", o)
		}
	}
	fmt.Println()

	fmt.Printf("### Line Numbers Stale (%d)\n", len(lineStaleList))
	if len(lineStaleList) == 0 {
		fmt.Println("None")
	} else {
		for _, s := range lineStaleList {
			note := "content matches, line shifted"
			if !s.hasHash {
				note = "line shifted, no hash"
			}
			fmt.Printf("- %s → line %d→%d (%s)\n", s.examplePath, s.oldLine, s.newLine, note)
		}
	}
	fmt.Println()

	// Deduplicate hash seed list by file (only report each file once)
	hashSeedFiles := make(map[string]bool)
	for _, f := range needsHashSeedList {
		hashSeedFiles[f.examplePath] = true
	}
	fmt.Printf("### Needs Hash Seeding (%d files)\n", len(hashSeedFiles))
	if len(hashSeedFiles) == 0 {
		fmt.Println("None")
	} else {
		sortedFiles := make([]string, 0, len(hashSeedFiles))
		for f := range hashSeedFiles {
			sortedFiles = append(sortedFiles, f)
		}
		sort.Strings(sortedFiles)
		for _, f := range sortedFiles {
			fmt.Printf("- %s\n", f)
		}
	}
	fmt.Println()

	fmt.Printf("### Example Broken — Compile Failure (%d)\n", boolToInt(buildErr != nil))
	if buildErr != nil {
		fmt.Println(buildOutput)
	} else {
		fmt.Println("None")
	}
	fmt.Println()

	fmt.Printf("### Matched and Compiling (%d)\n", matchedCount)
	fmt.Println()

	// --- Apply fixes ---
	fixedCount := 0
	if fixFlag {
		// Build unified fixup list from all fixable categories
		var allFixups []sourceFixup

		// Line stale entries: update line number, seed hash if missing
		for _, s := range lineStaleList {
			fx := sourceFixup{
				examplePath: s.examplePath,
				oldLine:     s.oldLine,
				newLine:     s.newLine,
			}
			// Find the MDX block hash for seeding
			for _, b := range blocks {
				if b.examplePath == s.examplePath && b.line == s.newLine {
					fx.newHash = b.contentHash
					break
				}
			}
			allFixups = append(allFixups, fx)
		}

		// Hash seed entries: keep line, add hash
		allFixups = append(allFixups, needsHashSeedList...)

		if len(allFixups) > 0 {
			fmt.Println("### Fixing Source Comments")
			fixedCount = applySourceFixups(srcDir, allFixups)
			fmt.Println()
		}
	}

	// Summary
	fixableCount := len(lineStaleList) + len(hashSeedFiles)
	unfixableCount := len(needsExampleList) + len(contentChangedList) + len(orphanedList) + boolToInt(buildErr != nil)
	totalIssues := unfixableCount + fixableCount
	if fixedCount > 0 {
		fmt.Println("---")
		fmt.Printf("%d fixed, %d remaining, %d matched and compiling\n", fixedCount, totalIssues-fixedCount, matchedCount)
	} else {
		fmt.Println("---")
		if fixableCount > 0 {
			fmt.Printf("%d issues (%d fixable with --fix), %d matched and compiling\n", totalIssues, fixableCount, matchedCount)
		} else {
			fmt.Printf("%d issues, %d matched and compiling\n", totalIssues, matchedCount)
		}
	}

	return nil
}

// parseMDXBlocks walks MDX files and extracts ```go blocks with verified comments.
func parseMDXBlocks(srcDir string) []mdxBlock {
	siteDocsRoot := filepath.Join(srcDir, "site", "src", "content", "docs")
	var blocks []mdxBlock

	_ = filepath.Walk(siteDocsRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || !strings.HasSuffix(path, ".mdx") {
			return walkErr
		}

		f, openErr := os.Open(path)
		if openErr != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		prevLine := ""

		relPath, _ := filepath.Rel(srcDir, path)

		var inCodeBlock bool
		var codeLines []string

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if inCodeBlock {
				if strings.HasPrefix(strings.TrimSpace(line), "```") {
					// End of code block — compute hash
					content := strings.Join(codeLines, "\n")
					content = strings.TrimSpace(content)
					h := crc32.ChecksumIEEE([]byte(content))

					// Set hash on the block we started
					if len(blocks) > 0 {
						last := &blocks[len(blocks)-1]
						last.contentHash = h
						last.contentLines = codeLines
					}
					inCodeBlock = false
					codeLines = nil
				} else {
					codeLines = append(codeLines, line)
				}
				prevLine = line
				continue
			}

			if strings.Contains(line, "```go") {
				verifiedPath := parseVerifiedComment(prevLine)
				if verifiedPath != "" {
					exPath := verifiedPath
					exLine := 0
					if idx := strings.LastIndex(exPath, ":"); idx >= 0 {
						if n, parseErr := strconv.Atoi(exPath[idx+1:]); parseErr == nil {
							exLine = n
							exPath = exPath[:idx]
						}
					}

					blocks = append(blocks, mdxBlock{
						mdxRelPath:  relPath,
						line:        lineNum,
						examplePath: exPath,
						exampleLine: exLine,
					})
					inCodeBlock = true
				}
			}
			prevLine = line
		}
		return nil
	})

	return blocks
}

// parseSourceComments walks example .go files and parses // Source: comments.
func parseSourceComments(srcDir string) map[string]sourceRef {
	examplesRoot := filepath.Join(srcDir, "site", "examples")
	refs := make(map[string]sourceRef)

	_ = filepath.Walk(examplesRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return walkErr
		}

		f, openErr := os.Open(path)
		if openErr != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "// Source:") {
				continue
			}

			rest := strings.TrimSpace(strings.TrimPrefix(line, "// Source:"))
			mdxPath, lineSpec := splitSourceRef(rest)
			if mdxPath == "" {
				continue
			}

			entries := parseLineSpec(lineSpec)
			relPath, _ := filepath.Rel(srcDir, path)
			refs[relPath] = sourceRef{
				mdxRelPath: mdxPath,
				entries:    entries,
			}
			break
		}
		return nil
	})

	return refs
}

// splitSourceRef splits "site/.../file.mdx:21#hash,50#hash" into path and line spec.
func splitSourceRef(s string) (string, string) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9' {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}

// parseLineSpec parses "21#a1b2c3f4,50#d4e5f678,205" into sourceEntry slices.
func parseLineSpec(spec string) []sourceEntry {
	if spec == "" {
		return nil
	}
	parts := strings.Split(spec, ",")
	entries := make([]sourceEntry, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var e sourceEntry
		if lineStr, hashStr, ok := strings.Cut(p, "#"); ok {
			n, parseErr := strconv.Atoi(lineStr)
			if parseErr != nil {
				continue
			}
			e.line = n
			h, parseErr := strconv.ParseUint(hashStr, 16, 32)
			if parseErr == nil {
				e.hash = uint32(h)
			}
		} else {
			n, parseErr := strconv.Atoi(p)
			if parseErr != nil {
				continue
			}
			e.line = n
		}
		entries = append(entries, e)
	}
	return entries
}

// runExamplesBuild runs "go build ./..." in the examples directory.
func runExamplesBuild(examplesDir string) (string, error) {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = examplesDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// applySourceFixups rewrites // Source: comments in example files to update
// line numbers and/or seed content hashes.
func applySourceFixups(srcDir string, fixups []sourceFixup) int {
	// Group by file
	byFile := make(map[string][]sourceFixup)
	for _, f := range fixups {
		byFile[f.examplePath] = append(byFile[f.examplePath], f)
	}

	fixed := 0
	for relPath, fileFixes := range byFile {
		absPath := filepath.Join(srcDir, relPath)
		data, err := os.ReadFile(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  warning: cannot read %s: %v\n", relPath, err)
			continue
		}

		content := string(data)
		lines := strings.SplitAfter(content, "\n")
		changed := false

		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "// Source:") {
				continue
			}

			newLine := line
			for _, fx := range fileFixes {
				newLine = applyOneFixup(newLine, fx)
			}

			if newLine != line {
				lines[i] = newLine
				changed = true
			}
		}

		if changed {
			result := strings.Join(lines, "")
			if err := os.WriteFile(absPath, []byte(result), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "  warning: cannot write %s: %v\n", relPath, err)
				continue
			}
			for _, fx := range fileFixes {
				if fx.newLine != 0 && fx.newLine != fx.oldLine {
					fmt.Printf("  %s: line %d → %d", relPath, fx.oldLine, fx.newLine)
					if fx.newHash != 0 {
						fmt.Printf(" + hash seeded")
					}
					fmt.Println()
				} else if fx.newHash != 0 {
					fmt.Printf("  %s: hash seeded for line %d\n", relPath, fx.oldLine)
				}
				fixed++
			}
		}
	}
	return fixed
}

// applyOneFixup modifies a single line entry in a Source comment.
func applyOneFixup(line string, fx sourceFixup) string {
	idx := strings.Index(line, "// Source:")
	if idx < 0 {
		return line
	}
	prefix := line[:idx+len("// Source:")]
	rest := line[idx+len("// Source:"):]

	// Find the colon that starts the line spec
	colonIdx := -1
	for i := len(rest) - 1; i >= 0; i-- {
		if rest[i] == ':' && i+1 < len(rest) && rest[i+1] >= '0' && rest[i+1] <= '9' {
			colonIdx = i
			break
		}
	}
	if colonIdx < 0 {
		return line
	}

	beforeSpec := rest[:colonIdx+1]
	spec := rest[colonIdx+1:]

	oldLineStr := strconv.Itoa(fx.oldLine)
	newLineStr := oldLineStr
	if fx.newLine != 0 {
		newLineStr = strconv.Itoa(fx.newLine)
	}

	parts := strings.Split(spec, ",")
	for i, p := range parts {
		// Preserve leading/trailing whitespace and newlines
		core := strings.TrimSpace(p)
		if core == "" {
			continue
		}
		leadIdx := strings.Index(p, core)
		lead := p[:leadIdx]
		trail := p[leadIdx+len(core):]

		numPart := core
		hashPart := ""
		if hashIdx := strings.Index(core, "#"); hashIdx >= 0 {
			numPart = core[:hashIdx]
			hashPart = core[hashIdx:]
		}

		if numPart != oldLineStr {
			continue
		}

		// Update line number
		newCore := newLineStr

		// Update or seed hash
		if fx.newHash != 0 {
			newCore += fmt.Sprintf("#%08x", fx.newHash)
		} else if hashPart != "" {
			newCore += hashPart
		}

		parts[i] = lead + newCore + trail
	}

	return prefix + beforeSpec + strings.Join(parts, ",")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
