// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

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
	boxMarginRect() BoxRect
	boxContentRect() BoxRect
	boxMargin() BoxEdges
	boxPadding() BoxEdges
	logicalWidth(writeMode writeMode) float64
	logicalHeight(writeMode writeMode) float64
	ChildBoxes() []box
	ChildTexts() []*text
	isWidthAuto() bool
	isHeightAuto() bool
	incrementSize(width, height float64)
	incrementIfNeeded(width, height float64) (widthDiff, heightDiff float64) // TODO: Return inline and block size, not width and height.
}
type boxCommon struct {
	nodeCommon
	elem       dom.Element
	marginRect BoxRect
	margin     BoxEdges
	padding    BoxEdges
	widthAuto  bool
	heightAuto bool
	childBoxes []box
	childTexts []*text
}

func (bx boxCommon) boxElement() dom.Element { return bx.elem }
func (bx boxCommon) boxMarginRect() BoxRect  { return bx.marginRect }                              // Rect containing margin area
func (bx boxCommon) boxPaddingRect() BoxRect { return bx.boxMarginRect().AddPadding(bx.margin) }   // Rect containing padding area
func (bx boxCommon) boxContentRect() BoxRect { return bx.boxPaddingRect().AddPadding(bx.padding) } // Rect containing content area
func (bx boxCommon) boxMargin() BoxEdges     { return bx.margin }
func (bx boxCommon) boxPadding() BoxEdges    { return bx.padding }

// https://www.w3.org/TR/css-writing-modes-4/#logical-width
func (bx boxCommon) logicalWidth(writeMode writeMode) float64 {
	if writeMode == writeModeHorizontal {
		return bx.boxContentRect().Width
	}
	return bx.boxContentRect().Height
}

// https://www.w3.org/TR/css-writing-modes-4/#logical-height
func (bx boxCommon) logicalHeight(writeMode writeMode) float64 {
	if writeMode == writeModeHorizontal {
		return bx.boxContentRect().Height
	}
	return bx.boxContentRect().Width
}
func (bx boxCommon) ChildBoxes() []box {
	return bx.childBoxes
}
func (bx boxCommon) ChildTexts() []*text {
	return bx.childTexts
}
func (bx boxCommon) isWidthAuto() bool  { return bx.widthAuto }
func (bx boxCommon) isHeightAuto() bool { return bx.heightAuto }
func (bx *boxCommon) incrementSize(width, height float64) {
	if width == 0 && height == 0 {
		return
	}
	if bx.marginRect.Width+width == 1355 {
		panic("?")
	}
	bx.marginRect.Width += width
	bx.marginRect.Height += height
	parent := bx.parentNode()
	if !util.IsNil(parent) {
		if parent, ok := parent.(box); ok {
			w := width
			h := height
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
func (bx *boxCommon) incrementIfNeeded(minWidth, minHeight float64) (widthDiff, heightDiff float64) {
	wDiff := max(minWidth-bx.boxContentRect().Width, 0)
	hDiff := max(minHeight-bx.boxContentRect().Height, 0)
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
	contentRect := image.Rect(
		int(bx.boxPaddingRect().Left),
		int(bx.boxPaddingRect().Top),
		int(bx.boxPaddingRect().Right()),
		int(bx.boxPaddingRect().Bottom()),
	)
	return paint.BoxPaint{Items: paintNodes, Color: col, Rect: contentRect}
}
