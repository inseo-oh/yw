package selector

import (
	"fmt"
	"log"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

// ComplexSelector represents a CSS complex selector (e.g. .foo > #bar)
//
// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector
type ComplexSelector struct {
	Base CompoundSelector      // Very first selector in the complex selector
	Rest []ComplexSelectorRest // Rest of selectors
}

// Value for [ComplexSelector]'s Rest field.
type ComplexSelectorRest struct {
	Combinator Combinator       // Relationship between Selector and previous ComplexSelectorRest(or ComplexSelector's Base, if not present)
	Selector   CompoundSelector // A selector
}

// Relationship between two selectors in [ComplexSelectorRest]
type Combinator uint8

const (
	ChildCombinator       Combinator = iota // A B (B is child of A)
	DirectChildCombinator                   // A > B (B is direct child of A)
	PlusCombinator                          // A + B
	TildeCombinator                         // A ~ B
	TwoBarsCombinator                       // A || B
)

func (sel ComplexSelector) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v", sel.Base))
	for _, v := range sel.Rest {
		switch v.Combinator {
		case ChildCombinator:
			sb.WriteString(" ")
		case DirectChildCombinator:
			sb.WriteString(">")
		case PlusCombinator:
			sb.WriteString("+")
		case TildeCombinator:
			sb.WriteString("~")
		case TwoBarsCombinator:
			sb.WriteString("||")
		}
		sb.WriteString(fmt.Sprintf("%v", v.Selector))
	}
	return sb.String()
}
func (sel ComplexSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(ComplexSelector); !ok {
		return false
	} else {
		if !sel.Base.Equals(otherSel.Base) {
			return false
		}
		for i := 0; i < len(sel.Rest); i++ {
			if sel.Rest[i].Combinator != otherSel.Rest[i].Combinator {
				return false
			}
			if !sel.Rest[i].Selector.Equals(otherSel.Rest[i].Selector) {
				return false
			}
		}
	}
	return true
}

// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-complex-selector-against-an-element
func (s ComplexSelector) MatchAgainst(element dom.Element) bool {
	// Test each compound selector, from right to left
	for i := len(s.Rest) - 1; 0 < i; i-- {
		prevSel := s.Base
		if i != 0 {
			prevSel = s.Rest[i-1].Selector
		}

		sel := s.Rest[i].Selector
		if sel.MatchAgainst(element) {
			return false
		}
		switch s.Rest[i].Combinator {
		case ChildCombinator:
			// A B
			currElem := element.Parent()
			found := false
			for !util.IsNil(currElem) {
				if _, ok := (currElem).(dom.Element); !ok {
					break
				}
				if prevSel.MatchAgainst(currElem.(dom.Element)) {
					found = true
					break
				}
				currElem = currElem.Parent()
			}
			if !found {
				return false
			}
		case DirectChildCombinator:
			// A > B
			if util.IsNil(element.Parent()) {
				return false
			}
			if parent, ok := element.Parent().(dom.Element); !ok {
				if !prevSel.MatchAgainst(parent) {
					return false
				}
			} else {
				return false
			}
		case PlusCombinator, TildeCombinator, TwoBarsCombinator:
			panic("TODO")
		default:
			log.Printf("BUG: bad Combinator %d while parsing selector: %v", s.Rest[i].Combinator, s)
			continue
		}
	}
	return s.Base.MatchAgainst(element)
}
