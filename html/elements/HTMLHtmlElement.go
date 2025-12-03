package elements

import "github.com/inseo-oh/yw/dom"

// NewHTMLHtmlElement constructs a new [HTMLElement] node for a [html] element.
//
// [html]: https://html.spec.whatwg.org/multipage/semantics.html#the-html-element
func NewHTMLHtmlElement(options dom.ElementCreationCommonOptions) HTMLElement {
	return NewHTMLElement(options)
}
