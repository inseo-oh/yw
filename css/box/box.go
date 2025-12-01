// Package box provide types and values for CSS Box Model Module Level 3
//
// Spec: https://www.w3.org/TR/css-box-3/
package box

import (
	"fmt"

	"github.com/inseo-oh/yw/css/values"
	"github.com/inseo-oh/yw/util"
)

// Represents margin of a box edge. Zero value for Margin means "auto".
//
//   - Spec: https://www.w3.org/TR/css-box-3/#margin-physical
//   - MDN: https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Box_model/Introduction#margin_area
type Margin struct {
	Value values.LengthResolvable // nil means auto
}

// IsAuto reports whether it's auto margin.
func (m Margin) IsAuto() bool { return util.IsNil(m.Value) }

func (m Margin) String() string {
	if m.IsAuto() {
		return "auto"
	}
	return fmt.Sprintf("%v", m.Value)
}
