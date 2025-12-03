package csssyntax

import (
	"github.com/inseo-oh/yw/css/sizing"
	"github.com/inseo-oh/yw/css/values"
	"github.com/inseo-oh/yw/util"
)

func (ts *tokenStream) parseSizeValueImpl(acceptAuto, acceptNone bool) (sizing.Size, bool) {
	if acceptAuto {
		if tk := ts.consumeIdentTokenWith("auto"); !util.IsNil(tk) {
			return sizing.Size{Type: sizing.Auto}, true
		}
	}
	if acceptNone {
		if tk := ts.consumeIdentTokenWith("none"); !util.IsNil(tk) {
			return sizing.Size{Type: sizing.Auto}, true
		}
	}
	if tk := ts.consumeIdentTokenWith("min-content"); !util.IsNil(tk) {
		return sizing.Size{Type: sizing.MinContent}, true
	}
	if tk := ts.consumeIdentTokenWith("max-content"); !util.IsNil(tk) {
		return sizing.Size{Type: sizing.MaxContent}, true
	}
	if tk, err := ts.consumeAstFuncWith("fit-content"); err == nil {
		ts := tokenStream{tokens: tk.value}
		var size values.LengthResolvable
		if v, err := ts.parseLengthOrPercentage(true); !util.IsNil(v) {
			size = v
		} else if err != nil {
			return sizing.Size{}, false
		}
		if !ts.isEnd() {
			return sizing.Size{}, false
		}
		return sizing.Size{Type: sizing.FitContent, Size: size}, true
	}
	if v, err := ts.parseLengthOrPercentage(true); !util.IsNil(v) {
		return sizing.Size{Type: sizing.ManualSize, Size: v}, true
	} else if err != nil {
		return sizing.Size{}, false
	}
	return sizing.Size{}, false
}
func (ts *tokenStream) parseSizeOrAuto() (sizing.Size, bool) {
	return ts.parseSizeValueImpl(true, false)
}
func (ts *tokenStream) parseSizeOrNone() (sizing.Size, bool) {
	return ts.parseSizeValueImpl(false, true)
}
