/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_common.h"
#include "yw_json.h"
#include "yw_tests.h"
#include <stdlib.h>
#include <string.h>

void yw_test_json_string_equals(YW_TestingContext *ctx)
{
    YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(NULL, "hello"), "%d", false);
    YW_JSONString test_str;

    yw_json_string_init(&test_str, "hello");
    YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(&test_str, "hello"), "%d", true);
    YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(&test_str, "hell"), "%d", false);
    yw_json_string_deinit(&test_str);

    yw_json_string_init(&test_str, "hel_lo");
    test_str.chars[3] = '\0';
    YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(&test_str, "hel\0lo"), "%d", false);
    yw_json_string_deinit(&test_str);
}

void yw_test_json_expect_object(YW_TestingContext *ctx)
{
    YW_JSONObjectEntry *ents = NULL;
    int ents_len = 0;
    int ents_cap = 0;
    YW_JSONValue *v1 = yw_json_new_number(34);
    YW_JSONValue *v2 = yw_json_new_number(35);
    yw_json_add_value_to_object_entry(&ents_cap, &ents_len, &ents, "x", &v1);
    yw_json_add_value_to_object_entry(&ents_cap, &ents_len, &ents, "y", &v2);
    YW_SHRINK_TO_FIT(YW_JSONObjectEntry, &ents_cap, ents_len, &ents);
    YW_JSONValue *v = yw_json_new_object(&ents, &ents_len);

    YW_JSONObjectValue const *ov = yw_json_expect_object(v);
    if (ov != NULL)
    {
        YW_TEST_EXPECT(int, ctx, ov->len, "%d", 2);

        if (ov->len == 2)
        {
            double v_num;

            /*
             * NOTE: We assume that entries have the same order as original object entry list.
             * So if we end up switching to hashmaps in the future, this test may fail.
             */

            YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(&ov->entries[0].name, "x"), "%d", true);
            YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, ov->entries[0].value), "%d", true);
            YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 34);

            YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(&ov->entries[1].name, "y"), "%d", true);
            YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, ov->entries[1].value), "%d", true);
            YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 35);
        }
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_json_expect_object returned NULL");
    }

    yw_json_value_free(v);
}

void yw_test_json_find_object_entry(YW_TestingContext *ctx)
{
    YW_JSONObjectEntry *ents = NULL;
    int ents_len = 0;
    int ents_cap = 0;
    YW_JSONValue *v1 = yw_json_new_number(34);
    YW_JSONValue *v2 = yw_json_new_number(35);
    yw_json_add_value_to_object_entry(&ents_cap, &ents_len, &ents, "x", &v1);
    yw_json_add_value_to_object_entry(&ents_cap, &ents_len, &ents, "y", &v2);
    YW_SHRINK_TO_FIT(YW_JSONObjectEntry, &ents_cap, ents_len, &ents);
    YW_JSONValue *v = yw_json_new_object(&ents, &ents_len);

    double v_num;

    YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, yw_json_find_object_entry(v, "x")), "%d", true);
    YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 34);

    YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, yw_json_find_object_entry(v, "y")), "%d", true);
    YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 35);

    yw_json_value_free(v);
}

void yw_test_json_expect_array(YW_TestingContext *ctx)
{
    YW_JSONValue **ents = NULL;
    int ents_len = 0;
    int ents_cap = 0;
    YW_JSONValue *v1 = yw_json_new_number(34);
    YW_JSONValue *v2 = yw_json_new_number(35);
    YW_PUSH(YW_JSONValue *, &ents_cap, &ents_len, &ents, v1);
    YW_PUSH(YW_JSONValue *, &ents_cap, &ents_len, &ents, v2);
    YW_SHRINK_TO_FIT(YW_JSONValue *, &ents_cap, ents_len, &ents);
    YW_JSONValue *v = yw_json_new_array(&ents, &ents_len);

    YW_JSONArrayValue const *av = yw_json_expect_array(v);
    if (av != NULL)
    {
        YW_TEST_EXPECT(int, ctx, av->len, "%d", 2);

        if (av->len == 2)
        {
            double v_num;

            YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, av->entries[0]), "%d", true);
            YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 34);

            YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, av->entries[1]), "%d", true);
            YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 35);
        }
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_json_expect_array returned NULL");
    }

    yw_json_value_free(v);
}

void yw_test_json_expect_string(YW_TestingContext *ctx)
{
    YW_JSONValue *v = yw_json_new_string("hello, world!");

    YW_JSONString const *sv = yw_json_expect_string(v);
    if (sv != NULL)
    {
        YW_TEST_EXPECT(bool, ctx, yw_json_string_equals(sv, "hello, world!"), "%d", true);
    }
    else
    {
        YW_FAILED_TEST(ctx, "yw_json_expect_string returned NULL");
    }

    yw_json_value_free(v);
}

void yw_test_json_expect_number(YW_TestingContext *ctx)
{
    YW_JSONValue *v = yw_json_new_number(34);

    double v_num;
    YW_TEST_EXPECT(bool, ctx, yw_json_expect_number(&v_num, v), "%d", true);
    YW_TEST_EXPECT(int, ctx, (int)v_num, "%d", 34);

    yw_json_value_free(v);
}

void yw_test_json_expect_boolean(YW_TestingContext *ctx)
{
    YW_JSONValue *v1 = yw_json_new_boolean(true);
    YW_JSONValue *v2 = yw_json_new_boolean(false);

    bool v_bol;
    YW_TEST_EXPECT(bool, ctx, yw_json_expect_boolean(&v_bol, v1), "%d", true);
    YW_TEST_EXPECT(bool, ctx, v_bol, "%d", true);
    YW_TEST_EXPECT(bool, ctx, yw_json_expect_boolean(&v_bol, v2), "%d", true);
    YW_TEST_EXPECT(bool, ctx, v_bol, "%d", false);

    yw_json_value_free(v1);
    yw_json_value_free(v2);
}

void yw_test_json_expect_null(YW_TestingContext *ctx)
{
    YW_JSONValue *v = yw_json_new_null();

    YW_TEST_EXPECT(bool, ctx, yw_json_expect_null(v), "%d", true);

    yw_json_value_free(v);
}
