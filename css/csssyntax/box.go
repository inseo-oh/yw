package csssyntax

import (
	"github.com/inseo-oh/yw/css/box"
	"github.com/inseo-oh/yw/css/values"
	"github.com/inseo-oh/yw/util"
)

func (ts *tokenStream) parseMargin() (box.Margin, bool) {
	if v, err := ts.parseLengthOrPercentage(true); !util.IsNil(v) {
		return box.Margin{Value: v}, true
	} else if err != nil {
		return box.Margin{}, false
	}
	if err := ts.consumeIdentTokenWith("auto"); err == nil {
		return box.Margin{Value: nil}, true
	}
	return box.Margin{}, false
}
func (ts *tokenStream) parsePadding() (values.LengthResolvable, bool) {
	v, _ := ts.parseLengthOrPercentage(true)
	if util.IsNil(v) {
		return nil, false
	}
	if len, ok := v.(values.Length); ok && len.Value.ToInt() < 0 {
		return nil, false
	}
	return v, true
}
