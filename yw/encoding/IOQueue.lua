--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

local object = require "yw.common.object"

---@alias Encoding_IOQueue_Item {type: "value", value: any}|{tyoe: "end"}
---@alias Encoding_IOQueue_InitMode
--- | "immediate"   Put end-of-queue at the end
--- | "streaming"   Do not put end-of-queue at the end

---https://encoding.spec.whatwg.org/#concept-stream
---@class Encoding_IOQueue
---@field list Encoding_IOQueue_Item[]
local IOQueue = {}


---@param initmode Encoding_IOQueue_InitMode
---@return Encoding_IOQueue
function IOQueue:new(initmode)
    local o = object.create(self)

    o.list = {}
    if initmode == "immediate" then
        table.insert(o.list, { type = "end" })
    end

    return o
end

---https://encoding.spec.whatwg.org/#concept-stream-read
---@param number number?
function IOQueue:read(number)
    if number ~= nil and 1 < number then
        -- 1. Let readItems be « ».
        local readItems = {}

        --  2. Perform the following step number times:
        for _ = 1, number do
            -- 1. Append to readItems the result of reading an item from ioQueue.
            -- NOTE: We also filter out end-of-queue instead of doing it in step 3.
            local item = self:read()
            if item.type ~= "end" then
                table.insert(readItems, item)
            end
        end

        -- 3. Remove end-of-queue from readItems.
        -- NOTE: We did this during above step 2.

        -- 4. Return readItems.
        return readItems
    end

    -- 1. If ioQueue is empty, then wait until its size is at least 1.
    while #self.list == 0 do
        coroutine.yield()
    end
    -- 2. If ioQueue[0] is end-of-queue, then return end-of-queue.
    if self.list[1].type == "end" then
        return { type = "end" }
    end
    -- 3. Remove ioQueue[0] and return it.
    local item = self.list[1]
    table.remove(self.list, 1)
    return item
end

---https://encoding.spec.whatwg.org/#i-o-queue-peek
---@return Encoding_IOQueue_Item
---@param number number
function IOQueue:peek(number)
    -- 1. Wait until either ioQueue’s size is equal to or greater than number, or ioQueue contains end-of-queue,  whichever comes first.
    while true do
        if #self.list >= number then
            break
        end
        local containsEndOfQueue = false
        for n = 1, #self.list do
            if self.list[n].type == "end" then
                containsEndOfQueue = true
                break
            end
        end
        if not containsEndOfQueue then
            break
        end
        coroutine.yield()
    end
    -- 2. Let prefix be « ».
    local prefix = {}

    -- 3. For each n in the range 1 to number, inclusive:
    for n = 1, number do
        -- 1. If ioQueue[n] is end-of-queue, break.
        if self.list[n].type == "end" then
            break
        end
        -- 2. Otherwise, append ioQueue[n] to prefix.
        table.insert(prefix, self.list[n])
    end

    -- 4. Return prefix.
    return prefix
end

---https://encoding.spec.whatwg.org/#concept-stream-push
---@param item Encoding_IOQueue_Item
function IOQueue:push(item)
    -- 1. If the last item in ioQueue is end-of-queue, then:
    if #self.list > 0 and self.list[#self.list].type == "end" then
        -- 1. If item is end-of-queue, do nothing.
        if item.type == "end" then
            return
        end
        -- 2. Otherwise, insert item before the last item in ioQueue.
        table.insert(self.list, #self.list, item)
    end
    -- 2. Otherwise, append item to ioQueue.
    table.insert(self.list, item)
end

---https://encoding.spec.whatwg.org/#concept-stream-push
---@param value any
function IOQueue:pushValue(value)
    self:push({ type = "value", value = value })
end

---https://encoding.spec.whatwg.org/#concept-stream-prepend
---@param item Encoding_IOQueue_Item
function IOQueue:restore(item)
    -- To restore an item other than end-of-queue to an I/O queue, perform the list prepend operation.
    table.insert(self.list, 1, item)
end

---https://encoding.spec.whatwg.org/#concept-stream-prepend
---@param value any
function IOQueue:restoreValue(value)
    self:restore({ type = "value", value = value })
end
