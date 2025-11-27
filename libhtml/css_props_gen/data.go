package main

// ==============================================================================
// List of CSS value types
// ==============================================================================
var (
	type_color            = css_type{"css_color", "parse_color"}
	type_size_or_auto     = css_type{"css_size_value", "parse_size_value_or_auto"}
	type_size_or_none     = css_type{"css_size_value", "parse_size_value_or_none"}
	type_display          = css_type{"css_display", "parse_display"}
	type_visibility       = css_type{"css_visibility", "parse_visibility"}
	type_line_style       = css_type{"css_line_style", "parse_line_style"}
	type_line_width       = css_type{"css_length", "parse_line_width"}
	type_margin           = css_type{"css_margin", "parse_margin"}
	type_padding          = css_type{"css_length_resolvable", "parse_padding"}
	type_font_family_list = css_type{"css_font_family_list", "parse_font_family"}
	type_font_weight      = css_type{"css_font_weight", "parse_font_weight"}
	type_font_stretch     = css_type{"css_font_stretch", "parse_font_stretch"}
	type_font_style       = css_type{"css_font_style", "parse_font_style"}
	type_font_size        = css_type{"css_font_size", "parse_font_size"}
)

// ==============================================================================
// List of CSS properties
// ==============================================================================
var (
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#border-color
	prop_border_top_color    = css_prop_simple{"border-top-color", type_color, "css_color_currentColor()", false}
	prop_border_right_color  = css_prop_simple{"border-right-color", type_color, "css_color_currentColor()", false}
	prop_border_bottom_color = css_prop_simple{"border-bottom-color", type_color, "css_color_currentColor()", false}
	prop_border_left_color   = css_prop_simple{"border-left-color", type_color, "css_color_currentColor()", false}
	// https://www.w3.org/TR/css-backgrounds-3/#border-style
	prop_border_top_style    = css_prop_simple{"border-top-style", type_line_style, "css_line_style_none", false}
	prop_border_right_style  = css_prop_simple{"border-right-style", type_line_style, "css_line_style_none", false}
	prop_border_bottom_style = css_prop_simple{"border-bottom-style", type_line_style, "css_line_style_none", false}
	prop_border_left_style   = css_prop_simple{"border-left-style", type_line_style, "css_line_style_none", false}
	// https://www.w3.org/TR/css-backgrounds-3/#border-width
	prop_border_top_width    = css_prop_simple{"border-top-width", type_line_width, "css_line_width_medium()", false}
	prop_border_right_width  = css_prop_simple{"border-right-width", type_line_width, "css_line_width_medium()", false}
	prop_border_bottom_width = css_prop_simple{"border-bottom-width", type_line_width, "css_line_width_medium()", false}
	prop_border_left_width   = css_prop_simple{"border-left-width", type_line_width, "css_line_width_medium()", false}
	// https://www.w3.org/TR/css-backgrounds-3/#border-shorthands
	prop_border_color = css_prop_shorthand_sides{"border-color", prop_border_top_color, prop_border_right_color, prop_border_bottom_color, prop_border_left_color, false}
	prop_border_style = css_prop_shorthand_sides{"border-style", prop_border_top_style, prop_border_right_style, prop_border_bottom_style, prop_border_left_style, false}
	prop_border_width = css_prop_shorthand_sides{"border-width", prop_border_top_width, prop_border_right_width, prop_border_bottom_width, prop_border_left_width, false}
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/#margin-physical
	prop_margin_top    = css_prop_simple{"margin-top", type_margin, "css_margin{css_length{css_number_from_int(0), css_length_unit_px}}", false}
	prop_margin_right  = css_prop_simple{"margin-right", type_margin, "css_margin{css_length{css_number_from_int(0), css_length_unit_px}}", false}
	prop_margin_bottom = css_prop_simple{"margin-bottom", type_margin, "css_margin{css_length{css_number_from_int(0), css_length_unit_px}}", false}
	prop_margin_left   = css_prop_simple{"margin-left", type_margin, "css_margin{css_length{css_number_from_int(0), css_length_unit_px}}", false}
	// https://www.w3.org/TR/css-box-3/#padding-physical
	prop_padding_top    = css_prop_simple{"padding-top", type_padding, "css_length{css_number_from_int(0), css_length_unit_px}", false}
	prop_padding_right  = css_prop_simple{"padding-right", type_padding, "css_length{css_number_from_int(0), css_length_unit_px}", false}
	prop_padding_bottom = css_prop_simple{"padding-bottom", type_padding, "css_length{css_number_from_int(0), css_length_unit_px}", false}
	prop_padding_left   = css_prop_simple{"padding-left", type_padding, "css_length{css_number_from_int(0), css_length_unit_px}", false}
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-family-prop
	prop_font_family = css_prop_simple{"font-family", type_font_family_list, "css_font_family_list{[]css_font_family{{css_font_family_type_sans_serif, \"\"}}}", true}
	// https://www.w3.org/TR/css-fonts-3/#font-weight-prop
	prop_font_weight = css_prop_simple{"font-weight", type_font_weight, "css_font_weight_normal", true}
	// https://www.w3.org/TR/css-fonts-3/#font-stretch-prop
	prop_font_stretch = css_prop_simple{"font-stretch", type_font_stretch, "css_font_stretch_normal", true}
	// https://www.w3.org/TR/css-fonts-3/#font-style-prop
	prop_font_style = css_prop_simple{"font-style", type_font_style, "css_font_style_normal", true}
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
	prop_font_size = css_prop_simple{"font-size", type_font_size, "css_absolute_size_medium", true}
)
var props = []css_prop{
	//==========================================================================
	// https://www.w3.org/TR/css-color-4/
	//==========================================================================
	// https://www.w3.org/TR/css-color-4/#the-color-property
	css_prop_simple{"color", type_color, "css_color_CanvasText()", true},
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/
	//==========================================================================
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#preferred-size-properties
	css_prop_simple{"width", type_size_or_auto, "css_size_value_auto()", false},
	css_prop_simple{"height", type_size_or_auto, "css_size_value_auto()", false},
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#min-size-properties
	css_prop_simple{"min-width", type_size_or_auto, "css_size_value_auto()", false},
	css_prop_simple{"min-height", type_size_or_auto, "css_size_value_auto()", false},
	// https://www.w3.org/TR/2021/WD-css-sizing-3-20211217/#max-size-properties
	css_prop_simple{"max-width", type_size_or_none, "css_size_value_none()", false},
	css_prop_simple{"max-height", type_size_or_none, "css_size_value_none()", false},
	//==========================================================================
	// https://www.w3.org/TR/css-display-3/
	//==========================================================================
	// https://www.w3.org/TR/css-display-3/#the-display-properties
	css_prop_simple{"display", type_display, "css_display{css_display_mode_outer_inner_mode, css_display_outer_mode_inline, css_display_inner_mode_flow}", false},
	// https://www.w3.org/TR/css-display-3/#visibility
	css_prop_simple{"visibility", type_visibility, "css_visibility_visible", true},
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/
	//==========================================================================
	// https://www.w3.org/TR/css-backgrounds-3/#background-color
	css_prop_simple{"background-color", type_color, "css_color_transparent()", false},
	// https://www.w3.org/TR/css-backgrounds-3/#border-color
	prop_border_top_color, prop_border_right_color, prop_border_bottom_color, prop_border_left_color, prop_border_color,
	// https://www.w3.org/TR/css-backgrounds-3/#border-style
	prop_border_top_style, prop_border_right_style, prop_border_bottom_style, prop_border_left_style, prop_border_style,
	// https://www.w3.org/TR/css-backgrounds-3/#border-width
	prop_border_top_width, prop_border_right_width, prop_border_bottom_width, prop_border_left_width, prop_border_width,
	// https://www.w3.org/TR/css-backgrounds-3/#border-shorthands
	css_prop_shorthand_any{"border-top", []css_prop{prop_border_top_width, prop_border_top_style, prop_border_top_color}, false},
	css_prop_shorthand_any{"border-right", []css_prop{prop_border_right_width, prop_border_right_style, prop_border_right_color}, false},
	css_prop_shorthand_any{"border-bottom", []css_prop{prop_border_bottom_width, prop_border_bottom_style, prop_border_bottom_color}, false},
	css_prop_shorthand_any{"border-left", []css_prop{prop_border_left_width, prop_border_left_style, prop_border_left_color}, false},
	css_prop_shorthand_any{"border", []css_prop{prop_border_width, prop_border_style, prop_border_color}, false},
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/
	//==========================================================================
	// https://www.w3.org/TR/css-box-3/#margin-physical
	prop_margin_top, prop_margin_right, prop_margin_bottom, prop_margin_left,
	css_prop_shorthand_sides{"margin", prop_margin_top, prop_margin_right, prop_margin_bottom, prop_margin_left, false},
	// https://www.w3.org/TR/css-box-3/#padding-physical
	prop_padding_top, prop_padding_right, prop_padding_bottom, prop_padding_left,
	css_prop_shorthand_sides{"padding", prop_padding_top, prop_padding_right, prop_padding_bottom, prop_padding_left, false},
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/
	//==========================================================================
	// https://www.w3.org/TR/css-fonts-3/#font-family-prop
	prop_font_family,
	// https://www.w3.org/TR/css-fonts-3/#font-weight-prop
	prop_font_weight,
	// https://www.w3.org/TR/css-fonts-3/#font-stretch-prop
	prop_font_stretch,
	// https://www.w3.org/TR/css-fonts-3/#font-style-prop
	prop_font_style,
	// https://www.w3.org/TR/css-fonts-3/#propdef-font-size
	prop_font_size,
	// https://www.w3.org/TR/css-fonts-3/#font-prop
	css_prop_shorthand_any{"font", []css_prop{prop_font_family, prop_font_weight, prop_font_stretch, prop_font_style, prop_font_size}, true},
}
