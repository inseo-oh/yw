/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_common.h"
#include "yw_css.h"
#include "yw_css_tokens.h"
#include <limits.h>
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*******************************************************************************
 *
 * CSS Values and Units
 *
 * https://www.w3.org/TR/css-values-3/
 *
 ******************************************************************************/

bool yw_css_parse_number(double *out, YW_CSSTokenStream *ts)
{
    YW_CSSToken const *tk = yw_expect_token(ts, YW_CSS_TOKEN_NUMBER);
    if (tk == NULL)
    {
        return false;
    }
    *out = tk->number_tk.value;
    return true;
}

bool yw_css_parse_length(YW_CSSLength *out, YW_CSSTokenStream *ts, YW_AllowZeroShorthand allow_zero_shorthand)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *dim_tk = yw_expect_token(ts, YW_CSS_TOKEN_DIMENSION);
    double len_value;
    YW_CSSLengthUnit len_unit;

    if (dim_tk == NULL)
    {
        if (allow_zero_shorthand == YW_ALLOW_ZERO_SHORTHAND)
        {
            YW_CSSToken const *num_tk = yw_expect_token(ts, YW_CSS_TOKEN_NUMBER);
            if (num_tk == NULL || num_tk->number_tk.value != 0)
            {
                ts->cursor = old_cursor;
            }
            else
            {
                out->value = 0;
                out->unit = YW_CSS_PX;
                return true;
            }
        }
        goto fail;
    }
    len_value = dim_tk->dimension_tk.value;
    if (strcmp(dim_tk->dimension_tk.unit, "em") == 0)
    {
        len_unit = YW_CSS_EM;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "ex") == 0)
    {
        len_unit = YW_CSS_EX;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "ch") == 0)
    {
        len_unit = YW_CSS_CH;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "rem") == 0)
    {
        len_unit = YW_CSS_REM;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "vw") == 0)
    {
        len_unit = YW_CSS_VW;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "vh") == 0)
    {
        len_unit = YW_CSS_VH;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "vmin") == 0)
    {
        len_unit = YW_CSS_VMIN;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "vmax") == 0)
    {
        len_unit = YW_CSS_VMAX;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "cm") == 0)
    {
        len_unit = YW_CSS_CM;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "mm") == 0)
    {
        len_unit = YW_CSS_MM;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "q") == 0)
    {
        len_unit = YW_CSS_Q;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "pc") == 0)
    {
        len_unit = YW_CSS_PC;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "pt") == 0)
    {
        len_unit = YW_CSS_PT;
    }
    else if (strcmp(dim_tk->dimension_tk.unit, "px") == 0)
    {
        len_unit = YW_CSS_PX;
    }
    else
    {
        /* Bad unit */
        goto fail;
    }
    out->unit = len_unit;
    out->value = len_value;
    return true;
fail:
    return false;
}

bool yw_css_parse_percentage(double *out, YW_CSSTokenStream *ts)
{
    YW_CSSToken const *tk = yw_expect_token(ts, YW_CSS_TOKEN_PERCENTAGE);
    if (tk == NULL)
    {
        return false;
    }
    *out = tk->percentage_tk.value;
    return true;
}

bool yw_css_parse_length_or_percentage(YW_CSSLengthOrPercentage *out, YW_CSSTokenStream *ts, YW_AllowZeroShorthand allow_zero_shorthand)
{
    if (yw_css_parse_length(&out->value.length, ts, allow_zero_shorthand))
    {
        out->is_percentage = false;
        return true;
    }
    if (yw_css_parse_percentage(&out->value.percentage, ts))
    {
        out->is_percentage = true;
        return true;
    }
    return false;
}

/*******************************************************************************
 *
 * CSS Backgrounds and Borders
 *
 * https://www.w3.org/TR/css-backgrounds-3/
 *
 ******************************************************************************/

bool yw_css_parse_line_style(YW_CSSLineStyle *out, YW_CSSTokenStream *ts)
{
    /*
     * https://www.w3.org/TR/css-backgrounds-3/#typedef-line-style
     */
    if (yw_expect_ident(ts, "none"))
    {
        *out = YW_CSS_NO_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "hidden"))
    {
        *out = YW_CSS_HIDDEN_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "dotted"))
    {
        *out = YW_CSS_DOTTED_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "dashed"))
    {
        *out = YW_CSS_DASHED_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "solid"))
    {
        *out = YW_CSS_SOLID_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "double"))
    {
        *out = YW_CSS_DOUBLE_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "groove"))
    {
        *out = YW_CSS_GROOVE_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "ridge"))
    {
        *out = YW_CSS_RIDGE_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "inset"))
    {
        *out = YW_CSS_INSET_LINE;
        return true;
    }
    if (yw_expect_ident(ts, "outset"))
    {
        *out = YW_CSS_OUTSET_LINE;
        return true;
    }
    return false;
}

bool yw_css_parse_line_width(YW_CSSLength *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "thin"))
    {
        out->value = YW_CSS_LINE_WIDTH_THIN;
        out->unit = YW_CSS_PX;
        return true;
    }
    if (yw_expect_ident(ts, "medium"))
    {
        out->value = YW_CSS_LINE_WIDTH_MEDIUM;
        out->unit = YW_CSS_PX;
        return true;
    }
    if (yw_expect_ident(ts, "thick"))
    {
        out->value = YW_CSS_LINE_WIDTH_THICK;
        out->unit = YW_CSS_PX;
        return true;
    }
    if (yw_css_parse_length(out, ts, YW_ALLOW_ZERO_SHORTHAND))
    {
        return true;
    }
    return false;
}

/*******************************************************************************
 *
 * CSS Box Model
 *
 * https://www.w3.org/TR/css-box-3/
 *
 ******************************************************************************/

bool yw_css_parse_margin(YW_CSSMargin *out, YW_CSSTokenStream *ts)
{
    if (yw_css_parse_length_or_percentage(&out->value, ts, YW_ALLOW_ZERO_SHORTHAND))
    {
        out->is_auto = false;
        return true;
    }
    else if (yw_expect_ident(ts, "auto"))
    {
        out->is_auto = true;
        return true;
    }
    return false;
}
bool yw_css_parse_padding(YW_CSSLengthOrPercentage *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    if (!yw_css_parse_length_or_percentage(out, ts, YW_ALLOW_ZERO_SHORTHAND))
    {
        goto fail;
    }
    if (!out->is_percentage && out->value.length.value < 0)
    {
        /* Out of range */
        goto fail;
    }
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

/*******************************************************************************
 *
 * CSS Color
 *
 * https://www.w3.org/TR/css-color-4/
 *
 ******************************************************************************/

static bool yw_parse_alpha(int *out, YW_CSSTokenStream *ts)
{
    double num;
    if (yw_css_parse_number(&num, ts))
    {
        *out = YW_CLAMP(num, 0, 255);
    }
    else if (yw_css_parse_percentage(&num, ts))
    {
        double per = YW_CLAMP(num, 0, 100);
        *out = (per / 100) * 255;
    }
    else
    {
        return false;
    }
    return true;
}

bool yw_css_parse_color(YW_CSSColor *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *temp_tk;
    YW_CSSTokenStream inner_ts;

    /* Try hex notation *******************************************************/
    temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_HASH);
    if (temp_tk != NULL)
    {
        char const *value = temp_tk->hash_tk.value;
        /*
         * https://www.w3.org/TR/css-color-4/#hex-notation
         */
        int chars_len = yw_utf8_strlen(value);
        if (8 < chars_len)
        {
            /* Bad length */
            goto fail;
        }
        /* Check if we have non-hex digit */
        {
            char const *cur = value;
            while (1)
            {
                YW_Char32 c = yw_utf8_next_char(&cur);
                if (c == 0)
                {
                    break;
                }
                else if (!yw_is_ascii_hex_digit(c))
                {
                    goto fail;
                }
            }
        }
        /* At this point we assume the string is entirely ASCII */
        switch (chars_len)
        {
        case 3: {
            /* #rgb */
            char red_str[] = {value[0], value[0], '\0'};
            char green_str[] = {value[1], value[1], '\0'};
            char blue_str[] = {value[2], value[2], '\0'};
            int red = strtol(red_str, NULL, 16);
            int green = strtol(green_str, NULL, 16);
            int blue = strtol(blue_str, NULL, 16);
            out->type = YW_CSS_RGB_COLOR;
            out->rgb.rgba = YW_CSS_RGBA(red, green, blue, 255);
            break;
        }
        case 4: {
            /* #rgba */
            char red_str[] = {value[0], value[0], '\0'};
            char green_str[] = {value[1], value[1], '\0'};
            char blue_str[] = {value[2], value[2], '\0'};
            char alpha_str[] = {value[3], value[3], '\0'};
            int red = strtol(red_str, NULL, 16);
            int green = strtol(green_str, NULL, 16);
            int blue = strtol(blue_str, NULL, 16);
            int alpha = strtol(alpha_str, NULL, 16);
            out->type = YW_CSS_RGB_COLOR;
            out->rgb.rgba = YW_CSS_RGBA(red, green, blue, alpha);
            break;
        }
        case 6: {
            /* #rrggbb */
            char red_str[] = {value[0], value[1], '\0'};
            char green_str[] = {value[2], value[3], '\0'};
            char blue_str[] = {value[4], value[5], '\0'};
            int red = strtol(red_str, NULL, 16);
            int green = strtol(green_str, NULL, 16);
            int blue = strtol(blue_str, NULL, 16);
            out->type = YW_CSS_RGB_COLOR;
            out->rgb.rgba = YW_CSS_RGBA(red, green, blue, 255);
            break;
        }
        case 8: {
            /* #rgba */
            char red_str[] = {value[0], value[1], '\0'};
            char green_str[] = {value[2], value[3], '\0'};
            char blue_str[] = {value[4], value[5], '\0'};
            char alpha_str[] = {value[6], value[7], '\0'};
            int red = strtol(red_str, NULL, 16);
            int green = strtol(green_str, NULL, 16);
            int blue = strtol(blue_str, NULL, 16);
            int alpha = strtol(alpha_str, NULL, 16);
            out->type = YW_CSS_RGB_COLOR;
            out->rgb.rgba = YW_CSS_RGBA(red, green, blue, alpha);
            break;
        }
        default:
            goto fail;
        }
        return true;
    }
    /* Try rgb()/rgba() function **********************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "rgb") || yw_expect_ast_func(&inner_ts, ts, "rgba"))
    {
        YW_CSSTokenStream *ts = &inner_ts;
        /*
         * https://www.w3.org/TR/css-color-4/#funcdef-rgb
         */
        int r, g, b, a;
        a = 255;

        /***********************************************************************
         * Try legacy syntax first
         * https://www.w3.org/TR/css-color-4/#typedef-legacy-rgb-syntax
         **********************************************************************/

        /* rgb(<  >r  ,  g  ,  b  ) *******************************************/
        /* rgb(<  >r  ,  g  ,  b  ,  a  ) *************************************/
        yw_skip_whitespaces(ts);
        int cursor_after_whitespaces = ts->cursor;
        /* rgb(  <r  ,  g  ,  b>  ) *******************************************/
        /* rgb(  <r  ,  g  ,  b>  ,  a  ) *************************************/
        double *values;
        int values_len;
        YW_CSS_PARSE_COMMA_SEPARATED_REPEATION(double, &values, &values_len, ts, 3, yw_css_parse_percentage);
        if (values_len == 3)
        {
            /* Percentage value */
            double r_per = YW_CLAMP(values[0], 0, 100);
            double g_per = YW_CLAMP(values[1], 0, 100);
            double b_per = YW_CLAMP(values[2], 0, 100);
            r = (r_per / 100) * 255;
            g = (g_per / 100) * 255;
            b = (b_per / 100) * 255;
            free(values);
        }
        else
        {
            free(values);
            ts->cursor = cursor_after_whitespaces;
            YW_CSS_PARSE_COMMA_SEPARATED_REPEATION(double, &values, &values_len, ts, 3, yw_css_parse_number);
            if (values_len == 3)
            {
                r = YW_CLAMP(values[0], 0, 255);
                g = YW_CLAMP(values[1], 0, 255);
                b = YW_CLAMP(values[2], 0, 255);
                free(values);
            }
            else
            {
                free(values);
                goto modern_syntax;
            }
        }
        /* rgb(  r  ,  g  ,  b<  >) *******************************************/
        /* rgb(  r  ,  g  ,  b<  >,  a  ) *************************************/
        yw_skip_whitespaces(ts);
        /* rgb(  r  ,  g  ,  b  <,>  a  ) *************************************/
        if (yw_expect_token(ts, YW_CSS_TOKEN_COMMA) != NULL)
        {
            /* rgb(  r  ,  g  ,  b  ,<  >a  ) *********************************/
            yw_skip_whitespaces(ts);
            /* rgb(  r  ,  g  ,  b  ,  <a>  ) *********************************/
            if (!yw_parse_alpha(&a, ts))
            {
                goto fail;
            }
            /* rgb(  r  ,  g  ,  b  ,  a<  >) *********************************/
            yw_skip_whitespaces(ts);
        }
        if (!yw_is_end_of_tokens(ts))
        {
            /* Extra junk at the end */
            goto fail;
        }
        out->type = YW_CSS_RGB_COLOR;
        out->rgb.rgba = YW_CSS_RGBA(r, g, b, a);
        return true;
    modern_syntax:
        ts->cursor = cursor_after_whitespaces;

        /***********************************************************************
         * Try modern syntax
         * https://www.w3.org/TR/css-color-4/#typedef-modern-rgb-syntax
         **********************************************************************/

        /* rgb(<  >r  g  b  ) *************************************************/
        /* rgb(<  >r  g  b  /  a  ) *******************************************/
        yw_skip_whitespaces(ts);
        /* rgb(  <r  g  b  >) *************************************************/
        /* rgb(  <r  g  b  >/  a  ) *******************************************/
        int chans[3];
        for (int i = 0; i < 3; i++)
        {
            /* rgb(  <r>  <g>  <b>  ) *****************************************/
            /* rgb(  <r>  <g>  <b>  /  a  ) ***********************************/
            double num;
            int v;
            if (yw_css_parse_number(&num, ts))
            {
                v = YW_CLAMP(num, 0, 255);
            }
            else if (yw_css_parse_percentage(&num, ts))
            {
                double per = YW_CLAMP(num, 0, 255);
                v = (per / 100) * 255;
            }
            else if (yw_expect_ident(ts, "none"))
            {
                YW_TODO();
            }
            else
            {
                goto fail;
            }
            chans[i] = v;
            /* rgb(  r<  >g<  >b<  >) *****************************************/
            /* rgb(  r<  >g<  >b<  >/  a  ) ***********************************/
            yw_skip_whitespaces(ts);
        }
        r = chans[0];
        g = chans[1];
        b = chans[2];
        a = 255;
        /* rgb(  r  g  b  </>  a  ) *******************************************/
        if (yw_expect_delim(ts, '/'))
        {
            /* rgb(  r  g  b  /<  >a  ) ***************************************/
            yw_skip_whitespaces(ts);
            /* rgb(  r  g  b  /  <a>  ) ***************************************/
            if (!yw_parse_alpha(&a, ts))
            {
                goto fail;
            }
            /* rgb(  r  g  b  /  a<  >) ***************************************/
            yw_skip_whitespaces(ts);
        }
        out->type = YW_CSS_RGB_COLOR;
        out->rgb.rgba = YW_CSS_RGBA(r, g, b, a);
        return true;
    }
    /* Try hsl()/hsla() function **********************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "hsl") || yw_expect_ast_func(&inner_ts, ts, "hsla"))
    {
        YW_TODO();
    }
    /* Try hwb() function *****************************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "hwb"))
    {
        YW_TODO();
    }
    /* Try lab() function *****************************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "lab"))
    {
        YW_TODO();
    }
    /* Try lch() function *****************************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "lch"))
    {
        YW_TODO();
    }
    /* Try oklab() function ***************************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "oklab"))
    {
        YW_TODO();
    }
    /* Try oklch() function ***************************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "oklch"))
    {
        YW_TODO();
    }
    /* Try color() function ***************************************************/
    if (yw_expect_ast_func(&inner_ts, ts, "color"))
    {
        YW_TODO();
    }
    /* Try named color, transparent, currentColor, system colors **************/
    temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (temp_tk != NULL)
    {
        /* Named colors *******************************************************/
        char const *ident = temp_tk->ident_tk.value;
        YW_CSSRgba rgba = yw_css_color_from_name(ident);
        if (rgba != 0)
        {
            out->type = YW_CSS_RGB_COLOR;
            out->rgb.rgba = rgba;
            return true;
        }
        /* Transparent ********************************************************/
        if (strcmp(ident, "transparent") == 0)
        {
            out->type = YW_CSS_RGB_COLOR;
            out->rgb.rgba = rgba;
            return true;
        }
        /* currentColor *******************************************************/
        if (strcmp(ident, "currentColor") == 0)
        {
            out->type = YW_CSS_CURRENT_COLOR;
            return true;
        }
        /* System colors ******************************************************/
        YW_TODO();
    }
fail:
    ts->cursor = old_cursor;
    return false;
}

/*******************************************************************************
 *
 * CSS Display
 *
 * https://www.w3.org/TR/css-display-3/
 *
 ******************************************************************************/

static bool yw_parse_display_outside(YW_CSSDisplay *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "block"))
    {
        *out = YW_CSS_DISPLAY_BLOCK;
        return true;
    }
    else if (yw_expect_ident(ts, "inline"))
    {
        *out = YW_CSS_DISPLAY_INLINE;
        return true;
    }
    else if (yw_expect_ident(ts, "run-in"))
    {
        *out = YW_CSS_DISPLAY_RUN_IN;
        return true;
    }
    return false;
}
static bool yw_parse_display_inside(YW_CSSDisplay *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "flow"))
    {
        *out = YW_CSS_DISPLAY_FLOW;
        return true;
    }
    else if (yw_expect_ident(ts, "flow-root"))
    {
        *out = YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table"))
    {
        *out = YW_CSS_DISPLAY_TABLE;
        return true;
    }
    else if (yw_expect_ident(ts, "flex"))
    {
        *out = YW_CSS_DISPLAY_FLEX;
        return true;
    }
    else if (yw_expect_ident(ts, "grid"))
    {
        *out = YW_CSS_DISPLAY_GRID;
        return true;
    }
    else if (yw_expect_ident(ts, "ruby"))
    {
        *out = YW_CSS_DISPLAY_RUBY;
        return true;
    }
    return false;
}
bool yw_css_parse_display(YW_CSSDisplay *out, YW_CSSTokenStream *ts)
{
    YW_CSSDisplay res;
    /* Try legacy keyword first ***********************************************/
    /*
     * https://www.w3.org/TR/css-display-3/#typedef-display-legacy
     */
    if (yw_expect_ident(ts, "inline-block"))
    {
        *out = YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "inline-table"))
    {
        *out = YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_TABLE;
        return true;
    }
    else if (yw_expect_ident(ts, "inline-flex"))
    {
        *out = YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_FLEX;
        return true;
    }
    else if (yw_expect_ident(ts, "inline-grid"))
    {
        *out = YW_CSS_DISPLAY_INLINE | YW_CSS_DISPLAY_GRID;
        return true;
    }
    /* Try <display-outside> <display-inside> *********************************/
    bool got_outer_mode = false, got_inner_mode = false;
    res = 0;
    while (!got_outer_mode || !got_inner_mode)
    {
        bool got_something = false;

        if (!got_outer_mode)
        {
            yw_skip_whitespaces(ts);
            YW_CSSDisplay temp;
            if (yw_parse_display_outside(&temp, ts))
            {
                got_something = true;
                got_outer_mode = true;
                res |= temp;
            }
        }
        if (!got_inner_mode)
        {
            yw_skip_whitespaces(ts);
            YW_CSSDisplay temp;
            if (yw_parse_display_inside(&temp, ts))
            {
                got_something = true;
                got_inner_mode = true;
                res |= temp;
            }
        }
        if (!got_something)
        {
            break;
        }
    }
    if (got_outer_mode || got_inner_mode)
    {
        if (!got_inner_mode)
        {
            res |= YW_CSS_DISPLAY_FLOW;
        }
        if (!got_outer_mode)
        {
            if ((res & YW_CSS_DISPLAY_INNER_MODE_MASK) == YW_CSS_DISPLAY_RUBY)
            {
                res |= YW_CSS_DISPLAY_INLINE;
            }
            else
            {
                res |= YW_CSS_DISPLAY_BLOCK;
            }
        }
        *out = res;
        return true;
    }

    /* Try display-listitem ***************************************************/
    /*
     * https://www.w3.org/TR/css-display-3/#typedef-display-listitem
     */
    /* TODO */

    /* Try display-internal ***************************************************/
    if (yw_expect_ident(ts, "table-row-group"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-row-group
         */
        *out = YW_CSS_DISPLAY_TABLE_ROW_GROUP | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-header-group"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-header-group
         */
        *out = YW_CSS_DISPLAY_TABLE_HEADER_GROUP | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-footer-group"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-footer-group
         */
        *out = YW_CSS_DISPLAY_TABLE_FOOTER_GROUP | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-row"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-row
         */
        *out = YW_CSS_DISPLAY_TABLE_ROW | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-cell"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-cell
         */
        *out = YW_CSS_DISPLAY_TABLE_CELL | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-column-group"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-column-group
         */
        *out = YW_CSS_DISPLAY_TABLE_COLUMN_GROUP | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-column"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-column
         */
        *out = YW_CSS_DISPLAY_TABLE_COLUMN | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "table-caption"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-table-caption
         */
        *out = YW_CSS_DISPLAY_TABLE_CAPTION | YW_CSS_DISPLAY_FLOW_ROOT;
        return true;
    }
    else if (yw_expect_ident(ts, "ruby-base"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-ruby-base
         */
        *out = YW_CSS_DISPLAY_RUBY_BASE | YW_CSS_DISPLAY_FLOW;
        return true;
    }
    else if (yw_expect_ident(ts, "ruby-text"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-ruby-text
         */
        *out = YW_CSS_DISPLAY_RUBY_TEXT | YW_CSS_DISPLAY_FLOW;
        return true;
    }
    else if (yw_expect_ident(ts, "ruby-base-container"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-ruby-base-container
         */
        *out = YW_CSS_DISPLAY_RUBY_BASE_CONTAINER | YW_CSS_DISPLAY_FLOW;
        return true;
    }
    else if (yw_expect_ident(ts, "ruby-text-container"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-ruby-text-container
         */
        *out = YW_CSS_DISPLAY_RUBY_TEXT_CONTAINER | YW_CSS_DISPLAY_FLOW;
        return true;
    }

    /* Try display-box ********************************************************/
    /*
     * https://www.w3.org/TR/css-display-3/#typedef-display-box
     */
    if (yw_expect_ident(ts, "contents"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-contents
         */
        *out = YW_CSS_DISPLAY_CONTENTS;
        return true;
    }
    else if (yw_expect_ident(ts, "none"))
    {
        /*
         * https://www.w3.org/TR/css-display-3/#valdef-display-none
         */
        *out = YW_CSS_DISPLAY_NONE;
        return true;
    }

    return false;
}
bool yw_css_parse_visibility(YW_CSSVisibility *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "visible"))
    {
        *out = YW_CSS_VISIBLE;
        return true;
    }
    else if (yw_expect_ident(ts, "hidden"))
    {
        *out = YW_CSS_HIDDEN;
        return true;
    }
    else if (yw_expect_ident(ts, "collapse"))
    {
        *out = YW_CSS_COLLAPSE;
        return true;
    }
    return false;
}

/*******************************************************************************
 *
 * CSS2 9.5 Floats
 *
 * https://www.w3.org/TR/CSS2/visuren.html#floats
 *
 ******************************************************************************/

bool yw_css_parse_float(YW_CSSFloat *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "none"))
    {
        *out = YW_CSS_NO_FLOAT;
        return true;
    }
    else if (yw_expect_ident(ts, "left"))
    {
        *out = YW_CSS_FLOAT_LEFT;
        return true;
    }
    else if (yw_expect_ident(ts, "right"))
    {
        *out = YW_CSS_FLOAT_RIGHT;
        return true;
    }
    return false;
}

/*******************************************************************************
 *
 * CSS Fonts
 *
 * https://www.w3.org/TR/css-fonts-3
 *
 ******************************************************************************/

static bool yw_parse_font_family_name_ident(char **out, YW_CSSTokenStream *ts)
{
    YW_CSSToken const *temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (temp_tk == NULL)
    {
        return false;
    }
    *out = yw_duplicate_str(temp_tk->ident_tk.value);
    return true;
}

static bool yw_parse_font_family_item(YW_CSSFontFamily *out, YW_CSSTokenStream *ts)
{
    /* Try generic font family ************************************************/
    /*
     * https://www.w3.org/TR/css-fonts-3/#generic-family-value
     */
    out->name = NULL;
    if (yw_expect_ident(ts, "serif"))
    {
        out->family = YW_CSS_SERIF;
        return true;
    }
    else if (yw_expect_ident(ts, "sans-serif"))
    {
        out->family = YW_CSS_SANS_SERIF;
        return true;
    }
    else if (yw_expect_ident(ts, "cursive"))
    {
        out->family = YW_CSS_CURSIVE;
        return true;
    }
    else if (yw_expect_ident(ts, "fantasy"))
    {
        out->family = YW_CSS_FANTASY;
        return true;
    }
    else if (yw_expect_ident(ts, "monospace"))
    {
        out->family = YW_CSS_MONOSPACE;
        return true;
    }

    YW_CSSToken const *temp_tk;

    /* Try string name ********************************************************/
    temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_STRING);
    if (temp_tk != NULL)
    {
        out->family = YW_CSS_NON_GENERIC_FONT_FAMILY;
        out->name = yw_duplicate_str(temp_tk->string_tk.value);
        return true;
    }
    /* Try series of identifier ***********************************************/
    char **names;
    int names_len;
    YW_CSS_PARSE_REPEATION(char *, &names, &names_len, ts, YW_CSS_NO_MAX_REPEATS, yw_parse_font_family_name_ident);
    if (names_len == 0)
    {
        return false;
    }
    char *final_name = NULL;
    for (int i = 0; i < names_len; i++)
    {
        yw_append_str(&final_name, names[i]);
        if (i != names_len - 1)
        {
            yw_append_str(&final_name, " ");
        }
        free(names[i]);
    }
    free(names);
    out->family = YW_CSS_NON_GENERIC_FONT_FAMILY;
    out->name = final_name;
    return true;
}

bool yw_css_parse_font_family(YW_CSSFontFamilies *out, YW_CSSTokenStream *ts)
{
    YW_CSSFontFamily *families;
    int families_len;
    YW_CSS_PARSE_COMMA_SEPARATED_REPEATION(YW_CSSFontFamily, &families, &families_len, ts, YW_CSS_NO_MAX_REPEATS, yw_parse_font_family_item);
    if (families_len == 0)
    {
        return false;
    }
    out->items = families;
    out->len = families_len;
    return true;
}

bool yw_css_parse_font_weight(YW_CSSFontWeight *out, YW_CSSTokenStream *ts)
{
    /*
     * https://www.w3.org/TR/css-fonts-3/#propdef-font-weight
     */

    int old_cursor = ts->cursor;
    /* Try predefined weights *************************************************/
    if (yw_expect_ident(ts, "normal"))
    {
        *out = YW_CSS_NORMAL_FONT_WEIGHT;
        return true;
    }
    else if (yw_expect_ident(ts, "bold"))
    {
        *out = YW_CSS_BOLD;
        return true;
    }
    /* Try numeric weight values **********************************************/
    double weight_d;
    if (yw_css_parse_number(&weight_d, ts))
    {
        /* FIXME: We shouldn't be accepting floating point values */
        int weight = weight_d;
        if (weight < 0 || 1000 < weight)
        {
            /* Out of range */
            goto fail;
        }
        *out = (YW_CSSFontWeight)weight;
        return true;
    }
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

bool yw_css_parse_font_stretch(YW_CSSFontStretch *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "ultra-condensed"))
    {
        *out = YW_CSS_ULTRA_CONDENSED;
        return true;
    }
    else if (yw_expect_ident(ts, "extra-condensed"))
    {
        *out = YW_CSS_EXTRA_CONDENSED;
        return true;
    }
    else if (yw_expect_ident(ts, "condensed"))
    {
        *out = YW_CSS_CONDENSED;
        return true;
    }
    else if (yw_expect_ident(ts, "semi-condensed"))
    {
        *out = YW_CSS_SEMI_CONDENSED;
        return true;
    }
    else if (yw_expect_ident(ts, "normal"))
    {
        *out = YW_CSS_NORMAL_FONT_STRETCH;
        return true;
    }
    else if (yw_expect_ident(ts, "semi-expanded"))
    {
        *out = YW_CSS_SEMI_EXPANDED;
        return true;
    }
    else if (yw_expect_ident(ts, "expanded"))
    {
        *out = YW_CSS_EXPANDED;
        return true;
    }
    else if (yw_expect_ident(ts, "extra-expanded"))
    {
        *out = YW_CSS_EXTRA_EXPANDED;
        return true;
    }
    else if (yw_expect_ident(ts, "ultra-expanded"))
    {
        *out = YW_CSS_ULTRA_EXPANDED;
        return true;
    }
    return false;
}

bool yw_css_parse_font_style(YW_CSSFontStyle *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "normal"))
    {
        *out = YW_CSS_NORMAL_FONT_STYLE;
        return true;
    }
    else if (yw_expect_ident(ts, "italic"))
    {
        *out = YW_CSS_ITALIC;
        return true;
    }
    else if (yw_expect_ident(ts, "oblique"))
    {
        *out = YW_CSS_OBLIQUE;
        return true;
    }
    return false;
}

bool yw_css_parse_font_size(YW_CSSFontSize *out, YW_CSSTokenStream *ts)
{
    /* Try absolute-size keywords *********************************************/
    /*
     * https://www.w3.org/TR/css-fonts-3/#absolute-size-value
     */
    if (yw_expect_ident(ts, "xx-small"))
    {
        out->type = YW_CSS_XX_SMALL;
        return true;
    }
    else if (yw_expect_ident(ts, "x-small"))
    {
        out->type = YW_CSS_X_SMALL;
        return true;
    }
    else if (yw_expect_ident(ts, "small"))
    {
        out->type = YW_CSS_SMALL;
        return true;
    }
    else if (yw_expect_ident(ts, "medium"))
    {
        out->type = YW_CSS_MEDIUM_FONT_SIZE;
        return true;
    }
    else if (yw_expect_ident(ts, "large"))
    {
        out->type = YW_CSS_LARGE;
        return true;
    }
    else if (yw_expect_ident(ts, "x-large"))
    {
        out->type = YW_CSS_X_LARGE;
        return true;
    }
    else if (yw_expect_ident(ts, "xx-large"))
    {
        out->type = YW_CSS_XX_LARGE;
        return true;
    }
    /* Try relative-size keywords *********************************************/
    /*
     * https://www.w3.org/TR/css-fonts-3/#relative-size-value
     */
    if (yw_expect_ident(ts, "larger"))
    {
        out->type = YW_CSS_LARGER;
        return true;
    }
    else if (yw_expect_ident(ts, "smaller"))
    {
        out->type = YW_CSS_SMALLER;
        return true;
    }
    /* Try length/percentage **************************************************/
    if (yw_css_parse_length_or_percentage(&out->size, ts, YW_ALLOW_ZERO_SHORTHAND))
    {
        out->type = YW_CSS_LENGTH_FONT_SIZE;
        return true;
    }
    return false;
}

/*******************************************************************************
 *
 * CSS Selectors
 *
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/
 *
 ******************************************************************************/

static bool yw_parse_ns_prefix(char **out, YW_CSSTokenStream *ts)
{
    /* STUB */
    (void)out;
    (void)ts;
    return false;
}

static bool yw_parse_wq_name(YW_CSSWqName *out, YW_CSSTokenStream *ts)
{
    char *ns_prefix = NULL;
    if (!yw_parse_ns_prefix(&ns_prefix, ts))
    {
        ns_prefix = NULL;
    }
    YW_CSSToken const *temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (temp_tk == NULL)
    {
        goto fail;
    }
    out->ns_prefix = ns_prefix;
    out->ident = yw_duplicate_str(temp_tk->ident_tk.value);
    return true;
fail:
    free(ns_prefix);
    return false;
}

static bool yw_parse_type_selector(YW_CSSSelector *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    if (yw_parse_wq_name(&out->type_sel.name, ts))
    {
        /* <wq-name> **********************************************************/
        out->type = YW_CSS_SELECTOR_TYPE;
        return true;
    }
    else
    {
        /* <ns-prefix?> * *****************************************************/
        if (!yw_parse_ns_prefix(&out->universal_sel.ns_prefix, ts))
        {
            out->universal_sel.ns_prefix = NULL;
        }
        /* ns-prefix? <*> *****************************************************/
        YW_CSSToken const *temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_DELIM);
        if (temp_tk == NULL || temp_tk->delim_tk.value != '*')
        {
            goto fail;
        }
        out->type = YW_CSS_SELECTOR_UNIVERSAL;
        return true;
    }
fail:
    ts->cursor = old_cursor;
    return false;
}

static bool yw_parse_pseudo_class_selector(YW_CSSPseudoClassSelector *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSToken const *temp_tk;

    /* <:>name ****************************************************************/
    /* <:>name(args) **********************************************************/
    if (yw_expect_token(ts, YW_CSS_TOKEN_COLON) == NULL)
    {
        goto fail;
    }
    /* :<name> ****************************************************************/
    temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (temp_tk != NULL)
    {
        out->type = YW_CSS_SELECTOR_PSEUDO_CLASS;
        out->name = yw_duplicate_str(temp_tk->ident_tk.value);
        return true;
    }
    /* :<name(args)> **********************************************************/
    temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_AST_FUNC);
    if (temp_tk != NULL)
    {
        YW_CSSTokenStream inner_ts;
        inner_ts.tokens = temp_tk->ast_func_tk.tokens;
        inner_ts.tokens_len = temp_tk->ast_func_tk.tokens_len;
        inner_ts.cursor = 0;
        YW_CSSTokenStream *ts = &inner_ts;

        YW_CSSToken *values;
        int values_len;
        if (!yw_consume_any_value(&values, &values_len, ts))
        {
            goto fail;
        }
        if (!yw_is_end_of_tokens(ts))
        {
            /* Extra junk at the end */
            goto fail;
        }
        out->type = YW_CSS_SELECTOR_PSEUDO_CLASS;
        out->name = yw_duplicate_str(temp_tk->ident_tk.value);
        /* XXX: We don't have args support yet */
        for (int i = 0; i < values_len; i++)
        {
            yw_token_deinit(&values[i]);
        }
        free(values);
        return true;
    }
fail:
    ts->cursor = old_cursor;
    return false;
}

static bool yw_parse_pseudo_element_selector(YW_CSSPseudoClassSelector *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    /* <:>:name ***************************************************************/
    /* <:>:name(args) *********************************************************/
    if (yw_expect_token(ts, YW_CSS_TOKEN_COLON) == NULL)
    {
        goto fail;
    }
    /* :<:name> ***************************************************************/
    /* :<:name(args)> *********************************************************/
    if (!yw_parse_pseudo_class_selector(out, ts))
    {
        goto fail;
    }
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

static bool yw_parse_subclass_selector(YW_CSSSelector *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    /* Try id selector ********************************************************/
    /*
     * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-id-selector
     */
    YW_CSSToken const *temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_HASH);
    if (temp_tk == NULL)
    {
        goto class_selector;
    }
    if (temp_tk->hash_tk.type != YW_HASH_ID)
    {
        goto fail;
    }
    out->type = YW_CSS_SELECTOR_ID;
    out->id_sel.id = yw_duplicate_str(temp_tk->hash_tk.value);
    return true;
class_selector:
    /* Try class selector *****************************************************/
    ts->cursor = old_cursor;
    /*
     * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-class-selector
     */
    if (!yw_expect_delim(ts, '.'))
    {
        goto attr_selector;
    }
    temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
    if (temp_tk == NULL)
    {
        goto fail;
    }
    out->type = YW_CSS_SELECTOR_CLASS;
    out->class_sel.class_name = yw_duplicate_str(temp_tk->ident_tk.value);
    return true;
attr_selector:
    /* Try attribute selector *************************************************/
    ts->cursor = old_cursor;
    YW_CSSTokenStream inner_ts;
    if (yw_expect_simple_block(&inner_ts, ts, YW_SIMPLE_BLOCK_SQUARE))
    {
        YW_CSSTokenStream *ts = &inner_ts;
        YW_CSSWqName name;
        char const *attr_value = NULL;
        YW_CSSValueMatchType match_type = YW_CSS_NO_VALUE_MATCH;
        /*
         * FIXME: Default case-sensitivity depends on the document language.
         */
        bool is_case_sensitive = true;

        /* [<  >attr  ] *******************************************************/
        /* [<  >attr  =  value  modifier  ] ***********************************/
        yw_skip_whitespaces(ts);
        /* [  <attr>  ] *******************************************************/
        /* [  <attr>  =  value  modifier  ] ***********************************/
        memset(&name, 0, sizeof(name));
        if (!yw_parse_wq_name(&name, ts))
        {
            goto attr_fail;
        }
        /* [  attr<  >] *******************************************************/
        /* [  attr<  >=  value  modifier  ] ***********************************/
        yw_skip_whitespaces(ts);
        if (!yw_is_end_of_tokens(ts))
        {
            /* [  attr  <=>  value  modifier  ] *******************************/
            if (yw_expect_delim(ts, '~'))
            {
                match_type = YW_CSS_VALUE_TILDE_EQUALS;
            }
            else if (yw_expect_delim(ts, '|'))
            {
                match_type = YW_CSS_VALUE_BAR_EQUALS;
            }
            else if (yw_expect_delim(ts, '^'))
            {
                match_type = YW_CSS_VALUE_CARET_EQUALS;
            }
            else if (yw_expect_delim(ts, '$'))
            {
                match_type = YW_CSS_VALUE_DOLLAR_EQUALS;
            }
            else if (yw_expect_delim(ts, '*'))
            {
                match_type = YW_CSS_VALUE_ASTERISK_EQUALS;
            }
            else
            {
                match_type = YW_CSS_VALUE_EQUALS;
            }
            if (!yw_expect_delim(ts, '='))
            {
                goto no_attr_value;
            }
            /* [  attr  =<  >value  modifier  ] *******************************/
            yw_skip_whitespaces(ts);
            /* [  attr  =  <value>  modifier  ] *******************************/
            temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_IDENT);
            if (temp_tk != NULL)
            {
                attr_value = temp_tk->ident_tk.value;
            }
            else
            {
                temp_tk = yw_expect_token(ts, YW_CSS_TOKEN_STRING);
                if (temp_tk != NULL)
                {
                    attr_value = temp_tk->string_tk.value;
                }
                else
                {
                    goto fail;
                }
            }
            /* [  attr  =  value<  >modifier  ] *******************************/
            yw_skip_whitespaces(ts);
            /* [  attr  =  value  <modifier>  ] *******************************/
            if (yw_expect_ident(ts, "i"))
            {
                is_case_sensitive = false;
            }
            else
            {
                is_case_sensitive = true;
            }
            /* [  attr  =  value  modifier<  >] *******************************/
            yw_skip_whitespaces(ts);
        }
    no_attr_value:
        out->type = YW_CSS_SELECTOR_ATTR;
        out->attr_sel.attr_name = name;
        out->attr_sel.value_match_type = match_type;
        out->attr_sel.attr_value = yw_duplicate_str(attr_value);
        out->attr_sel.is_case_sensitive = is_case_sensitive;
        return true;
    attr_fail:
        yw_css_wq_name_deinit(&name);
    }
    /* Try pseudo class selector **********************************************/
    ts->cursor = old_cursor;
    if (yw_parse_pseudo_class_selector(&out->pseudo_class_sel, ts))
    {
        return true;
    }
fail:
    ts->cursor = old_cursor;
    return false;
}

static bool yw_parse_compound_selector(YW_CSSCompoundSelector *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSSelector sel_temp;
    YW_CSSPseudoClassSelector class_sel_temp;
    YW_CSSSelector *type_sel = NULL;
    if (yw_parse_type_selector(&sel_temp, ts))
    {
        type_sel = YW_ALLOC(YW_CSSSelector, 1);
        *type_sel = sel_temp;
        memset(&sel_temp, 0, sizeof(sel_temp));
    }
    YW_CSSSelector *subclass_sels = NULL;
    int subclass_sels_len = 0;
    int subclass_sels_cap = 0;

    while (yw_parse_subclass_selector(&sel_temp, ts))
    {
        YW_PUSH(YW_CSSSelector, &subclass_sels_cap, &subclass_sels_len, &subclass_sels, sel_temp);
    }
    YW_SHRINK_TO_FIT(YW_CSSSelector, &subclass_sels_cap, subclass_sels_len, &subclass_sels);

    YW_CSSCompundSelectorPseudoItem *pseudo_items = NULL;
    int pseudo_items_len = 0;
    int pseudo_items_cap = 0;

    while (yw_parse_pseudo_element_selector(&class_sel_temp, ts))
    {
        YW_CSSSelector *elem_sel = YW_ALLOC(YW_CSSSelector, 1);
        elem_sel->pseudo_class_sel = class_sel_temp;
        memset(&class_sel_temp, 0, sizeof(class_sel_temp));

        YW_CSSSelector *class_sels = NULL;
        int class_sels_len = 0;
        int class_sels_cap = 0;
        while (yw_parse_pseudo_class_selector(&class_sel_temp, ts))
        {
            YW_GROW(YW_CSSSelector, &class_sels_cap, &class_sels_len, &class_sels, 1);
            class_sels[class_sels_len - 1].pseudo_class_sel = class_sel_temp;
        }
        YW_SHRINK_TO_FIT(YW_CSSSelector, &class_sels_cap, class_sels_len, &class_sels);

        YW_CSSCompundSelectorPseudoItem item;
        item.class_sels = class_sels;
        item.class_sels_len = class_sels_len;
        item.pseudo_elem_sel = elem_sel;
        YW_PUSH(YW_CSSCompundSelectorPseudoItem, &pseudo_items_cap, &pseudo_items_len, &pseudo_items, item);
    }
    YW_SHRINK_TO_FIT(YW_CSSCompundSelectorPseudoItem, &pseudo_items_cap, pseudo_items_len, &pseudo_items);
    if (type_sel == NULL && subclass_sels_len == 0 && pseudo_items_len == 0)
    {
        goto fail;
    }
    out->type = YW_CSS_SELECTOR_COMPOUND;
    out->type_sel = type_sel;
    out->subclass_sels = subclass_sels;
    out->subclass_sels_len = subclass_sels_len;
    out->pseudo_items = pseudo_items;
    out->pseudo_items_len = pseudo_items_len;
    return true;
fail:
    ts->cursor = old_cursor;
    return false;
}

static bool yw_parse_complex_selector(YW_CSSComplexSelector *out, YW_CSSTokenStream *ts)
{
    int old_cursor = ts->cursor;
    YW_CSSCompoundSelector base;
    YW_CSSComplexSelectorRest *rests = NULL;
    int rests_len = 0;
    int rests_cap = 0;
    if (!yw_parse_compound_selector(&base, ts))
    {
        goto fail;
    }
    while (1)
    {
        int cursor_before_comb = ts->cursor;
        YW_CSSCombinator comb;

        yw_skip_whitespaces(ts);
        if (yw_expect_delim(ts, '>'))
        {
            comb = YW_CSS_DIRECT_CHILD_COMBINATOR;
        }
        else if (yw_expect_delim(ts, '+'))
        {
            comb = YW_CSS_PLUS_COMBINATOR;
        }
        else if (yw_expect_delim(ts, '~'))
        {
            comb = YW_CSS_TILDE_COMBINATOR;
        }
        else if (yw_expect_delim(ts, '|'))
        {
            if (!yw_expect_delim(ts, '|'))
            {
                goto fail;
            }
            comb = YW_CSS_TWO_BARS_COMBINATOR;
        }
        else
        {
            comb = YW_CSS_CHILD_COMBINATOR;
        }
        yw_skip_whitespaces(ts);

        YW_CSSCompoundSelector compound_temp;
        if (!yw_parse_compound_selector(&compound_temp, ts))
        {
            ts->cursor = cursor_before_comb;
            break;
        }
        YW_CSSComplexSelectorRest rest_item;
        YW_CSSSelector *another = YW_ALLOC(YW_CSSSelector, 1);
        another->compound_sel = compound_temp;
        memset(&compound_temp, 0, sizeof(compound_temp));
        rest_item.selector = another;
        rest_item.combinator = comb;
        YW_PUSH(YW_CSSComplexSelectorRest, &rests_cap, &rests_len, &rests, rest_item);
    }
    YW_SHRINK_TO_FIT(YW_CSSComplexSelectorRest, &rests_cap, rests_len, &rests);
    out->type = YW_CSS_SELECTOR_COMPLEX;
    out->base = YW_ALLOC(YW_CSSSelector, 1);
    out->base->compound_sel = base;
    out->rests = rests;
    out->rests_len = rests_len;
    return true;
fail:
    for (int i = 0; i < rests_len; i++)
    {
        yw_css_selector_deinit((YW_CSSSelector *)rests[i].selector);
        free(rests[i].selector);
    }
    free(rests);
    ts->cursor = old_cursor;
    return false;
}

static bool yw_parse_complex_selector_list(YW_CSSComplexSelector **sels_out, int *len_out, YW_CSSTokenStream *ts)
{
    /*
     * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector-list
     */
    YW_CSS_PARSE_COMMA_SEPARATED_REPEATION(YW_CSSComplexSelector, sels_out, len_out, ts, YW_CSS_NO_MAX_REPEATS, yw_parse_complex_selector);
    if (*len_out == 0)
    {
        return false;
    }
    return true;
}

bool yw_css_parse_selector_list(YW_CSSSelector **sels_out, int *len_out, YW_CSSTokenStream *ts)
{
    YW_CSSComplexSelector *complex_sels;
    int complex_sels_len;
    if (!yw_parse_complex_selector_list(&complex_sels, &complex_sels_len, ts))
    {
        return false;
    }
    *sels_out = YW_ALLOC(YW_CSSSelector, complex_sels_len);
    *len_out = complex_sels_len;
    for (int i = 0; i < complex_sels_len; i++)
    {
        (*sels_out)[i].complex_sel = complex_sels[i];
    }
    free(complex_sels);
    return true;
}

bool yw_css_parse_selector(YW_CSSSelector **sels_out, int *len_out, uint8_t const *bytes, int bytes_len, const char *source_name)
{
    YW_CSSTokenStream ts;
    if (!yw_css_tokenize(&ts, bytes, bytes_len, source_name))
    {
        return false;
    }
    bool res = yw_css_parse_selector_list(sels_out, len_out, &ts);
    for (int i = 0; i < ts.tokens_len; i++)
    {
        yw_token_deinit(&ts.tokens[i]);
    }
    free(ts.tokens);
    return res;
}

/*******************************************************************************
 *
 * CSS Sizing
 *
 * https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/
 *
 ******************************************************************************/

typedef enum
{
    YW_ACCEPT_AUTO_SIZE = 1 << 0,
    YW_ACCEPT_NONE_SIZE = 1 << 1,
    YW_ACCEPT_AUTO_AND_NONE_SIZE = YW_ACCEPT_AUTO_SIZE | YW_ACCEPT_NONE_SIZE,
} YW_SizeAcceptFlags;

static bool yw_parse_size_value_impl(YW_CSSSize *out, YW_CSSTokenStream *ts, YW_SizeAcceptFlags accept_flags)
{
    int old_cursor = ts->cursor;
    if ((accept_flags & YW_ACCEPT_AUTO_SIZE) && yw_expect_ident(ts, "auto"))
    {
        out->type = YW_CSS_AUTO_SIZE;
        return true;
    }
    if ((accept_flags & YW_ACCEPT_NONE_SIZE) && yw_expect_ident(ts, "none"))
    {
        out->type = YW_CSS_NO_SIZE;
        return true;
    }
    if (yw_expect_ident(ts, "min-content"))
    {
        out->type = YW_CSS_MIN_CONTENT;
        return true;
    }
    if (yw_expect_ident(ts, "max-content"))
    {
        out->type = YW_CSS_MAX_CONTENT;
        return true;
    }
    YW_CSSTokenStream inner_ts;
    if (yw_expect_ast_func(&inner_ts, ts, "fit-content"))
    {
        YW_CSSTokenStream *ts = &inner_ts;
        out->type = YW_CSS_FIT_CONTENT;
        if (!yw_css_parse_length_or_percentage(&out->size, ts, YW_ALLOW_ZERO_SHORTHAND))
        {
            goto fail;
        }
        if (!yw_is_end_of_tokens(ts))
        {
            /* Extra junk at the end */
            goto fail;
        }
        return true;
    }
    if (yw_css_parse_length_or_percentage(&out->size, ts, YW_ALLOW_ZERO_SHORTHAND))
    {
        out->type = YW_CSS_MANUAL_SIZE;
        return true;
    }
fail:
    ts->cursor = old_cursor;
    return false;
}

bool yw_css_parse_size_or_auto(YW_CSSSize *out, YW_CSSTokenStream *ts)
{
    return yw_parse_size_value_impl(out, ts, YW_ACCEPT_AUTO_SIZE);
}
bool yw_css_parse_size_or_none(YW_CSSSize *out, YW_CSSTokenStream *ts)
{
    return yw_parse_size_value_impl(out, ts, YW_ACCEPT_NONE_SIZE);
}

/*******************************************************************************
 *
 * CSS Text
 *
 * https://www.w3.org/TR/css-text-3
 *
 ******************************************************************************/

bool yw_css_parse_text_transform(YW_CSSTextTransform *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "none"))
    {
        out = 0;
        return true;
    }
    YW_CSSTextTransform res = 0;
    bool got_type = false, got_is_full_width = false, got_is_full_kana = false;
    while (1)
    {
        bool got_something = false;
        if (!got_type)
        {
            yw_skip_whitespaces(ts);
            if (yw_expect_ident(ts, "capitalize"))
            {
                res |= YW_CSS_TEXT_TRANSFORM_CAPITALIZE;
                got_type = true;
                got_something = true;
            }
            else if (yw_expect_ident(ts, "uppercase"))
            {
                res |= YW_CSS_TEXT_TRANSFORM_UPPERCASE;
                got_type = true;
                got_something = true;
            }
            else if (yw_expect_ident(ts, "lowercase"))
            {
                res |= YW_CSS_TEXT_TRANSFORM_LOWERCASE;
                got_type = true;
                got_something = true;
            }
        }
        if (!got_is_full_width)
        {
            yw_skip_whitespaces(ts);
            if (yw_expect_ident(ts, "full-width"))
            {
                res |= YW_CSS_TEXT_TRANSFORM_FULL_WIDTH;
                got_is_full_width = true;
                got_something = true;
            }
        }
        if (!got_is_full_kana)
        {
            yw_skip_whitespaces(ts);
            if (yw_expect_ident(ts, "full-size-kana"))
            {
                res |= YW_CSS_TEXT_TRANSFORM_FULL_SIZE_KANA;
                got_is_full_kana = true;
                got_something = true;
            }
        }
        yw_skip_whitespaces(ts);
        if (!got_something)
        {
            break;
        }
    }
    if (res == 0)
    {
        return false;
    }
    *out = res;
    return true;
}

/*******************************************************************************
 *
 * CSS Text Decoration
 *
 * https://www.w3.org/TR/css-text-decor-3
 *
 ******************************************************************************/

bool yw_css_parse_text_decoration_line(YW_CSSTextDecorationLine *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "none"))
    {
        out = 0;
        return true;
    }
    YW_CSSTextTransform res = 0;
    while (1)
    {
        bool got_something = false;
        yw_skip_whitespaces(ts);
        if (yw_expect_ident(ts, "underline"))
        {
            res |= YW_CSS_TEXT_DECORATION_UNDERLINE;
            got_something = true;
        }
        if (yw_expect_ident(ts, "overline"))
        {
            res |= YW_CSS_TEXT_DECORATION_OVERLINE;
            got_something = true;
        }
        if (yw_expect_ident(ts, "line-through"))
        {
            res |= YW_CSS_TEXT_DECORATION_LINE_THROUGH;
            got_something = true;
        }
        if (yw_expect_ident(ts, "blink"))
        {
            res |= YW_CSS_TEXT_DECORATION_BLINK;
            got_something = true;
        }
        yw_skip_whitespaces(ts);
        if (!got_something)
        {
            break;
        }
    }
    if (*out == 0)
    {
        return false;
    }
    *out = res;
    return true;
}

bool yw_css_parse_text_decoration_style(YW_CSSTextDecorationStyle *out, YW_CSSTokenStream *ts)
{
    if (yw_expect_ident(ts, "solid"))
    {
        *out = YW_CSS_TEXT_DECORATION_SOLID;
        return true;
    }
    if (yw_expect_ident(ts, "double"))
    {
        *out = YW_CSS_TEXT_DECORATION_DOUBLE;
        return true;
    }
    if (yw_expect_ident(ts, "dotted"))
    {
        *out = YW_CSS_TEXT_DECORATION_DOTTED;
        return true;
    }
    if (yw_expect_ident(ts, "dashed"))
    {
        *out = YW_CSS_TEXT_DECORATION_DASHED;
        return true;
    }
    if (yw_expect_ident(ts, "wavy"))
    {
        *out = YW_CSS_TEXT_DECORATION_WAVY;
        return true;
    }
    return false;
}
