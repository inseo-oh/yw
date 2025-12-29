/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_dom.h"
#include "yw_common.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*******************************************************************************
 * Node
 ******************************************************************************/

#define YW_NODE_MAGIC 0xb1fedf1b

#define YW_VERIFY_NODE_MAGIC(_node)                                             \
    do                                                                          \
    {                                                                           \
        if ((_node) != NULL && (_node)->magic != YW_NODE_MAGIC)                 \
        {                                                                       \
            fprintf(stderr, "%s: %s has corrupted magic!\n", __func__, #_node); \
            abort();                                                            \
        }                                                                       \
    } while (0)

#define YW_VERIFY_NODE_TYPE(_node, _type_flags)                                           \
    do                                                                                    \
    {                                                                                     \
        if (!yw_dom_has_type((_node), (_type_flags)))                                     \
        {                                                                                 \
            fprintf(stderr, "%s: %s's type is not %s\n", __func__, #_node, #_type_flags); \
            abort();                                                                      \
        }                                                                                 \
    } while (0)

void yw_dom_node_init(YW_GC_PTR(YW_DOMNode) out)
{
    out->magic = YW_NODE_MAGIC;
}
void yw_dom_node_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;

    for (int i = 0; i < node->children.len; i++)
    {
        yw_gc_visit(node->children.items[i]);
    }
    yw_gc_visit(node->parent);
    yw_gc_visit(node->node_document);
}
void yw_dom_node_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_LIST_FREE(&node->children);
}

YW_GC_PTR(YW_DOMNode) yw_dom_first_child(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);
    if (node->children.len == 0)
    {
        return NULL;
    }
    return node->children.items[0];
}

YW_GC_PTR(YW_DOMNode) yw_dom_last_child(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);
    if (node->children.len == 0)
    {
        return NULL;
    }
    return node->children.items[node->children.len - 1];
}

YW_GC_PTR(YW_DOMNode) yw_dom_next_sibling(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    YW_GC_PTR(YW_DOMNode) parent = node->parent;
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

YW_GC_PTR(YW_DOMNode) yw_dom_prev_sibling(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    YW_GC_PTR(YW_DOMNode) parent = node->parent;
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

YW_GC_PTR(YW_DOMNode) yw_dom_root(void *node_v, YW_DOMSearchFlags flags)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    YW_GC_PTR(YW_DOMNode) res = node;
    while (res->parent != NULL)
    {
        if ((flags & YW_DOM_SHADOW_INCLUDING) && (yw_dom_has_type(res, YW_DOM_SHADOW_ROOT_NODE)))
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
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    YW_GC_PTR(YW_DOMNode) parent = node->parent;
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

bool yw_dom_has_type(void *node_v, YW_DOMNodeTypeFlags flags)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    if (node == NULL)
    {
        return false;
    }
    return (node->type_flags & flags) == flags;
}

bool yw_dom_is_in_same_tree(void *node_a_v, void *node_b_v)
{
    return yw_dom_root(node_a_v, YW_DOM_NO_SEARCH_FLAGS) == yw_dom_root(node_b_v, YW_DOM_NO_SEARCH_FLAGS);
}

bool yw_dom_is_connected(void *node_v)
{
    /* https://dom.spec.whatwg.org/#connected */

    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    return yw_dom_root(node_v, YW_DOM_SHADOW_INCLUDING) == (YW_GC_PTR(YW_DOMNode))node->node_document;
}

bool yw_dom_is_in_document_tree(void *node_v)
{
    return yw_dom_has_type(yw_dom_root(node_v, YW_DOM_NO_SEARCH_FLAGS), YW_DOM_DOCUMENT_NODE);
}

/* Caller owns the returned string. */
char *yw_dom_child_text(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    char *res_buf = NULL;

    for (int i = 0; i < node->children.len; i++)
    {
        if (yw_dom_has_type(node->children.items[i], YW_DOM_TEXT_NODE))
        {
            char *node_text = ((YW_GC_PTR(YW_DOMCharacterData))node->children.items[i])->text;
            yw_append_str(&res_buf, node_text);
        }
    }

    return res_buf;
}

YW_GC_PTR(YW_DOMNode) yw_dom_next_descendant(YW_DOMIter *iter)
{
    /* https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant */

    YW_GC_PTR(YW_DOMNode) curr_node = iter->last_node;
    YW_GC_PTR(YW_DOMNode) res = NULL;

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
    if (iter->shadow_including && yw_dom_has_type(res, YW_DOM_SHADOW_ROOT_NODE))
    {
        /* TODO: Support shadow root */
        YW_TODO();
    }

    iter->last_node = res;
    return res;
}

void yw_dom_inclusive_descendants_init(YW_DOMIter *out, void *root_node_v, YW_DOMSearchFlags flags)
{
    YW_GC_PTR(YW_DOMNode) root_node = (YW_GC_PTR(YW_DOMNode))root_node_v;
    YW_VERIFY_NODE_MAGIC(root_node);

    memset(out, 0, sizeof(*out));
    out->root_node = root_node;
    if (flags & YW_DOM_SHADOW_INCLUDING)
    {
        out->shadow_including = true;
    }
}

void yw_dom_descendants_init(YW_DOMIter *out, void *root_node_v, YW_DOMSearchFlags flags)
{
    yw_dom_inclusive_descendants_init(out, root_node_v, flags);
    yw_dom_next_descendant(out);
}

YW_GC_PTR(YW_DOMNode) yw_dom_next_ancestor(YW_DOMIter *iter)
{
    /* https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor */

    YW_GC_PTR(YW_DOMNode) curr_node = iter->last_node;
    YW_GC_PTR(YW_DOMNode) res = NULL;

    if (curr_node == NULL)
    {
        /* This is our first call */
        res = iter->root_node;
    }
    else if (curr_node->parent != NULL)
    {
        res = curr_node->parent;
    }
    else if (yw_dom_has_type(curr_node, YW_DOM_SHADOW_ROOT_NODE))
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

void yw_dom_inclusive_ancestors_init(YW_DOMIter *out, void *root_node_v, YW_DOMSearchFlags flags)
{
    YW_GC_PTR(YW_DOMNode) root_node = (YW_GC_PTR(YW_DOMNode))root_node_v;
    YW_VERIFY_NODE_MAGIC(root_node);

    memset(out, 0, sizeof(*out));
    out->root_node = root_node;
    if (flags & YW_DOM_SHADOW_INCLUDING)
    {
        out->shadow_including = true;
    }
}

void yw_dom_ancestors_init(YW_DOMIter *out, void *root_node_v, YW_DOMSearchFlags flags)
{
    yw_dom_inclusive_ancestors_init(out, root_node_v, flags);
    yw_dom_next_descendant(out);
}

void yw_dom_insert(void *node_v, void *parent_v, void *before_child_v, YW_DOMInsertFlag flags)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_GC_PTR(YW_DOMNode) parent = (YW_GC_PTR(YW_DOMNode))parent_v;
    YW_GC_PTR(YW_DOMNode) before_child = (YW_GC_PTR(YW_DOMNode))before_child_v;
    YW_VERIFY_NODE_MAGIC(node);
    YW_VERIFY_NODE_MAGIC(parent);
    YW_VERIFY_NODE_MAGIC(parent);

    YW_GC_PTR(YW_DOMNode) prev_sibling;

    /* NOTE: All the step numbers(S#.) are based on spec from when this was
     * initially written(2025.11.13) */

    /* S1 *********************************************************************/
    YW_DOMNodeList nodes;
    YW_LIST_INIT(&nodes);

    if (yw_dom_has_type(node, YW_DOM_DOCUMENT_FRAGMENT_NODE))
    {
        for (int i = 0; i < node->children.len; i++)
        {
            YW_LIST_PUSH(YW_GC_PTR(YW_DOMNode), &nodes, node->children.items[i]);
        }
    }
    else
    {
        YW_LIST_PUSH(YW_GC_PTR(YW_DOMNode), &nodes, node);
    }

    /* S2 *********************************************************************/
    int count = nodes.len;

    /* S3 *********************************************************************/
    if (count == 0)
    {
        goto out;
    }

    /* S4 *********************************************************************/
    if (yw_dom_has_type(node, YW_DOM_DOCUMENT_FRAGMENT_NODE))
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
        YW_GC_PTR(YW_DOMNode) node = nodes.items[i];
        /* S7-1 ***************************************************************/
        yw_dom_adopt_into(node, parent->node_document);

        if (before_child == NULL)
        {
            /* S7-2 ***********************************************************/
            YW_LIST_PUSH(YW_GC_PTR(YW_DOMNode), &parent->children, node);
        }
        else
        {
            /* S7-3 ***********************************************************/
            int index = yw_dom_index(before_child);
            YW_LIST_INSERT(YW_GC_PTR(YW_DOMNode), &parent->children, index, node);
        }

        /* S7-4 ***************************************************************/
        if (yw_dom_is_shadow_host(parent))
        {
            YW_TODO();
        }

        /* S7-5 ***************************************************************/
        YW_GC_PTR(YW_DOMNode) parent_root = yw_dom_root(parent, YW_DOM_NO_SEARCH_FLAGS);
        if (yw_dom_has_type(parent_root, YW_DOM_SHADOW_ROOT_NODE))
        {
            YW_TODO();
        }

        /* S7-6 ***************************************************************/
        /* TODO: Run assign slottables for a tree with node’s root. */

        /* S7-7 ***************************************************************/
        YW_DOMIter dscn_iter;
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(YW_DOMNode) dscn_node = yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }

            /* S7-7-1 *********************************************************/
            if (yw_dom_has_type(dscn_node, YW_DOM_ELEMENT_NODE))
            {
                YW_GC_PTR(YW_DOMElement) dscn_elem = (YW_GC_PTR(YW_DOMElement))dscn_node;
                /* S7-7-2 *****************************************************/
                YW_GC_PTR(YW_DOMCustomElementRegistry) reg = dscn_elem->custom_element_registry;
                if (reg == NULL)
                {
                    reg = yw_lookup_custom_element_registry(dscn_node->parent);
                    dscn_elem->custom_element_registry = reg;
                }
                else if (reg->is_scoped)
                {
                    YW_LIST_PUSH(YW_GC_PTR(YW_DOMDocument), &reg->scoped_document_set, dscn_node->node_document);
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
            else if (yw_dom_has_type(dscn_node, YW_DOM_SHADOW_ROOT_NODE))
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
    if (parent->callbacks != NULL && parent->callbacks->run_children_changed_steps != NULL)
    {
        parent->callbacks->run_children_changed_steps(parent);
    }

    /* S10 ********************************************************************/
    YW_DOMNodeList static_node_list;
    YW_LIST_INIT(&static_node_list);

    /* S11 ********************************************************************/
    for (int i = 0; i < nodes.len; i++)
    {
        YW_GC_PTR(YW_DOMNode) node = nodes.items[i];

        YW_DOMIter dscn_iter;
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(YW_DOMNode) dscn_node = yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            YW_LIST_PUSH(YW_GC_PTR(YW_DOMNode), &static_node_list, dscn_node);
        }
    }

    /* S12 ********************************************************************/
    for (int i = 0; i < static_node_list.len; i++)
    {
        YW_GC_PTR(YW_DOMNode) node = static_node_list.items[i];

        if (yw_dom_is_connected(node) && node->callbacks != NULL && node->callbacks->run_post_connection_steps != NULL)
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

void yw_dom_adopt_into(void *node_v, YW_GC_PTR(YW_DOMDocument) document)
{
    /* https://dom.spec.whatwg.org/#concept-node-adopt */

    /* NOTE: All the step numbers(S#.) are based on spec from when this was
     * initially written(2025.11.13) */

    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    /* S1 *********************************************************************/
    YW_GC_PTR(YW_DOMDocument) old_document = node->node_document;

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
        YW_DOMIter dscn_iter;
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(YW_DOMNode) dscn_node = yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            /* S3-1-1 *********************************************************/
            dscn_node->node_document = document;
            if (yw_dom_has_type(dscn_node, YW_DOM_SHADOW_ROOT_NODE))
            {
                /* S3-1-2 *****************************************************/
                YW_TODO();
            }
            else if (yw_dom_has_type(dscn_node, YW_DOM_ELEMENT_NODE))
            {
                /* S3-1-3 *****************************************************/
                YW_GC_PTR(YW_DOMElement) dscn_elem = (YW_GC_PTR(YW_DOMElement))dscn_node;

                /* S3-1-3-1 ***************************************************/
                for (int i = 0; i < dscn_elem->attrs.len; i++)
                {
                    YW_GC_PTR(YW_DOMAttr) dscn_attr = dscn_elem->attrs.items[i];
                    dscn_attr->_node.node_document = document;
                }
                /* S3-1-3-2 ***************************************************/
                if (yw_lookup_custom_element_registry(yw_lookup_custom_element_registry(dscn_node)))
                {
                    YW_TODO();
                }
            }
        }

        /* S3-2 ***************************************************************/
        yw_dom_descendants_init(&dscn_iter, node, YW_DOM_SHADOW_INCLUDING);
        while (1)
        {
            YW_GC_PTR(YW_DOMNode) dscn_node = yw_dom_next_descendant(&dscn_iter);
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
            YW_GC_PTR(YW_DOMNode) dscn_node = yw_dom_next_descendant(&dscn_iter);
            if (dscn_node == NULL)
            {
                break;
            }
            if (dscn_node->callbacks != NULL && dscn_node->callbacks->run_adopting_steps != NULL)
            {
                dscn_node->callbacks->run_adopting_steps(dscn_node, old_document);
            }
        }
    }
}

void yw_dom_print_tree(FILE *dest, void *node_v, int indent_level)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    /* Print indent */
    fprintf(dest, "%*s", indent_level * 4, "");

    if (yw_dom_has_type(node, YW_DOM_TEXT_NODE))
    {
        YW_GC_PTR(YW_DOMCharacterData) cdata = (YW_GC_PTR(YW_DOMCharacterData))node;
        fprintf(dest, "#text %s", cdata->text);
    }
    else if (yw_dom_has_type(node, YW_DOM_ELEMENT_NODE))
    {
        YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node;
        fprintf(dest, "<%s", elem->local_name);
        for (int i = 0; i < elem->attrs.len; i++)
        {
            fprintf(dest, " %s=%s", elem->attrs.items[i]->local_name, elem->attrs.items[i]->value);
        }
        fprintf(dest, ">");
    }
    else if (yw_dom_has_type(node, YW_DOM_DOCUMENT_NODE))
    {
        YW_GC_PTR(YW_DOMDocument) doc = (YW_GC_PTR(YW_DOMDocument))node;

        fprintf(dest, "#document(mode=");
        switch (doc->mode)
        {
        case YW_NO_QUIRKS:
            fprintf(dest, "no-quirks");
            break;
        case YW_QUIRKS:
            fprintf(dest, "quirks");
            break;
        case YW_LIMITED_QUIRKS:
            fprintf(dest, "limited-quirks");
            break;
        }
        fprintf(dest, ")");
    }
    else if (yw_dom_has_type(node, YW_DOM_DOCUMENT_TYPE_NODE))
    {
        YW_GC_PTR(YW_DOMDocumentType) doctype = (YW_GC_PTR(YW_DOMDocumentType))node;

        fprintf(dest, "<!DOCTYPE");
        if (doctype->name != NULL)
        {
            fprintf(dest, " %s", doctype->name);
        }
        if (doctype->public_id != NULL && doctype->system_id == NULL)
        {
            fprintf(dest, " PUBLIC \"%s\"", doctype->public_id);
        }
        else if (doctype->public_id == NULL && doctype->system_id != NULL)
        {
            fprintf(dest, " SYSTEM \"%s\"", doctype->system_id);
        }
        else if (doctype->public_id != NULL && doctype->system_id != NULL)
        {
            fprintf(dest, " PUBLIC \"%s\" \"%s\"", doctype->public_id, doctype->system_id);
        }
        fprintf(dest, ">");
    }
    else
    {
        fprintf(dest, "<unknown node with type_flags=%#x>", node->type_flags);
    }
    fprintf(dest, "\n");
    for (int i = 0; i < node->children.len; i++)
    {
        yw_dom_print_tree(dest, node->children.items[i], indent_level + 1);
    }
}

YW_GC_PTR(YW_DOMCustomElementRegistry)
yw_lookup_custom_element_registry(void *node_v)
{
    YW_GC_PTR(YW_DOMNode) node = (YW_GC_PTR(YW_DOMNode))node_v;
    YW_VERIFY_NODE_MAGIC(node);

    if (yw_dom_has_type(node, YW_DOM_ELEMENT_NODE))
    {
        return ((YW_GC_PTR(YW_DOMElement))node)->custom_element_registry;
    }
    else if (yw_dom_has_type(node, YW_DOM_DOCUMENT_NODE))
    {
        return ((YW_GC_PTR(YW_DOMDocument))node)->custom_element_registry;
    }
    else if (yw_dom_has_type(node, YW_DOM_SHADOW_ROOT_NODE))
    {
        YW_TODO();
    }
    return NULL;
}

/*******************************************************************************
 * Custom elements
 ******************************************************************************/

YW_DOMCustomElementDefinition const *yw_dom_lookup_custom_element_definition(YW_GC_PTR(YW_DOMCustomElementRegistry) registry, char const *namespace_, char const *local_name, char const *is)
{
    (void)registry;
    (void)namespace_;
    (void)local_name;
    (void)is;
    /* STUB */
    return NULL;
}

bool yw_dom_is_global_custom_element_reigstry(YW_GC_PTR(YW_DOMCustomElementRegistry) registry)
{
    /* https://dom.spec.whatwg.org/#is-a-global-custom-element-registry */

    if (registry == NULL)
    {
        return false;
    }
    return !registry->is_scoped;
}

void yw_dom_try_upgrade_element(void *node_v)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);
    YW_VERIFY_NODE_TYPE(elem, YW_DOM_ELEMENT_NODE);

    YW_DOMCustomElementDefinition const *definition = yw_dom_lookup_custom_element_definition(elem->custom_element_registry, elem->namespace_, elem->local_name, elem->is);
    if (definition != NULL)
    {
        YW_TODO();
    }
}

/*******************************************************************************
 * Document
 ******************************************************************************/

static YW_GcCallbacks yw_document_gc_callbacks = {
    .visit = yw_document_visit,
    .destroy = yw_document_destroy,
};

void yw_document_init(YW_GC_PTR(YW_DOMDocument) out)
{
    yw_dom_node_init(&out->_node);
    out->_node.type_flags |= YW_DOM_DOCUMENT_NODE;

    /* Node document of a document is itself */
    out->_node.node_document = out;
}

YW_GC_PTR(YW_DOMDocument) yw_document_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(YW_DOMDocument) document = YW_GC_ALLOC(YW_DOMDocument, heap, &yw_document_gc_callbacks, alloc_flags);
    yw_document_init(document);
    return document;
}
void yw_document_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMDocument) doc = (YW_GC_PTR(YW_DOMDocument))node_v;
    YW_VERIFY_NODE_MAGIC(&doc->_node);
    YW_VERIFY_NODE_TYPE(doc, YW_DOM_DOCUMENT_NODE);

    yw_dom_node_visit(&doc->_node);
    yw_gc_visit(doc->custom_element_registry);
}
void yw_document_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMDocument) doc = (YW_GC_PTR(YW_DOMDocument))node_v;
    YW_VERIFY_NODE_MAGIC(&doc->_node);
    YW_VERIFY_NODE_TYPE(doc, YW_DOM_DOCUMENT_NODE);

    yw_dom_node_destroy(doc);
}

YW_GC_PTR(YW_DOMCustomElementRegistry)
yw_document_effective_global_custom_element_registry(void *node_v)
{
    /* https://dom.spec.whatwg.org/#effective-global-custom-element-registry */

    YW_GC_PTR(YW_DOMDocument) doc = (YW_GC_PTR(YW_DOMDocument))node_v;
    YW_VERIFY_NODE_MAGIC(&doc->_node);
    YW_VERIFY_NODE_TYPE(doc, YW_DOM_DOCUMENT_NODE);

    if (yw_dom_is_global_custom_element_reigstry(doc->custom_element_registry))
    {
        return doc->custom_element_registry;
    }
    return NULL;
}

/*******************************************************************************
 * DocumentFragment
 ******************************************************************************/

static YW_GcCallbacks yw_document_fragment_gc_callbacks = {
    .visit = yw_document_fragment_visit,
    .destroy = yw_document_fragment_destroy,
};

void yw_document_fragment_init(YW_GC_PTR(YW_DOMDocumentFragment) out)
{
    yw_dom_node_init(&out->_node);
    out->_node.type_flags |= YW_DOM_DOCUMENT_FRAGMENT_NODE;
}

YW_GC_PTR(YW_DOMDocumentFragment) yw_document_fragment_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(YW_DOMDocumentFragment) docfrag = YW_GC_ALLOC(YW_DOMDocumentFragment, heap, &yw_document_fragment_gc_callbacks, alloc_flags);
    yw_document_fragment_init(docfrag);
    return docfrag;
}
void yw_document_fragment_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMDocumentFragment) docfrag = (YW_GC_PTR(YW_DOMDocumentFragment))node_v;
    YW_VERIFY_NODE_MAGIC(&docfrag->_node);
    YW_VERIFY_NODE_TYPE(docfrag, YW_DOM_DOCUMENT_FRAGMENT_NODE);

    yw_dom_node_visit(&docfrag->_node);
    yw_gc_visit(docfrag->host);
}
void yw_document_fragment_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMDocumentFragment) docfrag = (YW_GC_PTR(YW_DOMDocumentFragment))node_v;
    YW_VERIFY_NODE_MAGIC(&docfrag->_node);
    YW_VERIFY_NODE_TYPE(docfrag, YW_DOM_DOCUMENT_FRAGMENT_NODE);

    yw_dom_node_destroy(docfrag);
}

/*******************************************************************************
 * DocumentType
 ******************************************************************************/

static YW_GcCallbacks yw_document_type_gc_callbacks = {
    .visit = yw_document_type_visit,
    .destroy = yw_document_type_destroy,
};

void yw_document_type_init(YW_GC_PTR(YW_DOMDocumentType) out)
{
    yw_dom_node_init(&out->_node);
    out->_node.type_flags |= YW_DOM_DOCUMENT_TYPE_NODE;
}

YW_GC_PTR(YW_DOMDocumentType) yw_document_type_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(YW_DOMDocumentType) doctype = YW_GC_ALLOC(YW_DOMDocumentType, heap, &yw_document_type_gc_callbacks, alloc_flags);
    yw_document_type_init(doctype);
    return doctype;
}
void yw_document_type_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMDocumentType) doctype = (YW_GC_PTR(YW_DOMDocumentType))node_v;
    YW_VERIFY_NODE_MAGIC(&doctype->_node);
    YW_VERIFY_NODE_TYPE(doctype, YW_DOM_DOCUMENT_TYPE_NODE);

    yw_dom_node_visit(&doctype->_node);
}
void yw_document_type_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMDocumentType) doctype = (YW_GC_PTR(YW_DOMDocumentType))node_v;
    YW_VERIFY_NODE_MAGIC(&doctype->_node);
    YW_VERIFY_NODE_TYPE(doctype, YW_DOM_DOCUMENT_TYPE_NODE);

    yw_dom_node_destroy(doctype);
    free(doctype->name);
    free(doctype->public_id);
    free(doctype->system_id);
}

/*******************************************************************************
 * Attr
 ******************************************************************************/

void yw_dom_attr_data_deinit(YW_DOMAttrData *data)
{
    free(data->local_name);
    free(data->value);
    free(data->namespace_);
    free(data->namespace_prefix);
}

static YW_GcCallbacks yw_attr_gc_callbacks = {
    .visit = yw_attr_visit,
    .destroy = yw_attr_destroy,
};

void yw_attr_init(YW_GC_PTR(YW_DOMAttr) out)
{
    yw_dom_node_init(&out->_node);
    out->_node.type_flags |= YW_DOM_ATTR_NODE;
}

YW_GC_PTR(YW_DOMAttr) yw_attr_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(YW_DOMAttr) attr = YW_GC_ALLOC(YW_DOMAttr, heap, &yw_attr_gc_callbacks, alloc_flags);
    yw_attr_init(attr);
    return attr;
}
void yw_attr_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMAttr) attr = (YW_GC_PTR(YW_DOMAttr))node_v;
    YW_VERIFY_NODE_MAGIC(&attr->_node);
    YW_VERIFY_NODE_TYPE(attr, YW_DOM_ATTR_NODE);

    yw_dom_node_visit(&attr->_node);
    yw_gc_visit(attr->element);
}
void yw_attr_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMAttr) attr = (YW_GC_PTR(YW_DOMAttr))node_v;
    YW_VERIFY_NODE_MAGIC(&attr->_node);
    YW_VERIFY_NODE_TYPE(attr, YW_DOM_ATTR_NODE);

    free(attr->local_name);
    free(attr->namespace_);
    free(attr->namespace_prefix);
    free(attr->value);

    yw_dom_node_destroy(attr);
}

/*******************************************************************************
 * Element
 ******************************************************************************/

static YW_GcCallbacks yw_dom_element_gc_callbacks = {
    .visit = yw_dom_element_visit,
    .destroy = yw_dom_element_destroy,
};

void yw_dom_element_init(YW_GC_PTR(YW_DOMElement) out)
{
    yw_dom_node_init(&out->_node);
    out->_node.type_flags |= YW_DOM_ELEMENT_NODE;
}

YW_GC_PTR(YW_DOMElement) yw_dom_element_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(YW_DOMElement) elem = YW_GC_ALLOC(YW_DOMElement, heap, &yw_dom_element_gc_callbacks, alloc_flags);
    yw_dom_element_init(elem);
    return elem;
}
void yw_dom_element_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);
    YW_VERIFY_NODE_TYPE(elem, YW_DOM_ELEMENT_NODE);

    yw_dom_node_visit(&elem->_node);
    yw_gc_visit(elem->shadow_root);
    yw_gc_visit(elem->custom_element_registry);
}
void yw_dom_element_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);
    YW_VERIFY_NODE_TYPE(elem, YW_DOM_ELEMENT_NODE);

    yw_dom_node_destroy(elem);
    YW_LIST_FREE(&elem->attrs);
}

bool yw_dom_is_shadow_host(void *node_v)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);

    if (!yw_dom_has_type(elem, YW_DOM_ELEMENT_NODE))
    {
        return false;
    }
    return elem->shadow_root != NULL;
}

bool yw_dom_is_element_defined(void *node_v)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);

    if (!yw_dom_has_type(elem, YW_DOM_ELEMENT_NODE))
    {
        return false;
    }
    /* https://dom.spec.whatwg.org/#concept-element-defined */
    return elem->custom_element_state == YW_DOM_CUSTOM_ELEMENT_UNCUSTOMIZED || elem->custom_element_state == YW_DOM_CUSTOM_ELEMENT_CUSTOM;
}

bool yw_dom_is_element_custom(void *node_v)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);

    if (!yw_dom_has_type(elem, YW_DOM_ELEMENT_NODE))
    {
        return false;
    }
    /* https://dom.spec.whatwg.org/#concept-element-custom */
    return elem->custom_element_state == YW_DOM_CUSTOM_ELEMENT_CUSTOM;
}

bool yw_dom_is_element_inside(void *node_v, char const *namespace_, char const *local_name)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);

    if (!yw_dom_has_type(elem, YW_DOM_ELEMENT_NODE))
    {
        return false;
    }

    YW_GC_PTR(YW_DOMNode) current = elem->_node.parent;
    while (current != NULL)
    {
        if (yw_dom_is_element(current, namespace_, local_name))
        {
            return true;
        }
        current = current->parent;
    }
    return false;
}

bool yw_dom_is_element(void *node_v, char const *namespace_, char const *local_name)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);

    if (!yw_dom_has_type(elem, YW_DOM_ELEMENT_NODE))
    {
        return false;
    }
    return elem->namespace_ != NULL && strcmp(elem->namespace_, namespace_) == 0 && strcmp(elem->local_name, local_name) == 0;
}
bool yw_dom_is_html_element(void *node_v, char const *local_name)
{
    return yw_dom_is_element(node_v, YW_HTML_NAMESPACE, local_name);
}
bool yw_dom_is_mathml_element(void *node_v, char const *local_name)
{
    return yw_dom_is_element(node_v, YW_MATHML_NAMESPACE, local_name);
}
bool yw_dom_is_svg_element(void *node_v, char const *local_name)
{
    return yw_dom_is_element(node_v, YW_SVG_NAMESPACE, local_name);
}
void yw_dom_append_attr_to_element(void *node_v, YW_GcHeap *heap, YW_DOMAttrData const *data)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);
    YW_VERIFY_NODE_TYPE(elem, YW_DOM_ELEMENT_NODE);

    YW_GC_PTR(YW_DOMAttr) attr = yw_attr_alloc(heap, YW_NO_GC_ALLOC_FLAGS);
    attr->local_name = yw_duplicate_str(data->local_name);
    attr->value = yw_duplicate_str(data->value);
    attr->namespace_ = yw_duplicate_str(data->namespace_);
    attr->namespace_prefix = yw_duplicate_str(data->namespace_prefix);
    attr->element = elem;
    attr->_node.parent = &elem->_node;
    YW_LIST_PUSH(YW_GC_PTR(YW_DOMAttr), &elem->attrs, attr);
}

char const *yw_dom_attr_of_element(void *node_v, char const *namespace_, char const *local_name)
{
    YW_GC_PTR(YW_DOMElement) elem = (YW_GC_PTR(YW_DOMElement))node_v;
    YW_VERIFY_NODE_MAGIC(&elem->_node);
    YW_VERIFY_NODE_TYPE(elem, YW_DOM_ELEMENT_NODE);

    for (int i = 0; i < elem->attrs.len; i++)
    {
        YW_GC_PTR(YW_DOMAttr) attr = elem->attrs.items[i];
        bool ns_match = ((namespace_ == NULL && attr->namespace_ == NULL) || (namespace_ != NULL && attr->namespace_ != NULL && strcmp(attr->namespace_, namespace_) == 0));
        bool local_name_match = strcmp(attr->local_name, local_name) == 0;
        if (ns_match && local_name_match)
        {
            return attr->value;
        }
    }

    return NULL;
}

/*******************************************************************************
 * CharacterData
 ******************************************************************************/

void yw_dom_character_data_visit(void *node_v)
{
    YW_GC_PTR(YW_DOMCharacterData) text = (YW_GC_PTR(YW_DOMCharacterData))node_v;
    YW_VERIFY_NODE_MAGIC(&text->_node);
    YW_VERIFY_NODE_TYPE(text, YW_DOM_TEXT_NODE);

    yw_dom_node_visit(&text->_node);
}
void yw_dom_character_data_destroy(void *node_v)
{
    YW_GC_PTR(YW_DOMCharacterData) text = (YW_GC_PTR(YW_DOMCharacterData))node_v;
    YW_VERIFY_NODE_MAGIC(&text->_node);
    YW_VERIFY_NODE_TYPE(text, YW_DOM_TEXT_NODE);

    yw_dom_node_destroy(text);
    free(text->text);
}

static YW_GcCallbacks yw_character_data_gc_callbacks = {
    .visit = yw_dom_character_data_visit,
    .destroy = yw_dom_character_data_destroy,
};

void yw_dom_text_init(YW_GC_PTR(YW_DOMCharacterData) out)
{
    yw_dom_node_init(&out->_node);
    out->_node.type_flags |= YW_DOM_TEXT_NODE;
}

YW_GC_PTR(YW_DOMCharacterData) yw_text_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(YW_DOMCharacterData) text = YW_GC_ALLOC(YW_DOMCharacterData, heap, &yw_character_data_gc_callbacks, alloc_flags);
    yw_dom_text_init(text);
    return text;
}
