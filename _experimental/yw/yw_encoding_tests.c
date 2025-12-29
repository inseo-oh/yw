/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_common.h"
#include "yw_encoding.h"
#include "yw_tests.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void yw_test_encoding_from_label(YW_TestingContext *ctx)
{
    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_encoding_from_label("utf8"), "%d", YW_UTF8);
    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_encoding_from_label("shift-jis"), "%d", YW_SHIFT_JIS);
    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_encoding_from_label("ksc5601"), "%d", YW_EUC_KR);

    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_encoding_from_label("fox"), "%d", YW_INVALID_ENCODING);
}

void yw_test_bom_sniff(YW_TestingContext *ctx)
{
    YW_IOQueue queue;
    int utf8_bom[] = {0xef, 0xbb, 0xbf, 'A'};
    int utf16_be_bom[] = {0xfe, 0xff, 0x00, 0x00};
    int utf16_le_bom[] = {0xff, 0xfe, 0x00, 0x00};

    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, utf8_bom);
    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_bom_sniff(&queue), "%d", YW_UTF8);
    yw_io_queue_deinit(&queue);

    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, utf16_be_bom);
    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_bom_sniff(&queue), "%d", YW_UTF16_BE);
    yw_io_queue_deinit(&queue);

    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, utf16_le_bom);
    YW_TEST_EXPECT(YW_EncodingType, ctx, yw_bom_sniff(&queue), "%d", YW_UTF16_LE);
    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_item_list_to_items(YW_TestingContext *ctx)
{
    YW_IOQueueItemList list;
    YW_LIST_INIT(&list);
    YW_LIST_PUSH(YW_IOQueueItem, &list, 123);
    YW_LIST_PUSH(YW_IOQueueItem, &list, 456);
    YW_LIST_PUSH(YW_IOQueueItem, &list, 789);
    YW_LIST_PUSH(YW_IOQueueItem, &list, YW_END_OF_IO_QUEUE);
    YW_LIST_PUSH(YW_IOQueueItem, &list, 147);
    YW_LIST_PUSH(YW_IOQueueItem, &list, 258);
    YW_LIST_PUSH(YW_IOQueueItem, &list, 369);

    int *items, len;
    yw_io_queue_item_list_to_items(&items, &len, &list);

    YW_TEST_EXPECT_ARRAY(int, ctx, items, len, "%d", 123, 456, 789);

    YW_LIST_FREE(&list);
    free(items);
}

void yw_test_io_queue_from_items(YW_TestingContext *ctx)
{
    YW_IOQueue queue;
    int items[] = {123, 456, 789};

    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);
    YW_TEST_EXPECT_ARRAY(YW_IOQueueItem, ctx, queue.item_list.items, queue.item_list.len, "%d", 123, 456, 789, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_read_one(YW_TestingContext *ctx)
{
    YW_IOQueue queue;
    int items[] = {123, 456, 789};
    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);

    YW_TEST_EXPECT(YW_IOQueueItem, ctx, yw_io_queue_read_one(&queue), "%d", 123);
    YW_TEST_EXPECT(YW_IOQueueItem, ctx, yw_io_queue_read_one(&queue), "%d", 456);
    YW_TEST_EXPECT(YW_IOQueueItem, ctx, yw_io_queue_read_one(&queue), "%d", 789);
    YW_TEST_EXPECT(YW_IOQueueItem, ctx, yw_io_queue_read_one(&queue), "%d", YW_END_OF_IO_QUEUE);
    YW_TEST_EXPECT(YW_IOQueueItem, ctx, yw_io_queue_read_one(&queue), "%d", YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_read(YW_TestingContext *ctx)
{
    YW_IOQueue queue;
    int items[] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);

    int got_items[16];
    int len;

    len = yw_io_queue_read(&queue, got_items, 0);
    YW_TEST_EXPECT(int, ctx, len, "%d", 0);
    len = yw_io_queue_read(&queue, got_items, 1);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 1);
    len = yw_io_queue_read(&queue, got_items, 2);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 2, 3);
    len = yw_io_queue_read(&queue, got_items, 3);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 4, 5, 6);
    len = yw_io_queue_read(&queue, got_items, 5);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 7, 8, 9, 10);
    len = yw_io_queue_read(&queue, got_items, 5);
    YW_TEST_EXPECT(int, ctx, len, "%d", 0);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_peek(YW_TestingContext *ctx)
{
    YW_IOQueue queue;
    int items[] = {123, 456, 789};
    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);

    int got_items[16];
    int len;

    len = yw_io_queue_peek(&queue, got_items, 0);
    YW_TEST_EXPECT(int, ctx, len, "%d", 0);
    len = yw_io_queue_peek(&queue, got_items, 1);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 123);
    len = yw_io_queue_peek(&queue, got_items, 2);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 123, 456);
    len = yw_io_queue_peek(&queue, got_items, 10);
    YW_TEST_EXPECT_ARRAY(int, ctx, got_items, len, "%d", 123, 456, 789);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_push_one(YW_TestingContext *ctx)
{
    YW_IOQueue queue;

    yw_io_queue_init(&queue);
    yw_io_queue_push_one(&queue, 123);
    yw_io_queue_push_one(&queue, 456);
    yw_io_queue_push_one(&queue, 789);
    yw_io_queue_push_one(&queue, YW_END_OF_IO_QUEUE);
    yw_io_queue_push_one(&queue, 147);
    yw_io_queue_push_one(&queue, 258);
    yw_io_queue_push_one(&queue, 369);
    YW_TEST_EXPECT_ARRAY(YW_IOQueueItem, ctx, queue.item_list.items, queue.item_list.len, "%d", 123, 456, 789, 147, 258, 369, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_push(YW_TestingContext *ctx)
{
    YW_IOQueueItem items[] = {123, 456, 789, YW_END_OF_IO_QUEUE, 147, 258, 369};

    YW_IOQueue queue;

    yw_io_queue_init(&queue);
    YW_IO_QUEUE_PUSH_FROM_ARRAY(&queue, items);
    YW_TEST_EXPECT_ARRAY(YW_IOQueueItem, ctx, queue.item_list.items, queue.item_list.len, "%d", 123, 456, 789, 147, 258, 369, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_restore_one(YW_TestingContext *ctx)
{
    YW_IOQueue queue;

    yw_io_queue_init(&queue);
    yw_io_queue_push_one(&queue, 1000);
    yw_io_queue_restore_one(&queue, 123);
    yw_io_queue_restore_one(&queue, 456);
    yw_io_queue_restore_one(&queue, 789);
    yw_io_queue_restore_one(&queue, 147);
    yw_io_queue_restore_one(&queue, 258);
    yw_io_queue_restore_one(&queue, 369);
    yw_io_queue_push_one(&queue, 2000);
    YW_TEST_EXPECT_ARRAY(YW_IOQueueItem, ctx, queue.item_list.items, queue.item_list.len, "%d", 369, 258, 147, 789, 456, 123, 1000, 2000, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_restore(YW_TestingContext *ctx)
{
    YW_IOQueueItem items[] = {
        123,
        456,
        789,
        147,
        258,
        369,
    };

    YW_IOQueue queue;

    yw_io_queue_init(&queue);
    yw_io_queue_push_one(&queue, 1000);
    YW_IO_QUEUE_RESTORE_FROM_ARRAY(&queue, items);
    yw_io_queue_push_one(&queue, 2000);
    YW_TEST_EXPECT_ARRAY(YW_IOQueueItem, ctx, queue.item_list.items, queue.item_list.len, "%d", 369, 258, 147, 789, 456, 123, 1000, 2000, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

/*******************************************************************************
 * Tests for encoding implementations
 ******************************************************************************/

#define YW_TEST_DECODER(_encoding, _name, _input, ...)                       \
    do                                                                       \
    {                                                                        \
        YW_IOQueue input_queue;                                              \
        YW_IOQueue output_queue;                                             \
        yw_io_queue_init(&input_queue);                                      \
        yw_io_queue_init(&output_queue);                                     \
        for (int i = 0; i < (int)strlen((_input)); i++)                      \
        {                                                                    \
            yw_io_queue_push_one(&input_queue, (uint8_t)_input[i]);          \
        }                                                                    \
        yw_encoding_decode(&input_queue, (_encoding), &output_queue);        \
        int expected[] = {__VA_ARGS__};                                      \
        int expected_len = YW_SIZEOF_ARRAY(expected);                        \
        for (int i = 0; i < expected_len; i++)                               \
        {                                                                    \
            if (output_queue.item_list.items[i] == YW_END_OF_IO_QUEUE)       \
            {                                                                \
                printf("FAIL: %s[%s]: expected U+%04X at index %d, reached " \
                       "end of queue\n",                                     \
                       __func__, (_name), expected[i], i);                   \
                yw_failed_test_impl(ctx);                                    \
                break;                                                       \
            }                                                                \
            else                                                             \
            {                                                                \
                int res = (int)output_queue.item_list.items[i];              \
                if (res != expected[i])                                      \
                {                                                            \
                    printf("FAIL: %s[%s]: expected U+%04X at index %d, got " \
                           "U+%04X\n",                                       \
                           __func__, (_name), expected[i], i, res);          \
                    yw_failed_test_impl(ctx);                                \
                    break;                                                   \
                }                                                            \
            }                                                                \
        }                                                                    \
        yw_io_queue_deinit(&input_queue);                                    \
        yw_io_queue_deinit(&output_queue);                                   \
    } while (0)

void yw_test_utf8_decoder(YW_TestingContext *ctx)
{
    YW_TEST_DECODER(YW_UTF8, "Simple ASCII", "\x30\x31\x32\x33\x7e", '0', '1', '2', '3');
    YW_TEST_DECODER(YW_UTF8, "Two byte characters", "\xc2\xa0\xde\xb1", 0x00a0, 0x07b1);
    YW_TEST_DECODER(YW_UTF8, "Three byte characters", "\xe0\xa4\x80\xed\x9f\xbb\xef\xad\x8f", 0x0900, 0xd7fb, 0xfb4f);
    YW_TEST_DECODER(YW_UTF8, "Four byte characters", "\xf0\x90\x91\x90\xf0\x9f\x83\xb5\xf4\x81\x8a\x8f", 0x10450, 0x1f0f5, 0x10128f);
}
