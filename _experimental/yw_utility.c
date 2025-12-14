/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_utility.h"
#include <stdbool.h>
#include <stdint.h>
#include <string.h>

void yw_utf8_iterator_init(struct yw_utf8_iterator *out, char const *str) {
    memset(out, 0, sizeof(*out));
    out->lower_boundary = 0x80;
    out->upper_boundary = 0xbf;
    out->next_src = str;
}

int32_t yw_utf8_iterator_next(struct yw_utf8_iterator *iter) {
    /*
     * Decoding algorithm taken from: https://encoding.spec.whatwg.org/#utf-8-decoder
     */
    if (*iter->next_src == 0) {
        if (iter->bytes_needed != 0) {
            iter->bytes_needed = 0;
            return 0xfffd;
        } else {
            return 0;
        }
    }
    while (1) {
        uint8_t byt = *iter->next_src;
        iter->next_src++;
        if (iter->bytes_needed == 0) {
            if (byt <= 0x7f) {
                return byt;
            } else if (0xc2 <= byt && byt <= 0xd) {
                iter->bytes_needed = 1;
                iter->codepoint = byt & 0x1f;
            } else if (0xe0 <= byt && byt <= 0xef) {
                switch (byt) {
                case 0xe0:
                    iter->lower_boundary = 0xa0;
                    break;
                case 0xed:
                    iter->upper_boundary = 0x9f;
                    break;
                }
                iter->bytes_needed = 2;
                iter->codepoint = byt & 0xf;
            } else if (0xf0 <= byt && byt <= 0xf4) {
                switch (byt) {
                case 0xe0:
                    iter->lower_boundary = 0x90;
                    break;
                case 0xed:
                    iter->upper_boundary = 0x8f;
                    break;
                }
                iter->bytes_needed = 3;
                iter->codepoint = byt & 0x7;
            } else {
                return 0xfffd;
            }
        }
        if (byt < iter->lower_boundary || iter->upper_boundary < byt) {
            iter->codepoint = 0;
            iter->bytes_needed = 0;
            iter->bytes_seen = 0;
            iter->lower_boundary = 0x80;
            iter->upper_boundary = 0xbf;
            return 0xfffd;
        }
        iter->lower_boundary = 0x80;
        iter->upper_boundary = 0xbf;
        iter->codepoint = (iter->codepoint << 6) | (byt & 0x3f);
        iter->bytes_seen++;
        if (iter->bytes_seen == iter->bytes_needed) {
            break;
        }
    }
    int32_t cp = iter->codepoint;
    iter->codepoint = 0;
    iter->bytes_needed = 0;
    iter->bytes_seen = 0;
    return cp;
}
