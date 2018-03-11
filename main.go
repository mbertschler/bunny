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
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mbertschler/blocks/html"
	"github.com/mbertschler/bunny/pkg/data"
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
	log.Println(http.ListenAndServe(":"+config.port, router()))
}

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Mount("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(
				filepath.Join(config.root, "js", "node_modules")))))
	r.Mount("/js/",
		http.StripPrefix("/js/",
			http.FileServer(http.Dir(
				filepath.Join(config.root, "js", "src")))))
	r.Post("/gui/", guiAPI().ServeHTTP)
	r.Get("/item/{id}", renderItemPage)
	r.Get("/list/{id}", renderListPage)
	r.Get("/focus/", renderFocusPage)
	r.Get("/", renderAreaPage)
	return r
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

func renderItemPage(w http.ResponseWriter, r *http.Request) {

	ctx := chi.RouteContext(r.Context())
	idStr := ctx.URLParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
	}
	item, err := data.UserItemByID(1, id)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(pageBlock(displayItemBlock(item)), w)
	if err != nil {
		log.Println(err)
	}
}

func renderAreaPage(w http.ResponseWriter, r *http.Request) {
	_, things, err := data.UserArea(1, 1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(pageBlock(displayThingsBlock(things)), w)
	if err != nil {
		log.Println(err)
	}
}

func renderListPage(w http.ResponseWriter, r *http.Request) {
	list, err := data.UserItemList(1, 1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(pageBlock(displayListBlock(list)), w)
	if err != nil {
		log.Println(err)
	}
}

func renderFocusPage(w http.ResponseWriter, r *http.Request) {
	focus, err := data.FocusList(1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(pageBlock(displayFocusBlock(focus)), w)
	if err != nil {
		log.Println(err)
	}
}
