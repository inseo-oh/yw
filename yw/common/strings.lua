--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local codepoints = require "yw.common.codepoints"

local strings = {}

---https://infra.spec.whatwg.org/#strings


function strings.fromCharCodes(charCodes)
    local idx = 1
    local cps = {}
    -- https://tc39.es/ecma262/#sec-ecmascript-language-types-string-type
    while idx <= #charCodes do
        if not codepoints.isSurrogate(charCodes[idx]) then
            -- A code unit that is not a leading surrogate and not a trailing surrogate is interpreted as a code point with the same value.
            table.insert(cps, charCodes[idx])
            idx = idx + 1
        elseif codepoints.isLeadingSurrogate(charCodes[idx]) and codepoints.isTrailingSurrogate(charCodes[idx + 1] or 0) then
            -- A sequence of two code units, where the first code unit c1 is a leading surrogate and the second code unit c2 a trailing surrogate, is a surrogate pair and is interpreted as a code point with the value (c1 - 0xD800) × 0x400 + (c2 - 0xDC00) + 0x10000. (See 11.1.3) 
            local c1, c2 = charCodes[idx], charCodes[idx + 1]
            table.insert(cps, (c1 - 0xd800) * 0x400 + (c2 - 0xdc00) + 0x10000)
            idx = idx + 2
        elseif codepoints.isSurrogate(charCodes[idx]) then
            -- A code unit that is a leading surrogate or trailing surrogate, but is not part of a surrogate pair, is interpreted as a code point with the same value.
            table.insert(cps, charCodes[idx])
            idx = idx + 1
        else
            error("unreachable")
        end
    end
    return utf8.char(table.unpack(cps))
end

---https://infra.spec.whatwg.org/#code-point-substring
---@param s string
---@param start integer
---@param length integer
---@return string
---@nodiscard
function strings.codepointSubstring(s, start, length)
    -- Make start zero-indexed
    start = start - 1

    -- 1. Assert: start and length are nonnegative.
    assert(0 <= start)

    -- 2. Assert: start + length is less than or equal to string’s code point length.
    assert((start + length) <= utf8.len(s))

    -- 3. Let result be the empty string.

    -- 4. For each i in the range from start to start + length, exclusive: append the ith code point of string to result.
    local resultStartIdx = 0
    local resultEndIdx = 0

    local i = 1
    for p, c in utf8.codes(s) do
        if start == i then
            resultStartIdx = p
        elseif start + length == i then
            resultEndIdx = p
        end
        i = i + 1
    end
    local result = string.sub(resultStartIdx, resultEndIdx)

    -- 5. Return result.
    return result
end

---@param str string
---@param another string
---@return boolean
function strings.startsWith(str, another)
    return string.sub(str, 1, #another) == another
end

---@param str string
---@param another string
---@return boolean
function strings.endsWith(str, another)
    return string.sub(str, -#another) == another
end

---https://infra.spec.whatwg.org/#ascii-case-insensitive
---@param s string
---@param another string
---@return boolean
function strings.startsWithASCIICaseInsensitive(s, another)
    return strings.startsWith(strings.asciiLowercase(s), strings.asciiLowercase(another))
end

---https://infra.spec.whatwg.org/#ascii-case-insensitive
---@param s string
---@param another string
function strings.equalsASCIICaseInsensitive(s, another)
    return strings.asciiLowercase(s) == strings.asciiLowercase(another)
end

---https://infra.spec.whatwg.org/#ascii-lowercase
---@param s string
---@return string
---@nodiscard
function strings.asciiLowercase(s)
    local cps = {}
    for _, c in utf8.codes(s) do
        if codepoints.isAsciiUpperAlpha(c) then
            table.insert(cps, c + 0x20)
        end
    end
    return utf8.char(table.unpack(cps))

end

---https://infra.spec.whatwg.org/#ascii-uppercase
---@param s string
---@return string
---@nodiscard
function strings.asciiUppercase(s)
    local cps = {}
    for _, c in utf8.codes(s) do
        if codepoints.isAsciiLowerAlpha(c) then
            table.insert(cps, c - 0x20)
        end
    end
    return utf8.char(table.unpack(cps))
end

---https://infra.spec.whatwg.org/#normalize-newlines
---@param s string
---@return string
---@nodiscard
function strings.normalizeNewlines(s)
    -- To normalize newlines in a string,
    -- replace every U+000D CR U+000A LF code point pair with a single U+000A LF code point,
    local result = string.gsub(s, "\r\n", "\n")
    -- and then replace every remaining U+000D CR code point with a U+000A LF code point.
    result = string.gsub(result, "\r", "\n")
    return result
end

return strings
