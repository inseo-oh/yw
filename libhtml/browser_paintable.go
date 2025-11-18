package libhtml

import (
	"fmt"
	"image"
	"strings"
	"yw/libplatform"
)

type browser_paint_node interface {
	paint(dest *image.RGBA)
	String() string
}

type browser_text_paint_node struct {
	text_layout_node browser_layout_text_node
	font             libplatform.Font
}

func (t browser_text_paint_node) paint(dest *image.RGBA) {
	text := t.text_layout_node.text.get_text()
	// First we draw(to nowhere) with 0, 0 as baseline offset.
	rect := t.font.DrawText(text, nil, 0, 0)
	// Then we figure out where the baseline should be
	baseline_x := t.text_layout_node.rect.Left - rect.Left
	baseline_y := t.text_layout_node.rect.Top - rect.Top // Note that rect.Top would be a negative position
	// And finally we draw the text for real
	t.font.DrawText(text, dest, baseline_x, baseline_y)
}
func (t browser_text_paint_node) String() string {
	return fmt.Sprintf("text-paint %v", t.text_layout_node)
}

type browser_grouped_paint_node struct {
	items []browser_paint_node
}

func (g browser_grouped_paint_node) paint(dest *image.RGBA) {
	for _, item := range g.items {
		item.paint(dest)
	}
}
func (t browser_grouped_paint_node) String() string {
	return fmt.Sprintf("group-paint (%d items)", len(t.items))
}

func browser_print_paint_tree(node browser_paint_node) {
	var do_print func(node browser_paint_node, indent_level int)
	do_print = func(node browser_paint_node, indent_level int) {
		curr_node := node
		indent := strings.Repeat(" ", indent_level*4)
		fmt.Printf("%s%v\n", indent, node)
		if gpaint, ok := curr_node.(browser_grouped_paint_node); ok {
			for _, child := range gpaint.items {
				do_print(child, indent_level+1)
			}
		}
	}
	do_print(node, 0)

}
