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
	"sync"
)

var (
	dataLock sync.RWMutex
	dataList = listData{
		List: []int{1, 2, 3, 4},
	}
	dataMaxID = 4
	dataItems = map[int]itemData{
		1: itemData{
			ID:       1,
			Complete: false,
			Archived: false,
			Title:    "Hello world!",
			Body:     "Let's have some fun with bunny!",
			Focus:    true,
		},
		2: itemData{
			ID:       2,
			Complete: true,
			Archived: false,
			Title:    "Look at Bunny",
			Body:     "By reading this text you alredy completed this item.",
			Later:    true,
		},
		3: itemData{
			ID:       3,
			Complete: false,
			Archived: false,
			Title:    "Somebody else does it",
			Body:     "This is something that I am interested in. On the other hand I don't intend to work on it.",
			Watch:    true,
		},
		4: itemData{
			ID:       4,
			Complete: true,
			Archived: true,
			Title:    "Nevermind me, I'm old",
			Body:     "I am done and no longer relevant, so I got archived.",
		},
	}
)

type itemData struct {
	ID       int
	Complete bool
	Archived bool
	Title    string
	Body     string
	Focus    bool
	Later    bool
	Watch    bool
}

type listData struct {
	List []int
}

func getItemData(id int) itemData {
	dataLock.RLock()
	d := dataItems[id]
	dataLock.RUnlock()
	return d
}

func setItemData(id int, in itemData) {
	dataLock.Lock()
	_, ok := dataItems[id]
	if !ok {
		log.Println("setting nonexistent item?", id, in)
	}
	dataItems[id] = in
	dataLock.Unlock()
}

func newItem() itemData {
	item := itemData{}
	dataLock.Lock()
	dataMaxID++
	item.ID = dataMaxID
	dataItems[dataMaxID] = item
	dataList.List = append(dataList.List, dataMaxID)
	dataLock.Unlock()
	return item
}

func sortItem(old, new int) {
	dataLock.Lock()
	max := len(dataList.List)
	if old < max && old >= 0 &&
		new < max && new >= 0 {
		dataList.List = sortArray(dataList.List, old, new)
	} else {
		log.Println("invalid sorting attempt, from", old, "to", new, "max", max)
	}
	dataLock.Unlock()
}

// TODO, solve with a loop and benchmark solutions
func sortArray(arr []int, old, new int) []int {
	el := arr[old]
	arr = append(arr[:old], arr[old+1:]...)
	in := make([]int, len(arr))
	copy(in, arr)
	return append(append(arr[:new], el), in[new:]...)
}

func focusItem(id int, status string) itemData {
	dataLock.Lock()
	item := dataItems[id]
	switch status {
	case "later":
		if item.Later {
			item.Later = false
		} else {
			item.Later = true
			item.Focus = false
			item.Watch = false
		}
	case "focus":
		if item.Focus {
			item.Focus = false
		} else {
			item.Focus = true
			item.Later = false
			item.Watch = false
		}
	case "watch":
		if item.Watch {
			item.Watch = false
		} else {
			item.Watch = true
			item.Focus = false
			item.Later = false
		}
	}
	dataItems[id] = item
	dataLock.Unlock()
	return item
}

func deleteItem(id int) {
	dataLock.Lock()
	newList := make([]int, 0, len(dataList.List)-1)
	for _, e := range dataList.List {
		if e != id {
			newList = append(newList, e)
		}
	}
	dataList.List = newList
	delete(dataItems, id)
	dataLock.Unlock()
}

func getListData() []itemData {
	dataLock.RLock()
	out := make([]itemData, len(dataList.List))
	for i, id := range dataList.List {
		out[i] = dataItems[id]
	}
	dataLock.RUnlock()
	return out
}

func setListData(in []int) {
	dataLock.Lock()
	dataList.List = in
	dataLock.Unlock()
}
