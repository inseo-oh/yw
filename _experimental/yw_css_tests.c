/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_css.h"
#include "yw_css_tokens.h"
#include "yw_tests.h"
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static bool yw_tokenize_str(YW_CSSTokenStream *out, YW_TestingContext *ctx, char const *str)
{
    if (!yw_css_tokenize(out, (uint8_t *)str, strlen(str)))
    {
        YW_FAILED_TEST(ctx, "failed to tokenize");
        return false;
    }
    return true;
}

static void yw_css_free_tokens(YW_CSSToken *tokens, int len)
{
    for (int i = 0; i < len; i++)
    {
        yw_css_token_deinit(&tokens[i]);
    }
    free(tokens);
}

/*******************************************************************************
 * CSS Values and Units
 ******************************************************************************/

void yw_test_css_parse_number(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "31;"
                         "3.2;"
                         ".33;"
                         "3.4e1;"
                         "350e-1;"
                         "3.6E1;"
                         "370e-1"))
    {
        return;
    }
    double num;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 31);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num * 10, "%d", 32);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num * 100, "%d", 33);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 34);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 35);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 36);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_number(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 37);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_length(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "0;"
                         "1em;"
                         "2ex;"
                         "3ch;"
                         "4rem;"
                         "5vw;"
                         "6vh;"
                         "7vmin;"
                         "8vmax;"
                         "9cm;"
                         "10mm;"
                         "11q;"
                         "12pc;"
                         "13pt;"
                         "14px"))
    {
        return;
    }
    YW_CSSLength len;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_NO_ALLOW_ZERO_SHORTHAND), "%d", false);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 0);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 1);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_EM);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 2);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_EX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 3);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_CH);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 4);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_REM);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 5);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_VW);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 6);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_VH);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 7);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_VMIN);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length(&len, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(int, ctx, len.value, "%d", 8);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, len.unit, "%d", YW_CSS_VMAX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_percentage(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "25%;"
                         "50%"))
    {
        return;
    }
    double num;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_percentage(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 25);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_percentage(&num, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, num, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_length_or_percentage(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "25px;"
                         "50%"))
    {
        return;
    }
    YW_CSSLengthOrPercentage num;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length_or_percentage(&num, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(bool, ctx, num.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, num.value.length.value, "%d", 25);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, num.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_length_or_percentage(&num, &ts, YW_ALLOW_ZERO_SHORTHAND), "%d", true);
    YW_TEST_EXPECT(bool, ctx, num.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, num.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Backgrounds and Borders
 ******************************************************************************/

void yw_test_css_parse_line_style(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "none;"
                         "hidden;"
                         "dotted;"
                         "dashed;"
                         "solid;"
                         "double;"
                         "groove;"
                         "ridge;"
                         "inset"))
    {
        return;
    }
    YW_CSSLineStyle style;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_NO_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_HIDDEN_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_DOTTED_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_DASHED_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_SOLID_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_DOUBLE_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_GROOVE_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_RIDGE_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSLineStyle, ctx, style, "%d", YW_CSS_INSET_LINE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_line_width(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "thin;"
                         "medium;"
                         "thick;"
                         "10em"))
    {
        return;
    }
    YW_CSSLength length;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_width(&length, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, length.value, "%d", 1);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_width(&length, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, length.value, "%d", 3);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_width(&length, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, length.value, "%d", 5);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_line_width(&length, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, length.value, "%d", 10);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, length.unit, "%d", YW_CSS_EM);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Box Model
 ******************************************************************************/

void yw_test_css_parse_margin(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "10px;"
                         "50%;"
                         "auto"))
    {
        return;
    }
    YW_CSSMargin margin;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_margin(&margin, &ts), "%d", true);
    YW_TEST_EXPECT(bool, ctx, margin.value.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, margin.value.value.length.value, "%d", 10);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, margin.value.value.length.unit, "%d", YW_CSS_PX);
    YW_TEST_EXPECT(bool, ctx, margin.is_auto, "%d", false);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_margin(&margin, &ts), "%d", true);
    YW_TEST_EXPECT(bool, ctx, margin.value.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, margin.value.value.percentage, "%d", 50);
    YW_TEST_EXPECT(bool, ctx, margin.is_auto, "%d", false);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_margin(&margin, &ts), "%d", true);
    YW_TEST_EXPECT(bool, ctx, margin.is_auto, "%d", true);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_padding(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "10px;"
                         "50%;"
                         "-1px"))
    {
        return;
    }
    YW_CSSLengthOrPercentage padding;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_padding(&padding, &ts), "%d", true);
    YW_TEST_EXPECT(bool, ctx, padding.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, padding.value.length.value, "%d", 10);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, padding.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_padding(&padding, &ts), "%d", true);
    YW_TEST_EXPECT(bool, ctx, padding.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, padding.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_padding(&padding, &ts), "%d", false);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Color
 ******************************************************************************/
void yw_test_css_parse_color(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         /* Hexadecimal ***************************************/
                         "#12345678;"
                         "#123456;"
                         "#1234;"
                         "#123;"
                         /* rgb/rgba - legacy syntax **************************/
                         "rgb( 100%, 100%, 100%, 0%);"
                         "rgba(100%, 100%, 100%, 0%);"
                         "rgb( 100%, 100%, 100%);"
                         "rgba(100%, 100%, 100%);"
                         "rgb( 12, 34, 56, 78);"
                         "rgba(90, 12, 34, 56);"
                         "rgb( 78, 90, 12);"
                         "rgba(34, 56, 78);"
                         /* rgb/rgba - modern syntax **************************/
                         "rgb( 100% 100% 100% / 0%);"
                         "rgba(100% 100% 100% / 0%);"
                         "rgb( 100% 100% 100%);"
                         "rgba(100% 100% 100%);"
                         "rgb( 12 34 56 / 78);"
                         "rgba(90 12 34 / 56);"
                         "rgb( 78 90 12);"
                         "rgba(34 56 78);"
                         /* Named colors, transparent, currentColor ***********/
                         "blue;"
                         "transparent;"
                         "currentColor;"))
    {
        return;
    }
    YW_CSSColor color;

    /* Hexadecimal ************************************************************/
    /* #12345678 */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x12345678);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* #123456 */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x123456ff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* #1234 */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x11223344);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* #123 */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x112233ff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb/rgba - legacy syntax ***********************************************/
    /* rgb( 100%, 100%, 100%, 0%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffff00);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(100%, 100%, 100%, 0%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffff00);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb( 100%, 100%, 100%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffffff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(100%, 100%, 100%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffffff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb( 12, 34, 56, 78) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x0c22384e);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(90, 12, 34, 56) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x5a0c2238);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb( 78, 90, 12) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x4e5a0cff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(34, 56, 78) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x22384eff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb/rgba - modern syntax ***********************************************/
    /* rgb( 100% 100% 100% / 0%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffff00);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(100% 100% 100% / 0%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffff00);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb( 100%, 100%, 100%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffffff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(100% 100% 100%) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0xffffffff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb( 12 34 56 / 78) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x0c22384e);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(90 12 34 / 56) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x5a0c2238);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgb( 78 90 12) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x4e5a0cff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* rgba(34 56 78) */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x22384eff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* Named colors, transparent, currentColor ********************************/

    /* blue */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x0000ffff);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* transparent */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_RGB_COLOR);
    YW_TEST_EXPECT(YW_CSSRgba, ctx, color.rgb.rgba, "#%08x", 0x00000000);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* currentColor */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_color(&color, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSColorType, ctx, color.type, "%d", YW_CSS_CURRENT_COLOR);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Display
 ******************************************************************************/

void yw_test_css_parse_display(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "none;"
                         "contents;"
                         "block;"
                         "flow-root;"
                         "inline;"
                         "inline-block;"
                         "run-in;"
                         /* TODO */
                         /*
                         "list-item;"
                         "inline list-item;"
                         */
                         "flex;"
                         "inline-flex;"
                         "grid;"
                         "inline-grid;"
                         "ruby;"
                         "block ruby;"
                         "table;"
                         "inline-table;"))
    {
        return;
    }

    YW_CSSDisplay display;

    /* none */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_NONE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* contents */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_CONTENTS);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* block */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_FLOW);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* flow-root */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_FLOW_ROOT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* inline */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_FLOW);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* inline-block */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_FLOW_ROOT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* run-in */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_RUN_IN | YW_CSS_DISPLAY_FLOW);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* Below are disabled as they are not implemented yet*/
    if (0)
    {
        /* list-item */
        YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
        YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_FLOW | YW_CSS_DISPLAY_LIST_ITEM);
        yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

        /* inline list-item */
        YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
        YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_FLOW | YW_CSS_DISPLAY_LIST_ITEM);
        yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);
    }

    /* flex */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_FLEX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* inline-flex */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_FLEX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* grid */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_GRID);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* inline-grid */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_GRID);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* ruby */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_RUBY);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* block ruby */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_RUBY);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* table */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_BLOCK | YW_CSS_DISPLAY_TABLE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* inline-table */
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_display(&display, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSDisplay, ctx, display, "%d", YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_TABLE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    /* TODO: test display-internal values *************************************/

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS2 9.5 Floats
 ******************************************************************************/

void yw_test_css_parse_float(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "none;"
                         "left;"
                         "right"))
    {
        return;
    }
    YW_CSSFloat flo;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_float(&flo, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFloat, ctx, flo, "%d", YW_CSS_NO_FLOAT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_float(&flo, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFloat, ctx, flo, "%d", YW_CSS_FLOAT_LEFT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_float(&flo, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFloat, ctx, flo, "%d", YW_CSS_FLOAT_RIGHT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Fonts
 ******************************************************************************/

void yw_test_css_parse_font_family(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "serif, sans-serif, cursive, fantasy, monospace,"
                         "space separated,"
                         "\"quoted string\""))
    {
        return;
    }

    YW_CSSFontFamilies families;
    memset(&families, 0, sizeof(families));
    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_family(&families, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, families.len, "%d", 7);
    if (families.len == 7)
    {
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[0].family, "%d", YW_CSS_SERIF);
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[1].family, "%d", YW_CSS_SANS_SERIF);
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[2].family, "%d", YW_CSS_CURSIVE);
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[3].family, "%d", YW_CSS_FANTASY);
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[4].family, "%d", YW_CSS_MONOSPACE);
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[5].family, "%d", YW_CSS_NON_GENERIC_FONT_FAMILY);
        if (families.items[5].family == YW_CSS_NON_GENERIC_FONT_FAMILY)
        {
            YW_TEST_EXPECT_STR(ctx, families.items[5].name, "space separated");
        }
        YW_TEST_EXPECT(YW_CSSGenericFontFamily, ctx, families.items[6].family, "%d", YW_CSS_NON_GENERIC_FONT_FAMILY);
        if (families.items[6].family == YW_CSS_NON_GENERIC_FONT_FAMILY)
        {
            YW_TEST_EXPECT_STR(ctx, families.items[6].name, "quoted string");
        }
    }
    for (int i = 0; i < families.len; i++)
    {
        free(families.items[i].name);
    }
    free(families.items);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_font_weight(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "normal;"
                         "bold;"
                         "600"))
    {
        return;
    }
    YW_CSSFontWeight weight;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_weight(&weight, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, weight, "%d", (YW_CSSFontWeight)400);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_weight(&weight, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, weight, "%d", (YW_CSSFontWeight)700);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_weight(&weight, &ts), "%d", true);
    YW_TEST_EXPECT(int, ctx, weight, "%d", (YW_CSSFontWeight)600);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_font_stretch(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "ultra-condensed;"
                         "extra-condensed;"
                         "condensed;"
                         "semi-condensed;"
                         "normal;"
                         "semi-expanded;"
                         "expanded;"
                         "extra-expanded;"
                         "ultra-expanded;"))
    {
        return;
    }
    YW_CSSFontStretch stretch;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_ULTRA_CONDENSED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_EXTRA_CONDENSED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_CONDENSED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_SEMI_CONDENSED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_NORMAL_FONT_STRETCH);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_SEMI_EXPANDED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_EXPANDED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_EXTRA_EXPANDED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_stretch(&stretch, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStretch, ctx, stretch, "%d", YW_CSS_ULTRA_EXPANDED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_font_style(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "normal;"
                         "italic;"
                         "oblique"))
    {
        return;
    }
    YW_CSSFontStyle style;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStyle, ctx, style, "%d", YW_CSS_NORMAL_FONT_STYLE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStyle, ctx, style, "%d", YW_CSS_ITALIC);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontStyle, ctx, style, "%d", YW_CSS_OBLIQUE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_font_size(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "xx-small;"
                         "x-small;"
                         "small;"
                         "medium;"
                         "large;"
                         "x-large;"
                         "xx-large;"
                         "larger;"
                         "smaller;"
                         "50%;"
                         "50px;"))
    {
        return;
    }
    YW_CSSFontSize size;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_XX_SMALL);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_X_SMALL);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_SMALL);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_MEDIUM_FONT_SIZE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_LARGE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_X_LARGE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_XX_LARGE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_LARGER);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_SMALLER);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_LENGTH_FONT_SIZE);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, size.size.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_font_size(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSFontSizeType, ctx, size.type, "%d", YW_CSS_LENGTH_FONT_SIZE);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, size.size.value.length.value, "%d", 50);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, size.size.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Selectors
 ******************************************************************************/

static YW_CSSSelector const *yw_complex_sel_to_compound(YW_CSSComplexSelector const *sel)
{
    if (sel == NULL || sel->type != YW_CSS_SELECTOR_COMPLEX)
    {
        printf("%s: not a complex selector\n", __func__);
        return NULL;
    }
    if (sel->base == NULL)
    {
        printf("%s: sel->base is NULL\n", __func__);
        return NULL;
    }
    return sel->base;
}

static YW_CSSSelector const *yw_compound_sel_to_type(YW_CSSSelector const *sel)
{
    if (sel == NULL || sel->type != YW_CSS_SELECTOR_COMPOUND)
    {
        printf("%s: not a compound selector\n", __func__);
        return NULL;
    }
    if (sel->compound_sel.type_sel == NULL)
    {
        printf("%s: sel->compound_sel.type_sel is NULL\n", __func__);
        return NULL;
    }
    return sel->compound_sel.type_sel;
}

static YW_CSSSelector const *yw_compound_sel_to_subclass(YW_CSSSelector const *sel)
{
    if (sel == NULL || sel->type != YW_CSS_SELECTOR_COMPOUND)
    {
        printf("%s: not a compound selector\n", __func__);
        return NULL;
    }
    if (sel->compound_sel.subclass_sels_len == 0)
    {
        printf("%s: no subclass selectors\n", __func__);
        return NULL;
    }
    return &sel->compound_sel.subclass_sels[0];
}

static YW_CSSSelector const *yw_compound_sel_to_pseudo_elem(YW_CSSSelector const *sel)
{
    if (sel == NULL || sel->type != YW_CSS_SELECTOR_COMPOUND)
    {
        printf("%s: not a compound selector\n", __func__);
        return NULL;
    }
    if (sel->compound_sel.pseudo_items_len == 0)
    {
        printf("%s: no pseudo items\n", __func__);
        return NULL;
    }
    return sel->compound_sel.pseudo_items[0].pseudo_elem_sel;
}

static YW_CSSSelector const *yw_complex_sel_to_type(YW_CSSComplexSelector const *sel)
{
    return yw_compound_sel_to_type(yw_complex_sel_to_compound(sel));
}
static YW_CSSSelector const *yw_complex_sel_to_subclass(YW_CSSComplexSelector const *sel)
{
    return yw_compound_sel_to_subclass(yw_complex_sel_to_compound(sel));
}
static YW_CSSSelector const *yw_complex_sel_to_pseudo_elem(YW_CSSComplexSelector const *sel)
{
    return yw_compound_sel_to_pseudo_elem(yw_complex_sel_to_compound(sel));
}

void yw_test_css_parse_selector(YW_TestingContext *ctx)
{
    YW_CSSSelector *selectors;
    int selectors_len;
    YW_CSSSelector const *sel_temp;
    YW_CSSComplexSelector const *complex_sel_temp;

    char const *source =
        /* Type selector */
        "div,"
        /* Pseudo class selector */
        ":link, :nth-child(1),"
        /* Pseudo element selector */
        "::before,"
        /* Subclass selector */
        "#id, .class, [attr], [attr=value], [attr=\"string\"], [attr~=value],"
        "[attr|=value], [attr^=value], [attr$=value], [attr*=value],"
        /* Compound selector */
        "div.class, div.class::before, div.class::before:link , div:link,"
        /* Complex selector */
        "#a #b, #a>#b, #a+#b, #a~#b, #a||#b";

    if (!yw_css_parse_selector(&selectors, &selectors_len, (uint8_t const *)source, strlen(source)))
    {
        YW_FAILED_TEST(ctx, "failed to parse selector");
        return;
    }
    YW_TEST_EXPECT(int, ctx, selectors_len, "%d", 23);
    if (selectors_len != 23)
    {
        return;
    }
    /* Type selector **********************************************************/
    /* div */
    sel_temp = yw_complex_sel_to_type(&selectors[0].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_TYPE);
        YW_TEST_EXPECT_STR(ctx, sel_temp->type_sel.name.ident, "div");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_type() failed");
    }
    /* Pseudo class selector **************************************************/
    /* :link */
    sel_temp = yw_complex_sel_to_subclass(&selectors[1].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->pseudo_class_sel.name, "link");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* :nth-child(1) */
    /* TODO: Test the inner contents as well, when we implement it*/
    sel_temp = yw_complex_sel_to_subclass(&selectors[2].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->pseudo_class_sel.name, "nth-child");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* Pseudo element selector ************************************************/
    /* ::before */
    sel_temp = yw_complex_sel_to_pseudo_elem(&selectors[3].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->pseudo_class_sel.name, "before");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_pseudo_elem() failed");
    }
    /* Subclass selector ******************************************************/
    /* #id */
    sel_temp = yw_complex_sel_to_subclass(&selectors[4].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "id");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* .class */
    sel_temp = yw_complex_sel_to_subclass(&selectors[5].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->class_sel.class_name, "class");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[6].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_NO_VALUE_MATCH);
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr=value] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[7].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "value");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr="string"] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[8].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "string");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr~=value] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[9].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_TILDE_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "value");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr|=value] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[10].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_BAR_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "value");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr^=value] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[11].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_CARET_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "value");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr$=value] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[12].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_DOLLAR_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "value");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* [attr*=value] */
    sel_temp = yw_complex_sel_to_subclass(&selectors[13].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ATTR);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_name.ident, "attr");
        YW_TEST_EXPECT(YW_CSSValueMatchType, ctx, sel_temp->attr_sel.value_match_type, "%d", YW_CSS_VALUE_ASTERISK_EQUALS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->attr_sel.attr_value, "value");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }
    /* Compound selector ******************************************************/
    /* div.class */
    sel_temp = yw_complex_sel_to_compound(&selectors[14].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.subclass_sels_len, "%d", 1);
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items_len, "%d", 0);
        YW_TEST_EXPECT(bool, ctx, sel_temp->compound_sel.type_sel != NULL, "%d", true);
        if (sel_temp->compound_sel.type_sel != NULL)
        {
            YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.type_sel->type, "%d", YW_CSS_SELECTOR_TYPE);
            YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.type_sel->type_sel.name.ident, "div");
        }
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->compound_sel.subclass_sels[0].type, "%d", YW_CSS_SELECTOR_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.subclass_sels[0].class_sel.class_name, "class");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_compound() failed");
    }
    /* div.class::before */
    sel_temp = yw_complex_sel_to_compound(&selectors[15].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.subclass_sels_len, "%d", 1);
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items_len, "%d", 1);
        YW_TEST_EXPECT(bool, ctx, sel_temp->compound_sel.type_sel != NULL, "%d", true);
        if (sel_temp->compound_sel.type_sel != NULL)
        {
            YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.type_sel->type, "%d", YW_CSS_SELECTOR_TYPE);
            YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.type_sel->type_sel.name.ident, "div");
        }
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->compound_sel.subclass_sels[0].type, "%d", YW_CSS_SELECTOR_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.subclass_sels[0].class_sel.class_name, "class");
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items[0].class_sels_len, "%d", 0);
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items[0].pseudo_elem_sel->type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.pseudo_items[0].pseudo_elem_sel->pseudo_class_sel.name, "before");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_compound() failed");
    }
    /* div.class::before:link */
    sel_temp = yw_complex_sel_to_compound(&selectors[16].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.subclass_sels_len, "%d", 1);
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items_len, "%d", 1);
        YW_TEST_EXPECT(bool, ctx, sel_temp->compound_sel.type_sel != NULL, "%d", true);
        if (sel_temp->compound_sel.type_sel != NULL)
        {
            YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.type_sel->type, "%d", YW_CSS_SELECTOR_TYPE);
            YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.type_sel->type_sel.name.ident, "div");
        }
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->compound_sel.subclass_sels[0].type, "%d", YW_CSS_SELECTOR_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.subclass_sels[0].class_sel.class_name, "class");

        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items[0].class_sels_len, "%d", 1);
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->compound_sel.pseudo_items[0].pseudo_elem_sel->type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.pseudo_items[0].pseudo_elem_sel->pseudo_class_sel.name, "before");
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->compound_sel.pseudo_items[0].class_sels[0].type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.pseudo_items[0].class_sels[0].pseudo_class_sel.name, "link");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_compound() failed");
    }
    /* div:link */
    sel_temp = yw_complex_sel_to_compound(&selectors[17].complex_sel);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.subclass_sels_len, "%d", 1);
        YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.pseudo_items_len, "%d", 0);
        YW_TEST_EXPECT(bool, ctx, sel_temp->compound_sel.type_sel != NULL, "%d", true);
        if (sel_temp->compound_sel.type_sel != NULL)
        {
            YW_TEST_EXPECT(int, ctx, sel_temp->compound_sel.type_sel->type, "%d", YW_CSS_SELECTOR_TYPE);
            YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.type_sel->type_sel.name.ident, "div");
        }
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->compound_sel.subclass_sels[0].type, "%d", YW_CSS_SELECTOR_PSEUDO_CLASS);
        YW_TEST_EXPECT_STR(ctx, sel_temp->compound_sel.subclass_sels[0].class_sel.class_name, "link");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_compound() failed");
    }
    /* Complex selector *******************************************************/
    /* #a #b */
    complex_sel_temp = &selectors[18].complex_sel;
    YW_TEST_EXPECT(int, ctx, complex_sel_temp->type, "%d", YW_CSS_SELECTOR_COMPLEX);
    sel_temp = yw_complex_sel_to_subclass(complex_sel_temp);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "a");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_complex_sel_to_subclass() failed");
    }

    YW_TEST_EXPECT(int, ctx, complex_sel_temp->rests_len, "%d", 1);
    YW_TEST_EXPECT(YW_CSSCombinator, ctx, complex_sel_temp->rests[0].combinator, "%d", YW_CSS_CHILD_COMBINATOR);

    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->rests[0].selector);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "b");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }
    /* #a>#b */
    complex_sel_temp = &selectors[19].complex_sel;
    YW_TEST_EXPECT(YW_CSSSelectorType, ctx, complex_sel_temp->type, "%d", YW_CSS_SELECTOR_COMPLEX);
    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->base);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "a");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }

    YW_TEST_EXPECT(int, ctx, complex_sel_temp->rests_len, "%d", 1);
    YW_TEST_EXPECT(YW_CSSCombinator, ctx, complex_sel_temp->rests[0].combinator, "%d", YW_CSS_DIRECT_CHILD_COMBINATOR);

    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->rests[0].selector);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "b");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }
    /* #a+#b */
    complex_sel_temp = &selectors[20].complex_sel;
    YW_TEST_EXPECT(YW_CSSSelectorType, ctx, complex_sel_temp->type, "%d", YW_CSS_SELECTOR_COMPLEX);
    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->base);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "a");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }

    YW_TEST_EXPECT(int, ctx, complex_sel_temp->rests_len, "%d", 1);
    YW_TEST_EXPECT(YW_CSSCombinator, ctx, complex_sel_temp->rests[0].combinator, "%d", YW_CSS_PLUS_COMBINATOR);

    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->rests[0].selector);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "b");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }
    /* #a~#b */
    complex_sel_temp = &selectors[21].complex_sel;
    YW_TEST_EXPECT(YW_CSSSelectorType, ctx, complex_sel_temp->type, "%d", YW_CSS_SELECTOR_COMPLEX);
    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->base);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "a");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }

    YW_TEST_EXPECT(int, ctx, complex_sel_temp->rests_len, "%d", 1);
    YW_TEST_EXPECT(YW_CSSCombinator, ctx, complex_sel_temp->rests[0].combinator, "%d", YW_CSS_TILDE_COMBINATOR);

    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->rests[0].selector);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "b");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }
    /* #a||#b */
    complex_sel_temp = &selectors[22].complex_sel;
    YW_TEST_EXPECT(YW_CSSSelectorType, ctx, complex_sel_temp->type, "%d", YW_CSS_SELECTOR_COMPLEX);
    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->base);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "a");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }

    YW_TEST_EXPECT(int, ctx, complex_sel_temp->rests_len, "%d", 1);
    YW_TEST_EXPECT(YW_CSSCombinator, ctx, complex_sel_temp->rests[0].combinator, "%d", YW_CSS_TWO_BARS_COMBINATOR);

    sel_temp = yw_compound_sel_to_subclass(complex_sel_temp->rests[0].selector);
    if (sel_temp != NULL)
    {
        YW_TEST_EXPECT(YW_CSSSelectorType, ctx, sel_temp->type, "%d", YW_CSS_SELECTOR_ID);
        YW_TEST_EXPECT_STR(ctx, sel_temp->id_sel.id, "b");
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_compound_sel_to_subclass() failed");
    }

    for (int i = 0; i < selectors_len; i++)
    {
        yw_css_selector_deinit(&selectors[i]);
    }
    free(selectors);
}

/*******************************************************************************
 * CSS Sizing
 ******************************************************************************/
void yw_test_css_parse_size_or_auto(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "auto;"
                         "min-content;"
                         "max-content;"
                         "fit-content(50%);"
                         "fit-content(50px);"
                         "50%;"
                         "50px;"
                         "none"))
    {
        return;
    }
    YW_CSSSize size;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_AUTO_SIZE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MIN_CONTENT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MAX_CONTENT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_FIT_CONTENT);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, size.size.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_FIT_CONTENT);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, size.size.value.length.value, "%d", 50);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, size.size.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MANUAL_SIZE);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, size.size.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MANUAL_SIZE);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, size.size.value.length.value, "%d", 50);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, size.size.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_auto(&size, &ts), "%d", false);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_size_or_none(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "none;"
                         "min-content;"
                         "max-content;"
                         "fit-content(50%);"
                         "fit-content(50px);"
                         "50%;"
                         "50px;"
                         "auto"))
    {
        return;
    }
    YW_CSSSize size;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_NO_SIZE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MIN_CONTENT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MAX_CONTENT);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_FIT_CONTENT);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, size.size.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_FIT_CONTENT);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, size.size.value.length.value, "%d", 50);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, size.size.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MANUAL_SIZE);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", true);
    YW_TEST_EXPECT(int, ctx, size.size.value.percentage, "%d", 50);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSSizeType, ctx, size.type, "%d", YW_CSS_MANUAL_SIZE);
    YW_TEST_EXPECT(bool, ctx, size.size.is_percentage, "%d", false);
    YW_TEST_EXPECT(int, ctx, size.size.value.length.value, "%d", 50);
    YW_TEST_EXPECT(YW_CSSLengthUnit, ctx, size.size.value.length.unit, "%d", YW_CSS_PX);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_size_or_none(&size, &ts), "%d", false);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Text
 ******************************************************************************/

void yw_test_css_parse_text_transform(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "capitalize;"
                         "capitalize full-width;"
                         "capitalize full-width full-size-kana;"
                         "uppercase;"
                         "lowercase;"
                         "full-width full-size-kana;"))
    {
        return;
    }
    YW_CSSTextTransform trans;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_transform(&trans, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextTransform, ctx, trans, "%d", YW_CSS_TEXT_TRANSFORM_CAPITALIZE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_transform(&trans, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextTransform, ctx, trans, "%d", YW_CSS_TEXT_TRANSFORM_CAPITALIZE | YW_CSS_TEXT_TRANSFORM_FULL_WIDTH);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_transform(&trans, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextTransform, ctx, trans, "%d", YW_CSS_TEXT_TRANSFORM_CAPITALIZE | YW_CSS_TEXT_TRANSFORM_FULL_WIDTH | YW_CSS_TEXT_TRANSFORM_FULL_SIZE_KANA);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_transform(&trans, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextTransform, ctx, trans, "%d", YW_CSS_TEXT_TRANSFORM_UPPERCASE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_transform(&trans, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextTransform, ctx, trans, "%d", YW_CSS_TEXT_TRANSFORM_LOWERCASE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_transform(&trans, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextTransform, ctx, trans, "%d", YW_CSS_TEXT_TRANSFORM_ORIGINAL_CAPS | YW_CSS_TEXT_TRANSFORM_FULL_WIDTH | YW_CSS_TEXT_TRANSFORM_FULL_SIZE_KANA);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

/*******************************************************************************
 * CSS Text Decoration
 ******************************************************************************/

void yw_test_css_parse_text_decoration_line(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "none;"
                         " underline overline line-through blink;"))
    {
        return;
    }
    YW_CSSTextDecorationLine lines;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_line(&lines, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationLine, ctx, lines, "%d", 0);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_line(&lines, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationLine, ctx, lines, "%d", YW_CSS_TEXT_DECORATION_UNDERLINE | YW_CSS_TEXT_DECORATION_OVERLINE | YW_CSS_TEXT_DECORATION_LINE_THROUGH | YW_CSS_TEXT_DECORATION_BLINK);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}

void yw_test_css_parse_text_decoration_style(YW_TestingContext *ctx)
{
    YW_CSSTokenStream ts;
    if (!yw_tokenize_str(&ts, ctx,
                         "solid;"
                         "double;"
                         "dotted;"
                         "dashed;"
                         "wavy"))
    {
        return;
    }
    YW_CSSTextDecorationStyle style;

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationStyle, ctx, style, "%d", YW_CSS_TEXT_DECORATION_SOLID);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationStyle, ctx, style, "%d", YW_CSS_TEXT_DECORATION_DOUBLE);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationStyle, ctx, style, "%d", YW_CSS_TEXT_DECORATION_DOTTED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationStyle, ctx, style, "%d", YW_CSS_TEXT_DECORATION_DASHED);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    YW_TEST_EXPECT(bool, ctx, yw_css_parse_text_decoration_style(&style, &ts), "%d", true);
    YW_TEST_EXPECT(YW_CSSTextDecorationStyle, ctx, style, "%d", YW_CSS_TEXT_DECORATION_WAVY);
    yw_expect_token(&ts, YW_CSS_TOKEN_SEMICOLON);

    yw_css_free_tokens(ts.tokens, ts.tokens_len);
}
