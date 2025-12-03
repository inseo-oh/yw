// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package selector

import (
	"fmt"
	"slices"
	"strings"

	"github.com/inseo-oh/yw/dom"
)

// ClassSelector represents a [CSS class selector] (e.g. .foo)
//
// [CSS class selector]: https://www.w3.org/TR/2022/WD-selectors-4-20221111/#class-selector
type ClassSelector struct{ Class string }

func (sel ClassSelector) String() string { return fmt.Sprintf(".%s", sel.Class) }
func (sel ClassSelector) Equals(other Selector) bool {
	if otherSel, ok := other.(ClassSelector); !ok {
		return false
	} else {
		if sel.Class != otherSel.Class {
			return false
		}
	}
	return true
}
func (sel ClassSelector) MatchAgainst(element dom.Element) bool {
	class := sel.Class
	classes, ok := element.AttrWithoutNamespace("class")
	if !ok {
		return false
	}
	classList := strings.Split(classes, " ")
	return slices.Contains(classList, class)
}
