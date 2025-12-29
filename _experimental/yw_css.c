/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_css.h"
#include "yw_common.h"
#include "yw_dom.h"
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

char const *yw_css_length_unit_str(YW_CSSLengthUnit unit)
{
    switch (unit)
    {
    case YW_CSS_EM:
        return "em";
    case YW_CSS_EX:
        return "ex";
    case YW_CSS_CH:
        return "ch";
    case YW_CSS_REM:
        return "rem";
    case YW_CSS_VW:
        return "vw";
    case YW_CSS_VH:
        return "vh";
    case YW_CSS_VMIN:
        return "vmin";
    case YW_CSS_VMAX:
        return "vmax";
    case YW_CSS_CM:
        return "cm";
    case YW_CSS_MM:
        return "mm";
    case YW_CSS_Q:
        return "q";
    case YW_CSS_PC:
        return "pc";
    case YW_CSS_PT:
        return "pt";
    case YW_CSS_PX:
        return "px";
    }
    YW_ILLEGAL_VALUE(unit);
}

double yw_css_length_to_px(YW_CSSLength const *len, double font_size)
{
    switch (len->unit)
    {
    case YW_CSS_EM:
        return font_size * len->value;
    case YW_CSS_EX:
        YW_TODO();
    case YW_CSS_CH:
        YW_TODO();
    case YW_CSS_REM:
        YW_TODO();
    case YW_CSS_VW:
        YW_TODO();
    case YW_CSS_VH:
        YW_TODO();
    case YW_CSS_VMIN:
        YW_TODO();
    case YW_CSS_VMAX:
        YW_TODO();
    case YW_CSS_CM:
        YW_TODO();
    case YW_CSS_MM:
        YW_TODO();
    case YW_CSS_Q:
        YW_TODO();
    case YW_CSS_PC:
        YW_TODO();
    case YW_CSS_PT:
        YW_TODO();
    case YW_CSS_PX:
        return len->value;
    }
    YW_ILLEGAL_VALUE(len->unit);
}

double yw_css_percentage_to_px(double per, double container_size)
{
    return (per * container_size) / 100;
}
void yw_css_percentage_to_length(YW_CSSLength *out, double per, double container_size)
{
    out->unit = YW_CSS_PX;
    out->value = yw_css_percentage_to_px(per, container_size);
}

double yw_css_length_or_percentage_to_px(YW_CSSLengthOrPercentage const *len_or_per, double font_size, double container_size)
{
    if (len_or_per->is_percentage)
    {
        return yw_css_percentage_to_px(len_or_per->value.percentage, container_size);
    }
    return yw_css_length_to_px(&len_or_per->value.length, font_size);
}

/*******************************************************************************
 *
 * CSS Backgrounds and Borders
 *
 * https://www.w3.org/TR/css-backgrounds-3/
 *
 ******************************************************************************/

char const *yw_css_line_style_str(YW_CSSLineStyle style)
{
    switch (style)
    {
    case YW_CSS_NO_LINE:
        return "none";
    case YW_CSS_HIDDEN_LINE:
        return "hidden";
    case YW_CSS_DOTTED_LINE:
        return "dotted";
    case YW_CSS_DASHED_LINE:
        return "dashed";
    case YW_CSS_SOLID_LINE:
        return "solid";
    case YW_CSS_DOUBLE_LINE:
        return "double";
    case YW_CSS_GROOVE_LINE:
        return "grovve";
    case YW_CSS_RIDGE_LINE:
        return "ridge";
    case YW_CSS_INSET_LINE:
        return "inset";
    case YW_CSS_OUTSET_LINE:
        return "outset";
    }
    return "?";
}

/*******************************************************************************
 *
 * CSS Color
 *
 * https://www.w3.org/TR/css-color-4/
 *
 ******************************************************************************/

/* https://www.w3.org/TR/css-color-4/#named-colors */
static struct
{
    char const *name;
    YW_CSSRgba color;
} yw_named_colors[] = {
    {"aliceblue", YW_CSS_RGB(240, 248, 255)},
    {"antiquewhite", YW_CSS_RGB(250, 235, 215)},
    {"aqua", YW_CSS_RGB(0, 255, 255)},
    {"aquamarine", YW_CSS_RGB(127, 255, 212)},
    {"azure", YW_CSS_RGB(240, 255, 255)},
    {"beige", YW_CSS_RGB(245, 245, 220)},
    {"bisque", YW_CSS_RGB(255, 228, 196)},
    {"black", YW_CSS_RGB(0, 0, 0)},
    {"blanchedalmond", YW_CSS_RGB(255, 235, 205)},
    {"blue", YW_CSS_RGB(0, 0, 255)},
    {"blueviolet", YW_CSS_RGB(138, 43, 226)},
    {"brown", YW_CSS_RGB(165, 42, 42)},
    {"burlywood", YW_CSS_RGB(222, 184, 135)},
    {"cadetblue", YW_CSS_RGB(95, 158, 160)},
    {"chartreuse", YW_CSS_RGB(127, 255, 0)},
    {"chocolate", YW_CSS_RGB(210, 105, 30)},
    {"coral", YW_CSS_RGB(255, 127, 80)},
    {"cornflowerblue", YW_CSS_RGB(100, 149, 237)},
    {"cornsilk", YW_CSS_RGB(255, 248, 220)},
    {"crimson", YW_CSS_RGB(220, 20, 60)},
    {"cyan", YW_CSS_RGB(0, 255, 255)},
    {"darkblue", YW_CSS_RGB(0, 0, 139)},
    {"darkcyan", YW_CSS_RGB(0, 139, 139)},
    {"darkgoldenrod", YW_CSS_RGB(184, 134, 11)},
    {"darkgray", YW_CSS_RGB(169, 169, 169)},
    {"darkgreen", YW_CSS_RGB(0, 100, 0)},
    {"darkgrey", YW_CSS_RGB(169, 169, 169)},
    {"darkkhaki", YW_CSS_RGB(189, 183, 107)},
    {"darkmagenta", YW_CSS_RGB(139, 0, 139)},
    {"darkolivegreen", YW_CSS_RGB(85, 107, 47)},
    {"darkorange", YW_CSS_RGB(255, 140, 0)},
    {"darkorchid", YW_CSS_RGB(153, 50, 204)},
    {"darkred", YW_CSS_RGB(139, 0, 0)},
    {"darksalmon", YW_CSS_RGB(233, 150, 122)},
    {"darkseagreen", YW_CSS_RGB(143, 188, 143)},
    {"darkslateblue", YW_CSS_RGB(72, 61, 139)},
    {"darkslategray", YW_CSS_RGB(47, 79, 79)},
    {"darkslategrey", YW_CSS_RGB(47, 79, 79)},
    {"darkturquoise", YW_CSS_RGB(0, 206, 209)},
    {"darkviolet", YW_CSS_RGB(148, 0, 211)},
    {"deeppink", YW_CSS_RGB(255, 20, 147)},
    {"deepskyblue", YW_CSS_RGB(0, 191, 255)},
    {"dimgray", YW_CSS_RGB(105, 105, 105)},
    {"dimgrey", YW_CSS_RGB(105, 105, 105)},
    {"dodgerblue", YW_CSS_RGB(30, 144, 255)},
    {"firebrick", YW_CSS_RGB(178, 34, 34)},
    {"floralwhite", YW_CSS_RGB(255, 250, 240)},
    {"forestgreen", YW_CSS_RGB(34, 139, 34)},
    {"fuchsia", YW_CSS_RGB(255, 0, 255)},
    {"gainsboro", YW_CSS_RGB(220, 220, 220)},
    {"ghostwhite", YW_CSS_RGB(248, 248, 255)},
    {"gold", YW_CSS_RGB(255, 215, 0)},
    {"goldenrod", YW_CSS_RGB(218, 165, 32)},
    {"gray", YW_CSS_RGB(128, 128, 128)},
    {"green", YW_CSS_RGB(0, 128, 0)},
    {"greenyellow", YW_CSS_RGB(173, 255, 47)},
    {"grey", YW_CSS_RGB(128, 128, 128)},
    {"honeydew", YW_CSS_RGB(240, 255, 240)},
    {"hotpink", YW_CSS_RGB(255, 105, 180)},
    {"indianred", YW_CSS_RGB(205, 92, 92)},
    {"indigo", YW_CSS_RGB(75, 0, 130)},
    {"ivory", YW_CSS_RGB(255, 255, 240)},
    {"khaki", YW_CSS_RGB(240, 230, 140)},
    {"lavender", YW_CSS_RGB(230, 230, 250)},
    {"lavenderblush", YW_CSS_RGB(255, 240, 245)},
    {"lawngreen", YW_CSS_RGB(124, 252, 0)},
    {"lemonchiffon", YW_CSS_RGB(255, 250, 205)},
    {"lightblue", YW_CSS_RGB(173, 216, 230)},
    {"lightcoral", YW_CSS_RGB(240, 128, 128)},
    {"lightcyan", YW_CSS_RGB(224, 255, 255)},
    {"lightgoldenrodyellow", YW_CSS_RGB(250, 250, 210)},
    {"lightgray", YW_CSS_RGB(211, 211, 211)},
    {"lightgreen", YW_CSS_RGB(144, 238, 144)},
    {"lightgrey", YW_CSS_RGB(211, 211, 211)},
    {"lightpink", YW_CSS_RGB(255, 182, 193)},
    {"lightsalmon", YW_CSS_RGB(255, 160, 122)},
    {"lightseagreen", YW_CSS_RGB(32, 178, 170)},
    {"lightskyblue", YW_CSS_RGB(135, 206, 250)},
    {"lightslategray", YW_CSS_RGB(119, 136, 153)},
    {"lightslategrey", YW_CSS_RGB(119, 136, 153)},
    {"lightsteelblue", YW_CSS_RGB(176, 196, 222)},
    {"lightyellow", YW_CSS_RGB(255, 255, 224)},
    {"lime", YW_CSS_RGB(0, 255, 0)},
    {"limegreen", YW_CSS_RGB(50, 205, 50)},
    {"linen", YW_CSS_RGB(250, 240, 230)},
    {"magenta", YW_CSS_RGB(255, 0, 255)},
    {"maroon", YW_CSS_RGB(128, 0, 0)},
    {"mediumaquamarine", YW_CSS_RGB(102, 205, 170)},
    {"mediumblue", YW_CSS_RGB(0, 0, 205)},
    {"mediumorchid", YW_CSS_RGB(186, 85, 211)},
    {"mediumpurple", YW_CSS_RGB(147, 112, 219)},
    {"mediumseagreen", YW_CSS_RGB(60, 179, 113)},
    {"mediumslateblue", YW_CSS_RGB(123, 104, 238)},
    {"mediumspringgreen", YW_CSS_RGB(0, 250, 154)},
    {"mediumturquoise", YW_CSS_RGB(72, 209, 204)},
    {"mediumvioletred", YW_CSS_RGB(199, 21, 133)},
    {"midnightblue", YW_CSS_RGB(25, 25, 112)},
    {"mintcream", YW_CSS_RGB(245, 255, 250)},
    {"mistyrose", YW_CSS_RGB(255, 228, 225)},
    {"moccasin", YW_CSS_RGB(255, 228, 181)},
    {"navajowhite", YW_CSS_RGB(255, 222, 173)},
    {"navy", YW_CSS_RGB(0, 0, 128)},
    {"oldlace", YW_CSS_RGB(253, 245, 230)},
    {"olive", YW_CSS_RGB(128, 128, 0)},
    {"olivedrab", YW_CSS_RGB(107, 142, 35)},
    {"orange", YW_CSS_RGB(255, 165, 0)},
    {"orangered", YW_CSS_RGB(255, 69, 0)},
    {"orchid", YW_CSS_RGB(218, 112, 214)},
    {"palegoldenrod", YW_CSS_RGB(238, 232, 170)},
    {"palegreen", YW_CSS_RGB(152, 251, 152)},
    {"paleturquoise", YW_CSS_RGB(175, 238, 238)},
    {"palevioletred", YW_CSS_RGB(219, 112, 147)},
    {"papayawhip", YW_CSS_RGB(255, 239, 213)},
    {"peachpuff", YW_CSS_RGB(255, 218, 185)},
    {"peru", YW_CSS_RGB(205, 133, 63)},
    {"pink", YW_CSS_RGB(255, 192, 203)},
    {"plum", YW_CSS_RGB(221, 160, 221)},
    {"powderblue", YW_CSS_RGB(176, 224, 230)},
    {"purple", YW_CSS_RGB(128, 0, 128)},
    {"rebeccapurple", YW_CSS_RGB(102, 51, 153)},
    {"red", YW_CSS_RGB(255, 0, 0)},
    {"rosybrown", YW_CSS_RGB(188, 143, 143)},
    {"royalblue", YW_CSS_RGB(65, 105, 225)},
    {"saddlebrown", YW_CSS_RGB(139, 69, 19)},
    {"salmon", YW_CSS_RGB(250, 128, 114)},
    {"sandybrown", YW_CSS_RGB(244, 164, 96)},
    {"seagreen", YW_CSS_RGB(46, 139, 87)},
    {"seashell", YW_CSS_RGB(255, 245, 238)},
    {"sienna", YW_CSS_RGB(160, 82, 45)},
    {"silver", YW_CSS_RGB(192, 192, 192)},
    {"skyblue", YW_CSS_RGB(135, 206, 235)},
    {"slateblue", YW_CSS_RGB(106, 90, 205)},
    {"slategray", YW_CSS_RGB(112, 128, 144)},
    {"slategrey", YW_CSS_RGB(112, 128, 144)},
    {"snow", YW_CSS_RGB(255, 250, 250)},
    {"springgreen", YW_CSS_RGB(0, 255, 127)},
    {"steelblue", YW_CSS_RGB(70, 130, 180)},
    {"tan", YW_CSS_RGB(210, 180, 140)},
    {"teal", YW_CSS_RGB(0, 128, 128)},
    {"thistle", YW_CSS_RGB(216, 191, 216)},
    {"tomato", YW_CSS_RGB(255, 99, 71)},
    {"turquoise", YW_CSS_RGB(64, 224, 208)},
    {"violet", YW_CSS_RGB(238, 130, 238)},
    {"wheat", YW_CSS_RGB(245, 222, 179)},
    {"white", YW_CSS_RGB(255, 255, 255)},
    {"whitesmoke", YW_CSS_RGB(245, 245, 245)},
    {"yellow", YW_CSS_RGB(255, 255, 0)},
    {"yellowgreen", YW_CSS_RGB(154, 205, 50)},
};

/* Returns 0 if not found */
YW_CSSRgba yw_css_color_from_name(char const *name)
{
    for (int i = 0; i < (int)YW_SIZEOF_ARRAY(yw_named_colors); i++)
    {
        if (strcmp(yw_named_colors[i].name, name) == 0)
        {
            return yw_named_colors[i].color;
        }
    }
    return 0;
}

void yw_css_color_from_rgba(YW_CSSColor *out, YW_CSSRgba rgba)
{
    out->type = YW_CSS_RGB_COLOR;
    out->rgb.rgba = rgba;
}
YW_CSSRgba yw_css_color_to_rgba(YW_CSSColor const *color)
{
    switch (color->type)
    {
    case YW_CSS_RGB_COLOR:
        return color->rgb.rgba;
    case YW_CSS_CURRENT_COLOR:
        fprintf(stderr, "currentColor values must be handled by caller\n");
        abort();
    case YW_CSS_HSL_COLOR:
        YW_TODO();
    case YW_CSS_HWB_COLOR:
        YW_TODO();
    case YW_CSS_LAB_COLOR:
        YW_TODO();
    case YW_CSS_LCH_COLOR:
        YW_TODO();
    case YW_CSS_OKLAB_COLOR:
        YW_TODO();
    case YW_CSS_OKLCH_COLOR:
        YW_TODO();
    case YW_CSS_COLOR_FUNC_COLOR:
        YW_TODO();
    }
    YW_ILLEGAL_VALUE(color->type);
}

/*******************************************************************************
 *
 * CSS Display
 *
 * https://www.w3.org/TR/css-display-3/
 *
 ******************************************************************************/

char const *yw_css_visibility_str(YW_CSSVisibility vis)
{
    switch (vis)
    {
    case YW_CSS_VISIBLE:
        return "visible";
    case YW_CSS_HIDDEN:
        return "hidden";
    case YW_CSS_COLLAPSE:
        return "collapse";
    }
    YW_ILLEGAL_VALUE(vis);
}

/*******************************************************************************
 *
 * CSS Fonts
 *
 * https://www.w3.org/TR/css-fonts-3
 *
 ******************************************************************************/

char const *yw_css_generic_font_family_str(YW_CSSGenericFontFamily fam)
{
    switch (fam)
    {
    case YW_CSS_NON_GENERIC_FONT_FAMILY:
        return "<non-generic font-family>";
    case YW_CSS_SERIF:
        return "serif";
    case YW_CSS_SANS_SERIF:
        return "sans_serif";
    case YW_CSS_CURSIVE:
        return "cursive";
    case YW_CSS_FANTASY:
        return "fantasy";
    case YW_CSS_MONOSPACE:
        return "monospace";
    }
    YW_ILLEGAL_VALUE(fam);
}

char const *yw_css_generic_font_stretch_str(YW_CSSFontStretch str)
{
    switch (str)
    {
    case YW_CSS_ULTRA_CONDENSED:
        return "ultra-condensed";
    case YW_CSS_EXTRA_CONDENSED:
        return "extra-condensed";
    case YW_CSS_CONDENSED:
        return "condensed";
    case YW_CSS_SEMI_CONDENSED:
        return "semi-condensed";
    case YW_CSS_NORMAL_FONT_STRETCH:
        return "normal";
    case YW_CSS_SEMI_EXPANDED:
        return "semi-expanded";
    case YW_CSS_EXPANDED:
        return "expanded";
    case YW_CSS_EXTRA_EXPANDED:
        return "extra-expanded";
    case YW_CSS_ULTRA_EXPANDED:
        return "ultra-expanded";
    }
    YW_ILLEGAL_VALUE(str);
}

char const *yw_css_generic_font_style_str(YW_CSSFontStyle sty)
{
    switch (sty)
    {
    case YW_CSS_NORMAL_FONT_STYLE:
        return "normal";
    case YW_CSS_ITALIC:
        return "italic";
    case YW_CSS_OBLIQUE:
        return "oblique";
    }
    YW_ILLEGAL_VALUE(sty);
}

#define YW_CSS_XX_SMALL_PX ((YW_CSS_PREFERRED_FONT_SIZE * 3.0) / 5.0)
#define YW_CSS_X_SMALL_PX ((YW_CSS_PREFERRED_FONT_SIZE * 3.0) / 4.0)
#define YW_CSS_SMALL_PX ((YW_CSS_PREFERRED_FONT_SIZE * 8.0) / 9.0)
#define YW_CSS_MEDIUM_FONT_SIZE_PX (YW_CSS_PREFERRED_FONT_SIZE)
#define YW_CSS_LARGE_PX ((YW_CSS_PREFERRED_FONT_SIZE * 6.0) / 5.0)
#define YW_CSS_X_LARGE_PX ((YW_CSS_PREFERRED_FONT_SIZE * 3.0) / 2.0)
#define YW_CSS_XX_LARGE_PX ((YW_CSS_PREFERRED_FONT_SIZE * 2.0) / 1.0)

static struct
{
    YW_CSSFontSizeType type;
    double size;
} yw_absolute_sizes[] = {
    {YW_CSS_XX_SMALL, YW_CSS_XX_SMALL_PX},
    {YW_CSS_X_SMALL, YW_CSS_X_SMALL_PX},
    {YW_CSS_SMALL, YW_CSS_SMALL_PX},
    {YW_CSS_MEDIUM_FONT_SIZE, YW_CSS_MEDIUM_FONT_SIZE_PX},
    {YW_CSS_LARGE, YW_CSS_LARGE_PX},
    {YW_CSS_X_LARGE, YW_CSS_X_LARGE_PX},
    {YW_CSS_XX_LARGE, YW_CSS_XX_LARGE_PX},
};

static YW_CSSFontSizeType yw_px_to_absolute_size(double size)
{
    double min_diff = (double)INT_MAX;
    YW_CSSFontSizeType res = YW_CSS_XX_SMALL;

    for (int i = 0; i < (int)YW_SIZEOF_ARRAY(yw_absolute_sizes); i++)
    {
        double diff = fabs(size - yw_absolute_sizes[i].size);
        if (diff < min_diff)
        {
            res = yw_absolute_sizes[i].type;
            min_diff = diff;
        }
    }
    return res;
}

double yw_css_font_size_to_px(YW_CSSFontSize const *sz, double font_size, double parent_font_size)
{
    switch (sz->type)
    {
    case YW_CSS_XX_SMALL:
        return YW_CSS_XX_SMALL_PX;
    case YW_CSS_X_SMALL:
        return YW_CSS_X_SMALL_PX;
    case YW_CSS_SMALL:
        return YW_CSS_SMALL_PX;
    case YW_CSS_MEDIUM_FONT_SIZE:
        return YW_CSS_MEDIUM_FONT_SIZE_PX;
    case YW_CSS_LARGE:
        return YW_CSS_LARGE_PX;
    case YW_CSS_X_LARGE:
        return YW_CSS_X_LARGE_PX;
    case YW_CSS_XX_LARGE:
        return YW_CSS_XX_LARGE_PX;
    case YW_CSS_LARGER:
        switch (yw_px_to_absolute_size(parent_font_size))
        {
        case YW_CSS_XX_SMALL:
            return YW_CSS_X_SMALL_PX;
        case YW_CSS_X_SMALL:
            return YW_CSS_SMALL_PX;
        case YW_CSS_SMALL:
            return YW_CSS_MEDIUM_FONT_SIZE_PX;
        case YW_CSS_MEDIUM_FONT_SIZE:
            return YW_CSS_LARGE_PX;
        case YW_CSS_LARGE:
            return YW_CSS_X_LARGE_PX;
        case YW_CSS_X_LARGE:
        case YW_CSS_XX_LARGE:
            return YW_CSS_XX_LARGE_PX;
        default:
            YW_UNREACHABLE();
        }
        break;
    case YW_CSS_SMALLER:
        switch (yw_px_to_absolute_size(parent_font_size))
        {
        case YW_CSS_XX_SMALL:
        case YW_CSS_X_SMALL:
            return YW_CSS_XX_SMALL_PX;
        case YW_CSS_SMALL:
            return YW_CSS_X_SMALL_PX;
        case YW_CSS_MEDIUM_FONT_SIZE:
            return YW_CSS_SMALL_PX;
        case YW_CSS_LARGE:
            return YW_CSS_MEDIUM_FONT_SIZE_PX;
        case YW_CSS_X_LARGE:
            return YW_CSS_LARGE_PX;
        case YW_CSS_XX_LARGE:
            return YW_CSS_X_LARGE_PX;
        default:
            YW_UNREACHABLE();
        }
        break;
    case YW_CSS_LENGTH_FONT_SIZE:
        return yw_css_length_or_percentage_to_px(&sz->size, font_size, parent_font_size);
    }
    YW_ILLEGAL_VALUE(sz->type);
}

/*******************************************************************************
 *
 * CSS Selectors
 *
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/
 *
 ******************************************************************************/

void yw_css_wq_name_deinit(YW_CSSWqName *name)
{
    if (name == NULL)
    {
        return;
    }
    free(name->ident);
    free(name->ns_prefix);
}

void yw_css_selector_deinit(YW_CSSSelector *sel)
{
    if (sel == NULL)
    {
        return;
    }
    switch (sel->type)
    {
    case YW_CSS_SELECTOR_ATTR:
        yw_css_wq_name_deinit(&sel->attr_sel.attr_name);
        free(sel->attr_sel.attr_value);
        break;
    case YW_CSS_SELECTOR_CLASS:
        free(sel->class_sel.class_name);
        break;
    case YW_CSS_SELECTOR_ID:
        free(sel->id_sel.id);
        break;
    case YW_CSS_SELECTOR_TYPE:
        yw_css_wq_name_deinit(&sel->type_sel.name);
        break;
    case YW_CSS_SELECTOR_UNIVERSAL:
        free(sel->universal_sel.ns_prefix);
        break;
    case YW_CSS_SELECTOR_COMPOUND:
        yw_css_selector_deinit(sel->compound_sel.type_sel);
        free(sel->compound_sel.type_sel);
        for (int i = 0; i < sel->compound_sel.subclass_sels_len; i++)
        {
            yw_css_selector_deinit(&sel->compound_sel.subclass_sels[i]);
        }
        free(sel->compound_sel.subclass_sels);
        for (int i = 0; i < sel->compound_sel.pseudo_items_len; i++)
        {
            YW_CSSCompundSelectorPseudoItem *item = &sel->compound_sel.pseudo_items[i];
            yw_css_selector_deinit(item->pseudo_elem_sel);
            free(item->pseudo_elem_sel);
            for (int j = 0; j < item->class_sels_len; j++)
            {
                yw_css_selector_deinit(&item->class_sels[j]);
            }
            free(item->class_sels);
        }
        free(sel->compound_sel.pseudo_items);
        break;
    case YW_CSS_SELECTOR_PSEUDO_CLASS:
        free(sel->pseudo_class_sel.name);
        break;
    case YW_CSS_SELECTOR_COMPLEX:
        yw_css_selector_deinit(sel->complex_sel.base);
        free(sel->complex_sel.base);
        for (int i = 0; i < sel->complex_sel.rests_len; i++)
        {
            YW_CSSComplexSelectorRest *item = &sel->complex_sel.rests[i];
            yw_css_selector_deinit(item->selector);
            free(item->selector);
        }
        free(sel->complex_sel.rests);
        break;
    case YW_CSS_SELECTOR_NODE_PTR:
        break;
    default:
        YW_ILLEGAL_VALUE(sel->type);
    }
}

bool yw_css_selector_match_element(YW_CSSSelector const *sel, YW_GC_PTR(YW_DOMElement) elem)
{
    if (sel == NULL)
    {
        return false;
    }
    switch (sel->type)
    {
    case YW_CSS_SELECTOR_ATTR:
        YW_TODO();
    case YW_CSS_SELECTOR_CLASS: {
        char const *class_attr = yw_dom_attr_of_element(elem, NULL, "class");
        if (class_attr == NULL)
        {
            return false;
        }
        char const *next_class = strrchr(class_attr, ' ');
        next_class = next_class == NULL ? class_attr : next_class;
        while (next_class != NULL)
        {
            char const *next_space = strchr(next_class, ' ');
            size_t got_len = next_space == NULL ? strlen(next_class) : (size_t)(next_space - next_class);
            size_t expected_len = strlen(sel->class_sel.class_name);
            if (got_len != expected_len)
            {
                continue;
            }
            if (strncmp(next_class, sel->class_sel.class_name, got_len) == 0)
            {
                return true;
            }
            next_class = strrchr(next_space, ' ');
            next_class = next_class == NULL ? NULL : next_class + 1;
        }
        return false;
    }
    case YW_CSS_SELECTOR_ID: {
        char const *id_attr = yw_dom_attr_of_element(elem, NULL, "class");
        if (id_attr == NULL)
        {
            return false;
        }
        return strcmp(id_attr, sel->id_sel.id) == 0;
    }
    case YW_CSS_SELECTOR_TYPE:
        /* TODO: Handle namespace */
        return strcmp(elem->local_name, sel->type_sel.name.ident) == 0;
    case YW_CSS_SELECTOR_UNIVERSAL:
        /* TODO: Handle namespace */
        return true;
    case YW_CSS_SELECTOR_COMPOUND: {
        if (!yw_css_selector_match_element(sel->compound_sel.type_sel, elem))
        {
            return false;
        }
        for (int i = 0; i < sel->compound_sel.subclass_sels_len; i++)
        {
            if (!yw_css_selector_match_element(&sel->compound_sel.subclass_sels[i], elem))
            {
                return false;
            }
        }
        for (int i = 0; i < sel->compound_sel.pseudo_items_len; i++)
        {
            /* TODO */
        }
        return true;
    }
    case YW_CSS_SELECTOR_PSEUDO_CLASS:
        /* STUB */
        return false;
    case YW_CSS_SELECTOR_COMPLEX: {
        /* Test each compound selector, from right to left */

        for (int i = sel->complex_sel.rests_len - 1; 0 < i; i--)
        {
            YW_CSSSelector const *prev_sel = (i != 0) ? sel->complex_sel.rests[i - 1].selector : sel->complex_sel.base;
            YW_CSSSelector const *c_sel = sel->complex_sel.rests[i].selector;

            if (!yw_css_selector_match_element((YW_CSSSelector const *)c_sel, elem))
            {
                return false;
            }

            switch (sel->complex_sel.rests[i].combinator)
            {
            case YW_CSS_CHILD_COMBINATOR: {
                /* A B */
                YW_GC_PTR(YW_DOMNode) curr_elem = elem->_node.parent;
                bool found = false;
                while (curr_elem != NULL)
                {
                    if (!(curr_elem->type_flags & YW_DOM_ELEMENT_NODE))
                    {
                        break;
                    }
                    if (yw_css_selector_match_element((YW_CSSSelector const *)prev_sel, (YW_GC_PTR(YW_DOMElement))curr_elem))
                    {
                        found = true;
                        break;
                    }
                    curr_elem = curr_elem->parent;
                }
                return found;
            }
            case YW_CSS_DIRECT_CHILD_COMBINATOR: {
                /* A > B */
                if (elem->_node.parent == NULL || !(elem->_node.parent->type_flags & YW_DOM_ELEMENT_NODE) || !yw_css_selector_match_element((YW_CSSSelector const *)prev_sel, (YW_GC_PTR(YW_DOMElement))elem->_node.parent))
                {
                    return false;
                }
                return true;
            }
            case YW_CSS_PLUS_COMBINATOR:
            case YW_CSS_TILDE_COMBINATOR:
            case YW_CSS_TWO_BARS_COMBINATOR:
                YW_TODO();
            }
        }
        break;
    }
    case YW_CSS_SELECTOR_NODE_PTR:
        return elem == sel->node_ptr_sel.node_ptr;
    }
    YW_ILLEGAL_VALUE(sel->type);
}

/*******************************************************************************
 *
 * CSS2 9.5 Floats
 *
 * https://www.w3.org/TR/CSS2/visuren.html#floats
 *
 ******************************************************************************/

char const *yw_css_float_str(YW_CSSFloat flo)
{
    switch (flo)
    {
    case YW_CSS_NO_FLOAT:
        return "none";
    case YW_CSS_FLOAT_LEFT:
        return "left";
    case YW_CSS_FLOAT_RIGHT:
        return "right";
    }
    YW_ILLEGAL_VALUE(flo);
}
