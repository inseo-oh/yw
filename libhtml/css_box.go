// Implementation of the CSS Box Model Module 3 (https://www.w3.org/TR/css-box-3/)
package libhtml

import (
	"errors"
	"fmt"
	cm "yw/libcommon"
)

type css_padding_shorthand struct{ top, right, bottom, left css_length_resolvable }

func (p css_padding_shorthand) String() string {
	return fmt.Sprintf("%v %v %v %v", p.top, p.right, p.bottom, p.left)
}

type css_margin struct {
	v css_length_resolvable // nil means auto
}

func (m css_margin) is_auto() bool { return cm.IsNil(m.v) }
func (m css_margin) String() string {
	if m.is_auto() {
		return "auto"
	}
	return fmt.Sprintf("%v", m.v)
}

type css_margin_shorthand struct{ top, right, bottom, left css_length_resolvable }

func (m css_margin_shorthand) String() string {
	return fmt.Sprintf("%v %v %v %v", m.top, m.right, m.bottom, m.left)
}

func init() {
	// https://www.w3.org/TR/css-box-3/#margin-physical
	parse_margin := func(ts *css_token_stream) (css_property_value, error) {
		if v, err := ts.parse_length_or_percentage(true); !cm.IsNil(v) {
			return css_margin{v}, nil
		} else if err != nil {
			return nil, err
		}
		if ts.consume_ident_token_with("auto") != nil {
			return css_margin{nil}, nil
		}
		return nil, nil
	}
	css_property_descriptors_map["margin-top"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_margin,
	}
	css_property_descriptors_map["margin-right"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_margin,
	}
	css_property_descriptors_map["margin-bottom"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_margin,
	}
	css_property_descriptors_map["margin-left"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_margin,
	}
	// https://www.w3.org/TR/css-box-3/#margin-shorthand
	css_property_descriptors_map["margin"] = css_property_descriptor{
		initial: css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			items, err := css_accept_repetion(ts, 4, parse_margin)
			if items == nil {
				return nil, err
			}
			res := css_margin_shorthand{}
			switch len(items) {
			case 1:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[0].(css_length_resolvable)
				res.left = items[0].(css_length_resolvable)
				res.right = items[0].(css_length_resolvable)
			case 2:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[0].(css_length_resolvable)
				res.left = items[1].(css_length_resolvable)
				res.right = items[1].(css_length_resolvable)
			case 3:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[1].(css_length_resolvable)
				res.left = items[2].(css_length_resolvable)
				res.right = items[2].(css_length_resolvable)
			case 4:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[1].(css_length_resolvable)
				res.left = items[2].(css_length_resolvable)
				res.right = items[3].(css_length_resolvable)
			}
			return res, nil
		},
	}

	// https://www.w3.org/TR/css-box-3/#padding-physical
	parse_padding := func(ts *css_token_stream) (css_property_value, error) {
		v, err := ts.parse_length_or_percentage(true)
		if cm.IsNil(v) {
			return nil, err
		} else if len, ok := v.(css_length); ok && len.value.to_int() < 0 {
			return nil, errors.New("value of of range")
		}
		return v, nil
	}
	css_property_descriptors_map["padding-top"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_padding,
	}
	css_property_descriptors_map["padding-right"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_padding,
	}
	css_property_descriptors_map["padding-bottom"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_padding,
	}
	css_property_descriptors_map["padding-left"] = css_property_descriptor{
		initial:    css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: parse_padding,
	}
	// https://www.w3.org/TR/css-box-3/#padding-shorthand
	css_property_descriptors_map["padding"] = css_property_descriptor{
		initial: css_length{css_number_from_int(0), css_length_unit_px},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			items, err := css_accept_repetion(ts, 4, parse_padding)
			if items == nil {
				return nil, err
			}
			res := css_padding_shorthand{}
			switch len(items) {
			case 1:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[0].(css_length_resolvable)
				res.left = items[0].(css_length_resolvable)
				res.right = items[0].(css_length_resolvable)
			case 2:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[0].(css_length_resolvable)
				res.left = items[1].(css_length_resolvable)
				res.right = items[1].(css_length_resolvable)
			case 3:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[1].(css_length_resolvable)
				res.left = items[2].(css_length_resolvable)
				res.right = items[2].(css_length_resolvable)
			case 4:
				res.top = items[0].(css_length_resolvable)
				res.bottom = items[1].(css_length_resolvable)
				res.left = items[2].(css_length_resolvable)
				res.right = items[3].(css_length_resolvable)
			}
			return res, nil
		},
	}
}
