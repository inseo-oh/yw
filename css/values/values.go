// Implementation of the CSS Values and Units Module Level 3 (https://www.w3.org/TR/css-values-3/)
package values

import (
	"fmt"
	"log"
	"yw/css"
)

type LengthResolvable interface {
	AsLength(containerSize css.Num) Length
	String() string
}

// https://www.w3.org/TR/css-values-3/#length-value
type Length struct {
	value css.Num
	unit  LengthUnit
}

type LengthUnit uint8

func LengthFromPx(px css.Num) Length {
	return Length{px, Px}
}

func (l Length) String() string                        { return fmt.Sprintf("%v%v", l.value, l.unit) }
func (l Length) AsLength(containerSize css.Num) Length { return l }
func (l Length) ToPx(containerSize css.Num) float64 {
	switch l.unit {
	case Px:
		return l.value.ToFloat()
	case Em:
		return containerSize.ToFloat() * l.value.ToFloat()
	case Pt:
		// STUB -- For now we treat pt and px as the same thing.
		return l.value.ToFloat()
	default:
		log.Panicf("TODO: %v", l)
	}
	return 0
}

const (
	//==========================================================================
	// https://www.w3.org/TR/css-values-3/#relative-lengths
	//==========================================================================

	Em = LengthUnit(iota)
	Ex
	Ch
	Rem
	Vw
	Vh
	Vmin
	Vmax

	//==========================================================================
	// https://www.w3.org/TR/css-values-3/#absolute-lengths
	//==========================================================================

	Cm
	Mm
	Q
	Pc
	Pt
	Px
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

// https://www.w3.org/TR/css-values-3/#percentage-value
type Percentage struct {
	value css.Num
}

func (len Percentage) String() string { return fmt.Sprintf("%v%%", len.value) }

func (len Percentage) AsLength(containerSize css.Num) Length { panic("TODO") }
