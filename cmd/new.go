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
	exitIf(len(os.Args) < 2, "new: No subcommand specified")
	subcommand := os.Args[1]
	switch subcommand {
	case "sh":
		exitIf(len(os.Args) < 3, "new sh: No file name specified")
		n := os.Args[1]
		p := path.Join(bin, os.Args[2])
		s, _ := os.Stat(p)
		exitIf(s != nil, "new sh: %s: File exists", n)
		err := ioutil.WriteFile(p, []byte("#!/bin/bash\n\n"), 0744)
		exitIfError(err, "new sh: %s: Failed to create an shell script executable", n)
		err = exec.Command("open", "-a", visualStudioCode, p).Run()
		exitIfError(err, "new sh: %s: Failed to open with %s", n, visualStudioCode)
	default:
		exitIf(true, "new: %s: No such subcommand", subcommand)
	}
}

func exitIf(b bool, format string, a ...interface{}) {
	if b {
		fmt.Printf(format+"\n", a...)
		os.Exit(1)
	}
}

func exitIfError(err error, format string, a ...interface{}) {
	exitIf(err != nil, format, a...)
}
