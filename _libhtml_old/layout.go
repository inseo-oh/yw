package libhtml

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/platform"
	cm "github.com/inseo-oh/yw/util"
)

type layout_node interface {
	get_node_type() layout_node_type
	get_parent() layout_node
	make_paint_node() gfx.PaintNode
	is_block_level() bool
	String() string
}
type layout_node_common struct {
	parent layout_node
}

type layout_node_type uint8

const (
	layout_node_type_inline_box = layout_node_type(iota)
	layout_node_type_block_container
	layout_node_type_initial_containing_block
	layout_node_type_text
)

func (n layout_node_common) get_parent() layout_node {
	return n.parent
}

//==============================================================================
// Formatting contexts
//==============================================================================

type layout_formatting_context interface {
	get_formatting_context_type() layout_formatting_context_type
	get_current_natural_pos() float64
	increment_natural_pos(inc float64)
	get_creator_box() layout_box_node
}

type layout_formatting_context_common struct {
	is_dummy_context bool
	creator_box      layout_box_node
}
type layout_formatting_context_type uint8

const (
	layout_formatting_context_type_block = layout_formatting_context_type(iota)
	layout_formatting_context_type_inline
)

type layout_write_mode uint8

const (
	layout_write_mode_horizontal = layout_write_mode(iota)
	layout_write_mode_vertical
)

func (fc layout_formatting_context_common) get_creator_box() layout_box_node {
	return fc.creator_box
}

// Block Formatting Contexts(BFC for short) are responsible for tracking Y-axis,
// or more accurately, the opposite axis of writing mode.
// (English uses X-axis for writing text, so BFC's position grows Y-axis)
//
// https://www.w3.org/TR/CSS2/visuren.html#block-formatting
type layout_block_formatting_context struct {
	layout_formatting_context_common
	current_natural_pos float64
}

func (bfc layout_block_formatting_context) get_formatting_context_type() layout_formatting_context_type {
	return layout_formatting_context_type_block
}
func (bfc layout_block_formatting_context) get_current_natural_pos() float64 {
	return bfc.current_natural_pos
}
func (bfc *layout_block_formatting_context) increment_natural_pos(pos float64) {
	bfc.current_natural_pos += pos
}

// TODO: Use this thing for every BFC creation, and make similar function for IFC as well.
func layout_make_bfc(creator_box layout_box_node) *layout_block_formatting_context {
	bfc := layout_block_formatting_context{}
	bfc.creator_box = creator_box
	return &bfc
}

// Inline Formatting Contexts(IFC for short) are responsible for tracking X-axis,
// or more accurately, the primary axis of writing mode.
// (English uses X-axis for writing text, so IFC's position grows X-axis)
//
// or can be also thought as "The opposite axis of BFC", if you really want :D
//
// https://www.w3.org/TR/CSS2/visuren.html#inline-formatting
// https://www.w3.org/TR/css-inline-3/#inline-formatting-context
type layout_inline_formatting_context struct {
	layout_formatting_context_common

	bcon       *layout_block_container_node // Block container containing this inline node
	line_boxes []layout_line_box            // List of line boxes
}

func (ifc layout_inline_formatting_context) get_formatting_context_type() layout_formatting_context_type {
	return layout_formatting_context_type_inline
}
func (ifc *layout_inline_formatting_context) add_line_box(bfc *layout_block_formatting_context) {
	lb := layout_line_box{}
	lb.current_line_height = 0
	if len(ifc.line_boxes) != 0 {
		last_lb := ifc.get_current_line_box()
		lb.initial_block_pos = last_lb.initial_block_pos + last_lb.current_line_height
	} else {
		lb.initial_block_pos = bfc.get_current_natural_pos()
	}
	lb.available_width = ifc.bcon.rect.Width
	if !ifc.is_dummy_context && lb.available_width == 0 {
		browser_print_layout_tree(ifc.bcon)
	}
	ifc.line_boxes = append(ifc.line_boxes, lb)
}
func (ifc *layout_inline_formatting_context) get_current_line_box() *layout_line_box {
	return &ifc.line_boxes[len(ifc.line_boxes)-1]
}
func (ifc layout_inline_formatting_context) get_current_natural_pos() float64 {
	return ifc.get_current_line_box().current_natural_pos
}
func (ifc *layout_inline_formatting_context) increment_natural_pos(pos float64) {
	if len(ifc.line_boxes) == 0 {
		ifc.add_line_box(ifc.bcon.bfc)
	}
	lb := ifc.get_current_line_box()
	if lb.available_width < lb.current_natural_pos+pos && !ifc.is_dummy_context {
		panic("content overflow")
	}
	lb.current_natural_pos += pos
}

// Line box holds state needed for placing inline contents, such as next inline
// position(which gets reset when entering new line).
//
// https://www.w3.org/TR/css-inline-3/#line-box
type layout_line_box struct {
	available_width     float64
	current_natural_pos float64
	current_line_height float64
	initial_block_pos   float64
}

//==============================================================================
// Text
//==============================================================================

type layout_text_node struct {
	layout_node_common
	rect      gfx.Rect
	text      string
	font      gfx.Font
	font_size float64
	color     color.RGBA
}

func (txt layout_text_node) String() string {
	return fmt.Sprintf("text %s at [%v]", strconv.Quote(txt.text), txt.rect)
}
func (txt layout_text_node) get_node_type() layout_node_type {
	return layout_node_type_text
}
func (txt layout_text_node) make_paint_node() gfx.PaintNode {
	return gfx.TextPaint{
		Left:  txt.rect.Left,
		Top:   txt.rect.Top,
		Text:  txt.text,
		Font:  txt.font,
		Size:  txt.font_size,
		Color: txt.color,
	}
}
func (txt layout_text_node) is_block_level() bool { return false }

//==============================================================================
// Box
//==============================================================================

type layout_box_node interface {
	layout_node
	get_elem() dom_Element
	get_rect() gfx.Rect
	get_child_boxes() []layout_box_node
	get_child_texts() []*layout_text_node
	is_width_auto() bool
	is_height_auto() bool
	increment_size(width, height float64)
	increment_if_needed(width, height float64)
}
type layout_box_node_common struct {
	layout_node_common
	elem        dom_Element
	rect        gfx.Rect
	width_auto  bool
	height_auto bool
	child_boxes []layout_box_node
	child_texts []*layout_text_node
}

func (bx layout_box_node_common) get_elem() dom_Element { return bx.elem }
func (bx layout_box_node_common) get_rect() gfx.Rect    { return bx.rect }
func (bx layout_box_node_common) get_child_boxes() []layout_box_node {
	return bx.child_boxes
}
func (bx layout_box_node_common) get_child_texts() []*layout_text_node {
	return bx.child_texts
}
func (bx layout_box_node_common) is_width_auto() bool  { return bx.width_auto }
func (bx layout_box_node_common) is_height_auto() bool { return bx.height_auto }
func (bx *layout_box_node_common) increment_size(width, height float64) {
	if width == 0 && height == 0 {
		return
	}
	bx.rect.Width += width
	bx.rect.Height += height
	parent := bx.get_parent()
	if !cm.IsNil(parent) {
		if parent, ok := parent.(layout_box_node); ok {
			w := width
			h := height
			if !parent.is_width_auto() {
				w = 0
			}
			if !parent.is_height_auto() {
				h = 0
			}
			parent.increment_size(w, h)
		}
	}
}
func (bx *layout_box_node_common) increment_if_needed(min_width, min_height float64) {
	w_diff := max(min_width-bx.rect.Width, 0)
	h_diff := max(min_height-bx.rect.Height, 0)
	bx.increment_size(w_diff, h_diff)
}
func (bx layout_box_node_common) make_paint_node() gfx.PaintNode {
	paint_nodes := []gfx.PaintNode{}

	var color = css_color_transparent()
	if bx.elem != nil {
		color = bx.elem.get_computed_style_set().get_background_color()
	}
	rgba_color := color.to_rgba()

	for _, child := range bx.get_child_boxes() {
		paint_nodes = append(paint_nodes, child.make_paint_node())
	}
	for _, child := range bx.get_child_texts() {
		paint_nodes = append(paint_nodes, child.make_paint_node())
	}
	return gfx.BoxPaint{Items: paint_nodes, Color: rgba_color, Rect: bx.rect}
}

// https://www.w3.org/TR/css-display-3/#inline-box
type layout_inline_box_node struct {
	layout_box_node_common
	parent_bcon *layout_block_container_node
}

func (bx layout_inline_box_node) String() string {
	return fmt.Sprintf("inline-box %v at [%v]", bx.elem, bx.rect)
}
func (bx layout_inline_box_node) get_node_type() layout_node_type {
	return layout_node_type_inline_box
}
func (bx layout_inline_box_node) get_rect() gfx.Rect   { return bx.rect }
func (bx layout_inline_box_node) is_block_level() bool { return false }

// NOTE: This should *only* be called once after making layout node.
func (bx *layout_inline_box_node) init_children(
	tb layout_tree_builder,
	children []dom_Node,
	write_mode layout_write_mode,
) {
	if len(bx.child_boxes) != 0 || len(bx.child_texts) != 0 {
		panic("this method should be called only once")
	}
	for _, child_node := range children {
		nodes := tb.make_layout_for_node(bx.parent_bcon.ifc, bx.parent_bcon.bfc, bx.parent_bcon.ifc, write_mode, bx, child_node, false)
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if sub_bx, ok := node.(layout_box_node); ok {
				bx.child_boxes = append(bx.child_boxes, sub_bx)
			} else if txt, ok := node.(*layout_text_node); ok {
				bx.child_texts = append(bx.child_texts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}
	}
}

// https://www.w3.org/TR/css-display-3/#block-container
type layout_block_container_node struct {
	layout_box_node_common
	bfc         *layout_block_formatting_context
	ifc         *layout_inline_formatting_context
	parent_fctx layout_formatting_context
	created_bfc bool
	created_ifc bool
}

func (bcon layout_block_container_node) String() string {
	fc_str := ""
	if bcon.created_bfc {
		fc_str += "[BFC]"
	}
	if bcon.created_ifc {
		fc_str += "[IFC]"
	}
	return fmt.Sprintf("block-container [elem %v] at [%v] %s", bcon.elem, bcon.rect, fc_str)
}
func (bcon layout_block_container_node) get_node_type() layout_node_type {
	return layout_node_type_block_container
}
func (bcon layout_block_container_node) is_block_level() bool { return true }

// NOTE: This should *only* be called once after making layout node.
func (bcon *layout_block_container_node) init_children(
	tb layout_tree_builder,
	children []dom_Node,
	write_mode layout_write_mode,
) {
	if len(bcon.child_boxes) != 0 || len(bcon.child_texts) != 0 {
		panic("this method should be called only once")
	}

	// Check each children's display type - By running dry-run layout on each of them
	has_inline, has_block := false, false
	is_inline := make([]bool, len(children))
	for i, child_node := range children {
		nodes := tb.make_layout_for_node(bcon.parent_fctx, bcon.bfc, bcon.ifc, write_mode, bcon, child_node, true)
		is_inline[i] = false
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if node.is_block_level() {
				has_block = true
			} else {
				has_inline = true
				is_inline[i] = true
			}
		}
	}

	// If we have both inline and block-level, we need anonymous block container for inline nodes.
	// (This is actually part of CSS spec)
	need_anonymous_block_container := has_inline && has_block
	if has_inline && !has_block {
		bcon.ifc = &layout_inline_formatting_context{}
		bcon.ifc.creator_box = bcon
		bcon.ifc.bcon = bcon
		if bcon.bfc.is_dummy_context {
			bcon.ifc.is_dummy_context = true
		}
		bcon.created_ifc = true
	}

	// Now we can layout the children for real
	for i, child_node := range children {
		var nodes []layout_node
		if is_inline[i] && need_anonymous_block_container {
			// Create anonymous block container
			box_left, box_top := layout_calc_next_position(bcon.bfc, bcon.ifc, write_mode, false)
			box_rect := gfx.Rect{Left: box_left, Top: box_top, Width: bcon.rect.Width, Height: bcon.rect.Height}
			anon_bcon := tb.make_block_container(bcon.parent_fctx, bcon.ifc, bcon, nil, box_rect, false, false)
			anon_bcon.ifc = bcon.ifc
			anon_bcon.init_children(tb, []dom_Node{child_node}, write_mode)
			nodes = []layout_node{anon_bcon}
		} else {
			// Create layout node normally
			nodes = tb.make_layout_for_node(bcon.parent_fctx, bcon.bfc, bcon.ifc, write_mode, bcon, child_node, false)
		}
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if bx, ok := node.(layout_box_node); ok {
				bcon.child_boxes = append(bcon.child_boxes, bx)
			} else if txt, ok := node.(*layout_text_node); ok {
				bcon.child_texts = append(bcon.child_texts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}

	}
}

//==============================================================================
// The main layout code
//==============================================================================

type layout_tree_builder struct {
	font gfx.Font
}

func (tb layout_tree_builder) make_text(
	parent layout_box_node,
	text string,
	rect gfx.Rect,
	color color.RGBA,
	font_size float64,
) *layout_text_node {
	t := layout_text_node{}
	t.parent = parent
	t.text = text
	t.rect = rect
	t.font = tb.font
	t.color = color
	t.font_size = font_size
	return &t
}
func (tb layout_tree_builder) make_inline_box(
	parent_bcon *layout_block_container_node,
	elem dom_Element,
	rect gfx.Rect,
	width_auto, height_auto bool,
) *layout_inline_box_node {
	ibox := &layout_inline_box_node{}
	ibox.parent = parent_bcon
	ibox.elem = elem
	ibox.rect = rect
	ibox.width_auto = width_auto
	ibox.height_auto = height_auto
	ibox.parent_bcon = parent_bcon
	return ibox
}

// ifc may get overwritten during init_children()
func (tb layout_tree_builder) make_block_container(
	parent_fctx layout_formatting_context,
	ifc *layout_inline_formatting_context,
	parent layout_node,
	elem dom_Element,
	rect gfx.Rect,
	width_auto, height_auto bool,
) *layout_block_container_node {
	bcon := &layout_block_container_node{}
	bcon.parent = parent
	bcon.elem = elem
	bcon.rect = rect
	bcon.width_auto = width_auto
	bcon.height_auto = height_auto
	bcon.parent_fctx = parent_fctx
	bcon.ifc = ifc
	if cm.IsNil(parent_fctx) || parent_fctx.get_formatting_context_type() != layout_formatting_context_type_block {
		bcon.bfc = layout_make_bfc(bcon)
		bcon.created_bfc = true
	} else {
		bcon.bfc = parent_fctx.(*layout_block_formatting_context)
	}
	return bcon
}

func layout_calc_next_position(bfc *layout_block_formatting_context, ifc *layout_inline_formatting_context, write_mode layout_write_mode, is_inline bool) (float64, float64) {
	if is_inline {
		var inline_pos, block_pos float64
		if len(ifc.line_boxes) != 0 {
			inline_pos = ifc.get_current_natural_pos()
			block_pos = ifc.get_current_line_box().initial_block_pos
		} else {
			inline_pos = 0
			block_pos = bfc.get_current_natural_pos()
		}
		base_rect := bfc.get_creator_box().get_rect()
		if write_mode == layout_write_mode_vertical {
			return base_rect.Left + block_pos, base_rect.Top + inline_pos
		}
		return base_rect.Left + inline_pos, base_rect.Top + block_pos
	} else {
		var x, y float64
		if len(ifc.line_boxes) != 0 {
			x = ifc.get_current_natural_pos()
		} else {
			x = 0
		}
		y = bfc.get_current_natural_pos()
		base_rect := bfc.get_creator_box().get_rect()
		if write_mode == layout_write_mode_vertical {
			return base_rect.Left + y, base_rect.Top + x
		}
		return base_rect.Left + x, base_rect.Top + y
	}

}

// This function can be seen as heart of layout process.
//
// dry_run flag is intended for determine resulting box type. If dry_run is true:
//   - parent_fctx, parent_node will be internally replaced by dummy ones,
//     so that they don't affect actual parent context.
//   - will not layout its children, and so returned box will have empty children.
//   - New dummy formatting context will have its is_dummy_context set to true.
//     (As of writing this comment, this is mostly for debug prints.
//     outputs with dummy contexts can be confusing when mixed with real ones)
func (tb layout_tree_builder) make_layout_for_node(
	parent_fctx layout_formatting_context,
	bfc *layout_block_formatting_context,
	ifc *layout_inline_formatting_context,
	write_mode layout_write_mode,
	parent_node layout_box_node,
	dom_node dom_Node,
	dry_run bool,
) []layout_node {
	var parent_elem dom_Element
	{
		curr_node := parent_node
		for curr_node.get_elem() == nil {
			parent := curr_node.get_parent()
			if parent == nil {
				break
			}
			curr_node = parent.(layout_box_node)
		}
		parent_elem = curr_node.get_elem()
	}

	if dry_run {
		dummy_bcon := &layout_block_container_node{}
		dummy_bcon.elem = parent_node.get_elem()
		dummy_bcon.bfc = &layout_block_formatting_context{
			layout_formatting_context_common: layout_formatting_context_common{
				is_dummy_context: true,
				creator_box:      dummy_bcon,
			},
		}
		bfc = dummy_bcon.bfc
		ifc = &layout_inline_formatting_context{
			layout_formatting_context_common: layout_formatting_context_common{
				is_dummy_context: true,
				creator_box:      dummy_bcon,
			},
			bcon: dummy_bcon,
		}
		if parent_fctx.get_formatting_context_type() == layout_formatting_context_type_block {
			parent_fctx = bfc
		} else {
			parent_fctx = ifc
		}
		parent_node = dummy_bcon
	}
	if bfc == nil {
		panic("FFC should not be nil at this point")
	}
	if ifc == nil {
		panic("IFC should not be nil at this point")
	}

	if _, ok := dom_node.(*dom_Comment_s); ok {
		//======================================================================
		// Layout for Comment nodes
		//======================================================================

		// If you can see comments on screen without devtools, congraturations!
		// You are a very rare person with built-in devtools inside your brain.
		return nil
	}
	if text, ok := dom_node.(*dom_Text_s); ok {
		parent_css := parent_elem.get_computed_style_set()

		//======================================================================
		// Layout for Text nodes
		//======================================================================
		var text_node *layout_text_node
		str := text.get_text()

		// Apply text-transform
		if v := parent_css.get_text_transform(); !cm.IsNil(v) {
			str = v.apply(str)
		}

		if str == "" {
			// Nothing to display
			return nil
		}
		str = strings.TrimSpace(str)
		if str == "" {
			str = " "
		}

		// Create line box if needed
		if len(ifc.line_boxes) == 0 {
			ifc.add_line_box(bfc)
		}

		// Calculate the font size
		parent_font_size := css_number_from_float(css_preferred_font_size) // STUB
		font_size := parent_css.get_font_size().calculate_real_font_size(parent_font_size).length_to_px(parent_font_size)
		tb.font.SetTextSize(int(font_size)) // NOTE: Size we set here will only be used for measuring
		metrics := tb.font.Metrics()

		fragment_remaining := str
		text_nodes := []layout_node{}

		for 0 < len(fragment_remaining) {
			line_box := ifc.get_current_line_box()

			var rect gfx.Rect
			var inline_axis_size float64
			str_len := len(fragment_remaining)

			// Calculate left/top position
			left, top := layout_calc_next_position(bfc, ifc, write_mode, true)

			// Figure out where we should end current fragment, so that we don't overflow the line box.
			// TODO: We should not do this if we are not doing text wrapping(e.g. whitespace: nowrap).
			for {
				// FIXME: This is very brute-force way of fragmenting text.
				//        We need smarter way to handle this.

				// Calculate width/height using dimensions of the text
				width, _ := gfx.MeasureText(tb.font, fragment_remaining[:str_len])

				rect = gfx.Rect{Left: left, Top: top, Width: width, Height: metrics.LineHeight}

				// Check if parent's size is auto and we have to grow its size.
				inline_axis_size = rect.Width
				if write_mode == layout_write_mode_vertical {
					inline_axis_size = rect.Height
				}
				// Check if we overflow beyond available size
				if ifc.is_dummy_context || line_box.current_natural_pos+inline_axis_size <= line_box.available_width {
					// If not, we don't have to fragment text further.
					break
				}
				str_len-- // Decrement length and try again
			}
			fragment := fragment_remaining[:str_len]
			fragment_remaining = fragment_remaining[str_len:]

			// Make text node
			color := parent_css.get_color().to_rgba()
			text_node = tb.make_text(parent_node, fragment, rect, color, font_size)

			if parent_node.is_width_auto() {
				parent_node.increment_size(rect.Width, 0)
				// parent_node.get_rect().Width += rect.Width
			}

			line_box.current_line_height = max(line_box.current_line_height, rect.Height)
			if parent_node.is_height_auto() {
				// Increment parent's height if needed.
				parent_node.increment_if_needed(0, line_box.current_line_height)
			}
			ifc.increment_natural_pos(inline_axis_size)
			text_nodes = append(text_nodes, text_node)
			if len(fragment_remaining) != 0 {
				// Make a new line
				ifc.add_line_box(bfc)
			}
		}

		return text_nodes
	} else if elem, ok := dom_node.(dom_Element); ok {
		//======================================================================
		// Layout for Element nodes
		//======================================================================
		css := elem.get_computed_style_set()
		compute_box_rect := func(is_inline bool) (r gfx.Rect, width_auto, height_auto bool) {
			// Calculate left/top position
			box_left, box_top := layout_calc_next_position(bfc, ifc, write_mode, is_inline)

			// Calculate width/height using `width` and `height` property
			box_width := css.get_width()
			box_height := css.get_height()
			box_width_px := 0.0
			box_height_px := 0.0

			// If width or height is auto, we start from 0 and expand it as we layout the children.
			if box_width.tp != css_size_value_type_auto {
				parent_size := css_number_from_float(parent_node.get_rect().Width)
				box_width_px = box_width.compute_used_value(parent_size).length_to_px(parent_size)
			} else {
				width_auto = true
			}
			if box_height.tp != css_size_value_type_auto {
				parent_size := css_number_from_float(parent_node.get_rect().Height)
				box_height_px = box_height.compute_used_value(parent_size).length_to_px(parent_size)
			} else {
				height_auto = true
			}
			return gfx.Rect{Left: box_left, Top: box_top, Width: box_width_px, Height: box_height_px}, width_auto, height_auto
		}

		css_display := css.get_display()
		switch css_display.mode {
		case css_display_mode_none:
			return nil
		case css_display_mode_outer_inner_mode:
			is_inline := css_display.outer_mode == css_display_outer_mode_inline
			box_rect, width_auto, height_auto := compute_box_rect(is_inline)

			// Check if we have auto size on a block element. If so, use parent's size and unset auto.
			if css_display.outer_mode == css_display_outer_mode_block {
				if write_mode == layout_write_mode_horizontal && width_auto {
					box_rect.Width = parent_node.get_rect().Width
					width_auto = false
				} else if write_mode == layout_write_mode_vertical && height_auto {
					box_rect.Height = parent_node.get_rect().Height
					height_auto = false
				}
			}

			// Increment natural position(if it's auto)
			// XXX: Should we increment width/height if the element uses absolute positioning?
			switch css_display.outer_mode {
			case css_display_outer_mode_block:
				if parent_node.is_width_auto() {
					parent_node.increment_if_needed(box_rect.Width, 0)
				}
				if parent_node.is_height_auto() {
					parent_node.increment_size(0, box_rect.Height)
				}
			case css_display_outer_mode_inline:
				if parent_node.is_width_auto() {
					parent_node.increment_size(box_rect.Width, 0)
					if len(ifc.line_boxes) == 0 {
						ifc.add_line_box(bfc)
					}
				}
				if parent_node.is_height_auto() {
					if len(ifc.line_boxes) == 0 {
						ifc.add_line_box(bfc)
					}
					line_box := ifc.get_current_line_box()
					line_box.current_line_height = max(line_box.current_line_height, box_rect.Height)

					// Increment parent's height if needed.
					if parent_node.get_rect().Height < line_box.current_line_height {
						parent_node.increment_if_needed(0, line_box.current_line_height)
					}
				}
			}

			switch css_display.inner_mode {
			case css_display_inner_mode_flow:
				//==================================================================
				// "flow" mode (block, inline, run-in, list-item, inline list-item display modes)
				//==================================================================

				// https://www.w3.org/TR/css-display-3/#valdef-display-flow

				should_make_inline_box := false
				if css_display.outer_mode == css_display_outer_mode_inline ||
					css_display.outer_mode == css_display_outer_mode_run_in {
					if parent_fctx.get_formatting_context_type() == layout_formatting_context_type_block ||
						parent_fctx.get_formatting_context_type() == layout_formatting_context_type_inline {
						should_make_inline_box = true
					}
				}
				var box layout_box_node
				if should_make_inline_box {
					parent_bcon := parent_node
					for {
						if _, ok := parent_bcon.(*layout_block_container_node); ok {
							break
						}
						parent_bcon = parent_bcon.get_parent().(layout_box_node)
					}
					ibox := tb.make_inline_box(parent_bcon.(*layout_block_container_node), elem, box_rect, width_auto, height_auto)
					if !dry_run {
						ibox.init_children(tb, elem.get_children(), write_mode)
					}
					box = ibox
				} else {
					bcon := tb.make_block_container(parent_fctx, ifc, parent_node, elem, box_rect, width_auto, height_auto)
					if !dry_run {
						bcon.init_children(tb, elem.get_children(), write_mode)
					}
					box = bcon
				}

				// Increment natural position
				inline_axis_size := box.get_rect().Width
				block_axis_size := box.get_rect().Height
				if write_mode == layout_write_mode_vertical {
					inline_axis_size, block_axis_size = block_axis_size, inline_axis_size
				}
				switch css_display.outer_mode {
				case css_display_outer_mode_inline:
					_ = inline_axis_size
					// ifc.increment_natural_pos(inline_axis_size)
				case css_display_outer_mode_block:
					bfc.increment_natural_pos(block_axis_size)
				}

				return []layout_node{box}
			case css_display_inner_mode_flow_root:
				//==================================================================
				// "flow-root" mode (flow-root, inline-block display modes)
				//==================================================================
				// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
				bcon := tb.make_block_container(parent_fctx, ifc, parent_node, elem, box_rect, width_auto, height_auto)
				if !dry_run {
					bcon.init_children(tb, elem.get_children(), write_mode)
				}

				// Increment natural position
				inline_axis_size := bcon.get_rect().Width
				block_axis_size := bcon.get_rect().Height
				if write_mode == layout_write_mode_vertical {
					inline_axis_size, block_axis_size = block_axis_size, inline_axis_size
				}
				switch css_display.outer_mode {
				case css_display_outer_mode_inline:
					ifc.increment_natural_pos(inline_axis_size)
				case css_display_outer_mode_block:
					bfc.increment_natural_pos(block_axis_size)
				}

				return []layout_node{bcon}
			default:
				log.Panicf("TODO: Support display: %v", css_display)
			}

		default:
			log.Panicf("TODO: Support display: %v", css_display)
		}
	}

	panic("unreachable")
}

// https://www.w3.org/TR/css-display-3/#initial-containing-block
func make_layout(root dom_Element, viewport_width, viewport_height float64, plat platform.Platform) layout_node {
	tb := layout_tree_builder{}
	tb.font = plat.OpenFont("this_is_not_real_filename.ttf")
	tb.font.SetTextSize(32)
	bfc := &layout_block_formatting_context{}
	ifc := &layout_inline_formatting_context{}
	box_rect := gfx.Rect{Left: 0, Top: 0, Width: viewport_width, Height: viewport_height}
	icb := tb.make_block_container(bfc, ifc, nil, nil, box_rect, false, false)
	bfc.creator_box = icb
	ifc.creator_box = icb
	ifc.bcon = icb
	icb.init_children(tb, []dom_Node{root}, layout_write_mode_horizontal)
	return icb
}

func browser_print_layout_tree(node layout_node) {
	curr_node := node
	count := 0
	if !cm.IsNil(curr_node.get_parent()) {
		for n := curr_node.get_parent(); !cm.IsNil(n); n = n.get_parent() {
			count += 4
		}
	}
	indent := strings.Repeat(" ", count)
	fmt.Printf("%s%v\n", indent, node)
	if bx, ok := curr_node.(layout_box_node); ok {
		for _, child := range bx.get_child_boxes() {
			browser_print_layout_tree(child)
		}
		for _, child := range bx.get_child_texts() {
			browser_print_layout_tree(child)
		}
	}
}
