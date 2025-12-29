/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#include "yw_json.h"
#include "yw_common.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

union YW_JSONValue {
    YW_JSONValueType type;

    YW_JSONObjectValue object_val;
    YW_JSONArrayValue array_val;
    YW_JSONNumberValue number_val;
    YW_JSONStringValue string_val;
    YW_JSONBooleanValue boolean_val;
};

void yw_json_string_init(YW_JSONString *out, char const *s)
{
    /* FIXME: yw_duplicate_str will include NULL-terminator */
    out->chars = yw_duplicate_str(s);
    out->chars_len = strlen(s);
}
void yw_json_string_deinit(YW_JSONString *str)
{
    free(str->chars);
}
void yw_json_string_clone(YW_JSONString *dest, YW_JSONString const *str)
{
    dest->chars = YW_ALLOC(char, str->chars_len);
    for (int i = 0; i < str->chars_len; i++)
    {
        dest->chars[i] = str->chars[i];
    }
}
bool yw_json_string_equals(YW_JSONString const *s, char const *str)
{
    if (s == NULL)
    {
        return false;
    }
    bool eq = true;
    for (int i = 0; i < s->chars_len; i++)
    {
        if (str[i] == '\0')
        {
            if (i != s->chars_len)
            {
                /* End of str, but not the end of s->chars */
                eq = false;
            }
            break;
        }
        else if (str[i] != s->chars[i])
        {
            /* Character mismatch */
            eq = false;
            break;
        }
        else if (i == (s->chars_len - 1) && str[i + 1] != '\0')
        {
            /* End of s->chars, but not the end of str */
            eq = false;
            break;
        }
    }
    return eq;
}

void yw_json_object_entry_init(YW_JSONObjectEntry *out, char const *name, YW_JSONValue **v)
{
    yw_json_string_init(&out->name, name);
    out->value = *v;
    *v = NULL;
}
void yw_json_object_entry_deinit(YW_JSONObjectEntry *ent)
{
    yw_json_string_deinit(&ent->name);
    yw_json_value_free(ent->value);
}
void yw_json_add_value_to_object_entry(int *ents_cap, int *ents_len, YW_JSONObjectEntry **items, char const *name, YW_JSONValue **v)
{
    YW_JSONObjectEntry ent;
    yw_json_object_entry_init(&ent, name, v);
    YW_PUSH(YW_JSONObjectEntry, ents_cap, ents_len, items, ent);
}

static void yw_value_deinit(YW_JSONValue *value)
{
    if (value == NULL)
    {
        return;
    }
    switch (value->type)
    {
    case YW_JSON_NUMBER:
    case YW_JSON_BOOLEAN:
    case YW_JSON_NULL:
        break;
    case YW_JSON_OBJECT:
        for (int i = 0; i < value->object_val.len; i++)
        {
            yw_json_object_entry_deinit(&value->object_val.entries[i]);
        }
        free(value->object_val.entries);
        break;
    case YW_JSON_ARRAY:
        for (int i = 0; i < value->array_val.len; i++)
        {
            yw_json_value_free(value->array_val.items[i]);
        }
        free(value->array_val.items);
        break;
    case YW_JSON_STRING:
        yw_json_string_deinit(&value->string_val.str);
        break;
    }
}
void yw_json_value_free(YW_JSONValue *value)
{
    yw_value_deinit(value);
    free(value);
}

static void yw_value_clone(YW_JSONValue *dest, YW_JSONValue const *value)
{
    *dest = *value;
    switch (value->type)
    {
    case YW_JSON_NUMBER:
    case YW_JSON_BOOLEAN:
    case YW_JSON_NULL:
        break;
    case YW_JSON_OBJECT:
        dest->object_val.entries = YW_ALLOC(YW_JSONObjectEntry, value->object_val.len);
        for (int i = 0; i < value->object_val.len; i++)
        {
            yw_json_string_clone(&dest->object_val.entries[i].name, &value->object_val.entries[i].name);
            dest->object_val.entries[i].value = yw_json_value_clone(value->object_val.entries[i].value);
        }
        break;
    case YW_JSON_ARRAY:
        dest->array_val.items = YW_ALLOC(YW_JSONValue *, value->array_val.len);
        for (int i = 0; i < value->array_val.len; i++)
        {
            dest->array_val.items[i] = yw_json_value_clone(value->array_val.items[i]);
        }
        break;
    case YW_JSON_STRING:
        yw_json_string_clone(&dest->string_val.str, &value->string_val.str);
        break;
    default:
        YW_ILLEGAL_VALUE(value->type);
    }
}
YW_JSONValue *yw_json_value_clone(YW_JSONValue const *value)
{
    YW_JSONValue *dest = YW_ALLOC(YW_JSONValue, 1);
    yw_value_clone(dest, value);
    return dest;
}

YW_JSONValue *yw_json_new_object(YW_JSONObjectEntry **entries, int *entries_len)
{
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    res->type = YW_JSON_OBJECT;
    if (entries != NULL)
    {
        res->object_val.entries = *entries;
        res->object_val.len = *entries_len;
    }
    else
    {
        res->object_val.entries = NULL;
        res->object_val.len = 0;
    }
    *entries = NULL;
    *entries_len = 0;
    return res;
}
YW_JSONValue *yw_json_new_array(YW_JSONValue ***items, int *items_len)
{
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    res->type = YW_JSON_ARRAY;
    if (items != NULL)
    {
        res->array_val.items = *items;
        res->array_val.len = *items_len;
    }
    else
    {
        res->array_val.items = NULL;
        res->array_val.len = 0;
    }
    *items = NULL;
    *items_len = 0;
    return res;
}
YW_JSONValue *yw_json_new_number(double n)
{
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    res->type = YW_JSON_NUMBER;
    res->number_val.num = n;
    return res;
}
YW_JSONValue *yw_json_new_string(char const *s)
{
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    res->type = YW_JSON_STRING;
    yw_json_string_init(&res->string_val.str, s);
    return res;
}
YW_JSONValue *yw_json_new_boolean(bool b)
{
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    res->type = YW_JSON_BOOLEAN;
    res->boolean_val.bol = b;
    return res;
}
YW_JSONValue *yw_json_new_null(void)
{
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    res->type = YW_JSON_NULL;
    return res;
}

YW_JSONObjectValue const *yw_json_expect_object(YW_JSONValue const *value)
{
    if (value == NULL || value->type != YW_JSON_OBJECT)
    {
        return NULL;
    }
    return &value->object_val;
}
YW_JSONValue const *yw_json_find_object_entry(YW_JSONValue const *value, char const *name)
{
    YW_JSONObjectValue const *ov = yw_json_expect_object(value);
    if (ov == NULL)
    {
        return NULL;
    }
    for (int i = 0; i < ov->len; i++)
    {
        if (yw_json_string_equals(&ov->entries[i].name, name))
        {
            return ov->entries[i].value;
        }
    }
    return NULL;
}
YW_JSONArrayValue const *yw_json_expect_array(YW_JSONValue const *value)
{
    if (value == NULL || value->type != YW_JSON_ARRAY)
    {
        return NULL;
    }
    return &value->array_val;
}
YW_JSONString const *yw_json_expect_string(YW_JSONValue const *value)
{
    if (value == NULL || value->type != YW_JSON_STRING)
    {
        return NULL;
    }
    return &value->string_val.str;
}
bool yw_json_expect_number(double *out, YW_JSONValue const *value)
{
    if (value == NULL || value->type != YW_JSON_NUMBER)
    {
        return false;
    }
    *out = value->number_val.num;
    return true;
}
bool yw_json_expect_boolean(bool *out, YW_JSONValue const *value)
{
    if (value == NULL || value->type != YW_JSON_BOOLEAN)
    {
        return false;
    }
    *out = value->boolean_val.bol;
    return true;
}
bool yw_json_expect_null(YW_JSONValue const *value)
{
    if (value == NULL || value->type != YW_JSON_NULL)
    {
        return false;
    }
    return true;
}

/******************************************************************************
 *
 * JSON parser
 *
 *****************************************************************************/

typedef struct YW_JSONParser
{
    YW_TextReader tr;
} YW_JSONParser;

static bool yw_parse_value(YW_JSONValue *out, YW_JSONParser *par);

static void yw_skip_whitespaces(YW_JSONParser *par)
{
    while (1)
    {
        YW_TextCursor old_cursor = par->tr.cursor;
        if (yw_consume_one_of_chars(&par->tr, " \t\n\r") == -1)
        {
            par->tr.cursor = old_cursor;
            break;
        }
    }
}
static bool yw_parse_number(double *out, YW_JSONParser *par)
{
    /*
     * Note that we don't parse the number directly - We only check if it's a
     * valid  number. Rest of the job is handled by the standard library.
     */
    YW_TextCursor start_cursor = par->tr.cursor;
    YW_TextCursor end_cursor, cursor_before_exp;

    /***************************************************************************
     * Sign
     **************************************************************************/

    yw_consume_char(&par->tr, '-');

    /***************************************************************************
     * Integer part
     **************************************************************************/
    /* If we have 0, we cannot have any more digits */
    if (!yw_consume_char(&par->tr, '0'))
    {
        bool got_any_char = false;
        while (!yw_text_reader_is_eof(&par->tr))
        {
            YW_Char32 temp_char = yw_peek_char(&par->tr);
            if (yw_is_ascii_digit(temp_char))
            {
                yw_consume_any_char(&par->tr);
                got_any_char = true;
            }
            else
            {
                break;
            }
        }
        if (!got_any_char)
        {
            goto fail;
        }
    }

    /***************************************************************************
     * Decimal point
     **************************************************************************/
    if (yw_consume_char(&par->tr, '.'))
    {
        /***********************************************************************
         * Fractional part
         **********************************************************************/
        while (!yw_text_reader_is_eof(&par->tr))
        {
            YW_Char32 temp_char = yw_peek_char(&par->tr);
            if (yw_is_ascii_digit(temp_char))
            {
                yw_consume_any_char(&par->tr);
            }
            else
            {
                break;
            }
        }
    }

    /***************************************************************************
     * Exponent indicator
     **************************************************************************/
    cursor_before_exp = par->tr.cursor;
    if (yw_consume_one_of_chars(&par->tr, "eE") != -1)
    {
        int digit_count = 0;

        /***********************************************************************
         * Exponent sign
         **********************************************************************/
        yw_consume_one_of_chars(&par->tr, "+-");

        /***********************************************************************
         * Exponent
         **********************************************************************/
        while (!yw_text_reader_is_eof(&par->tr))
        {
            YW_Char32 temp_char = yw_peek_char(&par->tr);
            if (yw_is_ascii_digit(temp_char))
            {
                yw_consume_any_char(&par->tr);
                digit_count++;
            }
            else
            {
                break;
            }
        }
        if (digit_count == 0)
        {
            par->tr.cursor = cursor_before_exp;
        }
    }

    end_cursor = par->tr.cursor;

    /***************************************************************************
     * Now we parse the number
     **************************************************************************/
    {
        char *temp_buf = YW_ALLOC(char, end_cursor - start_cursor + 1);
        par->tr.cursor = start_cursor;
        while (par->tr.cursor < end_cursor)
        {
            int dest = par->tr.cursor - start_cursor;
            temp_buf[dest] = yw_consume_any_char(&par->tr);
        }
        temp_buf[par->tr.cursor - start_cursor] = '\0';
        char *nptr;
        double res = strtod(temp_buf, &nptr);
        if (*nptr != '\0')
        {
            fprintf(stderr, "%s: strtod() failed to parse some(or all) of %s\n", __func__, temp_buf);
        }
        free(temp_buf);
        *out = res;
    }

    return true;
fail:
    par->tr.cursor = start_cursor;
    return false;
}
static bool yw_parse_string(YW_JSONString *out, YW_JSONParser *par)
{
    YW_TextCursor old_cursor = par->tr.cursor;
    char *chars = NULL;
    int chars_len = 0;
    int chars_cap = 0;

    if (!yw_consume_char(&par->tr, '"'))
    {
        goto fail;
    }
    while (1)
    {
        YW_Char32 chr = yw_consume_any_char(&par->tr);
        if (chr == -1)
        {
            goto fail;
        }
        else if (chr == '"')
        {
            break;
        }
        else if (chr == '\\')
        {
            YW_Char32 escape_chr = yw_consume_one_of_chars(&par->tr, "\"\\/bfnrtu");
            switch (escape_chr)
            {
            case '\"':
            case '\\':
            case '/':
                chr = escape_chr;
                break;
            case 'b':
                chr = '\b';
                break;
            case 'f':
                chr = '\f';
                break;
            case 'n':
                chr = '\n';
                break;
            case 'r':
                chr = '\r';
                break;
            case 't':
                chr = '\t';
                break;
            case 'u': {
                chr = 0;
                for (int i = 0; i < 4; i++)
                {
                    YW_Char32 digit_chr = yw_consume_any_char(&par->tr);
                    if (yw_is_ascii_digit(digit_chr))
                    {
                        chr = (chr * 16) + (digit_chr - '0');
                    }
                    else if (yw_is_ascii_uppercase_hex_digit(digit_chr))
                    {
                        chr = (chr * 16) + (digit_chr - 'A' + 10);
                    }
                    else if (yw_is_ascii_lowercase_hex_digit(digit_chr))
                    {
                        chr = (chr * 16) + (digit_chr - 'A' + 10);
                    }
                    else
                    {
                        goto fail;
                    }
                }
                break;
            }
            case -1:
                goto fail;
            }
        }
        if (chr == '\0')
        {
            YW_PUSH(char, &chars_cap, &chars_len, &chars, '\0');
        }
        else
        {
            char *temp_str = yw_char_to_str(chr);
            for (int i = 0; temp_str[i] != '\0'; i++)
            {
                YW_PUSH(char, &chars_cap, &chars_len, &chars, temp_str[i]);
            }
            free(temp_str);
        }
    }
    YW_SHRINK_TO_FIT(char, &chars_cap, chars_len, &chars);
    out->chars = chars;
    out->chars_len = chars_len;
    return true;
fail:
    par->tr.cursor = old_cursor;
    return false;
}
static bool yw_parse_object(YW_JSONObjectEntry **items_out, int *len_out, YW_JSONParser *par)
{
    YW_TextCursor old_cursor = par->tr.cursor;
    YW_JSONObjectEntry *items = NULL;
    int items_len = 0;
    int items_cap = 0;

    if (!yw_consume_char(&par->tr, '{'))
    {
        goto fail;
    }
    while (1)
    {
        YW_JSONObjectEntry ent;
        /* < >name : value  ***************************************************/
        /* < >name : value , **************************************************/
        yw_skip_whitespaces(par);
        /*  <name> : value  ***************************************************/
        /*  <name> : value , **************************************************/
        if (!yw_parse_string(&ent.name, par))
        {
            goto fail;
        }
        /*  name< >: value  ***************************************************/
        yw_skip_whitespaces(par);
        /*  name <:> value  ***************************************************/
        /*  name <:> value , **************************************************/
        if (!yw_consume_char(&par->tr, ':'))
        {
            yw_json_string_deinit(&ent.name);
            goto fail;
        }
        /*  name :< >value  ***************************************************/
        yw_skip_whitespaces(par);
        /*  name : <value>  ***************************************************/
        /*  name : <value> , **************************************************/
        YW_JSONValue v;
        if (!yw_parse_value(&v, par))
        {
            yw_json_string_deinit(&ent.name);
            goto fail;
        }
        /*  name : value< > ***************************************************/
        /*  name : value< >, **************************************************/
        yw_skip_whitespaces(par);

        /*  name : value <,> **************************************************/
        bool has_more = yw_consume_char(&par->tr, ',');

        /* Add object entry to the result */
        ent.value = YW_ALLOC(YW_JSONValue, 1);
        *ent.value = v;
        YW_PUSH(YW_JSONObjectEntry, &items_cap, &items_len, &items, ent);

        if (!has_more)
        {
            break;
        }
    }
    YW_SHRINK_TO_FIT(YW_JSONObjectEntry, &items_cap, items_len, &items);
    if (!yw_consume_char(&par->tr, '}'))
    {
        goto fail;
    }
    *items_out = items;
    *len_out = items_len;
    return true;
fail:
    for (int i = 0; i < items_len; i++)
    {
        yw_json_object_entry_deinit(&items[i]);
    }
    free(items);
    par->tr.cursor = old_cursor;
    return false;
}
static bool yw_parse_array(YW_JSONValue ***items_out, int *len_out, YW_JSONParser *par)
{
    YW_TextCursor old_cursor = par->tr.cursor;
    YW_JSONValue **items = NULL;
    int items_len = 0;
    int items_cap = 0;

    if (!yw_consume_char(&par->tr, '['))
    {
        goto fail;
    }
    while (1)
    {
        /* < >value  **********************************************************/
        /* < >value , *********************************************************/
        yw_skip_whitespaces(par);
        /*  <value>  **********************************************************/
        /*  <value> , *********************************************************/
        YW_JSONValue temp_val;
        if (!yw_parse_value(&temp_val, par))
        {
            goto fail;
        }
        /*  value< > **********************************************************/
        /*  value< >, *********************************************************/
        yw_skip_whitespaces(par);

        /*  value <,> *********************************************************/
        bool has_more = yw_consume_char(&par->tr, ',');

        /* Add object entry to the result */
        YW_JSONValue *val = YW_ALLOC(YW_JSONValue, 1);
        *val = temp_val;
        YW_PUSH(YW_JSONValue *, &items_cap, &items_len, &items, val);

        if (!has_more)
        {
            break;
        }
    }
    YW_SHRINK_TO_FIT(YW_JSONValue *, &items_cap, items_len, &items);
    if (!yw_consume_char(&par->tr, ']'))
    {
        goto fail;
    }
    *items_out = items;
    *len_out = items_len;
    return true;
fail:
    for (int i = 0; i < items_len; i++)
    {
        yw_json_value_free(items[i]);
    }
    free(items);
    par->tr.cursor = old_cursor;
    return false;
}
static bool yw_parse_value(YW_JSONValue *out, YW_JSONParser *par)
{
    if (yw_parse_object(&out->object_val.entries, &out->object_val.len, par))
    {
        out->type = YW_JSON_OBJECT;
        return true;
    }
    else if (yw_parse_array(&out->array_val.items, &out->array_val.len, par))
    {
        out->type = YW_JSON_ARRAY;
        return true;
    }
    else if (yw_parse_number(&out->number_val.num, par))
    {
        out->type = YW_JSON_NUMBER;
        return true;
    }
    else if (yw_parse_string(&out->string_val.str, par))
    {
        out->type = YW_JSON_STRING;
        return true;
    }
    else if (yw_consume_str(&par->tr, "true", YW_NO_MATCH_FLAGS))
    {
        out->type = YW_JSON_BOOLEAN;
        out->boolean_val.bol = true;
        return true;
    }
    else if (yw_consume_str(&par->tr, "false", YW_NO_MATCH_FLAGS))
    {
        out->type = YW_JSON_BOOLEAN;
        out->boolean_val.bol = false;
        return true;
    }
    else if (yw_consume_str(&par->tr, "null", YW_NO_MATCH_FLAGS))
    {
        out->type = YW_JSON_NULL;
        return true;
    }
    return false;
}

YW_JSONValue *yw_parse_json(uint8_t const *chars, int chars_len)
{
    YW_JSONParser par;
    memset(&par, 0, sizeof(par));
    yw_text_reader_init(&par.tr, chars, chars_len);
    yw_skip_whitespaces(&par);
    YW_JSONValue temp;
    if (!yw_parse_value(&temp, &par))
    {
        return NULL;
    }
    YW_JSONValue *res = YW_ALLOC(YW_JSONValue, 1);
    *res = temp;
    return res;
}
YW_JSONValue *yw_parse_json_from_c_str(char const *s)
{
    return yw_parse_json((uint8_t const *)s, strlen(s));
}
