// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package css

import (
	"fmt"

	"github.com/inseo-oh/yw/util"
)

// CSS numeric value storing either integer or float.
type Num struct {
	Type  NumType // Type of number
	Value any     // Actual value of the number (either int64 or float64)
}

// Type of [Num]
type NumType uint8

const (
	NumTypeInt NumType = iota
	NumTypeFloat
)

// Converts to integer value
func (n Num) ToInt() int64 {
	if util.IsNil(n.Value) {
		return 0
	}
	if n.Type == NumTypeFloat {
		return int64(n.ToFloat())
	} else {
		return n.Value.(int64)
	}
}

// Converts to float value
func (n Num) ToFloat() float64 {
	if util.IsNil(n.Value) {
		return 0
	}
	if n.Type == NumTypeInt {
		return float64(n.ToInt())
	} else {
		return n.Value.(float64)
	}
}

// Compares two numeric value. If both values are integers, integer comparison is used.
// Otherwise, both values are converted to float before comparing them.
func (n Num) Equals(other Num) bool {
	if util.IsNil(n.Value) {
		return NumFromInt(0).Equals(other)
	}
	isFloat := (n.Type == NumTypeFloat) || (other.Type == NumTypeFloat)
	if isFloat {
		return n.ToFloat() == other.ToFloat()
	} else {
		return n.ToInt() == other.ToInt()
	}
}

func (n Num) Clamp(min, max Num) Num {
	if util.IsNil(n.Value) {
		return NumFromInt(0).Clamp(min, max)
	}
	// Using float all the time is probably fine, but let's avoid it if we can.
	isFloat := (n.Type == NumTypeFloat) || (min.Type == NumTypeFloat) || (max.Type == NumTypeFloat)
	if isFloat {
		if n.ToFloat() < min.ToFloat() {
			return min
		} else if max.ToFloat() < n.ToFloat() {
			return max
		}
	} else {
		if n.ToInt() < min.ToInt() {
			return min
		} else if max.ToInt() < n.ToInt() {
			return max
		}
	}
	return n
}

// Num represents the number in CSS form
func (n Num) String() string {
	if util.IsNil(n.Value) {
		return NumFromInt(0).String()
	}
	if n.Type == NumTypeFloat {
		return fmt.Sprintf("%f", n.ToFloat())
	}
	return fmt.Sprintf("%d", n.ToInt())
}

// Creates a new number from integer
func NumFromInt(v int64) Num {
	return Num{NumTypeInt, v}
}

// Creates a new number from floating point
func NumFromFloat(v float64) Num {
	return Num{NumTypeFloat, v}
}
