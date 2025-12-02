package layout

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/gfx/paint"
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
func (txt Text) MakePaintNode() paint.PaintNode {
	return paint.TextPaint{
		Left:  txt.rect.Left,
		Top:   txt.rect.Top,
		Text:  txt.text,
		Font:  txt.font,
		Size:  txt.fontSize,
		Color: txt.color,
	}
}
func (txt Text) IsBlockLevel() bool { return false }
