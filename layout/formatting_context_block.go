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
type blockFormattingContext struct {
	formattingContextCommon
	currentNaturalPos  LogicalPos
	leftFloatingBoxes  []Box
	rightFloatingBoxes []Box
}

func (bfc blockFormattingContext) naturalPos() LogicalPos {
	return bfc.currentNaturalPos
}
func (bfc *blockFormattingContext) incrementNaturalPos(pos LogicalPos) {
	bfc.currentNaturalPos += pos
	if pos < 0 {
		log.Printf("warning: attempted to increment natural position with negative value %g", pos)
	}
}
func (bfc *blockFormattingContext) leftFloatWidth(forLogicalY LogicalPos) LogicalPos {
	sum := LogicalPos(0.0)
	for _, bx := range bfc.leftFloatingBoxes {
		mRect := bx.boxMarginRect()
		if mRect.logicalY <= forLogicalY && forLogicalY <= (mRect.logicalY+mRect.logicalHeight) {
			sum += bx.logicalWidth()
		}
	}
	return sum
}
func (bfc *blockFormattingContext) rightFloatWidth(forLogicalY LogicalPos) LogicalPos {
	sum := LogicalPos(0.0)
	for _, bx := range bfc.rightFloatingBoxes {
		mRect := bx.boxMarginRect()
		if mRect.logicalY <= forLogicalY && forLogicalY <= (mRect.logicalY+mRect.logicalHeight) {
			sum += bx.logicalWidth()
		}
	}
	return sum
}
func (bfc *blockFormattingContext) floatWidth(forLogicalY LogicalPos) LogicalPos {
	return bfc.leftFloatWidth(forLogicalY) + bfc.rightFloatWidth(forLogicalY)
}
