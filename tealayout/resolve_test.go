package tealayout

import (
	"math/rand"
	"testing"
)

// mockHinter implements SizeHinter for testing.
type mockHinter struct {
	desiredWidth int
}

func (m mockHinter) SizeHint(availWidth, availHeight int) SizeHint {
	return SizeHint{
		Desired: Size{Width: m.desiredWidth},
		Max:     Size{Width: -1, Height: -1},
	}
}

func sum(s []int) int {
	total := 0
	for _, v := range s {
		total += v
	}
	return total
}

// --- Fixed-only tests ---

func TestResolveLinear_FixedOnly(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{fixedConstraint(20), fixedConstraint(20), fixedConstraint(20)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 20 || sizes[1] != 20 || sizes[2] != 20 {
		t.Errorf("got %v, want [20 20 20]", sizes)
	}
}

func TestResolveLinear_FixedOverflow(t *testing.T) {
	// Fixed children exceed available — no panic, just assigned their sizes
	sizes, err := resolveLinear(30, []constraint{fixedConstraint(20), fixedConstraint(20), fixedConstraint(20)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 20 || sizes[1] != 20 || sizes[2] != 20 {
		t.Errorf("got %v, want [20 20 20]", sizes)
	}
}

func TestResolveLinear_SingleFixed(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{fixedConstraint(30)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 30 {
		t.Errorf("got %v, want [30]", sizes)
	}
}

// --- Flex-only tests ---

func TestResolveLinear_FlexEqual(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{flexConstraint(1), flexConstraint(1), flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sum(sizes) != 80 {
		t.Errorf("sum = %d, want 80", sum(sizes))
	}
	// Each should be ~26-27
	for i, s := range sizes {
		if s < 26 || s > 27 {
			t.Errorf("sizes[%d] = %d, want 26-27", i, s)
		}
	}
}

func TestResolveLinear_FlexUnequal(t *testing.T) {
	sizes, err := resolveLinear(90, []constraint{flexConstraint(1), flexConstraint(2)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sum(sizes) != 90 {
		t.Errorf("sum = %d, want 90", sum(sizes))
	}
	if sizes[0] != 30 {
		t.Errorf("sizes[0] = %d, want 30", sizes[0])
	}
	if sizes[1] != 60 {
		t.Errorf("sizes[1] = %d, want 60", sizes[1])
	}
}

func TestResolveLinear_FlexGoldenRatio(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{flexConstraint(1.0), flexConstraint(1.618)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sum(sizes) != 80 {
		t.Errorf("sum = %d, want 80", sum(sizes))
	}
	// 1.0/(1.0+1.618) * 80 ≈ 30.6 → 31
	// 1.618/(1.0+1.618) * 80 ≈ 49.4 → 49
	if sizes[0] != 31 {
		t.Errorf("sizes[0] = %d, want 31", sizes[0])
	}
	if sizes[1] != 49 {
		t.Errorf("sizes[1] = %d, want 49", sizes[1])
	}
}

func TestResolveLinear_SingleFlex(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 80 {
		t.Errorf("sizes[0] = %d, want 80", sizes[0])
	}
}

// --- Mixed tests ---

func TestResolveLinear_FixedAndFlex(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{fixedConstraint(20), flexConstraint(1), flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 20 {
		t.Errorf("sizes[0] = %d, want 20", sizes[0])
	}
	if sizes[1] != 30 {
		t.Errorf("sizes[1] = %d, want 30", sizes[1])
	}
	if sizes[2] != 30 {
		t.Errorf("sizes[2] = %d, want 30", sizes[2])
	}
}

// --- Fit tests ---

func TestResolveLinear_Fit(t *testing.T) {
	hinters := []SizeHinter{mockHinter{desiredWidth: 25}}
	sizes, err := resolveLinear(80, []constraint{fitConstraint()}, 0, hinters)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 25 {
		t.Errorf("sizes[0] = %d, want 25", sizes[0])
	}
}

func TestResolveLinear_FitAndFlex(t *testing.T) {
	hinters := []SizeHinter{mockHinter{desiredWidth: 25}, nil}
	sizes, err := resolveLinear(80, []constraint{fitConstraint(), flexConstraint(1)}, 0, hinters)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 25 {
		t.Errorf("sizes[0] = %d, want 25", sizes[0])
	}
	if sizes[1] != 55 {
		t.Errorf("sizes[1] = %d, want 55", sizes[1])
	}
}

func TestResolveLinear_FitNoHinter(t *testing.T) {
	// Fit without a SizeHinter → size 0
	sizes, err := resolveLinear(80, []constraint{fitConstraint()}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 0 {
		t.Errorf("sizes[0] = %d, want 0", sizes[0])
	}
}

// --- Clamping tests ---

func TestResolveLinear_FlexMaxSize(t *testing.T) {
	c1 := flexConstraint(1)
	c1.maxSize = 20
	sizes, err := resolveLinear(80, []constraint{c1, flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 20 {
		t.Errorf("sizes[0] = %d, want 20 (capped)", sizes[0])
	}
	if sizes[1] != 60 {
		t.Errorf("sizes[1] = %d, want 60 (gets remainder)", sizes[1])
	}
}

func TestResolveLinear_FlexMinSize(t *testing.T) {
	// Two flex children in 40 space, one needs min 30
	c1 := flexConstraint(1)
	c1.minSize = 30
	sizes, err := resolveLinear(40, []constraint{c1, flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] < 30 {
		t.Errorf("sizes[0] = %d, want >= 30", sizes[0])
	}
}

func TestResolveLinear_MultipleClamped(t *testing.T) {
	c1 := flexConstraint(1)
	c1.maxSize = 15
	c2 := flexConstraint(1)
	c2.maxSize = 15
	sizes, err := resolveLinear(100, []constraint{c1, c2, flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 15 {
		t.Errorf("sizes[0] = %d, want 15", sizes[0])
	}
	if sizes[1] != 15 {
		t.Errorf("sizes[1] = %d, want 15", sizes[1])
	}
	if sizes[2] != 70 {
		t.Errorf("sizes[2] = %d, want 70", sizes[2])
	}
}

// --- Optional tests ---

func TestResolveLinear_OptionalRemoved(t *testing.T) {
	// 3 flex children in 30 space. Third is optional with min 20 — can't fit.
	c3 := flexConstraint(1)
	c3.minSize = 20
	c3.optional = true
	sizes, err := resolveLinear(30, []constraint{flexConstraint(1), flexConstraint(1), c3}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[2] != 0 {
		t.Errorf("sizes[2] = %d, want 0 (removed)", sizes[2])
	}
	if sizes[0]+sizes[1] != 30 {
		t.Errorf("remaining children should total 30, got %d", sizes[0]+sizes[1])
	}
}

func TestResolveLinear_OptionalKept(t *testing.T) {
	// Enough space for the optional child
	c3 := flexConstraint(1)
	c3.minSize = 20
	c3.optional = true
	sizes, err := resolveLinear(90, []constraint{flexConstraint(1), flexConstraint(1), c3}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[2] < 20 {
		t.Errorf("sizes[2] = %d, want >= 20 (should be kept)", sizes[2])
	}
}

// --- Gap tests ---

func TestResolveLinear_Gap(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{flexConstraint(1), flexConstraint(1), flexConstraint(1)}, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	// 80 - 2*2(gaps) = 76 available for children
	if sum(sizes) != 76 {
		t.Errorf("sum = %d, want 76 (80 - 4 gap)", sum(sizes))
	}
}

func TestResolveLinear_GapSingleChild(t *testing.T) {
	sizes, err := resolveLinear(80, []constraint{flexConstraint(1)}, 5, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Single child — no gap applied
	if sizes[0] != 80 {
		t.Errorf("sizes[0] = %d, want 80 (no gap for single child)", sizes[0])
	}
}

func TestResolveLinear_GapWithOptionalRemoval(t *testing.T) {
	// Third child is optional and will be removed → gap reclaimed
	c3 := flexConstraint(1)
	c3.minSize = 20
	c3.optional = true
	sizes, err := resolveLinear(40, []constraint{flexConstraint(1), flexConstraint(1), c3}, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[2] != 0 {
		t.Errorf("sizes[2] = %d, want 0 (removed)", sizes[2])
	}
	// 2 remaining children: 40 - 2(1 gap) = 38
	if sum(sizes) != 38 {
		t.Errorf("sum = %d, want 38", sum(sizes))
	}
}

// --- Edge cases ---

func TestResolveLinear_ZeroAvailable(t *testing.T) {
	sizes, err := resolveLinear(0, []constraint{flexConstraint(1), flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sum(sizes) != 0 {
		t.Errorf("sum = %d, want 0", sum(sizes))
	}
}

func TestResolveLinear_NegativeAvailable(t *testing.T) {
	sizes, err := resolveLinear(-10, []constraint{flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 0 {
		t.Errorf("sizes[0] = %d, want 0", sizes[0])
	}
}

func TestResolveLinear_Empty(t *testing.T) {
	sizes, err := resolveLinear(80, nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(sizes) != 0 {
		t.Errorf("len(sizes) = %d, want 0", len(sizes))
	}
}

// --- Invariant: sum(sizes) + totalGap == available (for random inputs) ---

func TestResolveLinear_SumInvariant(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	for trial := range 100 {
		n := rng.Intn(5) + 1
		available := rng.Intn(200) + 20
		gap := rng.Intn(4)

		constraints := make([]constraint, n)
		for i := range n {
			switch rng.Intn(3) {
			case 0:
				constraints[i] = fixedConstraint(rng.Intn(30) + 5)
			case 1:
				constraints[i] = flexConstraint(float64(rng.Intn(3) + 1))
			case 2:
				c := flexConstraint(float64(rng.Intn(3) + 1))
				c.minSize = 5
				c.maxSize = rng.Intn(50) + 20
				constraints[i] = c
			}
		}

		sizes, err := resolveLinear(available, constraints, gap, nil)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}

		// No negative sizes
		for i, s := range sizes {
			if s < 0 {
				t.Errorf("trial %d: sizes[%d] = %d (negative)", trial, i, s)
			}
		}

		// Min/max respected
		for i, s := range sizes {
			c := constraints[i]
			if s > 0 && s < c.minSize && !c.optional {
				t.Errorf("trial %d: sizes[%d] = %d < minSize %d", trial, i, s, c.minSize)
			}
			if c.maxSize >= 0 && s > c.maxSize {
				t.Errorf("trial %d: sizes[%d] = %d > maxSize %d", trial, i, s, c.maxSize)
			}
		}
	}
}

// --- No negatives ---

func TestResolveLinear_NoNegativeSizes(t *testing.T) {
	sizes, err := resolveLinear(10, []constraint{fixedConstraint(20), flexConstraint(1)}, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range sizes {
		if s < 0 {
			t.Errorf("sizes[%d] = %d (negative)", i, s)
		}
	}
}
