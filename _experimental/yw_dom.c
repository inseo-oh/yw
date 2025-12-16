/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_dom.h"
#include "yw_common.h"
#include "yw_namespaces.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*******************************************************************************
 * Node
 ******************************************************************************/

#define YW_DOM_NODE_MAGIC 0xb1fedf1b

static void yw_dom_node_check_magic(YW_GC_PTR(yw_dom_node) node)
{
    if (node == NULL)
    {
        return;
    }
    if (node->magic != YW_DOM_NODE_MAGIC)
    {
        fprintf(stderr, "%s: Node at %p has corrupted magic!\n", __func__,
                (void *)node);
        abort();
    }
}

void yw_dom_node_init(YW_GC_PTR(yw_dom_node) out)
{
    out->magic = YW_DOM_NODE_MAGIC;
}
void yw_dom_node_visit(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;

    for (int i = 0; i < node->children.len; i++)
    {
        yw_gc_visit(node->children.items[i]);
    }
    yw_gc_visit(node->parent);
    yw_gc_visit(node->node_document);
}
void yw_dom_node_destroy(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    YW_LIST_FREE(&node->children);
}

YW_GC_PTR(yw_dom_node) yw_dom_first_child(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);
    if (node->children.len == 0)
    {
        return NULL;
    }
    return node->children.items[0];
}

YW_GC_PTR(yw_dom_node) yw_dom_last_child(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);
    if (node->children.len == 0)
    {
        return NULL;
    }
    return node->children.items[node->children.len - 1];
}

YW_GC_PTR(yw_dom_node) yw_dom_next_sibling(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    YW_GC_PTR(yw_dom_node) parent = node->parent;
    if (parent == NULL)
    {
        return NULL;
    }
    int index = yw_dom_index(node_v);
    if (index == parent->children.len - 1)
    {
        return NULL;
    }
    return parent->children.items[index + 1];
}

YW_GC_PTR(yw_dom_node) yw_dom_prev_sibling(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    YW_GC_PTR(yw_dom_node) parent = node->parent;
    if (parent == NULL)
    {
        return NULL;
    }
    int index = yw_dom_index(node_v);
    if (index == 0)
    {
        return NULL;
    }
    return parent->children.items[index - 1];
}

YW_GC_PTR(yw_dom_node) yw_dom_root(void *node_v, yw_dom_search_flags flags)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    YW_GC_PTR(yw_dom_node) res = node;
    while (res->parent != NULL)
    {
        if ((flags & YW_DOM_SHADOW_INCLUDING) &&
            res->type == YW_DOM_SHADOW_ROOT_NODE)
        {
            /* TODO: Support shadow root */
            YW_TODO();
        }
        res = res->parent;
    }
    return res;
}

int yw_dom_index(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    YW_GC_PTR(yw_dom_node) parent = node->parent;
    if (parent == NULL)
    {
        return 0;
    }
    int idx = -1;
    for (int i = 0; i < parent->children.len; i++)
    {
        if (parent->children.items[i] == node)
        {
            idx = i;
        }
    }
    if (idx == -1)
    {
        fprintf(stderr, "%s: children could not be found", __func__);
        abort();
    }
    return idx;
}

bool yw_dom_is_in_the_same_tree_as(void *node_a_v, void *node_b_v)
{
    return yw_dom_root(node_a_v, YW_DOM_NO_SEARCH_FLAGS) ==
           yw_dom_root(node_b_v, YW_DOM_NO_SEARCH_FLAGS);
}

bool yw_dom_is_connected(void *node_v)
{
    /* https://dom.spec.whatwg.org/#connected */

    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    return yw_dom_root(node_v, YW_DOM_SHADOW_INCLUDING) ==
           (YW_GC_PTR(yw_dom_node))node->node_document;
}

bool yw_dom_is_in_document_tree(void *node_v)
{
    return yw_dom_root(node_v, YW_DOM_NO_SEARCH_FLAGS)->type ==
           YW_DOM_DOCUMENT_NODE;
}

/* Caller owns the returned string. */
char *yw_dom_child_text(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    char *res_buf = NULL;
    int res_cap = 0;
    int res_len = 0;

    for (int i = 0; i < node->children.len; i++)
    {
        if (node->children.items[i]->type == YW_DOM_TEXT_NODE)
        {
            char *node_text =
                ((YW_GC_PTR(yw_dom_character_data_node))node)->text;
            int len = strlen(node_text);
            res_buf = YW_GROW(char, &res_cap, &res_len, res_buf);
            res_buf[0] = '\0';
            strcat(res_buf, node_text);
        }
    }

    res_buf = YW_SHRINK_TO_FIT(char, &res_cap, res_len, res_buf);
    return res_buf;
}

YW_GC_PTR(yw_dom_node) yw_dom_next_descendant(yw_dom_iter *iter)
{
    /* https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant */

    YW_GC_PTR(yw_dom_node) curr_node = iter->last_node;
    YW_GC_PTR(yw_dom_node) res = NULL;

    if (curr_node == NULL)
    {
        /* This is our first call */
        res = iter->root_node;
    }
    else if (curr_node->children.len != 0)
    {
        /* Go to the first children */
        res = curr_node->children.items[0];
    }
    else
    {
        /* If we don't have more children, move to the next sibling. */
        while (res == NULL)
        {
            res = yw_dom_next_sibling(curr_node);
            if (res != NULL)
            {
                break;
            }
            /* We don't even have the next sibling -> Move to the parent */
            curr_node = curr_node->parent;
            if (curr_node == iter->root_node || curr_node == NULL)
            {
                /*
                 * We don't have parent, or we are currently at root.
                 * We have to stop here.
                 */
                res = NULL;
                break;
            }
        }
    }
    if (res == NULL)
    {
        return NULL;
    }
    if (iter->shadow_including && res->type == YW_DOM_SHADOW_ROOT_NODE)
    {
        /* TODO: Support shadow root */
        YW_TODO();
    }

    iter->last_node = res;
    return res;
}

void yw_dom_inclusive_descendants_init(yw_dom_iter *out, void *root_node_v,
                                       yw_dom_search_flags flags)
{
    YW_GC_PTR(yw_dom_node) root_node = (YW_GC_PTR(yw_dom_node))root_node_v;
    yw_dom_node_check_magic(root_node);

    memset(out, 0, sizeof(*out));
    out->root_node = root_node;
    if (flags & YW_DOM_SHADOW_INCLUDING)
    {
        out->shadow_including = true;
    }
}

void yw_dom_descendants_init(yw_dom_iter *out, void *root_node_v,
                             yw_dom_search_flags flags)
{
    yw_dom_inclusive_descendants_init(out, root_node_v, flags);
    yw_dom_next_descendant(out);
}

YW_GC_PTR(yw_dom_node) yw_dom_next_parent(yw_dom_iter *iter)
{
    /* https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor */

    YW_GC_PTR(yw_dom_node) curr_node = iter->last_node;
    YW_GC_PTR(yw_dom_node) res = NULL;

    if (curr_node == NULL)
    {
        /* This is our first call */
        res = iter->root_node;
    }
    else if (curr_node->parent != NULL)
    {
        res = curr_node->parent;
    }
    else if (curr_node->type == YW_DOM_SHADOW_ROOT_NODE)
    {
        /* TODO: Support shadow root */
        YW_TODO();
    }

    if (res == NULL)
    {
        return NULL;
    }
    iter->last_node = res;
    return res;
}

void yw_dom_inclusive_ancestors_init(yw_dom_iter *out, void *root_node_v,
                                     yw_dom_search_flags flags)
{
    YW_GC_PTR(yw_dom_node) root_node = (YW_GC_PTR(yw_dom_node))root_node_v;
    yw_dom_node_check_magic(root_node);

    memset(out, 0, sizeof(*out));
    out->root_node = root_node;
    if (flags & YW_DOM_SHADOW_INCLUDING)
    {
        out->shadow_including = true;
    }
}

void yw_dom_ancestors_init(yw_dom_iter *out, void *root_node_v,
                           yw_dom_search_flags flags)
{
    yw_dom_inclusive_ancestors_init(out, root_node_v, flags);
    yw_dom_next_descendant(out);
}

void yw_dom_insert(void *node_v, void *parent_v, void *before_child_v,
                   yw_dom_insert_flag flags)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    YW_GC_PTR(yw_dom_node) parent = (YW_GC_PTR(yw_dom_node))parent_v;
    YW_GC_PTR(yw_dom_node) before_child =
        (YW_GC_PTR(yw_dom_node))before_child_v;
    yw_dom_node_check_magic(node);
    yw_dom_node_check_magic(parent);
    yw_dom_node_check_magic(parent);

    YW_GC_PTR(yw_dom_node) prev_sibling;

    /* NOTE: All the step numbers(S#.) are based on spec from when this was
     * initially written(2025.11.13) */

    /* S1 *********************************************************************/
    yw_dom_node_list nodes;
    YW_LIST_INIT(&nodes);

    if (node->type == YW_DOM_DOCUMENT_FRAGMENT_NODE)
    {
        for (int i = 0; i < node->children.len; i++)
        {
            YW_LIST_PUSH(YW_GC_PTR(yw_dom_node), &nodes,
                         node->children.items[i]);
        }
    }
    else
    {
        YW_LIST_PUSH(YW_GC_PTR(yw_dom_node), &nodes, node);
    }

    /* S2 *********************************************************************/
    int count = nodes.len;

    /* S3 *********************************************************************/
    if (count == 0)
    {
        goto out;
    }

    /* S4 *********************************************************************/
    if (node->type == YW_DOM_DOCUMENT_FRAGMENT_NODE)
    {
        YW_TODO();
    }

    /* S5 *********************************************************************/
    if (before_child != NULL)
    {
        /* TODO */
    }

    /* S6 *********************************************************************/
    prev_sibling = yw_dom_last_child(parent);
    if (before_child != NULL)
    {
        prev_sibling = yw_dom_prev_sibling(before_child);
    }
    (void)prev_sibling;

    /* S7 *********************************************************************/
    for (int i = 0; i < nodes.len; i++)
    {
        YW_GC_PTR(yw_dom_node) node = nodes.items[i];
        /* S7-1 ***************************************************************/
        yw_dom_adopt_into(node, parent->node_document);

        if (before_child == NULL)
        {
            /* S7-2 ***********************************************************/
            YW_LIST_PUSH(YW_GC_PTR(yw_dom_node), &parent->children, node);
        }
        else
        {
            /* S7-3 ***********************************************************/
            int index = yw_dom_index(before_child);
            YW_LIST_INSERT(YW_GC_PTR(yw_dom_node), &parent->children, index,
                           node);
        }

        /* S7-4 ***************************************************************/
        if (yw_dom_is_shadow_host(parent))
        {
            YW_TODO();
        }

        /* S7-5 ***************************************************************/
        YW_GC_PTR(yw_dom_node) parent_root =
            yw_dom_root(parent, YW_DOM_NO_SEARCH_FLAGS);
        if (parent_root->type == YW_DOM_SHADOW_ROOT_NODE)
        {
            YW_TODO();
        }

        /* S7-6 ***************************************************************/
        /* TODO: Run assign slottables for a tree with node’s root. */

        /* S7-7 ***************************************************************/
        yw_dom_iter dscn_iter;
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(yw_dom_node) dscn_node =
                yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }

            /* S7-7-1 *********************************************************/
            if (dscn_node->type == YW_DOM_ELEMENT_NODE)
            {
                YW_GC_PTR(yw_dom_element_node) dscn_elem =
                    (YW_GC_PTR(yw_dom_element_node))dscn_node;
                /* S7-7-2 *****************************************************/
                yw_dom_custom_element_registry *reg =
                    dscn_elem->custom_element_registry;
                if (reg == NULL)
                {
                    reg = yw_dom_lookup_custom_element_registry(
                        dscn_node->parent);
                    dscn_elem->custom_element_registry = reg;
                }
                else if (reg->is_scoped)
                {
                    YW_LIST_PUSH(YW_GC_PTR(yw_dom_document),
                                 reg->scoped_document_set,
                                 dscn_node->node_document);
                }
                else if (yw_dom_is_element_custom(dscn_node))
                {
                    YW_TODO();
                }
                else
                {
                    yw_dom_try_upgrade_element(dscn_node);
                }
            }
            else if (dscn_node->type == YW_DOM_SHADOW_ROOT_NODE)
            {
                YW_TODO();
            }
        }
    }

    /* S8 *********************************************************************/
    if (!(flags & YW_DOM_SUPPRESS_OBSERVERS))
    {
        /* TODO */
    }

    /* S9 *********************************************************************/
    if (parent->callbacks != NULL &&
        parent->callbacks->run_children_changed_steps != NULL)
    {
        parent->callbacks->run_children_changed_steps(parent);
    }

    /* S10 ********************************************************************/
    yw_dom_node_list static_node_list;
    YW_LIST_INIT(&static_node_list);

    /* S11 ********************************************************************/
    for (int i = 0; i < nodes.len; i++)
    {
        YW_GC_PTR(yw_dom_node) node = nodes.items[i];

        yw_dom_iter dscn_iter;
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(yw_dom_node) dscn_node =
                yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            YW_LIST_PUSH(YW_GC_PTR(yw_dom_node), &static_node_list, dscn_node);
        }
    }

    /* S12 ********************************************************************/
    for (int i = 0; i < static_node_list.len; i++)
    {
        YW_GC_PTR(yw_dom_node) node = static_node_list.items[i];

        if (yw_dom_is_connected(node) && node->callbacks != NULL &&
            node->callbacks->run_post_connection_steps != NULL)
        {
            node->callbacks->run_post_connection_steps(node);
        }
    }

    node->parent = parent;

out:
    YW_LIST_FREE(&nodes);
    YW_LIST_FREE(&static_node_list);
}

void yw_dom_append_child(void *node_v, void *child_v)
{
    yw_dom_insert(child_v, node_v, NULL, YW_DOM_NO_INSERT_FLAGS);
}

void yw_dom_adopt_into(void *node_v, YW_GC_PTR(yw_dom_document) document)
{
    /* https://dom.spec.whatwg.org/#concept-node-adopt */

    /* NOTE: All the step numbers(S#.) are based on spec from when this was
     * initially written(2025.11.13) */

    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    /* S1 *********************************************************************/
    YW_GC_PTR(yw_dom_document) old_document = node->node_document;

    /* S2 *********************************************************************/
    if (node->parent != NULL)
    {
        /* TODO: Remove node */
        YW_TODO();
    }

    /* S3 *********************************************************************/
    if (document != old_document)
    {
        /* S3-1 ***************************************************************/
        yw_dom_iter dscn_iter;
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(yw_dom_node) dscn_node =
                yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            /* S3-1-1 *********************************************************/
            dscn_node->node_document = document;
            if (dscn_node->type == YW_DOM_SHADOW_ROOT_NODE)
            {
                /* S3-1-2 *****************************************************/
                YW_TODO();
            }
            else if (dscn_node->type == YW_DOM_ELEMENT_NODE)
            {
                /* S3-1-3 *****************************************************/
                YW_GC_PTR(yw_dom_element_node) dscn_elem =
                    (YW_GC_PTR(yw_dom_element_node))dscn_node;

                /* S3-1-3-1 ***************************************************/
                for (int i = 0; i < dscn_elem->attrs.len; i++)
                {
                    YW_GC_PTR(yw_dom_attr_node) dscn_attr =
                        dscn_elem->attrs.items[i];
                    dscn_attr->_node.node_document = document;
                }
                /* S3-1-3-2 ***************************************************/
                if (yw_dom_is_global_custom_element_registry(
                        yw_dom_lookup_custom_element_registry(dscn_node)))
                {
                    YW_TODO();
                }
            }
        }

        /* S3-2 ***************************************************************/
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(yw_dom_node) dscn_node =
                yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            if (!yw_dom_is_element_custom(dscn_node))
            {
                continue;
            }
            YW_TODO();
        }

        /* S3-3 ***************************************************************/
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(yw_dom_node) dscn_node =
                yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            if (dscn_node->callbacks != NULL &&
                dscn_node->callbacks->run_adopting_steps != NULL)
            {
                dscn_node->callbacks->run_adopting_steps(dscn_node,
                                                         old_document);
            }
        }
    }
}

void yw_dom_print_tree(FILE *dest, void *node_v, int indent_level)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    /* Print indent */
    fprintf(dest, "%*s", indent_level * 4, "");

    switch (node->type)
    {
    case YW_DOM_TEXT_NODE: {
        YW_GC_PTR(yw_dom_character_data_node) cdata =
            (YW_GC_PTR(yw_dom_character_data_node))node;
        fprintf(dest, "#text %s", cdata->text);
        break;
    }
    case YW_DOM_ELEMENT_NODE: {
        YW_GC_PTR(yw_dom_element_node) elem =
            (YW_GC_PTR(yw_dom_element_node))node;
        fprintf(dest, "<%s", elem->local_name);
        for (int i = 0; i < elem->attrs.len; i++)
        {
            fprintf(dest, " %s=%s", elem->attrs.items[i]->local_name,
                    elem->attrs.items[i]->value);
        }
        fprintf(dest, ">");
        break;
    }
    case YW_DOM_DOCUMENT_NODE:
        YW_TODO();
    case YW_DOM_DOCUMENT_FRAGMENT_NODE:
        YW_TODO();
    case YW_DOM_SHADOW_ROOT_NODE:
        YW_TODO();
    }
    fprintf(dest, "\n");
    for (int i = 0; i < node->children.len; i++)
    {
        yw_dom_print_tree(dest, node->children.items[i], indent_level + 1);
    }
}

YW_GC_PTR(yw_dom_custom_element_registry)
yw_dom_lookup_custom_element_registry(void *node_v)
{
    YW_GC_PTR(yw_dom_node) node = (YW_GC_PTR(yw_dom_node))node_v;
    yw_dom_node_check_magic(node);

    switch (node->type)
    {
    case YW_DOM_ELEMENT_NODE:
        return ((YW_GC_PTR(yw_dom_element_node))node)->custom_element_registry;
    case YW_DOM_DOCUMENT_NODE:
        return ((YW_GC_PTR(yw_dom_document_node))node)->custom_element_registry;
    case YW_DOM_SHADOW_ROOT_NODE:
        YW_TODO();
    }
}

/*******************************************************************************
 * Element
 ******************************************************************************/

static yw_gc_callbacks yw_dom_element_gc_callbacks = {
    .visit = yw_dom_element_visit,
};

void yw_dom_element_init(YW_GC_PTR(yw_dom_element_node) out)
{
    yw_dom_node_init(&out->_node);
}
YW_GC_PTR(yw_dom_element_node) yw_dom_element_alloc(
    yw_gc_heap *heap, yw_gc_alloc_flags alloc_flags)
{
    YW_GC_PTR(yw_dom_element_node) elem = YW_GC_ALLOC(
        yw_dom_element_node, heap, &yw_dom_element_gc_callbacks, alloc_flags);
    yw_dom_element_init(elem);
    return elem;
}
void yw_dom_element_visit(void *node_v)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);
    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        fprintf(stderr, "%s: non-element node found\n", __func__);
        abort();
    }
    yw_dom_node_visit(&elem->_node);
}
bool yw_dom_is_shadow_host(void *node_v)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);
    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        return false;
    }
    return elem->shadow_root != NULL;
}

bool yw_dom_is_element_defined(void *node_v)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);

    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        return false;
    }
    /* https://dom.spec.whatwg.org/#concept-element-defined */
    return elem->custom_element_state == YW_DOM_CUSTOM_ELEMENT_UNCUSTOMIZED ||
           elem->custom_element_state == YW_DOM_CUSTOM_ELEMENT_CUSTOM;
}

bool yw_dom_is_element_custom(void *node_v)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);

    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        return false;
    }
    /* https://dom.spec.whatwg.org/#concept-element-custom */
    return elem->custom_element_state == YW_DOM_CUSTOM_ELEMENT_CUSTOM;
}

bool yw_dom_is_element_inside(void *node_v, char const *namespace_,
                              char const *local_name)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);

    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        return false;
    }

    YW_GC_PTR(yw_dom_element_node) current = elem;
    while (current != NULL)
    {
        if (yw_dom_is_element(current, namespace_, local_name))
        {
            return true;
        }
        current = (YW_GC_PTR(yw_dom_element_node))current->_node.parent;
    }
    return false;
}

bool yw_dom_is_element(void *node_v, char const *namespace_,
                       char const *local_name)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);

    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        return false;
    }
    return elem->namespace_ != NULL &&
           strcmp(elem->namespace_, namespace_) == 0 &&
           strcmp(elem->local_name, local_name) == 0;
}
bool yw_dom_is_html_element(void *node_v, char const *local_name)
{
    yw_dom_is_element(node_v, YW_HTML_NAMESPACE, local_name);
}
bool yw_dom_is_mathml_element(void *node_v, char const *local_name)
{
    yw_dom_is_element(node_v, YW_MATHML_NAMESPACE, local_name);
}
bool yw_dom_is_svg_element(void *node_v, char const *local_name)
{
    yw_dom_is_element(node_v, YW_SVG_NAMESPACE, local_name);
}
void yw_dom_append_attr(void *node_v, yw_gc_heap *heap,
                        yw_dom_attr_data const *data)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);

    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        fprintf(stderr, "%s: non-element node found\n", __func__);
        abort();
    }

    YW_GC_PTR(yw_dom_attr_node) attr =
        yw_dom_attr_alloc(heap, YW_NO_GC_ALLOC_FLAGS);
    attr->local_name = data->local_name;
    attr->value = data->value;
    attr->namespace_ = data->namespace_;
    attr->namespace_prefix = data->namespace_prefix;
    attr->element = elem;
}

char const *yw_dom_attr(void *node_v, char const *namespace_,
                        char const *local_name)
{
    YW_GC_PTR(yw_dom_element_node) elem =
        (YW_GC_PTR(yw_dom_element_node))node_v;
    yw_dom_node_check_magic(&elem->_node);

    if (elem->_node.type != YW_DOM_ELEMENT_NODE)
    {
        fprintf(stderr, "%s: non-element node found\n", __func__);
        abort();
    }

    for (int i = 0; i < elem->attrs.len; i++)
    {
        YW_GC_PTR(yw_dom_attr_node) attr = elem->attrs.items[i];
        bool ns_match = ((namespace_ == NULL && attr->namespace_ == NULL) ||
                         (namespace_ != NULL && attr->namespace_ != NULL &&
                          strcmp(attr->namespace_, namespace_) == 0));
        bool local_name_match = strcmp(attr->local_name, local_name) == 0;
        if (ns_match && local_name_match)
        {
            return attr->value;
        }
    }

    return NULL;
}