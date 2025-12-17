/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_common.h"
#include "yw_dom.h"
#include "yw_tests.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static YW_GcCallbacks yw_test_node_gc_callbacks = {
    .visit = yw_dom_node_visit,
    .destroy = yw_dom_node_destroy,
};

static YW_GC_PTR(YW_DomNode) yw_test_node_alloc(YW_GcHeap *heap)
{
    YW_GC_PTR(YW_DomNode) test_node = YW_GC_ALLOC(
        YW_DomNode, heap, &yw_test_node_gc_callbacks, YW_GC_ROOT_OBJECT);
    yw_dom_node_init(test_node);
    return test_node;
}

/*
 * For simplicity we don't set parent pointer of children nodes manually.
 * This function is used to set those parent pointers.
 */
static void yw_fix_children_parent(YW_GC_PTR(YW_DomNode) root)
{
    if (root == NULL)
    {
        return;
    }
    for (int i = 0; i < root->children.len; i++)
    {
        root->children.items[i]->parent = root;
        yw_fix_children_parent(root->children.items[i]);
    }
}

/*******************************************************************************
 * Node
 ******************************************************************************/

void yw_test_dom_first_child(YW_TestingContext *ctx)
{
    YW_GcHeap heap;

    yw_gc_heap_init(&heap);
    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);

    YW_TEST_EXPECT((void *)yw_dom_first_child(root), "%p", (void *)NULL);

    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child1);
    yw_fix_children_parent(root);
    YW_TEST_EXPECT((void *)yw_dom_first_child(root), "%p", (void *)child1);

    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child2);
    yw_fix_children_parent(root);

    YW_TEST_EXPECT((void *)yw_dom_first_child(root), "%p", (void *)child1);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_last_child(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);

    YW_TEST_EXPECT((void *)yw_dom_last_child(root), "%p", (void *)NULL);

    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child1);
    yw_fix_children_parent(root);
    YW_TEST_EXPECT((void *)yw_dom_last_child(root), "%p", (void *)child1);

    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child2);
    yw_fix_children_parent(root);
    YW_TEST_EXPECT((void *)yw_dom_last_child(root), "%p", (void *)child2);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_next_sibling(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child2);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child3);
    yw_fix_children_parent(root);

    YW_TEST_EXPECT((void *)yw_dom_next_sibling(root), "%p", (void *)NULL);
    YW_TEST_EXPECT((void *)yw_dom_next_sibling(child1), "%p", (void *)child2);
    YW_TEST_EXPECT((void *)yw_dom_next_sibling(child2), "%p", (void *)child3);
    YW_TEST_EXPECT((void *)yw_dom_next_sibling(child3), "%p", (void *)NULL);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_prev_sibling(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child2);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child3);
    yw_fix_children_parent(root);

    YW_TEST_EXPECT((void *)yw_dom_prev_sibling(root), "%p", (void *)NULL);
    YW_TEST_EXPECT((void *)yw_dom_prev_sibling(child1), "%p", (void *)NULL);
    YW_TEST_EXPECT((void *)yw_dom_prev_sibling(child2), "%p", (void *)child1);
    YW_TEST_EXPECT((void *)yw_dom_prev_sibling(child3), "%p", (void *)child2);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_root(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &child1->children, child2);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &child2->children, child3);
    yw_fix_children_parent(root);

    YW_TEST_EXPECT((void *)yw_dom_root(root, YW_DOM_NO_SEARCH_FLAGS), "%p",
                   (void *)root);
    YW_TEST_EXPECT((void *)yw_dom_root(child1, YW_DOM_NO_SEARCH_FLAGS), "%p",
                   (void *)root);
    YW_TEST_EXPECT((void *)yw_dom_root(child2, YW_DOM_NO_SEARCH_FLAGS), "%p",
                   (void *)root);
    YW_TEST_EXPECT((void *)yw_dom_root(child3, YW_DOM_NO_SEARCH_FLAGS), "%p",
                   (void *)root);

    /* TODO: Test shadow-including root */

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_index(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child2);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children, child3);
    yw_fix_children_parent(root);

    YW_TEST_EXPECT(yw_dom_index(root), "%d", 0);
    YW_TEST_EXPECT(yw_dom_index(child1), "%d", 0);
    YW_TEST_EXPECT(yw_dom_index(child2), "%d", 1);
    YW_TEST_EXPECT(yw_dom_index(child3), "%d", 2);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_has_type(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) node = yw_test_node_alloc(&heap);
    node->type_flags |= YW_DOM_SHADOW_ROOT_NODE;

    YW_TEST_EXPECT(yw_dom_has_type((void *)NULL, YW_DOM_TEXT_NODE), "%d",
                   false);
    YW_TEST_EXPECT(yw_dom_has_type(node, YW_DOM_SHADOW_ROOT_NODE), "%d", true);
    YW_TEST_EXPECT(yw_dom_has_type(node, YW_DOM_DOCUMENT_FRAGMENT_NODE), "%d",
                   true);
    YW_TEST_EXPECT(yw_dom_has_type(node, YW_DOM_ELEMENT_NODE), "%d", false);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_is_in_same_tree(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) root2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root1->children, child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root1->children, child2);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root2->children, child3);
    yw_fix_children_parent(root1);
    yw_fix_children_parent(root2);

    YW_TEST_EXPECT(yw_dom_is_in_same_tree(root1, root1), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_in_same_tree(root1, child1), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_in_same_tree(root1, child2), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_in_same_tree(root1, child3), "%d", false);
    YW_TEST_EXPECT(yw_dom_is_in_same_tree(root2, child3), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_in_same_tree(child1, child2), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_in_same_tree(child1, child3), "%d", false);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_is_connected(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_Document) root1 =
        yw_document_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    YW_GC_PTR(YW_DomNode) root2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root1->_node.children, child1);
    child1->node_document = root1;
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &child1->children, child2);
    child2->node_document = root1;
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root2->children, child3);
    yw_fix_children_parent((YW_GC_PTR(YW_DomNode))root1);
    yw_fix_children_parent((YW_GC_PTR(YW_DomNode))root2);

    YW_TEST_EXPECT(yw_dom_is_connected(root1), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_connected(root2), "%d", false);
    YW_TEST_EXPECT(yw_dom_is_connected(child1), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_connected(child2), "%d", true);
    YW_TEST_EXPECT(yw_dom_is_connected(child3), "%d", false);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_child_text(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_CharacterData) child1 =
        yw_text_alloc(&heap, YW_GC_ROOT_OBJECT);
    yw_append_str(&child1->text, "123");
    YW_GC_PTR(YW_CharacterData) child2 =
        yw_text_alloc(&heap, YW_GC_ROOT_OBJECT);
    yw_append_str(&child2->text, "abc");
    YW_GC_PTR(YW_CharacterData) child3 =
        yw_text_alloc(&heap, YW_GC_ROOT_OBJECT);
    yw_append_str(&child3->text, "789");
    yw_fix_children_parent(root);

    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children,
                 (YW_GC_PTR(YW_DomNode))child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children,
                 (YW_GC_PTR(YW_DomNode))child2);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->children,
                 (YW_GC_PTR(YW_DomNode))child3);

    YW_TEST_EXPECT_STR(yw_dom_child_text(root), "123abc789");

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_iter(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    /*
     * The test tree (nodeN means Nth item in DFS order):
     *                      node0
     *                        |
     *                +-------+------+
     *                |       |      |
     *              node1   node6  node7
     *                |              |
     *             +--+---+       +--+---+
     *             |      |       |      |
     *           node2  node5   node8  node11
     *             |              |
     *          +--+---+       +--+---+
     *          |      |       |      |
     *        node3  node4   node9  node10
     *
     */
    YW_GC_PTR(YW_DomNode) nodes[12];
    for (int i = 0; i < (int)(sizeof(nodes) / sizeof(void *)); i++)
    {
        nodes[i] = yw_test_node_alloc(&heap);
    }
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[0]->children, nodes[1]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[0]->children, nodes[6]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[0]->children, nodes[7]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[1]->children, nodes[2]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[1]->children, nodes[5]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[7]->children, nodes[8]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[7]->children, nodes[11]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[2]->children, nodes[3]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[2]->children, nodes[4]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[8]->children, nodes[9]);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &nodes[8]->children, nodes[10]);
    yw_fix_children_parent(nodes[0]);

    YW_DomIter iter;
    yw_dom_inclusive_descendants_init(&iter, (void *)nodes[0],
                                      YW_DOM_NO_SEARCH_FLAGS);

    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[0]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[1]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[2]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[3]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[4]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[5]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[6]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[7]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[8]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[9]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[10]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[11]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p", (void *)NULL);

    yw_dom_inclusive_descendants_init(&iter, nodes[1], YW_DOM_NO_SEARCH_FLAGS);

    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[1]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[2]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[3]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[4]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[5]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p", (void *)NULL);

    yw_dom_descendants_init(&iter, nodes[0], YW_DOM_NO_SEARCH_FLAGS);

    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[1]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[2]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[3]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[4]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[5]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[6]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[7]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[8]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[9]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[10]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p",
                   (void *)nodes[11]);
    YW_TEST_EXPECT((void *)yw_dom_next_descendant(&iter), "%p", (void *)NULL);

    yw_dom_inclusive_ancestors_init(&iter, nodes[11], YW_DOM_NO_SEARCH_FLAGS);

    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p",
                   (void *)nodes[11]);
    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p", (void *)nodes[7]);
    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p", (void *)nodes[0]);
    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p", (void *)NULL);

    yw_dom_ancestors_init(&iter, nodes[11], YW_DOM_NO_SEARCH_FLAGS);

    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p", (void *)nodes[7]);
    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p", (void *)nodes[0]);
    YW_TEST_EXPECT((void *)yw_dom_next_ancestor(&iter), "%p", (void *)NULL);
}

void yw_test_dom_insert(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) root = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child1 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child2 = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_DomNode) child3 = yw_test_node_alloc(&heap);
    yw_dom_insert(child1, root, NULL, YW_DOM_NO_INSERT_FLAGS);
    yw_dom_insert(child2, root, NULL, YW_DOM_NO_INSERT_FLAGS);
    yw_dom_insert(child3, root, child2, YW_DOM_NO_INSERT_FLAGS);

    /*
     * We don't call yw_fix_children_parent() here, because yw_dom_insert() is
     * supposed to do that.
     */

    YW_TEST_EXPECT((void *)root->children.items[0], "%p", (void *)child1);
    YW_TEST_EXPECT((void *)root->children.items[0]->parent, "%p", (void *)root);
    YW_TEST_EXPECT((void *)root->children.items[1], "%p", (void *)child3);
    YW_TEST_EXPECT((void *)root->children.items[1]->parent, "%p", (void *)root);
    YW_TEST_EXPECT((void *)root->children.items[2], "%p", (void *)child2);
    YW_TEST_EXPECT((void *)root->children.items[2]->parent, "%p", (void *)root);

    yw_gc_heap_deinit(&heap);
}

/*******************************************************************************
 * Element
 ******************************************************************************/

/* TODO: Test yw_dom_is_shadow_host */

void yw_test_dom_is_element_defined(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) non_elem = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_Element) elem = yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_UNDEFINED;
    YW_TEST_EXPECT(yw_dom_is_element_defined(non_elem), "%d", false);
    YW_TEST_EXPECT(yw_dom_is_element_defined(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_FAILED;
    YW_TEST_EXPECT(yw_dom_is_element_defined(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_UNCUSTOMIZED;
    YW_TEST_EXPECT(yw_dom_is_element_defined(elem), "%d", true);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_PRECUSTOMIZED;
    YW_TEST_EXPECT(yw_dom_is_element_defined(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_CUSTOM;
    YW_TEST_EXPECT(yw_dom_is_element_defined(elem), "%d", true);
    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_is_element_custom(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) non_elem = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_Element) elem = yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_UNCUSTOMIZED;
    YW_TEST_EXPECT(yw_dom_is_element_custom(non_elem), "%d", false);
    YW_TEST_EXPECT(yw_dom_is_element_custom(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_FAILED;
    YW_TEST_EXPECT(yw_dom_is_element_custom(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_UNCUSTOMIZED;
    YW_TEST_EXPECT(yw_dom_is_element_custom(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_PRECUSTOMIZED;
    YW_TEST_EXPECT(yw_dom_is_element_custom(elem), "%d", false);
    elem->custom_element_state = YW_CUSTOM_ELEMENT_CUSTOM;
    YW_TEST_EXPECT(yw_dom_is_element_custom(elem), "%d", true);
    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_is_element_inside(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) non_elem = yw_test_node_alloc(&heap);

    YW_GC_PTR(YW_Element) root = yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    root->local_name = "div";
    root->namespace_ = YW_HTML_NAMESPACE;

    YW_GC_PTR(YW_Element) child1 =
        yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    child1->local_name = "p";
    child1->namespace_ = YW_HTML_NAMESPACE;

    YW_GC_PTR(YW_Element) child2 =
        yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    child2->local_name = "span";
    child2->namespace_ = YW_HTML_NAMESPACE;

    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->_node.children,
                 (YW_GC_PTR(YW_DomNode))child1);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &root->_node.children, non_elem);
    YW_LIST_PUSH(YW_GC_PTR(YW_DomNode), &child1->_node.children,
                 (YW_GC_PTR(YW_DomNode))child2);

    yw_fix_children_parent(&root->_node);

    YW_TEST_EXPECT(yw_dom_is_element_inside(non_elem, YW_HTML_NAMESPACE, "p"),
                   "%d", false);
    YW_TEST_EXPECT(yw_dom_is_element_inside(child2, YW_HTML_NAMESPACE, "p"),
                   "%d", true);
    YW_TEST_EXPECT(yw_dom_is_element_inside(child2, YW_HTML_NAMESPACE, "div"),
                   "%d", true);
    YW_TEST_EXPECT(yw_dom_is_element_inside(child1, YW_HTML_NAMESPACE, "p"),
                   "%d", false);
    YW_TEST_EXPECT(yw_dom_is_element_inside(child1, YW_HTML_NAMESPACE, "div"),
                   "%d", true);
    YW_TEST_EXPECT(yw_dom_is_element_inside(child1, YW_SVG_NAMESPACE, "div"),
                   "%d", false);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_is_element(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_DomNode) non_elem = yw_test_node_alloc(&heap);
    YW_GC_PTR(YW_Element) html_elem =
        yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    html_elem->local_name = "p";
    html_elem->namespace_ = YW_HTML_NAMESPACE;

    YW_GC_PTR(YW_Element) mathml_elem =
        yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    mathml_elem->local_name = "mi";
    mathml_elem->namespace_ = YW_MATHML_NAMESPACE;

    YW_GC_PTR(YW_Element) svg_elem =
        yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    svg_elem->local_name = "g";
    svg_elem->namespace_ = YW_SVG_NAMESPACE;

    YW_TEST_EXPECT(yw_dom_is_element(non_elem, YW_HTML_NAMESPACE, "p"), "%d",
                   false);

    YW_TEST_EXPECT(yw_dom_is_element(html_elem, YW_HTML_NAMESPACE, "p"), "%d",
                   true);
    YW_TEST_EXPECT(yw_dom_is_element(html_elem, YW_HTML_NAMESPACE, "li"), "%d",
                   false);
    YW_TEST_EXPECT(yw_dom_is_element(mathml_elem, YW_HTML_NAMESPACE, "p"), "%d",
                   false);
    YW_TEST_EXPECT(yw_dom_is_element(svg_elem, YW_HTML_NAMESPACE, "p"), "%d",
                   false);

    YW_TEST_EXPECT(yw_dom_is_element(html_elem, YW_SVG_NAMESPACE, "g"), "%d",
                   false);
    YW_TEST_EXPECT(yw_dom_is_element(mathml_elem, YW_SVG_NAMESPACE, "g"), "%d",
                   false);
    YW_TEST_EXPECT(yw_dom_is_element(svg_elem, YW_SVG_NAMESPACE, "g"), "%d",
                   true);
    YW_TEST_EXPECT(yw_dom_is_element(svg_elem, YW_SVG_NAMESPACE, "line"), "%d",
                   false);

    YW_TEST_EXPECT(yw_dom_is_element(html_elem, YW_MATHML_NAMESPACE, "mi"),
                   "%d", false);
    YW_TEST_EXPECT(yw_dom_is_element(mathml_elem, YW_MATHML_NAMESPACE, "mi"),
                   "%d", true);
    YW_TEST_EXPECT(yw_dom_is_element(mathml_elem, YW_MATHML_NAMESPACE, "foo"),
                   "%d", false);
    YW_TEST_EXPECT(yw_dom_is_element(svg_elem, YW_MATHML_NAMESPACE, "mi"), "%d",
                   false);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_append_attr(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_Element) elem = yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);

    YW_AttrData data;

    memset(&data, 0, sizeof(data));
    data.local_name = "name1";
    data.value = "value1";
    data.namespace_ = "ns1";
    data.namespace_prefix = "prefix1";
    yw_dom_append_attr(elem, &heap, &data);

    memset(&data, 0, sizeof(data));
    data.local_name = "name2";
    data.value = "value2";
    yw_dom_append_attr(elem, &heap, &data);

    /*
     * We don't call yw_fix_children_parent() here, because yw_dom_append_attr()
     * is supposed to do that.
     */

    YW_TEST_EXPECT(elem->attrs.len, "%d", 2);

    YW_TEST_EXPECT((void *)elem->attrs.items[0]->_node.parent, "%p",
                   (void *)elem);
    YW_TEST_EXPECT_STR(elem->attrs.items[0]->local_name, "name1");
    YW_TEST_EXPECT_STR(elem->attrs.items[0]->value, "value1");
    YW_TEST_EXPECT_STR(elem->attrs.items[0]->namespace_, "ns1");
    YW_TEST_EXPECT_STR(elem->attrs.items[0]->namespace_prefix, "prefix1");
    YW_TEST_EXPECT((void *)elem->attrs.items[0]->element, "%p", (void *)elem);

    YW_TEST_EXPECT((void *)elem->attrs.items[1]->_node.parent, "%p",
                   (void *)elem);
    YW_TEST_EXPECT_STR(elem->attrs.items[1]->local_name, "name2");
    YW_TEST_EXPECT_STR(elem->attrs.items[1]->value, "value2");
    YW_TEST_EXPECT_STR(elem->attrs.items[1]->namespace_, (char *)NULL);
    YW_TEST_EXPECT_STR(elem->attrs.items[1]->namespace_prefix, (char *)NULL);
    YW_TEST_EXPECT((void *)elem->attrs.items[1]->element, "%p", (void *)elem);

    yw_gc_heap_deinit(&heap);
}

void yw_test_dom_attr(YW_TestingContext *ctx)
{
    YW_GcHeap heap;
    yw_gc_heap_init(&heap);

    YW_GC_PTR(YW_Element) elem = yw_element_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);

    YW_GC_PTR(YW_Attr) attr = yw_attr_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    attr->local_name = yw_duplicate_str("name1");
    attr->value = yw_duplicate_str("value1");
    attr->namespace_ = yw_duplicate_str("ns1");
    YW_LIST_PUSH(YW_GC_PTR(YW_Attr), &elem->attrs, attr);

    attr = yw_attr_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    attr->local_name = yw_duplicate_str("name2");
    attr->value = yw_duplicate_str("value2");
    attr->namespace_ = yw_duplicate_str("ns1");
    YW_LIST_PUSH(YW_GC_PTR(YW_Attr), &elem->attrs, attr);

    attr = yw_attr_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    attr->local_name = yw_duplicate_str("name3");
    attr->value = yw_duplicate_str("value3");
    attr->namespace_ = yw_duplicate_str("ns2");
    YW_LIST_PUSH(YW_GC_PTR(YW_Attr), &elem->attrs, attr);

    attr = yw_attr_alloc(&heap, YW_NO_GC_ALLOC_FLAGS);
    attr->local_name = yw_duplicate_str("name4");
    attr->value = yw_duplicate_str("value4");
    YW_LIST_PUSH(YW_GC_PTR(YW_Attr), &elem->attrs, attr);

    yw_fix_children_parent(&elem->_node);

    YW_TEST_EXPECT_STR(yw_attr(elem, "ns1", "name1"), "value1");
    YW_TEST_EXPECT_STR(yw_attr(elem, "ns1", "name2"), "value2");
    YW_TEST_EXPECT_STR(yw_attr(elem, "ns1", "name3"), (char *)NULL);
    YW_TEST_EXPECT_STR(yw_attr(elem, "ns2", "name3"), "value3");
    YW_TEST_EXPECT_STR(yw_attr(elem, "ns2", "name4"), (char *)NULL);

    YW_TEST_EXPECT_STR(yw_attr(elem, NULL, "name1"), (char *)NULL);
    YW_TEST_EXPECT_STR(yw_attr(elem, NULL, "name2"), (char *)NULL);
    YW_TEST_EXPECT_STR(yw_attr(elem, NULL, "name3"), (char *)NULL);
    YW_TEST_EXPECT_STR(yw_attr(elem, NULL, "name4"), "value4");

    yw_gc_heap_deinit(&heap);
}
