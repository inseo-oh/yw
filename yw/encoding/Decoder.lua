--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

local object = require "yw.common.object"
local codepoints = require "yw.common.codepoints"

---@alias Encoding_HandlerResult
---| { type: "finished"}                   https://encoding.spec.whatwg.org/#finished
---| { type: "error", codepoint: number? }  https://encoding.spec.whatwg.org/#error
---| { type: "continue" }                  https://encoding.spec.whatwg.org/#continue
---| { type: "items", items: table }

---@alias Encoding_Result
---| { type: "finished" }
---| { type: "error", codepoint: number? }

---@alias Encoding_ResultOrContinue
---| { type: "finished" }
---| { type: "error", codepoint: number? }
---| { type: "continue" }

---@class Encoding_Decoder
local Decoder = {}

function Decoder:new()
    local o = object.create(self)

    return o
end

---@param queue Encoding_IOQueue
---@param item Encoding_IOQueue_Item
---@return Encoding_HandlerResult
function Decoder:handler(queue, item)
    error("Not implemented")
end

---https://encoding.spec.whatwg.org/#concept-encoding-run
---@param input Encoding_IOQueue
---@param output Encoding_IOQueue
---@param errorMode Encoding_ErrorMode
---@return Encoding_Result
function Decoder:processQueue(input, output, errorMode)
    -- 1. While true:
    while true do
        -- 1. Let result be the result of processing an item with the result of reading from input, encoderDecoder, input, output, and mode.
        local result = self:processItem(input:read(), input, output, errorMode)

        -- 2. If result is not continue, then return result.
        if result.type ~= "continue" then
            return result --[[@as Encoding_Result]]
        end
    end
end

---https://encoding.spec.whatwg.org/#concept-encoding-process
---@param item Encoding_IOQueue_Item
---@param input Encoding_IOQueue
---@param output Encoding_IOQueue
---@param errorMode Encoding_ErrorMode
---@return Encoding_ResultOrContinue
function Decoder:processItem(item, input, output, errorMode)
    -- 1. Assert: if encoderDecoder is an encoder instance, mode is not "replacement".

    -- 2. Assert: if encoderDecoder is a decoder instance, mode is not "html".
    assert(errorMode ~= "html")

    -- 3. Assert: if encoderDecoder is an encoder instance, item is not a surrogate.

    -- 4. Let result be the result of running encoderDecoder’s handler on input and item.
    local result = self:handler(input, item)

    if result.type == "finished" then
        -- [5. If result is finished:]

        -- 1. Push end-of-queue to output.
        table.insert(output.list, { type = "end" })

        -- 2. Return result.
        return result --[[@as Encoding_ResultOrContinue]]
    elseif result.type == "items" then
        -- [6. Otherwise, if result is one or more items:]

        for _, res in ipairs(result.items) do
            -- 1. Assert: if encoderDecoder is a decoder instance, result does not contain any surrogates.
            assert(not codepoints.isSurrogate(res.codepoint))
            -- 2. Push result to output.
            table.insert(output.list, res)
        end
    elseif result.type == "error" then
        -- [7. Otherwise, if result is an error, switch on mode and run the associated steps:]

        if errorMode == "replacement" then
            -- ["replacement"]

            -- 1. Push U+FFFD (�) to output.
            table.insert(output.list, { type = "codepoint", codepoint = 0xfffd })
        elseif errorMode == "html" then
            -- ["html"]

            -- Push 0x26 (&), 0x23 (#),
            output:pushValue(0x26)
            output:pushValue(0x23)
            -- followed by the shortest sequence of 0x30 (0) to 0x39 (9), inclusive, representing result’s code point’s value in base ten,
            local codepointStr = tostring(result.codepoint or 0)
            for i = 1, #codepointStr do
                output:pushValue(codepointStr:byte(i))
            end
            -- followed by 0x3B (;) to output.
            output:pushValue(0x3b)
        elseif errorMode == "fatal" then
            -- ["fatal"]
            
            -- 1. Return result.
            return result --[[@as Encoding_ResultOrContinue]]
        end
    end
    -- 8. Return continue.
    return { type = "continue" }
end

return Decoder
