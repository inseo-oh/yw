/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_HTML_TOKENS_H_
#define YW_HTML_TOKENS_H_

#include "yw_common.h"
#include "yw_dom.h"

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

struct YW_HTMLTokenizer;
typedef void(YW_HTMLTokenizerState)(struct YW_HTMLTokenizer *tkr);

typedef struct YW_HTMLTokenizer
{
    char *last_start_tag_name;
    YW_HTMLToken *current_token;

    char *temp_buf;
    /* Attributes that needs to be removed from current tag token */
    int *bad_attrs;
    int bad_attrs_len;
    int bad_attrs_cap;
    /* Capacity for current tag token's attribute array */
    int curr_attrs_cap;

    YW_TextReader tr;
    YW_HTMLTokenizerState *state;
    YW_HTMLTokenizerState *return_state;

    bool parser_pause_flag : 1;
} YW_HTMLTokenizer;

YW_HTMLTokenizerState yw_html_data_state;                  /* https://html.spec.whatwg.org/multipage/parsing.html#data-state */
YW_HTMLTokenizerState yw_html_rcdata_state;                /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-state */
YW_HTMLTokenizerState yw_html_rawtext_state;               /* https://html.spec.whatwg.org/multipage/parsing.html#rawtext-state */
YW_HTMLTokenizerState yw_html_plaintext_state;             /* https://html.spec.whatwg.org/multipage/parsing.html#plaintext-state */
YW_HTMLTokenizerState yw_html_tag_open_state;              /* https://html.spec.whatwg.org/multipage/parsing.html#tag-open-state */
YW_HTMLTokenizerState yw_html_end_tag_open_state;          /* https://html.spec.whatwg.org/multipage/parsing.html#end-tag-open-state */
YW_HTMLTokenizerState yw_html_tag_name_state;              /* https://html.spec.whatwg.org/multipage/parsing.html#tag-name-state */
YW_HTMLTokenizerState yw_html_rcdata_less_than_sign_state; /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-less-than-sign-state */
YW_HTMLTokenizerState yw_html_rcdata_end_tag_open_state;   /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-open-state */
YW_HTMLTokenizerState yw_html_rcdata_end_tag_name_state;   /* https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-name-state */
YW_HTMLTokenizerState yw_html_rawtext_less_than_sign_state;
YW_HTMLTokenizerState yw_html_rawtext_end_tag_open_state;
YW_HTMLTokenizerState yw_html_rawtext_end_tag_name_state;
YW_HTMLTokenizerState yw_html_before_attribute_name_state;
YW_HTMLTokenizerState yw_html_attribute_name_state;
YW_HTMLTokenizerState yw_html_after_attribute_name_state;
YW_HTMLTokenizerState yw_html_before_attribute_value_state;
YW_HTMLTokenizerState yw_html_attribute_value_double_quoted_state;
YW_HTMLTokenizerState yw_html_attribute_value_single_quoted_state;
YW_HTMLTokenizerState yw_html_attribute_value_unquoted_state;
YW_HTMLTokenizerState yw_html_after_attribute_value_quoted_state;
YW_HTMLTokenizerState yw_html_self_closing_start_tag_state;
YW_HTMLTokenizerState yw_html_bogus_comment_state;
YW_HTMLTokenizerState yw_html_markup_declaration_open_state;
YW_HTMLTokenizerState yw_html_comment_start_state;
YW_HTMLTokenizerState yw_html_comment_start_dash_state;
YW_HTMLTokenizerState yw_html_comment_state;
YW_HTMLTokenizerState yw_html_comment_less_than_sign_state;
YW_HTMLTokenizerState yw_html_comment_end_dash_state;
YW_HTMLTokenizerState yw_html_comment_end_state;
YW_HTMLTokenizerState yw_html_doctype_state;
YW_HTMLTokenizerState yw_html_before_doctype_name_state;
YW_HTMLTokenizerState yw_html_doctype_name_state;
YW_HTMLTokenizerState yw_html_after_doctype_name_state;
YW_HTMLTokenizerState yw_html_after_doctype_public_keyword_state;
YW_HTMLTokenizerState yw_html_before_doctype_public_identifier_state;
YW_HTMLTokenizerState yw_html_doctype_public_identifier_double_quoted_state;
YW_HTMLTokenizerState yw_html_doctype_public_identifier_single_quoted_state;
YW_HTMLTokenizerState yw_html_after_doctype_public_identifier_state;
YW_HTMLTokenizerState yw_html_between_doctype_public_and_system_identifiers_state;
YW_HTMLTokenizerState yw_html_after_doctype_system_keyword_state;
YW_HTMLTokenizerState yw_html_before_doctype_system_identifier_state;
YW_HTMLTokenizerState yw_html_doctype_system_identifier_double_quoted_state;
YW_HTMLTokenizerState yw_html_doctype_system_identifier_single_quoted_state;
YW_HTMLTokenizerState yw_html_after_doctype_system_identifier_state;
YW_HTMLTokenizerState yw_html_character_reference_state;
YW_HTMLTokenizerState yw_html_named_character_reference_state;
YW_HTMLTokenizerState yw_html_numeric_character_reference_state;
YW_HTMLTokenizerState yw_html_hexadecimal_character_reference_start_state;
YW_HTMLTokenizerState yw_html_decimal_character_reference_start_state;
YW_HTMLTokenizerState yw_html_hexadecimal_character_reference_state;
YW_HTMLTokenizerState yw_html_decimal_character_reference_state;
YW_HTMLTokenizerState yw_html_numeric_character_reference_end_state;

#endif
