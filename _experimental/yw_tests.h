/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#ifndef YW_TESTS_H_
#define YW_TESTS_H_

struct yw_testing_context {
    int failed_counter;
};

#define YW_ENUMERATE_TESTS(_x) \
    _x(yw_test_utf8_iterator)

#define YW_X(_name) void _name(struct yw_testing_context *ctx);
YW_ENUMERATE_TESTS(YW_X)
#undef YW_X

void yw_failed_test(struct yw_testing_context *ctx);
void yw_run_all_tests();

#endif /* #ifndef YW_TESTS_H_ */