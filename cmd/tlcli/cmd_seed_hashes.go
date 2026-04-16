package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func runSeedHashes(args []string) error {
	srcDir := "."
	if len(args) >= 1 {
		srcDir = args[0]
	}

	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return fmt.Errorf("resolving source dir: %w", err)
	}

	// Parse all MDX blocks to get content hashes keyed by mdxRelPath+line
	blocks := parseMDXBlocks(srcDir)

	// Build lookup: mdxRelPath -> line -> contentHash
	type hashKey struct {
		mdxPath string
		line    int
	}
	hashMap := make(map[hashKey]uint32)
	for _, b := range blocks {
		hashMap[hashKey{mdxPath: b.mdxRelPath, line: b.line}] = b.contentHash
	}

	// Walk example files and rewrite Source comments with hashes
	sources := parseSourceComments(srcDir)

	var updated, skipped, alreadyHashed int

	for exRelPath, ref := range sources {
		// Check if all entries already have hashes
		allHaveHash := true
		for _, e := range ref.entries {
			if e.hash == 0 {
				allHaveHash = false
				break
			}
		}
		if allHaveHash {
			alreadyHashed++
			continue
		}

		// Build new entries with hashes
		newEntries := make([]string, 0, len(ref.entries))
		anyChange := false
		for _, e := range ref.entries {
			if e.hash != 0 {
				// Already has a hash, keep it
				newEntries = append(newEntries, fmt.Sprintf("%d#%08x", e.line, e.hash))
				continue
			}
			// Look up the hash from the MDX block
			h, ok := hashMap[hashKey{mdxPath: ref.mdxRelPath, line: e.line}]
			if !ok || h == 0 {
				// Can't find the block — keep without hash
				newEntries = append(newEntries, strconv.Itoa(e.line))
				continue
			}
			newEntries = append(newEntries, fmt.Sprintf("%d#%08x", e.line, h))
			anyChange = true
		}

		if !anyChange {
			skipped++
			continue
		}

		// Rewrite the file's first line
		absPath := filepath.Join(srcDir, exRelPath)
		data, readErr := os.ReadFile(absPath)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "warning: cannot read %s: %v\n", exRelPath, readErr)
			skipped++
			continue
		}

		content := string(data)
		newSourceLine := fmt.Sprintf("// Source: %s:%s", ref.mdxRelPath, strings.Join(newEntries, ","))

		// Replace the first line
		if idx := strings.Index(content, "\n"); idx >= 0 {
			content = newSourceLine + content[idx:]
		} else {
			content = newSourceLine
		}

		if writeErr := os.WriteFile(absPath, []byte(content), 0644); writeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: cannot write %s: %v\n", exRelPath, writeErr)
			skipped++
			continue
		}
		updated++
		fmt.Printf("  updated %s\n", exRelPath)
	}

	fmt.Printf("\n%d files updated, %d already had hashes, %d skipped\n", updated, alreadyHashed, skipped)
	return nil
}
