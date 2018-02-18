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

	"github.com/mbertschler/bunny/pkg/data"
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
	res, err := replaceContainer(displayListBlock(data.ItemList()))
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
	data.SortItem(1, args.Old, args.New)
	return listViewHandler(nil)
}

func focusSortHandler(in json.RawMessage) (*Result, error) {
	var args = struct {
		Type string
		Old  int
		New  int
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	// var type
	// switch args.Type {
	// case "pause":
	// 	data.State = ItemOpen
	// case "later":
	// 	data.State = ItemComplete
	// 	data.Focus = FocusNone
	// case "watch":
	// 	data.State = ItemArchived
	// 	data.Focus = FocusNone
	// }
	// sortFocusItem(0, args.Old, args.New)
	return focusViewHandler(nil)
}

func itemNewHandler(in json.RawMessage) (*Result, error) {
	return replaceContainer(editItemBlock(data.Item{}, true))
}

func itemViewHandler(in json.RawMessage) (*Result, error) {
	var id int
	err := json.Unmarshal(in, &id)
	if err != nil {
		return nil, err
	}
	res, err := replaceContainer(displayItemBlock(data.ItemByID(id)))
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
	return replaceContainer(editItemBlock(data.ItemByID(id), false))
}

func itemSaveHandler(in json.RawMessage) (*Result, error) {
	var arg struct {
		ID    int
		New   bool
		Title string
		Body  string
	}
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	if arg.New {
		arg.ID = data.NewItem().ID
	}
	d := data.ItemByID(arg.ID)
	if len(arg.Title) > 0 {
		d.Title = arg.Title
	}
	if len(arg.Body) > 0 {
		d.Body = arg.Body
	}
	data.SetItem(d)
	return replaceContainer(displayItemBlock(d))
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
	d := data.ItemByID(args.ID)
	switch args.State {
	case "open":
		d.State = data.ItemOpen
	case "complete":
		d.State = data.ItemComplete
		d.Focus = data.FocusNone
	case "archived":
		d.State = data.ItemArchived
		d.Focus = data.FocusNone
	}
	data.SetItem(d)
	return replaceContainer(displayItemBlock(d))
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
	data.SetFocus(1, args.ID, 1)
	d := data.ItemByID(args.ID)
	return replaceContainer(displayItemBlock(d))
}

func itemDeleteHandler(in json.RawMessage) (*Result, error) {
	var arg int
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data.DeleteItem(arg)
	return replaceContainer(displayListBlock(data.ItemList()))
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
