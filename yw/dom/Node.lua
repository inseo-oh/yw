--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local iterutil     = require "yw.common.iterutil"
local Range        = require "yw.dom.Range"
local NodeIterator = require "yw.dom.NodeIterator"
local object       = require "yw.common.object"


---https://dom.spec.whatwg.org/#concept-node
---@class DOM_Node
---@field parentNode              DOM_Node|nil      https://dom.spec.whatwg.org/#concept-tree-parent
---@field firstChild              DOM_Node|nil      https://dom.spec.whatwg.org/#concept-tree-first-child
---@field lastChild               DOM_Node|nil      https://dom.spec.whatwg.org/#concept-tree-last-child
---@field nextSibling             DOM_Node|nil      https://dom.spec.whatwg.org/#concept-tree-next-sibling
---@field previousSibling         DOM_Node|nil      https://dom.spec.whatwg.org/#concept-tree-previous-sibling
---@field nodeDocument            DOM_Document      https://dom.spec.whatwg.org/#concept-node-document
---@field mutationObservers       any[]
---
---@field isText                  boolean
---@field isCDATASection          boolean
---@field isAttr                  boolean
---@field isElement               boolean
---@field isCharacterData         boolean
---@field isProcessingInstruction boolean
---@field isComment               boolean
---@field isDocument              boolean
---@field isDocumentType          boolean
---@field isDocumentFragment      boolean
---@field isShadowRoot            boolean
---@field isSlottable             boolean
local Node = {}

---@param document DOM_Document?
---@return DOM_Node
function Node:new(document)
    local o = object.create(self)


    o.parentNode              = nil
    o.firstChild              = nil
    o.lastChild               = nil
    o.nextSibling             = nil
    o.previousSibling         = nil
    o.nodeDocument            = document
    o.mutationObservers       = {}
    o.isText                  = false
    o.isCDATASection          = false
    o.isAttr                  = false
    o.isElement               = false
    o.isCharacterData         = false
    o.isProcessingInstruction = false
    o.isComment               = false
    o.isDocument              = false
    o.isDocumentType          = false
    o.isDocumentFragment      = false
    o.isShadowRoot            = false
    o.isSlottable             = false
    o.mutationObservers       = {} -- STUB
    return o
end

---@return DOM_Element
function Node:asElement()
    assert(self.isElement)
    return self --[[@as DOM_Element]]
end

---@return DOM_Attr
function Node:asAttr()
    assert(self.isAttr)
    return self --[[@as DOM_Attr]]
end

---@return DOM_DocumentFragment
function Node:asDocumentFragment()
    assert(self.isDocumentFragment)
    return self --[[@as DOM_DocumentFragment]]
end

---@return DOM_ShadowRoot
function Node:asShadowRoot()
    assert(self.isShadowRoot)
    return self --[[@as DOM_ShadowRoot]]
end

---@return DOM_CharacterData
function Node:asCharacterData()
    assert(self.isCharacterData)
    return self --[[@as DOM_CharacterData]]
end

function Node:asDocumentType()
    assert(self.isDocumentType)
    return self --[[@as DOM_DocumentType]]
end

---@param document DOM_Document?
function Node:runAdoptingSteps(document)
end

function Node:runChildrenChangedSteps()
end

function Node:runInsertionSteps()
end

---@param parent DOM_Node?
function Node:runRemovingSteps(parent)
end

function Node:isExclusiveTextNode()
    return self.isText and not self.isCDATASection
end

---@param parent DOM_Node
---@param childNode DOM_Node
---@param before DOM_Node | nil
local function insert(parent, childNode, before)
    -- Set parent
    childNode.parentNode = parent

    -- Set previous/next sibling
    childNode.nextSibling = before
    if before ~= nil then
        childNode.previousSibling = before.previousSibling
        before.previousSibling = childNode
    else
        childNode.previousSibling = nil
    end
    if childNode.previousSibling ~= nil then
        childNode.previousSibling.nextSibling = childNode
    end

    -- Set first child if needed
    if childNode.previousSibling == nil then
        parent.firstChild = childNode
    end
end

---https://dom.spec.whatwg.org/#concept-tree-preceding
---@return DOM_Node|nil
function Node:precedingNode()
    if self.previousSibling ~= nil then
        local current = self.previousSibling
        while (current ~= nil) and (current.lastChild ~= nil) do
            current = current.lastChild
        end
        return current
    end
    return self.parentNode
end

---@param test fun(DOM_Node):boolean
---@return boolean
function Node:isFollowedBy(test)
    local current = self:precedingNode()
    while current ~= nil do
        if test(current) then
            return true
        end
        current = current:precedingNode()
    end
    return false
end

---https://dom.spec.whatwg.org/#concept-tree-following
---@return DOM_Node|nil
function Node:followingNode()
    if self.firstChild == nil then
        local current = self
        while current ~= nil do
            if current.nextSibling ~= nil then
                return current.nextSibling
            end
            current = current.parentNode
        end
        return nil
    else
        return self.firstChild
    end
end

---@param test fun(DOM_Node):boolean
---@return boolean
function Node:isPrecededBy(test)
    local current = self:followingNode()
    while current ~= nil do
        if test(current) then
            return true
        end
        current = current:followingNode()
    end
    return false
end

---https://dom.spec.whatwg.org/#dom-node-text_node
---@return integer
function Node:nodeType()
    -- The nodeType getter steps are to return the first matching statement, switching on the interface this implements:

    -- Element
    if self.isElement then
        -- ELEMENT_NODE (1)
        return 1
    end
    -- Attr
    if self.isAttr then
        -- ATTRIBUTE_NODE (2);
        return 2
    end
    -- An exclusive Text node
    if self:isExclusiveTextNode() then
        -- TEXT_NODE (3);
        return 3
    end
    -- CDATASection
    if self.isCDATASection then
        -- CDATA_SECTION_NODE (4);
        return 4
    end
    -- ProcessingInstruction
    if self.isProcessingInstruction then
        -- PROCESSING_INSTRUCTION_NODE (7);
        return 7
    end
    -- Comment
    if self.isComment then
        -- COMMENT_NODE (8);
        return 8
    end
    -- Document
    if self.isDocument then
        -- DOCUMENT_NODE (9);
        return 9
    end
    -- DocumentType
    if self.isDocumentType then
        -- DOCUMENT_TYPE_NODE (10);
        return 10
    end
    -- DocumentFragment
    if self.isDocumentFragment then
        -- DOCUMENT_FRAGMENT_NODE (11).
        return 11
    end
    error("Unreachable")
end

---https://dom.spec.whatwg.org/#dom-node-nodevalue
---@return string|nil
function Node:nodeValue()
    -- https://dom.spec.whatwg.org/#dom-node-nodevalue
    if self.isAttr then
        -- this’s value.
        error("todo")
    elseif self.isCharacterData then
        -- this’s data.
        return self:asCharacterData().data
    else
        -- Null.
        return nil
    end
end

---https://dom.spec.whatwg.org/#dom-node-nodename
---@return string
function Node:nodeName()
    -- Attr
    if self.isElement then
        -- Its HTML-uppercased qualified name.
        return self:asElement():htmlUppercasedQualifiedName()
    elseif self.isAttr then
        -- Its qualified name.
        return self:asAttr():qualifiedName()
    elseif self:isExclusiveTextNode() then
        return "#text"
    elseif self.isCDATASection then
        return "#cdata-section"
    elseif self.isProcessingInstruction then
        error("TODO")
    elseif self.isComment then
        return "#comment"
    elseif self.isDocument then
        return "#document"
    elseif self.isDocumentType then
        return self:asDocumentType().name
    elseif self.isDocumentFragment then
        return "#document-fragment"
    end
    error("Unreachable")
end

---@param dest file*
---@param indent integer?
function Node:dump(dest, indent)
    if indent == nil then indent = 0 end
    for _ in 1, indent do
        dest:write(' ')
    end
    if self:nodeValue() ~= nil then
        dest:write(string.format("%s(%s) '%s'\n", self:nodeName(), self.nodeType, self:nodeValue()))
    else
        dest:write(string.format("%s(%s)\n", self:nodeName(), self.nodeType))
    end
    local child = self.firstChild
    while child ~= nil do
        child:dump(dest, indent + 2)
        child = child.nextSibling
    end
end

---https://dom.spec.whatwg.org/#concept-tree-root
---@return DOM_Node
function Node:root()
    local current = self
    while true do
        if current.parentNode == nil then
            return current
        end
        current = current.parentNode
    end
end

---https://dom.spec.whatwg.org/#concept-shadow-including-root
---@return DOM_Node
function Node:shadowIncludingRoot()
    local currentNode = self

    while true do
        if currentNode.isShadowRoot then
            currentNode = currentNode:asShadowRoot().host:shadowIncludingRoot()
        else
            return currentNode:root()
        end
    end
end

---@param node DOM_Node
---@param inclusive boolean
---@param includingShadow boolean
---@return fun():DOM_Node|nil
local function descendantsIter(node, inclusive, includingShadow)
    local current
    if inclusive then
        current = node
    else
        current = node.firstChild
    end
    return function()
        if current == nil then
            return nil
        end
        if includingShadow and current:isShadowRoot() then
            error("TODO")
        end
        return current.nextSibling
    end
end

---https://dom.spec.whatwg.org/#concept-tree-descendant
---@return fun():DOM_Node|nil
function Node:descendants()
    return descendantsIter(self, false, false)
end

---https://dom.spec.whatwg.org/#concept-tree-descendant
---@param n DOM_Node
---@return boolean
function Node:isDescendantOf(n)
    return iterutil.contains(n:descendants(), self)
end

---https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
---@return fun():DOM_Node|nil
function Node:inclusiveDescendants()
    return descendantsIter(self, true, false)
end

---https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
---@param n DOM_Node
---@return boolean
function Node:isInclusiveDescendantOf(n)
    return iterutil.contains(n:inclusiveDescendants(), self)
end

---https://dom.spec.whatwg.org/#concept-shadow-including-descendant
---@return fun():DOM_Node|nil
function Node:shadowIncludingDescendants()
    return descendantsIter(self, false, true)
end

---https://dom.spec.whatwg.org/#concept-tree-inclusive-descendant
---@param n DOM_Node
---@return boolean
function Node:isShadowIncludingDescendantOf(n)
    return iterutil.contains(n:shadowIncludingDescendants(), self)
end

---https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
function Node:shadowIncludingInclusiveDescendants()
    return descendantsIter(self, true, true)
end

----https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-descendant
---@param n DOM_Node
---@return boolean
function Node:isShadowIncludingInclusiveDescendantOf(n)
    return iterutil.contains(n:shadowIncludingDescendants(), self)
end

---@param node DOM_Node
---@param inclusive boolean
---@param includingShadow boolean
---@return fun():DOM_Node|nil
local function ancestorsIter(node, inclusive, includingShadow)
    local current
    if inclusive then
        current = node
    else
        current = node.parentNode
    end
    return function()
        if current == nil then
            return nil
        end
        local ret = current
        if includingShadow and current:isShadowRoot() then
            error("TODO")
        end
        current = current.parentNode
        return ret
    end
end

---https://dom.spec.whatwg.org/#concept-tree-ancestor
---@return fun():DOM_Node|nil
function Node:ancestors()
    return ancestorsIter(self, false, false)
end

---https://dom.spec.whatwg.org/#concept-tree-ancestor
---@param n DOM_Node
---@return boolean
function Node:isAncestorOf(n)
    return iterutil.contains(n:ancestors(), self)
end

---https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor
---@return fun():DOM_Node|nil
function Node:inclusiveAncestors()
    return ancestorsIter(self, true, false)
end

---https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor
---@param n DOM_Node
---@return boolean
function Node:isInclusiveAncestorOf(n)
    return iterutil.contains(n:inclusiveAncestors(), self)
end

---https://dom.spec.whatwg.org/#concept-shadow-including-ancestor
---@return fun():DOM_Node|nil
function Node:shadowIncludingAncestors()
    return ancestorsIter(self, false, true)
end

---https://dom.spec.whatwg.org/#concept-tree-inclusive-ancestor
---@param n DOM_Node
---@return boolean
function Node:isShadowIncludingAncestorOf(n)
    return iterutil.contains(n:shadowIncludingAncestors(), self)
end

---https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-ancestor
---@return fun():DOM_Node|nil
function Node:shadowIncludingInclusiveAncestors()
    return ancestorsIter(self, true, true)
end

---https://dom.spec.whatwg.org/#concept-shadow-including-inclusive-ancestor
---@param n DOM_Node
---@return boolean
function Node:isShadowIncludingInclusiveAncestorOf(n)
    return iterutil.contains(n:shadowIncludingAncestors(), self)
end

---https://dom.spec.whatwg.org/#concept-tree-host-including-inclusive-ancestor
---@return fun():DOM_Node|nil
function Node:hostIncludingInclusiveAncestors()
    local current = self

    return function()
        local ret = current
        if current.isDocumentFragment then
            current = current:asDocumentFragment().host
        else
            current = current.parentNode
        end
        return ret
    end
end

---https://dom.spec.whatwg.org/#concept-tree-host-including-inclusive-ancestor
---@param n DOM_Node
---@return boolean
function Node:isHostIncludingInclusiveAncestorOf(n)
    return iterutil.contains(n:hostIncludingInclusiveAncestors(), self)
end

---https://dom.spec.whatwg.org/#connected
---@return boolean
function Node:isConnected()
    -- A node is connected if its shadow-including root is a document.
    local root = self:shadowIncludingRoot()
    return root.isDocument
end

---https://dom.spec.whatwg.org/#concept-tree-following
---@return fun():DOM_Node|nil
function Node:followingNodes()
    local mCurrent = self:followingNode()

    return function()
        local ret = mCurrent
        if mCurrent == nil then
            return nil
        end
        mCurrent = mCurrent:followingNode()
        return ret
    end
end

---https://dom.spec.whatwg.org/#concept-tree-preceding
---@return fun():DOM_Node|nil
function Node:precedingNodes()
    local mCurrent = self:precedingNode()

    return function()
        local ret = mCurrent
        if mCurrent == nil then
            return nil
        end
        mCurrent = mCurrent:precedingNode()
        return ret
    end
end

---https://dom.spec.whatwg.org/#concept-tree-preceding
---@return fun():DOM_Node|nil
function Node:nextSiblings()
    local mCurrent = self.nextSibling

    return function()
        local ret = mCurrent
        if mCurrent == nil then
            return nil
        end
        mCurrent = mCurrent.nextSibling
        return ret
    end
end

---https://dom.spec.whatwg.org/#concept-tree-preceding
---@return fun():DOM_Node|nil
function Node:previousSiblings()
    local mCurrent = self.previousSibling

    return function()
        local ret = mCurrent
        if mCurrent == nil then
            return nil
        end
        mCurrent = mCurrent.previousSibling
        return ret
    end
end

---https://dom.spec.whatwg.org/#parent-element
---@return DOM_Node|nil
function Node:parentElement()
    if self.parentNode == nil then
        return nil
    end
    if self.parentNode.isElement then
        return self.parentNode
    end
    return nil
end

---https://dom.spec.whatwg.org/#concept-tree-index
---@return integer
function Node:index()
    local index = 0
    for s in self:previousSiblings() do
        index = index + 1
    end
    return index
end

---https://dom.spec.whatwg.org/#concept-tree-child
---@return fun():DOM_Node|nil
function Node:childNodes()
    local current = self.firstChild
    return function()
        local ret = current
        if current ~= nil then
            current = current.nextSibling
        end
        return ret
    end
end

---https://dom.spec.whatwg.org/#assign-slotables-for-a-tree
function Node:assignSlottablesForTree()
    for node in self:inclusiveDescendants() do
        if node.isElement and node:asElement():isHTMLElement("slot") then
            error("TODO")
        end
    end
end

---https://dom.spec.whatwg.org/#concept-node-adopt
function Node:adoptInto(document)
    -- 1. Let oldDocument be node’s node document.
    local oldDocument = self.nodeDocument

    -- 2. If node’s parent is non-null, then remove node.
    if self.parentNode ~= nil then
        error("TODO")
    end

    -- 3. If document is not oldDocument:
    if document ~= oldDocument then
        -- 1. For each inclusiveDescendant in node’s shadow-including inclusive descendants:
        for inclusiveDescendant in self:shadowIncludingInclusiveDescendants() do
            -- 1. Set inclusiveDescendant’s node document to document.
            inclusiveDescendant.nodeDocument = document

            -- 2. For each inclusiveDescendant in node’s shadow-including inclusive descendants that is custom, enqueue a custom element callback reaction with inclusiveDescendant, callback name "adoptedCallback", and « oldDocument, document ».
            if inclusiveDescendant.isElement then
                error("TODO")
            end
        end
        -- 2. For each inclusiveDescendant in node’s shadow-including inclusive descendants that is custom, enqueue a custom element callback reaction with inclusiveDescendant, callback name "adoptedCallback", and « oldDocument, document ».
        for inclusiveDescendant in self:shadowIncludingInclusiveDescendants() do
            if inclusiveDescendant.isElement then
                error("TODO")
            end
        end
        -- 3. For each inclusiveDescendant in node’s shadow-including inclusive descendants, in shadow-including tree order, run the adopting steps with inclusiveDescendant and oldDocument.
        for inclusiveDescendant in self:shadowIncludingInclusiveDescendants() do
            inclusiveDescendant:runAdoptingSteps(oldDocument)
        end
    end
end

---@param parent DOM_Node
---@param node DOM_Node
local function appendChild(parent, node)
    local prevChild = parent.lastChild

    -- Set parent
    node.parentNode = parent

    -- Set previous/next sibling
    node.nextSibling = nil
    node.previousSibling = prevChild
    if prevChild ~= nil then
        prevChild.nextSibling = node
    end

    -- Set first child if needed
    if parent.firstChild == nil then
        parent.firstChild = node
    end
    parent.lastChild = node
end

---https://dom.spec.whatwg.org/#concept-node-insert
---@param parent DOM_Node
---@param beforeChild DOM_Node?
---@param suppressObservers_ boolean?
function Node:insert(parent, beforeChild, suppressObservers_)
    local suppressObservers = suppressObservers_ or false

    -- 1. Let nodes be node’s children, if node is a DocumentFragment node; otherwise « node ».
    local nodes = {} ---@type DOM_Node[]
    if self.isDocumentFragment then
        nodes = iterutil.collect(self:childNodes())
    else
        nodes = { self }
    end

    -- 2. Let count be nodes’s size.
    local count = #nodes

    -- 3. If count is 0, then return.
    if count == 0 then
        return
    end

    -- 4. If node is a DocumentFragment node:
    if self.isDocumentFragment then
        error("TODO")
        -- 1. Remove its children with the suppress observers flag set.
        -- 2. Queue a tree mutation record for node with « », nodes, null, and null.
    end

    -- 5. If child is non-null, then:
    if beforeChild ~= nil then
        if #Range:liveRanges() ~= 0 then
            error("TODO")
        end

        -- 1. For each live range whose start node is parent and start offset is greater than child’s index, increase its start offset by count.

        -- 2. For each live range whose end node is parent and end offset is greater than child’s index, increase its end offset by count.
    end

    -- 6. Let previousSibling be child’s previous sibling or parent’s last child if child is null.
    local previousSibling
    if beforeChild == nil then
        previousSibling = parent.lastChild
    else
        previousSibling = beforeChild.previousSibling
    end

    -- 7. For each node in nodes, in tree order:
    for _, node in ipairs(nodes) do
        -- 1. Adopt node into parent’s node document.
        node:adoptInto(parent.nodeDocument)

        -- 2. If child is null, then append node to parent’s children.
        if beforeChild == nil then
            appendChild(parent, node)
            -- 3. Otherwise, insert node into parent’s children before child’s index.
        else
            insert(parent, node, beforeChild)
        end

        -- 4. If parent is a shadow host whose shadow root’s slot assignment is "named" and node is a slottable, then assign a slot for node.
        if parent.isElement then
            if parent:asElement():isShadowHost() then
                error("TODO")
            end
        end

        -- 5. If parent’s root is a shadow root, and parent is a slot whose assigned nodes is the empty list, then run signal a slot change for parent.
        if parent:root().isElement then
            error("TODO")
        end

        -- 6. Run assign slottables for a tree with node’s root.
        self:root():assignSlottablesForTree()

        -- 7. For each shadow-including inclusive descendant inclusiveDescendant of node, in shadow-including tree order:
        for inclusiveDescendant in self:shadowIncludingInclusiveDescendants() do
            -- 1. Run the insertion steps with inclusiveDescendant.
            inclusiveDescendant:runInsertionSteps()

            -- 2. If inclusiveDescendant is connected:
            if inclusiveDescendant:isConnected() then
                -- 1. If inclusiveDescendant is custom, then enqueue a custom element callback reaction with inclusiveDescendant, callback name "connectedCallback", and « ».
                if inclusiveDescendant.isElement and inclusiveDescendant:isElement():isCustom() then
                    error("TODO")
                    -- 2. Otherwise, try to upgrade inclusiveDescendant.
                else
                    -- TODO
                end
            end
        end
    end

    -- 8. If suppress observers flag is unset, then queue a tree mutation record for parent with nodes, « », previousSibling, and child.
    if not suppressObservers then
        -- TODO
    end

    -- 9. Run the children changed steps for parent.
    parent:runChildrenChangedSteps()

    -- 10. Let staticNodeList be a list of nodes, initially « ».
    local staticNodeList = {}

    -- 11. For each node of nodes, in tree order:
    for _, node in ipairs(nodes) do
        -- 1. For each shadow-including inclusive descendant inclusiveDescendant of node, in shadow-including tree order, append inclusiveDescendant to staticNodeList.
        for inclusiveDescendant in node:shadowIncludingInclusiveDescendants() do
            staticNodeList.add(inclusiveDescendant)
        end
    end

    -- 12. For each node of staticNodeList, if node.isconnected, then run the post-connection steps with node.
    for node in staticNodeList do
        if node:isConnected() then
            node:runPostConnectionSteps()
        end
    end
end

---https://dom.spec.whatwg.org/#concept-node-ensure-pre-insertion-validity
---@param parent DOM_Node
---@param beforeChild DOM_Node?
function Node:ensurePreInsertionValidity(parent, beforeChild)
    -- 1. If parent is not a Document, DocumentFragment, or Element node,
    if (parent.isDocument) and (parent.isDocumentFragment) and (parent.isElement) then
        -- then throw a "HierarchyRequestError" DOMException.
        error("DOMException.HierarchyRequestError")
    end

    -- 2. If node is a host-including inclusive ancestor of parent,
    if self:isHostIncludingInclusiveAncestorOf(parent) then
        -- then throw a "HierarchyRequestError" DOMException.
        error("DOMException.HierarchyRequestError")
    end

    -- 3. If child is non-null and its parent is not parent,
    if (beforeChild ~= nil) and (beforeChild.parentNode ~= parent) then
        -- then throw a "NotFoundError" DOMException.
        error("DOMException.NotFoundError")
    end

    -- 4. If node is not a DocumentFragment, DocumentType, Element, or CharacterData node,
    if (self.isDocumentFragment) and (self.isDocumentType) and (self.isCharacterData) then
        --  then throw a "HierarchyRequestError" DOMException.
        error("DOMException.HierarchyRequestError")
    end

    -- 5. If either node is a Text node and parent is a document, or node is a doctype and parent is not a document,
    if (self.isText and parent.isDocument) or not (self.isDocumentType and parent.isDocument) then
        -- then throw a "HierarchyRequestError" DOMException.
        error("DOMException.HierarchyRequestError")
    end

    -- 6. If parent is a document, and any of the statements below, switched on the interface node implements, are true, then throw a "HierarchyRequestError" DOMException
    if parent.isDocument then
        local cond = false
        -- DocumentFragment
        if self.isDocumentFragment then
            -- If node has more than one element child
            if (0 < iterutil.count(
                    self:childNodes(),
                    ---@param x DOM_Node
                    function(x) return x.isElement end
                )) or
                --  or has a Text node child.
                (0 < iterutil.count(
                    self:childNodes(),
                    ---@param x DOM_Node
                    function(x) return x.isText end) ~= 0
                )
            then
                cond = true
                -- Otherwise, if node has one element child and
            elseif (iterutil.count(
                    self:childNodes(),
                    ---@param x DOM_Node
                    function(x) return x.isElement end
                ) == 1) and (
                -- either parent has an element child,
                    iterutil.count(
                        parent:childNodes(),
                        ---@param x DOM_Node
                        function(x) return x.isElement end
                    ) or
                    -- child is a doctype,
                    ((beforeChild ~= nil) and (beforeChild.isDocumentType)) or
                    -- or child is non-null and a doctype is following child.
                    ((beforeChild ~= nil) and (beforeChild:followingNode() ~= nil) and (beforeChild:followingNode().isDocumentType))
                )
            then
                cond = true
            end
            -- Element
        elseif self.isElement then
            -- parent has an element child,
            if (0 < iterutil.count(
                    parent:childNodes(),
                    ---@param x DOM_Node
                    function(x) return x.isElement end
                )) or
                -- child is a doctype,
                ((beforeChild ~= nil) and (beforeChild.isDocumentType)) or
                -- or child is non-null and a doctype is following child.
                ((beforeChild ~= nil) and (beforeChild:followingNode() ~= nil) and (beforeChild:followingNode().isDocumentType))
            then
                cond = true
            end
            -- DocumentType
        elseif self.isDocumentType then
            -- parent has a doctype child,
            if (0 < iterutil.count(
                    parent:childNodes(),
                    ---@param x DOM_Node
                    function(x) return x.isDocumentType end
                )) or
                -- child is non-null and an element is preceding child,
                (beforeChild ~= nil) and (beforeChild:precedingNode() ~= nil) and (beforeChild:precedingNode().isElement) or
                -- or child is null and parent has an element child.
                (beforeChild == nil) and (0 < iterutil.count(
                    parent:childNodes(),
                    ---@param x DOM_Node
                    function(x) return x.isElement end
                ))
            then
                cond = true
            end
        end
        if cond then
            error("DOMException.HierarchyRequestError")
        end
    end
end

---https://dom.spec.whatwg.org/#concept-node-pre-insert
---@param parent DOM_Node
---@param beforeChild DOM_Node?
---@return DOM_Node
function Node:preInsert(parent, beforeChild)
    -- 1. Ensure pre-insertion validity of node into parent before child.
    self:ensurePreInsertionValidity(parent, beforeChild)

    -- 2. Let referenceChild be child.
    local referenceChild = beforeChild


    -- 3. If referenceChild is node, then set referenceChild to node’s next sibling.
    if referenceChild == self then
        referenceChild = self.nextSibling
    end

    -- 4. Insert node into parent before referenceChild.
    self:insert(parent, referenceChild)

    -- 5. Return node.
    return self
end

---https://dom.spec.whatwg.org/#concept-node-append
--- @param parent DOM_Node
function Node:append(parent)
    return self:preInsert(parent, nil)
end

---https://dom.spec.whatwg.org/#dom-node-appendchild
--- @param node DOM_Node
function Node:appendChild(node)
    return node:append(self)
end

local function remove(node)
    local parent = node.parentNode
    local nextSibling = node.nextSibling
    local previousSibling = node.previousSibling

    -- Set parent
    node.parentNode = nil

    -- Set previous/next sibling
    if previousSibling ~= nil then
        previousSibling.nextSibling = nextSibling
    end
    if nextSibling ~= nil then
        nextSibling.previousSibling = previousSibling
    end

    -- Set first/last child
    if previousSibling == nil then
        parent.firstChild = nextSibling
    end
    if nextSibling == nil then
        parent.lastChild = previousSibling
    end
end

---https://dom.spec.whatwg.org/#concept-node-remove
---@param suppressObservers boolean?
function Node:remove(suppressObservers)
    if suppressObservers == nil then suppressObservers = false end

    -- 1. Let parent be node’s parent.
    local parent = self.parentNode

    -- 2. Assert: parent is non-null.
    if parent == nil then
        error("parent shouldn't be nil here")
    end

    -- 3. Let index be node’s index.
    local index = self:index()

    if #Range:liveRanges() ~= 0 then
        error("TODO")
        -- 4. For each live range whose start node is an inclusive descendant of node, set its start to (parent, index).
        -- 5. For each live range whose end node is an inclusive descendant of node, set its end to (parent, index).
        -- 6. For each live range whose start node is parent and start offset is greater than index, decrease its start offset by 1.
        -- 7. For each live range whose end node is parent and end offset is greater than index, decrease its end offset by 1.
    end

    if #NodeIterator:nodeIterators() ~= 0 then
        error("TODO")
        -- 8. For each NodeIterator object iterator whose root’s node document is node’s node document, run the NodeIterator pre-removing steps given node and iterator.
    end

    -- 9. Let oldPreviousSibling be node’s previous sibling.
    local oldPreviousSibling = self.previousSibling

    -- 10. Let oldNextSibling be node’s next sibling.
    local oldNextSibling = self.nextSibling

    -- 11. Remove node from its parent’s children.
    remove(self)

    -- 12. If node is assigned, then run assign slottables for node’s assigned slot.
    if self.isSlottable then
        error("TODO")
    end

    -- 13. If parent’s root is a shadow root, and parent is a slot whose assigned nodes is the empty list, then run signal a slot change for parent.
    if parent:root().isShadowRoot then
        error("TODO")
    end

    -- 14. If node has an inclusive descendant that is a slot, then:
    if iterutil.contains(
            self:childNodes(),
            ---@param n DOM_Node
            function(n)
                return n.isElement and n:asElement():isHTMLElement("slot")
            end
        ) then
        error("TODO")
        -- 1. Run assign slottables for a tree with parent’s root.
        -- 2. Run assign slottables for a tree with node.
    end

    -- 15. Run the removing steps with node and parent.
    self:runRemovingSteps(parent)

    -- 16. Let isParentConnected be parent’s connected.
    local isParentConnected = parent:isConnected()

    -- 17. If node is custom and isParentConnected is true, then enqueue a custom element callback reaction with node, callback name "disconnectedCallback", and an empty argument list.
    if self.isElement and self:asElement():isCustom() and isParentConnected then
        error("TODO")
    end

    -- 18. For each shadow-including descendant descendant of node, in shadow-including tree order, then:
    for descendant in self:shadowIncludingDescendants() do
        -- 1. Run the removing steps with descendant and null.
        self:runRemovingSteps(nil)

        -- 2. If descendant is custom and isParentConnected is true, then enqueue a custom element callback reaction with descendant, callback name "disconnectedCallback", and an empty argument list.
        if descendant.isElement and descendant:asElement():isCustom() and isParentConnected then
            error("TODO")
        end
    end

    -- 19. For each inclusive ancestor inclusiveAncestor of parent,
    for inclusiveAncestor in self:inclusiveAncestors() do
        -- and then for each registered of inclusiveAncestor’s registered observer list,
        for _, registered in ipairs(inclusiveAncestor.mutationObservers) do
            -- if registered’s options["subtree"] is true, then append a new transient registered observer whose observer is registered’s observer, options is registered’s options, and source is registered to node’s registered observer list.
            error("TODO")
        end
    end

    -- 20. If suppress observers flag is unset,
    if not suppressObservers then
        -- then queue a tree mutation record for parent with « », « node », oldPreviousSibling, and oldNextSibling.
        -- TODO
    end

    -- 21. Run the children changed steps for parent.
    parent:runChildrenChangedSteps()
end

return Node
