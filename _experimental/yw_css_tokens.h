/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_CSS_TOKENS_H_
#define YW_CSS_TOKENS_H_
#include "yw_common.h"
#include "yw_css.h"
#include <stddef.h>
#include <stdint.h>

typedef enum
{
    YW_CSS_TOKEN_WHITESPACE,
    YW_CSS_TOKEN_LEFT_PAREN,
    YW_CSS_TOKEN_RIGHT_PAREN,
    YW_CSS_TOKEN_COMMA,
    YW_CSS_TOKEN_COLON,
    YW_CSS_TOKEN_SEMICOLON,
    YW_CSS_TOKEN_LEFT_SQUARE_BRACKET,
    YW_CSS_TOKEN_RIGHT_SQUARE_BRACKET,
    YW_CSS_TOKEN_LEFT_CURLY_BRACKET,
    YW_CSS_TOKEN_RIGHT_CURLY_BRACKET,
    YW_CSS_TOKEN_CDO,
    YW_CSS_TOKEN_CDC,
    YW_CSS_TOKEN_BAD_STRING,
    YW_CSS_TOKEN_BAD_URL,
    YW_CSS_TOKEN_NUMBER,
    YW_CSS_TOKEN_PERCENTAGE,
    YW_CSS_TOKEN_DIMENSION,
    YW_CSS_TOKEN_STRING,
    YW_CSS_TOKEN_URL,
    YW_CSS_TOKEN_AT_KEYWORD,
    YW_CSS_TOKEN_FUNC_KEYWORD,
    YW_CSS_TOKEN_IDENT,
    YW_CSS_TOKEN_HASH,
    YW_CSS_TOKEN_DELIM,

    /* High-level objects *****************************************************/

    YW_CSS_TOKEN_AST_SIMPLE_BLOCK,
    YW_CSS_TOKEN_AST_FUNC,
    YW_CSS_TOKEN_AST_QUALIFIED_RULE,
    YW_CSS_TOKEN_AST_AT_RULE,
    YW_CSS_TOKEN_AST_DECLARATION,
} YW_TokenType;

typedef struct YW_TokenCommon
{
    YW_TokenType type;
    int cursor_from, cursor_to;
} YW_TokenCommon;
typedef struct YW_CSSNumberToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_NUMBER */
    double value;
} YW_CSSNumberToken;
typedef struct YW_CSSPercentageToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_PERCENTAGE */
    double value;
} YW_CSSPercentageToken;
typedef struct YW_CSSDimensionToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_DIMENSION */
    char *unit;
    double value;
} YW_CSSDimensionToken;
typedef struct YW_CSSStringToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_STRING */
                           /* type = YW_CSS_TOKEN_URL */
                           /* type = YW_CSS_TOKEN_AT_KEYWORD */
                           /* type = YW_CSS_TOKEN_FUNC_KEYWORD */
                           /* type = YW_CSS_TOKEN_IDENT */
    char *value;
} YW_CSSStringToken;
typedef struct YW_CSSDelimToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_DELIM */
    YW_Char32 value;
} YW_CSSDelimToken;

typedef enum
{
    YW_HASH_ID,
    YW_HASH_UNRESTRICTED,
} YW_HashType;
typedef struct YW_CSSHashToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_HASH */
    char *value;
    YW_HashType type;
} YW_CSSHashToken;

typedef struct YW_AstFunctionToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_AST_FUNC */
    char *name;
    union YW_CSSToken *tokens;
    int tokens_len;
} YW_AstFunctionToken;
typedef struct YW_AstQualifiedRuleToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_AST_QUALIFIED_RULE */
    union YW_CSSToken *prelude_tokens, *body_tokens;
    int prelude_tokens_len, body_tokens_len;
} YW_AstQualifiedRuleToken;
typedef struct YW_AstAtRuleToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_AST_AT_RULE */
    union YW_CSSToken *prelude_tokens, *body_tokens;
    char *name;
    int prelude_tokens_len, body_tokens_len;
} YW_AstAtRuleToken;
typedef struct YW_AstDeclarationToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_AST_DECLARATION */
    char *name;
    union YW_CSSToken *value_tokens;
    int value_tokens_len;
    bool important;
} YW_AstDeclarationToken;

typedef enum
{
    YW_SIMPLE_BLOCK_CURLY,  /* { ... } */
    YW_SIMPLE_BLOCK_SQUARE, /* [ ... ] */
    YW_SIMPLE_BLOCK_PAREN,  /* ( ... ) */
} YW_SimpleBlockType;

typedef struct YW_AstSimpleBlockToken
{
    YW_TokenCommon common; /* type = YW_CSS_TOKEN_AST_SIMPLE_BLOCK */
    union YW_CSSToken *tokens;
    int tokens_len;
    YW_SimpleBlockType type;
} YW_AstSimpleBlockToken;

union YW_CSSToken {
    YW_TokenCommon common;

    YW_CSSNumberToken number_tk;
    YW_CSSPercentageToken percentage_tk;
    YW_CSSDimensionToken dimension_tk;
    YW_CSSStringToken string_tk, url_tk, at_keyword_tk, func_keyword_tk, ident_tk;
    YW_CSSDelimToken delim_tk;
    YW_CSSHashToken hash_tk;

    /***************************************************************************
     * High-level tokens
     **************************************************************************/

    YW_AstSimpleBlockToken ast_simple_block_tk;
    YW_AstFunctionToken ast_func_tk;
    YW_AstQualifiedRuleToken ast_qualified_rule_tk;
    YW_AstAtRuleToken ast_at_rule_tk;
    YW_AstDeclarationToken ast_declaration_tk;
};

bool yw_is_end_of_tokens(YW_CSSTokenStream const *ts);

/*
 * - expect~ functions return reference to an existing token.
 * - consume~ functions create a new token.
 */

/* Returns NULL on end of tokens */
YW_CSSToken const *yw_expect_any_token(YW_CSSTokenStream *ts);
/* Returns NULL if no such token was there */
YW_CSSToken const *yw_expect_token(YW_CSSTokenStream *ts, YW_TokenType type);
bool yw_expect_delim(YW_CSSTokenStream *ts, YW_Char32 d);
bool yw_expect_ident(YW_CSSTokenStream *ts, char const *i);
bool yw_expect_simple_block(YW_CSSTokenStream *inner_ts_out, YW_CSSTokenStream *ts, YW_SimpleBlockType type);
bool yw_expect_ast_func(YW_CSSTokenStream *inner_ts_out, YW_CSSTokenStream *ts, char const *f);
void yw_skip_whitespaces(YW_CSSTokenStream *ts);
bool yw_consume_declaration_value(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts);
bool yw_consume_any_value(YW_CSSToken **tokens_out, int *len_out, YW_CSSTokenStream *ts);

#define YW_CSS_NO_MAX_REPEATS -1

/*
 * This can be used to parse where a  syntax can be repeated separated
 * by comma.
 * If _max_repeats is YW_CSS_NO_MAX_REPEATS, repeat count is unlimited.
 *
 * https://www.w3.org/TR/-values-4/#mult-comma
 */
#define YW_CSS_PARSE_COMMA_SEPARATED_REPEATION(_type, _res_out, _len_out, _ts, _max_repeats, _parser) \
    do                                                                                                \
    {                                                                                                 \
        YW_CSSTokenStream *__ts = (_ts);                                                              \
        int __max_repeats = (_max_repeats);                                                           \
        _type *__res = NULL;                                                                          \
        int __res_len = 0;                                                                            \
        int __res_cap = 0;                                                                            \
        int __last_cursor_after_value = 0;                                                            \
        while (1)                                                                                     \
        {                                                                                             \
            _type __token;                                                                            \
            if (!(_parser)(&__token, __ts))                                                           \
            {                                                                                         \
                if (__res_len != 0)                                                                   \
                {                                                                                     \
                    __ts->cursor = __last_cursor_after_value;                                         \
                    break;                                                                            \
                }                                                                                     \
                break;                                                                                \
            }                                                                                         \
            YW_PUSH(_type, &__res_cap, &__res_len, &__res, __token);                                  \
            if (__max_repeats != YW_CSS_NO_MAX_REPEATS && __max_repeats <= __res_len)                 \
            {                                                                                         \
                break;                                                                                \
            }                                                                                         \
            yw_skip_whitespaces(__ts);                                                                \
            __last_cursor_after_value = __ts->cursor;                                                 \
            if (yw_expect_token(ts, YW_CSS_TOKEN_COMMA) == NULL)                                      \
            {                                                                                         \
                break;                                                                                \
            }                                                                                         \
            yw_skip_whitespaces(__ts);                                                                \
        }                                                                                             \
        YW_SHRINK_TO_FIT(_type, &__res_cap, __res_len, &__res);                                       \
        *(_res_out) = __res;                                                                          \
        *(_len_out) = __res_len;                                                                      \
    } while (0)

/*
 * This can be used to parse where a  syntax can be repeated multiple times.
 * If _max_repeats is YW_CSS_NO_MAX_REPEATS, repeat count is unlimited.
 *
 * https://www.w3.org/TR/-values-4/#mult-num-range
 */
#define YW_CSS_PARSE_REPEATION(_type, _res_out, _len_out, _ts, _max_repeats, _parser) \
    do                                                                                \
    {                                                                                 \
        YW_CSSTokenStream *__ts = (_ts);                                              \
        int __max_repeats = (_max_repeats);                                           \
        _type *__res = NULL;                                                          \
        int __res_len = 0;                                                            \
        int __res_cap = 0;                                                            \
        while (1)                                                                     \
        {                                                                             \
            _type __token;                                                            \
            if (!(_parser)(&__token, __ts))                                           \
            {                                                                         \
                break;                                                                \
            }                                                                         \
            YW_PUSH(_type, &__res_cap, &__res_len, &__res, __token);                  \
            if (__max_repeats != YW_CSS_NO_MAX_REPEATS && __max_repeats <= __res_len) \
            {                                                                         \
                break;                                                                \
            }                                                                         \
            yw_skip_whitespaces(__ts);                                                \
        }                                                                             \
        YW_SHRINK_TO_FIT(_type, &__res_cap, __res_len, &__res);                       \
        *(_res_out) = __res;                                                          \
        *(_len_out) = __res_len;                                                      \
    } while (0)

#endif /* #ifdef YW_CSS_TOKENS_H_ */
