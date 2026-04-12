package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func runAudit(args []string) error {
	srcDir := "."
	docsDir := ""

	if len(args) >= 1 {
		srcDir = args[0]
	}
	if len(args) >= 2 {
		docsDir = args[1]
	}

	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return fmt.Errorf("resolving source dir: %w", err)
	}

	// Auto-detect docs dir
	if docsDir == "" {
		docsDir = detectDocsDir(srcDir)
	} else {
		docsDir, err = filepath.Abs(docsDir)
		if err != nil {
			return fmt.Errorf("resolving docs dir: %w", err)
		}
	}

	// h2pp manages state and exports $H2PP_LAST_AUDIT
	lastAudit := os.Getenv("H2PP_LAST_AUDIT")
	auditKind := "full"
	if lastAudit != "" {
		auditKind = "incremental (last: " + lastAudit + ")"
	}

	// Load mapping file
	mappingFile := filepath.Join(srcDir, ".h2", "mappings", "doc-pages.yaml")
	mappings, mappingErr := loadMappings(mappingFile)
	hasMappings := mappingErr == nil && len(mappings) > 0

	// Discover packages
	packages := discoverPackages(srcDir)

	// Discover doc pages on disk
	var docPages map[string]bool
	if docsDir != "" {
		docPages = discoverDocPages(docsDir)
	} else {
		docPages = make(map[string]bool)
	}

	// Begin output
	fmt.Println("# Audit Inventory")
	fmt.Printf("Audit: %s\n", auditKind)
	if hasMappings {
		fmt.Printf("Mapping: loaded %d package entries from doc-pages.yaml\n", len(mappings))
	} else {
		fmt.Println("Mapping: doc-pages.yaml not found, skipping gap analysis")
	}
	fmt.Println()

	var totalIssues int

	// --- Gap Analysis ---
	fmt.Println("## Gap Analysis")
	fmt.Println()

	type staleEntry struct {
		pkg     string
		page    string
		srcTime time.Time
		docTime time.Time
	}

	var missing []string
	var stale []staleEntry
	var orphaned []string
	problemPkgs := make(map[string]bool)

	if hasMappings && docsDir != "" {
		// Track which doc pages are referenced
		referenced := make(map[string]bool)

		for pkg, pages := range mappings {
			for _, page := range pages {
				referenced[page] = true
				docPath := filepath.Join(docsDir, page)
				if _, statErr := os.Stat(docPath); os.IsNotExist(statErr) {
					missing = append(missing, fmt.Sprintf("- %s → %s", pkg, page))
					problemPkgs[pkg] = true
				} else if statErr == nil {
					// Check staleness
					docInfo, _ := os.Stat(docPath)
					srcNewest := newestGoFile(filepath.Join(srcDir, pkg))
					if !srcNewest.IsZero() && srcNewest.After(docInfo.ModTime()) {
						stale = append(stale, staleEntry{
							pkg:     pkg,
							page:    page,
							srcTime: srcNewest,
							docTime: docInfo.ModTime(),
						})
						problemPkgs[pkg] = true
					}
				}
			}
		}

		// Orphaned: doc pages not in any mapping
		for page := range docPages {
			if !referenced[page] {
				orphaned = append(orphaned, page)
			}
		}

		sort.Strings(missing)
		sort.Strings(orphaned)
		sort.Slice(stale, func(i, j int) bool {
			return stale[i].pkg < stale[j].pkg
		})
	}

	fmt.Println("### Missing Doc Pages")
	if len(missing) == 0 {
		fmt.Println("None")
	} else {
		for _, m := range missing {
			fmt.Println(m)
		}
		totalIssues += len(missing)
	}
	fmt.Println()

	fmt.Println("### Stale Doc Pages")
	if len(stale) == 0 {
		fmt.Println("None")
	} else {
		for _, s := range stale {
			fmt.Printf("- %s → %s (src: %s, doc: %s)\n", s.pkg, s.page,
				s.srcTime.Format(time.RFC3339), s.docTime.Format(time.RFC3339))
		}
		totalIssues += len(stale)
	}
	fmt.Println()

	fmt.Println("### Orphaned Doc Pages")
	if len(orphaned) == 0 {
		fmt.Println("None")
	} else {
		for _, o := range orphaned {
			fmt.Printf("- %s\n", o)
		}
		totalIssues += len(orphaned)
	}
	fmt.Println()

	// --- Exports for Problem Packages ---
	fmt.Println("## Exports for Problem Packages")
	fmt.Println()

	if len(problemPkgs) == 0 {
		fmt.Println("No problem packages found.")
	} else {
		sortedPkgs := make([]string, 0, len(problemPkgs))
		for pkg := range problemPkgs {
			sortedPkgs = append(sortedPkgs, pkg)
		}
		sort.Strings(sortedPkgs)

		exportRe := regexp.MustCompile(`^(type|func|var|const) [A-Z]`)

		for _, pkg := range sortedPkgs {
			fmt.Printf("### %s\n", pkg)
			exports := collectExports(filepath.Join(srcDir, pkg), exportRe)
			if len(exports) == 0 {
				fmt.Println("No exports found.")
			} else {
				for _, e := range exports {
					fmt.Printf("- %s\n", e)
				}
			}
			fmt.Println()
		}
	}

	// --- Code Example Verification ---
	fmt.Println("## Code Example Verification")
	fmt.Println()

	var totalBlocks, verifiedCount, unverifiedCount, staleCount int
	type blockInfo struct {
		file string
		line int
	}
	var unverified []blockInfo
	var staleBlocks []blockInfo

	siteDocsRoot := filepath.Join(srcDir, "site", "src", "content", "docs")
	if info, statErr := os.Stat(siteDocsRoot); statErr == nil && info.IsDir() {
		_ = filepath.Walk(siteDocsRoot, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || info.IsDir() {
				return walkErr
			}
			if !strings.HasSuffix(path, ".mdx") {
				return nil
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

			for scanner.Scan() {
				lineNum++
				line := scanner.Text()

				if strings.Contains(line, "```go") {
					totalBlocks++
					verifiedPath := parseVerifiedComment(prevLine)
					if verifiedPath == "" {
						unverifiedCount++
						unverified = append(unverified, blockInfo{file: relPath, line: lineNum})
					} else {
						verifiedCount++
						// Strip optional ":line" suffix before stat (e.g. "file.go:15" → "file.go")
						statPath := verifiedPath
						if idx := strings.LastIndex(statPath, ":"); idx >= 0 {
							if _, err := fmt.Sscanf(statPath[idx+1:], "%d", new(int)); err == nil {
								statPath = statPath[:idx]
							}
						}
						absVerified := filepath.Join(srcDir, statPath)
						if _, statErr2 := os.Stat(absVerified); os.IsNotExist(statErr2) {
							staleCount++
							staleBlocks = append(staleBlocks, blockInfo{file: relPath, line: lineNum})
						} else {
							// Check if source is newer than doc
							srcInfo, _ := os.Stat(absVerified)
							if srcInfo != nil && srcInfo.ModTime().After(info.ModTime()) {
								staleCount++
								staleBlocks = append(staleBlocks, blockInfo{file: relPath, line: lineNum})
							}
						}
					}
				}
				prevLine = line
			}
			return nil
		})
	}

	fmt.Printf("%d total: %d verified, %d unverified, %d stale\n",
		totalBlocks, verifiedCount, unverifiedCount, staleCount)
	fmt.Println()

	fmt.Println("### Unverified Code Blocks")
	if len(unverified) == 0 {
		fmt.Println("None")
	} else {
		for _, u := range unverified {
			fmt.Printf("- %s:%d\n", u.file, u.line)
		}
	}
	totalIssues += unverifiedCount
	fmt.Println()

	fmt.Println("### Stale Verified Blocks")
	if len(staleBlocks) == 0 {
		fmt.Println("None")
	} else {
		for _, s := range staleBlocks {
			fmt.Printf("- %s:%d\n", s.file, s.line)
		}
	}
	totalIssues += staleCount
	fmt.Println()

	// Summary
	fmt.Println("---")
	fmt.Printf("%d packages, %d doc pages, %d issues\n", len(packages), len(docPages), totalIssues)

	return nil
}

// loadMappings reads doc-pages.yaml. Values can be a single string or a list of strings.
func loadMappings(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]interface{})
	if err = yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	result := make(map[string][]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			result[k] = []string{val}
		case []interface{}:
			pages := make([]string, 0, len(val))
			for _, item := range val {
				if s, ok := item.(string); ok {
					pages = append(pages, s)
				}
			}
			result[k] = pages
		}
	}
	return result, nil
}

// detectDocsDir tries common locations for docs components directory.
func detectDocsDir(srcDir string) string {
	candidates := []string{
		filepath.Join(srcDir, "site", "src", "content", "docs", "components"),
		filepath.Join(srcDir, "docs", "components"),
		filepath.Join(srcDir, "site", "docs", "components"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}

// discoverPackages finds tea* directories containing .go files.
func discoverPackages(srcDir string) []string {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil
	}
	var pkgs []string
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "tea") {
			continue
		}
		goFiles, _ := filepath.Glob(filepath.Join(srcDir, e.Name(), "*.go"))
		if len(goFiles) > 0 {
			pkgs = append(pkgs, e.Name())
		}
	}
	sort.Strings(pkgs)
	return pkgs
}

// discoverDocPages finds .mdx and .md files in the docs directory.
func discoverDocPages(docsDir string) map[string]bool {
	pages := make(map[string]bool)
	entries, err := os.ReadDir(docsDir)
	if err != nil {
		return pages
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".mdx") || strings.HasSuffix(name, ".md") {
			pages[name] = true
		}
	}
	return pages
}

// newestGoFile returns the most recent ModTime of .go files (excluding _test.go) in a dir.
func newestGoFile(dir string) time.Time {
	var newest time.Time
	entries, err := os.ReadDir(dir)
	if err != nil {
		return newest
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}
	return newest
}

// collectExports reads .go files in a package dir and returns sorted exported declarations.
func collectExports(dir string, re *regexp.Regexp) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	seen := make(map[string]bool)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, openErr := os.Open(filepath.Join(dir, name))
		if openErr != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				// Trim to a reasonable declaration summary
				trimmed := strings.TrimSpace(line)
				seen[trimmed] = true
			}
		}
		f.Close()
	}
	result := make([]string, 0, len(seen))
	for line := range seen {
		result = append(result, line)
	}
	sort.Strings(result)
	return result
}

// parseVerifiedComment extracts the path from a "<!-- verified: path/to/file -->" comment.
func parseVerifiedComment(line string) string {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "<!-- verified:") {
		return ""
	}
	trimmed = strings.TrimPrefix(trimmed, "<!-- verified:")
	trimmed = strings.TrimSuffix(trimmed, "-->")
	trimmed = strings.TrimSpace(trimmed)
	return trimmed
}
