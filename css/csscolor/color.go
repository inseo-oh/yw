// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

// Package box provide types and values for [CSS Color Module Level 4].
//
// [CSS Color Module Level 4]: https://www.w3.org/TR/css-color-4/
package csscolor

import (
	"fmt"
	"image/color"
	"log"

	"github.com/inseo-oh/yw/css"
)

// NamedColors is map of [CSS named colors]
//
// [CSS named colors]: https://www.w3.org/TR/css-color-4/#named-colors
var NamedColors = map[string]color.Color{
	"aliceblue":            color.RGBA{240, 248, 255, 255},
	"antiquewhite":         color.RGBA{250, 235, 215, 255},
	"aqua":                 color.RGBA{0, 255, 255, 255},
	"aquamarine":           color.RGBA{127, 255, 212, 255},
	"azure":                color.RGBA{240, 255, 255, 255},
	"beige":                color.RGBA{245, 245, 220, 255},
	"bisque":               color.RGBA{255, 228, 196, 255},
	"black":                color.RGBA{0, 0, 0, 255},
	"blanchedalmond":       color.RGBA{255, 235, 205, 255},
	"blue":                 color.RGBA{0, 0, 255, 255},
	"blueviolet":           color.RGBA{138, 43, 226, 255},
	"brown":                color.RGBA{165, 42, 42, 255},
	"burlywood":            color.RGBA{222, 184, 135, 255},
	"cadetblue":            color.RGBA{95, 158, 160, 255},
	"chartreuse":           color.RGBA{127, 255, 0, 255},
	"chocolate":            color.RGBA{210, 105, 30, 255},
	"coral":                color.RGBA{255, 127, 80, 255},
	"cornflowerblue":       color.RGBA{100, 149, 237, 255},
	"cornsilk":             color.RGBA{255, 248, 220, 255},
	"crimson":              color.RGBA{220, 20, 60, 255},
	"cyan":                 color.RGBA{0, 255, 255, 255},
	"darkblue":             color.RGBA{0, 0, 139, 255},
	"darkcyan":             color.RGBA{0, 139, 139, 255},
	"darkgoldenrod":        color.RGBA{184, 134, 11, 255},
	"darkgray":             color.RGBA{169, 169, 169, 255},
	"darkgreen":            color.RGBA{0, 100, 0, 255},
	"darkgrey":             color.RGBA{169, 169, 169, 255},
	"darkkhaki":            color.RGBA{189, 183, 107, 255},
	"darkmagenta":          color.RGBA{139, 0, 139, 255},
	"darkolivegreen":       color.RGBA{85, 107, 47, 255},
	"darkorange":           color.RGBA{255, 140, 0, 255},
	"darkorchid":           color.RGBA{153, 50, 204, 255},
	"darkred":              color.RGBA{139, 0, 0, 255},
	"darksalmon":           color.RGBA{233, 150, 122, 255},
	"darkseagreen":         color.RGBA{143, 188, 143, 255},
	"darkslateblue":        color.RGBA{72, 61, 139, 255},
	"darkslategray":        color.RGBA{47, 79, 79, 255},
	"darkslategrey":        color.RGBA{47, 79, 79, 255},
	"darkturquoise":        color.RGBA{0, 206, 209, 255},
	"darkviolet":           color.RGBA{148, 0, 211, 255},
	"deeppink":             color.RGBA{255, 20, 147, 255},
	"deepskyblue":          color.RGBA{0, 191, 255, 255},
	"dimgray":              color.RGBA{105, 105, 105, 255},
	"dimgrey":              color.RGBA{105, 105, 105, 255},
	"dodgerblue":           color.RGBA{30, 144, 255, 255},
	"firebrick":            color.RGBA{178, 34, 34, 255},
	"floralwhite":          color.RGBA{255, 250, 240, 255},
	"forestgreen":          color.RGBA{34, 139, 34, 255},
	"fuchsia":              color.RGBA{255, 0, 255, 255},
	"gainsboro":            color.RGBA{220, 220, 220, 255},
	"ghostwhite":           color.RGBA{248, 248, 255, 255},
	"gold":                 color.RGBA{255, 215, 0, 255},
	"goldenrod":            color.RGBA{218, 165, 32, 255},
	"gray":                 color.RGBA{128, 128, 128, 255},
	"green":                color.RGBA{0, 128, 0, 255},
	"greenyellow":          color.RGBA{173, 255, 47, 255},
	"grey":                 color.RGBA{128, 128, 128, 255},
	"honeydew":             color.RGBA{240, 255, 240, 255},
	"hotpink":              color.RGBA{255, 105, 180, 255},
	"indianred":            color.RGBA{205, 92, 92, 255},
	"indigo":               color.RGBA{75, 0, 130, 255},
	"ivory":                color.RGBA{255, 255, 240, 255},
	"khaki":                color.RGBA{240, 230, 140, 255},
	"lavender":             color.RGBA{230, 230, 250, 255},
	"lavenderblush":        color.RGBA{255, 240, 245, 255},
	"lawngreen":            color.RGBA{124, 252, 0, 255},
	"lemonchiffon":         color.RGBA{255, 250, 205, 255},
	"lightblue":            color.RGBA{173, 216, 230, 255},
	"lightcoral":           color.RGBA{240, 128, 128, 255},
	"lightcyan":            color.RGBA{224, 255, 255, 255},
	"lightgoldenrodyellow": color.RGBA{250, 250, 210, 255},
	"lightgray":            color.RGBA{211, 211, 211, 255},
	"lightgreen":           color.RGBA{144, 238, 144, 255},
	"lightgrey":            color.RGBA{211, 211, 211, 255},
	"lightpink":            color.RGBA{255, 182, 193, 255},
	"lightsalmon":          color.RGBA{255, 160, 122, 255},
	"lightseagreen":        color.RGBA{32, 178, 170, 255},
	"lightskyblue":         color.RGBA{135, 206, 250, 255},
	"lightslategray":       color.RGBA{119, 136, 153, 255},
	"lightslategrey":       color.RGBA{119, 136, 153, 255},
	"lightsteelblue":       color.RGBA{176, 196, 222, 255},
	"lightyellow":          color.RGBA{255, 255, 224, 255},
	"lime":                 color.RGBA{0, 255, 0, 255},
	"limegreen":            color.RGBA{50, 205, 50, 255},
	"linen":                color.RGBA{250, 240, 230, 255},
	"magenta":              color.RGBA{255, 0, 255, 255},
	"maroon":               color.RGBA{128, 0, 0, 255},
	"mediumaquamarine":     color.RGBA{102, 205, 170, 255},
	"mediumblue":           color.RGBA{0, 0, 205, 255},
	"mediumorchid":         color.RGBA{186, 85, 211, 255},
	"mediumpurple":         color.RGBA{147, 112, 219, 255},
	"mediumseagreen":       color.RGBA{60, 179, 113, 255},
	"mediumslateblue":      color.RGBA{123, 104, 238, 255},
	"mediumspringgreen":    color.RGBA{0, 250, 154, 255},
	"mediumturquoise":      color.RGBA{72, 209, 204, 255},
	"mediumvioletred":      color.RGBA{199, 21, 133, 255},
	"midnightblue":         color.RGBA{25, 25, 112, 255},
	"mintcream":            color.RGBA{245, 255, 250, 255},
	"mistyrose":            color.RGBA{255, 228, 225, 255},
	"moccasin":             color.RGBA{255, 228, 181, 255},
	"navajowhite":          color.RGBA{255, 222, 173, 255},
	"navy":                 color.RGBA{0, 0, 128, 255},
	"oldlace":              color.RGBA{253, 245, 230, 255},
	"olive":                color.RGBA{128, 128, 0, 255},
	"olivedrab":            color.RGBA{107, 142, 35, 255},
	"orange":               color.RGBA{255, 165, 0, 255},
	"orangered":            color.RGBA{255, 69, 0, 255},
	"orchid":               color.RGBA{218, 112, 214, 255},
	"palegoldenrod":        color.RGBA{238, 232, 170, 255},
	"palegreen":            color.RGBA{152, 251, 152, 255},
	"paleturquoise":        color.RGBA{175, 238, 238, 255},
	"palevioletred":        color.RGBA{219, 112, 147, 255},
	"papayawhip":           color.RGBA{255, 239, 213, 255},
	"peachpuff":            color.RGBA{255, 218, 185, 255},
	"peru":                 color.RGBA{205, 133, 63, 255},
	"pink":                 color.RGBA{255, 192, 203, 255},
	"plum":                 color.RGBA{221, 160, 221, 255},
	"powderblue":           color.RGBA{176, 224, 230, 255},
	"purple":               color.RGBA{128, 0, 128, 255},
	"rebeccapurple":        color.RGBA{102, 51, 153, 255},
	"red":                  color.RGBA{255, 0, 0, 255},
	"rosybrown":            color.RGBA{188, 143, 143, 255},
	"royalblue":            color.RGBA{65, 105, 225, 255},
	"saddlebrown":          color.RGBA{139, 69, 19, 255},
	"salmon":               color.RGBA{250, 128, 114, 255},
	"sandybrown":           color.RGBA{244, 164, 96, 255},
	"seagreen":             color.RGBA{46, 139, 87, 255},
	"seashell":             color.RGBA{255, 245, 238, 255},
	"sienna":               color.RGBA{160, 82, 45, 255},
	"silver":               color.RGBA{192, 192, 192, 255},
	"skyblue":              color.RGBA{135, 206, 235, 255},
	"slateblue":            color.RGBA{106, 90, 205, 255},
	"slategray":            color.RGBA{112, 128, 144, 255},
	"slategrey":            color.RGBA{112, 128, 144, 255},
	"snow":                 color.RGBA{255, 250, 250, 255},
	"springgreen":          color.RGBA{0, 255, 127, 255},
	"steelblue":            color.RGBA{70, 130, 180, 255},
	"tan":                  color.RGBA{210, 180, 140, 255},
	"teal":                 color.RGBA{0, 128, 128, 255},
	"thistle":              color.RGBA{216, 191, 216, 255},
	"tomato":               color.RGBA{255, 99, 71, 255},
	"turquoise":            color.RGBA{64, 224, 208, 255},
	"violet":               color.RGBA{238, 130, 238, 255},
	"wheat":                color.RGBA{245, 222, 179, 255},
	"white":                color.RGBA{255, 255, 255, 255},
	"whitesmoke":           color.RGBA{245, 245, 245, 255},
	"yellow":               color.RGBA{255, 255, 0, 255},
	"yellowgreen":          color.RGBA{154, 205, 50, 255},
}

// CSS color value
type Color struct {
	Type       Type
	Components [4]css.Num
}

// Type of Color
type Type uint8

const (
	Rgb          Type = iota // rgb(), rgba(), hex colors, named colors
	CurrentColor             // currentColor
	Hsl                      // hsl(), hsla()
	Hwb                      // hwb()
	Lab                      // lab()
	Lch                      // lch()
	Oklab                    // oklab()
	Oklch                    // oklch()
	ColorFn                  // color()
	// TODO: System color
)

// Various predefined CSS color values
var (
	CanvasText  = FromStdColor(color.RGBA{0, 0, 0, 255}) // CanvasText
	Transparent = FromStdColor(color.RGBA{0, 0, 0, 0})   // transparent
)

// FromStdColor creates [Color] from [color.Color].
func FromStdColor(col color.Color) Color {
	r, g, b, a := col.RGBA()
	return Color{Rgb, [4]css.Num{
		css.NumFromInt(int64(r)),
		css.NumFromInt(int64(g)),
		css.NumFromInt(int64(b)),
		css.NumFromInt(int64(a)),
	}}
}

// ToStdColor returns a [color.Color] for the color.
// currentColor may be nil if the color is never going to be [CurrentColor].
//
// TODO(ois): Currently only RGBA and currentColor type is supported. Support more color types.
func (c Color) ToStdColor(currentColor color.Color) color.Color {
	switch c.Type {
	case Rgb:
		return color.RGBA{
			uint8(c.Components[0].ToInt() >> 8),
			uint8(c.Components[1].ToInt() >> 8),
			uint8(c.Components[2].ToInt() >> 8),
			uint8(c.Components[3].ToInt() >> 8),
		}
	case CurrentColor:
		return currentColor
	case Hsl:
		panic("TODO")
	case Hwb:
		panic("TODO")
	case Lab:
		panic("TODO")
	case Lch:
		panic("TODO")
	case Oklab:
		panic("TODO")
	case Oklch:
		panic("TODO")
	case ColorFn:
		panic("TODO")
	}
	log.Panicf("<bad Color type %v>", c.Type)
	panic("unreachable")
}

// Equals reports whether two Color values are equal or not.
//
// Note that colors with different types are considered as non-equal.
func (c Color) Equals(other Color) bool {
	if c.Type != other.Type {
		return false
	}
	if len(c.Components) != len(other.Components) {
		return false
	}
	for i := range len(c.Components) {
		if !c.Components[i].Equals(other.Components[i]) {
			return false
		}
	}
	return true
}
func (c Color) String() string {
	switch c.Type {
	case Rgb:
		return fmt.Sprintf("#%02x%02x%02x%02x", c.Components[0].ToInt()>>8, c.Components[1].ToInt()>>8, c.Components[2].ToInt()>>8, c.Components[3].ToInt()>>8)
	case CurrentColor:
		return "currentColor"
	case Hsl:
		return fmt.Sprintf("hsl(%v, %v, %v, %v)", c.Components[0], c.Components[1], c.Components[2], c.Components[3]) // STUB
	case Hwb:
		return fmt.Sprintf("hwb(%v, %v, %v, %v)", c.Components[0], c.Components[1], c.Components[2], c.Components[3]) // STUB
	case Lab:
		return fmt.Sprintf("lab(%v, %v, %v, %v)", c.Components[0], c.Components[1], c.Components[2], c.Components[3]) // STUB
	case Lch:
		return fmt.Sprintf("lch(%v, %v, %v, %v)", c.Components[0], c.Components[1], c.Components[2], c.Components[3]) // STUB
	case Oklab:
		return fmt.Sprintf("oklab(%v, %v, %v, %v)", c.Components[0], c.Components[1], c.Components[2], c.Components[3]) // STUB
	case Oklch:
		return fmt.Sprintf("oklch(%v, %v, %v, %v)", c.Components[0], c.Components[1], c.Components[2], c.Components[3]) // STUB
	case ColorFn:
		return "color(...)" // STUB
	}
	return fmt.Sprintf("<bad Color type %v>", c.Type)
}
