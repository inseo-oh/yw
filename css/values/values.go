// Implementation of the CSS Values and Units Module Level 3 (https://www.w3.org/TR/css-values-3/)
package values

import (
	"fmt"
	"log"

	"github.com/inseo-oh/yw/css"
)

type LengthResolvable interface {
	AsLength(containerSize css.Num) Length
	String() string
}

// https://www.w3.org/TR/css-values-3/#length-value
type Length struct {
	Value css.Num
	Unit  LengthUnit
}

type LengthUnit uint8

func LengthFromPx(px css.Num) Length {
	return Length{px, Px}
}

func (l Length) String() string                        { return fmt.Sprintf("%v%v", l.Value, l.Unit) }
func (l Length) AsLength(containerSize css.Num) Length { return l }
func (l Length) ToPx(parentFontSize css.Num) float64 {
	switch l.Unit {
	case Px:
		return l.Value.ToFloat()
	case Em:
		return parentFontSize.ToFloat() * l.Value.ToFloat()
	case Pt:
		// STUB -- For now we treat pt and px as the same thing.
		return l.Value.ToFloat()
	default:
		log.Panicf("TODO: %v", l)
	}
	return 0
}

const (
	//==========================================================================
	// https://www.w3.org/TR/css-values-3/#relative-lengths
	//==========================================================================

	Em LengthUnit = iota
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
	Value css.Num
}

func (len Percentage) String() string { return fmt.Sprintf("%v%%", len.Value) }

func (len Percentage) AsLength(containerSize css.Num) Length { panic("TODO") }
