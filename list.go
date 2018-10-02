package trig

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strings"

	"github.com/zengabor/skv"
)

func List() {
	store, err := skv.Open(dbFileName)
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
	if len(associations) == 0 {
		fmt.Println("No associations yet.")
		return
	}
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
	fmt.Println("Listing triggeringFile [dependent files]")
	for i, a := range associations {
		fmt.Printf(
			"%d. %s: %+v\n",
			i,
			strings.TrimPrefix(a.TriggeringFileName, base),
			strmap(a.DependentFileNames, func(s string) string { return strings.TrimPrefix(s, base) }),
		)
	}
}

func strmap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
