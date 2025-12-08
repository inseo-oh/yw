// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/selector"
)

// Returns nil if not found
func (ts *tokenStream) parseSelectorNsPrefix() (res selector.NsPrefix, err error) {
	// STUB
	return res, fmt.Errorf("%s: expected namespace prefix", ts.errorHeader())
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-wq-name
//
// Returns nil if not found
func (ts *tokenStream) parseSelectorWqName() (res selector.WqName, err error) {
	oldCursor := ts.cursor

	var nsPrefix *selector.NsPrefix
	if v, err := ts.parseSelectorNsPrefix(); err == nil {
		nsPrefix = &v
	}
	var identTk identToken
	if temp, err := ts.consumeTokenWith(tokenTypeIdent); err == nil {
		identTk = temp.(identToken)
	} else {
		ts.cursor = oldCursor
		return res, fmt.Errorf("%s: expected namespace prefix or identifier", ts.errorHeader())
	}
	return selector.WqName{NsPrefix: nsPrefix, Ident: identTk.value}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#ref-for-typedef-type-selector
//
// Returns nil if not found
func (ts *tokenStream) parseTypeSelector() (res selector.Selector, err error) {
	oldCursor := ts.cursor
	if typeName, err := ts.parseSelectorWqName(); err == nil {
		// <wq-name>
		return selector.TypeSelector{TypeName: typeName}, nil
	} else {
		// <ns-prefix?> *
		var nsPrefix *selector.NsPrefix
		if v, err := ts.parseSelectorNsPrefix(); err == nil {
			nsPrefix = &v
		}
		if tk, err := ts.consumeTokenWith(tokenTypeDelim); err == nil {
			if tk.(delimToken).value != '*' {
				ts.cursor = oldCursor
				return res, fmt.Errorf("%s: expected *", ts.errorHeader())
			}
			return selector.UniversalSelector{NsPrefix: nsPrefix}, nil
		} else {
			return res, err
		}
	}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-id-selector
func (ts *tokenStream) parseIdSelector() (res selector.IdSelector, err error) {
	var hashTk hashToken
	if temp, err := ts.consumeTokenWith(tokenTypeHash); err == nil {
		hashTk = temp.(hashToken)
	} else {
		return res, err
	}
	return selector.IdSelector{Id: hashTk.value}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-class-selector
func (ts *tokenStream) parseClassSelector() (res selector.ClassSelector, err error) {
	oldCursor := ts.cursor
	if err := ts.consumeDelimTokenWith('.'); err != nil {
		ts.cursor = oldCursor
		return res, err
	}
	var identTk identToken
	if temp, err := ts.consumeTokenWith(tokenTypeIdent); err == nil {
		identTk = temp.(identToken)
	} else {
		return res, err
	}
	return selector.ClassSelector{Class: identTk.value}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attribute-selector
func (ts *tokenStream) parseAttrSelector() (res selector.AttrSelector, err error) {
	oldCursor := ts.cursor
	blk, err := ts.consumeSimpleBlockWith(simpleBlockTypeSquare)
	if err != nil {
		return res, err
	}

	bodyStream := tokenStream{tokens: blk.body, tokenizerHelper: ts.tokenizerHelper}
	// [<  >attr  ] ------------------------------------------------------------
	// [<  >attr  =  value  modifier  ] ----------------------------------------
	bodyStream.skipWhitespaces()
	// [  <attr>  ] ------------------------------------------------------------
	// [  <attr>  =  value  modifier  ] ----------------------------------------
	wqName, err := bodyStream.parseSelectorWqName()
	if err != nil {
		return res, err
	}
	// [  attr<  >] ------------------------------------------------------------
	// [  attr<  >=  value  modifier  ] ----------------------------------------
	bodyStream.skipWhitespaces()
	if !bodyStream.isEnd() {
		// [  attr  <=>  value  modifier  ] ------------------------------------
		var matcher selector.Matcher
		// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-attr-matcher
		if err := bodyStream.consumeDelimTokenWith('~'); err == nil {
			matcher = selector.TildeMatcher
		} else if err := bodyStream.consumeDelimTokenWith('|'); err == nil {
			matcher = selector.BarMatcher
		} else if err := bodyStream.consumeDelimTokenWith('^'); err == nil {
			matcher = selector.CaretMatcher
		} else if err := bodyStream.consumeDelimTokenWith('$'); err == nil {
			matcher = selector.DollarMatcher
		} else if err := bodyStream.consumeDelimTokenWith('*'); err == nil {
			matcher = selector.AsteriskMatcher
		} else {
			matcher = selector.NormalMatcher
		}
		if err := bodyStream.consumeDelimTokenWith('='); err != nil {
			return res, err
		}
		// [  attr  =<  >value  modifier  ] ------------------------------------
		bodyStream.skipWhitespaces()
		// [  attr  =  <value>  modifier  ] ------------------------------------
		var attrValue string
		if n, err := bodyStream.consumeTokenWith(tokenTypeIdent); err == nil {
			attrValue = n.(identToken).value
		} else if n, err := bodyStream.consumeTokenWith(tokenTypeString); err == nil {
			attrValue = n.(stringToken).value
		} else {
			ts.cursor = oldCursor
			return res, fmt.Errorf("%s: expected attribute value after the operator", ts.errorHeader())
		}
		// [  attr  =  value<  >modifier  ] ------------------------------------
		bodyStream.skipWhitespaces()
		// [  attr  =  value  <modifier>  ] ------------------------------------
		isCaseSensitive := true
		if err := bodyStream.consumeIdentTokenWith("i"); err == nil {
			isCaseSensitive = false
		} else if err := bodyStream.consumeIdentTokenWith("s"); err == nil {
			isCaseSensitive = true
		}
		// [  attr  =  value  modifier<  >] ------------------------------------
		bodyStream.skipWhitespaces()
		if !bodyStream.isEnd() {
			return res, fmt.Errorf("%s: found junk after contents of the attribute selector", ts.errorHeader())
		}
		return selector.AttrSelector{AttrName: wqName, Matcher: matcher, AttrValue: attrValue, IsCaseSensitive: isCaseSensitive}, nil
	}
	return selector.AttrSelector{AttrName: wqName, Matcher: selector.NoMatcher, AttrValue: "", IsCaseSensitive: false}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-class-selector
func (ts *tokenStream) parsePseudoClassSelector() (res selector.PseudoClassSelector, err error) {
	oldCursor := ts.cursor

	// <:>name ----------------------------------------------------------------
	// <:>func(value) ----------------------------------------------------------
	if _, err := ts.consumeTokenWith(tokenTypeColon); err != nil {
		return res, err
	} else if identTk, err := ts.consumeTokenWith(tokenTypeIdent); err == nil {
		// :<name> ------------------------------------------------------------
		name := identTk.(identToken).value
		return selector.PseudoClassSelector{Name: name, Args: nil}, nil
	}
	if funcNode, err := ts.consumeTokenWith(tokenTypeAstFunc); err == nil {
		// :<func(value)> ------------------------------------------------------
		name := funcNode.(astFuncToken).name
		subStream := tokenStream{tokens: funcNode.(astFuncToken).value, tokenizerHelper: ts.tokenizerHelper}
		args := subStream.consumeAnyValue()
		if args == nil {
			ts.cursor = oldCursor
			return res, fmt.Errorf("%s: expected value after '('", ts.errorHeader())
		}
		if !subStream.isEnd() {
			return res, fmt.Errorf("%s: unexpected junk after arguments", ts.errorHeader())
		}
		argsNew := []any{}
		for _, arg := range args {
			argsNew = append(argsNew, arg)
		}
		return selector.PseudoClassSelector{Name: name, Args: argsNew}, nil
	} else {
		ts.cursor = oldCursor
		return res, err
	}
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-pseudo-element-selector
func (ts *tokenStream) parsePseudoElementSelector() (res selector.PseudoClassSelector, err error) {
	oldCursor := ts.cursor
	if _, err := ts.consumeTokenWith(tokenTypeColon); err != nil {
		ts.cursor = oldCursor
		return res, err
	}
	temp, err := ts.parsePseudoClassSelector()
	if err != nil {
		return res, err
	}
	return temp, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-subclass-selector
//
// Returns nil if not found
func (ts *tokenStream) parseSubclassSelector() (res selector.Selector, err error) {
	if sel, err := ts.parseIdSelector(); err == nil {
		return sel, nil
	}

	if sel, err := ts.parseClassSelector(); err == nil {
		return sel, nil
	}

	if sel, err := ts.parseAttrSelector(); err == nil {
		return sel, nil
	}

	if sel, err := ts.parsePseudoClassSelector(); err == nil {
		return sel, nil
	}

	return nil, fmt.Errorf("%s: expected a subclass selector", ts.errorHeader())
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-compound-selector
//
// Returns nil if not found
func (ts *tokenStream) parseCompoundSelector() (res selector.CompoundSelector, err error) {
	oldCursor := ts.cursor

	var typeSel selector.Selector
	gotTypeSel := false
	if v, err := ts.parseTypeSelector(); err == nil {
		typeSel = v
		gotTypeSel = true
	}
	subclassSels := []selector.Selector{}
	pseudoItems := []selector.CompundSelectorPseudoItem{}
	for {
		subclassSel, err := ts.parseSubclassSelector()
		if err != nil {
			break
		}
		subclassSels = append(subclassSels, subclassSel)
	}

	for {
		elemSel, err := ts.parsePseudoElementSelector()
		if err != nil {
			break
		}
		classSels := []selector.PseudoClassSelector{}
		for {
			classSel, err := ts.parsePseudoClassSelector()
			if err != nil {
				break
			}
			classSels = append(classSels, classSel)
		}
		pseudoItems = append(pseudoItems, selector.CompundSelectorPseudoItem{ElementSelector: elemSel, ClassSelector: classSels})

	}

	if !gotTypeSel && len(subclassSels) == 0 && len(pseudoItems) == 0 {
		ts.cursor = oldCursor
		return selector.CompoundSelector{}, fmt.Errorf("%s: expected a compound selector", ts.errorHeader())
	}
	return selector.CompoundSelector{TypeSelector: typeSel, SubclassSelector: subclassSels, PseudoItems: pseudoItems}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-compound-selector-list
//
// Returns nil if not found
func (ts *tokenStream) parseCompoundSelectorList() (res []selector.CompoundSelector, err error) {
	return parseCommaSeparatedRepeation(ts, 0, "compound selector", func(ts *tokenStream) (selector.CompoundSelector, error) {
		return ts.parseCompoundSelector()
	})
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector

func (ts *tokenStream) parseComplexSelector() (res selector.ComplexSelector, err error) {
	oldCursor := ts.cursor
	defer func() {
		if err != nil {
			ts.cursor = oldCursor
		}
	}()

	base, err := ts.parseCompoundSelector()
	if err != nil {
		return selector.ComplexSelector{}, err
	}
	rest := []selector.ComplexSelectorRest{}
	for {
		comb := selector.ChildCombinator
		if err := ts.consumeDelimTokenWith('>'); err == nil {
			comb = selector.DirectChildCombinator
		} else if err := ts.consumeDelimTokenWith('+'); err == nil {
			comb = selector.PlusCombinator
		} else if err := ts.consumeDelimTokenWith('~'); err == nil {
			comb = selector.TildeCombinator
		} else if err := ts.consumeDelimTokenWith('|'); err == nil {
			if err := ts.consumeDelimTokenWith('|'); err == nil {
				comb = selector.TwoBarsCombinator
			} else {
				ts.cursor -= 2
			}
		} else if _, err := ts.consumeTokenWith(tokenTypeWhitespace); err != nil {
			ts.skipWhitespaces()
			comb = selector.ChildCombinator
		}
		anotherUnit, err := ts.parseCompoundSelector()
		if err != nil {
			break
		}
		rest = append(rest, selector.ComplexSelectorRest{Combinator: comb, Selector: anotherUnit})
	}
	return selector.ComplexSelector{Base: base, Rest: rest}, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-complex-selector-list
//
// Returns nil if not found
func (ts *tokenStream) parseComplexSelectorList() (res []selector.Selector, err error) {
	selList, err := parseCommaSeparatedRepeation(ts, 0, "complex selector", func(ts *tokenStream) (selector.ComplexSelector, error) {
		return ts.parseComplexSelector()
	})
	if err != nil {
		return nil, err
	}
	for _, s := range selList {
		res = append(res, s)
	}
	return res, nil
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#typedef-selector-list
func (ts *tokenStream) parseSelectorList() (res []selector.Selector, err error) {
	return ts.parseComplexSelectorList()
}

// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#parse-a-selector
func parseSelector(src string, sourceStr string) (res []selector.Selector, err error) {
	ts, err := tokenize([]byte(src), sourceStr)
	if err != nil {
		return nil, err
	}
	return parse(&ts, func(ts *tokenStream) ([]selector.Selector, error) {
		return ts.parseSelectorList()
	})
}
