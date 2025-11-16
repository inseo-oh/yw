package main

import (
	"log"
	"yw/libes"
)

func main() {
	vm := libes.MakeVm()
	code, err := libes.Compile("(69);")
	if err != nil {
		log.Fatal(err)
	}
	res := vm.Exec(code)
	log.Println(res)
}
