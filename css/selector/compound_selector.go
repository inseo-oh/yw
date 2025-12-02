package selector

import (
	"fmt"
	"strings"

	"github.com/inseo-oh/yw/dom"
)

// CompoundSelector represents a [CSS compound selector] (e.g. div#foo.bar)
//
// [CSS compound selector]: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#compound
type CompoundSelector struct {
	TypeSelector     Selector // May be nil
	SubclassSelector []Selector
	PseudoItems      []CompundSelectorPseudoItem
}

// CompundSelectorPseudoItem is entry for [CompoundSelector]'s PseudoItems field.
type CompundSelectorPseudoItem struct {
	ElementSelector PseudoClassSelector
	ClassSelector   []PseudoClassSelector
}

func (sel CompoundSelector) String() string {
	sb := strings.Builder{}
	if sel.TypeSelector != nil {
		sb.WriteString(fmt.Sprintf("%v", sel.TypeSelector))
	}
	for _, v := range sel.SubclassSelector {
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	for _, p := range sel.PseudoItems {
		sb.WriteString(fmt.Sprintf("%v", p.ElementSelector))
		for _, v := range p.ClassSelector {
			sb.WriteString(fmt.Sprintf("%v", v))
		}
	}
	return sb.String()
}
func (sel CompoundSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(CompoundSelector); !ok {
		return false
	} else {
		if (sel.TypeSelector != nil) != (otherSel.TypeSelector != nil) {
			return false
		} else if (sel.TypeSelector != nil) && !sel.TypeSelector.Equals(otherSel.TypeSelector) {
			return false
		}
		if len(sel.SubclassSelector) != len(otherSel.SubclassSelector) {
			return false
		}
		for i := 0; i < len(sel.SubclassSelector); i++ {
			if !sel.SubclassSelector[i].Equals(otherSel.SubclassSelector[i]) {
				return false
			}
		}
		if len(sel.PseudoItems) != len(otherSel.PseudoItems) {
			return false
		}
		for i := 0; i < len(sel.PseudoItems); i++ {
			if !sel.PseudoItems[i].ElementSelector.Equals(otherSel.PseudoItems[i].ElementSelector) {
				return false
			}
			if len(sel.PseudoItems[i].ClassSelector) != len(otherSel.PseudoItems[i].ClassSelector) {
				return false
			}
			for j := 0; j < len(sel.PseudoItems[i].ClassSelector); j++ {
				if !sel.PseudoItems[i].ClassSelector[j].Equals(otherSel.PseudoItems[i].ClassSelector[j]) {
					return false
				}
			}
		}
	}
	return true
}
func (sel CompoundSelector) MatchAgainst(element dom.Element) bool {
	if sel.TypeSelector != nil && !sel.TypeSelector.MatchAgainst(element) {
		return false
	}
	for _, ss := range sel.SubclassSelector {
		if !ss.MatchAgainst(element) {
			return false
		}
	}
	if len(sel.PseudoItems) != 0 {
		for i := len(sel.PseudoItems) - 1; 0 < i; i-- {
			// TODO
		}
	}
	return true
}
