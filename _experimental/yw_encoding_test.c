/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_common.h"
#include "yw_encoding.h"
#include "yw_tests.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void yw_test_encoding_from_label(YW_TestingContext *ctx)
{
    YW_TEST_EXPECT(yw_encoding_from_label("utf8"), "%d", YW_UTF8);
    YW_TEST_EXPECT(yw_encoding_from_label("shift-jis"), "%d", YW_SHIFT_JIS);
    YW_TEST_EXPECT(yw_encoding_from_label("ksc5601"), "%d", YW_EUC_KR);

    YW_TEST_EXPECT(yw_encoding_from_label("fox"), "%d", YW_INVALID_ENCODING);
}

void yw_test_io_queue_item_list_to_items(YW_TestingContext *ctx)
{
    YW_IoQueueItemList list;
    YW_LIST_INIT(&list);
    YW_LIST_PUSH(YW_IoQueueItem, &list, (YW_IoQueueItem)123);
    YW_LIST_PUSH(YW_IoQueueItem, &list, (YW_IoQueueItem)456);
    YW_LIST_PUSH(YW_IoQueueItem, &list, (YW_IoQueueItem)789);
    YW_LIST_PUSH(YW_IoQueueItem, &list, YW_END_OF_IO_QUEUE);
    YW_LIST_PUSH(YW_IoQueueItem, &list, (YW_IoQueueItem)147);
    YW_LIST_PUSH(YW_IoQueueItem, &list, (YW_IoQueueItem)258);
    YW_LIST_PUSH(YW_IoQueueItem, &list, (YW_IoQueueItem)369);

    int *items, len;
    yw_io_queue_item_list_to_items(&items, &len, &list);

    YW_TEST_EXPECT_ARRAY(int, items, len, "%d", 123, 456, 789,
                         YW_END_OF_IO_QUEUE);

    free(items);
}

void yw_test_io_queue_from_items(YW_TestingContext *ctx)
{
    YW_IoQueue queue;
    int items[] = {123, 456, 789};

    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);
    YW_TEST_EXPECT_ARRAY(YW_IoQueueItem, queue.item_list.items,
                         queue.item_list.len, "%d", (YW_IoQueueItem)123,
                         (YW_IoQueueItem)456, (YW_IoQueueItem)789,
                         YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_read_one(YW_TestingContext *ctx)
{
    YW_IoQueue queue;
    int items[] = {123, 456, 789};
    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);

    YW_TEST_EXPECT(yw_io_queue_read_one(&queue), "%d", (YW_IoQueueItem)123);
    YW_TEST_EXPECT(yw_io_queue_read_one(&queue), "%d", (YW_IoQueueItem)456);
    YW_TEST_EXPECT(yw_io_queue_read_one(&queue), "%d", (YW_IoQueueItem)789);
    YW_TEST_EXPECT(yw_io_queue_read_one(&queue), "%d", YW_END_OF_IO_QUEUE);
    YW_TEST_EXPECT(yw_io_queue_read_one(&queue), "%d", YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_read(YW_TestingContext *ctx)
{
    YW_IoQueue queue;
    int items[] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);

    int got_items[16];
    int len;

    len = yw_io_queue_read(&queue, got_items, 0);
    YW_TEST_EXPECT(len, "%d", 0);
    len = yw_io_queue_read(&queue, got_items, 1);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 1);
    len = yw_io_queue_read(&queue, got_items, 2);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 2, 3);
    len = yw_io_queue_read(&queue, got_items, 3);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 4, 5, 6);
    len = yw_io_queue_read(&queue, got_items, 5);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 7, 8, 9, 10);
    len = yw_io_queue_read(&queue, got_items, 5);
    YW_TEST_EXPECT(len, "%d", 0);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_peek(YW_TestingContext *ctx)
{
    YW_IoQueue queue;
    int items[] = {123, 456, 789};
    YW_IO_QUEUE_INIT_FROM_ARRAY(&queue, items);

    int got_items[16];
    int len;

    len = yw_io_queue_peek(&queue, got_items, 0);
    YW_TEST_EXPECT(len, "%d", 0);
    len = yw_io_queue_peek(&queue, got_items, 1);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 123);
    len = yw_io_queue_peek(&queue, got_items, 2);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 123, 456);
    len = yw_io_queue_peek(&queue, got_items, 10);
    YW_TEST_EXPECT_ARRAY(int, got_items, len, "%d", 123, 456, 789);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_push_one(YW_TestingContext *ctx)
{
    YW_IoQueue queue;

    yw_io_queue_init(&queue);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)123);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)456);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)789);
    yw_io_queue_push_one(&queue, YW_END_OF_IO_QUEUE);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)147);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)258);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)369);
    YW_TEST_EXPECT_ARRAY(YW_IoQueueItem, queue.item_list.items,
                         queue.item_list.len, "%d", (YW_IoQueueItem)123,
                         (YW_IoQueueItem)456, (YW_IoQueueItem)789,
                         (YW_IoQueueItem)147, (YW_IoQueueItem)258,
                         (YW_IoQueueItem)369, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_push(YW_TestingContext *ctx)
{
    YW_IoQueueItem items[] = {(YW_IoQueueItem)123, (YW_IoQueueItem)456,
                              (YW_IoQueueItem)789, YW_END_OF_IO_QUEUE,
                              (YW_IoQueueItem)147, (YW_IoQueueItem)258,
                              (YW_IoQueueItem)369};

    YW_IoQueue queue;

    yw_io_queue_init(&queue);
    YW_IO_QUEUE_PUSH_FROM_ARRAY(&queue, items);
    YW_TEST_EXPECT_ARRAY(YW_IoQueueItem, queue.item_list.items,
                         queue.item_list.len, "%d", (YW_IoQueueItem)123,
                         (YW_IoQueueItem)456, (YW_IoQueueItem)789,
                         (YW_IoQueueItem)147, (YW_IoQueueItem)258,
                         (YW_IoQueueItem)369, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_restore_one(YW_TestingContext *ctx)
{
    YW_IoQueue queue;

    yw_io_queue_init(&queue);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)1000);
    yw_io_queue_restore_one(&queue, (YW_IoQueueItem)123);
    yw_io_queue_restore_one(&queue, (YW_IoQueueItem)456);
    yw_io_queue_restore_one(&queue, (YW_IoQueueItem)789);
    yw_io_queue_restore_one(&queue, (YW_IoQueueItem)147);
    yw_io_queue_restore_one(&queue, (YW_IoQueueItem)258);
    yw_io_queue_restore_one(&queue, (YW_IoQueueItem)369);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)2000);
    YW_TEST_EXPECT_ARRAY(
        YW_IoQueueItem, queue.item_list.items, queue.item_list.len, "%d",
        (YW_IoQueueItem)369, (YW_IoQueueItem)258, (YW_IoQueueItem)147,
        (YW_IoQueueItem)789, (YW_IoQueueItem)456, (YW_IoQueueItem)123,
        (YW_IoQueueItem)1000, (YW_IoQueueItem)2000, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}

void yw_test_io_queue_restore(YW_TestingContext *ctx)
{
    YW_IoQueueItem items[] = {
        (YW_IoQueueItem)123, (YW_IoQueueItem)456, (YW_IoQueueItem)789,
        YW_END_OF_IO_QUEUE,  (YW_IoQueueItem)147, (YW_IoQueueItem)258,
        (YW_IoQueueItem)369,
    };

    YW_IoQueue queue;

    yw_io_queue_init(&queue);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)1000);
    YW_IO_QUEUE_RESTORE_FROM_ARRAY(&queue, items);
    yw_io_queue_push_one(&queue, (YW_IoQueueItem)2000);
    YW_TEST_EXPECT_ARRAY(
        YW_IoQueueItem, queue.item_list.items, queue.item_list.len, "%d",
        (YW_IoQueueItem)369, (YW_IoQueueItem)258, (YW_IoQueueItem)147,
        (YW_IoQueueItem)789, (YW_IoQueueItem)456, (YW_IoQueueItem)123,
        (YW_IoQueueItem)1000, (YW_IoQueueItem)2000, YW_END_OF_IO_QUEUE);

    yw_io_queue_deinit(&queue);
}
