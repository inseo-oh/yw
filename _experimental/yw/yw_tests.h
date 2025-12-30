/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#ifndef YW_TESTS_H_
#define YW_TESTS_H_
#include <stdio.h>

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
    /* yw_css_tests */                                                         \
    _x(yw_test_css_parse_number)                                               \
    _x(yw_test_css_parse_length)                                               \
    _x(yw_test_css_parse_percentage)                                           \
    _x(yw_test_css_parse_line_style)                                           \
    _x(yw_test_css_parse_line_width)                                           \
    _x(yw_test_css_parse_margin)                                               \
    _x(yw_test_css_parse_padding)                                              \
    _x(yw_test_css_parse_color)                                                \
    _x(yw_test_css_parse_display)                                              \
    _x(yw_test_css_parse_float)                                                \
    _x(yw_test_css_parse_font_family)                                          \
    _x(yw_test_css_parse_font_weight)                                          \
    _x(yw_test_css_parse_font_stretch)                                         \
    _x(yw_test_css_parse_font_style)                                           \
    _x(yw_test_css_parse_font_size)                                            \
    _x(yw_test_css_parse_selector)                                             \
    _x(yw_test_css_parse_size_or_auto)                                         \
    _x(yw_test_css_parse_size_or_none)                                         \
    _x(yw_test_css_parse_text_transform)                                       \
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
    _x(yw_test_dom_append_attr_to_element)                                     \
    _x(yw_test_dom_attr_of_element)                                            \
    /* yw_encoding_tests */                                                    \
    _x(yw_test_bom_sniff)                                                      \
    _x(yw_test_encoding_from_label)                                            \
    _x(yw_test_io_queue_item_list_to_items)                                    \
    _x(yw_test_io_queue_from_items)                                            \
    _x(yw_test_io_queue_read_one)                                              \
    _x(yw_test_io_queue_read)                                                  \
    _x(yw_test_io_queue_peek)                                                  \
    _x(yw_test_io_queue_push_one)                                              \
    _x(yw_test_io_queue_push)                                                  \
    _x(yw_test_io_queue_restore_one)                                           \
    _x(yw_test_io_queue_restore)                                               \
    _x(yw_test_utf8_decoder)                                                   \
    /* yw_json_tests */                                                        \
    _x(yw_test_json_string_equals)                                             \
    _x(yw_test_json_string_to_c_str)                                           \
    _x(yw_test_json_expect_object)                                             \
    _x(yw_test_json_find_object_entry)                                         \
    _x(yw_test_json_expect_array)                                              \
    _x(yw_test_json_expect_string)                                             \
    _x(yw_test_json_expect_number)                                             \
    _x(yw_test_json_expect_boolean)                                            \
    _x(yw_test_json_expect_null)                                               \
    _x(yw_test_json_parse)
/* clang-format on */

typedef struct YW_TestingContext
{
    int failed_counter;
} YW_TestingContext;

#define YW_X(_name) void _name(YW_TestingContext *ctx);
YW_ENUMERATE_TESTS(YW_X)
#undef YW_X

#define YW_FAILED_TEST(_ctx, _msg)                  \
    do                                              \
    {                                               \
        YW_TestingContext *__ctx = (_ctx);          \
        printf("FAIL: %s(%s:%d): %s\n",             \
               __func__, __FILE__, __LINE__, _msg); \
        yw_failed_test_impl(__ctx);                 \
    } while (0)

#define YW_TEST_EXPECT(_type, _ctx, _got, _item_fmt, _expected)                        \
    do                                                                                 \
    {                                                                                  \
        YW_TestingContext *__ctx = (_ctx);                                             \
        _type __got = (_got);                                                          \
        _type __expected = (_expected);                                                \
        if (__got != __expected)                                                       \
        {                                                                              \
            printf("FAIL: %s(%s:%d): %s: expected " _item_fmt ", got " _item_fmt "\n", \
                   __func__, __FILE__, __LINE__, #_got, __expected, __got);            \
            yw_failed_test_impl(__ctx);                                                \
        }                                                                              \
    } while (0)

#define YW_TEST_EXPECT_STR(_ctx, _got, _expected)                                    \
    do                                                                               \
    {                                                                                \
        YW_TestingContext *__ctx = (_ctx);                                           \
        char const *__got = _got;                                                    \
        char const *__expected = _expected;                                          \
        if (__got == NULL && __expected == NULL)                                     \
        {                                                                            \
            break;                                                                   \
        }                                                                            \
        if (__got == NULL || (__expected != NULL && strcmp(__got, __expected) != 0)) \
        {                                                                            \
            printf("FAIL: %s(%s:%d): %s: expected [%s], got [%s]\n",                 \
                   __func__, __FILE__, __LINE__, #_got, __expected, __got);          \
            yw_failed_test_impl(__ctx);                                              \
        }                                                                            \
    } while (0)

#define YW_TEST_EXPECT_ARRAY(_item_type, _ctx, _got_array, _got_len, _item_fmt, ...)                       \
    do                                                                                                     \
    {                                                                                                      \
        YW_TestingContext *__ctx = (_ctx);                                                                 \
        _item_type __expected[] = {__VA_ARGS__};                                                           \
        _item_type *__got_array = (_got_array);                                                            \
        int __expected_len = YW_SIZEOF_ARRAY(__expected);                                                  \
        int __got_len = (_got_len);                                                                        \
        if (__expected_len != __got_len)                                                                   \
        {                                                                                                  \
            printf("FAIL: %s(%s:%d): %s: expected %d items, got %d\n",                                     \
                   __func__, __FILE__, __LINE__, #_got_array, __expected_len, __got_len);                  \
            yw_failed_test_impl(__ctx);                                                                    \
        }                                                                                                  \
        else                                                                                               \
        {                                                                                                  \
            for (int i = 0; i < __expected_len; i++)                                                       \
            {                                                                                              \
                if (__got_array[i] != __expected[i])                                                       \
                {                                                                                          \
                    printf("FAIL: %s(%s:%d): %s: expected " _item_fmt " at index %d, got " _item_fmt "\n", \
                           __func__, __FILE__, __LINE__, #_got_array, __expected[i], i, __got_array[i]);   \
                    yw_failed_test_impl(__ctx);                                                            \
                    break;                                                                                 \
                }                                                                                          \
            }                                                                                              \
        }                                                                                                  \
    } while (0)

void yw_failed_test_impl(YW_TestingContext *ctx);
void yw_run_all_tests(void);

#endif /* #ifndef YW_TESTS_H_ */
