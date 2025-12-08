// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// Package es implements [ECMA-262 standard], a.k.a. JavaScript.
// Documentation (and perheaps certain parts of API) refer to it as ES(Short for ECMAScript).
//
// # Relationship with sub-packages
//
// This package never references any of sub-packages, but sub-packages can
// freely import this package.
//
// [ECMA-262 standard]: https://tc39.es/ecma262/
package es

import (
	"fmt"
	"log"
)

// Value represents a [ES value].
//
// Numeric values are stored as either integer or float.
//
// [ES value]: https://tc39.es/ecma262/#sec-ecmascript-language-types
type Value struct {
	Type  ValueType
	Value any
}

// ValueType represents type of [Value].
type ValueType uint8

const (
	ValueTypeNull      ValueType = iota // Null
	ValueTypeUndefined                  // Undefined
	ValueTypeBoolean                    // Boolean value
	ValueTypeNumber                     // Number value
	ValueTypeString                     // String value
)

// NewNullValue creates a new null value.
func NewNullValue() Value { return Value{ValueTypeNull, nil} }

// NewUndefinedValue creates a new undefined value.
func NewUndefinedValue() Value { return Value{ValueTypeUndefined, nil} }

// NewBooleanValue creates a new boolean value.
func NewBooleanValue(v bool) Value { return Value{ValueTypeBoolean, v} }

// NewNumberValueF creates a new number value from float64.
func NewNumberValueF(v float64) Value { return Value{ValueTypeNumber, v} }

// NewNumberValueI creates a new number value from int64.
func NewNumberValueI(v int64) Value { return Value{ValueTypeNumber, v} }

// NewStringValue creates a new string value.
func NewStringValue(v string) Value { return Value{ValueTypeString, v} }

// ExpectNumberF returns number in float64, or panics if the value isn't number.
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

// ExpectNumberI returns number in int64, or panics if the value isn't number.
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

// ExpectNumberI returns number in boolean, or panics if the value isn't boolean.
func (v Value) ExpectBoolean() bool {
	if v.Type != ValueTypeBoolean {
		panic("the value is not a boolean")
	}
	return v.Value.(bool)
}

// String returns debug string for the value.
//
// TODO(ois): We probably could make Value.String represent value using ES syntax.
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

// Reference represents a [ES ReferenceRecord].
//
// TODO(ois): Reference is currently just STUB.
//
// [ES ReferenceRecord]: https://tc39.es/ecma262/#sec-reference-record-specification-type
type Reference struct{}

// GetValue takes either [Value] or [Reference] and returns a [Value].
//
//   - If input is [Value], input is returned.
//   - If input is [Reference], it tries to resolve the reference, and returns
//     the result.
//
// Spec: https://tc39.es/ecma262/#sec-getvalue
func GetValue(v any) Value {
	ref, isRef := v.(Reference)
	if !isRef {
		return v.(Value)
	}
	log.Println(ref)
	panic("TODO")
}
