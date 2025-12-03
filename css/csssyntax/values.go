// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"errors"
	"fmt"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/values"
)

// Returns nil if not found
func (ts *tokenStream) parseNumber() *css.Num {
	oldCursor := ts.cursor
	numTk, err := ts.consumeTokenWith(tokenTypeNumber)
	if err != nil {
		ts.cursor = oldCursor
		return nil
	}
	per := numTk.(numberToken)
	return &per.value
}

// allowZeroShorthand should not be set if the property(such as line-height) also accepts number token.
// (In that case, 0 should be parsed as <number 0>, not <length 0>)
func (ts *tokenStream) parseLength(allowZeroShorthand bool) (values.Length, error) {
	dimTk, err := ts.consumeTokenWith(tokenTypeDimension)
	if err != nil {
		if allowZeroShorthand {
			oldCursor := ts.cursor
			numTk, err := ts.consumeTokenWith(tokenTypeNumber)
			if err != nil || !numTk.(numberToken).value.Equals(css.NumFromInt(0)) {
				ts.cursor = oldCursor
			} else {
				return values.Length{Value: css.NumFromInt(0), Unit: values.Px}, nil
			}
		}
		return values.Length{}, fmt.Errorf("expected length")
	}
	dim := dimTk.(dimensionToken)
	var unit values.LengthUnit
	switch dim.unit {
	case "em":
		unit = values.Em
	case "ex":
		unit = values.Ex
	case "ch":
		unit = values.Ch
	case "rem":
		unit = values.Rem
	case "vw":
		unit = values.Vw
	case "vh":
		unit = values.Vh
	case "vmin":
		unit = values.Vmin
	case "vmax":
		unit = values.Vmax
	case "cm":
		unit = values.Cm
	case "mm":
		unit = values.Mm
	case "q":
		unit = values.Q
	case "pc":
		unit = values.Pc
	case "pt":
		unit = values.Pt
	case "px":
		unit = values.Px
	default:
		return values.Length{}, fmt.Errorf("<bad LengthUnit %s>", dim.unit)
	}
	return values.Length{Value: dim.value, Unit: unit}, nil
}

// Returns nil if not found
func (ts *tokenStream) parsePercentage() (values.Percentage, error) {
	perTk, err := ts.consumeTokenWith(tokenTypePercentage)
	if err != nil {
		return values.Percentage{}, err
	}
	per := perTk.(percentageToken)
	return values.Percentage{Value: per.value}, nil
}

// https://www.w3.org/TR/css-values-3/#typedef-length-percentage
func (ts *tokenStream) parseLengthOrPercentage(allowZeroShorthand bool) (values.LengthResolvable, error) {
	if len, err := ts.parseLength(allowZeroShorthand); err == nil {
		return len, nil
	}
	if per, err := ts.parsePercentage(); err == nil {
		return per, nil
	}
	return nil, errors.New("expected length or percentage")
}
