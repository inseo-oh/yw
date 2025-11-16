package main

import (
	"flag"
	"log"
	"os"
	"yw/libhtml"
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
	libhtml.DomDump(str)
}
