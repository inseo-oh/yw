// Implementation of the CSS Box Model Module 3 (https://www.w3.org/TR/css-box-3/)
package box

import (
	"fmt"
	"yw/css/values"
	cm "yw/libcommon"
)

type Margin struct {
	Value values.LengthResolvable // nil means auto
}

func (m Margin) IsAuto() bool { return cm.IsNil(m.Value) }
func (m Margin) String() string {
	if m.IsAuto() {
		return "auto"
	}
	return fmt.Sprintf("%v", m.Value)
}
