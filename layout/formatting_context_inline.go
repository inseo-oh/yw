// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

// Inline Formatting Contexts(IFC for short) are responsible for tracking X-axis,
// or more accurately, the primary axis of writing mode.
// (English uses X-axis for writing text, so IFC's position grows X-axis)
//
// or can be also thought as "The opposite axis of BFC", if you really want :D
//
// https://www.w3.org/TR/CSS2/visuren.html#inline-formatting
// https://www.w3.org/TR/css-inline-3/#inline-formatting-context
type InlineFormattingContext struct {
	formattingContextCommon

	BlockContainer        *BlockContainerBox // Block container containing this inline node
	LineBoxes             []lineBox          // List of line boxes
	InitialAvailableWidth LogicalPos
	InitialLogicalY       LogicalPos
	WrittenText           string
}

func (ifc *InlineFormattingContext) AddLineBox(lineHeight float64) {
	lb := lineBox{}
	lb.CurrentLineHeight = lineHeight
	if len(ifc.LineBoxes) != 0 {
		lastLb := ifc.CurrentLineBox()
		lb.InitialLogicalY = lastLb.InitialLogicalY + LogicalPos(lastLb.CurrentLineHeight)
	} else {
		lb.InitialLogicalY = ifc.InitialLogicalY
	}
	lb.AvailableWidth = ifc.InitialAvailableWidth - ifc.BlockContainer.Bfc.floatWidth(lb.InitialLogicalY)
	lb.leftOffset = PhysicalPos(ifc.BlockContainer.Bfc.leftFloatWidth(lb.InitialLogicalY))
	ifc.LineBoxes = append(ifc.LineBoxes, lb)
}
func (ifc *InlineFormattingContext) CurrentLineBox() *lineBox {
	return &ifc.LineBoxes[len(ifc.LineBoxes)-1]
}
func (ifc InlineFormattingContext) NaturalPos() LogicalPos {
	return ifc.CurrentLineBox().CurrentNaturalPos + LogicalPos(ifc.CurrentLineBox().leftOffset)
}
func (ifc *InlineFormattingContext) IncrementNaturalPos(pos LogicalPos) {
	if len(ifc.LineBoxes) == 0 {
		panic("attempted to increment natural position without creating lineBox")
	}
	lb := ifc.CurrentLineBox()
	if lb.AvailableWidth < lb.CurrentNaturalPos+pos {
		panic("content overflow")
	}
	lb.CurrentNaturalPos += pos
}

// Line box holds state needed for placing inline contents, such as next inline
// position(which gets reset when entering new line).
//
// https://www.w3.org/TR/css-inline-3/#line-box
type lineBox struct {
	leftOffset        PhysicalPos
	AvailableWidth    LogicalPos
	CurrentNaturalPos LogicalPos
	CurrentLineHeight float64
	InitialLogicalY   LogicalPos
}
