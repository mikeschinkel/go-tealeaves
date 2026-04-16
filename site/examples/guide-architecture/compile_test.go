// Source: site/src/content/docs/guides/architecture.mdx:67#03624845,88#9ae7a93f,100#9196c8c4
package examples_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// archModel is a local stand-in that demonstrates the ClearPath goto pattern
// from architecture.mdx line 67.
type archModel struct {
	value string
}

func (m archModel) withValue(v string) archModel {
	m.value = v
	return m
}

func computeValue() (string, error) {
	return "computed", nil
}

// TestCompile_ClearPathPattern verifies the ClearPath goto pattern from architecture.mdx:67.
func TestCompile_ClearPathPattern(t *testing.T) {
	doSomething := func(m archModel) (result archModel, err error) {
		var value string
		value, err = computeValue()
		if err != nil {
			goto end
		}
		result = m.withValue(value)
	end:
		return result, err
	}

	m := archModel{}
	result, err := doSomething(m)
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

// optionError is a structured error type that demonstrates the doterr metadata pattern.
// In the real codebase, doterr provides this generically; here we show the equivalent.
type optionError struct {
	option string
	cause  error
}

func (e *optionError) Error() string {
	return fmt.Sprintf("invalid option %q: %v", e.option, e.cause)
}

func (e *optionError) Unwrap() error { return e.cause }

// TestCompile_StructuredErrorPattern verifies structured error wrapping from architecture.mdx:88.
// doterr is an internal pattern; this demonstrates the equivalent using a custom error type + errors.As.
func TestCompile_StructuredErrorPattern(t *testing.T) {
	newOptionError := func(name string, cause error) error {
		return &optionError{option: name, cause: cause}
	}

	extractOption := func(err error) (string, bool) {
		var oe *optionError
		if errors.As(err, &oe) {
			return oe.option, true
		}
		return "", false
	}

	cause := errors.New("out of range")
	err := newOptionError("maxItems", cause)
	if err == nil {
		t.Fatal("expected error")
	}

	opt, ok := extractOption(err)
	if !ok {
		t.Fatal("expected to extract option from error")
	}
	_ = opt
}

// TestCompile_ImmutableStructPattern verifies immutable struct updates from architecture.mdx:100.
func TestCompile_ImmutableStructPattern(t *testing.T) {
	modal := teamodal.NewYesNoModal("Are you sure?", &teamodal.ConfirmModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	modal = modal.SetSize(80, 24)

	dropdown := teafields.NewDropdownModel(
		teafields.ToOptions([]string{"Option A", "Option B"}),
		&teafields.DropdownModelArgs{},
	)
	dropdown = dropdown.WithScreenSize(80, 24)

	_ = modal
	_ = dropdown
}
