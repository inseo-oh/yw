--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

local Logger    = require "yw.common.Logger"
local Tokenizer = require "yw.html.parser.Tokenizer"
local Parser    = require "yw.html.parser.Parser"

local profile   = require "thirdparty.profile.profile"

local L         = Logger:new("main")
local args      = { ... }

-- profile.start()
-- profile.stop()

local function dumpValue(value)
    if type(value) == "table" and value[1] == nil then
        local s = "{"
        for objK, objV in pairs(value) do
            s = s .. string.format("\"%s\":", objK) .. dumpValue(objV) .. "\n,"
        end
        s = s .. "}"
        return s
    elseif type(value) == "string" then
        return "\"" .. tostring(value) .. "\""
    elseif type(value) == "table" then
        local s = "["
        for _, arrayEntry in ipairs(value) do
            s = s .. dumpValue(arrayEntry) .. ",\n"
        end
        s = s .. "]"
        return s
    else
        return tostring(value)
    end
end

-- print(dumpValue(j))
