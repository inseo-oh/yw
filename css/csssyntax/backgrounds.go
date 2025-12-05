// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/backgrounds"
	"github.com/inseo-oh/yw/css/values"
)

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
func (ts *tokenStream) parseLineStyle() (res backgrounds.LineStyle, err error) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return backgrounds.NoLine, nil
	} else if err := ts.consumeIdentTokenWith("hidden"); err == nil {
		return backgrounds.HiddenLine, nil
	} else if err := ts.consumeIdentTokenWith("dotted"); err == nil {
		return backgrounds.DottedLine, nil
	} else if err := ts.consumeIdentTokenWith("dashed"); err == nil {
		return backgrounds.DashedLine, nil
	} else if err := ts.consumeIdentTokenWith("solid"); err == nil {
		return backgrounds.SolidLine, nil
	} else if err := ts.consumeIdentTokenWith("double"); err == nil {
		return backgrounds.DoubleLine, nil
	} else if err := ts.consumeIdentTokenWith("groove"); err == nil {
		return backgrounds.GrovveLine, nil
	} else if err := ts.consumeIdentTokenWith("ridge"); err == nil {
		return backgrounds.RidgeLine, nil
	} else if err := ts.consumeIdentTokenWith("inset"); err == nil {
		return backgrounds.InsetLine, nil
	} else if err := ts.consumeIdentTokenWith("outset"); err == nil {
		return backgrounds.OutsetLine, nil
	}
	return res, fmt.Errorf("%s: expected line-style", ts.errorHeader())
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
func (ts *tokenStream) parseLineWidth() (res values.Length, err error) {
	if err := ts.consumeIdentTokenWith("thin"); err == nil {
		return backgrounds.LineWidthThin, nil
	}
	if err := ts.consumeIdentTokenWith("medium"); err == nil {
		return backgrounds.LineWidthMedium, nil
	}
	if err := ts.consumeIdentTokenWith("thick"); err == nil {
		return backgrounds.LineWidthThick, nil
	}
	if len, err := ts.parseLength(true); err == nil {
		return len, nil
	}
	return res, fmt.Errorf("%s: expected line-width", ts.errorHeader())
}
