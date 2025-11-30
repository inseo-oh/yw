package cssom

import (
	"github.com/inseo-oh/yw/css/props"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

type DocumentOrShadowRootData struct {
	Stylesheets []*Stylesheet
}

func DocumentOrShadowRootDataOf(node dom.Node) *DocumentOrShadowRootData {
	if util.IsNil(node.CssData()) {
		node.SetCssData(&DocumentOrShadowRootData{})
	}
	return node.CssData().(*DocumentOrShadowRootData)
}

type ElementData struct {
	ComputedStyleSet props.ComputedStyleSet
}

func ElementDataOf(elem dom.Element) *ElementData {
	if util.IsNil(elem.CssData()) {
		elem.SetCssData(&ElementData{})
	}
	return elem.CssData().(*ElementData)
}

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

func ComputedStyleSetSourceOf(elem dom.Element) ComputedStyleSetSource {
	return ComputedStyleSetSource{elem}
}
