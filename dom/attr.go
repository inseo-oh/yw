package dom

import (
	"fmt"

	"github.com/inseo-oh/yw/namespaces"
)

type Attr interface {
	Node
	LocalName() string
	Value() string
	Namespace() (namespaces.Namespace, bool)
	NamespacePrefix() (string, bool)
	Element() Element
}
type AttrData struct {
	Node
	LocalName       string
	Value           string
	Namespace       *namespaces.Namespace // may be nil
	NamespacePrefix *string               // may be nil
}
type attrImpl struct {
	Node
	localName       string
	value           string
	namespace       *namespaces.Namespace // may be nil
	namespacePrefix *string               // may be nil
	element         Element
}

// namespace, namespacePrefix may be nil
func NewAttr(localName string, value string, namespace *namespaces.Namespace, namespacePrefix *string, element Element) Attr {
	return &attrImpl{
		NewNode(element.NodeDocument()),
		localName, value, namespace, namespacePrefix, element,
	}
}

func (at attrImpl) String() string {
	if at.namespace != nil {
		return fmt.Sprintf("#attr(%v:%s = %s)", *at.namespace, at.localName, at.value)
	} else {
		return fmt.Sprintf("#attr(%s = %s)", at.localName, at.value)
	}
}
func (at attrImpl) LocalName() string { return at.localName }
func (at attrImpl) Value() string     { return at.value }
func (at attrImpl) Element() Element  { return at.element }
func (at attrImpl) Namespace() (namespaces.Namespace, bool) {
	if at.namespace == nil {
		return "", false
	}
	return *at.namespace, true
}
func (at attrImpl) NamespacePrefix() (string, bool) {
	if at.namespacePrefix == nil {
		return "", false
	}
	return *at.namespacePrefix, true
}
