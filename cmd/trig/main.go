package main

import (
	"fmt"
	"log"
	"os"

	trig "github.com/zengabor/trig"
)

const (
	appName = "trig"
	version = "2.0"
)

func main() {
	command := "help"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}
	switch command {
	case "list":
		trig.List()
	case "set":
		if len(os.Args) < 3 {
			log.Fatal(appName + ": provide at least the dependent file")
		}
		trig.Set(os.Args[2], os.Args[3:])
	case "handle":
		if len(os.Args) != 3 {
			log.Fatal(appName + ": provide exactly one file path to handle")
		}
		trig.Handle(os.Args[2])
	case "help":
		help()
	default:
		log.Fatal(fmt.Sprintf(appName+": unknown command '%s'\n", os.Args[1]))
	}
}

func help() {
	fmt.Printf(`%s %s // github.com/zengabor/trig
This command-line tool extends CodeKit (https://codekitapp.com). It can set
associations between template files (triggering files) and dependent files
(the ones using those template files). For example, if you update a template
"footer.html" and call trig with the file path ("trig handle footer.html")
then trig will call CodeKit via AppleScript to process all the files that are
including "footer.html". This way you can be sure that everything is updated
when you modify a template.
Associations are stored in %s

Usage:      
  %[1]s <command> [<args>...]
  %[1]s set <dependent-file> <triggering-file-1> <triggering-file-2>...
  %[1]s set <dependent-file>
  %[1]s handle <triggering-file>
  %[1]s list
  %[1]s help

Commands:
  set       Associates a dependent file with triggering filess. If a currently
            associated triggering file is not mentioned, that association is
            removed. If no triggering files are provided then all are removed.
  handle    Calls CodeKit to process all dependent files that were associated
            with the triggering file.
  list      List associations.
  help      Prints this help screen.

Examples:
  %[1]s set www/index.go templates/one.gohtml templates/two.gohtml
  %[1]s set www/index.go
  %[1]s handle templates/two.gohtml
  %[1]s list

`, appName, version, trig.DBFileName())
}
