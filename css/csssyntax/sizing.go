package csssyntax

import (
	"errors"

	"github.com/inseo-oh/yw/css/sizing"
	"github.com/inseo-oh/yw/css/values"
)

func (ts *tokenStream) parseSizeValueImpl(acceptAuto, acceptNone bool) (sizing.Size, error) {
	if acceptAuto {
		if err := ts.consumeIdentTokenWith("auto"); err == nil {
			return sizing.Size{Type: sizing.Auto}, nil
		}
	}
	if acceptNone {
		if err := ts.consumeIdentTokenWith("none"); err == nil {
			return sizing.Size{Type: sizing.Auto}, nil
		}
	}
	if err := ts.consumeIdentTokenWith("min-content"); err == nil {
		return sizing.Size{Type: sizing.MinContent}, nil
	}
	if err := ts.consumeIdentTokenWith("max-content"); err == nil {
		return sizing.Size{Type: sizing.MaxContent}, nil
	}
	if tk, err := ts.consumeAstFuncWith("fit-content"); err == nil {
		ts := tokenStream{tokens: tk.value}
		var size values.LengthResolvable
		if v, err := ts.parseLengthOrPercentage(true); err == nil {
			size = v
		}
		if !ts.isEnd() {
			return sizing.Size{}, errors.New("expected end")
		}
		return sizing.Size{Type: sizing.FitContent, Size: size}, nil
	}
	if v, err := ts.parseLengthOrPercentage(true); err == nil {
		return sizing.Size{Type: sizing.ManualSize, Size: v}, nil
	}
	return sizing.Size{}, errors.New("expected size value")
}
func (ts *tokenStream) parseSizeOrAuto() (sizing.Size, error) {
	return ts.parseSizeValueImpl(true, false)
}
func (ts *tokenStream) parseSizeOrNone() (sizing.Size, error) {
	return ts.parseSizeValueImpl(false, true)
}
