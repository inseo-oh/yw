// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"errors"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/fonts"
)

// https://www.w3.org/TR/css-fonts-3/#propdef-font-family
func (ts *tokenStream) parseFamilyName() (string, error) {
	oldCursor := ts.cursor
	if tk, err := ts.consumeTokenWith(tokenTypeString); err == nil {
		return tk.(stringToken).value, nil
	} else {
		ts.cursor = oldCursor
	}
	out := ""
	identTks, err := parseRepeation(ts, 0, func(ts *tokenStream) (token, error) {
		return ts.consumeTokenWith(tokenTypeIdent)
	})
	if err != nil {
		return "", err
	}
	for _, tk := range identTks {
		out += tk.(identToken).value
	}
	return out, nil
}

// https://www.w3.org/TR/css-fonts-3/#generic-family-value
func (ts *tokenStream) parseGenericFamily() (fonts.FamilyType, error) {
	if err := ts.consumeIdentTokenWith("serif"); err == nil {
		return fonts.Serif, nil
	}
	if err := ts.consumeIdentTokenWith("sans-serif"); err == nil {
		return fonts.SansSerif, nil
	}
	if err := ts.consumeIdentTokenWith("cursive"); err == nil {
		return fonts.Cursive, nil
	}
	if err := ts.consumeIdentTokenWith("fantasy"); err == nil {
		return fonts.Fantasy, nil
	}
	if err := ts.consumeIdentTokenWith("monospace"); err == nil {
		return fonts.Monospace, nil
	}
	return 0, errors.New("invalid generic family")
}

// https://www.w3.org/TR/css-fonts-3/#font-family-prop
func (ts *tokenStream) parseFontFamily() (fonts.FamilyList, error) {
	familyPtrs, err := parseCommaSeparatedRepeation(ts, 0, func(ts *tokenStream) (*fonts.Family, error) {
		if tp, err := ts.parseGenericFamily(); err == nil {
			return &fonts.Family{Type: tp}, nil
		}
		if name, err := ts.parseFamilyName(); err == nil {
			return &fonts.Family{Type: fonts.NonGeneric, Name: name}, nil
		}
		return nil, errors.New("expected generic family or family name")
	})
	if err != nil {
		return fonts.FamilyList{}, err
	}
	families := []fonts.Family{}
	for _, f := range familyPtrs {
		families = append(families, *f)
	}
	return fonts.FamilyList{Families: families}, nil
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
func (ts *tokenStream) parseFontWeight() (fonts.Weight, error) {
	if err := ts.consumeIdentTokenWith("normal"); err == nil {
		return fonts.NormalWeight, nil
	}
	if err := ts.consumeIdentTokenWith("bold"); err == nil {
		return fonts.Bold, nil
	}
	if n := ts.parseNumber(); n != nil {
		if n.Type == css.NumTypeFloat {
			return 0, errors.New("floating point isn't accepted by font-weight")
		}
		intVal := n.ToInt()
		if intVal < 0 || 1000 < intVal {
			return 0, errors.New("font-weight value is out of range")
		}
		return fonts.Weight(n.ToInt()), nil
	}
	return 0, errors.New("invalid font-weight value")
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-stretch
func (ts *tokenStream) parseFontStretch() (fonts.Stretch, error) {
	if err := ts.consumeIdentTokenWith("ultra-condensed"); err == nil {
		return fonts.UltraCondensed, nil
	}
	if err := ts.consumeIdentTokenWith("extra-condensed"); err == nil {
		return fonts.ExtraCondensed, nil
	}
	if err := ts.consumeIdentTokenWith("condensed"); err == nil {
		return fonts.Condensed, nil
	}
	if err := ts.consumeIdentTokenWith("semi-condensed"); err == nil {
		return fonts.SemiCondensed, nil
	}
	if err := ts.consumeIdentTokenWith("normal"); err == nil {
		return fonts.NormalStretch, nil
	}
	if err := ts.consumeIdentTokenWith("semi-expanded"); err == nil {
		return fonts.SemiExpanded, nil
	}
	if err := ts.consumeIdentTokenWith("expanded"); err == nil {
		return fonts.Expanded, nil
	}
	if err := ts.consumeIdentTokenWith("extra-expanded"); err == nil {
		return fonts.ExtraExpanded, nil
	}
	if err := ts.consumeIdentTokenWith("ultra-expanded"); err == nil {
		return fonts.UltraExpanded, nil
	}
	return 0, errors.New("invalid font-stretch value")
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-style
func (ts *tokenStream) parseFontStyle() (fonts.Style, error) {
	if err := ts.consumeIdentTokenWith("normal"); err == nil {
		return fonts.NormalStyle, nil
	}
	if err := ts.consumeIdentTokenWith("italic"); err == nil {
		return fonts.Italic, nil
	}
	if err := ts.consumeIdentTokenWith("oblique"); err == nil {
		return fonts.Oblique, nil
	}
	return 0, errors.New("invalid font-style value")
}

// https://www.w3.org/TR/css-fonts-3/#absolute-size-value
func (ts *tokenStream) parseAbsoluteSize() (fonts.AbsoluteSize, error) {
	if err := ts.consumeIdentTokenWith("xx-small"); err == nil {
		return fonts.XXSmall, nil
	}
	if err := ts.consumeIdentTokenWith("x-small"); err == nil {
		return fonts.XSmall, nil
	}
	if err := ts.consumeIdentTokenWith("small"); err == nil {
		return fonts.Small, nil
	}
	if err := ts.consumeIdentTokenWith("medium"); err == nil {
		return fonts.MediumSize, nil
	}
	if err := ts.consumeIdentTokenWith("large"); err == nil {
		return fonts.Large, nil
	}
	if err := ts.consumeIdentTokenWith("x-large"); err == nil {
		return fonts.XLarge, nil
	}
	if err := ts.consumeIdentTokenWith("xx-large"); err == nil {
		return fonts.XXLarge, nil
	}
	return 0, errors.New("invalid absolute-size value")
}

// https://www.w3.org/TR/css-fonts-3/#relative-size-value
func (ts *tokenStream) parseRelativeSize() (fonts.RelativeSize, error) {
	if err := ts.consumeIdentTokenWith("larger"); err == nil {
		return fonts.Larger, nil
	} else if err := ts.consumeIdentTokenWith("smaller"); err == nil {
		return fonts.Smaller, nil
	}
	return 0, errors.New("invalid relative-size value")
}

// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
func (ts *tokenStream) parseFontSize() (fonts.Size, error) {
	if sz, err := ts.parseAbsoluteSize(); err == nil {
		return sz, nil
	}
	if sz, err := ts.parseRelativeSize(); err == nil {
		return sz, nil
	}
	if l, err := ts.parseLengthOrPercentage(true); err == nil {
		return fonts.LengthFontSize{LengthResolvable: l}, nil
	}
	return nil, errors.New("invalid font-size value")
}
