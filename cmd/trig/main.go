package main

import (
	"fmt"
	"log"
	"os"

	trig "github.com/zengabor/trig"
)

const (
	appName = "trig"
	version = "1.0"
)

func main() {
	command := "help"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}
	fmt.Printf("%v\n", os.Args)
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
Tool to set associations between  files and templates, so when it handles a triggering file it will "touch" all associated files (currently this is implemented by moving the file to a temporary directory, then moving it back after 2 seconds). Consequently a build tool can react and process those files. Associations are stored in %s

Usage:      
    %[1]s <command> [<args>...]
    %[1]s register <file-suffix> <command>
    %[1]s set <dependent-file> <triggering-file-1> <triggering-file-2>...
    %[1]s set <dependent-file>
    %[1]s list
    %[1]s handle <triggering-file>
    %[1]s help

Commands:

  register  Registers a command for a file suffix (typically the file extension). The string "$1" in the command is replaced with the file name.
  set       Associates triggering files with a dependent file (any existing triggering files not mentioned are removed).
  list      List associations.
  handle    Executes the registered command on all dependent files that were associated with the provided triggering file.
  help      Prints this help screen.

Examples:
  %[1]s register .go 'go run $1'
  %[1]s set www/index.go templates/one.gohtml templates/two.gohtml
  %[1]s set www/index.go
  %[1]s handle templates/two.gohtml
  %[1]s list

`, appName, version, trig.dbFileName)
}
