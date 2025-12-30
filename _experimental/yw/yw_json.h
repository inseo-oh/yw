/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.
 */
#ifndef YW_JSON_H_
#define YW_JSON_H_

#include <stdbool.h>
#include <stdint.h>

/*
 * Note that we don't use fancy reference counting or garbage collection here for simplicity.
 * This means when we clone an array, object, or string, everything inside it gets cloned.
 *
 * We may reconsider using reference counting if this ends up using too much memory though.
 */

typedef union YW_JSONValue YW_JSONValue;

/* JSON allows embedded NULLs inside strings, so we can't use normal C strings. */
typedef struct YW_JSONString
{
    char *chars; /* This is *not* NULL-terminated! */
    int chars_len;
} YW_JSONString;

void yw_json_string_init(YW_JSONString *out, char const *s);
void yw_json_string_deinit(YW_JSONString *str);
void yw_json_string_clone(YW_JSONString *dest, YW_JSONString const *str);

/*
 * NOTE: Embedded zeros will be replaced with U+FFFD character.
 * Also, it is safe to pass NULL to str.
 *
 * Caller owns the returned array.
 */
char *yw_json_string_to_c_str(YW_JSONString const *str);

/* NOTE: It is safe to pass NULL to s. */
bool yw_json_string_equals(YW_JSONString const *s, char const *str);

typedef struct YW_JSONObjectEntry
{
    YW_JSONString name;
    YW_JSONValue *value;
} YW_JSONObjectEntry;

/* NOTE: Given YW_JSONValue will be moved into the entry, and pointer v will be cleared. */
void yw_json_object_entry_init(YW_JSONObjectEntry *out, char const *name, YW_JSONValue **v);
void yw_json_object_entry_deinit(YW_JSONObjectEntry *ent);

/* NOTE: Given YW_JSONValue will be moved into the entry, and pointer v will be cleared. */
void yw_json_add_value_to_object_entry(int *ents_cap, int *ents_len, YW_JSONObjectEntry **items, char const *name, YW_JSONValue **v);

typedef enum
{
    YW_JSON_OBJECT,
    YW_JSON_ARRAY,
    YW_JSON_NUMBER,
    YW_JSON_STRING,
    YW_JSON_BOOLEAN,
    YW_JSON_NULL,
} YW_JSONValueType;

typedef struct YW_JSONObjectValue
{
    YW_JSONValueType type; /* YW_JSON_OBJECT */
    YW_JSONObjectEntry *entries;
    int len;
} YW_JSONObjectValue;

/*
 * NOTE: Given entries will be moved into the value, and the list will be cleared.
 *
 * If entries is NULL, empty object will be created (and entries_len is ignored).
 */
YW_JSONValue *yw_json_new_object(YW_JSONObjectEntry **entries, int *entries_len);

typedef struct YW_JSONArrayValue
{
    YW_JSONValueType type; /* YW_JSON_ARRAY */
    YW_JSONValue **items;
    int len;
} YW_JSONArrayValue;

/*
 * NOTE: Given items will be moved into the value, and the list will be cleared.
 *
 * If items is NULL, empty array will be created (and items_len is ignored).
 */
YW_JSONValue *yw_json_new_array(YW_JSONValue ***items, int *items_len);

typedef struct YW_JSONNumberValue
{
    YW_JSONValueType type; /* YW_JSON_NUMBER */
    double num;
} YW_JSONNumberValue;

YW_JSONValue *yw_json_new_number(double num);

typedef struct YW_JSONStringValue
{
    YW_JSONValueType type; /* YW_JSON_STRING */
    YW_JSONString str;
} YW_JSONStringValue;

YW_JSONValue *yw_json_new_string(char const *s);

typedef struct YW_JSONBooleanValue
{
    YW_JSONValueType type; /* YW_JSON_BOOLEAN */
    bool bol;
} YW_JSONBooleanValue;

YW_JSONValue *yw_json_new_boolean(bool b);
YW_JSONValue *yw_json_new_null(void);

void yw_json_value_free(YW_JSONValue *value);
YW_JSONValue *yw_json_value_clone(YW_JSONValue const *value);

/* NOTE: It is safe to pass NULL to expect~ and find~ functions. */
YW_JSONObjectValue const *yw_json_expect_object(YW_JSONValue const *value);
YW_JSONValue const *yw_json_find_object_entry(YW_JSONValue const *value, char const *name);
YW_JSONArrayValue const *yw_json_expect_array(YW_JSONValue const *value);
YW_JSONString const *yw_json_expect_string(YW_JSONValue const *value);
/* These take pointer to where to store the result instead of returning pointer. */
bool yw_json_expect_number(double *out, YW_JSONValue const *value);
bool yw_json_expect_boolean(bool *out, YW_JSONValue const *value);
bool yw_json_expect_null(YW_JSONValue const *value);

/* Caller owns the returned value. */
YW_JSONValue *yw_json_parse(uint8_t const *chars, int chars_len);

/* Caller owns the returned value. */
YW_JSONValue *yw_json_parse_from_c_str(char const *s);

#endif /* #ifndef YW_JSON_H_ */
