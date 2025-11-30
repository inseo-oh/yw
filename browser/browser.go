package browser

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
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/html/elements"
	"github.com/inseo-oh/yw/html/htmlparser"
	"github.com/inseo-oh/yw/layout"
	"github.com/inseo-oh/yw/namespaces"
	"github.com/inseo-oh/yw/platform"
)

type Browser struct{}

func (b *Browser) Run(urlStr string, plat platform.Platform, viewportImg *image.RGBA) {
	debugBuiltinStylesheet := false

	// Load the default CSS ----------------------------------------------------
	log.Println("Loading default CSS")
	bytes, err := os.ReadFile("res/default.css")
	if err != nil {
		log.Fatal(err)
	}
	tokens, err := csssyntax.Tokenize(csssyntax.DecodeBytes(bytes))
	if err != nil {
		log.Printf("<style>: failed to tokenize stylesheet: %v", err)
		return
	}
	// TODO: Can't we pass dom.Document instead?
	// Also, should <html> own the default stylesheet?
	initDefaultCss := func(htm elements.HTMLElement) cssom.Stylesheet {
		log.Println("Parsing default CSS")
		stylesheet := csssyntax.ParseStylesheet(tokens, nil)
		stylesheet.Type = "text/css"
		stylesheet.OwnerNode = htm
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
		return stylesheet
	}

	// Fetch the document ------------------------------------------------------
	log.Println("Loading document at", urlStr)
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Parse the HTML ----------------------------------------------------------
	html := string(bytes)
	par := htmlparser.NewParser(html)
	par.Document = dom.NewDocument()
	par.Document.SetBaseURL(*urlObj)
	doc := par.Run()

	// Find the <html> element -------------------------------------------------
	htmlElem := doc.FilterElementChildrenByLocalName(dom.NamePair{Namespace: namespaces.Html, LocalName: "html"})[0]
	uaStylesheet := initDefaultCss(htmlElem.(elements.HTMLElement))

	// Find the <head> element -------------------------------------------------
	headElem := htmlElem.FilterElementChildrenByLocalName(dom.NamePair{Namespace: namespaces.Html, LocalName: "head"})[0]

	// Apply style rules -------------------------------------------------------
	cascade.ApplyStyleRules(&uaStylesheet, doc)
	log.Println("Style rules applied")

	// Do something with it ----------------------------------------------------
	_ = headElem
	dom.PrintTree(doc)
	log.Println("Document loaded. Making layout tree...")
	viewportSize := viewportImg.Rect.Size()
	for y := range viewportSize.Y {
		for x := range viewportSize.X {
			viewportImg.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	icb := layout.Build(htmlElem, float64(viewportSize.X), float64(viewportSize.Y), plat)
	layout.PrintTree(icb)
	log.Println("Made layout. Making paint tree...")
	paint := icb.MakePaintNode()
	gfx.PrintPaintTree(paint)
	log.Println("Painting...")
	paint.Paint(viewportImg)

	log.Println("DONE")
}
