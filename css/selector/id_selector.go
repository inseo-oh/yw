package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-id-selector
type IdSelector struct{ Id string }

func (sel IdSelector) String() string { return fmt.Sprintf("#%s", sel.Id) }
func (sel IdSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(IdSelector); !ok {
		return false
	} else {
		if sel.Id != otherSel.Id {
			return false
		}
	}
	return true
}
func (sel IdSelector) MatchAgainst(element dom.Element) bool {
	id := sel.Id
	elemId, ok := element.AttrWithoutNamespace("id")
	if !ok || elemId != id {
		return false
	}
	return true
}
