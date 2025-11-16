package platform

type Platform interface {
	//==========================================================================
	// Fonts
	//==========================================================================
	MeasureText(text string) (width, height float64)
}
