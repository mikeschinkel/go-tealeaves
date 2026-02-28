.PHONY: help build test tidy fmt vet clean build-examples

# Modules in this repo
MODULES = teadd teadep teagrid teamodal teanotify teastatus teatextsel teatree teautils

# Example modules (each is an independent Go module)
EXAMPLES = \
	examples/teadd/simple \
	examples/teadd/demo \
	examples/teadep/treenav \
	examples/teamodal/choices \
	examples/teamodal/editlist \
	examples/teamodal/various \
	examples/teanotify/simple \
	examples/teastatus/statusbar \
	examples/teatextsel/editor \
	examples/teatree/filetree \
	examples/teagrid/simplest \
	examples/teagrid/sorting \
	examples/teagrid/filtering \
	examples/teagrid/scrolling \
	examples/teautils/keyhelp

# Default target
help:
	@echo "Available targets:"
	@echo "  make help            - Show this help message"
	@echo "  make build           - Build color-viewer to ./bin/"
	@echo "  make test            - Run tests across all modules"
	@echo "  make tidy            - Run go mod tidy across all modules and examples"
	@echo "  make fmt             - Format code with gofmt"
	@echo "  make vet             - Run go vet across all modules"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make build-examples  - Build all example programs to ./bin/examples/"

# Build color-viewer
build:
	@echo "Building color-viewer..."
	@mkdir -p bin
	@cd cmd/color-viewer && go build -o ../../bin/color-viewer . || exit 1
	@echo "Built to ./bin/color-viewer"

# Run tests across all modules
test:
	@for mod in $(MODULES); do \
		echo "Testing $$mod..."; \
		cd $$mod && go test ./... || exit 1; \
		cd ..; \
	done
	@echo "All tests passed!"

# Run go mod tidy across all modules and examples
tidy:
	@for mod in $(MODULES); do \
		echo "Tidying $$mod..."; \
		cd $$mod && go mod tidy || exit 1; \
		cd ..; \
	done
	@echo "Tidying cmd/color-viewer..."
	@cd cmd/color-viewer && go mod tidy || exit 1
	@for ex in $(EXAMPLES); do \
		echo "Tidying $$ex..."; \
		cd $$ex && go mod tidy || exit 1; \
		cd ../../..; \
	done
	@echo "All modules tidied!"

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# Run go vet across all modules
vet:
	@for mod in $(MODULES); do \
		echo "Vetting $$mod..."; \
		cd $$mod && go vet ./... || exit 1; \
		cd ..; \
	done
	@echo "All modules vetted!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin
	@for mod in $(MODULES); do \
		cd $$mod && go clean; \
		cd ..; \
	done
	@cd cmd/color-viewer && go clean
	@echo "Clean complete!"

# Build all example programs
build-examples:
	@mkdir -p bin/examples
	@for ex in $(EXAMPLES); do \
		name=$$(basename $$ex); \
		parent=$$(basename $$(dirname $$ex)); \
		echo "Building $$parent/$$name..."; \
		cd $$ex && go build -o ../../../bin/examples/$$parent-$$name . || exit 1; \
		cd ../../..; \
	done
	@echo "All examples built to ./bin/examples/"
