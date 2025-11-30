package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#ref-for-typedef-type-selector
type TypeSelector struct {
	TypeName WqName
}

func (sel TypeSelector) String() string { return fmt.Sprintf("%v", sel.TypeName) }
func (sel TypeSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(TypeSelector); !ok {
		return false
	} else {
		if !sel.TypeName.Equals(otherSel.TypeName) {
			return false
		}
	}
	return true
}
func (sel TypeSelector) MatchAgainst(element dom.Element) bool {
	// TODO: Handle namespace
	name := sel.TypeName.Ident
	if element.LocalName() != name {
		return false
	}
	return true
}
