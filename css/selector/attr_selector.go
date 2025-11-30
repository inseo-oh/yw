package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attribute-selector
type AttrSelector struct {
	AttrName WqName
	Matcher  Matcher
	// Below are valid only if the matcher isn't 'none'
	AttrValue       string
	IsCaseSensitive bool
}

type Matcher uint8

const (
	NoMatcher       = Matcher(iota)
	NormalMatcher   // =
	TildeMatcher    // ~=
	BarMatcher      // |=
	CaretMatcher    // ^=
	DollarMatcher   // $=
	AsteriskMatcher // *=
)

func (sel AttrSelector) String() string {
	flagStr := "s"
	if !sel.IsCaseSensitive {
		flagStr = "i"
	}
	switch sel.Matcher {
	case NoMatcher:
		return fmt.Sprintf("[%s]", sel.AttrName)
	case NormalMatcher:
		return fmt.Sprintf("[%s=%s %s]", sel.AttrName, sel.AttrValue, flagStr)
	case TildeMatcher:
		return fmt.Sprintf("[%s~=%s %s]", sel.AttrName, sel.AttrValue, flagStr)
	case BarMatcher:
		return fmt.Sprintf("[%s|=%s %s]", sel.AttrName, sel.AttrValue, flagStr)
	case CaretMatcher:
		return fmt.Sprintf("[%s^=%s %s]", sel.AttrName, sel.AttrValue, flagStr)
	case DollarMatcher:
		return fmt.Sprintf("[%s$=%s %s]", sel.AttrName, sel.AttrValue, flagStr)
	case AsteriskMatcher:
		return fmt.Sprintf("[%s*=%s %s]", sel.AttrName, sel.AttrValue, flagStr)
	}
	return fmt.Sprintf("[%s<bad matcher %d>%s %s]", sel.AttrName, sel.Matcher, sel.AttrValue, flagStr)
}
func (sel AttrSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(AttrSelector); !ok {
		return false
	} else {
		if !sel.AttrName.Equals(otherSel.AttrName) {
			return false
		}
		if sel.Matcher != otherSel.Matcher {
			return false
		}
		if sel.Matcher != NoMatcher {
			if sel.AttrValue != otherSel.AttrValue {
				return false
			}
			if sel.IsCaseSensitive != otherSel.IsCaseSensitive {
				return false
			}
		}
	}
	return true
}
func (sel AttrSelector) MatchAgainst(element dom.Element) bool {
	// STUB
	return false
}
