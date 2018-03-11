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

package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

var (
	Port string // $BUNNY_PORT
	Root string // $BUNNY_ROOT
)

func Setup() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setFromEnvironment(&Port, "BUNNY_PORT", "3080")
	setFromEnvironment(&Root, "BUNNY_ROOT", "")
	if Root == "" {
		var err error
		Root, err = findProjectFolder()
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
