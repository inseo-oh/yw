// Implementation of the CSS Box Model Module 3 (https://www.w3.org/TR/css-box-3/)
package libhtml

import (
	"fmt"

	cm "github.com/inseo-oh/yw/util"
)

type css_margin struct {
	v css_length_resolvable // nil means auto
}

func (m css_margin) is_auto() bool { return cm.IsNil(m.v) }
func (m css_margin) String() string {
	if m.is_auto() {
		return "auto"
	}
	return fmt.Sprintf("%v", m.v)
}

func (ts *css_token_stream) parse_margin() (css_margin, bool) {
	if v, err := ts.parse_length_or_percentage(true); !cm.IsNil(v) {
		return css_margin{v}, true
	} else if err != nil {
		return css_margin{}, false
	}
	if ts.consume_ident_token_with("auto") != nil {
		return css_margin{nil}, true
	}
	return css_margin{}, false
}
func (ts *css_token_stream) parse_padding() (css_length_resolvable, bool) {
	v, _ := ts.parse_length_or_percentage(true)
	if cm.IsNil(v) {
		return nil, false
	}
	if len, ok := v.(css_length); ok && len.value.to_int() < 0 {
		return nil, false
	}
	return v, true
}
