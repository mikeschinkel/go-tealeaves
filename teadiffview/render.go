package teadiffview

// RenderFileDiffs renders multiple file diffs into styled lines.
// Each file gets a header, followed by its blocks with context and change lines.
// Files are separated by the renderer's separator output.
func RenderFileDiffs(files []FileDiff, renderer DiffRenderer, width int) []string {
	var lines []string

	for i, file := range files {
		if i > 0 {
			sep := renderer.RenderSeparator()
			if sep != "" {
				lines = append(lines, sep)
			}
			lines = append(lines, "") // blank line between files
		}

		lines = append(lines, renderer.RenderFileHeader(file.Path, file.Status, width))

		for _, block := range file.Blocks {
			lines = append(lines, renderer.RenderBlockHeader(block.Type, block.LineCount))

			for _, line := range block.ContextBefore {
				lines = append(lines, renderer.RenderContextLine(line, file.Status, width))
			}

			for _, line := range block.ChangedLines {
				switch block.Type {
				case "added":
					lines = append(lines, renderer.RenderAddedLine(line, file.Status, width))
				case "deleted":
					lines = append(lines, renderer.RenderDeletedLine(line, file.Status, width))
				default:
					lines = append(lines, renderer.RenderContextLine(line, file.Status, width))
				}
			}

			for _, line := range block.ContextAfter {
				lines = append(lines, renderer.RenderContextLine(line, file.Status, width))
			}

			if block.IsTruncated {
				lines = append(lines, renderer.RenderTruncation(file.Status))
			}
		}
	}

	return lines
}
