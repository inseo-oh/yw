// Implementation of the CSS Display Module Level 3 (https://www.w3.org/TR/css-display-3/)
package libhtml

import (
	"fmt"
	cm "yw/util"
)

type css_display_outer_mode uint8

const (
	css_display_outer_mode_block = css_display_outer_mode(iota)
	css_display_outer_mode_inline
	css_display_outer_mode_run_in
)

func (m css_display_outer_mode) String() string {
	switch m {
	case css_display_outer_mode_block:
		return "block"
	case css_display_outer_mode_inline:
		return "inline"
	case css_display_outer_mode_run_in:
		return "run-in"
	}
	return fmt.Sprintf("unregognized css_display_outer_mode %d", m)
}

type css_display_inner_mode uint8

const (
	css_display_inner_mode_flow = css_display_inner_mode(iota)
	css_display_inner_mode_flow_root
	css_display_inner_mode_table
	css_display_inner_mode_flex
	css_display_inner_mode_grid
	css_display_inner_mode_ruby
)

func (m css_display_inner_mode) String() string {
	switch m {
	case css_display_inner_mode_flow:
		return "flow"
	case css_display_inner_mode_flow_root:
		return "flow-root"
	case css_display_inner_mode_table:
		return "table"
	case css_display_inner_mode_flex:
		return "flex"
	case css_display_inner_mode_grid:
		return "grid"
	case css_display_inner_mode_ruby:
		return "ruby"
	}
	return fmt.Sprintf("unregognized css_display_inner_mode %d", m)
}

type css_display struct {
	mode css_display_mode

	outer_mode css_display_outer_mode // Only valid when mode is css_display_mode_outer_mode_inner_mode
	inner_mode css_display_inner_mode // Only valid when mode is css_display_mode_outer_mode_inner_mode
}

func (d css_display) String() string {
	switch d.mode {
	case css_display_mode_outer_inner_mode:
		return fmt.Sprintf("%v %v", d.outer_mode, d.inner_mode)
	case css_display_mode_table_row_group:
		return "table-row-group"
	case css_display_mode_table_header_group:
		return "table-header-group"
	case css_display_mode_table_footer_group:
		return "table-footer-group"
	case css_display_mode_table_row:
		return "table-row"
	case css_display_mode_table_cell:
		return "table-cell"
	case css_display_mode_table_column_group:
		return "table-column-group"
	case css_display_mode_table_column:
		return "table-column"
	case css_display_mode_table_caption:
		return "table-caption"
	case css_display_mode_ruby_base:
		return "ruby-base"
	case css_display_mode_ruby_text:
		return "ruby-text"
	case css_display_mode_ruby_base_container:
		return "ruby-base-container"
	case css_display_mode_ruby_text_container:
		return "ruby-text-container"
	case css_display_mode_contents:
		return "contents"
	case css_display_mode_none:
		return "none"
	}
	return fmt.Sprintf("unregognized css_display_mode %d", d)
}

type css_display_mode uint8

const (
	// Display mode determined by `outer_mode` and `inner_mode` field.
	css_display_mode_outer_inner_mode = css_display_mode(iota)

	// https://www.w3.org/TR/css-display-3/#typedef-display-internal

	css_display_mode_table_row_group
	css_display_mode_table_header_group
	css_display_mode_table_footer_group
	css_display_mode_table_row
	css_display_mode_table_cell
	css_display_mode_table_column_group
	css_display_mode_table_column
	css_display_mode_table_caption
	css_display_mode_ruby_base
	css_display_mode_ruby_text
	css_display_mode_ruby_base_container
	css_display_mode_ruby_text_container

	// https://www.w3.org/TR/css-display-3/#typedef-display-box

	css_display_mode_contents
	css_display_mode_none
)

type css_visibility uint8

const (
	css_visibility_visible = css_visibility(iota)
	css_visibility_hidden
	css_visibility_collapse
)

func (m css_visibility) String() string {
	switch m {
	case css_visibility_visible:
		return "visible"
	case css_visibility_hidden:
		return "hidden"
	case css_visibility_collapse:
		return "collapse"
	}
	return fmt.Sprintf("unregognized css_visibility %d", m)
}

// https://www.w3.org/TR/css-display-3/#typedef-display-outside
func (ts *css_token_stream) parse_display_outside() (css_display_outer_mode, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("block")) {
		return css_display_outer_mode_block, true
	} else if !cm.IsNil(ts.consume_ident_token_with("inline")) {
		return css_display_outer_mode_inline, true
	} else if !cm.IsNil(ts.consume_ident_token_with("run-in")) {
		return css_display_outer_mode_run_in, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-display-3/#typedef-display-inside
func (ts *css_token_stream) parse_display_inside() (css_display_inner_mode, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("flow")) {
		return css_display_inner_mode_flow, true
	} else if !cm.IsNil(ts.consume_ident_token_with("flow-root")) {
		return css_display_inner_mode_flow_root, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table")) {
		return css_display_inner_mode_table, true
	} else if !cm.IsNil(ts.consume_ident_token_with("flex")) {
		return css_display_inner_mode_flex, true
	} else if !cm.IsNil(ts.consume_ident_token_with("grid")) {
		return css_display_inner_mode_grid, true
	} else if !cm.IsNil(ts.consume_ident_token_with("ruby")) {
		return css_display_inner_mode_ruby, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-display-3/#propdef-display
func (ts *css_token_stream) parse_display() (css_display, bool) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if !cm.IsNil(ts.consume_ident_token_with("inline-block")) {
		return css_display{css_display_mode_outer_inner_mode, css_display_outer_mode_inline, css_display_inner_mode_flow_root}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("inline-table")) {
		return css_display{css_display_mode_outer_inner_mode, css_display_outer_mode_inline, css_display_inner_mode_table}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("inline-flex")) {
		return css_display{css_display_mode_outer_inner_mode, css_display_outer_mode_inline, css_display_inner_mode_flex}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("inline-grid")) {
		return css_display{css_display_mode_outer_inner_mode, css_display_outer_mode_inline, css_display_inner_mode_grid}, true
	}
	// Try <display-outside> <display-inside> ------------------------------
	got_outer_mode, got_inner_mode := false, false
	var outer_mode css_display_outer_mode
	var inner_mode css_display_inner_mode
	for !got_outer_mode || !got_inner_mode {
		got_something := false
		if !got_outer_mode {
			outer_mode, got_outer_mode = ts.parse_display_outside()
			if got_outer_mode {
				got_something = true
			}
		}
		if !got_inner_mode {
			inner_mode, got_inner_mode = ts.parse_display_inside()
			if got_inner_mode {
				got_something = true
			}
		}
		if !got_something {
			break
		}
	}
	if got_outer_mode || got_inner_mode {
		if !got_inner_mode {
			inner_mode = css_display_inner_mode_flow
		} else if !got_outer_mode {
			if inner_mode == css_display_inner_mode_ruby {
				outer_mode = css_display_outer_mode_inline
			} else {
				outer_mode = css_display_outer_mode_block
			}
		}
		return css_display{css_display_mode_outer_inner_mode, outer_mode, inner_mode}, true
	}
	// Try display-listitem ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-listitem
	// TODO

	// Try display-internal ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-internal

	if !cm.IsNil(ts.consume_ident_token_with("table-row-group")) {
		return css_display{css_display_mode_table_row_group, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-header-group")) {
		return css_display{css_display_mode_table_header_group, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-footer-group")) {
		return css_display{css_display_mode_table_footer_group, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-row")) {
		return css_display{css_display_mode_table_row, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-cell")) {
		return css_display{css_display_mode_table_cell, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-column-group")) {
		return css_display{css_display_mode_table_column_group, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-column")) {
		return css_display{css_display_mode_table_column, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("table-caption")) {
		return css_display{css_display_mode_table_caption, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("ruby-base")) {
		return css_display{css_display_mode_ruby_base, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("ruby-text")) {
		return css_display{css_display_mode_ruby_text, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("ruby-base-container")) {
		return css_display{css_display_mode_ruby_base_container, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("ruby-text-container")) {
		return css_display{css_display_mode_ruby_text_container, 0, 0}, true
	}

	// Try display-box -----------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-box

	if !cm.IsNil(ts.consume_ident_token_with("contents")) {
		return css_display{css_display_mode_contents, 0, 0}, true
	} else if !cm.IsNil(ts.consume_ident_token_with("none")) {
		return css_display{css_display_mode_none, 0, 0}, true
	}

	return css_display{}, false
}

func (ts *css_token_stream) parse_visibility() (css_visibility, bool) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if !cm.IsNil(ts.consume_ident_token_with("visible")) {
		return css_visibility_visible, true
	} else if !cm.IsNil(ts.consume_ident_token_with("hidden")) {
		return css_visibility_hidden, true
	} else if !cm.IsNil(ts.consume_ident_token_with("collapse")) {
		return css_visibility_collapse, true
	}
	return 0, false
}
