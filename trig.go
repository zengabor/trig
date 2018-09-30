package trig

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

const (
	subdir = ".trig"
	file   = "associations.db"
)

var dbFileName = ""

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Join(usr.HomeDir, subdir)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(fmt.Sprintf("error creating '%s': %s", dir, err))
	}
	dbFileName = path.Join(dir, file)
}
