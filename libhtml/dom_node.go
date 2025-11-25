package libhtml

import (
	"fmt"
	"log"
	"slices"
	"strings"
	cm "yw/libcommon"
)

type dom_Node_s struct {
	children      []dom_Node
	parent        dom_Node
	node_document dom_Document
	callbacks     dom_node_callbacks
}
type dom_node_callbacks struct {
	run_insertion_steps                func()
	run_children_changed_steps         func()
	run_post_connection_steps          func()
	run_adopting_steps                 func(old_doc dom_Document)
	popped_from_stack_of_open_elements func()
}
type dom_Node interface {
	get_callbacks() *dom_node_callbacks
	set_parent(node dom_Node)
	get_parent() dom_Node
	get_node_document() dom_Document
	set_node_document(doc dom_Document)
	get_children() []dom_Node
	set_children(nodes []dom_Node)
	first_child() dom_Node
	last_child() dom_Node
	filter_children(filter func(n dom_Node) bool) []dom_Node
	filter_elem_children(filter func(e dom_Element) bool) []dom_Element
	filter_elem_children_by_local_name(p dom_name_pair) []dom_Element
	get_child_text_node() (string, bool)
	String() string

	run_insertion_steps()
	run_children_changed_steps()
	run_post_connection_steps()
	run_adopting_steps(old_doc dom_Document)
}

func dom_make_Node(doc dom_Document) dom_Node {
	return &dom_Node_s{node_document: doc}
}

func (n *dom_Node_s) get_callbacks() *dom_node_callbacks {
	return &n.callbacks
}
func (n dom_Node_s) run_insertion_steps() {
	if c := n.callbacks.run_insertion_steps; c != nil {
		c()
	}
}
func (n dom_Node_s) run_children_changed_steps() {
	if c := n.callbacks.run_children_changed_steps; c != nil {
		c()
	}
}
func (n dom_Node_s) run_post_connection_steps() {
	if c := n.callbacks.run_post_connection_steps; c != nil {
		c()
	}
}
func (n dom_Node_s) run_adopting_steps(old_doc dom_Document) {
	if c := n.callbacks.run_adopting_steps; c != nil {
		c(old_doc)
	}
}

func (n dom_Node_s) String() string {
	panic("not implemented")
}
func (n *dom_Node_s) set_parent(node dom_Node) {
	n.parent = node
}
func (n dom_Node_s) get_parent() dom_Node {
	return n.parent
}
func (n dom_Node_s) get_node_document() dom_Document {
	return n.node_document
}
func (n *dom_Node_s) set_node_document(doc dom_Document) {
	n.node_document = doc
}
func (n dom_Node_s) get_children() []dom_Node {
	return n.children
}
func (n *dom_Node_s) set_children(nodes []dom_Node) {
	n.children = nodes
}
func (n dom_Node_s) first_child() dom_Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[0]
}
func (n dom_Node_s) last_child() dom_Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[len(n.children)-1]
}
func (n dom_Node_s) filter_children(filter func(n dom_Node) bool) []dom_Node {
	out := []dom_Node{}
	for _, c := range n.children {
		if filter(c) {
			out = append(out, c)
		}
	}
	return out
}
func (n dom_Node_s) filter_elem_children(filter func(e dom_Element) bool) []dom_Element {
	res := n.filter_children(func(n dom_Node) bool {
		if e, ok := n.(dom_Element); ok {
			return filter(e)
		}
		return false
	})
	out := []dom_Element{}
	for _, n := range res {
		out = append(out, n.(dom_Element))
	}
	return out
}
func (n dom_Node_s) filter_elem_children_by_local_name(p dom_name_pair) []dom_Element {
	return n.filter_elem_children(func(e dom_Element) bool {
		return e.is_element(p)
	})
}

func (n dom_Node_s) get_child_text_node() (string, bool) {
	text_nodes := n.filter_children(func(n dom_Node) bool { _, ok := n.(dom_Text); return ok })
	if len(text_nodes) == 0 {
		return "", false
	}
	return text_nodes[0].(dom_Text).get_text(), true
}

// A lot of "methods" are not implemented as dom_Node_s's methods, because if we did,
// when the method gets called, the receiver would point to the dom_Node_s struct,
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

func dom_node_get_next_sibling(node dom_Node) dom_Node {
	if cm.IsNil(node.get_parent()) {
		return nil
	}
	p := node.get_parent()
	idx := dom_node_get_index(node)
	if idx == len(p.get_children())-1 {
		return nil
	}
	return p.get_children()[idx+1]
}
func dom_node_get_prev_sibling(node dom_Node) dom_Node {
	if cm.IsNil(node.get_parent()) {
		return nil
	}
	p := node.get_parent()
	idx := dom_node_get_index(node)
	if idx == 0 {
		return nil
	}
	return p.get_children()[idx-1]
}
func dom_node_root(node dom_Node) dom_Node {
	var p dom_Node = node
	for !cm.IsNil(p.get_parent()) {
		p = p.get_parent()
	}
	return p
}
func dom_node_in_the_same_tree_as(node, other_p dom_Node) bool {
	return dom_node_root(node) == dom_node_root(other_p)
}

// https://dom.spec.whatwg.org/#concept-tree-index
func dom_node_get_index(node dom_Node) int {
	p := node.get_parent()
	if cm.IsNil(p) {
		return 0
	}
	for i, child := range p.get_children() {
		if child == node {
			return i
		}
	}
	log.Panicf("%v is not children of %v", node, p)
	return -1
}

// https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
func dom_node_inclusive_descendants(root_node dom_Node) []dom_Node {
	// In a nutshell: It's just DFS search.
	out := []dom_Node{}
	var last_node dom_Node

	for {
		curr_node := last_node

		var res dom_Node
		if cm.IsNil(curr_node) {
			// This is our first call
			res = root_node
		} else {
			if len(curr_node.get_children()) == 0 {
				// We don't have any more children
				res = nil
			} else {
				// Go to the first children
				res = curr_node.get_children()[0]
			}
			// If we don't have more children, move to the next sibling
			for cm.IsNil(res) {
				res = dom_node_get_next_sibling(curr_node)
				if !cm.IsNil(res) {
					break
				}
				// We don't even have the next sibling -> Move to the parent
				curr_node = curr_node.get_parent()
				if curr_node == root_node || cm.IsNil(curr_node) {
					// We don't have parent, or we are currently at root. We stop here.
					res = nil
					break
				}
			}

		}
		if cm.IsNil(res) {
			break
		}
		last_node = res
		out = append(out, res)
	}
	return out
}

// https://dom.spec.whatwg.org/#concept-tree-descendant
func dom_node_descendants(root_node dom_Node) []dom_Node {
	return dom_node_inclusive_descendants(root_node)[1:]
}

// https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor
func dom_node_inclusive_ancestors(node dom_Node) []dom_Node {
	out := []dom_Node{node}
	p := node
	for !cm.IsNil(p.get_parent()) {
		p = p.get_parent()
		out = append(out, p)
	}
	return out
}

// https://dom.spec.whatwg.org/#concept-tree-ancestor
func dom_node_ancestors(root_node dom_Node) []dom_Node {
	return dom_node_inclusive_ancestors(root_node)[1:]
}

// https://dom.spec.whatwg.org/#concept-shadow-including-root
func dom_node_shadow_including_root(node dom_Node) dom_Node {
	root := dom_node_root(node)
	if sr, ok := root.(dom_ShadowRoot); ok {
		return dom_node_shadow_including_root(sr.get_host())
	}
	return root
}

// https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
func dom_node_shadow_including_inclusive_descendants(root_node dom_Node) []dom_Node {
	descendants := dom_node_inclusive_descendants(root_node)
	out := []dom_Node{}
	for _, d := range descendants {
		if sr, ok := d.(dom_ShadowRoot); ok {
			out = append(out, dom_node_shadow_including_inclusive_descendants(sr)...)
		} else {
			out = append(out, d)
		}
	}
	return out
}

// https://dom.spec.whatwg.org/#concept-shadow-including-descendant
func dom_node_shadow_including_descendants(root_node dom_Node) []dom_Node {
	return dom_node_shadow_including_inclusive_descendants(root_node)[1:]
}

func dom_node_look_up_custom_element_registry(node dom_Node) *html_custom_element_registry {
	if x, ok := node.(dom_Element); ok {
		return x.get_custom_element_registry()
	}
	if x, ok := node.(dom_Document); ok {
		return x.get_custom_element_registry()
	}
	if x, ok := node.(dom_ShadowRoot); ok {
		return x.get_custom_element_registry()
	}
	return nil
}

// https://dom.spec.whatwg.org/#connected
func dom_node_is_connected(node dom_Node) bool {
	return dom_node_shadow_including_root(node) == node.get_node_document()
}

// https://dom.spec.whatwg.org/#in-a-document-tree
func dom_node_is_in_document_tree(node dom_Node) bool {
	_, ok := dom_node_root(node).(dom_Document)
	return ok
}

// https://dom.spec.whatwg.org/#concept-node-insert
func dom_node_insert(node, parent, before_child dom_Node, suppress_observers bool) {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S1.
	nodes := []dom_Node{node}
	if _, ok := node.(dom_DocumentFragment); ok {
		nodes = node.get_children()
	}
	// S2.
	count := len(nodes)
	// S3.
	if count == 0 {
		return
	}
	// S4.
	if _, ok := node.(dom_DocumentFragment); ok {
		log.Panicf("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
	}
	// S5.
	if !cm.IsNil(before_child) {
		// TODO[https://dom.spec.whatwg.org/#concept-node-insert]
		// 1. For each live range whose start node is parent and start offset is greater than child’s index, increase its start offset by count.
		// 2. For each live range whose end node is parent and end offset is greater than child’s index, increase its end offset by count.
	}
	// S6.
	prev_sibling := parent.last_child()
	if !cm.IsNil(before_child) {
		prev_sibling = dom_node_get_prev_sibling(before_child)
	}
	_ = prev_sibling
	// S7.
	for _, node := range nodes {
		// S7-1.
		dom_adopt_node_into(node, parent.get_node_document())
		if cm.IsNil(before_child) {
			// S7-2.
			children := parent.get_children()
			children = append(children, node)
			parent.set_children(children)
		} else {
			// S7-3.
			children := parent.get_children()
			insert_index := slices.Index(children, before_child)
			children = append(append(children[:insert_index], node), children[insert_index:]...)
			parent.set_children(children)
		}
		// S7-4.
		if parent, ok := parent.(dom_Element); ok && parent.is_shadow_host() {
			panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
		}
		// S7-5.
		parent_root := dom_node_root(parent)
		if sr, ok := parent_root.(dom_ShadowRoot); ok {
			_ = sr
			panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
		}
		// S7-6.
		// TODO: Run assign slottables for a tree with node’s root.
		// S7-7.
		for _, inclusive_descendant := range dom_node_shadow_including_descendants(node) {
			// S7-7-1.
			inclusive_descendant.run_insertion_steps()
			if inclusive_descendant_elem, ok := inclusive_descendant.(dom_Element); ok {
				// S7-7-2.
				if reg := inclusive_descendant_elem.get_custom_element_registry(); reg == nil {
					reg = dom_node_look_up_custom_element_registry(inclusive_descendant.get_parent())
					inclusive_descendant_elem.set_custom_elem_registry(reg)
				} else if reg.is_scoped {
					reg.scoped_document_set = append(reg.scoped_document_set, inclusive_descendant.get_node_document())
				} else if inclusive_descendant_elem.is_custom() {
					// TODO: enqueue a custom element callback reaction with inclusiveDescendant, callback name "connectedCallback", and « ».
					panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
				} else {
					html_try_upgrade_element(inclusive_descendant_elem)
				}
			} else if inclusive_descendant_sr, ok := inclusive_descendant.(dom_ShadowRoot); ok {
				// S7-7-3.
				_ = inclusive_descendant_sr
				// TODO: If inclusiveDescendant’s custom element registry is null and inclusiveDescendant’s keep custom element registry null is false, then set inclusiveDescendant’s custom element registry to the result of looking up a custom element registry given inclusiveDescendant’s host.
				// TODO: Otherwise, if inclusiveDescendant’s custom element registry is non-null and inclusiveDescendant’s custom element registry’s is scoped is true, append inclusiveDescendant’s node document to inclusiveDescendant’s custom element registry’s scoped document set.
				panic("TODO[https://dom.spec.whatwg.org/#concept-node-insert]")
			}
		}
	}
	// S8.
	if !suppress_observers {
		// TODO: queue a tree mutation record for parent with nodes, « », previousSibling, and child.
	}
	// S9.
	parent.run_children_changed_steps()
	// S10.
	static_node_list := []dom_Node{}
	// S11.
	for _, node := range nodes {
		static_node_list = append(static_node_list, dom_node_shadow_including_descendants(node)...)
	}
	// S12.
	for _, node := range static_node_list {
		if dom_node_is_connected(node) {
			node.run_post_connection_steps()
		}
	}

	node.set_parent(parent)
}
func dom_node_append_child(node, child dom_Node) {
	dom_node_insert(child, node, nil, false)
}

// https://dom.spec.whatwg.org/#concept-node-adopt
func dom_adopt_node_into(node dom_Node, document dom_Document) {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S1.
	old_document := node.get_node_document()
	// S2.
	if !cm.IsNil(node.get_parent()) {
		// TODO: remove node
		panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
	}
	// S3.
	if document != old_document {
		// S3-1.
		for _, inclusive_descendant := range dom_node_shadow_including_descendants(node) {
			// S3-1-1.
			inclusive_descendant.set_node_document(document)
			if inclusive_descendant_sr, ok := inclusive_descendant.(dom_ShadowRoot); ok && dom_is_global_custom_element_registry(dom_node_look_up_custom_element_registry(inclusive_descendant)) {
				// S3-1-2.
				_ = inclusive_descendant_sr
				// TODO: set inclusiveDescendant’s custom element registry to document’s effective global custom element registry.
				inclusive_descendant_sr.set_custom_elem_registry(document.effective_global_custom_element_registry())
				panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
			} else if e, ok := inclusive_descendant.(dom_Element); ok {
				// S3-1-3.
				// S3-1-3-1.
				attrs := e.get_attrs()
				for i := range len(attrs) {
					attrs[i].set_node_document(document)
				}
				// S3-1-3-2.
				if dom_is_global_custom_element_registry(dom_node_look_up_custom_element_registry(inclusive_descendant)) {
					// TODO: set inclusiveDescendant’s custom element registry to document’s effective global custom element registry.
					panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
				}
			}

		}
		// S3-2.
		for _, inclusive_descendant := range dom_node_shadow_including_descendants(node) {
			if !inclusive_descendant.(dom_Element).is_custom() {
				continue
			}
			// TODO: enqueue a custom element callback reaction with inclusiveDescendant, callback name "adoptedCallback", and « oldDocument, document ».
			panic("TODO[https://dom.spec.whatwg.org/#concept-node-adopt]")
		}
		// S3-3.
		for _, inclusive_descendant := range dom_node_shadow_including_descendants(node) {
			inclusive_descendant.run_adopting_steps(old_document)
		}
	}
}

func dom_print_tree(node dom_Node) {
	curr_node := node
	count := 0
	if !cm.IsNil(curr_node.get_parent()) {
		for n := curr_node.get_parent(); !cm.IsNil(n); n = n.get_parent() {
			count += 4
		}
	}
	indent := strings.Repeat(" ", count)
	fmt.Printf("%s%v\n", indent, node)
	for _, child := range curr_node.get_children() {
		dom_print_tree(child)
	}
}

type dom_TestNode struct {
	dom_Node_s
	name string
}

func (n dom_TestNode) String() string {
	return fmt.Sprintf("TestNode %s", n.name)
}
