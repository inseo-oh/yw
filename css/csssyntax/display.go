package csssyntax

import (
	"yw/css/display"
	"yw/util"
)

// https://www.w3.org/TR/css-display-3/#typedef-display-outside
func (ts *tokenStream) parseDisplayOutside() (display.OuterMode, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("block")) {
		return display.Block, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("inline")) {
		return display.Inline, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("run-in")) {
		return display.RunIn, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-display-3/#typedef-display-inside
func (ts *tokenStream) parseDisplayInside() (display.InnerMode, bool) {
	if !util.IsNil(ts.consumeIdentTokenWith("flow")) {
		return display.Flow, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("flow-root")) {
		return display.FlowRoot, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table")) {
		return display.Table, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("flex")) {
		return display.Flex, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("grid")) {
		return display.Grid, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("ruby")) {
		return display.Ruby, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-display-3/#propdef-display
func (ts *tokenStream) parseDisplay() (display.Display, bool) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if !util.IsNil(ts.consumeIdentTokenWith("inline-block")) {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.FlowRoot}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("inline-table")) {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Table}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("inline-flex")) {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Flex}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("inline-grid")) {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Grid}, true
	}
	// Try <display-outside> <display-inside> ------------------------------
	gotOuterMode, gotInnerMode := false, false
	var outerMode display.OuterMode
	var innerMode display.InnerMode
	for !gotOuterMode || !gotInnerMode {
		gotSomething := false
		if !gotOuterMode {
			outerMode, gotOuterMode = ts.parseDisplayOutside()
			if gotOuterMode {
				gotSomething = true
			}
		}
		if !gotInnerMode {
			innerMode, gotInnerMode = ts.parseDisplayInside()
			if gotInnerMode {
				gotSomething = true
			}
		}
		if !gotSomething {
			break
		}
	}
	if gotOuterMode || gotInnerMode {
		if !gotInnerMode {
			innerMode = display.Flow
		} else if !gotOuterMode {
			if innerMode == display.Ruby {
				outerMode = display.Inline
			} else {
				outerMode = display.Block
			}
		}
		return display.Display{Mode: display.OuterInnerMode, OuterMode: outerMode, InnerMode: innerMode}, true
	}
	// Try display-listitem ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-listitem
	// TODO

	// Try display-internal ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-internal

	if !util.IsNil(ts.consumeIdentTokenWith("table-row-group")) {
		return display.Display{Mode: display.TableRowGroup}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-header-group")) {
		return display.Display{Mode: display.TableHeaderGroup}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-footer-group")) {
		return display.Display{Mode: display.TableFooterGroup}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-row")) {
		return display.Display{Mode: display.TableRow}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-cell")) {
		return display.Display{Mode: display.TableCell}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-column-group")) {
		return display.Display{Mode: display.TableColumnGroup}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-column")) {
		return display.Display{Mode: display.TableColumn}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("table-caption")) {
		return display.Display{Mode: display.TableCaption}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("ruby-base")) {
		return display.Display{Mode: display.RubyBase}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("ruby-text")) {
		return display.Display{Mode: display.RubyText}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("ruby-base-container")) {
		return display.Display{Mode: display.RubyBaseContainer}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("ruby-text-container")) {
		return display.Display{Mode: display.RubyTextContainer}, true
	}

	// Try display-box -----------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-box

	if !util.IsNil(ts.consumeIdentTokenWith("contents")) {
		return display.Display{Mode: display.Contents}, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("none")) {
		return display.Display{Mode: display.DisplayNone}, true
	}

	return display.Display{}, false
}

func (ts *tokenStream) parseVisibility() (display.Visibility, bool) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if !util.IsNil(ts.consumeIdentTokenWith("visible")) {
		return display.Visible, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("hidden")) {
		return display.Hidden, true
	} else if !util.IsNil(ts.consumeIdentTokenWith("collapse")) {
		return display.Collapse, true
	}
	return 0, false
}
