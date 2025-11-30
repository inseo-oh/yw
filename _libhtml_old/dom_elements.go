package libhtml

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"strconv"
	"strings"

	cm "github.com/inseo-oh/yw/util"
)

// https://dom.spec.whatwg.org/#documentorshadowroot
type dom_DocumentOrShadowRoot interface {
	dom_Node
	get_custom_element_registry() *html_custom_element_registry
	get_css_stylesheets() []*css_stylesheet
	set_css_stylesheets(sheets []*css_stylesheet)
}

//------------------------------------------------------------------------------
// DOM Document
//------------------------------------------------------------------------------

// STUB types
type document_origin struct{}
type document_environment_settings struct{}
type document_policy_container struct{}

type dom_Document interface {
	dom_Node

	get_mode() dom_Document_mode
	set_mode(mode dom_Document_mode)
	is_parser_cannot_change_mode() bool
	is_iframe_srcdoc_document() bool
	effective_global_custom_element_registry() *html_custom_element_registry
	get_custom_element_registry() *html_custom_element_registry
	get_css_stylesheets() []*css_stylesheet
	set_css_stylesheets(sheets []*css_stylesheet)

	get_origin() document_origin
	set_origin(origin document_origin)
	get_releavant_settings() document_environment_settings
	set_environment_settings(settings document_environment_settings)
	get_policy_container() document_policy_container
	set_policy_container(policy_container document_policy_container)
	get_base_url() url.URL
	set_base_url(url url.URL)
}

type dom_Document_s struct {
	dom_Node
	mode                      dom_Document_mode
	iframe_srcdoc_document    bool
	parser_cannot_change_mode bool
	custom_elem_registry      html_custom_element_registry
	stylesheets               []*css_stylesheet
	base_url                  url.URL

	// Below are STUB
	origin               document_origin
	environment_settings document_environment_settings
	policy_container     document_policy_container
}
type dom_Document_mode uint8

const (
	dom_Document_mode_no_quirks = dom_Document_mode(iota)
	dom_Document_mode_quirks
	dom_Document_mode_limited_quirks
)

func dom_make_Document() dom_Document {
	doc := &dom_Document_s{}
	doc.dom_Node = dom_make_Node(doc)
	return doc
}
func (doc dom_Document_s) String() string {
	mode := "?"
	switch doc.mode {
	case dom_Document_mode_no_quirks:
		mode = "no-quirks"
	case dom_Document_mode_quirks:
		mode = "quirks"
	case dom_Document_mode_limited_quirks:
		mode = "limited-quirks"
	}
	return fmt.Sprintf("#document(mode=%s)", mode)
}

func (doc dom_Document_s) get_origin() document_origin {
	return doc.origin
}
func (doc *dom_Document_s) set_origin(origin document_origin) {
	doc.origin = origin
}
func (doc dom_Document_s) get_releavant_settings() document_environment_settings {
	return doc.environment_settings
}
func (doc *dom_Document_s) set_environment_settings(settings document_environment_settings) {
	doc.environment_settings = settings
}
func (doc dom_Document_s) get_policy_container() document_policy_container {
	return doc.policy_container
}
func (doc *dom_Document_s) set_policy_container(policy_container document_policy_container) {
	doc.policy_container = policy_container
}
func (doc dom_Document_s) get_base_url() url.URL {
	return doc.base_url
}
func (doc *dom_Document_s) set_base_url(url url.URL) {
	doc.base_url = url
}

func (doc dom_Document_s) get_mode() dom_Document_mode      { return doc.mode }
func (doc *dom_Document_s) set_mode(mode dom_Document_mode) { doc.mode = mode }
func (d dom_Document_s) is_parser_cannot_change_mode() bool { return d.parser_cannot_change_mode }
func (d dom_Document_s) is_iframe_srcdoc_document() bool    { return d.iframe_srcdoc_document }
func (doc dom_Document_s) get_custom_element_registry() *html_custom_element_registry {
	return &doc.custom_elem_registry
}
func (doc dom_Document_s) get_css_stylesheets() []*css_stylesheet        { return doc.stylesheets }
func (doc *dom_Document_s) set_css_stylesheets(sheets []*css_stylesheet) { doc.stylesheets = sheets }

// https://dom.spec.whatwg.org/#effective-global-custom-element-registry
func (doc dom_Document_s) effective_global_custom_element_registry() *html_custom_element_registry {
	if dom_is_global_custom_element_registry(&doc.custom_elem_registry) {
		return &doc.custom_elem_registry
	}
	return nil
}

//------------------------------------------------------------------------------
// DOM DocumentFragment
//------------------------------------------------------------------------------

type dom_DocumentFragment interface {
	dom_Node
	get_host() dom_Node
}

// STUB

//------------------------------------------------------------------------------
// DOM CharacterData
//------------------------------------------------------------------------------

type dom_CharacterData interface {
	dom_Node
	get_text() string
	append_text(s string)
}
type dom_CharacterData_s struct {
	dom_Node
	text string
}

func dom_make_CharacterData(doc dom_Document, text string) dom_CharacterData {
	return &dom_CharacterData_s{dom_make_Node(doc), text}
}
func (c dom_CharacterData_s) get_text() string {
	return c.text
}
func (c *dom_CharacterData_s) append_text(s string) {
	c.text += s
}

//------------------------------------------------------------------------------
// DOM Comment
//------------------------------------------------------------------------------

type dom_Comment interface {
	dom_CharacterData
}
type dom_Comment_s struct {
	dom_CharacterData
}

func dom_make_Comment(doc dom_Document, text string) dom_Comment {
	return &dom_Comment_s{dom_make_CharacterData(doc, text)}
}
func (cm dom_Comment_s) String() string {
	return fmt.Sprintf("<!-- %s -->", cm.get_text())
}

//------------------------------------------------------------------------------
// DOM Text
//------------------------------------------------------------------------------

type dom_Text interface {
	dom_CharacterData
}
type dom_Text_s struct{ dom_CharacterData }

func dom_make_Text(doc dom_Document, text string) dom_Text {
	return &dom_Text_s{dom_make_CharacterData(doc, text)}
}
func (txt dom_Text_s) i_am_dom_Text() {}
func (txt dom_Text_s) String() string {
	return strconv.Quote(txt.get_text())
}

//------------------------------------------------------------------------------
// DOM ShadowRoot
//------------------------------------------------------------------------------

type dom_ShadowRoot interface {
	dom_DocumentFragment
	get_css_stylesheets() []*css_stylesheet
	set_css_stylesheets(sheets []*css_stylesheet)
	get_custom_element_registry() *html_custom_element_registry
	set_custom_elem_registry(registry *html_custom_element_registry)
}

// STUB

// ------------------------------------------------------------------------------
// DOM Attr
// ------------------------------------------------------------------------------
type dom_Attr interface {
	dom_Node
	get_local_name() string
	get_value() string
	get_namespace() (namespace, bool)
	get_namespace_prefix() (string, bool)
	get_element() dom_Element
}
type dom_Attr_s struct {
	dom_Node
	local_name       string
	value            string
	namespace        *namespace // may be nil
	namespace_prefix *string    // may be nil
	element          dom_Element
}

// namespace, namespace_prefix may be nil
func dom_make_Attr(local_name string, value string, namespace *namespace, namespace_prefix *string, element dom_Element) dom_Attr {
	return &dom_Attr_s{
		dom_make_Node(element.get_node_document()),
		local_name, value, namespace, namespace_prefix, element,
	}
}

func (at dom_Attr_s) String() string {
	if at.namespace != nil {
		return fmt.Sprintf("#attr(%v:%s = %s)", *at.namespace, at.local_name, at.value)
	} else {
		return fmt.Sprintf("#attr(%s = %s)", at.local_name, at.value)
	}
}
func (at dom_Attr_s) get_local_name() string   { return at.local_name }
func (at dom_Attr_s) get_value() string        { return at.value }
func (at dom_Attr_s) get_element() dom_Element { return at.element }
func (at dom_Attr_s) get_namespace() (namespace, bool) {
	if at.namespace == nil {
		return "", false
	}
	return *at.namespace, true
}
func (at dom_Attr_s) get_namespace_prefix() (string, bool) {
	if at.namespace_prefix == nil {
		return "", false
	}
	return *at.namespace_prefix, true
}

//------------------------------------------------------------------------------
// DOM DocumentType
//------------------------------------------------------------------------------

type dom_DocumentType interface {
	dom_Node
	get_name() string
	get_public_id() string
	get_system_id() string
}
type dom_DocumentType_s struct {
	dom_Node
	name      string
	public_id string
	system_id string
}

func dom_make_DocumentType(doc dom_Document, name, public_id, system_id string) dom_DocumentType {
	return dom_DocumentType_s{
		dom_make_Node(doc),
		name, public_id, system_id,
	}
}
func (dt dom_DocumentType_s) String() string {
	sb := strings.Builder{}
	sb.WriteString("<!DOCTYPE")
	if dt.name != "" {
		sb.WriteString(fmt.Sprintf(" %s", dt.name))
	}
	if dt.public_id != "" {
		sb.WriteString(fmt.Sprintf(" PUBLIC %s", dt.public_id))
	}
	if dt.system_id != "" {
		sb.WriteString(fmt.Sprintf(" SYSTEM %s", dt.system_id))
	}
	sb.WriteString(">")
	return sb.String()
}
func (dt dom_DocumentType_s) get_name() string      { return dt.name }
func (dt dom_DocumentType_s) get_public_id() string { return dt.public_id }
func (dt dom_DocumentType_s) get_system_id() string { return dt.system_id }

//------------------------------------------------------------------------------
// DOM Element
//------------------------------------------------------------------------------

type dom_Element interface {
	dom_Node
	is_block_level() bool
	get_local_name() string
	get_namespace() (namespace, bool)
	get_prefix() (string, bool)
	get_is() (string, bool)
	get_attrs() []dom_Attr
	get_tag_token() html_tag_token
	get_custom_element_registry() *html_custom_element_registry
	set_custom_elem_registry(reg *html_custom_element_registry)
	get_shadow_root() dom_ShadowRoot
	set_shadow_root(sr dom_ShadowRoot)
	is_shadow_host() bool
	is_custom() bool
	append_attr(attr_s dom_Attr_s)
	is_element(name dom_name_pair) bool
	is_html_element(local_name string) bool
	is_mathml_element(local_name string) bool
	is_svg_element(local_name string) bool
	get_attr_with_namespace(name dom_name_pair) (string, bool)
	get_attribute_without_namespace(name string) (string, bool)
	get_computed_style_set() *css_computed_style_set
	get_intrinsic_size() (width float64, height float64)
	is_inside(dom_name_pair dom_name_pair) bool

	// Some extensions for HTML parser

	is_html_special_element() bool
	is_html_formatting_element() bool
	is_html_ordinary_element() bool
	is_mathml_text_integration_point() bool
	is_html_integration_point() bool
}

type dom_Element_s struct {
	dom_Node

	namespace            *namespace // May be nil
	prefix               *string    // May be nil
	local_name           string
	is                   *string                       // May be nil
	custom_elem_registry *html_custom_element_registry // May be nil
	shadow_root          dom_ShadowRoot                // may be nil
	attrs                []dom_Attr
	tag_token            html_tag_token
	custom_element_state dom_custom_element_state
	callbacks            dom_element_callbacks
	computed_style_set   css_computed_style_set
	is_block_level_elem  bool
}
type dom_custom_element_state uint8

const (
	// https://dom.spec.whatwg.org/#concept-element-custom-element-state
	dom_custom_element_state_undefined = dom_custom_element_state(iota)
	dom_custom_element_state_failed
	dom_custom_element_state_uncustomized
	dom_custom_element_state_precustomized
	dom_custom_element_state_custom
)

type dom_element_callbacks struct {
	get_intrinsic_size func() (width float64, height float64)
}

func dom_make_Element(options dom_element_creation_common_options, callbacks dom_element_callbacks, is_block_level_elem bool) dom_Element {
	return &dom_Element_s{
		dom_Node:             dom_make_Node(options.node_document),
		namespace:            options.namespace,
		prefix:               options.prefix,
		local_name:           options.local_name,
		attrs:                []dom_Attr{},
		tag_token:            options.tag_token,
		custom_elem_registry: options.custom_elem_registry,
		custom_element_state: options.custom_element_state,
		is:                   options.is,
		callbacks:            callbacks,
		is_block_level_elem:  is_block_level_elem,
	}
}
func (n dom_Element_s) is_block_level() bool   { return n.is_block_level_elem }
func (n dom_Element_s) get_local_name() string { return n.local_name }
func (n dom_Element_s) get_namespace() (namespace, bool) {
	if n.namespace == nil {
		return "", false
	}
	return *n.namespace, true
}
func (n dom_Element_s) get_prefix() (string, bool) {
	if n.prefix == nil {
		return "", false
	}
	return *n.prefix, true
}
func (n dom_Element_s) get_is() (string, bool) {
	if n.is == nil {
		return "", false
	}
	return *n.is, true
}
func (n dom_Element_s) get_custom_element_registry() *html_custom_element_registry {
	return n.custom_elem_registry
}
func (n *dom_Element_s) set_custom_elem_registry(reg *html_custom_element_registry) {
	n.custom_elem_registry = reg
}
func (n dom_Element_s) get_shadow_root() dom_ShadowRoot {
	return n.shadow_root
}
func (n *dom_Element_s) set_shadow_root(sr dom_ShadowRoot) {
	n.shadow_root = sr
}
func (n dom_Element_s) is_shadow_host() bool {
	return !cm.IsNil(n.shadow_root)
}
func (n dom_Element_s) is_defined() bool {
	// https://dom.spec.whatwg.org/#concept-element-defined
	return n.custom_element_state == dom_custom_element_state_uncustomized ||
		n.custom_element_state == dom_custom_element_state_custom
}
func (n dom_Element_s) is_custom() bool {
	// https://dom.spec.whatwg.org/#concept-element-custom
	return n.custom_element_state == dom_custom_element_state_custom
}
func (n dom_Element_s) get_attrs() []dom_Attr {
	return n.attrs
}
func (n dom_Element_s) get_tag_token() html_tag_token {
	return n.tag_token
}
func (n *dom_Element_s) get_computed_style_set() *css_computed_style_set {
	return &n.computed_style_set
}
func (n dom_Element_s) get_intrinsic_size() (width float64, height float64) {
	if n.callbacks.get_intrinsic_size == nil {
		log.Printf("%v does not implement get_intrinsic_size(). Assuming size 0, 0", n)
		return 0, 0
	}
	return n.callbacks.get_intrinsic_size()
}

func (n dom_Element_s) String() string {
	sb := strings.Builder{}
	sb.WriteString("<")
	if n.namespace != nil {
		sb.WriteString(fmt.Sprintf("%v:", n.namespace))
	}
	sb.WriteString(n.local_name)
	for _, attr := range n.attrs {
		sb.WriteString(" ")
		val := strconv.Quote(attr.get_value())
		ns, has_ns := attr.get_namespace()
		if has_ns {
			sb.WriteString(fmt.Sprintf("%v:%s=%s", ns, attr.get_local_name(), val))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s", attr.get_local_name(), val))
		}
	}
	sb.WriteString(">")
	return sb.String()
}

type dom_name_pair struct {
	namespace  namespace
	local_name string
}

func (n dom_Element_s) is_inside(name dom_name_pair) bool {
	current := n.get_parent()
	for !cm.IsNil(current) {
		current_elem, ok := current.(dom_Element)
		if !ok {
			break
		} else if current_elem.is_element(name) {
			return true
		}
		current = current_elem.(dom_Node).get_parent()
	}
	return false
}

// Note that attr.element will be overwritten by this function, so don't set it.
func (n *dom_Element_s) append_attr(attr_s dom_Attr_s) {
	attr := dom_make_Attr(attr_s.local_name, attr_s.value, attr_s.namespace, attr_s.namespace_prefix, n)
	n.attrs = append(n.attrs, attr)
}

func (n dom_Element_s) is_element(name dom_name_pair) bool {
	return n.namespace != nil && *n.namespace == name.namespace && n.local_name == name.local_name
}
func (n dom_Element_s) is_html_element(local_name string) bool {
	return n.is_element(dom_name_pair{html_namespace, local_name})
}
func (n dom_Element_s) is_mathml_element(local_name string) bool {
	return n.is_element(dom_name_pair{mathml_namespace, local_name})
}
func (n dom_Element_s) is_svg_element(local_name string) bool {
	return n.is_element(dom_name_pair{svg_namespace, local_name})
}

// This will only match attributes with namespace specified.
func (n dom_Element_s) get_attr_with_namespace(name dom_name_pair) (string, bool) {
	for _, attr := range n.attrs {
		if ns, has_ns := attr.get_namespace(); has_ns && ns == name.namespace && attr.get_local_name() == name.local_name {
			return attr.get_value(), true
		}
	}
	return "", false
}

// Returns nil if not found. This will only match attributes without namespace specified.
func (n dom_Element_s) get_attribute_without_namespace(name string) (string, bool) {
	for _, attr := range n.attrs {
		if _, has_ns := attr.get_namespace(); !has_ns && attr.get_local_name() == name {
			return attr.get_value(), true
		}
	}
	return "", false
}

//------------------------------------------------------------------------------
// Some extensions for HTML parser
//------------------------------------------------------------------------------

// https://html.spec.whatwg.org/multipage/parsing.html#special
func (n dom_Element_s) is_html_special_element() bool {
	special_html_elems := []string{
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
	special_mathml_elems := []string{
		"mi", "mo", "mn", "ms", "mtext", "annotation-xml",
	}
	special_svg_elems := []string{
		"foreignObject", "desc", "title",
	}
	if slices.ContainsFunc(special_html_elems, n.is_html_element) {
		return true
	}
	if slices.ContainsFunc(special_mathml_elems, n.is_mathml_element) {
		return true
	}
	return slices.ContainsFunc(special_svg_elems, n.is_svg_element)
}

// https://html.spec.whatwg.org/multipage/parsing.html#formatting
func (n dom_Element_s) is_html_formatting_element() bool {
	special_html_elems := []string{
		"a", "b", "big", "code", "em", "font", "i", "nobr", "s", "small",
		"strike", "strong", "tt", "u",
	}
	return slices.ContainsFunc(special_html_elems, n.is_html_element)
}

// https://html.spec.whatwg.org/multipage/parsing.html#ordinary
func (n dom_Element_s) is_html_ordinary_element() bool {
	return !n.is_html_special_element() && !n.is_html_formatting_element()
}

// https://html.spec.whatwg.org/multipage/parsing.html#mathml-text-integration-point
func (n dom_Element_s) is_mathml_text_integration_point() bool {
	return slices.ContainsFunc([]string{"mi", "mo", "mn", "ms", "mtext"}, n.is_mathml_element)
}

// https://html.spec.whatwg.org/multipage/parsing.html#html-integration-point
func (n dom_Element_s) is_html_integration_point() bool {
	if n.is_mathml_element("annotation-xml") {
		if attr := n.tag_token.get_attr("encoding"); attr != nil &&
			(cm.ToAsciiLowercase(*attr) == "text/html" || cm.ToAsciiLowercase(*attr) == "application/xhtml+xml") {
			return true
		}
	}
	if n.is_svg_element("foreignObject") ||
		n.is_svg_element("desc") ||
		n.is_svg_element("title") {
		return true
	}
	return false
}
