// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package dom

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/inseo-oh/yw/util"
)

// Node represents a [DOM Node].
//
// [DOM Node]: https://dom.spec.whatwg.org/#concept-node
type Node interface {
	// NodeDocument returns [node document] of the node.
	//
	// [node document]: https://dom.spec.whatwg.org/#concept-node-document
	NodeDocument() Document

	// SetNodeDocument sets [node document] of the node to doc.
	//
	// [node document]: https://dom.spec.whatwg.org/#concept-node-document
	SetNodeDocument(doc Document)

	// Parent returns the parent element.
	Parent() Node

	// SetParent sets the parent element to node.
	SetParent(node Node)

	// Children returns children of the node.
	Children() []Node

	// SetChildren sets children of the node to nodes.
	SetChildren(nodes []Node)

	// FirstChild returns its first children.
	FirstChild() Node

	// LastChild returns its last children.
	LastChild() Node

	// FilterChildren returns list of nodes where node n of children satisfies filter(n).
	FilterChildren(filter func(n Node) bool) []Node

	// FilterElementChildren returns list of elements where element e of children
	// satisfies filter(e).
	FilterElementChildren(filter func(e Element) bool) []Element

	// FilterElementChildrenByLocalName returns list of elements where
	// element e's local name is namePair's Name and is in namePair's Namespace.
	FilterElementChildrenByLocalName(namePair NamePair) []Element

	// ChildTextNode returns contents of its first text node children.
	ChildTextNode() (string, bool)

	// Callbacks returns pointer to struct containng callbacks for node.
	Callbacks() *NodeCallbacks

	// RunInsertionSteps runs [insertion steps] of the Node specified by [NodeCallbacks], if present.
	//
	// [insertion steps]: https://dom.spec.whatwg.org/#concept-node-insert-ext
	RunInsertionSteps()

	// RunChildrenChangedSteps runs [children changed steps] of the Node specified by [NodeCallbacks], if present.
	//
	// [children changed steps]: https://dom.spec.whatwg.org/#concept-node-children-changed-ext
	RunChildrenChangedSteps()

	// RunPostConncectionSteps runs [post-connection steps] of the Node specified by [NodeCallbacks], if present.
	//
	// [post-connection steps]: https://dom.spec.whatwg.org/#concept-node-post-connection-ext
	RunPostConncectionSteps()

	// RunAdoptingSteps runs [adopting steps] of the Node specified by [NodeCallbacks], if present.
	//
	// [adopting steps]: https://dom.spec.whatwg.org/#concept-node-adopt-ext
	RunAdoptingSteps(oldDoc Document)

	// CssData returns CSS-specific data for this node. dom package doesn't do anything with this.
	CssData() any

	// SetCssData sets CSS-specific data for this node. dom package doesn't do anything with this.
	SetCssData(data any)

	// String returns description of the Node.
	String() string
}
type nodeImpl struct {
	children     []Node
	parent       Node
	nodeDocument Document
	callbacks    NodeCallbacks
	cssData      any
}

// NodeCallbacks holds callbacks needed for [Node]. All callback functions are optional.
type NodeCallbacks struct {
	RunInsertionSteps       func()                // Callback for [Node.RunChildrenChangedSteps]
	RunChildrenChangedSteps func()                // Callback for [Node.RunChildrenChangedSteps]
	RunPostConnectionSteps  func()                // Callback for [Node.RunPostConnectionSteps]
	RunAdoptingSteps        func(oldDoc Document) // Callback for [Node.RunAdoptingSteps]

	// Element callbacks -------------------------------------------------------

	IntrinsicSize                 func() (width float64, height float64) // Callback for [Element.IntrinsicSize]
	PoppedFromStackOfOpenElements func()                                 // Callback called when HTML parser pops node from stack of open elements.
	PresentationalHints           func() any                             // Callback called by CSS system to get presentational hints.
}

// NewNode constructs a new [Node].
func NewNode(doc Document) Node {
	return &nodeImpl{nodeDocument: doc}
}

func (n nodeImpl) CssData() any {
	return n.cssData
}
func (n *nodeImpl) SetCssData(data any) {
	n.cssData = data
}
func (n *nodeImpl) Callbacks() *NodeCallbacks {
	return &n.callbacks
}
func (n nodeImpl) RunInsertionSteps() {
	if c := n.callbacks.RunInsertionSteps; c != nil {
		c()
	}
}
func (n nodeImpl) RunChildrenChangedSteps() {
	if c := n.callbacks.RunChildrenChangedSteps; c != nil {
		c()
	}
}
func (n nodeImpl) RunPostConncectionSteps() {
	if c := n.callbacks.RunPostConnectionSteps; c != nil {
		c()
	}
}
func (n nodeImpl) RunAdoptingSteps(oldDoc Document) {
	if c := n.callbacks.RunAdoptingSteps; c != nil {
		c(oldDoc)
	}
}

func (n nodeImpl) String() string {
	panic("not implemented")
}
func (n *nodeImpl) SetParent(node Node) {
	n.parent = node
}
func (n nodeImpl) Parent() Node {
	return n.parent
}
func (n nodeImpl) NodeDocument() Document {
	return n.nodeDocument
}
func (n *nodeImpl) SetNodeDocument(doc Document) {
	n.nodeDocument = doc
}
func (n nodeImpl) Children() []Node {
	return n.children
}
func (n *nodeImpl) SetChildren(nodes []Node) {
	n.children = nodes
}
func (n nodeImpl) FirstChild() Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[0]
}
func (n nodeImpl) LastChild() Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[len(n.children)-1]
}
func (n nodeImpl) FilterChildren(filter func(n Node) bool) []Node {
	res := []Node{}
	for _, c := range n.children {
		if filter(c) {
			res = append(res, c)
		}
	}
	return res
}
func (n nodeImpl) FilterElementChildren(filter func(e Element) bool) []Element {
	children := n.FilterChildren(func(n Node) bool {
		if e, ok := n.(Element); ok {
			return filter(e)
		}
		return false
	})
	res := []Element{}
	for _, n := range children {
		res = append(res, n.(Element))
	}
	return res
}
func (n nodeImpl) FilterElementChildrenByLocalName(namePair NamePair) []Element {
	return n.FilterElementChildren(func(e Element) bool {
		return e.IsElement(namePair)
	})
}

func (n nodeImpl) ChildTextNode() (string, bool) {
	textNodes := n.FilterChildren(func(n Node) bool { _, ok := n.(Text); return ok })
	if len(textNodes) == 0 {
		return "", false
	}
	return textNodes[0].(Text).Text(), true
}

// A lot of "methods" are not implemented as nodeImpl's methods, because if we did,
// when the method gets called, the receiver would point to the nodeImpl struct,
// not the original node pointer (and we can't use interface as receivers in Go).
//
// This would become a problem for functions comparing node pointers, so we don't.
// And for the same reason, any function calling those functions also must not be
// a method.
// This is simply due to nature of tree operation, and the fact that we have
// hierarchical type system on a language that wasn't designed to do so.
// This (hopefully) should be rare outside of this file.
//
// Functions that are implemented as methods only deal with itself.

// NextSibling returns [next sibling] of the node.
//
// [next sibling]: https://dom.spec.whatwg.org/#concept-tree-next-sibling
func NextSibling(node Node) Node {
	if util.IsNil(node.Parent()) {
		return nil
	}
	p := node.Parent()
	idx := Index(node)
	if idx == len(p.Children())-1 {
		return nil
	}
	return p.Children()[idx+1]
}

// PrevSibling returns [previous sibling] of the node.
//
// [previous sibling]: https://dom.spec.whatwg.org/#concept-tree-previous-sibling
func PrevSibling(node Node) Node {
	if util.IsNil(node.Parent()) {
		return nil
	}
	p := node.Parent()
	idx := Index(node)
	if idx == 0 {
		return nil
	}
	return p.Children()[idx-1]
}

// Root returns [root] of the node.
//
// [root]: https://dom.spec.whatwg.org/#concept-tree-root
func Root(node Node) Node {
	var p Node = node
	for !util.IsNil(p.Parent()) {
		p = p.Parent()
	}
	return p
}

// InTheSameTreeAs reports whether two nodes share the same root.
func InTheSameTreeAs(node, other Node) bool {
	return Root(node) == Root(other)
}

// Index returns [index] of the node.
//
// [index]: https://dom.spec.whatwg.org/#concept-tree-index
func Index(node Node) int {
	p := node.Parent()
	if util.IsNil(p) {
		return 0
	}
	for i, child := range p.Children() {
		if child == node {
			return i
		}
	}
	log.Panicf("%v is not children of %v", node, p)
	return -1
}

// InclusiveDescendants returns [inclusive descendant] nodes of rootNode.
//
// [inclusive descendant]: https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
func InclusiveDescendants(rootNode Node) []Node {
	// In a nutshell: It's just DFS search.
	resNodes := []Node{}
	var lastNode Node

	for {
		currNode := lastNode

		var res Node
		if util.IsNil(currNode) {
			// This is our first call
			res = rootNode
		} else {
			if len(currNode.Children()) == 0 {
				// We don't have any more children
				res = nil
			} else {
				// Go to the first children
				res = currNode.Children()[0]
			}
			// If we don't have more children, move to the next sibling
			for util.IsNil(res) {
				res = NextSibling(currNode)
				if !util.IsNil(res) {
					break
				}
				// We don't even have the next sibling -> Move to the parent
				currNode = currNode.Parent()
				if currNode == rootNode || util.IsNil(currNode) {
					// We don't have parent, or we are currently at root. We stop here.
					res = nil
					break
				}
			}

		}
		if util.IsNil(res) {
			break
		}
		lastNode = res
		resNodes = append(resNodes, res)
	}
	return resNodes
}

// Descendants returns [descendant] nodes of rootNode.
//
// [descendant]: https://dom.spec.whatwg.org/#concept-tree-descendant
func Descendants(rootNode Node) []Node {
	return InclusiveDescendants(rootNode)[1:]
}

// InclusiveAncestors returns [inclusive ancestor] nodes of rootNode.
//
// [inclusive ancestor]: https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor
func InclusiveAncestors(node Node) []Node {
	res := []Node{node}
	p := node
	for !util.IsNil(p.Parent()) {
		p = p.Parent()
		res = append(res, p)
	}
	return res
}

// Ancestors returns [ancestor] nodes of rootNode.
//
// [ancestor]: https://dom.spec.whatwg.org/#concept-tree-ancestor
func Ancestors(rootNode Node) []Node {
	return InclusiveAncestors(rootNode)[1:]
}

// ShadowIncludingRoot returns [shadow-including root] of the node.
//
// [shadow-including root]: https://dom.spec.whatwg.org/#concept-shadow-including-root
func ShadowIncludingRoot(node Node) Node {
	root := Root(node)
	if sr, ok := root.(ShadowRoot); ok {
		return ShadowIncludingRoot(sr.Host())
	}
	return root
}

// ShadowIncludingInclusiveDescendants returns [shadow-including inclusive descendant] nodes of rootNode.
//
// [shadow-including inclusive descendant]: https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
func ShadowIncludingInclusiveDescendants(rootNode Node) []Node {
	descendants := InclusiveDescendants(rootNode)
	res := []Node{}
	for _, d := range descendants {
		if sr, ok := d.(ShadowRoot); ok {
			res = append(res, ShadowIncludingInclusiveDescendants(sr)...)
		} else {
			res = append(res, d)
		}
	}
	return res
}

// ShadowIncludingDescendants returns [shadow-including descendant] nodes of rootNode.
//
// [shadow-including descendant]: https://dom.spec.whatwg.org/#concept-shadow-including-descendant
func ShadowIncludingDescendants(rootNode Node) []Node {
	return ShadowIncludingInclusiveDescendants(rootNode)[1:]
}

// IsConnected reports whether node is [connected].
//
// [connected]: https://dom.spec.whatwg.org/#connected
func IsConnected(node Node) bool {
	return ShadowIncludingRoot(node) == node.NodeDocument()
}

// IsInDocumentTree reports whether node is [in a document tree].
//
// [in a document tree]: https://dom.spec.whatwg.org/#in-a-document-tree
func IsInDocumentTree(node Node) bool {
	_, ok := Root(node).(Document)
	return ok
}

// Insert inserts the node to parent before beforeChild.
// If beforeChild is nil, node is inserted at the end of parent's children instead.
//
// I'm not entirely sure what suppressObservers does yet (it's part of spec).
//
// Spec: https://dom.spec.whatwg.org/#concept-node-insert
func Insert(node, parent, beforeChild Node, suppressObservers bool) {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S1.
	nodes := []Node{node}
	if _, ok := node.(DocumentFragment); ok {
		nodes = node.Children()
	}
	// S2.
	count := len(nodes)
	// S3.
	if count == 0 {
		return
	}
	// S4.
	if _, ok := node.(DocumentFragment); ok {
		log.Panicf("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
	}
	// S5.
	if !util.IsNil(beforeChild) {
		// TODO[https://dom.spec.whatwg.org/#concept-node-insert]
		// 1. For each live range whose start node is parent and start offset is greater than child’s index, increase its start offset by count.
		// 2. For each live range whose end node is parent and end offset is greater than child’s index, increase its end offset by count.
	}
	// S6.
	prevSibling := parent.LastChild()
	if !util.IsNil(beforeChild) {
		prevSibling = PrevSibling(beforeChild)
	}
	_ = prevSibling
	// S7.
	for _, node := range nodes {
		// S7-1.
		AdoptNodeInto(node, parent.NodeDocument())
		if util.IsNil(beforeChild) {
			// S7-2.
			children := parent.Children()
			children = append(children, node)
			parent.SetChildren(children)
		} else {
			// S7-3.
			children := parent.Children()
			insertIndex := slices.Index(children, beforeChild)
			children = append(append(children[:insertIndex], node), children[insertIndex:]...)
			parent.SetChildren(children)
		}
		// S7-4.
		if parent, ok := parent.(Element); ok && parent.IsShadowHost() {
			panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
		}
		// S7-5.
		parentRoot := Root(parent)
		if sr, ok := parentRoot.(ShadowRoot); ok {
			_ = sr
			panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
		}
		// S7-6.
		// TODO: Run assign slottables for a tree with node’s root.
		// S7-7.
		for _, inclusiveDescendant := range ShadowIncludingDescendants(node) {
			// S7-7-1.
			inclusiveDescendant.RunInsertionSteps()
			if inclusiveDescendantElem, ok := inclusiveDescendant.(Element); ok {
				// S7-7-2.
				if reg := inclusiveDescendantElem.CustomElementRegistry(); reg == nil {
					reg = LookupCustomElementRegistry(inclusiveDescendant.Parent())
					inclusiveDescendantElem.SetCustomElementRegistry(reg)
				} else if reg.IsScoped {
					reg.ScopedDocumentSet = append(reg.ScopedDocumentSet, inclusiveDescendant.NodeDocument())
				} else if inclusiveDescendantElem.IsCustom() {
					// TODO: enqueue a custom element callback reaction with inclusiveDescendant, callback name "connectedCallback", and « ».
					panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
				} else {
					tryUpgradeElement(inclusiveDescendantElem)
				}
			} else if inclusiveDescendantSr, ok := inclusiveDescendant.(ShadowRoot); ok {
				// S7-7-3.
				_ = inclusiveDescendantSr
				// TODO: If inclusiveDescendant’s custom element registry is null and inclusiveDescendant’s keep custom element registry null is false, then set inclusiveDescendant’s custom element registry to the result of looking up a custom element registry given inclusiveDescendant’s host.
				// TODO: Otherwise, if inclusiveDescendant’s custom element registry is non-null and inclusiveDescendant’s custom element registry’s is scoped is true, append inclusiveDescendant’s node document to inclusiveDescendant’s custom element registry’s scoped document set.
				panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
			}
		}
	}
	// S8.
	if !suppressObservers {
		// TODO: queue a tree mutation record for parent with nodes, « », previousSibling, and child.
	}
	// S9.
	parent.RunChildrenChangedSteps()
	// S10.
	staticNodeList := []Node{}
	// S11.
	for _, node := range nodes {
		staticNodeList = append(staticNodeList, ShadowIncludingDescendants(node)...)
	}
	// S12.
	for _, node := range staticNodeList {
		if IsConnected(node) {
			node.RunPostConncectionSteps()
		}
	}

	node.SetParent(parent)
}

// AppendChild is shorthand for [Insert], that just adds child to the node.
func AppendChild(node, child Node) {
	Insert(child, node, nil, false)
}

// AdoptNodeInto adopts node into the document.
//
// Spec: https://dom.spec.whatwg.org/#concept-node-adopt
func AdoptNodeInto(node Node, document Document) {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S1.
	oldDocument := node.NodeDocument()
	// S2.
	if !util.IsNil(node.Parent()) {
		// TODO: remove node
		panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
	}
	// S3.
	if document != oldDocument {
		// S3-1.
		for _, inclusiveDescendant := range ShadowIncludingDescendants(node) {
			// S3-1-1.
			inclusiveDescendant.SetNodeDocument(document)
			if inclusiveDescendantSr, ok := inclusiveDescendant.(ShadowRoot); ok && IsGlobalCustomElementReigstry(LookupCustomElementRegistry(inclusiveDescendant)) {
				// S3-1-2.
				_ = inclusiveDescendantSr
				// TODO: set inclusiveDescendant’s custom element registry to document’s effective global custom element registry.
				inclusiveDescendantSr.SetCustomElementRegistry(document.EffectiveGlobalCustomElementRegistry())
				panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
			} else if e, ok := inclusiveDescendant.(Element); ok {
				// S3-1-3.
				// S3-1-3-1.
				attrs := e.Attrs()
				for i := range len(attrs) {
					attrs[i].SetNodeDocument(document)
				}
				// S3-1-3-2.
				if IsGlobalCustomElementReigstry(LookupCustomElementRegistry(inclusiveDescendant)) {
					// TODO: set inclusiveDescendant’s custom element registry to document’s effective global custom element registry.
					panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
				}
			}

		}
		// S3-2.
		for _, inclusiveDescendant := range ShadowIncludingDescendants(node) {
			if !inclusiveDescendant.(Element).IsCustom() {
				continue
			}
			// TODO: enqueue a custom element callback reaction with inclusiveDescendant, callback name "adoptedCallback", and « oldDocument, document ».
			panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
		}
		// S3-3.
		for _, inclusiveDescendant := range ShadowIncludingDescendants(node) {
			inclusiveDescendant.RunAdoptingSteps(oldDocument)
		}
	}
}

// PrintTree prints DOM tree to standard output.
func PrintTree(node Node) {
	currNode := node
	count := 0
	if !util.IsNil(currNode.Parent()) {
		for n := currNode.Parent(); !util.IsNil(n); n = n.Parent() {
			count += 4
		}
	}
	indent := strings.Repeat(" ", count)
	fmt.Printf("%s%v\n", indent, node)
	for _, child := range currNode.Children() {
		PrintTree(child)
	}
}

// LookupCustomElementRegistry returns custom element registry for the node, or
// nil if not applicable.
//
// Spec: https://html.spec.whatwg.org/multipage/custom-elements.html#look-up-a-custom-element-registry
func LookupCustomElementRegistry(node Node) *CustomElementRegistry {
	if x, ok := node.(Element); ok {
		return x.CustomElementRegistry()
	}
	if x, ok := node.(Document); ok {
		return x.CustomElementRegistry()
	}
	if x, ok := node.(ShadowRoot); ok {
		return x.CustomElementRegistry()
	}
	return nil
}

type testNode struct {
	nodeImpl
	name string
}

func (n testNode) String() string {
	return fmt.Sprintf("TestNode %s", n.name)
}
