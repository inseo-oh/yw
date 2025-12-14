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
	currentNaturalPos  float64
	leftFloatingBoxes  []Box
	rightFloatingBoxes []Box
}

func (bfc blockFormattingContext) naturalPos() float64 {
	return bfc.currentNaturalPos
}
func (bfc *blockFormattingContext) incrementNaturalPos(pos float64) {
	bfc.currentNaturalPos += pos
	if pos < 0 {
		log.Printf("warning: attempted to increment natural position with negative value %g", pos)
	}
}
func (bfc *blockFormattingContext) leftFloatWidth(forBlockPos float64) float64 {
	sum := 0.0
	for _, bx := range bfc.leftFloatingBoxes {
		mRect := bx.boxMarginRect()
		if mRect.blockPos <= forBlockPos && forBlockPos <= (mRect.blockPos+mRect.logicalHeight) {
			sum += bx.logicalWidth()
		}
	}
	return sum
}
func (bfc *blockFormattingContext) rightFloatWidth(forBlockPos float64) float64 {
	sum := 0.0
	for _, bx := range bfc.rightFloatingBoxes {
		mRect := bx.boxMarginRect()
		if mRect.blockPos <= forBlockPos && forBlockPos <= (mRect.blockPos+mRect.logicalHeight) {
			sum += bx.logicalWidth()
		}
	}
	return sum
}
func (bfc *blockFormattingContext) floatWidth(forBlockPos float64) float64 {
	return bfc.leftFloatWidth(forBlockPos) + bfc.rightFloatWidth(forBlockPos)
}
