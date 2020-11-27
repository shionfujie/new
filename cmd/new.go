package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

const bin = "/Users/shion.t.fujie/Desktop/machinery/bin"
const visualStudioCode = "Visual Studio Code"

const goMainFileTemplate = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Happy coding, %s!!!")
}
`

const chromeThemeManifestTemplate = `{
	"version": "0.1.0",
	"name": "%s",
	"manifest_version": 2,
	"theme": {
		"colors": {
			"frame": [255, 255, 255],
			"toolbar": [255, 255, 255],
			"ntp_text" : [255, 255, 255],
			"ntp_link" : [117, 117, 117],
			"ntp_section" : [207, 0, 192],
			"button_background" : [255, 0, 255]
		},
		"tints": {
			"buttons": [0, 0, 0.46]
		}
	}
}
`

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
		ensureFileNotExists(logger, p)

		err := ioutil.WriteFile(p, []byte("#!/bin/bash\n\n"), 0744)
		logger.FatalfIfError(err, "%s: Failed to create an shell script executable", n)

		fmt.Fprintf(logger.O, "A shell script executable has been created at '%s'\n", p) // Prints without the predefined format

		err = exec.Command("open", "-a", visualStudioCode, p).Run()
		logger.FatalfIfError(err, " %s: Failed to open with %s", n, visualStudioCode)
	case "go-cmd":
		logger.SetPrefix("new go-cmd: ")
		logger.FatalfIf(len(os.Args) < 3, "No command name specified. Specify a command name.")

		n := os.Args[2]
		ensureFileNotExists(logger, n)

		os.Mkdir(n, 0744)
		os.Chdir(n)
		goPackage := "sfujie.io/cli/" + n
		err := exec.Command("go", "mod", "init", goPackage).Run()
		logger.FatalfIfError(err, "%s: Failed to create the GO package '%s'", n, goPackage)

		os.Mkdir("cmd", 0744)
		os.Chdir("cmd")
		code := fmt.Sprintf(goMainFileTemplate, n)
		err = ioutil.WriteFile(n+".go", []byte(code), 0744)
		logger.FatalfIfError(err, "%s: Failed to create a main file at '%s', ")

		fmt.Fprintln(logger.O, "Have created a command-line project. Strive to code!!!")

		exec.Command("open", "-a", visualStudioCode, "../../"+n).Run() // Try to open the project
		exec.Command("open", "-a", visualStudioCode, n+".go").Run()
	case "chrome-theme":
		logger.SetPrefix("new chrome-theme: ")
		logger.FatalfIf(len(os.Args) < 3, "Extension name argument expected")

		n := os.Args[2]
		ensureFileNotExists(logger, n)
		os.Mkdir(n, 0744)
		os.Chdir(n)

		json := fmt.Sprintf(chromeThemeManifestTemplate, n)
		err := ioutil.WriteFile("manifest.json", []byte(json), 0744)
		logger.FatalfIfError(err, "%s: Failed to create a manifest file for a chrome extension", n)

		logger.Println("Have created a chrome theme. Excited to decorate!!!")

		exec.Command("open", "-a", visualStudioCode, "../../"+n).Run() // Try to open the project
		exec.Command("open", "-a", visualStudioCode, "manifest.json").Run()
	default:
		logger.Fatalf("%s: No such subcommand", subcommand)
	}
}

type sLogger struct {
	*log.Logger
	O io.Writer
}

func ensureFileNotExists(logger *sLogger, name string) {
	s, _ := os.Stat(name)
	logger.FatalfIf(s != nil, "%s: File exists", name)
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
