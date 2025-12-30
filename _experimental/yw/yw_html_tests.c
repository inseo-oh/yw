/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_common.h"
#include "yw_html_tokens.h"
#include "yw_tests.h"
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

typedef struct YW_TokenList
{
    YW_HTMLToken *items;
    int len, cap;
} YW_TokenList;

static void yw_emit_callback(YW_HTMLToken *token, void *callback_data)
{
    YW_TokenList *token_list = (YW_TokenList *)callback_data;
    YW_HTMLToken new_tk;
    yw_html_token_clone(&new_tk, token);
    YW_LIST_PUSH(YW_HTMLToken, token_list, new_tk);
}

static void yw_tokenize(YW_HTMLToken **tokens_out, int *len_out, char const *str)
{
    YW_TokenList tl;
    memset(&tl, 0, sizeof(tl));
    yw_html_tokenize((const uint8_t *)str, strlen(str), yw_emit_callback, &tl);
    YW_SHRINK_TO_FIT(YW_HTMLToken, &tl.cap, tl.len, &tl.items);
    *tokens_out = tl.items;
    *len_out = tl.len;
}

void yw_test_html_parse_character_reference(YW_TestingContext *ctx)
{
    YW_HTMLToken *tokens;
    int len;

    yw_tokenize(&tokens, &len,
                "&#44032;"
                "&#xac01;"
                "&#xAC02;"
                "&nbsp;");

    YW_TEST_EXPECT(int, ctx, len, "%d", 5);
    if (len == 5)
    {
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[0].type, "%d", YW_HTML_CHAR_TOKEN);
        YW_TEST_EXPECT(YW_Char32, ctx, tokens[0].char_tk.chr, "%d", 0xac00);
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[1].type, "%d", YW_HTML_CHAR_TOKEN);
        YW_TEST_EXPECT(YW_Char32, ctx, tokens[1].char_tk.chr, "%d", 0xac01);
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[2].type, "%d", YW_HTML_CHAR_TOKEN);
        YW_TEST_EXPECT(YW_Char32, ctx, tokens[2].char_tk.chr, "%d", 0xac02);
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[3].type, "%d", YW_HTML_CHAR_TOKEN);
        YW_TEST_EXPECT(YW_Char32, ctx, tokens[3].char_tk.chr, "%d", 0x00a0);
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[4].type, "%d", YW_HTML_EOF_TOKEN);
    }
    for (int i = 0; i < len; i++)
    {
        yw_html_token_deinit(&tokens[i]);
    }
    free(tokens);
}

void yw_test_html_parse_comment(YW_TestingContext *ctx)
{
    YW_HTMLToken *tokens;
    int len;

    yw_tokenize(&tokens, &len, "<!--this is comment-->");
    YW_TEST_EXPECT(int, ctx, len, "%d", 2);
    if (len == 2)
    {
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[0].type, "%d", YW_HTML_COMMENT_TOKEN);
        YW_TEST_EXPECT_STR(ctx, tokens[0].comment_tk.data, "this is comment");
        YW_TEST_EXPECT(YW_HTMLTokenType, ctx, tokens[1].type, "%d", YW_HTML_EOF_TOKEN);
    }
    for (int i = 0; i < len; i++)
    {
        yw_html_token_deinit(&tokens[i]);
    }
    free(tokens);
}
