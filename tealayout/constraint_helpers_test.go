package tealayout

// Test helpers for creating constraints in tests.

func fixedConstraint(n int) constraint {
	return constraint{
		kind:      constraintFixed,
		fixedSize: n,
		maxSize:   -1,
	}
}

func flexConstraint(weight float64) constraint {
	return constraint{
		kind:       constraintFlex,
		flexWeight: weight,
		maxSize:    -1,
	}
}

func fitConstraint() constraint {
	return constraint{
		kind:    constraintFit,
		maxSize: -1,
	}
}
