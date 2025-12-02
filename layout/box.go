package layout

import (
	"github.com/inseo-oh/yw/css/csscolor"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/gfx/paint"
	"github.com/inseo-oh/yw/util"
)

type box interface {
	Node
	boxElement() dom.Element
	boxMarginRect() gfx.Rect
	boxContentRect() gfx.Rect
	boxMargin() gfx.Edges
	logicalWidth(writeMode writeMode) float64
	logicalHeight(writeMode writeMode) float64
	ChildBoxes() []box
	ChildTexts() []*Text
	isWidthAuto() bool
	isHeightAuto() bool
	incrementSize(width, height float64)
	incrementIfNeeded(width, height float64)
}
type boxCommon struct {
	NodeCommon
	elem       dom.Element
	marginRect gfx.Rect
	margin     gfx.Edges
	widthAuto  bool
	heightAuto bool
	childBoxes []box
	childTexts []*Text
}

func (bx boxCommon) boxElement() dom.Element  { return bx.elem }
func (bx boxCommon) boxMarginRect() gfx.Rect  { return bx.marginRect }                       // Rect containing margin area
func (bx boxCommon) boxContentRect() gfx.Rect { return bx.marginRect.AddPadding(bx.margin) } // Rect containing content area
func (bx boxCommon) boxMargin() gfx.Edges     { return bx.margin }

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
func (bx boxCommon) ChildTexts() []*Text {
	return bx.childTexts
}
func (bx boxCommon) isWidthAuto() bool  { return bx.widthAuto }
func (bx boxCommon) isHeightAuto() bool { return bx.heightAuto }
func (bx *boxCommon) incrementSize(width, height float64) {
	if width == 0 && height == 0 {
		return
	}
	bx.marginRect.Width += width
	bx.marginRect.Height += height
	parent := bx.ParentNode()
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
func (bx *boxCommon) incrementIfNeeded(minWidth, minHeight float64) {
	wDiff := max(minWidth-bx.boxContentRect().Width, 0)
	hDiff := max(minHeight-bx.boxContentRect().Height, 0)
	bx.incrementSize(wDiff, hDiff)
}
func (bx boxCommon) MakePaintNode() paint.PaintNode {
	paintNodes := []paint.PaintNode{}

	var color = csscolor.Transparent
	if bx.elem != nil {
		color = cssom.ElementDataOf(bx.elem).ComputedStyleSet.BackgroundColor()
	}
	rgbaColor := color.ToRgba()

	for _, child := range bx.ChildBoxes() {
		paintNodes = append(paintNodes, child.MakePaintNode())
	}
	for _, child := range bx.ChildTexts() {
		paintNodes = append(paintNodes, child.MakePaintNode())
	}
	return paint.BoxPaint{Items: paintNodes, Color: rgbaColor, Rect: bx.boxContentRect()}
}
