package csssyntax

import (
	"github.com/inseo-oh/yw/css/backgrounds"
	"github.com/inseo-oh/yw/css/values"
)

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
func (ts *tokenStream) parseLineStyle() (backgrounds.LineStyle, bool) {
	if ts.consumeIdentTokenWith("none") {
		return backgrounds.NoLine, true
	} else if ts.consumeIdentTokenWith("hidden") {
		return backgrounds.HiddenLine, true
	} else if ts.consumeIdentTokenWith("dotted") {
		return backgrounds.DottedLine, true
	} else if ts.consumeIdentTokenWith("dashed") {
		return backgrounds.DashedLine, true
	} else if ts.consumeIdentTokenWith("solid") {
		return backgrounds.SolidLine, true
	} else if ts.consumeIdentTokenWith("double") {
		return backgrounds.DoubleLine, true
	} else if ts.consumeIdentTokenWith("groove") {
		return backgrounds.GrovveLine, true
	} else if ts.consumeIdentTokenWith("ridge") {
		return backgrounds.RidgeLine, true
	} else if ts.consumeIdentTokenWith("inset") {
		return backgrounds.InsetLine, true
	} else if ts.consumeIdentTokenWith("outset") {
		return backgrounds.OutsetLine, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
func (ts *tokenStream) parseLineWidth() (values.Length, bool) {
	if ts.consumeIdentTokenWith("thin") {
		return backgrounds.LineWidthThin(), true
	}
	if ts.consumeIdentTokenWith("medium") {
		return backgrounds.LineWidthMedium(), true
	}
	if ts.consumeIdentTokenWith("thick") {
		return backgrounds.LineWidthThick(), true
	}
	if len, _ := ts.parseLength(true); len != nil {
		return *len, true
	}
	return values.Length{}, false
}
