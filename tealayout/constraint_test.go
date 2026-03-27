package tealayout

import "testing"

func TestFixedConstraint(t *testing.T) {
	c := fixedConstraint(30)
	if c.kind != constraintFixed {
		t.Errorf("kind = %v, want constraintFixed", c.kind)
	}
	if c.fixedSize != 30 {
		t.Errorf("fixedSize = %d, want 30", c.fixedSize)
	}
	if c.maxSize != -1 {
		t.Errorf("maxSize = %d, want -1 (unbounded)", c.maxSize)
	}
	if c.minSize != 0 {
		t.Errorf("minSize = %d, want 0", c.minSize)
	}
	if c.optional {
		t.Error("optional = true, want false")
	}
}

func TestFlexConstraint(t *testing.T) {
	c := flexConstraint(1.618)
	if c.kind != constraintFlex {
		t.Errorf("kind = %v, want constraintFlex", c.kind)
	}
	if c.flexWeight != 1.618 {
		t.Errorf("flexWeight = %f, want 1.618", c.flexWeight)
	}
	if c.maxSize != -1 {
		t.Errorf("maxSize = %d, want -1", c.maxSize)
	}
}

func TestFitConstraint(t *testing.T) {
	c := fitConstraint()
	if c.kind != constraintFit {
		t.Errorf("kind = %v, want constraintFit", c.kind)
	}
	if c.maxSize != -1 {
		t.Errorf("maxSize = %d, want -1", c.maxSize)
	}
}

func TestConstraint_FieldModification(t *testing.T) {
	c := flexConstraint(1.0)
	c.minSize = 20
	if c.minSize != 20 {
		t.Errorf("minSize = %d, want 20", c.minSize)
	}

	c.maxSize = 60
	if c.maxSize != 60 {
		t.Errorf("maxSize = %d, want 60", c.maxSize)
	}

	c.optional = true
	if !c.optional {
		t.Error("optional = false, want true")
	}
}
