/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_DOM_H_
#define YW_DOM_H_

#include "yw_common.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

/*******************************************************************************
 * Node
 ******************************************************************************/

struct yw_dom_node_callbacks
{
    void (*run_insertion_steps)(void *self);
    void (*run_children_changed_steps)(void *self);
    void (*run_post_connection_steps)(void *self);
    void (*run_adopting_steps)(void *self,
                               YW_GC_PTR(struct yw_dom_document) old_document);

    /* Element callbacks ******************************************************/

    void (*intrinsic_size)(double *width_out, double *height_out, void *self);
    void (*popped_from_stack_of_open_elements)(void *self);
    void (*presentational_hints)(void *self);
    void (*run_form_reset_algorithm)(void *self);
};

struct yw_dom_node_list
{
    YW_GC_PTR(struct yw_dom_node) *items;
    int len, cap;
};

enum yw_dom_node_type
{
    YW_DOM_TEXT_NODE, /* yw_dom_character_data_node */

    YW_DOM_DOCUMENT_NODE,
    YW_DOM_DOCUMENT_FRAGMENT_NODE,
    YW_DOM_SHADOW_ROOT_NODE,
    YW_DOM_ELEMENT_NODE,
};

YW_GC_TYPE(struct yw_dom_node)
{
    struct yw_gc_object_header gc_header;
    uint32_t magic;

    YW_GC_PTR(struct yw_dom_node) parent;
    YW_GC_PTR(struct yw_dom_document) node_document;
    struct yw_dom_node_callbacks const *callbacks;
    struct yw_dom_node_list children;
    enum yw_dom_node_type type;
};

enum yw_dom_search_flags
{
    YW_DOM_SHADOW_INCLUDING = 1 << 0
};

void yw_dom_node_init(YW_GC_PTR(struct yw_dom_node) out);
void yw_dom_node_visit(void *node_v);
void yw_dom_node_destroy(void *node_v);

/* Returns NULL if there's no children. */
YW_GC_PTR(struct yw_dom_node) yw_dom_first_child(void *node_v);
/* Returns NULL if there's no children. */
YW_GC_PTR(struct yw_dom_node) yw_dom_last_child(void *node_v);
/* Returns NULL if there's no parent. */
YW_GC_PTR(struct yw_dom_node) yw_dom_next_sibling(void *node_v);
/* Returns NULL if there's no parent. */
YW_GC_PTR(struct yw_dom_node) yw_dom_prev_sibling(void *node_v);
YW_GC_PTR(struct yw_dom_node) yw_dom_root(void *node_v,
                                          enum yw_dom_search_flags flags);

int yw_dom_index(void *node_v);
bool yw_dom_is_in_the_same_tree_as(void *node_a_v, void *node_b_v);
bool yw_dom_is_connected(void *node_v);
char *yw_dom_child_text(void *node_v);

struct yw_dom_iter
{
    YW_GC_PTR(struct yw_dom_node) root_node;
    YW_GC_PTR(struct yw_dom_node) last_node;
    bool shadow_including;
};

YW_GC_PTR(struct yw_dom_node) yw_dom_next_descendant(struct yw_dom_iter *iter);

void yw_dom_inclusive_descendants_init(struct yw_dom_iter *out,
                                       void *root_node_v,
                                       enum yw_dom_search_flags flags);
void yw_dom_descendants_init(struct yw_dom_iter *out, void *root_node_v,
                             enum yw_dom_search_flags flags);

struct yw_dom_parents
{
    YW_GC_PTR(struct yw_dom_node) root_node;
    YW_GC_PTR(struct yw_dom_node) last_node;
    bool shadow_including;
};

YW_GC_PTR(struct yw_dom_node) yw_dom_next_parent(struct yw_dom_iter *iter);
void yw_dom_inclusive_ancestors_init(struct yw_dom_iter *out, void *root_node_v,
                                     enum yw_dom_search_flags flags);
void yw_dom_ancestors_init(struct yw_dom_iter *out, void *root_node_v,
                           enum yw_dom_search_flags flags);

enum yw_dom_insert_flag
{
    YW_DOM_SUPPRESS_OBSERVERS
};

void yw_dom_insert(void *node_v, void *parent_v, void *before_child_v,
                   enum yw_dom_insert_flag flags);
void yw_dom_append_child(void *node_v, void *child_v);
void yw_dom_adopt_into(void *node_v,
                       YW_GC_PTR(struct yw_dom_document) document);
void yw_dom_print_tree(FILE *dest, void *node_v, int indent_level);

/*******************************************************************************
 * Custom elements
 ******************************************************************************/
struct yw_dom_custom_element_registry
{
};

/* https://dom.spec.whatwg.org/#concept-element-custom-element-state */
enum yw_custom_element_state
{
    YW_DOM_CUSTOM_ELEMENT_UNDEFINED,
    YW_DOM_CUSTOM_ELEMENT_FAILED,
    YW_DOM_CUSTOM_ELEMENT_UNCUSTOMIZED,
    YW_DOM_CUSTOM_ELEMENT_PRECUSTOMIZED,
    YW_DOM_CUSTOM_ELEMENT_CUSTOM,
};

/*******************************************************************************
 * ShadowRoot
 ******************************************************************************/
YW_GC_TYPE(struct yw_dom_shadow_root_node){
    /* STUB */
};

/*******************************************************************************
 * Attr
 ******************************************************************************/

struct yw_dom_attr_list
{
    YW_GC_PTR(struct yw_dom_attr_node) *items;
    int len, cap;
};

/*******************************************************************************
 * Element
 ******************************************************************************/

YW_GC_TYPE(struct yw_dom_element_node)
{
    YW_GC_TYPE(struct yw_dom_node) _node;

    char const *namespace_;       /* May be NULL */
    char const *namespace_prefix; /* May be NULL */
    char const *is;               /* May be NULL */
    char const *local_name;
    YW_GC_PTR(struct yw_dom_shadow_root_node) shadow_root;
    void *tag_token;
    struct yw_dom_attr_list attrs;
    enum yw_custom_element_state custom_element_state;

    struct yw_dom_custom_element_registry custom_element_registry;
};

void yw_dom_element_init(YW_GC_PTR(struct yw_dom_element) out);
void yw_dom_element_visit(void *node_v);

/*******************************************************************************
 * CharacterData
 ******************************************************************************/

YW_GC_TYPE(struct yw_dom_character_data_node)
{
    YW_GC_TYPE(struct yw_dom_node) _node;
    char *text;
};

#endif /* #ifndef YW_DOM_H_ */
