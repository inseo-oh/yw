// Implementation of the CSS Sizing Module Level 3 (https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/)
package sizing

import (
	"fmt"
	"log"
	"yw/css"
	"yw/css/values"
)

// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#sizing-values
type Size struct {
	tp   SizeType
	size values.LengthResolvable
}
type SizeType uint8

const (
	NoneSize   = SizeType(iota) // none
	Auto                        // auto
	MinContent                  // min-content
	MaxContent                  // max-content
	FitContent                  // fit-content
	ManualSize                  // Manually specified size
)

func (s Size) String() string {
	switch s.tp {
	case NoneSize:
		return "none"
	case Auto:
		return "auto"
	case MinContent:
		return "min-content"
	case MaxContent:
		return "max-content"
	case FitContent:
		return fmt.Sprintf("fit-content(%v)", s.size)
	case ManualSize:
		return s.size.String()
	}
	return fmt.Sprintf("<bad Size type %v>", s.tp)
}

func (s Size) ComputeUsedValue(parentSize css.Num) values.Length {
	switch s.tp {
	case NoneSize:
		panic("TODO: NoneSize")
	case Auto:
		panic("Auto sizes must be calculated by caller")
	case MinContent:
		panic("TODO: MinContent")
	case MaxContent:
		panic("TODO: MaxContent")
	case FitContent:
		panic("TODO: FitContent")
	case ManualSize:
		return s.size.AsLength(parentSize)
	}
	log.Panicf("<bad Size type %v>", s.tp)
	return values.Length{}
}
