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

	"github.com/mbertschler/bunny/pkg/data/memory"
	"github.com/mbertschler/bunny/pkg/data/stored"
)

var db *memory.DB

func init() {
	db = memory.Open()
	db.SetItem(storedItem(Item{
		ID:    1,
		State: ItemOpen,
		Title: "Hello world!",
		Body:  "Let's have some fun with bunny!",
	}))
	db.SetItem(storedItem(Item{
		ID:    2,
		State: ItemComplete,
		Title: "Look at Bunny",
		Body:  "By reading this text you alredy completed this item.",
	}))
	db.SetItem(storedItem(Item{
		ID:    3,
		State: ItemOpen,
		Title: "Somebody else does it",
		Body:  "This is something that I am interested in. On the other hand I don't intend to work on it.",
	}))
	db.SetItem(storedItem(Item{
		ID:    4,
		State: ItemArchived,
		Title: "Nevermind me, I'm old",
		Body:  "I am done and no longer relevant, so I got archived.",
	}))
	db.SetItem(storedItem(Item{
		ID:    5,
		State: ItemOpen,
		Title: "I started it but don't know how to finish",
		Body:  "Somebody please help me so that I can complete this item.",
	}))
}

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
	dataItems = map[int]Item{}
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
	return restoreItem(db.ItemByID(id))
}

func SetItem(in Item) {
	db.SetItem(storedItem(in))
}

func storedItem(in Item) stored.Item {
	return stored.Item{
		ID:    in.ID,
		State: int(in.State),
		Title: in.Title,
		Body:  in.Body,
	}
}

func restoreItem(in stored.Item) Item {
	return Item{
		ID:    in.ID,
		State: ItemState(in.State),
		Title: in.Title,
		Body:  in.Body,
	}
}

func NewItem() Item {
	i := Item{}
	i.ID = db.NewItem(storedItem(i))
	return i
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
	db.DeleteItem(id)
}

func Items() []Item {
	var out []Item
	for _, i := range db.Items() {
		out = append(out, restoreItem(i))
	}
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
