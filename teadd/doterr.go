// MIT License
//
// Copyright (c) 2025 Mike Schinkel <mike@newclarity.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// doterr.go provides a tiny, stdlib-only way to attach structured metadata
// and sentinels to errors while staying fully composable with the Go standard
// library. The model is:
//
//   - Each function builds an entry with New(...) passing an optional trailing cause:
//     return NewErr(ErrRepo, "key", val, cause) // cause last
//
//   - With(...) is a flexible convenience for same-function enrichment. It can:
//
//   - Enrich the rightmost doterr entry within an existing joined error, or
//
//   - Join a new entry if no doterr entry exists, and
//
//   - (Optionally) treat a final trailing error as the cause and join it last.
//
//   - Combine([]error) bundles independent failures into a single error that unwraps
//     to its members, preserving order.

package teadd

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// ErrKV represents a key/value metadata pair. Keys are preserved in
// insertion order, and values may be of any type.
type ErrKV interface {
	Key() string
	Value() any
}

// Typed ErrKV constructors (inspired by log/slog.Attr constructors)

// StringKV creates a KV with a string value.
func StringKV(key, value string) ErrKV {
	return kv{k: key, v: value}
}

// IntKV creates a KV with an int value.
func IntKV(key string, value int) ErrKV {
	return kv{k: key, v: value}
}

// Int64KV creates a KV with an int64 value.
func Int64KV(key string, value int64) ErrKV {
	return kv{k: key, v: value}
}

// BoolKV creates a KV with a bool value.
func BoolKV(key string, value bool) ErrKV {
	return kv{k: key, v: value}
}

// Float64KV creates a KV with a float64 value.
func Float64KV(key string, value float64) ErrKV {
	return kv{k: key, v: value}
}

// AnyKV creates a KV with any value type.
// Use this for custom types or when the specific type constructor doesn't exist.
func AnyKV(key string, value any) ErrKV {
	return kv{k: key, v: value}
}

// ErrorKV creates a KV with an error value.
// Use this when you want to include an error as metadata (not as a cause).
func ErrorKV(key string, value error) ErrKV {
	return kv{k: key, v: value}
}

// Sentinel errors for validation failures
var (
	ErrMissingSentinel     = errors.New("missing required sentinel error")
	ErrTrailingKey         = errors.New("trailing key without value")
	ErrMisplacedError      = errors.New("error in wrong position")
	ErrInvalidArgumentType = errors.New("invalid argument type")
	ErrOddKeyValueCount    = errors.New("odd number of key-value arguments")
	ErrCrossPackageError   = errors.New("error from different doterr package")
)

// NewErr builds a standalone structured entry (no primary cause inside).
// Accepted parts:
//   - error         — sentinel/tag (required: at least one, must be first)
//   - KV{Key,Value} — explicit key/value
//   - "key", value  — implicit pair (value can be any type, including error)
//   - error         — optional trailing cause (joined last via errors.Join)
//
// Pattern: one or more sentinels (error), then zero or more key-value pairs,
// then optional trailing cause (error). After the first string key, all
// remaining args must form valid pairs, except for an optional final error.
// Returns nil if no meaningful parts are provided after validation.
// Returns a validation error joined with the partial entry if validation fails.
//
//goland:noinspection DuplicatedCode
func NewErr(parts ...any) error {
	// Separate optional trailing cause from the parts
	//goland:noinspection DuplicatedCode
	coreParts, cause := extractTrailingCause(parts)

	if validationErr := validateNewParts(coreParts); validationErr != nil {
		// Return validation error as first error in entry.errors
		var e entry
		e.id = uniqueId
		appendEntry(&e, coreParts...)
		if e.empty() {
			return validationErr
		}
		// Prepend validation error to entry.errors
		e.errors = append([]error{validationErr}, e.errors...)
		return e
	}

	var e entry
	e.id = uniqueId
	appendEntry(&e, coreParts...)
	if e.empty() {
		return cause // if we only had a cause, return it
	}

	// Append optional cause to entry.errors (cause last)
	if cause != nil {
		e.errors = append(e.errors, cause)
	}
	return e
}

// MsgErr creates an ad-hoc error message without requiring a sentinel error.
// This is a convenience for rapid development - use sentinels for production code.
//
// Accepts two forms:
//   - MsgErr("message") - creates error with the given message
//   - MsgErr(err) - wraps existing error, preserving error chain for errors.Is()
//
// Use this when you don't want to define a sentinel immediately. Future tooling
// can detect MsgErr usage and suggest/generate appropriate sentinel errors.
//
// For errors with metadata, use NewErr with a sentinel, or wrap with WithErr:
//
//	err := MsgErr("config invalid")
//	err = WithErr(err, "path", configPath)
func MsgErr(msg any) error {
	switch v := msg.(type) {
	case string:
		return msgErr{msg: v}
	case error:
		// Wrap the error, preserving error chain
		return msgErr{wrapped: v}
	default:
		// Invalid input - panic in development
		panic(fmt.Sprintf("MsgErr: invalid type %T, expected string or error", v))
	}
}

// WithErr is a flexible enrichment helper. Typical uses:
//
//	// Enrich an existing composite error (err may be an errors.Join tree):
//	err = WithErr(err, "Foo", 10)
//
//	// Build an entry and join a trailing cause in one shot:
//	err = WithErr("endpoint", ep, ErrTemplate, cause) // 'cause' is last
//
// Behavior:
//
//  1. If the FIRST arg is an error, it is treated as the base error to enrich:
//     • If it is a doterr entry, merge KVs/sentinels into that entry.
//     • If it is a multi-unwrap (errors.Join tree), find the RIGHTMOST
//     doterr entry, merge into it, and rebuild preserving order.
//     • If no doterr entry is found, a new entry will be joined in (see step 3).
//
//  2. After consuming the base (if present), if the LAST remaining arg is an error,
//     it is treated as the CAUSE and joined LAST.
//
//  3. The remaining middle args (if any) are collected into an entry. If we
//     enriched an existing doterr entry in step 1, that merged entry is used;
//     otherwise, a fresh entry is created. If there is a trailing CAUSE from step 2,
//     the result is errors.Join(entry, cause). If there is no cause, the entry is returned.
//
// Note: For inter-function composition, prefer New() with trailing cause:
//
//	return NewErr(ErrRepo, "key", val, cause) // cause last
//
//goland:noinspection DuplicatedCode
func WithErr(parts ...any) error {
	if len(parts) == 0 {
		return nil
	}

	i, j := 0, len(parts)-1

	// Optional base (first arg)
	var baseErr error
	firstErr, ok := parts[i].(error)
	if ok {
		baseErr = checkCrossPackage(firstErr)
		i++
	}

	// Optional cause (last arg)
	var cause error
	if j >= i {
		lastErr, ok := parts[j].(error)
		if ok {
			// Special case: if we only have ONE part left (a single error),
			// treat it as a sentinel for new context, NOT as a cause.
			// This enables: WithErr(err, ErrNewContext) to add context layers.
			// To add a cause without context, use: NewErr(..., cause) or errors.Join
			if i == j {
				// Single error - it's a context sentinel, not a cause
				// Don't extract it as cause; let it go into middle
			} else {
				// Multiple parts - last error is a cause
				cause = checkCrossPackage(lastErr)
				j--
			}
		}
	}

	// Middle segment are the metadata/sentinels to apply.
	middle := parts[i : j+1]

	// No base error: build entry from middle, then (if present) join cause LAST.
	if baseErr == nil {
		return handleCause(buildEntry(middle...), cause)
	}

	// Have a base error: try to enrich rightmost entry or join a fresh entry.
	err := buildErr(baseErr, middle)

	// Now handle the cause
	return handleCause(err, cause)

}

// CombineErrs bundles a slice of errors into a single composite error that unwraps
// to its members. Order is preserved and nils are skipped. Returns nil for an
// empty/fully-nil slice, or the sole error when there is exactly one.
func CombineErrs(errs []error) error {
	filtered := make([]error, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			filtered = append(filtered, e)
		}
	}
	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		return combined{errs: filtered}
	}
}

// ErrMeta returns the key/value pairs stored on all doterr entries in the error tree.
// If err is a doterr entry, returns its metadata.
// If err is a joined error (has Unwrap() []error), scans all children recursively
// and collects metadata from all doterr entries found, preserving order.
// Otherwise returns nil.
// The returned slice preserves insertion order across all entries.
//
//goland:noinspection DuplicatedCode
func ErrMeta(err error) []ErrKV {
	var ok bool
	var e, ce entry
	var cePtr *entry

	// Case (a): err is an entry → return its metadata plus metadata from errors in entry.errors
	//goland:noinspection GoTypeAssertionOnErrors,DuplicatedCode
	e, ok = err.(entry)
	if ok {
		out := make([]ErrKV, len(e.kvs))
		for i, pair := range e.kvs {
			out[i] = pair
		}
		// Also collect metadata from errors stored in entry.errors
		for _, childErr := range e.errors {
			childMeta := ErrMeta(childErr)
			if childMeta != nil {
				out = append(out, childMeta...)
			}
		}
		return out
	}

	// Case (b): err is a join (multi-unwrap) → scan all children and collect metadata
	type unwrapper interface{ Unwrap() []error }
	u, ok := err.(unwrapper)
	if !ok {
		return nil
	}

	// Collect metadata from all doterr entries in the tree
	var allKVs []ErrKV
	children := u.Unwrap()
	for _, child := range children {
		// Check if child is an entry
		//goland:noinspection GoTypeAssertionOnErrors
		ce, ok = child.(entry)
		if !ok {
			//goland:noinspection GoTypeAssertionOnErrors
			cePtr, ok = child.(*entry)
			if ok {
				ce = *cePtr
			}
		}
		if ok {
			// Add this entry's metadata
			for _, pair := range ce.kvs {
				allKVs = append(allKVs, pair)
			}
		} else {
			// Recursively check if child is itself a join
			childMeta := ErrMeta(child)
			if childMeta != nil {
				allKVs = append(allKVs, childMeta...)
			}
		}
	}

	if len(allKVs) > 0 {
		return allKVs
	}
	return nil
}

// ErrValue extracts a single metadata value by key with type safety.
// Returns the value and true if found and the value is of type T.
// Returns the zero value of T and false if not found or type mismatch.
//
// Example:
//
//	status, ok := ErrValue[int](err, "http_status")
//	if ok {
//	    fmt.Printf("Status: %d\n", status)
//	}
//
//	name, ok := ErrValue[string](err, "parameter_name")
func ErrValue[T any](err error, key string) (T, bool) {
	var zero T
	kvs := ErrMeta(err)
	if kvs == nil {
		return zero, false
	}

	for _, pair := range kvs {
		if pair.Key() == key {
			if val, ok := pair.Value().(T); ok {
				return val, true
			}
			return zero, false
		}
	}
	return zero, false
}

// Errors returns the errors stored on a doterr entry.
// If err is a doterr entry, returns its errors.
// If err is a joined error (has Unwrap() []error), scans immediate children
// left-to-right and returns errors from the first doterr entry found.
// Otherwise returns nil.
// The returned slice preserves insertion order and is a copy.
//
// Note: These errors may be sentinel errors (e.g., ErrRepo), custom error types
// (e.g., *rfc9457.Error), or any other error type stored in the entry.
//
//goland:noinspection DuplicatedCode
func Errors(err error) []error {
	// Case (a): err is an entry → return its errors
	//goland:noinspection GoTypeAssertionOnErrors
	e, ok := err.(entry)
	if ok {
		out := make([]error, len(e.errors))
		copy(out, e.errors)
		return out
	}

	// Case (b): err is a join (multi-unwrap) → scan immediate children left-to-right
	type unwrapper interface{ Unwrap() []error }
	u, ok := err.(unwrapper)
	if !ok {
		return nil
	}
	children := u.Unwrap()
	for _, child := range children {
		//goland:noinspection GoTypeAssertionOnErrors
		if ce, ok := child.(entry); ok {
			out := make([]error, len(ce.errors))
			copy(out, ce.errors)
			return out
		}
	}

	return nil
}

// FindErr walks an error tree (including errors.Join trees)
// and returns the first match for target (via errors.As).
func FindErr[T error](err error) (out T, ok bool) {
	var t T
	if errors.As(err, &t) {
		return t, true
	}
	return out, false
}

func AppendErr(errs []error, err error) []error {
	if err == nil {
		return errs
	}
	return append(errs, err)
}

// AppendKV appends key-value pairs to a slice of KV values.
// It accepts variadic arguments in three forms:
//  1. Individual KV values: AppendKV(kvs, StringKV("foo", "bar"))
//  2. String key-value pairs: AppendKV(kvs, "foo", "bar")
//  3. Lazy ErrKV functions: AppendKV(kvs, func()ErrKV{ return StringKV("expensive", compute()) })
//
// This function is useful for accumulating metadata throughout a function
// before creating an error:
//
//	var kvs []ErrKV
//	kvs = AppendKV(kvs, "user_id", userID)
//	if complexCondition {
//	    kvs = AppendKV(kvs, "reason", "complex")
//	}
//	return NewErr(ErrFailed, kvs, cause)
//
// String key-value pairs must come in even counts (odd counts panic in debug mode).
// Lazy functions (func()ErrKV) are wrapped and not evaluated until error creation.
func AppendKV(kvs []ErrKV, parts ...any) []ErrKV {
	for i := 0; i < len(parts); {
		switch v := parts[i].(type) {
		case ErrKV:
			kvs = append(kvs, v)
			i++
		case func() ErrKV:
			// Wrap in lazyKV - evaluation deferred until error creation
			kvs = append(kvs, &lazyKV{fn: v})
			i++
		case string:
			// String key must have a value
			if i+1 >= len(parts) {
				panic("AppendKV: trailing key without value: " + v)
			}
			kvs = append(kvs, kv{k: v, v: parts[i+1]})
			i += 2
		default:
			// Skip invalid types with panic in debug mode
			panic(fmt.Sprintf("AppendKV: invalid type %T at position %d", v, i))
		}
	}
	return kvs
}

//--------------------------------
// Unexported implementation types
//--------------------------------

// kv is the internal implementation of the KV interface.
type kv struct {
	k string
	v any
}

func (p kv) Key() string { return p.k }
func (p kv) Value() any  { return p.v }

// lazyKV defers evaluation of a KV until accessed.
// Evaluation happens at error creation time (when Key() or Value() is called),
// not at AppendKV time, to avoid expensive operations in the happy path.
type lazyKV struct {
	fn     func() ErrKV
	once   sync.Once
	result kv
}

func (l *lazyKV) Key() string {
	l.once.Do(l.evaluate)
	return l.result.k
}

func (l *lazyKV) Value() any {
	l.once.Do(l.evaluate)
	return l.result.v
}

func (l *lazyKV) evaluate() {
	kvResult := l.fn()
	l.result = kv{k: kvResult.Key(), v: kvResult.Value()}
}

// msgErr represents an ad-hoc error message created with MsgErr.
// It serves as a placeholder until sentinels are defined, and can be
// detected by tooling to suggest/generate sentinel errors.
type msgErr struct {
	msg     string // The error message (if created from string)
	wrapped error  // The wrapped error (if created from error)
}

func (m msgErr) Error() string {
	if m.wrapped != nil {
		return m.wrapped.Error()
	}
	return m.msg
}

func (m msgErr) Unwrap() error {
	return m.wrapped
}

var uniqueId = func() int {
	var b [8]byte
	// crypto/rand.Read never fails on supported platforms
	_, _ = rand.Read(b[:])
	return int(binary.LittleEndian.Uint64(b[:]))
}()

// entry represents one function's contribution to an error chain.
// Each function creates one entry with errors (sentinels, custom typed errors) and metadata.
// It implements error and Unwrap() []error.
type entry struct {
	id     int     // Unique ID
	errors []error // sentinels, custom typed errors (NOT the primary cause)
	kvs    []kv    // metadata in insertion order
}

func newEntry(errors []error, kvs []kv) *entry {
	return &entry{
		id:     uniqueId,
		errors: errors,
		kvs:    kvs,
	}
}

// Error implements the std lib error interface and returns the strings value
func (e entry) Error() (s string) {
	if len(e.errors) == 0 && len(e.kvs) == 0 {
		return "doterr{}"
	}

	sb := strings.Builder{}

	// Include sentinel errors first (recursively passing seen map)
	for i, err := range e.errors {
		if i != 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	if len(e.kvs) == 0 {
		return sb.String()
	}
	var meta strings.Builder

	// Then include metadata (skip if already seen)
	for i, pair := range e.kvs {
		if e.MatchKV(pair.k, pair.v) {
			// We've already added to our error string
			continue
		}
		if i != 0 {
			meta.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&meta, " %s=%v", pair.k, pair.v)
	}
	s = sb.String()
	if meta.Len() > 0 {
		s += "; meta:" + meta.String()
	}
	return
}

var _ IsDotErrEntry = (*entry)(nil)

type IsDotErrEntry = interface {
	IsDotErrEntry()
	HasKV(key string, value any) bool
	MatchKV(key string, value any) bool
}

func (e entry) IsDotErrEntry() {}

// HasKV checks if this entry contains a key-value pair.
// Zero allocation query method for cross-package metadata inspection.
func (e entry) HasKV(key string, value any) bool {
	for _, kv := range e.kvs {
		if kv.k == key && kv.v == value {
			return true
		}
	}
	return false
}

// MatchKV checks if a key-value pair exists in any descendant entry (not in self).
// Returns true if any child or deeper descendant has this exact key-value pair.
func (e entry) MatchKV(key string, value any) (matched bool) {
	for _, err := range e.errors {
		// Assert to full anonymous interface (works across packages)
		child, ok := err.(interface {
			IsDotErrEntry()
			HasKV(key string, value any) bool
			MatchKV(key string, value any) bool
		})
		if !ok {
			continue
		}
		// Check if child has this pair (zero allocation query)
		if child.HasKV(key, value) {
			return true
		}
		// Recursively check child's descendants
		if child.MatchKV(key, value) {
			return true
		}
	}
	return false
}

func (e entry) Unwrap() []error {
	if len(e.errors) == 0 {
		return nil
	}
	cp := make([]error, len(e.errors))
	copy(cp, e.errors)
	return cp
}

func (e entry) empty() bool { return len(e.errors) == 0 && len(e.kvs) == 0 }

//goland:noinspection DuplicatedCode
func appendEntry(e *entry, parts ...any) {
	for i := 0; i < len(parts); {
		switch v := parts[i].(type) {
		case ErrKV:
			// Handles both regular kv and lazyKV (both implement ErrKV interface)
			// Calling Key()/Value() triggers lazy evaluation if needed
			e.kvs = append(e.kvs, kv{k: v.Key(), v: v.Value()})
			i++
		case []ErrKV:
			// Handle slice of ErrKV (from AppendKV)
			for _, kvItem := range v {
				e.kvs = append(e.kvs, kv{k: kvItem.Key(), v: kvItem.Value()})
			}
			i++
		case func() ErrKV:
			// Handle direct func()ErrKV (not wrapped in slice)
			lazyKV := v()
			e.kvs = append(e.kvs, kv{k: lazyKV.Key(), v: lazyKV.Value()})
			i++
		case string:
			if i+1 < len(parts) {
				e.kvs = append(e.kvs, kv{k: v, v: parts[i+1]})
				i += 2
			} else {
				// Trailing key without value: skip it (validation should have caught this)
				i++
			}
		case error:
			if v != nil {
				e.errors = append(e.errors, v) // sentinel/tag
			}
			i++
		default:
			// Non-string, non-error, non-ErrKV: skip it (validation should have caught this)
			i++
		}
	}
}

// combined implements a composite error for Combine().
type combined struct{ errs []error }

func (c combined) Error() string {
	var messages []string
	for _, err := range c.errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}
	return strings.Join(messages, "\n")
}

func (c combined) Unwrap() []error {
	cp := make([]error, len(c.errs))
	copy(cp, c.errs)
	return cp
}

//------------------------
// Unexported helper funcs
//------------------------

// addToEntryErrors appends additional errors to an entry's errors field.
// If err is an entry, it appends additionalErrs to its errors slice.
// If err is not an entry, it creates a new entry containing err and additionalErrs.
// Returns nil if both err and additionalErrs are empty.
func addToEntryErrors(err error, additionalErrs ...error) error {
	// Filter out nil errors
	filtered := make([]error, 0, len(additionalErrs))
	for _, e := range additionalErrs {
		if e != nil {
			filtered = append(filtered, e)
		}
	}

	if err == nil {
		if len(filtered) == 0 {
			return nil
		}
		if len(filtered) == 1 {
			return filtered[0]
		}
		var e entry
		e.id = uniqueId
		e.errors = filtered
		return e
	}

	if len(filtered) == 0 {
		return err
	}

	// Try to add to existing entry
	if e, ok := err.(entry); ok {
		e.errors = append(e.errors, filtered...)
		return e
	}

	// err is not an entry - create new entry wrapping all
	var e entry
	e.id = uniqueId
	e.errors = make([]error, 0, len(filtered)+1)
	e.errors = append(e.errors, err)
	e.errors = append(e.errors, filtered...)
	return e
}

// extractTrailingCause checks if the last element in parts is an error that should
// be treated as a trailing cause. Returns (cause, remaining parts).
// A trailing error is considered a cause if:
//   - It's the last element, AND
//   - It comes after at least one sentinel, AND
//   - It's not part of an incomplete key-value pair (would leave odd count)
//
//goland:noinspection DuplicatedCode
func extractTrailingCause(parts []any) (_ []any, err error) {
	var lastIdx, sentinelCount, nonSentinelCount int
	var ok bool

	if len(parts) == 0 {
		goto end
	}

	lastIdx = len(parts) - 1
	err, ok = parts[lastIdx].(error)
	if !ok {
		goto end
	}

	// Count sentinels at the beginning (but don't count entry as sentinel)
	sentinelCount = 0
	for i := 0; i < len(parts); i++ {
		err, ok := parts[i].(error)
		if !ok {
			break
		}
		// entry should not count as a sentinel - it's a wrapped error
		//goland:noinspection GoTypeAssertionOnErrors
		_, ok = err.(entry)
		if ok {
			break
		}
		sentinelCount++
	}

	// If we only have errors (all sentinels), the last one is not a cause
	if sentinelCount == len(parts) {
		err = nil // Not a cause - clear err to avoid returning it twice
		goto end
	}

	// Check if removing the last error would leave an odd number of non-sentinel args
	// (which would mean the error is actually a value for a key)
	// Count non-sentinel args, treating ErrKV/[]ErrKV/func()ErrKV as single items
	nonSentinelCount = 0
	for i := sentinelCount; i < len(parts)-1; i++ {
		switch v := parts[i].(type) {
		case ErrKV, func() ErrKV:
			// These are "atomic" - count as single items
			nonSentinelCount += 2
		case []ErrKV:
			// These are "atomic" - count as single items
			nonSentinelCount += 2 * len(v)
		case string:
			// String key + value = 2 items
			nonSentinelCount++
			if i+1 < len(parts)-1 {
				i++ // Skip value
				nonSentinelCount++
			}
		}
	}

	if nonSentinelCount%2 != 0 {
		// Odd count - last error is a value, not a cause
		err = nil // Not a cause - clear err to avoid returning it twice
		goto end
	}

	parts = parts[:lastIdx]
end:
	// Last error is a trailing cause
	return parts, err
}

// validateNewParts validates that parts conforms to the expected pattern for NewErr:
//   - At least one sentinel error (required)
//   - Then zero or more KV or string key-value pairs
//   - After the first string key, remaining args must form valid pairs (even count)
//   - No non-string/non-error/non-ErrKV values allowed
//
//goland:noinspection DuplicatedCode
func validateNewParts(parts []any) error {
	if len(parts) == 0 {
		// Build entry manually to avoid recursion
		// TODO: We need to review the validity of this after introductin of MsgErr()
		e := newEntry([]error{ErrMissingSentinel}, []kv{
			{k: "message", v: "NewErr requires at least one sentinel error"},
		})
		return e
	}

	// Find where sentinels end (first non-error)
	sentinelCount := 0
	i := 0
	for i < len(parts) {
		if _, ok := parts[i].(error); ok {
			sentinelCount++
			i++
		} else {
			break
		}
	}

	// Must have at least one sentinel
	if sentinelCount == 0 {
		// Build entry manually to avoid recursion
		// TODO: We need to review the validity of this after introductin of MsgErr()
		e := newEntry([]error{ErrMissingSentinel}, []kv{
			{k: "message", v: "NewErr requires at least one sentinel error as the first argument"},
		})
		return e
	}

	// Validate remaining args are valid key-value pairs
	firstKeyIdx := -1
	for j := i; j < len(parts); j++ {
		switch v := parts[j].(type) {
		case ErrKV:
			// Explicit ErrKV is always valid
			continue
		case []ErrKV:
			// Slice of ErrKV is always valid (doesn't affect parity)
			continue
		case func() ErrKV:
			// Lazy ErrKV function is always valid (doesn't affect parity)
			continue
		case string:
			// Mark first key position if not yet found
			if firstKeyIdx == -1 {
				firstKeyIdx = j
			}
			// After first key, must have even number of remaining args
			if j+1 >= len(parts) {
				// Panic immediately - this is a programming error
				panic(fmt.Sprintf("doterr.NewErr(): trailing key %q at position %d without value. "+
					"Keys must be followed by values, or use an error as the final argument for a cause. "+
					"Fix: either add a value for %q, remove it, or if you meant to pass a cause error, ensure it's an error type not a string.",
					v, j, v))
			}
			j++ // Skip the value
		case error:
			// Errors after sentinels are not allowed
			// Build entry manually to avoid recursion
			e := newEntry([]error{ErrMisplacedError}, []kv{
				{k: "position", v: j},
				{k: "message", v: "errors must be first"},
			})
			return e
		default:
			// Non-string, non-error, non-ErrKV values are not allowed - panic immediately
			panic(fmt.Sprintf("doterr.NewErr: invalid argument type %T at position %d. "+
				"Only error, ErrKV, []ErrKV, func()ErrKV, or string keys are allowed. "+
				"Fix: remove the invalid argument or convert it to a valid type.",
				v, j))
		}
	}

	// If we found a first key, check that remaining count is even
	if firstKeyIdx >= 0 {
		remaining := len(parts) - firstKeyIdx
		if remaining%2 != 0 {
			// Build entry manually to avoid recursion
			e := newEntry([]error{ErrOddKeyValueCount}, []kv{
				{k: "position", v: firstKeyIdx},
				{k: "remaining_count", v: remaining},
			})
			return e
		}
	}

	return nil
}

// buildEntry creates an entry from parts without validation.
// Used internally by WithErr where sentinels are optional.
// Returns nil if parts are empty or result in an empty entry.
func buildEntry(parts ...any) error {
	if len(parts) == 0 {
		return nil
	}
	var e entry
	e.id = uniqueId
	appendEntry(&e, parts...)
	if e.empty() {
		return nil
	}
	return e
}

// handleCause inspects err and cause and, if cause is non-nil,
// appends cause to entry.errors (or creates new entry wrapping both).
func handleCause(err, cause error) error {
	if cause == nil {
		return err
	}
	if err == nil {
		return cause
	}
	// Append cause to entry.errors (cause LAST)
	return addToEntryErrors(err, cause)
}

// checkCrossPackage wraps an error with ErrCrossPackageError if it's an entry
// from a different doterr package (different uniqueId).
func checkCrossPackage(err error) error {
	//goland:noinspection GoTypeAssertionOnErrors
	e, isEntry := err.(entry)
	if isEntry && e.id != uniqueId {
		// Cross-package error detected - create entry with sentinel and cross-package error
		crossPkgErr := buildEntry(ErrCrossPackageError, "package_id", e.id, "expected_id", uniqueId)
		return addToEntryErrors(crossPkgErr, err)
	}
	return err
}

// containsSentinels checks if parts contains any error types (sentinels).
// Returns true if any errors are found, false if only metadata (KV, string pairs, etc.).
func containsSentinels(parts []any) bool {
	for _, part := range parts {
		if _, isError := part.(error); isError {
			return true
		}
	}
	return false
}

// buildErr tries to enrich the rightmost doterr entry inside baseErr.
// If none found, creates a fresh entry with baseErr in its errors slice,
// preserving baseErr's internals (including any existing cause).
//
// Key behavior: If middle contains any sentinels (errors), we create a NEW entry
// and store baseErr in entry.errors to preserve context layering. Metadata-only additions can enrich.
func buildErr(baseErr error, middle []any) error {
	// Check if middle contains any sentinel errors
	hasSentinels := containsSentinels(middle)

	// If we have new sentinels, create a new entry (don't enrich)
	// This preserves separate context levels with their own metadata
	if hasSentinels {
		e := buildEntry(middle...)
		if e != nil {
			// New context entry stores base error in its errors slice
			return addToEntryErrors(e, baseErr)
		}
		return baseErr
	}

	// No sentinels - just metadata enrichment, so try to enrich existing entry
	enriched, ok := enrichRightmost(baseErr, middle...)
	if ok {
		// Successfully merged metadata into an existing doterr entry
		return enriched
	}

	// No doterr entry found to enrich; create a fresh entry with baseErr in errors
	e := buildEntry(middle...)
	if e != nil {
		return addToEntryErrors(e, baseErr)
	}
	return baseErr
}

// enrichRightmost tries to merge parts into a doterr entry that is either:
//
//	(a) the error itself, or
//	(b) one of the immediate children of a multi-unwrap (errors.Join) tree.
//
// It does NOT recurse deeper than one join level.
//
//goland:noinspection DuplicatedCode
func enrichRightmost(err error, parts ...any) (error, bool) {
	// Case (a): err is an entry → merge directly.
	//goland:noinspection GoTypeAssertionOnErrors
	e, ok := err.(entry)
	if ok {
		tmp := e
		appendEntry(&tmp, parts...)
		return tmp, true
	}

	// Case (b): err is a join (multi-unwrap) → scan immediate children right-to-left.
	type unwrapper interface{ Unwrap() []error }
	u, ok := err.(unwrapper)
	if !ok {
		return err, false
	}
	children := u.Unwrap()
	if len(children) == 0 {
		return err, false
	}
	newKids := make([]error, len(children))
	copy(newKids, children)

	for i := len(newKids) - 1; i >= 0; i-- {
		//goland:noinspection GoTypeAssertionOnErrors
		e, ok := newKids[i].(entry)
		if ok {
			tmp := e
			appendEntry(&tmp, parts...)
			newKids[i] = tmp
			return errors.Join(newKids...), true
		}
		// NOTE: Deliberately NO recursion into nested joins.
	}

	// No doterr entry found at this level.
	return err, false
}
