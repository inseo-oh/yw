/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_html_tokens.h"
#include "yw_common.h"
#include "yw_dom.h"
#include <assert.h>
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
void yw_html_token_clone(YW_HTMLToken *dest, YW_HTMLToken const *tk)
{
    *dest = *tk;
    switch (tk->type)
    {
    case YW_HTML_EOF_TOKEN:
    case YW_HTML_CHAR_TOKEN:
        return;
    case YW_HTML_COMMENT_TOKEN:
        dest->comment_tk.data = yw_duplicate_str(tk->comment_tk.data);
        return;
    case YW_HTML_DOCTYPE_TOKEN:
        dest->doctype_tk.name = yw_duplicate_str(tk->doctype_tk.name);
        dest->doctype_tk.public_id = yw_duplicate_str(tk->doctype_tk.public_id);
        dest->doctype_tk.system_id = yw_duplicate_str(tk->doctype_tk.system_id);
        return;
    case YW_HTML_TAG_TOKEN:
        dest->tag_tk.name = yw_duplicate_str(tk->tag_tk.name);
        dest->tag_tk.attrs = YW_ALLOC(YW_DOMAttrData, tk->tag_tk.attrs_len);
        for (int i = 0; i < tk->tag_tk.attrs_len; i++)
        {
            yw_dom_attr_data_clone(&dest->tag_tk.attrs[i], &tk->tag_tk.attrs[i]);
        }
        return;
    }
    YW_UNREACHABLE();
}

static bool yw_is_start_tag(YW_HTMLToken const *tk)
{
    return tk->type == YW_HTML_TAG_TOKEN && !tk->tag_tk.is_end;
}
static bool yw_is_end_tag(YW_HTMLToken const *tk)
{
    return tk->type == YW_HTML_TAG_TOKEN && tk->tag_tk.is_end;
}

static YW_HTMLToken *yw_make_eof_token(void)
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
static YW_HTMLToken *yw_make_doctype_token(void)
{
    YW_HTMLToken *res = YW_ALLOC(YW_HTMLToken, 1);
    memset(res, 0, sizeof(*res));
    res->type = YW_HTML_DOCTYPE_TOKEN;
    return res;
}

static const bool YW_TRACE_TOKENIZER_STATE = true;

void yw_html_tokenize(
    uint8_t const *chars, int chars_len,
    void (*emit_callback)(YW_HTMLToken *token, void *callback_data), void *emit_callback_data)
{
    YW_HTMLTokenizer tkr;
    memset(&tkr, 0, sizeof(tkr));

    yw_text_reader_init(&tkr.tr, chars, chars_len);
    tkr.state = yw_html_data_state;
    tkr.emit_callback = emit_callback;
    tkr.emit_callback_data = emit_callback_data;

    while (!tkr.eof_emitted)
    {
        if (YW_TRACE_TOKENIZER_STATE)
        {
#define YW_X(_x)                                                  \
    do                                                            \
    {                                                             \
        if (tkr.state == _x)                                      \
        {                                                         \
            fprintf(stderr, "%s: NEXT STATE: %s\n", __func__, #_x); \
        }                                                         \
    } while (0);
            YW_HTML_ENUMERATE_TOKENIZER_STATE(YW_X);
#undef YW_X
        }
        tkr.state(&tkr);
    }

    free(tkr.temp_buf);
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
    if (tk->type == YW_HTML_EOF_TOKEN)
    {
        tkr->eof_emitted = true;
    }
    tkr->emit_callback(tk, tkr->emit_callback_data);
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
static void yw_emit_doctype_token(YW_HTMLTokenizer *tkr)
{
    YW_HTMLToken *token = yw_make_doctype_token();
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
        char const *next_char = tkr->temp_buf;
        while (1)
        {
            YW_Char32 c = yw_utf8_next_char(&next_char);
            if (c == 0)
            {
                break;
            }
            yw_emit_char_token(tkr, c);
        }
    }
}

static void yw_add_attr_to_current_tag(YW_HTMLTokenizer *tkr, char const *name)
{
    YW_HTMLTagToken *tag = yw_current_tag_token(tkr);
    YW_DOMAttrData attr;
    memset(&attr, 0, sizeof(attr));
    attr.local_name = yw_duplicate_str(name);
    YW_PUSH(YW_DOMAttrData, &tkr->curr_attrs_cap, &tag->attrs_len, &tag->attrs, attr);
    assert(tkr->bad_attrs == NULL);
    assert(tkr->bad_attrs_cap == 0);
    assert(tkr->bad_attrs_len == 0);
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
void yw_html_rcdata_state(YW_HTMLTokenizer *tkr)
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
void yw_html_rawtext_state(YW_HTMLTokenizer *tkr)
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
void yw_html_plaintext_state(YW_HTMLTokenizer *tkr)
{
    (void)tkr;
    YW_TODO();
}
void yw_html_tag_open_state(YW_HTMLTokenizer *tkr)
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
void yw_html_end_tag_open_state(YW_HTMLTokenizer *tkr)
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
void yw_html_tag_name_state(YW_HTMLTokenizer *tkr)
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
void yw_html_rcdata_less_than_sign_state(YW_HTMLTokenizer *tkr)
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
void yw_html_rcdata_end_tag_open_state(YW_HTMLTokenizer *tkr)
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
void yw_html_rcdata_end_tag_name_state(YW_HTMLTokenizer *tkr)
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
void yw_html_rawtext_less_than_sign_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '/':
        tkr->temp_buf = yw_duplicate_str("");
        tkr->state = yw_html_rawtext_end_tag_open_state;
        break;
    default:
        yw_emit_char_token(tkr, '<');
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_rawtext_state;
    }
}
void yw_html_rawtext_end_tag_open_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    if (yw_is_ascii_alpha(next_char))
    {
        yw_set_current_token(tkr, yw_make_tag_token(""));
        yw_current_tag_token(tkr)->is_end = true;
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_rawtext_end_tag_name_state;
    }
    else
    {
        yw_emit_char_token(tkr, '<');
        yw_emit_char_token(tkr, '/');
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_rawtext_state;
    }
}
void yw_html_rawtext_end_tag_name_state(YW_HTMLTokenizer *tkr)
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
    tkr->state = yw_html_rawtext_state;
}
void yw_html_before_attribute_name_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '/':
    case '>':
    case -1:
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_after_attribute_name_state;
        break;
    case '=': {
        char *name_str = yw_char_to_str(next_char);
        yw_add_attr_to_current_tag(tkr, name_str);
        free(name_str);
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_attribute_name_state;
        break;
    }
    default:
        yw_add_attr_to_current_tag(tkr, "");
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_attribute_name_state;
        break;
    }
}
void yw_html_attribute_name_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
    case '/':
    case '>':
    case -1:
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_after_attribute_name_state;
        break;
    case '=':
        tkr->state = yw_html_before_attribute_value_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_attr(tkr)->local_name, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case '"':
    case '\\':
    case '<':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_CHARACTER_IN_ATTRIBUTE_NAME_ERROR);
        goto anything_else;
    default:
        goto anything_else;
    }
    goto check_dupliate_attr_name;
anything_else:
    yw_append_char(&yw_current_attr(tkr)->local_name, next_char);
    return;
check_dupliate_attr_name: {
    YW_HTMLTagToken *tag = yw_current_tag_token(tkr);
    YW_DOMAttrData const *current_attr = yw_current_attr(tkr);
    for (int i = 0; i < tag->attrs_len; i++)
    {
        if (strcmp(current_attr->local_name, tag->attrs[i].local_name) == 0)
        {
            YW_PUSH(int, &tkr->bad_attrs_cap, &tkr->bad_attrs_len, &tkr->bad_attrs, i);
        }
    }
}
}
void yw_html_after_attribute_name_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '/':
        tkr->state = yw_html_self_closing_start_tag_state;
        break;
    case '=':
        tkr->state = yw_html_before_attribute_value_state;
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_add_attr_to_current_tag(tkr, "");
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_before_attribute_name_state;
        break;
    }
}
void yw_html_before_attribute_value_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '"':
        tkr->state = yw_html_attribute_value_double_quoted_state;
        break;
    case '\'':
        tkr->state = yw_html_attribute_value_single_quoted_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_ATTRIBUTE_VALUE_ERROR);
        tkr->state = yw_html_data_state;
        break;
    default:
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_attribute_value_unquoted_state;
    }
}
void yw_html_attribute_value_double_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '"':
        tkr->state = yw_html_after_attribute_value_quoted_state;
        break;
    case '&':
        tkr->return_state = yw_html_attribute_value_double_quoted_state;
        tkr->state = yw_html_character_reference_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_attr(tkr)->value, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_char(&yw_current_attr(tkr)->value, next_char);
        break;
    }
}
void yw_html_attribute_value_single_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\'':
        tkr->state = yw_html_after_attribute_value_quoted_state;
        break;
    case '&':
        tkr->return_state = yw_html_attribute_value_single_quoted_state;
        tkr->state = yw_html_character_reference_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_attr(tkr)->value, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_char(&yw_current_attr(tkr)->value, next_char);
        break;
    }
}
void yw_html_attribute_value_unquoted_state(YW_HTMLTokenizer *tkr)
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
    case '&':
        tkr->return_state = yw_html_attribute_value_unquoted_state;
        tkr->state = yw_html_character_reference_state;
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_attr(tkr)->value, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_char(&yw_current_attr(tkr)->value, next_char);
        break;
    }
}
void yw_html_after_attribute_value_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
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
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_BETWEEN_ATTRIBUTES_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_before_attribute_name_state;
        break;
    }
}
void yw_html_self_closing_start_tag_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '>':
        yw_current_tag_token(tkr)->is_self_closing = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_TAG_ERROR);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_SOLIDUS_IN_TAG_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_before_attribute_name_state;
        break;
    }
}
void yw_html_bogus_comment_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_comment_token(tkr)->data, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    default:
        yw_append_char(&yw_current_comment_token(tkr)->data, next_char);
        break;
    }
}
void yw_html_markup_declaration_open_state(YW_HTMLTokenizer *tkr)
{
    if (yw_consume_str(&tkr->tr, "--", YW_NO_MATCH_FLAGS))
    {
        yw_set_current_token(tkr, yw_make_comment_token(""));
        tkr->state = yw_html_comment_start_state;
    }
    else if (yw_consume_str(&tkr->tr, "DOCTYPE", YW_ASCII_CASE_INSENSITIVE))
    {
        tkr->state = yw_html_doctype_state;
    }
    else if (yw_consume_str(&tkr->tr, "[CDATA[", YW_NO_MATCH_FLAGS))
    {
        YW_TODO();
    }
    else
    {
        yw_parse_error_encountered(tkr, YW_INCORRECTLY_OPENED_COMMENT_ERROR);
        yw_set_current_token(tkr, yw_make_comment_token(""));
        tkr->state = yw_html_bogus_comment_state;
    }
}
void yw_html_comment_start_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '-':
        tkr->state = yw_html_comment_start_dash_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_ABRUPT_CLOSING_OF_EMPTY_COMMENT_ERROR);
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    default:
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_bogus_comment_state;
        break;
    }
}
void yw_html_comment_start_dash_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '-':
        tkr->state = yw_html_comment_end_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_ABRUPT_CLOSING_OF_EMPTY_COMMENT_ERROR);
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_COMMENT_ERROR);
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_str(&yw_current_comment_token(tkr)->data, "-");
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_comment_state;
        break;
    }
}
void yw_html_comment_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '<':
        yw_append_str(&yw_current_comment_token(tkr)->data, "<");
        tkr->state = yw_html_comment_less_than_sign_state;
        break;
    case '-':
        tkr->state = yw_html_comment_end_dash_state;
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_comment_token(tkr)->data, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_COMMENT_ERROR);
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_char(&yw_current_comment_token(tkr)->data, next_char);
        break;
    }
}
void yw_html_comment_less_than_sign_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '!':
        yw_append_char(&yw_current_comment_token(tkr)->data, next_char);
        YW_TODO();
        /* tkr->state = yw_html_comment_less_than_sign_bang_state; */
        break;
    case '<':
        yw_append_char(&yw_current_comment_token(tkr)->data, next_char);
        break;
    default:
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_comment_state;
        break;
    }
}
void yw_html_comment_end_dash_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '-':
        tkr->state = yw_html_comment_end_state;
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_COMMENT_ERROR);
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_str(&yw_current_comment_token(tkr)->data, "-");
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_comment_state;
        break;
    }
}
void yw_html_comment_end_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case '!':
        YW_TODO();
        /* tkr->state = yw_html_comment_end_bang_state; */
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_COMMENT_ERROR);
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_str(&yw_current_comment_token(tkr)->data, "--");
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_bogus_comment_state;
        break;
    }
}
void yw_html_doctype_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        tkr->state = yw_html_before_doctype_name_state;
        break;
    case '>':
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_before_doctype_name_state;
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_emit_doctype_token(tkr);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_BEFORE_DOCTYPE_NAME_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_before_doctype_name_state;
        break;
    }
}
void yw_html_before_doctype_name_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_set_current_token(tkr, yw_make_doctype_token());
        yw_current_doctype_token(tkr)->name = yw_duplicate_str(YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_DOCTYPE_NAME_ERROR);
        yw_set_current_token(tkr, yw_make_doctype_token());
        yw_current_doctype_token(tkr)->force_quirks = true;
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_emit_doctype_token(tkr);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_set_current_token(tkr, yw_make_doctype_token());
        yw_current_doctype_token(tkr)->name = yw_char_to_str(yw_to_ascii_lowercase(next_char));
        tkr->state = yw_html_doctype_name_state;
    }
}
void yw_html_doctype_name_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        tkr->state = yw_html_after_doctype_name_state;
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case '\0':
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_NULL_CHARACTER_ERROR);
        yw_append_str(&yw_current_doctype_token(tkr)->name, YW_UNICODE_REPLACEMENT_CHAR);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_char(&yw_current_doctype_token(tkr)->name, yw_to_ascii_lowercase(next_char));
        break;
    }
}
void yw_html_after_doctype_name_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        tkr->tr.cursor = old_cursor;
        if (yw_consume_str(&tkr->tr, "PUBLIC", YW_ASCII_CASE_INSENSITIVE))
        {
            tkr->state = yw_html_after_doctype_public_keyword_state;
        }
        else if (yw_consume_str(&tkr->tr, "SYSTEM", YW_ASCII_CASE_INSENSITIVE))
        {
            tkr->state = yw_html_after_doctype_system_keyword_state;
        }
        else
        {
            yw_parse_error_encountered(tkr, YW_INVALID_CHARACTER_SEQUENCE_AFTER_DOCTYPE_NAME_ERROR);
            yw_current_doctype_token(tkr)->force_quirks = true;
            YW_TODO();
            /* tkr->state = bogus_doctype_state; */
        }
    }
}
void yw_html_after_doctype_public_keyword_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        tkr->state = yw_html_before_doctype_public_identifier_state;
        break;
    case '"':
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_AFTER_DOCTYPE_PUBLIC_KEYWORD_ERROR);
        yw_current_doctype_token(tkr)->public_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_public_identifier_double_quoted_state;
        break;
    case '\'':
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_AFTER_DOCTYPE_PUBLIC_KEYWORD_ERROR);
        yw_current_doctype_token(tkr)->public_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_public_identifier_single_quoted_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_DOCTYPE_PUBLIC_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_QUOTE_BEFORE_DOCTYPE_PUBLIC_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_before_doctype_public_identifier_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '"':
        yw_current_doctype_token(tkr)->public_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_public_identifier_double_quoted_state;
        break;
    case '\'':
        yw_current_doctype_token(tkr)->public_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_public_identifier_single_quoted_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_DOCTYPE_PUBLIC_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_QUOTE_BEFORE_DOCTYPE_PUBLIC_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_doctype_public_identifier_double_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '"':
        tkr->state = yw_html_after_doctype_public_identifier_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_ABRUPT_DOCTYPE_PUBLIC_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_char(&yw_current_doctype_token(tkr)->public_id, next_char);
        break;
    }
}
void yw_html_doctype_public_identifier_single_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\'':
        tkr->state = yw_html_after_doctype_public_identifier_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_ABRUPT_DOCTYPE_PUBLIC_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_char(&yw_current_doctype_token(tkr)->public_id, next_char);
        break;
    }
}
void yw_html_after_doctype_public_identifier_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        tkr->state = yw_html_between_doctype_public_and_system_identifiers_state;
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case '"':
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_BETWEEN_DOCTYPE_PUBLIC_AND_SYSTEM_IDENTIFIERS_ERROR);
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_double_quoted_state;
        break;
    case '\'':
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_BETWEEN_DOCTYPE_PUBLIC_AND_SYSTEM_IDENTIFIERS_ERROR);
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_single_quoted_state;
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_QUOTE_BEFORE_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_between_doctype_public_and_system_identifiers_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case '"':
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_double_quoted_state;
        break;
    case '\'':
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_single_quoted_state;
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_QUOTE_BEFORE_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_after_doctype_system_keyword_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        tkr->state = yw_html_before_doctype_system_identifier_state;
        break;
    case '"':
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_AFTER_DOCTYPE_SYSTEM_KEYWORD_ERROR);
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_double_quoted_state;
        break;
    case '\'':
        yw_parse_error_encountered(tkr, YW_MISSING_WHITESPACE_AFTER_DOCTYPE_SYSTEM_KEYWORD_ERROR);
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_single_quoted_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_QUOTE_BEFORE_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_before_doctype_system_identifier_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '"':
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_double_quoted_state;
        break;
    case '\'':
        yw_current_doctype_token(tkr)->system_id = yw_duplicate_str("");
        tkr->state = yw_html_doctype_system_identifier_single_quoted_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_MISSING_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_MISSING_QUOTE_BEFORE_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_doctype_system_identifier_double_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '"':
        tkr->state = yw_html_after_doctype_system_identifier_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_ABRUPT_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_char(&yw_current_doctype_token(tkr)->system_id, next_char);
        break;
    }
}
void yw_html_doctype_system_identifier_single_quoted_state(YW_HTMLTokenizer *tkr)
{
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\'':
        tkr->state = yw_html_after_doctype_system_identifier_state;
        break;
    case '>':
        yw_parse_error_encountered(tkr, YW_ABRUPT_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_append_char(&yw_current_doctype_token(tkr)->system_id, next_char);
        break;
    }
}
void yw_html_after_doctype_system_identifier_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '\t':
    case '\n':
    case '\x0c':
    case ' ':
        break;
    case '>':
        tkr->state = yw_html_data_state;
        yw_emit_token(tkr, &tkr->current_token);
        break;
    case -1:
        yw_parse_error_encountered(tkr, YW_EOF_IN_DOCTYPE_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        yw_emit_token(tkr, &tkr->current_token);
        yw_emit_eof_token(tkr);
        break;
    default:
        yw_parse_error_encountered(tkr, YW_UNEXPECTED_CHARACTER_AFTER_DOCTYPE_SYSTEM_IDENTIFIER_ERROR);
        yw_current_doctype_token(tkr)->force_quirks = true;
        tkr->tr.cursor = old_cursor;
        YW_TODO();
        /* tkr->state = bogus_doctype_state; */
    }
}
void yw_html_character_reference_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case '#':
        yw_append_char(&tkr->temp_buf, next_char);
        tkr->state = yw_html_numeric_character_reference_state;
        break;
    default:
        if (yw_is_ascii_alphanumeric(next_char))
        {
            tkr->tr.cursor = old_cursor;
            tkr->state = yw_html_named_character_reference_state;
        }
        else
        {
            yw_flush_codepoints_consumed_as_char_reference(tkr);
            tkr->tr.cursor = old_cursor;
            tkr->state = tkr->return_state;
        }
    }
}

#include "yw_html_entities_autogen.i"

void yw_html_named_character_reference_state(YW_HTMLTokenizer *tkr)
{
    (void)tkr;
    int entities_count = YW_SIZEOF_ARRAY(yw_html_entities);
    char const *found_name = NULL;
    char const *found_str = NULL;
    YW_TextCursor cursor_after_found;
    for (int i = 0; i < entities_count; i++)
    {
        if (yw_html_entities[i].name[0] != '&')
        {
            fprintf(stderr, "%s: internal warning: key %s in yw_html_entities doesn't start with &\n",
                    __func__, yw_html_entities[i].name);
            continue;
        }
        YW_TextCursor cursor_before_str = tkr->tr.cursor;
        if (yw_consume_str(&tkr->tr, &yw_html_entities[i].name[1], YW_NO_MATCH_FLAGS))
        {
            if (found_name == NULL || strlen(found_name) < strlen(yw_html_entities[i].name))
            {
                found_name = yw_html_entities[i].name;
                found_str = yw_html_entities[i].str;
                cursor_after_found = tkr->tr.cursor;
            }
            tkr->tr.cursor = cursor_before_str;
        }
    }
    if (found_name != NULL)
    {
        tkr->tr.cursor = cursor_after_found;
        if (yw_is_consumed_as_part_of_attr(tkr) &&
            found_name[strlen(found_name) - 1] != ';' &&
            (yw_peek_char(&tkr->tr) == '=' || yw_is_ascii_alphanumeric(yw_peek_char(&tkr->tr))))
        {
            yw_flush_codepoints_consumed_as_char_reference(tkr);
            tkr->state = tkr->return_state;
        }
        else
        {
            if (found_name[strlen(found_name) - 1] != ';')
            {
                yw_parse_error_encountered(tkr, YW_MISSING_SEMICOLON_AFTER_CHARACTER_REFERENCE_ERROR);
            }
            tkr->temp_buf = yw_duplicate_str(found_str);
            yw_flush_codepoints_consumed_as_char_reference(tkr);
            tkr->state = tkr->return_state;
        }
    }
    else
    {
        yw_flush_codepoints_consumed_as_char_reference(tkr);
        tkr->state = tkr->return_state;
    }
}
void yw_html_numeric_character_reference_state(YW_HTMLTokenizer *tkr)
{
    tkr->character_reference_code = 0;

    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    switch (next_char)
    {
    case 'X':
    case 'x':
        yw_append_char(&tkr->temp_buf, next_char);
        tkr->state = yw_html_hexadecimal_character_reference_start_state;
        break;
    default:
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_decimal_character_reference_start_state;
        break;
    }
}
void yw_html_hexadecimal_character_reference_start_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    if (yw_is_ascii_hex_digit(next_char))
    {
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_hexadecimal_character_reference_state;
    }
    else
    {
        yw_parse_error_encountered(tkr, YW_ABSENCE_OF_DIGITS_IN_NUMERIC_CHARACTER_REFERENCE_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = tkr->return_state;
    }
}
void yw_html_decimal_character_reference_start_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    if (yw_is_ascii_digit(next_char))
    {
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_decimal_character_reference_state;
    }
    else
    {
        yw_parse_error_encountered(tkr, YW_ABSENCE_OF_DIGITS_IN_NUMERIC_CHARACTER_REFERENCE_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = tkr->return_state;
    }
}
void yw_html_hexadecimal_character_reference_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    if (yw_is_ascii_digit(next_char))
    {
        tkr->character_reference_code = (tkr->character_reference_code * 16) + (next_char - '0');
    }
    else if (yw_is_ascii_uppercase_hex_digit(next_char))
    {
        tkr->character_reference_code = (tkr->character_reference_code * 16) + (next_char - 'A' + 10);
    }
    else if (yw_is_ascii_lowercase_hex_digit(next_char))
    {
        tkr->character_reference_code = (tkr->character_reference_code * 16) + (next_char - 'a' + 10);
    }
    else if (next_char == ';')
    {
        tkr->state = yw_html_numeric_character_reference_end_state;
    }
    else
    {
        yw_parse_error_encountered(tkr, YW_MISSING_SEMICOLON_AFTER_CHARACTER_REFERENCE_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_numeric_character_reference_end_state;
    }
}
void yw_html_decimal_character_reference_state(YW_HTMLTokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    YW_Char32 next_char = yw_consume_any_char(&tkr->tr);
    if (yw_is_ascii_digit(next_char))
    {
        tkr->character_reference_code = (tkr->character_reference_code * 10) + (next_char - '0');
    }
    else if (next_char == ';')
    {
        tkr->state = yw_html_numeric_character_reference_end_state;
    }
    else
    {
        yw_parse_error_encountered(tkr, YW_MISSING_SEMICOLON_AFTER_CHARACTER_REFERENCE_ERROR);
        tkr->tr.cursor = old_cursor;
        tkr->state = yw_html_numeric_character_reference_end_state;
    }
}
void yw_html_numeric_character_reference_end_state(YW_HTMLTokenizer *tkr)
{
    if (tkr->character_reference_code == 0x0000)
    {
        yw_parse_error_encountered(tkr, YW_NULL_CHARACTER_REFERENCE_ERROR);
        tkr->character_reference_code = 0xfffd;
    }
    else if (0x10ffff < tkr->character_reference_code)
    {
        yw_parse_error_encountered(tkr, YW_CHARACTER_REFERENCE_OUTSIDE_UNICODE_RANGE_ERROR);
        tkr->character_reference_code = 0xfffd;
    }
    else if (yw_is_surrogate_char(tkr->character_reference_code))
    {
        yw_parse_error_encountered(tkr, YW_SURROGATE_CHARACTER_REFERENCE_ERROR);
        tkr->character_reference_code = 0xfffd;
    }
    else if (yw_is_noncharacter(tkr->character_reference_code))
    {
        yw_parse_error_encountered(tkr, YW_NONCHARACTER_REFERENCE_ERROR);
    }
    else if (
        tkr->character_reference_code == 0x0d ||
        (yw_is_control_char(tkr->character_reference_code) && !yw_is_ascii_whitespace(tkr->character_reference_code)))
    {
        yw_parse_error_encountered(tkr, YW_CONTROL_CHARACTER_REFERENCE_ERROR);
        switch (tkr->character_reference_code)
        {
        case 0x80:
            tkr->character_reference_code = 0x20ac;
            break;
        case 0x82:
            tkr->character_reference_code = 0x201a;
            break;
        case 0x83:
            tkr->character_reference_code = 0x0192;
            break;
        case 0x84:
            tkr->character_reference_code = 0x201e;
            break;
        case 0x85:
            tkr->character_reference_code = 0x2026;
            break;
        case 0x86:
            tkr->character_reference_code = 0x2020;
            break;
        case 0x87:
            tkr->character_reference_code = 0x2021;
            break;
        case 0x88:
            tkr->character_reference_code = 0x02c6;
            break;
        case 0x89:
            tkr->character_reference_code = 0x2030;
            break;
        case 0x8a:
            tkr->character_reference_code = 0x0160;
            break;
        case 0x8b:
            tkr->character_reference_code = 0x2039;
            break;
        case 0x8c:
            tkr->character_reference_code = 0x0152;
            break;
        case 0x8e:
            tkr->character_reference_code = 0x017d;
            break;
        case 0x91:
            tkr->character_reference_code = 0x2018;
            break;
        case 0x92:
            tkr->character_reference_code = 0x2019;
            break;
        case 0x93:
            tkr->character_reference_code = 0x201c;
            break;
        case 0x94:
            tkr->character_reference_code = 0x201d;
            break;
        case 0x95:
            tkr->character_reference_code = 0x2022;
            break;
        case 0x96:
            tkr->character_reference_code = 0x2013;
            break;
        case 0x97:
            tkr->character_reference_code = 0x2014;
            break;
        case 0x98:
            tkr->character_reference_code = 0x02dc;
            break;
        case 0x99:
            tkr->character_reference_code = 0x2122;
            break;
        case 0x9a:
            tkr->character_reference_code = 0x0161;
            break;
        case 0x9b:
            tkr->character_reference_code = 0x203a;
            break;
        case 0x9c:
            tkr->character_reference_code = 0x0153;
            break;
        case 0x9e:
            tkr->character_reference_code = 0x017e;
            break;
        case 0x9f:
            tkr->character_reference_code = 0x0178;
            break;
        }
    }
    free(tkr->temp_buf);
    tkr->temp_buf = yw_char_to_str(tkr->character_reference_code);
    yw_flush_codepoints_consumed_as_char_reference(tkr);
    tkr->state = tkr->return_state;
}
