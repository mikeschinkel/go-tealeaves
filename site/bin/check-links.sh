#!/usr/bin/env bash
# check-links.sh — Build the Astro site and verify all internal href links resolve.
# Exits with 0 if no broken links found, 1 if broken links are present.
#
# The site uses base: '/go-tealeaves' in astro.config.mjs, so all links are
# prefixed with /go-tealeaves/. We strip that prefix when checking files in dist/.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
SITE_DIR="${REPO_ROOT}/site"
DIST_DIR="${SITE_DIR}/dist"
BASE_PATH="/go-tealeaves"

echo "=== Tea Leaves link checker ==="
echo ""

# ── Step 1: Build the site ────────────────────────────────────────────────────
echo "Building site (just site-build)..."
cd "${REPO_ROOT}"
if ! just site-build; then
  echo "ERROR: Site build failed. Aborting link check."
  exit 1
fi
echo "Build complete. Dist at: ${DIST_DIR}"
echo ""

# ── Step 2: Collect all HTML files ───────────────────────────────────────────
mapfile -t HTML_FILES < <(find "${DIST_DIR}" -name "*.html" | sort)

if [[ ${#HTML_FILES[@]} -eq 0 ]]; then
  echo "ERROR: No HTML files found in ${DIST_DIR}"
  exit 1
fi

echo "Found ${#HTML_FILES[@]} HTML file(s) to scan."
echo "Using base path: ${BASE_PATH}"
echo ""

# ── Step 3: Extract and verify internal hrefs ─────────────────────────────────
broken_count=0
declare -A seen_broken  # dedup broken link reports

check_href() {
  local source_file="$1"
  local href="$2"

  # Strip query string and fragment
  local path="${href%%\?*}"
  path="${path%%#*}"

  # Skip empty, external links, mailto, tel, javascript
  if [[ -z "${path}" ]] || \
     [[ "${path}" == http* ]] || \
     [[ "${path}" == mailto:* ]] || \
     [[ "${path}" == tel:* ]] || \
     [[ "${path}" == javascript:* ]]; then
    return 0
  fi

  # Strip the base path prefix for resolution within dist/
  local dist_path="${path}"
  if [[ "${dist_path}" == "${BASE_PATH}"* ]]; then
    dist_path="${dist_path#${BASE_PATH}}"
    # Empty means root: /
    [[ -z "${dist_path}" ]] && dist_path="/"
  fi

  # Resolve to absolute path in dist
  local resolved
  if [[ "${dist_path}" == /* ]]; then
    resolved="${DIST_DIR}${dist_path}"
  else
    local source_dir
    source_dir="$(dirname "${source_file}")"
    resolved="${source_dir}/${dist_path}"
  fi

  # Accept if: exact file exists, OR path has an index.html, OR path.html exists
  if [[ -f "${resolved}" ]] || \
     [[ -f "${resolved%/}/index.html" ]] || \
     [[ -f "${resolved}/index.html" ]] || \
     [[ -f "${resolved}.html" ]]; then
    return 0
  fi

  # Also accept if it's a directory
  if [[ -d "${resolved}" ]]; then
    return 0
  fi

  local key="${source_file}::${path}"
  if [[ -z "${seen_broken[${key}]+x}" ]]; then
    seen_broken["${key}"]=1
    local rel_source="${source_file#${DIST_DIR}/}"
    echo "BROKEN  [${rel_source}]  →  ${href}"
    broken_count=$((broken_count + 1))
  fi
}

# Iterate over all HTML files
for html_file in "${HTML_FILES[@]}"; do
  # Extract href="..." values using grep + sed
  while IFS= read -r href; do
    check_href "${html_file}" "${href}"
  done < <(
    grep -oE 'href="[^"]*"' "${html_file}" 2>/dev/null \
      | sed 's/^href="//;s/"$//'
  )
done

echo ""
echo "=== Link check results ==="
if [[ ${broken_count} -eq 0 ]]; then
  echo "All internal links OK (0 broken)."
  exit 0
else
  echo "Found ${broken_count} broken internal link(s)."
  exit 1
fi
