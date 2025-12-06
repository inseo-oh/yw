// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

// Package csssyntax implements [CSS Syntax Module Level 3], as well as parsers
// for various properties.
//
// [CSS Syntax Module Level 3]: https://www.w3.org/TR/css-syntax-3/
package csssyntax

//go:generate go run ./gen

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/encoding"
	"github.com/inseo-oh/yw/util"
)

type tokenCommon struct{ cursorFrom, cursorTo int }
type token interface {
	tokenCursorFrom() int
	tokenCursorTo() int
	tokenType() tokenType
	String() string
}

type tokenType uint8

const (
	tokenTypeEof tokenType = iota // TODO: Remove this
	tokenTypeWhitespace
	tokenTypeLeftParen
	tokenTypeRightParen
	tokenTypeComma
	tokenTypeColon
	tokenTypeSemicolon
	tokenTypeLeftSquareBracket
	tokenTypeRightSquareBracket
	tokenTypeLeftCurlyBracket
	tokenTypeRightCurlyBracket
	tokenTypeCdo
	tokenTypeCdc
	tokenTypeBadString
	tokenTypeBadUrl
	tokenTypeNumber
	tokenTypePercentage
	tokenTypeDimension
	tokenTypeString
	tokenTypeUrl
	tokenTypeAtKeyword
	tokenTypeFuncKeyword
	tokenTypeIdent
	tokenTypeHash
	tokenTypeDelim
	// High-level objects ------------------------------------------------------
	tokenTypeSimpleBlock
	tokenTypeAstFunc
	tokenTypeQualifiedRule
	tokenTypeAtRule
	tokenTypeDeclaration
)

func (typ tokenType) String() string {
	switch typ {
	case tokenTypeWhitespace:
		return "whitespace"
	case tokenTypeLeftParen:
		return "left-paren"
	case tokenTypeRightParen:
		return "right-paren"
	case tokenTypeComma:
		return "comma"
	case tokenTypeColon:
		return "colon"
	case tokenTypeSemicolon:
		return "semicolon"
	case tokenTypeLeftSquareBracket:
		return "left-square-bracket"
	case tokenTypeRightSquareBracket:
		return "right-square-bracket"
	case tokenTypeLeftCurlyBracket:
		return "left-curly-bracket"
	case tokenTypeRightCurlyBracket:
		return "right-curly-bracket"
	case tokenTypeCdo:
		return "cdo"
	case tokenTypeCdc:
		return "cdc"
	case tokenTypeBadString:
		return "bad-string"
	case tokenTypeBadUrl:
		return "bad-url"
	case tokenTypeNumber:
		return "number"
	case tokenTypePercentage:
		return "percentage"
	case tokenTypeDimension:
		return "dimension"
	case tokenTypeString:
		return "string"
	case tokenTypeUrl:
		return "url"
	case tokenTypeAtKeyword:
		return "at-keyword"
	case tokenTypeFuncKeyword:
		return "function"
	case tokenTypeIdent:
		return "ident"
	case tokenTypeHash:
		return "hash"
	case tokenTypeDelim:
		return "delim"
	case tokenTypeSimpleBlock:
		return "simple-block"
	case tokenTypeAstFunc:
		return "function"
	case tokenTypeQualifiedRule:
		return "qualified-rule"
	case tokenTypeAtRule:
		return "at-rule"
	case tokenTypeDeclaration:
		return "declaration"
	}
	return fmt.Sprintf("<bad tokenType %d>", typ)
}
func (t tokenCommon) tokenCursorFrom() int {
	return t.cursorFrom
}
func (t tokenCommon) tokenCursorTo() int {
	return t.cursorTo
}

type simpleToken struct {
	tokenCommon
	tp tokenType
}

func (t simpleToken) tokenType() tokenType { return t.tp }
func (t simpleToken) String() string {
	switch t.tp {
	case tokenTypeWhitespace:
		return " "
	case tokenTypeLeftParen:
		return "("
	case tokenTypeRightParen:
		return ")"
	case tokenTypeComma:
		return ","
	case tokenTypeColon:
		return ":"
	case tokenTypeSemicolon:
		return ";"
	case tokenTypeLeftSquareBracket:
		return "["
	case tokenTypeRightSquareBracket:
		return "]"
	case tokenTypeLeftCurlyBracket:
		return "{"
	case tokenTypeRightCurlyBracket:
		return "}"
	case tokenTypeCdo:
		return "<!--"
	case tokenTypeCdc:
		return "-->"
	case tokenTypeBadString:
		return "/*bad-string*/"
	case tokenTypeBadUrl:
		return "/*bad-url*/"
	}
	return fmt.Sprintf("<bad simpleToken type %v>", t.tp)
}

type numberToken struct {
	tokenCommon
	value css.Num
}

func (t numberToken) tokenType() tokenType { return tokenTypeNumber }
func (t numberToken) String() string       { return fmt.Sprintf("%v", t.value) }

type percentageToken struct {
	tokenCommon
	value css.Num
}

func (t percentageToken) tokenType() tokenType { return tokenTypePercentage }
func (t percentageToken) String() string       { return fmt.Sprintf("%v%%", t.value) }

type dimensionToken struct {
	tokenCommon
	value css.Num
	unit  string
}

func (t dimensionToken) tokenType() tokenType { return tokenTypeDimension }
func (t dimensionToken) String() string       { return fmt.Sprintf("%v%s", t.value, t.unit) }

type stringToken struct {
	tokenCommon
	value string
}

func (t stringToken) tokenType() tokenType { return tokenTypeString }
func (t stringToken) String() string       { return strconv.Quote(t.value) }

type urlToken struct {
	tokenCommon
	value string
}

func (t urlToken) tokenType() tokenType { return tokenTypeUrl }
func (t urlToken) String() string       { return fmt.Sprintf("url(%s)", t.value) }

type atKeywordToken struct {
	tokenCommon
	name string
}

func (t atKeywordToken) tokenType() tokenType { return tokenTypeAtKeyword }
func (t atKeywordToken) String() string       { return fmt.Sprintf("@%s", t.name) }

type funcKeywordToken struct {
	tokenCommon
	value string
}

func (t funcKeywordToken) tokenType() tokenType { return tokenTypeFuncKeyword }
func (t funcKeywordToken) String() string       { return fmt.Sprintf("%s(", t.value) }

type identToken struct {
	tokenCommon
	value string
}

func (t identToken) tokenType() tokenType { return tokenTypeIdent }
func (t identToken) String() string       { return t.value }

type hashToken struct {
	tokenCommon
	tp    hashTokenType
	value string
}

func (t hashToken) tokenType() tokenType { return tokenTypeHash }
func (t hashToken) String() string       { return fmt.Sprintf("#%s/*%s*/", t.value, t.tp) }

type hashTokenType uint8

const (
	hashTokenTypeId hashTokenType = iota
	hashTokenTypeUnrestricted
)

func (typ hashTokenType) String() string {
	switch typ {
	case hashTokenTypeId:
		return "id"
	case hashTokenTypeUnrestricted:
		return "unrestricted"
	}
	return fmt.Sprintf("<bad hashTokenType %d>", typ)
}

type delimToken struct {
	tokenCommon
	value rune
}

func (t delimToken) tokenType() tokenType { return tokenTypeDelim }
func (t delimToken) String() string       { return fmt.Sprintf("%c", t.value) }

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-filter-code-points
func filterCodepoints(src string) string {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")
	src = strings.ReplaceAll(src, "\u000c", "\n")
	return src
}

func tokenize(bytes []byte, sourceName string) (res tokenStream, err error) {
	src := decodeBytes(bytes)
	src = filterCodepoints(src)
	tkh := util.TokenizerHelper{Str: []rune(src), SourceName: sourceName}

	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#ident-start-code-point
	isIdentStartCodepoint := func(chr rune) bool {
		return util.AsciiAlphaRegex.MatchString(string(chr)) ||
			0x80 <= chr ||
			chr == '_'
	}
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#ident-code-point
	isIdentCodepoint := func(chr rune) bool {
		return isIdentStartCodepoint(chr) ||
			util.AsciiDigitRegex.MatchString(string(chr)) ||
			chr == '-'
	}
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#check-if-three-code-points-would-start-an-ident-sequence
	isValidIdentStartSequence := func(s string) bool {
		cps := []rune(s)
		if len(cps) == 0 {
			return false
		}
		if isIdentStartCodepoint(cps[0]) {
			return true
		}
		switch cps[0] {
		case '-':
			if 1 < len(cps) && isIdentCodepoint(cps[1]) {
				return true
			}
			if 2 < len(cps) && cps[1] == '\\' && cps[2] != '\n' {
				return true
			}
			return false
		case '\\':
			if 1 < len(cps) && cps[1] != '\n' {
				return true
			}
			return false
		}
		return false
	}

	consumeComments := func() {
		endFound := false
		for !tkh.IsEof() {
			if tkh.ConsumeStrIfMatches("/*", 0) == "" {
				return
			}
			for !tkh.IsEof() {
				if tkh.ConsumeStrIfMatches("*/", 0) != "" {
					endFound = true
					break
				}
				tkh.ConsumeChar()
			}
			if endFound {
				continue
			}
			// PARSE ERROR: Reached EOF without closing the comment.
			return
		}
	}

	// Returns nil if not found
	consumeNumber := func() *css.Num {
		startCursor := tkh.Cursor
		haveIntegerPart := false
		haveFractionalPart := false
		res := css.Num{}

		// Note that we don't parse the number directly - We only check if it's a valid CSS number.
		// Rest of the job is handled by the standard library.

		// Sign ----------------------------------------------------------------

		// Integer part --------------------------------------------------------
		for !tkh.IsEof() {
			tempChar := tkh.PeekChar()
			if util.AsciiDigitRegex.MatchString(string(tempChar)) {
				tkh.ConsumeChar()
				haveIntegerPart = true
			} else {
				break
			}
		}
		// Decimal point -------------------------------------------------------
		oldCursor := tkh.Cursor
		if tkh.ConsumeCharIfMatches('.') != -1 {
			// Fractional part -------------------------------------------------
			digitCount := 0

			for !tkh.IsEof() {
				tempChar := tkh.PeekChar()
				if util.AsciiDigitRegex.MatchString(string(tempChar)) {
					tkh.ConsumeChar()
					digitCount++
				} else {
					break
				}
			}
			if !haveIntegerPart && digitCount == 0 {
				tkh.Cursor = oldCursor
				return nil
			}
			res.Type = css.NumTypeFloat
			haveFractionalPart = true
		}

		if !haveIntegerPart && !haveFractionalPart {
			// We have invalid number
			tkh.Cursor = startCursor
			return nil
		}

		// Exponent indicator --------------------------------------------------
		oldCursor = tkh.Cursor
		if tkh.ConsumeCharIfMatchesOneOf("eE") != -1 {
			digitCount := 0

			// Exponent sign ---------------------------------------------------
			tkh.ConsumeCharIfMatchesOneOf("+-")

			// Exponent --------------------------------------------------------
			for !tkh.IsEof() {
				tempChar := tkh.PeekChar()
				if util.AsciiDigitRegex.MatchString(string(tempChar)) {
					tkh.ConsumeChar()
				} else {
					break
				}
			}
			if digitCount == 0 {
				tkh.Cursor = oldCursor
			}
		}

		endCursor := tkh.Cursor

		// Now we parse the number ---------------------------------------------
		tempBuf := strings.Builder{}
		tkh.Cursor = startCursor
		for tkh.Cursor < endCursor {
			tempBuf.WriteRune(tkh.ConsumeChar())
		}
		// TODO: Check for range errors
		if res.Type == css.NumTypeFloat {
			v, err := strconv.ParseFloat(tempBuf.String(), 64)
			if err != nil {
				log.Panic(err)
			}
			res.Value = v
		} else {
			v, err := strconv.ParseInt(tempBuf.String(), 10, 64)
			if err != nil {
				log.Panic(err)
			}
			res.Value = v
		}

		return &res
	}

	// Returns nil if not found
	consumeEscapedCodepoint := func() *rune {
		oldCursor := tkh.Cursor
		if tkh.ConsumeCharIfMatches('\\') == -1 {
			return nil
		}
		isHexDigit := false
		hexDigitVal := 0
		hexDigitCount := 0

		if tkh.IsEof() {
			// PARSE ERROR: Unexpected EOF
			cp := rune(0xfffd)
			return &cp
		}
		if tkh.ConsumeCharIfMatches('\n') != -1 {
			tkh.Cursor = oldCursor
			return nil
		}
		// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-an-escaped-code-point
		for !tkh.IsEof() && hexDigitCount < 6 {
			tempChar := tkh.PeekChar()
			digit := 0
			if util.AsciiDigitRegex.MatchString(string(tempChar)) {
				digit = int(tempChar - '0')
			} else if util.AsciiLowerHexDigitRegex.MatchString(string(tempChar)) {
				digit = int(tempChar - 'a' + 10)
			} else if util.AsciiUpperHexDigitRegex.MatchString(string(tempChar)) {
				digit = int(tempChar - 'A' + 10)
			} else {
				break
			}
			tkh.ConsumeChar()
			hexDigitVal = hexDigitVal*16 + digit
			isHexDigit = true
			hexDigitCount++
		}
		var out rune
		if isHexDigit {
			out = rune(hexDigitVal)
		} else {
			out = tkh.ConsumeChar()
		}
		return &out
	}

	// Returns nil if not found
	consumeIdentSequence := func(mustStartWithIdentStart bool) *string {
		sb := strings.Builder{}

		for !tkh.IsEof() {
			oldCursor := tkh.Cursor

			var resultChr rune
			if temp := consumeEscapedCodepoint(); temp == nil {
				resultChr = tkh.ConsumeChar()
			} else {
				resultChr = *temp
			}
			if isIdentStartCodepoint(resultChr) ||
				((sb.Len() != 0 || !mustStartWithIdentStart) && isIdentCodepoint(resultChr)) {
				sb.WriteRune(resultChr)
			} else {
				tkh.Cursor = oldCursor
				break
			}
		}

		if sb.Len() == 0 {
			return nil
		}
		return util.MakeStrPtr(sb.String())
	}

	// Returns nil if not found
	consumeWhitespaceToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		found := false
		for !tkh.IsEof() {
			if tkh.ConsumeCharIfMatchesOneOf(" \t\n") == -1 {
				break
			}
			found = true
		}
		if !found {
			return nil, errors.New("expected a whitespace")
		}
		return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeWhitespace}, nil
	}

	// Returns nil if not found
	consumeSimpleToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		switch tkh.ConsumeCharIfMatchesOneOf("(),:;[]{}") {
		case '(':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeLeftParen}, nil
		case ')':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeRightParen}, nil
		case ',':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeComma}, nil
		case ':':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeColon}, nil
		case ';':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeSemicolon}, nil
		case '[':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeLeftSquareBracket}, nil
		case ']':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeRightSquareBracket}, nil
		case '{':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeLeftCurlyBracket}, nil
		case '}':
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeRightCurlyBracket}, nil
		case -1:
		default:
			panic("unreachable")
		}
		switch tkh.ConsumeStrIfMatchesOneOf([]string{"<!--", "-->"}, 0) {
		case "<!--":
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeCdo}, nil
		case "-->":
			return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeCdc}, nil
		case "":
		default:
			panic("unreachable")
		}
		return nil, errors.New("expected something") // TODO: Improve the error message?
	}

	// Returns nil if not found
	consumeStringToken := func() (res token, err error) {
		var endingChar rune
		sb := strings.Builder{}

		switch tkh.ConsumeCharIfMatchesOneOf("\"'") {
		case '"':
			endingChar = '"'
		case '\'':
			endingChar = '\''
		default:
			return nil, errors.New("expected string")
		}

	loop:
		for !tkh.IsEof() {
			switch tkh.ConsumeCharIfMatchesOneOf(string(endingChar) + "\n") {
			case endingChar:
				break loop
			case '\n':
				// PARSE ERROR: Unexpected newline -- Move the cursor back and exit
				tkh.Cursor--
				break loop
			default:
				if tkh.IsEof() {
					// PARSE ERROR: Unexpected EOF
					break
				} else if tkh.ConsumeCharIfMatchesOneOf("\\\n") != -1 {
					continue
				}
				var resultChr rune
				if temp := consumeEscapedCodepoint(); temp != nil {
					resultChr = *temp
				} else {
					resultChr = tkh.ConsumeChar()
				}
				sb.WriteRune(resultChr)
			}
		}
		return stringToken{
			value: sb.String(),
		}, nil
	}

	// Returns nil if not found
	consumeHashToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		if tkh.ConsumeCharIfMatches('#') == -1 {
			return nil, errors.New("expected '#'")
		}
		var hashValue string
		if temp := consumeIdentSequence(false); temp != nil {
			hashValue = *temp
		} else {
			return nil, errors.New("expected identifier after '#'")
		}
		if len(hashValue) == 0 {
			return nil, errors.New("expected identifier after '#'")
		}
		var subtype hashTokenType
		if isValidIdentStartSequence(hashValue) {
			subtype = hashTokenTypeId
		} else {
			subtype = hashTokenTypeUnrestricted
		}
		return hashToken{tokenCommon{cursorFrom, tkh.Cursor}, subtype, hashValue}, nil
	}

	// Returns nil if not found
	consumeAtToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		if tkh.ConsumeCharIfMatches('@') == -1 {
			return nil, errors.New("expected '@'")
		}
		var atValue string
		if temp := consumeIdentSequence(true); temp != nil {
			atValue = *temp
		} else {
			return nil, errors.New("expected identifier after '@'")
		}
		if len(atValue) == 0 || !isValidIdentStartSequence(atValue) {
			return nil, errors.New("expected identifier after '@'")
		}
		return atKeywordToken{tokenCommon{cursorFrom, tkh.Cursor}, atValue}, nil
	}

	// Returns nil if not found
	consumeNumericToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		var num css.Num
		if temp := consumeNumber(); temp != nil {
			num = *temp
		} else {
			return nil, errors.New("expected number")
		}
		oldCursor := tkh.Cursor

		if ident := consumeIdentSequence(true); ident != nil {
			if isValidIdentStartSequence(*ident) {
				return dimensionToken{tokenCommon{cursorFrom, tkh.Cursor}, num, *ident}, nil
			} else {
				tkh.Cursor = oldCursor
			}
		}
		if tkh.ConsumeCharIfMatches('%') != -1 {
			return percentageToken{tokenCommon{cursorFrom, tkh.Cursor}, num}, nil
		} else {
			return numberToken{tokenCommon{cursorFrom, tkh.Cursor}, num}, nil
		}
	}

	consumeRemnantsOfBadUrl := func() {
		for !tkh.IsEof() {
			if tkh.ConsumeCharIfMatches(')') != -1 {
				break
			}
			if consumeEscapedCodepoint() == nil {
				tkh.ConsumeChar()
			}
		}
	}

	// Returns nil if not found
	consumeIdentLikeToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		var ident string

		if temp := consumeIdentSequence(true); temp != nil {
			ident = *temp
		} else {
			return nil, errors.New("expected function, url, or identifier")
		}
		if util.ToAsciiLowercase(ident) == "url" && tkh.ConsumeCharIfMatches('(') != -1 {
			for tkh.ConsumeStrIfMatches("  ", 0) != "" {
			}
			oldCursor := tkh.Cursor
			if tkh.ConsumeCharIfMatchesOneOf("\"'") != -1 ||
				tkh.ConsumeStrIfMatches(" \"", 0) != "" ||
				tkh.ConsumeStrIfMatches(" '", 0) != "" {
				// Function token ----------------------------------------------
				tkh.Cursor = oldCursor
				return funcKeywordToken{tokenCommon{cursorFrom, tkh.Cursor}, ident}, nil
			} else {
				// URL token ---------------------------------------------------
				consumeWhitespaceToken()
				urlSb := strings.Builder{}
			loop:
				for {
					if tkh.IsEof() {
						// PARSE ERROR: Unexpected EOF
						break loop
					}
					switch tkh.ConsumeCharIfMatchesOneOf(")\"'(") {
					case ')':
						break loop
					case '"':
					case '\'':
					case '(':
						// PARSE ERROR: Unexpected character
						consumeRemnantsOfBadUrl()
						return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeBadUrl}, nil
					default:
						var escapedChr rune
						if temp := consumeEscapedCodepoint(); temp != nil {
							escapedChr = *temp
						} else if tkh.ConsumeCharIfMatches('\\') != -1 {
							// PARSE ERROR: Unexpected character after \
							consumeRemnantsOfBadUrl()
							return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeBadUrl}, nil
						} else {
							escapedChr = tkh.ConsumeChar()
						}
						urlSb.WriteRune(escapedChr)
					}
				}
				return urlToken{tokenCommon{cursorFrom, tkh.Cursor}, urlSb.String()}, nil
			}
		} else if tkh.ConsumeCharIfMatches('(') != -1 {
			return funcKeywordToken{tokenCommon{cursorFrom, tkh.Cursor}, ident}, nil
		} else {
			return identToken{tokenCommon{cursorFrom, tkh.Cursor}, ident}, nil
		}
	}

	consumeToken := func() (res token, err error) {
		cursorFrom := tkh.Cursor
		handlers := []func() (token, error){
			consumeWhitespaceToken,
			consumeStringToken,
			consumeHashToken,
			consumeAtToken,
			consumeSimpleToken,
			consumeNumericToken,
			consumeIdentLikeToken,
		}

		consumeComments()
		for _, h := range handlers {
			res, err := h()
			if err == nil {
				return res, nil
			}
		}
		if tkh.IsEof() {
			return nil, errors.New("reached end of input")
		} else {
			c := tkh.ConsumeChar()
			return delimToken{tokenCommon{cursorFrom, tkh.Cursor}, c}, nil
		}
	}

	tokenList := []token{}
	for {
		tk, err := consumeToken()
		if err != nil {
			if tkh.IsEof() {
				break
			}
			return res, err
		}
		tokenList = append(tokenList, tk)

	}

	return tokenStream{tokens: tokenList, tokenizerHelper: &tkh}, nil
}

type simpleBlockToken struct {
	tokenCommon
	tp   simpleBlockType
	body []token
}

func (t simpleBlockToken) String() string {
	sb := strings.Builder{}
	switch t.tp {
	case simpleBlockTypeCurly:
		sb.WriteRune('{')
	case simpleBlockTypeSquare:
		sb.WriteRune('[')
	case simpleBlockTypeParen:
		sb.WriteRune('(')
	}
	for _, tk := range t.body {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	switch t.tp {
	case simpleBlockTypeCurly:
		sb.WriteRune('}')
	case simpleBlockTypeSquare:
		sb.WriteRune(']')
	case simpleBlockTypeParen:
		sb.WriteRune(')')
	}
	return sb.String()
}

type simpleBlockType uint8

const (
	simpleBlockTypeSquare simpleBlockType = iota
	simpleBlockTypeCurly
	simpleBlockTypeParen
)

func (n simpleBlockToken) tokenType() tokenType {
	return tokenTypeSimpleBlock
}

type astFuncToken struct {
	tokenCommon
	name  string
	value []token
}

func (t astFuncToken) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s(", t.name))
	for _, tk := range t.value {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	sb.WriteString(")")
	return sb.String()
}

func (t astFuncToken) tokenType() tokenType {
	return tokenTypeAstFunc
}

type qualifiedRuleToken struct {
	tokenCommon
	prelude []token
	body    []token
}

func (t qualifiedRuleToken) String() string {
	sb := strings.Builder{}
	for _, tk := range t.prelude {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	sb.WriteString("{")
	for _, tk := range t.body {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	sb.WriteString("}")
	return sb.String()
}
func (t qualifiedRuleToken) tokenType() tokenType {
	return tokenTypeQualifiedRule
}

type atRuleToken struct {
	tokenCommon
	name    string
	prelude []token // NOTE: This is just STUB -- We would want actual parsed value.
	body    []token // NOTE: This is just STUB -- We would want actual parsed value.
}

func (t atRuleToken) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("@%s ", t.name))
	for _, tk := range t.prelude {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	sb.WriteString("{")
	for _, tk := range t.body {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	sb.WriteString("}")
	return sb.String()
}
func (t atRuleToken) tokenType() tokenType {
	return tokenTypeAtRule
}

type declarationToken struct {
	tokenCommon
	name      string
	value     []token
	important bool
}

func (t declarationToken) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s:", t.name))
	for _, tk := range t.value {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	if t.important {
		sb.WriteString("!important")
	}
	return sb.String()
}
func (t declarationToken) tokenType() tokenType {
	return tokenTypeDeclaration
}

type tokenStream struct {
	tokens          []token
	cursor          int
	tokenizerHelper *util.TokenizerHelper //  TokenizerHelper used to tokenize.
}

func (ts *tokenStream) errorHeader() string {
	cursor := min(ts.cursor, len(ts.tokens)-1)
	return ts.tokenizerHelper.ErrorHeader(ts.tokens[cursor].tokenCursorFrom())
}

func (ts *tokenStream) isEnd() bool {
	return len(ts.tokens) <= ts.cursor
}
func (ts *tokenStream) consumeToken() (res token, err error) {
	if ts.isEnd() {
		return nil, fmt.Errorf("%s: reached end of input", ts.errorHeader())
	}
	ts.cursor++
	return ts.tokens[ts.cursor-1], nil
}
func (ts *tokenStream) consumeTokenWith(tp tokenType) (res token, err error) {
	oldCursor := ts.cursor
	tk, err := ts.consumeToken()
	if err != nil {
		return nil, err
	} else if tk.tokenType() != tp {
		// TODO: Describe what token we want in more friendly way.
		ts.cursor = oldCursor
		return nil, fmt.Errorf("%s: expected token with type %v", ts.errorHeader(), tp)
	}
	return tk, nil
}

func (ts *tokenStream) skipWhitespaces() {
	for {
		oldCursor := ts.cursor
		_, err := ts.consumeTokenWith(tokenTypeWhitespace)
		if err != nil {
			ts.cursor = oldCursor
			break
		}
	}
}

func (ts *tokenStream) consumeDelimTokenWith(d rune) error {
	oldCursor := ts.cursor
	tk, err := ts.consumeTokenWith(tokenTypeDelim)
	if err != nil || tk.(delimToken).value != d {
		ts.cursor = oldCursor
		return fmt.Errorf("%s: expected %c", ts.errorHeader(), d)
	}
	return nil
}
func (ts *tokenStream) consumeIdentTokenWith(i string) error {
	oldCursor := ts.cursor
	tk, err := ts.consumeTokenWith(tokenTypeIdent)
	if err != nil || tk.(identToken).value != i {
		ts.cursor = oldCursor
		return fmt.Errorf("%s: expected %s", ts.errorHeader(), i)
	}
	return nil
}
func (ts *tokenStream) consumeSimpleBlockWith(tp simpleBlockType) (res simpleBlockToken, err error) {
	oldCursor := ts.cursor
	tk, err := ts.consumeTokenWith(tokenTypeSimpleBlock)
	if err != nil {
		return res, err
	}
	blk := tk.(simpleBlockToken)
	if blk.tp != tp {
		// TODO: Describe what token we want in more friendly way.
		ts.cursor = oldCursor
		return res, fmt.Errorf("%s: expected simple block with type %v", ts.errorHeader(), tp)
	}
	return blk, nil
}
func (ts *tokenStream) consumeAstFuncWith(name string) (res astFuncToken, err error) {
	oldCursor := ts.cursor
	tk, err := ts.consumeTokenWith(tokenTypeAstFunc)
	if err != nil {
		return res, err
	}
	if tk.(astFuncToken).name != name {
		ts.cursor = oldCursor
		return res, fmt.Errorf("%s: expected function %s()", ts.errorHeader(), name)
	}
	return tk.(astFuncToken), nil
}

func (ts *tokenStream) consumePreservedToken() (res token, err error) {
	oldCursor := ts.cursor
	tk, err := ts.consumeToken()
	if err != nil {
		return nil, err
	}
	switch tk.tokenType() {
	case tokenTypeFuncKeyword,
		tokenTypeLeftCurlyBracket,
		tokenTypeLeftSquareBracket,
		tokenTypeLeftParen:
		ts.cursor = oldCursor
		return nil, fmt.Errorf("%s: non-preserved token found", ts.errorHeader())
	}
	return tk, nil
}

func (ts *tokenStream) consumeSimpleBlock(openTokenType, closeTokenType tokenType) (res simpleBlockToken, err error) {
	resNodes := []token{}
	var blockType simpleBlockType
	switch openTokenType {
	case tokenTypeLeftCurlyBracket:
		blockType = simpleBlockTypeCurly
	case tokenTypeLeftSquareBracket:
		blockType = simpleBlockTypeSquare
	case tokenTypeLeftParen:
		blockType = simpleBlockTypeParen
	default:
		panic("unsupported openTokenType")
	}

	openToken, err := ts.consumeTokenWith(openTokenType)
	if err != nil {
		return res, err
	}
	var closeToken token
	for {
		tempTk, err := ts.consumeComponentValue()
		if err == nil && tempTk.tokenType() == closeTokenType {
			closeToken = tempTk
			break
		} else if err != nil {
			break
		}
		resNodes = append(resNodes, tempTk)
	}
	if util.IsNil(closeToken) {
		return res, fmt.Errorf("%s: expected closing character", ts.errorHeader())
	}
	return simpleBlockToken{
		tokenCommon{openToken.tokenCursorFrom(), closeToken.tokenCursorTo()},
		blockType, resNodes,
	}, nil
}

func (ts *tokenStream) consumeCurlyBlock() (res simpleBlockToken, err error) {
	return ts.consumeSimpleBlock(tokenTypeLeftCurlyBracket, tokenTypeRightCurlyBracket)
}

func (ts *tokenStream) consumeSquareBlock() (res simpleBlockToken, err error) {
	return ts.consumeSimpleBlock(tokenTypeLeftSquareBracket, tokenTypeRightSquareBracket)
}

func (ts *tokenStream) consumeParenBlock() (res simpleBlockToken, err error) {
	return ts.consumeSimpleBlock(tokenTypeLeftParen, tokenTypeRightParen)
}

func (ts *tokenStream) consumeFunc() (res astFuncToken, err error) {
	fnValueNodes := []token{}
	var fnToken funcKeywordToken
	if temp, err := ts.consumeTokenWith(tokenTypeFuncKeyword); err == nil {
		fnToken = temp.(funcKeywordToken)
	} else {
		return res, err
	}
	var closeToken token
	for {
		tempTk, err := ts.consumeComponentValue()
		if err == nil && tempTk.tokenType() == tokenTypeRightParen {
			closeToken = tempTk
			break
		} else if err != nil {
			panic("TODO: Handle error")
		}
		fnValueNodes = append(fnValueNodes, tempTk)
	}

	return astFuncToken{
		tokenCommon{fnToken.tokenCursorFrom(), closeToken.tokenCursorTo()},
		fnToken.value, fnValueNodes,
	}, nil
}

// Returns nil if not found
func (ts *tokenStream) consumeComponentValue() (res token, err error) {
	if res, err := ts.consumeCurlyBlock(); err == nil {
		return res, nil
	}
	if res, err := ts.consumeSquareBlock(); err == nil {
		return res, nil
	}
	if res, err := ts.consumeParenBlock(); err == nil {
		return res, nil
	}
	if res, err := ts.consumeFunc(); err == nil {
		return res, nil
	}
	if res, err := ts.consumePreservedToken(); err == nil {
		return res, nil
	}
	return nil, fmt.Errorf("%s: expected component value", ts.errorHeader())
}

func (ts *tokenStream) consumeQualifiedRule() (res qualifiedRuleToken, err error) {
	oldCursor := ts.cursor
	prelude := []token{}

	for {
		block, err := ts.consumeCurlyBlock()
		if err == nil {
			return qualifiedRuleToken{
				tokenCommon{block.cursorFrom, block.cursorTo},
				prelude,
				block.body,
			}, nil
		} else if ts.isEnd() {
			ts.cursor = oldCursor
			return res, fmt.Errorf("%s: expected qualified rule", ts.errorHeader())
		}
		comp, err := ts.consumeComponentValue()
		if err != nil {
			ts.cursor = oldCursor
			return res, err
		}
		prelude = append(prelude, comp)
	}
}

// Returns nil if not found
func (ts *tokenStream) consumeAtRule() (res atRuleToken, err error) {
	oldCursor := ts.cursor
	prelude := []token{}
	var kwdToken atKeywordToken
	if temp, err := ts.consumeTokenWith(tokenTypeAtKeyword); err == nil {
		kwdToken = temp.(atKeywordToken)
	} else {
		return res, err
	}

	for {
		block, err := ts.consumeCurlyBlock()
		if err == nil {
			return atRuleToken{
				tokenCommon{block.cursorFrom, block.cursorTo},
				kwdToken.name,
				prelude,
				block.body,
			}, nil
		} else if ts.isEnd() {
			ts.cursor = oldCursor
			return res, fmt.Errorf("%s: expected at-rule body", ts.errorHeader())
		}
		comp, err := ts.consumeComponentValue()
		if err != nil {
			panic("TODO: Handle error")
		}
		prelude = append(prelude, comp)
	}
}

// Returns nil if not found
func (ts *tokenStream) consumeDeclaration() (res declarationToken, err error) {
	// <name>  :  contents  !important -----------------------------------------
	var identTk identToken
	if temp, err := ts.consumeTokenWith(tokenTypeIdent); err == nil {
		identTk = temp.(identToken)
	} else {
		return res, err
	}
	declName := identTk.value
	declValue := []token{}
	declIsImportant := false
	// name<  >:  contents  !important -----------------------------------------
	ts.skipWhitespaces()
	// name  <:>  contents  !important -----------------------------------------
	if _, err := ts.consumeTokenWith(tokenTypeColon); err != nil {
		// Parse error
		return res, err
	}
	// name  :<  >contents  !important -----------------------------------------
	ts.skipWhitespaces()

	// name  :  <contents  !important> -----------------------------------------
	for {
		tempTk, err := ts.consumeComponentValue()
		if err != nil {
			break
		}
		declValue = append(declValue, tempTk)
	}
	lastNode := declValue[len(declValue)-1]
	if 2 <= len(declValue) {
		// See if we have !important
		ptk1 := declValue[len(declValue)-2]
		ptk2 := declValue[len(declValue)-1]
		if ptk1.tokenType() == tokenTypeDelim && ptk1.(delimToken).value == '!' &&
			ptk2.tokenType() == tokenTypeIdent && ptk2.(identToken).value == "important" {
			declValue = declValue[:len(declValue)-2]
			declIsImportant = true
		}
	}
	return declarationToken{
		tokenCommon{identTk.cursorFrom, lastNode.tokenCursorTo()},
		declName,
		declValue,
		declIsImportant,
	}, nil
}

func (ts *tokenStream) consumeDeclarationList() []token {
	decls := []token{}

	for {
		oldCursor := ts.cursor
		token, err := ts.consumeToken()
		if err != nil {
			ts.cursor = oldCursor
			break
		} else if token.tokenType() == tokenTypeWhitespace || token.tokenType() == tokenTypeSemicolon {
			continue
		} else if token.tokenType() == tokenTypeAtKeyword {
			ts.cursor--
			rule, err := ts.consumeAtRule()
			if err != nil {
				log.Fatal(err)
			}
			decls = append(decls, rule)
		} else if token.tokenType() == tokenTypeIdent {
			tempStream := tokenStream{tokenizerHelper: ts.tokenizerHelper}
			tempStream.tokens = append(tempStream.tokens, token)
			for {
				token, err = ts.consumeToken()
				oldCursor := ts.cursor
				if err != nil || token.tokenType() == tokenTypeSemicolon {
					ts.cursor = oldCursor
					break
				}
				tempStream.tokens = append(tempStream.tokens, token)
			}
			decl, err := tempStream.consumeDeclaration()
			if err != nil {
				break
			}
			decls = append(decls, decl)
		} else {
			// PARSE ERROR
			for {
				oldCursor := ts.cursor
				token, err = ts.consumeToken()
				if err != nil || token.tokenType() == tokenTypeSemicolon {
					ts.cursor = oldCursor
					break
				}
				ts.cursor = oldCursor
				ts.consumeComponentValue()
			}
		}
	}
	return decls
}
func (ts *tokenStream) consumeStyleBlockContents() []token {
	decls := []token{}
	rules := []qualifiedRuleToken{}

	for {
		oldCursor := ts.cursor
		token, err := ts.consumeToken()
		if err != nil {
			ts.cursor = oldCursor
			break
		} else if token.tokenType() == tokenTypeWhitespace || token.tokenType() == tokenTypeSemicolon {
			continue
		} else if token.tokenType() == tokenTypeAtKeyword {
			ts.cursor--
			rule, err := ts.consumeAtRule()
			if err != nil {
				log.Fatal(err)
			}
			decls = append(decls, rule)
		} else if token.tokenType() == tokenTypeIdent {
			tempStream := tokenStream{tokenizerHelper: ts.tokenizerHelper}
			tempStream.tokens = append(tempStream.tokens, token)
			for {
				oldCursor := ts.cursor
				token, err = ts.consumeToken()
				if err != nil || token.tokenType() == tokenTypeSemicolon {
					ts.cursor = oldCursor
					break
				}
				tempStream.tokens = append(tempStream.tokens, token)
			}
			decl, err := tempStream.consumeDeclaration()
			if err != nil {
				break
			}
			decls = append(decls, decl)
		} else if token.tokenType() == tokenTypeDelim && token.(delimToken).value == '&' {
			ts.cursor--
			if rule, err := ts.consumeQualifiedRule(); err == nil {
				rules = append(rules, rule)
			}
		} else {
			// PARSE ERROR
			for {
				oldCursor := ts.cursor
				token, err = ts.consumeToken()
				if err != nil || token.tokenType() == tokenTypeSemicolon {
					break
				}
				ts.cursor = oldCursor
				ts.consumeComponentValue()
			}
		}
	}
	for _, rule := range rules {
		decls = append(decls, rule)
	}
	return decls
}
func (ts *tokenStream) consumeListOfRules(topLevelFlag bool) []token {
	rules := []token{}

	for {
		oldCursor := ts.cursor
		token, err := ts.consumeToken()
		if err != nil {
			ts.cursor = oldCursor
			break
		} else if token.tokenType() == tokenTypeWhitespace || token.tokenType() == tokenTypeSemicolon {
			continue
		} else if token.tokenType() == tokenTypeCdo || token.tokenType() == tokenTypeCdc {
			if topLevelFlag {
				continue
			}
			ts.cursor--
			if rule, err := ts.consumeQualifiedRule(); err == nil {
				rules = append(rules, rule)
			}
		} else if token.tokenType() == tokenTypeAtKeyword {
			ts.cursor--
			rule, err := ts.consumeAtRule()
			if err != nil {
				log.Fatal(err)
			}
			rules = append(rules, rule)
		} else {
			ts.cursor--
			if rule, err := ts.consumeQualifiedRule(); err == nil {
				rules = append(rules, rule)
			} else {
				break
			}
		}
	}
	return rules
}

func parseCommaSeparatedListOfComponentValues(ts tokenStream) [][]token {
	valueLists := [][]token{}
	tempList := []token{}

	for {
		value, err := ts.consumeComponentValue()
		if err != nil || value.tokenType() == tokenTypeComma {
			valueLists = append(valueLists, tempList)
			tempList = tempList[:0]
			if value.tokenType() != tokenTypeComma {
				break
			}
			continue
		}
		tempList = append(tempList, value)
	}
	return valueLists
}
func parseListOfComponentValues(ts *tokenStream) []token {
	tempList := []token{}

	for {
		value, err := ts.consumeComponentValue()
		if err != nil || value.tokenType() == tokenTypeComma {
			break
		}
		tempList = append(tempList, value)
	}
	return tempList
}

func parseComponentValue(ts *tokenStream) token {
	ts.skipWhitespaces()
	if !ts.isEnd() {
		panic("TODO: syntax error: expected component value")
	}
	value, err := ts.consumeComponentValue()
	if err != nil {
		panic("TODO: Handle error")
	}
	ts.skipWhitespaces()
	if ts.isEnd() {
		panic("TODO: syntax error: expected eof")
	}
	return value
}

func parseListOfDeclarations(ts *tokenStream) []token {
	return ts.consumeDeclarationList()
}
func parseStyleBlockContents(ts *tokenStream) []token {
	return ts.consumeStyleBlockContents()
}
func parseDeclaration(ts *tokenStream) declarationToken {
	ts.skipWhitespaces()
	if _, err := ts.consumeTokenWith(tokenTypeIdent); err != nil {
		panic("TODO: syntax error: expected identifier")
	}
	ts.cursor--
	node, err := ts.consumeDeclaration()
	if err != nil {
		panic("TODO: syntax error: expected declaration")
	}
	return node
}
func parseRule(ts *tokenStream) token {
	ts.skipWhitespaces()
	var res token
	res, err := ts.consumeAtRule()
	if err != nil {
		res, err = ts.consumeQualifiedRule()
	}
	if err != nil {
		panic("TODO: syntax error: expected at-rule or qualified-rule")
	}
	if ts.isEnd() {
		panic("TODO: syntax error: expected eof")
	}
	return res
}
func parseListOfRules(ts *tokenStream) []token {
	return ts.consumeListOfRules(false)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#typedef-declaration-value
func (ts *tokenStream) _consumeDeclarationValue(anyValue bool) []token {
	res := []token{}
	openBlockTokens := []tokenType{}
	for {
		oldCursor := ts.cursor
		tk, err := ts.consumeToken()
		if err != nil {
			ts.cursor = oldCursor
			break
		}
		tkType := tk.tokenType()
		if tkType == tokenTypeBadString ||
			tkType == tokenTypeBadUrl ||
			// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#typedef-any-value
			(!anyValue && (tkType == tokenTypeSemicolon ||
				(tkType == tokenTypeDelim && tk.(delimToken).value == '!'))) {
			ts.cursor = oldCursor
			break
		}
		// If we have block opening token, push it to the stack.
		if tkType == tokenTypeLeftParen ||
			tkType == tokenTypeLeftSquareBracket ||
			tkType == tokenTypeLeftCurlyBracket {
			openBlockTokens = append(openBlockTokens, tkType)
		}
		// If we have block closing token, see if we have unmatched token.
		if (tkType == tokenTypeRightParen &&
			(len(openBlockTokens) == 0 ||
				openBlockTokens[len(openBlockTokens)-1] != tokenTypeLeftParen)) ||
			(tkType == tokenTypeRightSquareBracket &&
				(len(openBlockTokens) == 0 ||
					openBlockTokens[len(openBlockTokens)-1] != tokenTypeLeftSquareBracket)) ||
			(tkType == tokenTypeRightCurlyBracket &&
				(len(openBlockTokens) == 0 ||
					openBlockTokens[len(openBlockTokens)-1] != tokenTypeLeftCurlyBracket)) {
			ts.cursor = oldCursor
			break
		}
		res = append(res, tk)
	}
	if len(res) == 0 {
		return nil
	}
	return res
}

func (ts *tokenStream) consumeDeclarationValue() []token {
	return ts._consumeDeclarationValue(false)
}
func (ts *tokenStream) consumeAnyValue() []token {
	return ts._consumeDeclarationValue(true)
}

// This can be used to parse where a CSS syntax can be repeated separated by comma.
// If maxRepeats is 0, repeat count is unlimited.
//
// https://www.w3.org/TR/css-values-4/#mult-comma
func parseCommaSeparatedRepeation[T any](ts *tokenStream, maxRepeats int, description string, parser func(ts *tokenStream) (T, error)) ([]T, error) {
	res := []T{}
	for {
		oldCursor := ts.cursor
		x, err := parser(ts)
		if err != nil {
			if len(res) != 0 {
				// We encountered an error after ','
				return nil, err
			}
			ts.cursor = oldCursor
			break
		}
		// FIXME: Remove this when we know it's safe to remove this check.
		if util.IsNil(x) {
			panic("callback returned a nil value")
		}
		res = append(res, x)
		if maxRepeats != 0 && maxRepeats <= len(res) {
			break
		}
		ts.skipWhitespaces()
		if _, err := ts.consumeTokenWith(tokenTypeComma); err != nil {
			ts.cursor = oldCursor
			break
		}
		ts.skipWhitespaces()
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("%s: expected %s value", ts.errorHeader(), description)
	}
	return res, nil
}

// This can be used to parse where a CSS syntax can be repeated multiple times.
// If maxRepeats is 0, repeat count is unlimited.
//
// https://www.w3.org/TR/css-values-4/#mult-num-range
func parseRepeation[T any](ts *tokenStream, maxRepeats int, description string, parser func(ts *tokenStream) (T, error)) ([]T, error) {
	res := []T{}
	for {
		x, err := parser(ts)
		if err != nil {
			break
		}
		// FIXME: Remove this when we know it's safe to remove this check.
		if util.IsNil(x) {
			panic("callback returned a nil value")
		}
		res = append(res, x)
		if maxRepeats != 0 && maxRepeats <= len(res) {
			break
		}
		ts.skipWhitespaces()
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("%s: expected %s value", ts.errorHeader(), description)
	}
	return res, nil
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-parse-something-according-to-a-css-grammar
func parse[T any](ts *tokenStream, parser func(ts *tokenStream) (T, error)) (T, error) {
	compValues := parseListOfComponentValues(ts)
	subTs := tokenStream{tokens: compValues, tokenizerHelper: ts.tokenizerHelper}
	return parser(&subTs)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-decode-bytes
func decodeBytes(bytes []byte) string {
	fallback := determineFallbackEncoding(bytes)
	input := encoding.IoQueueFromSlice(bytes)
	output := encoding.IoQueueFromSlice[rune](nil)
	encoding.Decode(&input, fallback, &output)
	return string(encoding.IoQueueToSlice[rune](output))
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#determine-the-fallback-encoding
func determineFallbackEncoding(bytes []byte) encoding.Type {
	// Check if HTTP or equivalent protocol provides an encoding label ---------
	// TODO

	// Check '@charset "xxx";' byte sequence -----------------------------------
	if 1024 < len(bytes) {
		bytes = bytes[:1024]
	}
	// NOTE: Below sequence of bytes are '@charset "' in ASCII
	if len(bytes) < 10 || !slices.Equal([]byte{0x40, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x20, 0x22}, bytes[:10]) {
		remainingBytes := bytes
		foundEnd := false
		encodingName := []rune{}
		for 0 < len(remainingBytes) {
			// NOTE: Below sequence of bytes are '";' in ASCII
			if 2 <= len(bytes) && slices.Equal([]byte{0x22, 0x3b}, bytes[:2]) {
				foundEnd = true
				break
			}
			encodingName = append(encodingName, rune(remainingBytes[0]))
			remainingBytes = remainingBytes[1:]
		}
		if foundEnd {
			enc, err := encoding.GetEncodingFromLabel(string(encodingName))
			if err == nil {
				if enc == encoding.Utf16Be || enc == encoding.Utf16Le {
					// This is not a bug. The standard says to do this.
					return encoding.Utf8
				}
				return enc
			}
		}
	}
	// Check if environment encoding is provided -------------------------------
	// TODO

	// Fallback to UTF-8 -------------------------------------------------------
	return encoding.Utf8
}

// ParseStylesheet parses stylesheet from given input, with optional location.
func ParseStylesheet(input []byte, location *string, sourceOfInput string) (res cssom.Stylesheet, err error) {
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-stylesheets
	ts, err := tokenize(input, sourceOfInput)
	if err != nil {
		return res, err
	}
	stylesheet := cssom.Stylesheet{
		Location: location,
	}
	ruleNodes := ts.consumeListOfRules(true)

	// Parse top-level qualified rules as style rules
	stylesheet.StyleRules = parseStyleRulesFromNodes(ruleNodes, ts.tokenizerHelper)

	return stylesheet, nil
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#style-rules
//
// TODO(ois): Make parseStyleRulesFromNodes receive tokenStream instead of ruleNodes and tkh.
func parseStyleRulesFromNodes(ruleNodes []token, tkh *util.TokenizerHelper) []cssom.StyleRule {
	styleRules := []cssom.StyleRule{}
	printRawRuleNodes := false
	for _, n := range ruleNodes {
		if n.tokenType() != tokenTypeQualifiedRule {
			continue
		}
		qrule := n.(qualifiedRuleToken)
		preludeStream := tokenStream{tokens: qrule.prelude, tokenizerHelper: tkh}
		selectorList, err := preludeStream.parseSelectorList()
		if err != nil {
			// TODO: Report error
			log.Printf("selector parsing error: %v", err)
			continue
		}
		if len(selectorList) == 0 {
			log.Println("FIXME: having no selectors after successfully parsing it isn't possible")
			printRawRuleNodes = true
			continue
		}
		contentsStream := tokenStream{tokens: qrule.body, tokenizerHelper: tkh}
		contents := contentsStream.consumeStyleBlockContents()
		decls := []cssom.Declaration{}
		atRules := []cssom.AtRule{}
		for _, content := range contents {
			if content.tokenType() == tokenTypeDeclaration {
				declNode := content.(declarationToken)

				parseFunc, ok := parseFuncMap[util.ToAsciiLowercase(declNode.name)]
				if !ok {
					log.Printf("unknown property name: %v", declNode.name)
					continue
				}
				innerAs := tokenStream{tokens: declNode.value, tokenizerHelper: tkh}
				innerAs.skipWhitespaces()
				value, err := parseFunc(&innerAs)
				if err != nil {
					log.Printf("bad value for property: %v (%v)", declNode.name, err)
					continue
				}
				innerAs.skipWhitespaces()
				if !innerAs.isEnd() {
					log.Printf("extra junk at the end for property: %v (token list: %v)", declNode.name, innerAs.tokens[innerAs.cursor:])
					continue
				}
				decls = append(decls, cssom.Declaration{Name: declNode.name, Value: value, IsImportant: declNode.important})
			} else if content.tokenType() == tokenTypeAtRule {
				ruleNode := content.(atRuleToken)
				atRules = append(atRules, cssom.AtRule{Name: ruleNode.name, Prelude: ruleNode.prelude, Value: ruleNode.body})
			} else {
				log.Printf("warning: unexpected node with type %v found while parsing style block contents", content.tokenType())
			}
		}

		styleRules = append(styleRules, cssom.StyleRule{SelectorList: selectorList, Declarations: decls, AtRules: atRules})
	}
	if printRawRuleNodes {
		log.Println("=============== BEGIN: Raw rule nodes ===============")
		log.Println(ruleNodes)
		log.Println("=============== END:   Raw rule nodes ===============")
	}
	return styleRules
}
