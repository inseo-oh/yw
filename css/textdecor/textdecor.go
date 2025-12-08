// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// Package text provides types and values for [CSS Text Decoration Module Level 3].
//
// [CSS Text Decoration Module Level 3]: https://www.w3.org/TR/css-text-decor-3
package textdecor

import (
	"log"
	"strings"
)

// LineFlags represents value of [CSS text-decoration-line] property.
//
// [CSS text-decoration-line]: https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-line
type LineFlags uint8

const (
	NoLine      LineFlags = 0      // text-decoration-line: none
	Underline   LineFlags = 1 << 1 // text-decoration-line: underline
	Overline    LineFlags = 1 << 2 // text-decoration-line: overline
	LineThrough LineFlags = 1 << 3 // text-decoration-line: linethrough
	Blink       LineFlags = 1 << 4 // text-decoration-line: blink
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

// Style represents value of [CSS text-decoration-style] property.
//
// [CSS text-decoration-style]: https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-style
type Style uint8

const (
	Solid  Style = iota // text-decoration-style: solid
	Double              // text-decoration-style: double
	Dotted              // text-decoration-style: dotted
	Dashed              // text-decoration-style: dashed
	Wavy                // text-decoration-style: wavy
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

// PositionFlags represents value of [CSS text-underline-position] property.
//
// [CSS text-underline-position]: https://www.w3.org/TR/css-text-decor-3/#propdef-text-underline-position
type PositionFlags uint8

// Note that [SideLeft] and [SideRight] flag must be used with either [PositionAuto] or [Under].
// In CSS syntax, "auto left" and "auto right" isn't accepted, but writing just "left"
// or "right" implies "auto", so [PositionAuto] will be set.
const (
	PositionAuto PositionFlags = 1 << 1 // text-underline-position: auto
	Under        PositionFlags = 1 << 2 // text-underline-position: under
	SideLeft     PositionFlags = 1 << 3 // text-underline-position: left
	SideRight    PositionFlags = 2 << 3 // text-underline-position: right
	SideMask     PositionFlags = SideLeft | SideRight
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
