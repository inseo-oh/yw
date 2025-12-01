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

type Element interface {
	Node
	LocalName() string
	Namespace() (namespaces.Namespace, bool)
	Prefix() (string, bool)
	Is() (string, bool)
	Attrs() []Attr
	TagToken() TagToken
	CustomElementRegistry() *CustomElementRegistry
	SetCustomElementRegistry(reg *CustomElementRegistry)
	ShadowRoot() ShadowRoot
	SetShadowRoot(sr ShadowRoot)
	IsShadowHost() bool
	IsCustom() bool
	AppendAttr(attrData AttrData)
	IsElement(name NamePair) bool
	IsHtmlElement(localName string) bool
	IsMathmlElement(localName string) bool
	IsSvgElement(localName string) bool
	AttrWithNamespace(name NamePair) (string, bool)
	AttrWithoutNamespace(name string) (string, bool)
	IntrinsicSize() (width float64, height float64)
	IsInside(namePair NamePair) bool

	// Some extensions for HTML parser
	// TODO: Should we move this to HTML parser instead?

	IsHtmlSpecialElement() bool
	IsHtmlFormattingElement() bool
	IsHtmlOrdinaryElement() bool
	IsMathmlTextIntegrationPoint() bool
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

func (n elementImpl) IsInside(name NamePair) bool {
	current := n.Parent()
	for !util.IsNil(current) {
		currElem, ok := current.(Element)
		if !ok {
			break
		} else if currElem.IsElement(name) {
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

func (n elementImpl) IsElement(name NamePair) bool {
	return n.namespace != nil && *n.namespace == name.Namespace && n.localName == name.LocalName
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

// This will only match attributes with namespace specified.
func (n elementImpl) AttrWithNamespace(name NamePair) (string, bool) {
	for _, attr := range n.attrs {
		if ns, hasNs := attr.Namespace(); hasNs && ns == name.Namespace && attr.LocalName() == name.LocalName {
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
