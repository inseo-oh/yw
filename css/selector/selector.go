// Implementation of the CSS Selector Module Level 4 (https://www.w3.org/TR/2022/WD-selectors-4-20221111/)
package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

type Selector interface {
	String() string
	Equals(other Selector) bool
	MatchAgainst(element dom.Element) bool
}

type NsPrefix struct{ Ident string }

func (sel NsPrefix) Equals(other NsPrefix) bool {
	// STUB
	return true
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-wq-name
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

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-selector-against-an-element
func MatchAgainstElement(selector []Selector, element dom.Element) bool {
	for _, s := range selector {
		if s.MatchAgainst(element) {
			return true
		}
	}
	return false
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-selector-against-a-tree
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
