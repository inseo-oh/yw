package selector

import (
	"fmt"
	"yw/dom"
)

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-class-selector
type PseudoClassSelector struct {
	Name string
	Args []any
}

func (sel PseudoClassSelector) String() string {
	if len(sel.Args) != 0 {
		// TODO: Display arguments in better way
		return fmt.Sprintf(":%s(%v)", sel.Name, sel.Args)
	} else {
		return fmt.Sprintf(":%s", sel.Name)
	}
}
func (sel PseudoClassSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(PseudoClassSelector); !ok {
		return false
	} else {
		if sel.Name != otherSel.Name {
			return false
		}
		if len(sel.Args) != len(otherSel.Args) {
			return false
		}
		// TODO: Compare actual arguments
	}
	return true
}
func (sel PseudoClassSelector) MatchAgainst(element dom.Element) bool {
	// STUB
	return false
}
