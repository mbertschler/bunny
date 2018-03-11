package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

var config = struct {
	port string // $BUNNY_PORT
	root string // $BUNNY_ROOT
}{}

func setupConfig() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setFromEnvironment(&config.port, "BUNNY_PORT", "3080")
	setFromEnvironment(&config.root, "BUNNY_ROOT", "")
	if config.root == "" {
		var err error
		config.root, err = findProjectFolder()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func setFromEnvironment(target *string, name, fallback string) {
	val, ok := os.LookupEnv(name)
	if ok {
		*target = val
	} else {
		*target = fallback
	}
}

func findProjectFolder() (string, error) {
	gopath := os.Getenv("GOPATH")
	paths := filepath.SplitList(gopath)
	for _, p := range paths {
		project := filepath.Join(p, "src", "github.com", "mbertschler", "bunny")
		info, err := os.Stat(project)
		if err == nil && info.IsDir() {
			return project, nil
		}
	}
	return "", errors.New("couldn't find the project in GOPATH")
}
