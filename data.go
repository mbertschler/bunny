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
	dataLock      sync.RWMutex
	dataList      = []int{1, 2, 3, 5, 4}
	dataFocusList = []int{1, 5, 2, 3}
	dataMaxID     = 4
	dataItems     = map[int]itemData{
		1: itemData{
			ID:    1,
			State: ItemOpen,
			Focus: FocusNow,
			Title: "Hello world!",
			Body:  "Let's have some fun with bunny!",
		},
		2: itemData{
			ID:    2,
			State: ItemComplete,
			Focus: FocusLater,
			Title: "Look at Bunny",
			Body:  "By reading this text you alredy completed this item.",
		},
		3: itemData{
			ID:    3,
			State: ItemOpen,
			Focus: FocusWatch,
			Title: "Somebody else does it",
			Body:  "This is something that I am interested in. On the other hand I don't intend to work on it.",
		},
		4: itemData{
			ID:    4,
			State: ItemArchived,
			Title: "Nevermind me, I'm old",
			Body:  "I am done and no longer relevant, so I got archived.",
		},
		5: itemData{
			ID:    5,
			State: ItemOpen,
			Focus: FocusPause,
			Title: "I started it but don't know how to finish",
			Body:  "Somebody please help me so that I can complete this item.",
		},
	}
)

type itemData struct {
	ID    int
	State ItemState
	Focus FocusState
	Title string
	Body  string
}

type ItemState int8

const (
	ItemOpen ItemState = iota
	ItemComplete
	ItemArchived
)

type FocusState int8

const (
	FocusNone FocusState = iota
	FocusNow
	FocusPause
	FocusLater
	FocusWatch
)

type FocusListData struct {
	ID    int
	State FocusState
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
	dataList = append(dataList, dataMaxID)
	dataLock.Unlock()
	return item
}

func sortItem(old, new int) {
	dataLock.Lock()
	max := len(dataList)
	if old < max && old >= 0 &&
		new < max && new >= 0 {
		dataList = sortArray(dataList, old, new)
	} else {
		log.Println("invalid sorting attempt, from", old, "to", new, "max", max)
	}
	dataLock.Unlock()
}

func sortFocusItem(old, new int) {
	dataLock.Lock()
	max := len(dataFocusList)
	if old < max && old >= 0 &&
		new < max && new >= 0 {
		dataFocusList = sortArray(dataFocusList, old, new)
	} else {
		log.Println("invalid focus sorting attempt, from", old, "to", new, "max", max)
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
		if item.Focus == FocusLater {
			item.Focus = FocusNone
		} else {
			item.Focus = FocusLater
		}
	case "focus":
		if item.Focus == FocusNow {
			item.Focus = FocusNone
		} else {
			item.Focus = FocusNow
		}
	case "watch":
		if item.Focus == FocusWatch {
			item.Focus = FocusNone
		} else {
			item.Focus = FocusWatch
		}
	}
	dataItems[id] = item
	dataLock.Unlock()
	return item
}

func deleteItem(id int) {
	dataLock.Lock()
	newList := make([]int, 0, len(dataList)-1)
	for _, e := range dataList {
		if e != id {
			newList = append(newList, e)
		}
	}
	dataList = newList
	delete(dataItems, id)
	dataLock.Unlock()
}

func getListData() []itemData {
	dataLock.RLock()
	out := make([]itemData, len(dataList))
	for i, id := range dataList {
		out[i] = dataItems[id]
	}
	dataLock.RUnlock()
	return out
}

func setListData(in []int) {
	dataLock.Lock()
	dataList = in
	dataLock.Unlock()
}

func getFocusData() []itemData {
	dataLock.RLock()
	out := make([]itemData, len(dataFocusList))
	for i, id := range dataFocusList {
		out[i] = dataItems[id]
	}
	dataLock.RUnlock()
	return out
}

func setFocusData(in []int) {
	dataLock.Lock()
	dataFocusList = in
	dataLock.Unlock()
}
