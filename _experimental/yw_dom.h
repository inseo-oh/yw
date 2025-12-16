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

typedef YW_GC_TYPE(struct yw_dom_custom_element_registry)
    YW_GC_TYPE(yw_dom_custom_element_registry);

typedef struct yw_dom_node_callbacks yw_dom_node_callbacks;
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

typedef struct yw_dom_node_list yw_dom_node_list;
struct yw_dom_node_list
{
    YW_GC_PTR(struct yw_dom_node) *items;
    int len, cap;
};

typedef enum
{
    YW_DOM_TEXT_NODE,    /* yw_dom_character_data_node */
    YW_DOM_ELEMENT_NODE, /* yw_dom_element_node */

    YW_DOM_DOCUMENT_NODE,
    YW_DOM_DOCUMENT_FRAGMENT_NODE,
    YW_DOM_SHADOW_ROOT_NODE,
} yw_dom_node_type;

typedef YW_GC_TYPE(struct yw_dom_node) YW_GC_TYPE(yw_dom_node);
YW_GC_TYPE(struct yw_dom_node)
{
    struct yw_gc_object_header gc_header;
    uint32_t magic;

    YW_GC_PTR(yw_dom_node) parent;
    YW_GC_PTR(yw_dom_document) node_document;
    yw_dom_node_callbacks const *callbacks;
    yw_dom_node_list children;
    yw_dom_node_type type;
};

typedef enum
{
    YW_DOM_NO_SEARCH_FLAGS = 0,
    YW_DOM_SHADOW_INCLUDING = 1 << 0
} yw_dom_search_flags;

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
                                          yw_dom_search_flags flags);

int yw_dom_index(void *node_v);
bool yw_dom_is_in_the_same_tree_as(void *node_a_v, void *node_b_v);
bool yw_dom_is_connected(void *node_v);
char *yw_dom_child_text(void *node_v);

typedef struct yw_dom_iter yw_dom_iter;
struct yw_dom_iter
{
    YW_GC_PTR(yw_dom_node) root_node;
    YW_GC_PTR(yw_dom_node) last_node;
    bool shadow_including;
};

YW_GC_PTR(struct yw_dom_node) yw_dom_next_descendant(yw_dom_iter *iter);

void yw_dom_inclusive_descendants_init(yw_dom_iter *out, void *root_node_v,
                                       yw_dom_search_flags flags);
void yw_dom_descendants_init(yw_dom_iter *out, void *root_node_v,
                             yw_dom_search_flags flags);

typedef struct yw_dom_parents yw_dom_parents;
struct yw_dom_parents
{
    YW_GC_PTR(yw_dom_node) root_node;
    YW_GC_PTR(yw_dom_node) last_node;
    bool shadow_including;
};

YW_GC_PTR(struct yw_dom_node) yw_dom_next_parent(yw_dom_iter *iter);
void yw_dom_inclusive_ancestors_init(yw_dom_iter *out, void *root_node_v,
                                     yw_dom_search_flags flags);
void yw_dom_ancestors_init(yw_dom_iter *out, void *root_node_v,
                           yw_dom_search_flags flags);

typedef enum
{
    YW_DOM_NO_INSERT_FLAGS,
    YW_DOM_SUPPRESS_OBSERVERS = 1 << 0
} yw_dom_insert_flag;

void yw_dom_insert(void *node_v, void *parent_v, void *before_child_v,
                   yw_dom_insert_flag flags);
void yw_dom_append_child(void *node_v, void *child_v);
void yw_dom_adopt_into(void *node_v,
                       YW_GC_PTR(struct yw_dom_document) document);
void yw_dom_print_tree(FILE *dest, void *node_v, int indent_level);
YW_GC_PTR(yw_dom_custom_element_registry)
yw_dom_lookup_custom_element_registry(void *node_v);

/*******************************************************************************
 * Custom elements
 ******************************************************************************/
typedef YW_GC_TYPE(struct yw_dom_custom_element_registry)
    YW_GC_TYPE(yw_dom_custom_element_registry);
YW_GC_TYPE(struct yw_dom_custom_element_registry){};

/* https://dom.spec.whatwg.org/#concept-element-custom-element-state */
typedef enum
{
    YW_DOM_CUSTOM_ELEMENT_UNDEFINED,
    YW_DOM_CUSTOM_ELEMENT_FAILED,
    YW_DOM_CUSTOM_ELEMENT_UNCUSTOMIZED,
    YW_DOM_CUSTOM_ELEMENT_PRECUSTOMIZED,
    YW_DOM_CUSTOM_ELEMENT_CUSTOM,
} yw_custom_element_state;

/*******************************************************************************
 * ShadowRoot
 ******************************************************************************/
typedef YW_GC_TYPE(struct yw_dom_shadow_root_node)
    YW_GC_TYPE(yw_dom_shadow_root_node);
YW_GC_TYPE(struct yw_dom_shadow_root_node){
    /* STUB */
};

/*******************************************************************************
 * Attr
 ******************************************************************************/

typedef YW_GC_TYPE(struct yw_dom_element_node) YW_GC_TYPE(yw_dom_element_node);

typedef struct yw_dom_attr_data yw_dom_attr_data;
struct yw_dom_attr_data
{
    char const *local_name;
    char const *value;
    char const *namespace_;       /* May be NULL */
    char const *namespace_prefix; /* May be NULL */
};

typedef YW_GC_TYPE(struct yw_dom_attr_node) YW_GC_TYPE(yw_dom_attr_node);
YW_GC_TYPE(struct yw_dom_attr_node)
{
    YW_GC_TYPE(yw_dom_node) _node;

    char const *local_name;
    char const *value;
    char const *namespace_;       /* May be NULL */
    char const *namespace_prefix; /* May be NULL */
    YW_GC_PTR(yw_dom_element_node) element;
};

typedef struct yw_dom_attr_list yw_dom_attr_list;
struct yw_dom_attr_list
{
    YW_GC_PTR(yw_dom_attr_node) *items;
    int len, cap;
};

/*******************************************************************************
 * Element
 ******************************************************************************/

typedef YW_GC_TYPE(struct yw_dom_element_node) YW_GC_TYPE(yw_dom_element_node);
YW_GC_TYPE(struct yw_dom_element_node)
{
    YW_GC_TYPE(yw_dom_node) _node;

    char const *namespace_;       /* May be NULL */
    char const *namespace_prefix; /* May be NULL */
    char const *is;               /* May be NULL */
    char const *local_name;
    YW_GC_PTR(yw_dom_shadow_root_node) shadow_root;
    void *tag_token;
    yw_dom_attr_list attrs;
    yw_custom_element_state custom_element_state;

    YW_GC_PTR(yw_dom_custom_element_registry) custom_element_registry;
};

void yw_dom_element_init(YW_GC_PTR(yw_dom_element_node) out);
void yw_dom_element_visit(void *node_v);

bool yw_dom_is_shadow_host(void *node_v);
bool yw_dom_is_element_defined(void *node_v);
bool yw_dom_is_element_custom(void *node_v);
bool yw_dom_is_element_inside(void *node_v, char const *namespace_,
                              char const *local_name);
bool yw_dom_is_element(void *node_v, char const *namespace_,
                       char const *local_name);
bool yw_dom_is_html_element(void *node_v, char const *local_name);
bool yw_dom_is_mathml_element(void *node_v, char const *local_name);
bool yw_dom_is_svg_element(void *node_v, char const *local_name);
void yw_dom_append_attr(void *node_v, yw_gc_heap *heap,
                        yw_dom_attr_data const *data);

/*
 * If namespace is NULL, attributes with namespaces will not be matched.
 * If namespace is non-NULL, attributes without namespaces will not be matched.
 *
 */
char const *yw_dom_attr(void *node_v, char const *namespace_,
                        char const *local_name);

/*******************************************************************************
 * CharacterData
 ******************************************************************************/

typedef YW_GC_TYPE(struct yw_dom_character_data_node)
    YW_GC_TYPE(yw_dom_character_data_node);
YW_GC_TYPE(struct yw_dom_character_data_node)
{
    YW_GC_TYPE(yw_dom_node) _node;
    char *text;
};

#endif /* #ifndef YW_DOM_H_ */
