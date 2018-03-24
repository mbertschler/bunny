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

package router

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mbertschler/blocks/html"
	"github.com/mbertschler/bunny/pkg/blocks"
	"github.com/mbertschler/bunny/pkg/data"
	"github.com/mbertschler/bunny/pkg/guiapi"
)

func Router(root string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	mountFileServer(r, "/static/", root, "js", "node_modules")
	mountFileServer(r, "/js/", root, "js", "src")
	r.Method("POST", "/gui/", guiapi.Handlers())
	r.Mount("/", pages())
	return r
}

func pages() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/item/{id}", viewItemPage)
	r.Get("/list/{id}", viewListPage)
	r.Get("/focus/", viewFocusPage)
	r.Get("/", viewAreaPage)
	return r
}

func mountFileServer(router *chi.Mux, url string, path ...string) {
	dir := http.Dir(filepath.Join(path...))
	fs := http.StripPrefix(url,
		http.FileServer(dir))
	router.Mount(url, fs)
}

func viewItemPage(w http.ResponseWriter, r *http.Request) {
	id, err := intFromUrl(r, "id")
	if err != nil {
		log.Println(err)
	}
	item, err := data.UserItemByID(1, id)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(blocks.LayoutBlock(blocks.ViewItemPage(item)), w)
	if err != nil {
		log.Println(err)
	}
}

func viewAreaPage(w http.ResponseWriter, r *http.Request) {
	_, things, err := data.UserArea(1, 1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(blocks.LayoutBlock(blocks.ViewAreaPage(things)), w)
	if err != nil {
		log.Println(err)
	}
}

func viewListPage(w http.ResponseWriter, r *http.Request) {
	id, err := intFromUrl(r, "id")
	if err != nil {
		log.Println(err)
	}
	list, err := data.UserItemList(1, id)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(blocks.LayoutBlock(blocks.ViewListPage(list)), w)
	if err != nil {
		log.Println(err)
	}
}

func viewFocusPage(w http.ResponseWriter, r *http.Request) {
	focus, err := data.FocusList(1)
	if err != nil {
		log.Println(err)
	}
	err = html.Render(blocks.LayoutBlock(blocks.ViewFocusPage(focus)), w)
	if err != nil {
		log.Println(err)
	}
}

func intFromUrl(r *http.Request, name string) (int, error) {
	ctx := chi.RouteContext(r.Context())
	str := ctx.URLParam(name)
	return strconv.Atoi(str)
}
