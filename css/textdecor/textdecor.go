// Implementation of the CSS Text Decoration Module Level 3 (https://www.w3.org/TR/css-text-decor-3)
package textdecor

import (
	"log"
	"strings"
)

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-line
type LineFlags uint8

const (
	NoLine      = LineFlags(0)
	Underline   = LineFlags(1 << 1)
	Overline    = LineFlags(1 << 2)
	LineThrough = LineFlags(1 << 3)
	Blink       = LineFlags(1 << 4)
)

func (t LineFlags) String() string {
	if t == NoLine {
		return "none"
	}
	sb := strings.Builder{}
	if t&Underline != 0 {
		sb.WriteString("underline ")
	}
	if t&Overline != 0 {
		sb.WriteString("overline ")
	}
	if t&LineThrough != 0 {
		sb.WriteString("line-through ")
	}
	if t&Blink != 0 {
		sb.WriteString("blink ")
	}
	return strings.TrimSpace(sb.String())
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-style
type Style uint8

const (
	Solid Style = iota
	Double
	Dotted
	Dashed
	Wavy
)

func (s Style) String() string {
	switch s {
	case Solid:
		return "solid"
	case Double:
		return "double"
	case Dotted:
		return "dotted"
	case Dashed:
		return "dashed"
	case Wavy:
		return "wavy"
	default:
		log.Panicf("<bad Style %d>", s)
	}
	return ""
}

// https://www.w3.org/TR/css-text-decor-3/#propdef-text-underline-position
type PositionFlags uint8

const (
	PositionAuto = PositionFlags(1 << 1)
	Under        = PositionFlags(1 << 2)
	SideLeft     = PositionFlags(1 << 3) // Must be used with either PositionAuto or Under
	SideRight    = PositionFlags(2 << 3) // Must be used with either PositionAuto or Under
	SideMask     = SideLeft | SideRight
)

func (t PositionFlags) String() string {
	if t == PositionAuto {
		return "auto"
	}
	sb := strings.Builder{}
	if t&Under != 0 {
		sb.WriteString("under ")
	}
	switch t & SideMask {
	case SideLeft:
		sb.WriteString("left ")
	case SideRight:
		sb.WriteString("right ")
	}
	return strings.TrimSpace(sb.String())
}
