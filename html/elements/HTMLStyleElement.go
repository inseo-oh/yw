package elements

import (
	"log"

	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/csssyntax"
	"github.com/inseo-oh/yw/dom"
)

// HTMLStyleElement represents a [style] element.
//
// [style]: https://html.spec.whatwg.org/multipage/semantics.html#the-style-element
type HTMLStyleElement interface{ HTMLElement }
type htmlStyleElementImpl struct {
	HTMLElement
}

// NewHTMLStyleElement constructs a new [HTMLStyleElement] node.
//
// [html]: https://html.spec.whatwg.org/multipage/semantics.html#the-html-element
func NewHTMLStyleElement(options dom.ElementCreationCommonOptions) HTMLStyleElement {
	elem := htmlStyleElementImpl{
		HTMLElement: NewHTMLElement(options),
	}

	cbs := elem.Callbacks()

	// From 4.2.6. The style element(https://html.spec.whatwg.org/multipage/semantics.html#the-style-element)
	// The user agent must run the update a style block algorithm whenever any of the following conditions occur:
	//  - The element is popped off the stack of open elements of an HTML parser or XML parser.
	cbs.PoppedFromStackOfOpenElements = func() {
		elem.updateStyleBlock()
	}
	//  - The element is not on the stack of open elements of an HTML parser or XML parser, and it becomes connected or disconnected.
	//  - The element's children changed steps run.
	cbs.RunChildrenChangedSteps = func() {
		elem.updateStyleBlock()
	}
	return elem
}

// https://html.spec.whatwg.org/multipage/semantics.html#update-a-style-block
func (elem *htmlStyleElementImpl) updateStyleBlock() {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S2.
	if sheet := cssom.AssociatedStylesheet(elem); sheet != nil {
		cssom.RemoveStylesheet(sheet)
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
	stylesheet, err := csssyntax.ParseStylesheet([]byte(text), nil)
	if err != nil {
		log.Printf("<style>: failed to tokenize stylesheet: %v", err)
		return
	}
	stylesheet.Type = "text/css"
	stylesheet.OwnerNode = elem
	// TODO: Set stylesheet.media once we implement that
	if dom.IsInDocumentTree(elem) {
		if attr, ok := elem.AttrWithoutNamespace("title"); ok {
			stylesheet.Title = attr
		}
	}
	stylesheet.AlternateFlag = false
	stylesheet.OriginCleanFlag = true
	stylesheet.Location = nil
	stylesheet.ParentStylesheet = nil
	stylesheet.OwnerRule = nil
	cssom.AddStylesheet(&stylesheet)
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
