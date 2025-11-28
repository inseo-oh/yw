package dom

import (
	"fmt"
	"strings"
)

type DocumentType interface {
	Node
	Name() string
	PublicId() string
	SystemId() string
}
type documentTypeImpl struct {
	Node
	name     string
	publicId string
	systemId string
}

func NewDocumentType(doc Document, name, publicId, systemId string) DocumentType {
	return documentTypeImpl{
		NewNode(doc),
		name, publicId, systemId,
	}
}
func (dt documentTypeImpl) String() string {
	sb := strings.Builder{}
	sb.WriteString("<!DOCTYPE")
	if dt.name != "" {
		sb.WriteString(fmt.Sprintf(" %s", dt.name))
	}
	if dt.publicId != "" {
		sb.WriteString(fmt.Sprintf(" PUBLIC %s", dt.publicId))
	}
	if dt.systemId != "" {
		sb.WriteString(fmt.Sprintf(" SYSTEM %s", dt.systemId))
	}
	sb.WriteString(">")
	return sb.String()
}
func (dt documentTypeImpl) Name() string     { return dt.name }
func (dt documentTypeImpl) PublicId() string { return dt.publicId }
func (dt documentTypeImpl) SystemId() string { return dt.systemId }
