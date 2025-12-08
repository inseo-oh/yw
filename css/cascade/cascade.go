// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// Package cascade provide CSS cascading logic based on [CSS Cascading and Inheritance Level 4].
//
// [CSS Cascading and Inheritance Level 4]: https://www.w3.org/TR/css-cascade-4
package cascade

import (
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/selector"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

// ApplyStyleRules collects all relevant style rules from stylesheets associated
// with docOrSr, together with uaStylesheet, and calculates computed value of
// each descendant element of the docOrSr.
//
// Resulting style is saved to each element's ComputedStyleSet.
func ApplyStyleRules(uaStylesheet *cssom.Stylesheet, docOrSr dom.Node) {
	// https://www.w3.org/TR/css-cascade-4/#cascade-origin

	const (
		priorityTransition = iota // Highest priority
		priorityImportantUserAgent
		priorityImportantUser
		priorityImportantAuthor
		priorityAnimation
		priorityNormalAuthor
		priorityNormalUser
		priorityNormalUserAgent // Lowest priority
	)

	type declEntry struct {
		rule cssom.StyleRule
		decl cssom.Declaration
	}
	declGroups := [][]declEntry{
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
	addDecl := func(group *[]declEntry, rule cssom.StyleRule, decl cssom.Declaration) {
		*group = append(*group, declEntry{rule, decl})
	}
	elems := []dom.Element{}
	for _, n := range dom.InclusiveDescendants(docOrSr) {
		if elem, ok := n.(dom.Element); ok {
			elems = append(elems, elem)
		}
	}

	// User agent declarations -------------------------------------------------
	for _, rule := range uaStylesheet.StyleRules {
		for _, decl := range rule.Declarations {
			if decl.IsImportant {
				addDecl(&declGroups[priorityImportantUserAgent], rule, decl)
			} else {
				addDecl(&declGroups[priorityNormalUserAgent], rule, decl)
			}
		}
	}
	// Author declarations -----------------------------------------------------
	for _, sheet := range cssom.DocumentOrShadowRootDataOf(docOrSr).Stylesheets {
		for _, rule := range sheet.StyleRules {
			for _, decl := range rule.Declarations {
				if decl.IsImportant {
					addDecl(&declGroups[priorityImportantAuthor], rule, decl)
				} else {
					addDecl(&declGroups[priorityNormalAuthor], rule, decl)
				}
			}
		}
	}
	// Presentional hints ------------------------------------------------------
	for _, elem := range elems {
		cbs := elem.Callbacks()
		if cbs.PresentationalHints != nil {
			rules := cbs.PresentationalHints().([]cssom.StyleRule)
			for _, rule := range rules {
				for _, decl := range rule.Declarations {
					if decl.IsImportant {
						addDecl(&declGroups[priorityImportantAuthor], rule, decl)
					} else {
						addDecl(&declGroups[priorityNormalAuthor], rule, decl)
					}
				}
			}
		}
	}

	// Apply specificity -------------------------------------------------------
	// TODO

	// Now we apply rules ------------------------------------------------------
	for i := len(declGroups) - 1; 0 <= i; i-- {
		declGroup := declGroups[i]
		for _, declEntry := range declGroup {
			rule := declEntry.rule
			selectedElements := selector.MatchAgainstTree(rule.SelectorList, []dom.Node{docOrSr})
			for _, node := range selectedElements {
				elem := node.(dom.Element)
				declEntry.decl.ApplyStyleRules(elem)
			}
		}
	}
	// Inherit missing values from parent --------------------------------------
	for _, elem := range elems {
		if parent := elem.Parent(); !util.IsNil(parent) {
			if parentElem, ok := parent.(dom.Element); ok {
				cssom.ElementDataOf(elem).ComputedStyleSet.InheritPropertiesFromParent(cssom.ComputedStyleSetSourceOf(parentElem))
			}
		}
	}
}
