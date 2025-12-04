// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/csscolor"
	"github.com/inseo-oh/yw/css/values"
)

func (ts *tokenStream) parseColor() (csscolor.Color, error) {
	oldCursor := ts.cursor
	// Try hex notation --------------------------------------------------------
	if tk, err := ts.consumeTokenWith(tokenTypeHash); err == nil {
		// https://www.w3.org/TR/css-color-4/#hex-notation
		chrs := []rune(tk.(hashToken).value)
		var rStr, gStr, bStr, aStr string
		aStr = "ff"
		switch len(chrs) {
		case 3:
			// #rgb
			rStr = strings.Repeat(string(chrs[0]), 2)
			gStr = strings.Repeat(string(chrs[1]), 2)
			bStr = strings.Repeat(string(chrs[2]), 2)
		case 4:
			// #rgba
			rStr = strings.Repeat(string(chrs[0]), 2)
			gStr = strings.Repeat(string(chrs[1]), 2)
			bStr = strings.Repeat(string(chrs[2]), 2)
			aStr = strings.Repeat(string(chrs[3]), 2)
		case 6:
			// #rrggbb
			rStr = string(chrs[0:2])
			gStr = string(chrs[2:4])
			bStr = string(chrs[4:6])
		case 8:
			// #rrggbb
			rStr = string(chrs[0:2])
			gStr = string(chrs[2:4])
			bStr = string(chrs[4:6])
			aStr = string(chrs[6:8])
		default:
			return csscolor.Color{}, fmt.Errorf("invalid hex digit length: %s", tk.(hashToken).value)
		}
		r, err := strconv.ParseUint(rStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, err
		}
		g, err := strconv.ParseUint(gStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, err
		}
		b, err := strconv.ParseUint(bStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, err
		}
		a, err := strconv.ParseUint(aStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, err
		}
		return csscolor.Color{Type: csscolor.Rgb, Components: [4]css.Num{
			css.NumFromInt(int64(r)),
			css.NumFromInt(int64(g)),
			css.NumFromInt(int64(b)),
			css.NumFromInt(int64(a)),
		}}, nil
	} else {
		ts.cursor = oldCursor
	}
	// Try rgb()/rgba() function -----------------------------------------------
	fn, err := ts.consumeAstFuncWith("rgb")
	if err != nil {
		// rgba() is alias for rgb()
		fn, err = ts.consumeAstFuncWith("rgba")
	}
	if err == nil {
		ts := tokenStream{tokens: fn.value, tokenizerHelper: ts.tokenizerHelper}

		parseAlpha := func() (css.Num, error) {
			var v css.Num
			if num := ts.parseNumber(); num != nil {
				v = css.NumFromFloat(num.Clamp(css.NumFromInt(0), css.NumFromInt(1)).ToFloat() * 255)
			} else if num, err := ts.parsePercentage(); err == nil {
				aPer := num.Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
				v = css.NumFromFloat((aPer / 100) * 255)
			} else {
				return v, fmt.Errorf("%s: expected number or percentage", ts.errorHeader())
			}
			return v, nil
		}

		// https://www.w3.org/TR/css-color-4/#funcdef-rgb
		var r, g, b, a css.Num
		a = css.NumFromInt(255)
		oldCursor := ts.cursor

		//======================================================================
		// Try legacy syntax first
		//======================================================================
		// https://www.w3.org/TR/css-color-4/#typedef-legacy-rgb-syntax

		// rgb(<  >r  ,  g  ,  b  ) --------------------------------------------
		// rgb(<  >r  ,  g  ,  b  ,  a  ) --------------------------------------
		ts.skipWhitespaces()
		// rgb(  <r  ,  g  ,  b>  ) --------------------------------------------
		// rgb(  <r  ,  g  ,  b>  ,  a  ) --------------------------------------
		per, err := parseCommaSeparatedRepeation(&ts, 3, "percentage", func(ts *tokenStream) (values.Percentage, error) {
			return ts.parsePercentage()
		})
		if err != nil {
			return csscolor.Color{}, err
		} else if len(per) == 3 {
			// Percentage value
			rPer := per[0].Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
			gPer := per[1].Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
			bBer := per[2].Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
			r = css.NumFromFloat((rPer / 100) * 255)
			g = css.NumFromFloat((gPer / 100) * 255)
			b = css.NumFromFloat((bBer / 100) * 255)
		} else {
			num, err := parseCommaSeparatedRepeation(&ts, 3, "number", func(ts *tokenStream) (*css.Num, error) {
				n := ts.parseNumber()
				if n == nil {
					return nil, fmt.Errorf("%s: expected number", ts.errorHeader())
				}
				return ts.parseNumber(), nil
			})
			if err != nil {
				return csscolor.Color{}, err
			} else if len(num) == 3 {
				// Number value
				r = num[0].Clamp(css.NumFromInt(0), css.NumFromInt(255))
				g = num[1].Clamp(css.NumFromInt(0), css.NumFromInt(255))
				b = num[2].Clamp(css.NumFromInt(0), css.NumFromInt(255))
			} else {
				goto modernSyntax
			}
		}
		// rgb(  r  ,  g  ,  b<  >) --------------------------------------------
		// rgb(  r  ,  g  ,  b<  >,  a  ) --------------------------------------
		ts.skipWhitespaces()
		// rgb(  r  ,  g  ,  b  <,>  a  ) --------------------------------------
		if _, err := ts.consumeTokenWith(tokenTypeComma); err == nil {
			// rgb(  r  ,  g  ,  b  ,<  >a  ) ----------------------------------
			ts.skipWhitespaces()
			// rgb(  r  ,  g  ,  b  ,  <a>  ) ----------------------------------
			if v, err := parseAlpha(); err == nil {
				a = v
			} else {
				return csscolor.Color{}, err
			}
			// rgb(  r  ,  g  ,  b  ,  a<  >) ----------------------------------
			ts.skipWhitespaces()
		} else {
			ts.cursor = oldCursor
		}
		if !ts.isEnd() {
			return csscolor.Color{}, fmt.Errorf("%s: expected end", ts.errorHeader())
		}
		return csscolor.Color{Type: csscolor.Rgb, Components: [4]css.Num{r, g, b, a}}, nil
	modernSyntax:
		ts.cursor = oldCursor

		//======================================================================
		// Try modern syntax
		//======================================================================
		// https://www.w3.org/TR/css-color-4/#typedef-modern-rgb-syntax

		// rgb(<  >r  g  b  ) --------------------------------------------------
		// rgb(<  >r  g  b  /  a  ) --------------------------------------------
		ts.skipWhitespaces()
		// rgb(  <r  g  b  >) --------------------------------------------------
		// rgb(  <r  g  b  >/  a  ) --------------------------------------------
		components := []css.Num{}
		for range 3 {
			// rgb(  <r>  <g>  <b>  ) ------------------------------------------
			// rgb(  <r>  <g>  <b>  /  a  ) ------------------------------------
			var v css.Num
			if num := ts.parseNumber(); num != nil {
				v = num.Clamp(css.NumFromInt(0), css.NumFromInt(255))
			} else if num, err := ts.parsePercentage(); err == nil {
				per := num.Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
				v = css.NumFromFloat((per / 100) * 255)
			} else if err := ts.consumeIdentTokenWith("none"); err == nil {
				panic("TODO")
			} else {
				return csscolor.Color{}, fmt.Errorf("%s: expected number or percentage", ts.errorHeader())
			}
			components = append(components, v)
			// rgb(  r<  >g<  >b<  >) ------------------------------------------
			// rgb(  r<  >g<  >b<  >/  a  ) ------------------------------------
			ts.skipWhitespaces()
		}
		// rgb(  r  g  b  </>  a  ) --------------------------------------------
		a = css.NumFromInt(255)
		if err := ts.consumeDelimTokenWith('/'); err == nil {
			// rgb(  r  g  b  /<  >a  ) --------------------------------------------
			ts.skipWhitespaces()
			// rgb(  r  g  b  /  <a>  ) --------------------------------------------
			if v, err := parseAlpha(); err == nil {
				a = v
			} else {
				return csscolor.Color{}, err
			}
			// rgb(  r  g  b  /  a<  >) --------------------------------------------
			ts.skipWhitespaces()
		}
		if !ts.isEnd() {
			return csscolor.Color{}, fmt.Errorf("%s: expected end", ts.errorHeader())
		}
		return csscolor.Color{Type: csscolor.Rgb, Components: [4]css.Num{components[0], components[1], components[2], a}}, nil
	}
	ts.cursor = oldCursor
	// Try hsl()/hsla() function -----------------------------------------------
	fn, err = ts.consumeAstFuncWith("hsl")
	if err != nil {
		// hsla() is alias for hsl()
		fn, err = ts.consumeAstFuncWith("hsl")
	}
	if err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-hsl
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-hsl]")
	}
	// Try hwb() function ------------------------------------------------------
	if _, err := ts.consumeAstFuncWith("hwb"); err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-hwb
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-hwb]")
	}
	// Try lab() function ------------------------------------------------------
	if _, err := ts.consumeAstFuncWith("hwb"); err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-lab
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-lab]")
	}
	// Try lch() function ------------------------------------------------------
	if _, err := ts.consumeAstFuncWith("lch"); err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-lch
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-lch]")
	}
	// Try oklab() function ----------------------------------------------------
	if _, err := ts.consumeAstFuncWith("oklab"); err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-oklab
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-oklab]")
	}
	// Try oklch() function ----------------------------------------------------
	if _, err := ts.consumeAstFuncWith("oklch"); err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-oklch
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-oklch]")
	}
	// Try color() function ----------------------------------------------------
	if _, err := ts.consumeAstFuncWith("color"); err == nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-color
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-color]")
	}
	// Try named color ---------------------------------------------------------
	ident, err := ts.consumeTokenWith(tokenTypeIdent)
	if err == nil {
		rgba, ok := csscolor.NamedColors[ident.(identToken).value]
		if ok {
			return csscolor.Color{Type: csscolor.Rgb, Components: [4]css.Num{
				css.NumFromInt(int64(rgba.R)),
				css.NumFromInt(int64(rgba.G)),
				css.NumFromInt(int64(rgba.B)),
				css.NumFromInt(int64(rgba.A)),
			}}, nil
		}
	} else {
		ts.cursor = oldCursor
	}
	// Try transparent ---------------------------------------------------------
	if err := ts.consumeIdentTokenWith("transparent"); err == nil {
		c := csscolor.Transparent
		return c, nil
	}
	ts.cursor = oldCursor
	// Try currentColor --------------------------------------------------------
	if err := ts.consumeIdentTokenWith("currentColor"); err == nil {
		return csscolor.Color{Type: csscolor.CurrentColor}, nil
	}
	ts.cursor = oldCursor
	// TODO: Try system colors
	return csscolor.Color{}, fmt.Errorf("%s expected color", ts.errorHeader())
}
