package trig

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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
