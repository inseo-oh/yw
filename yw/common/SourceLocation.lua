--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object = require "yw.common.object"

---@class SourceLocation
---@field file   string
---@field line   number
---@field column number
local SourceLocation = {}

---@param file string
---@param line number
---@param column number
function SourceLocation:new(file, line, column)
    local o = object.create(self)

    o.file = file
    o.line = line
    o.column = column
    o.__tostring = function()
        return string.format("%s:%d:%d", o.file, o.line, o.column)
    end

    return o
end

return SourceLocation