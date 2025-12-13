// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/sizing"
	"github.com/inseo-oh/yw/css/values"
)

func (ts *tokenStream) parseSizeValueImpl(acceptAuto, acceptNone bool) (res sizing.Size, err error) {
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
		ts := tokenStream{tokens: tk.value, tokenizerHelper: ts.tokenizerHelper}
		var size values.LengthResolvable
		if v, err := ts.parseLengthOrPercentage(true); err == nil {
			size = v
		}
		if !ts.isEnd() {
			return res, fmt.Errorf("%s: expected end", ts.errorHeader())
		}
		return sizing.Size{Type: sizing.FitContent, Size: size}, nil
	}
	if v, err := ts.parseLengthOrPercentage(true); err == nil {
		return sizing.Size{Type: sizing.ManualSize, Size: v}, nil
	}
	return res, fmt.Errorf("%s: expected size value", ts.errorHeader())
}
func (ts *tokenStream) parseSizeOrAuto() (res sizing.Size, err error) {
	return ts.parseSizeValueImpl(true, false)
}
func (ts *tokenStream) parseSizeOrNone() (res sizing.Size, err error) {
	return ts.parseSizeValueImpl(false, true)
}
