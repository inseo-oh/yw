// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

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

type box interface {
	Node
	boxElement() dom.Element
	boxMarginRect() logicalRect
	boxContentRect() logicalRect
	boxMargin() physicalEdges
	boxPadding() physicalEdges
	logicalWidth() float64
	logicalHeight() float64
	ChildBoxes() []box
	ChildTexts() []*text
	isWidthAuto() bool
	isHeightAuto() bool
	incrementSize(width, height float64)
	incrementIfNeeded(width, height float64) (widthDiff, heightDiff float64) // TODO: Return inline and block size, not width and height.
}
type boxCommon struct {
	nodeCommon
	elem               dom.Element
	marginRect         logicalRect
	margin             physicalEdges
	padding            physicalEdges
	physicalWidthAuto  bool
	physicalHeightAuto bool
	childBoxes         []box
	childTexts         []*text
}

func (bx boxCommon) boxElement() dom.Element     { return bx.elem }
func (bx boxCommon) boxMarginRect() logicalRect  { return bx.marginRect }                              // Rect containing margin area
func (bx boxCommon) boxPaddingRect() logicalRect { return bx.boxMarginRect().addPadding(bx.margin) }   // Rect containing padding area
func (bx boxCommon) boxContentRect() logicalRect { return bx.boxPaddingRect().addPadding(bx.padding) } // Rect containing content area
func (bx boxCommon) boxMargin() physicalEdges    { return bx.margin }
func (bx boxCommon) boxPadding() physicalEdges   { return bx.padding }

// https://www.w3.org/TR/css-writing-modes-4/#logical-width
func (bx boxCommon) logicalWidth() float64 {
	return bx.boxContentRect().logicalWidth
}

// https://www.w3.org/TR/css-writing-modes-4/#logical-height
func (bx boxCommon) logicalHeight() float64 {
	return bx.boxContentRect().logicalHeight
}
func (bx boxCommon) ChildBoxes() []box {
	return bx.childBoxes
}
func (bx boxCommon) ChildTexts() []*text {
	return bx.childTexts
}
func (bx boxCommon) isWidthAuto() bool  { return bx.physicalWidthAuto }
func (bx boxCommon) isHeightAuto() bool { return bx.physicalHeightAuto }
func (bx *boxCommon) incrementSize(logicalWidth, logicalHeight float64) {
	if logicalWidth == 0 && logicalHeight == 0 {
		return
	}
	bx.marginRect.logicalWidth += logicalWidth
	bx.marginRect.logicalHeight += logicalHeight
	parent := bx.parentNode()
	if !util.IsNil(parent) {
		if parent, ok := parent.(box); ok {
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
}
func (bx *boxCommon) incrementIfNeeded(minLogicalWidth, minLogicalHeight float64) (widthDiff, heightDiff float64) {
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
