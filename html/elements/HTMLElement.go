// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package elements

import (
	"slices"

	"github.com/inseo-oh/yw/dom"
)

// HTMLElement represents a [HTML element].
//
// [HTML element]: https://html.spec.whatwg.org/multipage/dom.html#htmlelement
type HTMLElement interface {
	dom.Element

	// IsFormAssociatedCustomElement reports whether the element is [form-associated custom element].
	//
	// [form-associated custom element]: https://html.spec.whatwg.org/multipage/custom-elements.html#form-associated-custom-element
	IsFormAssociatedCustomElement() bool

	// IsFormAssociatedElement reports whether the element is [form-associated element].
	//
	// [form-associated element]: https://html.spec.whatwg.org/multipage/forms.html#form-associated-element
	IsFormAssociatedElement() bool

	// IsFormListedElement reports whether the element is [listed element].
	//
	// [listed element]: https://html.spec.whatwg.org/multipage/forms.html#category-listed
	IsFormListedElement() bool

	// IsFormSubmittableElement reports whether the element is [submittable element].
	//
	// [submittable element]: https://html.spec.whatwg.org/multipage/forms.html#category-submit
	IsFormSubmittableElement() bool

	// IsFormResettableElement reports whether the element is [resettable element].
	//
	// [resettable element]: https://html.spec.whatwg.org/multipage/forms.html#category-submit
	IsFormResettableElement() bool

	// IsFormAutocapitalizeAndAutocorrectInheritingElement reports whether the element is [autocapitalize-and-autocorrect-inheriting element].
	//
	// [autocapitalize-and-autocorrect-inheriting element]: https://html.spec.whatwg.org/multipage/forms.html#category-autocapitalize
	IsFormAutocapitalizeAndAutocorrectInheritingElement() bool

	// IsFormLabelableElement reports whether the element is [labelable element].
	//
	// [labelable element]: https://html.spec.whatwg.org/multipage/forms.html#category-label
	IsFormLabelableElement() bool

	// ContributesScriptBlockingStylesheet reports whether the element [contributes a script-blocking style sheet].
	//
	// [contributes a script-blocking style sheet]: https://html.spec.whatwg.org/multipage/semantics.html#contributes-a-script-blocking-style-sheet
	ContributesScriptBlockingStylesheet() bool
}
type htmlElementImpl struct{ dom.Element }

// NewHTMLElement constructs a new [HTMLElement] node.
func NewHTMLElement(options dom.ElementCreationCommonOptions) HTMLElement {
	return htmlElementImpl{dom.NewElement(options)}
}

func (elem htmlElementImpl) IsFormAssociatedCustomElement() bool {
	// STUB
	return false
}

func (elem htmlElementImpl) IsFormAssociatedElement() bool {
	htmlElems := []string{
		"button", "fieldset", "input", "object", "output", "select", "textarea", "img",
	}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(htmlElems, elem.IsHtmlElement)
}

func (elem htmlElementImpl) IsFormListedElement() bool {
	htmlElems := []string{
		"button", "fieldset", "input", "object", "output", "select", "textarea",
	}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(htmlElems, elem.IsHtmlElement)
}

func (elem htmlElementImpl) IsFormSubmittableElement() bool {
	htmlElems := []string{"button", "input", "select", "textarea"}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(htmlElems, elem.IsHtmlElement)
}

func (elem htmlElementImpl) IsFormResettableElement() bool {
	htmlElems := []string{"input", "output", "select", "textarea"}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(htmlElems, elem.IsHtmlElement)
}

func (elem htmlElementImpl) IsFormAutocapitalizeAndAutocorrectInheritingElement() bool {
	htmlElems := []string{
		"button", "fieldset", "input", "output", "select", "textarea",
	}
	return elem.IsFormAssociatedCustomElement() ||
		slices.ContainsFunc(htmlElems, elem.IsHtmlElement)
}

func (elem htmlElementImpl) IsFormLabelableElement() bool {
	htmlElems := []string{
		"button", "meter", "output", "progress", "select", "textarea",
	}
	if elem.IsFormAssociatedCustomElement() {
		return true
	}
	if slices.ContainsFunc(htmlElems, elem.IsHtmlElement) {
		return true
	}
	if elem.IsHtmlElement("input") {
		if attr, ok := elem.AttrWithoutNamespace("type"); ok && attr == "hidden" {
			return true
		}
	}
	return false
}

func (elem htmlElementImpl) ContributesScriptBlockingStylesheet() bool {
	// STUB
	return false
}
