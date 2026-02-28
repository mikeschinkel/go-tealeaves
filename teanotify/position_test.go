package teanotify

import "testing"

func TestPosition_IsValid(t *testing.T) {
	t.Run("TopLeft", func(t *testing.T) {
		if !TopLeftPosition.IsValid() {
			t.Error("expected TopLeftPosition to be valid")
		}
	})
	t.Run("TopCenter", func(t *testing.T) {
		if !TopCenterPosition.IsValid() {
			t.Error("expected TopCenterPosition to be valid")
		}
	})
	t.Run("TopRight", func(t *testing.T) {
		if !TopRightPosition.IsValid() {
			t.Error("expected TopRightPosition to be valid")
		}
	})
	t.Run("BottomLeft", func(t *testing.T) {
		if !BottomLeftPosition.IsValid() {
			t.Error("expected BottomLeftPosition to be valid")
		}
	})
	t.Run("BottomCenter", func(t *testing.T) {
		if !BottomCenterPosition.IsValid() {
			t.Error("expected BottomCenterPosition to be valid")
		}
	})
	t.Run("BottomRight", func(t *testing.T) {
		if !BottomRightPosition.IsValid() {
			t.Error("expected BottomRightPosition to be valid")
		}
	})
	t.Run("Unspecified", func(t *testing.T) {
		if UnspecifiedPosition.IsValid() {
			t.Error("expected UnspecifiedPosition to be invalid")
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		if Position("XX").IsValid() {
			t.Error("expected Position(\"XX\") to be invalid")
		}
	})
}

func TestPosition_String(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want string
	}{
		{"TopLeft", TopLeftPosition, "top-left"},
		{"TopCenter", TopCenterPosition, "top-center"},
		{"TopRight", TopRightPosition, "top-right"},
		{"BottomLeft", BottomLeftPosition, "bottom-left"},
		{"BottomCenter", BottomCenterPosition, "bottom-center"},
		{"BottomRight", BottomRightPosition, "bottom-right"},
		{"Unspecified", UnspecifiedPosition, "unknown"},
		{"Invalid", Position("XX"), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.String()
			if got != tt.want {
				t.Errorf("Position(%q).String() = %q, want %q", tt.pos, got, tt.want)
			}
		})
	}
}

func TestPosition_Label(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want string
	}{
		{"TopLeft", TopLeftPosition, "Top Left"},
		{"TopCenter", TopCenterPosition, "Top Center"},
		{"TopRight", TopRightPosition, "Top Right"},
		{"BottomLeft", BottomLeftPosition, "Bottom Left"},
		{"BottomCenter", BottomCenterPosition, "Bottom Center"},
		{"BottomRight", BottomRightPosition, "Bottom Right"},
		{"Unspecified", UnspecifiedPosition, "Unknown"},
		{"Invalid", Position("XX"), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.Label()
			if got != tt.want {
				t.Errorf("Position(%q).Label() = %q, want %q", tt.pos, got, tt.want)
			}
		})
	}
}
