--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local codepoints = {
    TAB   = 0x0009,
    LF    = 0x000a,
    FF    = 0x000c,
    CR    = 0x000d,
    SPACE = 0x0020,
}

function codepoints.toCodeUnits(c)
    if c < 0x10000 then
        return c
    end
    c = c - 0x10000
    local high = 0xd800 + (c >> 10)
    local low = 0xdc00 + (c & 0x3ff)
    return high, low
end

---https://infra.spec.whatwg.org/#leading-surrogate
---@param c string|integer
---@return boolean
function codepoints.isLeadingSurrogate(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- A leading surrogate is a code point that is in the range U+D800 to U+DBFF, inclusive.
    return (0xd800 <= c) and (c <= 0xdbff)
end

---https://infra.spec.whatwg.org/#trailing-surrogate
---@param c integer|string
---@return boolean
function codepoints.isTrailingSurrogate(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- A trailing surrogate is a code point that is in the range U+DC00 to U+DFFF, inclusive.
    return (0xdc00 <= c) and (c <= 0xdfff)
end

---https://infra.spec.whatwg.org/#surrogate
---@param c integer|string
---@return boolean
function codepoints.isSurrogate(c)
    return codepoints.isLeadingSurrogate(c) or codepoints.isTrailingSurrogate(c)
end

---https://infra.spec.whatwg.org/#noncharacter
---@param c integer|string
---@return boolean
function codepoints.scalarValue(c)
    return not codepoints.isSurrogate(c)
end

---https://infra.spec.whatwg.org/#noncharacter
---@param c integer|string
---@return boolean
function codepoints.isNonCharacter(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- A noncharacter is a code point that is in the range U+FDD0 to U+FDEF, inclusive,
    if (0xfdd0 <= c) and (c <= 0xfdef) then
        return true
    end
    -- or U+FFFE, U+FFFF, U+1FFFE, U+1FFFF, U+2FFFE, U+2FFFF, U+3FFFE, U+3FFFF, U+4FFFE, U+4FFFF, U+5FFFE, U+5FFFF, U+6FFFE, U+6FFFF, U+7FFFE, U+7FFFF, U+8FFFE, U+8FFFF, U+9FFFE, U+9FFFF, U+AFFFE, U+AFFFF, U+BFFFE, U+BFFFF, U+CFFFE, U+CFFFF, U+DFFFE, U+DFFFF, U+EFFFE, U+EFFFF, U+FFFFE, U+FFFFF, U+10FFFE, or U+10FFFF.
    local nonCharacters = {
        0xfffe, 0xffff, 0x1fffe, 0x1ffff, 0x2fffe, 0x2ffff, 0x3fffe, 0x3ffff, 0x4fffe, 0x4ffff, 0x5fffe, 0x5ffff,
        0x6fffe, 0x6ffff, 0x7fffe, 0x7ffff, 0x8fffe, 0x8ffff, 0x9fffe, 0x9ffff, 0xafffe, 0xaffff, 0xbfffe, 0xbffff,
        0xcfffe, 0xcffff, 0xdfffe, 0xdffff, 0xefffe, 0xeffff, 0xffffe, 0xfffff, 0x10fffe, 0x10ffff
    }
    for _, cp in ipairs(nonCharacters) do
        if cp == c then
            return true
        end
    end
    return false
end

---https://infra.spec.whatwg.org/#ascii-whitespace
---@param c integer|string
---@return boolean
function codepoints.isAsciiWhitespace(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- ASCII whitespace is U+0009 TAB, U+000A LF, U+000C FF, U+000D CR, or U+0020 SPACE.
    local whitespaces = { codepoints.TAB, codepoints.LF, codepoints.FF, codepoints.CR, codepoints.SPACE }
    for _, cp in ipairs(whitespaces) do
        if cp == c then
            return true
        end
    end
    return false
end

---https://infra.spec.whatwg.org/#ascii-upper-alpha
---@param c integer|string
---@return boolean
function codepoints.isAsciiUpperAlpha(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- An ASCII upper alpha is a code point in the range U+0041 (A) to U+005A (Z), inclusive.
    return (0x0041 <= c) and (c <= 0x005a)
end

---https://infra.spec.whatwg.org/#ascii-lower-alpha
---@param c integer|string
---@return boolean
function codepoints.isAsciiLowerAlpha(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- An ASCII lower alpha is a code point in the range U+0061 (a) to U+007A (z), inclusive.
    return (0x0061 <= c) and (c <= 0x007a)
end

---https://infra.spec.whatwg.org/#ascii-alpha
---@param c integer|string
---@return boolean
function codepoints.isAsciiAlpha(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- An ASCII alpha is an ASCII upper alpha or ASCII lower alpha.
    return codepoints.isAsciiUpperAlpha(c) or codepoints.isAsciiLowerAlpha(c)
end

---https://infra.spec.whatwg.org/#ascii-digit
---@param c integer|string
---@return boolean
function codepoints.isAsciiDigit(c)
    if type(c) == "string" then c = utf8.codepoint(c) end
    -- An ASCII digit is a code point in the range U+0030 (0) to U+0039 (9), inclusive.
    return (0x0030 <= c) and (c <= 0x0039)
end

---https://infra.spec.whatwg.org/#ascii-alphanumeric
---@param c integer|string
---@return boolean
function codepoints.isAsciiAlphanumeric(c)
    -- An ASCII alphanumeric is an ASCII digit or ASCII alpha.
    return codepoints.isAsciiDigit(c) or codepoints.isAsciiAlpha(c)
end

---https://infra.spec.whatwg.org/#ascii-upper-hex-digit
---@param c integer|string
---@return boolean
function codepoints.isAsciiUpperHexDigit(c)
    -- An ASCII upper hex digit is an ASCII digit or a code point in the range U+0041 (A) to U+0046 (F), inclusive.
    return codepoints.isAsciiDigit(c) or ((0x0041 <= c) and (c <= 0x0046))
end

---https://infra.spec.whatwg.org/#ascii-lower-hex-digit
---@param c integer|string
---@return boolean
function codepoints.isAsciiLowerHexDigit(c)
    -- An ASCII lower hex digit is an ASCII digit or a code point in the range U+0061 (a) to U+0066 (f), inclusive.
    return codepoints.isAsciiDigit(c) or ((0x0061 <= c) and (c <= 0x0066))
end

---https://infra.spec.whatwg.org/#ascii-hex-digit
---@param c integer|string
---@return boolean
function codepoints.isAsciiHexDigit(c)
    -- An ASCII hex digit is an ASCII upper hex digit or ASCII lower hex digit.
    return codepoints.isAsciiUpperHexDigit(c) or codepoints.isAsciiLowerHexDigit(c)
end

return codepoints
