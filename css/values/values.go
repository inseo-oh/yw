// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Implementation of the CSS Values and Units Module Level 3 (https://www.w3.org/TR/css-values-3/)
package values

import (
	"fmt"
	"log"

	"github.com/inseo-oh/yw/css"
)

// LengthResolvable represents a value that can be resolved to a [Length].
//
// [Length] and [Percentage] implements this type.
type LengthResolvable interface {
	AsLength(containerSize func() css.Num) Length
	String() string
}

// Length is a CSS number with [LengthUnit].
//
// https://www.w3.org/TR/css-values-3/#length-value
type Length struct {
	Value css.Num
	Unit  LengthUnit
}

// LengthFromPx creates a Px unit length from a number.
//
// TODO(ois): Wouldn't this be more useful if this just accepted float value?
func LengthFromPx(px css.Num) Length {
	return Length{px, Px}
}

func (l Length) String() string                               { return fmt.Sprintf("%v%v", l.Value, l.Unit) }
func (l Length) AsLength(containerSize func() css.Num) Length { return l }
func (l Length) ToPx(fontSize func() css.Num) float64 {
	switch l.Unit {
	case Px:
		return l.Value.ToFloat()
	case Em:
		return fontSize().ToFloat() * l.Value.ToFloat()
	case Pt:
		// STUB -- For now we treat pt and px as the same thing.
		return l.Value.ToFloat()
	default:
		log.Panicf("TODO: %v", l)
	}
	return 0
}

// Unit for [Length]
type LengthUnit uint8

const (
	//==========================================================================
	// Relative lengths
	//
	// https://www.w3.org/TR/css-values-3/#relative-lengths
	//==========================================================================

	Em   LengthUnit = iota // em
	Ex                     // ex
	Ch                     // ch
	Rem                    // rem
	Vw                     // vw
	Vh                     // vh
	Vmin                   // vmin
	Vmax                   // vmax

	//==========================================================================
	// Absolute lengths
	//
	// https://www.w3.org/TR/css-values-3/#absolute-lengths
	//==========================================================================

	Cm // cm
	Mm // mm
	Q  // q
	Pc // pc
	Pt // pt
	Px // px
)

func (u LengthUnit) String() string {
	switch u {
	case Em:
		return "em"
	case Ex:
		return "ex"
	case Ch:
		return "ch"
	case Rem:
		return "rem"
	case Vw:
		return "vw"
	case Vh:
		return "vh"
	case Vmin:
		return "vmin"
	case Vmax:
		return "vmax"
	case Cm:
		return "cm"
	case Mm:
		return "mm"
	case Q:
		return "q"
	case Pc:
		return "pc"
	case Pt:
		return "pt"
	case Px:
		return "px"
	}
	return fmt.Sprintf("<bad LengthUnit %d>", u)
}

// Percentage is a CSS number with %.
//
// https://www.w3.org/TR/css-values-3/#percentage-value
type Percentage struct {
	Value css.Num
}

func (len Percentage) String() string { return fmt.Sprintf("%v%%", len.Value) }

func (len Percentage) AsLength(containerSize func() css.Num) Length {
	return LengthFromPx(css.NumFromFloat((len.Value.ToFloat() * containerSize().ToFloat()) / 100))
}
