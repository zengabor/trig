package trig

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/zengabor/skv"
)

func Set(dependentFileName string, triggeringFileNames []string) {
	dependentFile := getFullPath(dependentFileName)
	triggeringFiles := getFullPaths(triggeringFileNames)
	store, err := skv.Open(dbFileName)
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
