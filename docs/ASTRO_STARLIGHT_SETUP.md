# Astro + Starlight Setup Instructions

## For Repository: go-tealeaves

### Target Directory: ./site

These instructions are written so an AI coding agent (e.g. Codex) can
execute them step‑by‑step.

------------------------------------------------------------------------

## Objective

Add an Astro site using the Starlight docs template inside the existing
repository at:

    ./site

The Astro site will: - Use **Starlight** - Prefer **Bun** as the
runtime - Be configured for deployment to **GitHub Pages** - Not
interfere with the existing Go modules

------------------------------------------------------------------------

# 1. Preconditions

Ensure the following are installed:

-   Bun (preferred)
-   Git
-   Existing local clone of `go-tealeaves`

If Bun is not available, fall back to Node.js (LTS).

------------------------------------------------------------------------

# 2. Create Astro + Starlight Project

From the repository root:

``` bash
cd go-tealeaves
bun create astro@latest site
```

When prompted:

-   Template → Select **Starlight**
-   TypeScript → Yes
-   Install dependencies → Yes

If Bun fails, fallback:

``` bash
npm create astro@latest site
```

------------------------------------------------------------------------

# 3. Project Structure (Expected)

After creation, the structure should resemble:

    go-tealeaves/
    ├── go.mod
    ├── ...
    ├── site/
    │   ├── package.json
    │   ├── astro.config.mjs
    │   ├── public/
    │   ├── src/
    │   │   └── content/
    │   │       └── docs/
    │   │           └── index.md
    │   └── tsconfig.json

Do NOT modify Go code or root files.

------------------------------------------------------------------------

# 4. Configure Base Path for GitHub Pages

Edit:

    site/astro.config.mjs

Add or update the `site` and `base` properties:

``` js
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://<github-username>.github.io',
  base: '/go-tealeaves/',
  integrations: [
    starlight({
      title: 'Go Tealeaves',
    }),
  ],
});
```

Replace `<github-username>` accordingly.

If deploying under a subpath later (e.g., `/astro/`), adjust `base`.

------------------------------------------------------------------------

# 5. Update .gitignore (Root)

Ensure the following is added to the root `.gitignore`:

    site/node_modules/
    site/dist/

------------------------------------------------------------------------

# 6. Add NPM/Bun Scripts (Verify)

In `site/package.json`, ensure scripts include:

``` json
"scripts": {
  "dev": "astro dev",
  "build": "astro build",
  "preview": "astro preview"
}
```

------------------------------------------------------------------------

# 7. Verify Local Development

From repository root:

``` bash
cd site
bun run dev
```

Open:

    http://localhost:4321

Confirm Starlight renders successfully.

------------------------------------------------------------------------

# 8. Add GitHub Actions Workflow

Create:

    .github/workflows/deploy-docs.yml

With contents:

``` yaml
name: Deploy Docs

on:
  push:
    branches: [ main ]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Bun
        uses: oven-sh/setup-bun@v1

      - name: Install dependencies
        working-directory: site
        run: bun install

      - name: Build
        working-directory: site
        run: bun run build

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: site/dist

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

------------------------------------------------------------------------

# 9. Enable GitHub Pages

In repository:

Settings → Pages → Source → GitHub Actions

------------------------------------------------------------------------

# 10. Post‑Setup Checklist

-   Site builds locally
-   Workflow runs successfully
-   Pages site loads
-   Sidebar navigation works
-   Search functions (default Starlight behavior)

------------------------------------------------------------------------

# Notes

-   Do NOT commit `node_modules` or `dist`
-   Keep Astro fully contained inside `./site`
-   No modifications should be made to existing Go modules
-   Bun is preferred but Node fallback is acceptable

------------------------------------------------------------------------

End of instructions.
