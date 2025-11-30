package libhtml

import (
	"log"
	"testing"
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
				&html_tag_token{tag_name: "tag-without-attr", attrs: []dom_Attr_s{}},
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
				&html_tag_token{tag_name: "tag-with-attrs", attrs: []dom_Attr_s{
					{local_name: "attr-without-value", value: ""},
					{local_name: "attr-with-unquoted-value", value: "abc"},
					{local_name: "attr-with-single-quote", value: "def\""},
					{local_name: "attr-with-double-quote", value: "ghi'"},
				}},
				&html_eof_token{},
			},
		},
		{
			"</end-tag>", []html_token{
				&html_tag_token{tag_name: "end-tag", is_end: true, attrs: []dom_Attr_s{}},
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

	make_doc := func(mode dom_Document_mode, init_children func(doc dom_Document) []dom_Node) dom_Document {
		doc := dom_make_Document()
		doc.set_mode(mode)
		children := init_children(doc)
		for _, c := range children {
			dom_node_append_child(doc, c)
		}
		return doc
	}
	make_elem := func(doc dom_Document, local_name string, namespace *namespace, children []dom_Node) dom_Element {
		elem := dom_make_Element(dom_element_creation_common_options{
			node_document: doc,
			local_name:    local_name, namespace: namespace,
		}, dom_element_callbacks{}, true)
		for _, c := range children {
			dom_node_append_child(elem, c)
		}
		return elem
	}
	test_cases := []struct {
		input    string
		dom_root dom_Node
	}{
		// Empty HTML --> Should generate basic DOM, with document being in quirks mode
		{"", make_doc(dom_Document_mode_quirks,
			func(doc dom_Document) []dom_Node {
				return []dom_Node{
					make_elem(doc, "html", html_namespace_p(), []dom_Node{
						make_elem(doc, "head", html_namespace_p(), nil),
						make_elem(doc, "body", html_namespace_p(), nil),
					}),
				}
			},
		)},
		// Empty HTML, but with DOCTYPE --> Should generate basic DOM, with document being in no-quirks mode
		{"<!doctype html>", make_doc(dom_Document_mode_no_quirks,
			func(doc dom_Document) []dom_Node {
				return []dom_Node{
					dom_make_DocumentType(doc, "html", "", ""),
					make_elem(doc, "html", html_namespace_p(), []dom_Node{
						make_elem(doc, "head", html_namespace_p(), nil),
						make_elem(doc, "body", html_namespace_p(), nil),
					}),
				}
			},
		)},
		// Normal basic HTML structure
		{"<!doctype html><html><head></head><body></body></html>", make_doc(
			dom_Document_mode_no_quirks,
			func(doc dom_Document) []dom_Node {
				return []dom_Node{
					dom_make_DocumentType(doc, "html", "", ""),
					make_elem(doc, "html", html_namespace_p(), []dom_Node{
						make_elem(doc, "head", html_namespace_p(), nil),
						make_elem(doc, "body", html_namespace_p(), nil),
					}),
				}
			},
		)},
		// Text node
		{"<!doctype html><body>abc", make_doc(dom_Document_mode_no_quirks,
			func(doc dom_Document) []dom_Node {
				return []dom_Node{
					dom_make_DocumentType(doc, "html", "", ""),
					make_elem(doc, "html", html_namespace_p(), []dom_Node{
						make_elem(doc, "head", html_namespace_p(), nil),
						make_elem(doc, "body", html_namespace_p(), []dom_Node{
							dom_make_Text(doc, "abc"),
						}),
					}),
				}
			},
		)},
		// <h1>~<h6>
		{"<!doctype html><body><h1><h2><h3><h4><h5><h6>", make_doc(dom_Document_mode_no_quirks,
			func(doc dom_Document) []dom_Node {
				return []dom_Node{
					dom_make_DocumentType(doc, "html", "", ""),
					make_elem(doc, "html", html_namespace_p(), []dom_Node{
						make_elem(doc, "head", html_namespace_p(), nil),
						make_elem(doc, "body", html_namespace_p(), []dom_Node{
							make_elem(doc, "h1", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "h2", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "h3", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "h4", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "h5", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "h6", html_namespace_p(), []dom_Node{}),
						}),
					}),
				}
			},
		)},
		// <li>,<dd>,<dt>
		{"<!doctype html><body><li><li><li><li><dd><dt><dd><dt><dd><dt>", make_doc(dom_Document_mode_no_quirks,
			func(doc dom_Document) []dom_Node {
				return []dom_Node{
					dom_make_DocumentType(doc, "html", "", ""),
					make_elem(doc, "html", html_namespace_p(), []dom_Node{
						make_elem(doc, "head", html_namespace_p(), []dom_Node{}),
						make_elem(doc, "body", html_namespace_p(), []dom_Node{
							make_elem(doc, "li", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "li", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "li", html_namespace_p(), []dom_Node{}),
							make_elem(doc, "li", html_namespace_p(), []dom_Node{
								make_elem(doc, "dd", html_namespace_p(), []dom_Node{}),
								make_elem(doc, "dt", html_namespace_p(), []dom_Node{}),
								make_elem(doc, "dd", html_namespace_p(), []dom_Node{}),
								make_elem(doc, "dt", html_namespace_p(), []dom_Node{}),
								make_elem(doc, "dd", html_namespace_p(), []dom_Node{}),
								make_elem(doc, "dt", html_namespace_p(), []dom_Node{}),
							}),
						}),
					}),
				}
			},
		)},
	}
	var fix_children_parent_ptr func(node dom_Node)
	fix_children_parent_ptr = func(node dom_Node) {
		for _, child_p := range node.get_children() {
			child_p.set_parent(node)
			fix_children_parent_ptr(child_p)
		}
	}
	for _, cs := range test_cases {
		fix_children_parent_ptr(cs.dom_root)

		par := html_make_parser(cs.input)
		par.Run()
		exp := dom_node_inclusive_descendants(cs.dom_root)
		got := dom_node_inclusive_descendants(par.document)
		failed := false
		if len(exp) != len(got) {
			log.Printf("Expected %v nodes, got %v nodes", len(exp), len(got))
			failed = true
		} else {
			for i := 0; i < len(exp); i++ {
				if exp_n, ok := exp[i].(dom_Document); ok {
					if got_n, ok := got[i].(dom_Document); ok {
						if exp_n.get_mode() != got_n.get_mode() {
							failed = true
						}
					} else {
						failed = true
					}
				} else if exp_n, ok := exp[i].(dom_DocumentType); ok {
					if got_n, ok := got[i].(dom_DocumentType); ok {
						if exp_n.get_name() != got_n.get_name() {
							failed = true
						}
						if exp_n.get_public_id() != got_n.get_public_id() {
							failed = true
						}
						if exp_n.get_system_id() != got_n.get_system_id() {
							failed = true
						}
					} else {
						failed = true
					}
				} else if exp_n, ok := exp[i].(dom_Element); ok {
					if got_n, ok := got[i].(dom_Element); ok {
						exp_ns, exp_has_ns := exp_n.get_namespace()
						got_ns, got_has_ns := got_n.get_namespace()
						if exp_has_ns != got_has_ns {
							failed = true
						}
						if exp_has_ns && exp_ns != got_ns {
							failed = true
						}
						if exp_n.get_local_name() != got_n.get_local_name() {
							failed = true
						}
						// TODO: Compare attributes

					} else {
						failed = true
					}
				} else if exp_n, ok := exp[i].(dom_Text); ok {
					if got_n, ok := got[i].(dom_Text); ok {
						if exp_n.get_text() != got_n.get_text() {
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
			dom_print_tree(par.document)
			log.Print("---------- Expected ----------")
			dom_print_tree(cs.dom_root)
			t.Fail()
		}

	}
}
