// Implementation of the CSS Fonts Module Level 3 (https://www.w3.org/TR/css-fonts-3/#font-family-prop)
package libhtml

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	cm "yw/libcommon"
)

type css_font_family struct {
	tp   css_font_family_type
	name string // Only valid if the type is "custom"
}

type css_font_family_type uint8

const (
	css_font_family_type_serif = css_font_family_type(iota)
	css_font_family_type_sans_serif
	css_font_family_type_cursive
	css_font_family_type_fantasy
	css_font_family_type_monospace
	css_font_family_type_custom
)

func (f css_font_family) String() string {
	switch f.tp {
	case css_font_family_type_serif:
		return "serif"
	case css_font_family_type_sans_serif:
		return "sans-serif"
	case css_font_family_type_cursive:
		return "cursive"
	case css_font_family_type_fantasy:
		return "fantasy"
	case css_font_family_type_monospace:
		return "monospace"
	case css_font_family_type_custom:
		return strconv.Quote(f.name)
	}
	return fmt.Sprintf("<unknown css_font_family type %d>", f.tp)
}

type css_font_family_list struct {
	families []css_font_family
}

func (f css_font_family_list) String() string {
	out_sb := strings.Builder{}
	for i, f := range f.families {
		if i != 0 {
			out_sb.WriteString(",")
		}
		out_sb.WriteString(fmt.Sprintf("%v", f))
	}
	return out_sb.String()
}

type css_font_weight uint16

const (
	css_font_weight_normal = css_font_weight(400) // https://www.w3.org/TR/css-fonts-3/#font-weight-normal-value
	css_font_weight_bold   = css_font_weight(700) // https://www.w3.org/TR/css-fonts-3/#bold
)

func (w css_font_weight) String() string {
	return fmt.Sprintf("%d", w)
}

type css_font_stretch uint8

const (
	css_font_stretch_ultra_condensed = css_font_stretch(iota)
	css_font_stretch_extra_condensed
	css_font_stretch_condensed
	css_font_stretch_semi_condensed
	css_font_stretch_normal
	css_font_stretch_semi_expanded
	css_font_stretch_expanded
	css_font_stretch_extra_expanded
	css_font_stretch_ultra_expanded
)

func (s css_font_stretch) String() string {
	switch s {
	case css_font_stretch_ultra_condensed:
		return "ultra-condensed"
	case css_font_stretch_extra_condensed:
		return "extra-condensed"
	case css_font_stretch_condensed:
		return "condensed"
	case css_font_stretch_semi_condensed:
		return "semi-condensed"
	case css_font_stretch_normal:
		return "normal"
	case css_font_stretch_semi_expanded:
		return "semi-expanded"
	case css_font_stretch_expanded:
		return "expanded"
	case css_font_stretch_extra_expanded:
		return "extra-expanded"
	case css_font_stretch_ultra_expanded:
		return "ultra-expanded"
	}
	return fmt.Sprintf("<unknown css_font_stretch value %d>", s)
}

type css_font_style uint8

const (
	css_font_style_normal = css_font_style(iota)
	css_font_style_italic
	css_font_style_oblique
)

func (s css_font_style) String() string {
	switch s {
	case css_font_style_normal:
		return "normal"
	case css_font_style_italic:
		return "italic"
	case css_font_style_oblique:
		return "oblique"
	}
	return fmt.Sprintf("<unknown css_font_style value %d>", s)
}

type css_font_size interface {
	calculate_real_font_size() css_length
	String() string
}

type css_absolute_size uint8

const (
	css_absolute_size_xx_small = css_absolute_size(iota)
	css_absolute_size_x_small
	css_absolute_size_small
	css_absolute_size_medium
	css_absolute_size_large
	css_absolute_size_x_large
	css_absolute_size_xx_large
)

func (s css_absolute_size) String() string {
	switch s {
	case css_absolute_size_xx_small:
		return "xx-small"
	case css_absolute_size_x_small:
		return "x-small"
	case css_absolute_size_small:
		return "small"
	case css_absolute_size_medium:
		return "medium"
	case css_absolute_size_large:
		return "large"
	case css_absolute_size_x_large:
		return "x-large"
	case css_absolute_size_xx_large:
		return "xx-large"
	}
	return fmt.Sprintf("<unknown css_absolute_size value %d>", s)
}

func (s css_absolute_size) calculate_real_font_size() css_length {
	panic("TODO")
}

type css_relative_size uint8

const (
	css_relative_size_larger = css_relative_size(iota)
	css_relative_size_smaller
)

func (s css_relative_size) String() string {
	switch s {
	case css_relative_size_larger:
		return "larger"
	case css_relative_size_smaller:
		return "smaller"
	}
	return fmt.Sprintf("<unknown css_relative_size value %d>", s)
}

func (s css_relative_size) calculate_real_font_size() css_length {
	panic("TODO")
}

type css_length_font_size struct{ css_length_resolvable }

func (l css_length_font_size) calculate_real_font_size() css_length {
	panic("TODO")
}

type css_font_shorthand struct {
	style css_font_style
	// TODO: https://www.w3.org/TR/css-fonts-3/#font-variant-css21-values
	weight  css_font_weight
	stretch css_font_stretch
	size    css_font_size
	// TODO: line-height
	family css_font_family_list
}

func (f css_font_shorthand) String() string {
	return fmt.Sprintf("%v %v %v %v %v", f.style, f.weight, f.stretch, f.size, f.family)
}

func init() {
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-family-prop
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-family
	parse_family_name := func(ts *css_token_stream) (*string, error) {
		if tk := ts.consume_token_with_type(css_token_type_string); !cm.IsNil(tk) {
			return cm.MakeStrPtr(tk.(css_string_token).value), nil
		}
		out := ""
		ident_tks, err := css_accept_repetion(ts, 0, func(ts *css_token_stream) (css_token, error) {
			return ts.consume_token_with_type(css_token_type_ident), nil
		})
		if ident_tks == nil {
			return nil, err
		}
		for _, tk := range ident_tks {
			out += tk.(css_ident_token).value
		}
		return &out, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#generic-family-value
	parse_generic_family := func(ts *css_token_stream) (*css_font_family_type, error) {
		var v css_font_family_type
		if !cm.IsNil(ts.consume_ident_token_with("serif")) {
			v = css_font_family_type_serif
		} else if !cm.IsNil(ts.consume_ident_token_with("sans-serif")) {
			v = css_font_family_type_sans_serif
		} else if !cm.IsNil(ts.consume_ident_token_with("cursive")) {
			v = css_font_family_type_cursive
		} else if !cm.IsNil(ts.consume_ident_token_with("fantasy")) {
			v = css_font_family_type_fantasy
		} else if !cm.IsNil(ts.consume_ident_token_with("monospace")) {
			v = css_font_family_type_monospace
		} else {
			return nil, nil
		}
		return &v, nil
	}
	parse_font_family := func(ts *css_token_stream) (*css_font_family_list, error) {
		family_ptrs, err := css_accept_comma_separated_repetion(ts, 0, func(ts *css_token_stream) (*css_font_family, error) {
			if tp, err := parse_generic_family(ts); tp != nil {
				return &css_font_family{*tp, ""}, nil
			} else if err != nil {
				return nil, err
			}
			if name, err := parse_family_name(ts); name != nil {
				return &css_font_family{css_font_family_type_custom, *name}, nil
			} else if err != nil {
				return nil, err
			}
			return nil, nil
		})
		if family_ptrs == nil {
			return nil, err
		}
		families := []css_font_family{}
		for _, f := range family_ptrs {
			families = append(families, *f)
		}
		return &css_font_family_list{families}, nil
	}

	// https://www.w3.org/TR/css-fonts-3/#font-family-prop
	css_property_descriptors_map["font-family"] = css_property_descriptor{
		initial: css_font_family_list{[]css_font_family{{css_font_family_type_sans_serif, ""}}},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return parse_font_family(ts)
		},
	}

	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-weight-prop
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
	parse_font_weight := func(ts *css_token_stream) (*css_font_weight, error) {
		var v css_font_weight
		if !cm.IsNil(ts.consume_ident_token_with("normal")) {
			v = css_font_weight_normal
		} else if !cm.IsNil(ts.consume_ident_token_with("bold")) {
			v = css_font_weight_bold
		} else if n := ts.parse_number(); n != nil {
			if n.tp == css_number_type_float {
				return nil, errors.New("floating point not allowed")
			}
			int_val := n.to_int()
			if int_val < 0 || 1000 < int_val {
				return nil, errors.New("value out of range")
			}
			v = css_font_weight(n.to_int())
		} else {
			return nil, nil
		}
		return &v, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
	css_property_descriptors_map["font-weight"] = css_property_descriptor{
		initial: css_font_weight_normal,
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return parse_font_weight(ts)
		},
	}
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-stretch-prop
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-stretch
	parse_font_stretch := func(ts *css_token_stream) (*css_font_stretch, error) {
		var v css_font_stretch
		if !cm.IsNil(ts.consume_ident_token_with("ultra-condensed")) {
			v = css_font_stretch_ultra_condensed
		} else if !cm.IsNil(ts.consume_ident_token_with("extra-condensed")) {
			v = css_font_stretch_extra_condensed
		} else if !cm.IsNil(ts.consume_ident_token_with("condensed")) {
			v = css_font_stretch_condensed
		} else if !cm.IsNil(ts.consume_ident_token_with("semi-condensed")) {
			v = css_font_stretch_semi_condensed
		} else if !cm.IsNil(ts.consume_ident_token_with("normal")) {
			v = css_font_stretch_normal
		} else if !cm.IsNil(ts.consume_ident_token_with("semi-expanded")) {
			v = css_font_stretch_semi_expanded
		} else if !cm.IsNil(ts.consume_ident_token_with("expanded")) {
			v = css_font_stretch_expanded
		} else if !cm.IsNil(ts.consume_ident_token_with("extra-expanded")) {
			v = css_font_stretch_extra_expanded
		} else if !cm.IsNil(ts.consume_ident_token_with("ultra-expanded")) {
			v = css_font_stretch_ultra_expanded
		} else {
			return nil, nil
		}
		return &v, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-stretch
	css_property_descriptors_map["font-stretch"] = css_property_descriptor{
		initial: css_font_stretch_normal,
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return parse_font_stretch(ts)
		},
	}
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-style-prop
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-style
	parse_font_style := func(ts *css_token_stream) (*css_font_style, error) {
		var v css_font_style
		if !cm.IsNil(ts.consume_ident_token_with("normal")) {
			v = css_font_style_normal
		} else if !cm.IsNil(ts.consume_ident_token_with("italic")) {
			v = css_font_style_italic
		} else if !cm.IsNil(ts.consume_ident_token_with("oblique")) {
			v = css_font_style_oblique
		} else {
			return nil, nil
		}
		return &v, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-style
	css_property_descriptors_map["font-style"] = css_property_descriptor{
		initial: css_font_style_normal,
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return parse_font_style(ts)
		},
	}
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-size-prop
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#absolute-size-value
	parse_absolute_size := func(ts *css_token_stream) (*css_absolute_size, error) {
		var v css_absolute_size
		if !cm.IsNil(ts.consume_ident_token_with("xx-small")) {
			v = css_absolute_size_xx_small
		} else if !cm.IsNil(ts.consume_ident_token_with("x-small")) {
			v = css_absolute_size_x_small
		} else if !cm.IsNil(ts.consume_ident_token_with("small")) {
			v = css_absolute_size_small
		} else if !cm.IsNil(ts.consume_ident_token_with("medium")) {
			v = css_absolute_size_medium
		} else if !cm.IsNil(ts.consume_ident_token_with("large")) {
			v = css_absolute_size_large
		} else if !cm.IsNil(ts.consume_ident_token_with("x-large")) {
			v = css_absolute_size_x_large
		} else if !cm.IsNil(ts.consume_ident_token_with("xx-large")) {
			v = css_absolute_size_xx_large
		} else {
			return nil, nil
		}
		return &v, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#relative-size-value
	parse_relative_size := func(ts *css_token_stream) (*css_relative_size, error) {
		var v css_relative_size
		if !cm.IsNil(ts.consume_ident_token_with("larger")) {
			v = css_relative_size_larger
		} else if !cm.IsNil(ts.consume_ident_token_with("smaller")) {
			v = css_relative_size_smaller
		} else {
			return nil, nil
		}
		return &v, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
	parse_font_size := func(ts *css_token_stream) (css_font_size, error) {
		if sz, err := parse_absolute_size(ts); sz != nil {
			return *sz, nil
		} else if err != nil {
			return nil, nil
		}
		if sz, err := parse_relative_size(ts); sz != nil {
			return *sz, nil
		} else if err != nil {
			return nil, nil
		}
		if l, err := ts.parse_length_or_percentage(true); !cm.IsNil(l) {
			return css_length_font_size{l}, nil
		} else if err != nil {
			return nil, nil
		}
		return nil, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
	css_property_descriptors_map["font-size"] = css_property_descriptor{
		initial: css_absolute_size_medium,
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return parse_font_size(ts)
		},
	}
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-size-adjust-prop
	//==========================================================================
	// TODO

	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-prop
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#propdef-font
	parse_font_shorthand := func(ts *css_token_stream) (*css_font_shorthand, error) {
		out := css_font_shorthand{
			css_property_descriptors_map["font-style"].initial.(css_font_style),
			css_property_descriptors_map["font-weight"].initial.(css_font_weight),
			css_property_descriptors_map["font-stretch"].initial.(css_font_stretch),
			css_property_descriptors_map["font-size"].initial.(css_font_size),
			css_property_descriptors_map["font-family"].initial.(css_font_family_list),
		}
		got_style := false
		got_weight := false
		got_stretch := false
		got_size := false
		got_family := false
		for !got_style || !got_weight || !got_stretch || !got_size || !got_family {
			invalid := true
			if !got_style {
				ts.skip_whitespaces()
				if res, err := parse_font_style(ts); !cm.IsNil(res) {
					out.style = *res
					invalid = false
					got_style = true
				} else if err != nil {
					return nil, err
				}
			}
			if !got_weight {
				ts.skip_whitespaces()
				if res, err := parse_font_weight(ts); !cm.IsNil(res) {
					out.weight = *res
					invalid = false
					got_weight = true
				} else if err != nil {
					return nil, err
				}
			}
			if !got_stretch {
				ts.skip_whitespaces()
				if res, err := parse_font_stretch(ts); !cm.IsNil(res) {
					out.stretch = *res
					invalid = false
					got_stretch = true
				} else if err != nil {
					return nil, err
				}
			}
			if !got_size {
				ts.skip_whitespaces()
				if res, err := parse_font_size(ts); !cm.IsNil(res) {
					out.size = res
					invalid = false
					got_size = true
				} else if err != nil {
					return nil, err
				}
			}
			if !got_family {
				ts.skip_whitespaces()
				if res, err := parse_font_family(ts); !cm.IsNil(res) {
					out.family = *res
					invalid = false
					got_family = true
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
		if !got_style && !got_weight && !got_stretch && !got_size && !got_family {
			return nil, nil
		}
		return &out, nil
	}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font
	css_property_descriptors_map["font"] = css_property_descriptor{
		initial: css_font_shorthand{
			css_property_descriptors_map["font-style"].initial.(css_font_style),
			css_property_descriptors_map["font-weight"].initial.(css_font_weight),
			css_property_descriptors_map["font-stretch"].initial.(css_font_stretch),
			css_property_descriptors_map["font-size"].initial.(css_font_size),
			css_property_descriptors_map["font-family"].initial.(css_font_family_list),
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return parse_font_shorthand(ts)
		},
	}
}
