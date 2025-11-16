// Implementation of the CSS Backgrounds and Borders Module Level 3 (https://www.w3.org/TR/css-backgrounds-3/)
package libhtml

import (
	"errors"
	"fmt"
	cm "yw/libcommon"
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

type css_border_color_shorthand struct{ top, right, bottom, left css_color }

func (c css_border_color_shorthand) String() string {
	return fmt.Sprintf("%v %v %v %v", c.top, c.right, c.bottom, c.left)
}

type css_border_style_shorthand struct{ top, right, bottom, left css_line_style }

func (c css_border_style_shorthand) String() string {
	return fmt.Sprintf("%v %v %v %v", c.top, c.right, c.bottom, c.left)
}

type css_border_width_shorthand struct{ top, right, bottom, left css_length }

func (c css_border_width_shorthand) String() string {
	return fmt.Sprintf("%v %v %v %v", c.top, c.right, c.bottom, c.left)
}

type css_border_shorthand struct {
	width css_length
	style css_line_style
	color css_color
}

func (c css_border_shorthand) String() string {
	return fmt.Sprintf("%v %v %v", c.width, c.style, c.color)
}

func init() {
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#background-color
	//==========================================================================

	// https://www.w3.org/TR/css-backgrounds-3/#propdef-background-color
	css_property_descriptors_map["background-color"] = css_property_descriptor{
		initial: css_color_transparent(),
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return ts.parse_color()
		},
	}
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#border-color
	//==========================================================================

	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-top-color
	css_property_descriptors_map["border-top-color"] = css_property_descriptor{
		initial: css_color{css_color_type_current_color, nil},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return ts.parse_color()
		},
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-right-color
	css_property_descriptors_map["border-right-color"] = css_property_descriptor{
		initial: css_color{css_color_type_current_color, nil},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return ts.parse_color()
		},
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-bottom-color
	css_property_descriptors_map["border-bottom-color"] = css_property_descriptor{
		initial: css_color{css_color_type_current_color, nil},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return ts.parse_color()
		},
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-left-color
	css_property_descriptors_map["border-left-color"] = css_property_descriptor{
		initial: css_color{css_color_type_current_color, nil},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return ts.parse_color()
		},
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-color
	css_property_descriptors_map["border-color"] = css_property_descriptor{
		initial: css_border_color_shorthand{
			css_property_descriptors_map["border-top-color"].initial.(css_color),
			css_property_descriptors_map["border-right-color"].initial.(css_color),
			css_property_descriptors_map["border-bottom-color"].initial.(css_color),
			css_property_descriptors_map["border-left-color"].initial.(css_color),
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			items, err := css_accept_repetion(ts, 4, (*css_token_stream).parse_color)
			if items == nil {
				return nil, err
			}
			res := css_border_color_shorthand{}
			switch len(items) {
			case 1:
				res.top = *items[0]
				res.bottom = *items[0]
				res.left = *items[0]
				res.right = *items[0]
			case 2:
				res.top = *items[0]
				res.bottom = *items[0]
				res.left = *items[1]
				res.right = *items[1]
			case 3:
				res.top = *items[0]
				res.bottom = *items[1]
				res.left = *items[2]
				res.right = *items[2]
			case 4:
				res.top = *items[0]
				res.bottom = *items[1]
				res.left = *items[2]
				res.right = *items[3]
			default:
				return nil, errors.New("too few values")
			}
			return res, nil
		},
	}

	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#border-style
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
	parse_line_style := func(ts *css_token_stream) (css_property_value, error) {
		if ts.consume_ident_token_with("none") != nil {
			return css_line_style_none, nil
		} else if ts.consume_ident_token_with("hidden") != nil {
			return css_line_style_hidden, nil
		} else if ts.consume_ident_token_with("hidden") != nil {
			return css_line_style_hidden, nil
		} else if ts.consume_ident_token_with("dotted") != nil {
			return css_line_style_dotted, nil
		} else if ts.consume_ident_token_with("dashed") != nil {
			return css_line_style_dashed, nil
		} else if ts.consume_ident_token_with("solid") != nil {
			return css_line_style_solid, nil
		} else if ts.consume_ident_token_with("double") != nil {
			return css_line_style_double, nil
		} else if ts.consume_ident_token_with("groove") != nil {
			return css_line_style_groove, nil
		} else if ts.consume_ident_token_with("ridge") != nil {
			return css_line_style_ridge, nil
		} else if ts.consume_ident_token_with("inset") != nil {
			return css_line_style_inset, nil
		} else if ts.consume_ident_token_with("outset") != nil {
			return css_line_style_outset, nil
		}
		return nil, nil
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-top-style
	css_property_descriptors_map["border-top-style"] = css_property_descriptor{
		initial:    css_line_style_none,
		parse_func: parse_line_style,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-right-style
	css_property_descriptors_map["border-right-style"] = css_property_descriptor{
		initial:    css_line_style_none,
		parse_func: parse_line_style,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-bottom-style
	css_property_descriptors_map["border-bottom-style"] = css_property_descriptor{
		initial:    css_line_style_none,
		parse_func: parse_line_style,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-left-style
	css_property_descriptors_map["border-left-style"] = css_property_descriptor{
		initial:    css_line_style_none,
		parse_func: parse_line_style,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-style
	css_property_descriptors_map["border-style"] = css_property_descriptor{
		initial: css_border_style_shorthand{
			css_property_descriptors_map["border-top-style"].initial.(css_line_style),
			css_property_descriptors_map["border-right-style"].initial.(css_line_style),
			css_property_descriptors_map["border-bottom-style"].initial.(css_line_style),
			css_property_descriptors_map["border-left-style"].initial.(css_line_style),
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			items, err := css_accept_repetion(ts, 4, parse_line_style)
			if items == nil {
				return nil, err
			}
			res := css_border_style_shorthand{}
			switch len(items) {
			case 1:
				res.top = items[0].(css_line_style)
				res.bottom = items[0].(css_line_style)
				res.left = items[0].(css_line_style)
				res.right = items[0].(css_line_style)
			case 2:
				res.top = items[0].(css_line_style)
				res.bottom = items[0].(css_line_style)
				res.left = items[1].(css_line_style)
				res.right = items[1].(css_line_style)
			case 3:
				res.top = items[0].(css_line_style)
				res.bottom = items[1].(css_line_style)
				res.left = items[2].(css_line_style)
				res.right = items[2].(css_line_style)
			case 4:
				res.top = items[0].(css_line_style)
				res.bottom = items[1].(css_line_style)
				res.left = items[2].(css_line_style)
				res.right = items[3].(css_line_style)
			default:
				return nil, errors.New("too few values")
			}
			return res, nil
		},
	}

	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#border-width
	//==========================================================================
	line_width_thin := css_length{css_number_from_int(1), css_length_unit_px}
	line_width_medium := css_length{css_number_from_int(3), css_length_unit_px}
	line_width_thick := css_length{css_number_from_int(5), css_length_unit_px}
	// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
	parse_line_width := func(ts *css_token_stream) (css_property_value, error) {
		if ts.consume_ident_token_with("thin") != nil {
			return line_width_thin, nil
		} else if ts.consume_ident_token_with("medium") != nil {
			return line_width_medium, nil
		} else if ts.consume_ident_token_with("thick") != nil {
			return line_width_thick, nil
		}
		if len, err := ts.parse_length(true); len != nil {
			return *len, nil
		} else {
			return nil, err
		}
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-top-width
	css_property_descriptors_map["border-top-width"] = css_property_descriptor{
		initial:    line_width_medium,
		parse_func: parse_line_width,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-right-width
	css_property_descriptors_map["border-right-width"] = css_property_descriptor{
		initial:    line_width_medium,
		parse_func: parse_line_width,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-bottom-width
	css_property_descriptors_map["border-bottom-width"] = css_property_descriptor{
		initial:    line_width_medium,
		parse_func: parse_line_width,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-left-width
	css_property_descriptors_map["border-left-width"] = css_property_descriptor{
		initial:    line_width_medium,
		parse_func: parse_line_width,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-width
	css_property_descriptors_map["border-width"] = css_property_descriptor{
		initial: css_border_width_shorthand{
			css_property_descriptors_map["border-top-width"].initial.(css_length),
			css_property_descriptors_map["border-right-width"].initial.(css_length),
			css_property_descriptors_map["border-bottom-width"].initial.(css_length),
			css_property_descriptors_map["border-left-width"].initial.(css_length),
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			items, err := css_accept_repetion(ts, 4, parse_line_width)
			if items == nil {
				return nil, err
			}
			res := css_border_width_shorthand{}
			switch len(items) {
			case 1:
				res.top = items[0].(css_length)
				res.bottom = items[0].(css_length)
				res.left = items[0].(css_length)
				res.right = items[0].(css_length)
			case 2:
				res.top = items[0].(css_length)
				res.bottom = items[0].(css_length)
				res.left = items[1].(css_length)
				res.right = items[1].(css_length)
			case 3:
				res.top = items[0].(css_length)
				res.bottom = items[1].(css_length)
				res.left = items[2].(css_length)
				res.right = items[2].(css_length)
			case 4:
				res.top = items[0].(css_length)
				res.bottom = items[1].(css_length)
				res.left = items[2].(css_length)
				res.right = items[3].(css_length)
			default:
				return nil, errors.New("too few values")
			}
			return res, nil
		},
	}

	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#border-shorthands
	//==========================================================================
	parse_border_shorthand := func(ts *css_token_stream) (css_property_value, error) {
		out := css_border_shorthand{
			css_property_descriptors_map["border-top-width"].initial.(css_length),
			css_property_descriptors_map["border-top-style"].initial.(css_line_style),
			css_property_descriptors_map["border-top-color"].initial.(css_color),
		}
		got_width := false
		got_style := false
		got_color := false
		for !got_width || !got_style || !got_color {
			invalid := true
			if !got_width {
				ts.skip_whitespaces()
				if res, err := parse_line_width(ts); !cm.IsNil(res) {
					out.width = res.(css_length)
					invalid = false
					got_width = true
				} else if err != nil {
					return nil, err
				}
			}
			if !got_style {
				ts.skip_whitespaces()
				if res, err := parse_line_style(ts); !cm.IsNil(res) {
					out.style = res.(css_line_style)
					invalid = false
					got_style = true
				} else if err != nil {
					return nil, err
				}
			}
			if !got_color {
				ts.skip_whitespaces()
				if res, err := ts.parse_color(); !cm.IsNil(res) {
					out.color = *res
					invalid = false
					got_color = true
				} else if err != nil {
					return nil, err
				}
			}
			ts.skip_whitespaces()
			if ts.is_end() {
				invalid = false
				break
			}
			if invalid {
				return nil, errors.New("got unexpected token")
			}
		}
		if !got_width && !got_style && !got_color {
			return nil, nil
		}
		return out, nil
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-top
	css_property_descriptors_map["border-top"] = css_property_descriptor{
		initial: css_border_shorthand{
			css_property_descriptors_map["border-top-width"].initial.(css_length),
			css_property_descriptors_map["border-top-style"].initial.(css_line_style),
			css_property_descriptors_map["border-top-color"].initial.(css_color),
		},
		parse_func: parse_border_shorthand,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-right
	css_property_descriptors_map["border-right"] = css_property_descriptor{
		initial: css_border_shorthand{
			css_property_descriptors_map["border-right-width"].initial.(css_length),
			css_property_descriptors_map["border-right-style"].initial.(css_line_style),
			css_property_descriptors_map["border-right-color"].initial.(css_color),
		},
		parse_func: parse_border_shorthand,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-bottom
	css_property_descriptors_map["border-bottom"] = css_property_descriptor{
		initial: css_border_shorthand{
			css_property_descriptors_map["border-bottom-width"].initial.(css_length),
			css_property_descriptors_map["border-bottom-style"].initial.(css_line_style),
			css_property_descriptors_map["border-bottom-color"].initial.(css_color),
		},
		parse_func: parse_border_shorthand,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border-left
	css_property_descriptors_map["border-left"] = css_property_descriptor{
		initial: css_border_shorthand{
			css_property_descriptors_map["border-left-width"].initial.(css_length),
			css_property_descriptors_map["border-left-style"].initial.(css_line_style),
			css_property_descriptors_map["border-left-color"].initial.(css_color),
		},
		parse_func: parse_border_shorthand,
	}
	// https://www.w3.org/TR/css-backgrounds-3/#propdef-border
	css_property_descriptors_map["border"] = css_property_descriptor{
		initial: css_border_shorthand{
			css_property_descriptors_map["border-top-width"].initial.(css_length),
			css_property_descriptors_map["border-top-style"].initial.(css_line_style),
			css_property_descriptors_map["border-top-color"].initial.(css_color),
		},
		parse_func: parse_border_shorthand,
	}

}
