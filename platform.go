package main

type LinuxPlatform struct {
}

func (plat *LinuxPlatform) MeasureText(text string) (width, height float64) {
	// STUB
	return 10, 10
}
