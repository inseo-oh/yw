/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_COMMON_H_
#define YW_COMMON_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef int32_t YW_CHAR32;

#define YW_TODO()                                                              \
    do                                                                         \
    {                                                                          \
        fprintf(stderr, "[%s:%d] %s: TODO", __FILE__, __LINE__, __func__);     \
        abort();                                                               \
    } while (0)

/*******************************************************************************
 * Testing support
 ******************************************************************************/

typedef struct yw_testing_context yw_testing_context;
struct yw_testing_context
{
    int failed_counter;
};

/* clang-format off */
#define YW_ENUMERATE_TESTS(_x)                                                 \
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
    _x(yw_test_consume_one_of_strs)
/* clang-format on */

#define YW_X(_name) void _name(yw_testing_context *ctx);
YW_ENUMERATE_TESTS(YW_X)
#undef YW_X

/* NOTE: _got will be re-evaulated multiple times! */
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

/* NOTE: _got_array will be re-evaulated multiple times! */
#define YW_TEST_EXPECT_ARRAY(_got_array, _got_len, _item_fmt, ...)             \
    do                                                                         \
    {                                                                          \
        YW_CHAR32 expected[] = {__VA_ARGS__};                                  \
        int expected_len = sizeof(expected) / sizeof(*expected);               \
        if (expected_len != (_got_len))                                        \
        {                                                                      \
            printf("FAIL: %s: %s: expected %d items, got %d\n", __func__,      \
                   #_got_array, expected_len, (_got_len));                     \
            yw_failed_test(ctx);                                               \
        }                                                                      \
        for (int i = 0; i < expected_len; i++)                                 \
        {                                                                      \
            if ((_got_array)[i] != expected[i])                                \
            {                                                                  \
                printf("FAIL: %s: %s: expected " _item_fmt                     \
                       " at index %d, got " _item_fmt "\n",                    \
                       __func__, #_got_array, expected[i], i,                  \
                       (_got_array)[i]);                                       \
                yw_failed_test(ctx);                                           \
                break;                                                         \
            }                                                                  \
        }                                                                      \
    } while (0)

void yw_failed_test(yw_testing_context *ctx);
void yw_run_all_tests();

/*******************************************************************************
 * ASCII character conversion & testing
 ******************************************************************************/

bool yw_is_leading_surrogate_char(YW_CHAR32 c);
bool yw_is_trailing_surrogate_char(YW_CHAR32 c);
bool yw_is_surrogate_char(YW_CHAR32 c);
bool yw_is_c0_control_char(YW_CHAR32 c);
bool yw_is_control_char(YW_CHAR32 c);
bool yw_is_ascii_digit(YW_CHAR32 c);
bool yw_is_ascii_uppercase(YW_CHAR32 c);
bool yw_is_ascii_lowercase(YW_CHAR32 c);
bool yw_is_ascii_alpha(YW_CHAR32 c);
bool yw_is_ascii_alphanumeric(YW_CHAR32 c);
bool yw_is_ascii_uppercase_hex_digit(YW_CHAR32 c);
bool yw_is_ascii_lowercase_hex_digit(YW_CHAR32 c);
bool yw_is_ascii_hex_digit(YW_CHAR32 c);
bool yw_is_ascii_whitespace(YW_CHAR32 c);
bool yw_is_noncharacter(YW_CHAR32 c);
YW_CHAR32 yw_to_ascii_lowercase(YW_CHAR32 c);
YW_CHAR32 yw_to_ascii_uppercase(YW_CHAR32 c);

/*******************************************************************************
 * Memory utilities
 ******************************************************************************/

void *yw_grow_impl(int *cap_inout, int *len_inout, void *old_buf,
                   size_t item_size);
void *yw_shrink_to_fit_impl(int *cap_inout, int len, void *old_buf,
                            size_t item_size);
#define YW_GROW(_type, _cap_inout, _len_inout, _old_buf)                       \
    (_type *)yw_grow_impl((_cap_inout), (_len_inout), (_old_buf), sizeof(_type))
#define YW_SHRINK_TO_FIT(_type, _cap_inout, _len, _old_buf)                    \
    (_type *)yw_shrink_to_fit_impl((_cap_inout), (_len), (_old_buf),           \
                                   sizeof(_type))

#define YW_LIST_INIT(_list)                                                    \
    do                                                                         \
    {                                                                          \
        memset((_list), 0, sizeof(*(_list)));                                  \
    } while (0)
#define YW_LIST_FREE(_list)                                                    \
    do                                                                         \
    {                                                                          \
        free((_list)->items);                                                  \
    } while (0)
#define YW_LIST_PUSH(_type, _list, _item)                                      \
    do                                                                         \
    {                                                                          \
        (_list)->items =                                                       \
            YW_GROW(_type, &(_list)->cap, &(_list)->len, (_list)->items);      \
        (_list)->items[(_list)->len - 1] = (_item);                            \
    } while (0)
#define YW_LIST_REMOVE(_type, _list, _index)                                   \
    do                                                                         \
    {                                                                          \
        if ((_list)->len <= _index)                                            \
        {                                                                      \
            fprintf(stderr, "illegal list item index %d\n", _index);           \
            abort();                                                           \
        }                                                                      \
        int copy_count = (_list)->len - 1 - (_index);                          \
        for (int i = 0; i < copy_count; i++)                                   \
        {                                                                      \
            (_list)->items[(_index) + i] = (_list)->items[(_index) + i + 1];   \
        }                                                                      \
        (_list)->len--;                                                        \
        (_list)->items = YW_SHRINK_TO_FIT(_type, &(_list)->cap, (_list)->len,  \
                                          (_list)->items);                     \
    } while (0)
#define YW_LIST_INSERT(_type, _list, _index, _item)                            \
    do                                                                         \
    {                                                                          \
        if ((_list)->len < _index)                                             \
        {                                                                      \
            fprintf(stderr, "illegal list item index %d\n", _index);           \
            abort();                                                           \
        }                                                                      \
        (_list)->items =                                                       \
            YW_GROW(_type, &(_list)->cap, &(_list)->len, (_list)->items);      \
        int copy_count = (_list)->len - (_index);                              \
        for (int i = copy_count - 1; 0 <= i; i--)                              \
        {                                                                      \
            (_list)->items[(_index) + i + 1] = (_list)->items[(_index) + i];   \
        }                                                                      \
        (_list)->items[(_index)] = (_item);                                    \
    } while (0)

/*******************************************************************************
 * UTF-8 character utility
 ******************************************************************************/

/*
 * Returns resulting codepoint, or 0 if end was reached.
 * If an error was encountered, U+FFFD is returned instead.
 */
YW_CHAR32 yw_utf8_next_char(char const **str);

/* Caller owns the returned array. */
void yw_utf8_to_char32(YW_CHAR32 **chars_out, int *chars_len_out,
                       char const *str);

/* UTF-8-aware version of strlen(). */
size_t yw_utf8_strlen(char const *s);

/*
 * UTF-8-aware version of strchr().
 *
 * NOTE: For searching non-unicode character, using normal strchr() is OK.
 *       (And may even be faster)
 */
char const *yw_utf8_strchr(char const *s, YW_CHAR32 c);

/*******************************************************************************
 * yw_text_reader
 ******************************************************************************/

typedef struct yw_text_reader yw_text_reader;
struct yw_text_reader
{
    char const *source_name;
    YW_CHAR32 const *chars;
    int chars_len;
    int cursor;
};

void yw_text_reader_init(yw_text_reader *out, char const *source_name,
                         YW_CHAR32 const *chars, int chars_len);
bool yw_text_reader_is_eof(yw_text_reader const *tr);

/* Returns -1 on EOF. */
YW_CHAR32 yw_peek_char(yw_text_reader const *tr);

/* Returns -1 on EOF. */
YW_CHAR32 yw_consume_any_char(yw_text_reader *tr);

/*
 * Returns -1 on EOF or when no match was found.
 * Also note that this function can only match ASCII characters.
 */
int yw_consume_one_of_chars(yw_text_reader *tr, char const *chars);
bool yw_consume_char(yw_text_reader *tr, YW_CHAR32 chr);

typedef enum
{
    YW_NO_MATCH_FLAGS = 0,
    YW_ASCII_CASE_INSENSITIVE = 1 << 0
} yw_match_flags;

/*
 * Returns index of matched string, or -1 if not found.
 *
 * strs must be NULL-terminated list!
 */
int yw_consume_one_of_strs(yw_text_reader *tr, char const **strs,
                           yw_match_flags flags);
bool yw_consume_str(yw_text_reader *tr, char const *str, yw_match_flags flags);

/*******************************************************************************
 * Garbage collector
 ******************************************************************************/

typedef void *YW_PTR_SLOT;

/*
 * Each slot may store either a pointer or NULL.
 * NULL means free "slot", and new pointers can be stored there.
 */
typedef struct yw_ptr_collection yw_ptr_collection;
struct yw_ptr_collection
{
    YW_PTR_SLOT *slots;
    int slots_len;
    int slots_cap;
};

/* Returns pointer to the slot, or NULL if there's not enough memory. */
YW_PTR_SLOT *yw_add_ptr_to_collection(yw_ptr_collection *coll, void *obj);

typedef struct yw_gc_callbacks yw_gc_callbacks;
struct yw_gc_callbacks
{
    void (*visit)(void *self);
    void (*destroy)(void *self);
};

typedef struct yw_gc_object_header yw_gc_object_header;
struct yw_gc_object_header
{
    uint64_t magic_and_marked_flag; /* LSB is used as marked flag. */
    yw_gc_callbacks const *callbacks;
};

typedef struct yw_gc_heap yw_gc_heap;
struct yw_gc_heap
{
    yw_ptr_collection all_objs;
    yw_ptr_collection root_objs;
};

typedef enum
{
    YW_NO_GC_ALLOC_FLAGS = 0,
    YW_ADD_TO_GC_ROOT = 1 << 0
} yw_gc_alloc_flags;

/* NOTE: It is safe to pass NULL pointer. */
void yw_gc_visit(void *obj_v);

#define YW_GC_TYPE(_x) _x##_GC
#define YW_GC_PTR(_x) YW_GC_TYPE(_x) *

void yw_gc_init_heap(yw_gc_heap *out);
void *yw_gc_alloc_impl(yw_gc_heap *heap, int size,
                       yw_gc_callbacks const *callbacks,
                       yw_gc_alloc_flags alloc_flags);
#define YW_GC_ALLOC(_type, _heap, _callbacks, _alloc_flags)                    \
    (YW_GC_PTR(_type))yw_gc_alloc_impl((_heap), sizeof(YW_GC_TYPE(_type)),     \
                                       (_callbacks), (_alloc_flags))

void yw_gc(yw_gc_heap *heap);

#endif /* #ifndef YW_COMMON_H_ */