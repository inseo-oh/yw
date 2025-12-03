// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// PseudoClassSelector represents a [CSS pseudo class selector] (e.g. :first-letter)
//
// [CSS pseudo class selector]: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#pseudo-class
type PseudoClassSelector struct {
	Name string
	Args []any
}

func (sel PseudoClassSelector) String() string {
	if len(sel.Args) != 0 {
		// TODO: Display arguments in better way
		return fmt.Sprintf(":%s(%v)", sel.Name, sel.Args)
	} else {
		return fmt.Sprintf(":%s", sel.Name)
	}
}
func (sel PseudoClassSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(PseudoClassSelector); !ok {
		return false
	} else {
		if sel.Name != otherSel.Name {
			return false
		}
		if len(sel.Args) != len(otherSel.Args) {
			return false
		}
		// TODO: Compare actual arguments
	}
	return true
}
func (sel PseudoClassSelector) MatchAgainst(element dom.Element) bool {
	// STUB
	return false
}
