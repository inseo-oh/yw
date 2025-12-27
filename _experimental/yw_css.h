/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_CSS_H_
#define YW_CSS_H_

#include <stdbool.h>
#include <stdint.h>

/*******************************************************************************
 *
 * CSS Syntax
 *
 * https://www.w3.org/TR/css-syntax-3/
 *
 ******************************************************************************/

typedef union YW_CSSToken YW_CSSToken;

void yw_token_deinit(YW_CSSToken *tk);

typedef struct YW_CSSTokenStream
{
    YW_CSSToken *tokens;
    int tokens_len;
    int cursor;
} YW_CSSTokenStream;

/* NOTE: out->tokens must be freed after using them */
bool yw_css_tokenize(YW_CSSTokenStream *out, uint8_t const *bytes, int bytes_len, const char *source_name);

/*******************************************************************************
 *
 * CSS Values and Units
 *
 * https://www.w3.org/TR/css-values-3/
 *
 ******************************************************************************/

typedef enum
{
    /***************************************************************************
     * Relative lengths
     * https://www.w3.org/TR/css-values-3/#relative-lengths
     **************************************************************************/

    YW_CSS_EM,
    YW_CSS_EX,
    YW_CSS_CH,
    YW_CSS_REM,
    YW_CSS_VW,
    YW_CSS_VH,
    YW_CSS_VMIN,
    YW_CSS_VMAX,

    /***************************************************************************
     * Absolute lengths
     * https://www.w3.org/TR/css-values-3/#absolute-lengths
     **************************************************************************/

    YW_CSS_CM,
    YW_CSS_MM,
    YW_CSS_Q,
    YW_CSS_PC,
    YW_CSS_PT,
    YW_CSS_PX,
} YW_CSSLengthUnit;

char const *yw_css_length_unit_str(YW_CSSLengthUnit unit);

/* https://www.w3.org/TR/css-values-3/#length-value */
typedef struct YW_CSSLength
{
    double value;
    YW_CSSLengthUnit unit;
} YW_CSSLength;

#define YW_CSS_LENGTH_FMT "%g%s"
#define YW_CSS_LENGTH_FMT_ARGS(_l) (_l).value, yw_css_length_unit_str((_l).unit)

/*
 * font_size is used as base when font size relative units like
 * YW_CSS_EM is used.
 */
double yw_css_length_to_px(YW_CSSLength const *len, double font_size);

double yw_css_percentage_to_px(double per, double container_size);
void yw_css_percentage_to_length(YW_CSSLength *out, double per, double container_size);

typedef struct YW_CSSLengthOrPercentage
{
    union {
        YW_CSSLength length;
        double percentage;
    } value;
    bool is_percentage;
} YW_CSSLengthOrPercentage;

/*
 * font_size is used as base when font size relative units like YW_CSS_EM
 * is used.
 */
double yw_css_length_or_percentage_to_px(YW_CSSLengthOrPercentage const *len_or_per, double font_size, double container_size);

bool yw_css_parse_number(double *out, YW_CSSTokenStream *ts);
bool yw_css_parse_percentage(double *out, YW_CSSTokenStream *ts);

/* YW_ALLOW_ZERO_SHORTHAND should not be used if the property also accepts
 * number token. (e.g. line-height)
 * In such case, 0 should be parsed as <number 0>, not <length 0>.
 */
typedef enum
{
    YW_NO_ALLOW_ZERO_SHORTHAND,
    YW_ALLOW_ZERO_SHORTHAND,
} YW_AllowZeroShorthand;

bool yw_css_parse_length(YW_CSSLength *out, YW_CSSTokenStream *ts, YW_AllowZeroShorthand allow_zero_shorthand);
bool yw_css_parse_length_or_percentage(YW_CSSLengthOrPercentage *out, YW_CSSTokenStream *ts, YW_AllowZeroShorthand allow_zero_shorthand);

/*******************************************************************************
 *
 * CSS Backgrounds and Borders
 *
 * https://www.w3.org/TR/css-backgrounds-3/
 *
 ******************************************************************************/

typedef enum
{
    YW_CSS_NO_LINE,
    YW_CSS_HIDDEN_LINE,
    YW_CSS_DOTTED_LINE,
    YW_CSS_DASHED_LINE,
    YW_CSS_SOLID_LINE,
    YW_CSS_DOUBLE_LINE,
    YW_CSS_GROOVE_LINE,
    YW_CSS_RIDGE_LINE,
    YW_CSS_INSET_LINE,
    YW_CSS_OUTSET_LINE,
} YW_CSSLineStyle;

char const *yw_css_line_style_str(YW_CSSLineStyle style);

/* https://www.w3.org/TR/css-backgrounds-3/#typedef-line-width */
enum
{
    YW_CSS_LINE_WIDTH_THIN = 1,
    YW_CSS_LINE_WIDTH_MEDIUM = 3,
    YW_CSS_LINE_WIDTH_THICK = 5,
};

bool yw_css_parse_line_style(YW_CSSLineStyle *out, YW_CSSTokenStream *ts);
bool yw_css_parse_line_width(YW_CSSLength *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Box Model
 *
 * https://www.w3.org/TR/css-box-3/
 *
 ******************************************************************************/

typedef struct YW_CSSMargin
{
    YW_CSSLengthOrPercentage value;
    bool is_auto;
} YW_CSSMargin;

bool yw_css_parse_margin(YW_CSSMargin *out, YW_CSSTokenStream *ts);
bool yw_css_parse_padding(YW_CSSLengthOrPercentage *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Color
 *
 * https://www.w3.org/TR/css-color-4/
 *
 ******************************************************************************/

typedef uint32_t YW_CSSRgba;
typedef union YW_CSSColor YW_CSSColor;

#define YW_CSS_RGBA(_r, _g, _b, _a) (((uint32_t)((_r) & 0xff) << 24) | ((uint32_t)((_g) & 0xff) << 16) | ((uint32_t)((_b) & 0xff) << 8) | (uint32_t)((_a) & 0xff))
#define YW_CSS_RGB(_r, _g, _b) YW_CSS_RGBA((_r), (_g), (_b), 255)

#define YW_CSS_RED(_rgba) (((uint32_t)(_rgba) >> 24) & 0xff)
#define YW_CSS_GREEN(_rgba) (((uint32_t)(_rgba) >> 16) & 0xff)
#define YW_CSS_BLUE(_rgba) (((uint32_t)(_rgba) >> 8) & 0xff)
#define YW_CSS_ALPHA(_rgba) ((uint32_t)(_rgba) & 0xff)

YW_CSSRgba yw_css_color_from_name(char const *name);

typedef enum
{
    YW_CSS_RGB_COLOR,        /* rgb(), rgba(), hex colors, named colors */
    YW_CSS_CURRENT_COLOR,    /* currentColor */
    YW_CSS_HSL_COLOR,        /* hsl(), hsla()*/
    YW_CSS_HWB_COLOR,        /* hwb() */
    YW_CSS_LAB_COLOR,        /* lab() */
    YW_CSS_LCH_COLOR,        /* lch() */
    YW_CSS_OKLAB_COLOR,      /* oklab() */
    YW_CSS_OKLCH_COLOR,      /* oklch() */
    YW_CSS_COLOR_FUNC_COLOR, /* color() */
} YW_CSSColorType;

union YW_CSSColor {
    YW_CSSColorType type;
    struct
    {
        YW_CSSColorType type; /* YW_CSS_RGB */
        YW_CSSRgba rgba;
    } rgb;
};

void yw_css_color_from_rgba(YW_CSSColor *out, YW_CSSRgba rgba);

/* NOTE: currentColor must be handled by caller. */
YW_CSSRgba yw_css_color_to_rgba(YW_CSSColor const *color);

bool yw_css_parse_color(YW_CSSColor *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Display
 *
 * https://www.w3.org/TR/css-display-3/
 *
 ******************************************************************************/

/*
 * Display mode bit flags
 *  - Bit 0 ~ 1: Outer mode
 *  - Bit 4 ~ 6: Inner mode
 *  - Bit 8 ~ 11 : Special mode
 */
typedef uint16_t YW_CSSDisplay;

/* Bit 0 ~ 1: Outer mode ******************************************************/
#define YW_CSS_DISPLAY_OUTER_MODE_MASK (0x3 << 0)
#define YW_CSS_DISPLAY_BLOCK (0x0 << 0)
#define YW_CSS_DISPLAY_INLINE (0x1 << 0)
#define YW_CSS_DISPLAY_RUN_IN (0x2 << 0)
/* Bit 4 ~ 6: Inner mode ******************************************************/
#define YW_CSS_DISPLAY_INNER_MODE_MASK (0x7 << 4)
#define YW_CSS_DISPLAY_FLOW (0x0 << 4)
#define YW_CSS_DISPLAY_FLOW_ROOT (0x1 << 4)
#define YW_CSS_DISPLAY_TABLE (0x2 << 4)
#define YW_CSS_DISPLAY_FLEX (0x3 << 4)
#define YW_CSS_DISPLAY_GRID (0x4 << 4)
#define YW_CSS_DISPLAY_RUBY (0x5 << 4)
/*  Bit 8 ~ 11 : Special mode *************************************************/
/* NOTE: Outer mode is ignored for below modes */
#define YW_CSS_DISPLAY_SPECIAL_NONE (0x0 << 8)
#define YW_CSS_DISPLAY_TABLE_ROW_GROUP (0x1 << 8)
#define YW_CSS_DISPLAY_TABLE_HEADER_GROUP (0x2 << 8)
#define YW_CSS_DISPLAY_TABLE_FOOTER_GROUP (0x3 << 8)
#define YW_CSS_DISPLAY_TABLE_ROW (0x4 << 8)
#define YW_CSS_DISPLAY_TABLE_CELL (0x5 << 8)
#define YW_CSS_DISPLAY_TABLE_COLUMN_GROUP (0x6 << 8)
#define YW_CSS_DISPLAY_TABLE_COLUMN (0x7 << 8)
#define YW_CSS_DISPLAY_TABLE_CAPTION (0x8 << 8)
#define YW_CSS_DISPLAY_RUBY_BASE (0x9 << 8)
#define YW_CSS_DISPLAY_RUBY_TEXT (0xa << 8)
#define YW_CSS_DISPLAY_RUBY_BASE_CONTAINER (0xb << 8)
#define YW_CSS_DISPLAY_RUBY_TEXT_CONTAINER (0xc << 8)
#define YW_CSS_DISPLAY_LIST_ITEM (0xd << 8)
/* NOTE: Outer/inner mode is ignored for below modes */
#define YW_CSS_DISPLAY_CONTENTS (0xe << 9)
#define YW_CSS_DISPLAY_NONE (0xf << 9)

typedef enum
{
    YW_CSS_VISIBLE,
    YW_CSS_HIDDEN,
    YW_CSS_COLLAPSE,
} YW_CSSVisibility;

char const *yw_css_visibility_str(YW_CSSVisibility vis);

bool yw_css_parse_display(YW_CSSDisplay *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS2 9.5 Floats
 *
 * https://www.w3.org/TR/CSS2/visuren.html#floats
 *
 ******************************************************************************/

typedef enum
{
    YW_CSS_NO_FLOAT,
    YW_CSS_FLOAT_LEFT,
    YW_CSS_FLOAT_RIGHT,
} YW_CSSFloat;

char const *yw_css_float_str(YW_CSSFloat flo);

bool yw_css_parse_float(YW_CSSFloat *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Fonts
 *
 * https://www.w3.org/TR/css-fonts-3
 *
 ******************************************************************************/

typedef enum
{
    YW_CSS_NON_GENERIC_FONT_FAMILY = -1,
    YW_CSS_SERIF = 0,
    YW_CSS_SANS_SERIF,
    YW_CSS_CURSIVE,
    YW_CSS_FANTASY,
    YW_CSS_MONOSPACE,
} YW_CSSGenericFontFamily;

char const *yw_css_generic_font_family_str(YW_CSSGenericFontFamily fam);

typedef struct YW_CSSFontFamily
{
    /*
     * Only valid if family is YW_CSS_NON_GENERIC_FONT_FAMILY (otherwise NULL)
     */
    char *name;

    YW_CSSGenericFontFamily family;
} YW_CSSFontFamily;

typedef struct YW_CSSFontFamilies
{
    YW_CSSFontFamily *items;
    int len;
} YW_CSSFontFamilies;

/* NOTE: YW_CSSFontWeight can be any integer between 0 and 1000(inclusive) */
typedef enum
{
    YW_CSS_NORMAL_FONT_WEIGHT = 400,
    YW_CSS_BOLD = 700,
} YW_CSSFontWeight;

typedef enum
{
    YW_CSS_ULTRA_CONDENSED,
    YW_CSS_EXTRA_CONDENSED,
    YW_CSS_CONDENSED,
    YW_CSS_SEMI_CONDENSED,
    YW_CSS_NORMAL_FONT_STRETCH,
    YW_CSS_SEMI_EXPANDED,
    YW_CSS_EXPANDED,
    YW_CSS_EXTRA_EXPANDED,
    YW_CSS_ULTRA_EXPANDED,
} YW_CSSFontStretch;

char const *yw_css_generic_font_stretch_str(YW_CSSFontStretch str);

typedef enum
{
    YW_CSS_NORMAL_FONT_STYLE,
    YW_CSS_ITALIC,
    YW_CSS_OBLIQUE,
} YW_CSSFontStyle;

char const *yw_css_generic_font_style_str(YW_CSSFontStyle sty);

/* XXX: Let user choose this size! */
#define YW_CSS_PREFERRED_FONT_SIZE 14.0

typedef enum
{
    YW_CSS_LENGTH_FONT_SIZE,

    /* Absolute sizes *********************************************************/

    YW_CSS_XX_SMALL,
    YW_CSS_X_SMALL,
    YW_CSS_SMALL,
    YW_CSS_MEDIUM_FONT_SIZE,
    YW_CSS_LARGE,
    YW_CSS_X_LARGE,
    YW_CSS_XX_LARGE,

    /* Relative sizes *********************************************************/

    YW_CSS_LARGER,
    YW_CSS_SMALLER,
} YW_CSSFontSizeType;

typedef struct YW_CSSFontSize
{
    /* Only valid when type is YW_CSS_LENGTH_FONT_SIZE */
    YW_CSSLengthOrPercentage size;

    YW_CSSFontSizeType type;
} YW_CSSFontSize;

/*
 * font_size is used as base when font size relative units like
 * YW_CSS_EM is used.
 */
double yw_css_font_size_to_px(YW_CSSFontSize const *sz, double font_size, double parent_font_size);

/* Returned families in output must be freed using free() */
bool yw_css_parse_font_family(YW_CSSFontFamilies *out, YW_CSSTokenStream *ts);
bool yw_css_parse_font_weight(YW_CSSFontWeight *out, YW_CSSTokenStream *ts);
bool yw_css_parse_font_stretch(YW_CSSFontStretch *out, YW_CSSTokenStream *ts);
bool yw_css_parse_font_style(YW_CSSFontStyle *out, YW_CSSTokenStream *ts);
bool yw_css_parse_font_size(YW_CSSFontSize *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Selectors
 *
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/
 *
 ******************************************************************************/

typedef struct YW_CSSWqName
{
    char *ns_prefix; /* Namespace prefix - May be nil */
    char *ident;
} YW_CSSWqName;

typedef enum
{
    YW_CSS_SELECTOR_ATTR,
    YW_CSS_SELECTOR_CLASS,
    YW_CSS_SELECTOR_ID,
    YW_CSS_SELECTOR_TYPE,
    YW_CSS_SELECTOR_UNIVERSAL,
    YW_CSS_SELECTOR_COMPOUND,
    YW_CSS_SELECTOR_PSEUDO_CLASS,
    YW_CSS_SELECTOR_COMPLEX,
    YW_CSS_SELECTOR_NODE_PTR,
} YW_CSSSelectorType;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#class-selector
 */
typedef struct YW_CSSClassSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_CLASS */
    char *class_name;
} YW_CSSClassSelector;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#id-selector
 */
typedef struct YW_CSSIdSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_ID */
    char *id;
} YW_CSSIdSelector;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#type-selector
 */
typedef struct YW_CSSTypeSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_TYPE */
    YW_CSSWqName name;
} YW_CSSTypeSelector;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#the-universal-selector
 */
typedef struct YW_CSSUniversalSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_TYPE */
    char *ns_prefix;         /* Namespace prefix - May be nil */
} YW_CSSUniversalSelector;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#pseudo-class
 */
typedef struct YW_CSSPseudoClassSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_PSEUDO_CLASS */
    char *name;
    /* TODO: Arguments */
} YW_CSSPseudoClassSelector;

typedef enum
{
    /* TODO: Give these better names */

    YW_CSS_NO_VALUE_MATCH,        /* [attr] */
    YW_CSS_VALUE_EQUALS,          /* [attr=value] */
    YW_CSS_VALUE_TILDE_EQUALS,    /* [attr~=value] */
    YW_CSS_VALUE_BAR_EQUALS,      /* [attr|=value] */
    YW_CSS_VALUE_CARET_EQUALS,    /* [attr^=value] */
    YW_CSS_VALUE_DOLLAR_EQUALS,   /* [attr$=value] */
    YW_CSS_VALUE_ASTERISK_EQUALS, /* [attr*=value] */
} YW_CSSValueMatchType;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#attribute-selector
 */
typedef struct YW_CSSAttrSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_ATTR */
    YW_CSSValueMatchType value_match_type;
    YW_CSSWqName attr_name;
    char *attr_value;
    bool is_case_sensitive;
} YW_CSSAttrSelector;

typedef struct YW_CSSCompundSelectorPseudoItem
{
    union YW_CSSSelector *pseudo_elem_sel; /* Pseudo element selector */
    union YW_CSSSelector *class_sels;      /* Pseudo class selectors */
    int class_sels_len;
} YW_CSSCompundSelectorPseudoItem;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#compound
 */
typedef struct YW_CSSCompoundSelector
{
    YW_CSSSelectorType type;        /* YW_CSS_SELECTOR_COMPOUND */
    union YW_CSSSelector *type_sel; /* Type selector. May be NULL */
    union YW_CSSSelector *subclass_sels;
    YW_CSSCompundSelectorPseudoItem *pseudo_items;
    int subclass_sels_len;
    int pseudo_items_len;
} YW_CSSCompoundSelector;

typedef enum
{
    YW_CSS_CHILD_COMBINATOR,        /* A B */
    YW_CSS_DIRECT_CHILD_COMBINATOR, /* A > B */

    /* TODO: Give these better names */

    YW_CSS_PLUS_COMBINATOR,     /* A + B */
    YW_CSS_TILDE_COMBINATOR,    /* A ~ B */
    YW_CSS_TWO_BARS_COMBINATOR, /* A || B */
} YW_CSSCombinator;

typedef struct YW_CSSComplexSelectorRest
{
    union YW_CSSSelector *selector; /* Compound selector */
    YW_CSSCombinator combinator;
} YW_CSSComplexSelectorRest;

/*
 * https://www.w3.org/TR/2022/WD-selectors-4-20221111/#complex
 */
typedef struct YW_CSSComplexSelector
{
    YW_CSSSelectorType type;    /* YW_CSS_SELECTOR_COMPLEX */
    union YW_CSSSelector *base; /* Compound selector */
    YW_CSSComplexSelectorRest *rests;
    int rests_len;
} YW_CSSComplexSelector;

/*
 * This selector matches directly using pointer to DOM node.
 * This is only for internal use, and not part of CSS spec.
 */
typedef struct YW_CSSNodePtrSelector
{
    YW_CSSSelectorType type; /* YW_CSS_SELECTOR_NODE_PTR */
    void *node_ptr;
} YW_CSSNodePtrSelector;

typedef union YW_CSSSelector {
    YW_CSSSelectorType type;

    YW_CSSClassSelector class_sel;
    YW_CSSIdSelector id_sel;
    YW_CSSTypeSelector type_sel;
    YW_CSSUniversalSelector universal_sel;
    YW_CSSPseudoClassSelector pseudo_class_sel;
    YW_CSSAttrSelector attr_sel;
    YW_CSSCompoundSelector compound_sel;
    YW_CSSComplexSelector complex_sel;
    YW_CSSNodePtrSelector node_ptr_sel;
} YW_CSSSelector;

void yw_css_wq_name_deinit(YW_CSSWqName *name);
void yw_css_selector_deinit(YW_CSSSelector *sel);
bool yw_css_parse_selector_list(YW_CSSSelector **sels_out, int *len_out, YW_CSSTokenStream *ts);
bool yw_css_parse_selector(YW_CSSSelector **sels_out, int *len_out, uint8_t const *bytes, int bytes_len, const char *source_name);

/*******************************************************************************
 *
 * CSS Sizing
 *
 * https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/
 *
 ******************************************************************************/

typedef enum
{
    YW_CSS_NO_SIZE,   /* none */
    YW_CSS_AUTO_SIZE, /* auto */
    YW_CSS_MIN_CONTENT,
    YW_CSS_MAX_CONTENT,
    YW_CSS_FIT_CONTENT,
    YW_CSS_MANUAL_SIZE,
} YW_CSSSizeType;

typedef struct YW_CSSSize
{
    /* Only valid if type is YW_CSS_MANUAL_SIZE */
    YW_CSSLengthOrPercentage size;

    YW_CSSSizeType type;
} YW_CSSSize;

bool yw_css_parse_size_or_auto(YW_CSSSize *out, YW_CSSTokenStream *ts);
bool yw_css_parse_size_or_none(YW_CSSSize *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Text
 *
 * https://www.w3.org/TR/css-text-3
 *
 ******************************************************************************/

/*
 * Transform mode bit flags
 *  - Bit 0 ~ 1: Caps transform
 *  - Bit 4 ~ 5: Other flags
 */
typedef uint8_t YW_CSSTextTransform;

/* Bit 0 ~ 1: Caps transform */
#define YW_CSS_TEXT_TRANSFORM_CAPS_MODE_MASK (0x3 << 0)
#define YW_CSS_TEXT_TRANSFORM_ORIGINAL_CAPS (0x0 << 0)
#define YW_CSS_TEXT_TRANSFORM_CAPITALIZE (0x1 << 0)
#define YW_CSS_TEXT_TRANSFORM_UPPERCASE (0x2 << 0)
#define YW_CSS_TEXT_TRANSFORM_LOWERCASE (0x3 << 0)
/* Bit 4 ~ 5: Other flags */
#define YW_CSS_TEXT_TRANSFORM_FULL_WIDTH (1 << 4)
#define YW_CSS_TEXT_TRANSFORM_FULL_SIZE_KANA (1 << 5)

bool yw_css_parse_text_transform(YW_CSSTextTransform *out, YW_CSSTokenStream *ts);

/*******************************************************************************
 *
 * CSS Text Decoration
 *
 * https://www.w3.org/TR/css-text-decor-3
 *
 ******************************************************************************/

/* https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-line */
typedef uint8_t YW_CSSTextDecorationLine;
#define YW_CSS_TEXT_DECORATION_UNDERLINE (1 << 0)
#define YW_CSS_TEXT_DECORATION_OVERLINE (1 << 1)
#define YW_CSS_TEXT_DECORATION_LINE_THROUGH (1 << 2)
#define YW_CSS_TEXT_DECORATION_BLINK (1 << 3)

/* Bit 4 ~ 6: Decoration style */
/* https://www.w3.org/TR/css-text-decor-3/#propdef-text-decoration-style */

typedef enum
{
    YW_CSS_TEXT_DECORATION_SOLID,
    YW_CSS_TEXT_DECORATION_DOUBLE,
    YW_CSS_TEXT_DECORATION_DOTTED,
    YW_CSS_TEXT_DECORATION_DASHED,
    YW_CSS_TEXT_DECORATION_WAVY,
} YW_CSSTextDecorationStyle;

bool yw_css_parse_text_decoration_line(YW_CSSTextDecorationLine *out, YW_CSSTokenStream *ts);
bool yw_css_parse_text_decoration_style(YW_CSSTextDecorationStyle *out, YW_CSSTokenStream *ts);

#endif /* #ifndef YW_CSS_H_ */
