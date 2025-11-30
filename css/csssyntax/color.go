package csssyntax

import (
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/csscolor"
	"github.com/inseo-oh/yw/css/values"
	"github.com/inseo-oh/yw/util"
)

func (ts *tokenStream) parseColor() (csscolor.Color, bool) {
	// Try hex notation --------------------------------------------------------
	if tk := ts.consumeTokenWith(tokenTypeHash); !util.IsNil(tk) {
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
			return csscolor.Color{}, false
		}
		r, err := strconv.ParseUint(rStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, false
		}
		g, err := strconv.ParseUint(gStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, false
		}
		b, err := strconv.ParseUint(bStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, false
		}
		a, err := strconv.ParseUint(aStr, 16, 8)
		if err != nil {
			return csscolor.Color{}, false
		}
		return csscolor.Color{Type: csscolor.TypeRgb, Components: []css.Num{
			css.NumFromInt(int64(r)),
			css.NumFromInt(int64(g)),
			css.NumFromInt(int64(b)),
			css.NumFromInt(int64(a)),
		}}, true
	}
	// Try rgb()/rgba() function -----------------------------------------------
	fn := ts.consumeAstFuncWith("rgb")
	if fn == nil {
		// rgba() is alias for rgb()
		fn = ts.consumeAstFuncWith("rgba")
	}
	if fn != nil {
		ts := tokenStream{tokens: fn.value}

		parseAlpha := func() *css.Num {
			var v css.Num
			if num := ts.parseNumber(); num != nil {
				v = css.NumFromFloat(num.Clamp(css.NumFromInt(0), css.NumFromInt(1)).ToFloat() * 255)
			} else if num := ts.parsePercentage(); num != nil {
				aPer := num.Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
				v = css.NumFromFloat((aPer / 100) * 255)
			} else {
				return nil
			}
			return &v
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
		per, err := parseCommaSeparatedRepeation(&ts, 3, func(ts *tokenStream) (*values.Percentage, error) {
			return ts.parsePercentage(), nil
		})
		if per == nil && err != nil {
			return csscolor.Color{}, false
		} else if len(per) == 3 {
			// Percentage value
			rPer := per[0].Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
			gPer := per[1].Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
			bBer := per[2].Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
			r = css.NumFromFloat((rPer / 100) * 255)
			g = css.NumFromFloat((gPer / 100) * 255)
			b = css.NumFromFloat((bBer / 100) * 255)
		} else {
			num, err := parseCommaSeparatedRepeation(&ts, 3, func(ts *tokenStream) (*css.Num, error) {
				return ts.parseNumber(), nil
			})
			if num == nil && err != nil {
				return csscolor.Color{}, false
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
		if tk := ts.consumeTokenWith(tokenTypeComma); !util.IsNil(tk) {
			// rgb(  r  ,  g  ,  b  ,<  >a  ) ----------------------------------
			ts.skipWhitespaces()
			// rgb(  r  ,  g  ,  b  ,  <a>  ) ----------------------------------
			if v := parseAlpha(); v != nil {
				a = *v
			} else {
				return csscolor.Color{}, false
			}
			// rgb(  r  ,  g  ,  b  ,  a<  >) ----------------------------------
			ts.skipWhitespaces()
		}
		if !ts.isEnd() {
			return csscolor.Color{}, false
		}
		return csscolor.Color{Type: csscolor.TypeRgb, Components: []css.Num{r, g, b, a}}, true
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
			} else if num := ts.parsePercentage(); num != nil {
				per := num.Value.Clamp(css.NumFromInt(0), css.NumFromInt(100)).ToFloat()
				v = css.NumFromFloat((per / 100) * 255)
			} else if tk := ts.consumeIdentTokenWith("none"); !util.IsNil(tk) {
				panic("TODO")
			} else {
				return csscolor.Color{}, false
			}
			components = append(components, v)
			// rgb(  r<  >g<  >b<  >) ------------------------------------------
			// rgb(  r<  >g<  >b<  >/  a  ) ------------------------------------
			ts.skipWhitespaces()
		}
		// rgb(  r  g  b  </>  a  ) --------------------------------------------
		a = css.NumFromInt(255)
		if tk := ts.consumeDelimTokenWith('/'); tk != nil {
			// rgb(  r  g  b  /<  >a  ) --------------------------------------------
			ts.skipWhitespaces()
			// rgb(  r  g  b  /  <a>  ) --------------------------------------------
			if v := parseAlpha(); v != nil {
				a = *v
			} else {
				return csscolor.Color{}, false
			}
			// rgb(  r  g  b  /  a<  >) --------------------------------------------
			ts.skipWhitespaces()
		}
		if !ts.isEnd() {
			return csscolor.Color{}, false
		}
		return csscolor.Color{Type: csscolor.TypeRgb, Components: []css.Num{components[0], components[1], components[2], a}}, true
	}
	// Try hsl()/hsla() function -----------------------------------------------
	fn = ts.consumeAstFuncWith("hsl")
	if fn == nil {
		// hsla() is alias for hsl()
		fn = ts.consumeAstFuncWith("hsl")
	}
	if fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-hsl
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-hsl]")
	}
	// Try hwb() function ------------------------------------------------------
	if fn = ts.consumeAstFuncWith("hwb"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-hwb
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-hwb]")
	}
	// Try lab() function ------------------------------------------------------
	if fn = ts.consumeAstFuncWith("hwb"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-lab
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-lab]")
	}
	// Try lch() function ------------------------------------------------------
	if fn = ts.consumeAstFuncWith("lch"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-lch
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-lch]")
	}
	// Try oklab() function ----------------------------------------------------
	if fn = ts.consumeAstFuncWith("oklab"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-oklab
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-oklab]")
	}
	// Try oklch() function ----------------------------------------------------
	if fn = ts.consumeAstFuncWith("oklch"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-oklch
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-oklch]")
	}
	// Try color() function ----------------------------------------------------
	if fn = ts.consumeAstFuncWith("color"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-color
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-color]")
	}
	// Try named color ---------------------------------------------------------
	oldCursor := ts.cursor
	ident := ts.consumeTokenWith(tokenTypeIdent)
	if !util.IsNil(ident) {
		rgba, ok := csscolor.NamedColors[ident.(identToken).value]
		if ok {
			return csscolor.Color{Type: csscolor.TypeRgb, Components: []css.Num{
				css.NumFromInt(int64(rgba.R)),
				css.NumFromInt(int64(rgba.G)),
				css.NumFromInt(int64(rgba.B)),
				css.NumFromInt(int64(rgba.A)),
			}}, true
		}
		ts.cursor = oldCursor
	}
	// Try transparent ---------------------------------------------------------
	if ts.consumeIdentTokenWith("transparent") != nil {
		c := csscolor.Transparent()
		return c, true
	}
	// Try currentColor --------------------------------------------------------
	if ts.consumeIdentTokenWith("currentColor") != nil {
		return csscolor.Color{Type: csscolor.TypeCurrentColor}, true
	}
	// TODO: Try system colors
	return csscolor.Color{}, false
}
