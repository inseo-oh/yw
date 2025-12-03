// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"errors"

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
	return box.Margin{}, errors.New("expected margin")
}
func (ts *tokenStream) parsePadding() (values.LengthResolvable, error) {
	v, err := ts.parseLengthOrPercentage(true)
	if err != nil {
		return nil, errors.New("expected length or percentage")
	}
	if len, ok := v.(values.Length); ok && len.Value.ToInt() < 0 {
		return nil, errors.New("length is out of range")
	}
	return v, nil
}
