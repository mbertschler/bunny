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
		List: []int{1, 2, 3},
	}
	dataMaxID = 3
	dataItems = map[int]itemData{
		1: itemData{
			ID:       1,
			Closed:   false,
			Archived: false,
			Title:    "Hello world!",
			Body:     "Let's have some fun with bunny!",
		},
		2: itemData{
			ID:       2,
			Closed:   true,
			Archived: false,
			Title:    "Look at Bunny",
			Body:     "By reading this text you alredy completed this item.",
		},
		3: itemData{
			ID:       3,
			Closed:   true,
			Archived: true,
			Title:    "Nevermind me",
			Body:     "I am done and no longer relevant, so I got archived.",
		},
	}
)

type itemData struct {
	ID       int
	Closed   bool
	Archived bool
	Title    string
	Body     string
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
