package libes

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"yw/libcommon"
)

type es_compiler struct {
	tkh libcommon.TokenizerHelper
}

//------------------------------------------------------------------------------
// Punctuators
//------------------------------------------------------------------------------

var es_other_puncts = []string{
	"{", "(", ")", "[", "]", ".", "...", ";", ",", "<", ">", "<=", ">=", "==",
	"!=", "===", "!==", "+", "-", "*", "%", "**", "++", "--", "<<", ">>", ">>>",
	"&", "|", "^", "!", "~", "&&", "||", "??", "?", ":", "=", "+=", "-=", "*=",
	"%=", "**=", "<<=", ">>=", ">>>=", "&=", "|=", "^=", "&&=", "||=", "??=",
	"=>",
}

// https://tc39.es/ecma262/#prod-Punctuator
func (comp *es_compiler) consume_punctuator() string {
	old_cursor := comp.tkh.Cursor
	candiates := []string{}
	if comp.tkh.ConsumeStrIfMatches("?.", 0) != "" {
		if comp.consume_decimal_digits(struct{ sep bool }{false}) == "" {
			candiates = append(candiates, "?")
		} else {
			comp.tkh.Cursor = old_cursor
		}
	}
	comp.tkh.Cursor = old_cursor
	candiates = append(candiates, comp.tkh.ConsumeStrIfMatchesOneOf(es_other_puncts, 0))
	return libcommon.ConsumeLongestString(candiates)
}

var es_div_puncts = []string{"/", "/="}

// https://tc39.es/ecma262/#prod-DivPunctuator
func (comp *es_compiler) consume_div_punctuator() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(es_div_puncts, 0)
}

// https://tc39.es/ecma262/#prod-RightBracePunctuator
func (comp *es_compiler) consume_right_brace_punctuator() string {
	return comp.tkh.ConsumeStrIfMatches("}", 0)
}

//------------------------------------------------------------------------------
// Line terminators and whitespaces
//------------------------------------------------------------------------------

var es_line_terminators = []string{"\n", "\r", "\u2028", "\u2029"}
var es_line_terminator_sequences = []string{"\n", "\r", "\u2028", "\u2029", "\r\n"}
var es_whitespaces = []string{" ", "\t", "\u000b", "\u000c", "\u00a0", "\uffef"}

func (comp *es_compiler) consume_line_terminator() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(es_line_terminators, 0)
}
func (comp *es_compiler) consume_line_terminator_sequence() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(es_line_terminator_sequences, 0)
}
func (comp *es_compiler) skip_whitespaces_and_line_terminators(line_terminator_allowed bool) bool {
	for {
		var found = false
		for comp.tkh.ConsumeStrIfMatchesOneOf(es_whitespaces, 0) != "" {
			found = true
		}
		for comp.consume_line_terminator() != "" {
			found = true
			if !line_terminator_allowed {
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
var es_boolean_literals = []string{"true", "false"}

func (comp *es_compiler) consume_null_literal() string {
	return comp.tkh.ConsumeStrIfMatches("null", 0)
}
func (comp *es_compiler) consume_boolean_literal() string {
	return comp.tkh.ConsumeStrIfMatchesOneOf(es_boolean_literals, 0)
}

// ------------------------------------------------------------------------------
// String Literals
// ------------------------------------------------------------------------------
var es_single_escape_chars = []string{
	"'", "\"", "\\", "b", "f", "n", "r", "t", "v",
}
var es_escape_chars = []string{
	"'", "\"", "\\", "b", "f", "n", "r", "t", "v", "0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "x", "u",
}

// https://tc39.es/ecma262/#prod-StringLiteral
func (comp *es_compiler) consume_string_literal() string {
	if comp.tkh.IsEof() {
		return ""
	}
	start_cursor_pos := comp.tkh.Cursor
	end_cursor_pos := 0
	cursor_pos := comp.tkh.Cursor
	open_char := comp.tkh.ConsumeCharIfMatchesOneOf("\"'")
	if open_char == -1 {
		return ""
	}
	for {
		cur_char := comp.tkh.PeekChar()
		if cur_char == open_char {
			// End of literal
			cursor_pos++
			break
		} else if cur_char == '\\' {
			// Escape sequence or line continuation
			comp.tkh.Cursor = cursor_pos + 1
			if comp.consume_line_terminator_sequence() != "" {
				// Line continuation
				cursor_pos = comp.tkh.Cursor
				continue
			} else if comp.consume_escape_sequence() != "" {
				// Escape sequence
				cursor_pos = comp.tkh.Cursor
				continue
			} else {
				break
			}
		} else if cur_char == 0x0028 || cur_char == 0x0029 {
			// <LS> and <PS>
			cursor_pos++
			continue
		} else {
			// Other characters
			comp.tkh.Cursor = cursor_pos
			if comp.consume_line_terminator() != "" {
				cursor_pos = comp.tkh.Cursor
				break
			}
			cursor_pos++
			continue
		}
	}
	end_cursor_pos = cursor_pos
	comp.tkh.Cursor = cursor_pos
	return string(comp.tkh.Str[start_cursor_pos:end_cursor_pos])
}

// https://tc39.es/ecma262/#prod-EscapeSequence
func (comp *es_compiler) consume_escape_sequence() string {
	old_cursor := comp.tkh.Cursor
	candidates := []string{}

	// https://tc39.es/ecma262/#prod-SingleEscapeCharacter
	candidates = append(candidates, comp.tkh.ConsumeStrIfMatchesOneOf(es_single_escape_chars, 0))
	comp.tkh.Cursor = old_cursor

	// https://tc39.es/ecma262/#prod-NonEscapeCharacter
	temp := comp.tkh.ConsumeStrIfMatchesOneOf(es_escape_chars, 0)
	if temp == "" {
		temp = comp.consume_line_terminator()
		if temp == "" {
			candidates = append(candidates, string(comp.tkh.PeekChar()))
		}
	}

	res := libcommon.ConsumeLongestString(candidates)
	comp.tkh.Cursor = old_cursor + len([]rune(res))
	return res
}

// ------------------------------------------------------------------------------
// Numeric Literals
// ------------------------------------------------------------------------------

var es_regexp_non_zero_digit = regexp.MustCompile("[1-9]")
var es_regexp_decimal_digits = regexp.MustCompile("[0-9]+")
var es_regexp_decimal_digits_with_sep = regexp.MustCompile("[0-9_]+")
var es_regexp_binary_digits = regexp.MustCompile("[0-1]+")
var es_regexp_binary_digits_with_sep = regexp.MustCompile("[0-1_]+")
var es_regexp_octal_digits = regexp.MustCompile("[0-7]+")
var es_regexp_octal_digits_with_sep = regexp.MustCompile("[0-7_]+")
var es_regexp_non_octal_digits = regexp.MustCompile("[8-9]+")
var es_regexp_hex_digits = regexp.MustCompile("[0-9A-Fa-f]+")
var es_regexp_hex_digits_with_sep = regexp.MustCompile("[0-9A-Fa-f_]+")
var es_regexp_exponent_part = regexp.MustCompile("[eE][+-]?[0-9]+")
var es_regexp_exponent_part_with_sep = regexp.MustCompile("[eE][+-]?[0-9_]+")

// https://tc39.es/ecma262/#sec-literals-numeric-literals
func (comp *es_compiler) consume_numeric_literal() string {
	old_cursor := comp.tkh.Cursor
	candidates := []string{}

	candidates = append(candidates, comp.consume_decimal_literal())
	comp.tkh.Cursor = old_cursor
	// TODO: https://tc39.es/ecma262/#prod-DecimalBigIntegerLiteral
	candidates = append(candidates, comp.consume_non_decimal_integer_literal(struct{ sep bool }{true}))
	comp.tkh.Cursor = old_cursor
	candidates = append(candidates, comp.consume_non_octal_decimal_integer_literal())
	comp.tkh.Cursor = old_cursor
	// TODO: https://tc39.es/ecma262/#prod-NonDecimalIntegerLiteral + BigIntLiteralSuffix
	candidates = append(candidates, comp.consume_legacy_octal_integer_literal())
	comp.tkh.Cursor = old_cursor

	res := libcommon.ConsumeLongestString(candidates)
	comp.tkh.Cursor = old_cursor + len([]rune(res))
	return res
}

// https://tc39.es/ecma262/#prod-DecimalLiteral
func (comp *es_compiler) consume_decimal_literal() string {
	old_cursor := comp.tkh.Cursor
	integer := comp.consume_decimal_integer_literal()

	if comp.tkh.ConsumeStrIfMatches(".", 0) != "" {
		digits := comp.consume_decimal_digits(struct{ sep bool }{true})
		exp := comp.consume_exponent_part(struct{ sep bool }{true})
		return integer + "." + digits + exp
	} else if integer == "" {
		comp.tkh.Cursor = old_cursor
		return ""
	}
	exp := comp.consume_exponent_part(struct{ sep bool }{true})
	return integer + exp
}

// https://tc39.es/ecma262/#prod-DecimalIntegerLiteral
func (comp *es_compiler) consume_decimal_integer_literal() string {
	old_cursor := comp.tkh.Cursor
	candidates := []string{}

	candidates = append(candidates, comp.tkh.ConsumeStrIfMatches("0", 0))
	comp.tkh.Cursor = old_cursor
	digits := comp.consume_non_zero_digit()
	if digits != "" {
		has_sep := comp.tkh.ConsumeStrIfMatches("_", 0) != ""
		digits_temp := comp.consume_decimal_digits(struct{ sep bool }{true})
		if has_sep && digits_temp == "" {
			comp.tkh.Cursor = old_cursor
			return ""
		}
		digits += digits_temp
		candidates = append(candidates, digits)
	}
	// TODO: https://tc39.es/ecma262/#prod-NonOctalDecimalIntegerLiteral
	res := libcommon.ConsumeLongestString(candidates)
	comp.tkh.Cursor = old_cursor + len([]rune(res))
	return res
}

// https://tc39.es/ecma262/#prod-NonZeroDigit
func (comp *es_compiler) consume_non_zero_digit() string {
	return comp.tkh.ConsumeStrIfMatchesR(*es_regexp_non_zero_digit)
}

// https://tc39.es/ecma262/#prod-DecimalDigits
func (comp *es_compiler) consume_decimal_digits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *es_regexp_decimal_digits_with_sep
	} else {
		pat = *es_regexp_decimal_digits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-ExponentPart
func (comp *es_compiler) consume_exponent_part(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *es_regexp_exponent_part_with_sep
	} else {
		pat = *es_regexp_exponent_part
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-NonDecimalIntegerLiteral
func (comp *es_compiler) consume_non_decimal_integer_literal(flags struct{ sep bool }) string {
	old_cursor := comp.tkh.Cursor
	candidates := []string{}

	candidates = append(candidates, comp.consume_binary_integer_literal(struct{ sep bool }{flags.sep}))
	comp.tkh.Cursor = old_cursor
	candidates = append(candidates, comp.consume_octal_integer_literal(struct{ sep bool }{flags.sep}))
	comp.tkh.Cursor = old_cursor
	candidates = append(candidates, comp.consume_hex_integer_literal(struct{ sep bool }{flags.sep}))
	comp.tkh.Cursor = old_cursor

	res := libcommon.ConsumeLongestString(candidates)
	comp.tkh.Cursor = old_cursor + len([]rune(res))
	return res
}

// https://tc39.es/ecma262/#prod-BinaryIntegerLiteral
func (comp *es_compiler) consume_binary_integer_literal(flags struct{ sep bool }) string {
	old_cursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatchesOneOf([]string{"0b", "0B"}, 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consume_binary_digits(struct{ sep bool }{flags.sep})
	if digits == "" {
		comp.tkh.Cursor = old_cursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-BinaryDigits
func (comp *es_compiler) consume_binary_digits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *es_regexp_binary_digits_with_sep
	} else {
		pat = *es_regexp_binary_digits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-OctalIntegerLiteral
func (comp *es_compiler) consume_octal_integer_literal(flags struct{ sep bool }) string {
	old_cursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatchesOneOf([]string{"0o", "0O"}, 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consume_octal_digits(struct{ sep bool }{flags.sep})
	if digits == "" {
		comp.tkh.Cursor = old_cursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-OctalDigits
func (comp *es_compiler) consume_octal_digits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *es_regexp_octal_digits_with_sep
	} else {
		pat = *es_regexp_octal_digits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-HexIntegerLiteral
func (comp *es_compiler) consume_hex_integer_literal(flags struct{ sep bool }) string {
	old_cursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatchesOneOf([]string{"0x", "0X"}, 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consume_hex_digits(struct{ sep bool }{flags.sep})
	if digits == "" {
		comp.tkh.Cursor = old_cursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-HexDigits
func (comp *es_compiler) consume_hex_digits(flags struct{ sep bool }) string {
	var pat regexp.Regexp

	if flags.sep {
		pat = *es_regexp_hex_digits_with_sep
	} else {
		pat = *es_regexp_hex_digits
	}
	return comp.tkh.ConsumeStrIfMatchesR(pat)
}

// https://tc39.es/ecma262/#prod-LegacyOctalIntegerLiteral
func (comp *es_compiler) consume_legacy_octal_integer_literal() string {
	old_cursor := comp.tkh.Cursor

	prefix := comp.tkh.ConsumeStrIfMatches("0", 0)
	if prefix == "" {
		return ""
	}
	digits := comp.consume_octal_digits(struct{ sep bool }{false})
	if digits == "" {
		comp.tkh.Cursor = old_cursor
		return ""
	}
	return prefix + digits
}

// https://tc39.es/ecma262/#prod-NonOctalDecimalIntegerLiteral
func (comp *es_compiler) consume_non_octal_decimal_integer_literal() string {
	old_cursor := comp.tkh.Cursor
	if comp.tkh.ConsumeStrIfMatches("0", 0) == "" {
		return ""
	}
	comp.tkh.Cursor = old_cursor
	// If we have octal digits at beginning, and no decimal follows, it is just an octal.
	if comp.consume_octal_digits(struct{ sep bool }{false}) != "" {
		if comp.consume_decimal_digits(struct{ sep bool }{false}) == "" {
			// This is an octal integer
			comp.tkh.Cursor = old_cursor
			return ""
		} else {
			// It starts with octal digits, but it is followed by non-octal digits
			comp.tkh.Cursor = old_cursor
		}
	}
	return comp.consume_decimal_digits(struct{ sep bool }{false})
}

// ------------------------------------------------------------------------------
// Identifiers
// ------------------------------------------------------------------------------

// XXX: We should be accepting any character that is in the unicode ID_Start and
//
//	ID_Continue category, but that would involve processing Unicode Character Database.
var es_regexp_ident = regexp.MustCompile("[a-zA-Z$_][a-zA-Z$_0-9]*")

// https://tc39.es/ecma262/#prod-IdentifierName
func (comp *es_compiler) consume_identifier_name() string {
	return comp.tkh.ConsumeStrIfMatchesR(*es_regexp_ident)
}

// ------------------------------------------------------------------------------
//
// BEGINNING OF PARSER CODE & AST
//
// <About error conditions>
// Many parser functions will return pair of (<some node>, *es_syntax_error)
// - If <some node> isn't nil, it's success.
// - If <some node> is nil, and es_syntax_error is also nil, it means there was no match.
// - If <some node> is nil, but es_syntax_error isn't nil, match was found but there is a syntax error.
//
// ------------------------------------------------------------------------------

type es_syntax_error struct {
	message                es_syntax_error_msg
	cursor_from, cursor_to int
}

type es_syntax_error_msg string

const (
	es_syntax_error_msg_missing_left_paren  = es_syntax_error_msg("missing (")
	es_syntax_error_msg_missing_right_paren = es_syntax_error_msg("missing )")
	es_syntax_error_msg_missing_left_brace  = es_syntax_error_msg("missing {")
	es_syntax_error_msg_missing_right_brace = es_syntax_error_msg("missing }")
	es_syntax_error_msg_missing_semicolon   = es_syntax_error_msg("missing ;")
	es_syntax_error_msg_missing_colon       = es_syntax_error_msg("missing :")
	es_syntax_error_msg_missing_expression  = es_syntax_error_msg("missing expression")
)

func es_choose_longest_node(nodes []es_ast_node) es_ast_node {
	var longest_node es_ast_node

	for _, node := range nodes {
		if libcommon.IsNil(longest_node) || es_get_node_length(longest_node) < es_get_node_length(node) {
			longest_node = node
		}
	}
	return longest_node
}
func es_get_node_length(node es_ast_node) int {
	return node.get_cursor_to() - node.get_cursor_from()
}

var es_reserved_words = []string{
	"await", "break", "case", "catch", "class", "const", "continue", "debugger",
	"default", "delete", "do", "else", "enum", "export", "extends", "false",
	"finally", "for", "function", "if", "import", "in", "instanceof", "new",
	"null", "return", "super", "switch", "this", "throw", "true", "try",
	"typeof", "var", "void", "while", "with", "yield",
}

func (comp *es_compiler) consume_keyword(keyword string) bool {
	cursor_from := comp.tkh.Cursor

	if temp := comp.consume_identifier_name(); temp == "" || temp != keyword {
		comp.tkh.Cursor = cursor_from
		return false
	}
	return true
}

func (comp *es_compiler) consume_punctuator_with(punct string) bool {
	cursor_from := comp.tkh.Cursor

	if temp := comp.consume_punctuator(); temp == "" || temp != punct {
		comp.tkh.Cursor = cursor_from
		return false
	}
	return true
}

func (comp *es_compiler) consume_semicolon() bool {
	// TODO: Handle automatic semicolon insertion
	return comp.consume_punctuator_with(";")
}

// https://tc39.es/ecma262/#prod-IdentifierReference
func (comp *es_compiler) consume_identifier_reference_node(flags struct{ yield, await bool }) (res *es_ast_identifier_reference_node) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	name := comp.consume_identifier()
	if name == "" {
		return nil
	}
	// TODO: Yield, Await
	cursor_to := comp.tkh.Cursor
	return &es_ast_identifier_reference_node{
		cursor_from: cursor_from,
		cursor_to:   cursor_to,
		name:        name,
	}
}

// https://tc39.es/ecma262/#prod-BindingIdentifier
func (comp *es_compiler) consume_binding_identifier_node(flags struct{ yield, await bool }) (res *es_ast_binding_identifier_node) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	name := comp.consume_identifier()
	if name == "" {
		return nil
	}
	// TODO: Yield, Await
	cursor_to := comp.tkh.Cursor
	return &es_ast_binding_identifier_node{
		cursor_from: cursor_from,
		cursor_to:   cursor_to,
		name:        name,
	}
}

// https://tc39.es/ecma262/#sec-identifiers
//
// Returns empty string if not found or it's a reserved word.
func (comp *es_compiler) consume_identifier() string {
	cursor_from := comp.tkh.Cursor
	name := comp.consume_identifier_name()
	if slices.Contains(es_reserved_words, name) {
		comp.tkh.Cursor = cursor_from
		return ""
	}
	if name == "" {
		comp.tkh.Cursor = cursor_from
		return ""
	}
	return name
}

// https://tc39.es/ecma262/#sec-primary-expression
func (comp *es_compiler) consume_primary_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	// TODO: this

	if temp := comp.consume_identifier_reference_node(struct{ yield, await bool }{flags.yield, flags.await}); !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}
	if temp := comp.consume_literal_node(); !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
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
	if comp.consume_punctuator_with("(") {
		comp.skip_whitespaces_and_line_terminators(true)
		expr, serr := comp.consume_expression_node(struct{ in, yield, await bool }{true, flags.yield, flags.await})
		if serr != nil {
			return nil, serr
		} else if libcommon.IsNil(expr) {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_left_paren}
		}
		comp.skip_whitespaces_and_line_terminators(true)
		if !comp.consume_punctuator_with(")") {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_right_paren}
		}
		return &es_ast_paren_expr_node{
			cursor_from: cursor_from,
			cursor_to:   comp.tkh.Cursor,
			node:        expr,
		}, nil
	}

	// TODO: CoverParenthesizedExpressionAndArrowParameterList
	return es_choose_longest_node(candidates), nil
}

var es_regexp_hex_literal = regexp.MustCompile("^0[xX]")
var es_regexp_octal_literal_1 = regexp.MustCompile("^0o")
var es_regexp_octal_literal_2 = regexp.MustCompile("^0[0-7]")
var es_regexp_bin_literal = regexp.MustCompile("^0b")

// https://tc39.es/ecma262/#prod-Literal
func (comp *es_compiler) consume_literal_node() *es_ast_literal_node {
	cursor_from := comp.tkh.Cursor

	if temp := comp.consume_null_literal(); temp != "" {
		return &es_ast_literal_node{
			cursor_from: cursor_from,
			cursor_to:   comp.tkh.Cursor,
			value:       es_make_null_value(),
		}
	}
	if temp := comp.consume_boolean_literal(); temp != "" {
		switch temp {
		case "true":
			return &es_ast_literal_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				value:       es_make_boolean_value(true),
			}
		case "false":
			return &es_ast_literal_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				value:       es_make_boolean_value(false),
			}
		}
		panic("unreachable")
	}

	if temp := comp.consume_numeric_literal(); temp != "" {
		// Remove seaprator
		temp = strings.ReplaceAll(temp, "_", "")
		var value_f float64
		var value_i int64
		var err error
		is_int := false
		if es_regexp_hex_literal.MatchString(temp) {
			value_i, err = strconv.ParseInt(strings.TrimPrefix(strings.TrimPrefix(temp, "0x"), "0X"), 16, 64)
			if err != nil {
				panic(err)
			}
			is_int = true
		} else if es_regexp_octal_literal_1.MatchString(temp) {
			value_i, err = strconv.ParseInt(strings.TrimPrefix(temp, "0o"), 8, 64)
			if err != nil {
				panic(err)
			}
			is_int = true
		} else if es_regexp_octal_literal_2.MatchString(temp) && !es_regexp_non_octal_digits.MatchString(temp) {
			value_i, err = strconv.ParseInt(strings.TrimPrefix(temp, "0"), 8, 64)
			if err != nil {
				panic(err)
			}
			is_int = true
		} else if es_regexp_bin_literal.MatchString(temp) {
			value_i, err = strconv.ParseInt(strings.TrimPrefix(temp, "0b"), 2, 64)
			if err != nil {
				panic(err)
			}
			is_int = true
		} else {
			value_f, err = strconv.ParseFloat(temp, 64)
			if err != nil {
				panic(err)
			}
			is_int = false
		}
		if is_int {
			return &es_ast_literal_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				value:       es_make_number_value_i(value_i),
			}
		} else {
			return &es_ast_literal_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				value:       es_make_number_value_f(value_f),
			}
		}

	}
	if temp := comp.consume_string_literal(); temp != "" {
		panic("TODO")
	}
	return nil
}

// https://tc39.es/ecma262/#prod-MemberExpression
func (comp *es_compiler) consume_member_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	if temp, serr := comp.consume_primary_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}

	// TODO: MemberExpression [ Expression ]
	// TODO: MemberExpression . IdentifierName
	// TODO: MemberExpression TemplateLiteral
	// TODO: SuperProperty
	// TODO: MetaProperty
	// TODO: new MemberExpression Arguments
	// TODO: MemberExpression . PrivateIdentifier

	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-NewExpression
func (comp *es_compiler) consume_new_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	// STUB
	return comp.consume_member_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-CallExpression
func (comp *es_compiler) consume_call_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	if callee, serr := comp.consume_member_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(callee) {
		comp.skip_whitespaces_and_line_terminators(true)
		args, serr := comp.consume_arguments_node(struct{ yield, await bool }{flags.yield, flags.await})
		if serr != nil {
			return nil, serr
		} else if args == nil {
			return nil, nil
		}
		cursor_to := comp.tkh.Cursor
		res = es_ast_call_expr_node{
			cursor_from: cursor_from,
			cursor_to:   cursor_to,
			callee:      callee,
			args:        *args,
		}
		// Handle additional things that might follow the function all, like:
		//  foo(bar)(baz)[qux]
		for {
			if another_args, serr := comp.consume_arguments_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if another_args != nil {
				another_cursor_to := comp.tkh.Cursor
				res = es_ast_call_expr_node{
					cursor_from: cursor_from,
					cursor_to:   another_cursor_to,
					callee:      res,
					args:        *args,
				}
			} else if comp.consume_punctuator_with("[") {
				panic("TODO")
			} else if comp.consume_punctuator_with(".") {
				panic("TODO")
			} else {
				break
			}
		}
		return res, nil
	} else if comp.consume_keyword("super") {
		panic("TODO")
	} else if comp.consume_keyword("import") {
		panic("TODO")
	}
	return nil, nil
}

// https://tc39.es/ecma262/#prod-Arguments
func (comp *es_compiler) consume_arguments_node(flags struct{ yield, await bool }) (res *es_ast_arguments, serr *es_syntax_error) {
	var rest_args es_ast_node
	cursor_from := comp.tkh.Cursor
	defer func() {
		if res == nil {
			comp.tkh.Cursor = cursor_from
		}
	}()

	if !comp.consume_punctuator_with("(") {
		return nil, nil
	}

	args := []es_ast_node{}
	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("...") {
			comp.skip_whitespaces_and_line_terminators(true)
			if temp, serr := comp.consume_assignment_expression_node(struct{ in, yield, await bool }{true, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rest_args) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				rest_args = temp
			}
		} else if temp, serr := comp.consume_assignment_expression_node(struct{ in, yield, await bool }{true, flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if !libcommon.IsNil(temp) {
			args = append(args, temp)
		} else {
			break
		}
		comp.skip_whitespaces_and_line_terminators(true)
		if !comp.consume_punctuator_with(",") {
			break
		}
	}
	if !comp.consume_punctuator_with(")") {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_right_paren}
	}
	return &es_ast_arguments{
		cursor_from: cursor_from,
		cursor_to:   comp.tkh.Cursor,
		args:        args,
		rest_args:   rest_args,
	}, nil
}

// https://tc39.es/ecma262/#prod-LeftHandSideExpression
func (comp *es_compiler) consume_left_hand_side_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	if temp, serr := comp.consume_new_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}
	if temp, serr := comp.consume_call_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}
	// TODO: https://tc39.es/ecma262/#prod-OptionalExpression
	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-UpdateExpression
func (comp *es_compiler) consume_update_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	// STUB
	return comp.consume_left_hand_side_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-UnaryExpression
func (comp *es_compiler) consume_unary_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	if comp.consume_keyword("delete") {
		panic("TODO")
	} else if comp.consume_keyword("void") {
		panic("TODO")
	} else if comp.consume_keyword("typeof") {
		panic("TODO")
	} else if comp.consume_punctuator_with("+") {
		comp.skip_whitespaces_and_line_terminators(true)

		if node, serr := comp.consume_unary_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if libcommon.IsNil(node) {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
		} else {
			return es_ast_unary_op_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				node:        node,
				tp:          es_ast_unary_op_type_plus,
			}, nil
		}

	} else if comp.consume_punctuator_with("-") {
		comp.skip_whitespaces_and_line_terminators(true)
		if node, serr := comp.consume_unary_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if libcommon.IsNil(node) {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
		} else {
			return es_ast_unary_op_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				node:        node,
				tp:          es_ast_unary_op_type_neg,
			}, nil
		}
	} else if comp.consume_punctuator_with("~") {
		comp.skip_whitespaces_and_line_terminators(true)
		if node, serr := comp.consume_unary_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if libcommon.IsNil(node) {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
		} else {
			return es_ast_unary_op_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				node:        node,
				tp:          es_ast_unary_op_type_bnot,
			}, nil
		}
	} else if comp.consume_punctuator_with("!") {
		comp.skip_whitespaces_and_line_terminators(true)
		if node, serr := comp.consume_unary_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if libcommon.IsNil(node) {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
		} else {
			return es_ast_unary_op_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				node:        node,
				tp:          es_ast_unary_op_type_lnot,
			}, nil
		}
	} else if flags.await && comp.consume_keyword("await") {
		comp.skip_whitespaces_and_line_terminators(true)
		if node, serr := comp.consume_unary_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if libcommon.IsNil(node) {
			return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
		} else {
			return es_ast_unary_op_node{
				cursor_from: cursor_from,
				cursor_to:   comp.tkh.Cursor,
				node:        node,
				tp:          es_ast_unary_op_type_await,
			}, nil
		}
	}
	return comp.consume_update_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-ExponentiationExpression
func (comp *es_compiler) consume_exponentiation_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_update_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return comp.consume_unary_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	}
	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_punctuator_with("**") {
		return lhs, nil
	}
	comp.skip_whitespaces_and_line_terminators(true)
	rhs, serr := comp.consume_exponentiation_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(rhs) {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
	}
	return es_ast_binary_op_node{
		cursor_from: cursor_from,
		cursor_to:   comp.tkh.Cursor,
		lhs_node:    lhs,
		rhs_node:    rhs,
		tp:          es_ast_binary_op_type_exponent,
	}, nil
}

// https://tc39.es/ecma262/#prod-MultiplicativeExpression
func (comp *es_compiler) consume_multiplicative_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_exponentiation_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("*") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_exponentiation_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_multiply,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_div_punctuator() != "" {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_exponentiation_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_divide,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("%") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_exponentiation_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_modulo,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-AdditiveExpression
func (comp *es_compiler) consume_additive_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_multiplicative_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("+") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_multiplicative_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_add,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("-") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_multiplicative_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_subtract,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-ShiftExpression
func (comp *es_compiler) consume_shift_expression_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_additive_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("<<") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_additive_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_left_shift,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with(">>>") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_additive_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_right_lshift,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with(">>") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_additive_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_right_ashift,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-RelationalExpression
func (comp *es_compiler) consume_relational_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("<") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_less_than,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with(">") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_greater_than,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("<=") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_less_than_or_equal,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with(">=") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_greater_than_or_equal,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("instanceof") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_instanceof,
					rhs_node:    rhs,
				}
			}
		} else if flags.in && comp.consume_punctuator_with("in") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_shift_expression_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
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
func (comp *es_compiler) consume_equality_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_relational_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		// NOTE: We try longer ones first
		if comp.consume_punctuator_with("===") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_relational_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_strict_equal,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("!==") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_relational_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_strict_not_equal,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("==") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_relational_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_equal,
					rhs_node:    rhs,
				}
			}
		} else if comp.consume_punctuator_with("!=") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_relational_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_not_equal,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-BitwiseANDExpression
func (comp *es_compiler) consume_band_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_equality_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("&") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_equality_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_band,
					rhs_node:    rhs,
				}
			}

		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-BitwiseXORExpression
func (comp *es_compiler) consume_bxor_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_band_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("^") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_band_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_bxor,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-BitwiseORExpression
func (comp *es_compiler) consume_bor_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_bxor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("|") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_bxor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_bor,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-LogicalANDExpression
func (comp *es_compiler) consume_land_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_bor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("&&") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_bor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_land,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-LogicalORExpression
func (comp *es_compiler) consume_lor_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_land_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("||") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_land_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_lor,
					rhs_node:    rhs,
				}
			}

		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-CoalesceExpression
func (comp *es_compiler) consume_coalesce_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_bor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("??") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_bor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_binary_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					tp:          es_ast_binary_op_type_coalesce,
					rhs_node:    rhs,
				}
			}

		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-ShortCircuitExpression
func (comp *es_compiler) consume_short_circuit_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	candidates := []es_ast_node{}
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	if temp, serr := comp.consume_lor_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}

	if temp, serr := comp.consume_coalesce_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}
	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-ConditionalExpression
func (comp *es_compiler) consume_conditional_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	cond, serr := comp.consume_short_circuit_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(cond) {
		return nil, nil
	}
	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_punctuator_with("?") {
		return cond, nil
	}

	comp.skip_whitespaces_and_line_terminators(true)
	true_expr, serr := comp.consume_assignment_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(true_expr) {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
	}

	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_punctuator_with(":") {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_colon}
	}

	comp.skip_whitespaces_and_line_terminators(true)
	false_expr, serr := comp.consume_assignment_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(false_expr) {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
	}

	return es_ast_cond_expr_node{
		cursor_from: cursor_from,
		cursor_to:   comp.tkh.Cursor,
		cond_node:   cond,
		true_node:   true_expr,
		false_node:  false_expr,
	}, nil
}

// https://tc39.es/ecma262/#prod-AssignmentExpression
func (comp *es_compiler) consume_assignment_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	// STUB
	return comp.consume_conditional_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-Expression
func (comp *es_compiler) consume_expression_node(flags struct{ in, yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	lhs, serr := comp.consume_assignment_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(lhs) {
		return nil, nil
	}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with(",") {
			comp.skip_whitespaces_and_line_terminators(true)
			if rhs, serr := comp.consume_assignment_expression_node(struct{ in, yield, await bool }{flags.in, flags.yield, flags.await}); serr != nil {
				return nil, serr
			} else if libcommon.IsNil(rhs) {
				return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_expression}
			} else {
				lhs = es_ast_comma_op_node{
					cursor_from: cursor_from,
					cursor_to:   comp.tkh.Cursor,
					lhs_node:    lhs,
					rhs_node:    rhs,
				}
			}
		} else {
			break
		}
	}
	return lhs, nil
}

// https://tc39.es/ecma262/#prod-Statement
func (comp *es_compiler) consume_statement_node(flags struct{ yield, await, retrn bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	if temp, serr := comp.consume_expression_statement_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}

	if temp, serr := comp.consume_return_statement_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}

	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-Declaration
func (comp *es_compiler) consume_declaration_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	if temp, serr := comp.consume_hoistable_declaration_node(struct{ yield, await, defult bool }{flags.yield, flags.await, false}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}

	// TODO: ClassDeclaration
	// TODO: LexicalDeclaration

	return es_choose_longest_node(candidates), nil

}

// https://tc39.es/ecma262/#prod-HoistableDeclaration
func (comp *es_compiler) consume_hoistable_declaration_node(flags struct{ yield, await, defult bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	if temp, serr := comp.consume_function_declaration_node(struct{ yield, await, defult bool }{flags.yield, flags.await, flags.defult}); serr != nil {
		return nil, serr
	} else if !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}

	// TODO: GeneratorDeclaration
	// TODO: AsyncFunctionDeclaration
	// TODO: AsyncGeneratorDeclaration

	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-StatementList
func (comp *es_compiler) consume_statement_list(flags struct{ yield, await, retrn bool }) (res []es_ast_node, serr *es_syntax_error) {
	statements := []es_ast_node{}

	for {
		cursor_from := comp.tkh.Cursor
		comp.skip_whitespaces_and_line_terminators(true)
		candidates := []es_ast_node{}

		if temp, serr := comp.consume_statement_node(struct{ yield, await, retrn bool }{flags.yield, flags.await, flags.retrn}); serr != nil {
			return nil, serr
		} else if !libcommon.IsNil(temp) {
			candidates = append(candidates, temp)
			comp.tkh.Cursor = cursor_from
		}

		if temp, serr := comp.consume_declaration_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if !libcommon.IsNil(temp) {
			candidates = append(candidates, temp)
			comp.tkh.Cursor = cursor_from
		}

		item := es_choose_longest_node(candidates)
		if !libcommon.IsNil(item) {
			comp.tkh.Cursor = item.get_cursor_to()
			statements = append(statements, item)
		} else {
			break
		}
	}

	return statements, nil
}

// https://tc39.es/ecma262/#prod-BindingElement
func (comp *es_compiler) consume_binding_element_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}

	// https://tc39.es/ecma262/#prod-SingleNameBinding
	if temp := comp.consume_binding_identifier_node(struct{ yield, await bool }{flags.yield, flags.await}); !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
		// TODO: Accept initializer after the identifier (https://tc39.es/ecma262/#prod-Initializer)
	}

	// TODO: BindingPattern Initializer(opt)

	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-BindingRestElement
func (comp *es_compiler) consume_binding_rest_element_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	candidates := []es_ast_node{}
	if !comp.consume_punctuator_with("...") {
		return nil, nil
	}
	comp.skip_whitespaces_and_line_terminators(true)
	if temp := comp.consume_binding_identifier_node(struct{ yield, await bool }{flags.yield, flags.await}); !libcommon.IsNil(temp) {
		candidates = append(candidates, temp)
		comp.tkh.Cursor = cursor_from
	}
	// TODO: ...BindingPattern

	return es_choose_longest_node(candidates), nil
}

// https://tc39.es/ecma262/#prod-ExpressionStatement
func (comp *es_compiler) consume_expression_statement_node(flags struct{ yield, await bool }) (res *es_ast_expr_statement_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	found := false
	comp.tkh.Lookahead(func() {
		comp.skip_whitespaces_and_line_terminators(true)
		if comp.consume_punctuator_with("{") ||
			comp.consume_keyword("function") ||
			comp.consume_keyword("class") ||
			comp.consume_keyword("let") ||
			comp.consume_punctuator_with("[") {
			found = true
		}
		if comp.consume_keyword("async") {
			if comp.consume_line_terminator() == "" {
				found = true
			}
		}
	})
	if found {
		return nil, nil
	}
	expr, serr := comp.consume_expression_node(struct{ in, yield, await bool }{true, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(expr) {
		return nil, nil
	}
	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_semicolon() {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_semicolon}
	}
	return &es_ast_expr_statement_node{
		cursor_from: cursor_from,
		cursor_to:   comp.tkh.Cursor,
		node:        expr,
	}, nil
}

// https://tc39.es/ecma262/#prod-ReturnStatement
func (comp *es_compiler) consume_return_statement_node(flags struct{ yield, await bool }) (res *es_ast_return_statement_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	if !comp.consume_keyword("return") {
		return nil, nil
	}
	comp.skip_whitespaces_and_line_terminators(true)
	if comp.consume_semicolon() {
		return &es_ast_return_statement_node{
			cursor_from: cursor_from,
			cursor_to:   comp.tkh.Cursor,
			node:        nil,
		}, nil
	}
	comp.skip_whitespaces_and_line_terminators(false)
	expr, serr := comp.consume_expression_node(struct{ in, yield, await bool }{true, flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	} else if libcommon.IsNil(expr) {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_semicolon}
	}
	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_semicolon() {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_semicolon}
	}
	return &es_ast_return_statement_node{
		cursor_from: cursor_from,
		cursor_to:   comp.tkh.Cursor,
		node:        expr,
	}, nil
}

// https://tc39.es/ecma262/#prod-FormalParameters
func (comp *es_compiler) consume_formal_parameters(flags struct{ yield, await bool }) (res *es_ast_formal_parameters, serr *es_syntax_error) {
	params := []es_ast_node{}

	for {
		comp.skip_whitespaces_and_line_terminators(true)
		if temp, serr := comp.consume_formal_parameter_node(struct{ yield, await bool }{flags.yield, flags.await}); serr != nil {
			return nil, serr
		} else if libcommon.IsNil(temp) {
			break
		} else {
			params = append(params, temp)
		}
		comp.skip_whitespaces_and_line_terminators(true)
		if !comp.consume_punctuator_with(",") {
			break
		}
	}
	rest_param, serr := comp.consume_function_rest_parameter_node(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	}
	return &es_ast_formal_parameters{
		params:     params,
		rest_param: rest_param,
	}, nil
}

// https://tc39.es/ecma262/#prod-FormalParameter
func (comp *es_compiler) consume_formal_parameter_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	return comp.consume_binding_element_node(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-FunctionRestParameter
func (comp *es_compiler) consume_function_rest_parameter_node(flags struct{ yield, await bool }) (res es_ast_node, serr *es_syntax_error) {
	return comp.consume_binding_rest_element_node(struct{ yield, await bool }{flags.yield, flags.await})
}

// https://tc39.es/ecma262/#prod-FunctionDeclaration
func (comp *es_compiler) consume_function_declaration_node(flags struct{ yield, await, defult bool }) (res es_ast_node, serr *es_syntax_error) {
	cursor_from := comp.tkh.Cursor
	defer func() { comp.set_final_cursor(res, cursor_from) }()

	if !comp.consume_keyword("function") {
		return nil, nil
	}

	comp.skip_whitespaces_and_line_terminators(true)
	func_ident := comp.consume_binding_identifier_node(struct{ yield, await bool }{flags.yield, flags.await})
	if !flags.defult && func_ident == nil {
		return nil, nil
	}

	if !comp.consume_punctuator_with("(") {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_left_paren}
	}

	comp.skip_whitespaces_and_line_terminators(true)
	params, serr := comp.consume_formal_parameters(struct{ yield, await bool }{flags.yield, flags.await})
	if serr != nil {
		return nil, serr
	}

	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_punctuator_with(")") {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_left_paren}
	}

	comp.skip_whitespaces_and_line_terminators(true)
	if !comp.consume_punctuator_with("{") {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_left_brace}
	}

	comp.skip_whitespaces_and_line_terminators(true)
	statements, serr := comp.consume_statement_list(struct{ yield, await, retrn bool }{flags.yield, flags.await, true})
	if serr != nil {
		return nil, serr
	}

	comp.skip_whitespaces_and_line_terminators(true)
	if comp.consume_right_brace_punctuator() == "" {
		return nil, &es_syntax_error{cursor_from: cursor_from, cursor_to: comp.tkh.Cursor, message: es_syntax_error_msg_missing_right_brace}
	}

	return es_ast_function_decl_node{
		cursor_from: cursor_from,
		cursor_to:   comp.tkh.Cursor,
		ident:       *func_ident,
		params:      *params,
		body:        statements,
	}, nil
}

func (comp *es_compiler) set_final_cursor(node es_ast_node, cursor_initial int) {
	if !libcommon.IsNil(node) {
		comp.tkh.Cursor = node.get_cursor_to()
	} else {
		comp.tkh.Cursor = cursor_initial
	}
}

// Compiles given ECMAScript code to series of VM instructions.
func Compile(str string) (res []es_vm_instr, serr *es_syntax_error) {
	compiler := es_compiler{tkh: libcommon.TokenizerHelper{Str: []rune(str)}}
	nodes, serr := compiler.consume_statement_list(struct{ yield, await, retrn bool }{false, false, false})
	if serr != nil {
		return nil, serr
	}
	return es_make_code_for_ast_nodes(nodes), nil
}
