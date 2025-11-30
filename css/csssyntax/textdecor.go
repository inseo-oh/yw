package csssyntax

import (
	"github.com/inseo-oh/yw/css/textdecor"
)

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-line
func (ts *tokenStream) parseTextDecorationLine() (textdecor.LineFlags, bool) {
	if ts.consumeIdentTokenWith("none") {
		return textdecor.NoLine, true
	}
	var out textdecor.LineFlags
	gotAny := false
	for {
		valid := false
		ts.skipWhitespaces()
		if ts.consumeIdentTokenWith("underline") {
			out |= textdecor.Underline
			valid = true
		} else if ts.consumeIdentTokenWith("overline") {
			out |= textdecor.Overline
			valid = true
		} else if ts.consumeIdentTokenWith("line-through") {
			out |= textdecor.LineThrough
			valid = true
		} else if ts.consumeIdentTokenWith("blink") {
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
	if ts.consumeIdentTokenWith("solid") {
		return textdecor.Solid, true
	} else if ts.consumeIdentTokenWith("double") {
		return textdecor.Double, true
	} else if ts.consumeIdentTokenWith("dotted") {
		return textdecor.Dotted, true
	} else if ts.consumeIdentTokenWith("dashed") {
		return textdecor.Dashed, true
	} else if ts.consumeIdentTokenWith("wavy") {
		return textdecor.Wavy, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-underline-position
func (ts *tokenStream) parseTextDecorationPosition() (textdecor.PositionFlags, bool) {
	if ts.consumeIdentTokenWith("auto") {
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
			if ts.consumeIdentTokenWith("under") {
				out |= textdecor.Under
				gotUnder = true
				valid = true
			}
		}
		if !gotSide {
			ts.skipWhitespaces()
			if ts.consumeIdentTokenWith("left") {
				out |= textdecor.SideLeft
				if out == textdecor.SideLeft {
					// If these were used alone, auto is also implied
					out |= textdecor.PositionAuto
				}
				gotSide = true
				valid = true
			} else if ts.consumeIdentTokenWith("right") {
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
