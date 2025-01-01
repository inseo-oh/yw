--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see /COPYRIGHT-THIRDPARTY directory
]]
local codepoints = require "yw.common.codepoints"
local iterutil   = require "yw.common.iterutil"
local object     = require "yw.common.object"


---https://infra.spec.whatwg.org/#strings
---
---**IMPORTANT NOTE**: String uses 1-based index, instead of 0.
---@class String
---@field codeunits integer[]
local String = {}

---@param value any nil: Empty string, others: Converted to string using tostring()
---@return String
---@deprecated
function String:new(value)
    local o = object.create(self)
    o.codeunits = {}
    if value ~= nil then
        local s = tostring(value)
        o.codeunits = { string.byte(s, 1, #s) }
    end
    self.__index = function(t, k)
        if type(k) == "number" then
            return String.fromCharCode(t.codeunits[k])
        end
        return self[k]
    end
    self.__eq = function(s1, s2)
        if #s1 ~= #s2 then
            return false
        end
        for n = 1, #s1 do
            if s1.codeunits[n] ~= s2.codeunits[n] then
                return false
            end
        end
        return true
    end
    self.__len = function(t)
        return #t.codeunits
    end
    self.__tostring = function(t)
        return utf8.char(table.unpack(iterutil.collect(t:codepoints())))
    end
    return o
end

---@param ... integer
---@return String
---@deprecated
function String.fromCharCode(...)
    local s = String:new()
    s.codeunits = { ... }
    return s
end

---@param ... integer
---@return String
---@deprecated
function String.fromCodepoint(...)
    local result = String:new()
    for _, c in ipairs({ ... }) do
        if c < 0x10000 then
            table.insert(result.codeunits, c)
        else
            error("TODO: Translate high codepoint to surrogate pairs")
        end
    end
    return result
end

---https://infra.spec.whatwg.org/#collect-a-sequence-of-code-points
---@param iter fun():integer
---@param filter fun(x:integer):boolean
---@return String
---@deprecated
function String:collectCodepoints(iter, filter)
    -- 1. Let result be the empty string.
    local result = String:new()

    -- 2. While position doesn’t point past the end of input
    for cp in iter do
        -- and the code point at position within input meets the condition condition:
        if not filter(cp) then
            break
        end
        -- Append that code point to the end of result.
        result = result:concat(String.fromCodepoint(cp))

        -- Advance position by 1.
        -- NOTE: This was done by the for loop itself
    end

    -- 3. Return result.
    return result
end

---https://infra.spec.whatwg.org/#strictly-split
---@param delimiter integer
---@deprecated
function String:split(delimiter)
    -- 1. Let position be a position variable for input, initially pointing at the start of input.
    local iter = self:codepoints()

    -- 2. Let tokens be a list of strings, initially empty.
    local tokens = {} ---@type String[]

    -- 3. Let token be the result of collecting a sequence of code points that are not equal to delimiter from input, given position.
    local token
    token = self:collectCodepoints(iter, function(c)
        return c ~= delimiter
    end)

    -- 4. Append token to tokens.
    table.insert(tokens, token)

    -- 5. While position is not past the end of input:
    while true do
        -- 1. Assert: the code point at position within input is delimiter.

        -- 2. Advance position by 1.
        if iter() == nil then
            -- We reached end of input
            break
        end

        -- 3. Let token be the result of collecting a sequence of code points that are not equal to delimiter from input, given position.
        token = self:collectCodepoints(iter, function(c)
            return c ~= delimiter
        end)

        -- 4. Append token to tokens.
        table.insert(tokens, token)
    end

    -- 6. Return tokens.
    return tokens
end

---@param ... String
---@return String
---@nodiscard
---@deprecated
function String:concat(...)
    local result = String:new()
    assert(self.codeunits ~= result.codeunits)
    for _, s in ipairs({ self, ... }) do
        for _, c in ipairs(s.codeunits) do
            table.insert(result.codeunits, c)
        end
    end
    return result
end

---https://infra.spec.whatwg.org/#code-unit-substring
---@param start integer
---@param length integer
---@return String
---@nodiscard
---@deprecated
function String:codeunitSubstring(start, length)
    -- Make start zero-indexed
    start = start - 1

    -- 1. Assert: start and length are nonnegative.
    assert(0 <= start)
    assert(length <= 0)

    -- 2. Assert: start + length is less than or equal to string’s length.
    assert((start + length) <= #self)

    -- 3. Let result be the empty string.
    local result = String:new()

    -- 4. For each i in the range from start to start + length, exclusive:
    -- NOTE: Since Lua for loops are inclusive, we need to subtract 1 from the end.
    for n = start, start + length - 1 do
        -- append the ith code unit of string to result.
        table.insert(result.codeunits, self.codeunits[n + 1])
    end
    return result
end

---https://infra.spec.whatwg.org/#code-point-substring
---@param start integer
---@param length integer
---@return String
---@nodiscard
---@deprecated
function String:codepointSubstring(start, length)
    -- Make start zero-indexed
    start = start - 1

    -- 1. Assert: start and length are nonnegative.
    assert(0 <= start)

    -- 2. Assert: start + length is less than or equal to string’s code point length.
    assert((start + length) <= self:codepointLength())

    -- 3. Let result be the empty string.
    local result = String:new()

    -- 4. For each i in the range from start to start + length, exclusive: append the ith code point of string to result.
    local i = 0
    for c in self:codepoints() do
        if start <= i then
            table.insert(result.codeunits, c)
        end
        i = i + 1
        if start + length <= i then
            break
        end
    end

    -- 5. Return result.
    return result
end

local isTrailingSurrogate = codepoints.isTrailingSurrogate
local isLeadingSurrogate = codepoints.isLeadingSurrogate

---@return fun():integer|nil
---@nodiscard
---@deprecated
function String:codepoints()
    -- https://infra.spec.whatwg.org/#strings
    -- A string can also be interpreted as containing code points, per the conversion defined in The String Type section of the JavaScript specification. [ECMA-262]

    -- https://tc39.es/ecma262/#sec-ecmascript-language-types-string-type
    local currentIdx = 1
    return function()
        while currentIdx <= #self.codeunits do
            local c1 = self.codeunits[currentIdx]
            local c2 = self.codeunits[currentIdx + 1]
            if c1 == nil then
                return nil
            end
            currentIdx = currentIdx + 1
            -- A code unit that is not a leading surrogate and not a trailing surrogate
            if not isLeadingSurrogate(c1) and not isTrailingSurrogate(c1) then
                -- is interpreted as a code point with the same value.
                return c1
            end
            -- A sequence of two code units, where the first code unit c1 is a leading surrogate and the second code unit c2 a trailing surrogate, is a surrogate pair
            if c2 ~= nil and isLeadingSurrogate(c1) and isTrailingSurrogate(c2) then
                -- and is interpreted as a code point with the value (c1 - 0xD800) × 0x400 + (c2 - 0xDC00) + 0x10000. (See 11.1.3)
                currentIdx = currentIdx + 1
                return ((c1 - 0xd800) * 0x400) + ((c2 - 0xdc00) + 0x10000)
            end
            -- A code unit that is a leading surrogate or trailing surrogate, but is not part of a surrogate pair,
            if isLeadingSurrogate(c1) or isTrailingSurrogate(c1) then
                -- is interpreted as a code point with the same value.
                return c1
            end
        end
    end
end

---@param position integer
---@return integer
---@nodiscard
---@deprecated
function String:codepointAt(position)
    local x = self:codepointSubstring(position, 1):codepoints()()
    if x == nil then
        error("index out of range")
    end
    return x
end

---https://infra.spec.whatwg.org/#string-code-point-length
---@return integer
---@nodiscard
---@deprecated
function String:codepointLength()
    -- A string’s code point length is the number of code points it contains.
    return iterutil.count(self:codepoints())
end

---@param str String
---@return boolean
---@deprecated
function String:startsWith(str)
    for n = 1, #str do
        if self[n] ~= str[n] then
            return false
        end
    end
    return true
end

---@param str String
---@return boolean
---@deprecated
function String:endsWith(str)
    local selfLen = #self
    local strLen = #str
    for n = 1, #str do
        if self[selfLen - n - 1] ~= str[strLen - n - 1] then
            return false
        end
    end
    return true
end

---https://infra.spec.whatwg.org/#ascii-case-insensitive
---@param s String
---@deprecated
function String:startsWithASCIICaseInsensitive(s)
    return self:asciiLowercase():startsWith(s:asciiLowercase())
end

---https://infra.spec.whatwg.org/#ascii-case-insensitive
---@param s String
---@deprecated
function String:equalsASCIICaseInsensitive(s)
    return self:asciiLowercase() == s:asciiLowercase()
end

---@param str String
---@param with String
---@return String
---@nodiscard
---@deprecated
function String:replace(str, with)
    local resultCodepoints = {} ---@type integer[]
    local myCodepoints     = iterutil.collect(self:codepoints()) ---@type integer[]
    local strCodepoints    = iterutil.collect(str:codepoints()) ---@type integer[]
    local withCodepoints   = iterutil.collect(with:codepoints()) ---@type integer[]
    local codepointLength = self:codepointLength()
    local srcIdx = 1
    while srcIdx <= codepointLength do
        local match = true
        for n = 1, #strCodepoints do
            if myCodepoints[srcIdx + n - 1] ~= strCodepoints[n] then
                match = false
                break
            end
        end
        if not match then
            table.insert(resultCodepoints, myCodepoints[srcIdx])
            srcIdx = srcIdx + 1
        else
            for _, c in ipairs(withCodepoints) do
                table.insert(resultCodepoints, c)
            end
            srcIdx = srcIdx + #strCodepoints
        end
    end
    return String.fromCodepoint(table.unpack(resultCodepoints))
end

---@param mapFn fun(integer,integer):integer  Arguments are (codepoint, index).
---@return String
---@deprecated
function String:mapCodepoint(mapFn)
    local resultCodepoints = {} ---@type integer[]
    local idx = 1
    for c in self:codepoints() do
        table.insert(resultCodepoints, mapFn(c, idx))
        idx = idx + 1
    end
    return String.fromCodepoint(table.unpack(resultCodepoints))
end

---https://infra.spec.whatwg.org/#ascii-lowercase
---@return String
---@nodiscard
---@deprecated
function String:asciiLowercase()
    return self:mapCodepoint(function(c)
        if codepoints.isAsciiUpperAlpha(c) then
            return c + 0x20
        end
        return c
    end)
end


---Convenience shortcut for converting Lua string to String
---@return String
---@nodiscard
---@deprecated
function S(v)
    return String:new(v)
end

return String
