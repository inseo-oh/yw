package csssyntax

import (
	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/fonts"
)

// https://www.w3.org/TR/css-fonts-3/#propdef-font-family
func (ts *tokenStream) parseFamilyName() (string, bool) {
	oldCursor := ts.cursor
	if tk, err := ts.consumeTokenWith(tokenTypeString); err == nil {
		return tk.(stringToken).value, true
	} else {
		ts.cursor = oldCursor
	}
	out := ""
	identTks, _ := parseRepeation(ts, 0, func(ts *tokenStream) (token, error) {
		return ts.consumeTokenWith(tokenTypeIdent)
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
	if err := ts.consumeIdentTokenWith("serif"); err == nil {
		return fonts.Serif, true
	}
	if err := ts.consumeIdentTokenWith("sans-serif"); err == nil {
		return fonts.SansSerif, true
	}
	if err := ts.consumeIdentTokenWith("cursive"); err == nil {
		return fonts.Cursive, true
	}
	if err := ts.consumeIdentTokenWith("fantasy"); err == nil {
		return fonts.Fantasy, true
	}
	if err := ts.consumeIdentTokenWith("monospace"); err == nil {
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
	if err := ts.consumeIdentTokenWith("normal"); err == nil {
		return fonts.NormalWeight, true
	}
	if err := ts.consumeIdentTokenWith("bold"); err == nil {
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
	if err := ts.consumeIdentTokenWith("ultra-condensed"); err == nil {
		return fonts.UltraCondensed, true
	}
	if err := ts.consumeIdentTokenWith("extra-condensed"); err == nil {
		return fonts.ExtraCondensed, true
	}
	if err := ts.consumeIdentTokenWith("condensed"); err == nil {
		return fonts.Condensed, true
	}
	if err := ts.consumeIdentTokenWith("semi-condensed"); err == nil {
		return fonts.SemiCondensed, true
	}
	if err := ts.consumeIdentTokenWith("normal"); err == nil {
		return fonts.NormalStretch, true
	}
	if err := ts.consumeIdentTokenWith("semi-expanded"); err == nil {
		return fonts.SemiExpanded, true
	}
	if err := ts.consumeIdentTokenWith("expanded"); err == nil {
		return fonts.Expanded, true
	}
	if err := ts.consumeIdentTokenWith("extra-expanded"); err == nil {
		return fonts.ExtraExpanded, true
	}
	if err := ts.consumeIdentTokenWith("ultra-expanded"); err == nil {
		return fonts.UltraExpanded, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-style
func (ts *tokenStream) parseFontStyle() (fonts.Style, bool) {
	if err := ts.consumeIdentTokenWith("normal"); err == nil {
		return fonts.NormalStyle, true
	}
	if err := ts.consumeIdentTokenWith("italic"); err == nil {
		return fonts.Italic, true
	}
	if err := ts.consumeIdentTokenWith("oblique"); err == nil {
		return fonts.Oblique, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#absolute-size-value
func (ts *tokenStream) parseAbsoluteSize() (fonts.AbsoluteSize, bool) {
	if err := ts.consumeIdentTokenWith("xx-small"); err == nil {
		return fonts.XXSmall, true
	}
	if err := ts.consumeIdentTokenWith("x-small"); err == nil {
		return fonts.XSmall, true
	}
	if err := ts.consumeIdentTokenWith("small"); err == nil {
		return fonts.Small, true
	}
	if err := ts.consumeIdentTokenWith("medium"); err == nil {
		return fonts.MediumSize, true
	}
	if err := ts.consumeIdentTokenWith("large"); err == nil {
		return fonts.Large, true
	}
	if err := ts.consumeIdentTokenWith("x-large"); err == nil {
		return fonts.XLarge, true
	}
	if err := ts.consumeIdentTokenWith("xx-large"); err == nil {
		return fonts.XXLarge, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-fonts-3/#relative-size-value
func (ts *tokenStream) parseRelativeSize() (fonts.RelativeSize, bool) {
	if err := ts.consumeIdentTokenWith("larger"); err == nil {
		return fonts.Larger, true
	} else if err := ts.consumeIdentTokenWith("smaller"); err == nil {
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
	if l, err := ts.parseLengthOrPercentage(true); err == nil {
		return fonts.LengthFontSize{LengthResolvable: l}, true
	}
	return nil, false
}
