// Program jediterm-bug reproduces a rendering bug in JediTerm (GoLand/IntelliJ)
// caused by JediTerm not supporting ultraviolet's use of differential renderer
// when transitioning between two fundamentally different view layouts.
//
// The two frames were recorded from a real gomion session using RecordingModel.
// Frame 1: wide two-pane layout (44-col left + 36-col right)
// Frame 2: narrow single-pane layout (33-col)
//
// Press space to toggle, q to quit.
// Writes trace.log next to the executable.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func main() {
	logPath := "trace.log"
	if exe, err := os.Executable(); err == nil {
		logPath = filepath.Join(filepath.Dir(exe), "trace.log")
	}

	tee, err := newTeeWriter(os.Stdout, logPath)
	if err != nil {
		stderrf("trace log: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		err := tee.Close()
		if err != nil {
			stderrf("Error: %v\n", err)
			os.Exit(1)
		}
	}()

	p := tea.NewProgram(model{}, tea.WithOutput(tee))
	if _, err := p.Run(); err != nil {
		stderrf("Error: %v\n", err)
		os.Exit(1)
	}
}

// ── Model ──

type model struct {
	viewB bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.viewB = !m.viewB
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var body string
	if m.viewB {
		body = frameB
	} else {
		body = frameA
	}
	v := tea.NewView(body)
	v.AltScreen = true
	return v
}

// ── Recorded frames ──

// frameA: two-pane layout (tea.KeyPressMsg)
var frameA = "\x1b[1;38;5;46mgithub.com/mikeschinkel/gomion\x1b[m\n\x1b[38;2;0;135;255m╭──────────────────────────────────────────╮\x1b[m\x1b[38;2;188;188;188m╭──────────────────────────────────╮\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m\x1b[1m Select a Commit Target:\x1b[m                  \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[1;38;5;130mModule\x1b[m                           \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m   git@github.com:mikeschinkel/gomion.git \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[38;5;238m  github.com/mikeschinkel/go-dt\x1b[m  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m ▶ github.com/mikeschinkel/gomion/cli     \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m                                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m   github.com/mikeschinkel/gomion/gommod  \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[38;5;244mDependencies: 0\x1b[m                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m ▶ github.com/mikeschinkel/go-cfgstore    \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m                                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m ▶ github.com/mikeschinkel/go-logutil     \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[1;38;5;130mStatus\x1b[m                           \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m ▶ github.com/mikeschinkel/go-cliutil     \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m   \x1b[1;38;5;166m● Uncommitted changes\x1b[m          \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m ▶ github.com/mikeschinkel/go-dt/dtx      \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m                                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m\x1b[1;38;5;230;48;5;62m ▶ github.com/mikeschinkel/go-dt \x1b[m\x1b[48;5;62m         \x1b[m\x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[1;38;5;130mCommit Group\x1b[m                     \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m   \x1b[38;5;166m✗ go-dt/\x1b[m                       \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m   \x1b[38;5;28m✓ go-dt/test/\x1b[m                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m   \x1b[38;5;28m✓ go-dt/examples/basic_usage/\x1b[m  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m                                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[1;38;5;130mLatest Tag\x1b[m                       \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m   \x1b[38;5;28mv0.6.0\x1b[m                         \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[38;5;244m  Up to date\x1b[m                     \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m                                  \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[1;38;5;130mScan Directories\x1b[m                 \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[38;5;244m  ~/Projects/go-pkgs\x1b[m             \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                                          \x1b[38;2;0;135;255m│\x1b[m\x1b[38;2;188;188;188m│\x1b[m \x1b[38;5;244m  ~/Projects/xmlui\x1b[m               \x1b[38;2;188;188;188m│\x1b[m\n\x1b[38;2;0;135;255m╰──────────────────────────────────────────╯\x1b[m\x1b[38;2;188;188;188m╰──────────────────────────────────╯\x1b[m\n \x1b[1;38;5;30m[?]\x1b[m \x1b[38;5;238mMenu\x1b[m  \x1b[1;38;5;30m[n]\x1b[m \x1b[38;5;238mGuide\x1b[m  \x1b[1;38;5;30m[enter]\x1b[m \x1b[38;5;238mSelect\x1b[m  \x1b[1;38;5;30m[q]\x1b[m \x1b[38;5;238mQuit\x1b[m              \x1b[1;38;5;166mUncommitted\x1b[m\x1b[38;5;250m | \x1b[m\x1b[38;5;244mv0.6.0\x1b[m"

// frameB: single-pane layout (gomtui.DrillDownMsg)
var frameB = "\x1b[38;5;34mgomion\x1b[m\x1b[38;5;51m > \x1b[m\x1b[1;38;5;46mgithub.com/mikeschinkel/go-dt\x1b[m\n\x1b[38;2;0;135;255m╭───────────────────────────────╮\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m \x1b[1;38;5;130mFiles to Commit:\x1b[m              \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m └─▼ \x1b[7;38;2;128;128;128mgo-dt\x1b[m [\x1b[38;2;128;128;128mo\x1b[m]                 \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m   └─ \x1b[38;2;128;128;128mpath_segments_ext.go\x1b[m [\x1b[38;2;128;128;128mo\x1b[m] \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m│\x1b[m                               \x1b[38;2;0;135;255m│\x1b[m\n\x1b[38;2;0;135;255m╰───────────────────────────────╯\x1b[m\n \x1b[1;38;5;30m[?]\x1b[m \x1b[38;5;238mMenu\x1b[m  \x1b[1;38;5;30m[n]\x1b[m \x1b[38;5;238mGuide\x1b[m  \x1b[1;38;5;30m[tab]\x1b[m \x1b[38;5;238mSwitch pane\x1b[m  \x1b[1;38;5;30m[esc]\x1b[m \x1b[38;5;238mBack\x1b[m  \x1b[1;38;5;30m[enter]\x1b[m \x1b[38;5;238mCommits\x1b[m            "

// ── TeeWriter ──

type teeWriter struct {
	file *os.File
	log  *os.File
}

func newTeeWriter(output *os.File, logPath string) (*teeWriter, error) {
	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}
	return &teeWriter{file: output, log: logFile}, nil
}

func (w *teeWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.Write(p)
	if err == nil {
		_, err = w.log.Write(p[:n])
	}
	return n, err
}

func (w *teeWriter) Read(p []byte) (int, error) { return w.file.Read(p) }
func (w *teeWriter) Fd() uintptr                { return w.file.Fd() }
func (w *teeWriter) Close() error               { return w.log.Close() }

func stderrf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}
