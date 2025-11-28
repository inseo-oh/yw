package parser

import (
	"log"
	"testing"
	"yw/dom"
	"yw/namespaces"
)

func TestHtmlTokenizer(t *testing.T) {
	test_cases := []struct {
		input  string
		tokens []html_token
	}{
		// Character token -----------------------------------------------------
		{
			"char", []html_token{
				&html_char_token{value: 'c'},
				&html_char_token{value: 'h'},
				&html_char_token{value: 'a'},
				&html_char_token{value: 'r'},
				&html_eof_token{},
			},
		},
		// Character reference -------------------------------------------------
		{
			"&#44032;", []html_token{
				&html_char_token{value: 0xac00},
				&html_eof_token{},
			},
		},
		{
			"&#xAc00;", []html_token{
				&html_char_token{value: 0xac00},
				&html_eof_token{},
			},
		},
		{
			"&lt;", []html_token{
				&html_char_token{value: '<'},
				&html_eof_token{},
			},
		},
		// Bogus comment -------------------------------------------------------
		{
			"<?$comment=bogus>", []html_token{
				&html_comment_token{data: "?$comment=bogus"},
				&html_eof_token{},
			},
		},
		// Proper comment ------------------------------------------------------
		{
			"<!--good-comment-->", []html_token{
				&html_comment_token{data: "good-comment"},
				&html_eof_token{},
			},
		},
		// DOCTYPE -------------------------------------------------------------
		{
			"<!DocType someName>", []html_token{
				&html_doctype_token{name: func() *string {
					s := "somename"
					return &s
				}()},
				&html_eof_token{},
			},
		},
		// Tag -----------------------------------------------------------------
		{
			"<tag-without-attr>", []html_token{
				&tagToken{tag_name: "tag-without-attr", attrs: []dom.AttrData{}},
				&html_eof_token{},
			},
		},
		{
			`<tag-with-attrs
				attr-without-value
				attr-with-unquoted-value=abc
				attr-with-single-quote='def"'
				attr-with-double-quote="ghi'"
				>`, []html_token{
				&tagToken{tag_name: "tag-with-attrs", attrs: []dom.AttrData{
					{LocalName: "attr-without-value", Value: ""},
					{LocalName: "attr-with-unquoted-value", Value: "abc"},
					{LocalName: "attr-with-single-quote", Value: "def\""},
					{LocalName: "attr-with-double-quote", Value: "ghi'"},
				}},
				&html_eof_token{},
			},
		},
		{
			"</end-tag>", []html_token{
				&tagToken{tag_name: "end-tag", is_end: true, attrs: []dom.AttrData{}},
				&html_eof_token{},
			},
		},
	}
	for _, cs := range test_cases {
		t.Run(cs.input, func(t *testing.T) {
			tk := html_make_tokenizer(cs.input)
			tk.on_token_emitted = func(got html_token) {
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

	make_doc := func(mode dom.DocumentMode, init_children func(doc dom.Document) []dom.Node) dom.Document {
		doc := dom.NewDocument()
		doc.SetMode(mode)
		children := init_children(doc)
		for _, c := range children {
			dom.AppendChild(doc, c)
		}
		return doc
	}
	make_elem := func(doc dom.Document, localName string, namespace *namespaces.Namespace, children []dom.Node) dom.Element {
		elem := dom.NewElement(dom.ElementCreationCommonOptions{
			NodeDocument: doc, LocalName: localName, Namespace: namespace,
		})
		for _, c := range children {
			dom.AppendChild(elem, c)
		}
		return elem
	}
	test_cases := []struct {
		input    string
		dom_root dom.Node
	}{
		// Empty HTML --> Should generate basic DOM, with document being in quirks mode
		{"", make_doc(dom.Quirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					make_elem(doc, "html", namespaces.HtmlP(), []dom.Node{
						make_elem(doc, "head", namespaces.HtmlP(), nil),
						make_elem(doc, "body", namespaces.HtmlP(), nil),
					}),
				}
			},
		)},
		// Empty HTML, but with DOCTYPE --> Should generate basic DOM, with document being in no-quirks mode
		{"<!doctype html>", make_doc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					make_elem(doc, "html", namespaces.HtmlP(), []dom.Node{
						make_elem(doc, "head", namespaces.HtmlP(), nil),
						make_elem(doc, "body", namespaces.HtmlP(), nil),
					}),
				}
			},
		)},
		// Normal basic HTML structure
		{"<!doctype html><html><head></head><body></body></html>", make_doc(
			dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					make_elem(doc, "html", namespaces.HtmlP(), []dom.Node{
						make_elem(doc, "head", namespaces.HtmlP(), nil),
						make_elem(doc, "body", namespaces.HtmlP(), nil),
					}),
				}
			},
		)},
		// Text node
		{"<!doctype html><body>abc", make_doc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					make_elem(doc, "html", namespaces.HtmlP(), []dom.Node{
						make_elem(doc, "head", namespaces.HtmlP(), nil),
						make_elem(doc, "body", namespaces.HtmlP(), []dom.Node{
							dom.NewText(doc, "abc"),
						}),
					}),
				}
			},
		)},
		// <h1>~<h6>
		{"<!doctype html><body><h1><h2><h3><h4><h5><h6>", make_doc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					make_elem(doc, "html", namespaces.HtmlP(), []dom.Node{
						make_elem(doc, "head", namespaces.HtmlP(), nil),
						make_elem(doc, "body", namespaces.HtmlP(), []dom.Node{
							make_elem(doc, "h1", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "h2", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "h3", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "h4", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "h5", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "h6", namespaces.HtmlP(), []dom.Node{}),
						}),
					}),
				}
			},
		)},
		// <li>,<dd>,<dt>
		{"<!doctype html><body><li><li><li><li><dd><dt><dd><dt><dd><dt>", make_doc(dom.NoQuirks,
			func(doc dom.Document) []dom.Node {
				return []dom.Node{
					dom.NewDocumentType(doc, "html", "", ""),
					make_elem(doc, "html", namespaces.HtmlP(), []dom.Node{
						make_elem(doc, "head", namespaces.HtmlP(), []dom.Node{}),
						make_elem(doc, "body", namespaces.HtmlP(), []dom.Node{
							make_elem(doc, "li", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "li", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "li", namespaces.HtmlP(), []dom.Node{}),
							make_elem(doc, "li", namespaces.HtmlP(), []dom.Node{
								make_elem(doc, "dd", namespaces.HtmlP(), []dom.Node{}),
								make_elem(doc, "dt", namespaces.HtmlP(), []dom.Node{}),
								make_elem(doc, "dd", namespaces.HtmlP(), []dom.Node{}),
								make_elem(doc, "dt", namespaces.HtmlP(), []dom.Node{}),
								make_elem(doc, "dd", namespaces.HtmlP(), []dom.Node{}),
								make_elem(doc, "dt", namespaces.HtmlP(), []dom.Node{}),
							}),
						}),
					}),
				}
			},
		)},
	}
	var fix_children_parent_ptr func(node dom.Node)
	fix_children_parent_ptr = func(node dom.Node) {
		for _, child_p := range node.Children() {
			child_p.SetParent(node)
			fix_children_parent_ptr(child_p)
		}
	}
	for _, cs := range test_cases {
		fix_children_parent_ptr(cs.dom_root)

		par := NewParser(cs.input)
		par.Run()
		exp := dom.InclusiveDescendants(cs.dom_root)
		got := dom.InclusiveDescendants(par.document)
		failed := false
		if len(exp) != len(got) {
			log.Printf("Expected %v nodes, got %v nodes", len(exp), len(got))
			failed = true
		} else {
			for i := 0; i < len(exp); i++ {
				if exp_n, ok := exp[i].(dom.Document); ok {
					if got_n, ok := got[i].(dom.Document); ok {
						if exp_n.Mode() != got_n.Mode() {
							failed = true
						}
					} else {
						failed = true
					}
				} else if exp_n, ok := exp[i].(dom.DocumentType); ok {
					if got_n, ok := got[i].(dom.DocumentType); ok {
						if exp_n.Name() != got_n.Name() {
							failed = true
						}
						if exp_n.PublicId() != got_n.PublicId() {
							failed = true
						}
						if exp_n.SystemId() != got_n.SystemId() {
							failed = true
						}
					} else {
						failed = true
					}
				} else if exp_n, ok := exp[i].(dom.Element); ok {
					if got_n, ok := got[i].(dom.Element); ok {
						exp_ns, exp_has_ns := exp_n.Namespace()
						got_ns, got_has_ns := got_n.Namespace()
						if exp_has_ns != got_has_ns {
							failed = true
						}
						if exp_has_ns && exp_ns != got_ns {
							failed = true
						}
						if exp_n.LocalName() != got_n.LocalName() {
							failed = true
						}
						// TODO: Compare attributes

					} else {
						failed = true
					}
				} else if exp_n, ok := exp[i].(dom.Text); ok {
					if got_n, ok := got[i].(dom.Text); ok {
						if exp_n.Text() != got_n.Text() {
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
			dom.PrintTree(par.document)
			log.Print("---------- Expected ----------")
			dom.PrintTree(cs.dom_root)
			t.Fail()
		}

	}
}
