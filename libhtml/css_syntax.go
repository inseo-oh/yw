// Implementation of the CSS Syntax Module Level 3 (https://www.w3.org/TR/css-syntax-3/)
package libhtml

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	cm "yw/libcommon"
	"yw/libencoding"
)

type css_token_common struct{ cursor_from, cursor_to int }
type css_token interface {
	get_cursor_from() int
	get_cursor_to() int
	get_token_type() css_token_type
	String() string
}

type css_token_type uint8

const (
	css_token_type_eof = css_token_type(iota) // TODO: Remove this
	css_token_type_whitespace
	css_token_type_left_paren
	css_token_type_right_paren
	css_token_type_comma
	css_token_type_colon
	css_token_type_semicolon
	css_token_type_left_square_bracket
	css_token_type_right_square_bracket
	css_token_type_left_curly_bracket
	css_token_type_right_curly_bracket
	css_token_type_cdo
	css_token_type_cdc
	css_token_type_bad_string
	css_token_type_bad_url
	css_token_type_number
	css_token_type_percentage
	css_token_type_dimension
	css_token_type_string
	css_token_type_url
	css_token_type_at_keyword
	css_token_type_function
	css_token_type_ident
	css_token_type_hash
	css_token_type_delim
	// High-level objects ------------------------------------------------------
	css_token_type_ast_simple_block
	css_token_type_ast_function
	css_token_type_ast_qualified_rule
	css_token_type_ast_at_rule
	css_token_type_ast_declaration
)

func (typ css_token_type) String() string {
	switch typ {
	case css_token_type_whitespace:
		return "whitespace"
	case css_token_type_left_paren:
		return "left-paren"
	case css_token_type_right_paren:
		return "right-paren"
	case css_token_type_comma:
		return "comma"
	case css_token_type_colon:
		return "colon"
	case css_token_type_semicolon:
		return "semicolon"
	case css_token_type_left_square_bracket:
		return "left-square-bracket"
	case css_token_type_right_square_bracket:
		return "right-square-bracket"
	case css_token_type_left_curly_bracket:
		return "left-curly-bracket"
	case css_token_type_right_curly_bracket:
		return "right-curly-bracket"
	case css_token_type_cdo:
		return "cdo"
	case css_token_type_cdc:
		return "cdc"
	case css_token_type_bad_string:
		return "bad-string"
	case css_token_type_bad_url:
		return "bad-url"
	case css_token_type_number:
		return "number"
	case css_token_type_percentage:
		return "percentage"
	case css_token_type_dimension:
		return "dimension"
	case css_token_type_string:
		return "string"
	case css_token_type_url:
		return "url"
	case css_token_type_at_keyword:
		return "at-keyword"
	case css_token_type_function:
		return "function"
	case css_token_type_ident:
		return "ident"
	case css_token_type_hash:
		return "hash"
	case css_token_type_delim:
		return "delim"
	case css_token_type_ast_simple_block:
		return "simple_block"
	case css_token_type_ast_function:
		return "function"
	case css_token_type_ast_qualified_rule:
		return "qualified_rule"
	case css_token_type_ast_at_rule:
		return "at_rule"
	case css_token_type_ast_declaration:
		return "declaration"
	}
	return fmt.Sprintf("[unknown css_token_type %d]", typ)
}
func (t css_token_common) get_cursor_from() int {
	return t.cursor_from
}
func (t css_token_common) get_cursor_to() int {
	return t.cursor_to
}

type css_simple_token struct {
	css_token_common
	tp css_token_type
}

func (t css_simple_token) get_token_type() css_token_type { return t.tp }
func (t css_simple_token) String() string {
	switch t.tp {
	case css_token_type_whitespace:
		return " "
	case css_token_type_left_paren:
		return "("
	case css_token_type_right_paren:
		return ")"
	case css_token_type_comma:
		return ","
	case css_token_type_colon:
		return ":"
	case css_token_type_semicolon:
		return ";"
	case css_token_type_left_square_bracket:
		return "["
	case css_token_type_right_square_bracket:
		return "]"
	case css_token_type_left_curly_bracket:
		return "{"
	case css_token_type_right_curly_bracket:
		return "}"
	case css_token_type_cdo:
		return "<!--"
	case css_token_type_cdc:
		return "-->"
	case css_token_type_bad_string:
		return "/*bad-string*/"
	case css_token_type_bad_url:
		return "/*bad-url*/"
	}
	return fmt.Sprintf("<unknown css_simple_token type %v>", t.tp)
}

type css_number_token struct {
	css_token_common
	value css_number
}

func (t css_number_token) get_token_type() css_token_type { return css_token_type_number }
func (t css_number_token) String() string                 { return fmt.Sprintf("%v", t.value) }

type css_percentage_token struct {
	css_token_common
	value css_number
}

func (t css_percentage_token) get_token_type() css_token_type { return css_token_type_percentage }
func (t css_percentage_token) String() string                 { return fmt.Sprintf("%v%%", t.value) }

type css_dimension_token struct {
	css_token_common
	value css_number
	unit  string
}

func (t css_dimension_token) get_token_type() css_token_type { return css_token_type_dimension }
func (t css_dimension_token) String() string                 { return fmt.Sprintf("%v%s", t.value, t.unit) }

type css_string_token struct {
	css_token_common
	value string
}

func (t css_string_token) get_token_type() css_token_type { return css_token_type_string }
func (t css_string_token) String() string                 { return strconv.Quote(t.value) }

type css_url_token struct {
	css_token_common
	value string
}

func (t css_url_token) get_token_type() css_token_type { return css_token_type_url }
func (t css_url_token) String() string                 { return fmt.Sprintf("url(%s)", t.value) }

type css_at_keyword_token struct {
	css_token_common
	name string
}

func (t css_at_keyword_token) get_token_type() css_token_type { return css_token_type_at_keyword }
func (t css_at_keyword_token) String() string                 { return fmt.Sprintf("@%s", t.name) }

type css_function_token struct {
	css_token_common
	value string
}

func (t css_function_token) get_token_type() css_token_type { return css_token_type_function }
func (t css_function_token) String() string                 { return fmt.Sprintf("%s(", t.value) }

type css_ident_token struct {
	css_token_common
	value string
}

func (t css_ident_token) get_token_type() css_token_type { return css_token_type_ident }
func (t css_ident_token) String() string                 { return t.value }

type css_hash_token struct {
	css_token_common
	tp    css_hash_token_type
	value string
}

func (t css_hash_token) get_token_type() css_token_type { return css_token_type_hash }
func (t css_hash_token) String() string                 { return fmt.Sprintf("#%s/*%s*/", t.value, t.tp) }

type css_hash_token_type uint8

const (
	css_hash_token_type_id = css_hash_token_type(iota)
	css_hash_token_type_unrestricted
)

func (typ css_hash_token_type) String() string {
	switch typ {
	case css_hash_token_type_id:
		return "id"
	case css_hash_token_type_unrestricted:
		return "unrestricted"
	}
	return fmt.Sprintf("[unknown css_hash_token_type %d]", typ)
}

type css_delim_token struct {
	css_token_common
	value rune
}

func (t css_delim_token) get_token_type() css_token_type { return css_token_type_delim }
func (t css_delim_token) String() string                 { return fmt.Sprintf("%c", t.value) }

func css_tokenize(src string) ([]css_token, error) {
	tkh := cm.TokenizerHelper{Str: []rune(src)}

	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#ident-start-code-point
	is_ident_start_code_point := func(chr rune) bool {
		return cm.AsciiAlphaRegex.MatchString(string(chr)) ||
			0x80 <= chr ||
			chr == '_'
	}
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#ident-code-point
	is_ident_code_point := func(chr rune) bool {
		return is_ident_start_code_point(chr) ||
			cm.AsciiDigitRegex.MatchString(string(chr)) ||
			chr == '-'
	}
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#check-if-three-code-points-would-start-an-ident-sequence
	is_valid_ident_start_sequence := func(s string) bool {
		cps := []rune(s)
		if len(cps) == 0 {
			return false
		}
		if is_ident_start_code_point(cps[0]) {
			return true
		}
		switch cps[0] {
		case '-':
			if 1 < len(cps) && is_ident_code_point(cps[1]) {
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

	consume_comments := func() {
		end_found := false
		for !tkh.IsEof() {
			if tkh.ConsumeStrIfMatches("/*", 0) == "" {
				return
			}
			for !tkh.IsEof() {
				if tkh.ConsumeStrIfMatches("*/", 0) != "" {
					end_found = true
					break
				}
				tkh.ConsumeChar()
			}
			if end_found {
				continue
			}
			// PARSE ERROR: Reached EOF without closing the comment.
			return
		}
	}

	// Returns nil if not found
	consume_number := func() *css_number {
		start_cursor := tkh.Cursor
		have_integer_part := false
		have_fractional_part := false
		out := css_number{}

		// Note that we don't parse the number directly - We only check if it's a valid CSS number.
		// Rest of the job is handled by the standard library.

		// Sign ----------------------------------------------------------------

		// Integer part --------------------------------------------------------
		for !tkh.IsEof() {
			temp_char := tkh.PeekChar()
			if cm.AsciiDigitRegex.MatchString(string(temp_char)) {
				tkh.ConsumeChar()
				have_integer_part = true
			} else {
				break
			}
		}
		// Decimal point -------------------------------------------------------
		old_cursor := tkh.Cursor
		if tkh.ConsumeCharIfMatches('.') != -1 {
			// Fractional part -------------------------------------------------
			digit_count := 0

			for !tkh.IsEof() {
				temp_char := tkh.PeekChar()
				if cm.AsciiDigitRegex.MatchString(string(temp_char)) {
					tkh.ConsumeChar()
					digit_count++
				} else {
					break
				}
			}
			if !have_integer_part && digit_count == 0 {
				tkh.Cursor = old_cursor
				return nil
			}
			out.tp = css_number_type_float
			have_fractional_part = true
		}

		if !have_integer_part && !have_fractional_part {
			// We have invalid number
			tkh.Cursor = start_cursor
			return nil
		}

		// Exponent indicator --------------------------------------------------
		old_cursor = tkh.Cursor
		if tkh.ConsumeCharIfMatchesOneOf("eE") != -1 {
			digit_count := 0

			// Exponent sign ---------------------------------------------------
			tkh.ConsumeCharIfMatchesOneOf("+-")

			// Exponent --------------------------------------------------------
			for !tkh.IsEof() {
				temp_char := tkh.PeekChar()
				if cm.AsciiDigitRegex.MatchString(string(temp_char)) {
					tkh.ConsumeChar()
				} else {
					break
				}
			}
			if digit_count == 0 {
				tkh.Cursor = old_cursor
			}
		}

		end_cursor := tkh.Cursor

		// Now we parse the number ---------------------------------------------
		temp_buf := strings.Builder{}
		tkh.Cursor = start_cursor
		for tkh.Cursor < end_cursor {
			temp_buf.WriteRune(tkh.ConsumeChar())
		}
		// TODO: Check for range errors
		if out.tp == css_number_type_float {
			v, err := strconv.ParseFloat(temp_buf.String(), 64)
			if err != nil {
				log.Panic(err)
			}
			out.value = v
		} else {
			v, err := strconv.ParseInt(temp_buf.String(), 10, 64)
			if err != nil {
				log.Panic(err)
			}
			out.value = v
		}

		return &out
	}

	// Returns nil if not found
	consume_escaped_code_point := func() *rune {
		old_cursor := tkh.Cursor
		if tkh.ConsumeCharIfMatches('\\') == -1 {
			return nil
		}
		is_hex_digit := false
		hex_digit_val := 0
		hex_digit_count := 0

		if tkh.IsEof() {
			// PARSE ERROR: Unexpected EOF
			cp := rune(0xfffd)
			return &cp
		}
		if tkh.ConsumeCharIfMatches('\n') != -1 {
			tkh.Cursor = old_cursor
			return nil
		}
		// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-an-escaped-code-point
		for !tkh.IsEof() && hex_digit_count < 6 {
			temp_char := tkh.PeekChar()
			digit := 0
			if cm.AsciiDigitRegex.MatchString(string(temp_char)) {
				digit = int(temp_char - '0')
			} else if cm.AsciiLowerHexDigitRegex.MatchString(string(temp_char)) {
				digit = int(temp_char - 'a' + 10)
			} else if cm.AsciiUpperHexDigitRegex.MatchString(string(temp_char)) {
				digit = int(temp_char - 'A' + 10)
			} else {
				break
			}
			tkh.ConsumeChar()
			hex_digit_val = hex_digit_val*16 + digit
			is_hex_digit = true
			hex_digit_count++
		}
		var out rune
		if is_hex_digit {
			out = rune(hex_digit_val)
		} else {
			out = tkh.ConsumeChar()
		}
		return &out
	}

	// Returns nil if not found
	consume_ident_sequence := func(must_start_with_ident_start bool) *string {
		sb := strings.Builder{}

		for !tkh.IsEof() {
			old_cursor := tkh.Cursor

			var result_chr rune
			if temp := consume_escaped_code_point(); temp == nil {
				result_chr = tkh.ConsumeChar()
			} else {
				result_chr = *temp
			}
			if is_ident_start_code_point(result_chr) ||
				((sb.Len() != 0 || !must_start_with_ident_start) && is_ident_code_point(result_chr)) {
				sb.WriteRune(result_chr)
			} else {
				tkh.Cursor = old_cursor
				break
			}
		}

		if sb.Len() == 0 {
			return nil
		}
		return cm.MakeStrPtr(sb.String())
	}

	// Returns nil if not found
	consume_whitespace_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
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
		return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_whitespace}, nil
	}

	// Returns nil if not found
	consume_simple_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
		switch tkh.ConsumeCharIfMatchesOneOf("(),:;[]{}") {
		case '(':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_left_paren}, nil
		case ')':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_right_paren}, nil
		case ',':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_comma}, nil
		case ':':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_colon}, nil
		case ';':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_semicolon}, nil
		case '[':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_left_square_bracket}, nil
		case ']':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_right_square_bracket}, nil
		case '{':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_left_curly_bracket}, nil
		case '}':
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_right_curly_bracket}, nil
		case -1:
		default:
			panic("unreachable")
		}
		switch tkh.ConsumeStrIfMatchesOneOf([]string{"<!--", "-->"}, 0) {
		case "<!--":
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_cdo}, nil
		case "-->":
			return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_cdc}, nil
		case "":
		default:
			panic("unreachable")
		}
		return nil, nil
	}

	// Returns nil if not found
	consume_string_token := func() (css_token, error) {
		var ending_char rune
		sb := strings.Builder{}

		switch tkh.ConsumeCharIfMatchesOneOf("\"'") {
		case '"':
			ending_char = '"'
		case '\'':
			ending_char = '\''
		default:
			return nil, nil
		}

	loop:
		for !tkh.IsEof() {
			switch tkh.ConsumeCharIfMatchesOneOf(string(ending_char) + "\n") {
			case ending_char:
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
				var result_chr rune
				if temp := consume_escaped_code_point(); temp != nil {
					result_chr = *temp
				} else {
					result_chr = tkh.ConsumeChar()
				}
				sb.WriteRune(result_chr)
			}
		}
		return css_string_token{
			value: sb.String(),
		}, nil
	}

	// Returns nil if not found
	consume_hash_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
		if tkh.ConsumeCharIfMatches('#') == -1 {
			return nil, nil
		}
		var hash_value string
		if temp := consume_ident_sequence(false); temp != nil {
			hash_value = *temp
		} else {
			return nil, errors.New("expected identifier after '#'")
		}
		if len(hash_value) == 0 {
			return nil, errors.New("expected identifier after '#'")
		}
		var subtype css_hash_token_type
		if is_valid_ident_start_sequence(hash_value) {
			subtype = css_hash_token_type_id
		} else {
			subtype = css_hash_token_type_unrestricted
		}
		return css_hash_token{css_token_common{cursor_from, tkh.Cursor}, subtype, hash_value}, nil
	}

	// Returns nil if not found
	consume_at_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
		if tkh.ConsumeCharIfMatches('@') == -1 {
			return nil, nil
		}
		var at_value string
		if temp := consume_ident_sequence(true); temp != nil {
			at_value = *temp
		} else {
			return nil, errors.New("expected identifier after '@'")
		}
		if len(at_value) == 0 || !is_valid_ident_start_sequence(at_value) {
			return nil, errors.New("expected identifier after '@'")
		}
		return css_at_keyword_token{css_token_common{cursor_from, tkh.Cursor}, at_value}, nil
	}

	// Returns nil if not found
	consume_numeric_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
		var num css_number
		if temp := consume_number(); temp != nil {
			num = *temp
		} else {
			return nil, nil
		}
		old_cursor := tkh.Cursor

		if ident := consume_ident_sequence(true); ident != nil {
			if is_valid_ident_start_sequence(*ident) {
				return css_dimension_token{css_token_common{cursor_from, tkh.Cursor}, num, *ident}, nil
			} else {
				tkh.Cursor = old_cursor
			}
		}
		if tkh.ConsumeCharIfMatches('%') != -1 {
			return css_percentage_token{css_token_common{cursor_from, tkh.Cursor}, num}, nil
		} else {
			return css_number_token{css_token_common{cursor_from, tkh.Cursor}, num}, nil
		}
	}

	consume_remnants_of_bad_url := func() {
		for !tkh.IsEof() {
			if tkh.ConsumeCharIfMatches(')') != -1 {
				break
			}
			if consume_escaped_code_point() == nil {
				tkh.ConsumeChar()
			}
		}
	}

	// Returns nil if not found
	consume_ident_like_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
		var ident string

		if temp := consume_ident_sequence(true); temp != nil {
			ident = *temp
		} else {
			return nil, nil
		}
		if cm.ToAsciiLowercase(ident) == "url" && tkh.ConsumeCharIfMatches('(') != -1 {
			for tkh.ConsumeStrIfMatches("  ", 0) != "" {
			}
			old_cursor := tkh.Cursor
			if tkh.ConsumeCharIfMatchesOneOf("\"'") != -1 ||
				tkh.ConsumeStrIfMatches(" \"", 0) != "" ||
				tkh.ConsumeStrIfMatches(" '", 0) != "" {
				// Function token ----------------------------------------------
				tkh.Cursor = old_cursor
				return css_function_token{css_token_common{cursor_from, tkh.Cursor}, ident}, nil
			} else {
				// URL token ---------------------------------------------------
				consume_whitespace_token()
				url_sb := strings.Builder{}
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
						consume_remnants_of_bad_url()
						return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_bad_url}, nil
					default:
						var escaped_chr rune
						if temp := consume_escaped_code_point(); temp != nil {
							escaped_chr = *temp
						} else if tkh.ConsumeCharIfMatches('\\') != -1 {
							// PARSE ERROR: Unexpected character after \
							consume_remnants_of_bad_url()
							return css_simple_token{css_token_common{cursor_from, tkh.Cursor}, css_token_type_bad_url}, nil
						} else {
							escaped_chr = tkh.ConsumeChar()
						}
						url_sb.WriteRune(escaped_chr)
					}
				}
				return css_url_token{css_token_common{cursor_from, tkh.Cursor}, url_sb.String()}, nil
			}
		} else if tkh.ConsumeCharIfMatches('(') != -1 {
			return css_function_token{css_token_common{cursor_from, tkh.Cursor}, ident}, nil
		} else {
			return css_ident_token{css_token_common{cursor_from, tkh.Cursor}, ident}, nil
		}
	}

	consume_token := func() (css_token, error) {
		cursor_from := tkh.Cursor
		handlers := []func() (css_token, error){
			consume_whitespace_token,
			consume_string_token,
			consume_hash_token,
			consume_at_token,
			consume_simple_token,
			consume_numeric_token,
			consume_ident_like_token,
		}

		consume_comments()
		for _, h := range handlers {
			res, err := h()
			if !cm.IsNil(res) {
				return res, nil
			} else if err != nil {
				return nil, err
			}
		}
		if tkh.IsEof() {
			return nil, nil
		} else {
			c := tkh.ConsumeChar()
			return css_delim_token{css_token_common{cursor_from, tkh.Cursor}, c}, nil
		}
	}

	out := []css_token{}
	for {
		tk, err := consume_token()
		if cm.IsNil(tk) && err != nil {
			return nil, err
		} else if cm.IsNil(tk) {
			break
		}
		out = append(out, tk)

	}

	return out, nil
}

type css_ast_simple_block_token struct {
	css_token_common
	tp   css_ast_simple_block_type
	body []css_token
}

func (t css_ast_simple_block_token) String() string {
	sb := strings.Builder{}
	switch t.tp {
	case css_ast_simple_block_type_curly:
		sb.WriteRune('{')
	case css_ast_simple_block_type_square:
		sb.WriteRune('[')
	case css_ast_simple_block_type_paren:
		sb.WriteRune('(')
	}
	for _, tk := range t.body {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	switch t.tp {
	case css_ast_simple_block_type_curly:
		sb.WriteRune('}')
	case css_ast_simple_block_type_square:
		sb.WriteRune(']')
	case css_ast_simple_block_type_paren:
		sb.WriteRune(')')
	}
	return sb.String()
}

type css_ast_simple_block_type uint8

const (
	css_ast_simple_block_type_square = css_ast_simple_block_type(iota)
	css_ast_simple_block_type_curly
	css_ast_simple_block_type_paren
)

func (n css_ast_simple_block_token) get_token_type() css_token_type {
	return css_token_type_ast_simple_block
}

type css_ast_function_token struct {
	css_token_common
	name  string
	value []css_token
}

func (t css_ast_function_token) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s(", t.name))
	for _, tk := range t.value {
		sb.WriteString(fmt.Sprintf("%v", tk))
	}
	sb.WriteString(")")
	return sb.String()
}

func (t css_ast_function_token) get_token_type() css_token_type {
	return css_token_type_ast_function
}

type css_ast_qualified_rule_token struct {
	css_token_common
	prelude []css_token
	body    []css_token
}

func (t css_ast_qualified_rule_token) String() string {
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
func (t css_ast_qualified_rule_token) get_token_type() css_token_type {
	return css_token_type_ast_qualified_rule
}

type css_ast_at_rule_token struct {
	css_token_common
	name    string
	prelude []css_token // NOTE: This is just STUB -- We would want actual parsed value.
	body    []css_token // NOTE: This is just STUB -- We would want actual parsed value.
}

func (t css_ast_at_rule_token) String() string {
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
func (t css_ast_at_rule_token) get_token_type() css_token_type {
	return css_token_type_ast_at_rule
}

type css_ast_declaration_token struct {
	css_token_common
	name      string
	value     []css_token
	important bool
}

func (t css_ast_declaration_token) String() string {
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
func (t css_ast_declaration_token) get_token_type() css_token_type {
	return css_token_type_ast_declaration
}

type css_token_stream struct {
	tokens []css_token
	cursor int
}

func (ts *css_token_stream) is_end() bool {
	return len(ts.tokens) <= ts.cursor
}
func (ts *css_token_stream) consume_token() css_token {
	if ts.is_end() {
		return nil
	}
	ts.cursor++
	res := ts.tokens[ts.cursor-1]
	return res
}
func (ts *css_token_stream) consume_token_with_type(tp css_token_type) css_token {
	if ts.is_end() {
		return nil
	}
	old_cursor := ts.cursor
	tk := ts.consume_token()
	if tk.get_token_type() != tp {
		ts.cursor = old_cursor
		return nil
	}
	return tk
}

func (ts *css_token_stream) skip_whitespaces() {
	for !cm.IsNil(ts.consume_token_with_type(css_token_type_whitespace)) {
	}
}

func (ts *css_token_stream) consume_delim_token_with(d rune) *css_delim_token {
	old_cursor := ts.cursor
	tk := ts.consume_token_with_type(css_token_type_delim)
	if tk == nil || tk.(css_delim_token).value != d {
		ts.cursor = old_cursor
		return nil
	}
	v := tk.(css_delim_token)
	return &v
}
func (ts *css_token_stream) consume_ident_token_with(i string) *css_ident_token {
	old_cursor := ts.cursor
	tk := ts.consume_token_with_type(css_token_type_ident)
	if tk == nil || tk.(css_ident_token).value != i {
		ts.cursor = old_cursor
		return nil
	}
	v := tk.(css_ident_token)
	return &v
}
func (ts *css_token_stream) consume_ast_simple_block_with(tp css_ast_simple_block_type) *css_ast_simple_block_token {
	old_cursor := ts.cursor
	n := ts.consume_token_with_type(css_token_type_ast_simple_block)
	if cm.IsNil(n) {
		ts.cursor = old_cursor
		return nil
	}
	blk := n.(css_ast_simple_block_token)
	if blk.tp != tp {
		ts.cursor = old_cursor
		return nil
	}
	return &blk
}
func (ts *css_token_stream) consume_ast_function_with(name string) *css_ast_function_token {
	old_cursor := ts.cursor
	tk := ts.consume_token_with_type(css_token_type_ast_function)
	if tk == nil || tk.(css_ast_function_token).name != name {
		ts.cursor = old_cursor
		return nil
	}
	v := tk.(css_ast_function_token)
	return &v
}

// Returns nil if not found
func (ts *css_token_stream) consume_preserved_token() css_token {
	old_cursor := ts.cursor

	tk := ts.consume_token()
	if cm.IsNil(tk) {
		return nil
	}
	switch tk.get_token_type() {
	case css_token_type_function,
		css_token_type_left_curly_bracket,
		css_token_type_left_square_bracket,
		css_token_type_left_paren:
		ts.cursor = old_cursor
		return nil
	}
	return tk
}

// Returns nil if not found
func (ts *css_token_stream) consume_simple_block(open_token_type, close_token_type css_token_type) *css_ast_simple_block_token {
	res_nodes := []css_token{}
	var block_type css_ast_simple_block_type
	switch open_token_type {
	case css_token_type_left_curly_bracket:
		block_type = css_ast_simple_block_type_curly
	case css_token_type_left_square_bracket:
		block_type = css_ast_simple_block_type_square
	case css_token_type_left_paren:
		block_type = css_ast_simple_block_type_paren
	default:
		panic("unsupported open_token_type")
	}

	open_token := ts.consume_token_with_type(open_token_type)
	if cm.IsNil(open_token) {
		return nil
	}
	var close_token css_token
	for {
		temp_tk := ts.consume_component_value()
		if cm.IsNil(temp_tk) || temp_tk.get_token_type() == close_token_type {
			close_token = temp_tk
			break
		}
		res_nodes = append(res_nodes, temp_tk)
	}
	if cm.IsNil(close_token) {
		return nil
	}
	return &css_ast_simple_block_token{
		css_token_common{open_token.get_cursor_from(), close_token.get_cursor_to()},
		block_type, res_nodes,
	}
}

// Returns nil if not found
func (ts *css_token_stream) consume_curly_block() *css_ast_simple_block_token {
	return ts.consume_simple_block(css_token_type_left_curly_bracket, css_token_type_right_curly_bracket)
}

// Returns nil if not found
func (ts *css_token_stream) consume_square_block() *css_ast_simple_block_token {
	return ts.consume_simple_block(css_token_type_left_square_bracket, css_token_type_right_square_bracket)
}

// Returns nil if not found
func (ts *css_token_stream) consume_paren_block() *css_ast_simple_block_token {
	return ts.consume_simple_block(css_token_type_left_paren, css_token_type_right_paren)
}

// Returns nil if not found
func (ts *css_token_stream) consume_function() *css_ast_function_token {
	fn_value_nodes := []css_token{}
	var fn_token css_function_token
	if temp := ts.consume_token_with_type(css_token_type_function); !cm.IsNil(temp) {
		fn_token = temp.(css_function_token)
	} else {
		return nil
	}
	var close_token css_token
	for {
		temp_tk := ts.consume_component_value()
		if cm.IsNil(temp_tk) || temp_tk.get_token_type() == css_token_type_right_paren {
			close_token = temp_tk
			break
		}
		fn_value_nodes = append(fn_value_nodes, temp_tk)
	}

	return &css_ast_function_token{
		css_token_common{fn_token.get_cursor_from(), close_token.get_cursor_to()},
		fn_token.value, fn_value_nodes,
	}
}

// Returns nil if not found
func (ts *css_token_stream) consume_component_value() css_token {
	if res := ts.consume_curly_block(); res != nil {
		return *res
	}
	if res := ts.consume_square_block(); res != nil {
		return *res
	}
	if res := ts.consume_paren_block(); res != nil {
		return *res
	}
	if res := ts.consume_function(); res != nil {
		return *res
	}
	if res := ts.consume_preserved_token(); res != nil {
		return res
	}
	return nil
}

// Returns nil if not found
func (ts *css_token_stream) consume_qualified_rule() *css_ast_qualified_rule_token {
	prelude := []css_token{}

	for {
		block := ts.consume_curly_block()
		if block != nil {
			return &css_ast_qualified_rule_token{
				css_token_common{block.cursor_from, block.cursor_to},
				prelude,
				block.body,
			}
		} else if ts.is_end() {
			return nil
		}
		prelude = append(prelude, ts.consume_component_value())
	}
}

// Returns nil if not found
func (ts *css_token_stream) consume_at_rule() *css_ast_at_rule_token {
	prelude := []css_token{}
	var kwd_token css_at_keyword_token
	if temp := ts.consume_token_with_type(css_token_type_at_keyword); !cm.IsNil(temp) {
		kwd_token = temp.(css_at_keyword_token)
	} else {
		return nil
	}

	for {
		block := ts.consume_curly_block()
		if block != nil {
			return &css_ast_at_rule_token{
				css_token_common{block.cursor_from, block.cursor_to},
				kwd_token.name,
				prelude,
				block.body,
			}
		} else if ts.is_end() {
			return nil
		}
		prelude = append(prelude, ts.consume_component_value())
	}
}

// Returns nil if not found
func (ts *css_token_stream) consume_declaration() *css_ast_declaration_token {
	// <name>  :  contents  !important -----------------------------------------
	var ident_token css_ident_token
	if temp := ts.consume_token_with_type(css_token_type_ident); temp != nil {
		ident_token = temp.(css_ident_token)
	} else {
		return nil
	}
	decl_name := ident_token.value
	decl_value := []css_token{}
	decl_is_important := false
	// name<  >:  contents  !important -----------------------------------------
	for !cm.IsNil(ts.consume_token_with_type(css_token_type_whitespace)) {
	}
	// name  <:>  contents  !important -----------------------------------------
	if ts.consume_token_with_type(css_token_type_colon) == nil {
		// Parse error
		return nil
	}
	// name  :<  >contents  !important -----------------------------------------
	for !cm.IsNil(ts.consume_token_with_type(css_token_type_whitespace)) {
	}

	// name  :  <contents  !important> -----------------------------------------
	for {
		temp_tk := ts.consume_component_value()
		if cm.IsNil(temp_tk) {
			break
		}
		decl_value = append(decl_value, temp_tk)
	}
	last_node := decl_value[len(decl_value)-1]
	if 2 <= len(decl_value) {
		// See if we have !important
		ptk1 := decl_value[len(decl_value)-2]
		ptk2 := decl_value[len(decl_value)-1]
		if ptk1.get_token_type() == css_token_type_delim && ptk1.(css_delim_token).value == '!' &&
			ptk2.get_token_type() == css_token_type_ident && ptk2.(css_ident_token).value == "important" {
			decl_value = decl_value[:len(decl_value)-2]
			decl_is_important = true
		}
	}
	return &css_ast_declaration_token{
		css_token_common{ident_token.cursor_from, last_node.get_cursor_to()},
		decl_name,
		decl_value,
		decl_is_important,
	}
}

func (ts *css_token_stream) consume_declaration_list() []css_token {
	decls := []css_token{}

	for {
		token := ts.consume_token()

		if cm.IsNil(token) {
			break
		} else if token.get_token_type() == css_token_type_whitespace || token.get_token_type() == css_token_type_semicolon {
			continue
		} else if token.get_token_type() == css_token_type_at_keyword {
			ts.cursor--
			rule := ts.consume_at_rule()
			if rule == nil {
				panic("unreachable")
			}
			decls = append(decls, rule)
		} else if token.get_token_type() == css_token_type_ident {
			temp_stream := css_token_stream{}
			temp_stream.tokens = append(temp_stream.tokens, token)
			for {
				token = ts.consume_token()
				if cm.IsNil(token) || token.get_token_type() == css_token_type_semicolon {
					break
				}
				temp_stream.tokens = append(temp_stream.tokens, token)
			}
			decl := temp_stream.consume_declaration()
			if decl == nil {
				break
			}
			decls = append(decls, *decl)
		} else {
			// PARSE ERROR
			for {
				old_cursor := ts.cursor
				token = ts.consume_token()
				if cm.IsNil(token) || token.get_token_type() == css_token_type_semicolon {
					break
				}
				ts.cursor = old_cursor
				ts.consume_component_value()
			}
		}
	}
	return decls
}
func (ts *css_token_stream) consume_style_block_contents() []css_token {
	decls := []css_token{}
	rules := []css_ast_qualified_rule_token{}

	for {
		token := ts.consume_token()

		if cm.IsNil(token) {
			break
		} else if token.get_token_type() == css_token_type_whitespace || token.get_token_type() == css_token_type_semicolon {
			continue
		} else if token.get_token_type() == css_token_type_at_keyword {
			ts.cursor--
			rule := ts.consume_at_rule()
			if rule == nil {
				panic("unreachable")
			}
			decls = append(decls, rule)
		} else if token.get_token_type() == css_token_type_ident {
			temp_stream := css_token_stream{}
			temp_stream.tokens = append(temp_stream.tokens, token)
			for {
				token = ts.consume_token()
				if cm.IsNil(token) || token.get_token_type() == css_token_type_semicolon {
					break
				}
				temp_stream.tokens = append(temp_stream.tokens, token)
			}
			decl := temp_stream.consume_declaration()
			if decl == nil {
				break
			}
			decls = append(decls, *decl)
		} else if token.get_token_type() == css_token_type_delim && token.(css_delim_token).value == '&' {
			ts.cursor--
			if rule := ts.consume_qualified_rule(); rule != nil {
				rules = append(rules, *rule)
			}
		} else {
			// PARSE ERROR
			for {
				old_cursor := ts.cursor
				token = ts.consume_token()
				if cm.IsNil(token) || token.get_token_type() == css_token_type_semicolon {
					break
				}
				ts.cursor = old_cursor
				ts.consume_component_value()
			}
		}
	}
	for _, rule := range rules {
		decls = append(decls, rule)
	}
	return decls
}
func (ts *css_token_stream) consume_list_of_rules(top_level_flag bool) []css_token {
	rules := []css_token{}

	for {
		token := ts.consume_token()

		if cm.IsNil(token) {
			break
		} else if token.get_token_type() == css_token_type_whitespace || token.get_token_type() == css_token_type_semicolon {
			continue
		} else if token.get_token_type() == css_token_type_cdo || token.get_token_type() == css_token_type_cdc {
			if top_level_flag {
				continue
			}
			ts.cursor--
			if rule := ts.consume_qualified_rule(); rule != nil {
				rules = append(rules, *rule)
			}
		} else if token.get_token_type() == css_token_type_at_keyword {
			ts.cursor--
			rule := ts.consume_at_rule()
			if rule == nil {
				panic("unreachable")
			}
			rules = append(rules, *rule)
		} else {
			ts.cursor--
			if rule := ts.consume_qualified_rule(); rule != nil {
				rules = append(rules, *rule)
			}
		}
	}
	return rules
}

func css_parse_comma_separated_list_of_component_values(tokens []css_token) [][]css_token {
	value_lists := [][]css_token{}
	temp_list := []css_token{}

	stream := css_token_stream{tokens: tokens}
	for {
		value := stream.consume_component_value()
		if cm.IsNil(value) || value.get_token_type() == css_token_type_comma {
			value_lists = append(value_lists, temp_list)
			temp_list = temp_list[:0]
			if value.get_token_type() != css_token_type_comma {
				break
			}
			continue
		}
		temp_list = append(temp_list, value)
	}
	return value_lists
}
func css_parse_list_of_component_values(tokens []css_token) []css_token {
	temp_list := []css_token{}

	stream := css_token_stream{tokens: tokens}
	for {
		value := stream.consume_component_value()
		if cm.IsNil(value) || value.get_token_type() == css_token_type_comma {
			break
		}
		temp_list = append(temp_list, value)
	}
	return temp_list
}

func css_parse_component_value(tokens []css_token) css_token {
	ts := css_token_stream{tokens: tokens}
	ts.skip_whitespaces()
	if !ts.is_end() {
		panic("TODO: syntax error: expected component value")
	}
	value := ts.consume_component_value()
	ts.skip_whitespaces()
	if ts.is_end() {
		panic("TODO: syntax error: expected eof")
	}
	return value
}

func css_parse_list_of_declarations(tokens []css_token) []css_token {
	stream := css_token_stream{tokens: tokens}
	return stream.consume_declaration_list()
}
func css_parse_style_block_contents(tokens []css_token) []css_token {
	stream := css_token_stream{tokens: tokens}
	return stream.consume_style_block_contents()
}
func css_parse_declaration(tokens []css_token) *css_ast_declaration_token {
	stream := css_token_stream{tokens: tokens}
	stream.skip_whitespaces()
	if cm.IsNil(stream.consume_token_with_type(css_token_type_ident)) {
		panic("TODO: syntax error: expected identifier")
	}
	stream.cursor--
	node := stream.consume_declaration()
	if node == nil {
		panic("TODO: syntax error: expected declaration")
	}
	return node
}
func css_parse_rule(tokens []css_token) css_token {
	ts := css_token_stream{tokens: tokens}
	ts.skip_whitespaces()
	var res css_token
	res = ts.consume_at_rule()
	if cm.IsNil(res) {
		res = ts.consume_qualified_rule()
	}
	if cm.IsNil(res) {
		panic("TODO: syntax error: expected at-rule or qualified-rule")
	}
	if ts.is_end() {
		panic("TODO: syntax error: expected eof")
	}
	return res
}
func css_parse_list_of_rules(tokens []css_token) []css_token {
	stream := css_token_stream{tokens: tokens}
	return stream.consume_list_of_rules(false)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#typedef-declaration-value
func (ts *css_token_stream) _consume_declaration_value(any_value bool) []css_token {
	out := []css_token{}
	open_block_tokens := []css_token_type{}
	for {
		old_cursor := ts.cursor
		tk := ts.consume_token()
		if cm.IsNil(tk) {
			break
		}
		tk_type := tk.get_token_type()
		if tk_type == css_token_type_bad_string ||
			tk_type == css_token_type_bad_url ||
			// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#typedef-any-value
			(!any_value && (tk_type == css_token_type_semicolon ||
				(tk_type == css_token_type_delim && tk.(css_delim_token).value == '!'))) {
			ts.cursor = old_cursor
			break
		}
		// If we have block opening token, push it to the stack.
		if tk_type == css_token_type_left_paren ||
			tk_type == css_token_type_left_square_bracket ||
			tk_type == css_token_type_left_curly_bracket {
			open_block_tokens = append(open_block_tokens, tk_type)
		}
		// If we have block closing token, see if we have unmatched token.
		if (tk_type == css_token_type_right_paren &&
			(len(open_block_tokens) == 0 ||
				open_block_tokens[len(open_block_tokens)-1] != css_token_type_left_paren)) ||
			(tk_type == css_token_type_right_square_bracket &&
				(len(open_block_tokens) == 0 ||
					open_block_tokens[len(open_block_tokens)-1] != css_token_type_left_square_bracket)) ||
			(tk_type == css_token_type_right_curly_bracket &&
				(len(open_block_tokens) == 0 ||
					open_block_tokens[len(open_block_tokens)-1] != css_token_type_left_curly_bracket)) {
			ts.cursor = old_cursor
			break
		}
		out = append(out, tk)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func (ts *css_token_stream) consume_declaration_value() []css_token {
	return ts._consume_declaration_value(false)
}
func (ts *css_token_stream) consume_any_value() []css_token {
	return ts._consume_declaration_value(true)
}

// This can be used to parse where a CSS syntax can be repeated separated by comma.
// If max_repeats is 0, repeat count is unlimited.
//
// https://www.w3.org/TR/css-values-4/#mult-comma
func css_accept_comma_separated_repetion[T any](ts *css_token_stream, max_repeats int, parser func(ts *css_token_stream) (T, error)) ([]T, error) {
	res := []T{}
	for {
		x, err := parser(ts)
		if cm.IsNil(x) {
			if err != nil {
				return nil, err
			} else if len(res) != 0 {
				return nil, errors.New("expected something after ','")
			} else {
				break
			}
		}
		res = append(res, x)
		if max_repeats != 0 && max_repeats <= len(res) {
			break
		}
		ts.skip_whitespaces()
		if cm.IsNil(ts.consume_token_with_type(css_token_type_comma)) {
			break
		}
		ts.skip_whitespaces()
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// This can be used to parse where a CSS syntax can be repeated multiple times.
// If max_repeats is 0, repeat count is unlimited.
//
// https://www.w3.org/TR/css-values-4/#mult-num-range
//
// TODO: There is a similar pattern in css_background.go.. replace it with this.
func css_accept_repetion[T any](ts *css_token_stream, max_repeats int, parser func(ts *css_token_stream) (T, error)) ([]T, error) {
	res := []T{}
	for {
		x, err := parser(ts)
		if cm.IsNil(x) {
			if err != nil {
				return nil, err
			} else {
				break
			}
		}
		res = append(res, x)
		if max_repeats != 0 && max_repeats <= len(res) {
			break
		}
		ts.skip_whitespaces()
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-parse-something-according-to-a-css-grammar
func css_parse[T any](tokens []css_token, parser func(ts *css_token_stream) (T, error)) (T, error) {
	comp_values := css_parse_list_of_component_values(tokens)
	stream := css_token_stream{tokens: comp_values}
	return parser(&stream)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-decode-bytes
func css_decode_bytes(bytes []byte) string {
	fallback := css_determine_fallback_encoding(bytes)
	input := libencoding.IoQueueFromSlice(bytes)
	output := libencoding.IoQueueFromSlice[rune](nil)
	libencoding.Decode(&input, fallback, &output)
	return string(libencoding.IoQueueToSlice[rune](output))
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#determine-the-fallback-encoding
func css_determine_fallback_encoding(bytes []byte) libencoding.Type {
	// Check if HTTP or equivalent protocol provides an encoding label ---------
	// TODO

	// Check '@charset "xxx";' byte sequence -----------------------------------
	if 1024 < len(bytes) {
		bytes = bytes[:1024]
	}
	// NOTE: Below sequence of bytes are '@charset "' in ASCII
	if 10 <= len(bytes) || !slices.Equal([]byte{0x40, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x20, 0x22}, bytes[:10]) {
		remaining_bytes := bytes
		found_end := false
		encoding_name := []rune{}
		for 0 < len(remaining_bytes) {
			// NOTE: Below sequence of bytes are '";' in ASCII
			if 2 <= len(bytes) && slices.Equal([]byte{0x22, 0x3b}, bytes[:2]) {
				found_end = true
				break
			}
			encoding_name = append(encoding_name, rune(remaining_bytes[0]))
			remaining_bytes = remaining_bytes[1:]
		}
		if found_end {
			encoding, err := libencoding.GetEncodingFromLabel(string(encoding_name))
			if err == nil {
				if encoding == libencoding.Utf16Be || encoding == libencoding.Utf16Le {
					// This is not a bug. The standard says to do this.
					return libencoding.Utf8
				}
				return encoding
			}
		}
	}
	// Check if environment encoding is provided -------------------------------
	// TODO

	// Fallback to UTF-8 -------------------------------------------------------
	return libencoding.Utf8
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#css-stylesheets
func css_parse_stylesheet(input []css_token, location *string) css_stylesheet {
	stylesheet := css_stylesheet{
		location: location,
	}
	ts := css_token_stream{tokens: input}
	rule_nodes := ts.consume_list_of_rules(true)

	// Parse top-level qualified rules as style rules
	stylesheet.style_rules = css_parse_style_rules_from_nodes(rule_nodes)

	return stylesheet
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#style-rules
func css_parse_style_rules_from_nodes(rule_nodes []css_token) []css_style_rule {
	style_rules := []css_style_rule{}
	print_raw_rule_nodes := false
	for _, n := range rule_nodes {
		if n.get_token_type() != css_token_type_ast_qualified_rule {
			continue
		}
		qrule := n.(css_ast_qualified_rule_token)
		prelude_stream := css_token_stream{tokens: qrule.prelude}
		selector_list, err := prelude_stream.parse_selector_list()
		if err != nil {
			// TODO: Report error
			log.Printf("selector parsing error: %v", err)
			continue
		}
		if len(selector_list) == 0 {
			log.Println("FIXME: having no selectors after successfully parsing it isn't possible")
			print_raw_rule_nodes = true
			continue
		}
		contents_stream := css_token_stream{tokens: qrule.body}
		contents := contents_stream.consume_style_block_contents()
		decls := []css_declaration{}
		at_rules := []css_at_rule{}
		for _, content := range contents {
			if content.get_token_type() == css_token_type_ast_declaration {
				decl_node := content.(css_ast_declaration_token)

				prop_desc, ok := css_property_descriptors_map[cm.ToAsciiLowercase(decl_node.name)]
				if !ok {
					log.Printf("unknown property name: %v", decl_node.name)
					continue
				}
				inner_as := css_token_stream{tokens: decl_node.value}
				inner_as.skip_whitespaces()
				value, ok := prop_desc.parse_func(&inner_as)
				if !ok {
					log.Printf("bad value for property: %v (token list: %v)", decl_node.name, inner_as.tokens)
					continue
				}
				inner_as.skip_whitespaces()
				if !inner_as.is_end() {
					log.Printf("extra junk at the end for property: %v (token list: %v)", decl_node.name, inner_as.tokens[inner_as.cursor:])
					continue
				}
				decls = append(decls, css_declaration{decl_node.name, value, decl_node.important})
			} else if content.get_token_type() == css_token_type_ast_at_rule {
				rule_node := content.(css_ast_at_rule_token)
				at_rules = append(at_rules, css_at_rule{rule_node.name, rule_node.prelude, rule_node.body})
			} else {
				log.Printf("warning: unexpected node with type %v found while parsing style block contents", content.get_token_type())
			}
		}

		style_rules = append(style_rules, css_style_rule{selector_list, decls, at_rules})
	}
	if print_raw_rule_nodes {
		log.Println("=============== BEGIN: Raw rule nodes ===============")
		log.Println(rule_nodes)
		log.Println("=============== END:   Raw rule nodes ===============")
	}
	return style_rules
}

type css_property_value interface {
	String() string
}
type css_property_descriptor struct {
	initial    css_property_value
	parse_func func(ts *css_token_stream) (css_property_value, bool)
	apply_func func(dest *css_computed_style_set, value any)
}

// var css_property_descriptors_map = map[string]css_property_descriptor{}
