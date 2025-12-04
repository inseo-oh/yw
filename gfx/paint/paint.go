// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

// Package paint provides paint tree and painting operations.
//
// Paint tree is generated during layout process(see layout package), and
// drawn to the destination image by Paint()ing the root node.
package paint

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"

	"github.com/inseo-oh/yw/gfx"
)

func fillRect(dest *image.RGBA, r image.Rectangle, color color.Color) {
	for y := r.Min.Y; y <= r.Max.Y; y++ {
		for x := r.Min.X; x <= r.Max.X; x++ {
			dest.Set(x, y, color)
		}
	}
}

// Node represents a node in the paint tree.
type Node interface {
	// Paint paints to the dest.
	Paint(dest *image.RGBA)

	// String returns decscription of the node.
	String() string
}

var (
	DashWidth = 10 // Width of dash
)

// Text is Node that paints a text.
type TextPaint struct {
	Left, Top float64
	Text      string
	Font      gfx.Font
	Size      float64
	Color     color.RGBA
	Decors    []gfx.TextDecorOptions
}

func (t TextPaint) Paint(dest *image.RGBA) {
	t.Font.SetTextSize(int(t.Size))
	metrics := t.Font.Metrics()
	x := t.Left
	baselineY := t.Top + metrics.Ascender
	text := t.Text
	textRect := t.Font.DrawText(text, dest, x, baselineY, t.Color)

	for _, decor := range t.Decors {
		yPos := 0.0
		thickness := metrics.UnderlineThickness
		switch decor.Type {
		case gfx.Underline:
			yPos = baselineY - metrics.UnderlinePosition - thickness/2
		case gfx.Overline:
			yPos = t.Top
		case gfx.ThroughText:
			yPos = textRect.Top
		}
		decorRect := image.Rect(int(textRect.Left), int(yPos), int(textRect.Left+textRect.Width)-1, int(yPos+thickness))
		switch decor.Style {
		case gfx.SolidLine:
			fillRect(dest, decorRect, decor.Color)
		case gfx.DoubleLine:
			fillRect(dest, decorRect, decor.Color)
			nextLineRect := decorRect.Add(image.Pt(0, int(thickness*2)))
			fillRect(dest, nextLineRect, decor.Color)
		case gfx.DottedLine, gfx.DashedLine:
			var width int
			if decor.Style == gfx.DashedLine {
				width = DashWidth
			} else {
				width = decorRect.Dy()
			}
			for x := decorRect.Min.X; x < decorRect.Max.X; x += width * 2 {
				dotRect := image.Rect(x, decorRect.Min.Y, x+width-1, decorRect.Max.Y)
				fillRect(dest, dotRect, decor.Color)
			}
		case gfx.WavyLine:
			sinInput := 0.0
			currYOff := 0.0
			for x := decorRect.Min.X; x < decorRect.Max.X; x++ {
				dotRect := image.Rect(x, int(float64(decorRect.Min.Y)+currYOff), x, int(float64(decorRect.Max.Y)+currYOff))
				fillRect(dest, dotRect, decor.Color)
				currYOff += math.Sin(sinInput)
				sinInput += 0.5
			}
		}
	}
}
func (t TextPaint) String() string {
	return fmt.Sprintf("text-paint(%s) %v %g", t.Text, t.Color, t.Size)
}

// Text is Node that paints a box.
type BoxPaint struct {
	Items []Node
	Rect  gfx.Rect
	Color color.RGBA
}

func (g BoxPaint) Paint(dest *image.RGBA) {
	colorImg := image.NewRGBA(image.Rect(0, 0, int(g.Rect.Width)-1, int(g.Rect.Height)-1))
	for yOff := range int(g.Rect.Height) {
		for xOff := range int(g.Rect.Width) {
			colorImg.Set(xOff, yOff, g.Color)
		}
	}
	draw.Draw(dest, image.Rect(int(g.Rect.Left), int(g.Rect.Top), int(g.Rect.Left+g.Rect.Width)-1, int(g.Rect.Top+g.Rect.Height)-1), colorImg, image.Point{0, 0}, draw.Over)
	for _, item := range g.Items {
		item.Paint(dest)
	}
}
func (t BoxPaint) String() string {
	return fmt.Sprintf("box-paint(color=%v, rect=%v, %d items)", t.Color, t.Rect, len(t.Items))
}

// PrintTree prints paint tree to the standard output.
func PrintTree(node Node) {
	var doPrint func(node Node, indentLevel int)
	doPrint = func(node Node, indentLevel int) {
		currNode := node
		indent := strings.Repeat(" ", indentLevel*4)
		fmt.Printf("%s%v\n", indent, node)
		if gpaint, ok := currNode.(BoxPaint); ok {
			for _, child := range gpaint.Items {
				doPrint(child, indentLevel+1)
			}
		}
	}
	doPrint(node, 0)

}
