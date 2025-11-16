package main

import (
	"flag"
	"log"
	"os"
	"yw/libhtml"
)

var infile = flag.String("file", "", "Name of the HTML file")
var outfile = flag.String("out", "out.md", "Name of the output file")

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
	res := libhtml.Html2Md(string(bytes))
	err = os.WriteFile(*outfile, []byte(res), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// libhtml.DomPrintTree(libhtml.DomMakePtr(doc))
}
