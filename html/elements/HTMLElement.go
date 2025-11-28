package elements

import (
	"slices"
	"yw/dom"
)

// ------------------------------------------------------------------------------
// HTMLElement
// ------------------------------------------------------------------------------

type HTMLElement interface {
	dom.Element
	IsFormAssociatedCustomElement() bool
	IsFormAssociatedElement() bool
	IsFormListedElement() bool
	IsFormSubmittableElement() bool
	IsFormResettableElement() bool
	IsFormAutocapitalizeAndAutocorrectInheritingElement() bool
	IsFormLabelableElement() bool
	ContributesScriptBlockingStylesheet() bool
}
type htmlElementImpl struct{ dom.Element }

func NewHTMLElement(options dom.ElementCreationCommonOptions) HTMLElement {
	return htmlElementImpl{dom.NewElement(options)}
}

// https://html.spec.whatwg.org/multipage/custom-elements.html#form-associated-custom-element
func (elem htmlElementImpl) IsFormAssociatedCustomElement() bool {
	// STUB
	return false
}

// https://html.spec.whatwg.org/multipage/forms.html#form-associated-element
func (elem htmlElementImpl) IsFormAssociatedElement() bool {
	html_elems := []string{
		"button", "fieldset", "input", "object", "output", "select", "textarea", "img",
	}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(html_elems, elem.IsHtmlElement)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-listed
func (elem htmlElementImpl) IsFormListedElement() bool {
	html_elems := []string{
		"button", "fieldset", "input", "object", "output", "select", "textarea",
	}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(html_elems, elem.IsHtmlElement)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-submit
func (elem htmlElementImpl) IsFormSubmittableElement() bool {
	html_elems := []string{"button", "input", "select", "textarea"}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(html_elems, elem.IsHtmlElement)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-reset
func (elem htmlElementImpl) IsFormResettableElement() bool {
	html_elems := []string{"input", "output", "select", "textarea"}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(html_elems, elem.IsHtmlElement)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-autocapitalize
func (elem htmlElementImpl) IsFormAutocapitalizeAndAutocorrectInheritingElement() bool {
	html_elems := []string{
		"button", "fieldset", "input", "output", "select", "textarea",
	}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(html_elems, elem.IsHtmlElement)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-label
func (elem htmlElementImpl) IsFormLabelableElement() bool {
	html_elems := []string{
		"button", "meter", "output", "progress", "select", "textarea",
	}
	if elem.IsFormAssociatedCustomElement() {
		return true
	}
	if slices.ContainsFunc(html_elems, elem.IsHtmlElement) {
		return true
	}
	if elem.IsHtmlElement("input") {
		if attr, ok := elem.AttrWithoutNamespace("type"); ok && attr == "hidden" {
			return true
		}
	}
	return false
}

// https://html.spec.whatwg.org/multipage/semantics.html#contributes-a-script-blocking-style-sheet
func (elem htmlElementImpl) ContributesScriptBlockingStylesheet() bool {
	// STUB
	return false
}
