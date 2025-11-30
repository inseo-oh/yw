package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"yw/browser"
)

var url = flag.String("url", "", "The URL")

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	br := browser.Browser{}
	viewportImg := image.NewRGBA(image.Rect(0, 0, 1280, 720))
	plat := initPlatform()
	br.Run(*url, plat, viewportImg)

	destFile, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(destFile, viewportImg)
	if err != nil {
		log.Fatal(err)
	}
	destFile.Close()
}
