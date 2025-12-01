//go:generate go run ./gen

// Implementation of the CSS Syntax Module Level 3 (https://www.w3.org/TR/css-syntax-3/)
package csssyntax

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

func tokenize(bytes []byte) ([]token, error) {
	src := decodeBytes(bytes)
	src = filterCodepoints(src)
	tkh := util.TokenizerHelper{Str: []rune(src)}

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
		out := css.Num{}

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
			out.Type = css.NumTypeFloat
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
		if out.Type == css.NumTypeFloat {
			v, err := strconv.ParseFloat(tempBuf.String(), 64)
			if err != nil {
				log.Panic(err)
			}
			out.Value = v
		} else {
			v, err := strconv.ParseInt(tempBuf.String(), 10, 64)
			if err != nil {
				log.Panic(err)
			}
			out.Value = v
		}

		return &out
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
	consumeWhitespaceToken := func() (token, error) {
		cursorFrom := tkh.Cursor
		found := false
		for !tkh.IsEof() {
			if tkh.ConsumeCharIfMatchesOneOf(" \t\n") == -1 {
				break
			}
			found = true
		}
		if !found {
			return nil, nil
		}
		return simpleToken{tokenCommon{cursorFrom, tkh.Cursor}, tokenTypeWhitespace}, nil
	}

	// Returns nil if not found
	consumeSimpleToken := func() (token, error) {
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
		return nil, nil
	}

	// Returns nil if not found
	consumeStringToken := func() (token, error) {
		var endingChar rune
		sb := strings.Builder{}

		switch tkh.ConsumeCharIfMatchesOneOf("\"'") {
		case '"':
			endingChar = '"'
		case '\'':
			endingChar = '\''
		default:
			return nil, nil
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
	consumeHashToken := func() (token, error) {
		cursorFrom := tkh.Cursor
		if tkh.ConsumeCharIfMatches('#') == -1 {
			return nil, nil
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
	consumeAtToken := func() (token, error) {
		cursorFrom := tkh.Cursor
		if tkh.ConsumeCharIfMatches('@') == -1 {
			return nil, nil
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
	consumeNumericToken := func() (token, error) {
		cursorFrom := tkh.Cursor
		var num css.Num
		if temp := consumeNumber(); temp != nil {
			num = *temp
		} else {
			return nil, nil
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
	consumeIdentLikeToken := func() (token, error) {
		cursorFrom := tkh.Cursor
		var ident string

		if temp := consumeIdentSequence(true); temp != nil {
			ident = *temp
		} else {
			return nil, nil
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

	consumeToken := func() (token, error) {
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
			if !util.IsNil(res) {
				return res, nil
			} else if err != nil {
				return nil, err
			}
		}
		if tkh.IsEof() {
			return nil, nil
		} else {
			c := tkh.ConsumeChar()
			return delimToken{tokenCommon{cursorFrom, tkh.Cursor}, c}, nil
		}
	}

	out := []token{}
	for {
		tk, err := consumeToken()
		if util.IsNil(tk) && err != nil {
			return nil, err
		} else if util.IsNil(tk) {
			break
		}
		out = append(out, tk)

	}

	return out, nil
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
	tokens []token
	cursor int
}

func (ts *tokenStream) isEnd() bool {
	return len(ts.tokens) <= ts.cursor
}
func (ts *tokenStream) consumeToken() token {
	if ts.isEnd() {
		return nil
	}
	ts.cursor++
	res := ts.tokens[ts.cursor-1]
	return res
}
func (ts *tokenStream) consumeTokenWith(tp tokenType) token {
	if ts.isEnd() {
		return nil
	}
	oldCursor := ts.cursor
	tk := ts.consumeToken()
	if tk.tokenType() != tp {
		ts.cursor = oldCursor
		return nil
	}
	return tk
}

func (ts *tokenStream) skipWhitespaces() {
	for !util.IsNil(ts.consumeTokenWith(tokenTypeWhitespace)) {
	}
}

func (ts *tokenStream) consumeDelimTokenWith(d rune) *delimToken {
	oldCursor := ts.cursor
	tk := ts.consumeTokenWith(tokenTypeDelim)
	if tk == nil || tk.(delimToken).value != d {
		ts.cursor = oldCursor
		return nil
	}
	v := tk.(delimToken)
	return &v
}
func (ts *tokenStream) consumeIdentTokenWith(i string) bool {
	oldCursor := ts.cursor
	tk := ts.consumeTokenWith(tokenTypeIdent)
	if tk == nil || tk.(identToken).value != i {
		ts.cursor = oldCursor
		return false
	}
	return true
}
func (ts *tokenStream) consumeSimpleBlockWith(tp simpleBlockType) *simpleBlockToken {
	oldCursor := ts.cursor
	n := ts.consumeTokenWith(tokenTypeSimpleBlock)
	if util.IsNil(n) {
		ts.cursor = oldCursor
		return nil
	}
	blk := n.(simpleBlockToken)
	if blk.tp != tp {
		ts.cursor = oldCursor
		return nil
	}
	return &blk
}
func (ts *tokenStream) consumeAstFuncWith(name string) *astFuncToken {
	oldCursor := ts.cursor
	tk := ts.consumeTokenWith(tokenTypeAstFunc)
	if tk == nil || tk.(astFuncToken).name != name {
		ts.cursor = oldCursor
		return nil
	}
	v := tk.(astFuncToken)
	return &v
}

// Returns nil if not found
func (ts *tokenStream) consumePreservedToken() token {
	oldCursor := ts.cursor

	tk := ts.consumeToken()
	if util.IsNil(tk) {
		return nil
	}
	switch tk.tokenType() {
	case tokenTypeFuncKeyword,
		tokenTypeLeftCurlyBracket,
		tokenTypeLeftSquareBracket,
		tokenTypeLeftParen:
		ts.cursor = oldCursor
		return nil
	}
	return tk
}

// Returns nil if not found
func (ts *tokenStream) consumeSimpleBlock(openTokenType, closeTokenType tokenType) *simpleBlockToken {
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

	openToken := ts.consumeTokenWith(openTokenType)
	if util.IsNil(openToken) {
		return nil
	}
	var closeToken token
	for {
		tempTk := ts.consumeComponentValue()
		if util.IsNil(tempTk) || tempTk.tokenType() == closeTokenType {
			closeToken = tempTk
			break
		}
		resNodes = append(resNodes, tempTk)
	}
	if util.IsNil(closeToken) {
		return nil
	}
	return &simpleBlockToken{
		tokenCommon{openToken.tokenCursorFrom(), closeToken.tokenCursorTo()},
		blockType, resNodes,
	}
}

// Returns nil if not found
func (ts *tokenStream) consumeCurlyBlock() *simpleBlockToken {
	return ts.consumeSimpleBlock(tokenTypeLeftCurlyBracket, tokenTypeRightCurlyBracket)
}

// Returns nil if not found
func (ts *tokenStream) consumeSquareBlock() *simpleBlockToken {
	return ts.consumeSimpleBlock(tokenTypeLeftSquareBracket, tokenTypeRightSquareBracket)
}

// Returns nil if not found
func (ts *tokenStream) consumeParenBlock() *simpleBlockToken {
	return ts.consumeSimpleBlock(tokenTypeLeftParen, tokenTypeRightParen)
}

// Returns nil if not found
func (ts *tokenStream) consumeFunc() *astFuncToken {
	fnValueNodes := []token{}
	var fnToken funcKeywordToken
	if temp := ts.consumeTokenWith(tokenTypeFuncKeyword); !util.IsNil(temp) {
		fnToken = temp.(funcKeywordToken)
	} else {
		return nil
	}
	var closeToken token
	for {
		tempTk := ts.consumeComponentValue()
		if util.IsNil(tempTk) || tempTk.tokenType() == tokenTypeRightParen {
			closeToken = tempTk
			break
		}
		fnValueNodes = append(fnValueNodes, tempTk)
	}

	return &astFuncToken{
		tokenCommon{fnToken.tokenCursorFrom(), closeToken.tokenCursorTo()},
		fnToken.value, fnValueNodes,
	}
}

// Returns nil if not found
func (ts *tokenStream) consumeComponentValue() token {
	if res := ts.consumeCurlyBlock(); res != nil {
		return *res
	}
	if res := ts.consumeSquareBlock(); res != nil {
		return *res
	}
	if res := ts.consumeParenBlock(); res != nil {
		return *res
	}
	if res := ts.consumeFunc(); res != nil {
		return *res
	}
	if res := ts.consumePreservedToken(); res != nil {
		return res
	}
	return nil
}

// Returns nil if not found
func (ts *tokenStream) consumeQualifiedRule() *qualifiedRuleToken {
	prelude := []token{}

	for {
		block := ts.consumeCurlyBlock()
		if block != nil {
			return &qualifiedRuleToken{
				tokenCommon{block.cursorFrom, block.cursorTo},
				prelude,
				block.body,
			}
		} else if ts.isEnd() {
			return nil
		}
		prelude = append(prelude, ts.consumeComponentValue())
	}
}

// Returns nil if not found
func (ts *tokenStream) consumeAtRule() *atRuleToken {
	prelude := []token{}
	var kwdToken atKeywordToken
	if temp := ts.consumeTokenWith(tokenTypeAtKeyword); !util.IsNil(temp) {
		kwdToken = temp.(atKeywordToken)
	} else {
		return nil
	}

	for {
		block := ts.consumeCurlyBlock()
		if block != nil {
			return &atRuleToken{
				tokenCommon{block.cursorFrom, block.cursorTo},
				kwdToken.name,
				prelude,
				block.body,
			}
		} else if ts.isEnd() {
			return nil
		}
		prelude = append(prelude, ts.consumeComponentValue())
	}
}

// Returns nil if not found
func (ts *tokenStream) consumeDeclaration() *declarationToken {
	// <name>  :  contents  !important -----------------------------------------
	var identTk identToken
	if temp := ts.consumeTokenWith(tokenTypeIdent); temp != nil {
		identTk = temp.(identToken)
	} else {
		return nil
	}
	declName := identTk.value
	declValue := []token{}
	declIsImportant := false
	// name<  >:  contents  !important -----------------------------------------
	for !util.IsNil(ts.consumeTokenWith(tokenTypeWhitespace)) {
	}
	// name  <:>  contents  !important -----------------------------------------
	if ts.consumeTokenWith(tokenTypeColon) == nil {
		// Parse error
		return nil
	}
	// name  :<  >contents  !important -----------------------------------------
	for !util.IsNil(ts.consumeTokenWith(tokenTypeWhitespace)) {
	}

	// name  :  <contents  !important> -----------------------------------------
	for {
		tempTk := ts.consumeComponentValue()
		if util.IsNil(tempTk) {
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
	return &declarationToken{
		tokenCommon{identTk.cursorFrom, lastNode.tokenCursorTo()},
		declName,
		declValue,
		declIsImportant,
	}
}

func (ts *tokenStream) consumeDeclarationList() []token {
	decls := []token{}

	for {
		token := ts.consumeToken()

		if util.IsNil(token) {
			break
		} else if token.tokenType() == tokenTypeWhitespace || token.tokenType() == tokenTypeSemicolon {
			continue
		} else if token.tokenType() == tokenTypeAtKeyword {
			ts.cursor--
			rule := ts.consumeAtRule()
			if rule == nil {
				panic("unreachable")
			}
			decls = append(decls, rule)
		} else if token.tokenType() == tokenTypeIdent {
			tempStream := tokenStream{}
			tempStream.tokens = append(tempStream.tokens, token)
			for {
				token = ts.consumeToken()
				if util.IsNil(token) || token.tokenType() == tokenTypeSemicolon {
					break
				}
				tempStream.tokens = append(tempStream.tokens, token)
			}
			decl := tempStream.consumeDeclaration()
			if decl == nil {
				break
			}
			decls = append(decls, *decl)
		} else {
			// PARSE ERROR
			for {
				oldCursor := ts.cursor
				token = ts.consumeToken()
				if util.IsNil(token) || token.tokenType() == tokenTypeSemicolon {
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
		token := ts.consumeToken()

		if util.IsNil(token) {
			break
		} else if token.tokenType() == tokenTypeWhitespace || token.tokenType() == tokenTypeSemicolon {
			continue
		} else if token.tokenType() == tokenTypeAtKeyword {
			ts.cursor--
			rule := ts.consumeAtRule()
			if rule == nil {
				panic("unreachable")
			}
			decls = append(decls, rule)
		} else if token.tokenType() == tokenTypeIdent {
			tempStream := tokenStream{}
			tempStream.tokens = append(tempStream.tokens, token)
			for {
				token = ts.consumeToken()
				if util.IsNil(token) || token.tokenType() == tokenTypeSemicolon {
					break
				}
				tempStream.tokens = append(tempStream.tokens, token)
			}
			decl := tempStream.consumeDeclaration()
			if decl == nil {
				break
			}
			decls = append(decls, *decl)
		} else if token.tokenType() == tokenTypeDelim && token.(delimToken).value == '&' {
			ts.cursor--
			if rule := ts.consumeQualifiedRule(); rule != nil {
				rules = append(rules, *rule)
			}
		} else {
			// PARSE ERROR
			for {
				oldCursor := ts.cursor
				token = ts.consumeToken()
				if util.IsNil(token) || token.tokenType() == tokenTypeSemicolon {
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
		token := ts.consumeToken()

		if util.IsNil(token) {
			break
		} else if token.tokenType() == tokenTypeWhitespace || token.tokenType() == tokenTypeSemicolon {
			continue
		} else if token.tokenType() == tokenTypeCdo || token.tokenType() == tokenTypeCdc {
			if topLevelFlag {
				continue
			}
			ts.cursor--
			if rule := ts.consumeQualifiedRule(); rule != nil {
				rules = append(rules, *rule)
			}
		} else if token.tokenType() == tokenTypeAtKeyword {
			ts.cursor--
			rule := ts.consumeAtRule()
			if rule == nil {
				panic("unreachable")
			}
			rules = append(rules, *rule)
		} else {
			ts.cursor--
			if rule := ts.consumeQualifiedRule(); rule != nil {
				rules = append(rules, *rule)
			}
		}
	}
	return rules
}

func parseCommaSeparatedListOfComponentValues(tokens []token) [][]token {
	valueLists := [][]token{}
	tempList := []token{}

	stream := tokenStream{tokens: tokens}
	for {
		value := stream.consumeComponentValue()
		if util.IsNil(value) || value.tokenType() == tokenTypeComma {
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
func parseListOfComponentValues(tokens []token) []token {
	tempList := []token{}

	stream := tokenStream{tokens: tokens}
	for {
		value := stream.consumeComponentValue()
		if util.IsNil(value) || value.tokenType() == tokenTypeComma {
			break
		}
		tempList = append(tempList, value)
	}
	return tempList
}

func parseComponentValue(tokens []token) token {
	ts := tokenStream{tokens: tokens}
	ts.skipWhitespaces()
	if !ts.isEnd() {
		panic("TODO: syntax error: expected component value")
	}
	value := ts.consumeComponentValue()
	ts.skipWhitespaces()
	if ts.isEnd() {
		panic("TODO: syntax error: expected eof")
	}
	return value
}

func parseListOfDeclarations(tokens []token) []token {
	stream := tokenStream{tokens: tokens}
	return stream.consumeDeclarationList()
}
func parseStyleBlockContents(tokens []token) []token {
	stream := tokenStream{tokens: tokens}
	return stream.consumeStyleBlockContents()
}
func parseDeclaration(tokens []token) *declarationToken {
	stream := tokenStream{tokens: tokens}
	stream.skipWhitespaces()
	if util.IsNil(stream.consumeTokenWith(tokenTypeIdent)) {
		panic("TODO: syntax error: expected identifier")
	}
	stream.cursor--
	node := stream.consumeDeclaration()
	if node == nil {
		panic("TODO: syntax error: expected declaration")
	}
	return node
}
func parseRule(tokens []token) token {
	ts := tokenStream{tokens: tokens}
	ts.skipWhitespaces()
	var res token
	res = ts.consumeAtRule()
	if util.IsNil(res) {
		res = ts.consumeQualifiedRule()
	}
	if util.IsNil(res) {
		panic("TODO: syntax error: expected at-rule or qualified-rule")
	}
	if ts.isEnd() {
		panic("TODO: syntax error: expected eof")
	}
	return res
}
func parseListOfRules(tokens []token) []token {
	stream := tokenStream{tokens: tokens}
	return stream.consumeListOfRules(false)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#typedef-declaration-value
func (ts *tokenStream) _consumeDeclarationValue(anyValue bool) []token {
	out := []token{}
	openBlockTokens := []tokenType{}
	for {
		oldCursor := ts.cursor
		tk := ts.consumeToken()
		if util.IsNil(tk) {
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
		out = append(out, tk)
	}
	if len(out) == 0 {
		return nil
	}
	return out
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
func parseCommaSeparatedRepeation[T any](ts *tokenStream, maxRepeats int, parser func(ts *tokenStream) (T, error)) ([]T, error) {
	res := []T{}
	for {
		x, err := parser(ts)
		if util.IsNil(x) {
			if err != nil {
				return nil, err
			} else if len(res) != 0 {
				return nil, errors.New("expected something after ','")
			} else {
				break
			}
		}
		res = append(res, x)
		if maxRepeats != 0 && maxRepeats <= len(res) {
			break
		}
		ts.skipWhitespaces()
		if util.IsNil(ts.consumeTokenWith(tokenTypeComma)) {
			break
		}
		ts.skipWhitespaces()
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// This can be used to parse where a CSS syntax can be repeated multiple times.
// If maxRepeats is 0, repeat count is unlimited.
//
// https://www.w3.org/TR/css-values-4/#mult-num-range
func parseRepeation[T any](ts *tokenStream, maxRepeats int, parser func(ts *tokenStream) (T, error)) ([]T, error) {
	res := []T{}
	for {
		x, err := parser(ts)
		if util.IsNil(x) {
			if err != nil {
				return nil, err
			} else {
				break
			}
		}
		res = append(res, x)
		if maxRepeats != 0 && maxRepeats <= len(res) {
			break
		}
		ts.skipWhitespaces()
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-parse-something-according-to-a-css-grammar
func parse[T any](tokens []token, parser func(ts *tokenStream) (T, error)) (T, error) {
	compValues := parseListOfComponentValues(tokens)
	stream := tokenStream{tokens: compValues}
	return parser(&stream)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-decode-bytes
func decodeBytes(bytes []byte) string {
	fallback := cssDetermineFallbackEncoding(bytes)
	input := encoding.IoQueueFromSlice(bytes)
	output := encoding.IoQueueFromSlice[rune](nil)
	encoding.Decode(&input, fallback, &output)
	return string(encoding.IoQueueToSlice[rune](output))
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#determine-the-fallback-encoding
func cssDetermineFallbackEncoding(bytes []byte) encoding.Type {
	// Check if HTTP or equivalent protocol provides an encoding label ---------
	// TODO

	// Check '@charset "xxx";' byte sequence -----------------------------------
	if 1024 < len(bytes) {
		bytes = bytes[:1024]
	}
	// NOTE: Below sequence of bytes are '@charset "' in ASCII
	if 10 <= len(bytes) || !slices.Equal([]byte{0x40, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x20, 0x22}, bytes[:10]) {
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
func ParseStylesheet(input []byte, location *string) (cssom.Stylesheet, error) {
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-stylesheets
	tokens, err := tokenize(input)
	if err != nil {
		return cssom.Stylesheet{}, err
	}
	stylesheet := cssom.Stylesheet{
		Location: location,
	}
	ts := tokenStream{tokens: tokens}
	ruleNodes := ts.consumeListOfRules(true)

	// Parse top-level qualified rules as style rules
	stylesheet.StyleRules = parseStyleRulesFromNodes(ruleNodes)

	return stylesheet, nil
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#style-rules
func parseStyleRulesFromNodes(ruleNodes []token) []cssom.StyleRule {
	styleRules := []cssom.StyleRule{}
	printRawRuleNodes := false
	for _, n := range ruleNodes {
		if n.tokenType() != tokenTypeQualifiedRule {
			continue
		}
		qrule := n.(qualifiedRuleToken)
		preludeStream := tokenStream{tokens: qrule.prelude}
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
		contentsStream := tokenStream{tokens: qrule.body}
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
				innerAs := tokenStream{tokens: declNode.value}
				innerAs.skipWhitespaces()
				value, ok := parseFunc(&innerAs)
				if !ok {
					log.Printf("bad value for property: %v (token list: %v)", declNode.name, innerAs.tokens)
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
