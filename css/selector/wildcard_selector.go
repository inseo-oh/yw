package selector

import (
	"fmt"
	"yw/dom"
)

type WildcardSelector struct {
	NsPrefix *NsPrefix
}

func (sel WildcardSelector) String() string { return fmt.Sprintf("%v*", sel.NsPrefix) }
func (sel WildcardSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(WildcardSelector); !ok {
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
func (sel WildcardSelector) MatchAgainst(element dom.Element) bool {
	return true
}
