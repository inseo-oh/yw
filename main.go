package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"yw/libhtml"
)

var url = flag.String("url", "", "The URL")

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	br := libhtml.Browser{}
	viewport_img := image.NewRGBA(image.Rect(0, 0, 1280, 720))
	plat := init_platform()
	br.Init(*url, plat, viewport_img)

	dest_file, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(dest_file, viewport_img)
	if err != nil {
		log.Fatal(err)
	}
	dest_file.Close()
}
