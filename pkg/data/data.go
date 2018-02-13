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

package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

var (
	dataLock  sync.RWMutex
	dataList  = []int{1, 2, 3, 5, 4}
	dataFocus = focusList{
		Focussed: true,
		Focus:    1,
		Pause:    []int{5},
		Later:    []int{2},
		Watch:    []int{3},
		Index: map[int]FocusState{
			1: FocusNow,
			2: FocusLater,
			3: FocusWatch,
			5: FocusPause,
		},
	}
	dataMaxID = 5
	dataItems = map[int]Item{
		1: Item{
			ID:    1,
			State: ItemOpen,
			Title: "Hello world!",
			Body:  "Let's have some fun with bunny!",
		},
		2: Item{
			ID:    2,
			State: ItemComplete,
			Title: "Look at Bunny",
			Body:  "By reading this text you alredy completed this item.",
		},
		3: Item{
			ID:    3,
			State: ItemOpen,
			Title: "Somebody else does it",
			Body:  "This is something that I am interested in. On the other hand I don't intend to work on it.",
		},
		4: Item{
			ID:    4,
			State: ItemArchived,
			Title: "Nevermind me, I'm old",
			Body:  "I am done and no longer relevant, so I got archived.",
		},
		5: Item{
			ID:    5,
			State: ItemOpen,
			Title: "I started it but don't know how to finish",
			Body:  "Somebody please help me so that I can complete this item.",
		},
	}
)

type focusList struct {
	Focussed bool
	Focus    int
	Pause    []int
	Later    []int
	Watch    []int
	Index    map[int]FocusState
}

type FocusList struct {
	Focus *Item
	Pause []Item
	Later []Item
	Watch []Item
}

type Item struct {
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

func ItemByID(id int) Item {
	dataLock.RLock()
	d := dataItems[id]
	d.Focus = dataFocus.Index[id]
	dataLock.RUnlock()
	return d
}

func SetItem(in Item) {
	dataLock.Lock()
	_, ok := dataItems[in.ID]
	if !ok {
		log.Println("setting nonexistent item?", in.ID, in)
	}
	dataItems[in.ID] = in
	dataLock.Unlock()
}

func NewItem() Item {
	item := Item{}
	dataLock.Lock()
	dataMaxID++
	item.ID = dataMaxID
	dataItems[dataMaxID] = item
	dataList = append(dataList, dataMaxID)
	dataLock.Unlock()
	return item
}

func SortItem(old, new int) {
	dataLock.Lock()
	var err error
	dataList, err = sortArray(dataList, old, new)
	if err != nil {
		log.Println(err)
	}
	dataLock.Unlock()
}

func sortFocusItem(typ FocusState, old, new int) {
	dataLock.Lock()
	var err error
	switch typ {
	case FocusPause:
		dataFocus.Pause, err = sortArray(dataFocus.Pause, old, new)
	case FocusLater:
		dataFocus.Later, err = sortArray(dataFocus.Later, old, new)
	case FocusWatch:
		dataFocus.Watch, err = sortArray(dataFocus.Watch, old, new)
	}
	if err != nil {
		log.Println(err)
	}
	dataLock.Unlock()
}

// TODO, solve with a loop and benchmark solutions
func sortArray(arr []int, old, new int) ([]int, error) {
	max := len(arr)
	if !(old < max && old >= 0 &&
		new < max && new >= 0) {
		return arr, errors.New(fmt.Sprintln(
			"invalid sorting from", old, "to", new, "max", max))
	}
	el := arr[old]
	arr = append(arr[:old], arr[old+1:]...)
	in := make([]int, len(arr))
	copy(in, arr)
	return append(append(arr[:new], el), in[new:]...), nil
}

func FocusItem(id int, status string) Item {
	// dataLock.Lock()
	// current := dataFocus.Index[id]
	// switch status {
	// case "later":
	// 	if current == FocusLater {
	// 		item.Focus = FocusNone
	// 	} else {
	// 		item.Focus = FocusLater
	// 	}
	// case "focus":
	// 	if item.Focus == FocusNow {
	// 		item.Focus = FocusNone
	// 	} else {
	// 		item.Focus = FocusNow
	// 	}
	// case "watch":
	// 	if item.Focus == FocusWatch {
	// 		item.Focus = FocusNone
	// 	} else {
	// 		item.Focus = FocusWatch
	// 	}
	// }
	// dataLock.Unlock()
	return ItemByID(id)
}

func DeleteItem(id int) {
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

func Items() []Item {
	dataLock.RLock()
	out := make([]Item, len(dataList))
	for i, id := range dataList {
		out[i] = dataItems[id]
		out[i].Focus = dataFocus.Index[id]
	}
	dataLock.RUnlock()
	return out
}

func setListData(in []int) {
	dataLock.Lock()
	dataList = in
	dataLock.Unlock()
}

func Focus() FocusList {
	dataLock.RLock()
	var out FocusList
	if dataFocus.Focussed {
		item := dataItems[dataFocus.Focus]
		out.Focus = &item
	}
	for _, id := range dataFocus.Pause {
		out.Pause = append(out.Pause, dataItems[id])
	}
	for _, id := range dataFocus.Later {
		out.Later = append(out.Later, dataItems[id])
	}
	for _, id := range dataFocus.Watch {
		out.Watch = append(out.Watch, dataItems[id])
	}
	dataLock.RUnlock()
	return out
}

func WriteDebugData(w io.Writer) {
	dataLock.RLock()
	var data = struct {
		List  []int
		Focus focusList
		Items map[int]Item
	}{
		List:  dataList,
		Focus: dataFocus,
		Items: dataItems,
	}
	w.Write([]byte("<html><body><pre>"))
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(data)
	w.Write([]byte("</pre></body></html>"))
	if err != nil {
		log.Println(err)
	}
	dataLock.RUnlock()
}
