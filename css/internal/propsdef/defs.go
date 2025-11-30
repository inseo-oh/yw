package propsdef

// ==============================================================================
// List of CSS value types
// ==============================================================================
var (
	typeColor          = CssType{"csscolor.Color", "parseColor"}
	typeSizeOrAuto     = CssType{"sizing.Size", "parseSizeOrAuto"}
	typeSizeOrNone     = CssType{"sizing.Size", "parseSizeOrNone"}
	typeDisplay        = CssType{"display.Display", "parseDisplay"}
	typeVisibility     = CssType{"display.Visibility", "parseVisibility"}
	typeLineStyle      = CssType{"backgrounds.LineStyle", "parseLineStyle"}
	typeLineWidth      = CssType{"values.Length", "parseLineWidth"}
	typeMargin         = CssType{"box.Margin", "parseMargin"}
	typePadding        = CssType{"values.LengthResolvable", "parsePadding"}
	typeFontFamilyList = CssType{"fonts.FamilyList", "parseFontFamily"}
	typeFontWeight     = CssType{"fonts.Weight", "parseFontWeight"}
	typeFontStretch    = CssType{"fonts.Stretch", "parseFontStretch"}
	typeFontStyle      = CssType{"fonts.Style", "parseFontStyle"}
	typeFontSize       = CssType{"fonts.Size", "parseFontSize"}
	typeTextTransform  = CssType{"text.Transform", "parseTextTransform"}
)

// ==============================================================================
// List of CSS properties
// ==============================================================================
var (
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#border-color
	propBorderTopColor    = SimpleProp{"border-top-color", typeColor, "csscolor.NewCurrentColor()", false}
	propBorderRightColor  = SimpleProp{"border-right-color", typeColor, "csscolor.NewCurrentColor()", false}
	propBorderBottomColor = SimpleProp{"border-bottom-color", typeColor, "csscolor.NewCurrentColor()", false}
	propBorderLeftColor   = SimpleProp{"border-left-color", typeColor, "csscolor.NewCurrentColor()", false}
	// https://www.w3.org/TR/css-backgrounds-3/#border-style
	propBorderTopStyle    = SimpleProp{"border-top-style", typeLineStyle, "backgrounds.NoLine", false}
	propBorderRightStyle  = SimpleProp{"border-right-style", typeLineStyle, "backgrounds.NoLine", false}
	propBorderBottomStyle = SimpleProp{"border-bottom-style", typeLineStyle, "backgrounds.NoLine", false}
	propBorderLeftStyle   = SimpleProp{"border-left-style", typeLineStyle, "backgrounds.NoLine", false}
	// https://www.w3.org/TR/css-backgrounds-3/#border-width
	propBorderTopWidth    = SimpleProp{"border-top-width", typeLineWidth, "backgrounds.LineWidthMedium()", false}
	propBorderRightWidth  = SimpleProp{"border-right-width", typeLineWidth, "backgrounds.LineWidthMedium()", false}
	propBorderBottomWidth = SimpleProp{"border-bottom-width", typeLineWidth, "backgrounds.LineWidthMedium()", false}
	propBorderLeftWidth   = SimpleProp{"border-left-width", typeLineWidth, "backgrounds.LineWidthMedium()", false}
	// https://www.w3.org/TR/css-backgrounds-3/#border-shorthands
	propBorderColor = ShorthandSidesProp{"border-color", propBorderTopColor, propBorderRightColor, propBorderBottomColor, propBorderLeftColor, false}
	propBorderStyle = ShorthandSidesProp{"border-style", propBorderTopStyle, propBorderRightStyle, propBorderBottomStyle, propBorderLeftStyle, false}
	propBorderWidth = ShorthandSidesProp{"border-width", propBorderTopWidth, propBorderRightWidth, propBorderBottomWidth, propBorderLeftWidth, false}
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/#margin-physical
	propMarginTop    = SimpleProp{"margin-top", typeMargin, "box.Margin{Value: values.LengthFromPx(css.NumFromInt(0))}", false}
	propMarginRight  = SimpleProp{"margin-right", typeMargin, "box.Margin{Value: values.LengthFromPx(css.NumFromInt(0))}", false}
	propMarginBottom = SimpleProp{"margin-bottom", typeMargin, "box.Margin{Value: values.LengthFromPx(css.NumFromInt(0))}", false}
	propMarginLeft   = SimpleProp{"margin-left", typeMargin, "box.Margin{Value: values.LengthFromPx(css.NumFromInt(0))}", false}
	// https://www.w3.org/TR/css-box-3/#padding-physical
	propPaddingTop    = SimpleProp{"padding-top", typePadding, "values.LengthFromPx(css.NumFromInt(0))", false}
	propPaddingRight  = SimpleProp{"padding-right", typePadding, "values.LengthFromPx(css.NumFromInt(0))", false}
	propPaddingBottom = SimpleProp{"padding-bottom", typePadding, "values.LengthFromPx(css.NumFromInt(0))", false}
	propPaddingLeft   = SimpleProp{"padding-left", typePadding, "values.LengthFromPx(css.NumFromInt(0))", false}
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-family-prop
	propFontFamily = SimpleProp{"font-family", typeFontFamilyList, "fonts.FamilyList{Families: []fonts.Family{{Type: fonts.SansSerif}}}", true}
	// https://www.w3.org/TR/css-fonts-3/#font-weight-prop
	propFontWeight = SimpleProp{"font-weight", typeFontWeight, "fonts.NormalWeight", true}
	// https://www.w3.org/TR/css-fonts-3/#font-stretch-prop
	propFontStretch = SimpleProp{"font-stretch", typeFontStretch, "fonts.NormalStretch", true}
	// https://www.w3.org/TR/css-fonts-3/#font-style-prop
	propFontStyle = SimpleProp{"font-style", typeFontStyle, "fonts.NormalStyle", true}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
	propFontSize = SimpleProp{"font-size", typeFontSize, "fonts.MediumSize", true}
)
var Props = []CssProp{
	//==========================================================================
	// https://www.w3.org/TR/css-color-4/
	//==========================================================================
	// https://www.w3.org/TR/css-color-4/#the-color-property
	SimpleProp{"color", typeColor, "csscolor.CanvasText()", true},
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#preferred-size-properties
	SimpleProp{"width", typeSizeOrAuto, "sizing.Size{Type: sizing.Auto}", false},
	SimpleProp{"height", typeSizeOrAuto, "sizing.Size{Type: sizing.Auto}", false},
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#min-size-properties
	SimpleProp{"min-width", typeSizeOrAuto, "sizing.Size{Type: sizing.Auto}", false},
	SimpleProp{"min-height", typeSizeOrAuto, "sizing.Size{Type: sizing.Auto}", false},
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#max-size-properties
	SimpleProp{"max-width", typeSizeOrNone, "sizing.Size{Type: sizing.NoneSize}", false},
	SimpleProp{"max-height", typeSizeOrNone, "sizing.Size{Type: sizing.NoneSize}", false},
	//==========================================================================
	// https://www.w3.org/TR/css-display-3/
	//==========================================================================
	// https://www.w3.org/TR/css-display-3/#the-display-properties
	SimpleProp{"display", typeDisplay, "display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Flow}", false},
	// https://www.w3.org/TR/css-display-3/#visibility
	SimpleProp{"visibility", typeVisibility, "display.Visible", true},
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#background-color
	SimpleProp{"background-color", typeColor, "csscolor.Transparent()", false},
	// https://www.w3.org/TR/css-backgrounds-3/#border-color
	propBorderTopColor, propBorderRightColor, propBorderBottomColor, propBorderLeftColor, propBorderColor,
	// https://www.w3.org/TR/css-backgrounds-3/#border-style
	propBorderTopStyle, propBorderRightStyle, propBorderBottomStyle, propBorderLeftStyle, propBorderStyle,
	// https://www.w3.org/TR/css-backgrounds-3/#border-width
	propBorderTopWidth, propBorderRightWidth, propBorderBottomWidth, propBorderLeftWidth, propBorderWidth,
	// https://www.w3.org/TR/css-backgrounds-3/#border-shorthands
	ShorthandAnyProp{"border-top", []CssProp{propBorderTopWidth, propBorderTopStyle, propBorderTopColor}, false},
	ShorthandAnyProp{"border-right", []CssProp{propBorderRightWidth, propBorderRightStyle, propBorderRightColor}, false},
	ShorthandAnyProp{"border-bottom", []CssProp{propBorderBottomWidth, propBorderBottomStyle, propBorderBottomColor}, false},
	ShorthandAnyProp{"border-left", []CssProp{propBorderLeftWidth, propBorderLeftStyle, propBorderLeftColor}, false},
	ShorthandAnyProp{"border", []CssProp{propBorderWidth, propBorderStyle, propBorderColor}, false},
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/#margin-physical
	propMarginTop, propMarginRight, propMarginBottom, propMarginLeft,
	ShorthandSidesProp{"margin", propMarginTop, propMarginRight, propMarginBottom, propMarginLeft, false},
	// https://www.w3.org/TR/css-box-3/#padding-physical
	propPaddingTop, propPaddingRight, propPaddingBottom, propPaddingLeft,
	ShorthandSidesProp{"padding", propPaddingTop, propPaddingRight, propPaddingBottom, propPaddingLeft, false},
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-family-prop
	propFontFamily,
	// https://www.w3.org/TR/css-fonts-3/#font-weight-prop
	propFontWeight,
	// https://www.w3.org/TR/css-fonts-3/#font-stretch-prop
	propFontStretch,
	// https://www.w3.org/TR/css-fonts-3/#font-style-prop
	propFontStyle,
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
	propFontSize,
	// https://www.w3.org/TR/css-fonts-3/#font-prop
	ShorthandAnyProp{"font", []CssProp{propFontFamily, propFontWeight, propFontStretch, propFontStyle, propFontSize}, true},
	//==========================================================================
	// https://www.w3.org/TR/css-text-3/
	//==========================================================================
	// https://www.w3.org/TR/css-text-3/#text-transform-property
	SimpleProp{"text-transform", typeTextTransform, "text.Transform{Type: text.NoTransform}", true},
}
