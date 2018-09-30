package trig

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path"
	"time"

	"github.com/zengabor/skv"
)

func Handle(triggeringFileName string) {
	store, err := skv.OpenReadOnly(dbFileName)
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
		// first move out the file to a new path (e.g., `mv a.go tmp/_a.go`),
		// wait 2s, then move the file back to its original path.
		dir, file := path.Split(g)
		tempDir := path.Join(dir, "tmp")
		exe("mkdir", "-p", tempDir)
		tmpPath := path.Join(tempDir, "_"+file)
		exe("mv", g, tmpPath)
		defer exe("mv", tmpPath, g)
		time.Sleep(100 * time.Microsecond)
	}
	// wait 2s before the deferred moves
	time.Sleep(2 * time.Second)
}

func exe(command string, args ...string) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Print(out.String())
	if err != nil {
		fmt.Printf("trig: could not execute %s: %s\n", args, err)
		panic(err)
	}
	fmt.Printf("trig: executed %s\n", args)
}
