# go-tealeaves

A Go component library for Bubble Tea TUI applications.

## Project Structure

- **Component packages**: `teacrumbs/`, `teacolor/`, `teadiff/`, `teafields/`,
  `teagrid/`, `teaguide/`, `teahelp/`, `teahilite/`, `tealayout/`, `teamodal/`,
  `teanotify/`, `teapane/`, `teastatus/`, `teaterm/`, `teatext/`, `teatree/`, `teautils/`
- **Documentation site**: `site/` — Astro Starlight, docs at `site/src/content/docs/`
- **Examples**: `examples/` — each is an independent Go module
- **Multi-module workspace**: uses `go.work` (gitignored, generated locally)

## Skills to Apply

When writing or reviewing Go code in this project, invoke these skills:

- `/go-house-rules` — ClearPath style, doterr errors, go-dt paths (mandatory for all Go code)
- `/go-bubbletea-tui-author` — Bubble Tea patterns and hard-won lessons (mandatory for component code)
- `/go-testing-assertion-author` — testing conventions (mandatory for test code)

When writing documentation:

- Match existing page structure in `site/src/content/docs/components/`
- Use MDX format with Starlight frontmatter
- Component docs should include: overview, props/fields, usage examples, related components

## Build & Test

- `just test` — run tests across all modules
- `just vet` — run go vet across all modules
- `just tidy` — go mod tidy across all modules and examples
- `just serve` — dev server for docs site
- `just site-build` — production build of docs site
