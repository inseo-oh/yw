// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package main

import (
	"log"

	"github.com/inseo-oh/yw/es/escompiler"
	"github.com/inseo-oh/yw/es/vm"
)

func main() {
	vm := vm.Vm{}
	code, err := escompiler.Compile("(69);")
	if err != nil {
		log.Fatal(err)
	}
	res := vm.Exec(code)
	log.Println(res)
}
