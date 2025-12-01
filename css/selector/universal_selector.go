package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// UniversalSelector represents CSS Universal selector, with optional namespace prefix (e.g. *)
//
// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#the-universal-selector
type UniversalSelector struct {
	NsPrefix *NsPrefix
}

func (sel UniversalSelector) String() string { return fmt.Sprintf("%v*", sel.NsPrefix) }
func (sel UniversalSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(UniversalSelector); !ok {
		return false
	} else {
		if (sel.NsPrefix == nil) != (otherSel.NsPrefix == nil) {
			return false
		}
		if sel.NsPrefix != nil && !sel.NsPrefix.Equals(*otherSel.NsPrefix) {
			return false
		}
	}
	return true
}
func (sel UniversalSelector) MatchAgainst(element dom.Element) bool {
	return true
}
