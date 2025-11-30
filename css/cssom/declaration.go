package cssom

import (
	"log"
	"yw/css/props"
	"yw/dom"
)

type Declaration struct {
	Name        string
	Value       props.PropertyValue
	IsImportant bool
}

func (d Declaration) ApplyStyleRules(elem dom.Element) {
	desc := props.DescriptorsMap[d.Name]
	if desc.ApplyFunc == nil {
		log.Printf("TODO: CSS Property %s is recognized but not supported yet. (Missing applyFunc() function)", d.Name)
		return
	}
	desc.ApplyFunc(&ElementDataOf(elem).ComputedStyleSet, d.Value)
}
