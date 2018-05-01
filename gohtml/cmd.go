package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zengabor/gohtml"
)

const (
	version = "1.0"
)

func main() {
	command := "help"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}
	switch command {
	case "list":
		for i, a := range gohtml.GetAll() {
			fmt.Printf("%d. %s <-> %+v\n", i, a.TemplateFileName, a.GoFileNames)
		}
	case "set":
		if len(os.Args) < 4 {
			log.Fatal("gohtml: provide at least one template file path to set")
		}
		gohtml.Set(os.Args[2], os.Args[3:])
	case "unset":
		if len(os.Args) != 3 {
			log.Fatal("gohtml: provide exactly one file path to unset")
		}
		gohtml.Unset(os.Args[2])
	case "handle":
		if len(os.Args) != 3 {
			log.Fatal("gohtml: provide exactly one template file path for handle")
		}
		gohtml.Handle(os.Args[2])
	case "help":
		help()
	default:
		log.Fatal(fmt.Sprintf("gohtml: unknown command '%s'\n", os.Args[1]))
	}
}

func help() {
	fmt.Printf(`gohtml %s // github.com/zengabor/gohtml

Sets associations between (go) files and templates, so when it handles a template it will "touch" the associated files (currently this is implemented by moving the file to a temporary directory, then moving it back after 2 seconds). Consequently a build process can react and process those files. Associations are stored in %s

Usage:    gohtml <command> <args>

Available commands:
  set     Associates a file with templates. (If a currently associated template is not mentioned, the association is removed.)
  unset   Removes all associations of a file or template.
  list    List associations.
  handle  Touches all files that are associated with the provided template.
  help    Prints this screen.

Examples:
  gothml set path/index.go one.gohtml b/two.gohtml
  gohtml list
  gohtml handle b/two.gohtml
  gothml unset one.gohtml
  gothml unset path/index.go

`, version, gohtml.DBFileName)
}
