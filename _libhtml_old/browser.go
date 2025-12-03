package libhtml

import (
	"image"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/platform"
)

type Browser struct{}

func (b *Browser) Init(url_str string, plat platform.Platform, viewport_img *image.RGBA) {
	debug_builtin_stylesheet := false

	// Show a cool banner ------------------------------------------------------
	log.Println("█ █ █ █████ █    █████ █████ █████ █████")
	log.Println("█ █ █ █     █    █     █   █ █ █ █ █    ")
	log.Println("█ █ █ █████ █    █     █   █ █ █ █ █████")
	log.Println("█ █ █ █     █    █     █   █ █ █ █ █    ")
	log.Println(" █ █  █████ ████ █████ █████ █ █ █ █████")

	// Load the default CSS ----------------------------------------------------
	log.Println("Loading default CSS")
	bytes, err := os.ReadFile("res/default.css")
	if err != nil {
		log.Fatal(err)
	}
	tokens, err := css_tokenize(css_decode_bytes(bytes))
	if err != nil {
		log.Printf("<style>: failed to tokenize stylesheet: %v", err)
		return
	}
	// TODO: Can't we pass dom_Document instead?
	// Also, should <html> own the default stylesheet?
	init_default_css := func(html html_HTMLElement) css_stylesheet {
		log.Println("Parsing default CSS")
		stylesheet := css_parse_stylesheet(tokens, nil)
		stylesheet.tp = "text/css"
		stylesheet.owner_node = html
		// TODO: Set stylesheet.media once we implement that
		if dom_node_is_in_document_tree(html) {
			if attr, ok := html.get_attribute_without_namespace("title"); ok {
				stylesheet.title = attr
			}
		}
		stylesheet.alternate_flag = false
		stylesheet.origin_clean_flag = true
		stylesheet.location = nil
		stylesheet.parent_stylesheet = nil
		stylesheet.owner_rule = nil
		// css_add_stylesheet(&stylesheet)

		if debug_builtin_stylesheet {
			log.Println("dump of builtin stylesheet")
			stylesheet.dump()
		}
		return stylesheet
	}

	// Fetch the document ------------------------------------------------------
	log.Println("Loading document at", url_str)
	url_obj, err := url.Parse(url_str)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(url_str)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Parse the HTML ----------------------------------------------------------
	html := string(bytes)
	par := html_make_parser(html)
	par.document = dom_make_Document()
	par.document.set_base_url(*url_obj)
	doc := par.Run()

	// Find the <html> element -------------------------------------------------
	html_elem := doc.filter_elem_children_by_local_name(dom_name_pair{html_namespace, "html"})[0]
	ua_stylesheet := init_default_css(html_elem.(html_HTMLElement))

	// Find the <head> element -------------------------------------------------
	head_elem := html_elem.filter_elem_children_by_local_name(dom_name_pair{html_namespace, "head"})[0]

	// Apply style rules -------------------------------------------------------
	css_apply_style_rules(&ua_stylesheet, doc)
	log.Println("Style rules applied")

	// Do something with it ----------------------------------------------------
	_ = head_elem
	dom_print_tree(doc)
	log.Println("Document loaded. Making layout tree...")
	viewport_size := viewport_img.Rect.Size()
	for y := range viewport_size.Y {
		for x := range viewport_size.X {
			viewport_img.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	icb := make_layout(html_elem, float64(viewport_size.X), float64(viewport_size.Y), plat)
	browser_print_layout_tree(icb)
	log.Println("Made layout. Making paint tree...")
	paint := icb.make_paint_node()
	gfx.PrintPaintTree(paint)
	log.Println("Painting...")
	paint.Paint(viewport_img)

	log.Println("DONE")
}
