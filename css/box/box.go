// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Package box provide types and values for [CSS Box Model Module Level 3].
//
// [CSS Box Model Module Level 3]: https://www.w3.org/TR/css-box-3/
package box

import (
	"fmt"

	"github.com/inseo-oh/yw/css/values"
	"github.com/inseo-oh/yw/util"
)

// Margin represents value of [CSS margin property].
//
// Zero value for Margin means "auto".
//
// [CSS margin property]: https://www.w3.org/TR/css-box-3/#margin-physical
type Margin struct {
	Value values.LengthResolvable // nil means auto
}

// IsAuto reports whether it's auto margin.
func (m Margin) IsAuto() bool { return util.IsNil(m.Value) }

func (m Margin) String() string {
	if m.IsAuto() {
		return "auto"
	}
	return fmt.Sprintf("%v", m.Value)
}
