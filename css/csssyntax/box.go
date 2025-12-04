// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/box"
	"github.com/inseo-oh/yw/css/values"
)

func (ts *tokenStream) parseMargin() (box.Margin, error) {
	if v, err := ts.parseLengthOrPercentage(true); err == nil {
		return box.Margin{Value: v}, nil
	}
	if err := ts.consumeIdentTokenWith("auto"); err == nil {
		return box.Margin{Value: nil}, nil
	}
	return box.Margin{}, fmt.Errorf("%s: expected margin", ts.errorHeader())
}
func (ts *tokenStream) parsePadding() (values.LengthResolvable, error) {
	v, err := ts.parseLengthOrPercentage(true)
	if err != nil {
		return nil, fmt.Errorf("%s: expected length or percentage", ts.errorHeader())
	}
	if len, ok := v.(values.Length); ok && len.Value.ToInt() < 0 {
		return nil, fmt.Errorf("%s: length is out of range", ts.errorHeader())
	}
	return v, nil
}
