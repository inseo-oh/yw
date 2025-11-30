package selector

import "yw/dom"

// Special CSS selector that matches DOM node pointer directly.
type NodePtrSelector struct {
	Element dom.Element
}

func (sel NodePtrSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(NodePtrSelector); !ok {
		return false
	} else {
		return otherSel.Element == sel.Element
	}
}
func (sel NodePtrSelector) MatchAgainst(element dom.Element) bool {
	return sel.Element == element
}
func (sel NodePtrSelector) String() string {
	return sel.Element.String()
}
