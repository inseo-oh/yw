package libhtml

import (
	"fmt"
	"log"
	"strings"
	cm "yw/libcommon"
	"yw/libgfx"
	"yw/libplatform"
)

type browser_layout_node interface {
	get_node_type() browser_layout_node_type
	get_parent() browser_layout_node
	make_paint_node() browser_paint_node
	String() string
}
type browser_layout_node_common struct {
	parent browser_layout_node
}

type browser_layout_node_type uint8

const (
	browser_layout_node_type_inline_box = browser_layout_node_type(iota)
	browser_layout_node_type_block_container
	browser_layout_node_type_initial_containing_block
	browser_layout_node_type_text
)

func (n browser_layout_node_common) get_parent() browser_layout_node {
	return n.parent
}

//==============================================================================
// Formatting contexts
//==============================================================================

type browser_layout_formatting_context interface {
	get_formatting_context_type() browser_layout_formatting_context_type
	get_current_natural_pos() float64
	set_current_natural_pos(pos float64)
	increment_natural_position(inc float64)
}

type browser_layout_formatting_context_common struct {
	// Current natural position.
	current_natural_pos float64
}
type browser_layout_formatting_context_type uint8

const (
	browser_layout_formatting_context_type_block = browser_layout_formatting_context_type(iota)
	browser_layout_formatting_context_type_inline
)

type browser_layout_write_mode uint8

const (
	browser_layout_write_mode_horizontal = browser_layout_write_mode(iota)
	browser_layout_write_mode_vertical
)

func (fctx browser_layout_formatting_context_common) get_current_natural_pos() float64 {
	return fctx.current_natural_pos
}
func (fctx *browser_layout_formatting_context_common) set_current_natural_pos(pos float64) {
	fctx.current_natural_pos = pos
}
func (fctx *browser_layout_formatting_context_common) increment_natural_position(pos float64) {
	fctx.current_natural_pos += pos
}

// Block Formatting Contexts(BFC for short) are responsible for tracking Y-axis,
// or more accurately, the opposite axis of writing mode.
// (English uses X-axis for writing text, so BFC's position grows Y-axis)
//
// https://www.w3.org/TR/CSS2/visuren.html#block-formatting
type browser_layout_block_formatting_context struct {
	browser_layout_formatting_context_common
	current_position float64
}

func (fctx browser_layout_block_formatting_context) get_formatting_context_type() browser_layout_formatting_context_type {
	return browser_layout_formatting_context_type_block
}

func browser_layout_make_bfc() *browser_layout_block_formatting_context {
	bfc := browser_layout_block_formatting_context{}
	return &bfc
}

// Inline Formatting Contexts(IFC for short) are responsible for tracking X-axis,
// or more accurately, the primary axis of writing mode.
// (English uses X-axis for writing text, so IFC's position grows X-axis)
//
// or can be also thought as "The opposite axis of BFC", if you really want :D
//
// https://www.w3.org/TR/CSS2/visuren.html#inline-formatting
type browser_layout_inline_formatting_context struct {
	browser_layout_formatting_context_common
	current_position float64
}

func (fctx browser_layout_inline_formatting_context) get_formatting_context_type() browser_layout_formatting_context_type {
	return browser_layout_formatting_context_type_inline
}

//==============================================================================
// Text
//==============================================================================

type browser_layout_text_node struct {
	browser_layout_node_common
	rect libgfx.Rect
	text dom_Text
	font libplatform.Font
}

func (txt browser_layout_text_node) String() string {
	return fmt.Sprintf("text %s at [%v]", txt.text, txt.rect)
}
func (txt browser_layout_text_node) get_node_type() browser_layout_node_type {
	return browser_layout_node_type_text
}
func (txt browser_layout_text_node) make_paint_node() browser_paint_node {
	return browser_text_paint_node{
		text_layout_node: txt,
		font:             txt.font,
	}
}

//==============================================================================
// Box
//==============================================================================

type browser_layout_box_node interface {
	browser_layout_node
	get_rect() libgfx.Rect
	get_rect_p() *libgfx.Rect
	get_child_boxes() []browser_layout_box_node
	get_child_texts() []browser_layout_text_node
	is_width_auto() bool
	is_height_auto() bool
}
type browser_layout_box_node_common struct {
	browser_layout_node_common
	elem        dom_Element
	rect        libgfx.Rect
	width_auto  bool
	height_auto bool
	child_boxes []browser_layout_box_node
	child_texts []browser_layout_text_node
}

func (bx browser_layout_box_node_common) get_rect() libgfx.Rect     { return bx.rect }
func (bx *browser_layout_box_node_common) get_rect_p() *libgfx.Rect { return &bx.rect }
func (bx browser_layout_box_node_common) get_child_boxes() []browser_layout_box_node {
	return bx.child_boxes
}
func (bx browser_layout_box_node_common) get_child_texts() []browser_layout_text_node {
	return bx.child_texts
}
func (bx *browser_layout_box_node_common) set_child_texts(texts []browser_layout_text_node) {
	bx.child_texts = texts
}
func (bx browser_layout_box_node_common) is_width_auto() bool  { return bx.width_auto }
func (bx browser_layout_box_node_common) is_height_auto() bool { return bx.height_auto }
func (bx browser_layout_box_node_common) make_paint_node() browser_paint_node {
	paintables := []browser_paint_node{}
	for _, child := range bx.get_child_boxes() {
		paintables = append(paintables, child.make_paint_node())
	}
	for _, child := range bx.get_child_texts() {
		paintables = append(paintables, child.make_paint_node())
	}
	return browser_grouped_paint_node{items: paintables}
}

// https://www.w3.org/TR/css-display-3/#inline-box
type browser_layout_inline_box_node struct {
	browser_layout_box_node_common
}

func (bx browser_layout_inline_box_node) String() string {
	return fmt.Sprintf("inline-box %v at [%v]", bx.elem, bx.rect)
}
func (bx browser_layout_inline_box_node) get_node_type() browser_layout_node_type {
	return browser_layout_node_type_inline_box
}
func (bx browser_layout_inline_box_node) get_rect() libgfx.Rect { return bx.rect }

// https://www.w3.org/TR/css-display-3/#block-container
type browser_layout_block_container_node struct {
	browser_layout_box_node_common
	bfc         *browser_layout_block_formatting_context
	ifc         *browser_layout_inline_formatting_context
	created_bfc bool
}

func (bcon browser_layout_block_container_node) String() string {
	bfc_str := ""
	if bcon.created_bfc {
		bfc_str = "[BFC]"
	}
	return fmt.Sprintf("block-container %v at [%v] %s", bcon.elem, bcon.rect, bfc_str)
}
func (bcon browser_layout_block_container_node) get_node_type() browser_layout_node_type {
	return browser_layout_node_type_block_container
}
func (bcon *browser_layout_block_container_node) get_or_make_ifc() *browser_layout_inline_formatting_context {
	if bcon.ifc == nil {
		bcon.ifc = &browser_layout_inline_formatting_context{}
	}
	return bcon.ifc
}

//==============================================================================
// The main layout code
//==============================================================================

type browser_layout_tree_builder struct {
	font libplatform.Font
}

func (tb browser_layout_tree_builder) make_text(
	parent browser_layout_box_node,
	text dom_Text,
	rect libgfx.Rect,
) browser_layout_text_node {
	t := browser_layout_text_node{}
	t.parent = parent
	t.text = text
	t.rect = rect
	t.font = tb.font
	return t
}
func (tb browser_layout_tree_builder) make_inline_box(
	parent_fctx *browser_layout_inline_formatting_context,
	write_mode browser_layout_write_mode,
	parent browser_layout_box_node,
	elem dom_Element,
	rect libgfx.Rect,
	width_auto, height_auto bool,
) *browser_layout_inline_box_node {
	ibox := browser_layout_inline_box_node{}
	ibox.parent = parent
	ibox.elem = elem
	ibox.rect = rect
	ibox.width_auto = width_auto
	ibox.height_auto = height_auto
	for _, child_node := range elem.get_children() {
		node := tb.make_layout_for_node(parent_fctx, write_mode, &ibox, child_node)
		if cm.IsNil(node) {
			continue
		}
		if bx, ok := node.(browser_layout_box_node); ok {
			ibox.child_boxes = append(ibox.child_boxes, bx)
		} else if bx, ok := node.(browser_layout_text_node); ok {
			ibox.child_texts = append(ibox.child_texts, bx)
		} else {
			log.Panicf("unknown node result %v", node)
		}
	}
	return &ibox
}
func (tb browser_layout_tree_builder) make_block_container(
	parent_fctx browser_layout_formatting_context,
	write_mode browser_layout_write_mode,
	parent browser_layout_node,
	elem dom_Element,
	rect libgfx.Rect,
	width_auto, height_auto bool,
) *browser_layout_block_container_node {
	bcon := browser_layout_block_container_node{}
	bcon.parent = parent
	bcon.elem = elem
	bcon.rect = rect
	bcon.width_auto = width_auto
	bcon.height_auto = height_auto
	if parent_fctx.get_formatting_context_type() != browser_layout_formatting_context_type_block {
		bcon.bfc = browser_layout_make_bfc()
		bcon.created_bfc = true
	} else {
		bcon.bfc = parent_fctx.(*browser_layout_block_formatting_context)
	}
	for _, child_node := range elem.get_children() {
		node := tb.make_layout_for_node(parent_fctx, write_mode, &bcon, child_node)
		if cm.IsNil(node) {
			continue
		}
		if bx, ok := node.(browser_layout_box_node); ok {
			bcon.child_boxes = append(bcon.child_boxes, bx)
		} else if bx, ok := node.(browser_layout_text_node); ok {
			bcon.child_texts = append(bcon.child_texts, bx)
		} else {
			log.Panicf("unknown node result %v", node)
		}
	}
	return &bcon
}

func (tb browser_layout_tree_builder) make_layout_for_node(
	parent_fctx browser_layout_formatting_context,
	write_mode browser_layout_write_mode,
	parent browser_layout_box_node,
	dom_node dom_Node,
) browser_layout_node {
	closest_block_container := func() *browser_layout_block_container_node {
		curr := parent
		for {
			if curr.get_node_type() == browser_layout_node_type_block_container {
				return curr.(*browser_layout_block_container_node)
			}
			parent := curr.get_parent()
			if cm.IsNil(parent) {
				break
			}
			new_curr, ok := parent.(browser_layout_box_node)
			if !ok {
				break
			}
			curr = new_curr
		}
		panic("could not find any block container")
	}
	get_closest_bfc := func() *browser_layout_block_formatting_context {
		if parent_fctx.get_formatting_context_type() == browser_layout_formatting_context_type_block {
			return parent_fctx.(*browser_layout_block_formatting_context)
		}
		closest := closest_block_container()
		return closest.bfc
	}
	get_closest_ifc := func() *browser_layout_inline_formatting_context {
		if parent_fctx.get_formatting_context_type() == browser_layout_formatting_context_type_inline {
			return parent_fctx.(*browser_layout_inline_formatting_context)
		}
		closest := closest_block_container()
		return closest.get_or_make_ifc()
	}
	get_next_position := func() (float64, float64) {
		x, y := get_closest_ifc().current_position, get_closest_bfc().current_position
		if write_mode == browser_layout_write_mode_vertical {
			return y, x
		}
		return x, y
	}
	increment_next_box_position := func(x, y float64) {
		if write_mode == browser_layout_write_mode_vertical {
			x, y = y, x
		}
		get_closest_ifc().current_position += x
		get_closest_bfc().current_position += y
	}
	_ = increment_next_box_position // FIXME: Remove this function if we don't need it.

	if _, ok := dom_node.(*dom_Comment_s); ok {
		//======================================================================
		// Layout for Comment nodes
		//======================================================================

		// If you can see comments on screen without devtools, congraturations!
		// You are a very rare person with built-in devtools inside your brain.
		return nil
	}
	if text, ok := dom_node.(*dom_Text_s); ok {
		//======================================================================
		// Layout for Text nodes
		//======================================================================
		left, top := get_next_position()
		width, height := libplatform.MeasureText(tb.font, text.get_text())
		if parent.is_width_auto() {
			parent.get_rect_p().Width += width
		}
		if parent.is_height_auto() {
			parent.get_rect_p().Height += height
		}
		rect := libgfx.Rect{Left: left, Top: top, Width: width, Height: height}
		text := tb.make_text(parent, text, rect)
		inline_axis_size := rect.Width
		if write_mode == browser_layout_write_mode_vertical {
			inline_axis_size = rect.Height
		}
		get_closest_ifc().increment_natural_position(inline_axis_size)
		return text
	} else if elem, ok := dom_node.(dom_Element); ok {
		//======================================================================
		// Layout for Element nodes
		//======================================================================
		css := elem.get_computed_style_set()
		compute_box_rect := func() (r libgfx.Rect, width_auto, height_auto bool) {
			box_left, box_top := get_next_position()
			box_left += parent.get_rect().Left
			box_top += parent.get_rect().Top
			box_width := css.get_width()
			box_height := css.get_height()
			box_width_px := 0.0
			box_height_px := 0.0
			// If width or height is auto, we start from 0 and expand it as we layout the children.
			if box_width.tp != css_size_value_type_auto {
				box_width_px = box_width.compute_used_value(css_number_from_float(parent.get_rect().Width)).length_to_px()
			} else {
				width_auto = true
			}
			if box_height.tp != css_size_value_type_auto {
				box_height_px = box_height.compute_used_value(css_number_from_float(parent.get_rect().Height)).length_to_px()
			} else {
				height_auto = true
			}
			return libgfx.Rect{Left: box_left, Top: box_top, Width: box_width_px, Height: box_height_px}, width_auto, height_auto
		}

		css_display := css.get_display()
		switch css_display.mode {
		case css_display_mode_none:
			return nil
		case css_display_mode_outer_inner_mode:
			box_rect, width_auto, height_auto := compute_box_rect()

			// Check if we have auto size on a block element. If so, use parent's size and unset auto.
			if css_display.outer_mode == css_display_outer_mode_block {
				if write_mode == browser_layout_write_mode_horizontal && width_auto {
					box_rect.Width = parent.get_rect().Width
					width_auto = false
				} else if write_mode == browser_layout_write_mode_vertical && height_auto {
					box_rect.Height = parent.get_rect().Height
					height_auto = false
				}
			}

			// Check if parent is auto and we have to grow its size.
			// XXX: Should we increment width/height if the element uses absolute positioning?
			if parent.is_width_auto() {
				parent.get_rect_p().Width += box_rect.Width
			}
			if parent.is_height_auto() {
				parent.get_rect_p().Height += box_rect.Height
			}
			inline_axis_size := box_rect.Width
			block_axis_size := box_rect.Height
			if write_mode == browser_layout_write_mode_vertical {
				inline_axis_size, block_axis_size = block_axis_size, inline_axis_size
			}
			switch css_display.outer_mode {
			case css_display_outer_mode_inline:
				get_closest_ifc().increment_natural_position(inline_axis_size)
			case css_display_outer_mode_block:
				get_closest_bfc().increment_natural_position(block_axis_size)
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
					if parent_fctx.get_formatting_context_type() == browser_layout_formatting_context_type_block ||
						parent_fctx.get_formatting_context_type() == browser_layout_formatting_context_type_inline {
						should_make_inline_box = true
					}
				}
				var box browser_layout_box_node
				if should_make_inline_box {
					box = tb.make_inline_box(get_closest_ifc(), write_mode, parent, elem, box_rect, width_auto, height_auto)
				} else {
					box = tb.make_block_container(parent_fctx, write_mode, parent, elem, box_rect, width_auto, height_auto)
				}
				return box
			case css_display_inner_mode_flow_root:
				//==================================================================
				// "flow-root" mode (flow-root, inline-block display modes)
				//==================================================================
				// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
				box := tb.make_block_container(parent_fctx, write_mode, parent, elem, box_rect, width_auto, height_auto)
				return box
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
func browser_make_layout(root dom_Element, viewport_width, viewport_height float64, plat libplatform.Platform) browser_layout_node {
	tb := browser_layout_tree_builder{}
	tb.font = plat.OpenFont("this_is_not_real_filename.ttf")
	tb.font.SetTextSize(32)
	bfc := browser_layout_block_formatting_context{}
	box_rect := libgfx.Rect{Left: 0, Top: 0, Width: viewport_width, Height: viewport_height}
	icb := tb.make_block_container(&bfc, browser_layout_write_mode_horizontal, nil, root, box_rect, false, false)
	return icb
}

func browser_print_layout_tree(node browser_layout_node) {
	curr_node := node
	count := 0
	if !cm.IsNil(curr_node.get_parent()) {
		for n := curr_node.get_parent(); !cm.IsNil(n); n = n.get_parent() {
			count += 4
		}
	}
	indent := strings.Repeat(" ", count)
	fmt.Printf("%s%v\n", indent, node)
	if bx, ok := curr_node.(browser_layout_box_node); ok {
		for _, child := range bx.get_child_boxes() {
			browser_print_layout_tree(child)
		}
		for _, child := range bx.get_child_texts() {
			browser_print_layout_tree(child)
		}
	}
}
