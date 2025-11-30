// Implementation of the CSS Box Model Module 3 (https://www.w3.org/TR/css-box-3/)
package box

import (
	"fmt"

	"github.com/inseo-oh/yw/css/values"
	"github.com/inseo-oh/yw/util"
)

type Margin struct {
	Value values.LengthResolvable // nil means auto
}

func (m Margin) IsAuto() bool { return util.IsNil(m.Value) }
func (m Margin) String() string {
	if m.IsAuto() {
		return "auto"
	}
	return fmt.Sprintf("%v", m.Value)
}
