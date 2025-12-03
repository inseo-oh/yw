// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package csssyntax

import (
	"errors"

	"github.com/inseo-oh/yw/css/backgrounds"
	"github.com/inseo-oh/yw/css/values"
)

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
func (ts *tokenStream) parseLineStyle() (backgrounds.LineStyle, error) {
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
	return 0, errors.New("expected line-style")
}

// https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width
func (ts *tokenStream) parseLineWidth() (values.Length, error) {
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
	return values.Length{}, errors.New("expected line-width")
}
