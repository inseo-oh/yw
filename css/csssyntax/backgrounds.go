package csssyntax

import (
	"github.com/inseo-oh/yw/css/backgrounds"
	"github.com/inseo-oh/yw/css/values"
)

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
func (ts *tokenStream) parseLineStyle() (backgrounds.LineStyle, bool) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return backgrounds.NoLine, true
	} else if err := ts.consumeIdentTokenWith("hidden"); err == nil {
		return backgrounds.HiddenLine, true
	} else if err := ts.consumeIdentTokenWith("dotted"); err == nil {
		return backgrounds.DottedLine, true
	} else if err := ts.consumeIdentTokenWith("dashed"); err == nil {
		return backgrounds.DashedLine, true
	} else if err := ts.consumeIdentTokenWith("solid"); err == nil {
		return backgrounds.SolidLine, true
	} else if err := ts.consumeIdentTokenWith("double"); err == nil {
		return backgrounds.DoubleLine, true
	} else if err := ts.consumeIdentTokenWith("groove"); err == nil {
		return backgrounds.GrovveLine, true
	} else if err := ts.consumeIdentTokenWith("ridge"); err == nil {
		return backgrounds.RidgeLine, true
	} else if err := ts.consumeIdentTokenWith("inset"); err == nil {
		return backgrounds.InsetLine, true
	} else if err := ts.consumeIdentTokenWith("outset"); err == nil {
		return backgrounds.OutsetLine, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
func (ts *tokenStream) parseLineWidth() (values.Length, bool) {
	if err := ts.consumeIdentTokenWith("thin"); err == nil {
		return backgrounds.LineWidthThin, true
	}
	if err := ts.consumeIdentTokenWith("medium"); err == nil {
		return backgrounds.LineWidthMedium, true
	}
	if err := ts.consumeIdentTokenWith("thick"); err == nil {
		return backgrounds.LineWidthThick, true
	}
	if len, _ := ts.parseLength(true); len != nil {
		return *len, true
	}
	return values.Length{}, false
}
