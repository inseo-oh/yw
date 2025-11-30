package csssyntax

import (
	"github.com/inseo-oh/yw/css/backgrounds"
	"github.com/inseo-oh/yw/css/values"
)

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
func (ts *tokenStream) parseLineStyle() (backgrounds.LineStyle, bool) {
	if ts.consumeIdentTokenWith("none") != nil {
		return backgrounds.NoLine, true
	} else if ts.consumeIdentTokenWith("hidden") != nil {
		return backgrounds.HiddenLine, true
	} else if ts.consumeIdentTokenWith("dotted") != nil {
		return backgrounds.DottedLine, true
	} else if ts.consumeIdentTokenWith("dashed") != nil {
		return backgrounds.DashedLine, true
	} else if ts.consumeIdentTokenWith("solid") != nil {
		return backgrounds.SolidLine, true
	} else if ts.consumeIdentTokenWith("double") != nil {
		return backgrounds.DoubleLine, true
	} else if ts.consumeIdentTokenWith("groove") != nil {
		return backgrounds.GrovveLine, true
	} else if ts.consumeIdentTokenWith("ridge") != nil {
		return backgrounds.RidgeLine, true
	} else if ts.consumeIdentTokenWith("inset") != nil {
		return backgrounds.InsetLine, true
	} else if ts.consumeIdentTokenWith("outset") != nil {
		return backgrounds.OutsetLine, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
func (ts *tokenStream) parseLineWidth() (values.Length, bool) {
	if ts.consumeIdentTokenWith("thin") != nil {
		return backgrounds.LineWidthThin(), true
	}
	if ts.consumeIdentTokenWith("medium") != nil {
		return backgrounds.LineWidthMedium(), true
	}
	if ts.consumeIdentTokenWith("thick") != nil {
		return backgrounds.LineWidthThick(), true
	}
	if len, _ := ts.parseLength(true); len != nil {
		return *len, true
	}
	return values.Length{}, false
}
