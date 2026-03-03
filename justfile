# Modules in this repo
modules := "teadrpdwn teadepview teagrid teamodal teanotify teastatus teatxtsnip teatree teautils"

# Example modules (each is an independent Go module)
examples := "examples/teadrpdwn/simple examples/teadrpdwn/demo examples/teadepview/treenav examples/teamodal/choices examples/teamodal/editlist examples/teamodal/various examples/teanotify/simple examples/teastatus/statusbar examples/teatxtsnip/editor examples/teatree/filetree examples/teagrid/simplest examples/teagrid/sorting examples/teagrid/filtering examples/teagrid/scrolling examples/teautils/keyhelp"

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
