// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Package backgrounds provide types and values for CSS Backgrounds and Borders Module Level 3
//
// Releavnt spec: https://www.w3.org/TR/css-backgrounds-3/
package backgrounds

import (
	"fmt"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/values"
)

// LineStyle represents CSS <line-style> type
//
// Releavnt spec: https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
type LineStyle uint8

const (
	NoLine     LineStyle = iota // none
	HiddenLine                  // hidden
	DottedLine                  // dotted
	DashedLine                  // dashed
	SolidLine                   // solid
	DoubleLine                  // double
	GrovveLine                  // groove
	RidgeLine                   // ridge
	InsetLine                   // inset
	OutsetLine                  // outset
)

// Returns LineStyle in CSS syntax.
func (l LineStyle) String() string {
	switch l {
	case NoLine:
		return "none"
	case HiddenLine:
		return "hidden"
	case DottedLine:
		return "dotted"
	case DashedLine:
		return "dashed"
	case SolidLine:
		return "solid"
	case DoubleLine:
		return "double"
	case GrovveLine:
		return "groove"
	case RidgeLine:
		return "ridge"
	case InsetLine:
		return "inset"
	case OutsetLine:
		return "outset"
	}
	return fmt.Sprintf("<bad LineStyle %d>", l)
}

// Various predefined <line-width> values.
//
// Releavnt spec: https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
var (
	LineWidthThin   = values.LengthFromPx(css.NumFromInt(1)) // thin
	LineWidthMedium = values.LengthFromPx(css.NumFromInt(3)) // medium
	LineWidthThick  = values.LengthFromPx(css.NumFromInt(5)) // thick
)
