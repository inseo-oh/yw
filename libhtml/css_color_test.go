package libhtml

import (
	"testing"
	cm "yw/libcommon"
)

func TestCssColor(t *testing.T) {
	cases := []struct {
		css      string
		expected css_color
	}{
		// Hex color -----------------------------------------------------------
		{"#89abcdef", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(0x89), css_number_from_int(0xab), css_number_from_int(0xcd), css_number_from_int(0xef),
		}}},
		// transparent ---------------------------------------------------------
		{"transparent", css_color_transparent()},
		// rgb()/rgba() (legacy syntax) ----------------------------------------
		{"rgb(10, 20, 30)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_int(255),
		}}},
		{"rgba(10, 20, 30)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_int(255),
		}}},
		{"rgb(50%, 25%, 100%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_float(127.5), css_number_from_float(63.75), css_number_from_int(255), css_number_from_int(255),
		}}},
		{"rgba(50%, 25%, 100%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_float(127.5), css_number_from_float(63.75), css_number_from_int(255), css_number_from_int(255),
		}}},
		{"rgb(10, 20, 30, 0.5)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		{"rgba(10, 20, 30, 0.5)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		{"rgb(10, 20, 30, 50%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		{"rgba(10, 20, 30, 50%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		// rgb()/rgba() (modern syntax) ----------------------------------------
		{"rgb(10 20 30)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_int(255),
		}}},
		{"rgba(10 20 30)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_int(255),
		}}},
		{"rgb(50% 25% 100%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_float(127.5), css_number_from_float(63.75), css_number_from_int(255), css_number_from_int(255),
		}}},
		{"rgba(50% 25% 100%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_float(127.5), css_number_from_float(63.75), css_number_from_int(255), css_number_from_int(255),
		}}},
		{"rgb(10 20 30 / 0.5)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		{"rgba(10 20 30 / 0.5)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		{"rgb(10 20 30 / 50%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
		{"rgba(10 20 30 / 50%)", css_color{css_color_type_rgb, []css_number{
			css_number_from_int(10), css_number_from_int(20), css_number_from_int(30), css_number_from_float(127.5),
		}}},
	}
	for _, cs := range cases {
		t.Run(cs.css, func(t *testing.T) {
			tokens, err := css_tokenize(cs.css)
			if tokens == nil && err != nil {
				t.Errorf("failed to tokenize: %v", err)
				return
			}
			t.Logf("Tokens: %v", tokens)
			got, err := css_parse(tokens, func(ts *css_token_stream) (*css_color, error) {
				return ts.parse_color()
			})
			if cm.IsNil(got) && err != nil {
				t.Errorf("failed to parse: %v", err)
				return
			}
			t.Logf("Parsed: %v", got)
			if cm.IsNil(got) || !got.equals(cs.expected) {
				t.Errorf("expected %v, got %v", cs.expected, got)
			}
		})

	}
}
