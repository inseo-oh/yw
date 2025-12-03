package csssyntax

import (
	"github.com/inseo-oh/yw/css/box"
	"github.com/inseo-oh/yw/css/values"
)

func (ts *tokenStream) parseMargin() (box.Margin, bool) {
	if v, err := ts.parseLengthOrPercentage(true); err == nil {
		return box.Margin{Value: v}, true
	}
	if err := ts.consumeIdentTokenWith("auto"); err == nil {
		return box.Margin{Value: nil}, true
	}
	return box.Margin{}, false
}
func (ts *tokenStream) parsePadding() (values.LengthResolvable, bool) {
	v, err := ts.parseLengthOrPercentage(true)
	if err != nil {
		return nil, false
	}
	if len, ok := v.(values.Length); ok && len.Value.ToInt() < 0 {
		return nil, false
	}
	return v, true
}
