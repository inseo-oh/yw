// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import "log"

// Block Formatting Contexts(BFC for short) are responsible for tracking Y-axis,
// or more accurately, the opposite axis of writing mode.
// (English uses X-axis for writing text, so BFC's position grows Y-axis)
//
// https://www.w3.org/TR/CSS2/visuren.html#block-formatting
type BlockFormattingContext struct {
	formattingContextCommon
	CurrentNaturalPos  LogicalPos
	LeftFloatingBoxes  []Box
	RightFloatingBoxes []Box
}

func (bfc BlockFormattingContext) NaturalPos() LogicalPos {
	return bfc.CurrentNaturalPos
}
func (bfc *BlockFormattingContext) IncrementNaturalPos(pos LogicalPos) {
	bfc.CurrentNaturalPos += pos
	if pos < 0 {
		log.Printf("warning: attempted to increment natural position with negative value %g", pos)
	}
}
func (bfc *BlockFormattingContext) leftFloatWidth(forLogicalY LogicalPos) LogicalPos {
	sum := LogicalPos(0.0)
	for _, bx := range bfc.LeftFloatingBoxes {
		mRect := bx.BoxMarginRect()
		if mRect.LogicalY <= forLogicalY && forLogicalY <= (mRect.LogicalY+mRect.LogicalHeight) {
			sum += bx.LogicalWidth()
		}
	}
	return sum
}
func (bfc *BlockFormattingContext) rightFloatWidth(forLogicalY LogicalPos) LogicalPos {
	sum := LogicalPos(0.0)
	for _, bx := range bfc.RightFloatingBoxes {
		mRect := bx.BoxMarginRect()
		if mRect.LogicalY <= forLogicalY && forLogicalY <= (mRect.LogicalY+mRect.LogicalHeight) {
			sum += bx.LogicalWidth()
		}
	}
	return sum
}
func (bfc *BlockFormattingContext) floatWidth(forLogicalY LogicalPos) LogicalPos {
	return bfc.leftFloatWidth(forLogicalY) + bfc.rightFloatWidth(forLogicalY)
}
