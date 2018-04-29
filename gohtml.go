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

	"github.com/zengabor/skv"
)

const (
	version = "1.0"
	subdir  = ".gohtml"
	file    = "gohtml.db"
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
	case "run":
		if len(os.Args) != 3 {
			log.Fatal("gohtml: provide exactly one template file name for run")
		}
		Run(os.Args[2])
	case "help":
		Help()
	default:
		fmt.Printf("gohtml: unknown command '%s'\n", os.Args[1])
	}
}

func Set(goFile string, templates []string) {
	store, err := skv.Open(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	templatesToProcess := append([]string(nil), templates...)
	var toUpdate []*Association
	store.ForEach(func(k, v []byte) error {
		templateFileName := string(k)
		templatesToProcess = cleanSlice(templatesToProcess, templateFileName)
		var goFiles []string
		d := gob.NewDecoder(bytes.NewReader(v))
		if err := d.Decode(&goFiles); err != nil {
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
				log.Fatal(err)
			}
			continue
		}
		if err := store.Put(t.TemplateFileName, t.GoFileNames); err != nil {
			log.Fatal(err)
		}
	}
}

func Run(templateFileName string) {
	store, err := skv.OpenReadOnly(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	var goFiles []string
	err = store.Get(templateFileName, &goFiles)
	if err == skv.ErrNotFound {
		store.Close()
		log.Fatal("gohtml: template not set")
	}
	if err != nil {
		store.Close()
		log.Fatal(err)
	}
	for _, g := range goFiles {
		fmt.Printf("running %+v", g)
		cmd := exec.Command("go " + g)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf(": %q\n", out.String())
	}
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

func Help() {
	fmt.Printf(`gohtml %s - Sets associations between go files and templates. Invoke gohtml with the path to a template and it will run the associated go files. Associations are stored in %s

Usage: gohtml <command> <args>

Available commands:
  set   Associates a go file with templates. Earlier associations are deleted.
  run   Runs the go files that are associated with the template.
  help  Prints this screen.

Examples:
  gothml set path/index.go one.gohtml b/two.gohtml
  gohtml run b/two.gohtml

`, version, dbFileName)
}
