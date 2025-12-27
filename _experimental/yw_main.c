/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_common.h"
#include "yw_tests.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

static void gc_test(void);

int main(void)
{
    printf("hello, world!\n");
    yw_run_all_tests();
    // gc_test();
    return 0;
}

YW_GC_TYPE(struct my_object)
{
    struct YW_GcObjectHeader gc_header;
    struct my_object *another;
    int counter;
};

static void my_object_visit(void *self_v)
{
    printf("%p visit\n", self_v);
    yw_gc_visit(((YW_GC_PTR(struct my_object))self_v)->another);
}
static void my_object_destroy(void *self_v)
{
    printf("%p destroy\n", self_v);
}
static struct YW_GcCallbacks my_object_callbacks = {
    .visit = my_object_visit,
    .destroy = my_object_destroy,
};

static YW_GC_PTR(struct my_object) my_object_alloc(struct YW_GcHeap *heap, YW_GcAllocFlags alloc_flags)
{
    YW_GC_PTR(struct my_object)
    obj = YW_GC_ALLOC(struct my_object, heap, &my_object_callbacks, alloc_flags);
    obj->counter = 1;
    return obj;
}

static void gc_test(void)
{
    struct YW_GcHeap heap;
    yw_gc_heap_init(&heap);
    printf("heap init\n");

    YW_GC_PTR(struct my_object) objs[10];

    for (int i = 0; i < 10; i++)
    {
        YW_GC_PTR(struct my_object)
        obj = my_object_alloc(&heap, YW_GC_ROOT_OBJECT);
        printf("my object[%d] allocated @ %p\n", i, (void *)obj);
        obj->counter = i + 1;
        // obj->another = my_object_alloc(&heap, false);
        objs[i] = obj;
    }

    int sum = 0;
    for (int i = 0; i < 10; i++)
    {
        sum += objs[i]->counter;
    }
    printf("sum = %d\n", sum);
    printf("run GC\n");
    yw_gc(&heap);

    printf("run GC again\n");
    yw_gc(&heap);

    sum = 0;
    for (int i = 0; i < 10; i++)
    {
        sum += objs[i]->counter;
    }
    printf("sum = %d\n", sum);
}
