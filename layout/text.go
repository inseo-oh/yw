package layout

import (
	"fmt"
	"image/color"
	"strconv"
	"yw/gfx"
)

type Text struct {
	NodeCommon
	rect     gfx.Rect
	text     string
	font     gfx.Font
	fontSize float64
	color    color.RGBA
}

func (txt Text) String() string {
	return fmt.Sprintf("text %s at [%v]", strconv.Quote(txt.text), txt.rect)
}
func (txt Text) NodeType() NodeType {
	return NodeTypeText
}
func (txt Text) MakePaintNode() gfx.PaintNode {
	return gfx.TextPaint{
		Left:  txt.rect.Left,
		Top:   txt.rect.Top,
		Text:  txt.text,
		Font:  txt.font,
		Size:  txt.fontSize,
		Color: txt.color,
	}
}
func (txt Text) IsBlockLevel() bool { return false }
