--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

local object      = require "yw.common.object"
local Codec       = require "yw.encoding.Codec"
local Decoder     = require "yw.encoding.Decoder"

---@class Encoding_UTF8Decoder : Encoding_Decoder
---@field codepoint      number  https://encoding.spec.whatwg.org/#utf-8-code-point
---@field bytesSeen      number  https://encoding.spec.whatwg.org/#utf-8-bytes-seen
---@field bytesNeeded    number  https://encoding.spec.whatwg.org/#utf-8-bytes-needed
---@field lowerBoundary  number  https://encoding.spec.whatwg.org/#utf-8-lower-boundary
---@field upperBoundary  number  https://encoding.spec.whatwg.org/#utf-8-upper-boundary
local UTF8Decoder = object.create(Decoder)

function UTF8Decoder:new()
    local o         = Decoder.new(self)

    o.codepoint     = 0
    o.bytesSeen     = 0
    o.bytesNeeded   = 0
    o.lowerBoundary = 0x80
    o.upperBoundary = 0xbf

    return o
end

function UTF8Decoder:handler(queue, item)
    -- 1. If byte is end-of-queue and UTF-8 bytes needed is not 0, set UTF-8 bytes needed to 0 and return error.
    if item.type == "end" and self.bytesNeeded ~= 0 then
        self.bytesNeeded = 0
        return { type = "error" }
    end

    -- 2. If byte is end-of-queue, return finished.
    if item.type == "end" then
        return { type = "finished" }
    end
    local b = item.value

    -- 3. If UTF-8 bytes needed is 0, based on byte:
    if self.bytesNeeded == 0 then
        if b < 0x80 then
            -- [0x00 to 0x7F]

            -- Return a code point whose value is byte.
            return { type = "items", items = { codepoint = b } }
        elseif (0xc2 <= b) and (b <= 0xdf) then
            -- [0xC2 to 0xDF]

            -- Set UTF-8 bytes needed to 1.
            self.bytesNeeded = 1
            -- 2. Set UTF-8 code point to byte & 0x1F.
            self.codepoint = b & 0x1f
        elseif (0xe0 <= b) and (b <= 0xef) then
            -- [0xE0 to 0xEF]

            -- 1. If byte is 0xE0, set UTF-8 lower boundary to 0xA0.
            if b == 0xe0 then
                self.lowerBoundary = 0xa0
            end

            -- 2. If byte is 0xED, set UTF-8 upper boundary to 0x9F.
            if b == 0xed then
                self.upperBoundary = 0x9f
            end

            -- 3. Set UTF-8 bytes needed to 2.
            self.bytesNeeded = 2

            -- 4. Set UTF-8 code point to byte & 0xF.
            self.codepoint = b & 0xf
        elseif (0xf0 <= b) and (b <= 0xf4) then
            -- [0xF0 to 0xF4]

            -- 1. If byte is 0xF0, set UTF-8 lower boundary to 0x90.
            if b == 0xf0 then
                self.lowerBoundary = 0x90
            end

            -- 2. If byte is 0xF4, set UTF-8 upper boundary to 0x8F.
            if b == 0xf4 then
                self.upperBoundary = 0x8f
            end

            -- 3. Set UTF-8 bytes needed to 3.
            self.bytesNeeded = 3

            -- 4. Set UTF-8 code point to byte & 0x7.
            self.codepoint = b & 0x7
        else
            -- [Otherwise]

            -- Return error.
            return { type = "error" }
        end

        -- Return continue.
        return { type = "continue" }
    end

    -- 4. If byte is not in the range UTF-8 lower boundary to UTF-8 upper boundary, inclusive, then:
    if not ((self.lowerBoundary <= b) and (b <= self.upperBoundary)) then
        -- 1. Set UTF-8 code point, UTF-8 bytes needed, and UTF-8 bytes seen to 0, set UTF-8 lower boundary to 0x80, and set UTF-8 upper boundary to 0xBF.
        self.codepoint = 0
        self.bytesNeeded = 0
        self.bytesSeen = 0
        self.lowerBoundary = 0x80
        self.upperBoundary = 0xbf

        -- 2. Prepend byte to stream.
        queue:restoreValue(b)

        -- 3. Return error.
        return { type = "error" }
    end

    -- 5. Set UTF-8 lower boundary to 0x80 and UTF-8 upper boundary to 0xBF.
    self.lowerBoundary = 0x80
    self.upperBoundary = 0xbf

    -- 6. Set UTF-8 code point to (UTF-8 code point << 6) | (byte & 0x3F)
    self.codepoint = (self.codepoint << 6) | (b & 0x3f)

    -- 7. Increase UTF-8 bytes seen by one.
    self.bytesSeen = self.bytesSeen + 1

    -- 8. If UTF-8 bytes seen is not equal to UTF-8 bytes needed, return continue.
    if self.bytesSeen ~= self.bytesNeeded then
        return { type = "continue" }
    end

    -- 9. Let code point be UTF-8 code point.
    local codepoint = self.codepoint

    -- 10. Set UTF-8 code point, UTF-8 bytes needed, and UTF-8 bytes seen to 0.
    self.codepoint = 0
    self.bytesNeeded = 0
    self.bytesSeen = 0

   -- 11. Return a code point whose value is code point.
    return { type = "items", items = { codepoint = codepoint } }
end

---@class Encoding_UTF8Codec : Encoding_Codec
local UTF8Codec = object.create(Codec)

function UTF8Codec:new()
    local o = Codec.new(self)

    return o
end

function UTF8Codec:createDecoder()
    -- https://encoding.spec.whatwg.org/#utf-8-decoder
    return UTF8Decoder:new()
end
