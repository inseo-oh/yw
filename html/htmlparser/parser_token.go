// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package htmlparser

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

type htmlToken interface {
	equals(other htmlToken) bool // XXX: Do we need this to be here?
	String() string
}

type eofToken struct{}

func (tk eofToken) equals(other htmlToken) bool {
	_, ok := other.(*eofToken)
	return ok
}
func (tk eofToken) String() string {
	return "<!EOF>"
}

type charToken struct{ value rune }

func (tk charToken) equals(other htmlToken) bool {
	otherTk, ok := other.(*charToken)
	if !ok {
		return false
	}
	if tk.value != otherTk.value {
		return false
	}
	return true
}
func (tk charToken) isCharTokenWithOneOf(chars string) bool {
	for _, c := range chars {
		if tk.value == c {
			return true
		}
	}
	return false
}

func (tk charToken) String() string {
	return fmt.Sprintf("%c", tk.value)
}

type commentToken struct{ data string }

func (tk commentToken) equals(other htmlToken) bool {
	otherTk, ok := other.(*commentToken)
	if !ok {
		return false
	}
	if tk.data != otherTk.data {
		return false
	}
	return true
}
func (tk commentToken) String() string {
	return fmt.Sprintf("<!-- %s -->", strconv.Quote(tk.data))
}

type doctypeToken struct {
	// nil value represents "missing"
	name        *string
	publicId    *string
	systemId    *string
	forceQuirks bool
}

func (tk doctypeToken) equals(other htmlToken) bool {
	otherTk, ok := other.(*doctypeToken)
	if !ok {
		return false
	}
	if (tk.name == nil) != (otherTk.name == nil) {
		return false
	}
	if tk.name != nil && *tk.name != *otherTk.name {
		return false
	}
	return true
}
func (tk doctypeToken) String() string {
	sb := strings.Builder{}
	sb.WriteString("<!DOCTYPE")
	if tk.name != nil {
		sb.WriteString(fmt.Sprintf(" %s", strconv.Quote(*tk.name)))
	}
	if tk.publicId != nil {
		sb.WriteString(fmt.Sprintf(" PUBLIC %s", strconv.Quote(*tk.publicId)))
	}
	if tk.systemId != nil {
		sb.WriteString(fmt.Sprintf(" SYSTEM %s", strconv.Quote(*tk.systemId)))
	}
	sb.WriteString(">")
	return sb.String()
}

type tagToken struct {
	isEnd         bool
	isSelfClosing bool
	tagName       string
	attrs         []dom.AttrData

	selfClosingAcknowledged bool
}

func (tk tagToken) equals(other htmlToken) bool {
	otherTk, ok := other.(*tagToken)
	if !ok {
		return false
	}
	if len(tk.attrs) != len(otherTk.attrs) {
		return false
	}
	if tk.isEnd != otherTk.isEnd {
		return false
	}
	for i := 0; i < len(tk.attrs); i++ {
		attr := tk.attrs[i]
		otherTk := otherTk.attrs[i]
		if attr.LocalName != otherTk.LocalName {
			return false
		}
		if attr.Value != otherTk.Value {
			return false
		}
	}
	if tk.isSelfClosing != otherTk.isSelfClosing {
		return false
	}
	if tk.tagName != otherTk.tagName {
		return false
	}
	return true
}

func (tk tagToken) isStartTag() bool {
	return !tk.isEnd
}
func (tk tagToken) isEndTag() bool {
	return tk.isEnd
}

// Returns nil if there's no such attribute
func (tk tagToken) Attr(name string) (string, bool) {
	for _, attr := range tk.attrs {
		if attr.LocalName == name {
			return attr.Value, true
		}
	}
	return "", false
}

func (tk tagToken) String() string {
	sb := strings.Builder{}
	if tk.isEnd {
		sb.WriteString("</")
	} else {
		sb.WriteString("<")
	}
	sb.WriteString(tk.tagName)
	for _, attr := range tk.attrs {
		sb.WriteString(" ")
		val := strconv.Quote(attr.Value)
		if ns := attr.Namespace; ns != nil {
			sb.WriteString(fmt.Sprintf("%v:%s=%s", *ns, attr.LocalName, val))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s", attr.LocalName, val))
		}
	}
	if tk.isSelfClosing {
		sb.WriteString("/")
	}
	sb.WriteString(">")
	return sb.String()
}

type tokenizerState uint8

const (
	dataState tokenizerState = iota
	rcdataState
	rawtextState
	plaintextState
	tagOpenState
	endTagOpenState
	tagNameState
	rcdataLessThanSignState
	rcdataEndTagOpenState
	rcdataEndTagNameState
	rawtextLessThanSignState
	rawtextEndTagOpenState
	rawtextEndTagNameState
	beforeAttributeNameState
	attributeNameState
	afterAttributeNameState
	beforeAttributeValueState
	attributeValueDoubleQuotedState
	attributeValueSingleQuotedState
	attributeValueUnquotedState
	afterAttributeValueQuotedState
	selfClosingStartTagState
	bogusCommentState
	markupDeclarationOpenState
	commentStartState
	commentStartDashState
	commentState
	commentLessThanSignState
	commentEndDashState
	commentEndState
	doctypeState
	beforeDoctypeNameState
	doctypeNameState
	afterDoctypeNameState
	characterReferenceState
	namedCharacterReferenceState
	numericCharacterReferenceState
	hexadecimalCharacterReferenceStartState
	decimalCharacterReferenceStartState
	hexadecimalCharacterReferenceState
	decimalCharacterReferenceState
	numericCharacterReferenceEndState
)

type parseError string

const (
	absence_of_digits_in_numeric_character_reference_error = parseError("absence-of-digits-in-numeric-character-reference")
	abrupt_closing_of_empty_comment_error                  = parseError("abrupt-closing-of-empty-comment")
	character_reference_outside_unicode_range_error        = parseError("character-reference-outside-unicode-range")
	control_character_reference_error                      = parseError("control-character-reference")
	eof_before_tag_name_error                              = parseError("eof-before-tag-name")
	eof_in_comment_error                                   = parseError("eof-in-comment")
	eof_in_doctype_error                                   = parseError("eof-in-doctype")
	eof_in_tag_error                                       = parseError("eof-in-tag")
	incorrectly_opened_comment_error                       = parseError("incorrectly-opened-comment")
	invalid_character_sequence_after_doctype_name_error    = parseError("invalid-character-sequence-after-doctype-name")
	invalid_first_character_of_tag_name_error              = parseError("invalid-first-character-of-tag-name")
	missing_attribute_value_error                          = parseError("missing-attribute-value")
	missing_doctype_name_error                             = parseError("missing-doctype-name")
	missing_end_tag_name_error                             = parseError("missing-end-tag-name")
	missing_semicolon_after_character_reference_error      = parseError("missing-semicolon-after-character-reference")
	missing_whitespace_before_doctype_name_error           = parseError("missing-whitespace-before-doctype-name")
	missing_whitepace_between_attributes_error             = parseError("missing-whitespace-between-attributes")
	noncharacter_reference_error                           = parseError("noncharacter-character-reference")
	null_character_reference_error                         = parseError("null-character-reference")
	surrogate_character_reference_error                    = parseError("surrogate-character-reference")
	unexpected_character_in_attribute_name_error           = parseError("unexpected-character-in-attribute-name")
	unexpected_character_in_unquoted_attribute_value_error = parseError("unexpected-character-in-unquoted-attribute-value")
	unexpected_equals_sign_before_attribute_name_error     = parseError("unexpected-equals-sign-before-attribute-name")
	unexpected_null_character_error                        = parseError("unexpected-null-character")
	unexpected_question_mark_instead_of_tag_name_error     = parseError("unexpected-question-mark-instead-of-tag-name")
	unexpected_solidus_in_tag_error                        = parseError("unexpected-solidus-in-tag")
)

type tokenizer struct {
	tkh              util.TokenizerHelper
	state            tokenizerState
	parserPauseFlag  bool
	lastStartTagName string
	onTokenEmitted   func(tk htmlToken)
}

func newTokenizer(str string) tokenizer {
	return tokenizer{
		tkh: util.TokenizerHelper{
			Str: []rune(str),
		},
	}
}

func (t *tokenizer) parseErrorEncountered(err parseError) {
	fmt.Println(err)
}

func (t *tokenizer) run() {
	var currTk htmlToken
	var returnState tokenizerState
	tempBuf := ""
	characterReferenceCode := 0

	attrsToRemove := []int{}
	currTagToken := func() *tagToken {
		tok := currTk.(*tagToken)
		return tok
	}
	currCommentToken := func() *commentToken {
		tok := currTk.(*commentToken)
		return tok
	}
	currDoctypeToken := func() *doctypeToken {
		tok := currTk.(*doctypeToken)
		return tok
	}
	currAttr := func() *dom.AttrData {
		attrs := currTagToken().attrs
		return &attrs[len(attrs)-1]
	}
	checkDuplicateAttrName := func() {
		attrs := currTagToken().attrs
		for i, attr := range attrs {
			if i == (len(attrs) - 1) {
				continue
			}
			if currAttr().LocalName == attr.LocalName {
				attrsToRemove = append(attrsToRemove, i)
			}
		}
	}
	emitToken := func(tk htmlToken) {
		if tagTk, ok := tk.(*tagToken); ok {
			finalAttrs := []dom.AttrData{}
			for i, attr := range currTagToken().attrs {
				badAttr := slices.Contains(attrsToRemove, i)
				if !badAttr {
					finalAttrs = append(finalAttrs, attr)
				}
			}
			currTagToken().attrs = finalAttrs
			if tagTk.isStartTag() {
				t.lastStartTagName = tagTk.tagName
			}
		}
		t.onTokenEmitted(tk)
	}
	isConsumedAsPartOfAttr := func() bool {
		switch returnState {
		case attributeValueDoubleQuotedState,
			attributeValueSingleQuotedState,
			attributeValueUnquotedState:
			return true
		default:
			return false
		}
	}
	flushCodepointsConsumedAsCharReference := func() {
		if isConsumedAsPartOfAttr() {
			for _, c := range tempBuf {
				currAttr().Value += string(c)
			}
		} else {
			for _, c := range tempBuf {
				emitToken(&charToken{value: c})
			}
		}
	}
	isAppropriateEndTagToken := func(tk tagToken) bool {
		return t.lastStartTagName == tk.tagName
	}

	for {
		if t.parserPauseFlag {
			return
		}

		switch t.state {
		// https://html.spec.whatwg.org/multipage/parsing.html#data-state
		case dataState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '&':
				returnState = dataState
				t.state = characterReferenceState
			case '<':
				t.state = tagOpenState
			case 0x0000:
				t.parseErrorEncountered(unexpected_null_character_error)
				emitToken(&charToken{value: nextChar})
			case -1:
				emitToken(&eofToken{})
				return
			default:
				emitToken(&charToken{value: nextChar})
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-state
		case rcdataState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '&':
				returnState = dataState
				t.state = characterReferenceState
			case '<':
				t.state = rcdataLessThanSignState
			case 0x0000:
				t.parseErrorEncountered(unexpected_null_character_error)
				emitToken(&charToken{value: 0xfffd})
			case -1:
				emitToken(&eofToken{})
				return
			default:
				emitToken(&charToken{value: nextChar})
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-state
		case rawtextState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '<':
				t.state = rawtextLessThanSignState
			case 0x0000:
				t.parseErrorEncountered(unexpected_null_character_error)
				emitToken(&charToken{value: 0xfffd})
			case -1:
				emitToken(&eofToken{})
				return
			default:
				emitToken(&charToken{value: nextChar})
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#plaintext-state
		case plaintextState:
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#plaintext-state]")
		// https://html.spec.whatwg.org/multipage/parsing.html#tag-open-state
		case tagOpenState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '!':
				t.state = markupDeclarationOpenState
			case '/':
				t.state = endTagOpenState
			case '?':
				t.parseErrorEncountered(unexpected_question_mark_instead_of_tag_name_error)
				currTk = &commentToken{data: ""}
				t.tkh.Cursor--
				t.state = bogusCommentState
			case -1:
				emitToken(&charToken{value: '<'})
				emitToken(&eofToken{})
				return
			default:
				if util.AsciiAlphaRegex.MatchString(string(nextChar)) {
					currTk = &tagToken{
						isEnd:         false,
						isSelfClosing: false,
						tagName:       "",
						attrs:         []dom.AttrData{},
					}
					t.tkh.Cursor--
					t.state = tagNameState
				} else {
					t.parseErrorEncountered(invalid_first_character_of_tag_name_error)
					emitToken(&charToken{value: '<'})
					t.tkh.Cursor--
					t.state = dataState
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#end-tag-open-state
		case endTagOpenState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '>':
				t.parseErrorEncountered(missing_end_tag_name_error)
				t.state = dataState
			case -1:
				emitToken(&charToken{value: '<'})
				emitToken(&charToken{value: '/'})
				emitToken(&eofToken{})
				return
			default:
				if util.AsciiAlphaRegex.MatchString(string(nextChar)) {
					currTk = &tagToken{
						isEnd:   true,
						tagName: "",
					}
					t.tkh.Cursor--
					t.state = tagNameState
				} else {
					t.parseErrorEncountered(invalid_first_character_of_tag_name_error)
					currTagToken().tagName += string(rune(0xfffd))
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#tag-name-state
		case tagNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				t.state = beforeAttributeNameState
			case '/':
				t.state = selfClosingStartTagState
			case '>':
				t.state = dataState
				emitToken(currTk)
			case 0x0000:
				t.parseErrorEncountered(unexpected_null_character_error)
				currTagToken().tagName += string(rune(0xfffd))
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
				return
			default:
				if util.AsciiUppercaseRegex.MatchString(string(nextChar)) {
					currTagToken().tagName += strings.ToLower(string(nextChar))
				} else {
					currTagToken().tagName += string(nextChar)
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-less-than-sign-state
		case rcdataLessThanSignState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '/':
				tempBuf = ""
				t.state = rcdataEndTagOpenState
			default:
				emitToken(&charToken{value: '<'})
				t.tkh.Cursor--
				t.state = rcdataState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-open-state
		case rcdataEndTagOpenState:
			nextChar := t.tkh.ConsumeChar()
			if util.AsciiAlphaRegex.MatchString(string(nextChar)) {
				currTk = &tagToken{
					isEnd:   true,
					tagName: "",
				}
				t.tkh.Cursor--
				t.state = rcdataEndTagNameState
			} else {
				emitToken(&charToken{value: '<'})
				emitToken(&charToken{value: '/'})
				t.tkh.Cursor--
				t.state = rcdataState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-name-state
		case rcdataEndTagNameState:
			nextChar := t.tkh.ConsumeChar()
			anythingElse := func() {
				emitToken(&charToken{value: '<'})
				emitToken(&charToken{value: '/'})
				for _, c := range tempBuf {
					emitToken(&charToken{value: c})
				}
				t.tkh.Cursor--
				t.state = rcdataState
			}
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				if isAppropriateEndTagToken(*currTagToken()) {
					t.state = beforeAttributeNameState
				} else {
					anythingElse()
				}
			case '/':
				if isAppropriateEndTagToken(*currTagToken()) {
					t.state = selfClosingStartTagState
				} else {
					anythingElse()
				}
			case '>':
				if isAppropriateEndTagToken(*currTagToken()) {
					t.state = dataState
					emitToken(currTk)
				} else {
					anythingElse()
				}
			default:
				if util.AsciiUppercaseRegex.MatchString(string(nextChar)) {
					currTagToken().tagName += strings.ToLower(string(nextChar))
					tempBuf += string(nextChar)
				} else if util.AsciiLowercaseRegex.MatchString(string(nextChar)) {
					currTagToken().tagName += string(nextChar)
					tempBuf += string(nextChar)
				} else {
					anythingElse()
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-less-than-sign-state
		case rawtextLessThanSignState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '/':
				tempBuf = ""
				t.state = rawtextEndTagOpenState
			default:
				emitToken(&charToken{value: '<'})
				t.tkh.Cursor--
				t.state = rawtextState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-end-tag-open-state
		case rawtextEndTagOpenState:
			nextChar := t.tkh.ConsumeChar()
			if util.AsciiAlphaRegex.MatchString(string(nextChar)) {
				currTk = &tagToken{
					isEnd:   true,
					tagName: "",
				}
				t.tkh.Cursor--
				t.state = rawtextEndTagNameState
			} else {
				emitToken(&charToken{value: '<'})
				emitToken(&charToken{value: '/'})
				t.tkh.Cursor--
				t.state = rawtextState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-end-tag-name-state
		case rawtextEndTagNameState:
			nextChar := t.tkh.ConsumeChar()
			anythingElse := func() {
				emitToken(&charToken{value: '<'})
				emitToken(&charToken{value: '/'})
				for _, c := range tempBuf {
					emitToken(&charToken{value: c})
				}
				t.tkh.Cursor--
				t.state = rawtextState
			}
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				if isAppropriateEndTagToken(*currTagToken()) {
					t.state = beforeAttributeNameState
				} else {
					anythingElse()
				}
			case '/':
				if isAppropriateEndTagToken(*currTagToken()) {
					t.state = selfClosingStartTagState
				} else {
					anythingElse()
				}
			case '>':
				if isAppropriateEndTagToken(*currTagToken()) {
					t.state = dataState
					emitToken(currTk)
				} else {
					anythingElse()
				}
			default:
				if util.AsciiUppercaseRegex.MatchString(string(nextChar)) {
					currTagToken().tagName += strings.ToLower(string(nextChar))
					tempBuf += string(nextChar)
				} else if util.AsciiLowercaseRegex.MatchString(string(nextChar)) {
					currTagToken().tagName += string(nextChar)
					tempBuf += string(nextChar)
				} else {
					anythingElse()
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-name-state
		case beforeAttributeNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
			case '/', '>', -1:
				t.tkh.Cursor--
				t.state = afterAttributeNameState
			case '=':
				t.parseErrorEncountered(unexpected_equals_sign_before_attribute_name_error)
				currTagToken().attrs = append(currTagToken().attrs, dom.AttrData{
					LocalName: string(nextChar),
					Value:     "",
				})
				attrsToRemove = []int{}
				t.state = attributeNameState
			default:
				currTagToken().attrs = append(currTagToken().attrs, dom.AttrData{
					LocalName: "",
					Value:     "",
				})
				attrsToRemove = []int{}
				t.tkh.Cursor--
				t.state = attributeNameState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-name-state
		case attributeNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ', '/', '>':
				t.tkh.Cursor--
				t.state = afterAttributeNameState
				checkDuplicateAttrName()
			case '=':
				t.state = beforeAttributeValueState
				checkDuplicateAttrName()
			case 0x0000:
				t.parseErrorEncountered(unexpected_null_character_error)
				currAttr().LocalName += string(rune(0xfffd))
			case '"', '\'', '<':
				t.parseErrorEncountered(unexpected_character_in_attribute_name_error)
				fallthrough
			default:
				currAttr().LocalName += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-name-state
		case afterAttributeNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
			case '/':
				t.state = selfClosingStartTagState
			case '=':
				t.state = beforeAttributeValueState
			case '>':
				t.state = dataState
				emitToken(currTk)
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
			default:
				currTagToken().attrs = append(currTagToken().attrs, dom.AttrData{
					LocalName: "",
					Value:     "",
				})
				attrsToRemove = []int{}
				t.tkh.Cursor--
				t.state = attributeNameState

			}
		// https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-value-state
		case beforeAttributeValueState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
			case '"':
				t.state = attributeValueDoubleQuotedState
			case '\'':
				t.state = attributeValueSingleQuotedState
			case '>':
				t.parseErrorEncountered(missing_attribute_value_error)
				t.state = dataState
			default:
				t.tkh.Cursor--
				t.state = attributeValueUnquotedState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(double-quoted)-state
		case attributeValueDoubleQuotedState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '"':
				t.state = afterAttributeValueQuotedState
			case '&':
				returnState = attributeValueDoubleQuotedState
				t.state = characterReferenceState
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				currAttr().Value += string(rune(0xfffd))
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
				return
			default:
				currAttr().Value += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(single-quoted)-state
		case attributeValueSingleQuotedState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\'':
				t.state = afterAttributeValueQuotedState
			case '&':
				returnState = attributeValueSingleQuotedState
				t.state = characterReferenceState
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				currAttr().Value += string(rune(0xfffd))
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
				return
			default:
				currAttr().Value += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(unquoted)-state
		case attributeValueUnquotedState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				t.state = beforeAttributeNameState
			case '&':
				returnState = attributeValueUnquotedState
				t.state = characterReferenceState
			case '>':
				t.state = dataState
				emitToken(currTk)
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				currAttr().Value += string(rune(0xfffd))
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
				return
			case '"', '\'', '<', '=', '`':
				t.parseErrorEncountered(unexpected_character_in_unquoted_attribute_value_error)
				fallthrough
			default:
				currAttr().Value += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-value-(quoted)-state
		case afterAttributeValueQuotedState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				t.state = beforeAttributeNameState
			case '/':
				t.state = selfClosingStartTagState
			case '>':
				t.state = dataState
				emitToken(currTk)
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
				return
			default:
				t.parseErrorEncountered(missing_whitepace_between_attributes_error)
				t.tkh.Cursor--
				t.state = beforeAttributeNameState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#self-closing-start-tag-state
		case selfClosingStartTagState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '>':
				currTagToken().isSelfClosing = true
				t.state = dataState
				emitToken(currTk)
			case -1:
				t.parseErrorEncountered(eof_in_tag_error)
				emitToken(&eofToken{})
				return
			default:
				t.parseErrorEncountered(unexpected_solidus_in_tag_error)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#bogus-comment-state
		case bogusCommentState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '>':
				t.state = dataState
				emitToken(currTk)
			case -1:
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				currCommentToken().data += string(rune(0xfffd))
			default:
				currCommentToken().data += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#markup-declaration-open-state
		case markupDeclarationOpenState:
			if t.tkh.ConsumeStrIfMatches("--", 0) != "" {
				currTk = &commentToken{data: ""}
				t.state = commentStartState
			} else if t.tkh.ConsumeStrIfMatches("DOCTYPE", util.AsciiCaseInsensitive) != "" {
				t.state = doctypeState
			} else if t.tkh.ConsumeStrIfMatches("[CDATA[", 0) != "" {
				panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#markup-declaration-open-state]")
			} else {
				t.parseErrorEncountered(incorrectly_opened_comment_error)
				currTk = &commentToken{data: ""}
				t.state = bogusCommentState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-start-state
		case commentStartState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '-':
				t.state = commentStartDashState
			case '>':
				t.parseErrorEncountered(abrupt_closing_of_empty_comment_error)
				t.state = dataState
				emitToken(currTk)
			default:
				t.tkh.Cursor--
				t.state = commentState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-start-dash-state
		case commentStartDashState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '-':
				t.state = commentEndState
			case '>':
				t.parseErrorEncountered(abrupt_closing_of_empty_comment_error)
				t.state = dataState
				emitToken(currTk)
			case -1:
				t.parseErrorEncountered(eof_in_comment_error)
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				currCommentToken().data += "-"
				t.tkh.Cursor--
				t.state = commentState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-state
		case commentState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '<':
				currCommentToken().data += string(nextChar)
				t.state = commentLessThanSignState
			case '-':
				t.state = commentEndDashState
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				currCommentToken().data += string(rune(0xfffd))
			case -1:
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				currCommentToken().data += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-state
		case commentLessThanSignState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '!':
				currCommentToken().data += string(nextChar)
				panic("TODO")
				// t.state = commentLessThanSignBangState
			case '<':
				currCommentToken().data += string(nextChar)
			default:
				t.tkh.Cursor--
				t.state = commentState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-end-dash-state
		case commentEndDashState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '-':
				t.state = commentEndState
			case -1:
				t.parseErrorEncountered(eof_in_comment_error)
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				currCommentToken().data += "-"
				t.tkh.Cursor--
				t.state = commentState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-end-state
		case commentEndState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '>':
				t.state = dataState
				emitToken(currTk)
			case '!':
				panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#comment-end-state]")
				// t.state = commentEndBangState
			case -1:
				t.parseErrorEncountered(eof_in_comment_error)
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				currCommentToken().data += "--"
				t.tkh.Cursor--
				t.state = commentState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#doctype-state
		case doctypeState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				t.state = beforeDoctypeNameState
			case '>':
				t.tkh.Cursor--
				t.state = beforeDoctypeNameState
			case -1:
				t.parseErrorEncountered(eof_in_doctype_error)
				currTk = &doctypeToken{}
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				t.parseErrorEncountered(missing_whitespace_before_doctype_name_error)
				t.tkh.Cursor--
				t.state = beforeDoctypeNameState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#before-doctype-name-state
		case beforeDoctypeNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				v := string(rune(0xfffd))
				currTk = &doctypeToken{
					name: &v,
				}
			case '>':
				t.parseErrorEncountered(missing_doctype_name_error)
				currTk = &doctypeToken{
					forceQuirks: true,
				}
				emitToken(currTk)
			case -1:
				t.parseErrorEncountered(eof_in_doctype_error)
				currTk = &doctypeToken{}
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				if util.AsciiUppercaseRegex.MatchString(string(nextChar)) {
					nextChar = nextChar - 'A' + 'a'
				}
				v := string(nextChar)
				currTk = &doctypeToken{
					name: &v,
				}
				t.state = doctypeNameState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#doctype-name-state
		case doctypeNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
				t.state = afterDoctypeNameState
			case '>':
				t.state = dataState
				emitToken(currTk)
			case 0:
				t.parseErrorEncountered(unexpected_null_character_error)
				*(currDoctypeToken().name) += string(rune(0xfffd))
			case -1:
				t.parseErrorEncountered(eof_in_doctype_error)
				currDoctypeToken().forceQuirks = true
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				if util.AsciiUppercaseRegex.MatchString(string(nextChar)) {
					nextChar = nextChar - 'A' + 'a'
				}
				*(currDoctypeToken().name) += string(nextChar)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state
		case afterDoctypeNameState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '\t', '\n', 0x000c, ' ':
			case '>':
				t.state = dataState
				emitToken(currTk)
			case -1:
				t.parseErrorEncountered(eof_in_doctype_error)
				currDoctypeToken().forceQuirks = true
				emitToken(currTk)
				emitToken(&eofToken{})
				return
			default:
				if t.tkh.ConsumeStrIfMatches("PUBLIC", util.AsciiCaseInsensitive) != "" {
					panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state]")
				} else if t.tkh.ConsumeStrIfMatches("SYSTEM", util.AsciiCaseInsensitive) != "" {
					panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state]")
				} else {
					t.parseErrorEncountered(invalid_character_sequence_after_doctype_name_error)
					currDoctypeToken().forceQuirks = true
					t.tkh.Cursor--
					panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state]")
					// t.state = bogusDoctypeState
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#character-reference-state
		case characterReferenceState:
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case '#':
				tempBuf += string(nextChar)
				t.state = numericCharacterReferenceState
			default:
				if util.AsciiAlphanumericRegex.MatchString(string(nextChar)) {
					t.tkh.Cursor--
					t.state = namedCharacterReferenceState
				} else {
					flushCodepointsConsumedAsCharReference()
					t.tkh.Cursor--
					t.state = returnState
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#named-character-reference-state
		case namedCharacterReferenceState:
			foundName := ""
			for name := range htmlEntities {
				oldCursor := t.tkh.Cursor
				if name[0] != '&' {
					log.Printf("internal warning: key %s in htmlEntities doesn't start with &", name)
					continue
				}
				if t.tkh.ConsumeStrIfMatches(name[1:], 0) != "" {
					if len(foundName) < len(name) {
						foundName = name
					}
					t.tkh.Cursor = oldCursor
				}
			}
			if foundName != "" {
				t.tkh.Cursor += len([]rune(foundName))
				entity := htmlEntities[foundName]
				if isConsumedAsPartOfAttr() &&
					!strings.HasSuffix(foundName, ";") &&
					(t.tkh.PeekChar() == '=' || util.AsciiAlphanumericRegex.MatchString(string(t.tkh.PeekChar()))) {
					flushCodepointsConsumedAsCharReference()
					t.state = returnState
				} else {
					if !strings.HasSuffix(foundName, ";") {
						t.parseErrorEncountered(missing_semicolon_after_character_reference_error)
					}
					tempBuf = entity.characters
					flushCodepointsConsumedAsCharReference()
					t.state = returnState
				}
			} else {
				flushCodepointsConsumedAsCharReference()
				t.state = returnState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-state
		case numericCharacterReferenceState:
			characterReferenceCode = 0
			nextChar := t.tkh.ConsumeChar()
			switch nextChar {
			case 'X', 'x':
				tempBuf += string(nextChar)
				t.state = hexadecimalCharacterReferenceStartState
			default:
				t.tkh.Cursor--
				t.state = decimalCharacterReferenceStartState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#hexadecimal-character-reference-start-state
		case hexadecimalCharacterReferenceStartState:
			nextChar := t.tkh.ConsumeChar()
			if util.AsciiHexDigitRegex.MatchString(string(nextChar)) {
				t.tkh.Cursor--
				t.state = hexadecimalCharacterReferenceState
			} else {
				t.parseErrorEncountered(absence_of_digits_in_numeric_character_reference_error)
				t.tkh.Cursor--
				t.state = returnState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-start-state
		case decimalCharacterReferenceStartState:
			nextChar := t.tkh.ConsumeChar()
			if util.AsciiDigitRegex.MatchString(string(nextChar)) {
				t.tkh.Cursor--
				t.state = decimalCharacterReferenceState
			} else {
				t.parseErrorEncountered(absence_of_digits_in_numeric_character_reference_error)
				t.tkh.Cursor--
				t.state = returnState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#hexadecimal-character-reference-state
		case hexadecimalCharacterReferenceState:
			nextChar := t.tkh.ConsumeChar()
			if util.AsciiDigitRegex.MatchString(string(nextChar)) {
				characterReferenceCode = (characterReferenceCode * 16) + int(nextChar-'0')
			} else if util.AsciiUpperHexDigitRegex.MatchString(string(nextChar)) {
				characterReferenceCode = (characterReferenceCode * 16) + int(nextChar-'A'+10)
			} else if util.AsciiLowerHexDigitRegex.MatchString(string(nextChar)) {
				characterReferenceCode = (characterReferenceCode * 16) + int(nextChar-'a'+10)
			} else if nextChar == ';' {
				t.state = numericCharacterReferenceEndState
			} else {
				t.parseErrorEncountered(missing_semicolon_after_character_reference_error)
				t.tkh.Cursor--
				t.state = numericCharacterReferenceEndState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-state
		case decimalCharacterReferenceState:
			nextChar := t.tkh.ConsumeChar()
			if util.AsciiDigitRegex.MatchString(string(nextChar)) {
				characterReferenceCode = (characterReferenceCode * 10) + int(nextChar-'0')
			} else if nextChar == ';' {
				t.state = numericCharacterReferenceEndState
			} else {
				t.parseErrorEncountered(missing_semicolon_after_character_reference_error)
				t.tkh.Cursor--
				t.state = numericCharacterReferenceEndState
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-end-state
		case numericCharacterReferenceEndState:
			chr := rune(characterReferenceCode)
			if chr == 0x0000 {
				t.parseErrorEncountered(null_character_reference_error)
				chr = 0xfffd
			} else if 0x10ffff < chr {
				t.parseErrorEncountered(character_reference_outside_unicode_range_error)
				chr = 0xfffd
			} else if util.IsSurrogateChar(chr) {
				t.parseErrorEncountered(surrogate_character_reference_error)
				chr = 0xfffd
			} else if util.IsNoncharacter(chr) {
				t.parseErrorEncountered(noncharacter_reference_error)
			} else if (chr == 0x0d) || (util.IsControlChar(chr) && !util.IsAsciiWhitespace(chr)) {
				t.parseErrorEncountered(control_character_reference_error)
				replaceTable := []struct {
					from rune
					to   rune
				}{
					{0x80, 0x20ac}, {0x82, 0x201a}, {0x83, 0x0192}, {0x84, 0x201e},
					{0x85, 0x2026}, {0x86, 0x2020}, {0x87, 0x2021}, {0x88, 0x02c6},
					{0x89, 0x2030}, {0x8a, 0x0160}, {0x8b, 0x2039}, {0x8c, 0x0152},
					{0x8e, 0x017d}, {0x91, 0x2018}, {0x92, 0x2019}, {0x93, 0x201c},
					{0x94, 0x201d}, {0x95, 0x2022}, {0x96, 0x2013}, {0x97, 0x2014},
					{0x98, 0x02dc}, {0x99, 0x2122}, {0x9a, 0x0161}, {0x9b, 0x203a},
					{0x9c, 0x0153}, {0x9e, 0x017e}, {0x9f, 0x0178},
				}
				for _, item := range replaceTable {
					if item.from == chr {
						chr = item.to
						break
					}
				}
			}
			tempBuf = ""
			tempBuf += string(chr)
			flushCodepointsConsumedAsCharReference()
			t.state = returnState
		default:
			fmt.Printf("Unhandled state %v", t.state)
			panic("unhandled state")
		}
	}

}
