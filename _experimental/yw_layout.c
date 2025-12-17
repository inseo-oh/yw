/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_layout.h"
#include "yw_common.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define YW_SEGMENT_BREAK_CHAR '\n'

char *yw_apply_segment_break_transform(char const *str)
{
    /* https://www.w3.org/TR/css-text-3/#line-break-transform */

    char *res_str = NULL;

    {
        /***********************************************************************
         * Remove segment breaks immediately following another.
         *
         * "foo\n\nbar" -> "foo\nbar"
         **********************************************************************/
        int res_cap = 0;
        int res_len = 0;
        char const *next_chr = str;
        while (*next_chr != '\0')
        {
            res_str = YW_GROW(char, &res_cap, &res_len, res_str);
            char got = *next_chr;
            next_chr++;
            res_str[res_len - 1] = got;
            if (got == YW_SEGMENT_BREAK_CHAR)
            {
                char const *last_space =
                    strrchr(next_chr, YW_SEGMENT_BREAK_CHAR);
                if (last_space != NULL)
                {
                    next_chr = last_space + 1;
                }
            }
        }
        YW_SHRINK_TO_FIT(char, &res_cap, res_len, res_str);
    }

    /***************************************************************************
     * Turn remaining segment breaks into spaces.
     *
     * "foo\nbar\njaz" -> "foo bar jaz"
     **************************************************************************/
    for (int i = 0; res_str[i] != '\0'; i++)
    {
        if (res_str[i] == YW_SEGMENT_BREAK_CHAR)
        {
            res_str[i] = ' ';
        }
    }

    return res_str;
}

char *yw_apply_whitespace_collapsing(char const *str,
                                     YW_InlineFormattingContext *ifc)
{
    /* https://www.w3.org/TR/css-text-3/#white-space-phase-1 */

    (void)ifc;
    char *res_str = NULL;
    {
        int res_cap = 0;
        int res_len = 0;

        /*
         * TODO: Add support for white-space: pre, white-space:pre-wrap, or
         * white-space: break-spaces
         */

        /***********************************************************************
         * Ignore collapsible spaces and tabs immediately following/preceding
         * segment break.
         *
         * "foo   \n   bar" --> "foo\nbar"
         **********************************************************************/
        char const *next_chr = str;
        while (*next_chr != '\0')
        {
            char chr = *next_chr;
            next_chr++;
            if (chr == ' ' || chr == '\t')
            {
                /*
                 * Check if segment break immediately follows sequence of
                 * space or tab.
                 */
                char const *inner_next_chr = next_chr;
                char const *segment_break = NULL;
                while (1)
                {
                    char const *curr_next_chr = inner_next_chr;
                    char got = *inner_next_chr;
                    inner_next_chr++;
                    if (got == YW_SEGMENT_BREAK_CHAR)
                    {
                        segment_break = curr_next_chr;
                        break;
                    }
                    else if (got != ' ' && got != '\t')
                    {
                        break;
                    }
                }
                if (segment_break != NULL)
                {
                    /*
                     * We found space/tab followed by segment break.
                     * Ignore those characters.
                     */
                    next_chr = segment_break;
                    continue;
                }
            }
            res_str = YW_GROW(char, &res_cap, &res_len, res_str);
            res_str[res_len - 1] = chr;

            if (chr == YW_SEGMENT_BREAK_CHAR)
            {
                /*
                 * Check if segment break is immediately followed by sequence of
                 * space or tab.
                 */
                char const *inner_next_chr = next_chr;
                char const *first_non_space_or_tab = NULL;
                while (1)
                {
                    char const *curr_next_chr = inner_next_chr;
                    char got = *inner_next_chr;
                    inner_next_chr++;
                    if (got != ' ' && got != '\t')
                    {
                        first_non_space_or_tab = curr_next_chr;
                        break;
                    }
                }
                if (first_non_space_or_tab != NULL)
                {
                    /*
                     * Skip to next first non-space / tab.
                     */
                    next_chr = first_non_space_or_tab;
                    continue;
                }
            }
        }

        YW_SHRINK_TO_FIT(char, &res_cap, res_len, res_str);
    }
    /***************************************************************************
     * Transform segment breaks according to segment break transform rules.
     *
     * "foo\n\nbar" -> "foo\nbar"
     **************************************************************************/
    {
        char *new_res_str = yw_apply_segment_break_transform(res_str);
        free(res_str);
        res_str = new_res_str;
    }

    /***************************************************************************
     * Replace tabs with spaces.
     *
     * "foo\t\tbar" -> "foo  bar"
     **************************************************************************/
    for (int i = 0; res_str[i] != '\0'; i++)
    {
        if (res_str[i] == '\t')
        {
            res_str[i] = ' ';
        }
    }

    /***************************************************************************
     * Ignore any space following the another, including the ones outside of
     * current text, as long as it's part of the same IFC.
     * "foo   bar" -> "foo bar"
     *
     * TODO: CSS says these extra sapces don't have zero-advance width, and thus
     *       invisible, but still retains its soft wrap opportunity, if any.
     **************************************************************************/
    if (ifc->written_text != NULL && ifc->written_text[0] != '\0' &&
        ifc->written_text[strlen(ifc->written_text) - 1] == ' ')
    {
        while (res_str[0] == ' ')
        {
            res_str++;
        }
    }

    char *old_res_str = res_str;
    res_str = NULL;
    {
        int res_cap = 0;
        int res_len = 0;
        char const *next_chr = old_res_str;
        while (*next_chr != '\0')
        {
            res_str = YW_GROW(char, &res_cap, &res_len, res_str);
            char got = *next_chr;
            next_chr++;
            res_str[res_len - 1] = got;
            if (got == ' ')
            {
                char const *last_space = strrchr(next_chr, ' ');
                if (last_space != NULL)
                {
                    next_chr = last_space + 1;
                }
            }
        }
        YW_SHRINK_TO_FIT(char, &res_cap, res_len, res_str);
    }
    return res_str;
}
