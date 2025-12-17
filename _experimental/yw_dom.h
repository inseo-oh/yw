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

typedef struct YW_Attr_Rec YW_GC_TYPE(YW_Attr);
typedef struct YW_CharacterData_Rec YW_GC_TYPE(YW_CharacterData);
typedef struct
    YW_CustomElementRegistry_Rec YW_GC_TYPE(YW_CustomElementRegistry);
typedef struct YW_Document_Rec YW_GC_TYPE(YW_Document);
typedef struct YW_DocumentFragment_Rec YW_GC_TYPE(YW_DocumentFragment);
typedef struct YW_DocumentType_Rec YW_GC_TYPE(YW_DocumentType);
typedef struct YW_DomNode_Rec YW_GC_TYPE(YW_DomNode);
typedef struct YW_Element_Rec YW_GC_TYPE(YW_Element);
typedef struct YW_ShadowRoot_Rec YW_GC_TYPE(YW_ShadowRoot);

typedef struct YW_DomNodeList YW_DomNodeList;
typedef struct YW_DocumentList YW_DocumentList;
typedef struct YW_AttrList YW_AttrList;

typedef struct YW_AttrData YW_AttrData;
typedef struct YW_CustomElementDefinition YW_CustomElementDefinition;
typedef struct YW_DomIter YW_DomIter;
typedef struct YW_DomNodeCallbacks_Rec YW_DomNodeCallbacks;

/*******************************************************************************
 * Lists
 ******************************************************************************/

struct YW_DomNodeList
{
    YW_GC_PTR(YW_DomNode) *items;
    int len, cap;
};

struct YW_DocumentList
{
    YW_GC_PTR(YW_Document) *items;
    int len, cap;
};

struct YW_AttrList
{
    YW_GC_PTR(YW_Attr) *items;
    int len, cap;
};

/*******************************************************************************
 * Node
 ******************************************************************************/

struct YW_DomNodeCallbacks_Rec
{
    void (*run_insertion_steps)(void *self);
    void (*run_children_changed_steps)(void *self);
    void (*run_post_connection_steps)(void *self);
    void (*run_adopting_steps)(void *self, YW_GC_PTR(YW_Document) old_document);

    /* Element callbacks ******************************************************/

    void (*intrinsic_size)(double *width_out, double *height_out, void *self);
    void (*popped_from_stack_of_open_elements)(void *self);
    void (*presentational_hints)(void *self);
    void (*run_form_reset_algorithm)(void *self);
};

typedef enum
{
    YW_DOM_TEXT_NODE = 1 << 0,
    YW_DOM_ELEMENT_NODE = 1 << 1,
    YW_DOM_DOCUMENT_NODE = 1 << 2,
    YW_DOM_DOCUMENT_FRAGMENT_NODE = 1 << 3,
    YW_DOM_SHADOW_ROOT_NODE = YW_DOM_DOCUMENT_FRAGMENT_NODE | (1 << 4),
    YW_DOM_DOCUMENT_TYPE_NODE = 1 << 5,
    YW_DOM_ATTR_NODE = 1 << 6,
} YW_DomNodeTypeFlags;

struct YW_DomNode_Rec
{
    struct YW_GcObjectHeader gc_header;
    uint32_t magic;

    YW_GC_PTR(YW_DomNode) parent;
    YW_GC_PTR(YW_Document) node_document;
    YW_DomNodeCallbacks const *callbacks;
    YW_DomNodeList children;
    uint8_t type_flags;
};

typedef enum
{
    YW_DOM_NO_SEARCH_FLAGS = 0,
    YW_DOM_SHADOW_INCLUDING = 1 << 0
} YW_DomSearchFlags;

void yw_dom_node_init(YW_GC_PTR(YW_DomNode) out);
void yw_dom_node_visit(void *node_v);
void yw_dom_node_destroy(void *node_v);

/* Returns NULL if there's no children. */
YW_GC_PTR(YW_DomNode) yw_dom_first_child(void *node_v);
/* Returns NULL if there's no children. */
YW_GC_PTR(YW_DomNode) yw_dom_last_child(void *node_v);
/* Returns NULL if there's no parent. */
YW_GC_PTR(YW_DomNode) yw_dom_next_sibling(void *node_v);
/* Returns NULL if there's no parent. */
YW_GC_PTR(YW_DomNode) yw_dom_prev_sibling(void *node_v);
YW_GC_PTR(YW_DomNode) yw_dom_root(void *node_v, YW_DomSearchFlags flags);

int yw_dom_index(void *node_v);
bool yw_dom_has_type(void *node_v, YW_DomNodeTypeFlags flags);
bool yw_dom_is_in_same_tree(void *node_a_v, void *node_b_v);
bool yw_dom_is_connected(void *node_v);
char *yw_dom_child_text(void *node_v);

struct YW_DomIter
{
    YW_GC_PTR(YW_DomNode) root_node;
    YW_GC_PTR(YW_DomNode) last_node;
    bool shadow_including;
};

YW_GC_PTR(YW_DomNode) yw_dom_next_descendant(YW_DomIter *iter);

void yw_dom_inclusive_descendants_init(YW_DomIter *out, void *root_node_v,
                                       YW_DomSearchFlags flags);
void yw_dom_descendants_init(YW_DomIter *out, void *root_node_v,
                             YW_DomSearchFlags flags);

YW_GC_PTR(YW_DomNode) yw_dom_next_ancestor(YW_DomIter *iter);
void yw_dom_inclusive_ancestors_init(YW_DomIter *out, void *root_node_v,
                                     YW_DomSearchFlags flags);
void yw_dom_ancestors_init(YW_DomIter *out, void *root_node_v,
                           YW_DomSearchFlags flags);

typedef enum
{
    YW_DOM_NO_INSERT_FLAGS,
    YW_DOM_SUPPRESS_OBSERVERS = 1 << 0
} YW_DomInsertFlag;

void yw_dom_insert(void *node_v, void *parent_v, void *before_child_v,
                   YW_DomInsertFlag flags);
void yw_dom_append_child(void *node_v, void *child_v);
void yw_dom_adopt_into(void *node_v, YW_GC_PTR(YW_Document) document);
void yw_dom_print_tree(FILE *dest, void *node_v, int indent_level);
YW_GC_PTR(YW_CustomElementRegistry)
yw_lookup_custom_element_registry(void *node_v);

/*******************************************************************************
 * Custom elements
 ******************************************************************************/
struct YW_CustomElementRegistry_Rec
{
    /*
     * https://html.spec.whatwg.org/multipage/custom-elements.html#scoped-document-set
     */
    YW_DocumentList scoped_document_set;

    /*
     * https://html.spec.whatwg.org/multipage/custom-elements.html#is-scoped
     */
    bool is_scoped;
};

/* https://dom.spec.whatwg.org/#concept-element-custom-element-state */
typedef enum
{
    YW_CUSTOM_ELEMENT_UNDEFINED,
    YW_CUSTOM_ELEMENT_FAILED,
    YW_CUSTOM_ELEMENT_UNCUSTOMIZED,
    YW_CUSTOM_ELEMENT_PRECUSTOMIZED,
    YW_CUSTOM_ELEMENT_CUSTOM,
} YW_CustomElementState;

/*
 * https://html.spec.whatwg.org/multipage/custom-elements.html#custom-element-definition
 */
struct YW_CustomElementDefinition
{
    int dummy;
};

YW_CustomElementDefinition const *yw_dom_lookup_custom_element_definition(
    YW_GC_PTR(YW_CustomElementRegistry) registry, char const *namespace_,
    char const *local_name, char const *is);
bool yw_dom_is_global_custom_element_reigstry(
    YW_GC_PTR(YW_CustomElementRegistry) registry);
void yw_dom_try_upgrade_element(void *node_v);

/*******************************************************************************
 * Document
 ******************************************************************************/

/*
 * https://dom.spec.whatwg.org/#concept-document-mode
 */
typedef enum
{
    YW_NO_QUIRKS,
    YW_QUIRKS,
    YW_LIMITED_QUIRKS,
} YW_DocumentMode;

struct YW_Document_Rec
{
    YW_GC_TYPE(YW_DomNode) _node;

    YW_GC_PTR(YW_CustomElementRegistry) custom_element_registry;

    /*
     * FIXME: This should not be a field!
     * See:
     * https://html.spec.whatwg.org/multipage/urls-and-fetching.html#document-base-url
     */
    char const *base_url;

    struct
    {
        int dummy;
    } origin, environment_settings, policy_container; /* STUB */
    YW_DocumentMode mode;

    bool iframe_srcdoc_document;
    bool parser_cannot_change_mode;
};

void yw_document_init(YW_GC_PTR(YW_Document) out);
YW_GC_PTR(YW_Document) yw_document_alloc(YW_GcHeap *heap,
                                         YW_GcAllocFlags alloc_flags);
void yw_document_visit(void *node_v);
void yw_document_destroy(void *node_v);

YW_GC_PTR(YW_CustomElementRegistry)
yw_document_effective_global_custom_element_registry(void *node_v);

/*******************************************************************************
 * DocumentFragment
 ******************************************************************************/

struct YW_DocumentFragment_Rec
{
    YW_GC_TYPE(YW_DomNode) _node;

    YW_GC_PTR(YW_DomNode) host;
};

void yw_document_fragment_init(YW_GC_PTR(YW_DocumentFragment) out);
YW_GC_PTR(YW_DocumentFragment) yw_document_fragment_alloc(
    YW_GcHeap *heap, YW_GcAllocFlags alloc_flags);
void yw_document_fragment_visit(void *node_v);
void yw_document_fragment_destroy(void *node_v);

/*******************************************************************************
 * DocumentType
 ******************************************************************************/

struct YW_DocumentType_Rec
{
    YW_GC_TYPE(YW_DomNode) _node;

    char *name;
    char *public_id;
    char *system_id;
};

void yw_document_type_init(YW_GC_PTR(YW_DocumentType) out);
YW_GC_PTR(YW_DocumentType) yw_document_type_alloc(YW_GcHeap *heap,
                                                  YW_GcAllocFlags alloc_flags);
void yw_document_type_visit(void *node_v);
void yw_document_type_destroy(void *node_v);

/*******************************************************************************
 * Attr
 ******************************************************************************/

struct YW_AttrData
{
    char const *local_name;
    char const *value;
    char const *namespace_;       /* May be NULL */
    char const *namespace_prefix; /* May be NULL */
};

struct YW_Attr_Rec
{
    YW_GC_TYPE(YW_DomNode) _node;

    YW_GC_PTR(YW_Element) element;

    char *local_name;
    char *value;
    char *namespace_;       /* May be NULL */
    char *namespace_prefix; /* May be NULL */
};

void yw_attr_init(YW_GC_PTR(YW_Attr) out);
YW_GC_PTR(YW_Attr) yw_attr_alloc(YW_GcHeap *heap, YW_GcAllocFlags alloc_flags);
void yw_attr_visit(void *node_v);
void yw_attr_destroy(void *node_v);

/*******************************************************************************
 * Element
 ******************************************************************************/

struct YW_Element_Rec
{
    YW_GC_TYPE(YW_DomNode) _node;

    YW_GC_PTR(YW_ShadowRoot) shadow_root;
    YW_GC_PTR(YW_CustomElementRegistry) custom_element_registry;

    char const *namespace_;       /* May be NULL */
    char const *namespace_prefix; /* May be NULL */
    char const *is;               /* May be NULL */
    char const *local_name;
    void *tag_token;

    YW_AttrList attrs;
    YW_CustomElementState custom_element_state;
};

void yw_element_init(YW_GC_PTR(YW_Element) out);
YW_GC_PTR(YW_Element) yw_element_alloc(YW_GcHeap *heap,
                                       YW_GcAllocFlags alloc_flags);
void yw_element_visit(void *node_v);
void yw_element_destroy(void *node_v);

/*
 * yw_dom_is~ functions return false if it's not an element.
 * (meaning it is safe to pass non-element nodes)
 */

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

void yw_dom_append_attr(void *node_v, YW_GcHeap *heap, YW_AttrData const *data);

/*
 * If namespace is NULL, attributes with namespaces will not be matched.
 * If namespace is non-NULL, attributes without namespaces will not be
 * matched.
 */
char const *yw_attr(void *node_v, char const *namespace_,
                    char const *local_name);

/*******************************************************************************
 * CharacterData
 ******************************************************************************/

struct YW_CharacterData_Rec
{
    YW_GC_TYPE(YW_DomNode) _node;

    char *text;
};

void yw_dom_character_data_visit(void *node_v);
void yw_dom_character_data_destroy(void *node_v);
void yw_dom_text_init(YW_GC_PTR(YW_CharacterData) out);
YW_GC_PTR(YW_CharacterData) yw_text_alloc(YW_GcHeap *heap,
                                          YW_GcAllocFlags alloc_flags);

#endif /* #ifndef YW_DOM_H_ */
