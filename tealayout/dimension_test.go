package tealayout

import "testing"

func TestPercent(t *testing.T) {
	d := Percent(50)
	if d.kind != dimensionPercent {
		t.Errorf("kind = %v, want dimensionPercent", d.kind)
	}
	if d.value != 50 {
		t.Errorf("value = %f, want 50", d.value)
	}
}

func TestFixed_Dimension(t *testing.T) {
	d := Fixed(30)
	if d.kind != dimensionFixed {
		t.Errorf("kind = %v, want dimensionFixed", d.kind)
	}
	if d.value != 30 {
		t.Errorf("value = %f, want 30", d.value)
	}
}

func TestFlex_Dimension(t *testing.T) {
	d := Flex(1.618)
	if d.kind != dimensionFlex {
		t.Errorf("kind = %v, want dimensionFlex", d.kind)
	}
	if d.value != 1.618 {
		t.Errorf("value = %f, want 1.618", d.value)
	}
}

func TestFit_Dimension(t *testing.T) {
	d := Fit()
	if d.kind != dimensionFit {
		t.Errorf("kind = %v, want dimensionFit", d.kind)
	}
}

func TestPercentConstants(t *testing.T) {
	tests := []struct {
		name string
		dim  Dimension
		want float64
	}{
		{"Percent100", Percent100, 100},
		{"Percent75", Percent75, 75},
		{"Percent50", Percent50, 50},
		{"Percent33", Percent33, 33},
		{"Percent25", Percent25, 25},
		{"Percent20", Percent20, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dim.kind != dimensionPercent {
				t.Errorf("kind = %v, want dimensionPercent", tt.dim.kind)
			}
			if tt.dim.value != tt.want {
				t.Errorf("value = %f, want %f", tt.dim.value, tt.want)
			}
		})
	}
}

func TestPercent_MapsToFlexConstraint(t *testing.T) {
	p := NewRow(Percent100)
	cs := p.toConstraint()
	if cs.kind != constraintFlex {
		t.Errorf("constraint kind = %v, want constraintFlex", cs.kind)
	}
	if cs.flexWeight != 100 {
		t.Errorf("flexWeight = %f, want 100", cs.flexWeight)
	}
}

func TestFixed_MapsToFixedConstraint(t *testing.T) {
	p := NewRow(Fixed(30))
	cs := p.toConstraint()
	if cs.kind != constraintFixed {
		t.Errorf("constraint kind = %v, want constraintFixed", cs.kind)
	}
	if cs.fixedSize != 30 {
		t.Errorf("fixedSize = %d, want 30", cs.fixedSize)
	}
}

func TestFlex_MapsToFlexConstraint(t *testing.T) {
	p := NewRow(Flex(1.618))
	cs := p.toConstraint()
	if cs.kind != constraintFlex {
		t.Errorf("constraint kind = %v, want constraintFlex", cs.kind)
	}
	if cs.flexWeight != 1.618 {
		t.Errorf("flexWeight = %f, want 1.618", cs.flexWeight)
	}
}

func TestFit_MapsToFitConstraint(t *testing.T) {
	p := NewRow(Fit())
	cs := p.toConstraint()
	if cs.kind != constraintFit {
		t.Errorf("constraint kind = %v, want constraintFit", cs.kind)
	}
}
