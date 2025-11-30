package elements

import (
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/css/csscolor"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/selector"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

// https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#rules-for-parsing-a-legacy-colour-value
func parseLegacyColor(input string) (csscolor.Color, bool) {
	if input == "" {
		return csscolor.Color{}, false
	}
	input = strings.Trim(input, " ")
	if util.ToAsciiLowercase(input) == "transparent" {
		// transparent
		return csscolor.Color{}, false
	}
	if col, ok := csscolor.NamedColors[util.ToAsciiLowercase(input)]; ok {
		// CSS named colors
		return csscolor.NewRgba(col.R, col.G, col.B, col.A), true
	}
	inputCps := []rune(input)
	if len(inputCps) == 4 && inputCps[0] == '#' {
		// #rgb
		red, err1 := strconv.ParseInt(string(inputCps[1]), 16, 8)
		green, err2 := strconv.ParseInt(string(inputCps[2]), 16, 8)
		blue, err3 := strconv.ParseInt(string(inputCps[3]), 16, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			return csscolor.NewRgba(uint8(red), uint8(green), uint8(blue), 255), true
		}
	}
	// Now we assume the format is #rrggbb -------------------------------------
	newInputCps := make([]rune, 0, len(inputCps))
	for i := range len(inputCps) {
		// Replace characters beyond BMP with "00"
		if inputCps[i] > 0xffff {
			newInputCps = append(newInputCps, '0')
			newInputCps = append(newInputCps, '0')
		} else {
			newInputCps = append(newInputCps, inputCps[i])
		}
	}
	inputCps = newInputCps
	if 128 < len(inputCps) {
		inputCps = inputCps[:128]
	}
	if inputCps[0] == '#' {
		inputCps = inputCps[1:]
	}
	for i := range len(inputCps) {
		// Replace non-hex characters with '0'
		if _, err := strconv.ParseInt(string(inputCps[i]), 16, 8); err != nil {
			inputCps[i] = '0'
		}
	}
	// Length must be nonzero, and multiple of 3. If not, append '0's.
	for len(inputCps) == 0 || len(inputCps)%3 != 0 {
		inputCps = append(inputCps, '0')
	}
	compLen := len(inputCps) / 3
	comps := [][]rune{
		inputCps[:compLen*1],
		inputCps[compLen*1 : compLen*2],
		inputCps[compLen*2 : compLen*3],
	}
	if compLen > 8 {
		for i := range 3 {
			comps[i] = comps[i][compLen-8:]
		}
		compLen = 8
	}
	for compLen > 2 {
		for i := range 3 {
			if comps[i][0] == '0' {
				comps[i] = comps[i][1:]
			}
		}
		compLen--
	}
	if compLen > 2 {
		for i := range 3 {
			comps[i] = comps[i][:2]
		}
		compLen = 2
	}
	red, err1 := strconv.ParseInt(string(comps[0]), 16, 16)
	green, err2 := strconv.ParseInt(string(comps[1]), 16, 16)
	blue, err3 := strconv.ParseInt(string(comps[2]), 16, 16)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("unreachable")
	}
	return csscolor.NewRgba(uint8(red), uint8(green), uint8(blue), 255), true
}

func NewHTMLBodyElement(options dom.ElementCreationCommonOptions) HTMLElement {
	elem := NewHTMLElement(options)

	cbs := elem.Callbacks()
	cbs.PresentationalHints = func() any {
		decls := []cssom.Declaration{}

		// https://html.spec.whatwg.org/multipage/rendering.html#the-page
		if attr, ok := elem.AttrWithoutNamespace("bgcolor"); ok {
			color, ok := parseLegacyColor(attr)
			if ok {
				decls = append(decls, cssom.Declaration{Name: "background-color", Value: color, IsImportant: false})
			}
		}
		if attr, ok := elem.AttrWithoutNamespace("text"); ok {
			color, ok := parseLegacyColor(attr)
			if ok {
				decls = append(decls, cssom.Declaration{Name: "color", Value: color, IsImportant: false})
			}
		}
		rule := cssom.StyleRule{
			SelectorList: []selector.Selector{selector.NodePtrSelector{Element: elem}},
			Declarations: decls,
		}
		return []cssom.StyleRule{rule}

	}
	return elem
}
