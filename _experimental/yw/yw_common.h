/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#ifndef YW_COMMON_H_
#define YW_COMMON_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

typedef int YW_Char32;

#define YW_TODO()                              \
    do                                         \
    {                                          \
        fprintf(stderr, "[%s:%d] %s: TODO\n",  \
                __FILE__, __LINE__, __func__); \
        abort();                               \
    } while (0)

#define YW_UNREACHABLE()                             \
    do                                               \
    {                                                \
        fprintf(stderr, "[%s:%d] %s: unreachable\n", \
                __FILE__, __LINE__, __func__);       \
        abort();                                     \
    } while (0)

#define YW_ILLEGAL_VALUE(_v)                                        \
    do                                                              \
    {                                                               \
        fprintf(stderr, "[%s:%d] %s: %d is illegal value for %s\n", \
                __FILE__, __LINE__, __func__, _v, #_v);             \
        abort();                                                    \
    } while (0)

/*******************************************************************************
 * Namespaces
 ******************************************************************************/

#define YW_HTML_NAMESPACE "http://www.w3.org/1999/xhtml"
#define YW_MATHML_NAMESPACE "http://www.w3.org/1998/Math/MathML"
#define YW_SVG_NAMESPACE "http://www.w3.org/2000/svg"
#define YW_XLINK_NAMESPACE "http://www.w3.org/1999/xlink"
#define YW_XML_NAMESPACE "http://www.w3.org/XML/1998/namespace"
#define YW_XMLNS_NAMESPACE "http://www.w3.org/2000/xmlns/"

/*******************************************************************************
 * ASCII character conversion & testing
 ******************************************************************************/

bool yw_is_leading_surrogate_char(YW_Char32 c);
bool yw_is_trailing_surrogate_char(YW_Char32 c);
bool yw_is_surrogate_char(YW_Char32 c);
bool yw_is_c0_control_char(YW_Char32 c);
bool yw_is_control_char(YW_Char32 c);
bool yw_is_ascii_digit(YW_Char32 c);
bool yw_is_ascii_uppercase(YW_Char32 c);
bool yw_is_ascii_lowercase(YW_Char32 c);
bool yw_is_ascii_alpha(YW_Char32 c);
bool yw_is_ascii_alphanumeric(YW_Char32 c);
bool yw_is_ascii_uppercase_hex_digit(YW_Char32 c);
bool yw_is_ascii_lowercase_hex_digit(YW_Char32 c);
bool yw_is_ascii_hex_digit(YW_Char32 c);
bool yw_is_ascii_whitespace(YW_Char32 c);
bool yw_is_noncharacter(YW_Char32 c);
YW_Char32 yw_to_ascii_lowercase(YW_Char32 c);
YW_Char32 yw_to_ascii_uppercase(YW_Char32 c);
int yw_strcmp_ascii_case_insensitive(char const *a, char const *b);

/*******************************************************************************
 * Memory utilities
 ******************************************************************************/

#define YW_SIZEOF_ARRAY(_x) (sizeof((_x)) / sizeof(*(_x)))

void *yw_grow_impl(int *cap_inout, int *len_inout, void *old_buf, size_t item_size, int add_len);
void *yw_alloc_impl(size_t item_size, int count);
void *yw_shrink_to_fit_impl(int *cap_inout, int len, void *old_buf, size_t item_size);

#define YW_GROW(_type, _cap_inout, _len_inout, _buf_inout, _add_len) (*(_buf_inout) = (_type *)yw_grow_impl((_cap_inout), (_len_inout), *(_buf_inout), sizeof(_type), (_add_len)))
#define YW_PUSH(_type, _cap_inout, _len_inout, _buf_inout, _value)   \
    do                                                               \
    {                                                                \
        YW_GROW(_type, (_cap_inout), (_len_inout), (_buf_inout), 1); \
        (*(_buf_inout))[*(_len_inout)-1] = (_value);                 \
    } while (0)
#define YW_ALLOC(_type, _count) (_type *)yw_alloc_impl(sizeof(_type), (_count))
#define YW_SHRINK_TO_FIT(_type, _cap_inout, _len, _buf_inout) (*(_buf_inout) = (_type *)yw_shrink_to_fit_impl((_cap_inout), (_len), *(_buf_inout), sizeof(_type)))

/* if another is NULL, this function doesn't do anything */
void yw_append_str(char **dest, char const *another);
void yw_append_char(char **dest, YW_Char32 chr);
/* if s is NULL, this function returns NULL */
char *yw_duplicate_str(char const *s);
char *yw_char_to_str(YW_Char32 chr);

#define YW_LIST_INIT(_list)                   \
    do                                        \
    {                                         \
        memset((_list), 0, sizeof(*(_list))); \
    } while (0)
#define YW_LIST_FREE(_list)   \
    do                        \
    {                         \
        free((_list)->items); \
    } while (0)
#define YW_LIST_PUSH(_type, _list, _item)                                 \
    do                                                                    \
    {                                                                     \
        YW_GROW(_type, &(_list)->cap, &(_list)->len, &(_list)->items, 1); \
        (_list)->items[(_list)->len - 1] = (_item);                       \
    } while (0)
#define YW_LIST_REMOVE(_type, _list, _index)                                   \
    do                                                                         \
    {                                                                          \
        int __index = (_index);                                                \
        if ((_list)->len <= __index)                                           \
        {                                                                      \
            fprintf(stderr, "illegal list item index %d\n", __index);          \
            abort();                                                           \
        }                                                                      \
        int copy_count = (_list)->len - 1 - __index;                           \
        for (int i = 0; i < copy_count; i++)                                   \
        {                                                                      \
            (_list)->items[__index + i] = (_list)->items[__index + i + 1];     \
        }                                                                      \
        (_list)->len--;                                                        \
        YW_SHRINK_TO_FIT(_type, &(_list)->cap, (_list)->len, &(_list)->items); \
    } while (0)
#define YW_LIST_INSERT(_type, _list, _index, _item)                        \
    do                                                                     \
    {                                                                      \
        int __index = (_index);                                            \
        if ((_list)->len < __index)                                        \
        {                                                                  \
            fprintf(stderr, "illegal list item index %d\n", __index);      \
            abort();                                                       \
        }                                                                  \
        int copy_count = (_list)->len - __index;                           \
        YW_GROW(_type, &(_list)->cap, &(_list)->len, &(_list)->items, 1);  \
        for (int i = copy_count - 1; 0 <= i; i--)                          \
        {                                                                  \
            (_list)->items[__index + i + 1] = (_list)->items[__index + i]; \
        }                                                                  \
        (_list)->items[__index] = (_item);                                 \
    } while (0)

/*******************************************************************************
 * UTF-8 character utility
 ******************************************************************************/

/*
 * Caller owns the returned string.
 * If an error was encountered, ? is returned instead.
 * If chr is 0, this function will abort().
 */
char *yw_char_to_utf8(YW_Char32 chr);

/*
 * Returns resulting codepoint, or 0 if end was reached.
 * If an error was encountered, U+FFFD is returned instead.
 */
YW_Char32 yw_utf8_next_char(char const **str);

/* Caller owns the returned array. */
void yw_utf8_to_char32(YW_Char32 **chars_out, int *chars_len_out, char const *str);

/* UTF-8-aware version of strlen(). */
size_t yw_utf8_strlen(char const *s);

/*
 * UTF-8-aware version of strchr().
 *
 * NOTE: For searching non-unicode character, using normal strchr() is OK.
 *       (And may even be faster)
 */
char const *yw_utf8_strchr(char const *s, YW_Char32 c);

/*******************************************************************************
 * YW_TextReader
 ******************************************************************************/

typedef uint8_t const *YW_TextCursor;

typedef struct YW_TextReader
{
    uint8_t const *chars;
    int chars_len;
    uint8_t const *cursor;
} YW_TextReader;

void yw_text_reader_init(YW_TextReader *out, uint8_t const *chars, int chars_len);
void yw_text_reader_init_from_str(YW_TextReader *out, char const *s);
int yw_text_reader_cursor(YW_TextReader const *tr);
bool yw_text_reader_is_eof(YW_TextReader const *tr);

/* Returns -1 on EOF. */
YW_Char32 yw_peek_char(YW_TextReader const *tr);

/* Returns -1 on EOF. */
YW_Char32 yw_consume_any_char(YW_TextReader *tr);

/*
 * Returns -1 on EOF or when no match was found.
 * Also note that this function can only match ASCII characters.
 */
int yw_consume_one_of_chars(YW_TextReader *tr, char const *chars);
bool yw_consume_char(YW_TextReader *tr, YW_Char32 chr);

typedef enum
{
    YW_NO_MATCH_FLAGS = 0,
    YW_ASCII_CASE_INSENSITIVE = 1 << 0
} YW_MatchFlags;

/*
 * Returns index of matched string, or -1 if not found.
 *
 * strs must be NULL-terminated list!
 */
int yw_consume_one_of_strs(YW_TextReader *tr, char const **strs, YW_MatchFlags flags);
bool yw_consume_str(YW_TextReader *tr, char const *str, YW_MatchFlags flags);

/*******************************************************************************
 * Garbage collector
 ******************************************************************************/

typedef void *YW_PtrSlot;

/*
 * Each slot may store either a pointer or NULL.
 * NULL means free "slot", and new pointers can be stored there.
 */
typedef struct YW_PtrCollection
{
    YW_PtrSlot *slots;
    int slots_len;
    int slots_cap;
} YW_PtrCollection;

void YW_PtrCollection_init(YW_PtrCollection *out);
void YW_PtrCollection_deinit(YW_PtrCollection *coll);

/* Returns pointer to the slot */
YW_PtrSlot *yw_add_ptr_to_collection(YW_PtrCollection *coll, void *obj);

typedef struct YW_GcCallbacks
{
    void (*visit)(void *self);
    void (*destroy)(void *self);
} YW_GcCallbacks;

typedef struct YW_GcObjectHeader
{
    uint32_t magic_and_marked_flag; /* LSB is used as marked flag. */
    YW_GcCallbacks const *callbacks;
} YW_GcObjectHeader;

typedef struct YW_GcHeap
{
    YW_PtrCollection all_objs;
    YW_PtrCollection root_objs;
} YW_GcHeap;

typedef enum
{
    YW_NO_GC_ALLOC_FLAGS = 0,
    YW_GC_ROOT_OBJECT = 1 << 0
} YW_GcAllocFlags;

/* NOTE: It is safe to pass NULL pointer. */
void yw_gc_visit(void *obj_v);

#define YW_GC_TYPE(_x) _x##_GC
#define YW_GC_PTR(_x) YW_GC_TYPE(_x) *

void yw_gc_heap_init(YW_GcHeap *out);
void yw_gc_heap_deinit(YW_GcHeap *heap);

void *yw_gc_alloc_impl(YW_GcHeap *heap, int size, YW_GcCallbacks const *callbacks, YW_GcAllocFlags alloc_flags);
#define YW_GC_ALLOC(_type, _heap, _callbacks, _alloc_flags) (YW_GC_PTR(_type))yw_gc_alloc_impl((_heap), sizeof(YW_GC_TYPE(_type)), (_callbacks), (_alloc_flags))

void yw_gc(YW_GcHeap *heap);

/*******************************************************************************
 * Other small utilities
 ******************************************************************************/

#define YW_CLAMP(_n, _min, _max) ((_n) < (_min)) ? (_min) : (((_max) < (_n)) ? (_max) : (_n))

#endif /* #ifndef YW_COMMON_H_ */
