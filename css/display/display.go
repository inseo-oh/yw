// Implementation of the CSS Display Module Level 3 (https://www.w3.org/TR/css-display-3/)
package display

import (
	"fmt"
)

// https://www.w3.org/TR/css-display-3/#outer-role
type OuterMode uint8

const (
	Block  = OuterMode(iota) // display: block
	Inline                   // display: inline
	RunIn                    // display: run-in
)

func (m OuterMode) String() string {
	switch m {
	case Block:
		return "block"
	case Inline:
		return "inline"
	case RunIn:
		return "run-in"
	}
	return fmt.Sprintf("<bad OuterMode %d>", m)
}

// https://www.w3.org/TR/css-display-3/#inner-model
type InnerMode uint8

const (
	Flow     = InnerMode(iota) // display: flow
	FlowRoot                   // display: flow-root
	Table                      // display: table
	Flex                       // display: flex
	Grid                       // display: grid
	Ruby                       // display: ruby
)

func (m InnerMode) String() string {
	switch m {
	case Flow:
		return "flow"
	case FlowRoot:
		return "flow-root"
	case Table:
		return "table"
	case Flex:
		return "flex"
	case Grid:
		return "grid"
	case Ruby:
		return "ruby"
	}
	return fmt.Sprintf("<bad InnerMode %d>", m)
}

// https://www.w3.org/TR/css-display-3/#propdef-display
type Display struct {
	Mode Mode

	OuterMode OuterMode // Only valid when mode is OuterInnerMode
	InnerMode InnerMode // Only valid when mode is OuterInnerMode
}

func (d Display) String() string {
	switch d.Mode {
	case OuterInnerMode:
		return fmt.Sprintf("%v %v", d.OuterMode, d.InnerMode)
	case TableRowGroup:
		return "table-row-group"
	case TableHeaderGroup:
		return "table-header-group"
	case TableFooterGroup:
		return "table-footer-group"
	case TableRow:
		return "table-row"
	case TableCell:
		return "table-cell"
	case TableColumnGroup:
		return "table-column-group"
	case TableColumn:
		return "table-column"
	case TableCaption:
		return "table-caption"
	case RubyBase:
		return "ruby-base"
	case RubyText:
		return "ruby-text"
	case RubyBaseContainer:
		return "ruby-base-container"
	case RubyTextContainer:
		return "ruby-text-container"
	case Contents:
		return "contents"
	case DisplayNone:
		return "none"
	}
	return fmt.Sprintf("bad Display %d", d)
}

type Mode uint8

const (
	// Display mode determined by `OuterMode` and `InnerMode` field.
	OuterInnerMode = Mode(iota)

	// https://www.w3.org/TR/css-display-3/#layout-specific-display

	TableRowGroup     // table-row-group
	TableHeaderGroup  // table-header-group
	TableFooterGroup  // table-footer-group
	TableRow          // table-row
	TableCell         // table-cell
	TableColumnGroup  // table-column-group
	TableColumn       // table-column
	TableCaption      // table-caption
	RubyBase          // ruby-base
	RubyText          // ruby-text
	RubyBaseContainer // ruby-base-container
	RubyTextContainer // ruby-text-container

	// https://www.w3.org/TR/css-display-3/#box-generation

	Contents
	DisplayNone
)

// https://www.w3.org/TR/css-display-3/#propdef-visibility
type Visibility uint8

const (
	Visible  = Visibility(iota) // visibility: visible
	Hidden                      // visibility: hidden
	Collapse                    // visibility: collapse
)

func (m Visibility) String() string {
	switch m {
	case Visible:
		return "visible"
	case Hidden:
		return "hidden"
	case Collapse:
		return "collapse"
	}
	return fmt.Sprintf("bad Visibility %d", m)
}
