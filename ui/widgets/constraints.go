package widgets

import "github.com/cloudboy-jh/bentotui/core"

// HeightConstraintType defines how a widget wants its height allocated.
type HeightConstraintType int

const (
	HeightFixed HeightConstraintType = iota
	HeightMin
	HeightMax
	HeightFlex
)

// HeightConstraint expresses how much vertical space a widget wants.
type HeightConstraint struct {
	Type  HeightConstraintType
	Value int
}

// FixedHeight returns a constraint for exactly n lines.
func FixedHeight(n int) HeightConstraint {
	return HeightConstraint{Type: HeightFixed, Value: n}
}

// MinHeight returns a constraint for at least n lines.
func MinHeight(n int) HeightConstraint {
	return HeightConstraint{Type: HeightMin, Value: n}
}

// MaxHeight returns a constraint for at most n lines.
func MaxHeight(n int) HeightConstraint {
	return HeightConstraint{Type: HeightMax, Value: n}
}

// FlexHeight returns a constraint to fill remaining space with given weight.
func FlexHeight(weight int) HeightConstraint {
	if weight < 1 {
		weight = 1
	}
	return HeightConstraint{Type: HeightFlex, Value: weight}
}

// HeightConstraintAware widgets can declare their preferred height.
type HeightConstraintAware interface {
	core.Component
	HeightConstraint(width int) HeightConstraint
}
