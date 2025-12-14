/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_tests.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct yw_inline_formatting_context {
    int dummy;
};

#define SEGMENT_BREAK_CHAR '\n'

/* https://www.w3.org/TR/css-text-3/#white-space-phase-1 */
static char const *apply_whitespace_collapsing(
    char const *str,
    struct yw_inline_formatting_context *ifc) {
    (void)ifc;
    size_t src_size = (strlen(str) + 1) * sizeof(*str);
    char *res_str = malloc(src_size);
    if (res_str == NULL) {
        goto fail;
    }

    /*
     * TODO: Add support for white-space: pre, white-space:pre-wrap, or
     * white-space: break-spaces
     */

    /***************************************************************************
     * Ignore collapsible spaces and tabs immediately following/preceding
     * segment break.
     *
     * "foo   \n   bar" --> "foo\nbar"
     **************************************************************************/

    /* TODO */
    goto out;
fail:
    free(res_str);
    res_str = NULL;
out:
    return res_str;
}

int main() {
    printf("hello, world!\n");
    yw_run_all_tests();
    return 0;
}