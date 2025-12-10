// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// yw is a web page rendering engine aimed at simplicity.
package yw

import (
	"image"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/inseo-oh/yw/css/cascade"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/csssyntax"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx/paint"
	"github.com/inseo-oh/yw/html/htmlparser"
	"github.com/inseo-oh/yw/layout"
	"github.com/inseo-oh/yw/namespaces"
	"github.com/inseo-oh/yw/platform"
)

func loadUserAgentCss() *cssom.Stylesheet {
	debugBuiltinStylesheet := false

	log.Println("Reading UA CSS")
	sheetBytes, err := os.ReadFile("res/default.css")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsing UA CSS")
	stylesheet, err := csssyntax.ParseStylesheet(sheetBytes, nil, "<UA stylesheet>")
	if err != nil {
		log.Panicf("failed to parse UA stylesheet: %v", err)
	}
	stylesheet.Type = "text/css"
	stylesheet.OwnerNode = nil
	// TODO: Set stylesheet.media once we implement that
	stylesheet.AlternateFlag = false
	stylesheet.OriginCleanFlag = true
	stylesheet.Location = nil
	stylesheet.ParentStylesheet = nil
	stylesheet.OwnerRule = nil

	if debugBuiltinStylesheet {
		log.Println("dump of builtin stylesheet")
		stylesheet.Dump()
	}
	return &stylesheet
}

// Browser represents state of the browser.
type Browser struct {
	DumpDom    bool // Dump DOM tree?
	DumpLayout bool // Dump layout tree?
	DumpPaint  bool // Dump paint tree?
}

// Run loads the document from urlStr URL, and renders resulting document to viewportImg.
func (b *Browser) Run(urlStr string, fontProvider platform.FontProvider, viewportImg *image.RGBA) {
	log.Println("= Loading user agent CSS ====================================")
	uaStylesheet := loadUserAgentCss()

	// Fetch the document ------------------------------------------------------
	log.Println("= Fetching document =========================================")
	log.Printf("Document URL: %s", urlStr)
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Parse the HTML ----------------------------------------------------------
	log.Println("= Parsing document ==========================================")
	html := string(htmlBytes)
	par := htmlparser.NewParser(html)
	par.Document = dom.NewDocument()
	par.Document.SetBaseURL(*urlObj)
	doc := par.Run()
	log.Println("= Document parsed ===========================================")
	if b.DumpDom {
		dom.PrintTree(doc, 0)
	}

	// Apply style rules -------------------------------------------------------
	log.Println("= Applying style rules ======================================")
	cascade.ApplyStyleRules(uaStylesheet, doc)

	// Do something with it ----------------------------------------------------
	log.Println("= Building layout tree ======================================")
	viewportSize := viewportImg.Rect.Size()
	for y := range viewportSize.Y {
		for x := range viewportSize.X {
			viewportImg.Set(x, y, color.White)
		}
	}

	htmlElem := doc.FilterElementChildrenByLocalName(dom.NamePair{Namespace: namespaces.Html, LocalName: "html"})[0]
	icb := layout.Build(htmlElem, float64(viewportSize.X), float64(viewportSize.Y), fontProvider)
	if b.DumpLayout {
		layout.PrintTree(icb, 0)
	}
	log.Println("= Building paint tree =======================================")
	paintNode := icb.MakePaintNode()
	if b.DumpPaint {
		paint.PrintTree(paintNode, 0)
	}
	log.Println("= Painting ==================================================")
	paintNode.Paint(viewportImg)

	log.Println("DONE")
}
