/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_ENCODING_H_
#define YW_ENCODING_H_
#include "yw_common.h"
#include <stdint.h>

typedef enum
{
    YW_INVALID_ENCODING,
    YW_UTF8,           /* https://encoding.spec.whatwg.org/#utf-8 */
    YW_IBM866,         /* https://encoding.spec.whatwg.org/#ibm866 */
    YW_ISO8859_2,      /* https://encoding.spec.whatwg.org/#iso-8859-2 */
    YW_ISO8859_3,      /* https://encoding.spec.whatwg.org/#iso-8859-3 */
    YW_ISO8859_4,      /* https://encoding.spec.whatwg.org/#iso-8859-4 */
    YW_ISO8859_5,      /* https://encoding.spec.whatwg.org/#iso-8859-5 */
    YW_ISO8859_6,      /* https://encoding.spec.whatwg.org/#iso-8859-6 */
    YW_ISO8859_7,      /* https://encoding.spec.whatwg.org/#iso-8859-7 */
    YW_ISO8859_8,      /* https://encoding.spec.whatwg.org/#iso-8859-8 */
    YW_ISO8859_8I,     /* https://encoding.spec.whatwg.org/#iso-8859-8-i */
    YW_ISO8859_10,     /* https://encoding.spec.whatwg.org/#iso-8859-10 */
    YW_ISO8859_13,     /* https://encoding.spec.whatwg.org/#iso-8859-13 */
    YW_ISO8859_14,     /* https://encoding.spec.whatwg.org/#iso-8859-14 */
    YW_ISO8859_15,     /* https://encoding.spec.whatwg.org/#iso-8859-15 */
    YW_ISO8859_16,     /* https://encoding.spec.whatwg.org/#iso-8859-16 */
    YW_KOI8R,          /* https://encoding.spec.whatwg.org/#koi8-r */
    YW_KOI8U,          /* https://encoding.spec.whatwg.org/#koi8-u */
    YW_MACINTOSH,      /* https://encoding.spec.whatwg.org/#macintosh */
    YW_WINDOWS874,     /* https://encoding.spec.whatwg.org/#windows-874 */
    YW_WINDOWS1250,    /* https://encoding.spec.whatwg.org/#windows-1250 */
    YW_WINDOWS1251,    /* https://encoding.spec.whatwg.org/#windows-1251 */
    YW_WINDOWS1252,    /* https://encoding.spec.whatwg.org/#windows-1252 */
    YW_WINDOWS1253,    /* https://encoding.spec.whatwg.org/#windows-1253 */
    YW_WINDOWS1254,    /* https://encoding.spec.whatwg.org/#windows-1254 */
    YW_WINDOWS1255,    /* https://encoding.spec.whatwg.org/#windows-1255 */
    YW_WINDOWS1256,    /* https://encoding.spec.whatwg.org/#windows-1256 */
    YW_WINDOWS1257,    /* https://encoding.spec.whatwg.org/#windows-1257 */
    YW_WINDOWS1258,    /* https://encoding.spec.whatwg.org/#windows-1258 */
    YW_X_MAC_CYRILLIC, /* https://encoding.spec.whatwg.org/#x-mac-cyrillic */
    YW_GBK,            /* https://encoding.spec.whatwg.org/#gbk */
    YW_GB18030,        /* https://encoding.spec.whatwg.org/#gb18030 */
    YW_BIG5,           /* https://encoding.spec.whatwg.org/#big5 */
    YW_EUC_JP,         /* https://encoding.spec.whatwg.org/#euc-jp */
    YW_ISO2022_JP,     /* https://encoding.spec.whatwg.org/#iso-2022-jp */
    YW_SHIFT_JIS,      /* https://encoding.spec.whatwg.org/#shift_jis */
    YW_EUC_KR,         /* https://encoding.spec.whatwg.org/#euc-kr */
    YW_REPLACEMENT,    /* https://encoding.spec.whatwg.org/#replacement */
    YW_UTF16_BE,       /* https://encoding.spec.whatwg.org/#utf-16be */
    YW_UTF16_LE,       /* https://encoding.spec.whatwg.org/#utf-16le */
    YW_X_USER_DEFINED, /* https://encoding.spec.whatwg.org/#x-user-defined */
} YW_EncodingType;

/*
 * Positive values are normal byte or codepoint values, and -1 is special
 * value for "end-of-queue".
 */
typedef int YW_IOQueueItem;
#define YW_END_OF_IO_QUEUE -1

typedef struct YW_IOQueueItemList
{
    YW_IOQueueItem *items;
    int len, cap;
} YW_IOQueueItemList;
typedef struct YW_IOQueue
{
    YW_IOQueueItemList item_list;
} YW_IOQueue;

/* Returns YW_INVALID_ENCODING if there's no corresponding encoding. */
YW_EncodingType yw_encoding_from_label(char const *label);

void yw_encoding_decode(YW_IOQueue *input, YW_EncodingType fallback_encoding_type, YW_IOQueue *output);

/* Caller owns the returned array. */
void yw_io_queue_item_list_to_items(int **items_out, int *len_out, YW_IOQueueItemList const *list);

/* Caller owns the returned array. */
void yw_io_queue_to_items(int **items_out, int *len_out, YW_IOQueue const *queue);

/* Caller owns the returned array. */
void yw_io_queue_to_utf8(uint8_t **chars_out, int *len_out, YW_IOQueue const *queue);

void yw_io_queue_init(YW_IOQueue *out);
void yw_io_queue_init_with_items(YW_IOQueue *out, int const *items, int items_len);
void yw_io_queue_deinit(YW_IOQueue *queue);
YW_IOQueueItem yw_io_queue_read_one(YW_IOQueue *queue);
int yw_io_queue_read(YW_IOQueue *queue, int *buf, int max_len);
int yw_io_queue_peek(YW_IOQueue const *queue, int *buf, int max_len);
void yw_io_queue_push_one(YW_IOQueue *queue, YW_IOQueueItem item);
void yw_io_queue_push(YW_IOQueue *queue, YW_IOQueueItem const *items, int len);
void yw_io_queue_restore_one(YW_IOQueue *queue, YW_IOQueueItem item);
void yw_io_queue_restore(YW_IOQueue *queue, YW_IOQueueItem const *items, int len);

#define YW_IO_QUEUE_INIT_FROM_ARRAY(_out, _array) yw_io_queue_init_with_items((_out), (_array), YW_SIZEOF_ARRAY(_array))
#define YW_IO_QUEUE_READ_TO_ARRAY(_queue, _array) yw_io_queue_read((_queue), (_array), YW_SIZEOF_ARRAY(_array))
#define YW_IO_QUEUE_PEEK_TO_ARRAY(_queue, _array) yw_io_queue_peek((_queue), (_array), YW_SIZEOF_ARRAY(_array))
#define YW_IO_QUEUE_PUSH_FROM_ARRAY(_queue, _array) yw_io_queue_push((_queue), (_array), YW_SIZEOF_ARRAY(_array))
#define YW_IO_QUEUE_RESTORE_FROM_ARRAY(_queue, _array) yw_io_queue_restore((_queue), (_array), YW_SIZEOF_ARRAY(_array))

typedef enum
{
    /*
     * Positive values are OK result(resulting codepoint or byte), and negative
     * values are special results. (see below)
     */
    /* https://encoding.spec.whatwg.org/#error */
    YW_ENCODING_RESULT_ERROR = -99,
    /* https://encoding.spec.whatwg.org/#finished */
    YW_ENCODING_RESULT_FINISHED,
    /* https://encoding.spec.whatwg.org/#continue */
    YW_ENCODING_RESULT_CONTINUE,
} YW_EncodingResult;

typedef struct YW_TextDecoderCallbacks
{
    YW_EncodingResult (*handler)(void *self_v, YW_IOQueue *queue, int byte_item);
    void (*destroy)(void *self_v);
} YW_TextDecoderCallbacks;
typedef struct YW_TextDecoder
{
    void *state; /* Decoder state -- will be free()d after use */
    YW_TextDecoderCallbacks const *callbacks;
} YW_TextDecoder;

typedef enum
{
    YW_ERROR_MODE_REPLACEMENT,
    YW_ERROR_MODE_HTML,
    YW_ERROR_MODE_FATAL,
} YW_EncodingErrorMode;

/* Returns YW_INVALID_ENCODING if no encoding was found */
YW_EncodingType yw_bom_sniff(YW_IOQueue const *queue);

/*******************************************************************************
 * Encoding implementations
 ******************************************************************************/

/* clang-format off */
#define YW_ENUMERATE_ENCODINGS(_x)                                              \
    _x(YW_UTF8, yw_init_utf8_decoder)
/* clang-format on */

#define YW_X(_name, _decoder) void _decoder(YW_TextDecoder *out);
YW_ENUMERATE_ENCODINGS(YW_X)
#undef YW_X

#endif /* #ifndef YW_ENCODING_H_ */