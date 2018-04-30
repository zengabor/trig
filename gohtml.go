package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/zengabor/skv"
)

const (
	version = "1.0"
	subdir  = ".gohtml"
	file    = "gohtml.db"
	execute = "/usr/bin/touch"
)

type Association struct {
	TemplateFileName string
	GoFileNames      []string
}

var dbFileName = ""

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Join(usr.HomeDir, subdir)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	dbFileName = path.Join(dir, file)
}

func main() {
	command := "help"
	if len(os.Args) == 2 && os.Args[1] != "help" {
		log.Fatal("gohtml: invalid number of arguments")
	}
	if len(os.Args) >= 3 {
		command = os.Args[1]
	}
	switch command {
	case "set":
		if len(os.Args) < 4 {
			log.Fatal("gohtml: provide at least one template file name to set")
		}
		Set(os.Args[2], os.Args[3:])
	case "handle":
		if len(os.Args) != 3 {
			log.Fatal("gohtml: provide exactly one template file name for handle")
		}
		Handle(os.Args[2])
	case "help":
		Help()
	default:
		log.Fatal("gohtml: unknown command '%s'\n", os.Args[1])
	}
}

func Set(goFile string, templates []string) {
	goFile = getFullPath(goFile)
	templates = getFullPaths(templates)
	store, err := skv.Open(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	templatesToProcess := append([]string(nil), templates...)
	var toUpdate []*Association
	store.ForEach(func(k, v []byte) error {
		templateFileName := string(k)
		templatesToProcess = cleanSlice(templatesToProcess, templateFileName)
		var goFiles []string
		d := gob.NewDecoder(bytes.NewReader(v))
		if err := d.Decode(&goFiles); err != nil {
			store.Close()
			log.Fatal(err)
		}
		if isInSlice(templates, templateFileName) {
			if added := appendIfNecessary(goFiles, goFile); len(added) > len(goFiles) {
				toUpdate = append(toUpdate, &Association{templateFileName, added})
			}
		} else {
			if cleaned := cleanSlice(goFiles, goFile); len(cleaned) < len(goFiles) {
				toUpdate = append(toUpdate, &Association{templateFileName, cleaned})
			}
		}
		return nil
	})
	for _, t := range templatesToProcess {
		toUpdate = append(toUpdate, &Association{t, []string{goFile}})
	}
	for _, t := range toUpdate {
		if len(t.GoFileNames) == 0 {
			if err := store.Delete(t.TemplateFileName); err != nil {
				store.Close()
				log.Fatal(err)
			}
			continue
		}
		if err := store.Put(t.TemplateFileName, t.GoFileNames); err != nil {
			store.Close()
			log.Fatal(err)
		}
	}
	store.Close()
}

func Handle(templateFileName string) {
	store, err := skv.OpenReadOnly(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	t := getFullPath(templateFileName)
	var goFiles []string
	err = store.Get(t, &goFiles)
	if err == skv.ErrNotFound {
		store.Close()
		log.Fatal("gohtml: template not set")
	}
	if err != nil {
		store.Close()
		log.Fatal(err)
	}
	for _, g := range goFiles {
		cmd := exec.Command(execute + " " + g)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			store.Close()
			log.Fatal(err)
		}
		fmt.Println(out.String())
	}
	store.Close()
}

func Help() {
	fmt.Printf(`gohtml %s - Sets associations between (go) files and templates, so when you invoke it to handle a template it will touch the associated files so that the build process can pick them up. Associations are stored in %s

Usage:    gohtml <command> <args>

Available commands:
  set     Associates a file with templates. Obsolete associations are removed.
  handle  Touches all files that are associated with the provided template.
  help    Prints this screen.

Examples:
  gothml set path/index.go one.gohtml b/two.gohtml
  gohtml handle b/two.gohtml

`, version, dbFileName)
}

func isInSlice(s []string, v string) bool {
	for _, ss := range s {
		if ss == v {
			return true
		}
	}
	return false
}

func appendIfNecessary(s []string, v string) []string {
	for _, ss := range s {
		if ss == v {
			return s
		}
	}
	return append(s, v)
}

func cleanSlice(s []string, v string) []string {
	for i, ss := range s {
		if ss == v {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func getFullPath(pathToFile string) string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasPrefix(pathToFile, "/") {
		pathToFile = path.Join(workingDirectory, pathToFile)
	}
	s, err := filepath.Abs(pathToFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("gohtml: could not resolve full path: %s\n", err))
	}
	return s
}

func getFullPaths(pathToFiles []string) []string {
	var results []string
	for _, p := range pathToFiles {
		results = append(results, getFullPath(p))
	}
	return results
}
