package teadiffview

// FileStatus indicates whether a file is new, deleted, or modified.
type FileStatus int

const (
	FileModified FileStatus = iota
	FileNew
	FileDeleted
)

// CondensedBlock represents a condensed view of a change block.
type CondensedBlock struct {
	Type          string   // "added" or "deleted"
	LineCount     int
	ContextBefore []string
	ChangedLines  []string
	ContextAfter  []string
	IsTruncated   bool
}

// FileDiff holds condensed diffs for a single file.
type FileDiff struct {
	Path   string
	Status FileStatus
	Blocks []CondensedBlock
}

// DiffRenderer formats diff content for a specific output medium.
type DiffRenderer interface {
	RenderFileHeader(path string, status FileStatus, width int) string
	RenderBlockHeader(blockType string, lineCount int) string
	RenderContextLine(line string, status FileStatus, width int) string
	RenderAddedLine(line string, status FileStatus, width int) string
	RenderDeletedLine(line string, status FileStatus, width int) string
	RenderTruncation(status FileStatus) string
	RenderSeparator() string
}
