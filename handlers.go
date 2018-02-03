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
	"sync"

	"github.com/mbertschler/blocks/html"
)

var (
	dataLock sync.RWMutex
	dataVar  = itemData{
		Closed:   false,
		Archived: false,
		Title:    "Item title",
		Body:     "Item body",
	}
)

type itemData struct {
	Closed   bool
	Archived bool
	Title    string
	Body     string
}

func getItemData() itemData {
	dataLock.RLock()
	d := dataVar
	dataLock.RUnlock()
	return d
}

func setItemData(in itemData) {
	dataLock.Lock()
	dataVar = in
	dataLock.Unlock()
}

func viewItemHandler(in json.RawMessage) (*Result, error) {
	var arg string
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	return replaceContainer(displayBlock(getItemData()))
}

func editItemHandler(in json.RawMessage) (*Result, error) {
	var arg string
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	return replaceContainer(editBlock(getItemData()))
}

func saveItemHandler(in json.RawMessage) (*Result, error) {
	var arg itemData
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data := getItemData()
	if len(arg.Title) > 0 {
		data.Title = arg.Title
	}
	if len(arg.Body) > 0 {
		data.Body = arg.Body
	}
	setItemData(data)
	return replaceContainer(displayBlock(data))
}

func editItemClosedHandler(in json.RawMessage) (*Result, error) {
	var arg bool
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data := getItemData()
	data.Closed = arg
	if !arg {
		data.Archived = false
	}
	setItemData(data)
	return replaceContainer(displayBlock(data))
}

func editItemArchivedHandler(in json.RawMessage) (*Result, error) {
	var arg bool
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data := getItemData()
	data.Archived = arg
	setItemData(data)
	return replaceContainer(displayBlock(data))
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
