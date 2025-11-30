package es

import (
	"fmt"
	"log"
)

type Value struct {
	Type  ValueType
	Value any
}
type ValueType uint8

const (
	ValueTypeNull = ValueType(iota)
	ValueTypeUndefined
	ValueTypeBoolean
	ValueTypeNumber
	ValueTypeString
)

func NewNullValue() Value             { return Value{ValueTypeNull, nil} }
func NewUndefinedValue() Value        { return Value{ValueTypeUndefined, nil} }
func NewBooleanValue(v bool) Value    { return Value{ValueTypeBoolean, v} }
func NewNumberValueF(v float64) Value { return Value{ValueTypeNumber, v} }
func NewNumberValueI(v int64) Value   { return Value{ValueTypeNumber, v} }
func NewStringValue(v string) Value   { return Value{ValueTypeString, v} }

func (v Value) ExpectNumberF() float64 {
	if v.Type != ValueTypeNumber {
		panic("the value is not a number")
	}
	res, ok := v.Value.(float64)
	if !ok {
		return float64(v.Value.(int64))
	}
	return res
}
func (v Value) ExpectNumberI() int64 {
	if v.Type != ValueTypeNumber {
		panic("the value is not a number")
	}
	res, ok := v.Value.(int64)
	if !ok {
		return int64(v.Value.(float64))
	}
	return res
}

func (v Value) ExpectBoolean() bool {
	if v.Type != ValueTypeBoolean {
		panic("the value is not a boolean")
	}
	return v.Value.(bool)
}
func (v Value) String() string {
	switch v.Type {
	case ValueTypeNull:
		return "<null>"
	case ValueTypeUndefined:
		return "<undefined>"
	case ValueTypeBoolean:
		return fmt.Sprintf("<boolean:%v>", v.Value)
	case ValueTypeNumber:
		if _, ok := v.Value.(float64); ok {
			return fmt.Sprintf("<number:float64(%v)>", v.Value)
		} else if _, ok := v.Value.(int64); ok {
			return fmt.Sprintf("<number:int64(%v)>", v.Value)
		} else {
			return fmt.Sprintf("<number:???(%v)>", v.Value)
		}
	case ValueTypeString:
		return fmt.Sprintf("<number:%v>", v.Value)
	}
	return fmt.Sprintf("<?unknown type %d, value=%v>", v.Type, v.Value)
}

type Reference struct {
	// STUB
}

func GetValue(v any) Value {
	ref, isRef := v.(Reference)
	if !isRef {
		return v.(Value)
	}
	log.Println(ref)
	panic("TODO")
}
