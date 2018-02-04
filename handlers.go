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
	"encoding/json"
	"fmt"
	"log"

	"github.com/mbertschler/blocks/html"
)

func guiAPI() Handler {
	handler := Handler{
		Functions: map[string]Callable{
			"hello":            helloHandler,
			"viewList":         viewListHandler,
			"sortList":         sortListHandler,
			"newItem":          newItemHandler,
			"viewItem":         viewItemHandler,
			"editItem":         editItemHandler,
			"saveItem":         saveItemHandler,
			"editItemClosed":   editItemClosedHandler,
			"editItemArchived": editItemArchivedHandler,
		},
	}
	return handler
}

func viewListHandler(_ json.RawMessage) (*Result, error) {
	res, err := replaceContainer(displayListBlock(getListData()))
	if res != nil {
		args, err := json.Marshal([]interface{}{nil, "Bunny List", "/"})
		if err != nil {
			log.Println(err)
		}
		// TODO make Result API nicer
		res.JS = append(res.JS, JSCall{
			Name:      "setURL",
			Arguments: args,
		})
		res.JS = append(res.JS, JSCall{
			Name: "enableSorting",
		})
	}
	return res, err
}

func sortListHandler(in json.RawMessage) (*Result, error) {
	var args = struct {
		Old int
		New int
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	sortItem(args.Old, args.New)
	return viewListHandler(nil)
}

func newItemHandler(in json.RawMessage) (*Result, error) {
	item := newItem()
	return replaceContainer(editItemBlock(item))
}

func viewItemHandler(in json.RawMessage) (*Result, error) {
	var id int
	err := json.Unmarshal(in, &id)
	if err != nil {
		return nil, err
	}
	res, err := replaceContainer(displayItemBlock(getItemData(id)))
	if res != nil {
		args, err := json.Marshal([]interface{}{nil, "Bunny Item", fmt.Sprint("/item/", id)})
		if err != nil {
			log.Println(err)
		}
		res.JS = append(res.JS, JSCall{
			Name:      "setURL",
			Arguments: args,
		})
	}
	return res, err
}

func editItemHandler(in json.RawMessage) (*Result, error) {
	var id int
	err := json.Unmarshal(in, &id)
	if err != nil {
		return nil, err
	}
	return replaceContainer(editItemBlock(getItemData(id)))
}

func saveItemHandler(in json.RawMessage) (*Result, error) {
	var arg itemData
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data := getItemData(arg.ID)
	if len(arg.Title) > 0 {
		data.Title = arg.Title
	}
	if len(arg.Body) > 0 {
		data.Body = arg.Body
	}
	setItemData(arg.ID, data)
	return replaceContainer(displayItemBlock(data))
}

func editItemClosedHandler(in json.RawMessage) (*Result, error) {
	var arg itemData
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data := getItemData(arg.ID)
	data.Closed = arg.Closed
	if !arg.Closed {
		data.Archived = false
	}
	setItemData(arg.ID, data)
	return replaceContainer(displayItemBlock(data))
}

func editItemArchivedHandler(in json.RawMessage) (*Result, error) {
	var arg itemData
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data := getItemData(arg.ID)
	data.Archived = arg.Archived
	setItemData(arg.ID, data)
	return replaceContainer(displayItemBlock(data))
}

func helloHandler(in json.RawMessage) (*Result, error) {
	type argType struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	args := argType{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	block := html.H1(nil, html.Text("Hello "+args.Name))
	return replaceContainer(block)
}

func replaceContainer(block html.Block) (*Result, error) {
	out, err := html.RenderString(block)
	if err != nil {
		return nil, err
	}
	ret := &Result{
		HTML: []HTMLUpdate{
			{
				Operation: HTMLReplace,
				Selector:  "#container",
				Content:   out,
			},
		},
	}
	return ret, nil
}
