// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package selector

import (
	"fmt"

	"github.com/inseo-oh/yw/dom"
)

// UniversalSelector represents [CSS Universal selector] (e.g. *)
//
// [CSS Universal selector]: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#the-universal-selector
type UniversalSelector struct {
	NsPrefix *NsPrefix
}

func (sel UniversalSelector) String() string { return fmt.Sprintf("%v*", sel.NsPrefix) }
func (sel UniversalSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(UniversalSelector); !ok {
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
func (sel UniversalSelector) MatchAgainst(element dom.Element) bool {
	return true
}
