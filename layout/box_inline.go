// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import (
	"fmt"
)

// https://www.w3.org/TR/css-display-3/#inline-box
type InlineBox struct {
	boxCommon
	ParentBcon *BlockContainerBox
}

func (bx InlineBox) String() string {
	physMarginRect := bx.MarginRect.ToPhysicalRect()
	leftStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Left, bx.Margin.Left, bx.Padding.Left)
	topStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Top, bx.Margin.Top, bx.Padding.Top)
	rightStr := fmt.Sprintf("%g-%g-%g", physMarginRect.right(), bx.Margin.Right, bx.Padding.Right)
	bottomStr := fmt.Sprintf("%g-%g-%g", physMarginRect.bottom(), bx.Margin.Bottom, bx.Padding.Bottom)
	return fmt.Sprintf(
		"inline-box [elem %v] at [LTRB %s %s %s %s (%gx%g)]",
		bx.Elem, leftStr, topStr, rightStr, bottomStr, physMarginRect.Width, physMarginRect.Height)
}
func (bx InlineBox) isBlockLevel() bool { return false }
