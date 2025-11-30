package layout

import (
	"yw/css/csscolor"
	"yw/css/cssom"
	"yw/dom"
	"yw/gfx"
	"yw/util"
)

type box interface {
	Node
	boxElement() dom.Element
	boxRect() gfx.Rect
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
	rect       gfx.Rect
	widthAuto  bool
	heightAuto bool
	childBoxes []box
	childTexts []*Text
}

func (bx boxCommon) boxElement() dom.Element { return bx.elem }
func (bx boxCommon) boxRect() gfx.Rect       { return bx.rect }
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
	bx.rect.Width += width
	bx.rect.Height += height
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
	wDiff := max(minWidth-bx.rect.Width, 0)
	hDiff := max(minHeight-bx.rect.Height, 0)
	bx.incrementSize(wDiff, hDiff)
}
func (bx boxCommon) MakePaintNode() gfx.PaintNode {
	paintNodes := []gfx.PaintNode{}

	var color = csscolor.Transparent()
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
	return gfx.BoxPaint{Items: paintNodes, Color: rgbaColor, Rect: bx.rect}
}
