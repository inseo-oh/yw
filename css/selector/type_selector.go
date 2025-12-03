// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// TypeSelector represents [CSS type selector] (e.g. div)
//
// [CSS type selector]: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#type-selector
type TypeSelector struct {
	TypeName WqName
}

func (sel TypeSelector) String() string { return fmt.Sprintf("%v", sel.TypeName) }
func (sel TypeSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(TypeSelector); !ok {
		return false
	} else {
		if !sel.TypeName.Equals(otherSel.TypeName) {
			return false
		}
	}
	return true
}
func (sel TypeSelector) MatchAgainst(element dom.Element) bool {
	// TODO: Handle namespace
	name := sel.TypeName.Ident
	if element.LocalName() != name {
		return false
	}
	return true
}
