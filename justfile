# Modules in this repo
modules := "teacolor teacrumbs teadiff teafields teagrid teaguide teahelp teahilite tealayout teamodal teanotify teapane teastatus teatext teatree teautils"

# Example modules (each is an independent Go module)
examples := "teadiff/examples/splitdiff teafields/examples/simple teafields/examples/demo teagrid/examples/filtering teagrid/examples/panning teagrid/examples/scrolling teagrid/examples/simplest teagrid/examples/sorting tealayout/examples/multipane teamodal/examples/choices teamodal/examples/editlist teamodal/examples/multiselect teamodal/examples/various teamodal/examples/vertical teanotify/examples/simple teastatus/examples/statusbar teatext/examples/editor teatree/examples/drilldown teatree/examples/filetree teautils/examples/keyhelp teautils/examples/theming teaguide/example"

# Show available recipes
help:
    @just --list

# Build color-viewer to ./bin/
build:
    @echo "Building color-viewer..."
    @mkdir -p bin
    cd cmd/color-viewer && go build -o ../../bin/color-viewer .
    @echo "Built to ./bin/color-viewer"

# Run tests across all modules
test:
    #!/usr/bin/env bash
    for mod in {{modules}}; do
        echo "Testing $mod..."
        (cd "$mod" && go test ./...) || exit 1
    done
    echo "All tests passed!"

# Run go mod tidy across all modules and examples
tidy:
    #!/usr/bin/env bash
    for mod in {{modules}}; do
        echo "Tidying $mod..."
        (cd "$mod" && go mod tidy) || exit 1
    done
    echo "Tidying cmd/color-viewer..."
    (cd cmd/color-viewer && go mod tidy) || exit 1
    for ex in {{examples}}; do
        echo "Tidying $ex..."
        (cd "$ex" && go mod tidy) || exit 1
    done
    echo "All modules tidied!"

# Format code with gofmt
fmt:
    @echo "Formatting code..."
    @gofmt -s -w .

# Run go vet across all modules
vet:
    #!/usr/bin/env bash
    for mod in {{modules}}; do
        echo "Vetting $mod..."
        (cd "$mod" && go vet ./...) || exit 1
    done
    echo "All modules vetted!"

# Linter tool
linter := "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2"

# Run golangci-lint across all modules
lint:
    #!/usr/bin/env bash
    for mod in {{modules}}; do
        echo "Linting $mod..."
        (cd "$mod" && go run {{linter}} run ./... --timeout=5m) || exit 1
    done
    echo "All modules linted!"

# Clean build artifacts
clean:
    #!/usr/bin/env bash
    echo "Cleaning build artifacts..."
    rm -rf bin
    for mod in {{modules}}; do
        (cd "$mod" && go clean)
    done
    (cd cmd/color-viewer && go clean)
    echo "Clean complete!"

# Build all example programs
build-examples:
    #!/usr/bin/env bash
    mkdir -p bin/examples
    for ex in {{examples}}; do
        name=$(basename "$ex")
        parent=$(basename "$(dirname "$ex")")
        echo "Building $parent/$name..."
        (cd "$ex" && go build -o ../../../bin/examples/"$parent-$name" .) || exit 1
    done
    echo "All examples built to ./bin/examples/"

# Build the documentation site
site-build:
    cd site && npx astro build

# Serve the documentation site locally (dev mode with hot reload)
serve:
    cd site && npx astro dev

# Preview the production build of the documentation site
site-preview:
    cd site && npx astro build && npx astro preview

# --- h2 multi-agent workflows ---

# Project-specific h2 settings
h2_language := "Go"
h2_doc_framework := "Astro Starlight (MDX)"
h2_agents := "manager auditor docs-updater docs-creator icon-designer"
h2_session := "h2-agents"
main_session := "active"

# Launch h2 agents — each in its own tmux window (idempotent)
h2-start:
    #!/usr/bin/env bash
    set -e
    wd="$(pwd)"
    lang="{{h2_language}}"
    framework="{{h2_doc_framework}}"
    session="{{h2_session}}"
    # Stop any running h2 agents and kill previous tmux session
    for agent in {{h2_agents}}; do
        h2 stop "$agent" 2>/dev/null || true
    done
    tmux kill-session -t "$session" 2>/dev/null || true
    # Build commands
    cmd_manager="h2 run manager --role manager --var working_dir='$wd'; read"
    cmd_auditor="h2 run auditor --role code-auditor --var working_dir='$wd' --var language='$lang'; read"
    cmd_docs_updater="h2 run docs-updater --role docs-writer --var working_dir='$wd' --var doc_framework='$framework'; read"
    cmd_docs_creator="h2 run docs-creator --role docs-writer --var working_dir='$wd' --var doc_framework='$framework'; read"
    cmd_icon_designer="h2 run icon-designer --role icon-designer --var working_dir='$wd'; read"
    # Launch in tmux
    tmux new-session -d -s "$session" -n manager "$cmd_manager"
    tmux new-window -t "$session" -n auditor "$cmd_auditor"
    tmux new-window -t "$session" -n docs-updater "$cmd_docs_updater"
    tmux new-window -t "$session" -n docs-creator "$cmd_docs_creator"
    tmux new-window -t "$session" -n icon-designer "$cmd_icon_designer"
    tmux select-window -t "$session":manager
    echo "Agents launched in tmux session '$session'."

# Run a named project or send a free-form message to the manager
h2-run +args:
    #!/usr/bin/env bash
    first="{{args}}"
    if [[ -f ".h2/projects/$first.md" ]]; then
        h2 send manager "Read and execute the project brief at .h2/projects/$first.md"
    else
        h2 send manager "$first"
    fi

# Pick a recent Claude plan and send it to the manager
h2-plan:
    #!/usr/bin/env bash
    plans_dir="$HOME/.claude/plans"
    if [[ ! -d "$plans_dir" ]]; then
        echo "No plans directory found at $plans_dir"
        exit 1
    fi
    # List 10 most recent plans with titles
    mapfile -t files < <(ls -t "$plans_dir"/*.md 2>/dev/null | head -10)
    if [[ ${#files[@]} -eq 0 ]]; then
        echo "No plans found."
        exit 1
    fi
    echo "Recent plans:"
    echo ""
    for i in "${!files[@]}"; do
        f="${files[$i]}"
        title=$(head -1 "$f" | sed 's/^#\+\s*//')
        date=$(stat -f '%Sm' -t '%Y-%m-%d %H:%M' "$f")
        printf "  %2d. [%s] %s\n" $((i+1)) "$date" "$title"
    done
    echo ""
    read -rp "Pick a plan (number): " choice
    if [[ -z "$choice" ]] || [[ "$choice" -lt 1 ]] || [[ "$choice" -gt ${#files[@]} ]]; then
        echo "Invalid choice."
        exit 1
    fi
    selected="${files[$((choice-1))]}"
    echo "Sending to manager: $(basename "$selected")"
    h2 send manager "Read and execute the plan at $selected"

# Send a message directly to a specific agent
h2-tell agent +message:
    h2 send {{agent}} "{{message}}"

# View h2-agents in current terminal (grouped session for independent window switching)
h2-attach:
    #!/usr/bin/env bash
    session="{{h2_session}}"
    if ! tmux has-session -t "$session" 2>/dev/null; then
        echo "No '$session' session found. Run 'just h2-start' first."
        exit 1
    fi
    if [[ -n "$TMUX" ]]; then
        tmux switch-client -t "$session"
    else
        # Create a grouped session so this terminal can view h2-agents
        # windows independently from any other terminal
        tmux new-session -t "$session" -s "${session}-view" 2>/dev/null \
            || tmux attach -t "${session}-view"
    fi

# View h2-agents in a separate terminal (run from a non-tmux terminal)
h2-view:
    #!/usr/bin/env bash
    session="{{h2_session}}"
    if ! tmux has-session -t "$session" 2>/dev/null; then
        echo "No '$session' session found. Run 'just h2-start' first."
        exit 1
    fi
    tmux attach -t "$session"

# Stop all h2 agents and kill the tmux session
h2-stop:
    #!/usr/bin/env bash
    for agent in {{h2_agents}}; do
        h2 stop "$agent" 2>/dev/null || true
    done
    tmux kill-session -t {{h2_session}}-view 2>/dev/null || true
    tmux kill-session -t {{h2_session}} 2>/dev/null || true

# List available projects
h2-list:
    @ls -1 .h2/projects/*.md 2>/dev/null | sed 's|.h2/projects/||;s|\.md$||' || echo "No projects found in .h2/projects/"
