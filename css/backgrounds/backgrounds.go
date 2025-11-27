// Implementation of the CSS Backgrounds and Borders Module Level 3 (https://www.w3.org/TR/css-backgrounds-3/)
package backgrounds

import (
	"fmt"
	"yw/css"
	"yw/css/values"
)

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
type LineStyle uint8

const (
	NoLine     = LineStyle(iota) // none
	HiddenLine                   // hidden
	DottedLine                   // dotted
	DashedLine                   // dashed
	SolidLine                    // solid
	DoubleLine                   // double
	GrovveLine                   // groove
	RidgeLine                    // ridge
	InsetLine                    // inset
	OutsetLine                   // outset
)

func (l LineStyle) String() string {
	switch l {
	case NoLine:
		return "none"
	case HiddenLine:
		return "hidden"
	case DottedLine:
		return "dotted"
	case DashedLine:
		return "dashed"
	case SolidLine:
		return "solid"
	case DoubleLine:
		return "double"
	case GrovveLine:
		return "groove"
	case RidgeLine:
		return "ridge"
	case InsetLine:
		return "inset"
	case OutsetLine:
		return "outset"
	}
	return fmt.Sprintf("<bad LineStyle %d>", l)
}

func LineWidthThin() values.Length {
	return values.LengthFromPx(css.NumFromInt(1))
}
func LineWidthMedium() values.Length {
	return values.LengthFromPx(css.NumFromInt(3))
}
func LineWidthThick() values.Length {
	return values.LengthFromPx(css.NumFromInt(5))
}
