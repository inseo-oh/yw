// Implementation of the CSS Sizing Module Level 3 (https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/)
package libhtml

import (
	"fmt"
	"log"

	cm "github.com/inseo-oh/yw/util"
)

// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#sizing-values
type css_size_value struct {
	tp   css_size_value_type
	size css_length_resolvable
}
type css_size_value_type uint8

const (
	css_size_value_type_none = css_size_value_type(iota)
	css_size_value_type_auto
	css_size_value_type_min_content
	css_size_value_type_max_content
	css_size_value_type_fit_content
	css_size_value_type_manual
)

func css_size_value_auto() css_size_value {
	return css_size_value{css_size_value_type_auto, nil}
}
func css_size_value_none() css_size_value {
	return css_size_value{css_size_value_type_none, nil}
}

func (s css_size_value) String() string {
	switch s.tp {
	case css_size_value_type_none:
		return "none"
	case css_size_value_type_auto:
		return "auto"
	case css_size_value_type_min_content:
		return "min-content"
	case css_size_value_type_max_content:
		return "max-content"
	case css_size_value_type_fit_content:
		return fmt.Sprintf("fit-content(%v)", s.size)
	case css_size_value_type_manual:
		return s.size.String()
	}
	return fmt.Sprintf("unregognized css_size_value type %v", s.tp)
}

func (s css_size_value) compute_used_value(parent_size css_number) css_length {
	switch s.tp {
	case css_size_value_type_none:
		panic("TODO: css_size_value_type_none")
	case css_size_value_type_auto:
		panic("Auto sizes must be calculated by caller")
	case css_size_value_type_min_content:
		panic("TODO: css_size_value_type_min_content")
	case css_size_value_type_max_content:
		panic("TODO: css_size_value_type_max_content")
	case css_size_value_type_fit_content:
		panic("TODO: css_size_value_type_fit_content")
	case css_size_value_type_manual:
		return s.size.as_length(parent_size)
	}
	log.Panicf("unregognized css_size_value type %v", s.tp)
	return css_length{}
}

func (ts *css_token_stream) parse_size_value_impl(accept_auto, accept_none bool) (css_size_value, bool) {
	if accept_auto {
		if tk := ts.consume_ident_token_with("auto"); !cm.IsNil(tk) {
			return css_size_value{css_size_value_type_auto, nil}, true
		}
	}
	if accept_none {
		if tk := ts.consume_ident_token_with("none"); !cm.IsNil(tk) {
			return css_size_value{css_size_value_type_auto, nil}, true
		}
	}
	if tk := ts.consume_ident_token_with("min-content"); !cm.IsNil(tk) {
		return css_size_value{css_size_value_type_min_content, nil}, true
	}
	if tk := ts.consume_ident_token_with("max-content"); !cm.IsNil(tk) {
		return css_size_value{css_size_value_type_max_content, nil}, true
	}
	if tk := ts.consume_ast_function_with("fit-content"); !cm.IsNil(tk) {
		ts := css_token_stream{tokens: tk.value}
		var size css_length_resolvable
		if v, err := ts.parse_length_or_percentage(true); !cm.IsNil(v) {
			size = v
		} else if err != nil {
			return css_size_value{}, false
		}
		if !ts.is_end() {
			return css_size_value{}, false
		}
		return css_size_value{css_size_value_type_fit_content, size}, true
	}
	if v, err := ts.parse_length_or_percentage(true); !cm.IsNil(v) {
		return css_size_value{css_size_value_type_manual, v}, true
	} else if err != nil {
		return css_size_value{}, false
	}
	return css_size_value{}, false
}
func (ts *css_token_stream) parse_size_value_or_auto() (css_size_value, bool) {
	return ts.parse_size_value_impl(true, false)
}
func (ts *css_token_stream) parse_size_value_or_none() (css_size_value, bool) {
	return ts.parse_size_value_impl(false, true)
}
