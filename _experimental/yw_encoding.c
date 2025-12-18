/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_encoding.h"
#include "yw_common.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static struct
{
    char const *label;
    YW_EncodingType encoding_type;
} yw_encoding_labels[] = {
    {"unicode-1-1-utf-8", YW_UTF8},
    {"unicode11utf8", YW_UTF8},
    {"unicode20utf8", YW_UTF8},
    {"utf-8", YW_UTF8},
    {"utf8", YW_UTF8},
    {"x-unicode20utf8", YW_UTF8},

    {"866", YW_IBM866},
    {"cp866", YW_IBM866},
    {"csibm866", YW_IBM866},
    {"ibm866", YW_IBM866},

    {"csisolatin2", YW_ISO8859_2},
    {"iso-8859-2", YW_ISO8859_2},
    {"iso-ir-101", YW_ISO8859_2},
    {"iso8859-2", YW_ISO8859_2},
    {"iso88592", YW_ISO8859_2},
    {"iso_8859-2", YW_ISO8859_2},
    {"iso_8859-2:1987", YW_ISO8859_2},
    {"l2", YW_ISO8859_2},
    {"latin2", YW_ISO8859_2},

    {"csisolatin3", YW_ISO8859_3},
    {"iso-8859-3", YW_ISO8859_3},
    {"iso-ir-109", YW_ISO8859_3},
    {"iso8859-3", YW_ISO8859_3},
    {"iso88593", YW_ISO8859_3},
    {"iso_8859-3", YW_ISO8859_3},
    {"iso_8859-3:1988", YW_ISO8859_3},
    {"l3", YW_ISO8859_3},
    {"latin3", YW_ISO8859_3},

    {"csisolatin4", YW_ISO8859_4},
    {"iso-8859-4", YW_ISO8859_4},
    {"iso-ir-110", YW_ISO8859_4},
    {"iso8859-4", YW_ISO8859_4},
    {"iso88594", YW_ISO8859_4},
    {"iso_8859-4", YW_ISO8859_4},
    {"iso_8859-4:1988", YW_ISO8859_4},
    {"l4", YW_ISO8859_4},
    {"latin4", YW_ISO8859_4},

    {"csisolatincyrillic", YW_ISO8859_5},
    {"cyrillic", YW_ISO8859_5},
    {"iso-8859-5", YW_ISO8859_5},
    {"iso-ir-144", YW_ISO8859_5},
    {"iso8859-5", YW_ISO8859_5},
    {"iso88595", YW_ISO8859_5},
    {"iso_8859-5", YW_ISO8859_5},
    {"iso_8859-5:1988", YW_ISO8859_5},

    {"arabic", YW_ISO8859_6},
    {"asmo-708", YW_ISO8859_6},
    {"csiso88596e", YW_ISO8859_6},
    {"csiso88596i", YW_ISO8859_6},
    {"csisolatinarabic", YW_ISO8859_6},
    {"ecma-114", YW_ISO8859_6},
    {"iso-8859-6", YW_ISO8859_6},
    {"iso-8859-6-e", YW_ISO8859_6},
    {"iso-8859-6-i", YW_ISO8859_6},
    {"iso-ir-127", YW_ISO8859_6},
    {"iso8859-6", YW_ISO8859_6},
    {"iso88596", YW_ISO8859_6},
    {"iso_8859-6", YW_ISO8859_6},
    {"iso_8859-6:1987", YW_ISO8859_6},

    {"csisolatingreek", YW_ISO8859_7},
    {"ecma-118", YW_ISO8859_7},
    {"elot_928", YW_ISO8859_7},
    {"greek", YW_ISO8859_7},
    {"greek8", YW_ISO8859_7},
    {"iso-8859-7", YW_ISO8859_7},
    {"iso-ir-126", YW_ISO8859_7},
    {"iso8859-7", YW_ISO8859_7},
    {"iso88597", YW_ISO8859_7},
    {"iso_8859-7", YW_ISO8859_7},
    {"iso_8859-7:1987", YW_ISO8859_7},
    {"sun_eu_greek", YW_ISO8859_7},

    {"csiso88598e", YW_ISO8859_8},
    {"csisolatinhebrew", YW_ISO8859_8},
    {"hebrew", YW_ISO8859_8},
    {"iso-8859-8", YW_ISO8859_8},
    {"iso-8859-8-e", YW_ISO8859_8},
    {"iso-ir-138", YW_ISO8859_8},
    {"iso8859-8", YW_ISO8859_8},
    {"iso88598", YW_ISO8859_8},
    {"iso_8859-8", YW_ISO8859_8},
    {"iso_8859-8:1988", YW_ISO8859_8},
    {"visual", YW_ISO8859_8},

    {"csiso88598i", YW_ISO8859_8I},
    {"iso-8859-8-i", YW_ISO8859_8I},
    {"logical", YW_ISO8859_8I},

    {"csisolatin6", YW_ISO8859_10},
    {"iso-8859-10", YW_ISO8859_10},
    {"iso-ir-157", YW_ISO8859_10},
    {"iso8859-10", YW_ISO8859_10},
    {"iso885910", YW_ISO8859_10},
    {"l6", YW_ISO8859_10},
    {"latin6", YW_ISO8859_10},

    {"iso-8859-13", YW_ISO8859_13},
    {"iso8859-13", YW_ISO8859_13},
    {"iso885913", YW_ISO8859_13},

    {"iso-8859-14", YW_ISO8859_14},
    {"iso8859-14", YW_ISO8859_14},
    {"iso885914", YW_ISO8859_14},

    {"csisolatin9", YW_ISO8859_15},
    {"iso-8859-15", YW_ISO8859_15},
    {"iso8859-15", YW_ISO8859_15},
    {"iso885915", YW_ISO8859_15},
    {"iso_8859-15", YW_ISO8859_15},
    {"l9", YW_ISO8859_15},

    {"iso-8859-16", YW_ISO8859_16},

    {"cskoi8r", YW_KOI8R},
    {"koi", YW_KOI8R},
    {"koi8", YW_KOI8R},
    {"koi8-r", YW_KOI8R},
    {"koi8_r", YW_KOI8R},

    {"koi8-ru", YW_KOI8U},
    {"koi8-u", YW_KOI8U},

    {"csmacintosh", YW_MACINTOSH},
    {"mac", YW_MACINTOSH},
    {"macintosh", YW_MACINTOSH},
    {"x-mac-roman", YW_MACINTOSH},

    {"dos-874", YW_WINDOWS874},
    {"iso-8859-11", YW_WINDOWS874},
    {"iso8859-11", YW_WINDOWS874},
    {"iso885911", YW_WINDOWS874},
    {"tis-620", YW_WINDOWS874},
    {"windows-874", YW_WINDOWS874},

    {"cp1250", YW_WINDOWS1250},
    {"windows-1250", YW_WINDOWS1250},
    {"x-cp1250", YW_WINDOWS1250},

    {"cp1251", YW_WINDOWS1251},
    {"windows-1251", YW_WINDOWS1251},
    {"x-cp1251", YW_WINDOWS1251},

    {"ansi_x3.4-1968", YW_WINDOWS1252},
    {"ascii", YW_WINDOWS1252},
    {"cp1252", YW_WINDOWS1252},
    {"cp819", YW_WINDOWS1252},
    {"csisolatin1", YW_WINDOWS1252},
    {"ibm819", YW_WINDOWS1252},
    {"iso-8859-1", YW_WINDOWS1252},
    {"iso-ir-100", YW_WINDOWS1252},
    {"iso8859-1", YW_WINDOWS1252},
    {"iso88591", YW_WINDOWS1252},
    {"iso_8859-1", YW_WINDOWS1252},
    {"iso_8859-1:1987", YW_WINDOWS1252},
    {"l1", YW_WINDOWS1252},
    {"latin1", YW_WINDOWS1252},
    {"us-ascii", YW_WINDOWS1252},
    {"windows-1252", YW_WINDOWS1252},
    {"x-cp1252", YW_WINDOWS1252},

    {"cp1253", YW_WINDOWS1253},
    {"windows-1253", YW_WINDOWS1253},
    {"x-cp1253", YW_WINDOWS1253},

    {"cp1254", YW_WINDOWS1254},
    {"csisolatin5", YW_WINDOWS1254},
    {"iso-8859-9", YW_WINDOWS1254},
    {"iso-ir-148", YW_WINDOWS1254},
    {"iso8859-9", YW_WINDOWS1254},
    {"iso88599", YW_WINDOWS1254},
    {"iso_8859-9", YW_WINDOWS1254},
    {"iso_8859-9:1989", YW_WINDOWS1254},
    {"l5", YW_WINDOWS1254},
    {"latin5", YW_WINDOWS1254},
    {"windows-1254", YW_WINDOWS1254},
    {"x-cp1254", YW_WINDOWS1254},

    {"cp1255", YW_WINDOWS1255},
    {"windows-1255", YW_WINDOWS1255},
    {"x-cp1255", YW_WINDOWS1255},

    {"cp1256", YW_WINDOWS1256},
    {"windows-1256", YW_WINDOWS1256},
    {"x-cp1256", YW_WINDOWS1256},

    {"cp1257", YW_WINDOWS1257},
    {"windows-1257", YW_WINDOWS1257},
    {"x-cp1257", YW_WINDOWS1257},

    {"cp1258", YW_WINDOWS1258},
    {"windows-1258", YW_WINDOWS1258},
    {"x-cp1258", YW_WINDOWS1258},

    {"x-mac-cyrillic", YW_X_MAC_CYRILLIC},
    {"x-mac-ukrainian", YW_X_MAC_CYRILLIC},

    {"chinese", YW_GBK},
    {"csgb2312", YW_GBK},
    {"csiso58gb231280", YW_GBK},
    {"gb2312", YW_GBK},
    {"gb_2312", YW_GBK},
    {"gb_2312-80", YW_GBK},
    {"gbk", YW_GBK},
    {"iso-ir-58", YW_GBK},
    {"x-gbk", YW_GBK},

    {"gb18030", YW_GB18030},

    {"big5", YW_BIG5},
    {"big5-hkscs", YW_BIG5},
    {"cn-big5", YW_BIG5},
    {"csbig5", YW_BIG5},
    {"x-x-big5", YW_BIG5},

    {"cseucpkdfmtjapanese", YW_EUC_JP},
    {"euc-jp", YW_EUC_JP},
    {"x-euc-jp", YW_EUC_JP},

    {"csiso2022jp", YW_ISO2022_JP},
    {"iso-2022-jp", YW_ISO2022_JP},

    {"csshiftjis", YW_SHIFT_JIS},
    {"ms932", YW_SHIFT_JIS},
    {"ms_kanji", YW_SHIFT_JIS},
    {"shift-jis", YW_SHIFT_JIS},
    {"shift_jis", YW_SHIFT_JIS},
    {"sjis", YW_SHIFT_JIS},
    {"windows-31j", YW_SHIFT_JIS},
    {"x-sjis", YW_SHIFT_JIS},

    {"cseuckr", YW_EUC_KR},
    {"csksc56011987", YW_EUC_KR},
    {"euc-kr", YW_EUC_KR},
    {"iso-ir-149", YW_EUC_KR},
    {"korean", YW_EUC_KR},
    {"ks_c_5601-1987", YW_EUC_KR},
    {"ks_c_5601-1989", YW_EUC_KR},
    {"ksc5601", YW_EUC_KR},
    {"ksc_5601", YW_EUC_KR},
    {"windows-949", YW_EUC_KR},

    {"csiso2022kr", YW_REPLACEMENT},
    {"hz-gb-2312", YW_REPLACEMENT},
    {"iso-2022-cn", YW_REPLACEMENT},
    {"iso-2022-cn-ext", YW_REPLACEMENT},
    {"iso-2022-kr", YW_REPLACEMENT},
    {"replacement", YW_REPLACEMENT},

    {"unicodefffe", YW_UTF16_BE},
    {"utf-16be", YW_UTF16_BE},

    {"csunicode", YW_UTF16_LE},
    {"iso-10646-ucs-2", YW_UTF16_LE},
    {"ucs-2", YW_UTF16_LE},
    {"unicode", YW_UTF16_LE},
    {"unicodefeff", YW_UTF16_LE},
    {"utf-16", YW_UTF16_LE},
    {"utf-16le", YW_UTF16_LE},

    {"x-user-defined", YW_X_USER_DEFINED},
};

YW_EncodingType yw_encoding_from_label(char const *label)
{
    for (int i = 0; i < (int)YW_SIZEOF_ARRAY(yw_encoding_labels); i++)
    {
        if (strcmp(yw_encoding_labels[i].label, label) == 0)
        {
            return yw_encoding_labels[i].encoding_type;
        }
    }
    return YW_INVALID_ENCODING;
}

typedef void(YW_TextDecoderFactory)(YW_TextDecoder *out);

static struct
{
    YW_EncodingType type;
    YW_TextDecoderFactory *decoder_factory;
} yw_encodings[] = {

};

static YW_EncodingResult yw_decode_item(YW_IoQueueItem item,
                                        YW_TextDecoder const *decoder,
                                        YW_IoQueue *input, YW_IoQueue *output,
                                        YW_EncodingErrorMode mode)
{
    if (mode == YW_ERROR_MODE_HTML)
    {
        fprintf(stderr, "%s: bad error mode", __func__);
        abort();
    }

    YW_EncodingResult res =
        decoder->callbacks->handler(decoder->state, input, item);
    if (res == YW_ENCODING_RESULT_FINISHED)
    {
        yw_io_queue_push_one(output, YW_END_OF_IO_QUEUE);
    }
    else if (0 <= res)
    {
        if (yw_is_surrogate_char(res))
        {
            fprintf(stderr, "%s: result cannot contain surrogate char",
                    __func__);
            abort();
        }
        yw_io_queue_push_one(output, (YW_IoQueueItem)res);
    }
    else if (res == YW_ENCODING_RESULT_ERROR)
    {
        switch (mode)
        {
        case YW_ERROR_MODE_REPLACEMENT:
            yw_io_queue_push_one(output, (YW_IoQueueItem)0xfffd);
            break;
        case YW_ERROR_MODE_HTML:
            YW_UNREACHABLE();
        case YW_ERROR_MODE_FATAL:
            return res;
        }
    }
    return YW_ENCODING_RESULT_CONTINUE;
}

static YW_EncodingResult yw_decode(YW_TextDecoder const *decoder,
                                   YW_IoQueue *input, YW_IoQueue *output,
                                   YW_EncodingErrorMode mode)
{
    while (1)
    {
        YW_IoQueueItem item = yw_io_queue_read_one(input);
        YW_EncodingResult res =
            yw_decode_item(item, decoder, input, output, mode);
        if (res != YW_ENCODING_RESULT_CONTINUE)
        {
            return res;
        }
    }

    YW_TODO();
}

void yw_encoding_decode(YW_IoQueue *input,
                        YW_EncodingType fallback_encoding_type,
                        YW_IoQueue *output)
{
    /* https://encoding.spec.whatwg.org/#decode */
    YW_EncodingType encoding_type = fallback_encoding_type;
    YW_EncodingType bom_encoding = yw_bom_sniff(input);
    if (bom_encoding != YW_INVALID_ENCODING)
    {
        encoding_type = bom_encoding;
        if (bom_encoding == YW_UTF8)
        {
            int buf[3];
            YW_IO_QUEUE_READ_TO_ARRAY(input, buf);
        }
        else
        {
            int buf[2];
            YW_IO_QUEUE_READ_TO_ARRAY(input, buf);
        }
    }

    YW_TextDecoder decoder;
    for (int i = 0; i < (int)YW_SIZEOF_ARRAY(yw_encodings); i++)
    {
        if (yw_encodings[i].type == encoding_type)
        {
            yw_encodings[i].decoder_factory(&decoder);
            if (decoder.callbacks == NULL || decoder.callbacks->handler != NULL)
            {
                fprintf(stderr,
                        "%s: BUG: returned decoder must have callbacks set, "
                        "with non-NULL handler callback",
                        __func__);
                abort();
            }
            yw_decode(&decoder, input, output, YW_ERROR_MODE_REPLACEMENT);
            if (decoder.callbacks->destroy != NULL)
            {
                decoder.callbacks->destroy(decoder.state);
            }
            free(decoder.state);
            return;
        }
    }
    fprintf(stderr, "%s: unsupported encoding", __func__);
    YW_TODO();
}

YW_EncodingType yw_bom_sniff(YW_IoQueue const *queue)
{
    /* https://encoding.spec.whatwg.org/#bom-sniff */

    int bytes[3];
    int len = YW_IO_QUEUE_PEEK_TO_ARRAY(queue, bytes);
    if (3 <= len && bytes[0] == 0xef && bytes[1] == 0xbb && bytes[2] == 0xbf)
    {
        return YW_UTF8;
    }
    else if (2 <= len && bytes[0] == 0xfe && bytes[1] == 0xff)
    {
        return YW_UTF16_BE;
    }
    else if (2 <= len && bytes[0] == 0xff && bytes[1] == 0xfe)
    {
        return YW_UTF16_LE;
    }
    return YW_INVALID_ENCODING;
}

void yw_io_queue_item_list_to_items(int **items_out, int *len_out,
                                    YW_IoQueueItemList const *list)
{
    int *res_buf = NULL;
    int len = 0;
    int cap = 0;

    for (int i = 0; i < list->len; i++)
    {
        if (list->items[i] == YW_END_OF_IO_QUEUE)
        {
            break;
        }
        res_buf = YW_GROW(int, &cap, &len, res_buf);
        res_buf[len - 1] = list->items[i];
    }
    res_buf = YW_SHRINK_TO_FIT(int, &cap, len, res_buf);
    *items_out = res_buf;
    *len_out = len;
}

void yw_io_queue_to_items(int **items_out, int *len_out,
                          YW_IoQueue const *queue)
{
    yw_io_queue_item_list_to_items(items_out, len_out, &queue->item_list);
}

void yw_io_queue_init(YW_IoQueue *out)
{
    yw_io_queue_init_with_items(out, NULL, 0);
}

void yw_io_queue_init_with_items(YW_IoQueue *out, int const *items,
                                 int items_len)
{
    memset(out, 0, sizeof(*out));
    YW_LIST_INIT(&out->item_list);
    for (int i = 0; i < items_len; i++)
    {
        YW_LIST_PUSH(YW_IoQueueItem, &out->item_list, (YW_IoQueueItem)items[i]);
    }
    YW_LIST_PUSH(YW_IoQueueItem, &out->item_list, YW_END_OF_IO_QUEUE);
}

void yw_io_queue_deinit(YW_IoQueue *queue)
{
    YW_LIST_FREE(&queue->item_list);
}

YW_IoQueueItem yw_io_queue_read_one(YW_IoQueue *queue)
{
    /* https://encoding.spec.whatwg.org/#concept-stream-read */

    if (queue->item_list.len == 0)
    {
        fprintf(stderr, "%s: queue is empty", __func__);
        abort();
    }
    YW_IoQueueItem item = queue->item_list.items[0];
    if (item == YW_END_OF_IO_QUEUE)
    {
        return item;
    }
    YW_LIST_REMOVE(YW_IoQueueItem, &queue->item_list, 0);
    return item;
}

int yw_io_queue_read(YW_IoQueue *queue, int *buf, int max_len)
{
    /* https://encoding.spec.whatwg.org/#concept-stream-read */
    int len = 0;
    for (int i = 0; i < max_len; i++)
    {
        YW_IoQueueItem item = yw_io_queue_read_one(queue);
        if (item == YW_END_OF_IO_QUEUE)
        {
            continue;
        }
        buf[len] = item;
        len++;
    }
    return len;
}

int yw_io_queue_peek(YW_IoQueue const *queue, int *buf, int max_len)
{
    /* https://encoding.spec.whatwg.org/#i-o-queue-peek */
    int len = 0;
    for (int i = 0; i < max_len; i++)
    {
        YW_IoQueueItem item = queue->item_list.items[i];
        if (item == YW_END_OF_IO_QUEUE)
        {
            break;
        }
        buf[len] = item;
        len++;
    }
    return len;
}

void yw_io_queue_push_one(YW_IoQueue *queue, YW_IoQueueItem item)
{
    if (queue->item_list.len == 0 ||
        queue->item_list.items[queue->item_list.len - 1] != YW_END_OF_IO_QUEUE)
    {
        fprintf(stderr, "%s: the last item must be end-of-queue", __func__);
        abort();
    }
    if (item == YW_END_OF_IO_QUEUE)
    {
        return;
    }
    YW_LIST_INSERT(YW_IoQueueItem, &queue->item_list, queue->item_list.len - 1,
                   item);
}

void yw_io_queue_push(YW_IoQueue *queue, YW_IoQueueItem const *items, int len)
{
    for (int i = 0; i < len; i++)
    {
        yw_io_queue_push_one(queue, items[i]);
    }
}

void yw_io_queue_restore_one(YW_IoQueue *queue, YW_IoQueueItem item)
{
    if (item == YW_END_OF_IO_QUEUE)
    {
        fprintf(stderr, "%s: attempted to restore end-of-queue item", __func__);
        abort();
    }
    YW_LIST_INSERT(YW_IoQueueItem, &queue->item_list, 0, item);
}

void yw_io_queue_restore(YW_IoQueue *queue, YW_IoQueueItem const *items,
                         int len)
{
    for (int i = 0; i < len; i++)
    {
        yw_io_queue_restore_one(queue, items[i]);
    }
}
