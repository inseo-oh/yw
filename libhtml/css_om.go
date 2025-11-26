// Implementation of the CSS Object Model (https://www.w3.org/TR/2021/WD-cssom-1-20210826/)
package libhtml

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"
)

type css_stylesheet struct {
	tp                         string           // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-type
	location                   *string          // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-location
	parent_stylesheet          *css_stylesheet  // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-parent-css-style-sheet
	owner_node                 html_HTMLElement // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-owner-node
	owner_rule                 *css_style_rule  // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-owner-css-rule
	media                      any              // [STUB] https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-media
	title                      string           // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-title
	alternate_flag             bool             // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-alternate-flag
	disabled_flag              bool             // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-disabled-flag
	style_rules                []css_style_rule // (STUB) https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-css-rules
	origin_clean_flag          bool             // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-origin-clean-flag
	constructed_flag           bool             // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-constructed-flag
	disallow_modification_flag bool             // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-disallow-modification-flag
	constructor_document       dom_Document     // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-constructor-document
	stylesheet_base_url        *url.URL         // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-stylesheet-base-url
}

func (sheet css_stylesheet) dump() {
	for i, rule := range sheet.style_rules {
		selector_list_str := strings.Builder{}
		for i, s := range rule.selector_list {
			if i != 0 {
				selector_list_str.WriteString(", ")
			}
			selector_list_str.WriteString(fmt.Sprintf("%v", s))
		}
		log.Printf("style-rule[%d](%s) {", i, selector_list_str.String())
		log.Printf("	declarations {")
		for _, decl := range rule.declarations {
			log.Printf("        %s : %v", decl.name, decl.value)
		}
		log.Printf("    }")
		log.Printf("	at-rules {")
		for _, rule := range rule.at_rules {
			log.Printf("		   <name>: %s", rule.name)
			log.Printf("		<prelude>: %s", rule.prelude)
			log.Printf("		  <value>: %s", rule.value)
		}
		log.Printf("    }")
		log.Printf("}")
	}
}

var (
	css_preferred_stylesheet_set_name         = ""  // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#preferred-css-style-sheet-set-name
	css_last_stylesheet_set_name      *string = nil // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#last-css-style-sheet-set-name
)

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#css-style-sheet-set
type css_stylesheet_set []*css_stylesheet

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#css-style-sheet-set-name
func (s css_stylesheet_set) name() string {
	return s[0].title
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#persistent-css-style-sheet
func css_persistent_stylesheets(doc_or_sr dom_DocumentOrShadowRoot) []*css_stylesheet {
	out := []*css_stylesheet{}
	sheets := doc_or_sr.get_css_stylesheets()
	for i, sheet := range doc_or_sr.get_css_stylesheets() {
		title := sheet.title
		if title != "" || sheet.alternate_flag {
			continue
		}
		out = append(out, sheets[i])
	}
	return out
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#css-style-sheet-set
func css_stylesheet_sets(doc_or_sr dom_DocumentOrShadowRoot) []css_stylesheet_set {
	sets := []css_stylesheet_set{}
	sheets := doc_or_sr.get_css_stylesheets()
	for i, sheet := range doc_or_sr.get_css_stylesheets() {
		title := sheet.title
		if title == "" {
			continue
		}
		set_index := slices.IndexFunc(sets, func(set css_stylesheet_set) bool { return set[0].title == title })
		if set_index != -1 {
			sets[set_index] = append(sets[set_index], sheets[i])
		} else {
			sets = append(sets, []*css_stylesheet{sheets[i]})
		}
	}
	return sets
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#enabled-css-style-sheet-set
func css_enabled_stylesheet_sets(doc_or_sr dom_DocumentOrShadowRoot) []css_stylesheet_set {
	sets := []css_stylesheet_set{}
	for _, set := range css_stylesheet_sets(doc_or_sr) {
		sheets := []*css_stylesheet{}
		for _, sheet := range set {
			if sheet.disabled_flag {
				continue
			}
			sheets = append(sheets, sheet)
		}
		sets = append(sets, sheets)
	}
	return sets
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#change-the-preferred-css-style-sheet-set-name
func css_change_preferred_stylesheet_set_name(doc_or_sr dom_DocumentOrShadowRoot, name string) {
	current := css_preferred_stylesheet_set_name
	css_preferred_stylesheet_set_name = name
	if name != current && css_last_stylesheet_set_name == nil {
		css_enable_stylesheet_set(doc_or_sr, name)
	}
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#enable-a-css-style-sheet-set
func css_enable_stylesheet_set(doc_or_sr dom_DocumentOrShadowRoot, name string) {
	for _, set := range css_enabled_stylesheet_sets(doc_or_sr) {
		for _, sheet := range set {
			sheet.disabled_flag = true
		}
	}
	if name == "" {
		return
	}
	for _, set := range css_enabled_stylesheet_sets(doc_or_sr) {
		for _, sheet := range set {
			sheet.disabled_flag = set.name() != name
		}
	}
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#add-a-css-style-sheet
func css_add_stylesheet(sheet *css_stylesheet) {
	doc_or_sr := dom_node_root(sheet.owner_node).(dom_DocumentOrShadowRoot)
	_ = doc_or_sr

	// S1.
	doc_or_sr.set_css_stylesheets(append(doc_or_sr.get_css_stylesheets(), sheet))
	// S2.
	if sheet.disabled_flag {
		return
	}
	// S3.
	if sheet.title != "" && !sheet.alternate_flag && css_preferred_stylesheet_set_name == "" {
		css_change_preferred_stylesheet_set_name(doc_or_sr, sheet.title)
	}
	// S4.
	if sheet.title == "" ||
		(css_last_stylesheet_set_name == nil && sheet.title == css_preferred_stylesheet_set_name) ||
		(css_last_stylesheet_set_name != nil && sheet.title == *css_last_stylesheet_set_name) {
		sheet.disabled_flag = false
		return
	}
	// S5
	sheet.disabled_flag = true
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#remove-a-css-style-sheet
func css_remove_stylesheet(sheet *css_stylesheet) {
	doc_or_sr := dom_node_root(sheet.owner_node).(dom_DocumentOrShadowRoot)

	sheets := doc_or_sr.get_css_stylesheets()
	idx := slices.Index(sheets, sheet)
	doc_or_sr.set_css_stylesheets(append(sheets[:idx], sheets[idx+1:]...))
	sheet.parent_stylesheet = nil
	sheet.owner_node = nil
	sheet.owner_rule = nil
}

// https://drafts.csswg.org/cssom/#associated-css-style-sheet
// This is part of Editor's Draft, but HTML spec needs it and it's easy to implement :D
func css_associated_stylesheet(node dom_Element) *css_stylesheet {
	doc_or_sr := dom_node_root(node).(dom_DocumentOrShadowRoot)

	for _, set := range css_enabled_stylesheet_sets(doc_or_sr) {
		for _, sheet := range set {
			if sheet.owner_node == node {
				return sheet
			}
		}
	}
	return nil
}

func css_apply_style_rules(doc_or_sr dom_DocumentOrShadowRoot) {
	rules := []css_style_rule{}
	for _, sheet := range doc_or_sr.get_css_stylesheets() {
		rules = append(rules, sheet.style_rules...)
	}
	// Apply styles from HTML tags first
	elems := []dom_Element{}
	for _, n := range dom_node_inclusive_descendants(doc_or_sr) {
		if elem, ok := n.(dom_Element); ok {
			elems = append(elems, elem)
		}
	}
	// Apply presentational hints
	for _, elem := range elems {
		cbs := elem.get_callbacks()
		if cbs.get_presentational_hints != nil {
			decls := cbs.get_presentational_hints()
			for _, decl := range decls {
				decl.apply_style_rules(elem)
			}
		}
	}

	// TODO: Apply specificity, !important, etc...
	for _, rule := range rules {
		rule.apply_style_rules([]dom_Node{doc_or_sr})
	}

}
