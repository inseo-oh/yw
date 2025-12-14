/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_tests.h"

#include "yw_utility.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

void yw_test_utf8_iterator(struct yw_testing_context *ctx) {
#define YW_RUN_TEST(_name, _input, ...)                                        \
    do {                                                                       \
        struct yw_utf8_iterator iter;                                          \
        int32_t expected[] = {__VA_ARGS__};                                    \
        int dest_len = sizeof(expected) / sizeof(*expected);                   \
        yw_utf8_iterator_init(&iter, _input);                                  \
        for (int i = 0; i < dest_len; i++) {                                   \
            int32_t res = yw_utf8_iterator_next(&iter);                        \
            if (res != expected[i]) {                                          \
                printf(                                                        \
                    "FAIL: %s[%s]: expected U+%04X at index %d, got U+%04X\n", \
                    __func__, _name, expected[i], i, res);                     \
                yw_failed_test(ctx);                                           \
                break;                                                         \
            }                                                                  \
        }                                                                      \
    } while (0)

    YW_RUN_TEST("Simple ASCII", "\x30\x31\x32\x33\x7e",
                '0', '1', '2', '3');
    YW_RUN_TEST("Two byte characters",
                "\xc2\xa0\xde\xb1",
                0x00a0, 0x07b1);
    YW_RUN_TEST("Three byte characters",
                "\xe0\xa4\x80\xed\x9f\xbb\xef\xad\x8f",
                0x0900, 0xd7fb, 0xfb4f);
    YW_RUN_TEST("Four byte characters",
                "\xf0\x90\x91\x90\xf0\x9f\x83\xb5\xf4\x81\x8a\x8f",
                0x10450, 0x1f0f5, 0x10128f);

#undef YW_RUN_TEST
}