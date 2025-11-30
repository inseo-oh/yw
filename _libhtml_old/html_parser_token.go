//go:generate go run ./html_entities_gen
package libhtml

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	cm "yw/util"
)

type html_token interface {
	equals(other html_token) bool // XXX: Do we need this to be here?
	String() string
}

type html_eof_token struct{}

func (tk html_eof_token) equals(other html_token) bool {
	_, ok := other.(*html_eof_token)
	return ok
}
func (tk html_eof_token) String() string {
	return "<!EOF>"
}

type html_char_token struct{ value rune }

func (tk html_char_token) equals(other html_token) bool {
	other_tk, ok := other.(*html_char_token)
	if !ok {
		return false
	}
	if tk.value != other_tk.value {
		return false
	}
	return true
}
func (tk html_char_token) is_char_token_with_one_of(chars string) bool {
	for _, c := range chars {
		if tk.value == c {
			return true
		}
	}
	return false
}

func (tk html_char_token) String() string {
	return fmt.Sprintf("%c", tk.value)
}

type html_comment_token struct{ data string }

func (tk html_comment_token) equals(other html_token) bool {
	other_tk, ok := other.(*html_comment_token)
	if !ok {
		return false
	}
	if tk.data != other_tk.data {
		return false
	}
	return true
}
func (tk html_comment_token) String() string {
	return fmt.Sprintf("<!-- %s -->", strconv.Quote(tk.data))
}

type html_doctype_token struct {
	// nil value represents "missing"
	name         *string
	public_id    *string
	system_id    *string
	force_quirks bool
}

func (tk html_doctype_token) equals(other html_token) bool {
	other_tk, ok := other.(*html_doctype_token)
	if !ok {
		return false
	}
	if (tk.name == nil) != (other_tk.name == nil) {
		return false
	}
	if tk.name != nil && *tk.name != *other_tk.name {
		return false
	}
	return true
}
func (tk html_doctype_token) String() string {
	sb := strings.Builder{}
	sb.WriteString("<!DOCTYPE")
	if tk.name != nil {
		sb.WriteString(fmt.Sprintf(" %s", strconv.Quote(*tk.name)))
	}
	if tk.public_id != nil {
		sb.WriteString(fmt.Sprintf(" PUBLIC %s", strconv.Quote(*tk.public_id)))
	}
	if tk.system_id != nil {
		sb.WriteString(fmt.Sprintf(" SYSTEM %s", strconv.Quote(*tk.system_id)))
	}
	sb.WriteString(">")
	return sb.String()
}

type html_tag_token struct {
	is_end          bool
	is_self_closing bool
	tag_name        string
	attrs           []dom_Attr_s

	self_closing_acknowledged bool
}

func (tk html_tag_token) equals(other html_token) bool {
	other_tk, ok := other.(*html_tag_token)
	if !ok {
		return false
	}
	if len(tk.attrs) != len(other_tk.attrs) {
		return false
	}
	if tk.is_end != other_tk.is_end {
		return false
	}
	for i := 0; i < len(tk.attrs); i++ {
		attr := tk.attrs[i]
		other_tk := other_tk.attrs[i]
		if attr.local_name != other_tk.local_name {
			return false
		}
		if attr.value != other_tk.value {
			return false
		}
	}
	if tk.is_self_closing != other_tk.is_self_closing {
		return false
	}
	if tk.tag_name != other_tk.tag_name {
		return false
	}
	return true
}

func (tk html_tag_token) is_start_tag() bool {
	return !tk.is_end
}
func (tk html_tag_token) is_end_tag() bool {
	return tk.is_end
}

// Returns nil if there's no such attribute
func (tk html_tag_token) get_attr(name string) *string {
	for _, attr := range tk.attrs {
		if attr.local_name == name {
			return cm.MakeStrPtr(attr.value)
		}
	}
	return nil
}

func (tk html_tag_token) String() string {
	sb := strings.Builder{}
	if tk.is_end {
		sb.WriteString("</")
	} else {
		sb.WriteString("<")
	}
	sb.WriteString(tk.tag_name)
	for _, attr := range tk.attrs {
		sb.WriteString(" ")
		val := strconv.Quote(attr.get_value())
		if ns, ok := attr.get_namespace(); ok {
			sb.WriteString(fmt.Sprintf("%v:%s=%s", ns, attr.get_local_name(), val))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s", attr.get_local_name(), val))
		}
	}
	if tk.is_self_closing {
		sb.WriteString("/")
	}
	sb.WriteString(">")
	return sb.String()
}

type html_tokenizer_state uint8

const (
	html_tokenizer_data_state = html_tokenizer_state(iota)
	html_tokenizer_rcdata_state
	html_tokenizer_rawtext_state
	html_tokenizer_plaintext_state
	html_tokenizer_tag_open_state
	html_tokenizer_end_tag_open_state
	html_tokenizer_tag_name_state
	html_tokenizer_rcdata_less_than_sign_state
	html_tokenizer_rcdata_end_tag_open_state
	html_tokenizer_rcdata_end_tag_name_state
	html_tokenizer_rawtext_less_than_sign_state
	html_tokenizer_rawtext_end_tag_open_state
	html_tokenizer_rawtext_end_tag_name_state
	html_tokenizer_before_attribute_name_state
	html_tokenizer_attribute_name_state
	html_tokenizer_after_attribute_name_state
	html_tokenizer_before_attribute_value_state
	html_tokenizer_attribute_value_double_quoted_state
	html_tokenizer_attribute_value_single_quoted_state
	html_tokenizer_attribute_value_unquoted_state
	html_tokenizer_after_attribute_value_quoted_state
	html_tokenizer_self_closing_start_tag_state
	html_tokenizer_bogus_comment_state
	html_tokenizer_markup_declaration_open_state
	html_tokenizer_comment_start_state
	html_tokenizer_comment_start_dash_state
	html_tokenizer_comment_state
	html_tokenizer_comment_less_than_sign_state
	html_tokenizer_comment_end_dash_state
	html_tokenizer_comment_end_state
	html_tokenizer_doctype_state
	html_tokenizer_before_doctype_name_state
	html_tokenizer_doctype_name_state
	html_tokenizer_after_doctype_name_state
	html_tokenizer_character_reference_state
	html_tokenizer_named_character_reference_state
	html_tokenizer_numeric_character_reference_state
	html_tokenizer_hexadecimal_character_reference_start_state
	html_tokenizer_decimal_character_reference_start_state
	html_tokenizer_hexadecimal_character_reference_state
	html_tokenizer_decimal_character_reference_state
	html_tokenizer_numeric_character_reference_end_state
)

type html_parse_error string

const (
	html_absence_of_digits_in_numeric_character_reference_error = html_parse_error("absence-of-digits-in-numeric-character-reference")
	html_abrupt_closing_of_empty_comment_error                  = html_parse_error("abrupt-closing-of-empty-comment")
	html_character_reference_outside_unicode_range_error        = html_parse_error("character-reference-outside-unicode-range")
	html_control_character_reference_error                      = html_parse_error("control-character-reference")
	html_eof_before_tag_name_error                              = html_parse_error("eof-before-tag-name")
	html_eof_in_comment_error                                   = html_parse_error("eof-in-comment")
	html_eof_in_doctype_error                                   = html_parse_error("eof-in-doctype")
	html_eof_in_tag_error                                       = html_parse_error("eof-in-tag")
	html_incorrectly_opened_comment_error                       = html_parse_error("incorrectly-opened-comment")
	html_invalid_character_sequence_after_doctype_name_error    = html_parse_error("invalid-character-sequence-after-doctype-name")
	html_invalid_first_character_of_tag_name_error              = html_parse_error("invalid-first-character-of-tag-name")
	html_missing_attribute_value_error                          = html_parse_error("missing-attribute-value")
	html_missing_doctype_name_error                             = html_parse_error("missing-doctype-name")
	html_missing_end_tag_name_error                             = html_parse_error("missing-end-tag-name")
	html_missing_semicolon_after_character_reference_error      = html_parse_error("missing-semicolon-after-character-reference")
	html_missing_whitespace_before_doctype_name_error           = html_parse_error("missing-whitespace-before-doctype-name")
	html_missing_whitepace_between_attributes_error             = html_parse_error("missing-whitespace-between-attributes")
	html_noncharacter_reference_error                           = html_parse_error("noncharacter-character-reference")
	html_null_character_reference_error                         = html_parse_error("null-character-reference")
	html_surrogate_character_reference_error                    = html_parse_error("surrogate-character-reference")
	html_unexpected_character_in_attribute_name_error           = html_parse_error("unexpected-character-in-attribute-name")
	html_unexpected_character_in_unquoted_attribute_value_error = html_parse_error("unexpected-character-in-unquoted-attribute-value")
	html_unexpected_equals_sign_before_attribute_name_error     = html_parse_error("unexpected-equals-sign-before-attribute-name")
	html_unexpected_null_character_error                        = html_parse_error("unexpected-null-character")
	html_unexpected_question_mark_instead_of_tag_name_error     = html_parse_error("unexpected-question-mark-instead-of-tag-name")
	html_unexpected_solidus_in_tag_error                        = html_parse_error("unexpected-solidus-in-tag")
)

type html_tokenizer struct {
	tkh                 cm.TokenizerHelper
	state               html_tokenizer_state
	parser_pause_flag   bool
	last_start_tag_name string
	on_token_emitted    func(tk html_token)
}

func html_make_tokenizer(str string) html_tokenizer {
	return html_tokenizer{
		tkh: cm.TokenizerHelper{
			Str: []rune(str),
		},
	}
}

func (t *html_tokenizer) parse_error_encountered(err html_parse_error) {
	fmt.Println(err)
}

func (t *html_tokenizer) run() {
	var curr_tok html_token
	var return_state html_tokenizer_state
	temp_buf := ""
	character_reference_code := 0

	attrs_to_remove := []int{}
	curr_tag_token := func() *html_tag_token {
		tok := curr_tok.(*html_tag_token)
		return tok
	}
	curr_comment_token := func() *html_comment_token {
		tok := curr_tok.(*html_comment_token)
		return tok
	}
	curr_doctype_token := func() *html_doctype_token {
		tok := curr_tok.(*html_doctype_token)
		return tok
	}
	curr_attr := func() *dom_Attr_s {
		attrs := curr_tag_token().attrs
		return &attrs[len(attrs)-1]
	}
	check_duplicate_attr_name := func() {
		attrs := curr_tag_token().attrs
		for i, attr := range attrs {
			if i == (len(attrs) - 1) {
				continue
			}
			if curr_attr().local_name == attr.local_name {
				attrs_to_remove = append(attrs_to_remove, i)
			}
		}
	}
	emit_token := func(tk html_token) {
		if tag_tk, ok := tk.(*html_tag_token); ok {
			final_attrs := []dom_Attr_s{}
			for i, attr := range curr_tag_token().attrs {
				bad_attr := slices.Contains(attrs_to_remove, i)
				if !bad_attr {
					final_attrs = append(final_attrs, attr)
				}
			}
			curr_tag_token().attrs = final_attrs
			if tag_tk.is_start_tag() {
				t.last_start_tag_name = tag_tk.tag_name
			}
		}
		t.on_token_emitted(tk)
	}
	is_consumed_as_part_of_attr := func() bool {
		switch return_state {
		case html_tokenizer_attribute_value_double_quoted_state,
			html_tokenizer_attribute_value_single_quoted_state,
			html_tokenizer_attribute_value_unquoted_state:
			return true
		default:
			return false
		}
	}
	flush_code_points_consumed_as_char_reference := func() {
		if is_consumed_as_part_of_attr() {
			for _, c := range temp_buf {
				curr_attr().value += string(c)
			}
		} else {
			for _, c := range temp_buf {
				emit_token(&html_char_token{value: c})
			}
		}
	}
	is_appropriate_end_tag_token := func(tk html_tag_token) bool {
		return t.last_start_tag_name == tk.tag_name
	}

	for {
		if t.parser_pause_flag {
			return
		}

		switch t.state {
		// https://html.spec.whatwg.org/multipage/parsing.html#data-state
		case html_tokenizer_data_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '&':
				return_state = html_tokenizer_data_state
				t.state = html_tokenizer_character_reference_state
			case '<':
				t.state = html_tokenizer_tag_open_state
			case 0x0000:
				t.parse_error_encountered(html_unexpected_null_character_error)
				emit_token(&html_char_token{value: next_char})
			case -1:
				emit_token(&html_eof_token{})
				return
			default:
				emit_token(&html_char_token{value: next_char})
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-state
		case html_tokenizer_rcdata_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '&':
				return_state = html_tokenizer_data_state
				t.state = html_tokenizer_character_reference_state
			case '<':
				t.state = html_tokenizer_rcdata_less_than_sign_state
			case 0x0000:
				t.parse_error_encountered(html_unexpected_null_character_error)
				emit_token(&html_char_token{value: 0xfffd})
			case -1:
				emit_token(&html_eof_token{})
				return
			default:
				emit_token(&html_char_token{value: next_char})
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-state
		case html_tokenizer_rawtext_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '<':
				t.state = html_tokenizer_rawtext_less_than_sign_state
			case 0x0000:
				t.parse_error_encountered(html_unexpected_null_character_error)
				emit_token(&html_char_token{value: 0xfffd})
			case -1:
				emit_token(&html_eof_token{})
				return
			default:
				emit_token(&html_char_token{value: next_char})
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#plaintext-state
		case html_tokenizer_plaintext_state:
			panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#plaintext-state]")
		// https://html.spec.whatwg.org/multipage/parsing.html#tag-open-state
		case html_tokenizer_tag_open_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '!':
				t.state = html_tokenizer_markup_declaration_open_state
			case '/':
				t.state = html_tokenizer_end_tag_open_state
			case '?':
				t.parse_error_encountered(html_unexpected_question_mark_instead_of_tag_name_error)
				curr_tok = &html_comment_token{data: ""}
				t.tkh.Cursor--
				t.state = html_tokenizer_bogus_comment_state
			case -1:
				emit_token(&html_char_token{value: '<'})
				emit_token(&html_eof_token{})
				return
			default:
				if cm.AsciiAlphaRegex.MatchString(string(next_char)) {
					curr_tok = &html_tag_token{
						is_end:          false,
						is_self_closing: false,
						tag_name:        "",
						attrs:           []dom_Attr_s{},
					}
					t.tkh.Cursor--
					t.state = html_tokenizer_tag_name_state
				} else {
					t.parse_error_encountered(html_invalid_first_character_of_tag_name_error)
					emit_token(&html_char_token{value: '<'})
					t.tkh.Cursor--
					t.state = html_tokenizer_data_state
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#end-tag-open-state
		case html_tokenizer_end_tag_open_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '>':
				t.parse_error_encountered(html_missing_end_tag_name_error)
				t.state = html_tokenizer_data_state
			case -1:
				emit_token(&html_char_token{value: '<'})
				emit_token(&html_char_token{value: '/'})
				emit_token(&html_eof_token{})
				return
			default:
				if cm.AsciiAlphaRegex.MatchString(string(next_char)) {
					curr_tok = &html_tag_token{
						is_end:   true,
						tag_name: "",
					}
					t.tkh.Cursor--
					t.state = html_tokenizer_tag_name_state
				} else {
					t.parse_error_encountered(html_invalid_first_character_of_tag_name_error)
					curr_tag_token().tag_name += string(rune(0xfffd))
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#tag-name-state
		case html_tokenizer_tag_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				t.state = html_tokenizer_before_attribute_name_state
			case '/':
				t.state = html_tokenizer_self_closing_start_tag_state
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case 0x0000:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_tag_token().tag_name += string(rune(0xfffd))
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
				return
			default:
				if cm.AsciiUppercaseRegex.MatchString(string(next_char)) {
					curr_tag_token().tag_name += strings.ToLower(string(next_char))
				} else {
					curr_tag_token().tag_name += string(next_char)
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-less-than-sign-state
		case html_tokenizer_rcdata_less_than_sign_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '/':
				temp_buf = ""
				t.state = html_tokenizer_rcdata_end_tag_open_state
			default:
				emit_token(&html_char_token{value: '<'})
				t.tkh.Cursor--
				t.state = html_tokenizer_rcdata_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-open-state
		case html_tokenizer_rcdata_end_tag_open_state:
			next_char := t.tkh.ConsumeChar()
			if cm.AsciiAlphaRegex.MatchString(string(next_char)) {
				curr_tok = &html_tag_token{
					is_end:   true,
					tag_name: "",
				}
				t.tkh.Cursor--
				t.state = html_tokenizer_rcdata_end_tag_name_state
			} else {
				emit_token(&html_char_token{value: '<'})
				emit_token(&html_char_token{value: '/'})
				t.tkh.Cursor--
				t.state = html_tokenizer_rcdata_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rcdata-end-tag-name-state
		case html_tokenizer_rcdata_end_tag_name_state:
			next_char := t.tkh.ConsumeChar()
			anything_else := func() {
				emit_token(&html_char_token{value: '<'})
				emit_token(&html_char_token{value: '/'})
				for _, c := range temp_buf {
					emit_token(&html_char_token{value: c})
				}
				t.tkh.Cursor--
				t.state = html_tokenizer_rcdata_state
			}
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				if is_appropriate_end_tag_token(*curr_tag_token()) {
					t.state = html_tokenizer_before_attribute_name_state
				} else {
					anything_else()
				}
			case '/':
				if is_appropriate_end_tag_token(*curr_tag_token()) {
					t.state = html_tokenizer_self_closing_start_tag_state
				} else {
					anything_else()
				}
			case '>':
				if is_appropriate_end_tag_token(*curr_tag_token()) {
					t.state = html_tokenizer_data_state
					emit_token(curr_tok)
				} else {
					anything_else()
				}
			default:
				if cm.AsciiUppercaseRegex.MatchString(string(next_char)) {
					curr_tag_token().tag_name += strings.ToLower(string(next_char))
					temp_buf += string(next_char)
				} else if cm.AsciiLowercaseRegex.MatchString(string(next_char)) {
					curr_tag_token().tag_name += string(next_char)
					temp_buf += string(next_char)
				} else {
					anything_else()
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-less-than-sign-state
		case html_tokenizer_rawtext_less_than_sign_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '/':
				temp_buf = ""
				t.state = html_tokenizer_rawtext_end_tag_open_state
			default:
				emit_token(&html_char_token{value: '<'})
				t.tkh.Cursor--
				t.state = html_tokenizer_rawtext_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-end-tag-open-state
		case html_tokenizer_rawtext_end_tag_open_state:
			next_char := t.tkh.ConsumeChar()
			if cm.AsciiAlphaRegex.MatchString(string(next_char)) {
				curr_tok = &html_tag_token{
					is_end:   true,
					tag_name: "",
				}
				t.tkh.Cursor--
				t.state = html_tokenizer_rawtext_end_tag_name_state
			} else {
				emit_token(&html_char_token{value: '<'})
				emit_token(&html_char_token{value: '/'})
				t.tkh.Cursor--
				t.state = html_tokenizer_rawtext_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#rawtext-end-tag-name-state
		case html_tokenizer_rawtext_end_tag_name_state:
			next_char := t.tkh.ConsumeChar()
			anything_else := func() {
				emit_token(&html_char_token{value: '<'})
				emit_token(&html_char_token{value: '/'})
				for _, c := range temp_buf {
					emit_token(&html_char_token{value: c})
				}
				t.tkh.Cursor--
				t.state = html_tokenizer_rawtext_state
			}
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				if is_appropriate_end_tag_token(*curr_tag_token()) {
					t.state = html_tokenizer_before_attribute_name_state
				} else {
					anything_else()
				}
			case '/':
				if is_appropriate_end_tag_token(*curr_tag_token()) {
					t.state = html_tokenizer_self_closing_start_tag_state
				} else {
					anything_else()
				}
			case '>':
				if is_appropriate_end_tag_token(*curr_tag_token()) {
					t.state = html_tokenizer_data_state
					emit_token(curr_tok)
				} else {
					anything_else()
				}
			default:
				if cm.AsciiUppercaseRegex.MatchString(string(next_char)) {
					curr_tag_token().tag_name += strings.ToLower(string(next_char))
					temp_buf += string(next_char)
				} else if cm.AsciiLowercaseRegex.MatchString(string(next_char)) {
					curr_tag_token().tag_name += string(next_char)
					temp_buf += string(next_char)
				} else {
					anything_else()
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-name-state
		case html_tokenizer_before_attribute_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
			case '/', '>', -1:
				t.tkh.Cursor--
				t.state = html_tokenizer_after_attribute_name_state
			case '=':
				t.parse_error_encountered(html_unexpected_equals_sign_before_attribute_name_error)
				curr_tag_token().attrs = append(curr_tag_token().attrs, dom_Attr_s{
					local_name: string(next_char),
					value:      "",
				})
				attrs_to_remove = []int{}
				t.state = html_tokenizer_attribute_name_state
			default:
				curr_tag_token().attrs = append(curr_tag_token().attrs, dom_Attr_s{
					local_name: "",
					value:      "",
				})
				attrs_to_remove = []int{}
				t.tkh.Cursor--
				t.state = html_tokenizer_attribute_name_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-name-state
		case html_tokenizer_attribute_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ', '/', '>':
				t.tkh.Cursor--
				t.state = html_tokenizer_after_attribute_name_state
				check_duplicate_attr_name()
			case '=':
				t.state = html_tokenizer_before_attribute_value_state
				check_duplicate_attr_name()
			case 0x0000:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_attr().local_name += string(rune(0xfffd))
			case '"', '\'', '<':
				t.parse_error_encountered(html_unexpected_character_in_attribute_name_error)
				fallthrough
			default:
				curr_attr().local_name += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-name-state
		case html_tokenizer_after_attribute_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
			case '/':
				t.state = html_tokenizer_self_closing_start_tag_state
			case '=':
				t.state = html_tokenizer_before_attribute_value_state
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
			default:
				curr_tag_token().attrs = append(curr_tag_token().attrs, dom_Attr_s{
					local_name: "",
					value:      "",
				})
				attrs_to_remove = []int{}
				t.tkh.Cursor--
				t.state = html_tokenizer_attribute_name_state

			}
		// https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-value-state
		case html_tokenizer_before_attribute_value_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
			case '"':
				t.state = html_tokenizer_attribute_value_double_quoted_state
			case '\'':
				t.state = html_tokenizer_attribute_value_single_quoted_state
			case '>':
				t.parse_error_encountered(html_missing_attribute_value_error)
				t.state = html_tokenizer_data_state
			default:
				t.tkh.Cursor--
				t.state = html_tokenizer_attribute_value_unquoted_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(double-quoted)-state
		case html_tokenizer_attribute_value_double_quoted_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '"':
				t.state = html_tokenizer_after_attribute_value_quoted_state
			case '&':
				return_state = html_tokenizer_attribute_value_double_quoted_state
				t.state = html_tokenizer_character_reference_state
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_attr().value += string(rune(0xfffd))
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
				return
			default:
				curr_attr().value += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(single-quoted)-state
		case html_tokenizer_attribute_value_single_quoted_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\'':
				t.state = html_tokenizer_after_attribute_value_quoted_state
			case '&':
				return_state = html_tokenizer_attribute_value_single_quoted_state
				t.state = html_tokenizer_character_reference_state
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_attr().value += string(rune(0xfffd))
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
				return
			default:
				curr_attr().value += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#attribute-value-(unquoted)-state
		case html_tokenizer_attribute_value_unquoted_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				t.state = html_tokenizer_before_attribute_name_state
			case '&':
				return_state = html_tokenizer_attribute_value_unquoted_state
				t.state = html_tokenizer_character_reference_state
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_attr().value += string(rune(0xfffd))
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
				return
			case '"', '\'', '<', '=', '`':
				t.parse_error_encountered(html_unexpected_character_in_unquoted_attribute_value_error)
				fallthrough
			default:
				curr_attr().value += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#after-attribute-value-(quoted)-state
		case html_tokenizer_after_attribute_value_quoted_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				t.state = html_tokenizer_before_attribute_name_state
			case '/':
				t.state = html_tokenizer_self_closing_start_tag_state
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
				return
			default:
				t.parse_error_encountered(html_missing_whitepace_between_attributes_error)
				t.tkh.Cursor--
				t.state = html_tokenizer_before_attribute_name_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#self-closing-start-tag-state
		case html_tokenizer_self_closing_start_tag_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '>':
				curr_tag_token().is_self_closing = true
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case -1:
				t.parse_error_encountered(html_eof_in_tag_error)
				emit_token(&html_eof_token{})
				return
			default:
				t.parse_error_encountered(html_unexpected_solidus_in_tag_error)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#bogus-comment-state
		case html_tokenizer_bogus_comment_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case -1:
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_comment_token().data += string(rune(0xfffd))
			default:
				curr_comment_token().data += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#markup-declaration-open-state
		case html_tokenizer_markup_declaration_open_state:
			if t.tkh.ConsumeStrIfMatches("--", 0) != "" {
				curr_tok = &html_comment_token{data: ""}
				t.state = html_tokenizer_comment_start_state
			} else if t.tkh.ConsumeStrIfMatches("DOCTYPE", cm.MatchFlagsAsciiCaseInsensitive) != "" {
				t.state = html_tokenizer_doctype_state
			} else if t.tkh.ConsumeStrIfMatches("[CDATA[", 0) != "" {
				panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#markup-declaration-open-state]")
			} else {
				t.parse_error_encountered(html_incorrectly_opened_comment_error)
				curr_tok = &html_comment_token{data: ""}
				t.state = html_tokenizer_bogus_comment_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-start-state
		case html_tokenizer_comment_start_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '-':
				t.state = html_tokenizer_comment_start_dash_state
			case '>':
				t.parse_error_encountered(html_abrupt_closing_of_empty_comment_error)
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			default:
				t.tkh.Cursor--
				t.state = html_tokenizer_comment_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-start-dash-state
		case html_tokenizer_comment_start_dash_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '-':
				t.state = html_tokenizer_comment_end_state
			case '>':
				t.parse_error_encountered(html_abrupt_closing_of_empty_comment_error)
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case -1:
				t.parse_error_encountered(html_eof_in_comment_error)
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				curr_comment_token().data += "-"
				t.tkh.Cursor--
				t.state = html_tokenizer_comment_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-state
		case html_tokenizer_comment_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '<':
				curr_comment_token().data += string(next_char)
				t.state = html_tokenizer_comment_less_than_sign_state
			case '-':
				t.state = html_tokenizer_comment_end_dash_state
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				curr_comment_token().data += string(rune(0xfffd))
			case -1:
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				curr_comment_token().data += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-less-than-sign-state
		case html_tokenizer_comment_less_than_sign_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '!':
				curr_comment_token().data += string(next_char)
				panic("TODO")
				// t.state = html_tokenizer_comment_less_than_sign_bang_state
			case '<':
				curr_comment_token().data += string(next_char)
			default:
				t.tkh.Cursor--
				t.state = html_tokenizer_comment_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-end-dash-state
		case html_tokenizer_comment_end_dash_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '-':
				t.state = html_tokenizer_comment_end_state
			case -1:
				t.parse_error_encountered(html_eof_in_comment_error)
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				curr_comment_token().data += "-"
				t.tkh.Cursor--
				t.state = html_tokenizer_comment_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#comment-end-state
		case html_tokenizer_comment_end_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case '!':
				panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#comment-end-state]")
				// t.state = html_tokenizer_comment_end_bang_state
			case -1:
				t.parse_error_encountered(html_eof_in_comment_error)
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				curr_comment_token().data += "--"
				t.tkh.Cursor--
				t.state = html_tokenizer_comment_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#doctype-state
		case html_tokenizer_doctype_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				t.state = html_tokenizer_before_doctype_name_state
			case '>':
				t.tkh.Cursor--
				t.state = html_tokenizer_before_doctype_name_state
			case -1:
				t.parse_error_encountered(html_eof_in_doctype_error)
				curr_tok = &html_doctype_token{}
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				t.parse_error_encountered(html_missing_whitespace_before_doctype_name_error)
				t.tkh.Cursor--
				t.state = html_tokenizer_before_doctype_name_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#before-doctype-name-state
		case html_tokenizer_before_doctype_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				v := string(rune(0xfffd))
				curr_tok = &html_doctype_token{
					name: &v,
				}
			case '>':
				t.parse_error_encountered(html_missing_doctype_name_error)
				curr_tok = &html_doctype_token{
					force_quirks: true,
				}
				emit_token(curr_tok)
			case -1:
				t.parse_error_encountered(html_eof_in_doctype_error)
				curr_tok = &html_doctype_token{}
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				if cm.AsciiUppercaseRegex.MatchString(string(next_char)) {
					next_char = next_char - 'A' + 'a'
				}
				v := string(next_char)
				curr_tok = &html_doctype_token{
					name: &v,
				}
				t.state = html_tokenizer_doctype_name_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#doctype-name-state
		case html_tokenizer_doctype_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
				t.state = html_tokenizer_after_doctype_name_state
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case 0:
				t.parse_error_encountered(html_unexpected_null_character_error)
				*(curr_doctype_token().name) += string(rune(0xfffd))
			case -1:
				t.parse_error_encountered(html_eof_in_doctype_error)
				curr_doctype_token().force_quirks = true
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				if cm.AsciiUppercaseRegex.MatchString(string(next_char)) {
					next_char = next_char - 'A' + 'a'
				}
				*(curr_doctype_token().name) += string(next_char)
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state
		case html_tokenizer_after_doctype_name_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '\t', '\n', 0x000c, ' ':
			case '>':
				t.state = html_tokenizer_data_state
				emit_token(curr_tok)
			case -1:
				t.parse_error_encountered(html_eof_in_doctype_error)
				curr_doctype_token().force_quirks = true
				emit_token(curr_tok)
				emit_token(&html_eof_token{})
				return
			default:
				if t.tkh.ConsumeStrIfMatches("PUBLIC", cm.MatchFlagsAsciiCaseInsensitive) != "" {
					panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state]")
				} else if t.tkh.ConsumeStrIfMatches("SYSTEM", cm.MatchFlagsAsciiCaseInsensitive) != "" {
					panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state]")
				} else {
					t.parse_error_encountered(html_invalid_character_sequence_after_doctype_name_error)
					curr_doctype_token().force_quirks = true
					t.tkh.Cursor--
					panic("TODO[https://html.spec.whatwg.org/multipage/parsing.html#after-doctype-name-state]")
					// t.state = html_tokenizer_bogus_doctype_state
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#character-reference-state
		case html_tokenizer_character_reference_state:
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case '#':
				temp_buf += string(next_char)
				t.state = html_tokenizer_numeric_character_reference_state
			default:
				if cm.AsciiAlphanumericRegex.MatchString(string(next_char)) {
					t.tkh.Cursor--
					t.state = html_tokenizer_named_character_reference_state
				} else {
					flush_code_points_consumed_as_char_reference()
					t.tkh.Cursor--
					t.state = return_state
				}
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#named-character-reference-state
		case html_tokenizer_named_character_reference_state:
			found_name := ""
			for name := range html_entities {
				old_cursor := t.tkh.Cursor
				if name[0] != '&' {
					log.Printf("internal warning: key %s in html_entities doesn't start with &", name)
					continue
				}
				if t.tkh.ConsumeStrIfMatches(name[1:], 0) != "" {
					if len(found_name) < len(name) {
						found_name = name
					}
					t.tkh.Cursor = old_cursor
				}
			}
			if found_name != "" {
				t.tkh.Cursor += len([]rune(found_name))
				entity := html_entities[found_name]
				if is_consumed_as_part_of_attr() &&
					!strings.HasSuffix(found_name, ";") &&
					(t.tkh.PeekChar() == '=' || cm.AsciiAlphanumericRegex.MatchString(string(t.tkh.PeekChar()))) {
					flush_code_points_consumed_as_char_reference()
					t.state = return_state
				} else {
					if !strings.HasSuffix(found_name, ";") {
						t.parse_error_encountered(html_missing_semicolon_after_character_reference_error)
					}
					temp_buf = entity.characters
					flush_code_points_consumed_as_char_reference()
					t.state = return_state
				}
			} else {
				flush_code_points_consumed_as_char_reference()
				t.state = return_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-state
		case html_tokenizer_numeric_character_reference_state:
			character_reference_code = 0
			next_char := t.tkh.ConsumeChar()
			switch next_char {
			case 'X', 'x':
				temp_buf += string(next_char)
				t.state = html_tokenizer_hexadecimal_character_reference_start_state
			default:
				t.tkh.Cursor--
				t.state = html_tokenizer_decimal_character_reference_start_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#hexadecimal-character-reference-start-state
		case html_tokenizer_hexadecimal_character_reference_start_state:
			next_char := t.tkh.ConsumeChar()
			if cm.AsciiHexDigitRegex.MatchString(string(next_char)) {
				t.tkh.Cursor--
				t.state = html_tokenizer_hexadecimal_character_reference_state
			} else {
				t.parse_error_encountered(html_absence_of_digits_in_numeric_character_reference_error)
				t.tkh.Cursor--
				t.state = return_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-start-state
		case html_tokenizer_decimal_character_reference_start_state:
			next_char := t.tkh.ConsumeChar()
			if cm.AsciiDigitRegex.MatchString(string(next_char)) {
				t.tkh.Cursor--
				t.state = html_tokenizer_decimal_character_reference_state
			} else {
				t.parse_error_encountered(html_absence_of_digits_in_numeric_character_reference_error)
				t.tkh.Cursor--
				t.state = return_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#hexadecimal-character-reference-state
		case html_tokenizer_hexadecimal_character_reference_state:
			next_char := t.tkh.ConsumeChar()
			if cm.AsciiDigitRegex.MatchString(string(next_char)) {
				character_reference_code = (character_reference_code * 16) + int(next_char-'0')
			} else if cm.AsciiUpperHexDigitRegex.MatchString(string(next_char)) {
				character_reference_code = (character_reference_code * 16) + int(next_char-'A'+10)
			} else if cm.AsciiLowerHexDigitRegex.MatchString(string(next_char)) {
				character_reference_code = (character_reference_code * 16) + int(next_char-'a'+10)
			} else if next_char == ';' {
				t.state = html_tokenizer_numeric_character_reference_end_state
			} else {
				t.parse_error_encountered(html_missing_semicolon_after_character_reference_error)
				t.tkh.Cursor--
				t.state = html_tokenizer_numeric_character_reference_end_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#decimal-character-reference-state
		case html_tokenizer_decimal_character_reference_state:
			next_char := t.tkh.ConsumeChar()
			if cm.AsciiDigitRegex.MatchString(string(next_char)) {
				character_reference_code = (character_reference_code * 10) + int(next_char-'0')
			} else if next_char == ';' {
				t.state = html_tokenizer_numeric_character_reference_end_state
			} else {
				t.parse_error_encountered(html_missing_semicolon_after_character_reference_error)
				t.tkh.Cursor--
				t.state = html_tokenizer_numeric_character_reference_end_state
			}
		// https://html.spec.whatwg.org/multipage/parsing.html#numeric-character-reference-end-state
		case html_tokenizer_numeric_character_reference_end_state:
			chr := rune(character_reference_code)
			if chr == 0x0000 {
				t.parse_error_encountered(html_null_character_reference_error)
				chr = 0xfffd
			} else if 0x10ffff < chr {
				t.parse_error_encountered(html_character_reference_outside_unicode_range_error)
				chr = 0xfffd
			} else if cm.IsSurrogateChar(chr) {
				t.parse_error_encountered(html_surrogate_character_reference_error)
				chr = 0xfffd
			} else if cm.IsNoncharacter(chr) {
				t.parse_error_encountered(html_noncharacter_reference_error)
			} else if (chr == 0x0d) || (cm.IsControlChar(chr) && !cm.IsAsciiWhitespace(chr)) {
				t.parse_error_encountered(html_control_character_reference_error)
				replace_table := []struct {
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
				for _, item := range replace_table {
					if item.from == chr {
						chr = item.to
						break
					}
				}
			}
			temp_buf = ""
			temp_buf += string(chr)
			flush_code_points_consumed_as_char_reference()
			t.state = return_state
		default:
			fmt.Printf("Unhandled state %v", t.state)
			panic("unhandled state")
		}
	}

}
