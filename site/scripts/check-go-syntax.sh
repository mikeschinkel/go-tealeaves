#!/usr/bin/env bash
# check-go-syntax.sh — validate Go code blocks in MDX docs with gofmt
# Finds all .mdx files under site/src/content/docs/, extracts fenced Go code
# blocks that contain a "package" declaration, runs gofmt -e on each block,
# and reports failures with file name and line number.
# Exits 0 if all blocks pass, 1 if any fail.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCS_DIR="$(cd "$SCRIPT_DIR/../src/content/docs" && pwd)"

failures=0
checked=0
tmpfile="$(mktemp /tmp/go-syntax-check-XXXXXX.go)"
trap 'rm -f "$tmpfile"' EXIT

while IFS= read -r -d '' mdx_file; do
    rel_path="${mdx_file#"$DOCS_DIR/"}"
    in_block=0
    block_start_line=0
    block_lines=()
    line_num=0

    while IFS= read -r line; do
        line_num=$((line_num + 1))

        if [[ $in_block -eq 0 ]]; then
            # Detect opening fence: ```go or ```Go (with optional trailing spaces)
            if [[ "$line" =~ ^(\`\`\`go[[:space:]]*)$ ]]; then
                in_block=1
                block_start_line=$line_num
                block_lines=()
            fi
        else
            # Detect closing fence
            if [[ "$line" =~ ^(\`\`\`[[:space:]]*)$ ]]; then
                in_block=0
                # Check if block contains a line starting with 'package'
                has_package=0
                for bl in "${block_lines[@]+"${block_lines[@]}"}"; do
                    if [[ "$bl" =~ ^[[:space:]]*package[[:space:]] ]]; then
                        has_package=1
                        break
                    fi
                done

                if [[ $has_package -eq 1 ]]; then
                    checked=$((checked + 1))
                    printf '%s\n' "${block_lines[@]+"${block_lines[@]}"}" > "$tmpfile"
                    if ! gofmt_out="$(gofmt -e "$tmpfile" 2>&1)"; then
                        echo "FAIL: $rel_path line $block_start_line"
                        echo "$gofmt_out" | sed 's|'"$tmpfile"'|  |g'
                        failures=$((failures + 1))
                    fi
                fi
                block_lines=()
            else
                block_lines+=("$line")
            fi
        fi
    done < "$mdx_file"

    # Unclosed block at EOF — skip silently
done < <(find "$DOCS_DIR" -name '*.mdx' -print0 | sort -z)

echo ""
echo "Checked $checked Go code block(s) across MDX files in $DOCS_DIR"

if [[ $failures -gt 0 ]]; then
    echo "FAILED: $failures block(s) have syntax errors"
    exit 1
else
    echo "OK: all blocks pass gofmt -e"
    exit 0
fi
