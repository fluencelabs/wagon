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

	ExportFunctionName := flag.String("func-name", "main", "calling function name")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	WasmFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer WasmFile.Close()

	WasmModule, err := wasm.ReadModule(WasmFile, nil)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	err = validate.VerifyModule(WasmModule)
	if err != nil {
		log.Fatalf("could not verify module: %v", err)
	}

	if WasmModule.Export == nil {
		log.Fatalf("module has no export section")
	}

	WasmFunctionId, res := WasmModule.Export.Entries[*ExportFunctionName]
	if !res {
		log.Fatalf("export function not found")
	}

	WasmVM, err := exec.NewVM(WasmModule)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	ret, err := WasmVM.ExecCode(int64(WasmFunctionId.Index))
	if err != nil {
		log.Fatalf("error while function executing: %v", err)
	}

	fmt.Fprintf(os.Stdout, "%[1]v (%[1]T)\n", ret)
}
