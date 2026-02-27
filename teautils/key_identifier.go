package teautils

import (
	"fmt"
	"strings"
	"unicode"
)

// KeyIdentifier is a unique identifier for a key binding in the registry
// Format: "namespace.key" or "namespace.subnamespace.key"
// Examples: "app.help", "file-intent.commit", "tree.nav.up"
type KeyIdentifier string

// ParseKeyIdentifier validates and creates a KeyIdentifier
// Returns error if format is invalid
func ParseKeyIdentifier(s string) (id KeyIdentifier, err error) {
	if s == "" {
		err = ErrEmptyKeyIdentifier
		goto end
	}

	// Must contain at least one dot
	if !strings.Contains(s, ".") {
		err = NewErr(ErrKeyIdentifierMissingDot, "identifier", s)
		goto end
	}

	// Split by dots and validate each part
	{
		parts := strings.Split(s, ".")
		for i, part := range parts {
			if part == "" {
				err = NewErr(ErrKeyIdentifierEmptyPart,
					"identifier", s,
					"position", i)
				goto end
			}

			// Validate each part is a valid identifier component
			if !isValidKeyIdentifierPart(part) {
				err = NewErr(ErrKeyIdentifierInvalidPart,
					"identifier", s,
					"part", part,
					"position", i)
				goto end
			}
		}
	}

	id = KeyIdentifier(s)

end:
	return id, err
}

// MustParseKeyIdentifier is like ParseKeyIdentifier but panics on error
// Use in initialization code where identifiers are known to be valid
func MustParseKeyIdentifier(s string) KeyIdentifier {
	id, err := ParseKeyIdentifier(s)
	if err != nil {
		panic(fmt.Sprintf("teautils: invalid key identifier %q: %v", s, err))
	}
	return id
}

// String returns the string representation
func (k KeyIdentifier) String() string {
	return string(k)
}

// isValidKeyIdentifierPart checks if a part is a valid identifier component
// Allows: alphanumeric, hyphens, underscores
// Must start with letter or digit
func isValidKeyIdentifierPart(s string) bool {
	if len(s) == 0 {
		return false
	}

	// First character must be alphanumeric
	first := rune(s[0])
	if !unicode.IsLetter(first) && !unicode.IsDigit(first) {
		return false
	}

	// Rest can be alphanumeric, hyphen, or underscore
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}

	return true
}
