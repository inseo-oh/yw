// Package display provides types and values for CSS Selector Module Level 4,
// as well as Element selection logic.
//
// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/
package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

type Selector interface {
	String() string

	// Equals reports whether two selectors are equal.
	// Selectors with different types are not considered as equal.
	Equals(other Selector) bool

	// MatchAgainst reports whether the selector matches given element.
	MatchAgainst(element dom.Element) bool
}

// NsPrefix represents CSS namespace prefix (e.g. foo|)
type NsPrefix struct{ Ident string }

func (sel NsPrefix) Equals(other NsPrefix) bool {
	// STUB
	return true
}

// WqName represents an CSS identifier with optional namespace.
//
// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-wq-name
type WqName struct {
	NsPrefix *NsPrefix // May be nil
	Ident    string
}

func (wqName WqName) String() string {
	if wqName.NsPrefix != nil {
		return fmt.Sprintf("%v%s", wqName.NsPrefix, wqName.Ident)
	} else {
		return wqName.Ident
	}
}

// Equals reports whether two names are identical.
func (wqName WqName) Equals(other WqName) bool {
	if (wqName.NsPrefix != nil) != (other.NsPrefix != nil) {
		return false
	} else if (wqName.NsPrefix != nil) && !wqName.NsPrefix.Equals(*other.NsPrefix) {
		return false
	}
	if wqName.Ident != other.Ident {
		return false
	}
	return true
}

// MatchAgainstElement matches given selectors against an DOM element.
//
// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-selector-against-an-element
func MatchAgainstElement(selector []Selector, element dom.Element) bool {
	for _, s := range selector {
		if s.MatchAgainst(element) {
			return true
		}
	}
	return false
}

// MatchAgainstElement matches given selectors against given DOM trees.
//
// Spec: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-selector-against-a-tree
func MatchAgainstTree(selector []Selector, roots []dom.Node) []dom.Node {
	selectorMatchList := []dom.Node{}
	for _, root := range roots {
		candiateElems := []dom.Node{}
		for _, n := range dom.InclusiveDescendants(root) {
			if _, ok := n.(dom.Element); ok {
				candiateElems = append(candiateElems, n)
			}
		}
		for _, n := range candiateElems {
			if MatchAgainstElement(selector, n.(dom.Element)) {
				selectorMatchList = append(selectorMatchList, n)
			}
			// TODO: Pseudo element
		}

	}
	return selectorMatchList
}
