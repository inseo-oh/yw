package elements

import (
	"log"
	"yw/dom"
)

type HTMLStyleElement interface{ HTMLElement }
type htmlStyleElementImpl struct {
	HTMLElement
}

func NewHTMLStyleElement(options dom.ElementCreationCommonOptions) HTMLStyleElement {
	elem := htmlStyleElementImpl{
		HTMLElement: NewHTMLElement(options),
	}

	cbs := elem.Callbacks()

	// From 4.2.6. The style element(https://html.spec.whatwg.org/multipage/semantics.html#the-style-element)
	// The user agent must run the update a style block algorithm whenever any of the following conditions occur:
	//  - The element is popped off the stack of open elements of an HTML parser or XML parser.
	cbs.PoppedFromStackOfOpenElements = func() {
		elem.update_style_block()
	}
	//  - The element is not on the stack of open elements of an HTML parser or XML parser, and it becomes connected or disconnected.
	//  - The element's children changed steps run.
	cbs.RunChildrenChangedSteps = func() {
		elem.update_style_block()
	}
	return elem
}

// https://html.spec.whatwg.org/multipage/semantics.html#update-a-style-block
func (elem *htmlStyleElementImpl) update_style_block() {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S2.
	if sheet := css_associated_stylesheet(elem); sheet != nil {
		css_remove_stylesheet(sheet)
	}
	// S3.
	if !dom.IsConnected(elem) {
		return
	}
	// S4.
	// TODO: If element's type attribute is present and its value is neither the empty string nor an ASCII case-insensitive match for "text/css", then return.
	// S5.
	// TODO: If the Should element's inline behavior be blocked by Content Security Policy? algorithm returns "Blocked" when executed upon the style element, "style", and the style element's child text content, then return. [CSP]
	// S6.
	text, ok := elem.ChildTextNode()
	if !ok {
		text = ""
	}
	tokens, err := css_tokenize(text)
	if err != nil {
		log.Printf("<style>: failed to tokenize stylesheet: %v", err)
		return
	}
	stylesheet := css_parse_stylesheet(tokens, nil)
	stylesheet.tp = "text/css"
	stylesheet.owner_node = elem
	// TODO: Set stylesheet.media once we implement that
	if dom.IsInDocumentTree(elem) {
		if attr, ok := elem.AttrWithoutNamespace("title"); ok {
			stylesheet.title = attr
		}
	}
	stylesheet.alternate_flag = false
	stylesheet.origin_clean_flag = true
	stylesheet.location = nil
	stylesheet.parent_stylesheet = nil
	stylesheet.owner_rule = nil
	css_add_stylesheet(&stylesheet)
	log.Printf("<style>: stylesheet loaded")

	// S7.
	if elem.ContributesScriptBlockingStylesheet() {
		// TODO: append element to its node document's script-blocking style sheet set.
		panic("TODO[https://html.spec.whatwg.org/multipage/semantics.html#update-a-style-block]")
	}
	// S8.
	// If element's media attribute's value matches the environment and element is potentially render-blocking, then block rendering on element.

	// TODO: Specs has extra steps after critical subresources has been loaded, but they don't seem *that* important right now
	// (Mostly related to render blocking)

}
