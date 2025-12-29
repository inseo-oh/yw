/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_css_tokens.h"
#include "yw_common.h"
#include "yw_css.h"
#include "yw_encoding.h"
#include <limits.h>
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static bool yw_is_ident_start_codepoint(YW_Char32 c)
{
    /* https://www.w3.org/TR/-syntax-3/#ident-start-code-point */
    return yw_is_ascii_alpha(c) || 0x80 <= c || c == '_';
}

static bool yw_is_ident_codepoint(YW_Char32 c)
{
    /* https://www.w3.org/TR/2021/CRD--syntax-3-20211224/#ident-code-point */
    return yw_is_ident_start_codepoint(c) || yw_is_ascii_digit(c) || c == '-';
}

static bool yw_is_valid_ident_start_sequence(char const *s)
{
    YW_Char32 *chars;
    int chars_len;
    bool res;
    yw_utf8_to_char32(&chars, &chars_len, s);
    if (chars_len == 0)
    {
        res = false;
        goto end;
    }
    if (yw_is_ident_start_codepoint(chars[0]))
    {
        res = true;
        goto end;
    }
    switch (chars[0])
    {
    case '-':
        res = (1 < chars_len && yw_is_ident_codepoint(chars[1])) || (2 < chars_len && chars[1] == '\\' && chars[2] != '\n');
        goto end;
    case '\\':
        res = 1 < chars_len && chars[1] != '\n';
        goto end;
    }
    res = false;
end:
    free(chars);
    return res;
}

void yw_css_token_deinit(YW_CSSToken *tk)
{
    switch (tk->common.type)
    {
    case YW_CSS_TOKEN_WHITESPACE:
    case YW_CSS_TOKEN_LEFT_PAREN:
    case YW_CSS_TOKEN_RIGHT_PAREN:
    case YW_CSS_TOKEN_COMMA:
    case YW_CSS_TOKEN_COLON:
    case YW_CSS_TOKEN_SEMICOLON:
    case YW_CSS_TOKEN_LEFT_SQUARE_BRACKET:
    case YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET:
    case YW_CSS_TOKEN_LEFT_CURLY_BRACKET:
    case YW_CSS_TOKEN_RIGHT_CURLY_BRACKET:
    case YW_CSS_TOKEN_CDO:
    case YW_CSS_TOKEN_CDC:
    case YW_CSS_TOKEN_BAD_STRING:
    case YW_CSS_TOKEN_BAD_URL:
    case YW_CSS_TOKEN_NUMBER:
    case YW_CSS_TOKEN_PERCENTAGE:
    case YW_CSS_TOKEN_DELIM:
        return;
    case YW_CSS_TOKEN_DIMENSION:
        free(tk->dimension_tk.unit);
        break;
    case YW_CSS_TOKEN_STRING:
        free(tk->string_tk.value);
        break;
    case YW_CSS_TOKEN_URL:
        free(tk->url_tk.value);
        break;
    case YW_CSS_TOKEN_AT_KEYWORD:
        free(tk->at_keyword_tk.value);
        break;
    case YW_CSS_TOKEN_FUNC_KEYWORD:
        free(tk->func_keyword_tk.value);
        break;
    case YW_CSS_TOKEN_IDENT:
        free(tk->ident_tk.value);
        break;
    case YW_CSS_TOKEN_HASH:
        free(tk->hash_tk.value);
        break;
    case YW_CSS_TOKEN_AST_SIMPLE_BLOCK:
        for (int i = 0; i < tk->ast_simple_block_tk.tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_simple_block_tk.tokens[i]);
        }
        free(tk->ast_simple_block_tk.tokens);
        break;
    case YW_CSS_TOKEN_AST_FUNC:
        for (int i = 0; i < tk->ast_func_tk.tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_func_tk.tokens[i]);
        }
        free(tk->ast_func_tk.tokens);
        free(tk->ast_func_tk.name);
        break;
    case YW_CSS_TOKEN_AST_QUALIFIED_RULE:
        for (int i = 0; i < tk->ast_qualified_rule_tk.prelude_tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_qualified_rule_tk.prelude_tokens[i]);
        }
        for (int i = 0; i < tk->ast_qualified_rule_tk.body_tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_qualified_rule_tk.body_tokens[i]);
        }
        free(tk->ast_qualified_rule_tk.prelude_tokens);
        free(tk->ast_qualified_rule_tk.body_tokens);
        break;
    case YW_CSS_TOKEN_AST_AT_RULE:
        for (int i = 0; i < tk->ast_at_rule_tk.prelude_tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_at_rule_tk.prelude_tokens[i]);
        }
        for (int i = 0; i < tk->ast_at_rule_tk.body_tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_at_rule_tk.body_tokens[i]);
        }
        free(tk->ast_at_rule_tk.prelude_tokens);
        free(tk->ast_at_rule_tk.body_tokens);
        free(tk->ast_at_rule_tk.name);
        break;
    case YW_CSS_TOKEN_AST_DECLARATION:
        for (int i = 0; i < tk->ast_declaration_tk.value_tokens_len; i++)
        {
            yw_css_token_deinit(&tk->ast_declaration_tk.value_tokens[i]);
        }
        free(tk->ast_declaration_tk.value_tokens);
        free(tk->ast_declaration_tk.name);
        break;
    }
}

static void yw_token_clone(YW_CSSToken *dest, YW_CSSToken const *tk)
{
    memcpy(dest, tk, sizeof(*tk));
    switch (tk->common.type)
    {
    case YW_CSS_TOKEN_WHITESPACE:
    case YW_CSS_TOKEN_LEFT_PAREN:
    case YW_CSS_TOKEN_RIGHT_PAREN:
    case YW_CSS_TOKEN_COMMA:
    case YW_CSS_TOKEN_COLON:
    case YW_CSS_TOKEN_SEMICOLON:
    case YW_CSS_TOKEN_LEFT_SQUARE_BRACKET:
    case YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET:
    case YW_CSS_TOKEN_LEFT_CURLY_BRACKET:
    case YW_CSS_TOKEN_RIGHT_CURLY_BRACKET:
    case YW_CSS_TOKEN_CDO:
    case YW_CSS_TOKEN_CDC:
    case YW_CSS_TOKEN_BAD_STRING:
    case YW_CSS_TOKEN_BAD_URL:
    case YW_CSS_TOKEN_NUMBER:
    case YW_CSS_TOKEN_PERCENTAGE:
    case YW_CSS_TOKEN_DELIM:
        return;
    case YW_CSS_TOKEN_DIMENSION:
        dest->dimension_tk.unit = yw_duplicate_str(tk->dimension_tk.unit);
        break;
    case YW_CSS_TOKEN_STRING:
        dest->string_tk.value = yw_duplicate_str(tk->string_tk.value);
        break;
    case YW_CSS_TOKEN_URL:
        dest->url_tk.value = yw_duplicate_str(tk->url_tk.value);
        break;
    case YW_CSS_TOKEN_AT_KEYWORD:
        dest->at_keyword_tk.value = yw_duplicate_str(tk->at_keyword_tk.value);
        break;
    case YW_CSS_TOKEN_FUNC_KEYWORD:
        dest->func_keyword_tk.value = yw_duplicate_str(tk->func_keyword_tk.value);
        break;
    case YW_CSS_TOKEN_IDENT:
        dest->ident_tk.value = yw_duplicate_str(tk->ident_tk.value);
        break;
    case YW_CSS_TOKEN_HASH:
        dest->hash_tk.value = yw_duplicate_str(tk->hash_tk.value);
        break;
    case YW_CSS_TOKEN_AST_SIMPLE_BLOCK:
        dest->ast_simple_block_tk.tokens = YW_ALLOC(YW_CSSToken, tk->ast_simple_block_tk.tokens_len);
        for (int i = 0; i < tk->ast_simple_block_tk.tokens_len; i++)
        {
            yw_token_clone(&dest->ast_simple_block_tk.tokens[i], &tk->ast_simple_block_tk.tokens[i]);
        }
        break;
    case YW_CSS_TOKEN_AST_FUNC:
        dest->ast_func_tk.tokens = YW_ALLOC(YW_CSSToken, tk->ast_func_tk.tokens_len);
        for (int i = 0; i < tk->ast_func_tk.tokens_len; i++)
        {
            yw_token_clone(&dest->ast_func_tk.tokens[i], &tk->ast_func_tk.tokens[i]);
        }
        dest->ast_func_tk.name = yw_duplicate_str(tk->ast_func_tk.name);
        break;
    case YW_CSS_TOKEN_AST_QUALIFIED_RULE:
        dest->ast_qualified_rule_tk.prelude_tokens = YW_ALLOC(YW_CSSToken, tk->ast_qualified_rule_tk.prelude_tokens_len);
        dest->ast_qualified_rule_tk.body_tokens = YW_ALLOC(YW_CSSToken, tk->ast_qualified_rule_tk.body_tokens_len);
        for (int i = 0; i < tk->ast_qualified_rule_tk.prelude_tokens_len; i++)
        {
            yw_token_clone(&dest->ast_qualified_rule_tk.prelude_tokens[i], &tk->ast_qualified_rule_tk.prelude_tokens[i]);
        }
        for (int i = 0; i < tk->ast_qualified_rule_tk.body_tokens_len; i++)
        {
            yw_token_clone(&dest->ast_qualified_rule_tk.body_tokens[i], &tk->ast_qualified_rule_tk.body_tokens[i]);
        }
        break;
    case YW_CSS_TOKEN_AST_AT_RULE:
        dest->ast_at_rule_tk.prelude_tokens = YW_ALLOC(YW_CSSToken, tk->ast_at_rule_tk.prelude_tokens_len);
        dest->ast_at_rule_tk.body_tokens = YW_ALLOC(YW_CSSToken, tk->ast_at_rule_tk.body_tokens_len);
        for (int i = 0; i < tk->ast_at_rule_tk.prelude_tokens_len; i++)
        {
            yw_token_clone(&dest->ast_at_rule_tk.prelude_tokens[i], &tk->ast_at_rule_tk.prelude_tokens[i]);
        }
        for (int i = 0; i < tk->ast_at_rule_tk.body_tokens_len; i++)
        {
            yw_token_clone(&dest->ast_at_rule_tk.body_tokens[i], &tk->ast_at_rule_tk.body_tokens[i]);
        }
        dest->ast_at_rule_tk.name = yw_duplicate_str(tk->ast_at_rule_tk.name);
        break;
    case YW_CSS_TOKEN_AST_DECLARATION:
        dest->ast_declaration_tk.value_tokens = YW_ALLOC(YW_CSSToken, tk->ast_declaration_tk.value_tokens_len);
        for (int i = 0; i < tk->ast_declaration_tk.value_tokens_len; i++)
        {
            yw_token_clone(&dest->ast_declaration_tk.value_tokens[i], &tk->ast_declaration_tk.value_tokens[i]);
        }
        dest->ast_declaration_tk.name = yw_duplicate_str(tk->ast_declaration_tk.name);
        break;
    default:
        YW_UNREACHABLE();
    }
}

/***************************************************************************
 * CSS tokenizer
 **************************************************************************/
typedef struct YW_Tokenizer YW_Tokenizer;
struct YW_Tokenizer
{
    YW_TextReader tr;
};

static void yw_consume_comments(YW_Tokenizer *tkr)
{
    bool end_found = false;
    while (yw_text_reader_is_eof(&tkr->tr))
    {
        if (!yw_consume_str(&tkr->tr, "/*", YW_NO_MATCH_FLAGS))
        {
            return;
        }
        while (!yw_text_reader_is_eof(&tkr->tr))
        {
            if (yw_consume_str(&tkr->tr, "*/", YW_NO_MATCH_FLAGS))
            {
                end_found = true;
                break;
            }
            yw_consume_any_char(&tkr->tr);
        }
        if (end_found)
        {
            continue;
        }
        /* PARSE ERROR: Reached EOF without closing the comment. */
        return;
    }
}

static bool yw_consume_number(double *out, YW_Tokenizer *tkr)
{
    YW_TextCursor start_cursor = tkr->tr.cursor;
    bool have_integer_part = false;
    bool have_fractional_part = false;

    /*
     * Note that we don't parse the number directly - We only check if it's a
     * valid  number. Rest of the job is handled by the standard library.
     */

    /***************************************************************************
     * Sign
     **************************************************************************/

    yw_consume_one_of_chars(&tkr->tr, "+-");

    /***************************************************************************
     * Integer part
     **************************************************************************/
    while (!yw_text_reader_is_eof(&tkr->tr))
    {
        YW_Char32 temp_char = yw_peek_char(&tkr->tr);
        if (yw_is_ascii_digit(temp_char))
        {
            yw_consume_any_char(&tkr->tr);
            have_integer_part = true;
        }
        else
        {
            break;
        }
    }
    /***************************************************************************
     * Decimal point
     **************************************************************************/
    YW_TextCursor cursor_before_exp = tkr->tr.cursor;

    if (yw_consume_char(&tkr->tr, '.'))
    {
        /***********************************************************************
         * Fractional part
         **********************************************************************/
        int digit_count = 0;

        while (!yw_text_reader_is_eof(&tkr->tr))
        {
            YW_Char32 temp_char = yw_peek_char(&tkr->tr);
            if (yw_is_ascii_digit(temp_char))
            {
                yw_consume_any_char(&tkr->tr);
                digit_count++;
            }
            else
            {
                break;
            }
        }
        if (!have_integer_part && digit_count == 0)
        {
            tkr->tr.cursor = cursor_before_exp;
            return false;
        }
        have_fractional_part = true;
    }

    if (!have_integer_part && !have_fractional_part)
    {
        // We have invalid number
        tkr->tr.cursor = start_cursor;
        return false;
    }

    /***************************************************************************
     * Exponent indicator
     **************************************************************************/
    cursor_before_exp = tkr->tr.cursor;
    if (yw_consume_one_of_chars(&tkr->tr, "eE") != -1)
    {
        int digit_count = 0;

        /***********************************************************************
         * Exponent sign
         **********************************************************************/
        yw_consume_one_of_chars(&tkr->tr, "+-");

        /***********************************************************************
         * Exponent
         **********************************************************************/
        while (!yw_text_reader_is_eof(&tkr->tr))
        {
            YW_Char32 temp_char = yw_peek_char(&tkr->tr);
            if (yw_is_ascii_digit(temp_char))
            {
                yw_consume_any_char(&tkr->tr);
                digit_count++;
            }
            else
            {
                break;
            }
        }
        if (digit_count == 0)
        {
            tkr->tr.cursor = cursor_before_exp;
        }
    }

    YW_TextCursor end_cursor = tkr->tr.cursor;

    /***************************************************************************
     * Now we parse the number
     **************************************************************************/
    char *temp_buf = YW_ALLOC(char, end_cursor - start_cursor + 1);
    tkr->tr.cursor = start_cursor;
    while (tkr->tr.cursor < end_cursor)
    {
        int dest = tkr->tr.cursor - start_cursor;
        temp_buf[dest] = yw_consume_any_char(&tkr->tr);
    }
    temp_buf[tkr->tr.cursor - start_cursor] = '\0';
    char *nptr;
    double res = strtod(temp_buf, &nptr);
    if (*nptr != '\0')
    {
        fprintf(stderr, "%s: strtod() failed to parse some(or all) of %s\n", __func__, temp_buf);
    }
    free(temp_buf);
    *out = res;
    return true;
}

static bool yw_consume_escaped_codepoint(YW_Char32 *out, YW_Tokenizer *tkr)
{
    YW_TextCursor old_cursor = tkr->tr.cursor;
    if (!yw_consume_char(&tkr->tr, '\\'))
    {
        return false;
    }
    bool is_hex_digit = false;
    int hex_digit_val = 0;
    int hex_digit_count = 0;

    if (yw_text_reader_is_eof(&tkr->tr))
    {
        /* PARSE ERROR: Unexpected EOF */
        *out = 0xfffd;
        return true;
    }
    if (yw_consume_char(&tkr->tr, '\n'))
    {
        tkr->tr.cursor = old_cursor;
        return false;
    }
    // https://www.w3.org/TR/2021/CRD--syntax-3-20211224/#consume-an-escaped-code-point
    while ((!yw_text_reader_is_eof(&tkr->tr)) && hex_digit_count < 6)
    {
        YW_Char32 temp_char = yw_peek_char(&tkr->tr);
        int digit = 0;
        if (yw_is_ascii_digit(temp_char))
        {
            digit = temp_char - '0';
        }
        else if (yw_is_ascii_lowercase_hex_digit(temp_char))
        {
            digit = temp_char - 'a' + 10;
        }
        else if (yw_is_ascii_uppercase_hex_digit(temp_char))
        {
            digit = temp_char - 'A' + 10;
        }
        else
        {
            break;
        }
        yw_consume_any_char(&tkr->tr);
        hex_digit_val = hex_digit_val * 16 + digit;
        is_hex_digit = true;
        hex_digit_count++;
    }
    if (is_hex_digit)
    {
        *out = hex_digit_val;
    }
    else
    {
        *out = yw_consume_any_char(&tkr->tr);
    }
    return true;
}

/* Caller owns the returned string */
static bool yw_consume_ident_sequence(char **str_out, YW_Tokenizer *tkr, bool must_start_with_ident_start)
{
    char *res = NULL;
    int res_len = 0;
    int res_cap = 0;
    while (!yw_text_reader_is_eof(&tkr->tr))
    {
        YW_TextCursor old_cursor = tkr->tr.cursor;

        YW_Char32 result_chr;
        if (!yw_consume_escaped_codepoint(&result_chr, tkr))
        {
            result_chr = yw_consume_any_char(&tkr->tr);
        }
        if (yw_is_ident_start_codepoint(result_chr) || ((res_len != 0 || !must_start_with_ident_start) && yw_is_ident_codepoint(result_chr)))
        {
            if (0x7f < result_chr)
            {
                fprintf(stderr, "%s: unicode character is not supported yet\n", __func__);
                result_chr = '?';
            }
            YW_PUSH(char, &res_cap, &res_len, &res, result_chr);
        }
        else
        {
            tkr->tr.cursor = old_cursor;
            break;
        }
    }

    if (res_len == 0)
    {
        goto fail;
    }
    YW_PUSH(char, &res_cap, &res_len, &res, '\0');
    YW_SHRINK_TO_FIT(char, &res_cap, res_len, &res);
    *str_out = res;
    return true;
fail:
    free(res);
    return false;
}

static bool yw_consume_whitespace_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    bool found = false;
    while (!yw_text_reader_is_eof(&tkr->tr))
    {
        if (yw_consume_one_of_chars(&tkr->tr, " \t\n") == -1)
        {
            break;
        }
        found = true;
    }
    if (!found)
    {
        return false;
    }
    out->common.type = YW_CSS_TOKEN_WHITESPACE;
    return true;
}

static bool yw_consume_simple_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    int res = yw_consume_one_of_chars(&tkr->tr, "(),:;[]{}");
    switch (res)
    {
    case '(':
        out->common.type = YW_CSS_TOKEN_LEFT_PAREN;
        return true;
    case ')':
        out->common.type = YW_CSS_TOKEN_RIGHT_PAREN;
        return true;
    case ',':
        out->common.type = YW_CSS_TOKEN_COMMA;
        return true;
    case ':':
        out->common.type = YW_CSS_TOKEN_COLON;
        return true;
    case ';':
        out->common.type = YW_CSS_TOKEN_SEMICOLON;
        return true;
    case '[':
        out->common.type = YW_CSS_TOKEN_LEFT_SQUARE_BRACKET;
        return true;
    case ']':
        out->common.type = YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET;
        return true;
    case '{':
        out->common.type = YW_CSS_TOKEN_LEFT_CURLY_BRACKET;
        return true;
    case '}':
        out->common.type = YW_CSS_TOKEN_RIGHT_CURLY_BRACKET;
        return true;
    case -1:
        break;
    default:
        YW_UNREACHABLE();
    }
    char const *strs[] = {"<!--", "-->", NULL};
    res = yw_consume_one_of_strs(&tkr->tr, strs, YW_NO_MATCH_FLAGS);
    switch (res)
    {
    case 0:
        out->common.type = YW_CSS_TOKEN_CDO;
        return true;
    case 1:
        out->common.type = YW_CSS_TOKEN_CDC;
        return true;
    case -1:
        break;
    default:
        YW_UNREACHABLE();
    }
    return false;
}

static bool yw_consume_string_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    char ending_char;
    char *res = NULL;
    int res_len = 0;
    int res_cap = 0;

    switch (yw_consume_one_of_chars(&tkr->tr, "\"'"))
    {
    case '"':
        ending_char = '"';
        break;
    case '\'':
        ending_char = '\'';
        break;
    default:
        return false;
    }

    while (!yw_text_reader_is_eof(&tkr->tr))
    {
        YW_TextCursor cursor_before_chr = tkr->tr.cursor;
        char chars[] = {ending_char, '\n', 0};
        int chr = yw_consume_one_of_chars(&tkr->tr, chars);
        if (chr == ending_char)
        {
            break;
        }
        else if (chr == '\n')
        {
            /* PARSE ERROR: Unexpected newline*/
            tkr->tr.cursor = cursor_before_chr;
            break;
        }
        else if (yw_text_reader_is_eof(&tkr->tr))
        {
            /* PARSE ERROR: Unexpected EOF */
            break;
        }
        else if (yw_consume_one_of_chars(&tkr->tr, "\\\n") != -1)
        {
            continue;
        }
        YW_Char32 result_chr;
        if (!yw_consume_escaped_codepoint(&result_chr, tkr))
        {
            result_chr = yw_consume_any_char(&tkr->tr);
        }
        if (0x7f < result_chr)
        {
            fprintf(stderr, "%s: unicode character is not supported yet\n", __func__);
            result_chr = '?';
        }
        YW_PUSH(char, &res_cap, &res_len, &res, result_chr);
    }
    YW_PUSH(char, &res_cap, &res_len, &res, '\0');
    YW_SHRINK_TO_FIT(char, &res_cap, res_len, &res);
    out->common.type = YW_CSS_TOKEN_STRING;
    out->string_tk.value = res;
    return true;
}

static bool yw_consume_hash_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    YW_TextCursor cursor_from = tkr->tr.cursor;
    char *hash_str = NULL;
    if (!yw_consume_char(&tkr->tr, '#'))
    {
        goto fail;
    }
    if (!yw_consume_ident_sequence(&hash_str, tkr, false))
    {
        goto fail;
    }
    if (hash_str[0] == '\0')
    {
        goto fail;
    }
    out->common.type = YW_CSS_TOKEN_HASH;
    out->hash_tk.value = hash_str;
    out->hash_tk.type = yw_is_valid_ident_start_sequence(hash_str) ? YW_HASH_ID : YW_HASH_UNRESTRICTED;
    return true;
fail:
    free(hash_str);
    tkr->tr.cursor = cursor_from;
    return false;
}

static bool yw_consume_at_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    YW_TextCursor cursor_from = tkr->tr.cursor;
    char *at_str = NULL;
    if (!yw_consume_char(&tkr->tr, '@'))
    {
        return false;
    }
    if (!yw_consume_ident_sequence(&at_str, tkr, true))
    {
        goto fail;
    }
    if (at_str[0] == '\0' || !yw_is_valid_ident_start_sequence(at_str))
    {
        goto fail;
    }
    out->common.type = YW_CSS_TOKEN_AT_KEYWORD;
    out->at_keyword_tk.value = at_str;
    return true;
fail:
    free(at_str);
    tkr->tr.cursor = cursor_from;
    return false;
}

static bool yw_consume_numeric_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    YW_TextCursor cursor_from = tkr->tr.cursor;
    YW_TextCursor cursor_before_ident;
    char *ident = NULL;
    double value;
    if (!yw_consume_number(&value, tkr))
    {
        goto fail;
    }
    cursor_before_ident = tkr->tr.cursor;
    if (yw_consume_ident_sequence(&ident, tkr, true))
    {
        if (yw_is_valid_ident_start_sequence(ident))
        {
            out->common.type = YW_CSS_TOKEN_DIMENSION;
            out->dimension_tk.unit = ident;
            out->dimension_tk.value = value;
            return true;
        }
        else
        {
            tkr->tr.cursor = cursor_before_ident;
            free(ident);
        }
    }
    if (yw_consume_char(&tkr->tr, '%'))
    {
        out->common.type = YW_CSS_TOKEN_PERCENTAGE;
        out->percentage_tk.value = value;
        return true;
    }
    out->common.type = YW_CSS_TOKEN_NUMBER;
    out->number_tk.value = value;
    return true;
fail:
    free(ident);
    tkr->tr.cursor = cursor_from;
    return false;
}

static void yw_consume_remnants_of_bad_url(YW_Tokenizer *tkr)
{
    while (!yw_text_reader_is_eof(&tkr->tr))
    {
        if (yw_consume_char(&tkr->tr, ')'))
        {
            break;
        }
        YW_Char32 cp;
        if (!yw_consume_escaped_codepoint(&cp, tkr))
        {
            yw_consume_any_char(&tkr->tr);
        }
    }
}

static bool yw_consume_ident_like_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    YW_TextCursor cursor_from = tkr->tr.cursor;
    YW_TextCursor cursor_after_ident;
    char *ident = NULL;
    char *url = NULL;
    if (!yw_consume_ident_sequence(&ident, tkr, true))
    {
        goto fail;
    }

    cursor_after_ident = tkr->tr.cursor;
    if (yw_strcmp_ascii_case_insensitive(ident, "url") == 0 && yw_consume_char(&tkr->tr, '('))
    {
        while (yw_consume_str(&tkr->tr, "  ", YW_NO_MATCH_FLAGS))
        {
        }
        YW_TextCursor old_cursor = tkr->tr.cursor;
        if (yw_consume_one_of_chars(&tkr->tr, "\"'") || yw_consume_str(&tkr->tr, " \"", YW_NO_MATCH_FLAGS) || yw_consume_str(&tkr->tr, " '", YW_NO_MATCH_FLAGS))
        {
            /* Function token *************************************************/
            tkr->tr.cursor = old_cursor;
            out->common.type = YW_CSS_TOKEN_FUNC_KEYWORD;
            out->func_keyword_tk.value = ident;
            return true;
        }
        else
        {
            free(ident);
            // URL token ****************************************/----------
            YW_CSSToken tk;
            yw_consume_whitespace_token(&tk, tkr);

            url = NULL;
            int url_len = 0;
            int url_cap = 0;
            while (1)
            {
                if (yw_text_reader_is_eof(&tkr->tr))
                {
                    /* PARSE ERROR: Unexpected EOF */
                    goto done;
                }
                switch (yw_consume_one_of_chars(&tkr->tr, ")\"'("))
                {
                case ')':
                    goto done;
                case '"':
                case '\'':
                case '(':
                    /* PARSE ERROR: Unexpected character */
                    yw_consume_remnants_of_bad_url(tkr);
                    out->common.type = YW_CSS_TOKEN_BAD_URL;
                    return true;
                default: {
                    YW_Char32 escaped_chr;
                    YW_Char32 temp;
                    if (!yw_consume_escaped_codepoint(&temp, tkr))
                    {
                        escaped_chr = temp;
                    }
                    else if (yw_consume_char(&tkr->tr, '\\'))
                    {
                        /* PARSE ERROR: Unexpected character after \ */
                        yw_consume_remnants_of_bad_url(tkr);
                        out->common.type = YW_CSS_TOKEN_BAD_URL;
                        return true;
                    }
                    else
                    {
                        escaped_chr = yw_consume_any_char(&tkr->tr);
                    }
                    if (0x7f < escaped_chr)
                    {
                        fprintf(stderr, "%s: unicode character is not supported yet\n", __func__);
                        escaped_chr = '?';
                    }
                    YW_PUSH(char, &url_cap, &url_len, &url, escaped_chr);
                }
                }
            }
        done:
            YW_SHRINK_TO_FIT(char, &url_cap, url_len, &url);
            out->common.type = YW_CSS_TOKEN_URL;
            out->url_tk.value = url;
            return true;
        }
    }

    tkr->tr.cursor = cursor_after_ident;
    if (yw_consume_char(&tkr->tr, '('))
    {
        out->common.type = YW_CSS_TOKEN_FUNC_KEYWORD;
        out->func_keyword_tk.value = ident;
        return true;
    }

    out->common.type = YW_CSS_TOKEN_IDENT;
    out->ident_tk.value = ident;
    return true;
fail:
    free(ident);
    free(url);
    tkr->tr.cursor = cursor_from;
    return false;
}

static bool yw_consume_any_token(YW_CSSToken *out, YW_Tokenizer *tkr)
{
    typedef bool(TokenFunc)(YW_CSSToken * out, YW_Tokenizer * tkr);
    static TokenFunc *funcs[] = {
        yw_consume_whitespace_token,
        yw_consume_string_token,
        yw_consume_hash_token,
        yw_consume_at_token,
        yw_consume_simple_token,
        yw_consume_numeric_token,
        yw_consume_ident_like_token,
    };

    yw_consume_comments(tkr);
    for (int i = 0; i < (int)YW_SIZEOF_ARRAY(funcs); i++)
    {
        if (funcs[i](out, tkr))
        {
            return true;
        }
    }
    if (yw_text_reader_is_eof(&tkr->tr))
    {
        return false;
    }
    out->common.type = YW_CSS_TOKEN_DELIM;
    out->delim_tk.value = yw_consume_any_char(&tkr->tr);
    return true;
}

static YW_EncodingType yw_css_determine_fallback_encoding(uint8_t const *bytes, int bytes_len)
{
    /***************************************************************************
     * Check if HTTP or equivalent protocol provides an encoding label
     **************************************************************************/

    /* TODO */

    /***************************************************************************
     * Check '@charset "xxx";' byte sequence
     **************************************************************************/
    bytes_len = (1024 < bytes_len) ? 1024 : bytes_len;

    if (bytes_len < 10 || memcmp(bytes, "@charset \"", 10) != 0)
    {
        char const *start = &((char const *)bytes)[10];
        char const *end = strstr((char *)bytes, "\";");
        if (end != NULL)
        {
            int len = (end - 1) - start;
            char *label = YW_ALLOC(char, len + 1);
            memcpy(label, start, sizeof(*label) * len);
            YW_EncodingType enc = yw_encoding_from_label(label);
            if (enc == YW_UTF16_BE || enc == YW_UTF16_LE)
            {
                /* This is not a bug. The standard says to do this. */
                return YW_UTF8;
            }
            else if (enc != YW_INVALID_ENCODING)
            {
                return enc;
            }
        }
    }

    /***************************************************************************
     * Check if environment encoding is provided
     **************************************************************************/

    /* TODO */

    /***************************************************************************
     * Fallback to UTF-8
     **************************************************************************/
    return YW_UTF8;
}

/* Caller owns the returned array */
static void yw_css_decode_bytes(uint8_t **buf_out, int *len_out, uint8_t const *bytes, int bytes_len)
{
    YW_IOQueue input, output;
    YW_EncodingType fallback = yw_css_determine_fallback_encoding(bytes, bytes_len);
    yw_io_queue_init(&input);
    yw_io_queue_init(&output);
    for (int i = 0; i < bytes_len; i++)
    {
        yw_io_queue_push_one(&input, bytes[i]);
    }
    yw_encoding_decode(&input, fallback, &output);

    yw_io_queue_to_utf8(buf_out, len_out, &output);
    yw_io_queue_deinit(&input);
    yw_io_queue_deinit(&output);
}

/* Caller owns the returned array */
static void yw_css_filter_codepoints(uint8_t **buf_out, int *len_out, uint8_t const *in, int in_len)
{
    /* https://www.w3.org/TR/-syntax-3/#-filter-code-points */
    uint8_t *res = NULL;
    int res_len = 0;
    int res_cap = 0;
    for (int i = 0; i < in_len; i++)
    {
        if ((i + 2 != in_len) && in[i] == '\r' && in[i + 1] == '\n')
        {
            /* CR followed by LF */
            continue;
        }
        else if (in[i] == '\r' || in[i] == '\x0c')
        {
            YW_PUSH(uint8_t, &res_cap, &res_len, &res, '\n');
        }
        else if (in[i] == '\x00' || yw_is_surrogate_char(in[i]))
        {
            /* Push U+FFFD in UTF-8 form */
            YW_PUSH(uint8_t, &res_cap, &res_len, &res, 0xef);
            YW_PUSH(uint8_t, &res_cap, &res_len, &res, 0xbf);
            YW_PUSH(uint8_t, &res_cap, &res_len, &res, 0xbd);
        }
        else
        {
            YW_PUSH(uint8_t, &res_cap, &res_len, &res, in[i]);
        }
    }
    YW_SHRINK_TO_FIT(uint8_t, &res_cap, res_len, &res);
    *buf_out = res;
    *len_out = res_len;
}

static void yw_parse_list_of_component_values(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts);

bool yw_css_tokenize(YW_CSSTokenStream *out, uint8_t const *bytes, int bytes_len)
{
    uint8_t *src;
    int src_len;
    YW_Tokenizer tkr;
    YW_CSSTokenStream temp_ts;
    YW_CSSToken *res = NULL, *new_res = NULL;
    int res_len = 0, new_res_len = 0;
    int res_cap = 0;

    memset(&tkr, 0, sizeof(tkr));
    /* Decode bytes and filter codepoints *************************************/
    yw_css_decode_bytes(&src, &src_len, bytes, bytes_len);
    {
        uint8_t *new_src;
        int new_src_len;
        yw_css_filter_codepoints(&new_src, &new_src_len, src, src_len);
        free(src);
        src = new_src;
        src_len = new_src_len;
    }
    /* Consume tokens *********************************************************/
    yw_text_reader_init(&tkr.tr, src, src_len);
    YW_CSSToken tk;
    while (yw_consume_any_token(&tk, &tkr))
    {
        YW_PUSH(YW_CSSToken, &res_cap, &res_len, &res, tk);
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &res_cap, res_len, &res);
    bool is_eof = yw_text_reader_is_eof(&tkr.tr);
    free(src);
    if (!is_eof)
    {
        goto fail;
    }
    /* Create temporary token stream ******************************************/
    temp_ts.tokens = res;
    temp_ts.tokens_len = res_len;
    temp_ts.cursor = 0;
    /* Parse component values *************************************************/
    /* This turns tokens into higher-level objects (if possible) */
    yw_parse_list_of_component_values(&new_res, &new_res_len, &temp_ts);
    /* Clear old tokens *******************************************************/
    for (int i = 0; i < res_len; i++)
    {
        yw_css_token_deinit(&res[i]);
    }
    free(res);
    /* Prepare output and exit ************************************************/
    memset(out, 0, sizeof(*out));
    out->tokens = new_res;
    out->tokens_len = new_res_len;
    return true;
fail:
    for (int i = 0; i < res_len; i++)
    {
        yw_css_token_deinit(&res[i]);
    }
    free(res);
    for (int i = 0; i < new_res_len; i++)
    {
        yw_css_token_deinit(&new_res[i]);
    }
    free(new_res);
    return false;
}

/***************************************************************************
 * Main parsing code
 **************************************************************************/

bool yw_is_end_of_tokens(YW_CSSTokenStream const *ts)
{
    return ts->tokens_len <= ts->cursor;
}

YW_CSSToken const *yw_expect_any_token(YW_CSSTokenStream *ts)
{
    if (yw_is_end_of_tokens(ts))
    {
        return NULL;
    }
    ts->cursor++;
    return &ts->tokens[ts->cursor - 1];
}

YW_CSSToken const *yw_expect_token(YW_CSSTokenStream *ts, YW_TokenType type)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *tk = yw_expect_any_token(ts);
    if (tk == NULL || tk->common.type != type)
    {
        goto fail;
    }
    return tk;
fail:
    ts->cursor = old_cursor;
    return NULL;
}

bool yw_expect_delim(YW_CSSTokenStream *ts, YW_Char32 d)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *token = yw_expect_token(ts, YW_CSS_TOKEN_DELIM);
    if (token == NULL || token->delim_tk.value != d)
    {
        goto fail;
    }
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

bool yw_expect_ident(YW_CSSTokenStream *ts, char const *i)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *token = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (token == NULL || strcmp(token->ident_tk.value, i) != 0)
    {
        goto fail;
    }
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

bool yw_expect_simple_block(YW_CSSTokenStream *inner_ts_out, YW_CSSTokenStream *ts, YW_SimpleBlockType type)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *token = yw_expect_token(ts, YW_CSS_TOKEN_AST_SIMPLE_BLOCK);
    if (token == NULL || token->ast_simple_block_tk.type != type)
    {
        goto fail;
    }
    inner_ts_out->tokens = token->ast_simple_block_tk.tokens;
    inner_ts_out->tokens_len = token->ast_simple_block_tk.tokens_len;
    inner_ts_out->cursor = 0;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

bool yw_expect_ast_func(YW_CSSTokenStream *inner_ts_out, YW_CSSTokenStream *ts, char const *f)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *token = yw_expect_token(ts, YW_CSS_TOKEN_AST_FUNC);
    if (token == NULL || strcmp(token->ast_func_tk.name, f) != 0)
    {
        goto fail;
    }
    inner_ts_out->tokens = token->ast_func_tk.tokens;
    inner_ts_out->tokens_len = token->ast_func_tk.tokens_len;
    inner_ts_out->cursor = 0;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

void yw_skip_whitespaces(YW_CSSTokenStream *ts)
{
    while (1)
    {
        int old_cursor = ts->cursor;
        if (yw_expect_token(ts, YW_CSS_TOKEN_WHITESPACE) == NULL)
        {
            ts->cursor = old_cursor;
            break;
        }
    }
}

static bool yw_consume_preserved_token(YW_CSSToken *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *token = yw_expect_any_token(ts);
    if (token == NULL ||
        token->common.type == YW_CSS_TOKEN_FUNC_KEYWORD ||
        token->common.type == YW_CSS_TOKEN_LEFT_CURLY_BRACKET ||
        token->common.type == YW_CSS_TOKEN_LEFT_SQUARE_BRACKET ||
        token->common.type == YW_CSS_TOKEN_LEFT_PAREN)
    {
        goto fail;
    }
    yw_token_clone(out, token);
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

static bool yw_consume_component_value(YW_CSSToken *out, YW_CSSTokenStream *ts);

static bool yw_consume_simple_block(YW_AstSimpleBlockToken *out, YW_CSSTokenStream *ts, YW_TokenType open_token_type)
{
    int old_cursor = ts->cursor;
    YW_TokenType close_token_type;
    YW_SimpleBlockType type;
    YW_CSSToken *res = NULL;
    int res_cap = 0;
    int res_len = 0;

    switch (open_token_type)
    {
    case YW_CSS_TOKEN_LEFT_CURLY_BRACKET:
        close_token_type = YW_CSS_TOKEN_RIGHT_CURLY_BRACKET;
        type = YW_SIMPLE_BLOCK_CURLY;
        break;
    case YW_CSS_TOKEN_LEFT_SQUARE_BRACKET:
        close_token_type = YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET;
        type = YW_SIMPLE_BLOCK_SQUARE;
        break;
    case YW_CSS_TOKEN_LEFT_PAREN:
        close_token_type = YW_CSS_TOKEN_RIGHT_PAREN;
        type = YW_SIMPLE_BLOCK_PAREN;
        break;
    default:
        fprintf(stderr, "unsupported open_token_type\n");
        abort();
    }

    YW_CSSToken const *open_token = yw_expect_token(ts, open_token_type);
    YW_CSSToken const *close_token = NULL;
    if (open_token == NULL)
    {
        goto fail;
    }
    while (1)
    {
        YW_CSSToken temp_tk;
        if (!yw_consume_component_value(&temp_tk, ts))
        {
            break;
        }
        if (temp_tk.common.type == close_token_type)
        {
            close_token = &res[res_len - 1];
            break;
        }
        YW_PUSH(YW_CSSToken, &res_cap, &res_len, &res, temp_tk);
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &res_cap, res_len, &res);
    if (close_token == NULL)
    {
        goto fail;
    }
    out->common.type = YW_CSS_TOKEN_AST_SIMPLE_BLOCK;
    out->type = type;
    out->tokens = res;
    out->tokens_len = res_len;
    return true;
fail:
    free(res);
    ts->cursor = old_cursor;
    return false;
}
static bool yw_consume_curly_block(YW_AstSimpleBlockToken *out, YW_CSSTokenStream *ts)
{
    return yw_consume_simple_block(out, ts, YW_CSS_TOKEN_LEFT_CURLY_BRACKET);
}
static bool yw_consume_square_block(YW_AstSimpleBlockToken *out, YW_CSSTokenStream *ts)
{
    return yw_consume_simple_block(out, ts, YW_CSS_TOKEN_LEFT_SQUARE_BRACKET);
}
static bool yw_consume_paren_block(YW_AstSimpleBlockToken *out, YW_CSSTokenStream *ts)
{
    return yw_consume_simple_block(out, ts, YW_CSS_TOKEN_LEFT_PAREN);
}
static bool yw_consume_func(YW_AstFunctionToken *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSToken *res = NULL;
    int res_cap = 0;
    int res_len = 0;

    YW_CSSToken const *func_token = yw_expect_token(ts, YW_CSS_TOKEN_FUNC_KEYWORD);
    YW_CSSToken const *close_token = NULL;
    if (func_token == NULL)
    {
        goto fail;
    }
    while (1)
    {
        YW_CSSToken temp_tk;
        if (!yw_consume_component_value(&temp_tk, ts))
        {
            break;
        }
        if (temp_tk.common.type == YW_CSS_TOKEN_RIGHT_PAREN)
        {
            close_token = &res[res_len - 1];
            break;
        }
        YW_PUSH(YW_CSSToken, &res_cap, &res_len, &res, temp_tk);
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &res_cap, res_len, &res);
    if (close_token == NULL)
    {
        goto fail;
    }
    out->common.type = YW_CSS_TOKEN_AST_FUNC;
    out->name = yw_duplicate_str(func_token->func_keyword_tk.value);
    out->tokens = res;
    out->tokens_len = res_len;
    return true;
fail:
    free(res);
    ts->cursor = old_cursor;
    return false;
}
static bool yw_consume_component_value(YW_CSSToken *out, YW_CSSTokenStream *ts)
{
    return yw_consume_curly_block(&out->ast_simple_block_tk, ts) ||
           yw_consume_square_block(&out->ast_simple_block_tk, ts) ||
           yw_consume_paren_block(&out->ast_simple_block_tk, ts) ||
           yw_consume_func(&out->ast_func_tk, ts) ||
           yw_consume_preserved_token(out, ts);
}
static bool yw_consume_qualified_rule(YW_AstQualifiedRuleToken *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSToken *prelude_tokens = NULL;
    int prelude_tokens_len = 0;
    int prelude_tokens_cap = 0;

    while (1)
    {
        YW_AstSimpleBlockToken temp_block;
        if (yw_consume_curly_block(&temp_block, ts))
        {
            out->common.type = YW_CSS_TOKEN_AST_QUALIFIED_RULE;
            out->body_tokens = temp_block.tokens;
            out->body_tokens_len = temp_block.tokens_len;
            out->prelude_tokens = prelude_tokens;
            out->prelude_tokens_len = prelude_tokens_len;

            temp_block.tokens = NULL;
            temp_block.tokens_len = 0;
            yw_css_token_deinit((YW_CSSToken *)&temp_block);
            return true;
        }
        YW_CSSToken temp_prelude;
        if (!yw_consume_component_value(&temp_prelude, ts))
        {
            goto fail;
        }
        YW_PUSH(YW_CSSToken, &prelude_tokens_cap, &prelude_tokens_len, &prelude_tokens, temp_prelude);
    }

fail:
    for (int i = 0; i < prelude_tokens_len; i++)
    {
        yw_css_token_deinit(&prelude_tokens[i]);
    }
    free(prelude_tokens);
    ts->cursor = old_cursor;
    return false;
}
static bool yw_consume_at_rule(YW_AstAtRuleToken *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSToken *prelude_tokens = NULL;
    int prelude_tokens_len = 0;
    int prelude_tokens_cap = 0;
    YW_CSSToken const *kwd_token = yw_expect_token(ts, YW_CSS_TOKEN_AT_KEYWORD);

    while (1)
    {
        YW_AstSimpleBlockToken temp_block;
        if (yw_consume_curly_block(&temp_block, ts))
        {
            out->common.type = YW_CSS_TOKEN_AST_AT_RULE;
            out->name = yw_duplicate_str(kwd_token->at_keyword_tk.value);
            out->body_tokens = temp_block.tokens;
            out->body_tokens_len = temp_block.tokens_len;
            out->prelude_tokens = prelude_tokens;
            out->prelude_tokens_len = prelude_tokens_len;

            temp_block.tokens = NULL;
            temp_block.tokens_len = 0;
            yw_css_token_deinit((YW_CSSToken *)&temp_block);
            return true;
        }
        YW_CSSToken temp_prelude;
        if (!yw_consume_component_value(&temp_prelude, ts))
        {
            goto fail;
        }
        YW_PUSH(YW_CSSToken, &prelude_tokens_cap, &prelude_tokens_len, &prelude_tokens, temp_prelude);
    }
fail:
    for (int i = 0; i < prelude_tokens_len; i++)
    {
        yw_css_token_deinit(&prelude_tokens[i]);
    }
    free(prelude_tokens);
    ts->cursor = old_cursor;
    return false;
}
static bool yw_consume_declaration(YW_AstDeclarationToken *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;

    YW_CSSToken *decl_values = NULL;
    int decl_values_len = 0;
    int decl_values_cap = 0;
    bool decl_is_important = false;
    char const *decl_name;

    /* <name>  :  contents  !important ****************************************/
    YW_CSSToken const *temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (temp_tk == NULL)
    {
        goto fail;
    }
    decl_name = temp_tk->ident_tk.value;

    /* name<  >:  contents  !important ****************************************/
    yw_skip_whitespaces(ts);
    /* name  <:>  contents  !important ****************************************/
    if (yw_expect_token(ts, YW_CSS_TOKEN_COLON) == NULL)
    {
        goto fail;
    }
    /* name  :<  >contents  !important ****************************************/
    yw_skip_whitespaces(ts);
    /* name  :  <contents  !important> ****************************************/
    while (1)
    {
        YW_CSSToken temp_tk;
        if (!yw_consume_component_value(&temp_tk, ts))
        {
            break;
        }
        YW_PUSH(YW_CSSToken, &decl_values_cap, &decl_values_len, &decl_values, temp_tk);
    }
    if (2 <= decl_values_len)
    {
        // See if we have !important
        YW_CSSToken const *ptk1 = &decl_values[decl_values_len - 2];
        YW_CSSToken const *ptk2 = &decl_values[decl_values_len - 1];
        if (ptk1->common.type == YW_CSS_TOKEN_DELIM && ptk1->delim_tk.value == '!' && ptk2->common.type == YW_CSS_TOKEN_IDENT && strcmp(ptk2->ident_tk.value, "important") == 0)
        {
            decl_values_len -= 2;
            decl_is_important = true;
        }
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &decl_values_cap, decl_values_len, &decl_values);
    out->common.type = YW_CSS_TOKEN_AST_DECLARATION;
    out->name = yw_duplicate_str(decl_name);
    out->value_tokens = decl_values;
    out->value_tokens_len = decl_values_len;
    out->important = decl_is_important;
    return true;
fail:
    for (int i = 0; i < decl_values_len; i++)
    {
        yw_css_token_deinit(&decl_values[i]);
    }
    free(decl_values);
    ts->cursor = old_cursor;
    return false;
}

static bool yw_consume_declaration_list(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;

    YW_CSSToken *decls = NULL;
    int decls_len = 0;
    int decls_cap = 0;
    while (1)
    {
        int old_cursor = ts->cursor;
        YW_CSSToken const *token = yw_expect_any_token(ts);
        if (token == NULL)
        {
            break;
        }
        if (token->common.type == YW_CSS_TOKEN_WHITESPACE)
        {
            continue;
        }
        else if (token->common.type == YW_CSS_TOKEN_AT_KEYWORD)
        {
            ts->cursor = old_cursor;
            YW_CSSToken res;
            if (!yw_consume_at_rule(&res.ast_at_rule_tk, ts))
            {
                YW_UNREACHABLE();
            }
            YW_PUSH(YW_CSSToken, &decls_cap, &decls_len, &decls, res);
        }
        else if (token->common.type == YW_CSS_TOKEN_IDENT)
        {
            YW_CSSToken *tokens = NULL;
            int tokens_len = 0;
            int tokens_cap = 0;

            YW_PUSH(YW_CSSToken, &tokens_cap, &tokens_len, &tokens, *token);
            while (1)
            {
                int old_cursor = ts->cursor;
                YW_CSSToken const *token = yw_expect_any_token(ts);
                if (token == NULL)
                {
                    break;
                }
                if (token->common.type == YW_CSS_TOKEN_SEMICOLON)
                {
                    ts->cursor = old_cursor;
                    break;
                }
                YW_PUSH(YW_CSSToken, &tokens_cap, &tokens_len, &tokens, *token);
            }
            YW_SHRINK_TO_FIT(YW_CSSToken, &tokens_cap, tokens_len, &tokens);
            YW_CSSTokenStream inner_ts;
            inner_ts.tokens = tokens;
            inner_ts.tokens_len = tokens_len;
            YW_CSSToken decl;
            bool ok = yw_consume_declaration(&decl.ast_declaration_tk, &inner_ts);
            for (int i = 0; i < tokens_len; i++)
            {
                yw_css_token_deinit(&tokens[i]);
            }
            free(tokens);
            if (!ok)
            {
                break;
            }
            else
            {
                YW_PUSH(YW_CSSToken, &decls_cap, &decls_len, &decls, decl);
            }
        }
        else
        {
            /* PARSE ERROR */
            while (1)
            {
                int old_cursor = ts->cursor;
                YW_CSSToken const *token = yw_expect_any_token(ts);
                if (token == NULL)
                {
                    break;
                }
                ts->cursor = old_cursor;
                if (token->common.type == YW_CSS_TOKEN_SEMICOLON)
                {
                    break;
                }
                YW_CSSToken unused;
                if (yw_consume_component_value(&unused, ts))
                {
                    yw_css_token_deinit(&unused);
                }
            }
        }
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &decls_cap, decls_len, &decls);
    if (decls_len == 0)
    {
        goto fail;
    }
    *tokens_out = decls;
    *len_out = decls_len;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}
static bool yw_consume_style_block_contents(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;

    YW_CSSToken *decls = NULL;
    int decls_len = 0;
    int decls_cap = 0;

    YW_AstQualifiedRuleToken *rules = NULL;
    int rules_len = 0;
    int rules_cap = 0;

    while (1)
    {
        int old_cursor = ts->cursor;
        YW_CSSToken const *token = yw_expect_any_token(ts);
        if (token == NULL)
        {
            break;
        }
        if (token->common.type == YW_CSS_TOKEN_WHITESPACE)
        {
            continue;
        }
        else if (token->common.type == YW_CSS_TOKEN_AT_KEYWORD)
        {
            ts->cursor = old_cursor;
            YW_CSSToken res;
            if (!yw_consume_at_rule(&res.ast_at_rule_tk, ts))
            {
                YW_UNREACHABLE();
            }
            YW_PUSH(YW_CSSToken, &decls_cap, &decls_len, &decls, res);
        }
        else if (token->common.type == YW_CSS_TOKEN_IDENT)
        {
            YW_CSSToken *tokens = NULL;
            int tokens_len = 0;
            int tokens_cap = 0;

            YW_PUSH(YW_CSSToken, &tokens_cap, &tokens_len, &tokens, *token);
            while (1)
            {
                int old_cursor = ts->cursor;
                YW_CSSToken const *token = yw_expect_any_token(ts);
                if (token == NULL)
                {
                    break;
                }
                if (token->common.type == YW_CSS_TOKEN_SEMICOLON)
                {
                    ts->cursor = old_cursor;
                    break;
                }
                YW_PUSH(YW_CSSToken, &tokens_cap, &tokens_len, &tokens, *token);
            }
            YW_SHRINK_TO_FIT(YW_CSSToken, &tokens_cap, tokens_len, &tokens);
            YW_CSSTokenStream inner_ts;
            inner_ts.tokens = tokens;
            inner_ts.tokens_len = tokens_len;
            YW_CSSToken decl;
            bool ok = yw_consume_declaration(&decl.ast_declaration_tk, &inner_ts);
            for (int i = 0; i < tokens_len; i++)
            {
                yw_css_token_deinit(&tokens[i]);
            }
            free(tokens);
            if (!ok)
            {
                break;
            }
            else
            {
                YW_PUSH(YW_CSSToken, &decls_cap, &decls_len, &decls, decl);
            }
        }
        else if (token->common.type == YW_CSS_TOKEN_DELIM && token->delim_tk.value == '&')
        {
            ts->cursor = old_cursor;
            YW_AstQualifiedRuleToken res;
            if (yw_consume_qualified_rule(&res, ts))
            {
                YW_PUSH(YW_AstQualifiedRuleToken, &rules_cap, &rules_len, &rules, res);
            }
        }
        else
        {
            /* PARSE ERROR */
            while (1)
            {
                int old_cursor = ts->cursor;
                YW_CSSToken const *token = yw_expect_any_token(ts);
                if (token == NULL)
                {
                    break;
                }
                ts->cursor = old_cursor;
                if (token->common.type == YW_CSS_TOKEN_SEMICOLON)
                {
                    break;
                }
                YW_CSSToken unused;
                if (yw_consume_component_value(&unused, ts))
                {
                    yw_css_token_deinit(&unused);
                }
            }
        }
    }
    YW_SHRINK_TO_FIT(YW_AstQualifiedRuleToken, &rules_cap, rules_len, &rules);
    for (int i = 0; i < rules_len; i++)
    {
        YW_GROW(YW_CSSToken, &decls_cap, &decls_len, &decls, 1);
        decls[decls_len - 1].ast_qualified_rule_tk = rules[i];
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &decls_cap, decls_len, &decls);
    if (decls_len == 0)
    {
        goto fail;
    }
    *tokens_out = decls;
    *len_out = decls_len;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

typedef enum
{
    YW_TOP_LEVEL,
    YW_NOT_TOP_LEVEL,
} YW_TopLevelFlag;

static bool yw_consume_list_of_rules(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts, YW_TopLevelFlag top_level)
{
    int old_cursor = ts->cursor;

    YW_CSSToken *rules = NULL;
    int rules_len = 0;
    int rules_cap = 0;

    while (1)
    {
        int old_cursor = ts->cursor;
        YW_CSSToken const *token = yw_expect_any_token(ts);
        if (token == NULL)
        {
            break;
        }
        if (token->common.type == YW_CSS_TOKEN_WHITESPACE)
        {
            continue;
        }
        else if (token->common.type == YW_CSS_TOKEN_CDO || token->common.type == YW_CSS_TOKEN_CDC)
        {
            if (top_level == YW_TOP_LEVEL)
            {
                continue;
            }
            ts->cursor = old_cursor;
            YW_CSSToken res;
            if (yw_consume_qualified_rule(&res.ast_qualified_rule_tk, ts))
            {
                YW_PUSH(YW_CSSToken, &rules_cap, &rules_len, &rules, res);
            }
        }
        else if (token->common.type == YW_CSS_TOKEN_AT_KEYWORD)
        {
            ts->cursor = old_cursor;
            YW_CSSToken res;
            if (!yw_consume_at_rule(&res.ast_at_rule_tk, ts))
            {
                YW_UNREACHABLE();
            }
            YW_PUSH(YW_CSSToken, &rules_cap, &rules_len, &rules, res);
        }
        else
        {
            YW_CSSToken res;
            if (yw_consume_qualified_rule(&res.ast_qualified_rule_tk, ts))
            {
                YW_PUSH(YW_CSSToken, &rules_cap, &rules_len, &rules, res);
            }
            else
            {
                break;
            }
        }
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &rules_cap, rules_len, &rules);
    if (rules_len == 0)
    {
        goto fail;
    }
    *tokens_out = rules;
    *len_out = rules_len;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

static void yw_parse_list_of_component_values(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts)
{
    YW_CSSToken *res = NULL;
    int res_len = 0;
    int res_cap = 0;

    while (1)
    {
        YW_CSSToken token;
        if (!yw_consume_component_value(&token, ts))
        {
            break;
        }
        YW_PUSH(YW_CSSToken, &res_cap, &res_len, &res, token);
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &res_cap, res_len, &res);
    *tokens_out = res;
    *len_out = res_len;
}

static bool yw_consume_declaration_value_impl(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts, bool any_value)
{
    /*
     * https://www.w3.org/TR/2021/CRD--syntax-3-20211224/#typedef-declaration-value
     */
    int old_cursor = ts->cursor;

    YW_CSSToken *res = NULL;
    int res_len = 0;
    int res_cap = 0;

    YW_TokenType *open_block_tokens = NULL;
    int open_block_tokens_len = 0;
    int open_block_tokens_cap = 0;

    while (1)
    {
        int old_cursor = ts->cursor;
        YW_CSSToken const *token = yw_expect_any_token(ts);
        if (token == NULL)
        {
            break;
        }
        if (token->common.type == YW_CSS_TOKEN_BAD_STRING || token->common.type == YW_CSS_TOKEN_BAD_URL ||
            /*
             * https://www.w3.org/TR/2021/CRD--syntax-3-20211224/#typedef-any-value
             */
            (!any_value && (token->common.type == YW_CSS_TOKEN_SEMICOLON || (token->common.type == YW_CSS_TOKEN_DELIM && token->delim_tk.value == '!'))))
        {
            ts->cursor = old_cursor;
            break;
        }
        /* If we have block opening token, push it to the stack. */
        if (token->common.type == YW_CSS_TOKEN_LEFT_PAREN || token->common.type == YW_CSS_TOKEN_LEFT_SQUARE_BRACKET || token->common.type == YW_CSS_TOKEN_LEFT_CURLY_BRACKET)
        {
            YW_PUSH(YW_TokenType, &open_block_tokens_cap, &open_block_tokens_len, &open_block_tokens, token->common.type);
        }
        /* If we have block closing token, see if we have unmatched token. */
        else if (token->common.type == YW_CSS_TOKEN_RIGHT_PAREN || token->common.type == YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET || token->common.type == YW_CSS_TOKEN_RIGHT_CURLY_BRACKET)
        {
            if (open_block_tokens_len == 0)
            {
                break;
            }
            YW_TokenType last = open_block_tokens[open_block_tokens_len - 1];
            if ((token->common.type == YW_CSS_TOKEN_RIGHT_PAREN && last != YW_CSS_TOKEN_LEFT_PAREN) || (token->common.type == YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET && last != YW_CSS_TOKEN_LEFT_SQUARE_BRACKET) || (token->common.type == YW_CSS_TOKEN_RIGHT_CURLY_BRACKET && last != YW_CSS_TOKEN_LEFT_CURLY_BRACKET))
            {
                break;
            }
        }
        YW_GROW(YW_CSSToken, &res_cap, &res_len, &res, 1);
        yw_token_clone(&res[res_len - 1], token);
    }
    YW_SHRINK_TO_FIT(YW_CSSToken, &res_cap, res_len, &res);
    if (res_len == 0)
    {
        goto fail;
    }
    *tokens_out = res;
    *len_out = res_len;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}
bool yw_consume_declaration_value(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts)
{
    return yw_consume_declaration_value_impl(tokens_out, len_out, ts, false);
}
bool yw_consume_any_value(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts)
{
    return yw_consume_declaration_value_impl(tokens_out, len_out, ts, true);
}
