/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_LAYOUT_H_
#define YW_LAYOUT_H_

struct yw_inline_formatting_context
{
    char const *written_text;
};

/*
 * CSS Text Module Level 3 - 4.1.1
 * Returns NULL when there's not enough memory.
 */
char *yw_apply_whitespace_collapsing(char const *str,
                                     struct yw_inline_formatting_context *ifc);

/*
 * CSS Text Module Level 3 - 4.1.3
 * Returns NULL when there's not enough memory.
 */
char *yw_apply_segment_break_transform(char const *str);

#endif /* #ifndef YW_LAYOUT_H_ */