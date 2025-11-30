package main

import (
	"flag"
	"log"
	"os"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/html/htmlparser"
)

var filename = flag.String("file", "", "Name of the HTML file")

func main() {
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	bytes, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}
	str := string(bytes)
	par := htmlparser.NewParser(str)
	doc := par.Run()
	dom.PrintTree(doc)
}
