// Implementation of Cascading and Inheritance Level 4 (https://www.w3.org/TR/css-cascade-4)
package libhtml

import cm "github.com/inseo-oh/yw/util"

func css_apply_style_rules(ua_stylesheet *css_stylesheet, doc_or_sr dom_DocumentOrShadowRoot) {
	// https://www.w3.org/TR/css-cascade-4/#cascade-origin

	const (
		priority_transition = iota
		priority_important_user_agent
		priority_important_user
		priority_important_author
		priority_animation
		priority_normal_author
		priority_normal_user
		priority_normal_user_agent
	)

	type decl_entry struct {
		rule css_style_rule
		decl css_declaration
	}
	decl_groups := [][]decl_entry{
		// Higher priority first, lower priority last

		{}, // Transition declarations
		{}, // Important user agent declarations
		{}, // Important user declarations
		{}, // Important author declarations
		{}, // Animation declarations
		{}, // Normal author declarations
		{}, // Normal user declarations
		{}, // Normal user agent declarations
	}
	add_decl := func(group *[]decl_entry, rule css_style_rule, decl css_declaration) {
		*group = append(*group, decl_entry{rule, decl})
	}
	elems := []dom_Element{}
	for _, n := range dom_node_inclusive_descendants(doc_or_sr) {
		if elem, ok := n.(dom_Element); ok {
			elems = append(elems, elem)
		}
	}

	// User agent declarations -------------------------------------------------
	for _, rule := range ua_stylesheet.style_rules {
		for _, decl := range rule.declarations {
			if decl.is_important {
				add_decl(&decl_groups[priority_important_user_agent], rule, decl)
			} else {
				add_decl(&decl_groups[priority_normal_user_agent], rule, decl)
			}
		}
	}
	// Author declarations -----------------------------------------------------
	for _, sheet := range doc_or_sr.get_css_stylesheets() {
		for _, rule := range sheet.style_rules {
			for _, decl := range rule.declarations {
				if decl.is_important {
					add_decl(&decl_groups[priority_important_author], rule, decl)
				} else {
					add_decl(&decl_groups[priority_normal_author], rule, decl)
				}
			}
		}
	}
	// Presentional hints ------------------------------------------------------
	for _, elem := range elems {
		cbs := elem.get_callbacks()
		if cbs.get_presentational_hints != nil {
			rules := cbs.get_presentational_hints()
			for _, rule := range rules {
				for _, decl := range rule.declarations {
					if decl.is_important {
						add_decl(&decl_groups[priority_important_author], rule, decl)
					} else {
						add_decl(&decl_groups[priority_normal_author], rule, decl)
					}
				}
			}
		}
	}

	// Apply specificity -------------------------------------------------------
	// TODO

	// Now we apply rules ------------------------------------------------------
	for i := len(decl_groups) - 1; 0 <= i; i-- {
		decl_group := decl_groups[i]
		for _, decl_entry := range decl_group {
			rule := decl_entry.rule
			selected_elements := css_match_selector_against_tree(rule.selector_list, []dom_Node{doc_or_sr})
			for _, node := range selected_elements {
				elem := node.(dom_Element)
				decl_entry.decl.apply_style_rules(elem)
			}
		}
	}
	// Inherit missing values from parent --------------------------------------
	for _, elem := range elems {
		if parent := elem.get_parent(); !cm.IsNil(parent) {
			if parent_elem, ok := parent.(dom_Element); ok {
				elem.get_computed_style_set().inherit_properties_from_parent(parent_elem)
			}
		}
	}
}
