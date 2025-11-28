package namespaces

import "fmt"

type Namespace string

const (
	Html   = Namespace("http://www.w3.org/1999/xhtml")
	Mathml = Namespace("http://www.w3.org/1998/Math/MathML")
	Svg    = Namespace("http://www.w3.org/2000/svg")
	Xlink  = Namespace("http://www.w3.org/1999/xlink")
	Xml    = Namespace("http://www.w3.org/XML/1998/namespace")
	Xmlns  = Namespace("http://www.w3.org/2000/xmlns/")
)

func HtmlP() *Namespace   { v := Html; return &v }
func MathmlP() *Namespace { v := Mathml; return &v }
func SvgP() *Namespace    { v := Svg; return &v }
func XlinkP() *Namespace  { v := Xlink; return &v }
func XmlP() *Namespace    { v := Xml; return &v }
func XmlnsP() *Namespace  { v := Xmlns; return &v }

func (n Namespace) String() string {
	switch n {
	case Html:
		return "html"
	case Mathml:
		return "mathml"
	case Svg:
		return "svg"
	case Xlink:
		return "xlink"
	case Xml:
		return "xml"
	case Xmlns:
		return "xmlns"
	default:
		return fmt.Sprintf("<namespace %s>", string(n))
	}
}
