/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#ifndef YW_HTML_TOKENS_H_
#define YW_HTML_TOKENS_H_

#include "yw_common.h"
#include "yw_dom.h"
#include <stdint.h>

typedef enum
{
    YW_HTML_EOF_TOKEN,
    YW_HTML_CHAR_TOKEN,
    YW_HTML_COMMENT_TOKEN,
    YW_HTML_DOCTYPE_TOKEN,
    YW_HTML_TAG_TOKEN,
} YW_HTMLTokenType;

typedef struct YW_HTMLCharToken
{
    YW_HTMLTokenType type; /* YW_HTML_CHAR_TOKEN */
    YW_Char32 chr;
} YW_HTMLCharToken;
typedef struct YW_HTMLCommentToken
{
    YW_HTMLTokenType type; /* YW_HTML_COMMENT_TOKEN */
    char *data;
} YW_HTMLCommentToken;
typedef struct YW_HTMLDoctypeToken
{
    YW_HTMLTokenType type; /* YW_HTML_DOCTYPE_TOKEN */
    char *name;            /* NULL = missing */
    char *public_id;       /* NULL = missing */
    char *system_id;       /* NULL = missing */

    bool force_quirks : 1;
} YW_HTMLDoctypeToken;
typedef struct YW_HTMLTagToken
{
    YW_HTMLTokenType type; /* YW_HTML_TAG_TOKEN */
    char *name;
    YW_DOMAttrData *attrs;
    int attrs_len;

    bool is_end : 1;
    bool is_self_closing : 1;
} YW_HTMLTagToken;

typedef union YW_HTMLToken {
    YW_HTMLTokenType type;

    YW_HTMLCharToken char_tk;
    YW_HTMLCommentToken comment_tk;
    YW_HTMLDoctypeToken doctype_tk;
    YW_HTMLTagToken tag_tk;
} YW_HTMLToken;

void yw_html_token_deinit(YW_HTMLToken *tk);
void yw_html_token_clone(YW_HTMLToken *dest, YW_HTMLToken const *tk);

struct YW_HTMLTokenizer;
typedef void(YW_HTMLTokenizerState)(struct YW_HTMLTokenizer *tkr);

typedef struct YW_HTMLTokenizer
{
    char *last_start_tag_name;
    YW_HTMLToken *current_token;

    void (*emit_callback)(YW_HTMLToken *token, void *callback_data);
    void *emit_callback_data;
    YW_HTMLTokenizerState *state;
    YW_HTMLTokenizerState *return_state;
    char *temp_buf;
    /* Attributes that needs to be removed from current tag token */
    int *bad_attrs;
    int bad_attrs_len;
    int bad_attrs_cap;
    /* Capacity for current tag token's attribute array */
    int curr_attrs_cap;

    YW_Char32 character_reference_code;

    YW_TextReader tr;

    bool parser_pause_flag : 1;
    bool eof_emitted : 1;
} YW_HTMLTokenizer;

void yw_html_tokenize(
    uint8_t const *chars, int chars_len,
    void (*emit_callback)(YW_HTMLToken *token, void *callback_data), void *emit_callback_data);

/*******************************************************************************
 * Below should not be called directly!
 ******************************************************************************/

#define YW_HTML_ENUMERATE_TOKENIZER_STATE(_x)                                                                                                                                         \
    _x(yw_html_data_state)                                              /* https://html.spec.whatwg.org/multipage/parsing.html#data-state */                                          \
        _x(yw_html_rcdata_state)                                        /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-state */                                        \
        _x(yw_html_rawtext_state)                                       /* https://html.spec.whatwg.org/multipage/parsing.html#rawtext-state */                                       \
        _x(yw_html_plaintext_state)                                     /* https://html.spec.whatwg.org/multipage/parsing.html#plaintext-state */                                     \
        _x(yw_html_tag_open_state)                                      /* https://html.spec.whatwg.org/multipage/parsing.html#tag-open-state */                                      \
        _x(yw_html_end_tag_open_state)                                  /* https://html.spec.whatwg.org/multipage/parsing.html#end-tag-open-state */                                  \
        _x(yw_html_tag_name_state)                                      /* https://html.spec.whatwg.org/multipage/parsing.html#tag-name-state */                                      \
        _x(yw_html_rcdata_less_than_sign_state)                         /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-less-than-sign-state */                         \
        _x(yw_html_rcdata_end_tag_open_state)                           /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-open-state */                           \
        _x(yw_html_rcdata_end_tag_name_state)                           /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-name-state */                           \
        _x(yw_html_rawtext_less_than_sign_state)                        /* https://html.spec.whatwg.org/multipage/parsing.html#rawtext-less-than-sign-state */                        \
        _x(yw_html_rawtext_end_tag_open_state)                          /* https://html.spec.whatwg.org/multipage/parsing.html#rawtext-end-tag-open-state */                          \
        _x(yw_html_rawtext_end_tag_name_state)                          /* https://html.spec.whatwg.org/multipage/parsing.html#rawtext-end-tag-name-state */                          \
        _x(yw_html_before_attribute_name_state)                         /* https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-name-state */                         \
        _x(yw_html_attribute_name_state)                                /* https://html.spec.whatwg.org/multipage/parsing.html#attribute-name-state */                                \
        _x(yw_html_after_attribute_name_state)                          /* https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-name-state */                          \
        _x(yw_html_before_attribute_value_state)                        /* https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-value-state */                        \
        _x(yw_html_attribute_value_double_quoted_state)                 /* https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(double-quoted)-state */               \
        _x(yw_html_attribute_value_single_quoted_state)                 /* https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(single-quoted)-state */               \
        _x(yw_html_attribute_value_unquoted_state)                      /* https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(unquoted)-state */                    \
        _x(yw_html_after_attribute_value_quoted_state)                  /* https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-value-(quoted)-state */                \
        _x(yw_html_self_closing_start_tag_state)                        /* https://html.spec.whatwg.org/multipage/parsing.html#self-closing-start-tag-state */                        \
        _x(yw_html_bogus_comment_state)                                 /* https://html.spec.whatwg.org/multipage/parsing.html#bogus-comment-state */                                 \
        _x(yw_html_markup_declaration_open_state)                       /* https://html.spec.whatwg.org/multipage/parsing.html#markup-declaration-open-state */                       \
        _x(yw_html_comment_start_state)                                 /* https://html.spec.whatwg.org/multipage/parsing.html#comment-start-state */                                 \
        _x(yw_html_comment_start_dash_state)                            /* https://html.spec.whatwg.org/multipage/parsing.html#comment-start-dash-state */                            \
        _x(yw_html_comment_state)                                       /* https://html.spec.whatwg.org/multipage/parsing.html#comment-state */                                       \
        _x(yw_html_comment_less_than_sign_state)                        /* https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-state */                        \
        _x(yw_html_comment_end_dash_state)                              /* https://html.spec.whatwg.org/multipage/parsing.html#comment-end-dash-state */                              \
        _x(yw_html_comment_end_state)                                   /* https://html.spec.whatwg.org/multipage/parsing.html#comment-end-state */                                   \
        _x(yw_html_doctype_state)                                       /* https://html.spec.whatwg.org/multipage/parsing.html#doctype-state */                                       \
        _x(yw_html_before_doctype_name_state)                           /* https://html.spec.whatwg.org/multipage/parsing.html#before-doctype-name-state */                           \
        _x(yw_html_doctype_name_state)                                  /* https://html.spec.whatwg.org/multipage/parsing.html#doctype-name-state */                                  \
        _x(yw_html_after_doctype_name_state)                            /* https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state */                            \
        _x(yw_html_after_doctype_public_keyword_state)                  /* https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-public-keyword-state */                  \
        _x(yw_html_before_doctype_public_identifier_state)              /* https://html.spec.whatwg.org/multipage/parsing.html#before-doctype-public-identifier-state */              \
        _x(yw_html_doctype_public_identifier_double_quoted_state)       /* https://html.spec.whatwg.org/multipage/parsing.html#doctype-public-identifier-(double-quoted)-state */     \
        _x(yw_html_doctype_public_identifier_single_quoted_state)       /* https://html.spec.whatwg.org/multipage/parsing.html#doctype-public-identifier-(single-quoted)-state */     \
        _x(yw_html_after_doctype_public_identifier_state)               /* https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-public-identifier-state */               \
        _x(yw_html_between_doctype_public_and_system_identifiers_state) /* https://html.spec.whatwg.org/multipage/parsing.html#between-doctype-public-and-system-identifiers-state */ \
        _x(yw_html_after_doctype_system_keyword_state)                  /* https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-system-keyword-state */                  \
        _x(yw_html_before_doctype_system_identifier_state)              /* https://html.spec.whatwg.org/multipage/parsing.html#before-doctype-system-identifier-state*/               \
        _x(yw_html_doctype_system_identifier_double_quoted_state)       /* https://html.spec.whatwg.org/multipage/parsing.html#doctype-system-identifier-(double-quoted)-state */     \
        _x(yw_html_doctype_system_identifier_single_quoted_state)       /* https://html.spec.whatwg.org/multipage/parsing.html#doctype-system-identifier-(single-quoted)-state */     \
        _x(yw_html_after_doctype_system_identifier_state)               /* https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-system-identifier-state */               \
        _x(yw_html_character_reference_state)                           /* https://html.spec.whatwg.org/multipage/parsing.html#character-reference-state */                           \
        _x(yw_html_named_character_reference_state)                     /* https://html.spec.whatwg.org/multipage/parsing.html#named-character-reference-state */                     \
        _x(yw_html_numeric_character_reference_state)                   /* https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-state */                   \
        _x(yw_html_hexadecimal_character_reference_start_state)         /* https://html.spec.whatwg.org/multipage/parsing.html#hexadecimal-character-reference-start-state */         \
        _x(yw_html_decimal_character_reference_start_state)             /* https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-start-state */             \
        _x(yw_html_hexadecimal_character_reference_state)               /* https://html.spec.whatwg.org/multipage/parsing.html#hexadecimal-character-reference-state */               \
        _x(yw_html_decimal_character_reference_state)                   /* https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-state */                   \
        _x(yw_html_numeric_character_reference_end_state)

#define YW_X(_x) YW_HTMLTokenizerState _x;
YW_HTML_ENUMERATE_TOKENIZER_STATE(YW_X)
#undef YW_X

#endif
