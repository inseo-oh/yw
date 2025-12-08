// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/textdecor"
)

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-line
func (ts *tokenStream) parseTextDecorationLine() (res textdecor.LineFlags, err error) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return textdecor.NoLine, nil
	}
	var out textdecor.LineFlags
	gotAny := false
	for {
		valid := false
		ts.skipWhitespaces()
		if err := ts.consumeIdentTokenWith("underline"); err == nil {
			out |= textdecor.Underline
			valid = true
		} else if err := ts.consumeIdentTokenWith("overline"); err == nil {
			out |= textdecor.Overline
			valid = true
		} else if err := ts.consumeIdentTokenWith("line-through"); err == nil {
			out |= textdecor.LineThrough
			valid = true
		} else if err := ts.consumeIdentTokenWith("blink"); err == nil {
			out |= textdecor.Blink
			valid = true
		}
		ts.skipWhitespaces()
		if !valid {
			break
		}
		gotAny = true
	}
	if !gotAny {
		return 0, fmt.Errorf("%s: invalid text-decoration-line value", ts.errorHeader())
	}
	return out, nil
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-style
func (ts *tokenStream) parseTextDecorationStyle() (res textdecor.Style, err error) {
	if err := ts.consumeIdentTokenWith("solid"); err == nil {
		return textdecor.Solid, nil
	} else if err := ts.consumeIdentTokenWith("double"); err == nil {
		return textdecor.Double, nil
	} else if err := ts.consumeIdentTokenWith("dotted"); err == nil {
		return textdecor.Dotted, nil
	} else if err := ts.consumeIdentTokenWith("dashed"); err == nil {
		return textdecor.Dashed, nil
	} else if err := ts.consumeIdentTokenWith("wavy"); err == nil {
		return textdecor.Wavy, nil
	}
	return 0, fmt.Errorf("%s: invalid text-decoration-style value", ts.errorHeader())
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-underline-position
func (ts *tokenStream) parseTextDecorationPosition() (res textdecor.PositionFlags, err error) {
	if err := ts.consumeIdentTokenWith("auto"); err == nil {
		return textdecor.PositionAuto, nil
	}
	var out textdecor.PositionFlags
	gotUnder := false
	gotSide := false
	gotAny := false
	for {
		valid := false
		if !gotUnder {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("under"); err == nil {
				out |= textdecor.Under
				gotUnder = true
				valid = true
			}
		}
		if !gotSide {
			ts.skipWhitespaces()
			if err := ts.consumeIdentTokenWith("left"); err == nil {
				out |= textdecor.SideLeft
				if out == textdecor.SideLeft {
					// If these were used alone, auto is also implied
					out |= textdecor.PositionAuto
				}
				gotSide = true
				valid = true
			} else if err := ts.consumeIdentTokenWith("right"); err == nil {
				out |= textdecor.SideRight
				if out == textdecor.SideRight {
					// If these were used alone, auto is also implied
					out |= textdecor.PositionAuto
				}
				gotSide = true
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
		return out, fmt.Errorf("%s: invalid text-underline-position value", ts.errorHeader())
	}
	return out, nil
}
