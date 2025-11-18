// Implementation of the CSS Fonts Module Level 3 (https://www.w3.org/TR/css-fonts-3/#font-family-prop)
package libhtml

import (
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

// https://www.w3.org/TR/css-fonts-3/#propdef-font-family
func (ts *css_token_stream) parse_family_name() (string, bool) {
	if tk := ts.consume_token_with_type(css_token_type_string); !cm.IsNil(tk) {
		return tk.(css_string_token).value, true
	}
	out := ""
	ident_tks, _ := css_accept_repetion(ts, 0, func(ts *css_token_stream) (css_token, error) {
		return ts.consume_token_with_type(css_token_type_ident), nil
	})
	if ident_tks == nil {
		return "", false
	}
	for _, tk := range ident_tks {
		out += tk.(css_ident_token).value
	}
	return out, true
}

// https://www.w3.org/TR/css-fonts-3/#generic-family-value
func (ts *css_token_stream) parse_generic_family() (css_font_family_type, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("serif")) {
		return css_font_family_type_serif, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("sans-serif")) {
		return css_font_family_type_sans_serif, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("cursive")) {
		return css_font_family_type_cursive, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("fantasy")) {
		return css_font_family_type_fantasy, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("monospace")) {
		return css_font_family_type_monospace, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#font-family-prop
func (ts *css_token_stream) parse_font_family() (css_font_family_list, bool) {
	family_ptrs, _ := css_accept_comma_separated_repetion(ts, 0, func(ts *css_token_stream) (*css_font_family, error) {
		if tp, ok := ts.parse_generic_family(); ok {
			return &css_font_family{tp, ""}, nil
		}
		if name, ok := ts.parse_family_name(); ok {
			return &css_font_family{css_font_family_type_custom, name}, nil
		}
		return nil, nil
	})
	if family_ptrs == nil {
		return css_font_family_list{}, false
	}
	families := []css_font_family{}
	for _, f := range family_ptrs {
		families = append(families, *f)
	}
	return css_font_family_list{families}, true
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
func (ts *css_token_stream) parse_font_weight() (css_font_weight, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("normal")) {
		return css_font_weight_normal, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("bold")) {
		return css_font_weight_bold, true
	}
	if n := ts.parse_number(); n != nil {
		if n.tp == css_number_type_float {
			return 0, false
		}
		int_val := n.to_int()
		if int_val < 0 || 1000 < int_val {
			return 0, false
		}
		return css_font_weight(n.to_int()), true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-stretch
func (ts *css_token_stream) parse_font_stretch() (css_font_stretch, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("ultra-condensed")) {
		return css_font_stretch_ultra_condensed, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("extra-condensed")) {
		return css_font_stretch_extra_condensed, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("condensed")) {
		return css_font_stretch_condensed, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("semi-condensed")) {
		return css_font_stretch_semi_condensed, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("normal")) {
		return css_font_stretch_normal, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("semi-expanded")) {
		return css_font_stretch_semi_expanded, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("expanded")) {
		return css_font_stretch_expanded, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("extra-expanded")) {
		return css_font_stretch_extra_expanded, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("ultra-expanded")) {
		return css_font_stretch_ultra_expanded, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-style
func (ts *css_token_stream) parse_font_style() (css_font_style, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("normal")) {
		return css_font_style_normal, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("italic")) {
		return css_font_style_italic, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("oblique")) {
		return css_font_style_oblique, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#absolute-size-value
func (ts *css_token_stream) parse_absolute_size() (css_absolute_size, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("xx-small")) {
		return css_absolute_size_xx_small, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("x-small")) {
		return css_absolute_size_x_small, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("small")) {
		return css_absolute_size_small, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("medium")) {
		return css_absolute_size_medium, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("large")) {
		return css_absolute_size_large, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("x-large")) {
		return css_absolute_size_x_large, true
	}
	if !cm.IsNil(ts.consume_ident_token_with("xx-large")) {
		return css_absolute_size_xx_large, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#relative-size-value
func (ts *css_token_stream) parse_relative_size() (css_relative_size, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("larger")) {
		return css_relative_size_larger, true
	} else if !cm.IsNil(ts.consume_ident_token_with("smaller")) {
		return css_relative_size_smaller, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
func (ts *css_token_stream) parse_font_size() (css_font_size, bool) {
	if sz, ok := ts.parse_absolute_size(); ok {
		return sz, true
	}
	if sz, ok := ts.parse_relative_size(); ok {
		return sz, true
	}
	if l, _ := ts.parse_length_or_percentage(true); !cm.IsNil(l) {
		return css_length_font_size{l}, true
	}
	return nil, false
}
