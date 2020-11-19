package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const bin = "/Users/shion.t.fujie/Desktop/room of machinery/bin"
const visualStudioCode = "Visual Studio Code"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("fatal: No file name specified")
		os.Exit(1)
	}
	fn := path.Join(bin, os.Args[1])
	if err := ioutil.WriteFile(fn, []byte("#!/bin/bash\n\n"), 0744); err != nil {
		fmt.Printf("new sh: %s: Failed to initialize an shell script executable", fn)
		os.Exit(1)
	}
	if err := exec.Command("open", "-a", visualStudioCode, fn).Run(); err != nil {
		fmt.Printf("new sh: %s: Failed to open with %s", fn, visualStudioCode)
	}
}
