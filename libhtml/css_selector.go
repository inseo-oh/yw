// Implementation of the CSS Selector Module Level 4 (https://www.w3.org/TR/2022/WD-selectors-4-20221111/)
package libhtml

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	cm "yw/libcommon"
)

type css_selector interface {
	String() string
	equals(other css_selector) bool
	match_against_element(element dom_Element) bool
}

type css_selector_ns_prefix struct{ ident string }

func (sel css_selector_ns_prefix) equals(other css_selector_ns_prefix) bool {
	// STUB
	return true
}

// Returns nil if not found
func (ts *css_token_stream) parse_selector_ns_prefix() *css_selector_ns_prefix {
	// STUB
	return nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-wq-name
type css_selector_wq_name struct {
	ns_prefix *css_selector_ns_prefix // May be nil
	ident     string
}

func (wq_name css_selector_wq_name) String() string {
	if wq_name.ns_prefix != nil {
		return fmt.Sprintf("%v%s", wq_name.ns_prefix, wq_name.ident)
	} else {
		return wq_name.ident
	}
}
func (sel css_selector_wq_name) equals(other css_selector_wq_name) bool {
	if (sel.ns_prefix != nil) != (other.ns_prefix != nil) {
		return false
	} else if (sel.ns_prefix != nil) && !sel.ns_prefix.equals(*other.ns_prefix) {
		return false
	}
	if sel.ident != other.ident {
		return false
	}
	return true
}

// Returns nil if not found
func (ts *css_token_stream) parse_selector_wq_name() *css_selector_wq_name {
	old_cursor := ts.cursor
	ns_prefix := ts.parse_selector_ns_prefix()
	var ident_tk css_ident_token
	if temp := ts.consume_token_with_type(css_token_type_ident); !cm.IsNil(temp) {
		ident_tk = temp.(css_ident_token)
	} else {
		ts.cursor = old_cursor
		return nil
	}
	return &css_selector_wq_name{ns_prefix, ident_tk.value}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#ref-for-typedef-type-selector
type css_type_selector struct {
	type_name css_selector_wq_name
}

func (sel css_type_selector) String() string { return fmt.Sprintf("%v", sel.type_name) }
func (sel css_type_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_type_selector); !ok {
		return false
	} else {
		if !sel.type_name.equals(other_sel.type_name) {
			return false
		}
	}
	return true
}
func (sel css_type_selector) match_against_element(element dom_Element) bool {
	// TODO: Handle namespace
	name := sel.type_name.ident
	if element.get_local_name() != name {
		return false
	}
	return true
}

type css_wildcard_selector struct {
	ns_prefix *css_selector_ns_prefix
}

func (sel css_wildcard_selector) String() string { return fmt.Sprintf("%v*", sel.ns_prefix) }
func (sel css_wildcard_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_wildcard_selector); !ok {
		return false
	} else {
		if (sel.ns_prefix == nil) != (other_sel.ns_prefix == nil) {
			return false
		}
		if sel.ns_prefix != nil && !sel.ns_prefix.equals(*other_sel.ns_prefix) {
			return false
		}
	}
	return true
}
func (sel css_wildcard_selector) match_against_element(element dom_Element) bool {
	return true
}

// Returns nil if not found
func (ts *css_token_stream) parse_type_selector() css_selector {
	old_cursor := ts.cursor
	if type_name := ts.parse_selector_wq_name(); type_name != nil {
		// <wq-name>
		return css_type_selector{*type_name}
	} else {
		// <ns-prefix?> *
		ns_prefix := ts.parse_selector_ns_prefix()
		if tk := ts.consume_token_with_type(css_token_type_delim); !cm.IsNil(tk) {
			if tk.(css_delim_token).value != '*' {
				ts.cursor = old_cursor
				return nil
			}
			return css_wildcard_selector{ns_prefix}
		} else {
			ts.cursor = old_cursor
			return nil
		}
	}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-id-selector
type css_id_selector struct{ id string }

func (sel css_id_selector) String() string { return fmt.Sprintf("#%s", sel.id) }
func (sel css_id_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_id_selector); !ok {
		return false
	} else {
		if sel.id != other_sel.id {
			return false
		}
	}
	return true
}
func (sel css_id_selector) match_against_element(element dom_Element) bool {
	id := sel.id
	elem_id, ok := element.get_attribute_without_namespace("id")
	if !ok || elem_id != id {
		return false
	}
	return true
}

func (ts *css_token_stream) parse_id_selector() *css_id_selector {
	old_cursor := ts.cursor
	var hash_tk css_hash_token
	if temp := ts.consume_token_with_type(css_token_type_hash); !cm.IsNil(temp) {
		hash_tk = temp.(css_hash_token)
	} else {
		ts.cursor = old_cursor
		return nil
	}
	return &css_id_selector{hash_tk.value}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-class-selector
type css_class_selector struct{ class string }

func (sel css_class_selector) String() string { return fmt.Sprintf(".%s", sel.class) }
func (sel css_class_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_class_selector); !ok {
		return false
	} else {
		if sel.class != other_sel.class {
			return false
		}
	}
	return true
}
func (sel css_class_selector) match_against_element(element dom_Element) bool {
	class := sel.class
	classes, ok := element.get_attribute_without_namespace("class")
	if !ok {
		return false
	}
	class_list := strings.Split(classes, " ")
	return slices.Contains(class_list, class)
}

func (ts *css_token_stream) parse_class_selector() (*css_class_selector, error) {
	old_cursor := ts.cursor
	if cm.IsNil(ts.consume_delim_token_with('.')) {
		ts.cursor = old_cursor
		return nil, nil
	}
	var ident_tk css_ident_token
	if temp := ts.consume_token_with_type(css_token_type_ident); !cm.IsNil(temp) {
		ident_tk = temp.(css_ident_token)
	} else {
		return nil, errors.New("expected identifier after '.'")
	}
	return &css_class_selector{ident_tk.value}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attribute-selector
type css_attribute_selector struct {
	attr_name css_selector_wq_name
	matcher   css_attribute_matcher
	// Below are valid only if the matcher isn't 'none'
	attr_value        string
	is_case_sensitive bool
}

type css_attribute_matcher uint8

const (
	css_attribute_matcher_none     = css_attribute_matcher(iota)
	css_attribute_matcher_normal   // =
	css_attribute_matcher_tilde    // ~=
	css_attribute_matcher_bar      // |=
	css_attribute_matcher_caret    // ^=
	css_attribute_matcher_dollar   // $=
	css_attribute_matcher_asterisk // *=
)

func (sel css_attribute_selector) String() string {
	flag_str := "s"
	if !sel.is_case_sensitive {
		flag_str = "i"
	}
	switch sel.matcher {
	case css_attribute_matcher_none:
		return fmt.Sprintf("[%s]", sel.attr_name)
	case css_attribute_matcher_normal:
		return fmt.Sprintf("[%s=%s %s]", sel.attr_name, sel.attr_value, flag_str)
	case css_attribute_matcher_tilde:
		return fmt.Sprintf("[%s~=%s %s]", sel.attr_name, sel.attr_value, flag_str)
	case css_attribute_matcher_bar:
		return fmt.Sprintf("[%s|=%s %s]", sel.attr_name, sel.attr_value, flag_str)
	case css_attribute_matcher_caret:
		return fmt.Sprintf("[%s^=%s %s]", sel.attr_name, sel.attr_value, flag_str)
	case css_attribute_matcher_dollar:
		return fmt.Sprintf("[%s$=%s %s]", sel.attr_name, sel.attr_value, flag_str)
	case css_attribute_matcher_asterisk:
		return fmt.Sprintf("[%s*=%s %s]", sel.attr_name, sel.attr_value, flag_str)
	}
	return fmt.Sprintf("[%s<unknown matcher %d>%s %s]", sel.attr_name, sel.matcher, sel.attr_value, flag_str)
}
func (sel css_attribute_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_attribute_selector); !ok {
		return false
	} else {
		if !sel.attr_name.equals(other_sel.attr_name) {
			return false
		}
		if sel.matcher != other_sel.matcher {
			return false
		}
		if sel.matcher != css_attribute_matcher_none {
			if sel.attr_value != other_sel.attr_value {
				return false
			}
			if sel.is_case_sensitive != other_sel.is_case_sensitive {
				return false
			}
		}
	}
	return true
}
func (sel css_attribute_selector) match_against_element(element dom_Element) bool {
	// STUB
	return false
}

func (ts *css_token_stream) parse_attribute_selector() (*css_attribute_selector, error) {
	old_cursor := ts.cursor
	blk := ts.consume_ast_simple_block_with(css_ast_simple_block_type_square)
	if blk == nil {
		ts.cursor = old_cursor
		return nil, nil
	}

	body_stream := css_token_stream{tokens: blk.body}
	// [<  >attr  ] ------------------------------------------------------------
	// [<  >attr  =  value  modifier  ] ----------------------------------------
	body_stream.skip_whitespaces()
	// [  <attr>  ] ------------------------------------------------------------
	// [  <attr>  =  value  modifier  ] ----------------------------------------
	wq_name := body_stream.parse_selector_wq_name()
	if wq_name == nil {
		return nil, errors.New("expected name after '['")
	}
	// [  attr<  >] ------------------------------------------------------------
	// [  attr<  >=  value  modifier  ] ----------------------------------------
	body_stream.skip_whitespaces()
	if !body_stream.is_end() {
		// [  attr  <=>  value  modifier  ] ------------------------------------
		var matcher css_attribute_matcher
		// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attr-matcher
		if !cm.IsNil(body_stream.consume_delim_token_with('~')) {
			matcher = css_attribute_matcher_tilde
		} else if !cm.IsNil(body_stream.consume_delim_token_with('|')) {
			matcher = css_attribute_matcher_bar
		} else if !cm.IsNil(body_stream.consume_delim_token_with('^')) {
			matcher = css_attribute_matcher_caret
		} else if !cm.IsNil(body_stream.consume_delim_token_with('$')) {
			matcher = css_attribute_matcher_dollar
		} else if !cm.IsNil(body_stream.consume_delim_token_with('*')) {
			matcher = css_attribute_matcher_asterisk
		} else {
			matcher = css_attribute_matcher_normal
		}
		if cm.IsNil(body_stream.consume_delim_token_with('=')) {
			return nil, errors.New("expected operator after the attribute name")
		}
		// [  attr  =<  >value  modifier  ] ------------------------------------
		body_stream.skip_whitespaces()
		// [  attr  =  <value>  modifier  ] ------------------------------------
		var attr_value string
		if n := body_stream.consume_token_with_type(css_token_type_ident); !cm.IsNil(n) {
			attr_value = n.(css_ident_token).value
		} else if n := body_stream.consume_token_with_type(css_token_type_string); !cm.IsNil(n) {
			attr_value = n.(css_string_token).value
		} else {
			return nil, errors.New("expected attribute value after the operator")
		}
		// [  attr  =  value<  >modifier  ] ------------------------------------
		body_stream.skip_whitespaces()
		// [  attr  =  value  <modifier>  ] ------------------------------------
		is_case_sensitive := true
		if !cm.IsNil(body_stream.consume_ident_token_with("i")) {
			is_case_sensitive = false
		} else if !cm.IsNil(body_stream.consume_ident_token_with("s")) {
			is_case_sensitive = true
		}
		// [  attr  =  value  modifier<  >] ------------------------------------
		body_stream.skip_whitespaces()
		if !body_stream.is_end() {
			return nil, errors.New("found junk after contents of the attribute selector")
		}
		return &css_attribute_selector{*wq_name, matcher, attr_value, is_case_sensitive}, nil
	}
	return &css_attribute_selector{*wq_name, css_attribute_matcher_none, "", false}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-class-selector
type css_pseudo_class_selector struct {
	name string
	args []css_token
}

func (sel css_pseudo_class_selector) String() string {
	if len(sel.args) != 0 {
		// TODO: Display arguments in better way
		return fmt.Sprintf(":%s(%v)", sel.name, sel.args)
	} else {
		return fmt.Sprintf(":%s", sel.name)
	}
}
func (sel css_pseudo_class_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_pseudo_class_selector); !ok {
		return false
	} else {
		if sel.name != other_sel.name {
			return false
		}
		if len(sel.args) != len(other_sel.args) {
			return false
		}
		// TODO: Compare actual arguments
	}
	return true
}
func (sel css_pseudo_class_selector) match_against_element(element dom_Element) bool {
	// STUB
	return false
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-class-selector
func (ts *css_token_stream) parse_pseudo_class_selector() (*css_pseudo_class_selector, error) {
	old_cursor := ts.cursor

	// <:>name ----------------------------------------------------------------
	// <:>func(value) ----------------------------------------------------------
	if cm.IsNil(ts.consume_token_with_type(css_token_type_colon)) {
		ts.cursor = old_cursor
		return nil, nil
	}
	if ident_tk := ts.consume_token_with_type(css_token_type_ident); !cm.IsNil(ident_tk) {
		// :<name> ------------------------------------------------------------
		name := ident_tk.(css_ident_token).value
		return &css_pseudo_class_selector{name, nil}, nil
	} else if func_node := ts.consume_token_with_type(css_token_type_ast_function); !cm.IsNil(func_node) {
		// :<func(value)> ------------------------------------------------------
		name := func_node.(css_ast_function_token).name
		sub_stream := css_token_stream{tokens: func_node.(css_ast_function_token).value}
		args := sub_stream.consume_any_value()
		if args == nil {
			ts.cursor = old_cursor
			return nil, errors.New("expected value after '('")
		}
		if !sub_stream.is_end() {
			return nil, errors.New("unexpected junk after arguments")
		}
		return &css_pseudo_class_selector{name, args}, nil
	} else {
		ts.cursor = old_cursor
		return nil, nil
	}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-element-selector
func (ts *css_token_stream) parse_pseudo_element_selector() (*css_pseudo_class_selector, error) {
	old_cursor := ts.cursor
	if cm.IsNil(ts.consume_token_with_type(css_token_type_colon)) {
		ts.cursor = old_cursor
		return nil, nil
	}
	if temp, err := ts.parse_pseudo_class_selector(); temp != nil {
		return temp, nil
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-subclass-selector
//
// Returns nil if not found
func (ts *css_token_stream) parse_subclass_selector() (css_selector, error) {
	if sel := ts.parse_id_selector(); sel != nil {
		return *sel, nil
	}

	if sel, err := ts.parse_class_selector(); sel != nil {
		return *sel, nil
	} else if err != nil {
		return nil, err
	}

	if sel, err := ts.parse_attribute_selector(); sel != nil {
		return *sel, nil
	} else if err != nil {
		return nil, err
	}

	if sel, err := ts.parse_pseudo_class_selector(); sel != nil {
		return *sel, nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-compound-selector
type css_compound_selector struct {
	type_selector     css_selector // May be nil
	subclass_selector []css_selector
	pseudo_items      []css_compound_selector_pseudo_item
}
type css_compound_selector_pseudo_item struct {
	element_selector css_pseudo_class_selector
	class_selector   []css_pseudo_class_selector
}

func (sel css_compound_selector) String() string {
	sb := strings.Builder{}
	if sel.type_selector != nil {
		sb.WriteString(fmt.Sprintf("%v", sel.type_selector))
	}
	for _, v := range sel.subclass_selector {
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	for _, p := range sel.pseudo_items {
		sb.WriteString(fmt.Sprintf("%v", p.element_selector))
		for _, v := range p.class_selector {
			sb.WriteString(fmt.Sprintf("%v", v))
		}
	}
	return sb.String()
}
func (sel css_compound_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_compound_selector); !ok {
		return false
	} else {
		if (sel.type_selector != nil) != (other_sel.type_selector != nil) {
			return false
		} else if (sel.type_selector != nil) && !sel.type_selector.equals(other_sel.type_selector) {
			return false
		}
		if len(sel.subclass_selector) != len(other_sel.subclass_selector) {
			return false
		}
		for i := 0; i < len(sel.subclass_selector); i++ {
			if !sel.subclass_selector[i].equals(other_sel.subclass_selector[i]) {
				return false
			}
		}
		if len(sel.pseudo_items) != len(other_sel.pseudo_items) {
			return false
		}
		for i := 0; i < len(sel.pseudo_items); i++ {
			if !sel.pseudo_items[i].element_selector.equals(other_sel.pseudo_items[i].element_selector) {
				return false
			}
			if len(sel.pseudo_items[i].class_selector) != len(other_sel.pseudo_items[i].class_selector) {
				return false
			}
			for j := 0; j < len(sel.pseudo_items[i].class_selector); j++ {
				if !sel.pseudo_items[i].class_selector[j].equals(other_sel.pseudo_items[i].class_selector[j]) {
					return false
				}
			}
		}
	}
	return true
}
func (sel css_compound_selector) match_against_element(element dom_Element) bool {
	if sel.type_selector != nil && !sel.type_selector.match_against_element(element) {
		return false
	}
	for _, ss := range sel.subclass_selector {
		if !ss.match_against_element(element) {
			return false
		}
	}
	if len(sel.pseudo_items) != 0 {
		for i := len(sel.pseudo_items) - 1; 0 < i; i-- {
			// TODO
		}
	}
	return true
}

// Returns nil if not found
func (ts *css_token_stream) parse_compound_selector() (*css_compound_selector, error) {
	old_cursor := ts.cursor
	type_sel := ts.parse_type_selector()
	subclass_sels := []css_selector{}
	pseudo_items := []css_compound_selector_pseudo_item{}
	for {
		subclass_sel, err := ts.parse_subclass_selector()
		if cm.IsNil(subclass_sel) {
			if err != nil {
				return nil, err
			}
			break
		}
		subclass_sels = append(subclass_sels, subclass_sel)
	}

	for {
		elem_sel, err := ts.parse_pseudo_element_selector()
		if elem_sel == nil {
			if err != nil {
				return nil, err
			}
			break
		}
		class_sels := []css_pseudo_class_selector{}
		for {
			class_sel, err := ts.parse_pseudo_class_selector()
			if cm.IsNil(class_sel) {
				if err != nil {
					ts.cursor = old_cursor
					return nil, err
				}
				break
			}
			class_sels = append(class_sels, *class_sel)
		}
		pseudo_items = append(pseudo_items, css_compound_selector_pseudo_item{*elem_sel, class_sels})

	}

	if type_sel == nil && len(subclass_sels) == 0 && len(pseudo_items) == 0 {
		ts.cursor = old_cursor
		return nil, nil
	}
	return &css_compound_selector{type_sel, subclass_sels, pseudo_items}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-compound-selector-list
//
// Returns nil if not found
func (ts *css_token_stream) parse_compound_selector_list() ([]*css_compound_selector, error) {
	return css_accept_comma_separated_repetion(ts, 0, func(ts *css_token_stream) (*css_compound_selector, error) {
		return ts.parse_compound_selector()
	})
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector
type css_complex_selector struct {
	base css_compound_selector
	rest []css_complex_selector_rest
}
type css_complex_selector_rest struct {
	combinator css_combinator
	selector   css_compound_selector
}
type css_combinator uint8

const (
	css_combinator_child        = css_combinator(iota)
	css_combinator_direct_child // >
	css_combinator_plus         // +
	css_combinator_tilde        // ~
	css_combinator_two_bars     // ||
)

func (sel css_complex_selector) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v", sel.base))
	for _, v := range sel.rest {
		switch v.combinator {
		case css_combinator_child:
			sb.WriteString(" ")
		case css_combinator_direct_child:
			sb.WriteString(">")
		case css_combinator_plus:
			sb.WriteString("+")
		case css_combinator_tilde:
			sb.WriteString("~")
		case css_combinator_two_bars:
			sb.WriteString("||")
		}
		sb.WriteString(fmt.Sprintf("%v", v.selector))
	}
	return sb.String()
}
func (sel css_complex_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_complex_selector); !ok {
		return false
	} else {
		if !sel.base.equals(other_sel.base) {
			return false
		}
		for i := 0; i < len(sel.rest); i++ {
			if sel.rest[i].combinator != other_sel.rest[i].combinator {
				return false
			}
			if !sel.rest[i].selector.equals(other_sel.rest[i].selector) {
				return false
			}
		}
	}
	return true
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-complex-selector-against-an-element
func (s css_complex_selector) match_against_element(element dom_Element) bool {
	// Test each compound selector, from right to left
	for i := len(s.rest) - 1; 0 < i; i-- {
		prev_sel := s.base
		if i != 0 {
			prev_sel = s.rest[i-1].selector
		}

		sel := s.rest[i].selector
		if sel.match_against_element(element) {
			return false
		}
		switch s.rest[i].combinator {
		case css_combinator_child:
			// A B
			curr_elem := element.get_parent()
			found := false
			for !cm.IsNil(curr_elem) {
				if _, ok := (curr_elem).(dom_Element); !ok {
					break
				}
				if prev_sel.match_against_element(curr_elem.(dom_Element)) {
					found = true
					break
				}
				curr_elem = curr_elem.get_parent()
			}
			if !found {
				return false
			}
		case css_combinator_direct_child:
			// A > B
			if cm.IsNil(element.get_parent()) {
				return false
			}
			if parent, ok := element.get_parent().(dom_Element); !ok {
				if !prev_sel.match_against_element(parent) {
					return false
				}
			} else {
				return false
			}
		case css_combinator_plus, css_combinator_tilde, css_combinator_two_bars:
			panic("TODO")
		default:
			log.Printf("BUG: unrecognized css_combinator %d while parsing selector: %v", s.rest[i].combinator, s)
			continue
		}
	}
	return s.base.match_against_element(element)
}

// Special internal CSS selector that matches DOM node pointer directly.
type css_node_ptr_selector struct {
	element dom_Element
}

func (sel css_node_ptr_selector) equals(other css_selector) bool {
	if other_sel, ok := other.(css_node_ptr_selector); !ok {
		return false
	} else {
		return other_sel.element == sel.element
	}
}
func (sel css_node_ptr_selector) match_against_element(element dom_Element) bool {
	return sel.element == element
}
func (sel css_node_ptr_selector) String() string {
	return sel.element.String()
}

func (ts *css_token_stream) parse_complex_selector() (*css_complex_selector, error) {
	old_cursor := ts.cursor
	base, err := ts.parse_compound_selector()
	if base == nil {
		ts.cursor = old_cursor
		return nil, err
	}
	rest := []css_complex_selector_rest{}
	for {
		comb := css_combinator_child
		if !cm.IsNil(ts.consume_delim_token_with('>')) {
			comb = css_combinator_direct_child
		} else if !cm.IsNil(ts.consume_delim_token_with('+')) {
			comb = css_combinator_plus
		} else if !cm.IsNil(ts.consume_delim_token_with('~')) {
			comb = css_combinator_tilde
		} else if !cm.IsNil(ts.consume_delim_token_with('|')) {
			if !cm.IsNil(ts.consume_delim_token_with('|')) {
				comb = css_combinator_two_bars
			} else {
				ts.cursor -= 2
			}
		} else if !cm.IsNil(ts.consume_token_with_type(css_token_type_whitespace)) {
			ts.skip_whitespaces()
			comb = css_combinator_child
		}
		another_unit, err := ts.parse_compound_selector()
		if cm.IsNil(another_unit) {
			if err != nil {
				return nil, err
			}
			break
		}
		rest = append(rest, css_complex_selector_rest{comb, *another_unit})
	}
	if base == nil && len(rest) == 0 {
		ts.cursor = old_cursor
		return nil, nil
	}
	return &css_complex_selector{*base, rest}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector-list
//
// Returns nil if not found
func (ts *css_token_stream) parse_complex_selector_list() ([]css_selector, error) {
	sel_list, err := css_accept_comma_separated_repetion(ts, 0, func(ts *css_token_stream) (*css_complex_selector, error) {
		return ts.parse_complex_selector()
	})
	if sel_list == nil {
		return nil, err
	}
	out := []css_selector{}
	for _, s := range sel_list {
		out = append(out, *s)
	}
	return out, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-selector-list
func (ts *css_token_stream) parse_selector_list() ([]css_selector, error) {
	return ts.parse_complex_selector_list()
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#parse-a-selector
func css_parse_selector(src string) ([]css_selector, error) {
	tokens, err := css_tokenize(src)
	if tokens == nil && err != nil {
		return nil, err
	}
	return css_parse(tokens, func(ts *css_token_stream) ([]css_selector, error) {
		return ts.parse_selector_list()
	})
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-selector-against-an-element
func css_match_selector_against_element(selector []css_selector, element dom_Element) bool {
	for _, s := range selector {
		if s.match_against_element(element) {
			return true
		}
	}
	return false
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#match-a-selector-against-a-tree
func css_match_selector_against_tree(selector []css_selector, roots []dom_Node) []dom_Node {
	selector_match_list := []dom_Node{}
	for _, root := range roots {
		candiate_elems := []dom_Node{}
		for _, n := range dom_node_inclusive_descendants(root) {
			if _, ok := n.(dom_Element); ok {
				candiate_elems = append(candiate_elems, n)
			}
		}
		for _, n := range candiate_elems {
			if css_match_selector_against_element(selector, n.(dom_Element)) {
				selector_match_list = append(selector_match_list, n)
			}
			// TODO: Pseudo element
		}

	}
	return selector_match_list
}
