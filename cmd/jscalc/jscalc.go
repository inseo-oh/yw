package main

import (
	"log"

	"github.com/inseo-oh/yw/es/escompiler"
	"github.com/inseo-oh/yw/es/vm"
)

func main() {
	vm := vm.MakeVm()
	code, err := escompiler.Compile("(69);")
	if err != nil {
		log.Fatal(err)
	}
	res := vm.Exec(code)
	log.Println(res)
}
