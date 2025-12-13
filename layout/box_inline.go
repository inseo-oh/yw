// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import (
	"fmt"
)

// https://www.w3.org/TR/css-display-3/#inline-box
type inlineBox struct {
	boxCommon
	parentBcon *blockContainer
}

func (bx inlineBox) String() string {
	physMarginRect := bx.marginRect.toPhysicalRect()
	leftStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Left, bx.margin.left, bx.padding.left)
	topStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Top, bx.margin.top, bx.padding.top)
	rightStr := fmt.Sprintf("%g-%g-%g", physMarginRect.right(), bx.margin.right, bx.padding.right)
	bottomStr := fmt.Sprintf("%g-%g-%g", physMarginRect.bottom(), bx.margin.bottom, bx.padding.bottom)
	return fmt.Sprintf(
		"inline-box [elem %v] at [LTRB %s %s %s %s (%gx%g)]",
		bx.elem, leftStr, topStr, rightStr, bottomStr, physMarginRect.Width, physMarginRect.Height)
}
func (bx inlineBox) isBlockLevel() bool { return false }
