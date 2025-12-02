package cssom

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/inseo-oh/yw/dom"
)

// Stylesheet represents a CSS stylesheet
type Stylesheet struct {
	Type                     string       // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-type
	Location                 *string      // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-location
	ParentStylesheet         *Stylesheet  // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-parent-css-style-sheet
	OwnerNode                dom.Node     // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-owner-node
	OwnerRule                *StyleRule   // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-owner-css-rule
	Media                    any          // [STUB] https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-media
	Title                    string       // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-title
	AlternateFlag            bool         // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-alternate-flag
	DisabledFlag             bool         // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-disabled-flag
	StyleRules               []StyleRule  // (STUB) https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-css-rules
	OriginCleanFlag          bool         // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-origin-clean-flag
	ConstructedFlag          bool         // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-constructed-flag
	DisallowModificationFlag bool         // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-disallow-modification-flag
	ConstructorDocument      dom.Document // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-constructor-document
	StylesheetBaseURL        *url.URL     // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#concept-css-style-sheet-stylesheet-base-url
}

// Dump prints CSS stylesheet to the standard logger.
func (sheet Stylesheet) Dump() {
	for i, rule := range sheet.StyleRules {
		selectorListStr := strings.Builder{}
		for i, s := range rule.SelectorList {
			if i != 0 {
				selectorListStr.WriteString(", ")
			}
			selectorListStr.WriteString(fmt.Sprintf("%v", s))
		}
		log.Printf("style-rule[%d](%s) {", i, selectorListStr.String())
		log.Printf("	declarations {")
		for _, decl := range rule.Declarations {
			log.Printf("        %s : %v", decl.Name, decl.Value)
		}
		log.Printf("    }")
		log.Printf("	at-rules {")
		for _, rule := range rule.AtRules {
			log.Printf("		   <name>: %s", rule.Name)
			log.Printf("		<prelude>: %s", rule.Prelude)
			log.Printf("		  <value>: %s", rule.Value)
		}
		log.Printf("    }")
		log.Printf("}")
	}
}

var (
	preferredStylesheetSetName         = ""  // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#preferred-css-style-sheet-set-name
	lastStylesheetSetName      *string = nil // https://www.w3.org/TR/2021/WD-cssom-1-20210826/#last-css-style-sheet-set-name
)

// StylesheetSet is set of [CSS stylesheet sets].
//
// [CSS stylesheet sets]: https://www.w3.org/TR/2021/WD-cssom-1-20210826/#css-style-sheet-set
type StylesheetSet []*Stylesheet

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#css-style-sheet-set-name
func (s StylesheetSet) name() string {
	return s[0].Title
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#persistent-css-style-sheet
func persistentStylesheets(domRoot dom.Node) []*Stylesheet {
	out := []*Stylesheet{}
	sheets := DocumentOrShadowRootDataOf(domRoot).Stylesheets
	for i, sheet := range DocumentOrShadowRootDataOf(domRoot).Stylesheets {
		title := sheet.Title
		if title != "" || sheet.AlternateFlag {
			continue
		}
		out = append(out, sheets[i])
	}
	return out
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#css-style-sheet-set
func stylesheetSets(domRoot dom.Node) []StylesheetSet {
	sets := []StylesheetSet{}
	sheets := DocumentOrShadowRootDataOf(domRoot).Stylesheets
	for i, sheet := range DocumentOrShadowRootDataOf(domRoot).Stylesheets {
		title := sheet.Title
		if title == "" {
			continue
		}
		setIndex := slices.IndexFunc(sets, func(set StylesheetSet) bool { return set[0].Title == title })
		if setIndex != -1 {
			sets[setIndex] = append(sets[setIndex], sheets[i])
		} else {
			sets = append(sets, []*Stylesheet{sheets[i]})
		}
	}
	return sets
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#enabled-css-style-sheet-set
func enabledStylesheetSets(domRoot dom.Node) []StylesheetSet {
	sets := []StylesheetSet{}
	for _, set := range stylesheetSets(domRoot) {
		sheets := []*Stylesheet{}
		for _, sheet := range set {
			if sheet.DisabledFlag {
				continue
			}
			sheets = append(sheets, sheet)
		}
		sets = append(sets, sheets)
	}
	return sets
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#change-the-preferred-css-style-sheet-set-name
func changePreferredStylesheetSetName(domRoot dom.Node, name string) {
	current := preferredStylesheetSetName
	preferredStylesheetSetName = name
	if name != current && lastStylesheetSetName == nil {
		enableStylesheetSet(domRoot, name)
	}
}

// https://www.w3.org/TR/2021/WD-cssom-1-20210826/#enable-a-css-style-sheet-set
func enableStylesheetSet(domRoot dom.Node, name string) {
	for _, set := range enabledStylesheetSets(domRoot) {
		for _, sheet := range set {
			sheet.DisabledFlag = true
		}
	}
	if name == "" {
		return
	}
	for _, set := range enabledStylesheetSets(domRoot) {
		for _, sheet := range set {
			sheet.DisabledFlag = set.name() != name
		}
	}
}

// AddStylesheet adds the sheet to sheet's OwnerNode.
//
// Spec: https://www.w3.org/TR/2021/WD-cssom-1-20210826/#add-a-css-style-sheet
func AddStylesheet(sheet *Stylesheet) {
	domRoot := dom.Root(sheet.OwnerNode)
	_ = domRoot

	// S1.
	DocumentOrShadowRootDataOf(domRoot).Stylesheets = append(DocumentOrShadowRootDataOf(domRoot).Stylesheets, sheet)
	// S2.
	if sheet.DisabledFlag {
		return
	}
	// S3.
	if sheet.Title != "" && !sheet.AlternateFlag && preferredStylesheetSetName == "" {
		changePreferredStylesheetSetName(domRoot, sheet.Title)
	}
	// S4.
	if sheet.Title == "" ||
		(lastStylesheetSetName == nil && sheet.Title == preferredStylesheetSetName) ||
		(lastStylesheetSetName != nil && sheet.Title == *lastStylesheetSetName) {
		sheet.DisabledFlag = false
		return
	}
	// S5
	sheet.DisabledFlag = true
}

// RemoveStylesheet removes the sheet from sheet's OwnerNode.
//
// Spec: https://www.w3.org/TR/2021/WD-cssom-1-20210826/#remove-a-css-style-sheet
//
// BUG(ois): RemoveStylesheet will cause panic if sheet is not actually part of the OwnerNode.
func RemoveStylesheet(sheet *Stylesheet) {
	domRoot := dom.Root(sheet.OwnerNode)

	sheets := DocumentOrShadowRootDataOf(domRoot).Stylesheets
	idx := slices.Index(sheets, sheet)
	DocumentOrShadowRootDataOf(domRoot).Stylesheets = append(sheets[:idx], sheets[idx+1:]...)
	sheet.ParentStylesheet = nil
	sheet.OwnerNode = nil
	sheet.OwnerRule = nil
}

// AssociatedStylesheet returns stylesheets associated with given element.
//
// Spec: https://drafts.csswg.org/cssom/#associated-css-style-sheet
//
// (This is part of Editor's Draft, but HTML spec needs it and it's easy to implement)
func AssociatedStylesheet(node dom.Element) *Stylesheet {
	domRoot := dom.Root(node)

	for _, set := range enabledStylesheetSets(domRoot) {
		for _, sheet := range set {
			if sheet.OwnerNode == node {
				return sheet
			}
		}
	}
	return nil
}
