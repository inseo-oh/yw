/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_html_tokens.h"
#include "yw_common.h"
#include "yw_dom.h"
#include "yw_encoding.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define YW_UNICODE_REPLACEMENT_CHAR "\xef\xbf\xbd" /* U+FFFD in UTF-8 */

typedef enum
{
    YW_ABSENCE_OF_DIGITS_IN_NUMERIC_CHARACTER_REFERENCE_ERROR,
    YW_ABRUPT_CLOSING_OF_EMPTY_COMMENT_ERROR,
    YW_ABRUPT_DOCTYPE_PUBLIC_IDENTIFIER_ERROR,
    YW_ABRUPT_DOCTYPE_SYSTEM_IDENTIFIER_ERROR,
    YW_CHARACTER_REFERENCE_OUTSIDE_UNICODE_RANGE_ERROR,
    YW_CONTROL_CHARACTER_REFERENCE_ERROR,
    YW_EOF_BEFORE_TAG_NAME_ERROR,
    YW_EOF_IN_COMMENT_ERROR,
    YW_EOF_IN_DOCTYPE_ERROR,
    YW_EOF_IN_TAG_ERROR,
    YW_INCORRECTLY_OPENED_COMMENT_ERROR,
    YW_INVALID_CHARACTER_SEQUENCE_AFTER_DOCTYPE_NAME_ERROR,
    YW_INVALID_FIRST_CHARACTER_OF_TAG_NAME_ERROR,
    YW_MISSING_ATTRIBUTE_VALUE_ERROR,
    YW_MISSING_DOCTYPE_NAME_ERROR,
    YW_MISSING_DOCTYPE_PUBLIC_IDENTIFIER_ERROR,
    YW_MISSING_DOCTYPE_SYSTEM_IDENTIFIER_ERROR,
    YW_MISSING_END_TAG_NAME_ERROR,
    YW_MISSING_SEMICOLON_AFTER_CHARACTER_REFERENCE_ERROR,
    YW_MISSING_QUOTE_BEFORE_DOCTYPE_PUBLIC_IDENTIFIER_ERROR,
    YW_MISSING_QUOTE_BEFORE_DOCTYPE_SYSTEM_IDENTIFIER_ERROR,
    YW_MISSING_WHITESPACE_AFTER_DOCTYPE_PUBLIC_KEYWORD_ERROR,
    YW_MISSING_WHITESPACE_AFTER_DOCTYPE_SYSTEM_KEYWORD_ERROR,
    YW_MISSING_WHITESPACE_BEFORE_DOCTYPE_NAME_ERROR,
    YW_MISSING_WHITESPACE_BETWEEN_ATTRIBUTES_ERROR,
    YW_MISSING_WHITESPACE_BETWEEN_DOCTYPE_PUBLIC_AND_SYSTEM_IDENTIFIERS_ERROR,
    YW_NONCHARACTER_REFERENCE_ERROR,
    YW_NULL_CHARACTER_REFERENCE_ERROR,
    YW_SURROGATE_CHARACTER_REFERENCE_ERROR,
    YW_UNEXPECTED_CHARACTER_IN_ATTRIBUTE_NAME_ERROR,
    YW_UNEXPECTED_CHARACTER_IN_UNQUOTED_ATTRIBUTE_VALUE_ERROR,
    YW_UNEXPECTED_EQUALS_SIGN_BEFORE_ATTRIBUTE_NAME_ERROR,
    YW_UNEXPECTED_NULL_CHARACTER_ERROR,
    YW_UNEXPECTED_QUESTION_MARK_INSTEAD_OF_TAG_NAME_ERROR,
    YW_UNEXPECTED_SOLIDUS_IN_TAG_ERROR,
    YW_UNEXPECTED_CHARACTER_AFTER_DOCTYPE_SYSTEM_IDENTIFIER_ERROR,
} YW_HTMLParseError;

void yw_html_token_deinit(YW_HTMLToken *tk)
{
    switch (tk->type)
    {
    case YW_HTML_EOF_TOKEN:
    case YW_HTML_CHAR_TOKEN:
        break;
    case YW_HTML_COMMENT_TOKEN:
        free(tk->comment_tk.data);
        break;
    case YW_HTML_DOCTYPE_TOKEN:
        free(tk->doctype_tk.name);
        free(tk->doctype_tk.public_id);
        free(tk->doctype_tk.system_id);
        break;
    case YW_HTML_TAG_TOKEN: {
        free(tk->tag_tk.name);
        for (int i = 0; i < tk->tag_tk.attrs_len; i++)
        {
            yw_dom_attr_data_deinit(&tk->tag_tk.attrs[i]);
        }
        free(tk->tag_tk.attrs);
        break;
    }
    }
}

static bool yw_is_start_tag(YW_HTMLToken const *tk)
{
    return tk->type == YW_HTML_TAG_TOKEN && !tk->tag_tk.is_end;
}
static bool yw_is_end_tag(YW_HTMLToken const *tk)
{
    return tk->type == YW_HTML_TAG_TOKEN && tk->tag_tk.is_end;
}

static YW_HTMLToken *yw_make_eof_token()
{
    YW_HTMLToken *res = YW_ALLOC(YW_HTMLToken, 1);
    memset(res, 0, sizeof(*res));
    res->type = YW_HTML_EOF_TOKEN;
    return res;
}
static YW_HTMLToken *yw_make_char_token(YW_Char32 chr)
{
    YW_HTMLToken *res = YW_ALLOC(YW_HTMLToken, 1);
    memset(res, 0, sizeof(*res));
    res->type = YW_HTML_CHAR_TOKEN;
    res->char_tk.chr = chr;
    return res;
}
static YW_HTMLToken *yw_make_comment_token(char const *s)
{
    YW_HTMLToken *res = YW_ALLOC(YW_HTMLToken, 1);
    memset(res, 0, sizeof(*res));
    res->type = YW_HTML_COMMENT_TOKEN;
    res->comment_tk.data = yw_duplicate_str(s);
    return res;
}
static YW_HTMLToken *yw_make_tag_token(char const *name)
{
    YW_HTMLToken *res = YW_ALLOC(YW_HTMLToken, 1);
    memset(res, 0, sizeof(*res));
    res->type = YW_HTML_TAG_TOKEN;
    res->tag_tk.name = yw_duplicate_str(name);
    return res;
}

void yw_html_tokenizer_init(YW_HTMLTokenizer *out, const uint8_t *chars, int chars_len)
{
    memset(out, 0, sizeof(*out));

    yw_text_reader_init(&out->tr, chars, chars_len);
}

static void yw_set_current_token(YW_HTMLTokenizer *tkr, YW_HTMLToken *tk)
{
    if (tkr->current_token != NULL)
    {
        YW_UNREACHABLE();
    }
    tkr->current_token = tk;
}
static YW_HTMLTagToken *yw_current_tag_token(YW_HTMLTokenizer *tkr)
{
    if (tkr->current_token->type != YW_HTML_TAG_TOKEN)
    {
        YW_UNREACHABLE();
    }
    return &tkr->current_token->tag_tk;
}
static YW_HTMLCommentToken *yw_current_comment_token(YW_HTMLTokenizer *tkr)
{
    if (tkr->current_token->type != YW_HTML_COMMENT_TOKEN)
    {
        YW_UNREACHABLE();
    }
    return &tkr->current_token->comment_tk;
}
static YW_HTMLDoctypeToken *yw_current_doctype_token(YW_HTMLTokenizer *tkr)
{
    if (tkr->current_token->type != YW_HTML_DOCTYPE_TOKEN)
    {
        YW_UNREACHABLE();
    }
    return &tkr->current_token->doctype_tk;
}
static YW_DOMAttrData *yw_current_attr(YW_HTMLTokenizer *tkr)
{
    YW_HTMLTagToken *tag = yw_current_tag_token(tkr);
    if (tag->attrs_len < 1)
    {
        YW_UNREACHABLE();
    }
    return &tag->attrs[tag->attrs_len - 1];
}
static void yw_check_duplicate_attr_name(YW_HTMLTokenizer *tkr, YW_DOMAttrData const *attr)
{
    YW_HTMLTagToken *tag = yw_current_tag_token(tkr);
    for (int i = 0; i < tag->attrs_len; i++)
    {
        if (strcmp(attr->local_name, tag->attrs[i].local_name) == 0)
        {
            YW_PUSH(int, &tkr->bad_attrs_cap, &tkr->bad_attrs_len, &tkr->bad_attrs, i);
        }
    }
}

/*
 * This function will write NULL to the pointer.
 * NOTE: If it is a tag token, it must be the current token.
 */
static void yw_emit_token(YW_HTMLTokenizer *tkr, YW_HTMLToken **tk_inout)
{
    YW_HTMLToken *tk = *tk_inout;
    if (tk->type == YW_HTML_TAG_TOKEN)
    {
        if (tk != tkr->current_token)
        {
            YW_UNREACHABLE();
        }
        if (tkr->bad_attrs_len != 0)
        {
            YW_DOMAttrData *new_attrs = NULL;
            int new_attrs_len = 0;
            int new_attrs_cap = 0;
            for (int i = 0; i < tk->tag_tk.attrs_len; i++)
            {
                bool is_bad_attr = false;
                for (int j = 0; j < tkr->bad_attrs_len; i++)
                {
                    if (tkr->bad_attrs[j] == i)
                    {
                        is_bad_attr = true;
                        break;
                    }
                }
                if (!is_bad_attr)
                {
                    YW_PUSH(YW_DOMAttrData, &new_attrs_cap, &new_attrs_len, &new_attrs, tk->tag_tk.attrs[i]);
                }
            }
            free(tk->tag_tk.attrs);
            tk->tag_tk.attrs = new_attrs;
            tk->tag_tk.attrs_len = new_attrs_len;
            tkr->curr_attrs_cap = new_attrs_cap;
        }
        YW_SHRINK_TO_FIT(YW_DOMAttrData, &tkr->curr_attrs_cap, tk->tag_tk.attrs_len, &tk->tag_tk.attrs);
        if (yw_is_start_tag(tk))
        {
            tkr->last_start_tag_name = yw_duplicate_str(tk->tag_tk.name);
        }
    }
    yw_html_token_deinit(tk);
    free(tk);
    *tk_inout = NULL;
}

static void yw_emit_eof_token(YW_HTMLTokenizer *tkr)
{
    YW_HTMLToken *token = yw_make_eof_token();
    yw_emit_token(tkr, &token);
}
static void yw_emit_char_token(YW_HTMLTokenizer *tkr, YW_Char32 chr)
{
    YW_HTMLToken *token = yw_make_char_token(chr);
    yw_emit_token(tkr, &token);
}

static bool yw_is_consumed_as_part_of_attr(YW_HTMLTokenizer const *tkr)
{
    if (tkr->return_state == yw_html_attribute_value_double_quoted_state ||
        tkr->return_state == yw_html_attribute_value_single_quoted_state ||
        tkr->return_state == yw_html_attribute_value_unquoted_state)
    {
        return true;
    }

    return false;
}

static void yw_flush_codepoints_consumed_as_char_reference(YW_HTMLTokenizer *tkr)
{
    if (yw_is_consumed_as_part_of_attr(tkr))
    {
        yw_append_str(&yw_current_attr(tkr)->value, tkr->temp_buf);
    }
    else
    {
        for (char const *src = tkr->temp_buf; *src != '\0'; src++)
        {
            yw_emit_char_token(tkr, *src);
        }
    }
}

static bool yw_is_appropriate_end_tag_token(YW_HTMLTokenizer const *tkr, YW_HTMLToken const *tk)
{
    if (!yw_is_end_tag(tk))
    {
        return false;
    }
    return strcmp(tkr->last_start_tag_name, tk->tag_tk.name) == 0;
}

static void yw_parse_error_encountered(YW_HTMLTokenizer const *tkr, YW_HTMLParseError err)
{
    (void)tkr;
    fprintf(stderr, "HTML parse error %d\n", err);
}

void yw_html_data_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '&':
        tkr->return_state = yw_html_data_state;
        tkr->state = yw_html_character_reference_state;
        break;
    case '<':
        tkr->state = yw_html_tag_open_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_emit_char_token(tkr, next_char);
        break;
    case -1:
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_emit_char_token(tkr, next_char);
        break;
    }
}
void yw_html_rcdata_state(struct YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '&':
        tkr->return_state = yw_html_rcdata_state;
        tkr->state = yw_html_character_reference_state;
        break;
    case '<':
        tkr->state = yw_html_rcdata_less_than_sign_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_emit_char_token(tkr, 0xfffd);
        break;
    case -1:
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_emit_char_token(tkr, next_char);
        break;
    }
}
void yw_html_rawtext_state(struct YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '<':
        tkr->state = yw_html_rawtext_less_than_sign_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_emit_char_token(tkr, 0xfffd);
        break;
    case -1:
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_emit_char_token(tkr, next_char);
        break;
    }
}
void yw_html_plaintext_state(struct YW_HTMLTokenizer *tkr)
{
    (void)tkr;
    YW_TODO();
}
void yw_html_tag_open_state(struct YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '!':
        tkr->state = yw_html_markup_declaration_open_state;
        break;
    case '/':
        tkr->state = yw_html_end_tag_open_state;
        break;
    case '?':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_QUESTION_MARK_INSTEAD_OF_TAG_NAME_ERROR);
        yw_set_current_token(tkr, yw_make_comment_token(""));
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_bogus_comment_state;
        break;
    case -1:
        yw_emit_char_token(tkr, '<');
        yw_emit_eof_token(tkr);
        break;
    default:
        if (yw_is_ascii_alpha(next_char))
        {
            yw_set_current_token(tkr, yw_make_tag_token(""));
            tkr->tr.cursor = old_cursor;
            tkr->state = yw_html_tag_name_state;
        }
        else
        {
            yw_parse_error_encountered(tkr, YW_INVALID_FIRST_CHARACTER_OF_TAG_NAME_ERROR);
            yw_emit_char_token(tkr, '<');
            tkr->tr.cursor = old_cursor;
            tkr->state = yw_html_data_state;
        }
    }
}
void yw_html_end_tag_open_state(struct YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_END_TAG_NAME_ERROR);
        tkr->state = yw_html_data_state;
        break;
    case -1:
        yw_emit_char_token(tkr, '<');
        yw_emit_char_token(tkr, '/');
        yw_emit_eof_token(tkr);
        break;
    default:
        if (yw_is_ascii_alpha(next_char))
        {
            yw_set_current_token(tkr, yw_make_tag_token(""));
            yw_current_tag_token(tkr)->is_end = true;
            tkr->tr.cursor = old_cursor;
            tkr->state = yw_html_tag_name_state;
        }
        else
        {
            yw_parse_error_encountered(tkr, YW_INVALID_FIRST_CHARACTER_OF_TAG_NAME_ERROR);
            yw_set_current_token(tkr, yw_make_comment_token(""));
            tkr->tr.cursor = old_cursor;
            tkr->state = yw_html_bogus_comment_state;
        }
    }
}
void yw_html_tag_name_state(struct YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        tkr->state = yw_html_before_attribute_name_state;
        break;
    case '/':
        tkr->state = yw_html_self_closing_start_tag_state;
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_tag_token(tkr)->name, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default: {
        YW_Char32 chr = yw_to_ascii_lowercase(next_char);
        yw_append_char(&yw_current_tag_token(tkr)->name, chr);
    }
    }
}
void yw_html_rcdata_less_than_sign_state(struct YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '/':
        tkr->temp_buf = yw_duplicate_str("");
        tkr->state = yw_html_rcdata_end_tag_open_state;
        break;
    default:
        yw_emit_char_token(tkr, '<');
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_rcdata_state;
    }
}
void yw_html_rcdata_end_tag_open_state(struct YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    if (yw_is_ascii_alpha(next_char))
    {
        yw_set_current_token(tkr, yw_make_tag_token(""));
        yw_current_tag_token(tkr)->is_end = true;
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_rcdata_end_tag_name_state;
    }
    else
    {
        yw_emit_char_token(tkr, '<');
        yw_emit_char_token(tkr, '/');
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_rcdata_state;
    }
}
void yw_html_rcdata_end_tag_name_state(struct YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        if (yw_is_appropriate_end_tag_token(tkr, tkr->current_token))
        {
            tkr->state = yw_html_before_attribute_name_state;
            break;
        }
        goto anything_else;
    case '/':
        if (yw_is_appropriate_end_tag_token(tkr, tkr->current_token))
        {
            tkr->state = yw_html_self_closing_start_tag_state;
            break;
        }
        goto anything_else;
    case '>':
        if (yw_is_appropriate_end_tag_token(tkr, tkr->current_token))
        {
            tkr->state = yw_html_data_state;
            yw_emit_token(tkr, &tkr->current_token);
            break;
        }
        goto anything_else;
    default:
        if (yw_is_ascii_alpha(next_char))
        {
            YW_Char32 chr = yw_to_ascii_lowercase(next_char);
            yw_append_char(&tkr->temp_buf, chr);
            break;
        }
        goto anything_else;
    }
    return;
anything_else:
    yw_emit_char_token(tkr, '<');
    yw_emit_char_token(tkr, '/');

    char const *next_temp_buf_chr = tkr->temp_buf;
    while (1)
    {
        YW_Char32 c = yw_utf8_next_char(&next_temp_buf_chr);
        if (c == -1)
        {
            break;
        }
        yw_emit_char_token(tkr, c);
    }
    tkr->tr.cursor = old_cursor;
    tkr->state = yw_html_rcdata_state;
}
