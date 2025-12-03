package csssyntax

import (
	"testing"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/selector"
	"github.com/inseo-oh/yw/util"
)

func selectorTestHelper(t *testing.T, css string, expected selector.Selector, parser func(ts *tokenStream) (selector.Selector, error)) {
	t.Run(css, func(t *testing.T) {
		tokens, err := tokenize([]byte(css))
		if tokens == nil && err != nil {
			t.Errorf("failed to tokenize: %v", err)
			return
		}
		t.Logf("Tokens: %v", tokens)
		got, err := parse(tokens, parser)
		if util.IsNil(got) && err != nil {
			t.Errorf("failed to parse: %v", err)
			return
		}
		t.Logf("Parsed: %v", got)
		if util.IsNil(got) || !got.Equals(expected) {
			t.Errorf("expected %v, got %v", expected, got)

		}
	})
}

func TestCssTypeSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.TypeSelector
	}{
		{"tt", selector.TypeSelector{TypeName: selector.WqName{Ident: "tt"}}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parseTypeSelector(), nil
		})
	}
}
func TestCssIdSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.IdSelector
	}{
		{"#tt", selector.IdSelector{Id: "tt"}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parseSubclassSelector()
		})
	}
}
func TestCssClassSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.ClassSelector
	}{
		{".tt", selector.ClassSelector{Class: "tt"}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, (*tokenStream).parseSubclassSelector)
	}
}
func TestCssPseudoClassSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.PseudoClassSelector
	}{
		{":link", selector.PseudoClassSelector{Name: "link"}},
		{":nth-child(1)", selector.PseudoClassSelector{Name: "nth-child", Args: []any{
			numberToken{tokenCommon{10, 11}, css.NumFromInt(1)},
		}}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parsePseudoClassSelector()
		})
	}
}
func TestCssPseudoElementSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.PseudoClassSelector
	}{
		{"::before", selector.PseudoClassSelector{Name: "before"}},
		{"::after", selector.PseudoClassSelector{Name: "after"}},
		{"::first-line", selector.PseudoClassSelector{Name: "first-line"}},
		{"::first-letter", selector.PseudoClassSelector{Name: "first-letter"}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parsePseudoElementSelector()
		})
	}
}
func TestCssAttributeSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.AttrSelector
	}{
		{"[tt]", selector.AttrSelector{AttrName: selector.WqName{Ident: "tt"}, Matcher: selector.NoMatcher, AttrValue: "", IsCaseSensitive: true}},
		{"[attr=value]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.NormalMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[attr=value s]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.NormalMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[attr=value i]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.NormalMatcher, AttrValue: "value", IsCaseSensitive: false}},
		{"[attr=\"quoted_value\"]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.NormalMatcher, AttrValue: "quoted_value", IsCaseSensitive: true}},
		{"[attr~=value]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.TildeMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[attr|=value]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.BarMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[attr^=value]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.CaretMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[attr$=value]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.DollarMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[attr*=value]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.AsteriskMatcher, AttrValue: "value", IsCaseSensitive: true}},
		{"[  attr  *=  'value'  ]", selector.AttrSelector{AttrName: selector.WqName{Ident: "attr"}, Matcher: selector.AsteriskMatcher, AttrValue: "value", IsCaseSensitive: true}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parseAttrSelector()
		})
	}
}

func TestCssCompoundSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.CompoundSelector
	}{
		{"type#id.class", selector.CompoundSelector{
			TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type"}},
			SubclassSelector: []selector.Selector{
				selector.IdSelector{Id: "id"},
				selector.ClassSelector{Class: "class"},
			},
		}},
		{"type::before", selector.CompoundSelector{
			TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type"}},

			PseudoItems: []selector.CompundSelectorPseudoItem{
				{
					ElementSelector: selector.PseudoClassSelector{Name: "before"},
					ClassSelector:   []selector.PseudoClassSelector{},
				},
			},
		}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parseCompoundSelector()
		})
	}
}

func TestCssComplexSelector(t *testing.T) {
	cases := []struct {
		css      string
		expected selector.ComplexSelector
	}{
		{"type", selector.ComplexSelector{
			Base: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type"}}},
		}},
		{"type1>type2+type3~type4||type5  type6", selector.ComplexSelector{
			Base: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type1"}}},
			Rest: []selector.ComplexSelectorRest{
				{Combinator: selector.DirectChildCombinator, Selector: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type2"}}}},
				{Combinator: selector.PlusCombinator, Selector: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type3"}}}},
				{Combinator: selector.TildeCombinator, Selector: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type4"}}}},
				{Combinator: selector.TwoBarsCombinator, Selector: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type5"}}}},
				{Combinator: selector.ChildCombinator, Selector: selector.CompoundSelector{TypeSelector: selector.TypeSelector{TypeName: selector.WqName{Ident: "type6"}}}},
			},
		}},
	}
	for _, cs := range cases {
		selectorTestHelper(t, cs.css, cs.expected, func(ts *tokenStream) (selector.Selector, error) {
			return ts.parseComplexSelector()
		})
	}
}
