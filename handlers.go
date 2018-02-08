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
			"listView":   listViewHandler,
			"listSort":   listSortHandler,
			"itemNew":    itemNewHandler,
			"itemView":   itemViewHandler,
			"itemEdit":   itemEditHandler,
			"itemSave":   itemSaveHandler,
			"itemState":  itemStateHandler,
			"itemFocus":  itemFocusHandler,
			"itemDelete": itemDeleteHandler,
			"focusView":  focusViewHandler,
			"focusSort":  focusSortHandler,
		},
	}
	return handler
}

func listViewHandler(_ json.RawMessage) (*Result, error) {
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

func focusViewHandler(_ json.RawMessage) (*Result, error) {
	res, err := replaceContainer(displayFocusBlock())
	if res != nil {
		args, err := json.Marshal([]interface{}{nil, "Bunny Focus", "/focus/"})
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

func listSortHandler(in json.RawMessage) (*Result, error) {
	var args = struct {
		Old int
		New int
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	sortItem(args.Old, args.New)
	return listViewHandler(nil)
}

func focusSortHandler(in json.RawMessage) (*Result, error) {
	var args = struct {
		Old int
		New int
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	sortFocusItem(args.Old, args.New)
	return focusViewHandler(nil)
}

func itemNewHandler(in json.RawMessage) (*Result, error) {
	item := newItem()
	return replaceContainer(editItemBlock(item))
}

func itemViewHandler(in json.RawMessage) (*Result, error) {
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

func itemEditHandler(in json.RawMessage) (*Result, error) {
	var id int
	err := json.Unmarshal(in, &id)
	if err != nil {
		return nil, err
	}
	return replaceContainer(editItemBlock(getItemData(id)))
}

func itemSaveHandler(in json.RawMessage) (*Result, error) {
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

func itemStateHandler(in json.RawMessage) (*Result, error) {
	var args = struct {
		ID    int
		State string
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	data := getItemData(args.ID)
	switch args.State {
	case "open":
		data.State = ItemOpen
	case "complete":
		data.State = ItemComplete
		data.Focus = FocusNone
	case "archived":
		data.State = ItemArchived
		data.Focus = FocusNone
	}
	setItemData(args.ID, data)
	return replaceContainer(displayItemBlock(data))
}

func itemFocusHandler(in json.RawMessage) (*Result, error) {
	var args = struct {
		ID    int
		Focus string
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	data := focusItem(args.ID, args.Focus)
	return replaceContainer(displayItemBlock(data))
}

func itemDeleteHandler(in json.RawMessage) (*Result, error) {
	var arg int
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	deleteItem(arg)
	return replaceContainer(displayListBlock(getListData()))
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
