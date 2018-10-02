package trig

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	for _, f := range depFiles {
		exe(
			"/usr/bin/osascript",
			"-e",
			fmt.Sprintf("tell application %q to process file at path %q", "CodeKit", f),
		)
	}
}

func exe(command string, args ...string) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		panic(fmt.Sprintf("trig: could not execute %s %s: %s\n", command, args, err))
	}
	fmt.Printf("trig: executed %s", command)
	fmt.Println(strings.Trim(fmt.Sprint(args), "[]"))
}
