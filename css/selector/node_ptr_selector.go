// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package selector

import "github.com/inseo-oh/yw/dom"

// NodePtrSelector is special CSS selector that matches DOM node pointer directly.
//
// This has not part of CSS spec, and has no CSS representation ([NodePtrSelector.String]
// just returns value of String() of the element)
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

// String returns value of String() of the element.
func (sel NodePtrSelector) String() string {
	return sel.Element.String()
}
