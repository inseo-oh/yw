// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package layout

// Block Formatting Contexts(BFC for short) are responsible for tracking Y-axis,
// or more accurately, the opposite axis of writing mode.
// (English uses X-axis for writing text, so BFC's position grows Y-axis)
//
// https://www.w3.org/TR/CSS2/visuren.html#block-formatting
type blockFormattingContext struct {
	formattingContextCommon
	currentNaturalPos float64
}

func (bfc blockFormattingContext) formattingContextType() formattingContextType {
	return formattingContextTypeBlock
}
func (bfc blockFormattingContext) naturalPos() float64 {
	return bfc.currentNaturalPos
}
func (bfc *blockFormattingContext) incrementNaturalPos(pos float64) {
	bfc.currentNaturalPos += pos
}

// TODO: Use this thing for every BFC creation, and make similar function for IFC as well.
func makeBfc(creatorBox box) *blockFormattingContext {
	bfc := blockFormattingContext{}
	bfc.creatorBox = creatorBox
	return &bfc
}
