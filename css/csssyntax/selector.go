package csssyntax

import (
	"errors"
	"yw/css/selector"
	"yw/util"
)

// Returns nil if not found
func (ts *tokenStream) parseSelectorNsPrefix() *selector.NsPrefix {
	// STUB
	return nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-wq-name
//
// Returns nil if not found
func (ts *tokenStream) parseSelectorWqName() *selector.WqName {
	oldCursor := ts.cursor
	nsPrefix := ts.parseSelectorNsPrefix()
	var identTk identToken
	if temp := ts.consumeTokenWith(tokenTypeIdent); !util.IsNil(temp) {
		identTk = temp.(identToken)
	} else {
		ts.cursor = oldCursor
		return nil
	}
	return &selector.WqName{NsPrefix: nsPrefix, Ident: identTk.value}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#ref-for-typedef-type-selector
//
// Returns nil if not found
func (ts *tokenStream) parseTypeSelector() selector.Selector {
	oldCursor := ts.cursor
	if typeName := ts.parseSelectorWqName(); typeName != nil {
		// <wq-name>
		return selector.TypeSelector{TypeName: *typeName}
	} else {
		// <ns-prefix?> *
		nsPrefix := ts.parseSelectorNsPrefix()
		if tk := ts.consumeTokenWith(tokenTypeDelim); !util.IsNil(tk) {
			if tk.(delimToken).value != '*' {
				ts.cursor = oldCursor
				return nil
			}
			return selector.WildcardSelector{NsPrefix: nsPrefix}
		} else {
			ts.cursor = oldCursor
			return nil
		}
	}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-id-selector
func (ts *tokenStream) parseIdSelector() *selector.IdSelector {
	oldCursor := ts.cursor
	var hashTk hashToken
	if temp := ts.consumeTokenWith(tokenTypeHash); !util.IsNil(temp) {
		hashTk = temp.(hashToken)
	} else {
		ts.cursor = oldCursor
		return nil
	}
	return &selector.IdSelector{Id: hashTk.value}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-class-selector
func (ts *tokenStream) parseClassSelector() (*selector.ClassSelector, error) {
	oldCursor := ts.cursor
	if util.IsNil(ts.consumeDelimTokenWith('.')) {
		ts.cursor = oldCursor
		return nil, nil
	}
	var identTk identToken
	if temp := ts.consumeTokenWith(tokenTypeIdent); !util.IsNil(temp) {
		identTk = temp.(identToken)
	} else {
		return nil, errors.New("expected identifier after '.'")
	}
	return &selector.ClassSelector{Class: identTk.value}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attribute-selector
func (ts *tokenStream) parseAttrSelector() (*selector.AttrSelector, error) {
	oldCursor := ts.cursor
	blk := ts.consumeSimpleBlockWith(simpleBlockTypeSquare)
	if blk == nil {
		ts.cursor = oldCursor
		return nil, nil
	}

	bodyStream := tokenStream{tokens: blk.body}
	// [<  >attr  ] ------------------------------------------------------------
	// [<  >attr  =  value  modifier  ] ----------------------------------------
	bodyStream.skipWhitespaces()
	// [  <attr>  ] ------------------------------------------------------------
	// [  <attr>  =  value  modifier  ] ----------------------------------------
	wqName := bodyStream.parseSelectorWqName()
	if wqName == nil {
		return nil, errors.New("expected name after '['")
	}
	// [  attr<  >] ------------------------------------------------------------
	// [  attr<  >=  value  modifier  ] ----------------------------------------
	bodyStream.skipWhitespaces()
	if !bodyStream.isEnd() {
		// [  attr  <=>  value  modifier  ] ------------------------------------
		var matcher selector.Matcher
		// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attr-matcher
		if !util.IsNil(bodyStream.consumeDelimTokenWith('~')) {
			matcher = selector.TildeMatcher
		} else if !util.IsNil(bodyStream.consumeDelimTokenWith('|')) {
			matcher = selector.BarMatcher
		} else if !util.IsNil(bodyStream.consumeDelimTokenWith('^')) {
			matcher = selector.CaretMatcher
		} else if !util.IsNil(bodyStream.consumeDelimTokenWith('$')) {
			matcher = selector.DollarMatcher
		} else if !util.IsNil(bodyStream.consumeDelimTokenWith('*')) {
			matcher = selector.AsteriskMatcher
		} else {
			matcher = selector.NormalMatcher
		}
		if util.IsNil(bodyStream.consumeDelimTokenWith('=')) {
			return nil, errors.New("expected operator after the attribute name")
		}
		// [  attr  =<  >value  modifier  ] ------------------------------------
		bodyStream.skipWhitespaces()
		// [  attr  =  <value>  modifier  ] ------------------------------------
		var attrValue string
		if n := bodyStream.consumeTokenWith(tokenTypeIdent); !util.IsNil(n) {
			attrValue = n.(identToken).value
		} else if n := bodyStream.consumeTokenWith(tokenTypeString); !util.IsNil(n) {
			attrValue = n.(stringToken).value
		} else {
			return nil, errors.New("expected attribute value after the operator")
		}
		// [  attr  =  value<  >modifier  ] ------------------------------------
		bodyStream.skipWhitespaces()
		// [  attr  =  value  <modifier>  ] ------------------------------------
		isCaseSensitive := true
		if !util.IsNil(bodyStream.consumeIdentTokenWith("i")) {
			isCaseSensitive = false
		} else if !util.IsNil(bodyStream.consumeIdentTokenWith("s")) {
			isCaseSensitive = true
		}
		// [  attr  =  value  modifier<  >] ------------------------------------
		bodyStream.skipWhitespaces()
		if !bodyStream.isEnd() {
			return nil, errors.New("found junk after contents of the attribute selector")
		}
		return &selector.AttrSelector{AttrName: *wqName, Matcher: matcher, AttrValue: attrValue, IsCaseSensitive: isCaseSensitive}, nil
	}
	return &selector.AttrSelector{AttrName: *wqName, Matcher: selector.NoMatcher, AttrValue: "", IsCaseSensitive: false}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-class-selector
func (ts *tokenStream) parsePseudoClassSelector() (*selector.PseudoClassSelector, error) {
	oldCursor := ts.cursor

	// <:>name ----------------------------------------------------------------
	// <:>func(value) ----------------------------------------------------------
	if util.IsNil(ts.consumeTokenWith(tokenTypeColon)) {
		ts.cursor = oldCursor
		return nil, nil
	}
	if identTk := ts.consumeTokenWith(tokenTypeIdent); !util.IsNil(identTk) {
		// :<name> ------------------------------------------------------------
		name := identTk.(identToken).value
		return &selector.PseudoClassSelector{Name: name, Args: nil}, nil
	} else if funcNode := ts.consumeTokenWith(tokenTypeAstFunc); !util.IsNil(funcNode) {
		// :<func(value)> ------------------------------------------------------
		name := funcNode.(astFuncToken).name
		subStream := tokenStream{tokens: funcNode.(astFuncToken).value}
		args := subStream.consumeAnyValue()
		if args == nil {
			ts.cursor = oldCursor
			return nil, errors.New("expected value after '('")
		}
		if !subStream.isEnd() {
			return nil, errors.New("unexpected junk after arguments")
		}
		argsNew := []any{}
		for _, arg := range args {
			argsNew = append(argsNew, arg)
		}
		return &selector.PseudoClassSelector{Name: name, Args: argsNew}, nil
	} else {
		ts.cursor = oldCursor
		return nil, nil
	}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-element-selector
func (ts *tokenStream) parsePseudoElementSelector() (*selector.PseudoClassSelector, error) {
	oldCursor := ts.cursor
	if util.IsNil(ts.consumeTokenWith(tokenTypeColon)) {
		ts.cursor = oldCursor
		return nil, nil
	}
	if temp, err := ts.parsePseudoClassSelector(); temp != nil {
		return temp, nil
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-subclass-selector
//
// Returns nil if not found
func (ts *tokenStream) parseSubclassSelector() (selector.Selector, error) {
	if sel := ts.parseIdSelector(); sel != nil {
		return *sel, nil
	}

	if sel, err := ts.parseClassSelector(); sel != nil {
		return *sel, nil
	} else if err != nil {
		return nil, err
	}

	if sel, err := ts.parseAttrSelector(); sel != nil {
		return *sel, nil
	} else if err != nil {
		return nil, err
	}

	if sel, err := ts.parsePseudoClassSelector(); sel != nil {
		return *sel, nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-compound-selector
//
// Returns nil if not found
func (ts *tokenStream) parseCompoundSelector() (*selector.CompoundSelector, error) {
	oldCursor := ts.cursor
	typeSel := ts.parseTypeSelector()
	subclassSels := []selector.Selector{}
	pseudoItems := []selector.CompundSelectorPseudoItem{}
	for {
		subclassSel, err := ts.parseSubclassSelector()
		if util.IsNil(subclassSel) {
			if err != nil {
				return nil, err
			}
			break
		}
		subclassSels = append(subclassSels, subclassSel)
	}

	for {
		elemSel, err := ts.parsePseudoElementSelector()
		if elemSel == nil {
			if err != nil {
				return nil, err
			}
			break
		}
		classSels := []selector.PseudoClassSelector{}
		for {
			classSel, err := ts.parsePseudoClassSelector()
			if util.IsNil(classSel) {
				if err != nil {
					ts.cursor = oldCursor
					return nil, err
				}
				break
			}
			classSels = append(classSels, *classSel)
		}
		pseudoItems = append(pseudoItems, selector.CompundSelectorPseudoItem{ElementSelector: *elemSel, ClassSelector: classSels})

	}

	if typeSel == nil && len(subclassSels) == 0 && len(pseudoItems) == 0 {
		ts.cursor = oldCursor
		return nil, nil
	}
	return &selector.CompoundSelector{TypeSelector: typeSel, SubclassSelector: subclassSels, PseudoItems: pseudoItems}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-compound-selector-list
//
// Returns nil if not found
func (ts *tokenStream) parseCompoundSelectorList() ([]*selector.CompoundSelector, error) {
	return parseCommaSeparatedRepeation(ts, 0, func(ts *tokenStream) (*selector.CompoundSelector, error) {
		return ts.parseCompoundSelector()
	})
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector

func (ts *tokenStream) parseComplexSelector() (*selector.ComplexSelector, error) {
	oldCursor := ts.cursor
	base, err := ts.parseCompoundSelector()
	if base == nil {
		ts.cursor = oldCursor
		return nil, err
	}
	rest := []selector.ComplexSelectorRest{}
	for {
		comb := selector.ChildCombinator
		if !util.IsNil(ts.consumeDelimTokenWith('>')) {
			comb = selector.DirectChildCombinator
		} else if !util.IsNil(ts.consumeDelimTokenWith('+')) {
			comb = selector.PlusCombinator
		} else if !util.IsNil(ts.consumeDelimTokenWith('~')) {
			comb = selector.TildeCombinator
		} else if !util.IsNil(ts.consumeDelimTokenWith('|')) {
			if !util.IsNil(ts.consumeDelimTokenWith('|')) {
				comb = selector.TwoBarsCombinator
			} else {
				ts.cursor -= 2
			}
		} else if !util.IsNil(ts.consumeTokenWith(tokenTypeWhitespace)) {
			ts.skipWhitespaces()
			comb = selector.ChildCombinator
		}
		anotherUnit, err := ts.parseCompoundSelector()
		if util.IsNil(anotherUnit) {
			if err != nil {
				return nil, err
			}
			break
		}
		rest = append(rest, selector.ComplexSelectorRest{Combinator: comb, Selector: *anotherUnit})
	}
	if base == nil && len(rest) == 0 {
		ts.cursor = oldCursor
		return nil, nil
	}
	return &selector.ComplexSelector{Base: *base, Rest: rest}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector-list
//
// Returns nil if not found
func (ts *tokenStream) parseComplexSelectorList() ([]selector.Selector, error) {
	selList, err := parseCommaSeparatedRepeation(ts, 0, func(ts *tokenStream) (*selector.ComplexSelector, error) {
		return ts.parseComplexSelector()
	})
	if selList == nil {
		return nil, err
	}
	out := []selector.Selector{}
	for _, s := range selList {
		out = append(out, *s)
	}
	return out, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-selector-list
func (ts *tokenStream) parseSelectorList() ([]selector.Selector, error) {
	return ts.parseComplexSelectorList()
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#parse-a-selector
func parseSelector(src string) ([]selector.Selector, error) {
	tokens, err := Tokenize(src)
	if tokens == nil && err != nil {
		return nil, err
	}
	return parse(tokens, func(ts *tokenStream) ([]selector.Selector, error) {
		return ts.parseSelectorList()
	})
}
