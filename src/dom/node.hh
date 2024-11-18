// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "../_utility/error.hh"
#include "../_utility/noncopymove.hh"
#include "../idl/domexception.hh"
#include <cassert>
#include <cstddef>
#include <memory>
#include <string>
#include <type_traits>
#include <vector>

namespace yw::dom {

class Document;
class Node;
class Shadow_Root;

// https://dom.spec.whatwg.org/#trees
class Node : public std::enable_shared_from_this<Node> {
    YW_NON_COPYABLE(Node)
    YW_NON_MOVEABLE(Node)
    struct Constructor_Badge { };

    template <typename T>
    using Shadow_Root_For_Node = std::conditional_t<std::is_const_v<T>,
        Shadow_Root const, Shadow_Root>;

    // https://dom.spec.whatwg.org/#concept-shadow-including-descendant
    template <typename T, typename C>
    static bool _shadow_including_inclusive_descendants(T& node, C callback)
    {
        if (!callback(node.shared_from_this())) {
            return false;
        }
        using S = Shadow_Root_For_Node<T>;

        std::shared_ptr<S> shadowroot
            = std::dynamic_pointer_cast<S>(node.shared_from_this());
        if (shadowroot) {
            assert(!"TODO");
        }
        for (const std::shared_ptr<Node>& node : node.child_nodes()) {
            if (!_shadow_including_inclusive_descendants(*node, callback)) {
                return false;
            }
        }
        return true;
    }

    // https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
    template <typename T, typename C>
    static bool _shadow_including_descendants(T& node, C callback)
    {
        for (const std::shared_ptr<Node>& node : node.child_nodes()) {
            if (!_shadow_including_inclusive_descendants(*node, callback)) {
                return false;
            }
        }
        return true;
    }

    // https://dom.spec.whatwg.org/#concept-tree-descendant
    template <typename T, typename C>
    static bool _inclusive_descendants(T& node, C callback)
    {
        if (!callback(node.shared_from_this())) {
            return false;
        }
        for (const std::shared_ptr<T>& node : node.child_nodes()) {
            if (!_inclusive_descendants(*node, callback)) {
                return false;
            }
        }
        return true;
    }

    // https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
    template <typename T, typename C>
    static bool descendants(T& node, C callback)
    {
        for (const std::shared_ptr<T>& node : node.child_nodes()) {
            if (!_inclusive_descendants(*node, callback)) {
                return false;
            }
        }
        return true;
    }

public:
    enum class Type {
        // https://dom.spec.whatwg.org/#dom-node-element_node
        ELEMENT = 1,
        // https://dom.spec.whatwg.org/#dom-node-attribute_node
        ATTRIBUTE = 2,
        // https://dom.spec.whatwg.org/#dom-node-text_node
        TEXT = 3,
        // https://dom.spec.whatwg.org/#dom-node-cdata_section_node
        CDATA_SECTION = 4,
        // https://dom.spec.whatwg.org/#dom-node-processing_instruction_node
        PROCESSING_INSTRUCTION = 7,
        // https://dom.spec.whatwg.org/#dom-node-comment_node
        COMMENT = 8,
        // https://dom.spec.whatwg.org/#dom-node-document_node
        DOCUMENT = 9,
        // https://dom.spec.whatwg.org/#dom-node-document_type_node
        DOCUMENT_TYPE = 10,
        // https://dom.spec.whatwg.org/#dom-node-document_fragment_node
        DOCUMENT_FRAGMENT = 11,
    };

    Node(Constructor_Badge, std::string debug_name, Type type,
        std::shared_ptr<Document> const& node_document);
    virtual ~Node() = default;

    // https://dom.spec.whatwg.org/#concept-tree-root
    [[nodiscard]] std::shared_ptr<Node> root();
    [[nodiscard]] std::shared_ptr<Node const> root() const;

    // https://dom.spec.whatwg.org/#concept-tree-descendant
    [[nodiscard]] bool is_descendant_of(
        std::shared_ptr<Node const> const& of) const;

    // https://dom.spec.whatwg.org/#concept-tree-ancestor
    [[nodiscard]] bool is_ancestor_of(
        std::shared_ptr<Node const> const& of) const;

    // https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
    [[nodiscard]] bool is_inclusive_descendant_of(
        std::shared_ptr<Node const> const& of);

    // https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor
    [[nodiscard]] bool is_inclusive_ancestor_of(
        std::shared_ptr<Node const> const& of);

    // https://dom.spec.whatwg.org/#connected
    [[nodiscard]] bool is_connected() const;

    // https://dom.spec.whatwg.org/#concept-node-document
    [[nodiscard]] std::shared_ptr<Document> node_document();

    // https://dom.spec.whatwg.org/#concept-node-document
    [[nodiscard]] std::shared_ptr<Document> node_document() const;

    // https://dom.spec.whatwg.org/#concept-shadow-including-descendant
    template <typename C>
    bool shadow_including_inclusive_descendants(C callback) const
    {
        return _shadow_including_inclusive_descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-shadow-including-descendant
    template <typename C>
    bool shadow_including_inclusive_descendants(C callback)
    {
        return _shadow_including_inclusive_descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
    template <typename C> bool shadow_including_descendants(C callback) const
    {
        return _shadow_including_descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
    template <typename C> bool shadow_including_descendants(C callback)
    {
        return _shadow_including_descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-tree-descendant
    template <typename C> bool inclusive_descendants(C callback) const
    {
        return _inclusive_descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-tree-descendant
    template <typename C> bool inclusive_descendants(C callback)
    {
        return _inclusive_descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
    template <typename C> bool descendants(C callback) const
    {
        return _descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
    template <typename C> bool descendants(C callback)
    {
        return _descendants(*this, callback);
    }

    // https://dom.spec.whatwg.org/#concept-tree-preceding
    [[nodiscard]] std::shared_ptr<Node> preceding() const;

    // https://dom.spec.whatwg.org/#concept-tree-following
    [[nodiscard]] std::shared_ptr<Node> following() const;

    // https://dom.spec.whatwg.org/#concept-tree-index
    [[nodiscard]] size_t index() const;

    // https://dom.spec.whatwg.org/#concept-shadow-including-root
    [[nodiscard]] std::shared_ptr<Node> shadow_including_root();
    [[nodiscard]] std::shared_ptr<Node const> shadow_including_root() const;

    // https://dom.spec.whatwg.org/#concept-tree-host-including-inclusive-ancestor
    [[nodiscard]] bool host_including_inclusive_ancestor_of(
        std::shared_ptr<Node const> const& of) const;

    // https://dom.spec.whatwg.org/#concept-node-insert
    void insert(std::shared_ptr<Node> const& parent,
        std::shared_ptr<Node> const& before_child,
        bool suppress_observers = false);

    // https://dom.spec.whatwg.org/#concept-node-ensure-pre-insertion-validity
    Error<idl::DOM_Exception> ensure_pre_insertion_validity(
        std::shared_ptr<Node const> const& parent,
        std::shared_ptr<Node const> const& before_child) const;

    // https://dom.spec.whatwg.org/#concept-node-pre-insert
    Error<idl::DOM_Exception> pre_insert(std::shared_ptr<Node> const& parent,
        std::shared_ptr<Node> const& before_child);

    // https://dom.spec.whatwg.org/#concept-node-append
    Error<idl::DOM_Exception> append(std::shared_ptr<Node> const& parent);

    // https://dom.spec.whatwg.org/#assign-slotables-for-a-tree
    void assign_slottables_for_a_tree();

    // https://dom.spec.whatwg.org/#concept-node-adopt
    void adopt_into(std::shared_ptr<Document> const& document);

    // https://dom.spec.whatwg.org/#concept-node-insert-ext
    virtual void run_insertion_steps();

    // https://dom.spec.whatwg.org/#concept-node-adopt-ext
    virtual void run_adopting_steps(
        std::shared_ptr<Document> const& old_document);

    // https://dom.spec.whatwg.org/#concept-node-children-changed-ext
    virtual void run_child_changed_steps();

    // https://dom.spec.whatwg.org/#concept-node-post-connection-ext
    virtual void run_post_connection_steps();

    //--------------------------------------------------------------------------
    // Functions with associated MDN docs
    //--------------------------------------------------------------------------

    // https://dom.spec.whatwg.org/#dom-node-nodetype
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeType
    [[nodiscard]] Type node_type() const;

    // https://dom.spec.whatwg.org/#concept-tree-parent
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/parentNode
    [[nodiscard]] std::shared_ptr<Node> parent_node() const;

    // https://dom.spec.whatwg.org/#parent-element
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/parentElement
    [[nodiscard]] std::shared_ptr<Node> parent_element() const;

    // https://dom.spec.whatwg.org/#dom-node-haschildnodes
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/hasChildNodes
    [[nodiscard]] bool has_child_nodes() const;

    // https://dom.spec.whatwg.org/#dom-node-childnodes
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/childNodes
    [[nodiscard]] std::vector<std::shared_ptr<Node>> child_nodes() const;

    // https://dom.spec.whatwg.org/#concept-tree-first-child
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/firstChild
    [[nodiscard]] std::shared_ptr<Node> first_child() const;

    // https://dom.spec.whatwg.org/#concept-tree-last-child
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/lastChild
    [[nodiscard]] std::shared_ptr<Node> last_child() const;

    // https://dom.spec.whatwg.org/#concept-tree-previous-sibling
    [[nodiscard]] std::shared_ptr<Node> previous_sibling() const;

    // https://dom.spec.whatwg.org/#concept-tree-next-sibling
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Node/nextSibling
    [[nodiscard]] std::shared_ptr<Node> next_sibling() const;

    //--------------------------------------------------------------------------
    // YW Internal functions
    //--------------------------------------------------------------------------

    static std::shared_ptr<Node> _create(std::string debug_name, Type type, std::shared_ptr<Document> const &node_document);

    [[nodiscard]] std::string _debug_name() const;

    void _insert_child_before(
        std::shared_ptr<Node> const& node, std::shared_ptr<Node> const& before);

    void _append_child(std::shared_ptr<Node>& node);

    void _remove_from_parent();

protected:
    static Constructor_Badge _constructor_badge();

    void _set_node_document(std::shared_ptr<Document> const& document);

private:
    std::weak_ptr<Node> m_parent;

    std::weak_ptr<Node> m_preceding;
    std::shared_ptr<Node> m_following;

    std::weak_ptr<Node> m_first_child;
    std::weak_ptr<Node> m_last_child;

    std::weak_ptr<Node> m_previous_sibling;
    std::weak_ptr<Node> m_next_sibling;

    std::weak_ptr<Document> m_document;

    Type m_node_type;
    std::string m_debug_name;

    template <typename T> static std::shared_ptr<T> _root(T& self);

    template <typename T>
    static std::shared_ptr<T> _shadow_including_root(T& self);

    // Finds the very last node in this subtree, appearing in the tree order.
    std::shared_ptr<Node> _last_node_in_tree();

    int _child_count_for(Type tp) const;

    bool _is_followed_by(Type tp) const;
    bool _is_preceded_by(Type tp) const;
};

} // namespace yw::dom