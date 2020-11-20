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
	require(len(os.Args) > 1, "new sh: No file name specified")

	n := os.Args[1]
	p := path.Join(bin, os.Args[1])
	if _, err := os.Stat(p); err == nil {
		fmt.Printf("new sh: %s: File exists\n", n)
		os.Exit(1)
	}
	err := ioutil.WriteFile(p, []byte("#!/bin/bash\n\n"), 0744)
	requireNoError(err, "new sh: %s: Failed to create an shell script executable", n)
	err = exec.Command("open", "-a", visualStudioCode, p).Run()
	requireNoError(err, "new sh: %s: Failed to open with %s", n, visualStudioCode)
}

func require(b bool, format string, a ...interface{}) {
	if !b {
		fmt.Printf(format + "\n", a...)
		os.Exit(1)
	}
}

func requireNoError(err error, format string, a ...interface{}) {
	require(err == nil, format, a...)
}