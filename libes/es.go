package libes

import (
	"fmt"
	"log"
)

type es_value struct {
	tp    es_value_type
	value any
}
type es_value_type uint8

const (
	es_value_type_null = es_value_type(iota)
	es_value_type_undefined
	es_value_type_boolean
	es_value_type_number
	es_value_type_string
)

func es_make_null_value() es_value              { return es_value{es_value_type_null, nil} }
func es_make_undefined_value() es_value         { return es_value{es_value_type_undefined, nil} }
func es_make_boolean_value(v bool) es_value     { return es_value{es_value_type_boolean, v} }
func es_make_number_value_f(v float64) es_value { return es_value{es_value_type_number, v} }
func es_make_number_value_i(v int64) es_value   { return es_value{es_value_type_number, v} }
func es_make_string_value(v string) es_value    { return es_value{es_value_type_string, v} }

func (v es_value) expect_number_f() float64 {
	if v.tp != es_value_type_number {
		panic("the value is not a number")
	}
	res, ok := v.value.(float64)
	if !ok {
		return float64(v.value.(int64))
	}
	return res
}
func (v es_value) expect_number_i() int64 {
	if v.tp != es_value_type_number {
		panic("the value is not a number")
	}
	res, ok := v.value.(int64)
	if !ok {
		return int64(v.value.(float64))
	}
	return res
}

func (v es_value) expect_boolean() bool {
	if v.tp != es_value_type_boolean {
		panic("the value is not a boolean")
	}
	return v.value.(bool)
}
func (v es_value) String() string {
	switch v.tp {
	case es_value_type_null:
		return "<null>"
	case es_value_type_undefined:
		return "<undefined>"
	case es_value_type_boolean:
		return fmt.Sprintf("<boolean:%v>", v.value)
	case es_value_type_number:
		if _, ok := v.value.(float64); ok {
			return fmt.Sprintf("<number:float64(%v)>", v.value)
		} else if _, ok := v.value.(int64); ok {
			return fmt.Sprintf("<number:int64(%v)>", v.value)
		} else {
			return fmt.Sprintf("<number:???(%v)>", v.value)
		}
	case es_value_type_string:
		return fmt.Sprintf("<number:%v>", v.value)
	}
	return fmt.Sprintf("<?unknown type %d, value=%v>", v.tp, v.value)
}

type es_reference struct {
	// STUB
}

func es_get_value(v any) es_value {
	ref, is_ref := v.(es_reference)
	if !is_ref {
		return v.(es_value)
	}
	log.Println(ref)
	panic("TODO")
}
