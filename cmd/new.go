package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"s19f.io/x/log"
)

const bin = "/Users/shion.t.fujie/Desktop/machinery/bin"
const visualStudioCode = "Visual Studio Code"

const fanfareTemplate = `Have created %s!
At %s
%s
`

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
			"button_background_hover" : [26, 115, 232],
			"frame": [255, 255, 255],
			"toolbar": [255, 255, 255],
			"tab_text" : [117, 117, 117],
			"tab_background_text" : [117, 117, 117],
			"bookmark_text" : [117, 117, 117],
			"toolbar_button_icon" : [117, 117, 117]
		}
	}
}
`

const chromeXManifestTemplate = `{
	"name": "%s",
	"description": "",
	"version": "0.1.0",
	"manifest_version": 2,
	"icons": {
		"16": "images/get_started16.png",
		"32": "images/get_started32.png",
		"48": "images/get_started48.png",
		"128": "images/get_started128.png"
	},
	"permissions": [],
	"content_scripts": [],
	"background": {}
}
`

const puppeteerPackageJSONTemplate = `{
	"name": "%s",
	"version": "0.0.1",
	"description": "",
	"scripts": {
		"test": "echo \"Error: no test specified\" && exit 1"
	},
	"keywords": [],
	"author": "Shion Fujie (https://github.com/shionfujie)",
	"dependencies": {
		"puppeteer": "latest"
	}
}
`

const puppeteerMainFileTemplate = `const puppeteer = require('puppeteer');

(async () => {
	const browser = await puppeteer.launch();
	const page = await browser.newPage();

	console.log("Hello, %s!!!")

	await browser.close();
})();`

const scalaRootPackage = "io.s19f"
const buildSbtTemplate = `scalaVersion := "2.12.12"`
const buildPropertiesTemplate = `sbt.version=1.2.8`
const projectPluginsSbtTemplate = `initialize ~= (_ => sys.props("scala.repl.maxprintstring") = "0" ) // Sets no limit to print a large string`
const scalaEntryFileTemplate = `package %s.%s

object %s
`
const scalafmtConfTemplate = `version = "2.7.4"
align.preset=most
maxColumn = 200
newlines.avoidForSimpleOverflow=[tooLong]
rewrite.rules = [PreferCurlyFors]
align.arrowEnumeratorGenerator = false
newlines.beforeCurlyLambdaParams = multilineWithCaseOnly`
const scalaGitignoreTemplate = `target
.metals
metals.sbt
.vscode
/project/project
.bloop
.DS_Store`

func main() {
	logger := log.New(os.Stdout, "new: ", 0)

	logger.FatalfIf(len(os.Args) < 2, "No subcommand specified")
	subcommand := os.Args[1]
	switch subcommand {
	case "sh":
		logger.SetPrefix("new sh: ")
		logger.FatalfIf(len(os.Args) < 3, "No file name specified")

		n := os.Args[2]
		p := path.Join(os.Getenv("SBIN"), n)
		ensureFileNotExists(logger, p)

		err := ioutil.WriteFile(p, []byte("#!/bin/bash\n\n"), 0744)
		logger.FatalfIfError(err, "%s: Failed to create an shell script executable", n)

		fmt.Fprintf(logger.O, "A shell script executable has been created at '%s'\n", p) // Prints without the predefined format

		err = exec.Command("open", "-a", visualStudioCode, p).Run()
		logger.FatalfIfError(err, " %s: Failed to open with %s", n, visualStudioCode)
	case "go-cmd":
		logger.SetPrefix("new go-cmd: ")
		logger.FatalfIf(len(os.Args) < 3, "No command name specified. Specify a command name.")

		cmdName := os.Args[2]
		ensureFileNotExists(logger, cmdName)

		os.Mkdir(cmdName, 0744)
		os.Chdir(cmdName)
		goPackage := "s19f.io/cli/" + cmdName
		err := exec.Command("go", "mod", "init", goPackage).Run()
		logger.FatalfIfError(err, "%s: Failed to create the GO package '%s'", cmdName, goPackage)

		os.Mkdir("cmd", 0744)
		os.Chdir("cmd")
		code := fmt.Sprintf(goMainFileTemplate, cmdName)
		err = ioutil.WriteFile(cmdName+".go", []byte(code), 0744)
		logger.FatalfIfError(err, "%s: Failed to create a main file at '%s', ")

		fmt.Fprintln(logger.O, "Have created a command-line project. Strive to code!!!")

		exec.Command("open", "-a", visualStudioCode, "../../"+cmdName, cmdName+".go").Run() // Try to open the project
	case "chrome-theme":
		logger.SetPrefix("new chrome-theme: ")
		logger.FatalfIf(len(os.Args) < 3, "Extension name argument expected")

		themeName := os.Args[2]
		projectPath := path.Join(os.Getenv("CHROME_HOME"), "src/theme", themeName)
		ensureFileNotExists(logger, projectPath)
		os.Mkdir(projectPath, 0744)

		manifestPath := path.Join(projectPath, "manifest.json")
		json := fmt.Sprintf(chromeThemeManifestTemplate, themeName)
		err := ioutil.WriteFile(manifestPath, []byte(json), 0744)
		logger.FatalfIfError(err, "%s: Failed to create a manifest file for a chrome extension", themeName)

		fmt.Fprintf(logger.O, fanfareTemplate, "a chrome theme", projectPath, "Excited to decorate!!!")

		exec.Command("open", "-a", visualStudioCode, projectPath, manifestPath).Run() // Try to open the project
	case "chrome", "chrome-x":
		logger.SetPrefix("new chrome-x: ")
		logger.FatalfIf(len(os.Args) < 3, "Extension name argument expected")

		n := os.Args[2]
		ensureFileNotExists(logger, n)

		jsDir := n + "/js/"
		imagesDir := n + "/images/"
		os.MkdirAll(jsDir, 0744)
		os.MkdirAll(imagesDir, 0744)

		imagesTemplateDir := os.Getenv("CHROME_HOME") + "/template/images/"
		images := []string{"get_started16.png", "get_started32.png", "get_started48.png", "get_started128.png"}
		for _, name := range images {
			err := copyFile(imagesDir+name, imagesTemplateDir+name)
			logger.FatalfIfError(err, "%s: Failed to copy %s to %s: %v\n", n, imagesTemplateDir, imagesDir, err)
		}

		json := fmt.Sprintf(chromeXManifestTemplate, n)
		manifestFile := n + "/manifest.json"
		err := ioutil.WriteFile(manifestFile, []byte(json), 0744)
		logger.FatalfIfError(err, "%s: Failed to create a manifest file for a chrome extension", n)

		exec.Command("open", "-a", visualStudioCode, n, manifestFile).Run() // Try to open the project
	case "puppeteer", "chrome-script", "web-script":
		logger.SetPrefix("new " + subcommand + ": ")
		logger.FatalfIf(len(os.Args) < 3, "Project name argument expected")

		projectName := os.Args[2]
		ensureFileNotExists(logger, projectName)
		os.Mkdir(projectName, 0744)

		packageJSONPath := path.Join(projectName, "package.json")
		packageJSON := fmt.Sprintf(puppeteerPackageJSONTemplate, projectName)
		err := ioutil.WriteFile(packageJSONPath, []byte(packageJSON), 0744)
		logger.FatalfIfError(err, "%s: Failed to create package.json", projectName)

		mainFilePath := path.Join(projectName, "main.js")
		mainFile := fmt.Sprintf(puppeteerMainFileTemplate, projectName)
		err = ioutil.WriteFile(mainFilePath, []byte(mainFile), 0744)
		logger.FatalfIfError(err, "%s: Failed to create main.js", projectName)

		fmt.Fprintf(logger.O, fanfareTemplate, "a web script project with Puppeteer", projectName, "Automate everything!!!")

		exec.Command("open", "-a", visualStudioCode, projectName, mainFilePath).Run() // Try to open the project
	case "scala", "scala-console":
		logger.SetPrefix("new " + subcommand + ": ")
		logger.FatalfIf(len(os.Args) < 3, "Project name argument expected")

		
		projectName := os.Args[2]
		ensureFileNotExists(logger, projectName)

		projectDir := path.Join(projectName, "project")
		pkgDir := strings.Join(strings.Split(scalaRootPackage, "."), "/")
		srcDir := path.Join(projectName, "src/main/scala", pkgDir, projectName)
		os.MkdirAll(projectDir, 0744)
		os.MkdirAll(srcDir, 0744)
		
		buildSbtPath := path.Join(projectName, "build.sbt")
		err := ioutil.WriteFile(buildSbtPath, []byte(buildSbtTemplate), 0744)
		logger.FatalfIfError(err, "Failed to create build.sbt")
		
		confFilepath := path.Join(projectName, ".scalafmt.conf")
		err = ioutil.WriteFile(confFilepath, []byte(scalafmtConfTemplate), 0744)
		logger.FatalfIfError(err, "Failed to create .scalafmt.conf")
		
		gitignorePath := path.Join(projectName, ".gitignore")
		err = ioutil.WriteFile(gitignorePath, []byte(scalaGitignoreTemplate), 0744)
		logger.FatalfIfError(err, "Failed to create .scalafmt.conf")
		
		buildPropertiesPath := path.Join(projectDir, "build.properties")
		err = ioutil.WriteFile(buildPropertiesPath, []byte(buildPropertiesTemplate), 0744)
		logger.FatalfIfError(err, "Failed to create build.properties")

		pluginsSbtPath := path.Join(projectDir, "plugins.sbt")
		err = ioutil.WriteFile(pluginsSbtPath, []byte(projectPluginsSbtTemplate), 0744)
		logger.FatalfIfError(err, "Failed to create plugins.sbt")

		entryName := strings.Title(projectName)
		entryFilepath := path.Join(srcDir, entryName+".scala")
		entryFile := fmt.Sprintf(scalaEntryFileTemplate, scalaRootPackage, projectName, entryName)
		err = ioutil.WriteFile(entryFilepath, []byte(entryFile), 0744)
		logger.FatalfIfError(err, "Failed to create an entry file: %s", entryFilepath)

		exec.Command("open", "-a", visualStudioCode, projectName, entryFilepath, buildSbtPath).Run()
	default:
		logger.Fatalf("%s: No such subcommand", subcommand)
	}
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Failed to open '%s'", src)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("Failed to create a new file '%s'", dst)
	}

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("Failed to copy %s to %s", src, dst)
	}

	if err = out.Close(); err != nil {
		return fmt.Errorf("Failed to copy %s to %s", src, dst)
	}
	return nil
}

func ensureFileNotExists(logger *log.SLogger, name string) {
	s, _ := os.Stat(name)
	logger.FatalfIf(s != nil, "%s: File exists", name)
}
