// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/inseo-oh/yw/css/csssyntax"
)

var filename = flag.String("file", "", "Name of the CSS file")

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
	stylesheet, err := csssyntax.ParseStylesheet(bytes, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("stylesheet location:", stylesheet.Location)
	for i, rule := range stylesheet.StyleRules {
		selectorListStr := strings.Builder{}
		for i, s := range rule.SelectorList {
			if i != 0 {
				selectorListStr.WriteString(", ")
			}
			selectorListStr.WriteString(fmt.Sprintf("%v", s))
		}
		log.Printf("style-rule[%d](%s) {", i, selectorListStr.String())
		log.Printf("	declarations {")
		for _, decl := range rule.Declarations {
			log.Printf("        %s : %v", decl.Name, decl.Value)
		}
		log.Printf("    }")
		log.Printf("	at-rules {")
		for _, rule := range rule.AtRules {
			log.Printf("		   <name>: %s", rule.Name)
			log.Printf("		<prelude>: %s", rule.Prelude)
			log.Printf("		  <value>: %s", rule.Value)
		}
		log.Printf("    }")
		log.Printf("}")
	}
}
