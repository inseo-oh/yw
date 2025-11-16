package main

import (
	"flag"
	"os"
	"yw/libhtml"
)

var url = flag.String("url", "", "The URL")

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	br := libhtml.Browser{}
	plat := init_platform()
	br.Init(*url, plat)
}
