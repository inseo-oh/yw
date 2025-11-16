// Implementation of the CSS Sizing Module Level 3 (https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/)
package libhtml

import (
	"errors"
	"fmt"
	"log"
	cm "yw/libcommon"
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

func init() {
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-width
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-height
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-max-width
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-max-height
	parse_size_value := func(ts *css_token_stream) (*css_size_value, error) {
		if tk := ts.consume_ident_token_with("auto"); !cm.IsNil(tk) {
			return &css_size_value{css_size_value_type_auto, nil}, nil
		} else if tk := ts.consume_ident_token_with("none"); !cm.IsNil(tk) {
			return &css_size_value{css_size_value_type_auto, nil}, nil
		} else if tk := ts.consume_ident_token_with("min-content"); !cm.IsNil(tk) {
			return &css_size_value{css_size_value_type_min_content, nil}, nil
		} else if tk := ts.consume_ident_token_with("max-content"); !cm.IsNil(tk) {
			return &css_size_value{css_size_value_type_max_content, nil}, nil
		} else if tk := ts.consume_ast_function_with("fit-content"); !cm.IsNil(tk) {
			ts := css_token_stream{tokens: tk.value}
			var size css_length_resolvable
			if v, err := ts.parse_length_or_percentage(true); !cm.IsNil(v) {
				size = v
			} else if err != nil {
				return nil, err
			}
			if !ts.is_end() {
				return nil, errors.New("unexpected junk at the end of function")
			}
			return &css_size_value{css_size_value_type_fit_content, size}, nil
		} else if v, err := ts.parse_length_or_percentage(true); !cm.IsNil(v) {
			return &css_size_value{css_size_value_type_manual, v}, nil
		} else if err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	}
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#preferred-size-properties
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-width
	css_property_descriptors_map["width"] = css_property_descriptor{
		initial: css_size_value{css_size_value_type_auto, nil},
		apply_func: func(dest *css_computed_style_set, value any) {
			v := value.(css_size_value)
			dest.width = &v
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			if v, err := parse_size_value(ts); v != nil {
				if v.tp == css_size_value_type_none {
					return nil, errors.New("size value 'none' is not accepted in this context")
				}
				return *v, nil
			} else {
				return nil, err
			}
		},
	}
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-height
	css_property_descriptors_map["height"] = css_property_descriptor{
		initial: css_size_value{css_size_value_type_auto, nil},
		apply_func: func(dest *css_computed_style_set, value any) {
			v := value.(css_size_value)
			dest.height = &v
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			if v, err := parse_size_value(ts); v != nil {
				if v.tp == css_size_value_type_none {
					return nil, errors.New("size value 'none' is not accepted in this context")
				}
				return *v, nil
			} else {
				return nil, err
			}
		},
	}

	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#min-size-properties
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-min-width
	css_property_descriptors_map["min-width"] = css_property_descriptor{
		initial: css_size_value{css_size_value_type_auto, nil},
		apply_func: func(dest *css_computed_style_set, value any) {
			panic("TODO")
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			if v, err := parse_size_value(ts); v != nil {
				if v.tp == css_size_value_type_none {
					return nil, errors.New("size value 'none' is not accepted in this context")
				}
				return *v, nil
			} else {
				return nil, err
			}
		},
	}
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#propdef-min-height
	css_property_descriptors_map["min-height"] = css_property_descriptor{
		initial: css_size_value{css_size_value_type_auto, nil},
		apply_func: func(dest *css_computed_style_set, value any) {
			panic("TODO")
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			if v, err := parse_size_value(ts); v != nil {
				if v.tp == css_size_value_type_none {
					return nil, errors.New("size value 'none' is not accepted in this context")
				}
				return *v, nil
			} else {
				return nil, err
			}
		},
	}
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#max-size-properties
	//==========================================================================
	css_property_descriptors_map["max-width"] = css_property_descriptor{
		initial: css_size_value{css_size_value_type_auto, nil},
		apply_func: func(dest *css_computed_style_set, value any) {
			panic("TODO")
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			if v, err := parse_size_value(ts); v != nil {
				if v.tp == css_size_value_type_auto {
					return nil, errors.New("size value 'auto' is not accepted in this context")
				}
				return *v, nil
			} else {
				return nil, err
			}
		},
	}
	css_property_descriptors_map["max-height"] = css_property_descriptor{
		initial: css_size_value{css_size_value_type_auto, nil},
		apply_func: func(dest *css_computed_style_set, value any) {
			panic("TODO")
		},
		parse_func: func(ts *css_token_stream) (css_property_value, error) {
			if v, err := parse_size_value(ts); v != nil {
				if v.tp == css_size_value_type_auto {
					return nil, errors.New("size value 'auto' is not accepted in this context")
				}
				return *v, nil
			} else {
				return nil, err
			}
		},
	}
}
