/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_common.h"
#include "yw_tests.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*******************************************************************************
 * Memory utilities
 ******************************************************************************/

void yw_test_grow(YW_TestingContext *ctx)
{
    int len = 0;
    int cap = 0;
    char *buf = NULL;

    /* First we start with very normal tests **********************************/
    buf = YW_GROW(char, &cap, &len, buf);
    YW_TEST_EXPECT(cap, "%d", 2);
    YW_TEST_EXPECT(len, "%d", 1);

    buf = YW_GROW(char, &cap, &len, buf);
    YW_TEST_EXPECT(cap, "%d", 2);
    YW_TEST_EXPECT(len, "%d", 2);

    buf = YW_GROW(char, &cap, &len, buf);
    YW_TEST_EXPECT(cap, "%d", 6);
    YW_TEST_EXPECT(len, "%d", 3);

    free(buf);
}

void yw_test_shrink_to_fit(YW_TestingContext *ctx)
{
    int len = 0;
    int cap = 0;
    char *buf = (char *)malloc(1); /* Doesn't matter what size it is */
    if (buf == NULL)
    {
        printf("FATAL ERROR: %s: malloc() failed\n", __func__);
        abort();
    }

    /* First we start with very normal tests **********************************/
    len = 10;
    cap = 10;
    char *new_buf = YW_SHRINK_TO_FIT(char, &cap, len, buf);
    buf = new_buf;
    YW_TEST_EXPECT(cap, "%d", 10);
    YW_TEST_EXPECT(len, "%d", 10);

    len = 20;
    cap = 100;
    new_buf = YW_SHRINK_TO_FIT(char, &cap, len, buf);
    buf = new_buf;
    YW_TEST_EXPECT(cap, "%d", 20);
    YW_TEST_EXPECT(len, "%d", 20);

    free(buf);
}

void yw_test_list(YW_TestingContext *ctx)
{
    struct my_list
    {
        int *items;
        int len, cap;
    };
    struct my_list lst;

    YW_LIST_INIT(&lst);

    /* Add ********************************************************************/
    for (int i = 1; i <= 10; i++)
    {
        YW_LIST_PUSH(int, &lst, i);
    }
    YW_TEST_EXPECT(lst.items[0], "%d", 1);
    YW_TEST_EXPECT(lst.items[1], "%d", 2);
    YW_TEST_EXPECT(lst.items[2], "%d", 3);
    YW_TEST_EXPECT(lst.items[3], "%d", 4);
    YW_TEST_EXPECT(lst.items[4], "%d", 5);
    YW_TEST_EXPECT(lst.items[5], "%d", 6);
    YW_TEST_EXPECT(lst.items[6], "%d", 7);
    YW_TEST_EXPECT(lst.items[7], "%d", 8);
    YW_TEST_EXPECT(lst.items[8], "%d", 9);
    YW_TEST_EXPECT(lst.items[9], "%d", 10);

    /* Remove *****************************************************************/
    YW_LIST_REMOVE(int, &lst, 8); /* 1 2 3 4 5 6 7 8 [9] 10 */
    YW_LIST_REMOVE(int, &lst, 6); /* 1 2 3 4 5 6 [7] 8 10 */
    YW_LIST_REMOVE(int, &lst, 4); /* 1 2 3 4 [5] 6 8 10 */
    YW_LIST_REMOVE(int, &lst, 2); /* 1 2 [3] 4 6 8 10 */
    YW_LIST_REMOVE(int, &lst, 0); /* [1] 2 4 6 8 10 */

    YW_TEST_EXPECT(lst.items[0], "%d", 2);
    YW_TEST_EXPECT(lst.items[1], "%d", 4);
    YW_TEST_EXPECT(lst.items[2], "%d", 6);
    YW_TEST_EXPECT(lst.items[3], "%d", 8);
    YW_TEST_EXPECT(lst.items[4], "%d", 10);

    /* Insert *****************************************************************/
    YW_LIST_INSERT(int, &lst, 0, 20);  /* [20] 2 4 6 8 10 */
    YW_LIST_INSERT(int, &lst, 2, 40);  /* 20 2 [40] 4 6 8 10 */
    YW_LIST_INSERT(int, &lst, 4, 60);  /* 20 2 40 4 [60] 6 8 10 */
    YW_LIST_INSERT(int, &lst, 6, 80);  /* 20 2 40 4 60 6 [80] 8 10 */
    YW_LIST_INSERT(int, &lst, 9, 100); /* Insert at the end */

    YW_TEST_EXPECT(lst.items[0], "%d", 20);
    YW_TEST_EXPECT(lst.items[1], "%d", 2);
    YW_TEST_EXPECT(lst.items[2], "%d", 40);
    YW_TEST_EXPECT(lst.items[3], "%d", 4);
    YW_TEST_EXPECT(lst.items[4], "%d", 60);
    YW_TEST_EXPECT(lst.items[5], "%d", 6);
    YW_TEST_EXPECT(lst.items[6], "%d", 80);
    YW_TEST_EXPECT(lst.items[7], "%d", 8);
    YW_TEST_EXPECT(lst.items[8], "%d", 10);
    YW_TEST_EXPECT(lst.items[9], "%d", 100);

    YW_LIST_FREE(&lst);
}

/*******************************************************************************
 * UTF-8 character utility
 ******************************************************************************/

void yw_test_utf8_next_char(YW_TestingContext *ctx)
{
#define YW_RUN_TEST(_name, _input, ...)                                        \
    do                                                                         \
    {                                                                          \
        YW_Char32 expected[] = {__VA_ARGS__};                                  \
        int dest_len = YW_SIZEOF_ARRAY(expected);                              \
        char const *next_str = (_input);                                       \
        for (int i = 0; i < dest_len; i++)                                     \
        {                                                                      \
            YW_Char32 res = yw_utf8_next_char(&next_str);                      \
            if (res != expected[i])                                            \
            {                                                                  \
                printf(                                                        \
                    "FAIL: %s[%s]: expected U+%04X at index %d, got U+%04X\n", \
                    __func__, (_name), expected[i], i, res);                   \
                yw_failed_test(ctx);                                           \
                break;                                                         \
            }                                                                  \
        }                                                                      \
    } while (0)

    YW_RUN_TEST("Simple ASCII", "\x30\x31\x32\x33\x7e", '0', '1', '2', '3');
    YW_RUN_TEST("Two byte characters", "\xc2\xa0\xde\xb1", 0x00a0, 0x07b1);
    YW_RUN_TEST("Three byte characters", "\xe0\xa4\x80\xed\x9f\xbb\xef\xad\x8f",
                0x0900, 0xd7fb, 0xfb4f);
    YW_RUN_TEST("Four byte characters",
                "\xf0\x90\x91\x90\xf0\x9f\x83\xb5\xf4\x81\x8a\x8f", 0x10450,
                0x1f0f5, 0x10128f);

#undef YW_RUN_TEST
}

void yw_test_utf8_to_char32(YW_TestingContext *ctx)
{
    YW_Char32 *chars;
    int chars_len;
    yw_utf8_to_char32(&chars, &chars_len, "hello");
    YW_TEST_EXPECT_ARRAY(YW_Char32, chars, chars_len, "U+%04X", 'h', 'e', 'l',
                         'l', 'o');

    free(chars);
}

void yw_test_utf8_strlen(YW_TestingContext *ctx)
{
    char const *str = "This is so "
                      /* Kanji string "Kawaii" - U+53EF U+611B U+3044 */
                      "\xe5\x8f\xaf\xe6\x84\x9b\xe3\x81\x84";
    YW_TEST_EXPECT(yw_utf8_strlen(str), "%ld", 14L);
}

void yw_test_utf8_strchr(YW_TestingContext *ctx)
{
    char const *str = "This is so "
                      /* Kanji string "Kawaii" - U+53EF U+611B U+3044 */
                      "\xe5\x8f\xaf\xe6\x84\x9b\xe3\x81\x84";
    char const *ch0 = &str[11]; /* 1st kanji (U+53EF) */
    char const *ch1 = &str[14]; /* 2nd kanji (U+611B) */
    char const *ch2 = &str[17]; /* 3rd kanji (U+3044) */
    YW_TEST_EXPECT(yw_utf8_strchr(str, 0x53ef), "%s", ch0);
    YW_TEST_EXPECT(yw_utf8_strchr(str, 0x611b), "%s", ch1);
    YW_TEST_EXPECT(yw_utf8_strchr(str, 0x3044), "%s", ch2);
    YW_TEST_EXPECT(yw_utf8_strchr(str, 0x3045), "%s", (char *)NULL);
    YW_TEST_EXPECT(yw_utf8_strchr(str, '\0'), "%s", &str[strlen(str)]);
}

/*******************************************************************************
 * YW_TextReader
 ******************************************************************************/

void yw_test_peek_char(YW_TestingContext *ctx)
{
    struct YW_TextReader tr;
    YW_Char32 *chars;
    int chars_len;
    yw_utf8_to_char32(&chars, &chars_len, "hello");
    YW_TextReader_init(&tr, __func__, chars, chars_len);
    YW_TEST_EXPECT(yw_peek_char(&tr), "%c", 'h');
    YW_TEST_EXPECT(tr.cursor, "%d", 0);

    tr.cursor++;
    YW_TEST_EXPECT(yw_peek_char(&tr), "%c", 'e');
    YW_TEST_EXPECT(tr.cursor, "%d", 1);

    tr.cursor++;
    YW_TEST_EXPECT(yw_peek_char(&tr), "%c", 'l');
    YW_TEST_EXPECT(tr.cursor, "%d", 2);

    tr.cursor++;
    YW_TEST_EXPECT(yw_peek_char(&tr), "%c", 'l');
    YW_TEST_EXPECT(tr.cursor, "%d", 3);

    tr.cursor++;
    YW_TEST_EXPECT(yw_peek_char(&tr), "%c", 'o');
    YW_TEST_EXPECT(tr.cursor, "%d", 4);

    tr.cursor++;
    YW_TEST_EXPECT(yw_peek_char(&tr), "%d", -1);

    free(chars);
}
void yw_test_consume_any_char(YW_TestingContext *ctx)
{
    struct YW_TextReader tr;
    YW_Char32 *chars;
    int chars_len;
    yw_utf8_to_char32(&chars, &chars_len, "hello");
    YW_TextReader_init(&tr, __func__, chars, chars_len);

    YW_TEST_EXPECT(yw_consume_any_char(&tr), "%c", 'h');
    YW_TEST_EXPECT(tr.cursor, "%d", 1);

    YW_TEST_EXPECT(yw_consume_any_char(&tr), "%c", 'e');
    YW_TEST_EXPECT(tr.cursor, "%d", 2);

    YW_TEST_EXPECT(yw_consume_any_char(&tr), "%c", 'l');
    YW_TEST_EXPECT(tr.cursor, "%d", 3);

    YW_TEST_EXPECT(yw_consume_any_char(&tr), "%c", 'l');
    YW_TEST_EXPECT(tr.cursor, "%d", 4);

    YW_TEST_EXPECT(yw_consume_any_char(&tr), "%c", 'o');
    YW_TEST_EXPECT(tr.cursor, "%d", 5);

    YW_TEST_EXPECT(yw_consume_any_char(&tr), "%d", -1);

    free(chars);
}
void yw_test_consume_one_of_chars(YW_TestingContext *ctx)
{
    struct YW_TextReader tr;
    YW_Char32 *chars;
    int chars_len;
    yw_utf8_to_char32(&chars, &chars_len, "hello");
    YW_TextReader_init(&tr, __func__, chars, chars_len);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "abcdefgh"), "%c", 'h');
    YW_TEST_EXPECT(tr.cursor, "%d", 1);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "abcdefgh"), "%c", 'e');
    YW_TEST_EXPECT(tr.cursor, "%d", 2);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "abcdefgh"), "%c", -1);
    YW_TEST_EXPECT(tr.cursor, "%d", 2);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "hijklmn"), "%c", 'l');
    YW_TEST_EXPECT(tr.cursor, "%d", 3);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "hijklmn"), "%c", 'l');
    YW_TEST_EXPECT(tr.cursor, "%d", 4);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "opqrstu"), "%c", 'o');
    YW_TEST_EXPECT(tr.cursor, "%d", 5);

    YW_TEST_EXPECT(yw_consume_one_of_chars(&tr, "opqrstu"), "%c", -1);

    free(chars);
}
void yw_test_consume_one_of_strs(YW_TestingContext *ctx)
{
    struct YW_TextReader tr;
    YW_Char32 *chars;
    int chars_len;
    yw_utf8_to_char32(&chars, &chars_len,
                      "a quick fox jumps OvEr THE LAZY dog");
    YW_TextReader_init(&tr, __func__, chars, chars_len);

    char const *strs1[] = {"a ", "quick ", "fox jumps ", NULL};
    char const *strs2[] = {"oVeR ", "the lazy ", "DOG", NULL};

    YW_TEST_EXPECT(yw_consume_one_of_strs(&tr, strs1, YW_NO_MATCH_FLAGS), "%d",
                   0);
    YW_TEST_EXPECT(tr.cursor, "%d", 2);
    YW_TEST_EXPECT(yw_consume_one_of_strs(&tr, strs1, YW_NO_MATCH_FLAGS), "%d",
                   1);
    YW_TEST_EXPECT(tr.cursor, "%d", 8);
    YW_TEST_EXPECT(yw_consume_one_of_strs(&tr, strs1, YW_NO_MATCH_FLAGS), "%d",
                   2);
    YW_TEST_EXPECT(tr.cursor, "%d", 18);
    YW_TEST_EXPECT(yw_consume_one_of_strs(&tr, strs2, YW_NO_MATCH_FLAGS), "%d",
                   -1);
    YW_TEST_EXPECT(tr.cursor, "%d", 18);
    YW_TEST_EXPECT(
        yw_consume_one_of_strs(&tr, strs2, YW_ASCII_CASE_INSENSITIVE), "%d", 0);
    YW_TEST_EXPECT(tr.cursor, "%d", 23);
    YW_TEST_EXPECT(
        yw_consume_one_of_strs(&tr, strs2, YW_ASCII_CASE_INSENSITIVE), "%d", 1);
    YW_TEST_EXPECT(tr.cursor, "%d", 32);
    YW_TEST_EXPECT(
        yw_consume_one_of_strs(&tr, strs2, YW_ASCII_CASE_INSENSITIVE), "%d", 2);
    YW_TEST_EXPECT(tr.cursor, "%d", 35);
    YW_TEST_EXPECT(
        yw_consume_one_of_strs(&tr, strs2, YW_ASCII_CASE_INSENSITIVE), "%d",
        -1);

    free(chars);
}
