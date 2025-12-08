// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package htmlparser

import (
	"log"
	"testing"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/namespaces"
)

func TestHtmlTokenizer(t *testing.T) {
	testCases := []struct {
		input  string
		tokens []htmlToken
	}{
		// Character token -----------------------------------------------------
		{
			"char", []htmlToken{
				&charToken{value: 'c'},
				&charToken{value: 'h'},
				&charToken{value: 'a'},
				&charToken{value: 'r'},
				&eofToken{},
			},
		},
		// Character reference -------------------------------------------------
		{
			"&#44032;", []htmlToken{
				&charToken{value: 0xac00},
				&eofToken{},
			},
		},
		{
			"&#xAc00;", []htmlToken{
				&charToken{value: 0xac00},
				&eofToken{},
			},
		},
		{
			"&lt;", []htmlToken{
				&charToken{value: '<'},
				&eofToken{},
			},
		},
		// Bogus comment -------------------------------------------------------
		{
			"<?$comment=bogus>", []htmlToken{
				&commentToken{data: "?$comment=bogus"},
				&eofToken{},
			},
		},
		// Proper comment ------------------------------------------------------
		{
			"<!--good-comment-->", []htmlToken{
				&commentToken{data: "good-comment"},
				&eofToken{},
			},
		},
		// DOCTYPE -------------------------------------------------------------
		{
			"<!DocType someName>", []htmlToken{
				&doctypeToken{name: func() *string {
					s := "somename"
					return &s
				}()},
				&eofToken{},
			},
		},
		// Tag -----------------------------------------------------------------
		{
			"<tag-without-attr>", []htmlToken{
				&tagToken{tagName: "tag-without-attr", attrs: []dom.AttrData{}},
				&eofToken{},
			},
		},
		{
			`<tag-with-attrs
				attr-without-value
				attr-with-unquoted-value=abc
				attr-with-single-quote='def"'
				attr-with-double-quote="ghi'"
				>`, []htmlToken{
				&tagToken{tagName: "tag-with-attrs", attrs: []dom.AttrData{
					{LocalName: "attr-without-value", Value: ""},
					{LocalName: "attr-with-unquoted-value", Value: "abc"},
					{LocalName: "attr-with-single-quote", Value: "def\""},
					{LocalName: "attr-with-double-quote", Value: "ghi'"},
				}},
				&eofToken{},
			},
		},
		{
			"</end-tag>", []htmlToken{
				&tagToken{tagName: "end-tag", isEnd: true, attrs: []dom.AttrData{}},
				&eofToken{},
			},
		},
	}
	for _, cs := range testCases {
		t.Run(cs.input, func(t *testing.T) {
			tk := newTokenizer(cs.input)
			tk.onTokenEmitted = func(got htmlToken) {
				if len(cs.tokens) == 0 {
					t.Errorf("[%s] too many tokens - got %v", cs.input, got)
					return
				} else if !got.equals(cs.tokens[0]) {
					t.Errorf("[%s] token mismatch - got %v, want %v", cs.input, got, cs.tokens[0])
				}
				cs.tokens = cs.tokens[1:]
			}
			tk.run()
			if len(cs.tokens) != 0 {
				t.Errorf("[%s] remaining %v token(s) found", cs.input, len(cs.tokens))
			}
		})
	}

}

func TestHtmlParser(t *testing.T) {
	makeDoc := func(mode dom.DocumentMode, initChildren func(doc dom.Document) []dom.Node) dom.Document {
		doc := dom.NewDocument()
		doc.SetMode(mode)
		children := initChildren(doc)
		for _, c := range children {
			dom.AppendChild(doc, c)
		}
		return doc
	}
	makeElem := func(doc dom.Document, localName string, namespace *namespaces.Namespace, children []dom.Node) dom.Element {
		elem := dom.NewElement(dom.ElementCreationCommonOptions{
			NodeDocument: doc, LocalName: localName, Namespace: namespace,
		})
		for _, c := range children {
			dom.AppendChild(elem, c)
		}
		return elem
	}
	testCases := []struct {
		input   string
		domToot dom.Node
	}{
		// Empty HTML --> Should generate basic DOM, with document being in quirks mode
		{"", makeDoc(dom.Quirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					makeElem(doc, "html", &namespaces.Html, []dom.Node{
						makeElem(doc, "head", &namespaces.Html, nil),
						makeElem(doc, "body", &namespaces.Html, nil),
					}),
				}
			},
		)},
		// Empty HTML, but with DOCTYPE --> Should generate basic DOM, with document being in no-quirks mode
		{"<!doctype html>", makeDoc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					makeElem(doc, "html", &namespaces.Html, []dom.Node{
						makeElem(doc, "head", &namespaces.Html, nil),
						makeElem(doc, "body", &namespaces.Html, nil),
					}),
				}
			},
		)},
		// Normal basic HTML structure
		{"<!doctype html><html><head></head><body></body></html>", makeDoc(
			dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					makeElem(doc, "html", &namespaces.Html, []dom.Node{
						makeElem(doc, "head", &namespaces.Html, nil),
						makeElem(doc, "body", &namespaces.Html, nil),
					}),
				}
			},
		)},
		// Text node
		{"<!doctype html><body>abc", makeDoc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					makeElem(doc, "html", &namespaces.Html, []dom.Node{
						makeElem(doc, "head", &namespaces.Html, nil),
						makeElem(doc, "body", &namespaces.Html, []dom.Node{
							dom.NewText(doc, "abc"),
						}),
					}),
				}
			},
		)},
		// <h1>~<h6>
		{"<!doctype html><body><h1><h2><h3><h4><h5><h6>", makeDoc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					makeElem(doc, "html", &namespaces.Html, []dom.Node{
						makeElem(doc, "head", &namespaces.Html, nil),
						makeElem(doc, "body", &namespaces.Html, []dom.Node{
							makeElem(doc, "h1", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "h2", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "h3", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "h4", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "h5", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "h6", &namespaces.Html, []dom.Node{}),
						}),
					}),
				}
			},
		)},
		// <li>,<dd>,<dt>
		{"<!doctype html><body><li><li><li><li><dd><dt><dd><dt><dd><dt>", makeDoc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					makeElem(doc, "html", &namespaces.Html, []dom.Node{
						makeElem(doc, "head", &namespaces.Html, []dom.Node{}),
						makeElem(doc, "body", &namespaces.Html, []dom.Node{
							makeElem(doc, "li", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "li", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "li", &namespaces.Html, []dom.Node{}),
							makeElem(doc, "li", &namespaces.Html, []dom.Node{
								makeElem(doc, "dd", &namespaces.Html, []dom.Node{}),
								makeElem(doc, "dt", &namespaces.Html, []dom.Node{}),
								makeElem(doc, "dd", &namespaces.Html, []dom.Node{}),
								makeElem(doc, "dt", &namespaces.Html, []dom.Node{}),
								makeElem(doc, "dd", &namespaces.Html, []dom.Node{}),
								makeElem(doc, "dt", &namespaces.Html, []dom.Node{}),
							}),
						}),
					}),
				}
			},
		)},
	}
	var fixChildrenParentPtr func(node dom.Node)
	fixChildrenParentPtr = func(node dom.Node) {
		for _, child := range node.Children() {
			child.SetParent(node)
			fixChildrenParentPtr(child)
		}
	}
	for _, cs := range testCases {
		fixChildrenParentPtr(cs.domToot)

		par := NewParser(cs.input)
		par.Run()
		exp := dom.InclusiveDescendants(cs.domToot)
		got := dom.InclusiveDescendants(par.Document)
		failed := false
		if len(exp) != len(got) {
			log.Printf("Expected %v nodes, got %v nodes", len(exp), len(got))
			failed = true
		} else {
			for i := range exp {
				if expN, ok := exp[i].(dom.Document); ok {
					if gotN, ok := got[i].(dom.Document); ok {
						if expN.Mode() != gotN.Mode() {
							failed = true
						}
					} else {
						failed = true
					}
				} else if expN, ok := exp[i].(dom.DocumentType); ok {
					if gotN, ok := got[i].(dom.DocumentType); ok {
						if expN.Name() != gotN.Name() {
							failed = true
						}
						if expN.PublicId() != gotN.PublicId() {
							failed = true
						}
						if expN.SystemId() != gotN.SystemId() {
							failed = true
						}
					} else {
						failed = true
					}
				} else if expN, ok := exp[i].(dom.Element); ok {
					if gotN, ok := got[i].(dom.Element); ok {
						expNs, expHasNs := expN.Namespace()
						gotNs, gotHasNs := gotN.Namespace()
						if expHasNs != gotHasNs {
							failed = true
						}
						if expHasNs && expNs != gotNs {
							failed = true
						}
						if expN.LocalName() != gotN.LocalName() {
							failed = true
						}
						// TODO: Compare attributes

					} else {
						failed = true
					}
				} else if expN, ok := exp[i].(dom.Text); ok {
					if gotN, ok := got[i].(dom.Text); ok {
						if expN.Text() != gotN.Text() {
							failed = true
						}
					} else {
						failed = true
					}
				} else {
					log.Printf("Unexpected node %v in the test case!", exp[i])
					failed = true
				}
				if failed {
					log.Printf("Expected %v, got %v", exp[i], got[i])
				}
			}
		}
		if failed {
			log.Print("---------- Got ----------")
			dom.PrintTree(par.Document)
			log.Print("---------- Expected ----------")
			dom.PrintTree(cs.domToot)
			t.Fail()
		}

	}
}
