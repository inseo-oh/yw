package libhtml

import "fmt"

type namespace string

const (
	html_namespace   = namespace("http://www.w3.org/1999/xhtml")
	mathml_namespace = namespace("http://www.w3.org/1998/Math/MathML")
	svg_namespace    = namespace("http://www.w3.org/2000/svg")
	xlink_namespace  = namespace("http://www.w3.org/1999/xlink")
	xml_namespace    = namespace("http://www.w3.org/XML/1998/namespace")
	xmlns_namespace  = namespace("http://www.w3.org/2000/xmlns/")
)

func html_namespace_p() *namespace   { v := html_namespace; return &v }
func mathml_namespace_p() *namespace { v := mathml_namespace; return &v }
func svg_namespace_p() *namespace    { v := svg_namespace; return &v }
func xlink_namespace_p() *namespace  { v := xlink_namespace; return &v }
func xml_namespace_p() *namespace    { v := xml_namespace; return &v }
func xmlns_namespace_p() *namespace  { v := xmlns_namespace; return &v }

func (n namespace) String() string {
	switch n {
	case html_namespace:
		return "html"
	case mathml_namespace:
		return "mathml"
	case svg_namespace:
		return "svg"
	case xlink_namespace:
		return "xlink"
	case xml_namespace:
		return "xml"
	case xmlns_namespace:
		return "xmlns"
	default:
		return fmt.Sprintf("<namespace %s>", string(n))
	}
}
