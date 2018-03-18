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
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	r := Router("/")
	shouldMatch(t, r, "GET", "/js/app.js")
	shouldMatch(t, r, "GET", "/static/jquery/dist/jquery.min.js")
	shouldMatch(t, r, "POST", "/gui/")
	shouldMatch(t, r, "GET", "/item/123")
	shouldMatch(t, r, "GET", "/list/123")
	shouldMatch(t, r, "GET", "/focus/")
	shouldNotMatch(t, r, "GET", "/x/focus/")
	shouldMatch(t, r, "GET", "/")
}

func shouldMatch(t *testing.T, r *chi.Mux, method, path string) {
	ctx := chi.NewRouteContext()
	if !r.Match(ctx, method, path) {
		t.Error(method, "route", path, "should match")
	}
}

func shouldNotMatch(t *testing.T, r *chi.Mux, method, path string) {
	ctx := chi.NewRouteContext()
	if r.Match(ctx, method, path) {
		t.Error(method, "route", path, "should not match")
	}
}

type testCase struct {
	method   string
	route    string
	isNil    bool
	funcName string
	typeName string
}

func TestHandlers(t *testing.T) {
	tree, err := getTree(Router("/"))
	if err != nil {
		t.Error(err)
	}
	cases := []testCase{
		testCase{
			method:   "GET",
			route:    "/js/*",
			typeName: "net/http.HandlerFunc",
		},
		testCase{
			method:   "GET",
			route:    "/static/*",
			typeName: "net/http.HandlerFunc",
		},
		testCase{
			method:   "GET",
			route:    "/focus/",
			funcName: "github.com/mbertschler/bunny/pkg/router.viewFocusPage",
		},
		testCase{
			method:   "GET",
			route:    "/item/{id}",
			funcName: "github.com/mbertschler/bunny/pkg/router.viewItemPage",
		},
		testCase{
			method:   "GET",
			route:    "/list/{id}",
			funcName: "github.com/mbertschler/bunny/pkg/router.viewListPage",
		},
		testCase{
			method:   "GET",
			route:    "/",
			funcName: "github.com/mbertschler/bunny/pkg/router.viewAreaPage",
		},
		testCase{
			method:   "POST",
			route:    "/gui/",
			typeName: "github.com/mbertschler/bunny/pkg/guiapi.Handler",
		},
	}
	for _, tc := range cases {
		testOneCase(t, tree, tc)
	}
}

func testOneCase(t *testing.T, m handlerMap, tc testCase) {
	h := m[methodRoute{method: tc.method, route: tc.route}]
	if h == nil {
		if tc.isNil {
			return
		}
		t.Error(tc, "shouldn't be nil")
		return
	}

	if tc.funcName != "" {
		fn := functionName(h)
		if fn == tc.funcName {
			return
		}
		t.Error(tc, "wrong function name", fn)
		return
	}

	if tc.typeName != "" {
		typ := reflect.ValueOf(h).Type()
		name := typ.PkgPath() + "." + typ.Name()
		if name == tc.typeName {
			return
		}
		t.Error(tc, "wrong type name", name)
		return
	}
	t.Error(tc, "no funcName or typeName specified")
}

func functionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type methodRoute struct {
	method string
	route  string
}
type handlerMap map[methodRoute]http.Handler

func getTree(r *chi.Mux) (handlerMap, error) {
	m := handlerMap{}
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// HACK: Walk puts /*/ into routes for subroutes that
		// are mounted using Mount()
		route = strings.Replace(route, "/*/", "/", -1)
		m[methodRoute{method: method, route: route}] = handler
		return nil
	}
	err := chi.Walk(r, walkFunc)
	return m, err
}
