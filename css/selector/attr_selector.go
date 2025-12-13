// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// AttrSelector represents a [CSS attribute selector] (e.g. [attr=value])
//
// [CSS attribute selector]: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#attribute-selector
type AttrSelector struct {
	AttrName WqName  // Name of attribute
	Matcher  Matcher // How attribute's value should be matched

	// Below are valid only if the matcher isn't 'none'

	AttrValue       string // Value of attribute
	IsCaseSensitive bool   // Is case-sensitive?
}

// Matcher represents how attribute's value should be matched in [AttrSelector]
type Matcher uint8

const (
	NoMatcher       Matcher = iota // [attr] (Does not match values and only checks attribute's presence)
	NormalMatcher                  // [attr=value]
	TildeMatcher                   // [attr~=value]
	BarMatcher                     // [attr|=value]
	CaretMatcher                   // [attr^=value]
	DollarMatcher                  // [attr$=value]
	AsteriskMatcher                // [attr*=value]
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
