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
	ChildTexts() []*Text
	AddChildBox(b Box)
	AddChildText(t *Text)

	// MakePaintNode creates a paint node for given node and its children.
	// (So calling this on the root node will generate paint tree for the whole page)
	MakePaintNode() paint.Node

	BoxParent() Box
	BoxElement() dom.Element
	BoxMarginRect() LogicalRect
	BoxContentRect() LogicalRect
	LogicalWidth() LogicalPos
	LogicalHeight() LogicalPos
	IsWidthAuto() bool
	IsHeightAuto() bool
	IncrementSize(logicalWidth, logicalHeight LogicalPos)
	IncrementIfNeeded(minLogicalWidth, minLogicalHeight LogicalPos) (logicalWidthDiff, logicalHeightDiff LogicalPos)
}
type boxCommon struct {
	Parent             Box
	Elem               dom.Element
	MarginRect         LogicalRect
	Margin             PhysicalEdges
	Padding            PhysicalEdges
	PhysicalWidthAuto  bool
	PhysicalHeightAuto bool
	childBoxes         []Box
	childTexts         []*Text
}

func (bx boxCommon) BoxParent() Box              { return bx.Parent }
func (bx boxCommon) BoxElement() dom.Element     { return bx.Elem }
func (bx boxCommon) BoxMarginRect() LogicalRect  { return bx.MarginRect }                              // Rect containing margin area
func (bx boxCommon) BoxPaddingRect() LogicalRect { return bx.BoxMarginRect().addPadding(bx.Margin) }   // Rect containing padding area
func (bx boxCommon) BoxContentRect() LogicalRect { return bx.BoxPaddingRect().addPadding(bx.Padding) } // Rect containing content area

// https://www.w3.org/TR/css-writing-modes-4/#logical-width
func (bx boxCommon) LogicalWidth() LogicalPos {
	return bx.BoxContentRect().LogicalWidth
}

// https://www.w3.org/TR/css-writing-modes-4/#logical-height
func (bx boxCommon) LogicalHeight() LogicalPos {
	return bx.BoxContentRect().LogicalHeight
}

func (bx boxCommon) ChildBoxes() []Box {
	return bx.childBoxes
}
func (bx boxCommon) AddChildBox(b Box) {
	bx.childBoxes = append(bx.childBoxes, b)
}
func (bx boxCommon) ChildTexts() []*Text {
	return bx.childTexts
}
func (bx boxCommon) AddChildText(t *Text) {
	bx.childTexts = append(bx.childTexts, t)
}

func (bx boxCommon) IsWidthAuto() bool  { return bx.PhysicalWidthAuto }
func (bx boxCommon) IsHeightAuto() bool { return bx.PhysicalHeightAuto }

func (bx *boxCommon) IncrementSize(logicalWidth, logicalHeight LogicalPos) {
	if logicalWidth == 0 && logicalHeight == 0 {
		return
	}
	bx.MarginRect.LogicalWidth += logicalWidth
	bx.MarginRect.LogicalHeight += logicalHeight
	parent := bx.Parent
	if !util.IsNil(parent) {
		w := logicalWidth
		h := logicalHeight
		if !parent.IsWidthAuto() {
			w = 0
		}
		if !parent.IsHeightAuto() {
			h = 0
		}
		parent.IncrementSize(w, h)
	}
}
func (bx *boxCommon) IncrementIfNeeded(minLogicalWidth, minLogicalHeight LogicalPos) (logicalWidthDiff, logicalHeightDiff LogicalPos) {
	wDiff := max(minLogicalWidth-bx.BoxContentRect().LogicalWidth, 0)
	hDiff := max(minLogicalHeight-bx.BoxContentRect().LogicalHeight, 0)
	bx.IncrementSize(wDiff, hDiff)
	return wDiff, hDiff
}

func (bx boxCommon) MakePaintNode() paint.Node {
	var col color.Color
	paintNodes := []paint.Node{}
	if !util.IsNil(bx.Elem) {
		var color = csscolor.Transparent
		styleSetSource := cssom.ComputedStyleSetSourceOf(bx.Elem)
		styleSet := styleSetSource.ComputedStyleSet()
		if bx.Elem != nil {
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
	paddingPhysRect := bx.BoxPaddingRect().ToPhysicalRect()
	contentRect := image.Rect(
		int(paddingPhysRect.Left),
		int(paddingPhysRect.Top),
		int(paddingPhysRect.right()),
		int(paddingPhysRect.bottom()),
	)
	return paint.BoxPaint{Items: paintNodes, Color: col, Rect: contentRect}
}
