package trig

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
	subdir = ".trig"
	file   = "associations.db"
)

type Association struct {
	TriggeringFileName string
	DependentFileNames []string
}

var DBFileName = ""

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Join(usr.HomeDir, subdir)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(fmt.Sprintf("error creating '%s': %s", dir, err))
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
		var depFiles []string
		d := gob.NewDecoder(bytes.NewReader(v))
		if err := d.Decode(&depFiles); err != nil {
			store.Close()
			log.Fatal(err)
		}
		associations = append(associations, &Association{string(k), depFiles})
		fileNames = append(fileNames, append(depFiles, string(k))...)
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
			strings.TrimPrefix(a.TriggeringFileName, base),
			strmap(a.DependentFileNames, func(s string) string { return strings.TrimPrefix(s, base) }),
		)
	}
}

func Set(dependentFileName string, triggeringFileNames []string) {
	dependentFile := getFullPath(dependentFileName)
	triggeringFiles := getFullPaths(triggeringFileNames)
	store, err := skv.Open(DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	triggeringFilesToGo := append([]string(nil), triggeringFiles...)
	var toUpdate []*Association
	store.ForEach(func(k, v []byte) error {
		templateFileName, depFiles, err := decode(k, v)
		if err != nil {
			store.Close()
			log.Fatal(err)
		}
		triggeringFilesToGo = cleanSlice(triggeringFilesToGo, templateFileName)
		if isInSlice(triggeringFiles, templateFileName) {
			if added := appendIfNecessary(depFiles, dependentFile); len(added) > len(depFiles) {
				toUpdate = append(toUpdate, &Association{templateFileName, added})
			}
		} else if cleaned := cleanSlice(depFiles, dependentFile); len(cleaned) < len(depFiles) {
			toUpdate = append(toUpdate, &Association{templateFileName, cleaned})
		}
		return nil
	})
	for _, t := range triggeringFilesToGo {
		toUpdate = append(toUpdate, &Association{t, []string{dependentFile}})
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
		triggeringFileName, depFiles, err := decode(k, v)
		if err != nil {
			store.Close()
			log.Fatal(err)
		}
		if triggeringFileName == tbr {
			toUpdate = append(toUpdate, &Association{triggeringFileName, []string{}})
		} else if isInSlice(depFiles, tbr) {
			toUpdate = append(toUpdate, &Association{triggeringFileName, cleanSlice(depFiles, tbr)})
		}
		return nil
	})
	err = updateAssociations(store, toUpdate)
	store.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Handle(triggeringFileName string) {
	store, err := skv.OpenReadOnly(DBFileName)
	if err != nil {
		fmt.Print("trig: could not open db (read-only)")
		log.Fatal(err)
	}
	t := getFullPath(triggeringFileName)
	var depFiles []string
	err = store.Get(t, &depFiles)
	store.Close()
	if err == skv.ErrNotFound {
		fmt.Printf("trig: no associations for %s\n", t)
		return
	}
	if err != nil {
		log.Print("trig: could not get associated files for this template")
		log.Fatal(err)
	}
	for _, g := range depFiles {
		// TODO: waiting on https://github.com/bdkjones/CodeKit/issues/463
		// and since `touch` doesn't work, here is a horrible temporary hack:
		// first move out the file to a new path (starting with _), wait 2s,
		// then move the file back to the original path
		dir, file := path.Split(g)
		tmpPath := path.Join(dir, "_"+file)
		exe("mv", g, tmpPath)
		defer exe("mv", tmpPath, g)
		time.Sleep(100 * time.Microsecond)
	}
	// wait 2s before the deferred moves
	time.Sleep(2 * time.Second)
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
		fmt.Print("trig: could not execute " + command + " ")
		fmt.Println(args)
		panic(err)
	}
}

func updateAssociations(store *skv.KVStore, associations []*Association) error {
	for _, a := range associations {
		if len(a.DependentFileNames) == 0 {
			if err := store.Delete(a.TriggeringFileName); err != nil {
				return err
			}
			continue
		}
		if err := store.Put(a.TriggeringFileName, a.DependentFileNames); err != nil {
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
		log.Fatal(fmt.Sprintf("trig: could not resolve full path: %s\n", err))
	}
	return s
}

func getFullPaths(pathToFiles []string) (results []string) {
	for _, p := range pathToFiles {
		results = append(results, getFullPath(p))
	}
	return
}
