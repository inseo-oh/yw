package libhtml

import (
	"fmt"
	"log"
	"strings"
)

type dom_to_md struct {
	blockquote_level   int
	list_nesting_level int
}

func (h *dom_to_md) run(sb *strings.Builder, node dom_Node) {
	if e, ok := node.(dom_Element); ok {
		make_content_for_children := func() {
			for _, child := range e.get_children() {
				h.run(sb, child)
			}
		}
		if 1 < h.list_nesting_level {
			sb.WriteString("    ")
		}
		if e.is_html_element("p") {
			if h.blockquote_level == 0 {
				sb.WriteString("\n")
				make_content_for_children()
				sb.WriteString("\n")
			} else {
				h.blockquote_level++
				sb.WriteString("\n")
				for i := 0; i < h.blockquote_level; i++ {
					sb.WriteString(">")
				}
				make_content_for_children()
				sb.WriteString("\n")
				h.blockquote_level--
			}
		} else if e.is_html_element("strong") || e.is_html_element("b") {
			sb.WriteString("**")
			make_content_for_children()
			sb.WriteString("**")
		} else if e.is_html_element("em") || e.is_html_element("i") {
			sb.WriteString("*")
			make_content_for_children()
			sb.WriteString("*")
		} else if e.is_html_element("h1") {
			sb.WriteString("\n\n# ")
			make_content_for_children()
			sb.WriteString("\n\n")
		} else if e.is_html_element("h2") {
			sb.WriteString("\n\n## ")
			make_content_for_children()
			sb.WriteString("\n\n")
		} else if e.is_html_element("h3") {
			sb.WriteString("\n\n### ")
			make_content_for_children()
			sb.WriteString("\n\n")
		} else if e.is_html_element("h4") {
			sb.WriteString("\n\n#### ")
			make_content_for_children()
			sb.WriteString("\n\n")
		} else if e.is_html_element("h5") {
			sb.WriteString("\n\n##### ")
			make_content_for_children()
			sb.WriteString("\n\n")
		} else if e.is_html_element("h6") {
			sb.WriteString("\n\n###### ")
			make_content_for_children()
			sb.WriteString("\n\n")
		} else if e.is_html_element("br") {
			sb.WriteString("  \n")
		} else if e.is_html_element("a") {
			sb.WriteString("[")
			make_content_for_children()
			sb.WriteString("](")
			href, ok := e.get_attribute_without_namespace("href")
			if !ok {
				sb.WriteString(href)
			} else {
				sb.WriteString("#")
			}
			sb.WriteString(")")
		} else if e.is_html_element("blockquote") {
			sb.WriteString("\n>")
			h.blockquote_level++
			make_content_for_children()
			h.blockquote_level--
			sb.WriteString("\n\n")
		} else if e.is_html_element("ul") {
			h.list_nesting_level++
			make_content_for_children()
			h.list_nesting_level--
		} else if e.is_html_element("li") {
			sb.WriteString("- ")
			make_content_for_children()
			sb.WriteString("\n")
		} else if e.is_html_element("head") {
			// Ignore
		} else {
			make_content_for_children()
		}
	} else if e, ok := node.(dom_Text); ok {
		text := e.get_text()
		text = strings.ReplaceAll(text, "\t", " ")
		text = strings.ReplaceAll(text, "\n", " ")
		for strings.Contains(text, "  ") {
			text = strings.ReplaceAll(text, "  ", " ")
		}
		if text == "" && e.get_text() != "" {
			text = " "
		}
		sb.WriteString(text)
	} else if e := node; e != nil {
		for _, child := range e.get_children() {
			h.run(sb, child)
		}
	}
}

// This is mostly a PoC, rather than "functional" one.
func Html2Md(html string) string {
	par := html_make_parser(html)
	doc := par.Run()
	sb := strings.Builder{}
	h := dom_to_md{}
	h.run(&sb, doc)
	return sb.String()
}

func DomDump(html string) {
	par := html_make_parser(html)
	doc := par.Run()
	dom_print_tree(doc)
}

func CssDump(css string) {
	tokens, err := css_tokenize(css)
	if err != nil {
		log.Fatal(err)
	}
	stylesheet := css_parse_stylesheet(tokens, nil)
	log.Println("stylesheet location:", stylesheet.location)
	for i, rule := range stylesheet.style_rules {
		selector_list_str := strings.Builder{}
		for i, s := range rule.selector_list {
			if i != 0 {
				selector_list_str.WriteString(", ")
			}
			selector_list_str.WriteString(fmt.Sprintf("%v", s))
		}
		log.Printf("style-rule[%d](%s) {", i, selector_list_str.String())
		log.Printf("	declarations {")
		for _, decl := range rule.declarations {
			log.Printf("        %s : %v", decl.name, decl.value)
		}
		log.Printf("    }")
		log.Printf("	at-rules {")
		for _, rule := range rule.at_rules {
			log.Printf("		   <name>: %s", rule.name)
			log.Printf("		<prelude>: %s", rule.prelude)
			log.Printf("		  <value>: %s", rule.value)
		}
		log.Printf("    }")
		log.Printf("}")
	}
}
