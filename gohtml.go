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
	"time"

	"github.com/zengabor/skv"
)

const (
	subdir = ".gohtml"
	file   = "gohtml.db"
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

func List() {
	store, err := skv.Open(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	var associations []*Association
	var fileNames []string
	store.ForEach(func(k, v []byte) error {
		var goFiles []string
		d := gob.NewDecoder(bytes.NewReader(v))
		if err := d.Decode(&goFiles); err != nil {
			store.Close()
			log.Fatal(err)
		}
		associations = append(associations, &Association{string(k), goFiles})
		fileNames = append(fileNames, append(goFiles, string(k))...)
		return nil
	})
	store.Close()
	base := fileNames[0]
	for i := 1; i < len(fileNames); i++ {
		f := fileNames[i]
		l := len(base)
		if len(f) < l {
			l = len(f)
		}
		for l > 1 {
			if strings.HasPrefix(f[:l], base[:l]) {
				base = base[:l]
				break
			}
			l--
		}
	}
	fmt.Printf("Common base directory: %s\n", base)
	for i, a := range associations {
		fmt.Printf(
			"%d. %s: %+v\n",
			i,
			strings.TrimPrefix(a.TemplateFileName, base),
			strmap(a.GoFileNames, func(s string) string { return strings.TrimPrefix(s, base) }),
		)
	}
}

func Set(goFile string, templates []string) {
	goFile = getFullPath(goFile)
	templates = getFullPaths(templates)
	store, err := skv.Open(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	templatesToGo := append([]string(nil), templates...)
	var toUpdate []*Association
	store.ForEach(func(k, v []byte) error {
		templateFileName, goFiles, err := decode(k, v)
		if err != nil {
			store.Close()
			log.Fatal(err)
		}
		templatesToGo = cleanSlice(templatesToGo, templateFileName)
		if isInSlice(templates, templateFileName) {
			if added := appendIfNecessary(goFiles, goFile); len(added) > len(goFiles) {
				toUpdate = append(toUpdate, &Association{templateFileName, added})
			}
		} else if cleaned := cleanSlice(goFiles, goFile); len(cleaned) < len(goFiles) {
			toUpdate = append(toUpdate, &Association{templateFileName, cleaned})
		}
		return nil
	})
	for _, t := range templatesToGo {
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
		if templateFileName == tbr {
			toUpdate = append(toUpdate, &Association{templateFileName, []string{}})
		} else if isInSlice(goFiles, tbr) {
			toUpdate = append(toUpdate, &Association{templateFileName, cleanSlice(goFiles, tbr)})
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
		fmt.Print("gohtml: could not open db (read-only)")
		log.Fatal(err)
	}
	t := getFullPath(templateFileName)
	var goFiles []string
	err = store.Get(t, &goFiles)
	store.Close()
	if err == skv.ErrNotFound {
		fmt.Printf("gohtml: no associations for %s\n", t)
		return
	}
	if err != nil {
		log.Print("gohtml: could not get associated files for this template")
		log.Fatal(err)
	}
	for _, g := range goFiles {
		// TODO: waiting on https://github.com/bdkjones/CodeKit/issues/463
		// and since `touch` doesn't work this, is a horrible temporary hack:
		// first move out the file to a temp dir, wait 3s, then move it back
		dir, file := path.Split(g)
		tempDir := path.Join(dir, "tmp")
		exe("mkdir", "-p", tempDir)
		exe("mv", g, path.Join(tempDir, "_"+file))
		defer exe("mv", path.Join(tempDir, "_"+file), g)
	}
	// wait a little bit more than 2s before the deferred moves
	time.Sleep(2*time.Second + 100*time.Microsecond)
}

func GetFullPath(pathToFile string) (string, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(pathToFile, "/") {
		pathToFile = path.Join(workingDirectory, pathToFile)
	}
	s, err := filepath.Abs(pathToFile)
	if err != nil {
		return "", err
	}
	return s, nil
}

func exe(command string, args ...string) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Print(out.String())
	if err != nil {
		fmt.Print("gohtml: could not execute " + command + " ")
		fmt.Println(args)
		panic(err)
	}
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

func strmap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
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
	s, err := GetFullPath(pathToFile)
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
