package main

import (
	"fmt"
	"io"
	"log"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const bin = "/Users/shion.t.fujie/Desktop/room of machinery/bin"
const visualStudioCode = "Visual Studio Code"

func main() {
	logger := New(os.Stdout, "new: ", 0)

	logger.FatalfIf(len(os.Args) < 2, "No subcommand specified")
	subcommand := os.Args[1]
	switch subcommand {
	case "sh":
		logger.SetPrefix("new sh: ")
		logger.FatalfIf(len(os.Args) < 3, "No file name specified")
		n := os.Args[2]
		p := path.Join(bin, n)
		s, _ := os.Stat(p)
		logger.FatalfIf(s != nil, "%s: File exists", n)
		err := ioutil.WriteFile(p, []byte("#!/bin/bash\n\n"), 0744)
		logger.FatalfIfError(err, "%s: Failed to create an shell script executable", n)
		fmt.Fprintf(logger.O, "At %s\n", p) // Prints without the predefined format
		err = exec.Command("open", "-a", visualStudioCode, p).Run()
		logger.FatalfIfError(err, " %s: Failed to open with %s", n, visualStudioCode)
	default:
		logger.Fatalf("%s: No such subcommand", subcommand)
	}
}

type sLogger struct {
	*log.Logger
	O io.Writer
}

func New(out io.Writer, prefix string, flag int) *sLogger {
	return &sLogger{log.New(out, prefix, flag), out}
}

func (l *sLogger) FatalfIf(b bool, format string, a ...interface{}) {
	if b {
		l.Fatalf(format+"\n", a...)
	}
}

func (l *sLogger) FatalfIfError(err error, format string, a ...interface{}) {
	l.FatalfIf(err != nil, format, a...)
}
