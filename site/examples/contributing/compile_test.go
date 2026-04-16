// Source: site/src/content/docs/contributing/index.mdx:84#cc88d1b3
package examples_test

import (
	"fmt"
	"testing"
)

func fetchData() (string, error) { return "data", nil }
func process(s string) string    { return s }

func doWork() (result string, err error) {
	var data string
	data, err = fetchData()
	if err != nil {
		err = fmt.Errorf("fetch failed: %w", err)
		goto end
	}
	result = process(data)
end:
	return result, err
}

// TestCompile_ClearPathContributing verifies the ClearPath pattern example from contributing/index.mdx line 84.
func TestCompile_ClearPathContributing(t *testing.T) {
	result, err := doWork()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}
