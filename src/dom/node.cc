// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "node.hh"
#include "../_utility/error.hh"
#include "../_utility/logging.hh"
#include "../idl/domexception.hh"
#include "element.hh"
#include "range.hh"
#include "shadowroot.hh"

#include <cassert>
#include <cmath>
#include <cstddef>
#include <memory>
#include <string>
#include <utility>
#include <vector>

namespace yw::dom {

Node::Constructor_Badge Node::_constructor_badge()
{
    return {};
}

Node::Node([[maybe_unused]] Constructor_Badge badge, std::string debug_name,
    Type type, std::shared_ptr<Document> const& node_document)
    : m_document(node_document)
    , m_node_type(type)
    , m_debug_name(std::move(debug_name))
{
}

template <typename T> std::shared_ptr<T> Node::_root(T& self)
{
    std::shared_ptr<T> current_node = self.shared_from_this();
    while (true) {
        if (!current_node->parent_node()) {
            return current_node;
        }
        current_node = current_node->parent_node();
    }
}

std::shared_ptr<Node> Node::root()
{
    return _root(*this);
}

std::shared_ptr<Node const> Node::root() const
{
    return _root(*this);
}

bool Node::is_descendant_of(std::shared_ptr<Node const> const& of) const
{
    std::shared_ptr<Node const> current = shared_from_this();
    while (current != of) {
        if (!current->parent_node()) {
            return false;
        }
        current = current->parent_node();
    }
    return true;
}

bool Node::is_ancestor_of(std::shared_ptr<Node const> const& of) const
{
    return of->is_descendant_of(shared_from_this());
}

bool Node::is_inclusive_descendant_of(std::shared_ptr<Node const> const& of)
{
    return (of == shared_from_this()) || is_descendant_of(of);
}

bool Node::is_inclusive_ancestor_of(std::shared_ptr<Node const> const& of)
{
    return (of == shared_from_this()) || is_ancestor_of(of);
}

bool Node::is_connected() const
{
    std::shared_ptr<Node const> root = shadow_including_root();
    if (!root) {
        return false;
    }
    return root->node_type() == Type::DOCUMENT;
}

[[nodiscard]] std::shared_ptr<Document> Node::node_document()
{
    return m_document.lock();
}

[[nodiscard]] std::shared_ptr<Document> Node::node_document() const
{
    return m_document.lock();
}

std::shared_ptr<Node> Node::preceding() const
{
    return m_preceding.lock();
}

std::shared_ptr<Node> Node::following() const
{
    return m_following;
}

[[nodiscard]] size_t Node::index() const
{
    size_t index = 0;
    std::shared_ptr<Node const> current = shared_from_this();
    while (current->previous_sibling()) {
        current = current->previous_sibling();
        index++;
    }
    return index;
}

template <typename T>
[[nodiscard]] std::shared_ptr<T> Node::_shadow_including_root(T& self)
{
    std::shared_ptr<T> current_node = self.shared_from_this();
    while (true) {
        using S = Shadow_Root_For_Node<T>;
        std::shared_ptr<S> shadowroot
            = std::dynamic_pointer_cast<S>(current_node);
        if (shadowroot) {
            current_node = shadowroot->host()->shadow_including_root();
        } else {
            return current_node->root();
        }
    }
}

[[nodiscard]] std::shared_ptr<Node> Node::shadow_including_root()
{
    return _shadow_including_root(*this);
}

[[nodiscard]] std::shared_ptr<Node const> Node::shadow_including_root() const
{
    return _shadow_including_root(*this);
}

[[nodiscard]] bool Node::host_including_inclusive_ancestor_of(
    std::shared_ptr<Node const> const& of) const
{
    std::shared_ptr<Node const> current_of = of;
    while (true) {
        if (is_ancestor_of(current_of)) {
            return true;
        }
        std::shared_ptr<Shadow_Root const> root
            = std::dynamic_pointer_cast<Shadow_Root const>(current_of->root());
        if (root && root->host()) {
            current_of = root->host();
            continue;
        }
        return false;
    }
}

void Node::insert(std::shared_ptr<Node> const& parent,
    std::shared_ptr<Node> const& before_child, bool suppress_observers)
{
    // https://dom.spec.whatwg.org/#concept-node-insert

    // 1. Let nodes be node’s children, if node is a DocumentFragment node;
    // otherwise « node ».
    std::vector<std::shared_ptr<Node>> nodes
        = (node_type() == Type::DOCUMENT_FRAGMENT)
        ? child_nodes()
        : std::vector { shared_from_this() };

    // 2. Let count be nodes’s size.
    size_t count = nodes.size();

    // 3. If count is 0, then return.
    if (count == 0) {
        return;
    }

    // 4. If node is a DocumentFragment node, then:
    if (node_type() == Type::DOCUMENT_FRAGMENT) {
        assert(!"TODO");
        // 1. Remove its children with the suppress observers flag set.
        // 2. Queue a tree mutation record for node with « », nodes, null, and
        // null.
    }

    // 5. If child is non-null, then:
    if (before_child) {
        if (Range::_live_ranges().size() != 0) {
            assert(!"TODO");
        }
        // 1. For each live range whose start node is parent and start offset is
        // greater than child’s index, increase its start offset by count.

        // 2. For each live range whose end node is parent and end offset is
        // greater than child’s index, increase its end offset by count.
    }

    // 6. Let previousSibling be child’s previous sibling or parent’s last child
    // if child is null.
    std::shared_ptr<Node> previous_sibling = !before_child
        ? parent->last_child()
        : before_child->previous_sibling();

    // 7. For each node in nodes, in tree order:
    for (std::shared_ptr<Node> node : nodes) {

        // 1. Adopt node into parent’s node document.
        node->adopt_into(parent->node_document());

        // 2. If child is null, then append node to parent’s children.
        if (!before_child) {
            parent->_append_child(node);
        }
        // 3. Otherwise, insert node into parent’s children before child’s
        // index.
        else {
            parent->_insert_child_before(node, before_child);
        }

        // 4. If parent is a shadow host whose shadow root’s slot assignment is
        // "named" and node is a slottable, then assign a slot for node.
        if (parent->node_type() == Type::ELEMENT) {
            std::shared_ptr<Element> parent_elem
                = std::dynamic_pointer_cast<Element>(parent);
            assert(parent_elem);
            if (parent_elem->is_shadow_host()) {
                assert(!"TODO");
            }
        }

        // 5. If parent’s root is a shadow root, and parent is a slot whose
        // assigned nodes is the empty list, then run signal a slot change for
        // parent.
        if (parent->root()->node_type() == Type::ELEMENT) {
            std::shared_ptr<Shadow_Root> parent_sroot
                = std::dynamic_pointer_cast<Shadow_Root>(parent);
            if (parent_sroot) {
                assert(!"TODO");
            }
        }

        // 6. Run assign slottables for a tree with node’s root.
        root()->assign_slottables_for_a_tree();

        // 7. For each shadow-including inclusive descendant inclusiveDescendant
        // of node, in shadow-including tree order:
        shadow_including_inclusive_descendants(
            [&](std::shared_ptr<Node> const& inclusiveDescendant) {
                // 1. Run the insertion steps with inclusiveDescendant.
                inclusiveDescendant->run_insertion_steps();

                // 2. If inclusiveDescendant is connected, then:
                if (inclusiveDescendant->is_connected()) {

                    std::shared_ptr<Element> elem
                        = std::dynamic_pointer_cast<Element>(
                            inclusiveDescendant);
                    // 1. If inclusiveDescendant is custom, then enqueue a
                    // custom element callback reaction with
                    // inclusiveDescendant, callback name "connectedCallback",
                    // and an empty argument list.
                    if (elem && elem->is_custom()) {
                        assert(!"TODO");
                    }
                    // 2. Otherwise, try to upgrade inclusiveDescendant.
                    else {
                        LOG_TODO << "Try to upgrade connected node '"
                                 << inclusiveDescendant->_debug_name() << "'\n";
                    }
                }
                return true;
            });
    }

    // 8. If suppress observers flag is unset, then queue a tree mutation record
    // for parent with nodes, « », previousSibling, and child.
    if (!suppress_observers) {
        LOG_TODO << "Queue a tree mutation record for parent with nodes, « », "
                    "previousSibling, and child.\n";
    }

    // 9. Run the children changed steps for parent.
    parent->run_child_changed_steps();

    // 10. Let staticNodeList be a list of nodes, initially « ».
    std::vector<std::shared_ptr<Node>> static_node_list;

    // 11. For each node of nodes, in tree order:
    for (std::shared_ptr<Node> node : nodes) {
        // 1. For each shadow-including inclusive descendant inclusiveDescendant
        // of node, in shadow-including tree order, append inclusiveDescendant
        // to staticNodeList.
        shadow_including_inclusive_descendants(
            [&](std::shared_ptr<Node> const& inclusiveDescendant) {
                static_node_list.push_back(inclusiveDescendant);
                return true;
            });
    }

    // 12. For each node of staticNodeList, if node is connected, then run the
    // post-connection steps with node.
    for (const std::shared_ptr<Node>& node : static_node_list) {
        if (node->is_connected()) {
            node->run_post_connection_steps();
        }
    }
}

Error<idl::DOM_Exception> Node::ensure_pre_insertion_validity(
    std::shared_ptr<Node const> const& parent,
    std::shared_ptr<Node const> const& before_child) const
{
    // https://dom.spec.whatwg.org/#concept-node-ensure-pre-insertion-validity

    // 1. If parent is not a Document, DocumentFragment, or Element node, then
    // throw a "HierarchyRequestError" DOMException.
    switch (parent->node_type()) {
    case Type::DOCUMENT:
    case Type::DOCUMENT_FRAGMENT:
    case Type::ELEMENT:
        break;
    default:
        return DOM_EXCEPTION("", idl::DOM_Exception::HIERARCHY_REQUEST_ERROR);
    }

    // 2. If node is a host-including inclusive ancestor of parent, then throw a
    // "HierarchyRequestError" DOMException.
    if (host_including_inclusive_ancestor_of(parent)) {
        return DOM_EXCEPTION("", idl::DOM_Exception::HIERARCHY_REQUEST_ERROR);
    }

    // 3. If child is non-null and its parent is not parent, then throw a
    // "NotFoundError" DOMException.
    if (before_child && before_child->parent_node() != parent) {
        return DOM_EXCEPTION("", idl::DOM_Exception::NOT_FOUND_ERROR);
    }

    // 4. If node is not a DocumentFragment, DocumentType, Element, or
    // CharacterData node, then throw a "HierarchyRequestError" DOMException.
    switch (node_type()) {
    case Type::DOCUMENT_FRAGMENT:
    case Type::DOCUMENT_TYPE:
    case Type::ELEMENT:
    // Beginning of CharacterData nodes
    case Type::TEXT:
    case Type::PROCESSING_INSTRUCTION:
    case Type::COMMENT:
        break;
    // End of CharacterData nodes
    default:
        return DOM_EXCEPTION("", idl::DOM_Exception::HIERARCHY_REQUEST_ERROR);
    }

    // 5. If either node is a Text node and parent is a document, or node is a
    // doctype and parent is not a document, then throw a
    // "HierarchyRequestError" DOMException.
    if (((node_type() == Type::TEXT) && (parent->node_type() == Type::DOCUMENT))
        || ((node_type() == Type::DOCUMENT_TYPE)
            && (parent->node_type() != Type::DOCUMENT))) {
        return DOM_EXCEPTION("", idl::DOM_Exception::HIERARCHY_REQUEST_ERROR);
    }

    // 6. If parent is a document, and any of the statements below, switched on
    // the interface node implements, are true, then throw a
    // "HierarchyRequestError"
    if (parent->node_type() == Type::DOCUMENT) {
        bool die = false;
        std::shared_ptr<Node const> node = shared_from_this();

        switch (node_type()) {
        case Type::DOCUMENT_FRAGMENT:
            // If node has more than one element child or has a Text node child.
            if ((0 < _child_count_for(Type::ELEMENT))
                || (_child_count_for(Type::TEXT) != 0)) {
                die = true;
            }
            // Otherwise, if node has one element child and either parent has an
            // element child, child is a doctype, or child is non-null and a
            // doctype is following child.
            else if (_child_count_for(Type::ELEMENT) == 1) {
                if ((parent->_child_count_for(Type::ELEMENT) != 0)
                    || (before_child
                        && (before_child->node_type() == Type::DOCUMENT_TYPE))
                    || (before_child
                        && before_child->_is_followed_by(
                            Type::DOCUMENT_TYPE))) {
                    die = true;
                }
            }
            break;
        case Type::ELEMENT:
            if ((parent->_child_count_for(Type::ELEMENT) != 0)
                || (before_child
                    && (before_child->node_type() == Type::DOCUMENT_TYPE))
                || (before_child
                    && before_child->_is_followed_by(Type::DOCUMENT_TYPE))) {
                die = true;
            }
        case Type::DOCUMENT_TYPE:
            if ((parent->_child_count_for(Type::DOCUMENT_TYPE) != 0)
                || (before_child
                    && before_child->_is_preceded_by(Type::ELEMENT))
                || (!before_child
                    && (parent->_child_count_for(Type::ELEMENT) != 0))) {
                die = true;
            }
            break;
        default:
            break;
        }
        if (die) {
            return DOM_EXCEPTION(
                "", idl::DOM_Exception::HIERARCHY_REQUEST_ERROR);
        }
    }

    return {};
}

Error<idl::DOM_Exception> Node::pre_insert(std::shared_ptr<Node> const& parent,
    std::shared_ptr<Node> const& before_child)
{
    // https://dom.spec.whatwg.org/#concept-node-pre-insert

    // 1. Ensure pre-insertion validity of node into parent before child.
    Error<idl::DOM_Exception> error
        = ensure_pre_insertion_validity(parent, before_child);

    // 2. Let referenceChild be child.
    std::shared_ptr<Node> reference_child = before_child;

    // 3. If referenceChild is node, then set referenceChild to node’s next
    // sibling.
    if (reference_child == shared_from_this()) {
        reference_child = next_sibling();
    }

    // 4. Insert node into parent before referenceChild.
    insert(parent, reference_child);

    // 5. Return node.
    return {};
}

Error<idl::DOM_Exception> Node::append(std::shared_ptr<Node> const& parent)
{
    return pre_insert(parent, {});
}

void Node::assign_slottables_for_a_tree()
{
    inclusive_descendants([&](std::shared_ptr<Node> const& node) {
        std::shared_ptr<Element> element
            = std::dynamic_pointer_cast<Element>(node);
        if (!element) {
            return true;
        }
        if (element->tag_name() == "SLOT") {
            assert(!"TODO");
        }
        return true;
    });
}

void Node::adopt_into(std::shared_ptr<Document> const& document)
{
    // https://dom.spec.whatwg.org/#concept-node-adopt

    // 1. Let oldDocument be node’s node document.
    std::shared_ptr<Document> old_document = node_document();

    // 2. If node’s parent is non-null, then remove node.
    if (parent_node()) {
        assert(!"TODO");
    }

    // 3. If document is not oldDocument, then:
    if (document != old_document) {
        // 1. For each inclusiveDescendant in node’s shadow-including inclusive
        // descendants:
        shadow_including_inclusive_descendants(
            [&](std::shared_ptr<Node> const& inclusiveDescendant) {
                // 1. Set inclusiveDescendant’s node document to document.
                inclusiveDescendant->m_document = document;

                // 2.  If inclusiveDescendant is an element, then set the node
                // document of each attribute in inclusiveDescendant’s attribute
                // list to document.
                if (inclusiveDescendant->node_type() == Type::ELEMENT) {
                    assert(!"TODO");
                }
                return true;
            });

        // 2. For each inclusiveDescendant in node’s shadow-including inclusive
        // descendants that is custom, enqueue a custom element callback
        // reaction with inclusiveDescendant, callback name "adoptedCallback",
        // and an argument list containing oldDocument and document.
        shadow_including_inclusive_descendants(
            [&](std::shared_ptr<Node> const& inclusiveDescendant) {
                if (inclusiveDescendant->node_type() == Type::ELEMENT) {
                    assert(!"TODO");
                }
                return true;
            });
        // 3. For each inclusiveDescendant in node’s shadow-including inclusive
        // descendants, in shadow-including tree order, run the adopting steps
        // with inclusiveDescendant and oldDocument.
        shadow_including_inclusive_descendants(
            [&](std::shared_ptr<Node> const& inclusiveDescendant) {
                inclusiveDescendant->run_adopting_steps(old_document);
                return true;
            });
    }
}

void Node::run_insertion_steps()
{
}

void Node::run_adopting_steps(
    [[maybe_unused]] std::shared_ptr<Document> const& old_document)
{
}

void Node::run_child_changed_steps()
{
}

void Node::run_post_connection_steps()
{
}

//--------------------------------------------------------------------------

[[nodiscard]] Node::Type Node::node_type() const
{
    return m_node_type;
}

std::shared_ptr<Node> Node::parent_node() const
{
    return m_parent.lock();
}

std::shared_ptr<Node> Node::parent_element() const
{
    if (parent_node()->node_type() == Type::ELEMENT) {
        return parent_node();
    }
    return {};
}

bool Node::has_child_nodes() const
{
    return !!m_first_child.lock();
}

std::vector<std::shared_ptr<Node>> Node::child_nodes() const
{
    std::vector<std::shared_ptr<Node>> result;
    std::shared_ptr<Node> current = m_first_child.lock();
    while (current) {
        result.push_back(current);
        current = current->next_sibling();
    }
    return result;
}

std::shared_ptr<Node> Node::first_child() const
{
    return m_first_child.lock();
}

std::shared_ptr<Node> Node::last_child() const
{
    return m_last_child.lock();
}

std::shared_ptr<Node> Node::previous_sibling() const
{
    return m_previous_sibling.lock();
}

std::shared_ptr<Node> Node::next_sibling() const
{
    return m_next_sibling.lock();
}

//--------------------------------------------------------------------------

void Node::set_node_document(std::shared_ptr<Document> const& document)
{
    m_document = document;
}

std::shared_ptr<Node> Node::_create(std::string debug_name, Type type,
    std::shared_ptr<Document> const& node_document)
{
    std::shared_ptr<Node> node = std::make_shared<Node>(
        Constructor_Badge {}, debug_name, type, node_document);
    return node;
}

std::string Node::_debug_name() const
{
    return m_debug_name;
}

void Node::_insert_child_before(
    std::shared_ptr<Node> const& node, std::shared_ptr<Node> const& before)
{
    std::shared_ptr<Node> preceding = before->preceding();
    std::shared_ptr<Node> previous_sibling = before->previous_sibling();
    const std::shared_ptr<Node>& new_following = before;
    std::shared_ptr<Node> last_reachable_following = node->_last_node_in_tree();

    assert(!node->parent_node());

    // Set parent
    node->m_parent = shared_from_this();

    // Set previous/next sibling
    node->m_next_sibling = new_following;
    node->m_previous_sibling = new_following->m_previous_sibling;
    new_following->m_previous_sibling = node;
    if (previous_sibling) {
        previous_sibling->m_next_sibling = node;
    }

    // Set preceding and following
    // (Number between () is the order)
    //
    // Example:
    //                 [this(1)]
    //                 |   |   |
    //      +----------+   |   +--------+
    //      |              | node       |  before(new_following)
    // [child(2)]  +-> [child(4)]   [child(7)]
    //      |      |    |            ^
    //      |      |    +-------+    +--------+
    //      |      |            |             |
    // [child(3)] -+   [child(5)] [child(6)] -+
    //  preceding                 last_reachable_following
    //
    // NOTE: Lines coming from bottom     -> Child connection
    //       Lines coming from right side -> Following node
    preceding->m_following = node;
    node->m_preceding = preceding;
    last_reachable_following->m_following = before;
    before->m_preceding = last_reachable_following;

    // Set first child if needed
    if (!node->m_previous_sibling.lock()) {
        m_first_child = node;
        m_following = node;
    }
}

void Node::_append_child(std::shared_ptr<Node>& node)
{
    std::shared_ptr<Node> prev_child = m_last_child.lock();
    std::shared_ptr<Node> next_sibling = m_next_sibling.lock();
    std::shared_ptr<Node> preceding;
    if (prev_child) {
        preceding = prev_child->_last_node_in_tree();
        assert(shared_from_this() != preceding);
    } else {
        preceding = shared_from_this();
    }
    std::shared_ptr<Node> last_reachable_following = node->_last_node_in_tree();

    assert(!node->parent_node());

    // Set parent
    node->m_parent = shared_from_this();

    // Set previous/next sibling
    node->m_next_sibling = {};
    node->m_previous_sibling = prev_child;
    if (m_last_child.lock()) {
        m_last_child.lock()->m_next_sibling = node;
    }

    // Set preceding and following
    // (Number between () is the order)
    //
    // Example:
    //               [this(1)] --> [next sibling(7)]
    //                 |   |                  |
    //      +----------+   |                  |
    //      |              | node             |
    // [child(2)]  +-> [child(4)]             |
    //      |      |    |                     |
    //      |      |    +-------+             |
    //      |      |            |             |
    // [child(3)] -+   [child(5)] [child(6)] -+
    //  preceding                last_reachable_following
    //
    // NOTE: Lines coming from bottom     -> Child connection
    //       Lines coming from right side -> Following node
    //
    // When child is the first child:
    //  [this(1)] preceding  --> [next sibling(7)]
    //        |                      |
    //        | node                 |
    //   [child(4)]                  |
    //     |                         |
    //     +-------+                 |
    //             |                 |
    //    [child(5)] [child(6)] -----+
    //              last_reachable_following
    preceding->m_following = node;
    node->m_preceding = preceding;
    last_reachable_following->m_following = next_sibling;
    if (next_sibling) {
        next_sibling->m_preceding = last_reachable_following;
    }
    // Set first/last child
    if (!node->m_previous_sibling.lock()) {
        m_first_child = node;
        m_following = node;
    }
    m_last_child = node;
}

void Node::_remove_from_parent()
{
    std::shared_ptr<Node> parent = m_parent.lock();
    std::shared_ptr<Node> preceding = m_preceding.lock();
    std::shared_ptr<Node> following = _last_node_in_tree()->following();
    std::shared_ptr<Node> prev_sibling = m_previous_sibling.lock();
    std::shared_ptr<Node> next_sibling = m_next_sibling.lock();

    // Set parent
    m_parent = {};

    // Set previous/next sibling
    if (prev_sibling) {
        prev_sibling->m_next_sibling = next_sibling;
    }
    if (next_sibling) {
        next_sibling->m_previous_sibling = prev_sibling;
    }

    // Set first/last child
    if (!prev_sibling) {
        parent->m_first_child = next_sibling;
    }
    if (!next_sibling) {
        parent->m_last_child = prev_sibling;
    }

    // Set preceding and following
    // (Number between () is the order)
    //
    // Example:
    //                [parent(1)]
    //                 |   |   |
    //      +----------+   |   +-----------+
    //      |              | this          |
    // [child(2)] [child(Removed)] +-> [child(4)]
    //      | prev_sibling         |   next_sibling
    //      |                      |   following
    //      |                      |
    // [child(3)] -----------------+
    //  preceding
    //
    // NOTE: Lines coming from bottom     -> Child connection
    //       Lines coming from right side -> Following node
    if (preceding) {
        preceding->m_following = following;
    }
    if (following) {
        following->m_preceding = preceding;
    }
}

std::shared_ptr<Node> Node::_last_node_in_tree()
{
    std::shared_ptr<Node> current_node = shared_from_this();
    while (true) {
        if (!current_node->m_following
            || !is_ancestor_of(current_node->m_following)) {
            return current_node;
        }
        current_node = current_node->m_following;
    }
}

int Node::_child_count_for(Type tp) const
{
    int count = 0;
    for (std::shared_ptr<Node const> node : child_nodes()) {
        if (node->node_type() == tp) {
            count++;
        }
    }
    return count;
}

bool Node::_is_followed_by(Type tp) const
{
    std::shared_ptr<Node const> current = shared_from_this()->preceding();
    while (current) {
        if (current->node_type() == tp) {
            return true;
        }
        current = current->preceding();
    }
    return false;
}

bool Node::_is_preceded_by(Type tp) const
{
    std::shared_ptr<Node const> current = shared_from_this()->following();
    while (current) {
        if (current->node_type() == tp) {
            return true;
        }
        current = current->following();
    }
    return false;
}

} // namespace yw::dom
