package teautils

import "errors"

var (
	// KeyIdentifier errors
	ErrEmptyKeyIdentifier       = errors.New("empty key identifier")
	ErrKeyIdentifierMissingDot  = errors.New("key identifier must contain at least one dot separator")
	ErrKeyIdentifierEmptyPart   = errors.New("key identifier contains empty part")
	ErrKeyIdentifierInvalidPart = errors.New("key identifier contains invalid part")

	// KeyRegistry errors
	ErrKeyNotFound = errors.New("key not found in registry")
	ErrEmptyKeyID  = errors.New("cannot register key with empty identifier")
)
