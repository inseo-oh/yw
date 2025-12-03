// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

// Package sizing provides types and values for [CSS Sizing Module Level 3]
//
// [CSS Sizing Module Level 3]: https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/
package sizing

import (
	"fmt"
	"log"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/values"
)

// Size represents a [CSS size] value.
//
// [CSS size]: https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#sizing-values
type Size struct {
	Type SizeType                // Type of size value
	Size values.LengthResolvable // Manual size value when Type is ManualSize
}

// Type of [Size] value.
type SizeType uint8

const (
	NoneSize   SizeType = iota // none
	Auto                       // auto
	MinContent                 // min-content
	MaxContent                 // max-content
	FitContent                 // fit-content
	ManualSize                 // Manually specified size
)

func (s Size) String() string {
	switch s.Type {
	case NoneSize:
		return "none"
	case Auto:
		return "auto"
	case MinContent:
		return "min-content"
	case MaxContent:
		return "max-content"
	case FitContent:
		return fmt.Sprintf("fit-content(%v)", s.Size)
	case ManualSize:
		return s.Size.String()
	}
	return fmt.Sprintf("<bad Size type %v>", s.Type)
}

// ComputeUsedValue computes length value for the size.
func (s Size) ComputeUsedValue(parentSize css.Num) values.Length {
	switch s.Type {
	case NoneSize:
		panic("TODO: NoneSize")
	case Auto:
		panic("Auto sizes must be calculated by caller")
	case MinContent:
		panic("TODO: MinContent")
	case MaxContent:
		panic("TODO: MaxContent")
	case FitContent:
		panic("TODO: FitContent")
	case ManualSize:
		return s.Size.AsLength(parentSize)
	}
	log.Panicf("<bad Size type %v>", s.Type)
	return values.Length{}
}
