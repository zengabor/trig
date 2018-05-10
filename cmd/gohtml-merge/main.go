package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/zengabor/gohtml"
)

const (
	appName = "gohtml-merge"
	version = "1.0"
)

var re = regexp.MustCompile(`{{\s*include\s+"([^"]+)"\s*}}`)

func main() {
	switch {
	case len(os.Args) < 2:
		help()
	case len(os.Args) == 2:
		mergeFiles(os.Args[1])
	default:
		log.Fatal(appName + ": provide exactly one file path to a template")
	}
}

func getFileContent(filePath string) (s string, err error) {
	p, err := gohtml.GetFullPath(filePath)
	if err != nil {
		return
	}
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	return string(b), nil
}

func mergeFiles(pathToTemplate string) {
	tmpl, err := getFileContent(pathToTemplate)
	if err != nil {
		log.Fatal(appName + ": could not open " + pathToTemplate)
	}
	fmt.Fprintf(os.Stdout, re.ReplaceAllStringFunc(tmpl, func(s string) string {
		c, err := getFileContent(re.ReplaceAllString(s, "$1"))
		if err != nil {
			p, err := gohtml.GetFullPath(pathToTemplate)
			if err != nil {
				return appName + ": error resolving " + pathToTemplate
			}
			return appName + ": error merging " + p
		}
		return c
	}))
}

func help() {
	fmt.Printf(`%s %s // github.com/zengabor/gohtml

Usage:    %[1]s <pathToTemplateFile>

Example:
  %[1]s templates/main.gohtml > main.html

`, appName, version)
}
