// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package elements

import "github.com/inseo-oh/yw/dom"

// NewHTMLHtmlElement constructs a new [HTMLElement] node for a [html] element.
//
// [html]: https://html.spec.whatwg.org/multipage/semantics.html#the-html-element
func NewHTMLHtmlElement(options dom.ElementCreationCommonOptions) HTMLElement {
	return NewHTMLElement(options)
}
