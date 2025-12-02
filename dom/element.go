package dom

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/namespaces"
	"github.com/inseo-oh/yw/util"
)

// Element represents a [DOM element].
//
// TODO(ois): Should we move some of Element's HTML related extensions to the html package instead?
//
// [DOM element]: https://dom.spec.whatwg.org/#concept-element
type Element interface {
	Node

	// Namespace returns [namespace] of the element. ok is set to false if it's absent.
	//
	// [namespace]: https://dom.spec.whatwg.org/#concept-element-namespace
	Namespace() (ns namespaces.Namespace, ok bool)

	// Prefix returns [namespace prefix] of the element. ok is set to false if it's absent.
	//
	// TODO(ois): Should we call this NamespacePrefix instead?
	// [namespace prefix]: https://dom.spec.whatwg.org/#concept-element-namespace-prefix.
	Prefix() (pr string, ok bool)

	// LocalName returns [local name] of the element.
	//
	// [local name]: https://dom.spec.whatwg.org/#concept-element-local-name
	LocalName() string

	// CustomElementRegistry returns [custom element registry] of the element.
	//
	// [custom element registry]: https://dom.spec.whatwg.org/#element-custom-element-registry
	CustomElementRegistry() *CustomElementRegistry

	// CustomElementRegistry sets [custom element registry] of the element to reg.
	//
	// [custom element registry]: https://dom.spec.whatwg.org/#element-custom-element-registry
	SetCustomElementRegistry(reg *CustomElementRegistry)

	// CustomElementRegistry returns [is value] of the element.
	//
	// [is value]: https://dom.spec.whatwg.org/#concept-element-is-value
	Is() (string, bool)

	// IsCustom reports whether element is [custom].
	//
	// [custom]: https://dom.spec.whatwg.org/#concept-element-custom
	IsCustom() bool

	// ShadowRoot returns [shadow root] of the element.
	//
	// [shadow root]: https://dom.spec.whatwg.org/#concept-element-shadow-root
	ShadowRoot() ShadowRoot

	// SetShadowRoot sets [shadow root] of the element to sr.
	//
	// [shadow root]: https://dom.spec.whatwg.org/#concept-element-shadow-root
	SetShadowRoot(sr ShadowRoot)

	// IsShadowHost reports whether element is [shadow host].
	//
	// [shadow host]: https://dom.spec.whatwg.org/#element-shadow-host
	IsShadowHost() bool

	// Attrs returns [attributes] of the element.
	//
	// [attributes]: https://dom.spec.whatwg.org/#concept-element-attribute
	Attrs() []Attr

	// Attrs appends new attribute to [attributes] of the element.
	//
	// [attributes]: https://dom.spec.whatwg.org/#concept-element-attribute
	AppendAttr(attrData AttrData)

	// AttrWithNamespace searches attribute from element's [attributes], that
	// matches namePair's Namespace and Name, and returns its value.
	// Attributes without namespace are ignored. ok is set to false if there's
	// no such attribute.
	//
	// [attributes]: https://dom.spec.whatwg.org/#concept-element-attribute
	AttrWithNamespace(namePair NamePair) (value string, ok bool)

	// AttrWithoutNamespace searches attribute from element's [attributes], that
	// matches the name, and returns its value.
	// Attributes with namespace are ignored. ok is set to false if there's no
	// such attribute.
	//
	// [attributes]: https://dom.spec.whatwg.org/#concept-element-attribute
	AttrWithoutNamespace(name string) (value string, ok bool)

	// IsElement reports whether the element's local name is equal to namePair's
	// LocalName and is in name's Namespace.
	IsElement(namePair NamePair) bool

	// IsHtmlElement reports whether the element's local name is equal to given
	// localName and is in HTML namespace.
	IsHtmlElement(localName string) bool

	// IsMathmlElement reports whether the element's local name is equal to given
	// localName and is in MathML namespace.
	IsMathmlElement(localName string) bool

	// IsSvgElement reports whether the element's local name is equal to given
	// localName and is in SVG namespace.
	IsSvgElement(localName string) bool

	// IntrinsicSize returns intrinsic size of the element.
	IntrinsicSize() (width float64, height float64)

	// IsInside reports whether the element is inside another element that matches
	// given namePair's Namespace and Name.
	IsInside(namePair NamePair) bool

	// TagToken returns associated HTML tag token for the element.
	TagToken() TagToken

	//==========================================================================
	// Some extensions for HTML parser
	//==========================================================================

	// Reports whether the element is in HTML parser's [special category].
	//
	// [special category]: https://html.spec.whatwg.org/multipage/parsing.html#special
	IsHtmlSpecialElement() bool

	// Reports whether the element is in HTML parser's [formatting category].
	//
	// [formatting category]: https://html.spec.whatwg.org/multipage/parsing.html#formatting
	IsHtmlFormattingElement() bool

	// Reports whether the element is in HTML parser's [ordinary category].
	//
	// [ordinary category]: https://html.spec.whatwg.org/multipage/parsing.html#ordinary
	IsHtmlOrdinaryElement() bool

	// Reports whether the element is [MathML text integration point].
	//
	// [MathML text integration point]: https://html.spec.whatwg.org/multipage/parsing.html#mathml-text-integration-point
	IsMathmlTextIntegrationPoint() bool

	// Reports whether the element is [HTML integration point].
	//
	// [HTML integration point]: https://html.spec.whatwg.org/multipage/parsing.html#html-integration-point
	IsHtmlIntegrationPoint() bool
}

type elementImpl struct {
	Node

	namespace             *namespaces.Namespace // May be nil
	prefix                *string               // May be nil
	localName             string
	is                    *string                // May be nil
	customElementReigstry *CustomElementRegistry // May be nil
	shadowRoot            ShadowRoot             // may be nil
	attrs                 []Attr
	tagToken              TagToken
	customElementState    CustomElementState
}
type CustomElementState uint8

// https://dom.spec.whatwg.org/#concept-element-custom-element-state
const (
	CustomElementUndefined CustomElementState = iota
	CustomElementFailed
	CustomElementUncustomized
	CustomElementPrecustomized
	CustomElementCustom
)

func NewElement(options ElementCreationCommonOptions) Element {
	return &elementImpl{
		Node:                  NewNode(options.NodeDocument),
		namespace:             options.Namespace,
		prefix:                options.Prefix,
		localName:             options.LocalName,
		attrs:                 []Attr{},
		tagToken:              options.TagToken,
		customElementReigstry: options.CustomElementRegistry,
		customElementState:    options.CustomElementState,
		is:                    options.Is,
	}
}
func (n elementImpl) LocalName() string { return n.localName }
func (n elementImpl) Namespace() (namespaces.Namespace, bool) {
	if n.namespace == nil {
		return "", false
	}
	return *n.namespace, true
}
func (n elementImpl) Prefix() (string, bool) {
	if n.prefix == nil {
		return "", false
	}
	return *n.prefix, true
}
func (n elementImpl) Is() (string, bool) {
	if n.is == nil {
		return "", false
	}
	return *n.is, true
}
func (n elementImpl) CustomElementRegistry() *CustomElementRegistry {
	return n.customElementReigstry
}
func (n *elementImpl) SetCustomElementRegistry(reg *CustomElementRegistry) {
	n.customElementReigstry = reg
}
func (n elementImpl) ShadowRoot() ShadowRoot {
	return n.shadowRoot
}
func (n *elementImpl) SetShadowRoot(sr ShadowRoot) {
	n.shadowRoot = sr
}
func (n elementImpl) IsShadowHost() bool {
	return !util.IsNil(n.shadowRoot)
}
func (n elementImpl) IsDefined() bool {
	// https://dom.spec.whatwg.org/#concept-element-defined
	return n.customElementState == CustomElementUncustomized ||
		n.customElementState == CustomElementCustom
}
func (n elementImpl) IsCustom() bool {
	// https://dom.spec.whatwg.org/#concept-element-custom
	return n.customElementState == CustomElementCustom
}
func (n elementImpl) Attrs() []Attr {
	return n.attrs
}
func (n elementImpl) TagToken() TagToken {
	return n.tagToken
}
func (n elementImpl) IntrinsicSize() (width float64, height float64) {
	if n.Callbacks().IntrinsicSize == nil {
		log.Printf("%v does not implement IntrinsicSize(). Assuming size 0, 0", n)
		return 0, 0
	}
	return n.Callbacks().IntrinsicSize()
}

func (n elementImpl) String() string {
	sb := strings.Builder{}
	sb.WriteString("<")
	if n.namespace != nil {
		sb.WriteString(fmt.Sprintf("%v:", n.namespace))
	}
	sb.WriteString(n.localName)
	for _, attr := range n.attrs {
		sb.WriteString(" ")
		val := strconv.Quote(attr.Value())
		if ns, ok := attr.Namespace(); ok {
			sb.WriteString(fmt.Sprintf("%v:%s=%s", ns, attr.LocalName(), val))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s", attr.LocalName(), val))
		}
	}
	sb.WriteString(">")
	return sb.String()
}

type NamePair struct {
	Namespace namespaces.Namespace
	LocalName string
}

func (n elementImpl) IsInside(namePair NamePair) bool {
	current := n.Parent()
	for !util.IsNil(current) {
		currElem, ok := current.(Element)
		if !ok {
			break
		} else if currElem.IsElement(namePair) {
			return true
		}
		current = currElem.(Node).Parent()
	}
	return false
}

// Note that attr.element will be overwritten by this function, so don't set it.
func (n *elementImpl) AppendAttr(attrData AttrData) {
	attr := NewAttr(attrData.LocalName, attrData.Value, attrData.Namespace, attrData.NamespacePrefix, n)
	n.attrs = append(n.attrs, attr)
}

func (n elementImpl) IsElement(namePair NamePair) bool {
	return n.namespace != nil && *n.namespace == namePair.Namespace && n.localName == namePair.LocalName
}
func (n elementImpl) IsHtmlElement(localName string) bool {
	return n.IsElement(NamePair{namespaces.Html, localName})
}
func (n elementImpl) IsMathmlElement(localName string) bool {
	return n.IsElement(NamePair{namespaces.Mathml, localName})
}
func (n elementImpl) IsSvgElement(localName string) bool {
	return n.IsElement(NamePair{namespaces.Svg, localName})
}

func (n elementImpl) AttrWithNamespace(namePair NamePair) (string, bool) {
	for _, attr := range n.attrs {
		if ns, hasNs := attr.Namespace(); hasNs && ns == namePair.Namespace && attr.LocalName() == namePair.LocalName {
			return attr.Value(), true
		}
	}
	return "", false
}

// Returns nil if not found. This will only match attributes without namespace specified.
func (n elementImpl) AttrWithoutNamespace(name string) (string, bool) {
	for _, attr := range n.attrs {
		if _, hasNs := attr.Namespace(); !hasNs && attr.LocalName() == name {
			return attr.Value(), true
		}
	}
	return "", false
}

//------------------------------------------------------------------------------
// Some extensions for HTML parser
//------------------------------------------------------------------------------

// https://html.spec.whatwg.org/multipage/parsing.html#special
func (n elementImpl) IsHtmlSpecialElement() bool {
	specialHtmlElems := []string{
		"address", "applet", "area", "article", "aside", "base", "basefont",
		"bgsound", "blockquote", "body", "br", "button", "caption", "center",
		"col", "colgroup", "dd", "details", "dir", "div", "dl", "dt", "embed",
		"fieldset", "figcaption", "figure", "footer", "form", "frame",
		"frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header",
		"hgroup", "hr", "html", "iframe", "img", "input", "keygen", "li",
		"link", "listing", "main", "marquee", "menu", "meta", "nav", "noembed",
		"noframes", "noscript", "object", "ol", "p", "param", "plaintext",
		"pre", "script", "search", "section", "select", "source", "style",
		"summary", "table", "tbody", "td", "template", "textarea", "tfoot",
		"th", "thead", "title", "tr", "track", "ul", "wbr", "xmp",
	}
	specialMathmlElems := []string{
		"mi", "mo", "mn", "ms", "mtext", "annotation-xml",
	}
	specialSvgElems := []string{
		"foreignObject", "desc", "title",
	}
	if slices.ContainsFunc(specialHtmlElems, n.IsHtmlElement) {
		return true
	}
	if slices.ContainsFunc(specialMathmlElems, n.IsMathmlElement) {
		return true
	}
	return slices.ContainsFunc(specialSvgElems, n.IsSvgElement)
}

// https://html.spec.whatwg.org/multipage/parsing.html#formatting
func (n elementImpl) IsHtmlFormattingElement() bool {
	specialHtmlElems := []string{
		"a", "b", "big", "code", "em", "font", "i", "nobr", "s", "small",
		"strike", "strong", "tt", "u",
	}
	return slices.ContainsFunc(specialHtmlElems, n.IsHtmlElement)
}

// https://html.spec.whatwg.org/multipage/parsing.html#ordinary
func (n elementImpl) IsHtmlOrdinaryElement() bool {
	return !n.IsHtmlSpecialElement() && !n.IsHtmlFormattingElement()
}

// https://html.spec.whatwg.org/multipage/parsing.html#mathml-text-integration-point
func (n elementImpl) IsMathmlTextIntegrationPoint() bool {
	return slices.ContainsFunc([]string{"mi", "mo", "mn", "ms", "mtext"}, n.IsMathmlElement)
}

// https://html.spec.whatwg.org/multipage/parsing.html#html-integration-point
func (n elementImpl) IsHtmlIntegrationPoint() bool {
	if n.IsMathmlElement("annotation-xml") {
		if attr, ok := n.tagToken.Attr("encoding"); ok &&
			(util.ToAsciiLowercase(attr) == "text/html" || util.ToAsciiLowercase(attr) == "application/xhtml+xml") {
			return true
		}
	}
	if n.IsSvgElement("foreignObject") ||
		n.IsSvgElement("desc") ||
		n.IsSvgElement("title") {
		return true
	}
	return false
}
