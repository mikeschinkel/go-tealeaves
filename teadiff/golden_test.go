package teadiff

import (
	"log/slog"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/mikeschinkel/go-diffutils"
)

func newTestSplitDiff(width, height int) SplitDiffModel {
	m := NewSplitDiffModel(&SplitDiffModelArgs{
		Width:  width / 2,
		Height: height,
		Logger: slog.New(slog.DiscardHandler),
	})

	content := &diffutils.DiffContent{
		OldLines: []string{
			"package main",
			"",
			"func main() {",
			"\tprintln(\"hello\")",
			"}",
		},
		NewLines: []string{
			"package main",
			"",
			"import \"fmt\"",
			"",
			"func main() {",
			"\tfmt.Println(\"hello, world\")",
			"}",
		},
		Changes: []diffutils.DiffChange{
			{Type: diffutils.LinesAdded, OldRange: diffutils.LineRange{Start: 3, Count: 0}, NewRange: diffutils.LineRange{Start: 3, Count: 2}},
			{Type: diffutils.LinesDeleted, OldRange: diffutils.LineRange{Start: 4, Count: 1}, NewRange: diffutils.LineRange{Start: 6, Count: 0}},
			{Type: diffutils.LinesAdded, OldRange: diffutils.LineRange{Start: 4, Count: 0}, NewRange: diffutils.LineRange{Start: 6, Count: 1}},
		},
		Label: "main.go",
	}

	m = m.SetContent(content)
	m = m.SetSize(width/2, height)
	return m
}

func TestSplitDiffModel_Golden_80x24(t *testing.T) {
	m := newTestSplitDiff(80, 24)
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}

func TestSplitDiffModel_Golden_120x40(t *testing.T) {
	m := newTestSplitDiff(120, 40)
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}
