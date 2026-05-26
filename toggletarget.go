package utils

// TargetValue represents a toggleable target state.
type TargetValue int

// ToggleTarget manages a current target that can be toggled
// between two predefined states.
type ToggleTarget struct {
	Active    TargetValue
	Primary   TargetValue // Primary is the default / initial state.
	Secondary TargetValue // Secondary is the alternate state.
}

// NewToggleTarget returns a ToggleTarget initialized with
// a primary state and its alternate state.
func NewToggleTarget(start, alternate TargetValue) *ToggleTarget {
	return &ToggleTarget{
		Active:    start,
		Primary:   start,
		Secondary: alternate,
	}
}

type IfChain struct {
	matched bool
}

func (t *ToggleTarget) IfActive(want TargetValue, fn func()) *IfChain {
	if t.Active == want {
		fn()
		return &IfChain{matched: true}
	}
	return &IfChain{matched: false}
}

func (c *IfChain) ElseIf(want TargetValue, t *ToggleTarget, fn func()) *IfChain {
	if !c.matched && t.Active == want {
		fn()
		c.matched = true
	}
	return c
}

func (c *IfChain) Else(fn func()) {
	if !c.matched {
		fn()
	}
}

// SetCurrent sets the current active state explicitly.
// If val is not equal to either Primary or Secondary,
// the call is ignored.
func (t *ToggleTarget) SetCurrent(val TargetValue) {
	if val != t.Primary && val != t.Secondary {
		return // or panic, depending on design choice
	}
	t.Active = val
}

// Toggle switches the active state between Primary and Secondary.
// If the current state is invalid, it falls back to Primary.
func (t *ToggleTarget) Toggle() {
	switch t.Active {
	case t.Primary:
		t.Active = t.Secondary
	case t.Secondary:
		t.Active = t.Primary
	default:
		t.Active = t.Primary // or panic / ignore
	}
}
