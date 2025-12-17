/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_TESTS_H_
#define YW_TESTS_H_

/* clang-format off */
#define YW_ENUMERATE_TESTS(_x)                                                 \
    /* yw_common_tests */                                                      \
    _x(yw_test_grow)                                                           \
    _x(yw_test_shrink_to_fit)                                                  \
    _x(yw_test_list)                                                           \
    _x(yw_test_utf8_next_char)                                                 \
    _x(yw_test_utf8_strlen)                                                    \
    _x(yw_test_utf8_strchr)                                                    \
    _x(yw_test_utf8_to_char32)                                                 \
    _x(yw_test_peek_char)                                                      \
    _x(yw_test_consume_any_char)                                               \
    _x(yw_test_consume_one_of_chars)                                           \
    _x(yw_test_consume_one_of_strs)                                            \
    /* yw_dom_tests */                                                         \
    _x(yw_test_dom_first_child)                                                \
    _x(yw_test_dom_last_child)                                                 \
    _x(yw_test_dom_next_sibling)                                               \
    _x(yw_test_dom_prev_sibling)                                               \
    _x(yw_test_dom_root)                                                       \
    _x(yw_test_dom_index)                                                      \
    _x(yw_test_dom_has_type)                                                   \
    _x(yw_test_dom_is_in_same_tree)                                            \
    _x(yw_test_dom_is_connected)                                               \
    _x(yw_test_dom_child_text)                                                 \
    _x(yw_test_dom_iter)                                                       \
    _x(yw_test_dom_insert)                                                     \
    _x(yw_test_dom_is_element_defined)                                         \
    _x(yw_test_dom_is_element_custom)                                          \
    _x(yw_test_dom_is_element_inside)                                          \
    _x(yw_test_dom_is_element)                                                 \
    _x(yw_test_dom_append_attr)                                                \
    _x(yw_test_dom_attr)
/* clang-format on */

typedef struct YW_TestingContext YW_TestingContext;
struct YW_TestingContext
{
    int failed_counter;
};

#define YW_X(_name) void _name(YW_TestingContext *ctx);
YW_ENUMERATE_TESTS(YW_X)
#undef YW_X

/* WARNING: _got and _expected will get re-evaulated multiple times! */
#define YW_TEST_EXPECT(_got, _item_fmt, _expected)                             \
    do                                                                         \
    {                                                                          \
        if ((_got) != (_expected))                                             \
        {                                                                      \
            printf("FAIL: %s: %s: expected " _item_fmt ", got " _item_fmt      \
                   "\n",                                                       \
                   __func__, #_got, (_expected), (_got));                      \
            yw_failed_test(ctx);                                               \
        }                                                                      \
    } while (0)

#define YW_TEST_EXPECT_STR(_got, _expected)                                    \
    do                                                                         \
    {                                                                          \
        char const *got = _got;                                                \
        char const *expected = _expected;                                      \
        if (got == NULL && expected == NULL)                                   \
        {                                                                      \
            break;                                                             \
        }                                                                      \
        if (got == NULL || (expected != NULL && strcmp(got, expected) != 0))   \
        {                                                                      \
            printf("FAIL: %s: %s: expected %s, got %s\n", __func__, #_got,     \
                   expected, got);                                             \
            yw_failed_test(ctx);                                               \
        }                                                                      \
    } while (0)

#define YW_TEST_EXPECT_ARRAY(_item_type, _got_array, _got_len, _item_fmt, ...) \
    do                                                                         \
    {                                                                          \
        _item_type expected[] = {__VA_ARGS__};                                 \
        _item_type *got_array = (_got_array);                                  \
        int expected_len = YW_SIZEOF_ARRAY(expected);                          \
        int got_len = (_got_len);                                              \
        if (expected_len != got_len)                                           \
        {                                                                      \
            printf("FAIL: %s: %s: expected %d items, got %d\n", __func__,      \
                   #_got_array, expected_len, got_len);                        \
            yw_failed_test(ctx);                                               \
        }                                                                      \
        for (int i = 0; i < expected_len; i++)                                 \
        {                                                                      \
            if (got_array[i] != expected[i])                                   \
            {                                                                  \
                printf("FAIL: %s: %s: expected " _item_fmt                     \
                       " at index %d, got " _item_fmt "\n",                    \
                       __func__, #_got_array, expected[i], i, got_array[i]);   \
                yw_failed_test(ctx);                                           \
                break;                                                         \
            }                                                                  \
        }                                                                      \
    } while (0)

void yw_failed_test(YW_TestingContext *ctx);
void yw_run_all_tests(void);

#endif /* #ifndef YW_TESTS_H_ */
