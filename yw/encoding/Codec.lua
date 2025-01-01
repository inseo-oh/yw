--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

local object = require "yw.common.object"

---https://encoding.spec.whatwg.org/#error-mode
---@alias Encoding_ErrorMode "replacement"|"fatal"|"html"

---@class Encoding_Codec
local Codec = {}

function Codec:new()
    local o = object.create(self)

    return o
end

---@return Encoding_Decoder
function Codec:createDecoder()
    error("Not implemented")
end

return Codec