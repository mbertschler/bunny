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

package guiapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mbertschler/blocks/html"

	"github.com/mbertschler/bunny/pkg/blocks"
	"github.com/mbertschler/bunny/pkg/data"
)

func Handlers() Handler {
	handler := Handler{
		Functions: map[string]Callable{
			"areaView":   areaViewHandler,
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

func areaViewHandler(_ json.RawMessage) (*Result, error) {
	_, things, err := data.UserArea(1, 1)
	if err != nil {
		return nil, err
	}
	res, err := replaceContainer(blocks.ViewThingsBlock(things))
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

func listViewHandler(_ json.RawMessage) (*Result, error) {
	list, err := data.UserItemList(1, 1)
	if err != nil {
		return nil, err
	}
	res, err := replaceContainer(blocks.ViewListBlock(list))
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
	focus, err := data.FocusList(1)
	if err != nil {
		log.Println(err)
	}
	res, err := replaceContainer(blocks.ViewFocusBlock(focus))
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
		Item int
		Pos  int
	}{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	data.SetListItemPosition(1, args.Item, args.Pos)
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
	return replaceContainer(blocks.EditItemPage(data.Item{}, true))
}

func itemViewHandler(in json.RawMessage) (*Result, error) {
	var id int
	err := json.Unmarshal(in, &id)
	if err != nil {
		return nil, err
	}
	ui, _ := data.UserItemByID(1, id)
	res, err := replaceContainer(blocks.ViewItemBlock(ui))
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
	ui, _ := data.UserItemByID(1, id)
	return replaceContainer(blocks.EditItemPage(ui, false))
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
		if len(arg.Title) == 0 {
			list, err := data.UserItemList(1, 1)
			if err != nil {
				return nil, err
			}
			return replaceContainer(blocks.ViewListBlock(list))
		}
		newItem, err := data.NewItem()
		if err != nil {
			return nil, err
		}
		arg.ID = newItem.ID
	}
	d, _ := data.UserItemByID(1, arg.ID)
	if len(arg.Title) > 0 {
		d.Title = arg.Title
	}
	if len(arg.Body) > 0 {
		d.Body = arg.Body
	}
	data.SetItem(d)
	return replaceContainer(blocks.ViewItemBlock(d))
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
	d, _ := data.UserItemByID(1, args.ID)
	switch args.State {
	case "open":
		d.State = data.ItemOpen
	case "complete":
		d.State = data.ItemComplete
	case "archived":
		d.State = data.ItemArchived
	}
	data.SetItem(d)
	return replaceContainer(blocks.ViewItemBlock(d))
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

	d, _ := data.UserItemByID(1, args.ID)
	switch args.Focus {
	case "later":
		if d.Focus == data.FocusLater {
			data.SetFocus(1, args.ID, data.FocusNone)
		} else {
			data.SetFocus(1, args.ID, data.FocusLater)
		}
	case "focus":
		if d.Focus == data.FocusNow {
			data.SetFocus(1, args.ID, data.FocusNone)
		} else {
			data.SetFocus(1, args.ID, data.FocusNow)
		}
	case "watch":
		if d.Focus == data.FocusWatch {
			data.SetFocus(1, args.ID, data.FocusNone)
		} else {
			data.SetFocus(1, args.ID, data.FocusWatch)
		}
	}
	d, _ = data.UserItemByID(1, args.ID)
	return replaceContainer(blocks.ViewItemBlock(d))
}

func itemDeleteHandler(in json.RawMessage) (*Result, error) {
	var arg int
	err := json.Unmarshal(in, &arg)
	if err != nil {
		return nil, err
	}
	data.DeleteItem(arg)
	list, err := data.UserItemList(1, 1)
	if err != nil {
		return nil, err
	}
	return replaceContainer(blocks.ViewListBlock(list))
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
