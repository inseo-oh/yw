package parser

import (
	"log"
	"runtime"
	"slices"
	"strings"
	"yw/dom"
	"yw/html/elements"
	cm "yw/libcommon"
	"yw/namespaces"
)

type html_parser struct {
	tokenizer html_tokenizer

	document dom.Document

	head_element_pointer dom.Element
	form_element_pointer dom.Element

	run_parser                    bool
	is_frameset_not_ok            bool
	is_fragment_parsing           bool
	enable_scripting              bool
	enable_foster_parenting       bool
	has_active_speculative_parser bool // We don't have speculative parsing support, so this is mostly just a placeholder, just in case decide to we support it later.

	insertion_mode          html_parser_insertion_mode
	original_insertion_mode html_parser_insertion_mode

	stack_of_open_elems               []dom.Element
	list_of_active_formatting_elems   []html_active_formatting_elem
	stack_of_template_insertion_modes []html_parser_insertion_mode

	on_next_token func(token html_token) html_parser_control

	pending_table_char_tokens []html_char_token // https://html.spec.whatwg.org/multipage/parsing.html#concept-pending-table-char-tokens
}

func NewParser(str string) html_parser {
	return html_parser{
		tokenizer: html_make_tokenizer(str),
	}
}

func (p *html_parser) Run() dom.Document {
	if p.document == nil {
		p.document = dom.NewDocument()
	}
	p.tokenizer.on_token_emitted = func(tk html_token) {
		if p.on_next_token != nil {
			switch p.on_next_token(tk) {
			case html_parser_control_ignore_token:
				return
			case html_parser_control_continue:
			default:
				panic("unknown result from on_next_token()")
			}
		}

		is_start_tag_token := func() bool {
			if _, ok := (tk).(tagToken); ok {
				return true
			}
			return false
		}
		is_start_tag_token_with := func(name string) bool {
			if tk, ok := (tk).(tagToken); ok {
				return tk.is_start_tag() && tk.tag_name == name
			}
			return false
		}
		is_char_token := func() bool {
			if _, ok := (tk).(html_char_token); ok {
				return true
			}
			return false
		}
		is_eof_token := func() bool {
			if _, ok := (tk).(html_eof_token); ok {
				return true
			}
			return false
		}

		// https://html.spec.whatwg.org/multipage/parsing.html#tree-construction-dispatcher
		if (len(p.stack_of_open_elems) == 0) ||
			(func() bool {
				n := p.get_adjusted_current_node()
				if ns, ok := n.Namespace(); ok && ns == namespaces.Html {
					return true
				}
				return false
			}()) ||
			(p.get_adjusted_current_node().IsMathmlTextIntegrationPoint() && !is_start_tag_token_with("mglyph") && !is_start_tag_token_with("malignmark")) ||
			(p.get_adjusted_current_node().IsMathmlTextIntegrationPoint() && is_char_token()) ||
			(p.get_adjusted_current_node().IsMathmlElement("annotation-xml") && is_start_tag_token_with("svg")) ||
			(p.get_adjusted_current_node().IsHtmlIntegrationPoint() && is_start_tag_token()) ||
			(p.get_adjusted_current_node().IsHtmlIntegrationPoint() && is_char_token()) ||
			(p.get_adjusted_current_node().IsHtmlIntegrationPoint() && is_eof_token()) {
			p.apply_current_insertion_mode_rules(tk)
		} else {
			// TODO: Process the token according to the rules given in the section for parsing tokens in foreign content.
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#tree-construction-dispatcher]")
		}
	}
	p.run_parser = true
	for p.run_parser {
		p.tokenizer.run()
	}
	return p.document
}
func (p *html_parser) parse_error_encountered(tk html_token) {
	log.Println("Parse error occured near", tk)

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("-> From %s:%d (%s)", file, line, runtime.FuncForPC(pc).Name())
	}
}

type html_parser_control uint8

const (
	html_parser_control_ignore_token = html_parser_control(iota)
	html_parser_control_continue
)

type html_parser_insertion_mode uint8

const (
	html_parser_insertion_mode_initial = html_parser_insertion_mode(iota)
	html_parser_insertion_mode_before_html
	html_parser_insertion_mode_before_head
	html_parser_insertion_mode_in_head
	html_parser_insertion_mode_in_head_noscript
	html_parser_insertion_mode_after_head
	html_parser_insertion_mode_in_body
	html_parser_insertion_mode_text
	html_parser_insertion_mode_in_table
	html_parser_insertion_mode_in_table_text
	html_parser_insertion_mode_in_caption
	html_parser_insertion_mode_in_column_group
	html_parser_insertion_mode_in_table_body
	html_parser_insertion_mode_in_row
	html_parser_insertion_mode_in_cell
	html_parser_insertion_mode_in_template
	html_parser_insertion_mode_after_body
	html_parser_insertion_mode_in_frameset
	html_parser_insertion_mode_after_frameset
	html_parser_insertion_mode_after_after_body
	html_parser_insertion_mode_after_after_frameset
)

type html_active_formatting_elem struct {
	elem  dom.Element // If elem is nil, this is a marker.
	token tagToken
}

func (e html_active_formatting_elem) is_marker() bool { return cm.IsNil(e.elem) }

var html_active_formatting_elem_marker html_active_formatting_elem

func (p *html_parser) last_marker_in_list_of_active_formatting_elems() (elem *html_active_formatting_elem, idx int) {
	idx = slices.IndexFunc(p.list_of_active_formatting_elems, html_active_formatting_elem.is_marker)
	if idx == -1 {
		return nil, -1
	} else {
		return &p.list_of_active_formatting_elems[idx], idx
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#push-onto-the-list-of-active-formatting-elements
func (p *html_parser) push_to_list_of_active_formatting_elems(elem dom.Element) {
	last_marker, last_marker_idx := p.last_marker_in_list_of_active_formatting_elems()
	check_fn := func(other_elem dom.Element) bool {
		if elem.LocalName() != other_elem.LocalName() {
			return false
		}
		elem_ns, elem_has_ns := elem.Namespace()
		other_ns, other_has_ns := other_elem.Namespace()
		if elem_has_ns != other_has_ns {
			return false
		}
		if elem_has_ns && elem_ns != other_ns {
			return false
		}
		attrs := elem.Attrs()
		other_attrs := other_elem.Attrs()
		if len(attrs) != len(other_attrs) {
			return false
		}
		for i := 0; i < len(attrs); i++ {
			if attrs[i].LocalName() == other_attrs[i].LocalName() &&
				attrs[i].Value() == other_attrs[i].Value() {
				return true
			}
		}
		return false
	}
	matching_item_indices := []int{}
	check_start_idx := 0
	if last_marker != nil {
		check_start_idx = last_marker_idx + 1
	}
	for i := check_start_idx; i < len(p.list_of_active_formatting_elems); i++ {
		if check_fn(p.list_of_active_formatting_elems[i].elem) {
			matching_item_indices = append(matching_item_indices, i)
		}
	}
	if 3 <= len(matching_item_indices) {
		p.list_of_active_formatting_elems = append(
			p.list_of_active_formatting_elems[:matching_item_indices[0]],
			p.list_of_active_formatting_elems[matching_item_indices[0]+1:]...,
		)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#reconstruct-the-active-formatting-elements
func (p *html_parser) reconstruct_active_formatting_elems() {
	if len(p.list_of_active_formatting_elems) == 0 {
		return
	}
	last_entry := p.list_of_active_formatting_elems[len(p.list_of_active_formatting_elems)-1]
	if last_entry.is_marker() || slices.Contains(p.stack_of_open_elems, last_entry.elem) {
		return
	}
	entry_idx := len(p.list_of_active_formatting_elems) - 1
	for {
		entry := func() *html_active_formatting_elem { return &p.list_of_active_formatting_elems[entry_idx] }
	rewind:
		if entry_idx == 0 {
			goto create
		}
		entry_idx = entry_idx - 1
		if !entry().is_marker() && slices.Contains(p.stack_of_open_elems, entry().elem) {
			goto rewind
		}
	advance:
		entry_idx = entry_idx + 1
	create:
		new_elem := p.insert_html_element(entry().token)
		entry().elem = new_elem
		if entry_idx != len(p.list_of_active_formatting_elems)-1 {
			goto advance
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#clear-the-list-of-active-formatting-elements-up-to-the-last-marker
func (p *html_parser) clear_list_of_active_formatting_elems_up_to_last_marker() {
	for {
		last_entry := p.list_of_active_formatting_elems[len(p.list_of_active_formatting_elems)-1]
		p.list_of_active_formatting_elems = p.list_of_active_formatting_elems[:len(p.list_of_active_formatting_elems)-1]
		if last_entry.is_marker() {
			break
		}
	}
}

// Returns an item from stack of open elements.
// - Positive index starts from the top of the stack (first pushed item first).
// - Negative index starts from the bottom of the stack (most recent item first).
func (p *html_parser) get_soe_node(idx int) dom.Element {
	if 0 < idx {
		return p.stack_of_open_elems[idx]
	} else if idx < 0 {
		return p.stack_of_open_elems[len(p.stack_of_open_elems)+idx]
	} else {
		panic("zero index is not allowed")
	}
}
func (p *html_parser) get_current_node() dom.Element {
	return p.get_soe_node(-1)
}
func (p *html_parser) get_adjusted_current_node() dom.Element {
	if p.is_fragment_parsing {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#adjusted-current-node]")
	} else {
		return p.get_current_node()
	}
}
func (p *html_parser) push_node_to_soe(node dom.Element) {
	p.stack_of_open_elems = append(p.stack_of_open_elems, node)
}
func (p *html_parser) pop_node_from_soe() dom.Element {
	// TODO: https://html.spec.whatwg.org/multipage/parsing.html#the-stack-of-open-elements
	// When the current node is removed from the stack of open elements, process internal resource links given the current node's node document.
	node := p.stack_of_open_elems[len(p.stack_of_open_elems)-1]
	p.stack_of_open_elems = p.stack_of_open_elems[:len(p.stack_of_open_elems)-1]
	cb := node.Callbacks().PoppedFromStackOfOpenElements
	if cb != nil {
		cb()
	}
	return node
}
func (p *html_parser) remove_from_soe(idx int) {
	p.stack_of_open_elems = append(p.stack_of_open_elems[:idx], p.stack_of_open_elems[idx+1:]...)
}

func (p *html_parser) push_node_to_sot(mode html_parser_insertion_mode) {
	p.stack_of_template_insertion_modes = append(p.stack_of_template_insertion_modes, mode)
}
func (p *html_parser) pop_node_from_sot() html_parser_insertion_mode {
	// TODO: https://html.spec.whatwg.org/multipage/parsing.html#the-stack-of-open-elements
	// When the current node is removed from the stack of open elements, process internal resource links given the current node's node document.
	node := p.stack_of_template_insertion_modes[len(p.stack_of_template_insertion_modes)-1]
	p.stack_of_template_insertion_modes = p.stack_of_template_insertion_modes[:len(p.stack_of_template_insertion_modes)-1]
	return node
}

func (p *html_parser) have_element_in_specific_scope(is_target_node func(n dom.Element) bool, elem_types []dom.NamePair) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-the-specific-scope
	node_idx := len(p.stack_of_open_elems) - 1
	for {
		node := p.stack_of_open_elems[node_idx]
		if is_target_node(node) {
			return true
		}
		if slices.ContainsFunc(elem_types, node.IsElement) {
			return false
		}
		node_idx--
	}
}
func (p *html_parser) have_element_in_scope(is_target_node func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-scope
	return p.have_element_in_specific_scope(is_target_node, []dom.NamePair{
		{Namespace: namespaces.Html, LocalName: "applet"}, {Namespace: namespaces.Html, LocalName: "caption"},
		{Namespace: namespaces.Html, LocalName: "html"}, {Namespace: namespaces.Html, LocalName: "table"},
		{Namespace: namespaces.Html, LocalName: "td"}, {Namespace: namespaces.Html, LocalName: "th"},
		{Namespace: namespaces.Html, LocalName: "marquee"}, {Namespace: namespaces.Html, LocalName: "object"},
		{Namespace: namespaces.Html, LocalName: "select"}, {Namespace: namespaces.Html, LocalName: "template"},
		{Namespace: namespaces.Mathml, LocalName: "mi"}, {Namespace: namespaces.Mathml, LocalName: "mo"},
		{Namespace: namespaces.Mathml, LocalName: "mn"}, {Namespace: namespaces.Mathml, LocalName: "ms"},
		{Namespace: namespaces.Mathml, LocalName: "mtext"}, {Namespace: namespaces.Mathml, LocalName: "annotation-xml"},
		{Namespace: namespaces.Svg, LocalName: "foreignObject"}, {Namespace: namespaces.Svg, LocalName: "desc"},
		{Namespace: namespaces.Svg, LocalName: "title"},
	})
}
func (p *html_parser) have_element_in_list_item_scope(is_target_node func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-list-item-scope
	return p.have_element_in_specific_scope(is_target_node, []dom.NamePair{
		{Namespace: namespaces.Html, LocalName: "ol"}, {Namespace: namespaces.Html, LocalName: "ul"},
		// Below are the same as above "element scope"
		{Namespace: namespaces.Html, LocalName: "applet"}, {Namespace: namespaces.Html, LocalName: "caption"},
		{Namespace: namespaces.Html, LocalName: "html"}, {Namespace: namespaces.Html, LocalName: "table"},
		{Namespace: namespaces.Html, LocalName: "td"}, {Namespace: namespaces.Html, LocalName: "th"},
		{Namespace: namespaces.Html, LocalName: "marquee"}, {Namespace: namespaces.Html, LocalName: "object"},
		{Namespace: namespaces.Html, LocalName: "select"}, {Namespace: namespaces.Html, LocalName: "template"},
		{Namespace: namespaces.Mathml, LocalName: "mi"}, {Namespace: namespaces.Mathml, LocalName: "mo"},
		{Namespace: namespaces.Mathml, LocalName: "mn"}, {Namespace: namespaces.Mathml, LocalName: "ms"},
		{Namespace: namespaces.Mathml, LocalName: "mtext"}, {Namespace: namespaces.Mathml, LocalName: "annotation-xml"},
		{Namespace: namespaces.Svg, LocalName: "foreignObject"}, {Namespace: namespaces.Svg, LocalName: "desc"},
		{Namespace: namespaces.Svg, LocalName: "title"},
	})
}
func (p *html_parser) have_element_in_button_scope(is_target_node func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-button-scope
	return p.have_element_in_specific_scope(is_target_node, []dom.NamePair{
		{Namespace: namespaces.Html, LocalName: "button"},
		// Below are the same as above "element scope"
		{Namespace: namespaces.Html, LocalName: "applet"}, {Namespace: namespaces.Html, LocalName: "caption"},
		{Namespace: namespaces.Html, LocalName: "html"}, {Namespace: namespaces.Html, LocalName: "table"},
		{Namespace: namespaces.Html, LocalName: "td"}, {Namespace: namespaces.Html, LocalName: "th"},
		{Namespace: namespaces.Html, LocalName: "marquee"}, {Namespace: namespaces.Html, LocalName: "object"},
		{Namespace: namespaces.Html, LocalName: "select"}, {Namespace: namespaces.Html, LocalName: "template"},
		{Namespace: namespaces.Mathml, LocalName: "mi"}, {Namespace: namespaces.Mathml, LocalName: "mo"},
		{Namespace: namespaces.Mathml, LocalName: "mn"}, {Namespace: namespaces.Mathml, LocalName: "ms"},
		{Namespace: namespaces.Mathml, LocalName: "mtext"}, {Namespace: namespaces.Mathml, LocalName: "annotation-xml"},
		{Namespace: namespaces.Svg, LocalName: "foreignObject"}, {Namespace: namespaces.Svg, LocalName: "desc"},
		{Namespace: namespaces.Svg, LocalName: "title"},
	})
}
func (p *html_parser) have_element_in_table_scope(is_target_node func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-table-scope
	return p.have_element_in_specific_scope(is_target_node, []dom.NamePair{
		{Namespace: namespaces.Html, LocalName: "html"}, {Namespace: namespaces.Html, LocalName: "table"}, {Namespace: namespaces.Html, LocalName: "template"},
	})
}

type html_parser_insertion_location struct {
	parent_node dom.Node
	tp          html_adjusted_insertion_location_type
}
type html_adjusted_insertion_location_type uint8

const (
	html_adjusted_insertion_location_type_after_last_child = html_adjusted_insertion_location_type(iota)
)

// https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node
//
// override_target may be nil pointer
func (p *html_parser) get_appropriate_place_for_inserting_node(override_target dom.Element) html_parser_insertion_location {
	var res html_parser_insertion_location
	target := override_target
	if cm.IsNil(target) {
		target = p.get_current_node()
	}

	if target_elem := target; p.enable_foster_parenting && (target_elem.IsHtmlElement("table") ||
		target_elem.IsHtmlElement("tbody") ||
		target_elem.IsHtmlElement("tfoot") ||
		target_elem.IsHtmlElement("thead") ||
		target_elem.IsHtmlElement("tr")) {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node]")
	} else {
		res = html_parser_insertion_location{target, html_adjusted_insertion_location_type_after_last_child}
	}
	if target_elem := target; target_elem.IsInside(dom.NamePair{Namespace: namespaces.Html, LocalName: "template"}) {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node]")
	}
	return res
}

// https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token
func (p *html_parser) create_element_for_token(token tagToken, namespace namespaces.Namespace, intended_parent dom.Node) dom.Element {
	if p.has_active_speculative_parser {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	document := intended_parent.NodeDocument()
	localName := token.tag_name
	is_val, has_is := token.Attr("is")
	var is *string
	if has_is {
		is = &is_val
	}
	registry := dom.LookupCustomElementRegistry(intended_parent)
	definition := registry.LookupCustomElementDefinition(&namespace, localName, is)
	will_execute_script := false
	if definition != nil && !p.is_fragment_parsing {
		will_execute_script = true
	}
	if will_execute_script {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	elem := dom.CreateElement(document, localName, &namespace, nil, is, will_execute_script, registry, token, func(namespace *namespaces.Namespace, localName string) func(opt dom.ElementCreationCommonOptions) dom.Element {
		factory_fn := func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLElement(opt) }
		if namespace != nil && *namespace == namespaces.Html && localName == "html" {
			factory_fn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLHtmlElement(opt) }
		} else if namespace != nil && *namespace == namespaces.Html && localName == "body" {
			factory_fn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLBodyElement(opt) }
		} else if namespace != nil && *namespace == namespaces.Html && localName == "link" {
			factory_fn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLLinkElement(opt) }
		} else if namespace != nil && *namespace == namespaces.Html && localName == "style" {
			factory_fn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLStyleElement(opt) }
		}
		return factory_fn
	})
	for _, attr := range token.attrs {
		elem.AppendAttr(attr)
	}
	if will_execute_script {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	if attr, ok := elem.AttrWithNamespace(dom.NamePair{Namespace: namespaces.Xmlns, LocalName: "xmlns"}); ok {
		if ns, ok := elem.Namespace(); !ok || (attr != string(ns)) {
			p.parse_error_encountered(token)
		}
	}
	if attr, ok := elem.AttrWithNamespace(dom.NamePair{Namespace: namespaces.Xmlns, LocalName: "xmlns:xlink"}); ok && attr != string(namespaces.Xlink) {
		p.parse_error_encountered(token)
	}
	if elem.(elements.HTMLElement).IsFormResettableElement() && !elem.(elements.HTMLElement).IsFormAssociatedCustomElement() {
		// TODO: Invoke reset algorithm
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	has_attr := func(name string) bool {
		_, ok := elem.AttrWithoutNamespace("form")
		return ok
	}
	if elem.(elements.HTMLElement).IsFormAssociatedElement() &&
		!cm.IsNil(p.form_element_pointer) &&
		!slices.ContainsFunc(p.stack_of_open_elems, func(n dom.Element) bool { return n.IsHtmlElement("template") }) &&
		(elem.(elements.HTMLElement).IsFormListedElement() || !has_attr("form")) &&
		dom.InTheSameTreeAs(intended_parent, p.form_element_pointer) {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	return elem
}

func (p *html_parser) insert_at_location(elem dom.Node, position html_parser_insertion_location) {
	switch position.tp {
	case html_adjusted_insertion_location_type_after_last_child:
		dom.AppendChild(position.parent_node, elem)
	default:
		log.Panicf("unknown insertion mode %v", position.tp)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-an-element-at-the-adjusted-insertion-location
func (p *html_parser) insert_element_at_adjusted_insertion_location(elem dom.Node) {
	insertion_location := p.get_appropriate_place_for_inserting_node(nil)
	if !p.is_fragment_parsing {
		// TODO: push a new element queue onto element's relevant agent's custom element reactions stack.
	}
	p.insert_at_location(elem, insertion_location)
	if !p.is_fragment_parsing {
		// TODO: pop the element queue from element's relevant agent's custom element reactions stack,
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-a-foreign-element
func (p *html_parser) insert_foreign_element(token tagToken, namespace namespaces.Namespace, only_add_element_to_stack bool) dom.Element {
	insertion_location := p.get_appropriate_place_for_inserting_node(nil)
	elem := p.create_element_for_token(token, namespace, insertion_location.parent_node)
	if !only_add_element_to_stack {
		p.insert_element_at_adjusted_insertion_location(elem)
	}
	p.push_node_to_soe(elem)
	return elem
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-an-html-element
func (p *html_parser) insert_html_element(token tagToken) dom.Element {
	return p.insert_foreign_element(token, namespaces.Html, false)
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-a-comment
//
// position may be nil(= insert_comment() will figure it out)
func (p *html_parser) insert_comment(data string, position *html_parser_insertion_location) {
	if position == nil {
		position = new(html_parser_insertion_location)
		*position = p.get_appropriate_place_for_inserting_node(nil)
	}
	comment := dom.NewComment(position.parent_node.NodeDocument(), data)
	p.insert_at_location(comment, *position)
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-a-character
func (p *html_parser) insert_character(data rune) {
	insertion_location := p.get_appropriate_place_for_inserting_node(nil)
	if _, ok := insertion_location.parent_node.(dom.Document); ok {
		// Document node cannot have text as children
		return
	}
	switch insertion_location.tp {
	case html_adjusted_insertion_location_type_after_last_child:
		parent_node := insertion_location.parent_node
		parent_children := parent_node.Children()
		var existing_text dom.Text
		if len(parent_children) != 0 {
			if t, ok := parent_children[len(parent_children)-1].(dom.Text); ok {
				existing_text = t
			}
		}

		if !cm.IsNil(existing_text) {
			existing_text.AppendText(string(data))
		} else {
			text := dom.NewText(parent_node.NodeDocument(), string(data))
			p.insert_at_location(text, insertion_location)
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#generic-raw-text-element-parsing-algorithm
func (p *html_parser) parse_generic_raw_text_element(token tagToken) {
	p.insert_html_element(token)
	p.tokenizer.state = html_tokenizer_rawtext_state
	p.original_insertion_mode = p.insertion_mode
	p.insertion_mode = html_parser_insertion_mode_text
}

// https://html.spec.whatwg.org/multipage/parsing.html#generic-raw-text-element-parsing-algorithm
func (p *html_parser) parse_generic_rcdata_element(token tagToken) {
	p.insert_html_element(token)
	p.tokenizer.state = html_tokenizer_rcdata_state
	p.original_insertion_mode = p.insertion_mode
	p.insertion_mode = html_parser_insertion_mode_text
}

// https://html.spec.whatwg.org/multipage/parsing.html#generate-implied-end-tags
func (p *html_parser) generate_implied_end_tags(exclude_filter func(node dom.Element) bool) {
	html_elems := []string{
		"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp", "rt", "rtc",
	}
	for {
		current_node := p.get_current_node()
		if slices.ContainsFunc(html_elems, current_node.IsHtmlElement) &&
			(exclude_filter == nil || !exclude_filter(current_node)) {
			p.pop_node_from_soe()
		} else {
			break
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#generate-all-implied-end-tags-thoroughly
func (p *html_parser) generate_all_implied_end_tags_thoroughly(exclude_filter func(n dom.Element) bool) {
	html_elems := []string{
		"caption", "colgroup", "dd", "dt", "li", "optgroup", "option", "p",
		"rb", "rp", "rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
	}
	for {
		current_node_p := p.get_current_node()
		current_node := current_node_p
		if slices.ContainsFunc(html_elems, current_node.IsHtmlElement) &&
			!exclude_filter(current_node_p) {
			p.pop_node_from_soe()
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#reset-the-insertion-mode-appropriately
func (p *html_parser) reset_insertion_mode_appropriately() {
	panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#reset-the-insertion-mode-appropriately]")
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-initial-insertion-mode
func (p *html_parser) apply_initial_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		return
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, &html_parser_insertion_location{p.document, html_adjusted_insertion_location_type_after_last_child})
	} else if tk, ok := token.(*html_doctype_token); ok {
		if tk.name == nil || *tk.name != "html" || tk.public_id != nil || (tk.system_id != nil && *tk.system_id != "about:legacy-compat") {
			p.parse_error_encountered(token)
		}
		var name, public_id, system_id string = "", "", ""
		if tk.name != nil {
			name = *tk.name
		}
		if tk.public_id != nil {
			public_id = *tk.public_id
		}
		if tk.system_id != nil {
			system_id = *tk.system_id
		}

		doctype_node := dom.NewDocumentType(p.document, name, public_id, system_id)
		dom.AppendChild(p.document, doctype_node)

		p.document.SetMode(dom.NoQuirks)
		if !p.document.IsIframeSrcdocDocument() && !p.document.IsParserCannotChangeMode() {
			if tk.force_quirks ||
				(tk.name == nil || *tk.name != "html") ||
				(tk.public_id != nil && cm.ToAsciiLowercase(*tk.public_id) == "-//w3o//dtd w3 html strict 3.0//en//") ||
				(tk.public_id != nil && cm.ToAsciiLowercase(*tk.public_id) == "-/w3c/dtd html 4.0 transitional/en") ||
				(tk.public_id != nil && cm.ToAsciiLowercase(*tk.public_id) == "html") ||
				(tk.system_id != nil && cm.ToAsciiLowercase(*tk.system_id) == "http://www.ibm.com/data/dtd/v11/ibmxhtml1-transitional.dtd") ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "+//silmaril//dtd html pro v0r11 19970101//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//as//dtd html 3.0 aswedit + extensions//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//advasoft ltd//dtd html 3.0 aswedit + extensions//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.0 level 1//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.0 level 2//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.0 strict level 1//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.0 strict level 2//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.0 strict//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 2.1e//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 3.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 3.2 final//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 3.2//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html 3//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html level 0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html level 1//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html level 2//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html level 3//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html strict level 0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html strict level 1//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html strict level 2//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html strict level 3//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html strict//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//ietf//dtd html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//metrius//dtd metrius presentational//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//microsoft//dtd internet explorer 2.0 html strict//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//microsoft//dtd internet explorer 2.0 html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//microsoft//dtd internet explorer 2.0 tables//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//microsoft//dtd internet explorer 3.0 html strict//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//microsoft//dtd internet explorer 3.0 html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//microsoft//dtd internet explorer 3.0 tables//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//netscape comm. corp.//dtd html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//netscape comm. corp.//dtd strict html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//o'reilly and associates//dtd html 2.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//o'reilly and associates//dtd html extended 1.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//o'reilly and associates//dtd html extended relaxed 1.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//sq//dtd html 2.0 hotmetal + extensions//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//softquad software//dtd hotmetal pro 6.0::19990601::extensions to html 4.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//softquad//dtd hotmetal pro 4.0::19971010::extensions to html 4.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//spyglass//dtd html 2.0 extended//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//sun microsystems corp.//dtd hotjava html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//sun microsystems corp.//dtd hotjava strict html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 3 1995-03-24//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 3.2 draft//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 3.2 final//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 3.2//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 3.2s draft//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 4.0 frameset//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 4.0 transitional//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html experimental 19960712//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html experimental 970421//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd w3 html//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3o//dtd w3 html 3.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//webtechs//dtd mozilla html 2.0//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//webtechs//dtd mozilla html//")) ||
				(tk.system_id == nil && tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 4.01 frameset//")) ||
				(tk.system_id == nil && tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 4.01 transitional//")) {
				p.document.SetMode(dom.Quirks)
			} else if (tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd xhtml 1.0 frameset//")) ||
				(tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd xhtml 1.0 transitional//")) ||
				(tk.system_id != nil && tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 4.01 frameset//")) ||
				(tk.system_id != nil && tk.public_id != nil && strings.HasPrefix(cm.ToAsciiLowercase(*tk.public_id), "-//w3c//dtd html 4.01 transitional//")) {
				p.document.SetMode(dom.LimitedQuirks)
			}
		}
		p.insertion_mode = html_parser_insertion_mode_before_html
	} else {
		if !p.document.IsIframeSrcdocDocument() {
			p.parse_error_encountered(token)
			if !p.document.IsParserCannotChangeMode() {
				p.document.SetMode(dom.Quirks)
			}
		}
		p.insertion_mode = html_parser_insertion_mode_before_html
		p.apply_before_html_insertion_mode_rules(token)
		return
	}

}

// https://html.spec.whatwg.org/multipage/parsing.html#the-before-html-insertion-mode
func (p *html_parser) apply_before_html_insertion_mode_rules(token html_token) {
	if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, &html_parser_insertion_location{p.document, html_adjusted_insertion_location_type_after_last_child})
	} else if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		elem := p.create_element_for_token(*tk, namespaces.Html, p.document)
		dom.AppendChild(p.document, elem)
		p.push_node_to_soe(elem)
		p.insertion_mode = html_parser_insertion_mode_before_head
	} else if tk, ok := token.(*tagToken); ok && tk.is_end && !slices.Contains([]string{"head", "body", "html", "br"}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else {
		elem := elements.NewHTMLHtmlElement(dom.ElementCreationCommonOptions{
			NodeDocument: p.document,
			Namespace:    namespaces.HtmlP(),
			LocalName:    "html",
		})
		dom.AppendChild(p.document, elem)
		p.push_node_to_soe(elem)
		p.insertion_mode = html_parser_insertion_mode_before_head
		p.apply_before_head_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-before-head-insertion-mode
func (p *html_parser) apply_before_head_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		return
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "head" {
		elem := p.insert_html_element(*tk)
		p.head_element_pointer = elem
		p.insertion_mode = html_parser_insertion_mode_in_head
	} else if tk, ok := token.(*tagToken); ok && tk.is_end && !slices.Contains([]string{"head", "body", "html", "br"}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else {
		elem := p.insert_html_element(tagToken{tag_name: "head"})
		p.head_element_pointer = elem
		p.insertion_mode = html_parser_insertion_mode_in_head
		p.apply_in_head_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead
func (p *html_parser) apply_in_head_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok {
		p.insert_character(tk.value)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"base", "basefont", "bgsound", "link"}, tk.tag_name) {
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		if tk.is_self_closing {
			tk.self_closing_acknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "meta" {
		elem := p.insert_html_element(*tk)
		p.pop_node_from_soe()
		if !p.has_active_speculative_parser {
			elem := elem
			if attr, ok := elem.AttrWithoutNamespace("charset"); ok {
				_ = attr
				// TODO: Set encoding based on charset
			}
			if attr, ok := elem.AttrWithoutNamespace("http-equiv"); ok && cm.ToAsciiLowercase(attr) == "content-type" {
				_ = attr
				// TODO: Set encoding based on http-equiv Content-Type value
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "title" {
		p.parse_generic_rcdata_element(*tk)
	} else if tk, ok := token.(*tagToken); ok &&
		(((tk.is_start_tag() && tk.tag_name == "title") && p.enable_scripting) ||
			(tk.is_start_tag() && slices.Contains([]string{"noframes", "style"}, tk.tag_name))) {
		p.parse_generic_raw_text_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "noscript" && !p.enable_scripting {
		p.insert_html_element(*tk)
		p.insertion_mode = html_parser_insertion_mode_in_head_noscript
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "script" {
		// STUB
		p.parse_generic_raw_text_element(*tk)
		// panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead]")
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "head" {
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_after_head
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "template" {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead]")
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "template" {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead]")
	} else if tk, ok := token.(*tagToken); ok &&
		((tk.is_end && !slices.Contains([]string{"body", "html", "br"}, tk.tag_name)) ||
			tk.is_start_tag() && tk.tag_name == "head") {
		p.parse_error_encountered(token)
		return
	} else {
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_after_head
		p.apply_after_head_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inheadnoscript
func (p *html_parser) apply_in_head_noscript_insertion_mode_rules(token html_token) {
	if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "noscript" {
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_head
	} else if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.apply_in_head_insertion_mode_rules(token)
	} else if _, ok := token.(*html_comment_token); ok {
		p.apply_in_head_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"basefont", "bgsound", "link", "meta", "noframes", "style"}, tk.tag_name) {
		p.apply_in_head_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_end && !tk.is_end_tag() && tk.tag_name == "br") ||
		(tk.is_start_tag() && slices.Contains([]string{"head", "noscript"}, tk.tag_name)) {
		p.parse_error_encountered(token)
		return
	} else {
		p.parse_error_encountered(token)
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_head
		p.apply_in_head_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-after-head-insertion-mode
func (p *html_parser) apply_after_head_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.insert_character(tk.value)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "body" {
		p.insert_html_element(*tk)
		p.is_frameset_not_ok = true
		p.insertion_mode = html_parser_insertion_mode_in_body
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "body" {
		p.insert_html_element(*tk)
		p.insertion_mode = html_parser_insertion_mode_in_frameset
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{
		"base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title",
	}, tk.tag_name) {
		p.parse_error_encountered(token)
		p.push_node_to_soe(p.head_element_pointer)
		p.insertion_mode = html_parser_insertion_mode_in_head
		remove_idx := slices.Index(p.stack_of_open_elems, p.head_element_pointer)
		p.remove_from_soe(remove_idx)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "template" {
		p.apply_in_head_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok &&
		((tk.is_end && !slices.Contains([]string{"body", "html", "br"}, tk.tag_name)) ||
			(tk.is_start_tag() && tk.tag_name == "head")) {
		p.parse_error_encountered(token)
		return
	} else {
		elem := p.insert_html_element(tagToken{tag_name: "body"})
		p.head_element_pointer = elem
		p.insertion_mode = html_parser_insertion_mode_in_body
		p.apply_in_body_insertion_mode_rules(token)
	}
}

func (p *html_parser) soe_has_one_of_elems(elems []string) bool {
	return slices.ContainsFunc(p.stack_of_open_elems, func(n dom.Element) bool {
		return !slices.ContainsFunc(elems, n.IsHtmlElement)
	})
}
func (p *html_parser) soe_has_elem(elem string) bool {
	return p.soe_has_one_of_elems([]string{elem})
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody
func (p *html_parser) apply_in_body_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\u0000") {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.reconstruct_active_formatting_elems()
		p.insert_character(tk.value)
	} else if tk, ok := token.(*html_char_token); ok {
		p.reconstruct_active_formatting_elems()
		p.insert_character(tk.value)
		p.is_frameset_not_ok = true
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.parse_error_encountered(token)
		if p.soe_has_elem("template") {
			return
		} else {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && slices.Contains([]string{"base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title"}, tk.tag_name) ||
			(tk.is_end_tag() && tk.tag_name == "template")) {
		p.apply_in_head_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "body" {
		p.parse_error_encountered(token)
		if len(p.stack_of_open_elems) == 1 ||
			!p.stack_of_open_elems[1].IsHtmlElement("body") ||
			p.soe_has_elem("template") {
			return
		} else {
			p.is_frameset_not_ok = true
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "frameset" {
		p.parse_error_encountered(token)
		if len(p.stack_of_open_elems) == 1 ||
			!p.stack_of_open_elems[1].IsHtmlElement("body") {
			return
		} else if !p.is_frameset_not_ok {
			return
		} else {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
	} else if _, ok := token.(*html_eof_token); ok {
		if len(p.stack_of_template_insertion_modes) != 0 {
			p.apply_in_template_insertion_mode_rules(token)
		} else {
			if p.soe_has_one_of_elems([]string{
				"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
				"rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
				"body", "html",
			}) {
				p.parse_error_encountered(token)
			}
			p.stop_parsing()
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "body" {
		if !p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("body") }) {
			p.parse_error_encountered(token)
			return
		} else if p.soe_has_one_of_elems([]string{
			"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
			"rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
			"body", "html",
		}) {
			p.parse_error_encountered(token)
		}
		p.insertion_mode = html_parser_insertion_mode_after_body
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "html" {
		if !p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("body") }) {
			p.parse_error_encountered(token)
			return
		} else if p.soe_has_one_of_elems([]string{
			"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
			"rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
			"body", "html",
		}) {
			p.parse_error_encountered(token)
		}
		p.insertion_mode = html_parser_insertion_mode_after_body
		p.apply_after_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{
		"address", "article", "aside", "blockquote", "center", "details",
		"dialog", "dir", "div", "dl", "fieldset", "figcaption", "figure",
		"footer", "header", "hgroup", "main", "menu", "nav", "ol", "p",
		"search", "section", "summary", "ul",
	}, tk.tag_name) {
		if p.have_element_in_button_scope(func(n dom.Element) bool {
			return n.IsHtmlElement("p")
		}) {
			p.close_p_element()
		}
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, tk.tag_name) {
		if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		if slices.ContainsFunc([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, p.get_current_node().IsHtmlElement) {
			p.parse_error_encountered(token)
			p.pop_node_from_soe()
		}
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"pre", "listing"}, tk.tag_name) {
		if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		p.insert_html_element(*tk)
		p.on_next_token = func(token html_token) html_parser_control {
			if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\n") {
				return html_parser_control_ignore_token
			}
			return html_parser_control_continue
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "form" {
		if !cm.IsNil(p.form_element_pointer) && !p.soe_has_elem("template") {
			p.parse_error_encountered(token)
			return
		} else {
			if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
				p.close_p_element()
			}
			elem := p.insert_html_element(*tk)
			if !p.soe_has_elem("template") {
				p.form_element_pointer = elem
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "li" {
		p.is_frameset_not_ok = true
		node := p.get_current_node()
		for {
			if node.IsHtmlElement("li") {
				p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement("li") })
				if !p.get_current_node().IsHtmlElement("li") {
					p.parse_error_encountered(token)
				}
				for {
					popped_elem := p.pop_node_from_soe()
					if popped_elem.IsHtmlElement("li") {
						break
					}
				}
				break
			}
			if node.IsHtmlSpecialElement() &&
				!slices.ContainsFunc([]string{"address", "div", "p"}, node.IsHtmlElement) {
				break
			} else {
				node_idx := slices.Index(p.stack_of_open_elems, node) - 1
				node = p.stack_of_open_elems[node_idx]
			}
		}
		if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"dt", "dd"}, tk.tag_name) {
		p.is_frameset_not_ok = true
		node := p.get_current_node()
		for {
			if node.IsHtmlElement("dd") {
				p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement("dd") })
				if !p.get_current_node().IsHtmlElement("dd") {
					p.parse_error_encountered(token)
				}
				for {
					popped_elem := p.pop_node_from_soe()
					if popped_elem.IsHtmlElement("dd") {
						break
					}
				}
				break
			} else if node.IsHtmlElement("dt") {
				p.generate_implied_end_tags(func(node dom.Element) bool { return node.IsHtmlElement("dt") })
				if !p.get_current_node().IsHtmlElement("dt") {
					p.parse_error_encountered(token)
				}
				for {
					popped_elem := p.pop_node_from_soe()
					if popped_elem.IsHtmlElement("dt") {
						break
					}
				}
				break
			}
			if node.IsHtmlSpecialElement() &&
				!slices.ContainsFunc([]string{"address", "div", "p"}, node.IsHtmlElement) {
				break
			} else {
				node_idx := slices.Index(p.stack_of_open_elems, node) - 1
				node = p.stack_of_open_elems[node_idx]
			}
		}
		if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "plaintext" {
		if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		p.insert_html_element(*tk)
		p.tokenizer.state = html_tokenizer_plaintext_state
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "button" {
		if !p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("button") }) {
			p.parse_error_encountered(token)
			p.generate_implied_end_tags(func(node dom.Element) bool { return true })
			for {
				popped_elem := p.pop_node_from_soe()
				if popped_elem.IsHtmlElement("button") {
					break
				}
			}
		}
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
		p.is_frameset_not_ok = true
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{
		"address", "article", "aside", "blockquote", "button", "center",
		"details", "dialog", "dir", "div", "dl", "fieldset", "figcaption",
		"figure", "footer", "header", "hgroup", "listing", "main", "menu",
		"nav", "ol", "pre", "search", "section", "select", "summary", "ul",
	}, tk.tag_name) {
		if !p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) }) {
			p.parse_error_encountered(token)
			return
		} else {
			p.generate_implied_end_tags(nil)
			if !p.get_current_node().IsHtmlElement(tk.tag_name) {
				p.parse_error_encountered(token)
			}
			for {
				popped_elem := p.pop_node_from_soe()
				if popped_elem.IsHtmlElement(tk.tag_name) {
					break
				}
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "form" {
		if p.soe_has_elem("template") {
			node := p.form_element_pointer
			p.form_element_pointer = nil
			if cm.IsNil(node) || !p.have_element_in_scope(func(n dom.Element) bool { return n == node }) {
				p.parse_error_encountered(token)
				return
			}
			p.generate_implied_end_tags(nil)
			if p.get_current_node() != node {
				p.parse_error_encountered(token)
			}
			remove_idx := slices.Index(p.stack_of_open_elems, node)
			p.remove_from_soe(remove_idx)
		} else {
			if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("form") }) {
				p.parse_error_encountered(token)
				return
			}
			p.generate_implied_end_tags(nil)
			if !p.get_current_node().IsHtmlElement("form") {
				p.parse_error_encountered(token)
			}
			for {
				popped_elem := p.pop_node_from_soe()
				if popped_elem.IsHtmlElement("form") {
					break
				}
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "p" {
		if !p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.parse_error_encountered(token)
			p.insert_html_element(tagToken{tag_name: "p"})
		}
		p.close_p_element()
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "li" {
		if !p.have_element_in_list_item_scope(func(n dom.Element) bool { return n.IsHtmlElement("li") }) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement("li") })
		if !p.get_current_node().IsHtmlElement("li") {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement("li") {
				break
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{"dd", "dt"}, tk.tag_name) {
		if !p.have_element_in_list_item_scope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) }) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) })
		if !p.get_current_node().IsHtmlElement(tk.tag_name) {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement(tk.tag_name) {
				break
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, tk.tag_name) {
		if !p.have_element_in_list_item_scope(func(n dom.Element) bool {
			return slices.ContainsFunc([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, n.IsHtmlElement)
		}) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(nil)
		if !p.get_current_node().IsHtmlElement(tk.tag_name) {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if slices.ContainsFunc([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, popped_elem.IsHtmlElement) {
				break
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "a" {
		{
			last_marker_idx := slices.IndexFunc(p.list_of_active_formatting_elems, html_active_formatting_elem.is_marker)
			check_start_idx := 0
			if 0 <= last_marker_idx {
				check_start_idx = last_marker_idx + 1
			}
			var a_elem dom.Element
			for i := check_start_idx; i < len(p.list_of_active_formatting_elems); i++ {
				if p.list_of_active_formatting_elems[i].elem.IsHtmlElement("a") {
					a_elem = p.list_of_active_formatting_elems[i].elem
				}
			}
			if !cm.IsNil(a_elem) {
				p.parse_error_encountered(token)
				p.adoption_agency_algorithm(*tk)
				remove_idx := slices.IndexFunc(p.list_of_active_formatting_elems, func(e html_active_formatting_elem) bool { return e.elem == a_elem })
				p.list_of_active_formatting_elems = append(p.list_of_active_formatting_elems[:remove_idx], p.list_of_active_formatting_elems[remove_idx+1:]...)
				remove_idx = slices.Index(p.stack_of_open_elems, a_elem)
				p.remove_from_soe(remove_idx)
			}
		}
		p.reconstruct_active_formatting_elems()
		elem := p.insert_html_element(*tk)
		p.push_to_list_of_active_formatting_elems(elem)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{
		"b", "big", "code", "em", "font", "i", "s", "small", "strike", "strong", "tt", "u",
	}, tk.tag_name) {
		p.reconstruct_active_formatting_elems()
		elem := p.insert_html_element(*tk)
		p.push_to_list_of_active_formatting_elems(elem)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "nobr" {
		p.reconstruct_active_formatting_elems()
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("node_p") }) {
			p.parse_error_encountered(token)
			p.adoption_agency_algorithm(*tk)
			p.reconstruct_active_formatting_elems()
		}
		elem := p.insert_html_element(*tk)
		p.push_to_list_of_active_formatting_elems(elem)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{
		"a", "b", "big", "code", "em", "font", "i", "nobr", "s", "small", "strike", "strong", "tt", "u",
	}, tk.tag_name) {
		p.adoption_agency_algorithm(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"applet", "marquee", "object"}, tk.tag_name) {
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
		p.list_of_active_formatting_elems = append(p.list_of_active_formatting_elems, html_active_formatting_elem_marker)
		p.is_frameset_not_ok = true
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{"applet", "marquee", "object"}, tk.tag_name) {
		if !p.get_current_node().IsHtmlElement(tk.tag_name) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(nil)
		if !p.get_current_node().IsHtmlElement(tk.tag_name) {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement(tk.tag_name) {
				break
			}
		}
		p.clear_list_of_active_formatting_elems_up_to_last_marker()
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "table" {
		if (p.document.Mode() != dom.Quirks) &&
			p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		p.insert_html_element(*tk)
		p.is_frameset_not_ok = true
		p.insertion_mode = html_parser_insertion_mode_in_table
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_end_tag() && tk.tag_name == "br") ||
		(tk.is_start_tag() && slices.Contains([]string{"area", "br", "embed", "img", "keygen", "wbr"}, tk.tag_name)) {
		if tk.is_end_tag() && tk.tag_name == "br" {
			p.parse_error_encountered(token)
			tk.attrs = []dom.AttrData{}
			tk.is_end = false
		}
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		tk.self_closing_acknowledged = true
		p.is_frameset_not_ok = true
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "input" {
		if p.is_fragment_parsing {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.parse_error_encountered(token)
			for {
				popped_elem := p.pop_node_from_soe()
				if popped_elem.IsHtmlElement("select") {
					break
				}
			}
		}
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		tk.self_closing_acknowledged = true
		if type_attr, ok := tk.Attr("type"); !ok || cm.ToAsciiLowercase(type_attr) != "hidden" {
			p.is_frameset_not_ok = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "hr" {
		if p.have_element_in_button_scope(func(n dom.Element) bool {
			return n.IsHtmlElement("p")
		}) {
			p.close_p_element()
		}
		if p.have_element_in_scope(func(n dom.Element) bool {
			return n.IsHtmlElement("select")
		}) {
			p.generate_implied_end_tags(nil)
			if p.have_element_in_scope(func(n dom.Element) bool {
				return n.IsHtmlElement("option") ||
					n.IsHtmlElement("optgroup")
			}) {
				p.parse_error_encountered(token)
			}
		}
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		tk.self_closing_acknowledged = true
		p.is_frameset_not_ok = true
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "image" {
		p.parse_error_encountered(token)
		tk.tag_name = "img"
		p.apply_in_body_insertion_mode_rules(tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "textarea" {
		p.insert_html_element(*tk)
		p.on_next_token = func(token html_token) html_parser_control {
			if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\n") {
				return html_parser_control_ignore_token
			}
			return html_parser_control_continue
		}
		p.tokenizer.state = html_tokenizer_rcdata_state
		p.original_insertion_mode = p.insertion_mode
		p.is_frameset_not_ok = true
		p.insertion_mode = html_parser_insertion_mode_text
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "xmp" {
		if p.have_element_in_button_scope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.close_p_element()
		}
		p.reconstruct_active_formatting_elems()
		p.is_frameset_not_ok = true
		p.parse_generic_raw_text_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "iframe" {
		p.is_frameset_not_ok = true
		p.parse_generic_raw_text_element(*tk)
	} else if tk, ok := token.(*tagToken); ok &&
		((tk.is_start_tag() && tk.tag_name == "noembed") ||
			(tk.is_start_tag() && tk.tag_name == "noscript" && !p.enable_scripting)) {
		p.parse_generic_raw_text_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "select" {
		if p.is_fragment_parsing {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.parse_error_encountered(token)
			for {
				popped_elem := p.pop_node_from_soe()
				if popped_elem.IsHtmlElement("select") {
					break
				}
			}
			return
		}
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
		p.is_frameset_not_ok = true
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "option" {
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement("opgroup") })
			if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("option") }) {
				p.parse_error_encountered(token)
			}
		} else {
			if p.get_current_node().IsHtmlElement("option") {
				p.pop_node_from_soe()
			}
		}
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "optgroup" {
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.generate_implied_end_tags(nil)
			if p.have_element_in_scope(func(n dom.Element) bool {
				return n.IsHtmlElement("option") ||
					n.IsHtmlElement("optgroup")
			}) {
				p.parse_error_encountered(token)
			}
		} else {
			if p.get_current_node().IsHtmlElement("option") {
				p.pop_node_from_soe()
			}
		}
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"rb", "rtc"}, tk.tag_name) {
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("ruby") }) {
			p.generate_implied_end_tags(nil)
			if !p.get_current_node().IsHtmlElement("ruby") {
				p.parse_error_encountered(token)
			}
		}
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"rp", "rt"}, tk.tag_name) {
		if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("ruby") }) {
			p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement("rtc") })
			if !p.get_current_node().IsHtmlElement("rtc") &&
				!p.get_current_node().IsHtmlElement("ruby") {
				p.parse_error_encountered(token)
			}
		}
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "math" {
		p.reconstruct_active_formatting_elems()
		html_parser_adjust_mathml_attrs(tk)
		parserAdjustForeignAttrs(tk)
		p.insert_foreign_element(*tk, namespaces.Mathml, false)
		if tk.is_self_closing {
			p.pop_node_from_soe()
			tk.self_closing_acknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "svg" {
		p.reconstruct_active_formatting_elems()
		adjustSvgAttrs(tk)
		parserAdjustForeignAttrs(tk)
		p.insert_foreign_element(*tk, namespaces.Svg, false)
		if tk.is_self_closing {
			p.pop_node_from_soe()
			tk.self_closing_acknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{
		"caption", "col", "colgroup", "frame", "head", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && !tk.is_end {
		p.reconstruct_active_formatting_elems()
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end {
		node_idx := len(p.stack_of_open_elems) - 1
		node := func() dom.Element {
			return p.stack_of_open_elems[node_idx]
		}
		for {
			if node().IsHtmlElement(tk.tag_name) {
				p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) })
				if node() != p.get_current_node() {
					p.parse_error_encountered(token)
				}
				target_node := node()
				for p.pop_node_from_soe() != target_node {
				}
				return
			}
			if node().IsHtmlSpecialElement() {
				p.parse_error_encountered(token)
				return
			}
			node_idx--
		}
	} else {
		log.Printf("[in-body insertion mode] Unrecognized token %v", token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#close-a-p-element
func (p *html_parser) close_p_element() {
	p.generate_implied_end_tags(func(n dom.Element) bool { return n.IsHtmlElement("p") })
	if !p.get_current_node().IsHtmlElement("p") {
		p.parse_error_encountered(p.get_current_node().TagToken().(*tagToken))
	}
	for {
		popped_elem := p.pop_node_from_soe()
		if popped_elem.IsHtmlElement("p") {
			break
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#adoption-agency-algorithm
func (p *html_parser) adoption_agency_algorithm(token tagToken) {
	subject := token.tag_name
	if p.get_current_node().IsHtmlElement(subject) &&
		!slices.ContainsFunc(p.list_of_active_formatting_elems, func(e html_active_formatting_elem) bool {
			return e.elem == p.get_current_node()
		}) {
		p.pop_node_from_soe()
		return
	}
	panic("TODO")
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incdata
func (p *html_parser) apply_text_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok {
		p.insert_character(tk.value)
	} else if _, ok := token.(*html_eof_token); ok {
		p.parse_error_encountered(token)
		if p.get_current_node().IsHtmlElement("script") {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incdata]")
		}
		p.pop_node_from_soe()
		p.insertion_mode = p.original_insertion_mode
		p.apply_current_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "script" {
		// STUB
		// panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incdata]")
	} else if tk, ok := token.(*tagToken); ok && tk.is_end {
		p.pop_node_from_soe()
		p.insertion_mode = p.original_insertion_mode
	} else {
		log.Printf("[text insertion mode] Unrecognized token %v", token)
	}
}

func (p *html_parser) apply_in_table_insertion_mode_rules(token html_token) {
	clear_stack_back_to_table_context := func() {
		for !slices.ContainsFunc([]string{"table", "template", "html"}, p.get_current_node().IsHtmlElement) {
			p.pop_node_from_soe()
		}
	}

	if _, ok := token.(*html_char_token); ok && slices.ContainsFunc([]string{
		"table", "tbody", "template", "tfoot", "thead", "tr",
	}, p.get_current_node().IsHtmlElement) {
		p.pending_table_char_tokens = []html_char_token{}
		p.original_insertion_mode = p.insertion_mode
		p.insertion_mode = html_parser_insertion_mode_in_table_text
		p.apply_in_table_text_insertion_mode_rules(token)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "caption" {
		clear_stack_back_to_table_context()
		p.list_of_active_formatting_elems = append(p.list_of_active_formatting_elems, html_active_formatting_elem_marker)
		p.insert_html_element(*tk)
		p.insertion_mode = html_parser_insertion_mode_in_caption
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "colgroup" {
		clear_stack_back_to_table_context()
		p.insert_html_element(*tk)
		p.insertion_mode = html_parser_insertion_mode_in_column_group
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "col" {
		clear_stack_back_to_table_context()
		p.insert_html_element(tagToken{tag_name: "colgroup"})
		p.insertion_mode = html_parser_insertion_mode_in_column_group
		p.apply_in_column_group_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"tbody", "tfoot", "thead"}, tk.tag_name) {
		clear_stack_back_to_table_context()
		p.insertion_mode = html_parser_insertion_mode_in_table_body
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"td", "th", "tr"}, tk.tag_name) {
		clear_stack_back_to_table_context()
		p.insert_html_element(tagToken{tag_name: "tbody"})
		p.apply_in_table_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "table" {
		p.parse_error_encountered(token)
		if !p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("table") }) {
			return
		} else {
			for {
				popped_elem := p.pop_node_from_soe()
				if popped_elem.IsHtmlElement("table") {
					break
				}
			}
			p.reset_insertion_mode_appropriately()
			p.apply_current_insertion_mode_rules(token)
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "table" {
		if !p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("table") }) {
			p.parse_error_encountered(token)
			return
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement("table") {
				break
			}
		}
		p.reset_insertion_mode_appropriately()
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{
		"body", "caption", "col", "colgroup", "html", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && slices.Contains([]string{"style", "script", "template"}, tk.tag_name)) ||
		(tk.is_end_tag() && tk.tag_name == "template") {
		p.apply_in_head_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "input" &&
		func() bool {
			if attr, ok := tk.Attr("type"); ok && cm.ToAsciiLowercase(attr) == "hidden" {
				return true
			}
			return false
		}() {
		p.parse_error_encountered(token)
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		tk.self_closing_acknowledged = true
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "form" {
		if p.soe_has_elem("template") {
			node := p.form_element_pointer
			p.form_element_pointer = nil
			if cm.IsNil(node) || !p.have_element_in_scope(func(n dom.Element) bool { return n == node }) {
				p.parse_error_encountered(token)
				return
			}
			p.generate_implied_end_tags(nil)
			if p.get_current_node() != node {
				p.parse_error_encountered(token)
			}
			remove_idx := slices.Index(p.stack_of_open_elems, node)
			p.remove_from_soe(remove_idx)
		} else {
			p.parse_error_encountered(token)
			if p.have_element_in_scope(func(n dom.Element) bool { return n.IsHtmlElement("form") }) ||
				!cm.IsNil(p.form_element_pointer) {
				return
			}
			p.insert_html_element(*tk)
			p.pop_node_from_soe()
		}
	} else if _, ok := token.(*html_eof_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else {
		p.parse_error_encountered(token)
		p.enable_foster_parenting = true
		p.apply_in_body_insertion_mode_rules(token)
		p.enable_foster_parenting = false
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intabletext
func (p *html_parser) apply_in_table_text_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.value == 0x0000 {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*html_char_token); ok {
		p.pending_table_char_tokens = append(p.pending_table_char_tokens, *tk)
	} else {
		if slices.ContainsFunc(p.pending_table_char_tokens, func(t html_char_token) bool { return !cm.IsAsciiWhitespace(t.value) }) {
			p.parse_error_encountered(token)
			// Below do the same thing as "else" in "in table" insertion mode.
			p.enable_foster_parenting = true
			for _, tk := range p.pending_table_char_tokens {
				p.apply_in_body_insertion_mode_rules(tk)
			}
			p.enable_foster_parenting = false
		} else {
			for _, tk := range p.pending_table_char_tokens {
				p.insert_character(tk.value)
			}
		}
		p.insertion_mode = p.original_insertion_mode
		p.apply_current_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incaption
func (p *html_parser) apply_in_caption_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "caption" {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement("caption") }) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(nil)
		if !p.get_current_node().IsHtmlElement("caption") {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement("caption") {
				break
			}
		}
		p.clear_list_of_active_formatting_elems_up_to_last_marker()
		p.insertion_mode = html_parser_insertion_mode_in_table
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && slices.Contains([]string{"caption", "col", "colgroup", "tbody", "td", "tfoot", "th", "thead", "tr"}, tk.tag_name)) ||
		(tk.is_end_tag() && tk.tag_name == "table") {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement("caption") }) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(nil)
		if !p.get_current_node().IsHtmlElement("caption") {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement("caption") {
				break
			}
		}
		p.clear_list_of_active_formatting_elems_up_to_last_marker()
		p.insertion_mode = html_parser_insertion_mode_in_table
		p.apply_in_table_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && (tk.is_end_tag() && slices.Contains([]string{
		"body", "col", "colgroup", "html", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tag_name)) {
		p.parse_error_encountered(token)
		return
	} else {
		p.apply_in_body_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incolgroup
func (p *html_parser) apply_in_column_group_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.insert_character(tk.value)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "col" {
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		if tk.is_self_closing {
			tk.self_closing_acknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "colgroup" {
		if !p.get_current_node().IsHtmlElement("colgroup") {
			p.parse_error_encountered(token)
			return
		}
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_table
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "col" {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && tk.tag_name == "template") ||
		(tk.is_end_tag() && tk.tag_name == "template") {
		p.apply_in_head_insertion_mode_rules(token)
	} else if _, ok := token.(*html_eof_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else {
		if !p.get_current_node().IsHtmlElement("colgroup") {
			p.parse_error_encountered(token)
			return
		}
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_table
		p.apply_in_table_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intbody
func (p *html_parser) apply_in_table_body_insertion_mode_rules(token html_token) {
	clear_stack_back_to_table_body_context := func() {
		for slices.ContainsFunc([]string{"tbody", "tfoot", "thead", "template", "html"}, p.get_current_node().IsHtmlElement) {
			p.pop_node_from_soe()
		}
	}
	if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "tr" {
		clear_stack_back_to_table_body_context()
		p.insert_html_element(*tk)
		p.insertion_mode = html_parser_insertion_mode_in_row
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"th", "td"}, tk.tag_name) {
		p.parse_error_encountered(token)
		clear_stack_back_to_table_body_context()
		p.insert_html_element(tagToken{tag_name: "tr"})
		p.insertion_mode = html_parser_insertion_mode_in_row
		p.apply_in_row_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{"tbody", "tfoot", "thead"}, tk.tag_name) {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) }) {
			p.parse_error_encountered(token)
			return
		}
		clear_stack_back_to_table_body_context()
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_table
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && slices.Contains([]string{"caption", "col", "colgroup", "tbody", "tfoot", "thead"}, tk.tag_name)) ||
		(tk.is_end_tag() && tk.tag_name == "table") {
		if !p.have_element_in_table_scope(func(n dom.Element) bool {
			return slices.ContainsFunc([]string{"tbody", "thead", "tfoot"}, n.IsHtmlElement)
		}) {
			p.parse_error_encountered(token)
			return
		}
		clear_stack_back_to_table_body_context()
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_table
	} else if tk.is_end_tag() && slices.Contains([]string{"body", "caption", "col", "colgroup", "html", "td", "th", "tr"}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else {
		p.apply_in_table_insertion_mode_rules(token)
	}

}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intr
func (p *html_parser) apply_in_row_insertion_mode_rules(token html_token) {
	clear_stack_back_to_table_row_context := func() {
		for slices.ContainsFunc([]string{"tr", "template", "html"}, p.get_current_node().IsHtmlElement) {
			p.pop_node_from_soe()
		}
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intr]")
	}
	if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"th", "td"}, tk.tag_name) {
		clear_stack_back_to_table_row_context()
		p.insert_html_element(*tk)
		p.insertion_mode = html_parser_insertion_mode_in_cell
		p.list_of_active_formatting_elems = append(p.list_of_active_formatting_elems, html_active_formatting_elem_marker)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "tr" {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement("tr") }) {
			p.parse_error_encountered(token)
			return
		}
		clear_stack_back_to_table_row_context()
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_table_body
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && slices.Contains([]string{"caption", "col", "colgroup", "tbody", "tfoot", "thead", "tr"}, tk.tag_name)) ||
		(tk.is_end_tag() && tk.tag_name == "table") {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement("tr") }) {
			p.parse_error_encountered(token)
			return
		}
		clear_stack_back_to_table_row_context()
		p.pop_node_from_soe()
		p.insertion_mode = html_parser_insertion_mode_in_table_body
		p.apply_in_table_body_insertion_mode_rules(token)
	} else if tk.is_end_tag() && slices.Contains([]string{"tbody", "tfoot", "thead"}, tk.tag_name) {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) }) {
			p.parse_error_encountered(token)
			return
		}
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement("tr") }) {
			return
		} else {
			clear_stack_back_to_table_row_context()
			p.pop_node_from_soe()
			p.insertion_mode = html_parser_insertion_mode_in_table_body
		}
	} else if tk.is_end_tag() && slices.Contains([]string{"body", "caption", "col", "colgroup", "html", "td", "th"}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else {
		p.apply_in_table_insertion_mode_rules(token)
	}
}

func (p *html_parser) apply_in_cell_insertion_mode_rules(token html_token) {
	close_cell := func() {
		p.generate_implied_end_tags(nil)
		if !slices.ContainsFunc([]string{"td", "th"}, p.get_current_node().IsHtmlElement) {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if slices.ContainsFunc([]string{"td", "th"}, popped_elem.IsHtmlElement) {
				break
			}
		}
		p.clear_list_of_active_formatting_elems_up_to_last_marker()
		p.insertion_mode = html_parser_insertion_mode_in_row
	}

	if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && slices.Contains([]string{"th", "td"}, tk.tag_name) {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) }) {
			p.parse_error_encountered(token)
			return
		}
		p.generate_implied_end_tags(nil)
		if !p.get_current_node().IsHtmlElement(tk.tag_name) {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement(tk.tag_name) {
				break
			}
		}
		p.clear_list_of_active_formatting_elems_up_to_last_marker()
		p.insertion_mode = html_parser_insertion_mode_in_row
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{
		"caption", "col", "colgroup", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tag_name) {
		if !p.have_element_in_table_scope(func(n dom.Element) bool {
			return slices.ContainsFunc([]string{"td", "th"}, n.IsHtmlElement)
		}) {
			panic("we should have td or th in SOE at this point")
		}
		close_cell()
		p.apply_in_row_insertion_mode_rules(token)
	} else if tk.is_end_tag() && slices.Contains([]string{"body", "caption", "col", "colgroup", "html"}, tk.tag_name) {
		p.parse_error_encountered(token)
		return
	} else if tk.is_end_tag() && slices.Contains([]string{"table", "tbody", "tfoot", "thead", "tr"}, tk.tag_name) {
		if !p.have_element_in_table_scope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tag_name) }) {
			p.parse_error_encountered(token)
			return
		}
		close_cell()
		p.apply_in_row_insertion_mode_rules(token)
	} else {
		p.apply_in_body_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intemplate
func (p *html_parser) apply_in_template_insertion_mode_rules(token html_token) {
	if _, ok := token.(*html_char_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else if _, ok := token.(*html_comment_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.is_start_tag() && slices.Contains([]string{
			"base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title",
		}, tk.tag_name)) ||
		(tk.is_end_tag() && tk.tag_name == "template") {
		p.apply_in_head_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{
		"caption", "colgroup", "tbody", "tfoot", "thead",
	}, tk.tag_name) {
		p.pop_node_from_sot()
		p.push_node_to_sot(html_parser_insertion_mode_in_table)
		p.insertion_mode = html_parser_insertion_mode_in_table
		p.apply_in_table_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "col" {
		p.pop_node_from_sot()
		p.push_node_to_sot(html_parser_insertion_mode_in_column_group)
		p.insertion_mode = html_parser_insertion_mode_in_column_group
		p.apply_in_column_group_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "tr" {
		p.pop_node_from_sot()
		p.push_node_to_sot(html_parser_insertion_mode_in_table_body)
		p.insertion_mode = html_parser_insertion_mode_in_table_body
		p.apply_in_table_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && slices.Contains([]string{"td", "th"}, tk.tag_name) {
		p.pop_node_from_sot()
		p.push_node_to_sot(html_parser_insertion_mode_in_row)
		p.insertion_mode = html_parser_insertion_mode_in_row
		p.apply_in_row_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() {
		p.pop_node_from_sot()
		p.push_node_to_sot(html_parser_insertion_mode_in_body)
		p.insertion_mode = html_parser_insertion_mode_in_body
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() {
		p.parse_error_encountered(token)
		return
	} else if _, ok := token.(*html_eof_token); ok {
		if !p.soe_has_elem("template") {
			p.stop_parsing()
		} else {
			p.parse_error_encountered(token)
		}
		for {
			popped_elem := p.pop_node_from_soe()
			if popped_elem.IsHtmlElement("template") {
				break
			}
		}
		p.clear_list_of_active_formatting_elems_up_to_last_marker()
		p.pop_node_from_sot()
		p.reset_insertion_mode_appropriately()
		p.apply_current_insertion_mode_rules(token)
	} else {
		log.Printf("[in template insertion mode] Unrecognized token %v", token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-afterbody
func (p *html_parser) apply_after_body_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, &html_parser_insertion_location{
			parent_node: p.stack_of_open_elems[0],
			tp:          html_adjusted_insertion_location_type_after_last_child,
		})
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "html" {
		if p.is_fragment_parsing {
			p.parse_error_encountered(token)
			return
		}
		p.insertion_mode = html_parser_insertion_mode_after_after_body
	} else if _, ok := token.(*html_eof_token); ok {
		p.stop_parsing()
	} else {
		p.parse_error_encountered(token)
		p.insertion_mode = html_parser_insertion_mode_in_body
		p.apply_in_body_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inframeset
func (p *html_parser) apply_in_frameset_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.insert_character(tk.value)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "frameset" {
		p.insert_html_element(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "framesets" {
		if p.get_current_node().IsHtmlElement("html") &&
			cm.IsNil(p.get_current_node().Parent()) {
			// current node is root html node
			p.parse_error_encountered(token)
			return
		}
		p.pop_node_from_soe()
		if !p.is_fragment_parsing && !p.get_current_node().IsHtmlElement("frameset") {
			p.insertion_mode = html_parser_insertion_mode_after_frameset
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "frame" {
		p.insert_html_element(*tk)
		p.pop_node_from_soe()
		if tk.is_self_closing {
			tk.self_closing_acknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "noframes" {
		p.apply_in_head_insertion_mode_rules(token)
	} else if _, ok := token.(*html_eof_token); ok {
		if !p.get_current_node().IsHtmlElement("html") ||
			!cm.IsNil(p.get_current_node().Parent()) {
			// current node is NOT root html node
			p.parse_error_encountered(token)
		}
		p.stop_parsing()
	} else {
		p.parse_error_encountered(token)
		return
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-afterframeset
func (p *html_parser) apply_after_frameset_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.insert_character(tk.value)
	} else if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, nil)
	} else if _, ok := token.(*html_doctype_token); ok {
		p.parse_error_encountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_end_tag() && tk.tag_name == "html" {
		p.insertion_mode = html_parser_insertion_mode_after_after_frameset
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "noframes" {
		p.apply_in_head_insertion_mode_rules(token)
	} else if _, ok := token.(*html_eof_token); ok {
		p.stop_parsing()
	} else {
		p.parse_error_encountered(token)
		return
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-after-after-body-insertion-mode
func (p *html_parser) apply_after_after_body_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, &html_parser_insertion_location{
			parent_node: p.document,
			tp:          html_adjusted_insertion_location_type_after_last_child,
		})
	} else if _, ok := token.(*html_doctype_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if _, ok := token.(*html_eof_token); ok {
		p.stop_parsing()
	} else {
		p.parse_error_encountered(token)
		p.insertion_mode = html_parser_insertion_mode_in_body
		p.apply_in_body_insertion_mode_rules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-after-after-frameset-insertion-mode
func (p *html_parser) apply_after_after_frameset_insertion_mode_rules(token html_token) {
	if tk, ok := token.(*html_comment_token); ok {
		p.insert_comment(tk.data, &html_parser_insertion_location{
			parent_node: p.document,
			tp:          html_adjusted_insertion_location_type_after_last_child,
		})
	} else if _, ok := token.(*html_doctype_token); ok {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*html_char_token); ok && tk.is_char_token_with_one_of("\t\n\u000c\r ") {
		p.apply_in_body_insertion_mode_rules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "html" {
		p.apply_in_body_insertion_mode_rules(token)
	} else if _, ok := token.(*html_eof_token); ok {
		p.stop_parsing()
	} else if tk, ok := token.(*tagToken); ok && tk.is_start_tag() && tk.tag_name == "noframes" {
		p.apply_in_head_insertion_mode_rules(token)
	} else {
		p.parse_error_encountered(token)
		return
	}
}

func (p *html_parser) apply_current_insertion_mode_rules(token html_token) {
	insertion_mode_funcs := map[html_parser_insertion_mode]func(token html_token){
		html_parser_insertion_mode_initial:              p.apply_initial_insertion_mode_rules,
		html_parser_insertion_mode_before_html:          p.apply_before_html_insertion_mode_rules,
		html_parser_insertion_mode_before_head:          p.apply_before_head_insertion_mode_rules,
		html_parser_insertion_mode_in_head:              p.apply_in_head_insertion_mode_rules,
		html_parser_insertion_mode_in_head_noscript:     p.apply_in_head_noscript_insertion_mode_rules,
		html_parser_insertion_mode_after_head:           p.apply_after_head_insertion_mode_rules,
		html_parser_insertion_mode_in_body:              p.apply_in_body_insertion_mode_rules,
		html_parser_insertion_mode_text:                 p.apply_text_insertion_mode_rules,
		html_parser_insertion_mode_in_table:             p.apply_in_table_insertion_mode_rules,
		html_parser_insertion_mode_in_table_text:        p.apply_in_table_text_insertion_mode_rules,
		html_parser_insertion_mode_in_caption:           p.apply_in_caption_insertion_mode_rules,
		html_parser_insertion_mode_in_column_group:      p.apply_in_column_group_insertion_mode_rules,
		html_parser_insertion_mode_in_table_body:        p.apply_in_table_body_insertion_mode_rules,
		html_parser_insertion_mode_in_row:               p.apply_in_row_insertion_mode_rules,
		html_parser_insertion_mode_in_cell:              p.apply_in_cell_insertion_mode_rules,
		html_parser_insertion_mode_in_template:          p.apply_in_template_insertion_mode_rules,
		html_parser_insertion_mode_after_body:           p.apply_after_body_insertion_mode_rules,
		html_parser_insertion_mode_in_frameset:          p.apply_in_frameset_insertion_mode_rules,
		html_parser_insertion_mode_after_frameset:       p.apply_after_frameset_insertion_mode_rules,
		html_parser_insertion_mode_after_after_body:     p.apply_after_after_body_insertion_mode_rules,
		html_parser_insertion_mode_after_after_frameset: p.apply_after_after_frameset_insertion_mode_rules,
	}
	insertion_mode_funcs[p.insertion_mode](token)
}

// https://html.spec.whatwg.org/multipage/parsing.html#stop-parsing
func (p *html_parser) stop_parsing() {
	p.run_parser = false
	// TODO
}

var html_parser_mathml_attr_adjust_map = map[string]string{
	"definitionurl": "definitionURL",
}

// https://html.spec.whatwg.org/multipage/parsing.html#adjust-mathml-attributes
func html_parser_adjust_mathml_attrs(token *tagToken) {
	for i, attr := range token.attrs {
		if new_name, ok := html_parser_mathml_attr_adjust_map[attr.LocalName]; ok {
			token.attrs[i].LocalName = new_name
		}
	}
}

var parserSvgAttrAdjustMap = map[string]string{
	"attributename":       "attributeName",
	"attributetype":       "attributeType",
	"basefrequency":       "baseFrequency",
	"baseprofile":         "baseProfile",
	"calcmode":            "calcMode",
	"clippathunits":       "clipPathUnits",
	"diffuseconstant":     "diffuseConstant",
	"edgemode":            "edgeMode",
	"filterunits":         "filterUnits",
	"glyphref":            "glyphRef",
	"gradienttransform":   "gradientTransform",
	"gradientunits":       "gradientUnits",
	"kernelmatrix":        "kernelMatrix",
	"kernelunitlength":    "kernelUnitLength",
	"keypoints":           "keyPoints",
	"keysplines":          "keySplines",
	"keytimes":            "keyTimes",
	"lengthadjust":        "lengthAdjust",
	"limitingconeangle":   "limitingConeAngle",
	"markerheight":        "markerHeight",
	"markerunits":         "markerUnits",
	"markerwidth":         "markerWidth",
	"maskcontentunits":    "maskContentUnits",
	"maskunits":           "maskUnits",
	"numoctaves":          "numOctaves",
	"pathlength":          "pathLength",
	"patterncontentunits": "patternContentUnits",
	"patterntransform":    "patternTransform",
	"patternunits":        "patternUnits",
	"pointsatx":           "pointsAtX",
	"pointsaty":           "pointsAtY",
	"pointsatz":           "pointsAtZ",
	"preservealpha":       "preserveAlpha",
	"preserveaspectratio": "preserveAspectRatio",
	"primitiveunits":      "primitiveUnits",
	"refx":                "refX",
	"refy":                "refY",
	"repeatcount":         "repeatCount",
	"repeatdur":           "repeatDur",
	"requiredextensions":  "requiredExtensions",
	"requiredfeatures":    "requiredFeatures",
	"specularconstant":    "specularConstant",
	"specularexponent":    "specularExponent",
	"spreadmethod":        "spreadMethod",
	"startoffset":         "startOffset",
	"stddeviation":        "stdDeviation",
	"stitchtiles":         "stitchTiles",
	"surfacescale":        "surfaceScale",
	"systemlanguage":      "systemLanguage",
	"tablevalues":         "tableValues",
	"targetx":             "targetX",
	"targety":             "targetY",
	"textlength":          "textLength",
	"viewbox":             "viewBox",
	"viewtarget":          "viewTarget",
	"xchannelselector":    "xChannelSelector",
	"ychannelselector":    "yChannelSelector",
	"zoomandpan":          "zoomAndPan",
}

// https://html.spec.whatwg.org/multipage/parsing.html#adjust-svg-attributes
func adjustSvgAttrs(token *tagToken) {
	for i, attr := range token.attrs {
		if newName, ok := parserSvgAttrAdjustMap[attr.LocalName]; ok {
			token.attrs[i].LocalName = newName
		}
	}
}

var parserForeignAttrAdjustMap = map[string]struct {
	prefix    *string
	localName string
	namespace namespaces.Namespace
}{
	"xlink:actuate": {cm.MakeStrPtr("xlink"), "actuate", namespaces.Xlink},
	"xlink:arcrole": {cm.MakeStrPtr("xlink"), "arcrole", namespaces.Xlink},
	"xlink:href":    {cm.MakeStrPtr("xlink"), "href", namespaces.Xlink},
	"xlink:role":    {cm.MakeStrPtr("xlink"), "role", namespaces.Xlink},
	"xlink:show":    {cm.MakeStrPtr("xlink"), "show", namespaces.Xlink},
	"xlink:title":   {cm.MakeStrPtr("xlink"), "title", namespaces.Xlink},
	"xlink:type":    {cm.MakeStrPtr("xlink"), "type", namespaces.Xlink},
	"xml:lang":      {cm.MakeStrPtr("xml"), "lang", namespaces.Xml},
	"xml:space":     {cm.MakeStrPtr("xml"), "space", namespaces.Xml},
	"xmlns":         {nil, "xmlns", namespaces.Xmlns},
	"xmlns:xlink":   {cm.MakeStrPtr("xmlns"), "xlink", namespaces.Xmlns},
}

// https://html.spec.whatwg.org/multipage/parsing.html#adjust-foreign-attributes
func parserAdjustForeignAttrs(token *tagToken) {
	for i, attr := range token.attrs {
		if newAttrData, ok := parserForeignAttrAdjustMap[attr.LocalName]; ok {
			ns := newAttrData.namespace
			token.attrs[i].NamespacePrefix = newAttrData.prefix
			token.attrs[i].LocalName = newAttrData.localName
			token.attrs[i].Namespace = &ns
		}
	}
}
