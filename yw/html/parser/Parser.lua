--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Tokenizer    = require "yw.html.parser.Tokenizer"
local DocumentType = require "yw.dom.DocumentType"
local Element      = require "yw.dom.Element"
local token        = require "yw.html.parser.token"
local strings      = require "yw.common.strings"
local Logger       = require "yw.common.Logger"
local namespaces   = require "yw.common.namespaces"
local HTMLElement  = require "yw.html.HTMLElement"


local L = Logger:new("yw.html.parser.Parser")


local SCRIPTING_ENABLED = false
-- We don't have speculative HTML parsing, and this just acts as marker for if we ever decide to implement it.
local ACTIVE_SPECULATIVE_HTML_PARSER = nil

---@class HTML_Parser_Parser
---@field tokenizer                                                  HTML_Parser_Tokenizer
---@field document                                                   HTML_Document
---@field currentInsertionMode                                       HTML_Parser_InsertionMode
---@field originalInsertionMode                                      HTML_Parser_InsertionMode?
---@field sourceCode                                                 SourceCode
---@field createdAsPartOfHTMLFragmentParsingAlgorithm                boolean
---@field invokedViaDocumentWrite                                    boolean
---@field enableFoesterParenting                                     boolean                                     https://html.spec.whatwg.org/multipage/parsing.html#foster-parent
---@field headElementPointer                                         HTML_HTMLHeadElement                        https://html.spec.whatwg.org/multipage/parsing.html#head-element-pointer
---@field formElementPointer                                         HTML_HTMLFormElement                        https://html.spec.whatwg.org/multipage/parsing.html#form-element-pointer
---@field stackOfOpenElements                                        HTML_Parser_StackOfOpenElements             https://html.spec.whatwg.org/multipage/parsing.html#stack-of-open-elements
---@field stackOfTemplateInsertionModes                              HTML_Parser_InsertionMode[]                 https://html.spec.whatwg.org/multipage/parsing.html#stack-of-template-insertion-modes
---@field listOfActiveFormattingElements                             HTML_Parser_ListOfActiveFormattingElements  https://html.spec.whatwg.org/multipage/parsing.html#list-of-active-formatting-elements
---@field framesetOKFlag                                             "ok"|"not ok"                               https://html.spec.whatwg.org/multipage/parsing.html#frameset-ok-flag
---@field listOfScriptsThatWillExecuteWhenDocumentHasFinishedParsing DOM_Element[]                               https://html.spec.whatwg.org/multipage/parsing.html#list-of-scripts-that-will-execute-when-the-document-has-finished-parsing
local Parser = {}

---@alias HTML_Parser_InsertionMode fun(p: HTML_Parser_Parser, tk: HTML_Parser_Token_Union)


local InitialInsertionMode ---@type HTML_Parser_InsertionMode
local BeforeHTMLInsertionMode ---@type HTML_Parser_InsertionMode
local BeforeHeadInsertionMode ---@type HTML_Parser_InsertionMode
local AfterHeadInsertionMode ---@type HTML_Parser_InsertionMode
local InBodyInsertionMode ---@type HTML_Parser_InsertionMode
local InHeadInsertionMode ---@type HTML_Parser_InsertionMode
local InHeadNoScriptInsertionMode ---@type HTML_Parser_InsertionMode
local InTemplateInsertionMode ---@type HTML_Parser_InsertionMode

---@param node DOM_Element
---@param localName string
local function isNodeInsideHTMLElementInclusive(node, localName)
    if node:isHTMLElement(localName) then
        return true
    end
    for ancestor in node:ancestors() do
        if ancestor.isElement and (ancestor --[[@as DOM_Element]]):isHTMLElement(localName) then
            return true
        end
    end
    return false
end

---https://html.spec.whatwg.org/multipage/parsing.html#adjusted-current-node
---@param p HTML_Parser_Parser
local function adjustedCurrentNode(p)
    -- The adjusted current node is the context element if the parser was created as part of the HTML fragment parsing algorithm and the stack of open elements has only one element in it (fragment case);
    if p.createdAsPartOfHTMLFragmentParsingAlgorithm then
        error("todo")
    end
    -- otherwise, the adjusted current node is the current node.
    return p.stackOfOpenElements:currentNode()
end

---https://html.spec.whatwg.org/multipage/parsing.html#create-an-element-for-the-token
---@param p HTML_Parser_Parser
---@param tk HTML_Parser_Token_Union
---@param namespace string
---@param parent DOM_Node
---@return DOM_Element
local function createElementForToken(p, tk, namespace, parent)
    -- 1. If the active speculative HTML parser is not null,
    if ACTIVE_SPECULATIVE_HTML_PARSER then
        -- then return the result of creating a speculative mock element given given namespace, the tag name of the given token, and the attributes of the given token.
        error("not implemented")
    else
        -- 2. Otherwise, optionally create a speculative mock element given given namespace, the tag name of the given token, and the attributes of the given token.
    end

    -- NOTE: We don't have speculative HTML parsing support, so above steps are not applicable.

    -- 3. Let document be intended parent's node document.
    local document = parent.nodeDocument --[[@as HTML_Document]]
    assert(document.type == "html")


    -- 4. Let local name be the tag name of the token.
    local localName = token.tag

    -- 5. Let is be the value of the "is" attribute in the given token, if such an attribute exists, or null otherwise.
    local is = tk.attributes["is"]

    -- 6. Let definition be the result of looking up a custom element definition given document, given namespace, local name, and is.
    local definition = document:lookupCustomElementDefinition(namespace, localName, is)

    -- 7. Let willExecuteScript be true if definition is non-null and the parser was not created as part of the HTML fragment parsing algorithm; otherwise false.
    local willExecuteScript = definition ~= nil and not p.createdAsPartOfHTMLFragmentParsingAlgorithm

    -- 8. If willExecuteScript is true:
    if willExecuteScript then
        -- 1. Increment document's throw-on-dynamic-markup-insertion counter.
        document.throwOnDynamicMarkupInsertionCounter = document.throwOnDynamicMarkupInsertionCounter + 1
        -- 2. If the JavaScript execution context stack is empty, then perform a microtask checkpoint.
        error("todo")
        -- 3. Push a new element queue onto document's relevant agent's custom element reactions stack.
    end

    -- 9. Let element be the result of creating an element given document, localName, given namespace, null, is, and willExecuteScript.
    local element = Element.create(document, localName, namespace, nil, is, willExecuteScript)

    -- 10. Append each attribute in the given token to element.
    for _, attr in ipairs(tk.attributes) do
        element:setAttributeValue(attr.name, attr.value)
    end

    -- 11. If willExecuteScript is true:
    if willExecuteScript then
        -- 1. Let queue be the result of popping from document's relevant agent's custom element reactions stack. (This will be the same element queue as was pushed above.)
        error("todo")
        -- 2. Invoke custom element reactions in queue.
        -- 3. Decrement document's throw-on-dynamic-markup-insertion counter.
        document.throwOnDynamicMarkupInsertionCounter = document.throwOnDynamicMarkupInsertionCounter - 1
    end

    -- 12. If element has an xmlns attribute in the XMLNS namespace whose value is not exactly the same as the element's namespace, that is a parse error. Similarly, if element has an xmlns:xlink attribute in the XMLNS namespace whose value is not the XLink Namespace,
    if element:getAttribute(namespaces.XMLNS_NAMESPACE, "xmlns")
        and element:getAttributeValue(namespaces.XMLNS_NAMESPACE, "xmlns") ~= namespace
    then
        -- that is a parse error.
        reportError(p, tk)
    end
    -- Similarly, if element has an xmlns:xlink attribute in the XMLNS namespace whose value is not the XLink Namespace,
    if element:getAttribute(namespaces.XMLNS_NAMESPACE, "xmlns:xlink")
        and element:getAttributeValue(namespaces.XMLNS_NAMESPACE, "xmlns:xlink") ~= namespaces.XLINK_NAMESPACE
    then
        -- that is a parse error.
        reportError(p, tk)
    end
    -- 13. If element is a resettable element, invoke its reset algorithm. (This initializes the element's value and checkedness based on the element's attributes.)
    if element:isHTMLElement() and (element --[[@as HTML_HTMLElement]]):isFormAssociatedElement(HTMLElement.FORM_RESETTABLE_ELEMENTS) then
        (element --[[@as HTML_HTMLElement]]):reset()
    end
    -- 14. If element is a form-associated element and not a form-associated custom element,
    if element:isHTMLElement() and (
            (element --[[@as HTML_HTMLElement]]):isFormAssociatedElement(HTMLElement.FORM_RESETTABLE_ELEMENTS)
            and not (element --[[@as HTML_HTMLElement]]):isFormAssociatedCustomElement()
        )
        -- the form element pointer is not null,
        and p.formElementPointer ~= nil
        -- there is no template element on the stack of open elements,
        and p.stackOfOpenElements:hasHTMLElement("template")
    -- element is either not listed or doesn't have a form attribute,
    -- TODO
    -- and the intended parent is in the same tree as the element pointed to by the form element pointer,
    -- TODO
    then
        -- then associate element with the form element pointed to by the form element pointer and set element's parser inserted flag.
        error("todo")
    end
    -- 15. Return element
    return element
end

---@alias HTML_Parser_AdjustedInsertionLocation { insertionParent: DOM_Element, beforeChild: DOM_Element }

---https://html.spec.whatwg.org/multipage/parsing.html#appropriate-place-for-inserting-a-node
---@param p HTML_Parser_Parser
---@param overrideTarget DOM_Element?
---@return HTML_Parser_AdjustedInsertionLocation
local function appropriatePlaceForInsertingNode(p, overrideTarget)
    -- 1. If there was an override target specified, then let target be the override target. Otherwise, let target be the current node.
    local target = overrideTarget or p.stackOfOpenElements:currentNode()

    -- 2. Determine the adjusted insertion location using the first matching steps from the following list:
    local adjustedInsertionLocation

    if p.enableFoesterParenting              -- If foster parenting is enabled
        and (                                --  and target is a
            target:isHTMLElement("table")    -- table,
            or target:isHTMLElement("tbody") -- tbody,
            or target:isHTMLElement("tfoot") -- tfoot,
            or target:isHTMLElement("thead") -- thead,
            or target:isHTMLElement("tr")    -- or tr element
        ) then
        -- Run these substeps:

        -- 1. Let last template be the last template element in the stack of open elements, if any.
        error("TODO")

        -- 2. Let last table be the last table element in the stack of open elements, if any.

        -- 3. If there is a last template and either there is no last table, or there is one, but last template is lower (more recently added) than last table in the stack of open elements, then: let adjusted insertion location be inside last template's template contents, after its last child (if any), and abort these steps.

        -- 4. If there is no last table, then let adjusted insertion location be inside the first element in the stack of open elements (the html element), after its last child (if any), and abort these steps. (fragment case)

        -- 5. If last table has a parent node, then let adjusted insertion location be inside last table's parent node, immediately before last table, and abort these steps.

        -- 6. Let previous element be the element immediately above last table in the stack of open elements.

        -- 7. Let adjusted insertion location be inside previous element, after its last child (if any).
    else -- Otherwise
        -- Let adjusted insertion location be
        return {
            insertionParent = target, -- inside target,
            beforeChild = nil,        -- after its last child (if any).
        }
    end

    -- 3. If the adjusted insertion location is inside a template element,
    if isNodeInsideHTMLElementInclusive(adjustedInsertionLocation.insertionParent, "template") then
        -- let it instead be inside the template element's template contents, after its last child (if any).
        error("todo")
    end

    -- 4. Return the adjusted insertion location.
    return adjustedInsertionLocation
end

---https://html.spec.whatwg.org/multipage/parsing.html#insert-an-element-at-the-adjusted-insertion-location
---@param p HTML_Parser_Parser
---@param element DOM_Element
local function insertElementAtAdjustedInsertionLocation(p, element)
    -- 1. Let the adjusted insertion location be the appropriate place for inserting a node.
    local adjustedInsertionLocation = appropriatePlaceForInsertingNode(p)

    -- 2. If it is not possible to insert element at the adjusted insertion location,
    if not pcall(
            element.ensurePreInsertionValidity,
            element,
            adjustedInsertionLocation.insertionParent,
            adjustedInsertionLocation.beforeChild
        ) then
        -- abort these steps.
        return
    end

    -- 3. If the parser was not created as part of the HTML fragment parsing algorithm,
    if not p.createdAsPartOfHTMLFragmentParsingAlgorithm then
        -- then push a new element queue onto element's relevant agent's custom element reactions stack.
        -- TODO
    end

    -- 4. Insert element at the adjusted insertion location.
    element:preInsert(adjustedInsertionLocation.insertionParent, adjustedInsertionLocation.beforeChild)

    -- 5. If the parser was not created as part of the HTML fragment parsing algorithm,
    if not p.createdAsPartOfHTMLFragmentParsingAlgorithm then
        -- then pop the element queue from element's relevant agent's custom element reactions stack, and invoke custom element reactions in that queue.
        -- TODO
    end
end

---https://html.spec.whatwg.org/multipage/parsing.html#insert-a-foreign-element
---@param p HTML_Parser_Parser
---@param tk HTML_Parser_Token_Union
---@param namespace string
---@param onlyAddElementToStack boolean
---@return DOM_Element
local function insertForeignElement(p, tk, namespace, onlyAddElementToStack)
    -- 1. Let the adjusted insertion location be the appropriate place for inserting a node.
    local adjustedInsertionLocation = appropriatePlaceForInsertingNode(p)

    -- 2. Let element be the result of creating an element for the token in the given namespace, with the intended parent being the element in which the adjusted insertion location finds itself.
    local element = createElementForToken(p, tk, namespace, adjustedInsertionLocation.insertionParent)

    -- 3. If onlyAddToElementStack is false,
    if not onlyAddElementToStack then
        -- then run insert an element at the adjusted insertion location with element
        insertElementAtAdjustedInsertionLocation(p, element)
    end

    -- 4. Push element onto the stack of open elements so that it is the new current node.
    p.stackOfOpenElements:push(element)

    -- 5. Return element.
    return element
end

---https://html.spec.whatwg.org/multipage/parsing.html#insert-an-html-element
---@param p HTML_Parser_Parser
---@param tk HTML_Parser_Token_Union
---@return HTML_HTMLElement
local function insertHTMLElement(p, tk)
    return insertForeignElement(p, tk, namespaces.HTML_NAMESPACE, false) --[[@as HTML_HTMLElement]]
end

---@param tk HTML_Parser_Token_Union
---@param t string[]
---@return boolean
local function isCharacterTokenWith(tk, t)
    if tk.type ~= "character" then
        return false
    end
    for _, c in ipairs(t) do
        if tk.char == utf8.codepoint(c) then
            return true
        end
    end
    return false
end

---@param tk HTML_Parser_Token_Union
---@return boolean
local function isStartTag(tk)
    if tk.type ~= "tag" or tk.kind ~= "start" then
        return false
    end
    return false
end

---@param tk HTML_Parser_Token_Union
---@param t string[]
---@return boolean
local function isStartTagWithName(tk, t)
    if tk.type ~= "tag" or tk.kind ~= "start" then
        return false
    end
    for _, name in ipairs(t) do
        if tk.name == name then
            return true
        end
    end
    return false
end

---@param tk HTML_Parser_Token_Union
---@return boolean
local function isEndTag(tk)
    if tk.type ~= "tag" or tk.kind ~= "end" then
        return false
    end
    return false
end

---@param tk HTML_Parser_Token_Union
---@param t string[]
---@return boolean
local function isEndTagWithName(tk, t)
    if tk.type ~= "tag" or tk.kind ~= "end" then
        return false
    end
    for _, name in ipairs(t) do
        if tk.name == name then
            return true
        end
    end
    return false
end

---https://html.spec.whatwg.org/multipage/parsing.html#the-initial-insertion-mode
InitialInsertionMode = function(p, tk)
    -- When the user agent is to apply the rules for the "initial" insertion mode, the user agent must handle the token as follows:

    --> A character token that is one of U+0009 CHARACTER TABULATION, U+000A LINE FEED (LF), U+000C FORM FEED (FF), U+000D CARRIAGE RETURN (CR), or U+0020 SPACE
    if isCharacterTokenWith(tk, { "\t", "\n", "\f", "\r", " " }) then
        -- Ignore the token.
        return
    end
    --> A comment token
    if tk.type == "comment" then
        -- Insert a comment as the last child of the Document object.
        error("todo")
        return
    end
    --> A DOCTYPE token
    if tk.type == "doctype" then
        -- If the DOCTYPE token's name is not "html", or the token's public identifier is not missing, or the token's system identifier is neither missing nor "about:legacy-compat",
        if tk.name ~= "html"
            or tk.publicIdentifier ~= nil
            or (
                tk.systemIdentifier ~= nil and
                tk.systemIdentifier ~= "about:legacy-compat"
            )
        then
            -- then there is a parse error.
            reportError(p, tk)
        end
        -- Append a DocumentType node to the Document node,
        p.document:appendChild(
            DocumentType:new(
                p.document,
                -- with its name set to the name given in the DOCTYPE token, or the empty string if the name was missing;
                tk.name or "",
                -- its public ID set to the public identifier given in the DOCTYPE token, or the empty string if the public identifier was missing;
                tk.publicIdentifier or "",
                --  and its system ID set to the system identifier given in the DOCTYPE token, or the empty string if the system identifier was missing.
                tk.systemIdentifier or ""
            )
        )
        -- Then,

        -- NOTE: The standard says that system and public IDs must be compared in ASCII case-insensitive manner.
        local function isPublicIDSetTo(t)
            if tk.publicIdentifier == nil then
                return false
            end
            local id = strings.asciiLowercase(tk.publicIdentifier)
            for _, case in ipairs(t) do
                if id == case then
                    return true
                end
            end
            return false
        end

        local function isSystemIDSetTo(t)
            if tk.systemIdentifier == nil then
                return false
            end
            local id = strings.asciiLowercase(tk.systemIdentifier)
            for _, case in ipairs(t) do
                if id == case then
                    return true
                end
            end
            return false
        end

        local function doesPublicIDStartWith(t)
            if tk.publicIdentifier == nil then
                return false
            end
            local id = strings.asciiLowercase(tk.publicIdentifier)
            for _, case in ipairs(t) do
                if strings.startsWith(t, case) then
                    return true
                end
            end
            return false
        end

        if not p.document.isIFrameSrcdocDocument             -- if the document is not an iframe srcdoc document,
            and not p.document.parserCannotChangeTheModeFlag -- and the parser cannot change the mode flag is false,
            -- and the DOCTYPE token matches one of the conditions in the following list,
            and (
                tk.forceQuirks                              -- The force-quirks flag is set to on.
                or tk.name ~= "html"                        -- The name is not "html".
                or isPublicIDSetTo {
                    "-//w3o//dtd w3 html strict 3.0//en//", -- The public identifier is set to: "-//W3O//DTD W3 HTML Strict 3.0//EN//"
                    "-/w3c/dtd html 4.0 transitional/en",   -- The public identifier is set to: "-/W3C/DTD HTML 4.0 Transitional/EN"
                    "html"                                  -- The public identifier is set to: "HTML"
                }
                or isSystemIDSetTo {
                    "http://www.ibm.com/data/dtd/v11/ibmxhtml1-transitional.dtd", -- The system identifier is set to: "http://www.ibm.com/data/dtd/v11/ibmxhtml1-transitional.dtd"
                }
                or doesPublicIDStartWith {
                    "+//silmaril//dtd html pro v0r11 19970101//",                                     -- The public identifier starts with: "+//Silmaril//dtd html Pro v0r11 19970101//"
                    "-//as//dtd html 3.0 aswedit + extensions//",                                     -- The public identifier starts with: "-//AS//DTD HTML 3.0 asWedit + extensions//"
                    "-//advasoft ltd//dtd html 3.0 aswedit + extensions//",                           -- The public identifier starts with: "-//AdvaSoft Ltd//DTD HTML 3.0 asWedit + extensions//"
                    "-//ietf//dtd html 2.0 level 1//",                                                -- The public identifier starts with: "-//IETF//DTD HTML 2.0 Level 1//"
                    "-//ietf//dtd html 2.0 level 2//",                                                -- The public identifier starts with: "-//IETF//DTD HTML 2.0 Level 2//"
                    "-//ietf//dtd html 2.0 strict level 1//",                                         -- The public identifier starts with: "-//IETF//DTD HTML 2.0 Strict Level 1//"
                    "-//ietf//dtd html 2.0 strict level 2//",                                         -- The public identifier starts with: "-//IETF//DTD HTML 2.0 Strict Level 2//"
                    "-//ietf//dtd html 2.0 strict//",                                                 -- The public identifier starts with: "-//IETF//DTD HTML 2.0 Strict//"
                    "-//ietf//dtd html 2.0//",                                                        -- The public identifier starts with: "-//IETF//DTD HTML 2.0//"
                    "-//ietf//dtd html 2.1e//",                                                       -- The public identifier starts with: "-//IETF//DTD HTML 2.1E//"
                    "-//ietf//dtd html 3.0//",                                                        -- The public identifier starts with: "-//IETF//DTD HTML 3.0//"
                    "-//ietf//dtd html 3.2 final//",                                                  -- The public identifier starts with: "-//IETF//DTD HTML 3.2 Final//"
                    "-//ietf//dtd html 3.2//",                                                        -- The public identifier starts with: "-//IETF//DTD HTML 3.2//"
                    "-//ietf//dtd html 3//",                                                          -- The public identifier starts with: "-//IETF//DTD HTML 3//"
                    "-//ietf//dtd html level 0//",                                                    -- The public identifier starts with: "-//IETF//DTD HTML Level 0//"
                    "-//ietf//dtd html level 1//",                                                    -- The public identifier starts with: "-//IETF//DTD HTML Level 1//"
                    "-//ietf//dtd html level 2//",                                                    -- The public identifier starts with: "-//IETF//DTD HTML Level 2//"
                    "-//ietf//dtd html level 3//",                                                    -- The public identifier starts with: "-//IETF//DTD HTML Level 3//"
                    "-//ietf//dtd html strict level 0//",                                             -- The public identifier starts with: "-//IETF//DTD HTML Strict Level 0//"
                    "-//ietf//dtd html strict level 1//",                                             -- The public identifier starts with: "-//IETF//DTD HTML Strict Level 1//"
                    "-//ietf//dtd html strict level 2//",                                             -- The public identifier starts with: "-//IETF//DTD HTML Strict Level 2//"
                    "-//ietf//dtd html strict level 3//",                                             -- The public identifier starts with: "-//IETF//DTD HTML Strict Level 3//"
                    "-//ietf//dtd html strict//",                                                     -- The public identifier starts with: "-//IETF//DTD HTML Strict//"
                    "-//ietf//dtd html//",                                                            -- The public identifier starts with: "-//IETF//DTD HTML//"
                    "-//metrius//dtd metrius presentational//",                                       -- The public identifier starts with: "-//Metrius//DTD Metrius Presentational//"
                    "-//microsoft//dtd internet explorer 2.0 html strict//",                          -- The public identifier starts with: "-//Microsoft//DTD Internet Explorer 2.0 HTML Strict//"
                    "-//microsoft//dtd internet explorer 2.0 html//",                                 -- The public identifier starts with: "-//Microsoft//DTD Internet Explorer 2.0 HTML//"
                    "-//microsoft//dtd internet explorer 2.0 tables//",                               -- The public identifier starts with: "-//Microsoft//DTD Internet Explorer 2.0 Tables//"
                    "-//microsoft//dtd internet explorer 3.0 html strict//",                          -- The public identifier starts with: "-//Microsoft//DTD Internet Explorer 3.0 HTML Strict//"
                    "-//microsoft//dtd internet explorer 3.0 html//",                                 -- The public identifier starts with: "-//Microsoft//DTD Internet Explorer 3.0 HTML//"
                    "-//microsoft//dtd internet explorer 3.0 tables//",                               -- The public identifier starts with: "-//Microsoft//DTD Internet Explorer 3.0 Tables//"
                    "-//netscape comm. corp.//dtd html//",                                            -- The public identifier starts with: "-//Netscape Comm. Corp.//DTD HTML//"
                    "-//netscape comm. corp.//dtd strict html//",                                     -- The public identifier starts with: "-//Netscape Comm. Corp.//DTD Strict HTML//"
                    "-//o'reilly and associates//dtd html 2.0//",                                     -- The public identifier starts with: "-//O'Reilly and Associates//DTD HTML 2.0//"
                    "-//o'reilly and associates//dtd html extended 1.0//",                            -- The public identifier starts with: "-//O'Reilly and Associates//DTD HTML Extended 1.0//"
                    "-//o'reilly and associates//dtd html extended relaxed 1.0//",                    -- The public identifier starts with: "-//O'Reilly and Associates//DTD HTML Extended Relaxed 1.0//"
                    "-//sq//dtd html 2.0 hotmetal + extensions//",                                    -- The public identifier starts with: "-//SQ//DTD HTML 2.0 HoTMetaL + extensions//"
                    "-//softquad software//dtd hotmetal pro 6.0::19990601::extensions to html 4.0//", -- The public identifier starts with: "-//SoftQuad Software//DTD HoTMetaL PRO 6.0::19990601::extensions to HTML 4.0//"
                    "-//softquad//dtd hotmetal pro 4.0::19971010::extensions to html 4.0//",          -- The public identifier starts with: "-//SoftQuad//DTD HoTMetaL PRO 4.0::19971010::extensions to HTML 4.0//"
                    "-//spyglass//dtd html 2.0 extended//",                                           -- The public identifier starts with: "-//Spyglass//DTD HTML 2.0 Extended//"
                    "-//sun microsystems corp.//dtd hotjava html//",                                  -- The public identifier starts with: "-//Sun Microsystems Corp.//DTD HotJava HTML//"
                    "-//sun microsystems corp.//dtd hotjava strict html//",                           -- The public identifier starts with: "-//Sun Microsystems Corp.//DTD HotJava Strict HTML//"
                    "-//w3c//dtd html 3 1995-03-24//",                                                -- The public identifier starts with: "-//W3C//DTD HTML 3 1995-03-24//"
                    "-//w3c//dtd html 3.2 draft//",                                                   -- The public identifier starts with: "-//W3C//DTD HTML 3.2 Draft//"
                    "-//w3c//dtd html 3.2 final//",                                                   -- The public identifier starts with: "-//W3C//DTD HTML 3.2 Final//"
                    "-//w3c//dtd html 3.2//",                                                         -- The public identifier starts with: "-//W3C//DTD HTML 3.2//"
                    "-//w3c//dtd html 3.2s draft//",                                                  -- The public identifier starts with: "-//W3C//DTD HTML 3.2S Draft//"
                    "-//w3c//dtd html 4.0 frameset//",                                                -- The public identifier starts with: "-//W3C//DTD HTML 4.0 Frameset//"
                    "-//w3c//dtd html 4.0 transitional//",                                            -- The public identifier starts with: "-//W3C//DTD HTML 4.0 Transitional//"
                    "-//w3c//dtd html experimental 19960712//",                                       -- The public identifier starts with: "-//W3C//DTD HTML Experimental 19960712//"
                    "-//w3c//dtd html experimental 970421//",                                         -- The public identifier starts with: "-//W3C//DTD HTML Experimental 970421//"
                    "-//w3c//dtd w3 html//",                                                          -- The public identifier starts with: "-//W3C//DTD W3 HTML//"
                    "-//w3o//dtd w3 html 3.0//",                                                      -- The public identifier starts with: "-//W3O//DTD W3 HTML 3.0//"
                    "-//webtechs//dtd mozilla html 2.0//",                                            -- The public identifier starts with: "-//WebTechs//DTD Mozilla HTML 2.0//"
                    "-//webtechs//dtd mozilla html//",                                                -- The public identifier starts with: "-//WebTechs//DTD Mozilla HTML//"
                }
                or tk.systemIdentifier == nil and doesPublicIDStartWith {
                    "-//w3c//dtd html 4.01 frameset//",     -- The system identifier is missing and the public identifier starts with: "-//W3C//DTD HTML 4.01 Frameset//"
                    "-//w3c//dtd html 4.01 transitional//", -- The system identifier is missing and the public identifier starts with: "-//W3C//DTD HTML 4.01 Transitional//"
                }
            )
        then
            -- then set the Document to quirks mode:
            p.document.mode = "quirks"
        elseif not p.document.isIFrameSrcdocDocument         -- Otherwise, if the document is not an iframe srcdoc document,
            and not p.document.parserCannotChangeTheModeFlag -- and the parser cannot change the mode flag is false,
            -- and the DOCTYPE token matches one of the conditions in the following list,
            and (
                doesPublicIDStartWith {
                    "-//w3c//dtd xhtml 1.0 frameset//",     -- The public identifier starts with: "-//W3C//DTD XHTML 1.0 Frameset//"
                    "-//w3c//dtd xhtml 1.0 transitional//", -- The public identifier starts with: "-//W3C//DTD XHTML 1.0 Transitional//"
                }
                or tk.systemIdentifier ~= nil and doesPublicIDStartWith {
                    "-//w3c//dtd html 4.01 frameset//",     -- The system identifier is not missing and the public identifier starts with: "-//W3C//DTD HTML 4.01 Frameset//"
                    "-//w3c//dtd html 4.01 transitional//", -- The system identifier is not missing and the public identifier starts with: "-//W3C//DTD HTML 4.01 Transitional//"
                }
            )
        then
            -- then set the Document to limited-quirks mode:
            p.document.mode = "limited-quirks"
        end
        -- Then, switch the insertion mode to "before html".
        p.currentInsertionMode = BeforeHTMLInsertionMode
        return
    end
    --> Anything else

    -- If the document is not an iframe srcdoc document,
    if not p.document.isIFrameSrcdocDocument then
        -- then this is a parse error;
        reportError(p, tk)
        -- if the parser cannot change the mode flag is false, set the Document to quirks mode.
        if not p.document.parserCannotChangeTheModeFlag then
            p.document.mode = "quirks"
        end
    end
end

---https://html.spec.whatwg.org/multipage/parsing.html#the-before-html-insertion-mode
BeforeHTMLInsertionMode = function(p, tk)
    -- When the user agent is to apply the rules for the "before html" insertion mode, the user agent must handle the token as follows:
    if tk.type == "doctype" then -- A DOCTYPE token
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> A comment token
    if tk.type == "comment" then
        -- Insert a comment as the last child of the Document object.
        error("todo")
        return
    end
    --> A character token that is one of U+0009 CHARACTER TABULATION, U+000A LINE FEED (LF), U+000C FORM FEED (FF), U+000D CARRIAGE RETURN (CR), or U+0020 SPACE
    if isCharacterTokenWith(tk, { "\t", "\n", "\f", "\r", " ", }) then
        -- Ignore the token.
        return
    end
    --> A start tag whose tag name is "html"
    if isStartTagWithName(tk, { "html" }) then
        -- Create an element for the token in the HTML namespace, with the Document as the intended parent.
        local element = createElementForToken(p, tk, namespaces.HTML_NAMESPACE, p.document)
        -- Append it to the Document object.
        p.document:appendChild(element)
        -- Put this element in the stack of open elements.
        p.stackOfOpenElements:push(element)
        -- Switch the insertion mode to "before head".
        p.currentInsertionMode = BeforeHeadInsertionMode
        return
    end
    --> An end tag whose tag name is one of: "head", "body", "html", "br"
    if isEndTagWithName(tk, { "head", "body", "html", "br" }) then
        -- Act as described in the "anything else" entry below.
    end
    --> Any other end tag
    if isEndTag(tk) then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end

    --> Anything else
    -- Create an html element whose node document is the Document object.
    local element = p.document:createElement("html", nil)
    -- Append it to the Document object. Put this element in the stack of open elements.
    p.document:appendChild(element)
    p.stackOfOpenElements:push(element)
    -- Switch the insertion mode to "before head", then reprocess the token.
    p.currentInsertionMode = BeforeHeadInsertionMode
    p.currentInsertionMode(p, tk)
end

---https://html.spec.whatwg.org/multipage/parsing.html#the-before-head-insertion-mode
BeforeHeadInsertionMode = function(p, tk)
    -- When the user agent is to apply the rules for the "before head" insertion mode, the user agent must handle the token as follows:

    --> A character token that is one of U+0009 CHARACTER TABULATION, U+000A LINE FEED (LF), U+000C FORM FEED (FF), U+000D CARRIAGE RETURN (CR), or U+0020 SPACE
    if isCharacterTokenWith(tk, { "\t", "\n", "\f", "\r", " ", }) then
        -- Ignore the token.
        return
    end
    --> A comment token
    if tk.type == "comment" then
        -- Insert a comment
        insertComment(p, tk)
        return
    end
    --> A DOCTYPE token
    if tk.type == "doctype" then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> A start tag whose tag name is "html"
    if isStartTagWithName(tk, { "html" }) then
        -- Process the token using the rules for the "in body" insertion mode.
        return InBodyInsertionMode(p, tk)
    end
    --> A start tag whose tag name is "head"
    if isStartTagWithName(tk, { "head" }) then
        -- Insert an HTML element for the token.
        local element = insertHTMLElement(p, tk)

        -- Set the head element pointer to the newly created head element.
        p.headElementPointer = element --[[@as HTML_HTMLHeadElement]]

        -- Switch the insertion mode to "in head".
        p.currentInsertionMode = InHeadInsertionMode
        return
    end
    --> An end tag whose tag name is one of: "head", "body", "html", "br"
    if isEndTagWithName(tk, { "head", "body", "html", "br" }) then
        -- Act as described in the "anything else" entry below.
    end
    --> Any other end tag
    if isEndTag(tk) then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> Anything else

    -- Insert an HTML element for a "head" start tag token with no attributes.
    local element = insertHTMLElement(p, token.TagToken:new("head", "start", tk.startLocation, tk.endLocation))

    -- Set the head element pointer to the newly created head element.
    p.headElementPointer = element --[[@as HTML_HTMLHeadElement]]

    -- Switch the insertion mode to "in head".
    p.currentInsertionMode = InHeadInsertionMode

    -- Reprocess the current token.
    p.currentInsertionMode(p, tk)
end

---https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhead
InHeadInsertionMode = function(p, tk)
    -- When the user agent is to apply the rules for the "in head" insertion mode, the user agent must handle the token as follows:

    --> A character token that is one of U+0009 CHARACTER TABULATION, U+000A LINE FEED (LF), U+000C FORM FEED (FF), U+000D CARRIAGE RETURN (CR), or U+0020 SPACE
    if isCharacterTokenWith(tk, { "\t", "\n", "\f", "\r", " " }) then
        -- Insert the character.
        insertCharacter(tk)
        return
    end
    --> A comment token
    if tk.type == "comment" then
        -- -> Insert a comment
        insertComment(p, tk)
        return
    end
    --> A DOCTYPE token
    if tk.type == "doctype" then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> A start tag whose tag name is "html"
    if isStartTagWithName(tk, { "html" }) then
        -- Process the token using the rules for the "in body" insertion mode.
        InBodyInsertionMode(p, tk)
        return
    end
    --> A start tag whose tag name is one of: "base", "basefont", "bgsound", "link"
    if isStartTagWithName(tk, { "base", "basefont", "bgsound", "link" }) then
        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)
        -- Immediately pop the current node off the stack of open elements.
        p.stackOfOpenElements:pop()
        return
    end
    --> A start tag whose tag name is "meta"
    if isStartTagWithName(tk, { "meta" }) then
        -- Insert an HTML element for the token.
        local element = insertHTMLElement(p, tk)
        -- Immediately pop the current node off the stack of open elements.
        p.stackOfOpenElements:pop()

        -- Acknowledge the token's self-closing flag, if it is set.
        tk --[[@as HTML_Parser_TagToken]]:acknowledgeSelfClosingTag()

        -- If the active speculative HTML parser is null, then:
        if not ACTIVE_SPECULATIVE_HTML_PARSER then
            -- 1. If the element has a charset attribute, and getting an encoding from its value results in an encoding, and the confidence is currently tentative, then change the encoding to the resulting encoding.
            if element:getAttribute("charset") then
                error("todo")
            end

            -- 2. Otherwise, if the element has an http-equiv attribute whose value is an ASCII case-insensitive match for the string "Content-Type", and the element has a content attribute, and applying the algorithm for extracting a character encoding from a meta element to that attribute's value returns an encoding, and the confidence is currently tentative, then change the encoding to the extracted encoding.
            if element:getAttribute("http-equiv") then
                error("todo")
            end
        end
        return
    end
    --> A start tag whose tag name is "title"
    if isStartTagWithName(tk, { "title" }) then
        -- Follow the generic RCDATA element parsing algorithm.
        genericRCDATAElementParsingAlgorithm(tk)
        return
    end
    --> A start tag whose tag name is "noscript", if the scripting flag is enabled
    --> A start tag whose tag name is one of: "noframes", "style"
    if (isStartTagWithName(tk, { "noscript" }) and SCRIPTING_ENABLED) or
        isStartTagWithName(tk, { "noframes", "style" })
    then
        -- Follow the generic raw text element parsing algorithm.
        genericRawTextParsingAlgorithm(tk)
        return
    end
    --> A start tag whose tag name is "noscript", if the scripting flag is disabled
    if (isStartTagWithName(tk, { "noscript" }) and not SCRIPTING_ENABLED) then
        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)
        -- Switch the insertion mode to "in head noscript".
        p.currentInsertionMode = InHeadNoScript
    end
    -- A start tag whose tag name is "script"
    if isStartTagWithName(tk, { "script" }) then
        -- Run these steps:

        -- 1. Let the adjusted insertion location be the appropriate place for inserting a node.
        local adjustedInsertionLocation = appropriatePlaceForInsertingNode(p)

        -- 2. Create an element for the token in the HTML namespace, with the intended parent being the element in which the adjusted insertion location finds itself.
        local element = createElementForToken(p, tk, namespaces.HTML_NAMESPACE, adjustedInsertionLocation
            .insertionParent) --[[@as HTML_HTMLScriptElement]]

        -- 3. Set the element's parser document to the Document, and set the element's force async to false.
        element.parserDocument = p.document
        element.forceAsync = false

        -- 4. If the parser was created as part of the HTML fragment parsing algorithm,
        if p.createdAsPartOfHTMLFragmentParsingAlgorithm then
            -- then set the script element's already started to true. (fragment case)
            element.alreadyStarted = true
        end

        -- 5. If the parser was invoked via the document.write() or document.writeln() methods,
        if p.invokedViaDocumentWrite then
            -- then optionally set the script element's already started to true. (For example, the user agent might use this clause to prevent execution of cross-origin scripts inserted via document.write() under slow network conditions, or when the page has already taken a long time to load.)
            element.alreadyStarted = true
        end

        -- 6. Insert the newly created element at the adjusted insertion location.
        element:preInsert(adjustedInsertionLocation.insertionParent, adjustedInsertionLocation.beforeChild)

        -- 7. Push the element onto the stack of open elements so that it is the new current node.
        p.stackOfOpenElements:push(element)

        -- 8. Switch the tokenizer to the script data state.
        p.tokenizer:switchToState(Tokenizer.ScriptDataState)

        -- 9. Let the original insertion mode be the current insertion mode.
        p.originalInsertionMode = p.currentInsertionMode

        -- 10. Switch the insertion mode to "text".
        p.currentInsertionMode = TextInsertionMode
        return
    end
    --> An end tag whose tag name is "head"
    if isEndTagWithName(tk, { "head" }) then
        -- Pop the current node (which will be the head element) off the stack of open elements.
        p.stackOfOpenElements:pop()
        -- Switch the insertion mode to "after head".
        p.currentInsertionMode = AfterHeadInsertionMode
        return
    end
    --> An end tag whose tag name is one of: "body", "html", "br"
    if isEndTagWithName(tk, { "body", "html", "br" }) then
        -- Act as described in the "anything else" entry below.
    end
    --> A start tag whose tag name is "template"
    if isStartTagWithName(tk, { "template" }) then
        -- Let template start tag be the start tag.
        local templateStartTag = tk --[[@as HTML_Parser_TagToken]]
        -- Insert a marker at the end of the list of active formatting elements.
        p.listOfActiveFormattingElements:insertMarkerAtEnd()
        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"
        -- Switch the insertion mode to "in template".
        p.currentInsertionMode = InTemplateInsertionMode
        -- Push "in template" onto the stack of template insertion modes so that it is the new current template insertion mode.
        table.insert(p.stackOfTemplateInsertionModes, InTemplateInsertionMode)
        -- Let the adjusted insertion location be the appropriate place for inserting a node.
        local adjustedInsertionLocation = appropriatePlaceForInsertingNode(p)
        -- Let intended parent be the element in which the adjusted insertion location finds itself.
        local intendedParent = adjustedInsertionLocation.insertionParent
        -- Let document be intended parent's node document.
        local document = intendedParent.nodeDocument
        -- If any of the following are false:
        if not (
            -- template start tag's shadowrootmode is not in the none state;
                templateStartTag:getAttribute("shadowrootmode") ~= "none"
                -- document's allow declarative shadow roots is true; or
                or document.allowDeclarativeShadowRoots
                -- the adjusted current node is not the topmost element in the stack of open elements,
                or adjustedCurrentNode(p) ~= p.stackOfOpenElements:currentNode()
            ) then
            -- then insert an HTML element for the token.
            insertHTMLElement(p, tk)
        else -- Otherwise:
            -- 1. Let declarative shadow host element be adjusted current node.
            local declarativeShadowHostElement = adjustedCurrentNode(p)

            -- 2. Let template be the result of insert a foreign element for template start tag, with HTML namespace and true.
            local template = insertForeignElement(p, templateStartTag, namespaces.HTML_NAMESPACE, true) --[[@as HTML_HTMLTemplateElement]]

            -- 3. Let mode be template start tag's shadowrootmode attribute's value.
            local mode = templateStartTag:getAttribute("shadowrootmode")

            -- 4. Let clonable be true if template start tag has a shadowrootclonable attribute; otherwise false.
            local clonable = templateStartTag:getAttribute("shadowrootclonable") ~= nil

            -- 5. Let serializable be true if template start tag has a shadowrootserializable attribute; otherwise false.
            local serializable = templateStartTag:getAttribute("shadowrootserializable") ~= nil

            -- 6. Let delegatesFocus be true if template start tag has a shadowrootdelegatesfocus attribute; otherwise false.
            local delegatesFocus = templateStartTag:getAttribute("shadowrootdelegatesfocus") ~= nil

            -- 7. If declarative shadow host element is a shadow host,
            if declarativeShadowHostElement:isShadowHost() then
                -- then insert an element at the adjusted insertion location with template.
                insertElementAtAdjustedInsertionLocation(p, template)
            else -- 8. Otherwise:
                -- Attach a shadow root with declarative shadow host element, mode, clonable, serializable, delegatesFocus, and "named".
                local ok, err = pcall(function()
                    declarativeShadowHostElement:attachShadowRoot(mode, clonable, serializable, delegatesFocus, "named")
                end)
                -- If an exception is thrown, then catch it and:
                if not ok then
                    -- 1. Insert an element at the adjusted insertion location with template.
                    insertElementAtAdjustedInsertionLocation(p, template)

                    -- 2. The user agent may report an error to the developer console.
                    L:e("Failed to attach shadow root: %s", tostring(err))

                    -- 3. Return.
                    return
                end

                -- 2. Let shadow be declarative shadow host element's shadow root.
                local shadow = declarativeShadowHostElement.shadowRoot --[[@as DOM_ShadowRoot]]

                -- 3. Set shadow's declarative to true.
                shadow.declarative = true

                -- 4. Set template's template contents property to shadow.
                template.templateContents = shadow

                -- 5. Set shadow's available to element internals to true.
                shadow.availableToElementInternals = true
            end
        end
        return
    end
    --> An end tag whose tag name is "template"
    if isEndTagWithName(tk, { "template" }) then
        -- If there is no template element on the stack of open elements,
        if p.stackOfTemplateInsertionModes:hasHTMLElement("template") then
            -- then this is a parse error;
            reportError(p, tk)
            -- ignore the token.
            return
        end
        -- Otherwise, run these steps:
        -- 1. Generate all implied end tags thoroughly.
        generateImpliedEndTagsThroughly()
        -- 2. If the current node is not a template element,
        if not p.stackOfOpenElements:currentNode():isHTMLElement("template") then
            -- then this is a parse error.
            reportError(p, tk)
        end
        -- 3. Pop elements from the stack of open elements until a template element has been popped from the stack.
        while not p.stackOfOpenElements:pop():isHTMLElement("template") do end
        -- 4. Clear the list of active formatting elements up to the last marker.
        clearListOfActiveFormattingElementsUpToLastMarker(p)
        -- 5. Pop the current template insertion mode off the stack of template insertion modes.
        table.remove(p.stackOfTemplateInsertionModes, #p.stackOfTemplateInsertionModes)
        -- 6. Reset the insertion mode appropriately.
        resetIntertionModeAppropriately(p)
        return
    end
    --> A start tag whose tag name is "head"
    --> Any other end tag
    if isStartTagWithName(tk, { "head" })
        or isEndTag(tk)
    then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> Anything else

    -- Pop the current node (which will be the head element) off the stack of open elements.
    p.stackOfOpenElements:pop()

    -- Switch the insertion mode to "after head".
    p.currentInsertionMode = AfterHeadInsertionMode

    -- Reprocess the token.
    return p.currentInsertionMode(p, tk)
end

---https://html.spec.whatwg.org/multipage/parsing.html#the-after-head-insertion-mode
AfterHeadInsertionMode = function(p, tk)
    -- When the user agent is to apply the rules for the "after head" insertion mode, the user agent must handle the token as follows:

    --> A character token that is one of U+0009 CHARACTER TABULATION, U+000A LINE FEED (LF), U+000C FORM FEED (FF), U+000D CARRIAGE RETURN (CR), or U+0020 SPACE
    if isCharacterTokenWith(tk, { "\t", "\n", "\f", "\r", " " }) then
        -- Insert the character.
        insertCharacter(p, tk)
        return
    end
    --> A comment token
    if tk.type == "comment" then
        -- Insert a comment.
        insertComment(p, tk)
        return
    end
    --> A DOCTYPE token
    if tk.type == "doctype" then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> A start tag whose tag name is "html"
    if isStartTagWithName(tk, { "html" }) then
        -- Process the token using the rules for the "in body" insertion mode.
        return InBodyInsertionMode(p, tk)
    end
    --> A start tag whose tag name is "body"
    if isStartTagWithName(tk, { "body" }) then
        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- Switch the insertion mode to "in body".
        p.currentInsertionMode = InBodyInsertionMode
        return
    end
    --> A start tag whose tag name is "frameset"
    if isStartTagWithName(tk, { "frameset" }) then
        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Switch the insertion mode to "in frameset".
        p.currentInsertionMode = InFramesetInsertionMode
        return
    end
    --> A start tag whose tag name is one of: "base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title"
    if isStartTagWithName(tk, {
            "base", "basefont", "bgsound", "link", "meta", "noframes", "script",
            "style", "template", "title"
        })
    then
        -- Parse error.
        reportError(p, tk)

        -- Push the node pointed to by the head element pointer onto the stack of open elements.
        p.stackOfOpenElements:push(p.headElementPointer)

        -- Process the token using the rules for the "in head" insertion mode.
        InHeadInsertionMode(p, tk)

        -- Remove the node pointed to by the head element pointer from the stack of open elements. (It might not be the current node at this point.)
        p.stackOfOpenElements:remove(p.headElementPointer)
        return
    end
    --> An end tag whose tag name is "template"
    if isEndTagWithName(tk, { "template" }) then
        -- Process the token using the rules for the "in head" insertion mode.
        return InHeadInsertionMode(p, tk)
    end
    --> An end tag whose tag name is one of: "body", "html", "br"
    if isEndTagWithName(tk, { "body", "html", "br" }) then
        -- Act as described in the "anything else" entry below.
    end
    --> A start tag whose tag name is "head"
    --> Any other end tag
    if isStartTagWithName(tk, { "head" })
        or isEndTag(tk)
    then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> Anything else

    -- Insert an HTML element for a "body" start tag token with no attributes.
    insertHTMLElement(p, token.TagToken:new("body", "start", tk.startLocation, tk.endLocation))

    -- Switch the insertion mode to "in body".
    p.currentInsertionMode = InBodyInsertionMode

    -- Reprocess the current token
    return p.currentInsertionMode(p, tk)
end

---https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inbody
InBodyInsertionMode = function(p, tk)
    -- When the user agent is to apply the rules for the "in body" insertion mode, the user agent must handle the token as follows:

    --> A character token that is U+0000 NULL
    if isCharacterTokenWith(tk, { "\0" }) then
        -- Parse error.
        reportError(p, tk)
        --Ignore the token.
        return
    end
    --> A character token that is one of U+0009 CHARACTER TABULATION, U+000A LINE FEED (LF), U+000C FORM FEED (FF), U+000D CARRIAGE RETURN (CR), or U+0020 SPACE
    if isCharacterTokenWith(tk, { "\t", "\n", "\f", "\r", " " }) then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)
        -- Insert the token's character.
        insertCharacter(tk)
        return
    end
    --> Any other character token
    if tk.type == "character" then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Insert the token's character.
        insertCharacter(tk)

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"
        return
    end
    --> A comment token
    if tk.type == "comment" then
        -- Insert a comment.
        insertComment(tk)
        return
    end
    --> A DOCTYPE token
    if tk.type == "doctype" then
        -- Parse error.
        reportError(p, tk)
        -- Ignore the token.
        return
    end
    --> A start tag whose tag name is "html"
    if isStartTagWithName(tk, { "html" }) then
        -- Parse error.
        reportError(p, tk)

        -- If there is a template element on the stack of open elements, then ignore the token.
        if p.stackOfOpenElements:hasHTMLElement("template") then
            return
        else
            -- Otherwise, for each attribute on the token, check to see if the attribute is already present on the top element of the stack of open elements. If it is not, add the attribute and its corresponding value to that element.
            error("todo")
        end
        return
    end
    --> A start tag whose tag name is one of: "base", "basefont", "bgsound", "link", "meta", "noframes", "script", "style", "template", "title"
    --> An end tag whose tag name is "template"
    if isStartTagWithName(tk, {
            "base", "basefont", "bgsound", "link", "meta", "noframes", "script",
            "style", "template", "title"
        }) or isEndTagWithName(tk, { "template" })
    then
        -- Process the token using the rules for the "in head" insertion mode.
        return InHeadInsertionMode(p, tk)
    end
    --> A start tag whose tag name is "body"
    if isStartTagWithName(tk, { "body" }) then
        -- Parse error.
        reportError(p, tk)

        -- If the stack of open elements has only one node on it,
        if #p.stackOfOpenElements.elements == 1
            -- if the second element on the stack of open elements is not a body element,
            or not p.stackOfOpenElements.elements[1]:isHTMLElement("body")
            -- or if there is a template element on the stack of open elements,
            or not p.stackOfOpenElements:hasHTMLElement("template")
        then
            -- then ignore the token. (fragment case or there is a template element on the stack)
            return
        else
            -- Otherwise, set the frameset-ok flag to "not ok";
            p.framesetOKFlag = "not ok"
            -- then, for each attribute on the token, check to see if the attribute is already present on the body element (the second element) on the stack of open elements, and if it is not, add the attribute and its corresponding value to that element.
            error("todo")
        end
        return
    end
    --> A start tag whose tag name is "frameset"
    if isStartTagWithName(tk, { "frameset" }) then
        -- Parse error.
        reportError(p, tk)

        -- If the stack of open elements has only one node on it,
        if #p.stackOfOpenElements.elements == 1
            -- or if the second element on the stack of open elements is not a body element,
            or not p.stackOfOpenElements.elements[1]:isHTMLElement("body")
        then
            -- then ignore the token. (fragment case or there is a template element on the stack)
            return
        end
        -- If the frameset-ok flag is set to "not ok",
        if p.framesetOKFlag == "not ok" then
            -- ignore the token.
            return
        end
        -- Otherwise, run the following steps:
        error("todo")
        -- 1. Remove the second element on the stack of open elements from its parent node, if it has one.
        -- 2. Pop all the nodes from the bottom of the stack of open elements, from the current node up to, but not including, the root html element.
        -- 3. Insert an HTML element for the token.
        insertHTMLElement(p, tk)
        -- 4. Switch the insertion mode to "in frameset".
        p.currentInsertionMode = InFramesetInsertionMode
        return
    end
    --> An end-of-file token
    if tk.type == "eof" then
        -- If the stack of template insertion modes is not empty,
        if #p.stackOfTemplateInsertionModes ~= 0 then
            -- then process the token using the rules for the "in template" insertion mode.
            return InTemplateInsertionMode(p, tk)
        end
        -- Otherwise, follow these steps:

        -- 1. If there is a node in the stack of open elements that is not either a dd element, a dt element, an li element, an optgroup element, an option element, a p element, an rb element, an rp element, an rt element, an rtc element, a tbody element, a td element, a tfoot element, a th element, a thead element, a tr element, the body element, or the html element,
        for _, node in ipairs(p.stackOfOpenElements.elements) do
            if not node:isOneOfHTMLElements {
                    "dd", "dt", "li", "optgroup", "option", "p", "rb", "rp",
                    "rt", "rtc", "tbody", "td", "tfoot", "th", "thead", "tr",
                    "body", "html"
                }
            then
                -- then this is a parse error.
                reportError(p, tk)
            end
        end
        error("todo")
        -- 2. Stop parsing.
        stopParsing(p)
        return
    end
    --> An end tag whose tag name is "body"
    if isEndTagWithName(tk, { "body" }) then
        -- If the stack of open elements does not have a body element in scope,
        if not p.stackOfOpenElements:hasElementInScope { "body" } then
            -- this is a parse error;
            reportError(tk)
            -- ignore the token.
            return
        end

        -- Otherwise, if there is a node in the stack of open elements that is not either a dd element, a dt element, an li element, an optgroup element, an option element, a p element, an rb element, an rp element, an rt element, an rtc element, a tbody element, a td element, a tfoot element, a th element, a thead element, a tr element, the body element, or the html element,
        for _, node in ipairs(p.stackOfOpenElements.elements) do
            if not node:isOneOfHTMLElements {
                    "dd", "dt", "li", "optgroup", "option", "a", "p", "rb",
                    "rp", "rt", "rtc", "tbody", "td", "tfoot", "th", "thead",
                    "tr", "body", "html"
                }
            then
                -- then this is a parse error.
                reportError(p, tk)
            end
        end
        -- Switch the insertion mode to "after body".
        p.currentInsertionMode = AfterBodyInsertionMode
        return
    end
    --> An end tag whose tag name is "html"
    if isEndTagWithName(tk, { "html" }) then
        -- If the stack of open elements does not have a body element in scope,
        if not p.stackOfOpenElements:hasElementInScope { "body" } then
            -- this is a parse error;
            reportError(tk)
            -- ignore the token.
            return
        end

        -- Otherwise, If there is a node in the stack of open elements that is not either a dd element, a dt element, an li element, an optgroup element, an option element, a p element, an rb element, an rp element, an rt element, an rtc element, a tbody element, a td element, a tfoot element, a th element, a thead element, a tr element, the body element, or the html element,
        for _, node in ipairs(p.stackOfOpenElements.elements) do
            if not node:isOneOfHTMLElements {
                    "dd", "dt", "li", "optgroup", "option", "a", "p", "rb",
                    "rp", "rt", "rtc", "tbody", "td", "tfoot", "th", "thead",
                    "tr", "body", "html"
                } then
                -- then this is a parse error.
                reportError(p, tk)
            end
        end
        -- Switch the insertion mode to "after body".
        p.currentInsertionMode = AfterBodyInsertionMode

        -- Reprocess the token.
        return p.currentInsertionMode(p, tk)
    end
    --> A start tag whose tag name is one of: "address", "article", "aside", "blockquote", "center", "details", "dialog", "dir", "div", "dl", "fieldset", "figcaption", "figure", "footer", "header", "hgroup", "main", "menu", "nav", "ol", "p", "search", "section", "summary", "ul"
    if isStartTagWithName(tk, {
            "address", "article", "aside", "blockquote", "center", "details",
            "dialog", "dir", "div", "dl", "fieldset", "figcaption", "figure",
            "footer", "header", "hgroup", "main", "menu", "nav", "ol", "p",
            "search", "section", "summary", "ul"
        })
    then
        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end
        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)
        return
    end
    --> A start tag whose tag name is one of: "h1", "h2", "h3", "h4", "h5", "h6"
    if isStartTagWithName(tk, { "h1", "h2", "h3", "h4", "h5", "h6" }) then
        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- If the current node is an HTML element whose tag name is one of "h1", "h2", "h3", "h4", "h5", or "h6",
        if p.stackOfOpenElements:currentNode():isOneOfHTMLElements {
                "h1", "h2", "h3", "h4", "h5", "h6"
            }
        then
            -- then this is a parse error;
            reportError(p, tk)
            -- pop the current node off the stack of open elements.
            p.stackOfOpenElements:pop()
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)
        return
    end
    --> A start tag whose tag name is one of: "pre", "listing"
    if isStartTagWithName(tk, { "pre", "listing" }) then
        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- If the next token is a U+000A LINE FEED (LF) character token, then ignore that token and move on to the next one. (Newlines at the start of pre blocks are ignored as an authoring convenience.)
        p.ignoreNextLF = true

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        return
    end
    --> A start tag whose tag name is "form"
    if isStartTagWithName(tk, { "form" }) then
        -- If the form element pointer is not null, and there is no template element on the stack of open elements,
        if p.formElementPointer ~= nil and not p.stackOfOpenElements:hasHTMLElement("template") then
            -- then this is a parse error;
            reportError(p, tk)
            -- ignore the token.
            return
        end
        -- Otherwise:

        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- Insert an HTML element for the token,
        local element = insertHTMLElement(p, tk)

        -- and, if there is no template element on the stack of open elements,
        if not p.stackOfOpenElements:hasHTMLElement("template") then
            -- set the form element pointer to point to the element created.
            p.formElementPointer = element
        end
        return
    end
    --> A start tag whose tag name is "li"
    if isStartTagWithName(tk, { "li" }) then
        -- Run these steps:

        -- 1. Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- 2. Initialize node to be the current node (the bottommost node of the stack).
        local nodeIndex = #p.stackOfOpenElements.elements

        while true do
            local node = p.stackOfOpenElements.elements[nodeIndex]

            -- 3. Loop: If node is an li element, then run these substeps:
            if node:isHTMLElement("li") then
                -- 1. Generate implied end tags, except for li elements.
                generateImpliedEndTags(p, { "li" })

                -- 2. If the current node is not an li element,
                if not p.stackOfOpenElements:currentNode():isHTMLElement("li") then
                    -- then this is a parse error.
                    reportError(p, tk)
                end

                -- 3. Pop elements from the stack of open elements until an li element has been popped from the stack.
                while not p.stackOfOpenElements:pop():isHTMLElement("li") do end

                -- 4. Jump to the step labeled done below.
                break
            end
            -- 4. If node is in the special category, but is not an address, div, or p element,
            if node:isInSpecialCategory()
                and not node:isOneOfHTMLElements { "address", "div", "p" }
            then
                -- then jump to the step labeled done below.
                break
            end
            -- 5. Otherwise, set node to the previous entry in the stack of open elements and return to the step labeled loop.
            nodeIndex = nodeIndex - 1
        end
        -- 6. Done: If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- 7. Finally, insert an HTML element for the token.
        insertHTMLElement(p, tk)
        return
    end
    --> A start tag whose tag name is one of: "dd", "dt"
    if isStartTagWithName(tk, { "dd", "dt" }) then
        -- Run these steps:

        -- 1. Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- 2. Initialize node to be the current node (the bottommost node of the stack).
        local nodeIndex = #p.stackOfOpenElements.elements

        while true do
            local node = p.stackOfOpenElements.elements[nodeIndex]

            -- 3. Loop: If node is a dd element, then run these substeps:
            if node:isHTMLElement("dd") then
                -- 1. Generate implied end tags, except for dd elements.
                generateImpliedEndTags(p, { "dd" })

                -- 2. If the current node is not an dd element,
                if not p.stackOfOpenElements:currentNode():isHTMLElement("li") then
                    -- then this is a parse error.
                    reportError(p, tk)
                end

                -- 3. Pop elements from the stack of open elements until an dd element has been popped from the stack.
                while not p.stackOfOpenElements:pop():isHTMLElement("li") do end

                -- 4. Jump to the step labeled done below.
                break
            end
            -- 3. If node is a dt element, then run these substeps:
            if node:isHTMLElement("dt") then
                -- 1. Generate implied end tags, except for dt elements.
                generateImpliedEndTags(p, { "dt" })

                -- 2. If the current node is not an dt element,
                if not p.stackOfOpenElements:currentNode():isHTMLElement("li") then
                    -- then this is a parse error.
                    reportError(p, tk)
                end

                -- 3. Pop elements from the stack of open elements until an dt element has been popped from the stack.
                while not p.stackOfOpenElements:pop():isHTMLElement("li") do end

                -- 4. Jump to the step labeled done below.
                break
            end

            -- 5. If node is in the special category, but is not an address, div, or p element,
            if node:isInSpecialCategory()
                and not node:isOneOfHTMLElements { "address", "div", "p" }
            then
                -- then jump to the step labeled done below.
                break
            end
            -- 6. Otherwise, set node to the previous entry in the stack of open elements and return to the step labeled loop.
            nodeIndex = nodeIndex - 1
        end
        -- 7. Done: If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- 8. Finally, insert an HTML element for the token.
        insertHTMLElement(p, tk)
        return
    end
    --> A start tag whose tag name is "plaintext"
    if isStartTagWithName(tk, { "plaintext" }) then
        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Switch the tokenizer to the PLAINTEXT state.
        p.tokenizer:switchToState(p.PlainText)
    end
    --> A start tag whose tag name is "button"
    if isStartTagWithName(tk, { "button" }) then
        -- 1. If the stack of open elements has a button element in scope, then run these substeps:
        if p.stackOfOpenElements:hasElementInScope { "button" } then
            -- 1. Parse error.
            reportError(p, tk)

            -- 2. Generate implied end tags.
            generateImpliedEndTags(p)

            -- 3. Pop elements from the stack of open elements until a button element has been popped from the stack.
            while not p.stackOfOpenElements:pop():isHTMLElement("button") do end
        end

        -- 2. Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- 3. Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- 4. Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"
        return
    end
    --> An end tag whose tag name is one of: "address", "article", "aside", "blockquote", "button", "center", "details", "dialog", "dir", "div", "dl", "fieldset", "figcaption", "figure", "footer", "header", "hgroup", "listing", "main", "menu", "nav", "ol", "pre", "search", "section", "summary", "ul"
    if isEndTagWithName(tk, {
            "article", "aside", "blockquote", "button", "center", "details",
            "dialog", "dir", "div", "dl", "fieldset", "figcaption", "figure",
            "footer", "header", "hgroup", "listing", "main", "menu", "nav",
            "ol", "pre", "search", "section", "summary", "ul"
        })
    then
        -- If the stack of open elements does not have an element in scope that is an HTML element with the same tag name as that of the token,
        if p.stackOfOpenElements:hasElementInScope { tk.name } then
            -- then this is a parse error;
            reportError(p, tk)
            -- ignore the token.
            return
        end
        -- Otherwise, run these steps:

        -- 1. Generate implied end tags.
        generateImpliedEndTags(p)

        -- 2. If the current node is not an HTML element with the same tag name as that of the token,
        if not p.stackOfOpenElements:currentNode():isHTMLElement(tk.name) then
            -- then this is a parse error.
            reportError(p, tk)
        end

        -- 3. Pop elements from the stack of open elements until an HTML element with the same tag name as the token has been popped from the stack.
        while not p.stackOfOpenElements:pop():isHTMLElement(tk.name) do end

        return
    end

    --> An end tag whose tag name is "form"
    if isEndTagWithName(tk, { "form" }) then
        -- If there is no template element on the stack of open elements,
        if not p.stackOfOpenElements:hasHTMLElement("template") then
            -- then run these substeps:

            -- 1. Let node be the element that the form element pointer is set to, or null if it is not set to an element.
            local node = p.formElementPointer

            -- 2. Set the form element pointer to null.
            p.formElementPointer = nil

            -- 3. If node is null or if the stack of open elements does not have node in scope,
            if node == nil or p.stackOfOpenElements:hasElementInScope { node } then
                -- then this is a parse error;
                reportError(p, tk)
                -- return and ignore the token.
                return
            end

            -- 4. Generate implied end tags.
            generateImpliedEndTags(p)

            -- 5. If the current node is not node,
            if p.stackOfOpenElements:currentNode() ~= node then
                -- then this is a parse error.
                reportError(p, tk)
            end

            -- 6. Remove node from the stack of open elements.
            p.stackOfOpenElements:remove(node)
        else
            -- If there is a template element on the stack of open elements, then run these substeps instead:

            -- 1. If the stack of open elements does not have a form element in scope,
            if not p.stackOfOpenElements:hasElementInScope { "form" } then
                -- then this is a parse error;
                reportError(p, tk)
                -- return and ignore the token.
                return
            end

            -- 2. Generate implied end tags.
            generateImpliedEndTags(p)

            -- 3. If the current node is not a form element,
            if not p.stackOfOpenElements:currentNode():isHTMLElement("form") then
                -- then this is a parse error.
                reportError(p, tk)
            end

            -- 4. Pop elements from the stack of open elements until a form element has been popped from the stack.
            while not p.stackOfOpenElements:pop():isHTMLElement("form") do end
        end
        return
    end
    --> An end tag whose tag name is "p"
    if isEndTagWithName(tk, { "p" }) then
        -- If the stack of open elements does not have a p element in button scope,
        if not p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then this is a parse error;
            reportError(tk)
            -- insert an HTML element for a "p" start tag token with no attributes.
            insertHTMLElement(p, token.TagToken:new("p", "start", tk.startLocation, tk.endLocation))
            return
        end
        -- Close a p element.
        closePElement(p)
        return
    end
    --> An end tag whose tag name is "li"
    if isEndTagWithName(tk, { "li" }) then
        -- If the stack of open elements does not have an li element in list item scope,
        if not p.stackOfOpenElements:hasElementInListItemScope { "li" } then
            -- then this is a parse error;
            reportError(tk)
            -- ignore the token.
            return
        end

        -- 1. Generate implied end tags, except for li elements.
        generateImpliedEndTags(p, { "li" })

        -- 2. If the current node is not an li element,
        if not p.stackOfOpenElements:currentNode():isHTMLElement("li") then
            -- then this is a parse error.
            reportError(p, tk)
        end

        -- 3. Pop elements from the stack of open elements until an li element has been popped from the stack.
        while not p.stackOfOpenElements:pop():isHTMLElement("li") do end

        return
    end
    --> An end tag whose tag name is "dd", "dt"
    if isEndTagWithName(tk, { "dd", "dt" }) then
        -- If the stack of open elements does not have an element in scope that is an HTML element with the same tag name as that of the token,
        if not p.stackOfOpenElements:hasElementInScope { tk.name } then
            -- then this is a parse error;
            reportError(tk)
            -- ignore the token.
            return
        end
        -- Otherwise, run these steps:

        -- 1. Generate implied end tags, except for HTML elements with the same tag name as the token.
        generateImpliedEndTags(p, { tk.name })

        -- 2. If the current node is not an HTML element with the same tag name as that of the token,
        if not p.stackOfOpenElements:currentNode():isHTMLElement(tk.name) then
            -- then this is a parse error.
            reportError(p, tk)
        end

        -- 3. Pop elements from the stack of open elements until an HTML element with the same tag name as the token has been popped from the stack.
        while not p.stackOfOpenElements:pop():isHTMLElement(tk.name) do end

        return
    end
    --> An end tag whose tag name is one of: "h1", "h2", "h3", "h4", "h5", "h6"
    if isEndTagWithName(tk, { "h1", "h2", "h3", "h4", "h5", "h6" }) then
        -- If the stack of open elements does not have an element in scope that is an HTML element and whose tag name is one of "h1", "h2", "h3", "h4", "h5", or "h6",
        if not p.stackOfOpenElements:hasElementInScope { "h1", "h2", "h3", "h4", "h5", "h6" } then
            -- then this is a parse error;
            reportError(tk)
            -- ignore the token.
            return
        end
        -- Otherwise, run these steps:

        -- 1. Generate implied end tags
        generateImpliedEndTags(p)

        -- 2. If the current node is not an HTML element with the same tag name as that of the token,
        if not p.stackOfOpenElements:currentNode():isHTMLElement(tk.name) then
            -- then this is a parse error.
            reportError(p, tk)
        end

        -- 3. Pop elements from the stack of open elements until an HTML element whose tag name is one of "h1", "h2", "h3", "h4", "h5", or "h6" has been popped from the stack.
        while not p.stackOfOpenElements:pop():isOneOfHTMLElements { "h1", "h2", "h3", "h4", "h5", "h6" } do end

        return
    end
    --> An end tag whose tag name is "sarcasm"
    if isEndTagWithName(tk, { "sarcasm" }) then
        -- Take a deep breath, then act as described in the "any other end tag" entry below.
        local messages = {
            "It's time to take a deep breath",
            "심호흡을 할 시간입니다",
            "深呼吸をする時間です",
        };
        L:w(messages[math.random(3)] or messages[1])
    end
    --> A start tag whose tag name is "a"
    if isStartTagWithName(tk, { "a" }) then
        -- If the list of active formatting elements contains an a element between the end of the list and the last marker on the list (or the start of the list if there is no marker on the list),
        if p.listOfActiveFormattingElements:containsElementSinceLastMarker { "a" }
            or p.listOfActiveFormattingElements.elements[1]:isHTMLElement("a")
        then
            -- then this is a parse error;
            reportError(p, tk)
            -- run the adoption agency algorithm for the token,
            adoptionAgencyAlgorithm(p, tk)
            -- then remove that element from the list of active formatting elements and the stack of open elements if the adoption agency algorithm didn't already remove it (it might not have if the element is not in table scope).
            p.listOfActiveFormattingElements:remove(element)
        end
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)
        -- Insert an HTML element for the token.
        local element = insertHTMLElement(p, tk)
        -- Push onto the list of active formatting elements that element.
        p.listOfActiveFormattingElements:push(element, tk)
        return
    end
    --> A start tag whose tag name is one of: "b", "big", "code", "em", "font", "i", "s", "small", "strike", "strong", "tt", "u"
    if isStartTagWithName(tk, {
            "b", "big", "code", "em", "font", "i", "s", "small", "strike",
            "strong", "tt", "u"
        })
    then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)
        -- Insert an HTML element for the token.
        local element = insertHTMLElement(p, tk)
        -- Push onto the list of active formatting elements that element.
        p.listOfActiveFormattingElements:push(element, tk)
        return
    end
    --> A start tag whose tag name is "nobr"
    if isStartTagWithName(tk, { "nobr" }) then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- If the stack of open elements has a nobr element in scope,
        if p.stackOfOpenElements.hasElementInScope("nobr") then
            -- then this is a parse error;
            reportError(p, tk)
            -- run the adoption agency algorithm for the token,
            adoptionAgencyAlgorithm(p, tk)
            -- then once again reconstruct the active formatting elements, if any.
            reconstructActiveFormattingElements(p)
        end
        -- Insert an HTML element for the token.
        local element = insertHTMLElement(p, tk)
        -- Push onto the list of active formatting elements that element.
        p.listOfActiveFormattingElements:push(element, tk)
        return
    end
    --> An end tag whose tag name is one of: "a", "b", "big", "code", "em", "font", "i", "nobr", "s", "small", "strike", "strong", "tt", "u"
    if isStartTagWithName(tk, {
            "a", "b", "big", "code", "em", "font", "i", "nobr", "s", "small",
            "strike", "strong", "tt", "u"
        })
    then
        -- Run the adoption agency algorithm for the token.
        adoptionAgencyAlgorithm(p, tk)
        return
    end
    --> A start tag whose tag name is one of: "applet", "marquee", "object"
    if isStartTagWithName(tk, { "applet", "marquee", "object" }) then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Insert a marker at the end of the list of active formatting elements.
        table.insert(p.listOfActiveFormattingElements.elements, { type = "marker" })

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"
        return
    end
    --> An end tag token whose tag name is one of: "applet", "marquee", "object"
    if isEndTagWithName(tk, { "applet", "marquee", "object" }) then
        -- If the stack of open elements does not have an element in scope that is an HTML element with the same tag name as that of the token,
        if p.stackOfOpenElements:hasElementInScope { tk.name } then
            -- then this is a parse error;
            reportError(p, tk)
            -- ignore the token.
            return
        end
        -- Otherwise, run these steps:

        -- 1. Generate implied end tags.
        generateImpliedEndTags(p)

        -- 2. If the current node is not an HTML element with the same tag name as that of the token,
        if not p.stackOfOpenElements:currentNode():isHTMLElement(tk.name) then
            -- then this is a parse error.
            reportError(p, tk)
        end

        -- 3. Pop elements from the stack of open elements until an HTML element with the same tag name as the token has been popped from the stack.
        while not p.stackOfOpenElements:pop():isHTMLElement(tk.name) do end

        -- 4. Clear the list of active formatting elements up to the last marker.
        clearListOfActiveFormattingElementsUpToLastMarker()
        return
    end
    --> A start tag whose tag name is "table"
    if isStartTagWithName(tk, { "table" }) then
        -- If the Document is not set to quirks mode, and the stack of open elements has a p element in button scope,
        if p.document.mode ~= "quirks"
            and p.stackOfOpenElements:hasElementInButtonScope { "p" }
        then
            -- then close a p element.
            closePElement(p)
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- Switch the insertion mode to "in table".
        p.currentInsertionMode = InTableInsertionMode
        return
    end
    --> An end tag whose tag name is "br"
    --> A start tag whose tag name is one of: "area", "br", "embed", "img", "keygen", "wbr"
    local isEndBr = isEndTagWithName(tk, { "br" })
    if isEndBr or isStartTagWithName(tk, {
            "area", "br", "embed", "img", "keygen", "wbr"
        })
    then
        --> An end tag whose tag name is "br"
        if isEndBr then
            -- Parse error.
            reportError(p, tk)
            -- Drop the attributes from the token,
            tk.attributes = {}
            -- and act as described in the next entry; i.e. act as if this was a "br" start tag token with no attributes, rather than the end tag token that it actually is.
        end
        --> A start tag whose tag name is one of: "area", "br", "embed", "img", "keygen", "wbr"

        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Immediately pop the current node off the stack of open elements.
        p.stackOfOpenElements:pop()

        -- Acknowledge the token's self-closing flag, if it is set.
        if tk.selfClosing then tk.selfClosingAcknowledged = true end

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"
        return
    end
    --> A start tag whose tag name is "input"
    if isStartTagWithName(tk, { "input" }) then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Immediately pop the current node off the stack of open elements.
        p.stackOfOpenElements:pop()

        -- Acknowledge the token's self-closing flag, if it is set.
        if tk.selfClosing then tk.selfClosingAcknowledged = true end

        -- If the token does not have an attribute with the name "type", or if it does, but that attribute's value is not an ASCII case-insensitive match for the string "hidden", then:
        local nameAttr = tk --[[@as HTML_Parser_TagToken]]:getAttribute("name")
        if nameAttr == nil or nameAttr:lower() ~= "hidden" then
            -- set the frameset-ok flag to "not ok".
            p.framesetOKFlag = "not ok"
        end
        return
    end
    --> A start tag whose tag name is one of: "param", "source", "track"
    if isStartTagWithName(tk, { "param", "source", "track" }) then
        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Immediately pop the current node off the stack of open elements.
        p.stackOfOpenElements:pop()

        -- Acknowledge the token's self-closing flag, if it is set.
        if tk.selfClosing then tk.selfClosingAcknowledged = true end

        return
    end
    --> A start tag whose tag name is "hr"
    if isStartTagWithName(tk, { "hr" }) then
        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Immediately pop the current node off the stack of open elements.
        p.stackOfOpenElements:pop()

        -- Acknowledge the token's self-closing flag, if it is set.
        if tk.selfClosing then tk.selfClosingAcknowledged = true end

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        return
    end

    --> A start tag whose tag name is "image"
    if isStartTagWithName(tk, { "image" }) then
        -- Parse error.
        reportError(p, tk)

        -- Change the token's tag name to "img" and reprocess it. (Don't ask.)
        tk.name = "img"
        return InBodyInsertionMode(p, tk)
    end

    --> A start tag whose tag name is "textarea"
    if isStartTagWithName(tk, { "textarea" }) then
        -- Run these steps:

        -- 1. Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- 2. If the next token is a U+000A LINE FEED (LF) character token, then ignore that token and move on to the next one. (Newlines at the start of textarea elements are ignored as an authoring convenience.)
        p.ignoreNextLF = true

        -- 3. Switch the tokenizer to the RCDATA state.
        p.tokenizer:switchToState(Tokenizer.RCDATAState)

        -- 4. Let the original insertion mode be the current insertion mode.
        p.originalInsertionMode = p.currentInsertionMode

        -- 5. Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- 6. Switch the insertion mode to "text".
        p.currentInsertionMode = TextInsertionMode
        return
    end

    --> A start tag whose tag name is "xmp"
    if isStartTagWithName(tk, { "xmp" }) then
        -- If the stack of open elements has a p element in button scope,
        if p.stackOfOpenElements:hasElementInButtonScope { "p" } then
            -- then close a p element.
            closePElement(p)
        end

        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- Follow the generic raw text element parsing algorithm.
        genericRawTextParsingAlgorithm(p)
        return
    end

    --> A start tag whose tag name is "iframe"
    if isStartTagWithName(tk, { "iframe" }) then
        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- Follow the generic raw text element parsing algorithm.
        genericRawTextParsingAlgorithm(p)
        return
    end

    --> A start tag whose tag name is "noembed"
    --> A start tag whose tag name is "noscript", if the scripting flag is enabled
    if isStartTagWithName(tk, { "noembed" })
        or (SCRIPTING_ENABLED and isStartTagWithName(tk, { "noscript" }))
    then
        -- Follow the generic raw text element parsing algorithm.

        genericRawTextParsingAlgorithm(p)
        return
    end

    --> A start tag whose tag name is "select"
    if isStartTagWithName(tk, { "select" }) then
        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        -- Set the frameset-ok flag to "not ok".
        p.framesetOKFlag = "not ok"

        -- If the insertion mode is one of "in table", "in caption", "in table body", "in row", or "in cell",
        if p.currentInsertionMode == InTableInsertionMode
            or p.currentInsertionMode == InCaptionInsertionMode
            or p.currentInsertionMode == InTableBodyInsertionMode
            or p.currentInsertionMode == InRowInsertionMode
            or p.currentInsertionMode == InCellInsertionMode
        then
            -- then switch the insertion mode to "in select in table".
            p.currentInsertionMode = InSelectInTableInsertionMode
        else
            -- Otherwise, switch the insertion mode to "in select".
            p.currentInsertionMode = InSelectInsertionMode
        end

        return
    end

    --> A start tag whose tag name is one of: "optgroup", "option"
    if isStartTagWithName(tk, { "optgroup", "option" }) then
        -- If the current node is an option element,
        if p.stackOfOpenElements:currentNode():isHTMLElement("option") then
            -- then pop the current node off the stack of open elements.
            p.stackOfOpenElements:pop()
        end

        -- Reconstruct the active formatting elements, if any.
        reconstructActiveFormattingElements(p)

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        return
    end
    --> A start tag whose tag name is one of: "rb", "rtc"
    if isStartTagWithName(tk, { "rb", "rtc" }) then
        -- If the stack of open elements has a ruby element in scope,
        if p.stackOfOpenElements:hasElementInScope({ "ruby" }) then
            -- then generate implied end tags.
            generateImpliedEndTags(p)

            -- If the current node is not now a ruby element,
            if not p.stackOfOpenElements:currentNode():isHTMLElement("ruby") then
                -- this is a parse error.
                reportError(p, tk)
            end
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        return
    end
    --> A start tag whose tag name is one of: "rp", "rt"
    if isStartTagWithName(tk, { "rp", "rt" }) then
        -- If the stack of open elements has a ruby element in scope,
        if p.stackOfOpenElements:hasElementInScope({ "ruby" }) then
            -- then generate implied end tags, except for rtc elements.
            generateImpliedEndTags(p, { "rtc" })

            -- If the current node is not now a rtc element or a ruby element
            if not p.stackOfOpenElements:currentNode():isOneOfHTMLElements({ "rtc", "ruby" }) then
                -- this is a parse error.
                reportError(p, tk)
            end
        end

        -- Insert an HTML element for the token.
        insertHTMLElement(p, tk)

        return
    end
    --> A start tag whose tag name is "math"
    if isStartTagWithName(tk, { "math" }) then
        
    end


    error("unreachable")
end
