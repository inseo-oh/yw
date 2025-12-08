// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/inseo-oh/yw"
	"github.com/inseo-oh/yw/platform/stdplatform"
)

var url = flag.String("url", "", "The URL")

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	br := yw.Browser{}
	viewportImg := image.NewRGBA(image.Rect(0, 0, 1280, 720))
	fontProvider := stdplatform.NewDefaultFontProvider()
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
