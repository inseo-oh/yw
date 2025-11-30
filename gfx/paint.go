package gfx

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
)

type PaintNode interface {
	Paint(dest *image.RGBA)
	String() string
}

type TextPaint struct {
	Left, Top float64
	Text      string
	Font      Font
	Size      float64
	Color     color.RGBA
}

func (t TextPaint) Paint(dest *image.RGBA) {
	t.Font.SetTextSize(int(t.Size))
	metrics := t.Font.Metrics()
	x := t.Left
	baselineY := t.Top + metrics.Ascender
	text := t.Text
	t.Font.DrawText(text, dest, x, baselineY, t.Color)
}
func (t TextPaint) String() string {
	return fmt.Sprintf("text-paint(%s) %v %g", t.Text, t.Color, t.Size)
}

type BoxPaint struct {
	Items []PaintNode
	Rect  Rect
	Color color.RGBA
}

func (g BoxPaint) Paint(dest *image.RGBA) {
	colorImg := image.NewRGBA(image.Rect(0, 0, int(g.Rect.Width)-1, int(g.Rect.Height)-1))
	for yOff := range int(g.Rect.Height) {
		for xOff := range int(g.Rect.Width) {
			colorImg.Set(xOff, yOff, g.Color)
		}
	}
	draw.Draw(dest, image.Rect(int(g.Rect.Left), int(g.Rect.Top), int(g.Rect.Left+g.Rect.Width)-1, int(g.Rect.Top+g.Rect.Height)-1), colorImg, image.Point{0, 0}, draw.Over)
	for _, item := range g.Items {
		item.Paint(dest)
	}
}
func (t BoxPaint) String() string {
	return fmt.Sprintf("box-paint(color=%v, rect=%v, %d items)", t.Color, t.Rect, len(t.Items))
}

func PrintPaintTree(node PaintNode) {
	var doPrint func(node PaintNode, indentLevel int)
	doPrint = func(node PaintNode, indentLevel int) {
		currNode := node
		indent := strings.Repeat(" ", indentLevel*4)
		fmt.Printf("%s%v\n", indent, node)
		if gpaint, ok := currNode.(BoxPaint); ok {
			for _, child := range gpaint.Items {
				doPrint(child, indentLevel+1)
			}
		}
	}
	doPrint(node, 0)

}
