package gohtml

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
	subdir        = ".gohtml"
	file          = "gohtml.db"
	handleCommand = "/usr/bin/touch"
)

type Association struct {
	TemplateFileName string
	GoFileNames      []string
}

var DBFileName = ""

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Join(usr.HomeDir, subdir)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	DBFileName = path.Join(dir, file)
}

func GetAll() (associations []*Association) {
	store, err := skv.Open(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	store.ForEach(func(k, v []byte) error {
		var goFiles []string
		d := gob.NewDecoder(bytes.NewReader(v))
		if err := d.Decode(&goFiles); err != nil {
			store.Close()
			log.Fatal(err)
		}
		associations = append(associations, &Association{string(k), goFiles})
		return nil
	})
	store.Close()
	return
}

func Set(goFile string, templates []string) {
	goFile = getFullPath(goFile)
	templates = getFullPaths(templates)
	store, err := skv.Open(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	templatesToProcess := append([]string(nil), templates...)
	var toUpdate []*Association
	store.ForEach(func(k, v []byte) error {
		templateFileName, goFiles, err := decode(k, v)
		templatesToProcess = cleanSlice(templatesToProcess, templateFileName)
		if err != nil {
			store.Close()
			log.Fatal(err)
		}
		if isInSlice(templates, templateFileName) {
			if added := appendIfNecessary(goFiles, goFile); len(added) > len(goFiles) {
				toUpdate = append(toUpdate, &Association{templateFileName, added})
			}
		} else if cleaned := cleanSlice(goFiles, goFile); len(cleaned) < len(goFiles) {
			toUpdate = append(toUpdate, &Association{templateFileName, cleaned})
		}
		return nil
	})
	for _, t := range templatesToProcess {
		toUpdate = append(toUpdate, &Association{t, []string{goFile}})
	}
	err = updateAssociations(store, toUpdate)
	store.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Unset(toBeRemoved string) {
	tbr := getFullPath(toBeRemoved)
	store, err := skv.Open(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	var toUpdate []*Association
	store.ForEach(func(k, v []byte) error {
		templateFileName, goFiles, err := decode(k, v)
		if err != nil {
			store.Close()
			log.Fatal(err)
		}
		log.Printf(">>>>>> %v >>> %+v\n", templateFileName, goFiles)
		if templateFileName == tbr {
			toUpdate = append(toUpdate, &Association{templateFileName, []string{}})
			log.Println("removing " + templateFileName)
		} else if isInSlice(goFiles, tbr) {
			toUpdate = append(toUpdate, &Association{templateFileName, cleanSlice(goFiles, tbr)})
			log.Printf("cleaned: %+v\n", cleanSlice(goFiles, tbr))
		}
		return nil
	})
	err = updateAssociations(store, toUpdate)
	store.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Handle(templateFileName string) {
	store, err := skv.OpenReadOnly(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	t := getFullPath(templateFileName)
	var goFiles []string
	err = store.Get(t, &goFiles)
	if err == skv.ErrNotFound {
		store.Close()
		log.Fatal("gohtml: no associasions for " + t)
	}
	if err != nil {
		store.Close()
		log.Fatal(err)
	}
	for _, g := range goFiles {
		cmd := exec.Command(handleCommand + " " + g)
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

func updateAssociations(store *skv.KVStore, associations []*Association) error {
	for _, a := range associations {
		if len(a.GoFileNames) == 0 {
			if err := store.Delete(a.TemplateFileName); err != nil {
				return err
			}
			continue
		}
		if err := store.Put(a.TemplateFileName, a.GoFileNames); err != nil {
			return err
		}
		log.Printf("gohtml: '%s' is associated with %+v\n", a.TemplateFileName, a.GoFileNames)
	}
	return nil
}

func decode(k, v []byte) (key string, value []string, err error) {
	key = string(k)
	d := gob.NewDecoder(bytes.NewReader(v))
	err = d.Decode(&value)
	return
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

func getFullPaths(pathToFiles []string) (results []string) {
	for _, p := range pathToFiles {
		results = append(results, getFullPath(p))
	}
	return
}
