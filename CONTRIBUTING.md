# Contributing to go-tealeaves

Thank you for your interest in contributing to go-tealeaves!

## Reporting Issues

- Use [GitHub Issues](https://github.com/mikeschinkel/go-tealeaves/issues) to report bugs or request features.
- Include the package name (e.g., `teadd`, `teamodal`) in the issue title.
- For bugs, include Go version, OS, and a minimal reproduction case.

## Submitting Pull Requests

1. Fork the repository and create a feature branch.
2. Make your changes in the relevant package(s).
3. Run checks before submitting:
   ```bash
   make test && make vet
   ```
4. Write clear commit messages describing the change.
5. Open a PR against `main`.

## Multi-Module Structure

This repo contains 7 independent Go modules (one per package). When making changes:

- Only modify `go.mod`/`go.sum` in the package you're changing.
- Run `make tidy` to ensure all modules are consistent.
- If your change spans multiple packages, note the dependency order:
  - `teautils` must be updated before `teamodal`
  - `teadd` must be updated before `teadep`
  - All other packages are independent.

## Code Style

- Follow patterns documented in `docs/WHAT_WE_HAVE_LEARNED.md`.
- Use `goto end` with named returns instead of early returns (ClearPath style).
- Use `ansi.StringWidth()` for width calculations, never `len()`.
- Non-nil `tea.Cmd` return signals message consumption in modal components.
- See existing code in the package you're modifying for conventions.

## Testing

- Add tests for new functionality.
- Ensure existing tests pass: `make test`.
- Run `make vet` to catch common issues.
- Example programs in `examples/` serve as integration tests — verify they still build with `make build-examples`.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
