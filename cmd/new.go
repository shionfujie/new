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
		fmt.Println("new sh: No file name specified")
		os.Exit(1)
	}
	n := os.Args[1]
	p := path.Join(bin, os.Args[1])
	if err := ioutil.WriteFile(p, []byte("#!/bin/bash\n\n"), 0744); err != nil {
		fmt.Printf("new sh: %s: Failed to create an shell script executable", n)
		os.Exit(1)
	}
	if err := exec.Command("open", "-a", visualStudioCode, p).Run(); err != nil {
		fmt.Printf("new sh: %s: Failed to open with %s", n, visualStudioCode)
		os.Exit(1)
	}
}
