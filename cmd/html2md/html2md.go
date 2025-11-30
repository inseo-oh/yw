package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/html/htmlparser"
)

var (
	infile  = flag.String("file", "", "Name of the HTML file")
	outfile = flag.String("out", "out.md", "Name of the output file")
)

type domToMd struct {
	blockquoteLevel  int
	listNestingLevel int
}

func (h *domToMd) run(sb *strings.Builder, node dom.Node) {
	if e, ok := node.(dom.Element); ok {
		makeContentForChildren := func() {
			for _, child := range e.Children() {
				h.run(sb, child)
			}
		}
		if 1 < h.listNestingLevel {
			sb.WriteString("    ")
		}
		if e.IsHtmlElement("p") {
			if h.blockquoteLevel == 0 {
				sb.WriteString("\n")
				makeContentForChildren()
				sb.WriteString("\n")
			} else {
				h.blockquoteLevel++
				sb.WriteString("\n")
				for i := 0; i < h.blockquoteLevel; i++ {
					sb.WriteString(">")
				}
				makeContentForChildren()
				sb.WriteString("\n")
				h.blockquoteLevel--
			}
		} else if e.IsHtmlElement("strong") || e.IsHtmlElement("b") {
			sb.WriteString("**")
			makeContentForChildren()
			sb.WriteString("**")
		} else if e.IsHtmlElement("em") || e.IsHtmlElement("i") {
			sb.WriteString("*")
			makeContentForChildren()
			sb.WriteString("*")
		} else if e.IsHtmlElement("h1") {
			sb.WriteString("\n\n# ")
			makeContentForChildren()
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("h2") {
			sb.WriteString("\n\n## ")
			makeContentForChildren()
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("h3") {
			sb.WriteString("\n\n### ")
			makeContentForChildren()
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("h4") {
			sb.WriteString("\n\n#### ")
			makeContentForChildren()
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("h5") {
			sb.WriteString("\n\n##### ")
			makeContentForChildren()
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("h6") {
			sb.WriteString("\n\n###### ")
			makeContentForChildren()
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("br") {
			sb.WriteString("  \n")
		} else if e.IsHtmlElement("a") {
			sb.WriteString("[")
			makeContentForChildren()
			sb.WriteString("](")
			href, ok := e.AttrWithoutNamespace("href")
			if !ok {
				sb.WriteString(href)
			} else {
				sb.WriteString("#")
			}
			sb.WriteString(")")
		} else if e.IsHtmlElement("blockquote") {
			sb.WriteString("\n>")
			h.blockquoteLevel++
			makeContentForChildren()
			h.blockquoteLevel--
			sb.WriteString("\n\n")
		} else if e.IsHtmlElement("ul") {
			h.listNestingLevel++
			makeContentForChildren()
			h.listNestingLevel--
		} else if e.IsHtmlElement("li") {
			sb.WriteString("- ")
			makeContentForChildren()
			sb.WriteString("\n")
		} else if e.IsHtmlElement("head") {
			// Ignore
		} else {
			makeContentForChildren()
		}
	} else if e, ok := node.(dom.Text); ok {
		text := e.Text()
		text = strings.ReplaceAll(text, "\t", " ")
		text = strings.ReplaceAll(text, "\n", " ")
		for strings.Contains(text, "  ") {
			text = strings.ReplaceAll(text, "  ", " ")
		}
		if text == "" && e.Text() != "" {
			text = " "
		}
		sb.WriteString(text)
	} else if e := node; e != nil {
		for _, child := range e.Children() {
			h.run(sb, child)
		}
	}
}

func main() {
	flag.Parse()

	if *infile == "" {
		flag.Usage()
		os.Exit(1)
	}
	bytes, err := os.ReadFile(*infile)
	if err != nil {
		log.Fatal(err)
	}

	par := htmlparser.NewParser(string(bytes))
	doc := par.Run()
	sb := strings.Builder{}
	h := domToMd{}
	h.run(&sb, doc)
	res := sb.String()

	err = os.WriteFile(*outfile, []byte(res), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
