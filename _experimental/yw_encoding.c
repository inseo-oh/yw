/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_encoding.h"
#include "yw_common.h"
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

typedef YW_TextDecoder *(YW_TextDecoderFactory)();

static struct
{
    YW_EncodingType encoding;
    YW_TextDecoder *decoder_factory;
} yw_encodings[] = {};

void yw_encoding_decode(YW_IoQueue *input,
                        YW_EncodingType fallback_encoding_type,
                        YW_IoQueue *output)
{
    /* https://encoding.spec.whatwg.org/#decode */

    YW_TODO();
}

void yw_bom_sniff(YW_IoQueue *output)
{
    /* https://encoding.spec.whatwg.org/#bom-sniff */

    YW_TODO();
}

/* Caller owns the returned array. */
int *yw_io_queue_items_to_array(YW_IoQueueItems const *items)
{
    int *res_buf = NULL;
    int len = 0;
    int cap = 0;

    for (int i = 0; i < items->len; i++)
    {
        res_buf = YW_GROW(int, &cap, &len, res_buf);
        res_buf[len - 1] = items->items[i];
    }
    res_buf = YW_SHRINK_TO_FIT(int, &cap, len, res_buf);
    return res_buf;
}

/* Caller owns the returned array. */
int *yw_io_queue_to_array(YW_IoQueue const *queue)
{
    int *res_buf = NULL;
    int len = 0;
    int cap = 0;

    for (int i = 0; i < queue->len; i++)
    {
        res_buf = YW_GROW(int, &cap, &len, res_buf);
        res_buf[len - 1] = queue->items[i];
    }
    res_buf = YW_SHRINK_TO_FIT(int, &cap, len, res_buf);
    return res_buf;
}

void yw_io_queue_from_array(YW_IoQueue *out, int const *items, int items_len)
{
    memset(out, 0, sizeof(*out));
    YW_LIST_INIT(out);
    for (int i = 0; i < items_len; i++)
    {
        YW_LIST_PUSH(YW_IoQueueItem, out, (YW_IoQueueItem)items[i]);
    }
    YW_LIST_PUSH(YW_IoQueueItem, out, YW_END_OF_IO_QUEUE);
}
