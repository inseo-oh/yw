package libhtml

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"yw/libgfx"
	"yw/libplatform"
)

type browser_paint_node interface {
	paint(dest *image.RGBA)
	String() string
}

type browser_text_paint_node struct {
	text_layout_node browser_layout_text_node
	font             libplatform.Font
	font_size        float64
	color            color.RGBA
}

func (t browser_text_paint_node) paint(dest *image.RGBA) {
	t.font.SetTextSize(int(t.font_size))
	metrics := t.font.Metrics()
	x := t.text_layout_node.rect.Left
	baseline_y := t.text_layout_node.rect.Top + metrics.Ascender
	text := t.text_layout_node.text
	t.font.DrawText(text, dest, x, baseline_y, t.color)
}
func (t browser_text_paint_node) String() string {
	return fmt.Sprintf("text-paint(%v) %v %g", t.text_layout_node, t.color, t.font_size)
}

type browser_box_paint_node struct {
	items            []browser_paint_node
	rect             libgfx.Rect
	background_color color.RGBA
}

func (g browser_box_paint_node) paint(dest *image.RGBA) {
	color_img := image.NewRGBA(image.Rect(0, 0, int(g.rect.Width)-1, int(g.rect.Height)-1))
	for y_off := range int(g.rect.Height) {
		for x_off := range int(g.rect.Width) {
			color_img.Set(x_off, y_off, g.background_color)
		}
	}
	draw.Draw(dest, image.Rect(0, 0, int(g.rect.Width)-1, int(g.rect.Height)-1), color_img, image.Point{int(g.rect.Left), int(g.rect.Top)}, draw.Over)
	for _, item := range g.items {
		item.paint(dest)
	}
}
func (t browser_box_paint_node) String() string {
	return fmt.Sprintf("box-paint (%d items)", len(t.items))
}

func browser_print_paint_tree(node browser_paint_node) {
	var do_print func(node browser_paint_node, indent_level int)
	do_print = func(node browser_paint_node, indent_level int) {
		curr_node := node
		indent := strings.Repeat(" ", indent_level*4)
		fmt.Printf("%s%v\n", indent, node)
		if gpaint, ok := curr_node.(browser_box_paint_node); ok {
			for _, child := range gpaint.items {
				do_print(child, indent_level+1)
			}
		}
	}
	do_print(node, 0)

}
