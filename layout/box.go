// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import (
	"image"
	"image/color"

	"github.com/inseo-oh/yw/css/csscolor"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx/paint"
	"github.com/inseo-oh/yw/util"
)

// Box represents a box in the box tree
type Box interface {
	ChildBoxes() []Box
	ChildTexts() []*text

	// MakePaintNode creates a paint node for given node and its children.
	// (So calling this on the root node will generate paint tree for the whole page)
	MakePaintNode() paint.Node

	boxParent() Box
	boxElement() dom.Element
	boxMarginRect() logicalRect
	boxContentRect() logicalRect
	boxMargin() physicalEdges
	boxPadding() physicalEdges
	logicalWidth() LogicalPos
	logicalHeight() LogicalPos
	isWidthAuto() bool
	isHeightAuto() bool
	incrementSize(logicalWidth, logicalHeight LogicalPos)
	incrementIfNeeded(minLogicalWidth, minLogicalHeight LogicalPos) (logicalWidthDiff, logicalHeightDiff LogicalPos)
}
type boxCommon struct {
	parent             Box
	elem               dom.Element
	marginRect         logicalRect
	margin             physicalEdges
	padding            physicalEdges
	physicalWidthAuto  bool
	physicalHeightAuto bool
	childBoxes         []Box
	childTexts         []*text
}

func (bx boxCommon) boxParent() Box              { return bx.parent }
func (bx boxCommon) boxElement() dom.Element     { return bx.elem }
func (bx boxCommon) boxMarginRect() logicalRect  { return bx.marginRect }                              // Rect containing margin area
func (bx boxCommon) boxPaddingRect() logicalRect { return bx.boxMarginRect().addPadding(bx.margin) }   // Rect containing padding area
func (bx boxCommon) boxContentRect() logicalRect { return bx.boxPaddingRect().addPadding(bx.padding) } // Rect containing content area
func (bx boxCommon) boxMargin() physicalEdges    { return bx.margin }
func (bx boxCommon) boxPadding() physicalEdges   { return bx.padding }

// https://www.w3.org/TR/css-writing-modes-4/#logical-width
func (bx boxCommon) logicalWidth() LogicalPos {
	return bx.boxContentRect().logicalWidth
}

// https://www.w3.org/TR/css-writing-modes-4/#logical-height
func (bx boxCommon) logicalHeight() LogicalPos {
	return bx.boxContentRect().logicalHeight
}

func (bx boxCommon) ChildBoxes() []Box {
	return bx.childBoxes
}
func (bx boxCommon) ChildTexts() []*text {
	return bx.childTexts
}
func (bx boxCommon) isWidthAuto() bool  { return bx.physicalWidthAuto }
func (bx boxCommon) isHeightAuto() bool { return bx.physicalHeightAuto }

func (bx *boxCommon) incrementSize(logicalWidth, logicalHeight LogicalPos) {
	if logicalWidth == 0 && logicalHeight == 0 {
		return
	}
	bx.marginRect.logicalWidth += logicalWidth
	bx.marginRect.logicalHeight += logicalHeight
	parent := bx.parent
	if !util.IsNil(parent) {
		w := logicalWidth
		h := logicalHeight
		if !parent.isWidthAuto() {
			w = 0
		}
		if !parent.isHeightAuto() {
			h = 0
		}
		parent.incrementSize(w, h)
	}
}
func (bx *boxCommon) incrementIfNeeded(minLogicalWidth, minLogicalHeight LogicalPos) (logicalWidthDiff, logicalHeightDiff LogicalPos) {
	wDiff := max(minLogicalWidth-bx.boxContentRect().logicalWidth, 0)
	hDiff := max(minLogicalHeight-bx.boxContentRect().logicalHeight, 0)
	bx.incrementSize(wDiff, hDiff)
	return wDiff, hDiff
}

func (bx boxCommon) MakePaintNode() paint.Node {
	var col color.Color
	paintNodes := []paint.Node{}
	if !util.IsNil(bx.elem) {
		var color = csscolor.Transparent
		styleSetSource := cssom.ComputedStyleSetSourceOf(bx.elem)
		styleSet := styleSetSource.ComputedStyleSet()
		if bx.elem != nil {
			color = styleSet.BackgroundColor()
		}
		col = color.ToStdColor(styleSetSource.CurrentColor())
	}
	for _, child := range bx.ChildBoxes() {
		paintNodes = append(paintNodes, child.MakePaintNode())
	}
	for _, child := range bx.ChildTexts() {
		paintNodes = append(paintNodes, child.MakePaintNode())
	}
	paddingPhysRect := bx.boxPaddingRect().toPhysicalRect()
	contentRect := image.Rect(
		int(paddingPhysRect.Left),
		int(paddingPhysRect.Top),
		int(paddingPhysRect.right()),
		int(paddingPhysRect.bottom()),
	)
	return paint.BoxPaint{Items: paintNodes, Color: col, Rect: contentRect}
}
