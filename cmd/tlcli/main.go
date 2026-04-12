package main

import (
	"fmt"
	"os"
)

const usageText = `tlcli — Tea Leaves CLI

Usage:
  tlcli <command> [flags]

Commands:
  audit     Run documentation gap analysis and code example verification
  models    List types implementing the tea.Model component pattern
  colors    Interactive 256-color palette viewer

Run 'tlcli <command> -help' for details on a specific command.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usageText)
		os.Exit(1)
	}

	var err error

	switch os.Args[1] {
	case "audit":
		err = runAudit(os.Args[2:])
	case "models":
		err = runModels(os.Args[2:])
	case "colors":
		err = runColors(os.Args[2:])
	case "-help", "--help", "help":
		fmt.Print(usageText)
		return
	default:
		fmt.Fprintf(os.Stderr, "tlcli: unknown command %q\n\n%s", os.Args[1], usageText)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "tlcli %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
}
