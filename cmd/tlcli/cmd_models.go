package main

import (
	"flag"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/tools/go/packages"
)

const defaultCheckPattern = `\w+Model\b`

func runModels(args []string) error {
	fs := flag.NewFlagSet("models", flag.ExitOnError)

	var checkFlag string
	var checkBare bool

	// Custom flag parsing: --check without value defaults to the pattern
	fs.BoolVar(&checkBare, "check", false, "check naming convention (default pattern: "+defaultCheckPattern+")")
	fs.StringVar(&checkFlag, "check-pattern", "", "custom regex pattern for naming check")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `tlcli models — List types implementing the tea.Model component pattern

Scans all tea* packages under the current directory for exported types that
have Init() tea.Cmd, Update(tea.Msg) (Self, tea.Cmd), and View() tea.View
methods — the Bubble Tea v2 component pattern.

Usage:
  tlcli models [flags]

Flags:
  -check            Check naming convention (default pattern: %s)
  -check-pattern    Custom regex pattern for naming check
  -help             Show this help

Examples:
  tlcli models                          # list all component types
  tlcli models -check                   # check default naming convention
  tlcli models -check-pattern="Mdl$"    # check custom naming pattern

`, defaultCheckPattern)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	pattern := checkFlag
	if checkBare && pattern == "" {
		pattern = defaultCheckPattern
	}

	components, err := findComponents()
	if err != nil {
		return err
	}

	if len(components) == 0 {
		fmt.Println("No tea.Model component types found.")
		return nil
	}

	var checkRe *regexp.Regexp
	if pattern != "" {
		checkRe, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid check pattern: %w", err)
		}
	}

	// Print table
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	if checkRe != nil {
		fmt.Fprintln(w, "PACKAGE\tTYPE\tCHECK")
	} else {
		fmt.Fprintln(w, "PACKAGE\tTYPE")
	}

	var violations []componentInfo
	for _, c := range components {
		if checkRe != nil {
			status := "ok"
			if !checkRe.MatchString(c.typeName) {
				status = "FAIL"
				violations = append(violations, c)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", c.pkg, c.typeName, status)
		} else {
			fmt.Fprintf(w, "%s\t%s\n", c.pkg, c.typeName)
		}
	}
	w.Flush()

	if len(violations) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d naming violation(s) found:\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "  %s.%s does not match /%s/\n", v.pkg, v.typeName, pattern)
		}
		os.Exit(1)
	}

	return nil
}

type componentInfo struct {
	pkg      string
	typeName string
}

func findComponents() ([]componentInfo, error) {
	// Find all tea* directories
	dirs, err := filepath.Glob("tea*")
	if err != nil {
		return nil, fmt.Errorf("globbing tea* dirs: %w", err)
	}

	var patterns []string
	for _, d := range dirs {
		info, err := os.Stat(d)
		if err != nil || !info.IsDir() {
			continue
		}
		patterns = append(patterns, "./"+d)
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf("no tea* directories found (run from repo root)")
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  ".",
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}

	var results []componentInfo

	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}

		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			tn, ok := obj.(*types.TypeName)
			if !ok || !tn.Exported() {
				continue
			}

			t := tn.Type()

			// Skip interfaces — we want concrete types only
			if _, isIface := t.Underlying().(*types.Interface); isIface {
				continue
			}

			// For generic types, check the origin (uninstantiated) type's
			// method set, since methods are defined on the origin.
			checkType := t
			if named, ok := t.(*types.Named); ok && named.TypeParams() != nil {
				checkType = named.Origin()
			}

			if hasComponentPattern(checkType) {
				shortPkg := pkg.PkgPath
				// Show just the last segment for readability
				if i := strings.LastIndex(shortPkg, "/"); i >= 0 {
					shortPkg = shortPkg[i+1:]
				}
				typeName := tn.Name()
				// Annotate generic types with [T]
				if named, ok := t.(*types.Named); ok && named.TypeParams() != nil {
					typeName += "[T]"
				}
				results = append(results, componentInfo{
					pkg:      shortPkg,
					typeName: typeName,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].pkg != results[j].pkg {
			return results[i].pkg < results[j].pkg
		}
		return results[i].typeName < results[j].typeName
	})

	return results, nil
}

// hasComponentPattern checks if a type has the three methods of the
// Bubble Tea v2 component pattern:
//
//	Init() tea.Cmd
//	Update(tea.Msg) (Self, tea.Cmd)
//	View() tea.View
//
// The Update method may return the concrete type (not tea.Model), which
// is the standard pattern for BubbleTea v2 components.
func hasComponentPattern(t types.Type) bool {
	mset := types.NewMethodSet(t)

	var hasInit, hasUpdate, hasView bool

	for i := 0; i < mset.Len(); i++ {
		sel := mset.At(i)
		fn, ok := sel.Obj().(*types.Func)
		if !ok {
			continue
		}

		sig := fn.Type().(*types.Signature)

		switch fn.Name() {
		case "Init":
			hasInit = isInitSignature(sig)
		case "Update":
			hasUpdate = isUpdateSignature(sig, t)
		case "View":
			hasView = isViewSignature(sig)
		}
	}

	return hasInit && hasUpdate && hasView
}

// isInitSignature checks: Init() tea.Cmd
// - zero params, one result named Cmd from the bubbletea package
func isInitSignature(sig *types.Signature) bool {
	if sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	return isTeaType(sig.Results().At(0).Type(), "Cmd")
}

// isUpdateSignature checks: Update(tea.Msg) (Self, tea.Cmd)
// - one param (tea.Msg), two results (concrete type or tea.Model, tea.Cmd)
func isUpdateSignature(sig *types.Signature, selfType types.Type) bool {
	if sig.Params().Len() != 1 || sig.Results().Len() != 2 {
		return false
	}

	// Check param is tea.Msg (which is an alias for uv.Event)
	paramType := sig.Params().At(0).Type()
	if !isTeaMsgType(paramType) {
		return false
	}

	// Check second result is tea.Cmd
	if !isTeaType(sig.Results().At(1).Type(), "Cmd") {
		return false
	}

	// Check first result is either the concrete type or tea.Model
	resultType := sig.Results().At(0).Type()
	if types.Identical(resultType, selfType) {
		return true
	}
	// For generic types, the return type might be the origin
	if named, ok := resultType.(*types.Named); ok {
		if namedSelf, ok2 := selfType.(*types.Named); ok2 {
			if named.Origin() == namedSelf.Origin() {
				return true
			}
		}
	}
	// Also accept tea.Model (the interface itself)
	if isTeaType(resultType, "Model") {
		return true
	}

	return false
}

// isViewSignature checks: View() tea.View
// - zero params, one result named View from the bubbletea package
func isViewSignature(sig *types.Signature) bool {
	if sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	return isTeaType(sig.Results().At(0).Type(), "View")
}

// isTeaType checks if a type is a named type from the bubbletea package
func isTeaType(t types.Type, name string) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj.Name() != name {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}
	return strings.HasSuffix(pkg.Path(), "bubbletea/v2")
}

// isTeaMsgType checks if a type is tea.Msg, which is an alias for
// github.com/charmbracelet/ultraviolet.Event
func isTeaMsgType(t types.Type) bool {
	// tea.Msg is defined as: type Msg = uv.Event
	// It's a type alias, so the underlying type is uv.Event (an interface).
	// Check if it's the aliased type from ultraviolet or bubbletea.
	named, ok := t.(*types.Named)
	if ok {
		obj := named.Obj()
		pkg := obj.Pkg()
		if pkg == nil {
			return false
		}
		path := pkg.Path()
		// Accept both the direct bubbletea Msg and the ultraviolet Event
		if strings.HasSuffix(path, "bubbletea/v2") && obj.Name() == "Msg" {
			return true
		}
		if strings.HasSuffix(path, "ultraviolet") && obj.Name() == "Event" {
			return true
		}
	}

	// tea.Msg is a type alias to an interface, so it might appear as the
	// underlying interface type directly
	iface, ok := t.Underlying().(*types.Interface)
	if ok && iface.NumMethods() == 0 {
		// empty interface — could be the Event alias (any)
		// This is a fallback; the named check above should catch most cases
		return true
	}

	return false
}
