// Package display provides types and values for [CSS2 9.5 Floats].
//
// [CSS2 9.5 Floats]: https://www.w3.org/TR/CSS2/visuren.html#floats
package float

import (
	"fmt"
)

// Float represents value of [CSS float] property.
//
// [CSS float]: https://www.w3.org/TR/CSS2/visuren.html#propdef-float
type Float uint8

const (
	None  Float = iota // float: none
	Left               // float: left
	Right              // float: right
)

func (f Float) String() string {
	switch f {
	case None:
		return "none"
	case Left:
		return "left"
	case Right:
		return "right"
	}
	return fmt.Sprintf("<bad Float type %d>", f)
}
