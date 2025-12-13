// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Package htmlparser provides HTML parser.
package htmlparser

//go:generate go run ./entities_gen

import (
	"log"
	"runtime"
	"slices"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/html/elements"
	"github.com/inseo-oh/yw/namespaces"
	"github.com/inseo-oh/yw/util"
)

type stackOfOpenElements []dom.Element

// Returns an item from stack of open elements.
// - Positive index starts from the top of the stack (first pushed item first).
// - Negative index starts from the bottom of the stack (most recent item first).
func (soe *stackOfOpenElements) nodeAt(idx int) dom.Element {
	if 0 < idx {
		return (*soe)[idx]
	} else if idx < 0 {
		return (*soe)[len(*soe)+idx]
	} else {
		panic("zero index is not allowed")
	}
}
func (soe *stackOfOpenElements) push(node dom.Element) {
	*soe = append((*soe), node)
}
func (soe *stackOfOpenElements) pop() dom.Element {
	// TODO: https://html.spec.whatwg.org/multipage/parsing.html#the-stack-of-open-elements
	// When the current node is removed from the stack of open elements, process internal resource links given the current node's node document.
	node := (*soe)[len(*soe)-1]
	*soe = (*soe)[:len(*soe)-1]
	cb := node.Callbacks().PoppedFromStackOfOpenElements
	if cb != nil {
		cb()
	}
	return node
}
func (soe *stackOfOpenElements) remove(idx int) {
	*soe = append((*soe)[:idx], (*soe)[idx+1:]...)
}
func (soe stackOfOpenElements) hasOneOfElems(elems []string) bool {
	return slices.ContainsFunc(soe, func(n dom.Element) bool {
		return !slices.ContainsFunc(elems, n.IsHtmlElement)
	})
}
func (soe stackOfOpenElements) hasElem(elem string) bool {
	return soe.hasOneOfElems([]string{elem})
}

type listOfActiveFormattingElements []activeFormattingElement
type activeFormattingElement struct {
	elem  dom.Element // If elem is nil, this is a marker.
	token tagToken
}

func (e activeFormattingElement) isMarker() bool { return util.IsNil(e.elem) }

var activeFormattingElemMarker activeFormattingElement

func (laf listOfActiveFormattingElements) lastMarker() (elem *activeFormattingElement, idx int) {
	idx = slices.IndexFunc(laf, activeFormattingElement.isMarker)
	if idx == -1 {
		return nil, -1
	} else {
		return &laf[idx], idx
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#push-onto-the-list-of-active-formatting-elements
func (laf *listOfActiveFormattingElements) push(elem dom.Element) {
	lastMarker, lastMarkerIdx := laf.lastMarker()
	checkFn := func(otherElem dom.Element) bool {
		if elem.LocalName() != otherElem.LocalName() {
			return false
		}
		elemNs, elemHasNs := elem.Namespace()
		otherNs, otherHasNs := otherElem.Namespace()
		if elemHasNs != otherHasNs {
			return false
		}
		if elemHasNs && elemNs != otherNs {
			return false
		}
		attrs := elem.Attrs()
		otherAttrs := otherElem.Attrs()
		if len(attrs) != len(otherAttrs) {
			return false
		}
		for i := range attrs {
			if attrs[i].LocalName() == otherAttrs[i].LocalName() &&
				attrs[i].Value() == otherAttrs[i].Value() {
				return true
			}
		}
		return false
	}
	matchingItemIndices := []int{}
	checkStartIdx := 0
	if lastMarker != nil {
		checkStartIdx = lastMarkerIdx + 1
	}
	for i := checkStartIdx; i < len(*laf); i++ {
		if checkFn((*laf)[i].elem) {
			matchingItemIndices = append(matchingItemIndices, i)
		}
	}
	if 3 <= len(matchingItemIndices) {
		*laf = append(
			(*laf)[:matchingItemIndices[0]],
			(*laf)[matchingItemIndices[0]+1:]...,
		)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#clear-the-list-of-active-formatting-elements-up-to-the-last-marker
func (laf *listOfActiveFormattingElements) clearUpToLastMarker() {
	for {
		lastEntry := (*laf)[len(*laf)-1]
		*laf = (*laf)[:len(*laf)-1]
		if lastEntry.isMarker() {
			break
		}
	}
}

type stackOfTemplateInsertionModes []insertionMode

func (sot *stackOfTemplateInsertionModes) push(mode insertionMode) {
	*sot = append(*sot, mode)
}
func (sot *stackOfTemplateInsertionModes) pop() insertionMode {
	// TODO: https://html.spec.whatwg.org/multipage/parsing.html#the-stack-of-open-elements
	// When the current node is removed from the stack of open elements, process internal resource links given the current node's node document.
	node := (*sot)[len(*sot)-1]
	*sot = (*sot)[:len(*sot)-1]
	return node
}

// Parser holds state of a HTML parser.
//
// Empty value won't do anything useful - Use [NewParser] to create one.
type Parser struct {
	tokenizer tokenizer

	Document dom.Document

	headElementPointer dom.Element
	formElementPointer dom.Element

	runParser                  bool
	isFramesetNotOk            bool
	isFragmentParsing          bool
	enableScripting            bool
	enableFosterParenting      bool
	hasActiveSpeculativeParser bool // We don't have speculative parsing support, so this is mostly just a placeholder, just in case decide to we support it later.

	insertionMode         insertionMode
	originalInsertionMode insertionMode

	stackOfOpenElements            stackOfOpenElements
	listOfActiveFormattingElements listOfActiveFormattingElements
	stackOfTemplateInsertionModes  stackOfTemplateInsertionModes

	onNextToken func(token htmlToken) parserControl

	pendingTableCharTokens []charToken // https://html.spec.whatwg.org/multipage/parsing.html#concept-pending-table-char-tokens
}

// NewParser creates new parser for given sourceCode.
func NewParser(sourceCode string) Parser {
	return Parser{tokenizer: newTokenizer(sourceCode)}
}

// Run runs the parser, and returns resulting [dom.Document].
func (p *Parser) Run() dom.Document {
	if p.Document == nil {
		p.Document = dom.NewDocument()
	}
	p.tokenizer.onTokenEmitted = func(tk htmlToken) {
		if p.onNextToken != nil {
			switch p.onNextToken(tk) {
			case parserControlIgnoreToken:
				return
			case parserControlContinue:
			default:
				panic("unknown result from onNextToken()")
			}
		}

		isStartTagToken := func() bool {
			if _, ok := (tk).(tagToken); ok {
				return true
			}
			return false
		}
		isStartTagTokenWith := func(name string) bool {
			if tk, ok := (tk).(tagToken); ok {
				return tk.isStartTag() && tk.tagName == name
			}
			return false
		}
		isCharToken := func() bool {
			if _, ok := (tk).(charToken); ok {
				return true
			}
			return false
		}
		isEofToken := func() bool {
			if _, ok := (tk).(eofToken); ok {
				return true
			}
			return false
		}

		// https://html.spec.whatwg.org/multipage/parsing.html#tree-construction-dispatcher
		if (len(p.stackOfOpenElements) == 0) ||
			(func() bool {
				n := p.adjustedCurrentNode()
				if ns, ok := n.Namespace(); ok && ns == namespaces.Html {
					return true
				}
				return false
			}()) ||
			(p.adjustedCurrentNode().IsMathmlTextIntegrationPoint() && !isStartTagTokenWith("mglyph") && !isStartTagTokenWith("malignmark")) ||
			(p.adjustedCurrentNode().IsMathmlTextIntegrationPoint() && isCharToken()) ||
			(p.adjustedCurrentNode().IsMathmlElement("annotation-xml") && isStartTagTokenWith("svg")) ||
			(p.adjustedCurrentNode().IsHtmlIntegrationPoint() && isStartTagToken()) ||
			(p.adjustedCurrentNode().IsHtmlIntegrationPoint() && isCharToken()) ||
			(p.adjustedCurrentNode().IsHtmlIntegrationPoint() && isEofToken()) {
			p.applyCurrentInsertionModeRules(tk)
		} else {
			// TODO: Process the token according to the rules given in the section for parsing tokens in foreign content.
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#tree-construction-dispatcher]")
		}
	}
	p.runParser = true
	for p.runParser {
		p.tokenizer.run()
	}
	return p.Document
}
func (p *Parser) parseErrorEncountered(tk htmlToken) {
	log.Println("Parse error occured near", tk)

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("-> From %s:%d (%s)", file, line, runtime.FuncForPC(pc).Name())
	}
}

type parserControl uint8

const (
	parserControlIgnoreToken parserControl = iota
	parserControlContinue
)

type insertionMode uint8

const (
	initialInsertionMode insertionMode = iota
	beforeHtmlInsertionMode
	beforeHeadInsertionMode
	inHeadInsertionMode
	inHeadNoscriptInsertionMode
	afterHeadInsertionMode
	inBodyInsertionMode
	textInsertionMode
	inTableInsertionMode
	inTableTextInsertionMode
	inCaptionInsertionMode
	inColumnGroupInsertionMode
	inTableBodyInsertionMode
	inRowInsertionMode
	inCellInsertionMode
	inTemplateInsertionMode
	afterBodyInsertionMode
	inFramesetInsertionMode
	afterFramesetInsertionMode
	afterAfterBodyInsertionMode
	afterAfterFramesetInsertionMode
)

// https://html.spec.whatwg.org/multipage/parsing.html#reconstruct-the-active-formatting-elements
func (p *Parser) reconstructActiveFormattingElems() {
	if len(p.listOfActiveFormattingElements) == 0 {
		return
	}
	lastEntry := p.listOfActiveFormattingElements[len(p.listOfActiveFormattingElements)-1]
	if lastEntry.isMarker() || slices.Contains(p.stackOfOpenElements, lastEntry.elem) {
		return
	}
	entryIdx := len(p.listOfActiveFormattingElements) - 1
	for {
		entry := func() *activeFormattingElement { return &p.listOfActiveFormattingElements[entryIdx] }
	rewind:
		if entryIdx == 0 {
			goto create
		}
		entryIdx = entryIdx - 1
		if !entry().isMarker() && slices.Contains(p.stackOfOpenElements, entry().elem) {
			goto rewind
		}
	advance:
		entryIdx = entryIdx + 1
	create:
		newElem := p.insertHtmlElement(entry().token)
		entry().elem = newElem
		if entryIdx != len(p.listOfActiveFormattingElements)-1 {
			goto advance
		}
	}
}

func (p *Parser) currentNode() dom.Element {
	return p.stackOfOpenElements.nodeAt(-1)
}
func (p *Parser) adjustedCurrentNode() dom.Element {
	if p.isFragmentParsing {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#adjusted-current-node]")
	} else {
		return p.currentNode()
	}
}

func (p *Parser) haveElementInSpecificScope(isTargetNode func(n dom.Element) bool, elemTypes []dom.NamePair) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-the-specific-scope
	nodeIdx := len(p.stackOfOpenElements) - 1
	for {
		node := p.stackOfOpenElements[nodeIdx]
		if isTargetNode(node) {
			return true
		}
		if slices.ContainsFunc(elemTypes, node.IsElement) {
			return false
		}
		nodeIdx--
	}
}
func (p *Parser) haveElementInScope(isTargetNode func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-scope
	return p.haveElementInSpecificScope(isTargetNode, []dom.NamePair{
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
func (p *Parser) haveElementInListItemScope(isTargetNode func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-list-item-scope
	return p.haveElementInSpecificScope(isTargetNode, []dom.NamePair{
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
func (p *Parser) haveElementInButtonScope(isTargetNode func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-button-scope
	return p.haveElementInSpecificScope(isTargetNode, []dom.NamePair{
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
func (p *Parser) haveElementInTableScope(isTargetNode func(n dom.Element) bool) bool {
	// https://html.spec.whatwg.org/multipage/parsing.html#has-an-element-in-table-scope
	return p.haveElementInSpecificScope(isTargetNode, []dom.NamePair{
		{Namespace: namespaces.Html, LocalName: "html"}, {Namespace: namespaces.Html, LocalName: "table"}, {Namespace: namespaces.Html, LocalName: "template"},
	})
}

type insertionLocation struct {
	parentNode dom.Node
	tp         insertionLocationType
}
type insertionLocationType uint8

const (
	insertionLocationAfterLastChild insertionLocationType = iota
)

// https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node
//
// overrideTarget may be nil pointer
func (p *Parser) appropriatePlaceForInsertionNode(overrideTarget dom.Element) insertionLocation {
	var res insertionLocation
	target := overrideTarget
	if util.IsNil(target) {
		target = p.currentNode()
	}

	if targetElem := target; p.enableFosterParenting && (targetElem.IsHtmlElement("table") ||
		targetElem.IsHtmlElement("tbody") ||
		targetElem.IsHtmlElement("tfoot") ||
		targetElem.IsHtmlElement("thead") ||
		targetElem.IsHtmlElement("tr")) {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node]")
	} else {
		res = insertionLocation{target, insertionLocationAfterLastChild}
	}
	if targetElem := target; targetElem.IsInside(dom.NamePair{Namespace: namespaces.Html, LocalName: "template"}) {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node]")
	}
	return res
}

// https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token
func (p *Parser) createElementForToken(token tagToken, namespace namespaces.Namespace, intendedParent dom.Node) dom.Element {
	if p.hasActiveSpeculativeParser {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	document := intendedParent.NodeDocument()
	localName := token.tagName
	isVal, hasIs := token.Attr("is")
	var is *string
	if hasIs {
		is = &isVal
	}
	registry := dom.LookupCustomElementRegistry(intendedParent)
	definition := registry.LookupCustomElementDefinition(&namespace, localName, is)
	willExecuteScript := false
	if definition != nil && !p.isFragmentParsing {
		willExecuteScript = true
	}
	if willExecuteScript {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	elem := dom.CreateElement(document, localName, &namespace, nil, is, willExecuteScript, registry, token, func(namespace *namespaces.Namespace, localName string) func(opt dom.ElementCreationCommonOptions) dom.Element {
		factoryFn := func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLElement(opt) }
		if namespace != nil && *namespace == namespaces.Html && localName == "html" {
			factoryFn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLHtmlElement(opt) }
		} else if namespace != nil && *namespace == namespaces.Html && localName == "body" {
			factoryFn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLBodyElement(opt) }
		} else if namespace != nil && *namespace == namespaces.Html && localName == "link" {
			factoryFn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLLinkElement(opt) }
		} else if namespace != nil && *namespace == namespaces.Html && localName == "style" {
			factoryFn = func(opt dom.ElementCreationCommonOptions) dom.Element { return elements.NewHTMLStyleElement(opt) }
		}
		return factoryFn
	})
	for _, attr := range token.attrs {
		elem.AppendAttr(attr)
	}
	if willExecuteScript {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	if attr, ok := elem.AttrWithNamespace(dom.NamePair{Namespace: namespaces.Xmlns, LocalName: "xmlns"}); ok {
		if ns, ok := elem.Namespace(); !ok || (attr != string(ns)) {
			p.parseErrorEncountered(token)
		}
	}
	if attr, ok := elem.AttrWithNamespace(dom.NamePair{Namespace: namespaces.Xmlns, LocalName: "xmlns:xlink"}); ok && attr != string(namespaces.Xlink) {
		p.parseErrorEncountered(token)
	}
	if elem.(elements.HTMLElement).IsFormResettableElement() && !elem.(elements.HTMLElement).IsFormAssociatedCustomElement() {
		// TODO: Invoke reset algorithm
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	hasAttr := func(name string) bool {
		_, ok := elem.AttrWithoutNamespace(name)
		return ok
	}
	if elem.(elements.HTMLElement).IsFormAssociatedElement() &&
		!util.IsNil(p.formElementPointer) &&
		!slices.ContainsFunc(p.stackOfOpenElements, func(n dom.Element) bool { return n.IsHtmlElement("template") }) && // TODO: replace with hasOneOfElems()?
		(elem.(elements.HTMLElement).IsFormListedElement() || !hasAttr("form")) &&
		dom.InTheSameTreeAs(intendedParent, p.formElementPointer) {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token]")
	}
	return elem
}

func (p *Parser) insertAtLocation(elem dom.Node, position insertionLocation) {
	switch position.tp {
	case insertionLocationAfterLastChild:
		dom.AppendChild(position.parentNode, elem)
	default:
		log.Panicf("unknown insertion mode %v", position.tp)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-an-element-at-the-adjusted-insertion-location
func (p *Parser) insertElementAtAdjustedInsertionLocation(elem dom.Node) {
	insertionLocation := p.appropriatePlaceForInsertionNode(nil)
	if !p.isFragmentParsing {
		// TODO: push a new element queue onto element's relevant agent's custom element reactions stack.
	}
	p.insertAtLocation(elem, insertionLocation)
	if !p.isFragmentParsing {
		// TODO: pop the element queue from element's relevant agent's custom element reactions stack,
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-a-foreign-element
func (p *Parser) insertForeignElement(token tagToken, namespace namespaces.Namespace, onlyAddElementToStack bool) dom.Element {
	insertionLocation := p.appropriatePlaceForInsertionNode(nil)
	elem := p.createElementForToken(token, namespace, insertionLocation.parentNode)
	if !onlyAddElementToStack {
		p.insertElementAtAdjustedInsertionLocation(elem)
	}
	p.stackOfOpenElements.push(elem)
	return elem
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-an-html-element
func (p *Parser) insertHtmlElement(token tagToken) dom.Element {
	return p.insertForeignElement(token, namespaces.Html, false)
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-a-comment
//
// position may be nil(= insertComment() will figure it out)
func (p *Parser) insertComment(data string, position *insertionLocation) {
	if position == nil {
		position = new(insertionLocation)
		*position = p.appropriatePlaceForInsertionNode(nil)
	}
	comment := dom.NewComment(position.parentNode.NodeDocument(), data)
	p.insertAtLocation(comment, *position)
}

// https://html.spec.whatwg.org/multipage/parsing.html#insert-a-character
func (p *Parser) insertCharacter(data rune) {
	insertionLocation := p.appropriatePlaceForInsertionNode(nil)
	if _, ok := insertionLocation.parentNode.(dom.Document); ok {
		// Document node cannot have text as children
		return
	}
	switch insertionLocation.tp {
	case insertionLocationAfterLastChild:
		parentNode := insertionLocation.parentNode
		parentChildren := parentNode.Children()
		var existingText dom.Text
		if len(parentChildren) != 0 {
			if t, ok := parentChildren[len(parentChildren)-1].(dom.Text); ok {
				existingText = t
			}
		}

		if !util.IsNil(existingText) {
			existingText.AppendText(string(data))
		} else {
			text := dom.NewText(parentNode.NodeDocument(), string(data))
			p.insertAtLocation(text, insertionLocation)
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#generic-raw-text-element-parsing-algorithm
func (p *Parser) parseGenericRawTextElement(token tagToken) {
	p.insertHtmlElement(token)
	p.tokenizer.state = rawtextState
	p.originalInsertionMode = p.insertionMode
	p.insertionMode = textInsertionMode
}

// https://html.spec.whatwg.org/multipage/parsing.html#generic-raw-text-element-parsing-algorithm
func (p *Parser) parseGenericRcdataElement(token tagToken) {
	p.insertHtmlElement(token)
	p.tokenizer.state = rcdataState
	p.originalInsertionMode = p.insertionMode
	p.insertionMode = textInsertionMode
}

// https://html.spec.whatwg.org/multipage/parsing.html#generate-implied-end-tags
func (p *Parser) generateImpliedEndTags(excludeFilter func(node dom.Element) bool) {
	htmlElems := []string{
		"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp", "rt", "rtc",
	}
	for {
		currentNode := p.currentNode()
		if slices.ContainsFunc(htmlElems, currentNode.IsHtmlElement) &&
			(excludeFilter == nil || !excludeFilter(currentNode)) {
			p.stackOfOpenElements.pop()
		} else {
			break
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#generate-all-implied-end-tags-thoroughly
func (p *Parser) generateAllImpliedEndTagsThroughly(excludeFilter func(n dom.Element) bool) {
	htmlElems := []string{
		"caption", "colgroup", "dd", "dt", "li", "optgroup", "option", "p",
		"rb", "rp", "rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
	}
	for {
		currentNode := p.currentNode()
		if slices.ContainsFunc(htmlElems, currentNode.IsHtmlElement) &&
			!excludeFilter(currentNode) {
			p.stackOfOpenElements.pop()
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#reset-the-insertion-mode-appropriately
func (p *Parser) resetInsertionModeAppropriately() {
	panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#reset-the-insertion-mode-appropriately]")
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-initial-insertion-mode
func (p *Parser) applyInitialInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		return
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, &insertionLocation{p.Document, insertionLocationAfterLastChild})
	} else if tk, ok := token.(*doctypeToken); ok {
		if tk.name == nil || *tk.name != "html" || tk.publicId != nil || (tk.systemId != nil && *tk.systemId != "about:legacy-compat") {
			p.parseErrorEncountered(token)
		}
		var name, publicId, systemId string = "", "", ""
		if tk.name != nil {
			name = *tk.name
		}
		if tk.publicId != nil {
			publicId = *tk.publicId
		}
		if tk.systemId != nil {
			systemId = *tk.systemId
		}

		doctypeNode := dom.NewDocumentType(p.Document, name, publicId, systemId)
		dom.AppendChild(p.Document, doctypeNode)

		p.Document.SetMode(dom.NoQuirks)
		if !p.Document.IsIframeSrcdocDocument() && !p.Document.IsParserCannotChangeMode() {
			if tk.forceQuirks ||
				(tk.name == nil || *tk.name != "html") ||
				(tk.publicId != nil && util.ToAsciiLowercase(*tk.publicId) == "-//w3o//dtd w3 html strict 3.0//en//") ||
				(tk.publicId != nil && util.ToAsciiLowercase(*tk.publicId) == "-/w3c/dtd html 4.0 transitional/en") ||
				(tk.publicId != nil && util.ToAsciiLowercase(*tk.publicId) == "html") ||
				(tk.systemId != nil && util.ToAsciiLowercase(*tk.systemId) == "http://www.ibm.com/data/dtd/v11/ibmxhtml1-transitional.dtd") ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "+//silmaril//dtd html pro v0r11 19970101//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//as//dtd html 3.0 aswedit + extensions//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//advasoft ltd//dtd html 3.0 aswedit + extensions//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.0 level 1//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.0 level 2//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.0 strict level 1//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.0 strict level 2//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.0 strict//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 2.1e//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 3.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 3.2 final//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 3.2//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html 3//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html level 0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html level 1//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html level 2//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html level 3//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html strict level 0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html strict level 1//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html strict level 2//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html strict level 3//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html strict//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//ietf//dtd html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//metrius//dtd metrius presentational//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//microsoft//dtd internet explorer 2.0 html strict//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//microsoft//dtd internet explorer 2.0 html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//microsoft//dtd internet explorer 2.0 tables//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//microsoft//dtd internet explorer 3.0 html strict//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//microsoft//dtd internet explorer 3.0 html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//microsoft//dtd internet explorer 3.0 tables//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//netscape comm. corp.//dtd html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//netscape comm. corp.//dtd strict html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//o'reilly and associates//dtd html 2.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//o'reilly and associates//dtd html extended 1.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//o'reilly and associates//dtd html extended relaxed 1.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//sq//dtd html 2.0 hotmetal + extensions//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//softquad software//dtd hotmetal pro 6.0::19990601::extensions to html 4.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//softquad//dtd hotmetal pro 4.0::19971010::extensions to html 4.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//spyglass//dtd html 2.0 extended//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//sun microsystems corp.//dtd hotjava html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//sun microsystems corp.//dtd hotjava strict html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 3 1995-03-24//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 3.2 draft//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 3.2 final//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 3.2//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 3.2s draft//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 4.0 frameset//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 4.0 transitional//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html experimental 19960712//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html experimental 970421//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd w3 html//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3o//dtd w3 html 3.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//webtechs//dtd mozilla html 2.0//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//webtechs//dtd mozilla html//")) ||
				(tk.systemId == nil && tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 4.01 frameset//")) ||
				(tk.systemId == nil && tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 4.01 transitional//")) {
				p.Document.SetMode(dom.Quirks)
			} else if (tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd xhtml 1.0 frameset//")) ||
				(tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd xhtml 1.0 transitional//")) ||
				(tk.systemId != nil && tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 4.01 frameset//")) ||
				(tk.systemId != nil && tk.publicId != nil && strings.HasPrefix(util.ToAsciiLowercase(*tk.publicId), "-//w3c//dtd html 4.01 transitional//")) {
				p.Document.SetMode(dom.LimitedQuirks)
			}
		}
		p.insertionMode = beforeHtmlInsertionMode
	} else {
		if !p.Document.IsIframeSrcdocDocument() {
			p.parseErrorEncountered(token)
			if !p.Document.IsParserCannotChangeMode() {
				p.Document.SetMode(dom.Quirks)
			}
		}
		p.insertionMode = beforeHtmlInsertionMode
		p.applyBeforeHtmlInsertionModeRules(token)
		return
	}

}

// https://html.spec.whatwg.org/multipage/parsing.html#the-before-html-insertion-mode
func (p *Parser) applyBeforeHtmlInsertionModeRules(token htmlToken) {
	if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, &insertionLocation{p.Document, insertionLocationAfterLastChild})
	} else if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		elem := p.createElementForToken(*tk, namespaces.Html, p.Document)
		dom.AppendChild(p.Document, elem)
		p.stackOfOpenElements.push(elem)
		p.insertionMode = beforeHeadInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isEnd && !slices.Contains([]string{"head", "body", "html", "br"}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else {
		elem := elements.NewHTMLHtmlElement(dom.ElementCreationCommonOptions{
			NodeDocument: p.Document,
			Namespace:    &namespaces.Html,
			LocalName:    "html",
		})
		dom.AppendChild(p.Document, elem)
		p.stackOfOpenElements.push(elem)
		p.insertionMode = beforeHeadInsertionMode
		p.applyBeforeHeadInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-before-head-insertion-mode
func (p *Parser) applyBeforeHeadInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		return
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "head" {
		elem := p.insertHtmlElement(*tk)
		p.headElementPointer = elem
		p.insertionMode = inHeadInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isEnd && !slices.Contains([]string{"head", "body", "html", "br"}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else {
		elem := p.insertHtmlElement(tagToken{tagName: "head"})
		p.headElementPointer = elem
		p.insertionMode = inHeadInsertionMode
		p.applyInHeadInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead
func (p *Parser) applyInHeadInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok {
		p.insertCharacter(tk.value)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"base", "basefont", "bgsound", "link"}, tk.tagName) {
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		if tk.isSelfClosing {
			tk.selfClosingAcknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "meta" {
		elem := p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		if !p.hasActiveSpeculativeParser {
			elem := elem
			if attr, ok := elem.AttrWithoutNamespace("charset"); ok {
				_ = attr
				// TODO: Set encoding based on charset
			}
			if attr, ok := elem.AttrWithoutNamespace("http-equiv"); ok && util.ToAsciiLowercase(attr) == "content-type" {
				_ = attr
				// TODO: Set encoding based on http-equiv Content-Type value
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "title" {
		p.parseGenericRcdataElement(*tk)
	} else if tk, ok := token.(*tagToken); ok &&
		(((tk.isStartTag() && tk.tagName == "title") && p.enableScripting) ||
			(tk.isStartTag() && slices.Contains([]string{"noframes", "style"}, tk.tagName))) {
		p.parseGenericRawTextElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "noscript" && !p.enableScripting {
		p.insertHtmlElement(*tk)
		p.insertionMode = inHeadNoscriptInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "script" {
		// STUB
		p.parseGenericRawTextElement(*tk)
		// panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead]")
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "head" {
		p.stackOfOpenElements.pop()
		p.insertionMode = afterHeadInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "template" {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead]")
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "template" {
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead]")
	} else if tk, ok := token.(*tagToken); ok &&
		((tk.isEnd && !slices.Contains([]string{"body", "html", "br"}, tk.tagName)) ||
			tk.isStartTag() && tk.tagName == "head") {
		p.parseErrorEncountered(token)
		return
	} else {
		p.stackOfOpenElements.pop()
		p.insertionMode = afterHeadInsertionMode
		p.applyAfterHeadInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inheadnoscript
func (p *Parser) applyInHeadNoscriptInsertionModeRules(token htmlToken) {
	if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "noscript" {
		p.stackOfOpenElements.pop()
		p.insertionMode = inHeadInsertionMode
	} else if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.applyInHeadInsertionModeRules(token)
	} else if _, ok := token.(*commentToken); ok {
		p.applyInHeadInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"basefont", "bgsound", "link", "meta", "noframes", "style"}, tk.tagName) {
		p.applyInHeadInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isEnd && !tk.isEndTag() && tk.tagName == "br") ||
		(tk.isStartTag() && slices.Contains([]string{"head", "noscript"}, tk.tagName)) {
		p.parseErrorEncountered(token)
		return
	} else {
		p.parseErrorEncountered(token)
		p.stackOfOpenElements.pop()
		p.insertionMode = inHeadInsertionMode
		p.applyInHeadInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-after-head-insertion-mode
func (p *Parser) applyAfterHeadInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.insertCharacter(tk.value)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "body" {
		p.insertHtmlElement(*tk)
		p.isFramesetNotOk = true
		p.insertionMode = inBodyInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "body" {
		p.insertHtmlElement(*tk)
		p.insertionMode = inFramesetInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{
		"base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title",
	}, tk.tagName) {
		p.parseErrorEncountered(token)
		p.stackOfOpenElements.push(p.headElementPointer)
		p.insertionMode = inHeadInsertionMode
		removeIdx := slices.Index(p.stackOfOpenElements, p.headElementPointer)
		p.stackOfOpenElements.remove(removeIdx)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "template" {
		p.applyInHeadInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok &&
		((tk.isEnd && !slices.Contains([]string{"body", "html", "br"}, tk.tagName)) ||
			(tk.isStartTag() && tk.tagName == "head")) {
		p.parseErrorEncountered(token)
		return
	} else {
		elem := p.insertHtmlElement(tagToken{tagName: "body"})
		p.headElementPointer = elem
		p.insertionMode = inBodyInsertionMode
		p.applyInBodyInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody
func (p *Parser) applyInBodyInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\u0000") {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.reconstructActiveFormattingElems()
		p.insertCharacter(tk.value)
	} else if tk, ok := token.(*charToken); ok {
		p.reconstructActiveFormattingElems()
		p.insertCharacter(tk.value)
		p.isFramesetNotOk = true
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.parseErrorEncountered(token)
		if p.stackOfOpenElements.hasElem("template") {
			return
		} else {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && slices.Contains([]string{"base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title"}, tk.tagName) ||
			(tk.isEndTag() && tk.tagName == "template")) {
		p.applyInHeadInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "body" {
		p.parseErrorEncountered(token)
		if len(p.stackOfOpenElements) == 1 ||
			!p.stackOfOpenElements[1].IsHtmlElement("body") ||
			p.stackOfOpenElements.hasElem("template") {
			return
		} else {
			p.isFramesetNotOk = true
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "frameset" {
		p.parseErrorEncountered(token)
		if len(p.stackOfOpenElements) == 1 ||
			!p.stackOfOpenElements[1].IsHtmlElement("body") {
			return
		} else if !p.isFramesetNotOk {
			return
		} else {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
	} else if _, ok := token.(*eofToken); ok {
		if len(p.stackOfTemplateInsertionModes) != 0 {
			p.applyInTemplateInsertionModeRules(token)
		} else {
			if p.stackOfOpenElements.hasOneOfElems([]string{
				"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
				"rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
				"body", "html",
			}) {
				p.parseErrorEncountered(token)
			}
			p.stopParsing()
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "body" {
		if !p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("body") }) {
			p.parseErrorEncountered(token)
			return
		} else if p.stackOfOpenElements.hasOneOfElems([]string{
			"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
			"rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
			"body", "html",
		}) {
			p.parseErrorEncountered(token)
		}
		p.insertionMode = afterBodyInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "html" {
		if !p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("body") }) {
			p.parseErrorEncountered(token)
			return
		} else if p.stackOfOpenElements.hasOneOfElems([]string{
			"dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
			"rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
			"body", "html",
		}) {
			p.parseErrorEncountered(token)
		}
		p.insertionMode = afterBodyInsertionMode
		p.applyAfterBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{
		"address", "article", "aside", "blockquote", "center", "details",
		"dialog", "dir", "div", "dl", "fieldset", "figcaption", "figure",
		"footer", "header", "hgroup", "main", "menu", "nav", "ol", "p",
		"search", "section", "summary", "ul",
	}, tk.tagName) {
		if p.haveElementInButtonScope(func(n dom.Element) bool {
			return n.IsHtmlElement("p")
		}) {
			p.closePElement()
		}
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, tk.tagName) {
		if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		if slices.ContainsFunc([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, p.currentNode().IsHtmlElement) {
			p.parseErrorEncountered(token)
			p.stackOfOpenElements.pop()
		}
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"pre", "listing"}, tk.tagName) {
		if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		p.insertHtmlElement(*tk)
		p.onNextToken = func(token htmlToken) parserControl {
			if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\n") {
				return parserControlIgnoreToken
			}
			return parserControlContinue
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "form" {
		if !util.IsNil(p.formElementPointer) && !p.stackOfOpenElements.hasElem("template") {
			p.parseErrorEncountered(token)
			return
		} else {
			if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
				p.closePElement()
			}
			elem := p.insertHtmlElement(*tk)
			if !p.stackOfOpenElements.hasElem("template") {
				p.formElementPointer = elem
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "li" {
		p.isFramesetNotOk = true
		node := p.currentNode()
		for {
			if node.IsHtmlElement("li") {
				p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement("li") })
				if !p.currentNode().IsHtmlElement("li") {
					p.parseErrorEncountered(token)
				}
				for {
					poppedElem := p.stackOfOpenElements.pop()
					if poppedElem.IsHtmlElement("li") {
						break
					}
				}
				break
			}
			if node.IsHtmlSpecialElement() &&
				!slices.ContainsFunc([]string{"address", "div", "p"}, node.IsHtmlElement) {
				break
			} else {
				nodeIdx := slices.Index(p.stackOfOpenElements, node) - 1
				node = p.stackOfOpenElements[nodeIdx]
			}
		}
		if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"dt", "dd"}, tk.tagName) {
		p.isFramesetNotOk = true
		node := p.currentNode()
		for {
			if node.IsHtmlElement("dd") {
				p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement("dd") })
				if !p.currentNode().IsHtmlElement("dd") {
					p.parseErrorEncountered(token)
				}
				for {
					poppedElem := p.stackOfOpenElements.pop()
					if poppedElem.IsHtmlElement("dd") {
						break
					}
				}
				break
			} else if node.IsHtmlElement("dt") {
				p.generateImpliedEndTags(func(node dom.Element) bool { return node.IsHtmlElement("dt") })
				if !p.currentNode().IsHtmlElement("dt") {
					p.parseErrorEncountered(token)
				}
				for {
					poppedElem := p.stackOfOpenElements.pop()
					if poppedElem.IsHtmlElement("dt") {
						break
					}
				}
				break
			}
			if node.IsHtmlSpecialElement() &&
				!slices.ContainsFunc([]string{"address", "div", "p"}, node.IsHtmlElement) {
				break
			} else {
				nodeIdx := slices.Index(p.stackOfOpenElements, node) - 1
				node = p.stackOfOpenElements[nodeIdx]
			}
		}
		if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "plaintext" {
		if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		p.insertHtmlElement(*tk)
		p.tokenizer.state = plaintextState
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "button" {
		if !p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("button") }) {
			p.parseErrorEncountered(token)
			p.generateImpliedEndTags(func(node dom.Element) bool { return true })
			for {
				poppedElem := p.stackOfOpenElements.pop()
				if poppedElem.IsHtmlElement("button") {
					break
				}
			}
		}
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
		p.isFramesetNotOk = true
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{
		"address", "article", "aside", "blockquote", "button", "center",
		"details", "dialog", "dir", "div", "dl", "fieldset", "figcaption",
		"figure", "footer", "header", "hgroup", "listing", "main", "menu",
		"nav", "ol", "pre", "search", "section", "select", "summary", "ul",
	}, tk.tagName) {
		if !p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) }) {
			p.parseErrorEncountered(token)
			return
		} else {
			p.generateImpliedEndTags(nil)
			if !p.currentNode().IsHtmlElement(tk.tagName) {
				p.parseErrorEncountered(token)
			}
			for {
				poppedElem := p.stackOfOpenElements.pop()
				if poppedElem.IsHtmlElement(tk.tagName) {
					break
				}
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "form" {
		if p.stackOfOpenElements.hasElem("template") {
			node := p.formElementPointer
			p.formElementPointer = nil
			if util.IsNil(node) || !p.haveElementInScope(func(n dom.Element) bool { return n == node }) {
				p.parseErrorEncountered(token)
				return
			}
			p.generateImpliedEndTags(nil)
			if p.currentNode() != node {
				p.parseErrorEncountered(token)
			}
			removeIdx := slices.Index(p.stackOfOpenElements, node)
			p.stackOfOpenElements.remove(removeIdx)
		} else {
			if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("form") }) {
				p.parseErrorEncountered(token)
				return
			}
			p.generateImpliedEndTags(nil)
			if !p.currentNode().IsHtmlElement("form") {
				p.parseErrorEncountered(token)
			}
			for {
				poppedElem := p.stackOfOpenElements.pop()
				if poppedElem.IsHtmlElement("form") {
					break
				}
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "p" {
		if !p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.parseErrorEncountered(token)
			p.insertHtmlElement(tagToken{tagName: "p"})
		}
		p.closePElement()
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "li" {
		if !p.haveElementInListItemScope(func(n dom.Element) bool { return n.IsHtmlElement("li") }) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement("li") })
		if !p.currentNode().IsHtmlElement("li") {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement("li") {
				break
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{"dd", "dt"}, tk.tagName) {
		if !p.haveElementInListItemScope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) }) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) })
		if !p.currentNode().IsHtmlElement(tk.tagName) {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement(tk.tagName) {
				break
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, tk.tagName) {
		if !p.haveElementInListItemScope(func(n dom.Element) bool {
			return slices.ContainsFunc([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, n.IsHtmlElement)
		}) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(nil)
		if !p.currentNode().IsHtmlElement(tk.tagName) {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if slices.ContainsFunc([]string{"h1", "h2", "h3", "h4", "h5", "h6"}, poppedElem.IsHtmlElement) {
				break
			}
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "a" {
		{
			lastMarkerIdx := slices.IndexFunc(p.listOfActiveFormattingElements, activeFormattingElement.isMarker)
			checkStartIdx := 0
			if 0 <= lastMarkerIdx {
				checkStartIdx = lastMarkerIdx + 1
			}
			var aElem dom.Element
			for i := checkStartIdx; i < len(p.listOfActiveFormattingElements); i++ {
				if p.listOfActiveFormattingElements[i].elem.IsHtmlElement("a") {
					aElem = p.listOfActiveFormattingElements[i].elem
				}
			}
			if !util.IsNil(aElem) {
				p.parseErrorEncountered(token)
				p.adoptionAgencyAlgorithm(*tk)
				removeIdx := slices.IndexFunc(p.listOfActiveFormattingElements, func(e activeFormattingElement) bool { return e.elem == aElem })
				p.listOfActiveFormattingElements = append(p.listOfActiveFormattingElements[:removeIdx], p.listOfActiveFormattingElements[removeIdx+1:]...)
				removeIdx = slices.Index(p.stackOfOpenElements, aElem)
				p.stackOfOpenElements.remove(removeIdx)
			}
		}
		p.reconstructActiveFormattingElems()
		elem := p.insertHtmlElement(*tk)
		p.listOfActiveFormattingElements.push(elem)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{
		"b", "big", "code", "em", "font", "i", "s", "small", "strike", "strong", "tt", "u",
	}, tk.tagName) {
		p.reconstructActiveFormattingElems()
		elem := p.insertHtmlElement(*tk)
		p.listOfActiveFormattingElements.push(elem)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "nobr" {
		p.reconstructActiveFormattingElems()
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("nobr") }) {
			p.parseErrorEncountered(token)
			p.adoptionAgencyAlgorithm(*tk)
			p.reconstructActiveFormattingElems()
		}
		elem := p.insertHtmlElement(*tk)
		p.listOfActiveFormattingElements.push(elem)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{
		"a", "b", "big", "code", "em", "font", "i", "nobr", "s", "small", "strike", "strong", "tt", "u",
	}, tk.tagName) {
		p.adoptionAgencyAlgorithm(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"applet", "marquee", "object"}, tk.tagName) {
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
		p.listOfActiveFormattingElements = append(p.listOfActiveFormattingElements, activeFormattingElemMarker)
		p.isFramesetNotOk = true
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{"applet", "marquee", "object"}, tk.tagName) {
		if !p.currentNode().IsHtmlElement(tk.tagName) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(nil)
		if !p.currentNode().IsHtmlElement(tk.tagName) {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement(tk.tagName) {
				break
			}
		}
		p.listOfActiveFormattingElements.clearUpToLastMarker()
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "table" {
		if (p.Document.Mode() != dom.Quirks) &&
			p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		p.insertHtmlElement(*tk)
		p.isFramesetNotOk = true
		p.insertionMode = inTableInsertionMode
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isEndTag() && tk.tagName == "br") ||
		(tk.isStartTag() && slices.Contains([]string{"area", "br", "embed", "img", "keygen", "wbr"}, tk.tagName)) {
		if tk.isEndTag() && tk.tagName == "br" {
			p.parseErrorEncountered(token)
			tk.attrs = []dom.AttrData{}
			tk.isEnd = false
		}
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		tk.selfClosingAcknowledged = true
		p.isFramesetNotOk = true
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "input" {
		if p.isFragmentParsing {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.parseErrorEncountered(token)
			for {
				poppedElem := p.stackOfOpenElements.pop()
				if poppedElem.IsHtmlElement("select") {
					break
				}
			}
		}
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		tk.selfClosingAcknowledged = true
		if typeAttr, ok := tk.Attr("type"); !ok || util.ToAsciiLowercase(typeAttr) != "hidden" {
			p.isFramesetNotOk = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "hr" {
		if p.haveElementInButtonScope(func(n dom.Element) bool {
			return n.IsHtmlElement("p")
		}) {
			p.closePElement()
		}
		if p.haveElementInScope(func(n dom.Element) bool {
			return n.IsHtmlElement("select")
		}) {
			p.generateImpliedEndTags(nil)
			if p.haveElementInScope(func(n dom.Element) bool {
				return n.IsHtmlElement("option") ||
					n.IsHtmlElement("optgroup")
			}) {
				p.parseErrorEncountered(token)
			}
		}
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		tk.selfClosingAcknowledged = true
		p.isFramesetNotOk = true
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "image" {
		p.parseErrorEncountered(token)
		tk.tagName = "img"
		p.applyInBodyInsertionModeRules(tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "textarea" {
		p.insertHtmlElement(*tk)
		p.onNextToken = func(token htmlToken) parserControl {
			if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\n") {
				return parserControlIgnoreToken
			}
			return parserControlContinue
		}
		p.tokenizer.state = rcdataState
		p.originalInsertionMode = p.insertionMode
		p.isFramesetNotOk = true
		p.insertionMode = textInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "xmp" {
		if p.haveElementInButtonScope(func(n dom.Element) bool { return n.IsHtmlElement("p") }) {
			p.closePElement()
		}
		p.reconstructActiveFormattingElems()
		p.isFramesetNotOk = true
		p.parseGenericRawTextElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "iframe" {
		p.isFramesetNotOk = true
		p.parseGenericRawTextElement(*tk)
	} else if tk, ok := token.(*tagToken); ok &&
		((tk.isStartTag() && tk.tagName == "noembed") ||
			(tk.isStartTag() && tk.tagName == "noscript" && !p.enableScripting)) {
		p.parseGenericRawTextElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "select" {
		if p.isFragmentParsing {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody]")
		}
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.parseErrorEncountered(token)
			for {
				poppedElem := p.stackOfOpenElements.pop()
				if poppedElem.IsHtmlElement("select") {
					break
				}
			}
			return
		}
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
		p.isFramesetNotOk = true
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "option" {
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement("opgroup") })
			if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("option") }) {
				p.parseErrorEncountered(token)
			}
		} else {
			if p.currentNode().IsHtmlElement("option") {
				p.stackOfOpenElements.pop()
			}
		}
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "optgroup" {
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("select") }) {
			p.generateImpliedEndTags(nil)
			if p.haveElementInScope(func(n dom.Element) bool {
				return n.IsHtmlElement("option") ||
					n.IsHtmlElement("optgroup")
			}) {
				p.parseErrorEncountered(token)
			}
		} else {
			if p.currentNode().IsHtmlElement("option") {
				p.stackOfOpenElements.pop()
			}
		}
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"rb", "rtc"}, tk.tagName) {
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("ruby") }) {
			p.generateImpliedEndTags(nil)
			if !p.currentNode().IsHtmlElement("ruby") {
				p.parseErrorEncountered(token)
			}
		}
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"rp", "rt"}, tk.tagName) {
		if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("ruby") }) {
			p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement("rtc") })
			if !p.currentNode().IsHtmlElement("rtc") &&
				!p.currentNode().IsHtmlElement("ruby") {
				p.parseErrorEncountered(token)
			}
		}
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "math" {
		p.reconstructActiveFormattingElems()
		adjustMathmlAttrs(tk)
		parserAdjustForeignAttrs(tk)
		p.insertForeignElement(*tk, namespaces.Mathml, false)
		if tk.isSelfClosing {
			p.stackOfOpenElements.pop()
			tk.selfClosingAcknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "svg" {
		p.reconstructActiveFormattingElems()
		adjustSvgAttrs(tk)
		parserAdjustForeignAttrs(tk)
		p.insertForeignElement(*tk, namespaces.Svg, false)
		if tk.isSelfClosing {
			p.stackOfOpenElements.pop()
			tk.selfClosingAcknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{
		"caption", "col", "colgroup", "frame", "head", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && !tk.isEnd {
		p.reconstructActiveFormattingElems()
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isEnd {
		nodeIdx := len(p.stackOfOpenElements) - 1
		node := func() dom.Element {
			return p.stackOfOpenElements[nodeIdx]
		}
		for {
			if node().IsHtmlElement(tk.tagName) {
				p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) })
				if node() != p.currentNode() {
					p.parseErrorEncountered(token)
				}
				targetNode := node()
				for p.stackOfOpenElements.pop() != targetNode {
				}
				return
			}
			if node().IsHtmlSpecialElement() {
				p.parseErrorEncountered(token)
				return
			}
			nodeIdx--
		}
	} else {
		log.Printf("[in-body insertion mode] Unrecognized token %v", token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#close-a-p-element
func (p *Parser) closePElement() {
	p.generateImpliedEndTags(func(n dom.Element) bool { return n.IsHtmlElement("p") })
	if !p.currentNode().IsHtmlElement("p") {
		p.parseErrorEncountered(p.currentNode().TagToken().(*tagToken))
	}
	for {
		poppedElem := p.stackOfOpenElements.pop()
		if poppedElem.IsHtmlElement("p") {
			break
		}
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#adoption-agency-algorithm
func (p *Parser) adoptionAgencyAlgorithm(token tagToken) {
	subject := token.tagName
	if p.currentNode().IsHtmlElement(subject) &&
		!slices.ContainsFunc(p.listOfActiveFormattingElements, func(e activeFormattingElement) bool {
			return e.elem == p.currentNode()
		}) {
		p.stackOfOpenElements.pop()
		return
	}
	panic("TODO")
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incdata
func (p *Parser) applyTextInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok {
		p.insertCharacter(tk.value)
	} else if _, ok := token.(*eofToken); ok {
		p.parseErrorEncountered(token)
		if p.currentNode().IsHtmlElement("script") {
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incdata]")
		}
		p.stackOfOpenElements.pop()
		p.insertionMode = p.originalInsertionMode
		p.applyCurrentInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "script" {
		// STUB
		// panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incdata]")
	} else if tk, ok := token.(*tagToken); ok && tk.isEnd {
		p.stackOfOpenElements.pop()
		p.insertionMode = p.originalInsertionMode
	} else {
		log.Printf("[text insertion mode] Unrecognized token %v", token)
	}
}

func (p *Parser) applyInTableInsertionModeRules(token htmlToken) {
	clearStackBackToTableContext := func() {
		for !slices.ContainsFunc([]string{"table", "template", "html"}, p.currentNode().IsHtmlElement) {
			p.stackOfOpenElements.pop()
		}
	}

	if _, ok := token.(*charToken); ok && slices.ContainsFunc([]string{
		"table", "tbody", "template", "tfoot", "thead", "tr",
	}, p.currentNode().IsHtmlElement) {
		p.pendingTableCharTokens = []charToken{}
		p.originalInsertionMode = p.insertionMode
		p.insertionMode = inTableTextInsertionMode
		p.applyInTableTextInsertionModeRules(token)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "caption" {
		clearStackBackToTableContext()
		p.listOfActiveFormattingElements = append(p.listOfActiveFormattingElements, activeFormattingElemMarker)
		p.insertHtmlElement(*tk)
		p.insertionMode = inCaptionInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "colgroup" {
		clearStackBackToTableContext()
		p.insertHtmlElement(*tk)
		p.insertionMode = inColumnGroupInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "col" {
		clearStackBackToTableContext()
		p.insertHtmlElement(tagToken{tagName: "colgroup"})
		p.insertionMode = inColumnGroupInsertionMode
		p.applyInColumnGroupInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"tbody", "tfoot", "thead"}, tk.tagName) {
		clearStackBackToTableContext()
		p.insertionMode = inTableBodyInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"td", "th", "tr"}, tk.tagName) {
		clearStackBackToTableContext()
		p.insertHtmlElement(tagToken{tagName: "tbody"})
		p.inTableBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "table" {
		p.parseErrorEncountered(token)
		if !p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("table") }) {
			return
		} else {
			for {
				poppedElem := p.stackOfOpenElements.pop()
				if poppedElem.IsHtmlElement("table") {
					break
				}
			}
			p.resetInsertionModeAppropriately()
			p.applyCurrentInsertionModeRules(token)
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "table" {
		if !p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("table") }) {
			p.parseErrorEncountered(token)
			return
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement("table") {
				break
			}
		}
		p.resetInsertionModeAppropriately()
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{
		"body", "caption", "col", "colgroup", "html", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && slices.Contains([]string{"style", "script", "template"}, tk.tagName)) ||
		(tk.isEndTag() && tk.tagName == "template") {
		p.applyInHeadInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "input" &&
		func() bool {
			if attr, ok := tk.Attr("type"); ok && util.ToAsciiLowercase(attr) == "hidden" {
				return true
			}
			return false
		}() {
		p.parseErrorEncountered(token)
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		tk.selfClosingAcknowledged = true
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "form" {
		if p.stackOfOpenElements.hasElem("template") {
			node := p.formElementPointer
			p.formElementPointer = nil
			if util.IsNil(node) || !p.haveElementInScope(func(n dom.Element) bool { return n == node }) {
				p.parseErrorEncountered(token)
				return
			}
			p.generateImpliedEndTags(nil)
			if p.currentNode() != node {
				p.parseErrorEncountered(token)
			}
			removeIdx := slices.Index(p.stackOfOpenElements, node)
			p.stackOfOpenElements.remove(removeIdx)
		} else {
			p.parseErrorEncountered(token)
			if p.haveElementInScope(func(n dom.Element) bool { return n.IsHtmlElement("form") }) ||
				!util.IsNil(p.formElementPointer) {
				return
			}
			p.insertHtmlElement(*tk)
			p.stackOfOpenElements.pop()
		}
	} else if _, ok := token.(*eofToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else {
		p.parseErrorEncountered(token)
		p.enableFosterParenting = true
		p.applyInBodyInsertionModeRules(token)
		p.enableFosterParenting = false
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intabletext
func (p *Parser) applyInTableTextInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.value == 0x0000 {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*charToken); ok {
		p.pendingTableCharTokens = append(p.pendingTableCharTokens, *tk)
	} else {
		if slices.ContainsFunc(p.pendingTableCharTokens, func(t charToken) bool { return !util.IsAsciiWhitespace(t.value) }) {
			p.parseErrorEncountered(token)
			// Below do the same thing as "else" in "in table" insertion mode.
			p.enableFosterParenting = true
			for _, tk := range p.pendingTableCharTokens {
				p.applyInBodyInsertionModeRules(tk)
			}
			p.enableFosterParenting = false
		} else {
			for _, tk := range p.pendingTableCharTokens {
				p.insertCharacter(tk.value)
			}
		}
		p.insertionMode = p.originalInsertionMode
		p.applyCurrentInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incaption
func (p *Parser) applyInCaptionInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "caption" {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement("caption") }) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(nil)
		if !p.currentNode().IsHtmlElement("caption") {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement("caption") {
				break
			}
		}
		p.listOfActiveFormattingElements.clearUpToLastMarker()
		p.insertionMode = inTableInsertionMode
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && slices.Contains([]string{"caption", "col", "colgroup", "tbody", "td", "tfoot", "th", "thead", "tr"}, tk.tagName)) ||
		(tk.isEndTag() && tk.tagName == "table") {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement("caption") }) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(nil)
		if !p.currentNode().IsHtmlElement("caption") {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement("caption") {
				break
			}
		}
		p.listOfActiveFormattingElements.clearUpToLastMarker()
		p.insertionMode = inTableInsertionMode
		p.inTableBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && (tk.isEndTag() && slices.Contains([]string{
		"body", "col", "colgroup", "html", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tagName)) {
		p.parseErrorEncountered(token)
		return
	} else {
		p.applyInBodyInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-incolgroup
func (p *Parser) applyInColumnGroupInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.insertCharacter(tk.value)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "col" {
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		if tk.isSelfClosing {
			tk.selfClosingAcknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "colgroup" {
		if !p.currentNode().IsHtmlElement("colgroup") {
			p.parseErrorEncountered(token)
			return
		}
		p.stackOfOpenElements.pop()
		p.insertionMode = inTableInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "col" {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && tk.tagName == "template") ||
		(tk.isEndTag() && tk.tagName == "template") {
		p.applyInHeadInsertionModeRules(token)
	} else if _, ok := token.(*eofToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else {
		if !p.currentNode().IsHtmlElement("colgroup") {
			p.parseErrorEncountered(token)
			return
		}
		p.stackOfOpenElements.pop()
		p.insertionMode = inTableInsertionMode
		p.applyInTableInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intbody
func (p *Parser) inTableBodyInsertionModeRules(token htmlToken) {
	clearStackBackToTableBodyContext := func() {
		for slices.ContainsFunc([]string{"tbody", "tfoot", "thead", "template", "html"}, p.currentNode().IsHtmlElement) {
			p.stackOfOpenElements.pop()
		}
	}
	if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "tr" {
		clearStackBackToTableBodyContext()
		p.insertHtmlElement(*tk)
		p.insertionMode = inRowInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"th", "td"}, tk.tagName) {
		p.parseErrorEncountered(token)
		clearStackBackToTableBodyContext()
		p.insertHtmlElement(tagToken{tagName: "tr"})
		p.insertionMode = inRowInsertionMode
		p.applyInRowInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{"tbody", "tfoot", "thead"}, tk.tagName) {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) }) {
			p.parseErrorEncountered(token)
			return
		}
		clearStackBackToTableBodyContext()
		p.stackOfOpenElements.pop()
		p.insertionMode = inTableInsertionMode
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && slices.Contains([]string{"caption", "col", "colgroup", "tbody", "tfoot", "thead"}, tk.tagName)) ||
		(tk.isEndTag() && tk.tagName == "table") {
		if !p.haveElementInTableScope(func(n dom.Element) bool {
			return slices.ContainsFunc([]string{"tbody", "thead", "tfoot"}, n.IsHtmlElement)
		}) {
			p.parseErrorEncountered(token)
			return
		}
		clearStackBackToTableBodyContext()
		p.stackOfOpenElements.pop()
		p.insertionMode = inTableInsertionMode
	} else if tk.isEndTag() && slices.Contains([]string{"body", "caption", "col", "colgroup", "html", "td", "th", "tr"}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else {
		p.applyInTableInsertionModeRules(token)
	}

}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intr
func (p *Parser) applyInRowInsertionModeRules(token htmlToken) {
	clearStackBackToTableRowContext := func() {
		for slices.ContainsFunc([]string{"tr", "template", "html"}, p.currentNode().IsHtmlElement) {
			p.stackOfOpenElements.pop()
		}
		panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intr]")
	}
	if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"th", "td"}, tk.tagName) {
		clearStackBackToTableRowContext()
		p.insertHtmlElement(*tk)
		p.insertionMode = inCellInsertionMode
		p.listOfActiveFormattingElements = append(p.listOfActiveFormattingElements, activeFormattingElemMarker)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "tr" {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement("tr") }) {
			p.parseErrorEncountered(token)
			return
		}
		clearStackBackToTableRowContext()
		p.stackOfOpenElements.pop()
		p.insertionMode = inTableBodyInsertionMode
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && slices.Contains([]string{"caption", "col", "colgroup", "tbody", "tfoot", "thead", "tr"}, tk.tagName)) ||
		(tk.isEndTag() && tk.tagName == "table") {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement("tr") }) {
			p.parseErrorEncountered(token)
			return
		}
		clearStackBackToTableRowContext()
		p.stackOfOpenElements.pop()
		p.insertionMode = inTableBodyInsertionMode
		p.inTableBodyInsertionModeRules(token)
	} else if tk.isEndTag() && slices.Contains([]string{"tbody", "tfoot", "thead"}, tk.tagName) {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) }) {
			p.parseErrorEncountered(token)
			return
		}
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement("tr") }) {
			return
		} else {
			clearStackBackToTableRowContext()
			p.stackOfOpenElements.pop()
			p.insertionMode = inTableBodyInsertionMode
		}
	} else if tk.isEndTag() && slices.Contains([]string{"body", "caption", "col", "colgroup", "html", "td", "th"}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else {
		p.applyInTableInsertionModeRules(token)
	}
}

func (p *Parser) applyInCellInsertionModeRules(token htmlToken) {
	closeCell := func() {
		p.generateImpliedEndTags(nil)
		if !slices.ContainsFunc([]string{"td", "th"}, p.currentNode().IsHtmlElement) {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if slices.ContainsFunc([]string{"td", "th"}, poppedElem.IsHtmlElement) {
				break
			}
		}
		p.listOfActiveFormattingElements.clearUpToLastMarker()
		p.insertionMode = inRowInsertionMode
	}

	if tk, ok := token.(*tagToken); ok && tk.isEndTag() && slices.Contains([]string{"th", "td"}, tk.tagName) {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) }) {
			p.parseErrorEncountered(token)
			return
		}
		p.generateImpliedEndTags(nil)
		if !p.currentNode().IsHtmlElement(tk.tagName) {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement(tk.tagName) {
				break
			}
		}
		p.listOfActiveFormattingElements.clearUpToLastMarker()
		p.insertionMode = inRowInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{
		"caption", "col", "colgroup", "tbody", "td", "tfoot", "th", "thead", "tr",
	}, tk.tagName) {
		if !p.haveElementInTableScope(func(n dom.Element) bool {
			return slices.ContainsFunc([]string{"td", "th"}, n.IsHtmlElement)
		}) {
			panic("we should have td or th in SOE at this point")
		}
		closeCell()
		p.applyInRowInsertionModeRules(token)
	} else if tk.isEndTag() && slices.Contains([]string{"body", "caption", "col", "colgroup", "html"}, tk.tagName) {
		p.parseErrorEncountered(token)
		return
	} else if tk.isEndTag() && slices.Contains([]string{"table", "tbody", "tfoot", "thead", "tr"}, tk.tagName) {
		if !p.haveElementInTableScope(func(n dom.Element) bool { return n.IsHtmlElement(tk.tagName) }) {
			p.parseErrorEncountered(token)
			return
		}
		closeCell()
		p.applyInRowInsertionModeRules(token)
	} else {
		p.applyInBodyInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-intemplate
func (p *Parser) applyInTemplateInsertionModeRules(token htmlToken) {
	if _, ok := token.(*charToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else if _, ok := token.(*commentToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else if _, ok := token.(*doctypeToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok &&
		(tk.isStartTag() && slices.Contains([]string{
			"base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title",
		}, tk.tagName)) ||
		(tk.isEndTag() && tk.tagName == "template") {
		p.applyInHeadInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{
		"caption", "colgroup", "tbody", "tfoot", "thead",
	}, tk.tagName) {
		p.stackOfTemplateInsertionModes.pop()
		p.stackOfTemplateInsertionModes.push(inTableInsertionMode)
		p.insertionMode = inTableInsertionMode
		p.applyInTableInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "col" {
		p.stackOfTemplateInsertionModes.pop()
		p.stackOfTemplateInsertionModes.push(inColumnGroupInsertionMode)
		p.insertionMode = inColumnGroupInsertionMode
		p.applyInColumnGroupInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "tr" {
		p.stackOfTemplateInsertionModes.pop()
		p.stackOfTemplateInsertionModes.push(inTableBodyInsertionMode)
		p.insertionMode = inTableBodyInsertionMode
		p.inTableBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && slices.Contains([]string{"td", "th"}, tk.tagName) {
		p.stackOfTemplateInsertionModes.pop()
		p.stackOfTemplateInsertionModes.push(inRowInsertionMode)
		p.insertionMode = inRowInsertionMode
		p.applyInRowInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() {
		p.stackOfTemplateInsertionModes.pop()
		p.stackOfTemplateInsertionModes.push(inBodyInsertionMode)
		p.insertionMode = inBodyInsertionMode
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() {
		p.parseErrorEncountered(token)
		return
	} else if _, ok := token.(*eofToken); ok {
		if !p.stackOfOpenElements.hasElem("template") {
			p.stopParsing()
		} else {
			p.parseErrorEncountered(token)
		}
		for {
			poppedElem := p.stackOfOpenElements.pop()
			if poppedElem.IsHtmlElement("template") {
				break
			}
		}
		p.listOfActiveFormattingElements.clearUpToLastMarker()
		p.stackOfTemplateInsertionModes.pop()
		p.resetInsertionModeAppropriately()
		p.applyCurrentInsertionModeRules(token)
	} else {
		log.Printf("[in template insertion mode] Unrecognized token %v", token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-afterbody
func (p *Parser) applyAfterBodyInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, &insertionLocation{
			parentNode: p.stackOfOpenElements[0],
			tp:         insertionLocationAfterLastChild,
		})
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "html" {
		if p.isFragmentParsing {
			p.parseErrorEncountered(token)
			return
		}
		p.insertionMode = afterAfterBodyInsertionMode
	} else if _, ok := token.(*eofToken); ok {
		p.stopParsing()
	} else {
		p.parseErrorEncountered(token)
		p.insertionMode = inBodyInsertionMode
		p.applyInBodyInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inframeset
func (p *Parser) applyInFramesetInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.insertCharacter(tk.value)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "frameset" {
		p.insertHtmlElement(*tk)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "framesets" {
		if p.currentNode().IsHtmlElement("html") &&
			util.IsNil(p.currentNode().Parent()) {
			// current node is root html node
			p.parseErrorEncountered(token)
			return
		}
		p.stackOfOpenElements.pop()
		if !p.isFragmentParsing && !p.currentNode().IsHtmlElement("frameset") {
			p.insertionMode = afterFramesetInsertionMode
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "frame" {
		p.insertHtmlElement(*tk)
		p.stackOfOpenElements.pop()
		if tk.isSelfClosing {
			tk.selfClosingAcknowledged = true
		}
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "noframes" {
		p.applyInHeadInsertionModeRules(token)
	} else if _, ok := token.(*eofToken); ok {
		if !p.currentNode().IsHtmlElement("html") ||
			!util.IsNil(p.currentNode().Parent()) {
			// current node is NOT root html node
			p.parseErrorEncountered(token)
		}
		p.stopParsing()
	} else {
		p.parseErrorEncountered(token)
		return
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-afterframeset
func (p *Parser) applyAfterFramesetInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.insertCharacter(tk.value)
	} else if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, nil)
	} else if _, ok := token.(*doctypeToken); ok {
		p.parseErrorEncountered(token)
		return
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isEndTag() && tk.tagName == "html" {
		p.insertionMode = afterAfterFramesetInsertionMode
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "noframes" {
		p.applyInHeadInsertionModeRules(token)
	} else if _, ok := token.(*eofToken); ok {
		p.stopParsing()
	} else {
		p.parseErrorEncountered(token)
		return
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-after-after-body-insertion-mode
func (p *Parser) applyAfterAfterBodyInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, &insertionLocation{
			parentNode: p.Document,
			tp:         insertionLocationAfterLastChild,
		})
	} else if _, ok := token.(*doctypeToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if _, ok := token.(*eofToken); ok {
		p.stopParsing()
	} else {
		p.parseErrorEncountered(token)
		p.insertionMode = inBodyInsertionMode
		p.applyInBodyInsertionModeRules(token)
	}
}

// https://html.spec.whatwg.org/multipage/parsing.html#the-after-after-frameset-insertion-mode
func (p *Parser) applyAfterAfterFramesetInsertionModeRules(token htmlToken) {
	if tk, ok := token.(*commentToken); ok {
		p.insertComment(tk.data, &insertionLocation{
			parentNode: p.Document,
			tp:         insertionLocationAfterLastChild,
		})
	} else if _, ok := token.(*doctypeToken); ok {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*charToken); ok && tk.isCharTokenWithOneOf("\t\n\u000c\r ") {
		p.applyInBodyInsertionModeRules(token)
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "html" {
		p.applyInBodyInsertionModeRules(token)
	} else if _, ok := token.(*eofToken); ok {
		p.stopParsing()
	} else if tk, ok := token.(*tagToken); ok && tk.isStartTag() && tk.tagName == "noframes" {
		p.applyInHeadInsertionModeRules(token)
	} else {
		p.parseErrorEncountered(token)
		return
	}
}

func (p *Parser) applyCurrentInsertionModeRules(token htmlToken) {
	insertionModeFuncs := map[insertionMode]func(token htmlToken){
		initialInsertionMode:            p.applyInitialInsertionModeRules,
		beforeHtmlInsertionMode:         p.applyBeforeHtmlInsertionModeRules,
		beforeHeadInsertionMode:         p.applyBeforeHeadInsertionModeRules,
		inHeadInsertionMode:             p.applyInHeadInsertionModeRules,
		inHeadNoscriptInsertionMode:     p.applyInHeadNoscriptInsertionModeRules,
		afterHeadInsertionMode:          p.applyAfterHeadInsertionModeRules,
		inBodyInsertionMode:             p.applyInBodyInsertionModeRules,
		textInsertionMode:               p.applyTextInsertionModeRules,
		inTableInsertionMode:            p.applyInTableInsertionModeRules,
		inTableTextInsertionMode:        p.applyInTableTextInsertionModeRules,
		inCaptionInsertionMode:          p.applyInCaptionInsertionModeRules,
		inColumnGroupInsertionMode:      p.applyInColumnGroupInsertionModeRules,
		inTableBodyInsertionMode:        p.inTableBodyInsertionModeRules,
		inRowInsertionMode:              p.applyInRowInsertionModeRules,
		inCellInsertionMode:             p.applyInCellInsertionModeRules,
		inTemplateInsertionMode:         p.applyInTemplateInsertionModeRules,
		afterBodyInsertionMode:          p.applyAfterBodyInsertionModeRules,
		inFramesetInsertionMode:         p.applyInFramesetInsertionModeRules,
		afterFramesetInsertionMode:      p.applyAfterFramesetInsertionModeRules,
		afterAfterBodyInsertionMode:     p.applyAfterAfterBodyInsertionModeRules,
		afterAfterFramesetInsertionMode: p.applyAfterAfterFramesetInsertionModeRules,
	}
	insertionModeFuncs[p.insertionMode](token)
}

// https://html.spec.whatwg.org/multipage/parsing.html#stop-parsing
func (p *Parser) stopParsing() {
	p.runParser = false
	// TODO
}

var mathmlAttrAdjustMap = map[string]string{
	"definitionurl": "definitionURL",
}

// https://html.spec.whatwg.org/multipage/parsing.html#adjust-mathml-attributes
func adjustMathmlAttrs(token *tagToken) {
	for i, attr := range token.attrs {
		if newName, ok := mathmlAttrAdjustMap[attr.LocalName]; ok {
			token.attrs[i].LocalName = newName
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
	"xlink:actuate": {util.MakeStrPtr("xlink"), "actuate", namespaces.Xlink},
	"xlink:arcrole": {util.MakeStrPtr("xlink"), "arcrole", namespaces.Xlink},
	"xlink:href":    {util.MakeStrPtr("xlink"), "href", namespaces.Xlink},
	"xlink:role":    {util.MakeStrPtr("xlink"), "role", namespaces.Xlink},
	"xlink:show":    {util.MakeStrPtr("xlink"), "show", namespaces.Xlink},
	"xlink:title":   {util.MakeStrPtr("xlink"), "title", namespaces.Xlink},
	"xlink:type":    {util.MakeStrPtr("xlink"), "type", namespaces.Xlink},
	"xml:lang":      {util.MakeStrPtr("xml"), "lang", namespaces.Xml},
	"xml:space":     {util.MakeStrPtr("xml"), "space", namespaces.Xml},
	"xmlns":         {nil, "xmlns", namespaces.Xmlns},
	"xmlns:xlink":   {util.MakeStrPtr("xmlns"), "xlink", namespaces.Xmlns},
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
