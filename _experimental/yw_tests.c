/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_tests.h"
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

void yw_failed_test(YW_TestingContext *ctx)
{
    ctx->failed_counter++;
}

void yw_run_all_tests(void)
{
    YW_TestingContext ctx;
    memset(&ctx, 0, sizeof(ctx));

#define YW_X(_name)                                                            \
    printf("Running test: %s\n", #_name);                                      \
    _name(&ctx);

    YW_ENUMERATE_TESTS(YW_X);
#undef YW_X

    if (ctx.failed_counter != 0)
    {
        printf("%d failed tests\n", ctx.failed_counter);
    }
    else
    {
        printf("ALL TESTS PASSED\n");
    }
}
