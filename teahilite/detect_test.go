package teahilite

import "testing"

func TestDetectLanguage_KnownExtensions(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"main.go", "go"},
		{"script.py", "python"},
		{"style.css", "css"},
		{"index.html", "html"},
		{"app.js", "js"},
		{"data.json", "json"},
		{"config.yaml", "yaml"},
		{"README.md", "md"},
		{"query.sql", "mysql"},
		{"main.rs", "rust"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := DetectLanguage(tt.path)
			if got != tt.want {
				t.Errorf("DetectLanguage(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestDetectLanguage_FallbackToText(t *testing.T) {
	tests := []string{
		"noext",
		"file.unknownext123",
		"",
	}
	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			got := DetectLanguage(path)
			if got != "text" {
				t.Errorf("DetectLanguage(%q) = %q, want %q", path, got, "text")
			}
		})
	}
}

// stringAlias verifies the generic ~string constraint works with custom types.
type stringAlias string

func TestDetectLanguage_GenericStringType(t *testing.T) {
	path := stringAlias("main.go")
	got := DetectLanguage(path)
	if got != "go" {
		t.Errorf("DetectLanguage(stringAlias(%q)) = %q, want %q", path, got, "go")
	}
}
