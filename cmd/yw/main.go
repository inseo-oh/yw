// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/inseo-oh/yw"
	"github.com/inseo-oh/yw/platform/linux"
)

var (
	url        = flag.String("url", "", "The URL")
	dumpDom    = flag.Bool("dumpdom", false, "Dump DOM tree")
	dumpLayout = flag.Bool("dumplayout", false, "Dump layout tree")
	dumpPaint  = flag.Bool("dumppaint", false, "Dump paint tree")
)

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	br := yw.Browser{
		DumpDom:    *dumpDom,
		DumpLayout: *dumpLayout,
		DumpPaint:  *dumpPaint,
	}
	viewportImg := image.NewRGBA(image.Rect(0, 0, 1280, 720))
	fontProvider := linux.NewDefaultFontProvider()
	br.Run(*url, fontProvider, viewportImg)

	destFile, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(destFile, viewportImg)
	if err != nil {
		log.Fatal(err)
	}
	destFile.Close()
}
