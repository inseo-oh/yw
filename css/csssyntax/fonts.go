package csssyntax

import (
	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/fonts"
	"github.com/inseo-oh/yw/util"
)

// https://www.w3.org/TR/css-fonts-3/#propdef-font-family
func (ts *tokenStream) parseFamilyName() (string, bool) {
	if tk := ts.consumeTokenWith(tokenTypeString); !util.IsNil(tk) {
		return tk.(stringToken).value, true
	}
	out := ""
	identTks, _ := parseRepeation(ts, 0, func(ts *tokenStream) (token, error) {
		return ts.consumeTokenWith(tokenTypeIdent), nil
	})
	if identTks == nil {
		return "", false
	}
	for _, tk := range identTks {
		out += tk.(identToken).value
	}
	return out, true
}

// https://www.w3.org/TR/css-fonts-3/#generic-family-value
func (ts *tokenStream) parseGenericFamily() (fonts.FamilyType, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("serif")) {
		return fonts.Serif, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("sans-serif")) {
		return fonts.SansSerif, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("cursive")) {
		return fonts.Cursive, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("fantasy")) {
		return fonts.Fantasy, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("monospace")) {
		return fonts.Monospace, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#font-family-prop
func (ts *tokenStream) parseFontFamily() (fonts.FamilyList, bool) {
	familyPtrs, _ := parseCommaSeparatedRepeation(ts, 0, func(ts *tokenStream) (*fonts.Family, error) {
		if tp, ok := ts.parseGenericFamily(); ok {
			return &fonts.Family{Type: tp}, nil
		}
		if name, ok := ts.parseFamilyName(); ok {
			return &fonts.Family{Type: fonts.NonGeneric, Name: name}, nil
		}
		return nil, nil
	})
	if familyPtrs == nil {
		return fonts.FamilyList{}, false
	}
	families := []fonts.Family{}
	for _, f := range familyPtrs {
		families = append(families, *f)
	}
	return fonts.FamilyList{Families: families}, true
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
func (ts *tokenStream) parseFontWeight() (fonts.Weight, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("normal")) {
		return fonts.NormalWeight, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("bold")) {
		return fonts.Bold, true
	}
	if n := ts.parseNumber(); n != nil {
		if n.Type == css.NumTypeFloat {
			return 0, false
		}
		intVal := n.ToInt()
		if intVal < 0 || 1000 < intVal {
			return 0, false
		}
		return fonts.Weight(n.ToInt()), true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-stretch
func (ts *tokenStream) parseFontStretch() (fonts.Stretch, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("ultra-condensed")) {
		return fonts.UltraCondensed, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("extra-condensed")) {
		return fonts.ExtraCondensed, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("condensed")) {
		return fonts.Condensed, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("semi-condensed")) {
		return fonts.SemiCondensed, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("normal")) {
		return fonts.NormalStretch, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("semi-expanded")) {
		return fonts.SemiExpanded, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("expanded")) {
		return fonts.Expanded, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("extra-expanded")) {
		return fonts.ExtraExpanded, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("ultra-expanded")) {
		return fonts.UltraExpanded, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-style
func (ts *tokenStream) parseFontStyle() (fonts.Style, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("normal")) {
		return fonts.NormalStyle, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("italic")) {
		return fonts.Italic, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("oblique")) {
		return fonts.Oblique, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#absolute-size-value
func (ts *tokenStream) parseAbsoluteSize() (fonts.AbsoluteSize, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("xx-small")) {
		return fonts.XXSmall, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("x-small")) {
		return fonts.XSmall, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("small")) {
		return fonts.Small, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("medium")) {
		return fonts.MediumSize, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("large")) {
		return fonts.Large, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("x-large")) {
		return fonts.XLarge, true
	}
	if !util.IsNil(ts.consumeIdentTokenWith("xx-large")) {
		return fonts.XXLarge, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#relative-size-value
func (ts *tokenStream) parseRelativeSize() (fonts.RelativeSize, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("larger")) {
		return fonts.Larger, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("smaller")) {
		return fonts.Smaller, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
func (ts *tokenStream) parseFontSize() (fonts.Size, bool) {
	if sz, ok := ts.parseAbsoluteSize(); ok {
		return sz, true
	}
	if sz, ok := ts.parseRelativeSize(); ok {
		return sz, true
	}
	if l, _ := ts.parseLengthOrPercentage(true); !util.IsNil(l) {
		return fonts.LengthFontSize{LengthResolvable: l}, true
	}
	return nil, false
}
