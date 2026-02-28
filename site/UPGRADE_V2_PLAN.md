# Upgrading go-tealeaves Documentation Site for v0.2.0

## Context

The `site/` directory contains an Astro + Starlight documentation site that is currently a default Starlight starter with placeholder content. This plan covers upgrading the site to reflect the Charm v2 migration and serve as the primary documentation for go-tealeaves.

**Prerequisite:** This plan executes AFTER `UPGRADE_V2_PLAN.md` (root) reaches at least the DOCS phase. All Go code must be migrated, tested, and validated before documentation can accurately describe v2 APIs.

**Project location:** `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/site/`

---

## Scope

In scope:
1. Replace all Starlight placeholder content with real go-tealeaves documentation.
2. Document every module with v2-correct API references and examples.
3. Publish the migration guide (`UPGRADE_GUIDE_V2.md`) as a site page.
4. Update Astro/Starlight config for proper site metadata, sidebar structure, and navigation.
5. Verify build and deploy pipeline still works after content changes.

Out of scope:
1. Custom Starlight theme or component overrides (defer to post-v1.0).
2. API docs auto-generation from Go source (defer to post-v1.0).
3. Custom domain setup (deferred until public announcement).
4. Astro or Starlight major version upgrades (only if required by breakage).

---

## SITEBASE — Baseline and Cleanup

1. Verify current site builds locally:
   ```bash
   cd site && bun install && bun run build
   ```
2. Remove Starlight starter placeholder content:
   - `src/content/docs/index.mdx` — replace entirely
   - `src/content/docs/guides/example.md` — delete
   - `src/content/docs/reference/example.md` — delete
   - `src/assets/houston.webp` — delete (Starlight mascot)
3. Replace `site/README.md` with project-specific content (not the Starlight starter kit boilerplate).
4. Verify site still builds after cleanup (empty content is fine at this stage).

**Acceptance gate:** Site builds with no placeholder content; no Starlight starter boilerplate remains.

---

## SITECONFIG — Astro and Starlight Configuration

1. Update `astro.config.mjs`:
   - Confirm `site` and `base` are correct for current deployment target.
   - Update Starlight `title` if needed.
   - Restructure `sidebar` to match the documentation structure defined in SITECONTENT.
   - Add `favicon` if a project logo/icon exists.
2. Review `package.json` dependencies:
   - Run `bun update` to pick up latest compatible Astro/Starlight patches.
   - Verify no security advisories in dependencies.
3. Confirm `tsconfig.json` and `content.config.ts` need no changes (they are generic Starlight config).

**Acceptance gate:** Config reflects real project structure; `bun run build` succeeds; sidebar matches planned content hierarchy.

---

## SITECONTENT — Write Documentation Pages

All content uses Charm v2 APIs exclusively. No v1 code examples in the main docs.

### Landing page (`src/content/docs/index.mdx`)

- Project name, tagline, and brief description.
- Feature highlights (one card per module).
- Quick install snippet (`go get`).
- Links to Getting Started guide and module reference pages.

### Getting Started guide (`src/content/docs/guides/getting-started.md`)

- Prerequisites (Go version, Charm v2 dependencies).
- Installation instructions.
- Minimal "hello world" example using one or two modules.
- Link to examples directory in the repo.

### Module reference pages (one per module under `src/content/docs/reference/`)

Each page covers:
- Package import path.
- Purpose and when to use it.
- Key types and functions with v2-correct signatures.
- Usage example (working code snippet).
- Links to the module's README and example apps.

Pages to create:
- `src/content/docs/reference/teadd.md` — Dropdown component
- `src/content/docs/reference/teadep.md` — Dependency dropdown
- `src/content/docs/reference/teamodal.md` — Modal dialogs (confirm, choice, list, progress)
- `src/content/docs/reference/teastatus.md` — Status bar
- `src/content/docs/reference/teatextsel.md` — Text selection
- `src/content/docs/reference/teatree.md` — Tree navigation
- `src/content/docs/reference/teautils.md` — Utility components (key registry, help visor)
- `src/content/docs/reference/teanotify.md` — Notifications
- `src/content/docs/reference/teagrid.md` — Grid/table

### Migration guide (`src/content/docs/guides/charm-v2-migration.md`)

- Publish the content from `UPGRADE_GUIDE_V2.md` (created in the DOCS phase of the root plan) as a Starlight page.
- Add Starlight frontmatter (title, description).
- Ensure all code examples render correctly in Starlight's code blocks.

### Architecture / design decisions (`src/content/docs/guides/architecture.md`)

- Module dependency graph (which modules depend on which).
- Design philosophy (ClearPath style, doterr errors, separation of concerns).
- How modules compose together in a Bubble Tea app.

**Acceptance gate:** All pages render correctly; code examples use v2 APIs exclusively; sidebar navigation matches content structure; no broken internal links.

---

## SITESIDEBAR — Finalize Sidebar Configuration

Update `astro.config.mjs` sidebar to match final content:

```js
sidebar: [
  {
    label: 'Getting Started',
    items: [
      { label: 'Installation', slug: 'guides/getting-started' },
      { label: 'Architecture', slug: 'guides/architecture' },
    ],
  },
  {
    label: 'Components',
    autogenerate: { directory: 'reference' },
  },
  {
    label: 'Migration',
    items: [
      { label: 'Charm v2 Migration', slug: 'guides/charm-v2-migration' },
    ],
  },
],
```

Adjust labels and groupings based on actual content created in SITECONTENT.

**Acceptance gate:** Sidebar reflects all published pages; navigation is intuitive; no orphaned or unreachable pages.

---

## SITEVERIFY — Build, Deploy, and Smoke Test

1. Clean build:
   ```bash
   cd site && rm -rf dist && bun run build
   ```
2. Local preview:
   ```bash
   bun run preview
   ```
3. Verify in browser:
   - Landing page loads with correct project info.
   - All sidebar links work.
   - All code examples render with syntax highlighting.
   - Starlight search indexes all pages.
   - No 404s on internal navigation.
4. Push to `main` and verify GitHub Actions deploy succeeds.
5. Verify live site at `https://mikeschinkel.github.io/go-tealeaves/`.

**Acceptance gate:** Site builds cleanly; local preview shows all content; GitHub Actions deploy succeeds; live site is accessible and fully navigable.

---

## Dependency on Root Plan Phases

| Site Phase | Depends On (Root Plan) | Reason |
|---|---|---|
| SITEBASE | None | Cleanup is independent of Go migration |
| SITECONFIG | None | Config is independent of Go migration |
| SITECONTENT | DOCS | Content must describe v2 APIs accurately; `UPGRADE_GUIDE_V2.md` must exist |
| SITESIDEBAR | SITECONTENT | Sidebar matches final content |
| SITEVERIFY | RELEASE | Final verification should happen against released v0.2.0 code |

SITEBASE and SITECONFIG can begin in parallel with Go migration work. SITECONTENT is blocked until the root plan's DOCS phase produces `UPGRADE_GUIDE_V2.md` and all module APIs are finalized.

---

## Risk Notes

1. **Content accuracy** — Writing docs against in-progress v2 migration risks documenting APIs that change. Mitigate by waiting for DOCS phase completion.
2. **Astro/Starlight breaking changes** — Current versions (`astro@^5.6.1`, `@astrojs/starlight@^0.37.6`) may receive updates. Pin versions if instability appears; otherwise let `^` ranges pick up patches.
3. **Code example drift** — Standalone code snippets in docs can diverge from actual working examples. Mitigate by cross-referencing against `examples/` directory code.
4. **Search indexing** — Pagefind (Starlight's default search) indexes at build time. Verify search works for module names and key API terms after content is finalized.
