package csssyntax

import (
	"github.com/inseo-oh/yw/css/textdecor"
)

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-line
func (ts *tokenStream) parseTextDecorationLine() (textdecor.LineFlags, bool) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return textdecor.NoLine, true
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
		return out, false
	}
	return out, true
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-style
func (ts *tokenStream) parseTextDecorationStyle() (textdecor.Style, bool) {
	if err := ts.consumeIdentTokenWith("solid"); err == nil {
		return textdecor.Solid, true
	} else if err := ts.consumeIdentTokenWith("double"); err == nil {
		return textdecor.Double, true
	} else if err := ts.consumeIdentTokenWith("dotted"); err == nil {
		return textdecor.Dotted, true
	} else if err := ts.consumeIdentTokenWith("dashed"); err == nil {
		return textdecor.Dashed, true
	} else if err := ts.consumeIdentTokenWith("wavy"); err == nil {
		return textdecor.Wavy, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-underline-position
func (ts *tokenStream) parseTextDecorationPosition() (textdecor.PositionFlags, bool) {
	if err := ts.consumeIdentTokenWith("auto"); err == nil {
		return textdecor.PositionAuto, true
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
		return out, false
	}
	return out, true
}
