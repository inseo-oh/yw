package csssyntax

import (
	"errors"

	"github.com/inseo-oh/yw/css/text"
)

// https://www.w3.org/TR/css-text-3/#propdef-text-transform
func (ts *tokenStream) parseTextTransform() (text.Transform, error) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return text.Transform{Type: text.NoTransform}, nil
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
			if err := ts.consumeIdentTokenWith("capitalize"); err == nil {
				out.Type = text.Capitalize
				gotType = true
				valid = true
			} else if err := ts.consumeIdentTokenWith("uppercase"); err == nil {
				out.Type = text.Uppercase
				gotType = true
				valid = true
			} else if err := ts.consumeIdentTokenWith("lowercase"); err == nil {
				out.Type = text.Lowercase
				gotType = true
				valid = true
			}
		}
		if !gotIsFullWidth {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("full-width"); err == nil {
				out.FullWidth = true
				gotIsFullWidth = true
				valid = true
			}
		}
		if !gotIsFullKana {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("full-size-kana"); err == nil {
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
		return out, errors.New("invalid text-transform value")
	}
	return out, nil
}
