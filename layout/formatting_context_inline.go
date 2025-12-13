// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package layout

// Inline Formatting Contexts(IFC for short) are responsible for tracking X-axis,
// or more accurately, the primary axis of writing mode.
// (English uses X-axis for writing text, so IFC's position grows X-axis)
//
// or can be also thought as "The opposite axis of BFC", if you really want :D
//
// https://www.w3.org/TR/CSS2/visuren.html#inline-formatting
// https://www.w3.org/TR/css-inline-3/#inline-formatting-context
type inlineFormattingContext struct {
	formattingContextCommon

	bcon                  *blockContainer // Block container containing this inline node
	lineBoxes             []lineBox       // List of line boxes
	initialAvailableWidth float64
	initialBlockPos       float64
	writtenText           string
}

func (ifc *inlineFormattingContext) addLineBox(lineHeight float64) {
	lb := lineBox{}
	lb.currentLineHeight = lineHeight
	if len(ifc.lineBoxes) != 0 {
		lastLb := ifc.currentLineBox()
		lb.initialBlockPos = lastLb.initialBlockPos + lastLb.currentLineHeight
	} else {
		lb.initialBlockPos = ifc.initialBlockPos
	}
	lb.availableWidth = ifc.initialAvailableWidth
	ifc.lineBoxes = append(ifc.lineBoxes, lb)
}
func (ifc *inlineFormattingContext) currentLineBox() *lineBox {
	return &ifc.lineBoxes[len(ifc.lineBoxes)-1]
}
func (ifc inlineFormattingContext) naturalPos() float64 {
	return ifc.currentLineBox().currentNaturalPos
}
func (ifc *inlineFormattingContext) incrementNaturalPos(pos float64) {
	if len(ifc.lineBoxes) == 0 {
		panic("attempted to increment natural position without creating lineBox")
	}
	lb := ifc.currentLineBox()
	if lb.availableWidth < lb.currentNaturalPos+pos {
		panic("content overflow")
	}
	lb.currentNaturalPos += pos
}

// Line box holds state needed for placing inline contents, such as next inline
// position(which gets reset when entering new line).
//
// https://www.w3.org/TR/css-inline-3/#line-box
type lineBox struct {
	availableWidth    float64
	currentNaturalPos float64
	currentLineHeight float64
	initialBlockPos   float64
}
