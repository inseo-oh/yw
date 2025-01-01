--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local token      = require "yw.html.parser.token"
local codepoints = require "yw.common.codepoints"
local Logger     = require "yw.common.Logger"
local strings    = require "yw.common.strings"
local json       = require "yw.json"
local ioutil     = require "yw.common.ioutil"


local L = Logger:new("yw.html.parser.Tokenizer")

---@return table<string, { codepoints: integer[], characters: string }>
local function loadNamedCharReferences()
    L:i("Loading named character reference table")
    local j = json.parse(ioutil.readFile(L, "yw/res/htmlNamedCharRefs.json"))
    L:i("Loaded named character reference table")
    return j
end

local NAMED_CHARACTER_REFERENCES = loadNamedCharReferences()

---@alias HTML_Parser_TokenizerState fun(t: HTML_Parser_Tokenizer)

---https://html.spec.whatwg.org/multipage/parsing.html#parse-errors
---@alias HTML_Parser_TokenizationError
---| "abrupt-closing-of-empty-comment"
---| "abrupt-doctype-public-identifier"
---| "abrupt-doctype-system-identifier"
---| "absence-of-digits-in-numeric-character-reference"
---| "cdata-in-html-content"
---| "character-reference-outside-unicode-range"
---| "control-character-in-input-stream"
---| "control-character-reference"
---| "duplicate-attribute"
---| "end-tag-with-attributes"
---| "end-tag-with-trailing-solidus"
---| "eof-before-tag-name"
---| "eof-in-cdata"
---| "eof-in-comment"
---| "eof-in-doctype"
---| "eof-in-script-html-comment-like-text"
---| "eof-in-tag"
---| "incorrectly-closed-comment"
---| "incorrectly-opened-comment"
---| "invalid-character-sequence-after-doctype-name"
---| "invalid-first-character-of-tag-name"
---| "missing-attribute-value"
---| "missing-doctype-name"
---| "missing-doctype-public-identifier"
---| "missing-doctype-system-identifier"
---| "missing-end-tag-name"
---| "missing-quote-before-doctype-public-identifier"
---| "missing-quote-before-doctype-system-identifier"
---| "missing-semicolon-after-character-reference"
---| "missing-whitespace-after-doctype-public-keyword"
---| "missing-whitespace-after-doctype-system-keyword"
---| "missing-whitespace-before-doctype-name"
---| "missing-whitespace-between-attributes"
---| "missing-whitespace-between-doctype-public-and-system-identifiers"
---| "nested-comment"
---| "noncharacter-character-reference"
---| "noncharacter-in-input-stream"
---| "non-void-html-element-start-tag-with-trailing-solidus"
---| "null-character-reference"
---| "surrogate-character-reference"
---| "surrogate-in-input-stream"
---| "unexpected-character-after-doctype-system-identifier"
---| "unexpected-character-in-attribute-name"
---| "unexpected-character-in-unquoted-attribute-value"
---| "unexpected-equals-sign-before-attribute-name"
---| "unexpected-null-character"
---| "unexpected-question-mark-instead-of-tag-name"
---| "unexpected-solidus-in-tag"
---| "unknown-named-character-reference"

---@class HTML_Parser_Tokenizer
---@field sourceCode             SourceCode
---@field currentTagToken        HTML_Parser_TagToken?
---@field currentCommentToken    HTML_Parser_CommentToken?
---@field currentDoctypeToken    HTML_Parser_DoctypeToken?
---@field currentTagAttribute    HTML_Parser_TagAttr?
---@field isDuplicateAttribute   boolean
---@field lastStartTagToken      HTML_Parser_TagToken?
---@field currentState           HTML_Parser_TokenizerState
---@field returnState            HTML_Parser_TokenizerState?
---@field onTokenEmitted         fun(token:HTML_Parser_Token)
---@field currentInputChar       string   https://html.spec.whatwg.org/multipage/parsing.html#current-input-character
---@field characterReferenceCode integer  https://html.spec.whatwg.org/multipage/parsing.html#character-reference-code
---@field temporaryBuffer        string   https://html.spec.whatwg.org/multipage/parsing.html#temporary-buffer
local Tokenizer = {}


local DataState ---@type HTML_Parser_TokenizerState
local TagOpenState ---@type HTML_Parser_TokenizerState
local EndTagOpenState ---@type HTML_Parser_TokenizerState
local TagNameState ---@type HTML_Parser_TokenizerState
local RCDATALessThanSignState ---@type HTML_Parser_TokenizerState
local RCDATAEndTagOpenState ---@type HTML_Parser_TokenizerState
local RCDATAEndTagNameState ---@type HTML_Parser_TokenizerState
local BeforeAttributeNameState ---@type HTML_Parser_TokenizerState
local AttributeNameState ---@type HTML_Parser_TokenizerState
local AfterAttributeNameState ---@type HTML_Parser_TokenizerState
local BeforeAttributeValueState ---@type HTML_Parser_TokenizerState
local AttributeValueDoubleQuotedState ---@type HTML_Parser_TokenizerState
local AttributeValueSingleQuotedState ---@type HTML_Parser_TokenizerState
local AttributeValueUnquotedState ---@type HTML_Parser_TokenizerState
local AfterAttributeValueQuotedState ---@type HTML_Parser_TokenizerState
local MarkupDeclarationOpenState ---@type HTML_Parser_TokenizerState
local CommentStartState ---@type HTML_Parser_TokenizerState
local CommentStartDashState ---@type HTML_Parser_TokenizerState
local CommentState ---@type HTML_Parser_TokenizerState
local CommentLessThanSignState ---@type HTML_Parser_TokenizerState
local CommentLessThanSignBangState ---@type HTML_Parser_TokenizerState
local CommentLessThanSignBangDashState ---@type HTML_Parser_TokenizerState
local CommentLessThanSignBangDashDashState ---@type HTML_Parser_TokenizerState
local CommentEndDashState ---@type HTML_Parser_TokenizerState
local CommentEndState ---@type HTML_Parser_TokenizerState
local CommentEndBangState ---@type HTML_Parser_TokenizerState
local DOCTYPEState ---@type HTML_Parser_TokenizerState
local BeforeDOCTYPENameState ---@type HTML_Parser_TokenizerState
local DOCTYPENameState ---@type HTML_Parser_TokenizerState
local AfterDOCTYPENameState ---@type HTML_Parser_TokenizerState
local CharacterReferenceState ---@type HTML_Parser_TokenizerState
local NamedCharacterReferenceState ---@type HTML_Parser_TokenizerState
local NumericCharacterReferenceState ---@type HTML_Parser_TokenizerState
local DecimalCharacterReferenceStartState ---@type HTML_Parser_TokenizerState
local DecimalCharacterReferenceState ---@type HTML_Parser_TokenizerState
local NumericCharacterReferenceEndState ---@type HTML_Parser_TokenizerState


-- 1. https://html.spec.whatwg.org/multipage/parsing.html#data-state
DataState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    if currentChar == "&" then
        -- [U+0026 AMPERSAND (&)]

        -- Set the return state to the data state.
        t.returnState = DataState
        -- Switch to the character reference state.
        t:switchToState(CharacterReferenceState)
    elseif currentChar == "<" then
        -- [U+003C LESS-THAN SIGN (<)]

        -- Switch to the tag open state.
        t:switchToState(TagOpenState)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Emit the current input character as a character token.
        t:emitCharacterToken(currentChar)
    elseif currentChar == nil then
        -- [EOF]

        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Emit the current input character as a character token.
        t:emitCharacterToken(currentChar)
    end
end


-- 2. https://html.spec.whatwg.org/multipage/parsing.html#rcdata-state
RCDATAState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    if currentChar == "&" then
        -- [U+0026 AMPERSAND (&)]

        -- Set the return state to the RCDATA state. Switch to the character reference state.
        error("TODO")
    elseif currentChar == "<" then
        -- [U+003C LESS-THAN SIGN (<)]

        -- Switch to the RCDATA less-than sign state.
        t:switchToState(RCDATALessThanSignState)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Emit a U+FFFD REPLACEMENT CHARACTER character token.
        t:emitCharacterToken(0xfffd)
    elseif currentChar == nil then
        -- [EOF]

        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Emit the current input character as a character token.
        t:emitCharacterToken(currentChar)
    end
end


-- 6. https://html.spec.whatwg.org/multipage/parsing.html#tag-open-state
TagOpenState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    if currentChar == "!" then
        -- [U+0021 EXCLAMATION MARK (!)]

        -- Switch to the markup declaration open state.
        t:switchToState(MarkupDeclarationOpenState)
    elseif currentChar == "/" then
        -- [U+002F SOLIDUS (/)]

        -- Switch to the end tag open state.
        t:switchToState(EndTagOpenState)
    elseif currentChar ~= nil and codepoints.isAsciiAlpha(currentChar) then
        -- [ASCII alpha]

        -- Create a new start tag token, set its tag name to the empty string.
        t.currentTagToken =
            token.TagToken:new("", "start", t:getSourceLocation())
        -- Reconsume in the tag name state.
        t:reconsumeIn(TagNameState)
    elseif currentChar == "?" then
        -- [U+003F QUESTION MARK (?)]

        -- This is an unexpected-question-mark-instead-of-tag-name parse error.
        t:reportError("unexpected-question-mark-instead-of-tag-name")
        -- Create a comment token whose data is the empty string.
        t.currentCommentToken = token.CommentToken:new("", t:getSourceLocation())
        -- Reconsume in the bogus comment state.
        error("TODO")
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-before-tag-name parse error.
        t:reportError("eof-before-tag-name")
        -- Emit a U+003C LESS-THAN SIGN character token and an end-of-file token.
        t:emitCharacterToken(0x003c)
        t:emitEofToken()
    else
        -- [Anything else]

        -- This is an invalid-first-character-of-tag-name parse error.
        t:reportError("invalid-first-character-of-tag-name")
        -- Emit a U+003C LESS-THAN SIGN character token.
        t:emitCharacterToken(0x003c)
        -- Reconsume in the data state.
        t:reconsumeIn(DataState)
    end
end


-- 7. https://html.spec.whatwg.org/multipage/parsing.html#end-tag-open-state
EndTagOpenState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    if currentChar ~= nil and codepoints.isAsciiAlpha(currentChar) then
        -- [ASCII alpha]

        -- Create a new end tag token, set its tag name to the empty string. Reconsume in the
        -- tag name state.
        t:reconsumeIn(TagNameState)
        t.currentTagToken = token.TagToken:new("", "end", t:getSourceLocation())
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- This is a missing-end-tag-name parse error.
        t:reportError("missing-end-tag-name")
        -- Switch to the data state.
        t:switchToState(DataState)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-before-tag-name parse error.
        t:reportError("eof-before-tag-name")
        -- Emit a U+003C LESS-THAN SIGN character token, a U+002F SOLIDUS character token and an end-of-file token.
        t:emitCharacterToken(0x003c)
        t:emitCharacterToken(0x002f)
        t:emitEofToken()
    else
        -- [Anything else]

        -- This is an invalid-first-character-of-tag-name parse error. Create a comment token whose data
        -- is the empty string. Reconsume in the bogus comment state.
        t:reportError("invalid-first-character-of-tag-name")
        error("TODO")
    end
end


-- 8. https://html.spec.whatwg.org/multipage/parsing.html#tag-name-state
TagNameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- [U+0009 CHARACTER TABULATION (tab)]
        -- [U+000A LINE FEED (LF)]
        -- [U+000C FORM FEED (FF)]
        -- [U+0020 SPACE]

        -- Switch to the before attribute name state.
        t:switchToState(BeforeAttributeNameState)
    elseif currentChar == "/" then
        -- [U+002F SOLIDUS (/)]

        -- Switch to the self-closing start tag state.
        error("TODO")
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current tag token.
        t:emitCurrentTagToken()
    elseif currentChar ~= nil and codepoints.isAsciiUpperAlpha(currentChar) then
        -- [ASCII upper alpha]

        -- Append the lowercase version of the current input character (add 0x0020 to the character's code point) to the current tag token's tag name.
        t.currentTagToken.name = t.currentTagToken.name .. utf8.char(currentChar + 0x20)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the current tag token's tag name.
        t.currentTagToken.name = t.currentTagToken.name .. utf8.char(0xfffd)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-tag parse error.
        t:reportError("eof-in-tag")
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- Anything else

        -- Append the current input character to the current tag token's tag name.
        t.currentTagToken.name = t.currentTagToken.name + currentChar
    end
end


-- 9. https://html.spec.whatwg.org/multipage/parsing.html#rcdata-less-than-sign-state
RCDATALessThanSignState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002F SOLIDUS (/)
    if currentChar == "/" then
        -- Set the temporary buffer to the empty string.
        t.temporaryBuffer = ""
        -- Switch to the RCDATA end tag open state.
        t:switchToState(RCDATAEndTagOpenState)
    else
        -- [Anything else]
        -- Emit a U+003C LESS-THAN SIGN character token.
        t:emitCharacterToken(0x003c)
        -- Reconsume in the RCDATA state.
        t:reconsumeIn(RCDATAState)
    end
end


-- 10. https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-open-state
RCDATAEndTagOpenState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()
    if currentChar ~= nil and codepoints.isAsciiAlpha(currentChar) then
        -- Create a new end tag token, set its tag name to the empty string.
        t.currentTagToken =
            token.TagToken:new("", "end", t:getSourceLocation())
        -- Reconsume in the RCDATA end tag name state.
        t:reconsumeIn(RCDATAEndTagNameState)
    else
        -- Emit a U+003C LESS-THAN SIGN character token and a U+002F SOLIDUS character token.
        t:emitCharacterToken(0x003c)
        t:emitCharacterToken(0x002f)
        -- Reconsume in the RCDATA state.
        t:reconsumeIn(RCDATAState)
    end
end


-- 11. https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-name-state
RCDATAEndTagNameState = function(t)
    local anythingElse = function()
        -- Emit a U+003C LESS-THAN SIGN character token, a U+002F SOLIDUS character token, and a character token for
        -- each of the characters in the temporary buffer (in the order they were added to the buffer).
        t:emitCharacterToken(0x003c)
        t:emitCharacterToken(0x002f)
        for _, c in utf8.codes(t.temporaryBuffer) do
            t:emitCharacterToken(c)
        end
        -- Reconsume in the RCDATA state.
        t:reconsumeIn(RCDATAState)
    end
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- [U+0009 CHARACTER TABULATION (tab)]
        -- [U+000A LINE FEED (LF)]
        -- [U+000C FORM FEED (FF)]
        -- [U+0020 SPACE]

        -- If the current end tag token is an appropriate end tag token,
        if t:isAppropriateEndTagToken(t.currentTagToken) then
            --  then switch to the before attribute name state.
            t:switchToState(BeforeAttributeNameState)
            --  Otherwise, treat it as per the "anything else" entry below.
        else
            anythingElse()
        end
    elseif currentChar == "/" then
        -- [U+002F SOLIDUS (/)]

        -- If the current end tag token is an appropriate end tag token,
        if t:isAppropriateEndTagToken(t.currentTagToken) then
            -- then switch to the self-closing start tag state.
            error("TODO")
        else
            -- Otherwise, treat it as per the "anything else" entry below.
            anythingElse()
        end
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- If the current end tag token is an appropriate end tag token,
        if t:isAppropriateEndTagToken(t.currentTagToken) then
            -- then switch to the data state and emit the current tag token
            t:switchToState(DataState)
            t:emitCurrentTagToken()
        else
            -- Otherwise, treat it as per the "anything else" entry below.
            anythingElse()
        end
    elseif currentChar ~= nil and codepoints.isAsciiUpperAlpha(currentChar) then
        -- [ASCII upper alpha]

        -- Append the lowercase version of the current input character
        -- (add 0x0020 to the character's code point) to the current tag token's tag name.
        t.currentTagToken.name = t.currentTagToken.name .. utf8.char(currentChar + 0x20)
        -- Append the current input character to the temporary buffer.
        t.temporaryBuffer = t.temporaryBuffer .. currentChar
    elseif currentChar ~= nil and codepoints.isAsciiLowerAlpha(currentChar) then
        -- [ASCII lower alpha]

        -- Append the current input character to the current tag token's tag name.
        t.currentTagToken.name = t.currentTagToken.name .. currentChar
        -- Append the current input character to the temporary buffer.
        t.temporaryBuffer = t.temporaryBuffer .. currentChar
    else
        -- [Anything else]
        anythingElse()
    end
end


-- 32. https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-name-state
BeforeAttributeNameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()


    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- U+0009 CHARACTER TABULATION (tab)
        -- U+000A LINE FEED (LF)
        -- U+000C FORM FEED (FF)
        -- U+0020 SPACE

        -- Ignore the character.
    elseif currentChar == "/" or
        currentChar == ">" or
        currentChar == nil
    then
        -- [U+002F SOLIDUS (/)]
        -- [U+003E GREATER-THAN SIGN (>)]
        -- [EOF]

        -- Reconsume in the after attribute name state.
        t:reconsumeIn(AfterAttributeNameState)
    elseif currentChar == "=" then
        -- U+003D EQUALS SIGN (=)

        -- This is an unexpected-equals-sign-before-attribute-name parse error.
        t:reportError("unexpected-equals-sign-before-attribute-name")
        -- Start a new attribute in the current tag token.
        t:beginNewAttribute(token.TagToken.makeAttr(t:getSourceLocation()))
        -- Set that attribute's name to the current input character, and its value to the empty string.
        t.currentTagAttribute.name = "" .. currentChar
        t.currentTagAttribute.value = ""
        -- Switch to the attribute name state.
        t:switchToState(AttributeNameState)
        -- Anything else
    else
        -- Start a new attribute in the current tag token.
        t:beginNewAttribute(token.TagToken.makeAttr(t:getSourceLocation()))
        -- Set that attribute name and value to the empty string.
        t.currentTagAttribute.name = ""
        t.currentTagAttribute.value = ""
        -- Reconsume in the attribute name state.
        t:reconsumeIn(AttributeNameState)
    end
end


-- 33. https://html.spec.whatwg.org/multipage/parsing.html#attribute-name-state
AttributeNameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    local anythingElse = function()
        -- Append the current input character to the current attribute's name.
        t.currentTagAttribute.name = t.currentTagAttribute.name .. currentChar
    end

    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " " or
        currentChar == "/" or
        currentChar == ">" or
        currentChar == nil
    then
        -- [U+0009 CHARACTER TABULATION (tab)]
        -- [U+000A LINE FEED (LF)]
        -- [U+000C FORM FEED (FF)]
        -- [U+0020 SPACE]
        -- [U+002F SOLIDUS (/)]
        -- [U+003E GREATER-THAN SIGN (>)]
        -- [EOF]

        -- Reconsume in the after attribute name state.
        t:reconsumeIn(AfterAttributeNameState)
    elseif currentChar == "=" then
        -- [U+003D EQUALS SIGN (=)]
        -- Switch to the before attribute value state.
        t:switchToState(BeforeAttributeValueState)
    elseif codepoints.isAsciiUpperAlpha(currentChar) then
        -- [ASCII upper alpha]

        -- Append the lowercase version of the current input character (add 0x0020 to the character's
        -- code point) to the current attribute's name.
        t.currentTagAttribute.name = t.currentTagAttribute.name .. utf8.char(currentChar + 0x20)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the current attribute's name.
        t.currentTagAttribute.name = t.currentTagAttribute.name .. utf8.char(0xfffd)
    end

    -- U+0022 QUOTATION MARK (")
    -- U+0027 APOSTROPHE (')
    -- U+003C LESS-THAN SIGN (<)
    if currentChar == "\"" or
        currentChar == "'" or
        currentChar == "<" then
        -- This is an unexpected-character-in-attribute-name parse error.
        t:reportError("unexpected-character-in-attribute-name")
        -- Treat it as per the "anything else" entry below.
        anythingElse()
    else
        -- [Anything else]
        anythingElse()
    end
end


-- 34. https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-name-state
AfterAttributeNameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Ignore the character.
    elseif currentChar == "/" then
        -- [U+002F SOLIDUS (/)]

        -- Switch to the self-closing start tag state.
        error("TODO")
    elseif currentChar == "=" then
        -- [U+003D EQUALS SIGN (=)]

        -- Switch to the before attribute value state.
        t:switchToState(BeforeAttributeValueState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current tag token.
        t:emitCurrentTagToken()
    elseif currentChar == nil then
        -- [EOF]
        -- This is an eof-in-tag parse error.
        t:reportError("eof-in-tag")
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]
        -- Start a new attribute in the current tag token.
        t:beginNewAttribute(token.TagToken.makeAttr(t:getSourceLocation()))
        -- Set that attribute name and value to the empty string.
        t.currentTagAttribute.name = ""
        t.currentTagAttribute.value = ""
        -- Reconsume in the attribute name state.
        t:reconsumeIn(AttributeNameState)
    end
end


-- 35. https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-value-state
BeforeAttributeValueState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Ignore the character.
        return
    elseif currentChar == "\"" then
        -- [U+0022 QUOTATION MARK (")]

        -- Switch to the attribute value (double-quoted) state.
        t:switchToState(AttributeValueDoubleQuotedState)
    elseif currentChar == "'" then
        -- [U+0027 APOSTROPHE (')]

        -- Switch to the attribute value (single-quoted) state.
        t:switchToState(AttributeValueSingleQuotedState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- This is a missing-attribute-value parse error.
        t:reportError("missing-attribute-value")
        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current tag token.
        t:emitCurrentTagToken()
    else
        -- [Anything else]

        -- Reconsume in the attribute value (unquoted) state.
        t:reconsumeIn(AttributeValueUnquotedState)
    end
end


-- 36. https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(double-quoted)-state
AttributeValueDoubleQuotedState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0022 QUOTATION MARK (")
    if currentChar == "\"" then
        -- Switch to the after attribute value (quoted) state.
        t:switchToState(AfterAttributeValueQuotedState)
    elseif currentChar == "&" then
        -- [U+0026 AMPERSAND (&)]
        -- Set the return state to the attribute value (double-quoted) state.
        t.returnState = AttributeValueDoubleQuotedState
        -- Switch to the character reference state.
        t:switchToState(CharacterReferenceState)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the current attribute's value.
        t.currentTagAttribute.value = t.currentTagAttribute.value .. utf8.char(0xfffd)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-tag parse error.
        t:reportError("eof-in-tag")
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Append the current input character to the current attribute's value.
        t.currentTagAttribute.value = t.currentTagAttribute.value .. currentChar
    end
end


-- 37. https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(single-quoted)-state
AttributeValueSingleQuotedState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0027 APOSTROPHE (')
    if currentChar == "'" then
        -- Switch to the after attribute value (quoted) state.
        t:switchToState(AfterAttributeValueQuotedState)
    elseif currentChar == "&" then
        -- [U+0026 AMPERSAND (&)]

        -- Set the return state to the attribute value (signle-quoted) state.
        -- Switch to the character reference state.
        error("TODO")
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the current attribute's value.
        t.currentTagAttribute.value = t.currentTagAttribute.value .. utf8.char(0xfffd)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-tag parse error.
        t:reportError("eof-in-tag")
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]
        -- Append the current input character to the current attribute's value.
        t.currentTagAttribute.value = t.currentTagAttribute.value .. currentChar
    end
end


-- 38. https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(unquoted)-state
AttributeValueUnquotedState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Switch to the before attribute name state.
        t:switchToState(BeforeAttributeNameState)
    elseif currentChar == "&" then
        -- [U+0026 AMPERSAND (&)]

        -- Set the return state to the attribute value (unquoted) state.
        -- Switch to the character reference state.
        error("TODO")
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current tag token.
        t:emitCurrentTagToken()
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the current attribute's value.
        t.currentTagAttribute.value = t.currentTagAttribute.value .. utf8.char(0xfffd)
    end
    -- U+0022 QUOTATION MARK (")
    -- U+0027 APOSTROPHE (')
    -- U+003C LESS-THAN SIGN (<)
    -- U+003D EQUALS SIGN (=)
    -- U+0060 GRAVE ACCENT (`)
    if
        currentChar == "\"" or
        currentChar == "'" or
        currentChar == "<" or
        currentChar == "=" or
        currentChar == "`"
    then
        -- This is an unexpected-character-in-unquoted-attribute-value parse error.
        t:reportError("unexpected-character-in-unquoted-attribute-value")
        -- Append the current input character to the current attribute's value.
        t.currentTagAttribute.value = t.currentTagAttribute.value .. currentChar
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-tag parse error.
        t:reportError("eof-in-tag")
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]
        t.currentTagAttribute.value = t.currentTagAttribute.value .. currentChar
    end
end


-- 39. https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-value-(quoted)-state
AfterAttributeValueQuotedState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Switch to the before attribute name state.
        t:switchToState(BeforeAttributeNameState)
    elseif currentChar == "/" then
        -- [U+002F SOLIDUS (/)]

        -- Switch to the self-closing start tag state.
        error("TODO")
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Switch to the data state. Emit the current tag token.
        t:switchToState(DataState)
        t:emitCurrentTagToken()
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-tag parse error. Emit an end-of-file token.
        t:reportError("eof-in-tag")
        t:emitEofToken()
    else
        -- [Anything else]

        -- This is a missing-whitespace-between-attributes parse error.
        t:reportError("missing-whitespace-between-attributes")
        -- Reconsume in the before attribute name state.
        t:reconsumeIn(BeforeAttributeNameState)
    end
end


-- 42. https://html.spec.whatwg.org/multipage/parsing.html#markup-declaration-open-state
MarkupDeclarationOpenState = function(t)
    local remaining = t:peekAll()
    if strings.startsWith(remaining, "--") then
        -- [Two U+002D HYPHEN-MINUS characters (-)]

        -- Consume those two characters,
        t:consumeChars(#"--")
        -- create a comment token whose data is the empty string,
        t.currentCommentToken = token.CommentToken:new("", t:getSourceLocation())
        -- and switch to the comment start state.
        t:switchToState(CommentStartState)
    elseif strings.startsWith(strings.asciiLowercase(remaining), "--") then
        -- [ASCII case-insensitive match for the word "DOCTYPE"]

        -- Consume those characters
        t:consumeChars(#"DOCTYPE")
        -- and switch to the DOCTYPE state.
        t:switchToState(DOCTYPEState)
    elseif strings.startsWith(remaining, "[CDATA[") then
        -- Consume those characters.
        t:consumeChars(#"[CDATA[")
        -- If there is an adjusted current node and it is not an element in the HTML namespace, then switch to
        -- the CDATA section state. Otherwise, this is a cdata-in-html-content parse error.
        -- Create a comment token whose data is the "[CDATA[" string. Switch to the bogus comment state.
        error("TODO")
    else
        -- This is an incorrectly-opened-comment parse error.
        t:reportError("incorrectly-opened-comment")
        -- Create a comment token whose data is the empty string.
        t.currentCommentToken = token.CommentToken:new("", t:getSourceLocation())
        -- Switch to the bogus comment state (don't consume anything in the current state).
        error("TODO")
    end
end


-- 43. https://html.spec.whatwg.org/multipage/parsing.html#comment-start-state
CommentStartState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002D HYPHEN-MINUS (-)
    if currentChar == "-" then
        -- Switch to the comment start dash state.
        t:switchToState(CommentStartDashState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- This is an abrupt-closing-of-empty-comment parse error.
        t:reportError("abrupt-closing-of-empty-comment")
        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
        error("TODO")
    else
        -- [Anything else]

        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 44. https://html.spec.whatwg.org/multipage/parsing.html#comment-start-dash-state
CommentStartDashState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002D HYPHEN-MINUS (-)
    if currentChar == "-" then
        -- Switch to the comment end state.
        t:switchToState(CommentEndState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- This is an abrupt-closing-of-empty-comment parse error.
        t:reportError("abrupt-closing-of-empty-comment")
        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-comment parse error.
        t:reportError("eof-in-comment")
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
        return
    else
        -- [Anything else]

        -- Append a U+002D HYPHEN-MINUS character (-) to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. "-"
        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 45. https://html.spec.whatwg.org/multipage/parsing.html#comment-state
CommentState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+003C LESS-THAN SIGN (<)
    if currentChar == "<" then
        -- Append the current input character to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. currentChar
        -- Switch to the comment less-than sign state.
        t:switchToState(CommentLessThanSignState)
    elseif currentChar == "-" then
        -- [U+002D HYPHEN-MINUS (-)]

        -- Switch to the comment end dash state.
        t:switchToState(CommentEndDashState)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. utf8.char(0xfffd)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-comment parse error.
        t:reportError("eof-in-comment")
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
        return
    else
        -- [Anything else]

        -- Append the current input character to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. currentChar
    end
end


-- 46. https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-state
CommentLessThanSignState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0021 EXCLAMATION MARK (!)
    if currentChar == "!" then
        -- Append the current input character to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. currentChar
        -- Switch to the comment less-than sign bang state.
        t:switchToState(CommentLessThanSignBangState)
    elseif currentChar == "<" then
        -- [U+003C LESS-THAN SIGN (<)]
        -- Append the current input character to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. currentChar
    else
        -- [Anything else]

        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 47. https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-bang-state
CommentLessThanSignBangState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002D HYPHEN-MINUS (-)
    if currentChar == "-" then
        -- Switch to the comment less-than sign bang dash state.
        t:switchToState(CommentLessThanSignBangDashState)
    else
        -- [Anything else]

        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 48. https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-bang-dash-state
CommentLessThanSignBangDashState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002D HYPHEN-MINUS (-)
    if currentChar == "-" then
        -- Switch to the comment less-than sign bang dash dash state.
        t:switchToState(CommentLessThanSignBangDashDashState)
    else
        -- [Anything else]

        -- Reconsume in the comment end dash state.
        t:reconsumeIn(CommentEndDashState)
    end
end


-- 49. https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-bang-dash-dash-state
CommentLessThanSignBangDashDashState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+003E GREATER-THAN SIGN (>)
    if currentChar == ">" then
        -- Reconsume in the comment end state.
        t:reconsumeIn(CommentEndState)
    else
        -- [Anything else]

        -- This is a nested-comment parse error.
        t:reportError("nested-comment")
        -- Reconsume in the comment end state.
        t:reconsumeIn(CommentEndState)
    end
end


-- 50. https://html.spec.whatwg.org/multipage/parsing.html#comment-end-dash-state
CommentEndDashState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002D HYPHEN-MINUS (-)
    if currentChar == "-" then
        -- Reconsume in the comment end state.
        t:reconsumeIn(CommentEndState)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-comment parse error.
        t:reportError("eof-in-comment")
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Append a U+002D HYPHEN-MINUS character (-) to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. "-"
        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 51. https://html.spec.whatwg.org/multipage/parsing.html#comment-end-state
CommentEndState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+003E GREATER-THAN SIGN (>)
    if currentChar == ">" then
        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
    elseif currentChar == "!" then
        -- [U+0021 EXCLAMATION MARK (!)]

        -- Switch to the comment end bang state.
        t:switchToState(CommentEndBangState)
    elseif currentChar == "-" then
        -- [U+002D HYPHEN-MINUS (-)]

        -- Append a U+002D HYPHEN-MINUS character (-) to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. "-"
    elseif currentChar == nil then
        -- [EOF]
        -- This is an eof-in-comment parse error.
        t:reportError("eof-in-comment")
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Append two U+002D HYPHEN-MINUS characters (-) to the comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. "--"
        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 52. https://html.spec.whatwg.org/multipage/parsing.html#comment-end-bang-state
CommentEndBangState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+002D HYPHEN-MINUS (-)
    if currentChar == "-" then
        -- Append two U+002D HYPHEN-MINUS characters (-) and a U+0021 EXCLAMATION MARK character (!) to the
        -- comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. "--!"
        -- Switch to the comment end dash state.
        t:switchToState(CommentEndDashState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- This is an incorrectly-closed-comment parse error.
        t:reportError("incorrectly-closed-comment")
        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-comment parse error.
        t:reportError("eof-in-comment")
        -- Emit the current comment token.
        t:emitCurrentCommentToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Append two U+002D HYPHEN-MINUS characters (-) and a U+0021 EXCLAMATION MARK character (!) to the
        -- comment token's data.
        t.currentCommentToken.data = t.currentCommentToken.data .. "--!"
        -- Reconsume in the comment state.
        t:reconsumeIn(CommentState)
    end
end


-- 53. https://html.spec.whatwg.org/multipage/parsing.html#doctype-state
DOCTYPEState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Switch to the before DOCTYPE name state.
        t:switchToState(BeforeDOCTYPENameState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Reconsume in the before DOCTYPE name state.
        t:reconsumeIn(BeforeDOCTYPENameState)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-doctype parse error.
        t:reportError("eof-in-doctype")
        -- Create a new DOCTYPE token.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        -- Set its force-quirks flag to on.
        t.currentDoctypeToken.forceQuirks = true
        -- Emit the current token.
        t:emitCurrentDoctypeToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- This is a missing-whitespace-before-doctype-name parse error.
        t:reportError("missing-whitespace-before-doctype-name")
        -- Reconsume in the before DOCTYPE name state.
        t:reconsumeIn(BeforeDOCTYPENameState)
    end
end


-- 54. https://html.spec.whatwg.org/multipage/parsing.html#before-doctype-name-state
BeforeDOCTYPENameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Ignore the character.
    elseif currentChar ~= nil and codepoints.isAsciiUpperAlpha(currentChar) then
        -- [ASCII upper alpha]
        -- Create a new DOCTYPE token. Set the token's name to the lowercase version of the
        -- current input character (add 0x0020 to the character's code point). Switch to the
        -- DOCTYPE name state.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        t.currentDoctypeToken.name = "" + (currentChar + 0x20)
        t:switchToState(DOCTYPENameState)
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Create a new DOCTYPE token.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        -- Set the token's name to a U+FFFD REPLACEMENT CHARACTER character.
        t.currentDoctypeToken.name = "\u{fffd}"
        -- Switch to the DOCTYPE name state.
        t:switchToState(DOCTYPENameState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- This is a missing-doctype-name parse error.
        t:reportError("missing-doctype-name")
        -- Create a new DOCTYPE token.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        -- Set its force-quirks flag to on.
        t.currentDoctypeToken.forceQuirks = true
        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current token.
        t:emitCurrentDoctypeToken()
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-doctype parse error.
        t:reportError("eof-in-doctype")
        -- Create a new DOCTYPE token.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        -- Set its force-quirks flag to on.
        t.currentDoctypeToken.forceQuirks = true
        -- Emit the current token.
        t:emitCurrentDoctypeToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
    else
        -- [Anything else]

        -- Create a new DOCTYPE token. Set the token's name to the current input character. Switch to
        -- the DOCTYPE name state.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        t.currentDoctypeToken.name = "" + currentChar
        t:switchToState(DOCTYPENameState)
    end
end


-- 55. https://html.spec.whatwg.org/multipage/parsing.html#doctype-name-state
DOCTYPENameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Switch to the after DOCTYPE name state.
        t:switchToState(AfterDOCTYPENameState)
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current DOCTYPE token.
        t:emitCurrentDoctypeToken()
    elseif currentChar == "\u{0000}" then
        -- [U+0000 NULL]

        -- This is an unexpected-null-character parse error.
        t:reportError("unexpected-null-character")
        -- Append a U+FFFD REPLACEMENT CHARACTER character to the current DOCTYPE token's name.
        t.currentDoctypeToken.name = t.currentDoctypeToken.name .. utf8.char(0xfffd)
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-doctype parse error.
        t:reportError("eof-in-comment")
        -- Create a new DOCTYPE token.
        t.currentDoctypeToken = token.DoctypeToken:new(t:getSourceLocation())
        -- Set its force-quirks flag to on.
        t.currentDoctypeToken.forceQuirks = true
        -- Emit the current token.
        t:emitCurrentDoctypeToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
        return
    else
        -- [ASCII upper alpha]
        if codepoints.isAsciiUpperAlpha(currentChar) then
            -- Append the lowercase version of the current input character
            -- (add 0x0020 to the character's code point) to the current DOCTYPE token's name.
            t.currentDoctypeToken.name = t.currentDoctypeToken.name .. utf8.char(currentChar + 0x20)
        else
            -- [Anything else]
            -- Append the current input character to the current DOCTYPE token's name.
            t.currentDoctypeToken.name = t.currentDoctypeToken.name .. currentChar
        end
    end
end


-- 56. https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state
AfterDOCTYPENameState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0009 CHARACTER TABULATION (tab)
    -- U+000A LINE FEED (LF)
    -- U+000C FORM FEED (FF)
    -- U+0020 SPACE
    if
        currentChar == "\t" or
        currentChar == "\n" or
        currentChar == "\f" or
        currentChar == " "
    then
        -- Ignore the character.
    elseif currentChar == ">" then
        -- [U+003E GREATER-THAN SIGN (>)]

        -- Switch to the data state.
        t:switchToState(DataState)
        -- Emit the current DOCTYPE token.
        t:emitCurrentDoctypeToken()
    elseif currentChar == nil then
        -- [EOF]

        -- This is an eof-in-doctype parse error.
        t:reportError("eof-in-comment")
        -- Set the current DOCTYPE token's force-quirks flag to on.
        t.currentDoctypeToken.forceQuirks = true
        -- Emit the current DOCTYPE token.
        t:emitCurrentDoctypeToken()
        -- Emit an end-of-file token.
        t:emitEofToken()
        -- Anything else
    else
        local s = t:peek(6)
        -- If the six characters starting from the current input character are an ASCII case-insensitive
        -- match for the word "PUBLIC",
        if strings.asciiLowercase(s) == "public" then
            -- then consume those characters and switch to the after DOCTYPE public keyword state.
            t:consumeChars(#"public")
            error("TODO")
        elseif strings.asciiLowercase(s) == "system" then
            -- then consume those characters and switch to the after DOCTYPE system keyword state.
            t:consumeChars(#"system")
            error("TODO")
        else
            -- Set the current DOCTYPE token's force-quirks flag to on.
            t.currentDoctypeToken.forceQuirks = true
            -- Reconsume in the bogus DOCTYPE state.
            error("TODO")
        end
    end
end


-- 72. https://html.spec.whatwg.org/multipage/parsing.html#character-reference-state
CharacterReferenceState = function(t)
    -- Set the temporary buffer to the empty string.
    t.temporaryBuffer = ""
    -- Append a U+0026 AMPERSAND (&) character to the temporary buffer.
    t.temporaryBuffer = t.temporaryBuffer .. "&"
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0023 NUMBER SIGN (#)
    if currentChar == "#" then
        -- Append the current input character to the temporary buffer.
        t.currentDoctypeToken.name = t.temporaryBuffer .. currentChar
        -- Switch to the numeric character reference state.
        t:switchToState(NumericCharacterReferenceState)
    else
        -- [ASCII alphanumeric]
        if currentChar ~= nil and codepoints.isAsciiAlphanumeric(currentChar) then
            -- Reconsume in the named character reference state.
            t:reconsumeIn(NamedCharacterReferenceState)
        else
            -- [Anything else]

            -- Flush code points consumed as a character reference.
            t:flushCodepointsConsumedAsCharacterReference()
            -- Reconsume in the return state.
            t:reconsumeIn(t.returnState)
        end
    end
end


-- 73. https://html.spec.whatwg.org/multipage/parsing.html#named-character-reference-state
NamedCharacterReferenceState = function(t)
    -- Consume the maximum number of characters possible, where the consumed characters are one of the identifiers in the first column of the named character references table.
    local matchedWithoutSemicolon = false
    local matchedEntry = nil
    local matchedName = nil
    for name, entry in pairs(NAMED_CHARACTER_REFERENCES) do
        if t:peek(#name) == name then
            matchedName = name
            matchedEntry = entry.characters
            if strings.endsWith(name, ";") then
                matchedWithoutSemicolon = true
            end
            break
        end
    end
    --
    if matchedName ~= nil then
        t:consumeChars(matchedName.length)
        -- Append each character to the temporary buffer when it's consumed.
        t.temporaryBuffer = t.temporaryBuffer .. matchedName
    end
    -- If there is a match
    if matchedName ~= nil then
        -- If the character reference was consumed as part of an attribute,
        if
            t:consumedAsPartOfAttribute() and
            -- and the last character matched is not a U+003B SEMICOLON character (;),
            not matchedWithoutSemicolon and
            -- and the next input character is either a U+003D EQUALS SIGN character (=)
            (t:peek(1) ~= nil and (t:peek(1) == '=' or
                -- or an ASCII alphanumeric
                codepoints.isAsciiAlphanumeric(t:peek(1))))
        then
            -- then, for historical reasons, flush code points consumed as a character reference and switch
            -- to the return state.
            t:flushCodepointsConsumedAsCharacterReference()
            t:switchToState(t.returnState)
        else
            -- [Otherwise]:

            -- 1. If the last character matched is not a U+003B SEMICOLON character (;), then this is a missing-semicolon-after-character-reference parse error.
            if matchedWithoutSemicolon then
                t:reportError("missing-semicolon-after-character-reference")
            end
            -- 2. Set the temporary buffer to the empty string.
            t.temporaryBuffer = "";
            -- Append one or two characters corresponding to the character reference name (as given by the
            -- second column of the named character references table) to the temporary buffer.
            t.temporaryBuffer = t.temporaryBuffer .. matchedEntry;
            -- 3. Flush code points consumed as a character reference.
            t:flushCodepointsConsumedAsCharacterReference()
            -- Switch to the return state.
            t:switchToState(t.returnState)
        end
    else
        -- [Otherwise]

        -- Flush code points consumed as a character reference.
        t:flushCodepointsConsumedAsCharacterReference()
        -- Switch to the ambiguous ampersand state.
        error("TODO")
    end
end


-- 75. https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-state
NumericCharacterReferenceState = function(t)
    -- Set the character reference code to zero (0).
    t.characterReferenceCode = 0
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+0078 LATIN SMALL LETTER X
    -- U+0058 LATIN CAPITAL LETTER X
    if currentChar == "x" or currentChar == "X" then
        -- Append the current input character to the temporary buffer.
        t.currentDoctypeToken.name = t.temporaryBuffer .. currentChar
        -- Switch to the hexadecimal character reference start state.
        -- t:switchToState(HexadecimalCharacterReferenceStartState)
        error("TODO")
    else
        -- Reconsume in the decimal character reference start state.
        t:reconsumeIn(DecimalCharacterReferenceStartState)
    end
end


-- 77. https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-start-state
DecimalCharacterReferenceStartState = function(t)
    -- Consume the next input character:
    local c = t:consumeChar()
    -- ASCII digit
    if c ~= nil and codepoints.isAsciiDigit(c) then
        -- Reconsume in the decimal character reference state.
        t:reconsumeIn(DecimalCharacterReferenceState)
    else
        -- [Anything else]

        -- This is an absence-of-digits-in-numeric-character-reference parse error.
        t:reportError("absence-of-digits-in-numeric-character-reference")
        -- Flush code points consumed as a character reference.
        error("TODO")
        -- Reconsume in the return state.
        t:reconsumeIn(t.returnState)
    end
end


-- 79. https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-state
DecimalCharacterReferenceState = function(t)
    -- Consume the next input character:
    local currentChar = t:consumeChar()

    -- U+003B SEMICOLON
    if currentChar == ";" then
        -- Switch to the numeric character reference end state.
        t:switchToState(NumericCharacterReferenceEndState)
    else
        -- [ASCII digit]
        if currentChar ~= nil and codepoints.isAsciiDigit(currentChar) then
            -- Multiply the character reference code by 10.
            t.characterReferenceCode = t.characterReferenceCode * 10
            -- Add a numeric version of the current input character (subtract 0x0030 from the
            -- character's code point) to the character reference code.
            t.characterReferenceCode = t.characterReferenceCode + (currentChar - 0x30)
        else
            -- [Anything else]

            -- This is a missing-semicolon-after-character-reference parse error.
            t:reportError("missing-semicolon-after-character-reference")
            -- Reconsume in the numeric character reference end state.
            t:switchToState(NumericCharacterReferenceEndState)
        end
    end
end


-- 80. https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-end-state
NumericCharacterReferenceEndState = function(t)
    -- Check the character reference code:
    -- If the number is 0x00,
    if t.characterReferenceCode == 0x00 then
        -- then this is a null-character-reference parse error.
        t:reportError("null-character-reference")
        -- Set the character reference code to 0xFFFD.
        t.characterReferenceCode = 0xfffd
        -- If the number is greater than 0x10FFFF,
    elseif 0x10ffff < t.characterReferenceCode then
        -- then this is a character-reference-outside-unicode-range parse error.
        t:reportError("character-reference-outside-unicode-range")
        -- Set the character reference code to 0xFFFD.
        t.characterReferenceCode = 0xfffd
        -- If the number is a surrogate,
    elseif codepoints.isSurrogate(t.characterReferenceCode) then
        -- then this is a surrogate-character-reference parse error.
        t:reportError("surrogate-character-reference")
        -- Set the character reference code to 0xFFFD.
        t.characterReferenceCode = 0xfffd
        -- If the number is a noncharacter,
    elseif codepoints.isNonCharacter(t.characterReferenceCode) then
        -- then this is a noncharacter-character-reference parse error.
        t:reportError("noncharacter-character-reference")
        -- If the number is 0x0D, or a control that's not ASCII whitespace,
    elseif t.characterReferenceCode == 0x0d or
        codepoints.isISOControl(t.characterReferenceCode) and not codepoints.isAsciiWhitespace(t.characterReferenceCode)
    then
        -- then this is a control-character-reference parse error.
        t:reportError("control-character-reference")
        -- If the number is one of the numbers in the first column of the following table,
        -- then find the row with that number in the first column,
        local tbl = {
            { 0x80, 0x20AC }, -- EURO SIGN (€)
            { 0x82, 0x201A }, -- SINGLE LOW-9 QUOTATION MARK (‚)
            { 0x83, 0x0192 }, -- LATIN SMALL LETTER F WITH HOOK (ƒ)
            { 0x84, 0x201E }, -- DOUBLE LOW-9 QUOTATION MARK („)
            { 0x85, 0x2026 }, -- HORIZONTAL ELLIPSIS (…)
            { 0x86, 0x2020 }, -- DAGGER (†)
            { 0x87, 0x2021 }, -- DOUBLE DAGGER (‡)
            { 0x88, 0x02C6 }, -- MODIFIER LETTER CIRCUMFLEX ACCENT (ˆ)
            { 0x89, 0x2030 }, -- PER MILLE SIGN (‰)
            { 0x8A, 0x0160 }, -- LATIN CAPITAL LETTER S WITH CARON (Š)
            { 0x8B, 0x2039 }, -- SINGLE LEFT-POINTING ANGLE QUOTATION MARK (‹)
            { 0x8C, 0x0152 }, -- LATIN CAPITAL LIGATURE OE (Œ)
            { 0x8E, 0x017D }, -- LATIN CAPITAL LETTER Z WITH CARON (Ž)
            { 0x91, 0x2018 }, -- LEFT SINGLE QUOTATION MARK (‘)
            { 0x92, 0x2019 }, -- RIGHT SINGLE QUOTATION MARK (’)
            { 0x93, 0x201C }, -- LEFT DOUBLE QUOTATION MARK (“)
            { 0x94, 0x201D }, -- RIGHT DOUBLE QUOTATION MARK (”)
            { 0x95, 0x2022 }, -- BULLET (•)
            { 0x96, 0x2013 }, -- ,	EN DASH (–)
            { 0x97, 0x2014 }, -- ,	EM DASH (—)
            { 0x98, 0x02DC }, -- SMALL TILDE (˜)
            { 0x99, 0x2122 }, -- TRADE MARK SIGN (™)
            { 0x9A, 0x0161 }, -- LATIN SMALL LETTER S WITH CARON (š)
            { 0x9B, 0x203A }, -- SINGLE RIGHT-POINTING ANGLE QUOTATION MARK (›)
            { 0x9C, 0x0153 }, -- LATIN SMALL LIGATURE OE (œ)
            { 0x9E, 0x017E }, -- LATIN SMALL LETTER Z WITH CARON (ž)
            { 0x9F, 0x0178 }  -- LATIN CAPITAL LETTER Y WITH DIAERESIS (Ÿ)
        }
        for _, pair in ipairs(tbl) do
            if pair[1] == t.characterReferenceCode then
                -- and set the character reference code to the number in the second column of that row.
                t.characterReferenceCode = pair[2]
                break
            end
        end
    end
    -- Set the temporary buffer to the empty string.
    t.temporaryBuffer = ""
    -- Append a code point equal to the character reference code to the temporary buffer.
    t.temporaryBuffer = t.temporaryBuffer .. utf8.char(t.characterReferenceCode)
    -- Flush code points consumed as a character reference.
    t:flushCodepointsConsumedAsCharacterReference()
    -- Switch to the return state.
    t:switchToState(t.returnState)
end

---@param state HTML_Parser_TokenizerState
function Tokenizer:switchToState(state)
    if self.currentState == AttributeNameState then
        -- https://html.spec.whatwg.org/multipage/parsing.html#attribute-name-state
        -- When the user agent leaves the attribute name state (and before emitting the tag token, if appropriate),
        for _, attr in ipairs(self.currentTagToken.attributes) do
            -- the complete attribute's name must be compared to the other attributes on the same token;
            -- if there is already an attribute on the token with the exact same name,
            if attr.name == self.currentTagAttribute.name then
                -- then this is a duplicate-attribute parse error and the new attribute must be removed from the token.
                self:reportError("duplicate-attribute")
            end
        end
    end
    self.currentState = state
end

---@return SourceLocation
function Tokenizer:getSourceLocation()
    return {
        file = self.sourceCode.file,
        line = self.sourceCode.line,
        column = self.sourceCode.column,
    }
end

---@param error HTML_Parser_TokenizationError
function Tokenizer:reportError(error)
    ---@type HTML_Parser_TagAttr|HTML_Parser_Token
    local errorOn
    if self.currentTagAttribute ~= nil then
        errorOn = self.currentTagAttribute
    elseif self.currentTagToken ~= nil then
        errorOn = self.currentTagToken
    elseif self.currentCommentToken ~= nil then
        errorOn = self.currentCommentToken
    elseif self.currentDoctypeToken ~= nil then
        errorOn = self.currentDoctypeToken
    end

    local file, startLine, startCol, endLine, endCol
    file = self.sourceCode.file
    if errorOn ~= nil then
        startLine = errorOn.startLocation.line
        startCol = errorOn.startLocation.column
        endLine = errorOn.endLocation.line
        endCol = errorOn.endLocation.column
        if endLine == 0 or endCol == 0 then
            endLine = self.sourceCode.line
            endCol = self.sourceCode.column
        end
    else
        startLine = self.sourceCode.line
        startCol = self.sourceCode.column
        endLine = self.sourceCode.line
        endCol = self.sourceCode.column
    end

    L:e("%s: Error %s occured", tostring(file), error)
    self.sourceCode:printRange(L, "error", startLine, startCol, endLine, endCol)
end

function Tokenizer:addCurrentAttributeIfNeeded()
    if self.currentTagAttribute == nil then
        return
    end
    if self.isDuplicateAttribute then
        self.isDuplicateAttribute = false
    end
end

---@param attr HTML_Parser_TagAttr
function Tokenizer:beginNewAttribute(attr)
    self:addCurrentAttributeIfNeeded()
    self.currentTagAttribute = attr
end

---https://html.spec.whatwg.org/multipage/parsing.html#appropriate-end-tag-token
---@param tk HTML_Parser_TagToken
---@return boolean
function Tokenizer:isAppropriateEndTagToken(tk)
    assert(tk.kind == "end")
    return tk.name == self.lastStartTagToken.name
end

---https://html.spec.whatwg.org/multipage/parsing.html#charref-in-attribute
function Tokenizer:consumedAsPartOfAttribute()
    -- A character reference is said to be consumed as part of an attribute if the return state is either
    return
    -- attribute value (double-quoted) state,
        self.returnState == AttributeValueDoubleQuotedState or
        -- attribute value (single-quoted) state,
        self.returnState == AttributeValueSingleQuotedState or
        -- or attribute value (unquoted) state.
        self.returnState == AttributeValueUnquotedState
end

---https://html.spec.whatwg.org/multipage/parsing.html#flush-code-points-consumed-as-a-character-reference
function Tokenizer:flushCodepointsConsumedAsCharacterReference()
    -- When a state says to flush code points consumed as a character reference,
    -- it means that for each code point in the temporary buffer (in the order they were added to the buffer)
    for _, c in utf8.codes(self.temporaryBuffer) do
        if self:consumedAsPartOfAttribute() then
            -- user agent must append the code point from the buffer to the current attribute's value
            -- if the character reference was consumed as part of an attribute,
            self.currentTagAttribute.value = self.currentTagAttribute.value .. utf8.char(c)
        else
            -- or emit the code point as a character token otherwise.
            self:emitCharacterToken(c)
        end
    end
end

---@param tok HTML_Parser_Token
function Tokenizer:emitToken(tok)
    if tok.type == "tag" and (tok --[[@as HTML_Parser_TagToken]]).kind == "start" then
        self.lastStartTagToken = tok --[[@as HTML_Parser_TagToken]]
    end
    self.onTokenEmitted(tok)
end

function Tokenizer:emitEofToken()
    self:emitToken(token.EofToken:new(self:getSourceLocation()))
end

---@param char integer|string
function Tokenizer:emitCharacterToken(char)
    if type(char) == "string" then char = utf8.codepoint(char) end
    self:emitToken(token.CharacterToken:new(char, self:getSourceLocation()))
end

function Tokenizer:emitCurrentCommentToken()
    assert(self.currentCommentToken ~= nil)
    self:emitToken(self.currentCommentToken)
    self.currentCommentToken = nil
end

function Tokenizer:emitCurrentDoctypeToken()
    assert(self.currentDoctypeToken ~= nil)
    self:emitToken(self.currentDoctypeToken)
    self.currentDoctypeToken = nil
end

function Tokenizer:emitCurrentTagToken()
    assert(self.currentTagToken ~= nil)
    self:addCurrentAttributeIfNeeded()
    self:emitToken(self.currentTagToken)
    self.currentTagToken = nil
end

---@return string|nil
function Tokenizer:consumeChar()
    -- https://html.spec.whatwg.org/multipage/parsing.html#reconsume
    if self.shouldReconsume then
        -- When a state says to reconsume a matched character in a specified state, that means to switch to that state, but when it attempts to consume the next input character, provide it with the current input character instead.
        self.shouldReconsume = false
        return self.currentInputChar
    end

    local c = self.sourceCode:consume()
    if c == nil then
        return nil
    end
    self.currentInputChar = c
    return self.currentInputChar
end

---@param count integer
---@return string
function Tokenizer:peek(count)
    local result = ""
    if self.shouldReconsume then
        result = self.currentInputChar
        count = count - 1
    end
    return result .. self.sourceCode:peek(count)
end

---@return string
function Tokenizer:peekAll()
    local result = ""
    if self.shouldReconsume then
        result = self.currentInputChar
    end
    return result .. self.sourceCode:peekAll()
end

---@param count number
function Tokenizer:consumeChars(count)
    for _ = 1, count do
        self:consumeChar()
    end
end

---https://html.spec.whatwg.org/multipage/parsing.html#reconsume
---@param state HTML_Parser_TokenizerState
function Tokenizer:reconsumeIn(state)
    -- When a state says to reconsume a matched character in a specified state, that means to switch to that state,
    self:switchToState(state)
    -- but when it attempts to consume the next input character, provide it with the current input character instead.
    self.shouldReconsume = true
end

function Tokenizer:new()

end

return Tokenizer
