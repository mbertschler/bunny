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
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mbertschler/blocks/html"
	"github.com/mbertschler/bunny/pkg/data"
)

func main() {
	setupConfig()
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
	err = html.Render(layoutBlock(viewItemBlock(item)), w)
	if err != nil {
		log.Println(err)
	}
}

func renderAreaPage(w http.ResponseWriter, r *http.Request) {
	_, things, err := data.UserArea(1, 1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(layoutBlock(viewThingsBlock(things)), w)
	if err != nil {
		log.Println(err)
	}
}

func renderListPage(w http.ResponseWriter, r *http.Request) {
	list, err := data.UserItemList(1, 1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(layoutBlock(viewListBlock(list)), w)
	if err != nil {
		log.Println(err)
	}
}

func renderFocusPage(w http.ResponseWriter, r *http.Request) {
	focus, err := data.FocusList(1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(layoutBlock(displayFocusBlock(focus)), w)
	if err != nil {
		log.Println(err)
	}
}
