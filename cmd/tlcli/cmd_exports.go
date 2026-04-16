package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// exPathFilter is a repeatable flag for glob patterns.
type exPathFilter []string

func (f *exPathFilter) String() string {
	if f == nil {
		return ""
	}
	return strings.Join(*f, ",")
}

func (f *exPathFilter) Set(value string) error {
	v := strings.TrimSpace(value)
	if v != "" {
		*f = append(*f, v)
	}
	return nil
}

func runExports(args []string) error {
	fset := flag.NewFlagSet("exports", flag.ExitOnError)

	var includeDoterr, apiOnly bool
	var excludePaths, includePaths exPathFilter
	fset.BoolVar(&includeDoterr, "include-doterr", false, "Include doterr.go (default is excluded)")
	fset.BoolVar(&apiOnly, "api-only", false, "Exclude cmd/, site/, and examples/ (library API only)")
	fset.Var(&excludePaths, "exclude-path", "Glob pattern for module paths to exclude (repeatable)")
	fset.Var(&includePaths, "include-path", "Glob pattern for module paths to include (repeatable)")

	fset.Usage = func() {
		fmt.Fprintf(os.Stderr, `tlcli exports — List exported API for all packages in a repo

Finds all Go modules under the given root (or current directory), lists their
packages, and prints the exported API using AST parsing. Excludes doterr.go
and test files by default.

Usage:
  tlcli exports [flags] [repo-root]

Flags:
  -api-only                    Exclude cmd/, site/, and examples/ (library API only)
  -include-doterr              Include doterr.go (default is excluded)
  -exclude-path <glob>         Exclude modules/packages matching glob (repeatable)
  -include-path <glob>         Only include modules/packages matching glob (repeatable)
  -help                        Show this help

Path filters match against the repo-relative path (e.g. "./teagrid", "./cmd/tlcli").
When -include-path is set, only matching paths are included. -exclude-path
takes precedence over -include-path.

Examples:
  tlcli exports -api-only                # library packages only
  tlcli exports -include-path="teagrid"  # single package
`)
	}
	if err := fset.Parse(args); err != nil {
		return err
	}

	root := "."
	if fset.NArg() > 0 {
		root = fset.Arg(0)
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve root: %w", err)
	}
	info, err := os.Stat(absRoot)
	if err != nil {
		return fmt.Errorf("stat root: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("root is not a directory: %s", absRoot)
	}

	modules, err := exFindModules(absRoot)
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		return errors.New("no go.mod files found")
	}

	if apiOnly {
		excludePaths = append(excludePaths, "cmd/*", "site/*", "*/examples", "*/examples/*")
	}

	filter := exSourceFilter{
		exclude: []string{"doterr.go"},
	}
	if includeDoterr {
		filter.include = []string{"doterr.go"}
	}

	var out strings.Builder
	fmt.Fprintf(&out, "# Repo: `%s`\n", filepath.Base(absRoot))
	fmt.Fprintf(&out, "- Path: `%s`\n", exDisplayPath(absRoot))

	if repoMod := exRepoModulePath(absRoot, modules); repoMod != "" {
		fmt.Fprintf(&out, "- Module: `%s`\n", repoMod)
	}
	out.WriteString("\n")

	for _, mod := range modules {
		modRelPath := exRepoRelativePath(absRoot, mod.Dir)
		if !exPathAllowed(modRelPath, includePaths, excludePaths) {
			continue
		}

		fmt.Fprintf(&out, "## Module: `%s`\n", modRelPath)
		out.WriteString("\n")

		pkgs, listErr := exListModulePackages(mod.Dir)
		if listErr != nil {
			fmt.Fprintf(os.Stderr, "warning: %s: %v\n", mod.Dir, listErr)
			out.WriteString("_Unable to list packages._\n\n")
			continue
		}

		for _, pkg := range pkgs {
			pkgRelPath := exRepoRelativePath(absRoot, pkg.Dir)
			if !exPathAllowed(pkgRelPath, includePaths, excludePaths) {
				continue
			}

			fmt.Fprintf(&out, "### Package: `%s`\n", exPackageDisplayName(pkg))
			fmt.Fprintf(&out, "- Path: `%s`\n\n", pkgRelPath)

			docData, docErr := exBuildPackageDoc(pkg, filter)
			if docErr != nil {
				fmt.Fprintf(os.Stderr, "warning: %s: %v\n", pkg.ImportPath, docErr)
				out.WriteString("_Unable to parse package._\n\n")
				continue
			}

			exRenderSymbolGroup(&out, "Consts", docData.Consts)
			exRenderSymbolGroup(&out, "Vars", docData.Vars)
			exRenderSymbolGroup(&out, "Funcs", docData.Funcs)
			exRenderTypes(&out, docData.Types)
		}
	}

	_, err = os.Stdout.WriteString(out.String())
	return err
}

// ---------------------------------------------------------------------------
// Data types
// ---------------------------------------------------------------------------

type exListedPackage struct {
	ImportPath string
	Dir        string
	Name       string
}

type exModuleInfo struct {
	ImportPath string
	Dir        string
}

type exPackageDoc struct {
	Consts []string
	Vars   []string
	Funcs  []string
	Types  []exTypeDoc
}

type exTypeDoc struct {
	Signature  string
	Kind       string
	Consts     []string
	Vars       []string
	Funcs      []string
	Properties []string
	Methods    []string
}

type exSourceFilter struct {
	include []string
	exclude []string
}

// ---------------------------------------------------------------------------
// Rendering
// ---------------------------------------------------------------------------

func exRenderSymbolGroup(out *strings.Builder, heading string, symbols []string) {
	if len(symbols) == 0 {
		return
	}
	fmt.Fprintf(out, "#### %s\n", heading)
	for _, s := range symbols {
		fmt.Fprintf(out, "- `%s`\n", s)
	}
	out.WriteString("\n")
}

func exRenderTypes(out *strings.Builder, types []exTypeDoc) {
	if len(types) == 0 {
		return
	}
	out.WriteString("#### Types\n")
	out.WriteString("\n")
	for _, t := range types {
		fmt.Fprintf(out, "- `%s`\n", t.Signature)

		exRenderTypeSubgroup(out, "Consts", t.Consts)
		exRenderTypeSubgroup(out, "Vars", t.Vars)
		exRenderTypeSubgroup(out, "Funcs", t.Funcs)

		switch t.Kind {
		case "struct":
			exRenderTypeSubgroup(out, "Properties", t.Properties)
			exRenderTypeSubgroup(out, "Methods", t.Methods)
		case "interface":
			exRenderTypeSubgroup(out, "Methods", t.Methods)
		}
		out.WriteString("\n")
	}
}

func exRenderTypeSubgroup(out *strings.Builder, heading string, items []string) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintf(out, "  - %s\n", heading)
	for _, item := range items {
		fmt.Fprintf(out, "    - `%s`\n", item)
	}
}

// ---------------------------------------------------------------------------
// Module discovery
// ---------------------------------------------------------------------------

func exFindModules(root string) ([]exModuleInfo, error) {
	var modules []exModuleInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor":
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "go.mod" {
			modDir := filepath.Dir(path)
			modPath, modErr := exModuleImportPath(modDir)
			if modErr != nil {
				return modErr
			}
			modules = append(modules, exModuleInfo{
				ImportPath: modPath,
				Dir:        modDir,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk repo: %w", err)
	}

	sort.Slice(modules, func(i, j int) bool {
		iKey := exModuleSortKey(root, modules[i].Dir)
		jKey := exModuleSortKey(root, modules[j].Dir)
		if iKey == jKey {
			return modules[i].Dir < modules[j].Dir
		}
		return iKey < jKey
	})
	return modules, nil
}

func exModuleImportPath(modDir string) (string, error) {
	cmd := exec.Command("go", "mod", "edit", "-json")
	cmd.Dir = modDir
	output, err := cmd.Output()
	if err != nil {
		if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
			return "", fmt.Errorf("read module path for %s: %s", modDir, strings.TrimSpace(string(ee.Stderr)))
		}
		return "", fmt.Errorf("read module path for %s: %w", modDir, err)
	}

	var parsed struct {
		Module struct {
			Path string
		}
	}
	if err := json.Unmarshal(output, &parsed); err != nil {
		return "", fmt.Errorf("decode module json for %s: %w", modDir, err)
	}
	if strings.TrimSpace(parsed.Module.Path) == "" {
		return "", fmt.Errorf("module path not found in %s/go.mod", modDir)
	}
	return strings.TrimSpace(parsed.Module.Path), nil
}

func exListModulePackages(modDir string) ([]exListedPackage, error) {
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Dir = modDir
	output, err := cmd.Output()
	if err != nil {
		if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
			return nil, fmt.Errorf("go list failed: %s", strings.TrimSpace(string(ee.Stderr)))
		}
		return nil, fmt.Errorf("go list failed: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(output))
	pkgs := make([]exListedPackage, 0)
	for {
		var p exListedPackage
		if err := dec.Decode(&p); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode go list json: %w", err)
		}
		if p.ImportPath == "" || p.Dir == "" {
			continue
		}
		pkgs = append(pkgs, p)
	}

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ImportPath < pkgs[j].ImportPath
	})
	return pkgs, nil
}

// ---------------------------------------------------------------------------
// Source filtering
// ---------------------------------------------------------------------------

func (f exSourceFilter) allowFileName(name string) bool {
	if !strings.HasSuffix(name, ".go") {
		return false
	}
	if strings.HasSuffix(name, "_test.go") {
		return false
	}
	if exMatchesAny(f.exclude, name) && !exMatchesAny(f.include, name) {
		return false
	}
	return true
}

func exMatchesAny(patterns []string, value string) bool {
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		ok, err := filepath.Match(p, value)
		if err != nil {
			continue
		}
		if ok {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// AST parsing (ported from doc-go-repo)
// ---------------------------------------------------------------------------

func exBuildPackageDoc(pkg exListedPackage, filterRules exSourceFilter) (exPackageDoc, error) {
	fset := token.NewFileSet()
	filter := func(info fs.FileInfo) bool {
		return filterRules.allowFileName(info.Name())
	}

	parsed, err := parser.ParseDir(fset, pkg.Dir, filter, parser.ParseComments)
	if err != nil {
		return exPackageDoc{}, fmt.Errorf("parse dir: %w", err)
	}
	if len(parsed) == 0 {
		return exPackageDoc{}, nil
	}

	astPkg := exPickASTPackage(parsed, pkg.Name)
	docPkg := doc.New(astPkg, pkg.ImportPath, 0)

	result := exPackageDoc{
		Consts: exCollectValues(fset, docPkg.Consts),
		Vars:   exCollectValues(fset, docPkg.Vars),
		Funcs:  exCollectFuncs(fset, docPkg.Funcs),
		Types:  exCollectTypes(fset, docPkg.Types),
	}

	sort.Strings(result.Consts)
	sort.Strings(result.Vars)
	sort.Strings(result.Funcs)
	sort.Slice(result.Types, func(i, j int) bool {
		return result.Types[i].Signature < result.Types[j].Signature
	})

	return result, nil
}

func exPickASTPackage(pkgs map[string]*ast.Package, preferred string) *ast.Package {
	if p, ok := pkgs[preferred]; ok {
		return p
	}
	names := make([]string, 0, len(pkgs))
	for name := range pkgs {
		names = append(names, name)
	}
	sort.Strings(names)
	return pkgs[names[0]]
}

func exCollectValues(fset *token.FileSet, values []*doc.Value) []string {
	out := make([]string, 0)
	for _, value := range values {
		if value == nil || value.Decl == nil {
			continue
		}
		for _, spec := range value.Decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			typeSig := strings.TrimSpace(exNodeString(fset, vs.Type))
			for i, nameIdent := range vs.Names {
				name := strings.TrimSpace(nameIdent.Name)
				if name == "" || !ast.IsExported(name) {
					continue
				}
				var valueSig string
				switch {
				case len(vs.Values) == 1 && len(vs.Names) > 1:
					valueSig = strings.TrimSpace(exNodeString(fset, vs.Values[0]))
				case i < len(vs.Values):
					valueSig = strings.TrimSpace(exNodeString(fset, vs.Values[i]))
				}
				sig := name
				if typeSig != "" {
					sig += " " + typeSig
				}
				if valueSig != "" {
					sig += " = " + valueSig
				}
				out = append(out, strings.TrimSpace(sig))
			}
		}
	}
	return out
}

func exCollectFuncs(fset *token.FileSet, funcs []*doc.Func) []string {
	out := make([]string, 0, len(funcs))
	for _, fn := range funcs {
		if fn == nil || fn.Decl == nil {
			continue
		}
		if fn.Decl.Name == nil || !ast.IsExported(fn.Decl.Name.Name) {
			continue
		}
		out = append(out, strings.TrimSpace(exFormatFuncDeclSignature(fset, fn.Decl)))
	}
	sort.Strings(out)
	return out
}

func exCollectTypes(fset *token.FileSet, types []*doc.Type) []exTypeDoc {
	out := make([]exTypeDoc, 0, len(types))

	for _, t := range types {
		if t == nil || t.Decl == nil {
			continue
		}

		sig, ts, kind := exTypeSignature(fset, t.Decl)
		sig = strings.TrimSpace(sig)
		if sig == "" {
			continue
		}
		if ts == nil || ts.Name == nil || !ast.IsExported(ts.Name.Name) {
			continue
		}
		td := exTypeDoc{
			Signature: sig,
			Kind:      kind,
			Consts:    exCollectValues(fset, t.Consts),
			Vars:      exCollectValues(fset, t.Vars),
			Funcs:     exCollectFuncs(fset, t.Funcs),
		}

		switch concrete := ts.Type.(type) {
		case *ast.StructType:
			td.Properties = exCollectStructProperties(fset, concrete)
			td.Methods = exCollectFuncs(fset, t.Methods)
		case *ast.InterfaceType:
			td.Methods = exCollectInterfaceMethods(fset, concrete)
		}

		sort.Strings(td.Consts)
		sort.Strings(td.Vars)
		sort.Strings(td.Funcs)
		sort.Strings(td.Properties)
		sort.Strings(td.Methods)
		out = append(out, td)
	}

	return out
}

func exTypeSignature(fset *token.FileSet, decl *ast.GenDecl) (sig string, spec *ast.TypeSpec, kind string) {
	for _, s := range decl.Specs {
		ts, ok := s.(*ast.TypeSpec)
		if !ok {
			continue
		}
		name := strings.TrimSpace(ts.Name.Name)
		switch ts.Type.(type) {
		case *ast.StructType:
			return name + " struct{}", ts, "struct"
		case *ast.InterfaceType:
			return name + " interface{}", ts, "interface"
		}
		rhs := strings.TrimSpace(exNodeString(fset, ts.Type))
		if ts.Assign.IsValid() {
			return strings.TrimSpace(name + " = " + rhs), ts, "other"
		}
		return strings.TrimSpace(name + " " + rhs), ts, "other"
	}
	return "", &ast.TypeSpec{Name: ast.NewIdent("")}, "other"
}

func exCollectStructProperties(fset *token.FileSet, st *ast.StructType) []string {
	if st == nil || st.Fields == nil {
		return nil
	}
	out := make([]string, 0)
	for _, field := range st.Fields.List {
		typeSig := strings.TrimSpace(exNodeString(fset, field.Type))
		if len(field.Names) == 0 {
			if typeSig != "" && exIsExportedTypeExpr(field.Type) {
				out = append(out, typeSig)
			}
			continue
		}
		for _, name := range field.Names {
			sig := strings.TrimSpace(name.Name)
			if sig == "" || !ast.IsExported(sig) {
				continue
			}
			if typeSig != "" {
				sig += " " + typeSig
			}
			out = append(out, strings.TrimSpace(sig))
		}
	}
	return out
}

func exCollectInterfaceMethods(fset *token.FileSet, it *ast.InterfaceType) []string {
	if it == nil || it.Methods == nil {
		return nil
	}
	out := make([]string, 0)
	for _, field := range it.Methods.List {
		if len(field.Names) == 0 {
			embedded := strings.TrimSpace(exNodeString(fset, field.Type))
			if embedded != "" && exIsExportedTypeExpr(field.Type) {
				out = append(out, embedded)
			}
			continue
		}
		for _, name := range field.Names {
			method := strings.TrimSpace(name.Name)
			if method == "" || !ast.IsExported(method) {
				continue
			}
			ft, ok := field.Type.(*ast.FuncType)
			if !ok {
				typeSig := strings.TrimSpace(exNodeString(fset, field.Type))
				if typeSig != "" {
					out = append(out, strings.TrimSpace(method+" "+typeSig))
				} else {
					out = append(out, method)
				}
				continue
			}
			out = append(out, strings.TrimSpace(method+exFormatFuncTypeSuffix(fset, ft)))
		}
	}
	sort.Strings(out)
	return out
}

// ---------------------------------------------------------------------------
// Signature formatting
// ---------------------------------------------------------------------------

func exFormatFuncDeclSignature(fset *token.FileSet, fn *ast.FuncDecl) string {
	if fn == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(strings.TrimSpace(fn.Name.Name))
	b.WriteString(exFormatFuncTypeSuffix(fset, fn.Type))
	return strings.TrimSpace(b.String())
}

func exFormatFuncTypeSuffix(fset *token.FileSet, ft *ast.FuncType) string {
	if ft == nil {
		return "()"
	}
	var b strings.Builder
	if ft.TypeParams != nil && len(ft.TypeParams.List) > 0 {
		b.WriteString("[")
		b.WriteString(exFieldListToString(fset, ft.TypeParams))
		b.WriteString("]")
	}
	b.WriteString("(")
	b.WriteString(exFieldListToString(fset, ft.Params))
	b.WriteString(")")
	results := exFormatResults(fset, ft.Results)
	if results != "" {
		b.WriteString(" ")
		b.WriteString(results)
	}
	return b.String()
}

func exFormatResults(fset *token.FileSet, fl *ast.FieldList) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}
	if len(fl.List) == 1 && len(fl.List[0].Names) == 0 {
		return strings.TrimSpace(exNodeString(fset, fl.List[0].Type))
	}
	return "(" + exFieldListToString(fset, fl) + ")"
}

func exFieldListToString(fset *token.FileSet, fl *ast.FieldList) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}
	parts := make([]string, 0)
	for _, field := range fl.List {
		typeSig := strings.TrimSpace(exNodeString(fset, field.Type))
		if len(field.Names) == 0 {
			if typeSig != "" {
				parts = append(parts, typeSig)
			}
			continue
		}
		for _, name := range field.Names {
			n := strings.TrimSpace(name.Name)
			if n == "" {
				continue
			}
			if typeSig != "" {
				parts = append(parts, n+" "+typeSig)
			} else {
				parts = append(parts, n)
			}
		}
	}
	return strings.Join(parts, ", ")
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

func exNodeString(fset *token.FileSet, node any) string {
	if node == nil {
		return ""
	}
	var b bytes.Buffer
	if err := format.Node(&b, fset, node); err != nil {
		return ""
	}
	return b.String()
}

func exIsExportedTypeExpr(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		return ast.IsExported(t.Name)
	case *ast.StarExpr:
		return exIsExportedTypeExpr(t.X)
	case *ast.SelectorExpr:
		return t.Sel != nil && ast.IsExported(t.Sel.Name)
	case *ast.IndexExpr:
		return exIsExportedTypeExpr(t.X)
	case *ast.IndexListExpr:
		return exIsExportedTypeExpr(t.X)
	case *ast.ParenExpr:
		return exIsExportedTypeExpr(t.X)
	default:
		return false
	}
}

// exPathAllowed checks whether a repo-relative path (e.g. "./teagrid") passes
// the include/exclude filters. Exclude takes precedence. Matching is done against
// the path with the "./" prefix stripped, using filepath.Match glob semantics.
func exPathAllowed(relPath string, include, exclude exPathFilter) bool {
	// Strip "./" prefix for matching
	clean := strings.TrimPrefix(relPath, "./")
	if clean == "" {
		clean = "."
	}

	// Exclude takes precedence
	if exMatchesAny(exclude, clean) {
		return false
	}

	// If no include filters, everything is included
	if len(include) == 0 {
		return true
	}

	return exMatchesAny(include, clean)
}

func exPackageDisplayName(pkg exListedPackage) string {
	name := strings.TrimSpace(filepath.Base(pkg.ImportPath))
	if name != "" && name != "." && name != string(filepath.Separator) {
		return name
	}
	name = strings.TrimSpace(filepath.Base(pkg.Dir))
	if name != "" && name != "." && name != string(filepath.Separator) {
		return name
	}
	return strings.TrimSpace(pkg.Name)
}

func exRepoModulePath(repoRoot string, modules []exModuleInfo) string {
	for _, mod := range modules {
		if mod.Dir == repoRoot {
			return strings.TrimSpace(mod.ImportPath)
		}
	}
	if len(modules) == 0 {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(modules[0].ImportPath), "/")
	if len(parts) == 0 {
		return ""
	}
	for i := 1; i < len(modules); i++ {
		curr := strings.Split(strings.TrimSpace(modules[i].ImportPath), "/")
		max := len(parts)
		if len(curr) < max {
			max = len(curr)
		}
		j := 0
		for ; j < max; j++ {
			if parts[j] != curr[j] {
				break
			}
		}
		parts = parts[:j]
		if len(parts) == 0 {
			return ""
		}
	}
	return strings.Join(parts, "/")
}

func exRepoRelativePath(repoRoot, target string) string {
	rel, err := filepath.Rel(repoRoot, target)
	if err != nil {
		return exDisplayPath(target)
	}
	rel = filepath.ToSlash(rel)
	if rel == "." {
		return "./"
	}
	return "./" + strings.TrimPrefix(rel, "./")
}

func exModuleSortKey(repoRoot, moduleDir string) string {
	rel := exRepoRelativePath(repoRoot, moduleDir)
	if rel == "./examples" || strings.HasPrefix(rel, "./examples/") {
		return "1:" + rel
	}
	return "0:" + rel
}

func exDisplayPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return abs
	}
	if abs == home {
		return "~"
	}
	prefix := home + string(filepath.Separator)
	if strings.HasPrefix(abs, prefix) {
		return "~" + string(filepath.Separator) + strings.TrimPrefix(abs, prefix)
	}
	return abs
}
