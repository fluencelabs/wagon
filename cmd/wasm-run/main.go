// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/validate"
	"github.com/go-interpreter/wagon/wasm"
)

func main() {
	log.SetPrefix("wasm-run: ")
	log.SetFlags(0)

	wasmFuncName := flag.String("func-name", "main", "called function name")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	wasmFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer wasmFile.Close()

	wasmModule, err := wasm.ReadModule(wasmFile, nil)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	err = validate.VerifyModule(wasmModule)
	if err != nil {
		log.Fatalf("could not verify module: %v", err)
	}

	if wasmModule.Export == nil {
		log.Fatalf("module has no export section")
	}

	wasmFuncId, res := wasmModule.Export.Entries[*wasmFuncName]
	if !res {
		log.Fatalf("could not find export function")
	}

	wasmVM, err := exec.NewVM(wasmModule)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	ret, err := wasmVM.ExecCode(int64(wasmFuncId.Index))
	if err != nil {
		log.Fatalf("could not execute requested function: %v", err)
	}

	fmt.Fprintf(os.Stdout, "%[1]v (%[1]T)\n", ret)
}
