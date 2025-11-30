package csssyntax

import (
	"github.com/inseo-oh/yw/css/text"
	"github.com/inseo-oh/yw/util"
)

// https://www.w3.org/TR/css-text-3/#propdef-text-transform
func (ts *tokenStream) parseTextTransform() (text.Transform, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("none")) {
		return text.Transform{Type: text.NoTransform}, true
	}
	out := text.Transform{Type: text.OriginalCaps}
	gotType := false
	gotIsFullWidth := false
	gotIsFullKana := false
	gotAny := false
	for {
		valid := false
		if !gotType {
			ts.skipWhitespaces()
			if !util.IsNil(ts.consumeIdentTokenWith("capitalize")) {
				out.Type = text.Capitalize
				gotType = true
				valid = true
			} else if !util.IsNil(ts.consumeIdentTokenWith("uppercase")) {
				out.Type = text.Uppercase
				gotType = true
				valid = true
			} else if !util.IsNil(ts.consumeIdentTokenWith("lowercase")) {
				out.Type = text.Lowercase
				gotType = true
				valid = true
			}
		}
		if !gotIsFullWidth {
			ts.skipWhitespaces()
			if !util.IsNil(ts.consumeIdentTokenWith("full-width")) {
				out.FullWidth = true
				gotIsFullWidth = true
				valid = true
			}
		}
		if !gotIsFullKana {
			ts.skipWhitespaces()
			if !util.IsNil(ts.consumeIdentTokenWith("full-size-kana")) {
				out.FullSizeKana = true
				gotIsFullWidth = true
				valid = true
			}
		}
		ts.skipWhitespaces()
		if !valid {
			break
		}
		gotAny = true
	}
	if !gotAny {
		return out, false
	}
	return out, true
}
