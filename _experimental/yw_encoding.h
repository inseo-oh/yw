/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#ifndef YW_ENCODING_H_
#define YW_ENCODING_H_

typedef struct YW_TextDecoder YW_TextDecoder;
typedef struct YW_TextDecoderCallbacks YW_TextDecoderCallbacks;
typedef struct YW_IoQueue YW_IoQueue;
typedef struct YW_IoQueueItems YW_IoQueueItems;

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

struct YW_TextDecoderCallbacks
{
    YW_EncodingResult (*handler)(YW_IoQueue *queue, int byte_item);
    void (*destroy)(void *self_v);
};
struct YW_TextDecoder
{
    void *data;
    YW_TextDecoderCallbacks const *callbacks;
};

typedef enum
{
    /*
     * Positive values are normal byte or codepoint values, and -1 is special
     * value for "end-of-queue".
     */
    YW_END_OF_IO_QUEUE = -1
} YW_IoQueueItem;

struct YW_IoQueue
{
    YW_IoQueueItem *items;
    int len, cap;
};
struct YW_IoQueueItems
{
    YW_IoQueueItem *items;
    int len, cap;
};

/* Returns YW_INVALID_ENCODING if there's no corresponding encoding. */
YW_EncodingType yw_encoding_from_label(char const *label);

#endif /* #ifndef YW_ENCODING_H_ */