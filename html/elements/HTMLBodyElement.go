package elements

import (
	"strconv"
	"strings"
	css_color "yw/css/color"
	"yw/dom"
	cm "yw/libcommon"
)

// https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#rules-for-parsing-a-legacy-colour-value
func html_parse_legacy_color_value(input string) (css_color.Color, bool) {
	if input == "" {
		return css_color.Color{}, false
	}
	input = strings.Trim(input, " ")
	if cm.ToAsciiLowercase(input) == "transparent" {
		// transparent
		return css_color.Color{}, false
	}
	if col, ok := css_color.NamedColors[cm.ToAsciiLowercase(input)]; ok {
		// CSS named colors
		return css_color.NewRgba(col.R, col.G, col.B, col.A), true
	}
	input_cps := []rune(input)
	if len(input_cps) == 4 && input_cps[0] == '#' {
		// #rgb
		red, err1 := strconv.ParseInt(string(input_cps[1]), 16, 8)
		green, err2 := strconv.ParseInt(string(input_cps[2]), 16, 8)
		blue, err3 := strconv.ParseInt(string(input_cps[3]), 16, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			return css_color.NewRgba(uint8(red), uint8(green), uint8(blue), 255), true
		}
	}
	// Now we assume the format is #rrggbb -------------------------------------
	new_input_cps := make([]rune, 0, len(input_cps))
	for i := range len(input_cps) {
		// Replace characters beyond BMP with "00"
		if input_cps[i] > 0xffff {
			new_input_cps = append(new_input_cps, '0')
			new_input_cps = append(new_input_cps, '0')
		} else {
			new_input_cps = append(new_input_cps, input_cps[i])
		}
	}
	input_cps = new_input_cps
	if 128 < len(input_cps) {
		input_cps = input_cps[:128]
	}
	if input_cps[0] == '#' {
		input_cps = input_cps[1:]
	}
	for i := range len(input_cps) {
		// Replace non-hex characters with '0'
		if _, err := strconv.ParseInt(string(input_cps[i]), 16, 8); err != nil {
			input_cps[i] = '0'
		}
	}
	// Length must be nonzero, and multiple of 3. If not, append '0's.
	for len(input_cps) == 0 || len(input_cps)%3 != 0 {
		input_cps = append(input_cps, '0')
	}
	comp_len := len(input_cps) / 3
	comps := [][]rune{
		input_cps[:comp_len*1],
		input_cps[comp_len*1 : comp_len*2],
		input_cps[comp_len*2 : comp_len*3],
	}
	if comp_len > 8 {
		for i := range 3 {
			comps[i] = comps[i][comp_len-8:]
		}
		comp_len = 8
	}
	for comp_len > 2 {
		for i := range 3 {
			if comps[i][0] == '0' {
				comps[i] = comps[i][1:]
			}
		}
		comp_len--
	}
	if comp_len > 2 {
		for i := range 3 {
			comps[i] = comps[i][:2]
		}
		comp_len = 2
	}
	red, err1 := strconv.ParseInt(string(comps[0]), 16, 16)
	green, err2 := strconv.ParseInt(string(comps[1]), 16, 16)
	blue, err3 := strconv.ParseInt(string(comps[2]), 16, 16)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("unreachable")
	}
	return css_color.NewRgba(uint8(red), uint8(green), uint8(blue), 255), true
}

func NewHTMLBodyElement(options dom.ElementCreationCommonOptions) HTMLElement {
	elem := NewHTMLElement(options)

	cbs := elem.Callbacks()
	cbs.PresentationalHints = func() any {
		decls := []css_declaration{}

		// https://html.spec.whatwg.org/multipage/rendering.html#the-page
		if attr, ok := elem.AttrWithoutNamespace("bgcolor"); ok {
			color, ok := html_parse_legacy_color_value(attr)
			if ok {
				decls = append(decls, css_declaration{"background-color", color, false})
			}
		}
		if attr, ok := elem.AttrWithoutNamespace("text"); ok {
			color, ok := html_parse_legacy_color_value(attr)
			if ok {
				decls = append(decls, css_declaration{"color", color, false})
			}
		}
		rule := css_style_rule{
			selector_list: []css_selector{css_node_ptr_selector{elem}},
			declarations:  decls,
		}
		return []css_style_rule{rule}

	}
	return elem
}
