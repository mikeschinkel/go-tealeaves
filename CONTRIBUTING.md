# Contributing to go-tealeaves

Thank you for your interest in contributing to go-tealeaves!

## Reporting Issues

- Use [GitHub Issues](https://github.com/mikeschinkel/go-tealeaves/issues) to report bugs or request features.
- Include the package name (e.g., `teagrid`, `teamodal`) in the issue title.
- For bugs, include Go version, OS, and a minimal reproduction case.

## Submitting Pull Requests

1. Fork the repository and create a feature branch.
2. Make your changes in the relevant package(s).
3. Run checks before submitting:
   ```bash
   just test
   ```
4. Write clear commit messages describing the change.
5. Open a PR against `main`.

## Multi-Module Structure

This repo contains 15 independent Go modules (one per package). When making changes:

- Only modify `go.mod`/`go.sum` in the package you're changing.
- Run `just tidy` to ensure all modules are consistent.
- If your change spans multiple packages, note that `teautils` is foundational — most other packages depend on it.

## Code Style

- Follow ClearPath conventions: named returns, `goto end` pattern, no `else` chains.
- Use `doterr` for structured error handling (see `doterr.go` in any package).
- Use `ansi.StringWidth()` for width calculations, never `len()`.
- Non-nil `tea.Cmd` return signals message consumption in modal components.
- See existing code in the package you're modifying for conventions.
- See `docs/BEST_PRACTICES_CHARM_V2.md` for Bubble Tea v2 patterns.

## Testing

- Add tests for new functionality.
- Ensure existing tests pass: `just test`.
- Example programs in `<module>/examples/` serve as integration tests — verify they still build.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
