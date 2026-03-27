package tealayout

import "math"

// resolveLinear resolves a single axis of layout: given available space,
// constraints, gap between children, and optional SizeHinters (may be nil
// or contain nil entries), it returns the assigned size for each child.
func resolveLinear(available int, constraints []constraint, gap int, hinters []SizeHinter) ([]int, error) {
	n := len(constraints)
	if n == 0 {
		return nil, nil
	}

	// Normalize hinters slice
	if len(hinters) < n {
		padded := make([]SizeHinter, n)
		copy(padded, hinters)
		hinters = padded
	}

	sizes := make([]int, n)
	active := make([]bool, n)
	for i := range active {
		active[i] = true
	}

	sizes, _ = resolveWithOptionalRemoval(available, constraints, gap, hinters, sizes, active)
	return sizes, nil
}

// resolveWithOptionalRemoval runs the full 5-phase algorithm including
// optional child removal and retry.
func resolveWithOptionalRemoval(available int, constraints []constraint, gap int, hinters []SizeHinter, sizes []int, active []bool) ([]int, []bool) {
	n := len(constraints)

	for {
		// Count active children for gap calculation
		activeCount := 0
		for i := range n {
			if active[i] {
				activeCount++
			}
		}

		gapTotal := 0
		if activeCount > 1 {
			gapTotal = gap * (activeCount - 1)
		}

		remaining := available - gapTotal
		if remaining < 0 {
			remaining = 0
		}

		// Phase 1: Resolve Fixed and Fit children
		remaining = resolveFixedAndFit(remaining, constraints, hinters, sizes, active)

		// Phase 2+3: Distribute remaining to Flex with clamping loop
		distributeFlexWithClamping(remaining, constraints, sizes, active)

		// Phase 4: Check optional children — remove any below MinSize
		removed := false
		for i := range n {
			if !active[i] {
				continue
			}
			c := constraints[i]
			if c.optional && c.minSize > 0 && sizes[i] < c.minSize {
				active[i] = false
				sizes[i] = 0
				removed = true
			}
		}

		if !removed {
			break
		}
		// Reset non-removed sizes and retry from Phase 1
		for i := range n {
			if active[i] {
				sizes[i] = 0
			}
		}
	}

	return sizes, active
}

// resolveFixedAndFit handles Phase 1: assign sizes to Fixed and Fit children,
// subtract from remaining. Returns updated remaining space.
func resolveFixedAndFit(remaining int, constraints []constraint, hinters []SizeHinter, sizes []int, active []bool) int {
	for i, c := range constraints {
		if !active[i] {
			continue
		}
		switch c.kind {
		case constraintFixed:
			size := clampConstraint(c.fixedSize, c)
			sizes[i] = size
			remaining -= size

		case constraintFit:
			desired := 0
			if hinters[i] != nil {
				hint := hinters[i].SizeHint(remaining, 0)
				desired = hint.Desired.Width
			}
			size := clampConstraint(desired, c)
			sizes[i] = size
			remaining -= size
		}
	}
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// distributeFlexWithClamping handles Phases 2+3: distribute remaining space
// to Flex children using cumulative rounding, then clamp and redistribute
// until stable.
func distributeFlexWithClamping(remaining int, constraints []constraint, sizes []int, active []bool) {
	n := len(constraints)

	// Track which flex children are frozen (clamped at min/max)
	frozen := make([]bool, n)

	// Clamping loop — max N iterations (one child freezes per iteration worst case)
	for iter := range n {
		_ = iter

		// Collect unfrozen flex children
		totalWeight := 0.0
		flexAvailable := remaining
		for i, c := range constraints {
			if !active[i] || c.kind != constraintFlex {
				continue
			}
			if frozen[i] {
				flexAvailable -= sizes[i]
				continue
			}
			totalWeight += c.flexWeight
		}

		if totalWeight == 0 || flexAvailable <= 0 {
			break
		}

		// Cumulative rounding distribution
		cumulative := 0.0
		prevPos := 0
		anyClamped := false

		for i, c := range constraints {
			if !active[i] || c.kind != constraintFlex || frozen[i] {
				continue
			}
			fraction := c.flexWeight / totalWeight
			cumulative += fraction * float64(flexAvailable)
			pos := int(math.Round(cumulative))
			rawSize := pos - prevPos
			prevPos = pos

			// For optional children, skip minSize clamping during distribution.
			// Phase 4 will check the raw size and remove them if too small.
			clamped := clampConstraintForFlex(rawSize, c)
			sizes[i] = clamped
			if clamped != rawSize {
				frozen[i] = true
				anyClamped = true
			}
		}

		if !anyClamped {
			break
		}
	}
}

// clampConstraint applies min/max bounds from a constraint.
func clampConstraint(val int, c constraint) int {
	if val < c.minSize {
		val = c.minSize
	}
	if c.maxSize >= 0 && val > c.maxSize {
		val = c.maxSize
	}
	return val
}

// clampConstraintForFlex applies bounds during flex distribution. For optional
// children, minSize is not enforced — Phase 4 handles removal instead.
func clampConstraintForFlex(val int, c constraint) int {
	if !c.optional && val < c.minSize {
		val = c.minSize
	}
	if c.maxSize >= 0 && val > c.maxSize {
		val = c.maxSize
	}
	return val
}
