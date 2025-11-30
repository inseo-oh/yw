// Implementation of the CSS Backgrounds and Borders Module Level 3 (https://www.w3.org/TR/css-backgrounds-3/)
package libhtml

import (
	"fmt"
)

type css_line_style uint8

const (
	css_line_style_none = css_line_style(iota)
	css_line_style_hidden
	css_line_style_dotted
	css_line_style_dashed
	css_line_style_solid
	css_line_style_double
	css_line_style_groove
	css_line_style_ridge
	css_line_style_inset
	css_line_style_outset
)

func (l css_line_style) String() string {
	switch l {
	case css_line_style_none:
		return "none"
	case css_line_style_hidden:
		return "hidden"
	case css_line_style_dotted:
		return "dotted"
	case css_line_style_dashed:
		return "dashed"
	case css_line_style_solid:
		return "solid"
	case css_line_style_double:
		return "double"
	case css_line_style_groove:
		return "groove"
	case css_line_style_ridge:
		return "ridge"
	case css_line_style_inset:
		return "inset"
	case css_line_style_outset:
		return "outset"
	}
	return fmt.Sprintf("<unknown css_line_style %d>", l)
}

func css_line_width_thin() css_length {
	return css_length{css_number_from_int(1), css_length_unit_px}
}
func css_line_width_medium() css_length {
	return css_length{css_number_from_int(3), css_length_unit_px}
}
func css_line_width_thick() css_length {
	return css_length{css_number_from_int(5), css_length_unit_px}
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
func (ts *css_token_stream) parse_line_style() (css_line_style, bool) {
	if ts.consume_ident_token_with("none") != nil {
		return css_line_style_none, true
	} else if ts.consume_ident_token_with("hidden") != nil {
		return css_line_style_hidden, true
	} else if ts.consume_ident_token_with("hidden") != nil {
		return css_line_style_hidden, true
	} else if ts.consume_ident_token_with("dotted") != nil {
		return css_line_style_dotted, true
	} else if ts.consume_ident_token_with("dashed") != nil {
		return css_line_style_dashed, true
	} else if ts.consume_ident_token_with("solid") != nil {
		return css_line_style_solid, true
	} else if ts.consume_ident_token_with("double") != nil {
		return css_line_style_double, true
	} else if ts.consume_ident_token_with("groove") != nil {
		return css_line_style_groove, true
	} else if ts.consume_ident_token_with("ridge") != nil {
		return css_line_style_ridge, true
	} else if ts.consume_ident_token_with("inset") != nil {
		return css_line_style_inset, true
	} else if ts.consume_ident_token_with("outset") != nil {
		return css_line_style_outset, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
func (ts *css_token_stream) parse_line_width() (css_length, bool) {
	if ts.consume_ident_token_with("thin") != nil {
		return css_line_width_thin(), true
	}
	if ts.consume_ident_token_with("medium") != nil {
		return css_line_width_medium(), true
	}
	if ts.consume_ident_token_with("thick") != nil {
		return css_line_width_thick(), true
	}
	if len, _ := ts.parse_length(true); len != nil {
		return *len, true
	}
	return css_length{}, false
}
