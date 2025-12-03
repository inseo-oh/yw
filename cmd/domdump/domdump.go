// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

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
