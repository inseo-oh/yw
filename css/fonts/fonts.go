// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Package display provides types and values for [CSS Fonts Module Level 3]
//
// [CSS Fonts Module Level 3]: https://www.w3.org/TR/css-fonts-3
package fonts

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/values"
)

// Family represents single entry in [CSS font-family] property.
//
// [CSS font-family]: https://www.w3.org/TR/css-fonts-3/#propdef-font-family
type Family struct {
	Type FamilyType // Font family type
	Name string     // Only valid if the Type is NonGeneric
}

// FamilyType represents one of generic families, or NonGeneric, meaning Name should be used instead.
type FamilyType uint8

const (
	// https://www.w3.org/TR/css-fonts-3/#generic-family-value
	Serif      FamilyType = iota // font-family: serif
	SansSerif                    // font-family: sans-serif
	Cursive                      // font-family: cursive
	Fantasy                      // font-family: fantasy
	Monospace                    // font-family: monospace
	NonGeneric                   // font-family: <anything other than above>
)

func (f Family) String() string {
	switch f.Type {
	case Serif:
		return "serif"
	case SansSerif:
		return "sans-serif"
	case Cursive:
		return "cursive"
	case Fantasy:
		return "fantasy"
	case Monospace:
		return "monospace"
	case NonGeneric:
		return strconv.Quote(f.Name)
	}
	return fmt.Sprintf("<bad Family type %d>", f.Type)
}

// FamilyList represents value of [CSS font-family] property.
//
// [CSS font-family]: https://www.w3.org/TR/css-fonts-3/#propdef-font-family
type FamilyList struct {
	Families []Family
}

func (f FamilyList) String() string {
	sb := strings.Builder{}
	for i, f := range f.Families {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%v", f))
	}
	return sb.String()
}

// Weight represents value of [CSS font-weight] property.
//
// [CSS font-weight]: https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
type Weight uint16

// Some predefined font-weight values
const (
	NormalWeight Weight = 400 // font-weight: normal
	Bold         Weight = 700 // font-weight: bold
)

func (w Weight) String() string {
	return fmt.Sprintf("%d", w)
}

// Stretch represents value of [CSS font-stretch] property.
//
// [CSS font-stretch]: https://www.w3.org/TR/css-fonts-3/#propdef-font-stretch
type Stretch uint8

const (
	UltraCondensed Stretch = iota // font-stretch: ultra-condensed
	ExtraCondensed                // font-stretch: extra-condensed
	Condensed                     // font-stretch: condensed
	SemiCondensed                 // font-stretch: semi-condensed
	NormalStretch                 // font-stretch: normal
	SemiExpanded                  // font-stretch: semi-expanded
	Expanded                      // font-stretch: expanded
	ExtraExpanded                 // font-stretch: extra-expanded
	UltraExpanded                 // font-stretch: ultra-expanded
)

func (s Stretch) String() string {
	switch s {
	case UltraCondensed:
		return "ultra-condensed"
	case ExtraCondensed:
		return "extra-condensed"
	case Condensed:
		return "condensed"
	case SemiCondensed:
		return "semi-condensed"
	case NormalStretch:
		return "normal"
	case SemiExpanded:
		return "semi-expanded"
	case Expanded:
		return "expanded"
	case ExtraExpanded:
		return "extra-expanded"
	case UltraExpanded:
		return "ultra-expanded"
	}
	return fmt.Sprintf("<bad Stretch %d>", s)
}

// Style represents value of [CSS font-style] property.
//
// [CSS font-style]: https://www.w3.org/TR/css-fonts-3/#font-style-prop
type Style uint8

const (
	NormalStyle Style = iota // font-style: normal
	Italic                   // font-style: italic
	Oblique                  // font-style: oblique
)

func (s Style) String() string {
	switch s {
	case NormalStyle:
		return "normal"
	case Italic:
		return "italic"
	case Oblique:
		return "oblique"
	}
	return fmt.Sprintf("<bad Style %d>", s)
}

// Size represents value of [CSS font-size] property.
//
// [CSS font-size]: https://www.w3.org/TR/css-fonts-3/#propdef-font-size
type Size interface {
	// CalculateRealFontSize calculates real size.
	CalculateRealFontSize(parentFontSize css.Num) float64
	String() string
}

const (
	PreferredFontSize = 14 // XXX: Let user choose this size!
)

// AbsoluteSize represents [CSS absolute-size] value.
// Exact pixel sizes are determined based on [PreferredFontSize].
//
// [CSS <absolute-size>]: https://www.w3.org/TR/css-fonts-3/#absolute-size-value
type AbsoluteSize uint8

const (
	XXSmall    AbsoluteSize = iota // font-size: xx-small
	XSmall                         // font-size: x-small
	Small                          // font-size: small
	MediumSize                     // font-size: medium
	Large                          // font-size: large
	XLarge                         // font-size: x-large
	XXLarge                        // font-size: xx-large
)

func (s AbsoluteSize) String() string {
	switch s {
	case XXSmall:
		return "xx-small"
	case XSmall:
		return "x-small"
	case Small:
		return "small"
	case MediumSize:
		return "medium"
	case Large:
		return "large"
	case XLarge:
		return "x-large"
	case XXLarge:
		return "xx-large"
	}
	return fmt.Sprintf("<unknown AbsoluteSize value %d>", s)
}

// https://www.w3.org/TR/css-fonts-3/#absolute-size-value
var absoluteSizeMap = map[AbsoluteSize]float64{
	XXSmall:    (PreferredFontSize * 3) / 5,
	XSmall:     (PreferredFontSize * 3) / 4,
	Small:      (PreferredFontSize * 8) / 9,
	MediumSize: PreferredFontSize,
	Large:      (PreferredFontSize * 6) / 5,
	XLarge:     (PreferredFontSize * 3) / 2,
	XXLarge:    (PreferredFontSize * 2) / 1,
}

func (s AbsoluteSize) CalculateRealFontSize(parentFontSize css.Num) float64 {
	return absoluteSizeMap[s]
}

func pxToAbsoluteSize(size css.Num) AbsoluteSize {
	sizeNum := size.ToFloat()
	minDiff := float64(^0)
	var resSize AbsoluteSize
	for k, v := range absoluteSizeMap {
		diff := math.Abs(sizeNum - v)
		if diff < minDiff {
			resSize = k
			minDiff = diff
		}
	}
	return resSize
}

// RelativeSize represents [CSS relative-size] value.
// Font size is converted to [AbsoluteSize] first, and then returns
// next larger/smaller size for it.
//
// [CSS relative-size]: https://www.w3.org/TR/css-fonts-3/#relative-size-value
type RelativeSize uint8

const (
	Larger  RelativeSize = iota // font-size: larger
	Smaller                     // font-size: smaller
)

func (s RelativeSize) String() string {
	switch s {
	case Larger:
		return "larger"
	case Smaller:
		return "smaller"
	}
	return fmt.Sprintf("<bad RelativeSize %d>", s)
}

func (s RelativeSize) CalculateRealFontSize(parentFontSize css.Num) float64 {
	parentAbsSize := pxToAbsoluteSize(parentFontSize)
	var resultAbsSize AbsoluteSize
	switch s {
	case Larger:
		switch parentAbsSize {
		case XXSmall:
			resultAbsSize = XSmall
		case XSmall:
			resultAbsSize = Small
		case Small:
			resultAbsSize = MediumSize
		case MediumSize:
			resultAbsSize = Large
		case Large:
			resultAbsSize = XLarge
		case XLarge:
			resultAbsSize = XXLarge
		case XXLarge:
			resultAbsSize = XXLarge
		default:
			log.Panicf("<bad AbsoluteSize value %d>", s)
		}
		return resultAbsSize.CalculateRealFontSize(parentFontSize)
	case Smaller:
		switch parentAbsSize {
		case XXSmall:
			resultAbsSize = XXSmall
		case XSmall:
			resultAbsSize = XXSmall
		case Small:
			resultAbsSize = XSmall
		case MediumSize:
			resultAbsSize = Small
		case Large:
			resultAbsSize = MediumSize
		case XLarge:
			resultAbsSize = Large
		case XXLarge:
			resultAbsSize = XLarge
		default:
			log.Panicf("<bad AbsoluteSize %d>", s)
		}
		return resultAbsSize.CalculateRealFontSize(parentFontSize)
	}
	log.Panicf("<bad RelativeSize %d>", s)
	return 0
}

// LengthFontSize represents font size specified directly using [values.LengthResolvable].
type LengthFontSize struct{ values.LengthResolvable }

func (l LengthFontSize) CalculateRealFontSize(parentFontSize css.Num) float64 {
	return l.AsLength(parentFontSize).ToPx(parentFontSize)
}
