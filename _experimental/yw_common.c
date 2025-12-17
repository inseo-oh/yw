/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_common.h"
#include <limits.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*******************************************************************************
 * ASCII character conversion & testing
 ******************************************************************************/

bool yw_is_leading_surrogate_char(YW_Char32 c)
{
    return 0xd800 <= c && c <= 0xdbff;
}
bool yw_is_trailing_surrogate_char(YW_Char32 c)
{
    return 0xdc00 <= c && c <= 0xdfff;
}
bool yw_is_surrogate_char(YW_Char32 c)
{
    return yw_is_leading_surrogate_char(c) || yw_is_trailing_surrogate_char(c);
}
bool yw_is_c0_control_char(YW_Char32 c)
{
    return 0x0000 <= c && c <= 0x001f;
}
bool yw_is_control_char(YW_Char32 c)
{
    return yw_is_c0_control_char(c) || (0x007f <= c && c <= 0x009f);
}
bool yw_is_ascii_digit(YW_Char32 c)
{
    return '0' <= c && c <= '9';
}
bool yw_is_ascii_uppercase(YW_Char32 c)
{
    return 'A' <= c && c <= 'Z';
}
bool yw_is_ascii_lowercase(YW_Char32 c)
{
    return 'a' <= c && c <= 'z';
}
bool yw_is_ascii_alpha(YW_Char32 c)
{
    return yw_is_ascii_uppercase(c) || yw_is_ascii_lowercase(c);
}
bool yw_is_ascii_alphanumeric(YW_Char32 c)
{
    return yw_is_ascii_alpha(c) || yw_is_ascii_digit(c);
}
bool yw_is_ascii_uppercase_hex_digit(YW_Char32 c)
{
    return 'A' <= c && c <= 'F';
}
bool yw_is_ascii_lowercase_hex_digit(YW_Char32 c)
{
    return 'a' <= c && c <= 'f';
}
bool yw_is_ascii_hex_digit(YW_Char32 c)
{
    return yw_is_ascii_uppercase_hex_digit(c) ||
           yw_is_ascii_lowercase_hex_digit(c) || yw_is_ascii_digit(c);
}
bool yw_is_ascii_whitespace(YW_Char32 c)
{
    switch (c)
    {
    case 0x0009:
    case 0x000a:
    case 0x000c:
    case 0x000d:
        return true;
    }
    return false;
}
bool yw_is_noncharacter(YW_Char32 c)
{
    switch (c)
    {
    case 0xfffe:
    case 0xffff:
    case 0x1fffe:
    case 0x1ffff:
    case 0x2fffe:
    case 0x2ffff:
    case 0x3fffe:
    case 0x3ffff:
    case 0x4fffe:
    case 0x4ffff:
    case 0x5fffe:
    case 0x5ffff:
    case 0x6fffe:
    case 0x6ffff:
    case 0x7fffe:
    case 0x7ffff:
    case 0x8fffe:
    case 0x8ffff:
    case 0x9fffe:
    case 0x9ffff:
    case 0xafffe:
    case 0xaffff:
    case 0xbfffe:
    case 0xbffff:
    case 0xcfffe:
    case 0xcffff:
    case 0xdfffe:
    case 0xdffff:
    case 0xefffe:
    case 0xeffff:
    case 0xffffe:
    case 0xfffff:
    case 0x10fffe:
    case 0x10ffff:
        return true;
    }
    return false;
}
YW_Char32 yw_to_ascii_lowercase(YW_Char32 c)
{
    if (!yw_is_ascii_uppercase(c))
    {
        return c;
    }
    return c - 'A' + 'a';
}
YW_Char32 yw_to_ascii_uppercase(YW_Char32 c)
{
    if (!yw_is_ascii_lowercase(c))
    {
        return c;
    }
    return c - 'a' + 'A';
}

/*******************************************************************************
 * Memory utilities
 ******************************************************************************/

void *yw_grow_impl(int *cap_inout, int *len_inout, void *old_buf,
                   size_t item_size)
{
    if (*cap_inout < 0 || *len_inout < 0 || item_size == 0)
    {
        fprintf(stderr, "%s: Illegal item size, length or capacity detected",
                __func__);
        abort();
    }
    int new_len = *len_inout + 1;
    if (new_len <= *cap_inout)
    {
        *len_inout = new_len;
        return old_buf;
    }
    int new_cap = new_len * 2;
    void *new_buf = realloc(old_buf, new_cap * item_size);
    if (new_buf == NULL)
    {
        /* If we don't have enough space for that much memory, give space for
         * at least one new item. */
        new_cap = new_len;
        new_buf = realloc(old_buf, new_cap * item_size);
    }
    if (new_buf == NULL)
    {
        fprintf(stderr, "%s: out of memory", __func__);
        abort();
    }
    *len_inout = new_len;
    *cap_inout = new_cap;
    return new_buf;
}
void *yw_shrink_to_fit_impl(int *cap_inout, int len, void *old_buf,
                            size_t item_size)
{
    if (*cap_inout < 0 || len < 0 || item_size == 0)
    {
        fprintf(stderr, "%s: Illegal item size, length or capacity detected",
                __func__);
        abort();
    }
    if (len == *cap_inout)
    {
        return old_buf;
    }
    int new_cap = len;
    void *new_buf = realloc(old_buf, new_cap * item_size);
    if (new_buf == NULL)
    {
        return old_buf;
    }
    *cap_inout = new_cap;
    return new_buf;
}

void yw_append_str(char **dest, char const *another)
{
    bool was_dest_null = false;
    if (another == NULL)
    {
        return;
    }
    int len = 0;
    if (*dest != NULL)
    {
        len = strlen(*dest) + 1;
    }
    else
    {
        len = 1;
        was_dest_null = true;
    }
    int cap = len;
    int a_len = len - 1;
    size_t b_len = strlen(another);
    if (INT_MAX < b_len)
    {
        fprintf(stderr, "%s: String is too long", __func__);
        abort();
    }
    while (len < a_len + (int)b_len + 1)
    {
        *dest = YW_GROW(char, &cap, &len, *dest);
    }
    if (was_dest_null)
    {
        (*dest)[0] = '\0';
    }
    strcat(*dest, another);
}

char *yw_duplicate_str(char const *s)
{
    if (s == NULL)
    {
        return NULL;
    }
    char *res = NULL;
    yw_append_str(&res, s);
    return res;
}

/*******************************************************************************
 * UTF-8 character utility
 ******************************************************************************/

YW_Char32 yw_utf8_next_char(char const **s)
{
    uint8_t bytes_seen = 0;
    uint8_t bytes_needed = 0;
    uint8_t lower_boundary = 0x80;
    uint8_t upper_boundary = 0xbf;
    uint32_t codepoint;

    while (1)
    {
        /*
         * Decoding algorithm taken from:
         * https://encoding.spec.whatwg.org/#utf-8-decoder
         */

        if (**s == 0)
        {
            if (bytes_needed != 0)
            {
                bytes_needed = 0;
                return 0xfffd;
            }
            else
            {
                return 0;
            }
        }

        uint8_t byte = **s;
        (*s)++;
        if (bytes_needed == 0)
        {
            if (byte <= 0x7f)
            {
                return byte;
            }
            else if (0xc2 <= byte && byte <= 0xdf)
            {
                bytes_needed = 1;
                codepoint = byte & 0x1f;
            }
            else if (0xe0 <= byte && byte <= 0xef)
            {
                switch (byte)
                {
                case 0xe0:
                    lower_boundary = 0xa0;
                    break;
                case 0xed:
                    upper_boundary = 0x9f;
                    break;
                }
                bytes_needed = 2;
                codepoint = byte & 0xf;
            }
            else if (0xf0 <= byte && byte <= 0xf4)
            {
                switch (byte)
                {
                case 0xe0:
                    lower_boundary = 0x90;
                    break;
                case 0xed:
                    upper_boundary = 0x8f;
                    break;
                }
                bytes_needed = 3;
                codepoint = byte & 0x7;
            }
            else
            {
                return 0xfffd;
            }
            continue;
        }
        if (byte < lower_boundary || upper_boundary < byte)
        {
            codepoint = 0;
            bytes_needed = 0;
            bytes_seen = 0;
            lower_boundary = 0x80;
            upper_boundary = 0xbf;
            return 0xfffd;
        }
        lower_boundary = 0x80;
        upper_boundary = 0xbf;
        codepoint = (codepoint << 6) | (byte & 0x3f);
        bytes_seen++;
        if (bytes_seen == bytes_needed)
        {
            break;
        }
    }
    bytes_needed = 0;
    bytes_seen = 0;
    return codepoint;
}

void yw_utf8_to_char32(YW_Char32 **chars_out, int *chars_len_out,
                       char const *str)
{
    YW_Char32 *res_buf = NULL;
    int res_len = 0;
    int res_cap = 0;
    char const *next_str = str;

    while (1)
    {
        YW_Char32 chr = yw_utf8_next_char(&next_str);
        if (chr == 0)
        {
            break;
        }
        res_buf = YW_GROW(YW_Char32, &res_cap, &res_len, res_buf);
        res_buf[res_len - 1] = chr;
    }
    res_buf = YW_SHRINK_TO_FIT(YW_Char32, &res_cap, res_len, res_buf);
    *chars_out = res_buf;
    *chars_len_out = res_len;
}

size_t yw_utf8_strlen(char const *s)
{
    char const *next_str = s;
    size_t len = 0;
    while (1)
    {
        YW_Char32 got = yw_utf8_next_char(&next_str);
        if (got == 0)
        {
            break;
        }
        len++;
    }
    return len;
}

char const *yw_utf8_strchr(char const *s, YW_Char32 c)
{
    char const *next_str = s;
    while (1)
    {
        char const *res_str = next_str;
        YW_Char32 got = yw_utf8_next_char(&next_str);
        if (got == 0 && c != '\0')
        {
            return NULL;
        }
        else if (got == c)
        {
            return res_str;
        }
    }
    return NULL;
}

/*******************************************************************************
 * YW_TextReader
 ******************************************************************************/

void YW_TextReader_init(YW_TextReader *out, char const *source_name,
                         YW_Char32 const *chars, int chars_len)
{
    memset(out, 0, sizeof(*out));
    out->source_name = source_name;
    out->chars = chars;
    out->chars_len = chars_len;
}

bool YW_TextReader_is_eof(YW_TextReader const *tr)
{
    return tr->chars_len <= tr->cursor;
}

YW_Char32 yw_peek_char(YW_TextReader const *tr)
{
    if (YW_TextReader_is_eof(tr))
    {
        return -1;
    }
    return tr->chars[tr->cursor];
}

YW_Char32
yw_consume_any_char(YW_TextReader *tr)
{
    if (YW_TextReader_is_eof(tr))
    {
        return -1;
    }
    YW_Char32 res = yw_peek_char(tr);
    tr->cursor++;
    return res;
}

int yw_consume_one_of_chars(YW_TextReader *tr, char const *chars)
{
    if (YW_TextReader_is_eof(tr))
    {
        return -1;
    }
    YW_Char32 got = yw_peek_char(tr);
    for (char const *char_src = chars; *char_src != 0; char_src++)
    {
        if (*char_src == got)
        {
            yw_consume_any_char(tr);
            return got;
        }
    }
    return -1;
}

bool yw_consume_char(YW_TextReader *tr, YW_Char32 chr)
{
    if (YW_TextReader_is_eof(tr))
    {
        return false;
    }
    YW_Char32 got = yw_peek_char(tr);
    if (got == chr)
    {
        yw_consume_any_char(tr);
        return true;
    }
    return false;
}

int yw_consume_one_of_strs(YW_TextReader *tr, char const **strs,
                           YW_MatchFlags flags)
{
    if (YW_TextReader_is_eof(tr))
    {
        return -1;
    }

    bool found = false;
    int match_idx = -1;
    int match_len = 0;

    for (char const **src_str = strs; !found && (*src_str != NULL); src_str++)
    {
        match_idx++;
        found = true;
        match_len = 0;
        for (int i = tr->cursor; i < tr->chars_len; i++)
        {
            char src_chr = (*src_str)[i - tr->cursor];
            char got_chr = tr->chars[i];

            if (src_chr == 0)
            {
                break;
            }
            if (flags & YW_ASCII_CASE_INSENSITIVE)
            {
                src_chr = yw_to_ascii_lowercase(src_chr);
                got_chr = yw_to_ascii_lowercase(got_chr);
            }
            if (src_chr != got_chr)
            {
                /* Mismatch found */
                found = false;
                break;
            }
            match_len++;
        }
    }
    if (!found)
    {
        return -1;
    }
    tr->cursor += match_len;
    return match_idx;
}

bool yw_consume_str(YW_TextReader *tr, char const *str, YW_MatchFlags flags)
{
    char const *strs[] = {str, NULL};
    int res = yw_consume_one_of_strs(tr, strs, flags);
    return res == 0;
}

/*******************************************************************************
 * Garbage collector
 ******************************************************************************/

void YW_PtrCollection_init(YW_PtrCollection *out)
{
    memset(out, 0, sizeof(*out));
}

void YW_PtrCollection_deinit(YW_PtrCollection *coll)
{
    free(coll->slots);
}

YW_PtrSlot *yw_add_ptr_to_collection(YW_PtrCollection *coll, void *obj)
{
    for (int i = 0; i < coll->slots_len; i++)
    {
        if (coll->slots[i] == NULL)
        {
            coll->slots[i] = obj;
            return &coll->slots[i];
        }
    }
    coll->slots =
        YW_GROW(YW_PtrSlot, &coll->slots_cap, &coll->slots_len, coll->slots);
    coll->slots[coll->slots_len - 1] = obj;
    return &coll->slots[coll->slots_len - 1];
}

/* NOTE: LSB must be zero -- It is used as "marked" flag for GC. */
#define YW_GC_MAGIC 0x21b0fb278bf5e5ce

static bool yw_gc_is_marked(YW_GcObjectHeader const *obj)
{
    return (obj->magic_and_marked_flag & 0x1) != 0;
}
static void yw_gc_mark_object(YW_GcObjectHeader *obj)
{
    obj->magic_and_marked_flag |= 0x1;
}
static void yw_gc_unmark_object(YW_GcObjectHeader *obj)
{
    obj->magic_and_marked_flag &= ~0x1;
}
void yw_gc_visit(void *obj_v)
{
    YW_GcObjectHeader *obj = (YW_GcObjectHeader *)obj_v;
    if (obj == NULL)
    {
        return;
    }
    if ((obj->magic_and_marked_flag & ~0x1) != YW_GC_MAGIC)
    {
        fprintf(stderr, "%s: Object at %p has corrupted magic!\n", __func__,
                (void *)obj);
        abort();
    }
    yw_gc_unmark_object(obj);
    if (obj->callbacks != NULL && obj->callbacks->visit != NULL)
    {
        obj->callbacks->visit(obj);
    }
}

void yw_gc_heap_init(YW_GcHeap *out)
{
    memset(out, 0, sizeof(*out));
    YW_PtrCollection_init(&out->all_objs);
    YW_PtrCollection_init(&out->root_objs);
}
void yw_gc_heap_deinit(YW_GcHeap *heap)
{
    YW_PtrCollection_deinit(&heap->root_objs);
    YW_PtrCollection_deinit(&heap->all_objs);
}

void *yw_gc_alloc_impl(YW_GcHeap *heap, int size,
                       YW_GcCallbacks const *callbacks,
                       YW_GcAllocFlags alloc_flags)
{
    if (size < (int)sizeof(YW_GcObjectHeader))
    {
        printf("%s: illegal size %d!\n", __func__, size);
        abort();
    }
    void *mem = malloc(size);
    memset(mem, 0, size);

    YW_GcObjectHeader *mem_header = (YW_GcObjectHeader *)mem;
    mem_header->magic_and_marked_flag = YW_GC_MAGIC;
    mem_header->callbacks = callbacks;

    yw_add_ptr_to_collection(&heap->all_objs, mem);
    if (alloc_flags & YW_GC_ROOT_OBJECT)
    {
        yw_add_ptr_to_collection(&heap->root_objs, mem);
    }
    return mem;
}

void yw_gc(YW_GcHeap *heap)
{
    /* 1. Mark all objects ****************************************************/
    for (int i = 0; i < heap->all_objs.slots_len; i++)
    {
        if (heap->all_objs.slots[i] == NULL)
        {
            continue;
        }
        YW_GcObjectHeader *obj = (YW_GcObjectHeader *)heap->all_objs.slots[i];
        if ((obj->magic_and_marked_flag & ~0x1) != YW_GC_MAGIC)
        {
            printf("WARNING: %s: Object at %p has corrupted magic!\n", __func__,
                   (void *)obj);
            continue;
        }
        yw_gc_mark_object(obj);
    }
    /* 2. Visit root objects **************************************************/
    for (int i = 0; i < heap->root_objs.slots_len; i++)
    {
        if (heap->root_objs.slots[i] == NULL)
        {
            continue;
        }
        yw_gc_visit(heap->root_objs.slots[i]);
    }
    /* 3. Destroy still marked objects ****************************************/
    for (int i = 0; i < heap->all_objs.slots_len; i++)
    {
        if (heap->all_objs.slots[i] == NULL)
        {
            continue;
        }
        YW_GcObjectHeader *obj = (YW_GcObjectHeader *)heap->all_objs.slots[i];
        if ((obj->magic_and_marked_flag & ~0x1) != YW_GC_MAGIC)
        {
            printf("WARNING: %s: Object at %p has corrupted magic!\n", __func__,
                   (void *)obj);
            continue;
        }
        if (yw_gc_is_marked(obj))
        {
            if (obj->callbacks != NULL && obj->callbacks->destroy != NULL)
            {
                obj->callbacks->destroy(obj);
            }
            free(obj);
            heap->all_objs.slots[i] = NULL;
        }
    }
}
