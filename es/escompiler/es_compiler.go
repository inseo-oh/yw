// Package escompiler provides compiler that turns input ES source code into
// series of VM instructions.
package escompiler

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/inseo-oh/yw/es"
	"github.com/inseo-oh/yw/es/vm"
	"github.com/inseo-oh/yw/util"
)

type compiler struct {
	tkh util.TokenizerHelper
}

//------------------------------------------------------------------------------
// Punctuators
//------------------------------------------------------------------------------

var otherPuncts = []string{
	"{", "(", ")", "[", "]", ".", "...", ";", ",", "<", ">", "<=", ">=", "==",
	"!=", "===", "!==", "+", "-", "*", "%", "**", "++", "--", "<<", ">>", ">>>",
	"&", "|", "^", "!", "~", "&&", "||", "??", "?", ":", "=", "+=", "-=", "*=",
	"%=", "**=", "<<=", ">>=", ">>>=", "&=", "|=", "^=", "&&=", "||=", "??=",
	"=>",
}

// https://tc39.es/ecma262/#prod-Punctuator
func (comp *compiler) consumePunctuator() string {
	oldCursor := comp.tkh.Cursor
	candiates := []string{}
	if comp.tkh.ConsumeStrIfMatches("?.", 0) != "" {
		if comp.consumeDecimalDigits(struct{ sep bool }{false}) == "" {
			candiates = append(candiates, "?")
		} else {
			comp.tkh.Cursor = oldCursor
		}
	}
	comp.tkh.Cursor = oldCursor
	candiates = append(candiates, comp.tkh.ConsumeStrIfMatchesOneOf(otherPuncts, 0))
	return util.LongestString(candiates)
}

var divPuncts = []string{"/", "/="}

// https://tc39.es/ecma262/#prod-DivPunctuator
func (comp *compiler) consumeDivPunctuator() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(divPuncts, 0)
}

// https://tc39.es/ecma262/#prod-RightBracePunctuator
func (comp *compiler) consumeRightBracePunctuator() string {
	return comp.tkh.ConsumeStrIfMatches("}", 0)
}

//------------------------------------------------------------------------------
// Line terminators and whitespaces
//------------------------------------------------------------------------------

var lineTerminators = []string{"\n", "\r", "\u2028", "\u2029"}
var lineTerminatorSequences = []string{"\n", "\r", "\u2028", "\u2029", "\r\n"}
var whitespaces = []string{" ", "\t", "\u000b", "\u000c", "\u00a0", "\uffef"}

func (comp *compiler) consumeLineTerminator() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(lineTerminators, 0)
}
func (comp *compiler) consumeLineTerminatorSequence() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(lineTerminatorSequences, 0)
}
func (comp *compiler) skipWhitespacesAndLineTerminators(lineTerminatorAllowed bool) bool {
	for {
		var found = false
		for comp.tkh.ConsumeStrIfMatchesOneOf(whitespaces, 0) != "" {
			found = true
		}
		for comp.consumeLineTerminator() != "" {
			found = true
			if !lineTerminatorAllowed {
				return false
			}
		}
		if !found {
			break
		}
	}
	return true
}

// ------------------------------------------------------------------------------
// Simple Literals
// ------------------------------------------------------------------------------
var booleanLiterals = []string{"true", "false"}

func (comp *compiler) consumeNullLiteral() string {
	return comp.tkh.ConsumeStrIfMatches("null", 0)
}
func (comp *compiler) consumeBooleanLiteral() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(booleanLiterals, 0)
}

// ------------------------------------------------------------------------------
// String Literals
// ------------------------------------------------------------------------------
var singleEscapeChars = []string{
	"'", "\"", "\\", "b", "f", "n", "r", "t", "v",
}
var escapeChars = []string{
	"'", "\"", "\\", "b", "f", "n", "r", "t", "v", "0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "x", "u",
}

// https://tc39.es/ecma262/#prod-StringLiteral
func (comp *compiler) consumeStringLiteral() string {
	if comp.tkh.IsEof() {
		return ""
	}
	startCursorPos := comp.tkh.Cursor
	endCursorPos := 0
	cursorPos := comp.tkh.Cursor
	openChar := comp.tkh.ConsumeCharIfMatchesOneOf("\"'")
	if openChar == -1 {
		return ""
	}
	for {
		currChar := comp.tkh.PeekChar()
		if currChar == openChar {
			// End of literal
			cursorPos++
			break
		} else if currChar == '\\' {
			// Escape sequence or line continuation
			comp.tkh.Cursor = cursorPos + 1
			if comp.consumeLineTerminatorSequence() != "" {
				// Line continuation
				cursorPos = comp.tkh.Cursor
				continue
			} else if comp.consumeEscapeSequence() != "" {
				// Escape sequence
				cursorPos = comp.tkh.Cursor
				continue
			} else {
				break
			}
		} else if currChar == 0x0028 || currChar == 0x0029 {
			// <LS> and <PS>
			cursorPos++
			continue
		} else {
			// Other characters
			comp.tkh.Cursor = cursorPos
			if comp.consumeLineTerminator() != "" {
				cursorPos = comp.tkh.Cursor
				break
			}
			cursorPos++
			continue
		}
	}
	endCursorPos = cursorPos
	comp.tkh.Cursor = cursorPos
	return string(comp.tkh.Str[startCursorPos:endCursorPos])
}

// https://tc39.es/ecma262/#prod-EscapeSequence
func (comp *compiler) consumeEscapeSequence() string {
	oldCursor := comp.tkh.Cursor
	candidates := []string{}

	// https://tc39.es/ecma262/#prod-SingleEscapeCharacter
	candidates = append(candidates, comp.tkh.ConsumeStrIfMatchesOneOf(singleEscapeChars, 0))
	comp.tkh.Cursor = oldCursor

	// https://tc39.es/ecma262/#prod-NonEscapeCharacter
	temp := comp.tkh.ConsumeStrIfMatchesOneOf(escapeChars, 0)
	if temp == "" {
		temp = comp.consumeLineTerminator()
		if temp == "" {
			candidates = append(candidates, string(comp.tkh.PeekChar()))
		}
	}

	res := util.LongestString(candidates)
	comp.tkh.Cursor = oldCursor + len([]rune(res))
	return res
}

// ------------------------------------------------------------------------------
// Numeric Literals
// ------------------------------------------------------------------------------

var (
	regexpNonZeroDigit         = regexp.MustCompile("[1-9]")
	regexpDecimalDigits        = regexp.MustCompile("[0-9]+")
	regexpDecimalDigitsWithSep = regexp.MustCompile("[0-9_]+")
	regexpBinaryDigits         = regexp.MustCompile("[0-1]+")
	regexpBinaryDigitsWithSep  = regexp.MustCompile("[0-1_]+")
	regexpOctalDigits          = regexp.MustCompile("[0-7]+")
	regexpOctalDigitsWithSep   = regexp.MustCompile("[0-7_]+")
	regexpNonOctalDigits       = regexp.MustCompile("[8-9]+")
	regexpHexDigits            = regexp.MustCompile("[0-9A-Fa-f]+")
	regexpHexDigitsWithSep     = regexp.MustCompile("[0-9A-Fa-f_]+")
	regexpExponentPart         = regexp.MustCompile("[eE][+-]?[0-9]+")
	regexpExponentPartWithSep  = regexp.MustCompile("[eE][+-]?[0-9_]+")
)

// https://tc39.es/ecma262/#sec-literals-numeric-literals
func (comp *compiler) consumeNumericLiteral() string {
	oldCursor := comp.tkh.Cursor
	candidates := []string{}

	candidates = append(candidates, comp.consumeDecimalLiteral())
	comp.tkh.Cursor = oldCursor
	// TODO: https://tc39.es/ecma262/#prod-DecimalBigIntegerLiteral
	candidates = append(candidates, comp.consumeNonDecimalIntegerLiteral(struct{ sep bool }{true}))
	comp.tkh.Cursor = oldCursor
	candidates = append(candidates, comp.consumeNonOctalDecimalIntegerLiteral())
	comp.tkh.Cursor = oldCursor
	// TODO: https://tc39.es/ecma262/#prod-NonDecimalIntegerLiteral + BigIntLiteralSuffix
	candidates = append(candidates, comp.consumeLegacyOctalIntegerLiteral())
	comp.tkh.Cursor = oldCursor

	res := util.LongestString(candidates)
	comp.tkh.Cursor = oldCursor + len([]rune(res))
	return res
}

// https://tc39.es/ecma262/#prod-DecimalLiteral
func (comp *compiler) consumeDecimalLiteral() string {
	oldCursor := comp.tkh.Cursor
	integer := comp.consumeDecimalIntegerLiteral()

	if comp.tkh.ConsumeStrIfMatches(".", 0) != "" {
		digits := comp.consumeDecimalDigits(struct{ sep bool }{true})
		exp := comp.consumeExponentPart(struct{ sep bool }{true})
		return integer + "." + digits + exp
	} else if integer == "" {
		comp.tkh.Cursor = oldCursor
		return ""
	}
	exp := comp.consumeExponentPart(struct{ sep bool }{true})
	return integer + exp
}

// https://tc39.es/ecma262/#prod-DecimalIntegerLiteral
func (comp *compiler) consumeDecimalIntegerLiteral() string {
	oldCursor := comp.tkh.Cursor
	candidates := []string{}

	candidates = append(candidates, comp.tkh.ConsumeStrIfMatches("0", 0))
	comp.tkh.Cursor = oldCursor
	digits := comp.consumeNonZeroDigit()
	if digits != "" {
		hasSep := comp.tkh.ConsumeStrIfMatches("_", 0) != ""
		digitsTemp := comp.consumeDecimalDigits(struct{ sep bool }{true})
		if hasSep && digitsTemp == "" {
			comp.tkh.Cursor = oldCursor
			return ""
		}
		digits += digitsTemp
		candidates = append(candidates, digits)
	}
	// TODO: https://tc39.es/ecma262/#prod-NonOctalDecimalIntegerLiteral
	res := util.LongestString(candidates)
	comp.tkh.Cursor = oldCursor + len([]rune(res))
	return res
}

// https://tc39.es/ecma262/#prod-NonZeroDigit
func (comp *compiler) consumeNonZeroDigit() string {
	return comp.tkh.ConsumeStrIfMatchesR(*regexpNonZeroDigit)
}

// https://tc39.es/ecma262/#prod-DecimalDigits
func (comp *compiler) consumeDecimalDigits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *regexpDecimalDigitsWithSep
	} else {
		pat = *regexpDecimalDigits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-ExponentPart
func (comp *compiler) consumeExponentPart(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *regexpExponentPartWithSep
	} else {
		pat = *regexpExponentPart
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-NonDecimalIntegerLiteral
func (comp *compiler) consumeNonDecimalIntegerLiteral(flags struct{ sep bool }) string {
	oldCursor := comp.tkh.Cursor
	candidates := []string{}

	candidates = append(candidates, comp.consumeBinaryIntegerLiteral(struct{ sep bool }{flags.sep}))
	comp.tkh.Cursor = oldCursor
	candidates = append(candidates, comp.consumeOctalIntegerLiteral(struct{ sep bool }{flags.sep}))
	comp.tkh.Cursor = oldCursor
	candidates = append(candidates, comp.consumeHexIntegerLiteral(struct{ sep bool }{flags.sep}))
	comp.tkh.Cursor = oldCursor

	res := util.LongestString(candidates)
	comp.tkh.Cursor = oldCursor + len([]rune(res))
	return res
}

// https://tc39.es/ecma262/#prod-BinaryIntegerLiteral
func (comp *compiler) consumeBinaryIntegerLiteral(flags struct{ sep bool }) string {
	oldCursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatchesOneOf([]string{"0b", "0B"}, 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consumeBinaryDigits(struct{ sep bool }{flags.sep})
	if digits == "" {
		comp.tkh.Cursor = oldCursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-BinaryDigits
func (comp *compiler) consumeBinaryDigits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *regexpBinaryDigitsWithSep
	} else {
		pat = *regexpBinaryDigits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-OctalIntegerLiteral
func (comp *compiler) consumeOctalIntegerLiteral(flags struct{ sep bool }) string {
	oldCursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatchesOneOf([]string{"0o", "0O"}, 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consumeOctalDigits(struct{ sep bool }{flags.sep})
	if digits == "" {
		comp.tkh.Cursor = oldCursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-OctalDigits
func (comp *compiler) consumeOctalDigits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *regexpOctalDigitsWithSep
	} else {
		pat = *regexpOctalDigits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-HexIntegerLiteral
func (comp *compiler) consumeHexIntegerLiteral(flags struct{ sep bool }) string {
	oldCursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatchesOneOf([]string{"0x", "0X"}, 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consumeHexDigits(struct{ sep bool }{flags.sep})
	if digits == "" {
		comp.tkh.Cursor = oldCursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-HexDigits
func (comp *compiler) consumeHexDigits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *regexpHexDigitsWithSep
	} else {
		pat = *regexpHexDigits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-LegacyOctalIntegerLiteral
func (comp *compiler) consumeLegacyOctalIntegerLiteral() string {
	oldCursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatches("0", 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consumeOctalDigits(struct{ sep bool }{false})
	if digits == "" {
		comp.tkh.Cursor = oldCursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-NonOctalDecimalIntegerLiteral
func (comp *compiler) consumeNonOctalDecimalIntegerLiteral() string {
	oldCursor := comp.tkh.Cursor
	if comp.tkh.ConsumeStrIfMatches("0", 0) == "" {
		return ""
	}
	comp.tkh.Cursor = oldCursor
	// If we have octal digits at beginning, and no decimal follows, it is just an octal.
	if comp.consumeOctalDigits(struct{ sep bool }{false}) != "" {
		if comp.consumeDecimalDigits(struct{ sep bool }{false}) == "" {
			// This is an octal integer
			comp.tkh.Cursor = oldCursor
			return ""
		} else {
			// It starts with octal digits, but it is followed by non-octal digits
			comp.tkh.Cursor = oldCursor
		}
	}
	return comp.consumeDecimalDigits(struct{ sep bool }{false})
}

// ------------------------------------------------------------------------------
// Identifiers
// ------------------------------------------------------------------------------

// TODO: We should be accepting any character that is in the unicode ID_Start and ID_Continue category
var regexpIdent = regexp.MustCompile("[a-zA-Z$_][a-zA-Z$_0-9]*")

// https://tc39.es/ecma262/#prod-IdentifierName
func (comp *compiler) consumeIdentifierName() string {
	return comp.tkh.ConsumeStrIfMatchesR(*regexpIdent)
}

// ------------------------------------------------------------------------------
//
// BEGINNING OF PARSER CODE & AST
//
// <About error conditions>
// Many parser functions will return pair of (<some node>, *syntaxError)
// - If <some node> isn't nil, it's success.
// - If <some node> is nil, and syntaxError is also nil, it means there was no match.
// - If <some node> is nil, but syntaxError isn't nil, match was found but there is a syntax error.
//
// ------------------------------------------------------------------------------

type syntaxError struct {
	message              syntaxErrorMsg
	cursorFrom, cursorTo int
}

type syntaxErrorMsg string

const (
	error_msg_missing_left_paren  = syntaxErrorMsg("missing (")
	error_msg_missing_right_paren = syntaxErrorMsg("missing )")
	error_msg_missing_left_brace  = syntaxErrorMsg("missing {")
	error_msg_missing_right_brace = syntaxErrorMsg("missing }")
	error_msg_missing_semicolon   = syntaxErrorMsg("missing ;")
	error_msg_missing_colon       = syntaxErrorMsg("missing :")
	error_msg_missing_expression  = syntaxErrorMsg("missing expression")
)

func chooseLongestNode(nodes []astNode) astNode {
	var longestNode astNode

	for _, node := range nodes {
		if util.IsNil(longestNode) || nodeLength(longestNode) < nodeLength(node) {
			longestNode = node
		}
	}
	return longestNode
}
func nodeLength(node astNode) int {
	return node.nodeCursorTo() - node.nodeCursorFrom()
}

var reservedWords = []string{
	"await", "break", "case", "catch", "class", "const", "continue", "debugger",
	"default", "delete", "do", "else", "enum", "export", "extends", "false",
	"finally", "for", "function", "if", "import", "in", "instanceof", "new",
	"null", "return", "super", "switch", "this", "throw", "true", "try",
	"typeof", "var", "void", "while", "with", "yield",
}

func (comp *compiler) consumeKeyword(keyword string) bool {
	cursorFrom := comp.tkh.Cursor

	if temp := comp.consumeIdentifierName(); temp == "" || temp != keyword {
		comp.tkh.Cursor = cursorFrom
		return false
	}
	return true
}

func (comp *compiler) consumePunctuatorWith(punct string) bool {
	cursorFrom := comp.tkh.Cursor

	if temp := comp.consumePunctuator(); temp == "" || temp != punct {
		comp.tkh.Cursor = cursorFrom
		return false
	}
	return true
}

func (comp *compiler) consumeSemicolon() bool {
	// TODO: Handle automatic semicolon insertion
	return comp.consumePunctuatorWith(";")
}

// https://tc39.es/ecma262/#prod-IdentifierReference
func (comp *compiler) consumeIdentifierReferenceNode(flags struct{ yield, await bool }) (res *astIdentifierReferenceNode) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	name := comp.consumeIdentifier()
	if name == "" {
		return nil
	}
	// TODO: Yield, Await
	cursorTo := comp.tkh.Cursor
	return &astIdentifierReferenceNode{
		cursorFrom: cursorFrom,
		cursorTo:   cursorTo,
		name:       name,
	}
}

// https://tc39.es/ecma262/#prod-BindingIdentifier
func (comp *compiler) consumeBindingIdentifierNode(flags struct{ yield, await bool }) (res *astBindingIdentifierNode) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	name := comp.consumeIdentifier()
	if name == "" {
		return nil
	}
	// TODO: Yield, Await
	cursorTo := comp.tkh.Cursor
	return &astBindingIdentifierNode{
		cursorFrom: cursorFrom,
		cursorTo:   cursorTo,
		name:       name,
	}
}

// https://tc39.es/ecma262/#sec-identifiers
//
// Returns empty string if not found or it's a reserved word.
func (comp *compiler) consumeIdentifier() string {
	cursorFrom := comp.tkh.Cursor
	name := comp.consumeIdentifierName()
	if slices.Contains(reservedWords, name) {
		comp.tkh.Cursor = cursorFrom
		return ""
	}
	if name == "" {
		comp.tkh.Cursor = cursorFrom
		return ""
	}
	return name
}

// https://tc39.es/ecma262/#sec-primary-expression
func (comp *compiler) consumePrimaryExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	// TODO: this

	if temp := comp.consumeIdentifierReferenceNode(struct{ yield, await bool }{flags.yield, flags.await}); !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}
	if temp := comp.consumeLiteralNode(); !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	// TODO: ArrayLiteral
	// TODO: ObjectLiteral
	// TODO: FunctionExpression
	// TODO: ClassExpression
	// TODO: GeneratorExpression
	// TODO: AsyncFunctionExpression
	// TODO: AsyncGeneratorExpression
	// TODO: RegularExpressionLiteral
	// TODO: TemplateLiteral

	// https://tc39.es/ecma262/#prod-ParenthesizedExpression
	if comp.consumePunctuatorWith("(") {
		comp.skipWhitespacesAndLineTerminators(true)
		expr, serr := comp.consumeExpressionNode(struct{ in, yield, await bool }{true, flags.yield, flags.await})
		if serr != nil {
			return nil, serr
		} else if util.IsNil(expr) {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_left_paren}
		}
		comp.skipWhitespacesAndLineTerminators(true)
		if !comp.consumePunctuatorWith(")") {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_right_paren}
		}
		return &astParenExprNode{
			cursorFrom: cursorFrom,
			cursorTo:   comp.tkh.Cursor,
			node:       expr,
		}, nil
	}

	// TODO: CoverParenthesizedExpressionAndArrowParameterList
	return chooseLongestNode(candidates), nil
}

var (
	regexpHexLiteral    = regexp.MustCompile("^0[xX]")
	regexpOctalLiteralA = regexp.MustCompile("^0o")
	regexpOctalLiteralB = regexp.MustCompile("^0[0-7]")
	regexpBinLiteral    = regexp.MustCompile("^0b")
)

// https://tc39.es/ecma262/#prod-Literal
func (comp *compiler) consumeLiteralNode() *astLiteralNode {
	cursorFrom := comp.tkh.Cursor

	if temp := comp.consumeNullLiteral(); temp != "" {
		return &astLiteralNode{
			cursorFrom: cursorFrom,
			cursorTo:   comp.tkh.Cursor,
			value:      es.NewNullValue(),
		}
	}
	if temp := comp.consumeBooleanLiteral(); temp != "" {
		switch temp {
		case "true":
			return &astLiteralNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				value:      es.NewBooleanValue(true),
			}
		case "false":
			return &astLiteralNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				value:      es.NewBooleanValue(false),
			}
		}
		panic("unreachable")
	}

	if temp := comp.consumeNumericLiteral(); temp != "" {
		// Remove seaprator
		temp = strings.ReplaceAll(temp, "_", "")
		var valueF float64
		var valueI int64
		var err error
		isInt := false
		if regexpHexLiteral.MatchString(temp) {
			valueI, err = strconv.ParseInt(strings.TrimPrefix(strings.TrimPrefix(temp, "0x"), "0X"), 16, 64)
			if err != nil {
				panic(err)
			}
			isInt = true
		} else if regexpOctalLiteralA.MatchString(temp) {
			valueI, err = strconv.ParseInt(strings.TrimPrefix(temp, "0o"), 8, 64)
			if err != nil {
				panic(err)
			}
			isInt = true
		} else if regexpOctalLiteralB.MatchString(temp) && !regexpNonOctalDigits.MatchString(temp) {
			valueI, err = strconv.ParseInt(strings.TrimPrefix(temp, "0"), 8, 64)
			if err != nil {
				panic(err)
			}
			isInt = true
		} else if regexpBinLiteral.MatchString(temp) {
			valueI, err = strconv.ParseInt(strings.TrimPrefix(temp, "0b"), 2, 64)
			if err != nil {
				panic(err)
			}
			isInt = true
		} else {
			valueF, err = strconv.ParseFloat(temp, 64)
			if err != nil {
				panic(err)
			}
			isInt = false
		}
		if isInt {
			return &astLiteralNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				value:      es.NewNumberValueI(valueI),
			}
		} else {
			return &astLiteralNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				value:      es.NewNumberValueF(valueF),
			}
		}

	}
	if temp := comp.consumeStringLiteral(); temp != "" {
		panic("TODO")
	}
	return nil
}

// https://tc39.es/ecma262/#prod-MemberExpression
func (comp *compiler) consumeMemberExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	if temp, serr := comp.consumePrimaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	// TODO: MemberExpression [ Expression ]
	// TODO: MemberExpression . IdentifierName
	// TODO: MemberExpression TemplateLiteral
	// TODO: SuperProperty
	// TODO: MetaProperty
	// TODO: new MemberExpression Arguments
	// TODO: MemberExpression . PrivateIdentifier

	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-NewExpression
func (comp *compiler) consumeNewExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	// STUB
	return comp.consumeMemberExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-CallExpression
func (comp *compiler) consumeCallExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	if callee, serr := comp.consumeMemberExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(callee) {
		comp.skipWhitespacesAndLineTerminators(true)
		args, serr := comp.consumeArgumentsNode(struct{ yield, await bool }{flags.yield, flags.await})
		if serr != nil {
			return nil, serr
		} else if args == nil {
			return nil, nil
		}
		cursorTo := comp.tkh.Cursor
		res = astCallExprNode{
			cursorFrom: cursorFrom,
			cursorTo:   cursorTo,
			callee:     callee,
			args:       *args,
		}
		// Handle additional things that might follow the function all, like:
		//  foo(bar)(baz)[qux]
		for {
			if anotherArgs, serr := comp.consumeArgumentsNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if anotherArgs != nil {
				anotherCursorTo := comp.tkh.Cursor
				res = astCallExprNode{
					cursorFrom: cursorFrom,
					cursorTo:   anotherCursorTo,
					callee:     res,
					args:       *args,
				}
			} else if comp.consumePunctuatorWith("[") {
				panic("TODO")
			} else if comp.consumePunctuatorWith(".") {
				panic("TODO")
			} else {
				break
			}
		}
		return res, nil
	} else if comp.consumeKeyword("super") {
		panic("TODO")
	} else if comp.consumeKeyword("import") {
		panic("TODO")
	}
	return nil, nil
}

// https://tc39.es/ecma262/#prod-Arguments
func (comp *compiler) consumeArgumentsNode(flags struct{ yield, await bool }) (res *astArguments, serr *syntaxError) {
	var restArgs astNode
	cursorFrom := comp.tkh.Cursor
	defer func() {
		if res == nil {
			comp.tkh.Cursor = cursorFrom
		}
	}()

	if !comp.consumePunctuatorWith("(") {
		return nil, nil
	}

	args := []astNode{}
	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("...") {
			comp.skipWhitespacesAndLineTerminators(true)
			if temp, serr := comp.consumeAssignmentExpressionNode(struct{ in, yield, await bool }{true, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(restArgs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				restArgs = temp
			}
		} else if temp, serr := comp.consumeAssignmentExpressionNode(struct{ in, yield, await bool }{true, flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if !util.IsNil(temp) {
			args = append(args, temp)
		} else {
			break
		}
		comp.skipWhitespacesAndLineTerminators(true)
		if !comp.consumePunctuatorWith(",") {
			break
		}
	}
	if !comp.consumePunctuatorWith(")") {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_right_paren}
	}
	return &astArguments{
		cursorFrom: cursorFrom,
		cursorTo:   comp.tkh.Cursor,
		args:       args,
		restArgs:   restArgs,
	}, nil
}

// https://tc39.es/ecma262/#prod-LeftHandSideExpression
func (comp *compiler) consumeLeftHandSideExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	if temp, serr := comp.consumeNewExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}
	if temp, serr := comp.consumeCallExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}
	// TODO: https://tc39.es/ecma262/#prod-OptionalExpression
	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-UpdateExpression
func (comp *compiler) consumeUpdateExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	// STUB
	return comp.consumeLeftHandSideExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-UnaryExpression
func (comp *compiler) consumeUnaryExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	if comp.consumeKeyword("delete") {
		panic("TODO")
	} else if comp.consumeKeyword("void") {
		panic("TODO")
	} else if comp.consumeKeyword("typeof") {
		panic("TODO")
	} else if comp.consumePunctuatorWith("+") {
		comp.skipWhitespacesAndLineTerminators(true)

		if node, serr := comp.consumeUnaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if util.IsNil(node) {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
		} else {
			return astUnaryOpNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				node:       node,
				tp:         astOpTypePlus,
			}, nil
		}

	} else if comp.consumePunctuatorWith("-") {
		comp.skipWhitespacesAndLineTerminators(true)
		if node, serr := comp.consumeUnaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if util.IsNil(node) {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
		} else {
			return astUnaryOpNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				node:       node,
				tp:         astOpTypeNeg,
			}, nil
		}
	} else if comp.consumePunctuatorWith("~") {
		comp.skipWhitespacesAndLineTerminators(true)
		if node, serr := comp.consumeUnaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if util.IsNil(node) {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
		} else {
			return astUnaryOpNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				node:       node,
				tp:         astOpTypeBNot,
			}, nil
		}
	} else if comp.consumePunctuatorWith("!") {
		comp.skipWhitespacesAndLineTerminators(true)
		if node, serr := comp.consumeUnaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if util.IsNil(node) {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
		} else {
			return astUnaryOpNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				node:       node,
				tp:         astOpTypeLNot,
			}, nil
		}
	} else if flags.await && comp.consumeKeyword("await") {
		comp.skipWhitespacesAndLineTerminators(true)
		if node, serr := comp.consumeUnaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if util.IsNil(node) {
			return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
		} else {
			return astUnaryOpNode{
				cursorFrom: cursorFrom,
				cursorTo:   comp.tkh.Cursor,
				node:       node,
				tp:         astOpTypeAwait,
			}, nil
		}
	}
	return comp.consumeUpdateExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-ExponentiationExpression
func (comp *compiler) consumeExponentationExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeUpdateExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return comp.consumeUnaryExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	}
	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumePunctuatorWith("**") {
		return lhs, nil
	}
	comp.skipWhitespacesAndLineTerminators(true)
	rhs, serr := comp.consumeExponentationExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(rhs) {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
	}
	return astBinaryOpNode{
		cursorFrom: cursorFrom,
		cursorTo:   comp.tkh.Cursor,
		lhsNode:    lhs,
		rhsNode:    rhs,
		tp:         astOpTypeExponent,
	}, nil
}

// https://tc39.es/ecma262/#prod-MultiplicativeExpression
func (comp *compiler) consumeMultiplicativeExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeExponentationExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("*") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeExponentationExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeMultiply,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumeDivPunctuator() != "" {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeExponentationExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeDivide,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("%") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeExponentationExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeModulo,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-AdditiveExpression
func (comp *compiler) consumeAdditiveExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeMultiplicativeExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("+") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeMultiplicativeExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeAdd,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("-") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeMultiplicativeExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeSubtract,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-ShiftExpression
func (comp *compiler) consumeShiftExpressionNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeAdditiveExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("<<") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeAdditiveExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeLeftShift,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith(">>>") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeAdditiveExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeRightLShift,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith(">>") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeAdditiveExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeRightAShift,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-RelationalExpression
func (comp *compiler) consumeRelationalExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("<") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeLessThan,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith(">") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeGreaterThan,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("<=") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeLessThanOrEqual,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith(">=") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeGreaterThanOrEqual,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("instanceof") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeInstanceof,
					rhsNode:    rhs,
				}
			}
		} else if flags.in && comp.consumePunctuatorWith("in") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeShiftExpressionNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				panic("TODO")
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-EqualityExpression
func (comp *compiler) consumeEqualityExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeRelationalExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		// NOTE: We try longer ones first
		if comp.consumePunctuatorWith("===") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeRelationalExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeStrictEqual,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("!==") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeRelationalExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeStrictNotEqual,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("==") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeRelationalExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeEqual,
					rhsNode:    rhs,
				}
			}
		} else if comp.consumePunctuatorWith("!=") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeRelationalExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeNotEqual,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-BitwiseANDExpression
func (comp *compiler) consumeBAndExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeEqualityExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("&") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeEqualityExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeBAnd,
					rhsNode:    rhs,
				}
			}

		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-BitwiseXORExpression
func (comp *compiler) consumeBXorExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeBAndExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("^") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeBAndExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeBXor,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-BitwiseORExpression
func (comp *compiler) consumeBOrExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeBXorExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("|") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeBXorExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOPTypeBOr,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-LogicalANDExpression
func (comp *compiler) consumeLAndExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeBOrExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("&&") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeBOrExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeLAnd,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-LogicalORExpression
func (comp *compiler) consumeLOrExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeLAndExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("||") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeLAndExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeLOr,
					rhsNode:    rhs,
				}
			}

		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-CoalesceExpression
func (comp *compiler) consumeCoalesceExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeBOrExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("??") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeBOrExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astBinaryOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					tp:         astOpTypeCoalesce,
					rhsNode:    rhs,
				}
			}

		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-ShortCircuitExpression
func (comp *compiler) consumeShortCircuitExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	candidates := []astNode{}
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	if temp, serr := comp.consumeLOrExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	if temp, serr := comp.consumeCoalesceExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}
	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-ConditionalExpression
func (comp *compiler) consumeConditionalExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	cond, serr := comp.consumeShortCircuitExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(cond) {
		return nil, nil
	}
	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumePunctuatorWith("?") {
		return cond, nil
	}

	comp.skipWhitespacesAndLineTerminators(true)
	trueExpr, serr := comp.consumeAssignmentExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(trueExpr) {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
	}

	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumePunctuatorWith(":") {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_colon}
	}

	comp.skipWhitespacesAndLineTerminators(true)
	falseExpr, serr := comp.consumeAssignmentExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(falseExpr) {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
	}

	return astCondExprNode{
		cursorFrom: cursorFrom,
		cursorTo:   comp.tkh.Cursor,
		condNode:   cond,
		trueNode:   trueExpr,
		falseNode:  falseExpr,
	}, nil
}

// https://tc39.es/ecma262/#prod-AssignmentExpression
func (comp *compiler) consumeAssignmentExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	// STUB
	return comp.consumeConditionalExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-Expression
func (comp *compiler) consumeExpressionNode(flags struct{ in, yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	lhs, serr := comp.consumeAssignmentExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith(",") {
			comp.skipWhitespacesAndLineTerminators(true)
			if rhs, serr := comp.consumeAssignmentExpressionNode(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if util.IsNil(rhs) {
				return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_expression}
			} else {
				lhs = astCommaOpNode{
					cursorFrom: cursorFrom,
					cursorTo:   comp.tkh.Cursor,
					lhsNode:    lhs,
					rhsNode:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-Statement
func (comp *compiler) consumeStatementNode(flags struct{ yield, await, retrn bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	if temp, serr := comp.consumeExpressionStatementNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	if temp, serr := comp.consumeReturnStatementNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-Declaration
func (comp *compiler) consumeDeclarationNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	if temp, serr := comp.consumeHoistableDeclarationNode(struct{ yield, await, defult bool }{flags.yield, flags.await, false}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	// TODO: ClassDeclaration
	// TODO: LexicalDeclaration

	return chooseLongestNode(candidates), nil

}

// https://tc39.es/ecma262/#prod-HoistableDeclaration
func (comp *compiler) consumeHoistableDeclarationNode(flags struct{ yield, await, defult bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	if temp, serr := comp.consumeFunctionDeclarationNode(struct{ yield, await, defult bool }{flags.yield, flags.await, flags.defult}); serr != nil {
		return nil, serr
	} else if !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}

	// TODO: GeneratorDeclaration
	// TODO: AsyncFunctionDeclaration
	// TODO: AsyncGeneratorDeclaration

	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-StatementList
func (comp *compiler) consumeStatementList(flags struct{ yield, await, retrn bool }) (res []astNode, serr *syntaxError) {
	statements := []astNode{}

	for {
		cursorFrom := comp.tkh.Cursor
		comp.skipWhitespacesAndLineTerminators(true)
		candidates := []astNode{}

		if temp, serr := comp.consumeStatementNode(struct{ yield, await, retrn bool }{flags.yield, flags.await, flags.retrn}); serr != nil {
			return nil, serr
		} else if !util.IsNil(temp) {
			candidates = append(candidates, temp)
			comp.tkh.Cursor = cursorFrom
		}

		if temp, serr := comp.consumeDeclarationNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if !util.IsNil(temp) {
			candidates = append(candidates, temp)
			comp.tkh.Cursor = cursorFrom
		}

		item := chooseLongestNode(candidates)
		if !util.IsNil(item) {
			comp.tkh.Cursor = item.nodeCursorTo()
			statements = append(statements, item)
		} else {
			break
		}
	}

	return statements, nil
}

// https://tc39.es/ecma262/#prod-BindingElement
func (comp *compiler) consumeBindingElementNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}

	// https://tc39.es/ecma262/#prod-SingleNameBinding
	if temp := comp.consumeBindingIdentifierNode(struct{ yield, await bool }{flags.yield, flags.await}); !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
		// TODO: Accept initializer after the identifier (https://tc39.es/ecma262/#prod-Initializer)
	}

	// TODO: BindingPattern Initializer(opt)

	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-BindingRestElement
func (comp *compiler) consumeBindingRestElementNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	candidates := []astNode{}
	if !comp.consumePunctuatorWith("...") {
		return nil, nil
	}
	comp.skipWhitespacesAndLineTerminators(true)
	if temp := comp.consumeBindingIdentifierNode(struct{ yield, await bool }{flags.yield, flags.await}); !util.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursorFrom
	}
	// TODO: ...BindingPattern

	return chooseLongestNode(candidates), nil
}

// https://tc39.es/ecma262/#prod-ExpressionStatement
func (comp *compiler) consumeExpressionStatementNode(flags struct{ yield, await bool }) (res *astExprStatementNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	found := false
	comp.tkh.Lookahead(func() {
		comp.skipWhitespacesAndLineTerminators(true)
		if comp.consumePunctuatorWith("{") ||
			comp.consumeKeyword("function") ||
			comp.consumeKeyword("class") ||
			comp.consumeKeyword("let") ||
			comp.consumePunctuatorWith("[") {
			found = true
		}
		if comp.consumeKeyword("async") {
			if comp.consumeLineTerminator() == "" {
				found = true
			}
		}
	})
	if found {
		return nil, nil
	}
	expr, serr := comp.consumeExpressionNode(struct{ in, yield, await bool }{true, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(expr) {
		return nil, nil
	}
	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumeSemicolon() {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_semicolon}
	}
	return &astExprStatementNode{
		cursorFrom: cursorFrom,
		cursorTo:   comp.tkh.Cursor,
		node:       expr,
	}, nil
}

// https://tc39.es/ecma262/#prod-ReturnStatement
func (comp *compiler) consumeReturnStatementNode(flags struct{ yield, await bool }) (res *astReturnStatementNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	if !comp.consumeKeyword("return") {
		return nil, nil
	}
	comp.skipWhitespacesAndLineTerminators(true)
	if comp.consumeSemicolon() {
		return &astReturnStatementNode{
			cursorFrom: cursorFrom,
			cursorTo:   comp.tkh.Cursor,
			node:       nil,
		}, nil
	}
	comp.skipWhitespacesAndLineTerminators(false)
	expr, serr := comp.consumeExpressionNode(struct{ in, yield, await bool }{true, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if util.IsNil(expr) {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_semicolon}
	}
	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumeSemicolon() {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_semicolon}
	}
	return &astReturnStatementNode{
		cursorFrom: cursorFrom,
		cursorTo:   comp.tkh.Cursor,
		node:       expr,
	}, nil
}

// https://tc39.es/ecma262/#prod-FormalParameters
func (comp *compiler) consumeFormalParameters(flags struct{ yield, await bool }) (res *astFormalParameters, serr *syntaxError) {
	params := []astNode{}

	for {
		comp.skipWhitespacesAndLineTerminators(true)
		if temp, serr := comp.consumeFormalParameterNode(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if util.IsNil(temp) {
			break
		} else {
			params = append(params, temp)
		}
		comp.skipWhitespacesAndLineTerminators(true)
		if !comp.consumePunctuatorWith(",") {
			break
		}
	}
	restParam, serr := comp.consumeFunctionRestParameterNode(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	}
	return &astFormalParameters{
		params:    params,
		restParam: restParam,
	}, nil
}

// https://tc39.es/ecma262/#prod-FormalParameter
func (comp *compiler) consumeFormalParameterNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	return comp.consumeBindingElementNode(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-FunctionRestParameter
func (comp *compiler) consumeFunctionRestParameterNode(flags struct{ yield, await bool }) (res astNode, serr *syntaxError) {
	return comp.consumeBindingRestElementNode(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-FunctionDeclaration
func (comp *compiler) consumeFunctionDeclarationNode(flags struct{ yield, await, defult bool }) (res astNode, serr *syntaxError) {
	cursorFrom := comp.tkh.Cursor
	defer func() { comp.setFinalCursor(res, cursorFrom) }()

	if !comp.consumeKeyword("function") {
		return nil, nil
	}

	comp.skipWhitespacesAndLineTerminators(true)
	funcIdent := comp.consumeBindingIdentifierNode(struct{ yield, await bool }{flags.yield, flags.await})
	if !flags.defult && funcIdent == nil {
		return nil, nil
	}

	if !comp.consumePunctuatorWith("(") {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_left_paren}
	}

	comp.skipWhitespacesAndLineTerminators(true)
	params, serr := comp.consumeFormalParameters(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	}

	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumePunctuatorWith(")") {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_left_paren}
	}

	comp.skipWhitespacesAndLineTerminators(true)
	if !comp.consumePunctuatorWith("{") {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_left_brace}
	}

	comp.skipWhitespacesAndLineTerminators(true)
	statements, serr := comp.consumeStatementList(struct{ yield, await, retrn bool }{flags.yield, flags.await, true})
	if serr != nil {
		return nil, serr
	}

	comp.skipWhitespacesAndLineTerminators(true)
	if comp.consumeRightBracePunctuator() == "" {
		return nil, &syntaxError{cursorFrom: cursorFrom, cursorTo: comp.tkh.Cursor, message: error_msg_missing_right_brace}
	}

	return astFunctionDeclNode{
		cursorFrom: cursorFrom,
		cursorTo:   comp.tkh.Cursor,
		ident:      *funcIdent,
		params:     *params,
		body:       statements,
	}, nil
}

func (comp *compiler) setFinalCursor(node astNode, cursorInitial int) {
	if !util.IsNil(node) {
		comp.tkh.Cursor = node.nodeCursorTo()
	} else {
		comp.tkh.Cursor = cursorInitial
	}
}

// Compile compiles given ECMAScript code to series of VM instructions.
//
// TODO(ois): Compile should return normal Go error, instead of private syntaxError.
func Compile(str string) (res []vm.Instr, serr *syntaxError) {
	compiler := compiler{tkh: util.TokenizerHelper{Str: []rune(str)}}
	nodes, serr := compiler.consumeStatementList(struct{ yield, await, retrn bool }{false, false, false})
	if serr != nil {
		return nil, serr
	}
	return makeCodeForAstNodes(nodes), nil
}
