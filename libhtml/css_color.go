// Implementation of the CSS Color Module Level 4 (https://www.w3.org/TR/css-color-4/)
package libhtml

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	cm "yw/libcommon"
)

type css_rgba_color struct {
	red, green, blue, alpha uint8
}

// https://www.w3.org/TR/css-color-4/#named-colors
var css_named_colors_map = map[string]css_rgba_color{
	"aliceblue":            {240, 248, 255, 255},
	"antiquewhite":         {250, 235, 215, 255},
	"aqua":                 {0, 255, 255, 255},
	"aquamarine":           {127, 255, 212, 255},
	"azure":                {240, 255, 255, 255},
	"beige":                {245, 245, 220, 255},
	"bisque":               {255, 228, 196, 255},
	"black":                {0, 0, 0, 255},
	"blanchedalmond":       {255, 235, 205, 255},
	"blue":                 {0, 0, 255, 255},
	"blueviolet":           {138, 43, 226, 255},
	"brown":                {165, 42, 42, 255},
	"burlywood":            {222, 184, 135, 255},
	"cadetblue":            {95, 158, 160, 255},
	"chartreuse":           {127, 255, 0, 255},
	"chocolate":            {210, 105, 30, 255},
	"coral":                {255, 127, 80, 255},
	"cornflowerblue":       {100, 149, 237, 255},
	"cornsilk":             {255, 248, 220, 255},
	"crimson":              {220, 20, 60, 255},
	"cyan":                 {0, 255, 255, 255},
	"darkblue":             {0, 0, 139, 255},
	"darkcyan":             {0, 139, 139, 255},
	"darkgoldenrod":        {184, 134, 11, 255},
	"darkgray":             {169, 169, 169, 255},
	"darkgreen":            {0, 100, 0, 255},
	"darkgrey":             {169, 169, 169, 255},
	"darkkhaki":            {189, 183, 107, 255},
	"darkmagenta":          {139, 0, 139, 255},
	"darkolivegreen":       {85, 107, 47, 255},
	"darkorange":           {255, 140, 0, 255},
	"darkorchid":           {153, 50, 204, 255},
	"darkred":              {139, 0, 0, 255},
	"darksalmon":           {233, 150, 122, 255},
	"darkseagreen":         {143, 188, 143, 255},
	"darkslateblue":        {72, 61, 139, 255},
	"darkslategray":        {47, 79, 79, 255},
	"darkslategrey":        {47, 79, 79, 255},
	"darkturquoise":        {0, 206, 209, 255},
	"darkviolet":           {148, 0, 211, 255},
	"deeppink":             {255, 20, 147, 255},
	"deepskyblue":          {0, 191, 255, 255},
	"dimgray":              {105, 105, 105, 255},
	"dimgrey":              {105, 105, 105, 255},
	"dodgerblue":           {30, 144, 255, 255},
	"firebrick":            {178, 34, 34, 255},
	"floralwhite":          {255, 250, 240, 255},
	"forestgreen":          {34, 139, 34, 255},
	"fuchsia":              {255, 0, 255, 255},
	"gainsboro":            {220, 220, 220, 255},
	"ghostwhite":           {248, 248, 255, 255},
	"gold":                 {255, 215, 0, 255},
	"goldenrod":            {218, 165, 32, 255},
	"gray":                 {128, 128, 128, 255},
	"green":                {0, 128, 0, 255},
	"greenyellow":          {173, 255, 47, 255},
	"grey":                 {128, 128, 128, 255},
	"honeydew":             {240, 255, 240, 255},
	"hotpink":              {255, 105, 180, 255},
	"indianred":            {205, 92, 92, 255},
	"indigo":               {75, 0, 130, 255},
	"ivory":                {255, 255, 240, 255},
	"khaki":                {240, 230, 140, 255},
	"lavender":             {230, 230, 250, 255},
	"lavenderblush":        {255, 240, 245, 255},
	"lawngreen":            {124, 252, 0, 255},
	"lemonchiffon":         {255, 250, 205, 255},
	"lightblue":            {173, 216, 230, 255},
	"lightcoral":           {240, 128, 128, 255},
	"lightcyan":            {224, 255, 255, 255},
	"lightgoldenrodyellow": {250, 250, 210, 255},
	"lightgray":            {211, 211, 211, 255},
	"lightgreen":           {144, 238, 144, 255},
	"lightgrey":            {211, 211, 211, 255},
	"lightpink":            {255, 182, 193, 255},
	"lightsalmon":          {255, 160, 122, 255},
	"lightseagreen":        {32, 178, 170, 255},
	"lightskyblue":         {135, 206, 250, 255},
	"lightslategray":       {119, 136, 153, 255},
	"lightslategrey":       {119, 136, 153, 255},
	"lightsteelblue":       {176, 196, 222, 255},
	"lightyellow":          {255, 255, 224, 255},
	"lime":                 {0, 255, 0, 255},
	"limegreen":            {50, 205, 50, 255},
	"linen":                {250, 240, 230, 255},
	"magenta":              {255, 0, 255, 255},
	"maroon":               {128, 0, 0, 255},
	"mediumaquamarine":     {102, 205, 170, 255},
	"mediumblue":           {0, 0, 205, 255},
	"mediumorchid":         {186, 85, 211, 255},
	"mediumpurple":         {147, 112, 219, 255},
	"mediumseagreen":       {60, 179, 113, 255},
	"mediumslateblue":      {123, 104, 238, 255},
	"mediumspringgreen":    {0, 250, 154, 255},
	"mediumturquoise":      {72, 209, 204, 255},
	"mediumvioletred":      {199, 21, 133, 255},
	"midnightblue":         {25, 25, 112, 255},
	"mintcream":            {245, 255, 250, 255},
	"mistyrose":            {255, 228, 225, 255},
	"moccasin":             {255, 228, 181, 255},
	"navajowhite":          {255, 222, 173, 255},
	"navy":                 {0, 0, 128, 255},
	"oldlace":              {253, 245, 230, 255},
	"olive":                {128, 128, 0, 255},
	"olivedrab":            {107, 142, 35, 255},
	"orange":               {255, 165, 0, 255},
	"orangered":            {255, 69, 0, 255},
	"orchid":               {218, 112, 214, 255},
	"palegoldenrod":        {238, 232, 170, 255},
	"palegreen":            {152, 251, 152, 255},
	"paleturquoise":        {175, 238, 238, 255},
	"palevioletred":        {219, 112, 147, 255},
	"papayawhip":           {255, 239, 213, 255},
	"peachpuff":            {255, 218, 185, 255},
	"peru":                 {205, 133, 63, 255},
	"pink":                 {255, 192, 203, 255},
	"plum":                 {221, 160, 221, 255},
	"powderblue":           {176, 224, 230, 255},
	"purple":               {128, 0, 128, 255},
	"rebeccapurple":        {102, 51, 153, 255},
	"red":                  {255, 0, 0, 255},
	"rosybrown":            {188, 143, 143, 255},
	"royalblue":            {65, 105, 225, 255},
	"saddlebrown":          {139, 69, 19, 255},
	"salmon":               {250, 128, 114, 255},
	"sandybrown":           {244, 164, 96, 255},
	"seagreen":             {46, 139, 87, 255},
	"seashell":             {255, 245, 238, 255},
	"sienna":               {160, 82, 45, 255},
	"silver":               {192, 192, 192, 255},
	"skyblue":              {135, 206, 235, 255},
	"slateblue":            {106, 90, 205, 255},
	"slategray":            {112, 128, 144, 255},
	"slategrey":            {112, 128, 144, 255},
	"snow":                 {255, 250, 250, 255},
	"springgreen":          {0, 255, 127, 255},
	"steelblue":            {70, 130, 180, 255},
	"tan":                  {210, 180, 140, 255},
	"teal":                 {0, 128, 128, 255},
	"thistle":              {216, 191, 216, 255},
	"tomato":               {255, 99, 71, 255},
	"turquoise":            {64, 224, 208, 255},
	"violet":               {238, 130, 238, 255},
	"wheat":                {245, 222, 179, 255},
	"white":                {255, 255, 255, 255},
	"whitesmoke":           {245, 245, 245, 255},
	"yellow":               {255, 255, 0, 255},
	"yellowgreen":          {154, 205, 50, 255},
}

type css_color struct {
	tp     css_color_type
	values []css_number
}
type css_color_type uint8

func css_color_CanvasText() css_color {
	return css_color_from_rgba(0, 0, 0, 255)
}
func css_color_transparent() css_color {
	return css_color_from_rgba(0, 0, 0, 0)
}
func css_color_from_rgba(r, g, b, a uint8) css_color {
	return css_color{css_color_type_rgb, []css_number{
		css_number_from_int(int64(r)),
		css_number_from_int(int64(g)),
		css_number_from_int(int64(b)),
		css_number_from_int(int64(a)),
	}}
}

const (
	css_color_type_rgb           = css_color_type(iota) // rgb(), rgba(), hex colors, named colors
	css_color_type_current_color                        // currentColor
	css_color_type_hsl                                  // hsl(), hsla()
	css_color_type_hwb                                  // hwb()
	css_color_type_lab                                  // lab()
	css_color_type_lch                                  // lch()
	css_color_type_oklab                                // oklab()
	css_color_type_oklch                                // oklch()
	css_color_type_color                                // color()
	// TODO: System color
)

func (c css_color) equals(other css_color) bool {
	if c.tp != other.tp {
		return false
	}
	if len(c.values) != len(other.values) {
		return false
	}
	for i := range len(c.values) {
		if !c.values[i].equals(other.values[i]) {
			return false
		}
	}
	return true
}
func (c css_color) String() string {
	switch c.tp {
	case css_color_type_rgb:
		return fmt.Sprintf("#%02x%02x%02x%02x", c.values[0].to_int(), c.values[1].to_int(), c.values[2].to_int(), c.values[3].to_int())
	case css_color_type_current_color:
		return "currentColor"
	case css_color_type_hsl:
		return fmt.Sprintf("hsl(%v, %v, %v, %v)", c.values[0], c.values[1], c.values[2], c.values[3]) // STUB
	case css_color_type_hwb:
		return fmt.Sprintf("hwb(%v, %v, %v, %v)", c.values[0], c.values[1], c.values[2], c.values[3]) // STUB
	case css_color_type_lab:
		return fmt.Sprintf("lab(%v, %v, %v, %v)", c.values[0], c.values[1], c.values[2], c.values[3]) // STUB
	case css_color_type_lch:
		return fmt.Sprintf("lch(%v, %v, %v, %v)", c.values[0], c.values[1], c.values[2], c.values[3]) // STUB
	case css_color_type_oklab:
		return fmt.Sprintf("oklab(%v, %v, %v, %v)", c.values[0], c.values[1], c.values[2], c.values[3]) // STUB
	case css_color_type_oklch:
		return fmt.Sprintf("oklch(%v, %v, %v, %v)", c.values[0], c.values[1], c.values[2], c.values[3]) // STUB
	case css_color_type_color:
		return "color(...)" // STUB
	}
	return fmt.Sprintf("<unknown css_color type %v>", c.tp)
}

func (ts *css_token_stream) parse_color() (*css_color, error) {
	// Try hex notation --------------------------------------------------------
	if tk := ts.consume_token_with_type(css_token_type_hash); !cm.IsNil(tk) {
		// https://www.w3.org/TR/css-color-4/#hex-notation
		chrs := []rune(tk.(css_hash_token).value)
		var r_str, g_str, b_str, a_str string
		a_str = "ff"
		switch len(chrs) {
		case 3:
			// #rgb
			r_str = strings.Repeat(string(chrs[0]), 2)
			g_str = strings.Repeat(string(chrs[1]), 2)
			b_str = strings.Repeat(string(chrs[2]), 2)
		case 4:
			// #rgba
			r_str = strings.Repeat(string(chrs[0]), 2)
			g_str = strings.Repeat(string(chrs[1]), 2)
			b_str = strings.Repeat(string(chrs[2]), 2)
			a_str = strings.Repeat(string(chrs[3]), 2)
		case 6:
			// #rrggbb
			r_str = string(chrs[0:2])
			g_str = string(chrs[2:4])
			b_str = string(chrs[4:6])
		case 8:
			// #rrggbb
			r_str = string(chrs[0:2])
			g_str = string(chrs[2:4])
			b_str = string(chrs[4:6])
			a_str = string(chrs[6:8])
		default:
			return nil, errors.New("illegal hex color length")
		}
		r, err := strconv.ParseUint(r_str, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("illegal hex color digit: %w", err)
		}
		g, err := strconv.ParseUint(g_str, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("illegal hex color digit: %w", err)
		}
		b, err := strconv.ParseUint(b_str, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("illegal hex color digit: %w", err)
		}
		a, err := strconv.ParseUint(a_str, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("illegal hex color digit: %w", err)
		}
		return &css_color{css_color_type_rgb, []css_number{
			css_number_from_int(int64(r)),
			css_number_from_int(int64(g)),
			css_number_from_int(int64(b)),
			css_number_from_int(int64(a)),
		}}, nil
	}
	// Try rgb()/rgba() function -----------------------------------------------
	fn := ts.consume_ast_function_with("rgb")
	if fn == nil {
		// rgba() is alias for rgb()
		fn = ts.consume_ast_function_with("rgba")
	}
	if fn != nil {
		ts := css_token_stream{tokens: fn.value}

		parse_alpha := func() *css_number {
			var v css_number
			if num := ts.parse_number(); num != nil {
				v = css_number_from_float(num.clamp(css_number_from_int(0), css_number_from_int(1)).to_float() * 255)
			} else if num := ts.parse_percentage(); num != nil {
				a_per := num.value.clamp(css_number_from_int(0), css_number_from_int(100)).to_float()
				v = css_number_from_float((a_per / 100) * 255)
			} else {
				return nil
			}
			return &v
		}

		// https://www.w3.org/TR/css-color-4/#funcdef-rgb
		var r, g, b, a css_number
		a = css_number_from_int(255)
		old_cursor := ts.cursor

		//======================================================================
		// Try legacy syntax first
		//======================================================================
		// https://www.w3.org/TR/css-color-4/#typedef-legacy-rgb-syntax

		// rgb(<  >r  ,  g  ,  b  ) --------------------------------------------
		// rgb(<  >r  ,  g  ,  b  ,  a  ) --------------------------------------
		ts.skip_whitespaces()
		// rgb(  <r  ,  g  ,  b>  ) --------------------------------------------
		// rgb(  <r  ,  g  ,  b>  ,  a  ) --------------------------------------
		per, err := css_accept_comma_separated_repetion(&ts, 3, func(ts *css_token_stream) (*css_percentage, error) {
			return ts.parse_percentage(), nil
		})
		if per == nil && err != nil {
			return nil, err
		} else if len(per) == 3 {
			// Percentage value
			r_per := per[0].value.clamp(css_number_from_int(0), css_number_from_int(100)).to_float()
			g_per := per[1].value.clamp(css_number_from_int(0), css_number_from_int(100)).to_float()
			b_per := per[2].value.clamp(css_number_from_int(0), css_number_from_int(100)).to_float()
			r = css_number_from_float((r_per / 100) * 255)
			g = css_number_from_float((g_per / 100) * 255)
			b = css_number_from_float((b_per / 100) * 255)
		} else {
			num, err := css_accept_comma_separated_repetion(&ts, 3, func(ts *css_token_stream) (*css_number, error) {
				return ts.parse_number(), nil
			})
			if num == nil && err != nil {
				return nil, err
			} else if len(num) == 3 {
				// Number value
				r = num[0].clamp(css_number_from_int(0), css_number_from_int(255))
				g = num[1].clamp(css_number_from_int(0), css_number_from_int(255))
				b = num[2].clamp(css_number_from_int(0), css_number_from_int(255))
			} else {
				goto modern_syntax
			}
		}
		// rgb(  r  ,  g  ,  b<  >) --------------------------------------------
		// rgb(  r  ,  g  ,  b<  >,  a  ) --------------------------------------
		ts.skip_whitespaces()
		// rgb(  r  ,  g  ,  b  <,>  a  ) --------------------------------------
		if tk := ts.consume_token_with_type(css_token_type_comma); !cm.IsNil(tk) {
			// rgb(  r  ,  g  ,  b  ,<  >a  ) ----------------------------------
			ts.skip_whitespaces()
			// rgb(  r  ,  g  ,  b  ,  <a>  ) ----------------------------------
			if v := parse_alpha(); v != nil {
				a = *v
			} else {
				return nil, errors.New("expected alpha value after ','")
			}
			// rgb(  r  ,  g  ,  b  ,  a<  >) ----------------------------------
			ts.skip_whitespaces()
		}
		if !ts.is_end() {
			return nil, errors.New("unexpected junk at the end of function")
		}
		return &css_color{css_color_type_rgb, []css_number{r, g, b, a}}, nil
	modern_syntax:
		ts.cursor = old_cursor

		//======================================================================
		// Try modern syntax
		//======================================================================
		// https://www.w3.org/TR/css-color-4/#typedef-modern-rgb-syntax

		// rgb(<  >r  g  b  ) --------------------------------------------------
		// rgb(<  >r  g  b  /  a  ) --------------------------------------------
		ts.skip_whitespaces()
		// rgb(  <r  g  b  >) --------------------------------------------------
		// rgb(  <r  g  b  >/  a  ) --------------------------------------------
		components := []css_number{}
		for range 3 {
			// rgb(  <r>  <g>  <b>  ) ------------------------------------------
			// rgb(  <r>  <g>  <b>  /  a  ) ------------------------------------
			var v css_number
			if num := ts.parse_number(); num != nil {
				v = num.clamp(css_number_from_int(0), css_number_from_int(255))
			} else if num := ts.parse_percentage(); num != nil {
				per := num.value.clamp(css_number_from_int(0), css_number_from_int(100)).to_float()
				v = css_number_from_float((per / 100) * 255)
			} else if tk := ts.consume_ident_token_with("none"); !cm.IsNil(tk) {
				panic("TODO")
			} else {
				return nil, errors.New("unexpected value in rgb/rgba() function")
			}
			components = append(components, v)
			// rgb(  r<  >g<  >b<  >) ------------------------------------------
			// rgb(  r<  >g<  >b<  >/  a  ) ------------------------------------
			ts.skip_whitespaces()
		}
		// rgb(  r  g  b  </>  a  ) --------------------------------------------
		a = css_number_from_int(255)
		if tk := ts.consume_delim_token_with('/'); tk != nil {
			// rgb(  r  g  b  /<  >a  ) --------------------------------------------
			ts.skip_whitespaces()
			// rgb(  r  g  b  /  <a>  ) --------------------------------------------
			if v := parse_alpha(); v != nil {
				a = *v
			} else {
				return nil, errors.New("expected alpha value after '/'")
			}
			// rgb(  r  g  b  /  a<  >) --------------------------------------------
			ts.skip_whitespaces()
		}
		if !ts.is_end() {
			return nil, errors.New("unexpected junk at the end of function")
		}
		return &css_color{css_color_type_rgb, []css_number{components[0], components[1], components[2], a}}, nil
	}
	// Try hsl()/hsla() function -----------------------------------------------
	fn = ts.consume_ast_function_with("hsl")
	if fn == nil {
		// hsla() is alias for hsl()
		fn = ts.consume_ast_function_with("hsl")
	}
	if fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-hsl
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-hsl]")
	}
	// Try hwb() function ------------------------------------------------------
	if fn = ts.consume_ast_function_with("hwb"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-hwb
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-hwb]")
	}
	// Try lab() function ------------------------------------------------------
	if fn = ts.consume_ast_function_with("hwb"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-lab
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-lab]")
	}
	// Try lch() function ------------------------------------------------------
	if fn = ts.consume_ast_function_with("lch"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-lch
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-lch]")
	}
	// Try oklab() function ----------------------------------------------------
	if fn = ts.consume_ast_function_with("oklab"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-oklab
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-oklab]")
	}
	// Try oklch() function ----------------------------------------------------
	if fn = ts.consume_ast_function_with("oklch"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-oklch
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-oklch]")
	}
	// Try color() function ----------------------------------------------------
	if fn = ts.consume_ast_function_with("color"); fn != nil {
		// https://www.w3.org/TR/css-color-4/#funcdef-color
		panic("TODO[https://www.w3.org/TR/css-color-4/#funcdef-color]")
	}
	// Try named color ---------------------------------------------------------
	old_cursor := ts.cursor
	ident := ts.consume_token_with_type(css_token_type_ident)
	if !cm.IsNil(ident) {
		rgba, ok := css_named_colors_map[ident.(css_ident_token).value]
		if ok {
			return &css_color{css_color_type_rgb, []css_number{
				css_number_from_int(int64(rgba.red)),
				css_number_from_int(int64(rgba.green)),
				css_number_from_int(int64(rgba.blue)),
				css_number_from_int(int64(rgba.alpha)),
			}}, nil
		}
		ts.cursor = old_cursor
	}
	// Try transparent ---------------------------------------------------------
	if ts.consume_ident_token_with("transparent") != nil {
		c := css_color_transparent()
		return &c, nil
	}
	// Try currentColor --------------------------------------------------------
	if ts.consume_ident_token_with("currentColor") != nil {
		return &css_color{css_color_type_current_color, nil}, nil
	}
	// TODO: Try system colors

	return nil, nil
}

func init() {
	//==========================================================================
	// https://www.w3.org/TR/css-color-4/#the-color-property
	//==========================================================================
	// https://www.w3.org/TR/css-color-4/#propdef-color
	css_property_descriptors_map["color"] = css_property_descriptor{
		initial: css_color_CanvasText(),
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			return ts.parse_color()
		},
	}
}
