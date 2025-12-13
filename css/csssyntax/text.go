// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/text"
)

// https://www.w3.org/TR/css-text-3/#propdef-text-transform
func (ts *tokenStream) parseTextTransform() (res text.Transform, err error) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return text.Transform{Type: text.NoTransform}, nil
	}
	res = text.Transform{Type: text.OriginalCaps}
	gotType := false
	gotIsFullWidth := false
	gotIsFullKana := false
	gotAny := false
	for {
		valid := false
		if !gotType {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("capitalize"); err == nil {
				res.Type = text.Capitalize
				gotType = true
				valid = true
			} else if err := ts.consumeIdentTokenWith("uppercase"); err == nil {
				res.Type = text.Uppercase
				gotType = true
				valid = true
			} else if err := ts.consumeIdentTokenWith("lowercase"); err == nil {
				res.Type = text.Lowercase
				gotType = true
				valid = true
			}
		}
		if !gotIsFullWidth {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("full-width"); err == nil {
				res.FullWidth = true
				gotIsFullWidth = true
				valid = true
			}
		}
		if !gotIsFullKana {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("full-size-kana"); err == nil {
				res.FullSizeKana = true
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
		return res, fmt.Errorf("%s: invalid text-transform value", ts.errorHeader())
	}
	return res, nil
}
