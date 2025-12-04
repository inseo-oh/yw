// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package cssom

import (
	"image/color"

	"github.com/inseo-oh/yw/css/csscolor"
	"github.com/inseo-oh/yw/css/props"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

// DocumentOrShadowRootData holds CSS-specific data for DOM Document or ShadowRoot nodes.
type DocumentOrShadowRootData struct {
	Stylesheets []*Stylesheet // List of stylesheets
}

// DocumentOrShadowRootDataOf returns DocumentOrShadowRootData for given node.
func DocumentOrShadowRootDataOf(node dom.Node) *DocumentOrShadowRootData {
	if util.IsNil(node.CssData()) {
		node.SetCssData(&DocumentOrShadowRootData{})
	}
	return node.CssData().(*DocumentOrShadowRootData)
}

// ElementData holds CSS-specific data for DOM Element nodes.
type ElementData struct {
	ComputedStyleSet props.ComputedStyleSet
}

// ElementDataOf returns DocumentOrShadowRootData for given node.
func ElementDataOf(elem dom.Element) *ElementData {
	if util.IsNil(elem.CssData()) {
		elem.SetCssData(&ElementData{})
	}
	return elem.CssData().(*ElementData)
}

// ComputedStyleSetSource is a wrapper around DOM Element that is used by
// other CSS packages to read/write to the ComputedStyleSet.
type ComputedStyleSetSource struct {
	elem dom.Element
}

func (src ComputedStyleSetSource) ComputedStyleSet() *props.ComputedStyleSet {
	return &ElementDataOf(src.elem).ComputedStyleSet
}
func (src ComputedStyleSetSource) ParentSource() props.ComputedStyleSetSource {
	parent := src.elem.Parent()
	if util.IsNil(parent) {
		return nil
	}
	if parentElem, ok := parent.(dom.Element); ok {
		return ComputedStyleSetSourceOf(parentElem)
	}
	return nil
}
func (src ComputedStyleSetSource) CurrentColor() color.RGBA {
	colorVal := src.ComputedStyleSet().Color()
	if colorVal.Type == csscolor.CurrentColor {
		parentSrc := src.ParentSource()
		if util.IsNil(parentSrc) {
			return props.DescriptorsMap["color"].Initial.(csscolor.Color).ToRgba(color.RGBA{})
		}
		return parentSrc.CurrentColor()
	}
	return colorVal.ToRgba(color.RGBA{})
}

// ComputedStyleSetSourceOf creates new ComputedStyleSetSource for given elem.
func ComputedStyleSetSourceOf(elem dom.Element) ComputedStyleSetSource {
	return ComputedStyleSetSource{elem}
}
