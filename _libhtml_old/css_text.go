// Implementation of the CSS Text Module Level 3 (https://www.w3.org/TR/css-text-3)
package libhtml

import (
	"fmt"
	"log"
	"strings"

	cm "github.com/inseo-oh/yw/util"
)

type css_text_transform struct {
	tp             css_text_transform_type
	full_width     bool
	full_size_kana bool
}

type css_text_transform_type uint8

const (
	css_text_transform_type_none          = css_text_transform_type(iota) // No transform
	css_text_transform_type_original_caps                                 // Don't change cases
	css_text_transform_type_capitalize                                    // Captialize text
	css_text_transform_type_uppercase                                     // MAKE TEXT UPPERCASE
	css_text_transform_type_lowercase                                     // make text lowercase
)

// https://www.w3.org/TR/css-text-3/#propdef-text-transform
func (ts *css_token_stream) parse_text_transform() (css_text_transform, bool) {
	if !cm.IsNil(ts.consume_ident_token_with("none")) {
		return css_text_transform{css_text_transform_type_none, false, false}, true
	}
	out := css_text_transform{tp: css_text_transform_type_original_caps}
	got_type := false
	got_is_full_width := false
	got_full_kana := false
	got_any := false
	for {
		valid := false
		if !got_type {
			ts.skip_whitespaces()
			if !cm.IsNil(ts.consume_ident_token_with("capitalize")) {
				out.tp = css_text_transform_type_capitalize
				got_type = true
				valid = true
			} else if !cm.IsNil(ts.consume_ident_token_with("uppercase")) {
				out.tp = css_text_transform_type_uppercase
				got_type = true
				valid = true
			} else if !cm.IsNil(ts.consume_ident_token_with("lowercase")) {
				out.tp = css_text_transform_type_lowercase
				got_type = true
				valid = true
			}
		}
		if !got_is_full_width {
			ts.skip_whitespaces()
			if !cm.IsNil(ts.consume_ident_token_with("full-width")) {
				out.full_width = true
				got_is_full_width = true
				valid = true
			}
		}
		if !got_full_kana {
			ts.skip_whitespaces()
			if !cm.IsNil(ts.consume_ident_token_with("full-size-kana")) {
				out.full_width = true
				got_is_full_width = true
				valid = true
			}
		}
		ts.skip_whitespaces()
		if !valid {
			break
		}
		got_any = true
	}
	if !got_any {
		return out, false
	}
	return out, true
}

func (s css_text_transform) String() string {
	sb := strings.Builder{}
	switch s.tp {
	case css_text_transform_type_none:
		return "none"
	case css_text_transform_type_original_caps:
		break
	case css_text_transform_type_capitalize:
		sb.WriteString("capitalize ")
	case css_text_transform_type_lowercase:
		sb.WriteString("lowercase ")
	case css_text_transform_type_uppercase:
		sb.WriteString("uppercase ")
	default:
		return fmt.Sprintf("unregognized css_text_transform type %v", s.tp)
	}
	if s.full_width {
		sb.WriteString("full-width ")
	}
	if s.full_size_kana {
		sb.WriteString("full-size-kana ")
	}
	return strings.TrimSpace(sb.String())
}
func (s css_text_transform) apply(text string) string {
	out_text := text
	switch s.tp {
	case css_text_transform_type_none:
		return text
	case css_text_transform_type_original_caps:
		break
	case css_text_transform_type_capitalize:
		runes := []rune(text)
		out_text = strings.ToUpper(string(runes[0])) + string(runes[1:])
	case css_text_transform_type_lowercase:
		out_text = strings.ToLower(text)
	case css_text_transform_type_uppercase:
		out_text = strings.ToUpper(text)
	default:
		log.Panicf("unregognized css_text_transform type %v", s.tp)
	}
	if s.full_width {
		log.Printf("TODO: %v: Support full-width", s)
	}
	if s.full_size_kana {
		log.Printf("TODO: %v: Support full-size-kana", s)
	}
	return out_text
}
