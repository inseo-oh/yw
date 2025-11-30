package libhtml

import (
	"reflect"
	"testing"

	cm "github.com/inseo-oh/yw/util"
)

func css_selector_test_helper(t *testing.T, css string, expected css_selector, parser func(ts *css_token_stream) (css_selector, error)) {
	t.Run(css, func(t *testing.T) {
		tokens, err := css_tokenize(css)
		if tokens == nil && err != nil {
			t.Errorf("failed to tokenize: %v", err)
			return
		}
		t.Logf("Tokens: %v", tokens)
		got, err := css_parse(tokens, parser)
		if cm.IsNil(got) && err != nil {
			t.Errorf("failed to parse: %v", err)
			return
		}
		t.Logf("Parsed: %v", got)
		t.Logf("RUN: %v(%v)", reflect.TypeOf(got), reflect.TypeOf(expected))
		if cm.IsNil(got) || !got.equals(expected) {
			t.Errorf("expected %v, got %v", expected, got)

		}
	})
}

func TestCssTypeSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_type_selector
	}{
		{"tt", css_type_selector{css_selector_wq_name{nil, "tt"}}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_type_selector(), nil
		})
	}
}
func TestCssIdSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_id_selector
	}{
		{"#tt", css_id_selector{"tt"}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_subclass_selector()
		})
	}
}
func TestCssClassSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_class_selector
	}{
		{".tt", css_class_selector{"tt"}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, (*css_token_stream).parse_subclass_selector)
	}
}
func TestCssPseudoClassSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_pseudo_class_selector
	}{
		{":link", css_pseudo_class_selector{"link", nil}},
		{":nth-child(1)", css_pseudo_class_selector{"nth-child", []css_token{
			css_number_token{css_token_common{10, 11}, css_number_from_int(1)},
		}}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_pseudo_class_selector()
		})
	}
}
func TestCssPseudoElementSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_pseudo_class_selector
	}{
		{"::before", css_pseudo_class_selector{"before", nil}},
		{"::after", css_pseudo_class_selector{"after", nil}},
		{"::first-line", css_pseudo_class_selector{"first-line", nil}},
		{"::first-letter", css_pseudo_class_selector{"first-letter", nil}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_pseudo_element_selector()
		})
	}
}
func TestCssAttributeSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_attribute_selector
	}{
		{"[tt]", css_attribute_selector{css_selector_wq_name{nil, "tt"}, css_attribute_matcher_none, "", true}},
		{"[attr=value]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_normal, "value", true}},
		{"[attr=value s]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_normal, "value", true}},
		{"[attr=value i]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_normal, "value", false}},
		{"[attr=\"quoted_value\"]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_normal, "quoted_value", true}},
		{"[attr~=value]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_tilde, "value", true}},
		{"[attr|=value]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_bar, "value", true}},
		{"[attr^=value]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_caret, "value", true}},
		{"[attr$=value]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_dollar, "value", true}},
		{"[attr*=value]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_asterisk, "value", true}},
		{"[  attr  *=  'value'  ]", css_attribute_selector{css_selector_wq_name{nil, "attr"}, css_attribute_matcher_asterisk, "value", true}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_attribute_selector()
		})
	}
}

func TestCssCompoundSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_compound_selector
	}{
		{"type#id.class", css_compound_selector{
			css_type_selector{css_selector_wq_name{nil, "type"}},
			[]css_selector{
				css_id_selector{"id"},
				css_class_selector{"class"},
			},
			nil,
		}},
		{"type::before", css_compound_selector{
			css_type_selector{css_selector_wq_name{nil, "type"}},
			nil,
			[]css_compound_selector_pseudo_item{
				{
					css_pseudo_class_selector{"before", nil},
					[]css_pseudo_class_selector{},
				},
			},
		}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_compound_selector()
		})
	}
}

func TestCssComplexSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected css_complex_selector
	}{
		{"type", css_complex_selector{
			css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type"}}, nil, nil},
			nil,
		}},
		{"type1>type2+type3~type4||type5  type6", css_complex_selector{
			css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type1"}}, nil, nil},
			[]css_complex_selector_rest{
				{css_combinator_direct_child, css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type2"}}, nil, nil}},
				{css_combinator_plus, css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type3"}}, nil, nil}},
				{css_combinator_tilde, css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type4"}}, nil, nil}},
				{css_combinator_two_bars, css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type5"}}, nil, nil}},
				{css_combinator_child, css_compound_selector{css_type_selector{css_selector_wq_name{nil, "type6"}}, nil, nil}},
			},
		}},
	}
	for _, cs := range cases {
		css_selector_test_helper(t, cs.css, cs.expected, func(ts *css_token_stream) (css_selector, error) {
			return ts.parse_complex_selector()
		})
	}
}
