// Package paint provides paint tree and painting operations.
//
// Paint tree is generated during layout process(see layout package), and
// drawn to the destination image by Paint()ing the root node.
package paint

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/inseo-oh/yw/gfx"
)

// Node represents a node in the paint tree.
type PaintNode interface {
	// Paint paints to the dest.
	Paint(dest *image.RGBA)

	// String returns decscription of the node.
	String() string
}

// Text is Node that paints a text.
type TextPaint struct {
	Left, Top float64
	Text      string
	Font      gfx.Font
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

// Text is Node that paints a box.
type BoxPaint struct {
	Items []PaintNode
	Rect  gfx.Rect
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

// PrintTree prints paint tree to the standard output.
func PrintTree(node PaintNode) {
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
