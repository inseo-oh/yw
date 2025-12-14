/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#ifndef YW_UTILITY_H_
#define YW_UTILITY_H_

#include <stdint.h>

struct yw_utf8_iterator {
    char const *next_src;
    uint32_t codepoint;
    uint8_t bytes_seen;
    uint8_t bytes_needed;
    uint8_t lower_boundary;
    uint8_t upper_boundary;
};

void yw_utf8_iterator_init(struct yw_utf8_iterator *out, char const *str);

/*
 * Returns resulting codepoint, or 0 if end was reached.
 * If an error was encountered, U+FFFD is returned instead.
 */
int32_t yw_utf8_iterator_next(struct yw_utf8_iterator *iter);

#endif /* #ifndef YW_UTILITY_H_ */