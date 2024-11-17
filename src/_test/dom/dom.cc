// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "../tests.hh" // IWYU pragma: associated

#include "../../dom/node.hh"
#include "../../dom/shadowroot.hh"
#include "../testlib.hh"
#include <array>
#include <cstddef>
#include <iostream>
#include <memory>
#include <string>

namespace yw {

using yw::dom::Node;
using yw::dom::Shadow_Root;

namespace {

    [[nodiscard]] bool verify_node_parent_link(
        std::shared_ptr<Node> const& node, std::shared_ptr<Node> const& parent)
    {
        static constexpr char const* HEADER
            = "\x1b[35;1mBAD PARENT LINK\x1b[0m: ";

        bool ok = true;
        bool is_child = false;
        size_t expected_child_index = 0;
        for (std::shared_ptr<Node> child = parent->first_child(); child;
             child = child->next_sibling()) {
            if (child == node) {
                is_child = true;
                break;
            }
            expected_child_index++;
        }
        if (!is_child) {
            std::cerr << HEADER << "Node " << node->_debug_name()
                      << "'s parent " << parent->_debug_name()
                      << " does not include the node as a child\n";
            ok = false;
        } else if (expected_child_index != node->index()) {
            std::cerr
                << HEADER << "Node " << node->_debug_name() << "'s index "
                << node->index()
                << " does not match with the index we checked using parent "
                << parent->_debug_name() << "\n";
            ok = false;
        }
        // Check ancestor/descendant
        if (!node->is_descendant_of(parent)) {
            std::cerr << HEADER << "Node " << node->_debug_name()
                      << " is not descendant of its parent "
                      << parent->_debug_name() << "\n";
            ok = false;
        }
        if (!parent->is_ancestor_of(node)) {
            std::cerr << HEADER << "Parent " << parent->_debug_name()
                      << " is not"
                         " ancestor of its child "
                      << node->_debug_name() << "\n";
            ok = false;
        }
        return ok;
    }

    [[nodiscard]] bool verify_node_preceding_following_link(
        std::shared_ptr<Node> const& node)
    {
        static constexpr char const* HEADER
            = "\x1b[35;1mBAD PRECEDING/FOLLOWING LINK\x1b[0m: ";
        bool ok = true;
        std::shared_ptr<Node> const next = node->following();
        std::shared_ptr<Node> const prev = node->preceding();
        if (next && !next->preceding()) {
            std::cerr << HEADER << "Following node " << next->_debug_name()
                      << "'s"
                         " preceding node is not "
                      << node->_debug_name() << "(No preceding link)\n";
            ok = false;
        } else if (next && (next->preceding() != node)) {
            std::cerr << HEADER << "Following node " << next->_debug_name()
                      << "'s"
                         " preceding node is not "
                      << node->_debug_name() << "(Got "
                      << next->preceding()->_debug_name() << " instead)\n";
            ok = false;
        }
        if (prev && !prev->following()) {
            std::cerr << HEADER << "Preceding node " << prev->_debug_name()
                      << "'s"
                         " following node is not "
                      << node->_debug_name() << "(No following link)\n";
            ok = false;
        } else if (prev && (prev->following() != node)) {
            std::cerr << HEADER << "Preceding node " << prev->_debug_name()
                      << "'s"
                         " following node is not "
                      << node->_debug_name() << "(Got "
                      << prev->following()->_debug_name() << " instead)\n";
            ok = false;
        }
        return ok;
    }

    [[nodiscard]] bool verify_node_sibling_link(
        std::shared_ptr<Node> const& node)
    {
        static constexpr char const* HEADER
            = "\x1b[35;1mBAD SIBLING LINK\x1b[0m: ";
        bool ok = true;
        std::shared_ptr<Node> next = node->next_sibling();
        std::shared_ptr<Node> prev = node->previous_sibling();
        if (next && !next->previous_sibling()) {
            std::cerr << HEADER << "Next sibling node " << next->_debug_name()
                      << "'s"
                         " previous sibling node is not "
                      << node->_debug_name() << "(No previous sibling link)\n";
            ok = false;
        } else if (next && (next->previous_sibling() != node)) {
            std::cerr << HEADER << "Next sibling node " << next->_debug_name()
                      << "'s"
                         " previous sibling node is not "
                      << node->_debug_name() << "(Got "
                      << next->previous_sibling()->_debug_name()
                      << " instead)\n";
            ok = false;
        }
        if (prev && !prev->next_sibling()) {
            std::cerr << HEADER << "Previous sibling node "
                      << prev->_debug_name() << "'s next sibling node is not "
                      << node->_debug_name() << "(No next sibling link)\n";
            ok = false;
        } else if (prev && (prev->next_sibling() != node)) {
            std::cerr << HEADER << "Previous sibling node "
                      << prev->_debug_name() << "'s next sibling node is not "
                      << node->_debug_name() << "(Got "
                      << prev->next_sibling()->_debug_name() << " instead)\n";
            ok = false;
        }
        return ok;
    }

    [[nodiscard]] bool verify_node_link(std::shared_ptr<Node> const& node)
    {
        bool ok = true;
        std::shared_ptr<Node> const parent = node->parent_node();
        if (parent) {
            if (!verify_node_parent_link(node, parent)) {
                ok = false;
            }
        }
        if (!verify_node_sibling_link(node)) {
            ok = false;
        }
        if (!verify_node_preceding_following_link(node)) {
            ok = false;
        }

        return ok;
    }

    void test_create_node_internal()
    {
        std::shared_ptr<Node> const root
            = Node::_create("R", Node::Type::ELEMENT, {});
        TEST_ASSERT(!root->first_child());
        TEST_ASSERT(!root->last_child());
        TEST_ASSERT(!root->next_sibling());
        TEST_ASSERT(!root->previous_sibling());
        TEST_ASSERT(!root->preceding());
        TEST_ASSERT(!root->following());
        TEST_ASSERT(verify_node_link(root));
    }

    void test_append_child_internal()
    {
        std::shared_ptr<Node> const root
            = Node::_create("R", Node::Type::ELEMENT, {});
        //    [ R ] (1)
        //     |||
        //  +--+|+--+
        //  |   |   |
        // [N0][N1][N2]
        // (2) (3) (4)
        std::shared_ptr<Node> node0 = Node::_create("N0", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node1 = Node::_create("N1", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node2 = Node::_create("N2", Node::Type::ELEMENT, {});
        root->_append_child(node0);
        root->_append_child(node1);
        root->_append_child(node2);
        TEST_ASSERT(verify_node_link(node0));
        TEST_ASSERT(verify_node_link(node1));
        TEST_ASSERT(verify_node_link(node2));
        TEST_ASSERT(root->first_child() == node0);
        TEST_ASSERT(root->last_child() == node2);
        TEST_ASSERT(root->following() == node0);
        TEST_ASSERT(node0->following() == node1);
        TEST_ASSERT(node1->following() == node2);
        TEST_ASSERT(!node2->following());
    }

    void test_insert_child_before_internal()
    {
        std::shared_ptr<Node> root = Node::_create("R", Node::Type::ELEMENT, {});
        //    [ R ] (1)
        //     |||
        //  +--+|+--+
        //  |   |   |
        // [N2][N1][N0]
        // (2) (3) (4)
        std::shared_ptr<Node> node0 = Node::_create("N0", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node1 = Node::_create("N1", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node2 = Node::_create("N2", Node::Type::ELEMENT, {});
        root->_append_child(node0);
        // Note that we insert the node 2 first, then we insert
        // node 1 in between node 2 and node 0.
        root->_insert_child_before(node2, node0);
        root->_insert_child_before(node1, node0);
        TEST_ASSERT(verify_node_link(node0));
        TEST_ASSERT(verify_node_link(node1));
        TEST_ASSERT(verify_node_link(node2));
        TEST_ASSERT(root->first_child() == node2);
        TEST_ASSERT(root->last_child() == node0);
        TEST_ASSERT(root->following() == node2);
        TEST_ASSERT(node2->following() == node1);
        TEST_ASSERT(node1->following() == node0);
        TEST_ASSERT(!node0->following());
    }

    struct TestTree {
        std::shared_ptr<Node> root;
        std::array<std::shared_ptr<Node>, 3> nodes;
    };

    TestTree make_test_tree(char const* root_name, char const* child_prefix)
    {
        std::shared_ptr<Node> const root
            = Node::_create(root_name, Node::Type::ELEMENT, {});
        //    [ R ] (1)
        //     |||
        //  +--+|+--+
        //  |   |   |
        // [N0][N1][N2]
        // (2) (3) (4)
        std::string prefix(child_prefix);
        std::shared_ptr<Node> node0
            = Node::_create(prefix + "0", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node1
            = Node::_create(prefix + "1", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node2
            = Node::_create(prefix + "2", Node::Type::ELEMENT, {});
        root->_append_child(node0);
        root->_append_child(node1);
        root->_append_child(node2);
        return { root, { node0, node1, node2 } };
    }

    void test_append_to_child_internal()
    {
        //    [ R ] (1)
        //     |||
        //  +--+|+--+
        //  |   |   |
        // [N0][N1][N2]
        // (2) (3) (4)
        //  |
        // [N3]
        TestTree root_tree = make_test_tree("R", "N");
        std::shared_ptr<Node> node3 = Node::_create("N3", Node::Type::ELEMENT, {});
        root_tree.nodes[0]->_append_child(node3);
        TEST_ASSERT(verify_node_link(root_tree.nodes[0]));
        TEST_ASSERT(verify_node_link(node3));
        TEST_ASSERT(verify_node_link(root_tree.nodes[1]));
        TEST_ASSERT(root_tree.nodes[0]->following() == node3);
        TEST_ASSERT(node3->following() == root_tree.nodes[1]);
    }

    void test_append_tree_child_internal()
    {
        TestTree root_tree = make_test_tree("R", "N");
        TestTree sub_tree = make_test_tree("N3", "N3.");
        //    [  R  ] (1)
        //     ||||
        //     |||+-----+
        //  +--+|+--+   |
        //  |   |   |   |
        // [N0][N1][N2][N3]
        // (2) (3) (4) (5)
        //             |||
        // +-----------+||
        // |      +-----+|
        // |      |      |
        // [N3.0][N3.1][N3.2]
        // (6)    (7)   (8)
        root_tree.root->_append_child(sub_tree.root);
        TEST_ASSERT(verify_node_link(root_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(sub_tree.root));
        TEST_ASSERT(verify_node_link(root_tree.nodes[1]));
        TEST_ASSERT(root_tree.nodes[2]->following() == sub_tree.root);
        TEST_ASSERT(sub_tree.root->following() == sub_tree.nodes[0]);
        TEST_ASSERT(sub_tree.nodes[0]->following() == sub_tree.nodes[1]);
    }

    void test_insert_tree_child_before_internal()
    {
        TestTree root_tree = make_test_tree("R", "N");
        TestTree sub_tree = make_test_tree("N3", "N3.");

        //    [  R  ] (1)
        //     ||||
        //     |||+-----+
        //  +--+|+--+   |
        //  |   |   |   |
        // [N0][N3][N1][N2]
        // (2) (3) (7) (8)
        //     |||
        // +---+|+------+
        // |    ++      |
        // |     |      |
        // [N3.0][N3.1][N3.2]
        // (4)    (5)   (6)
        root_tree.root->_insert_child_before(sub_tree.root, root_tree.nodes[1]);
        TEST_ASSERT(verify_node_link(root_tree.nodes[0]));
        TEST_ASSERT(verify_node_link(sub_tree.root));
        TEST_ASSERT(verify_node_link(sub_tree.nodes[0]));
        TEST_ASSERT(verify_node_link(sub_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(root_tree.nodes[1]));
        TEST_ASSERT(root_tree.nodes[0]->following() == sub_tree.root);
        TEST_ASSERT(sub_tree.root->following() == sub_tree.nodes[0]);
        TEST_ASSERT(sub_tree.nodes[0]->following() == sub_tree.nodes[1]);
        TEST_ASSERT(sub_tree.nodes[2]->following() == root_tree.nodes[1]);
    }

    void test_remove_internal()
    {
        // Initial state:
        //    [  R  ] (1)
        //     ||||
        //     |||+-----+
        //  +--+|+--+   |
        //  |   |   |   |
        // [N0][N1][N2][N3]
        // (2) (3) (4) (5)
        //             |||
        // +-----------+||
        // |      +-----+|
        // |      |      |
        // [N3.0][N3.1][N3.2]
        // (6)    (7)   (8)

        TestTree root_tree = make_test_tree("R", "N");
        TestTree sub_tree_1 = make_test_tree("N3", "N3.");
        TestTree sub_tree_2 = make_test_tree("N4", "N4.");
        root_tree.root->_append_child(sub_tree_1.root);
        root_tree.root->_append_child(sub_tree_2.root);

        // Remove N3.1
        //    [  R  ] (1)
        //     ||||
        //     |||+-----------+
        //     |||+-----+     |
        //  +--+|+--+   |     |
        //  |   |   |   |     |
        // [N0][N1][N2][N3]  [N4]
        // (2) (3) (4) (5)   (8)
        //             | |   |||
        // +-----------+ |   ||+----------+
        // |             |   |+-----+     |
        // |             |   |      |     |
        // [N3.0]      [N3.2][N4.0][N4.1][N4.2]
        // (6)          (7)    (9)  (10)  (11)
        sub_tree_1.nodes[1]->_remove_from_parent();
        TEST_ASSERT(verify_node_link(root_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(sub_tree_1.root));
        TEST_ASSERT(verify_node_link(sub_tree_1.nodes[0]));
        TEST_ASSERT(verify_node_link(sub_tree_1.nodes[2]));
        TEST_ASSERT(root_tree.nodes[2]->following() == sub_tree_1.root);
        TEST_ASSERT(sub_tree_1.root->first_child() == sub_tree_1.nodes[0]);
        TEST_ASSERT(sub_tree_1.root->last_child() == sub_tree_1.nodes[2]);
        TEST_ASSERT(sub_tree_1.root->following() == sub_tree_1.nodes[0]);
        TEST_ASSERT(sub_tree_1.nodes[0]->following() == sub_tree_1.nodes[2]);
        TEST_ASSERT(sub_tree_1.nodes[2]->following() == sub_tree_2.root);

        // Remove N3.2
        //    [  R  ] (1)
        //     ||||
        //     |||+-----------+
        //     |||+-----+     |
        //  +--+|+--+   |     |
        //  |   |   |   |     |
        // [N0][N1][N2][N3]  [N4]
        // (2) (3) (4) (5)   (7)
        //             |     |||
        // +-----------+     ||+----------+
        // |                 |+-----+     |
        // |                 |      |     |
        // [N3.0]            [N4.0][N4.1][N4.2]
        // (6)               (8)   (9)   (10)

        sub_tree_1.nodes[2]->_remove_from_parent();
        TEST_ASSERT(verify_node_link(root_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(sub_tree_1.root));
        TEST_ASSERT(verify_node_link(sub_tree_1.nodes[0]));
        TEST_ASSERT(verify_node_link(sub_tree_2.root));
        TEST_ASSERT(root_tree.nodes[2]->following() == sub_tree_1.root);
        TEST_ASSERT(sub_tree_1.root->first_child() == sub_tree_1.nodes[0]);
        TEST_ASSERT(sub_tree_1.root->last_child() == sub_tree_1.nodes[0]);
        TEST_ASSERT(sub_tree_1.root->following() == sub_tree_1.nodes[0]);
        TEST_ASSERT(sub_tree_1.nodes[0]->following() == sub_tree_2.root);

        // Remove N3.0
        //    [  R  ] (1)
        //     ||||
        //     |||+-----------+
        //     |||+-----+     |
        //  +--+|+--+   |     |
        //  |   |   |   |     |
        // [N0][N1][N2][N3]  [N4]
        // (2) (3) (4) (5)   (6)
        //                   |||
        //                   ||+----------+
        //                   |+-----+     |
        //                   |      |     |
        //                   [N4.0][N4.1][N4.2]
        //                   (7)   (8)   (9)
        sub_tree_1.nodes[0]->_remove_from_parent();
        TEST_ASSERT(verify_node_link(root_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(sub_tree_1.root));
        TEST_ASSERT(verify_node_link(sub_tree_2.root));
        TEST_ASSERT(root_tree.nodes[2]->following() == sub_tree_1.root);
        TEST_ASSERT(!sub_tree_1.root->first_child());
        TEST_ASSERT(!sub_tree_1.root->last_child());
        TEST_ASSERT(sub_tree_1.root->following() == sub_tree_2.root);

        // Remove N4.2
        //    [  R  ] (1)
        //     ||||
        //     |||+-----------+
        //     |||+-----+     |
        //  +--+|+--+   |     |
        //  |   |   |   |     |
        // [N0][N1][N2][N3]  [N4]
        // (2) (3) (4) (5)   (6)
        //                   ||
        //                   ||
        //                   |+-----+
        //                   |      |
        //                   [N4.0][N4.1]
        //                   (7)   (8)
        sub_tree_2.nodes[2]->_remove_from_parent();
        TEST_ASSERT(verify_node_link(root_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(sub_tree_1.root));
        TEST_ASSERT(verify_node_link(sub_tree_2.root));
        TEST_ASSERT(verify_node_link(sub_tree_2.nodes[0]));
        TEST_ASSERT(verify_node_link(sub_tree_2.nodes[1]));
        TEST_ASSERT(root_tree.nodes[2]->following() == sub_tree_1.root);
        TEST_ASSERT(!sub_tree_1.root->first_child());
        TEST_ASSERT(!sub_tree_1.root->last_child());
        TEST_ASSERT(sub_tree_1.root->following() == sub_tree_2.root);
        TEST_ASSERT(sub_tree_2.root->following() == sub_tree_2.nodes[0]);
        TEST_ASSERT(sub_tree_2.nodes[0]->following() == sub_tree_2.nodes[1]);
        TEST_ASSERT(!sub_tree_2.nodes[1]->following());

        // Remove N4.0
        //    [  R  ] (1)
        //     ||||
        //     |||+-----------+
        //     |||+-----+     |
        //  +--+|+--+   |     |
        //  |   |   |   |     |
        // [N0][N1][N2][N3]  [N4]
        // (2) (3) (4) (5)   (6)
        //                    |
        //                    |
        //                    +-----+
        //                          |
        //                         [N4.1]
        //                          (7)
        sub_tree_2.nodes[0]->_remove_from_parent();
        TEST_ASSERT(verify_node_link(root_tree.nodes[2]));
        TEST_ASSERT(verify_node_link(sub_tree_1.root));
        TEST_ASSERT(verify_node_link(sub_tree_2.root));
        TEST_ASSERT(verify_node_link(sub_tree_2.nodes[1]));
        TEST_ASSERT(root_tree.nodes[2]->following() == sub_tree_1.root);
        TEST_ASSERT(!sub_tree_1.root->first_child());
        TEST_ASSERT(!sub_tree_1.root->last_child());
        TEST_ASSERT(sub_tree_1.root->following() == sub_tree_2.root);
        TEST_ASSERT(sub_tree_2.root->following() == sub_tree_2.nodes[1]);
        TEST_ASSERT(!sub_tree_2.nodes[1]->following());
    }

    void test_shadow_including_root()
    {
        std::shared_ptr<Node> root = Node::_create("R", Node::Type::ELEMENT, {});
        std::shared_ptr<Shadow_Root> sroot = Shadow_Root::_create("SR", {});
        sroot->_set_host(root);
        TEST_ASSERT(sroot->shadow_including_root() == root);
    }

    void test_parent_element()
    {
        std::shared_ptr<Node> elem = Node::_create("EP", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> doc = Node::_create("DP", Node::Type::DOCUMENT, {});
        std::shared_ptr<Node> elem_child
            = Node::_create("EPC", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> doc_child
            = Node::_create("DPC", Node::Type::ELEMENT, {});
        elem->_append_child(elem_child);
        doc->_append_child(doc_child);
        TEST_ASSERT(elem_child->parent_element() == elem);
        TEST_ASSERT(!doc_child->parent_element());
    }

    void test_host_including_inclusive_ancestor_of()
    {
        std::shared_ptr<Node> root = Node::_create("R", Node::Type::ELEMENT, {});
        std::shared_ptr<Node> node = Node::_create("N1", Node::Type::ELEMENT, {});
        root->_append_child(node);
        std::shared_ptr<Shadow_Root> sroot = Shadow_Root::_create("SR", {});
        sroot->_set_host(root);
        std::shared_ptr<Node> snode
            = Node::_create("SN1", Node::Type::ELEMENT, {});
        sroot->_append_child(snode);
        TEST_ASSERT(root->host_including_inclusive_ancestor_of(node));
        TEST_ASSERT(root->host_including_inclusive_ancestor_of(snode));
        TEST_ASSERT(!sroot->host_including_inclusive_ancestor_of(node));
        TEST_ASSERT(sroot->host_including_inclusive_ancestor_of(snode));
    }

} // namespace

void tests_init_dom(TestManager& tm)
{
    tm.register_test("_create_node", test_create_node_internal);
    tm.register_test("_append_child", test_append_child_internal);
    tm.register_test(
        "_append_child (to a child)", test_append_to_child_internal);
    tm.register_test("_append_child (tree)", test_append_tree_child_internal);
    tm.register_test("_insert_child_before", test_insert_child_before_internal);
    tm.register_test(
        "_insert_child_before (tree)", test_insert_tree_child_before_internal);
    tm.register_test("_remove", test_remove_internal);
    tm.register_test("_remove", test_remove_internal);
    tm.register_test("shadow_including_root", test_shadow_including_root);
    tm.register_test("parent_element", test_parent_element);
    tm.register_test("host_including_inclusive_ancestor_of",
        test_host_including_inclusive_ancestor_of);
}

} // namespace yw
