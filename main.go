// Copyright 2018 Martin Bertschler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var config = struct {
	port string // $BUNNY_PORT
	root string // $BUNNY_ROOT
}{}

func main() {
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
	log.Println("Bunny :) running at port", config.port)
	http.Handle("/", http.FileServer(http.Dir(config.root)))
	log.Println(http.ListenAndServe(":"+config.port, nil))
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
