--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local iterutil = {}

---@param iter fun():any
---@param count integer?
function iterutil.skip(iter, count)
    if count == nil then count = 1 end
    for _ = 1, count do
        iter()
    end
end

---@param iter fun():any
---@param filter fun(arg:any):boolean
---@return fun():any
function iterutil.filter(iter, filter)
    return function()
        for x in iter do
            if filter(x) then
                return x
            end
        end
    end
end

---@param iter fun():any|nil
---@param filter any|fun(arg:any):boolean If it's function, it is used as filter function. Otherwise, it simply tests whether this value is produced by `iter`.
---@return boolean
function iterutil.contains(iter, filter)
    if type(filter) ~= "function" then
        return iterutil.contains(iter, function(x) return x == filter end)
    end
    for x in iter do
        if filter(x) then
            return true
        end
    end
    return false
end

---@param iter fun():any
---@param filter nil|fun(arg:any):boolean If nil, no filter is applied.
---@return integer
function iterutil.count(iter, filter)
    if filter == nil then
        return iterutil.count(iter, function() return true end)
    end
    local currentCount = 0
    for x in iter do
        if filter(x) then
            currentCount = currentCount + 1
        end
    end
    return currentCount
end

---@param iter fun():any|nil
---@return any[]
function iterutil.collect(iter)
    local result = {}
    for x in iter do
        table.insert(result, x)
    end
    return result
end

return iterutil