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
	"errors"
	"fmt"
	"io"

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

	db.SetList(storedList(List{
		ID: 1,
	}))

	db.SortListItemAfter(1, 1, 0)
	db.SortListItemAfter(1, 2, 1)
	db.SortListItemAfter(1, 3, 2)
	db.SortListItemAfter(1, 5, 3)
	db.SortListItemAfter(1, 4, 5)

	db.SetUserFocus(1, 1, int(FocusNow))
	db.SetUserFocus(1, 5, int(FocusPause))
	db.SetUserFocus(1, 2, int(FocusLater))
	db.SetUserFocus(1, 3, int(FocusWatch))
}

type focusList struct {
	Focussed bool
	Focus    int
	Pause    []int
	Later    []int
	Watch    []int
	Index    map[int]FocusState
}

type FocusData struct {
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

type List struct {
	ID    int
	State ItemState
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
	i := restoreItem(db.ItemByID(id))
	return i
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

func storedList(in List) stored.List {
	return stored.List{
		ID:    in.ID,
		State: int(in.State),
		Title: in.Title,
		Body:  in.Body,
	}
}

func restoreList(in stored.List) List {
	return List{
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

func SortItem(listID, itemID, after int) {
	db.SortListItemAfter(listID, itemID, after)
}

func SortFocusItem(user, id, after int) {
	db.SortUserFocusAfter(user, id, after)
}

// TODO, solve with a loop and benchmark solutions
func sortArrayOld(arr []int, old, new int) ([]int, error) {
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

func sortArray(in []int, old, new int) ([]int, error) {
	max := len(in)
	if !(old < max && old >= 0 &&
		new < max && new >= 0) {
		return in, errors.New(fmt.Sprintln(
			"invalid sorting from", old, "to", new, "max", max))
	}
	out := make([]int, len(in))
	i, j := 0, 0
	for j < max {
		if j == new {
			out[j] = in[old]
			j++
			continue
		}
		if i != old {
			out[j] = in[i]
			j++
		}
		i++
	}
	return out, nil
}

func SetFocus(user, id int, focus FocusState) {
	db.SetUserFocus(user, id, int(focus))
}

func FocusByUserItem(user, item int) FocusState {
	// db.UserFocus(user, id, int(focus))
	return 0
}

func DeleteItem(id int) {
	db.DeleteItem(id)
}

func ItemList() []Item {
	var out []Item
	for _, i := range db.ItemList(1) {
		out = append(out, restoreItem(i))
	}
	return out
}

func FocusList() FocusData {
	var out FocusData
	for _, i := range db.FocusList(1) {
		switch FocusState(i.Focus) {
		case FocusNow:
			item := restoreItem(i)
			out.Focus = &item
		case FocusPause:
			out.Pause = append(out.Pause, restoreItem(i))
		case FocusLater:
			out.Later = append(out.Later, restoreItem(i))
		case FocusWatch:
			out.Watch = append(out.Watch, restoreItem(i))
		}
	}
	return out
}

func WriteDebugData(w io.Writer) {
	// dataLock.RLock()
	// var data = struct {
	// 	List  []int
	// 	Focus focusList
	// 	Items map[int]Item
	// }{
	// 	List:  dataList,
	// 	Focus: dataFocus,
	// 	Items: dataItems,
	// }
	// w.Write([]byte("<html><body><pre>"))
	// enc := json.NewEncoder(w)
	// enc.SetIndent("", "    ")
	// err := enc.Encode(data)
	// w.Write([]byte("</pre></body></html>"))
	// if err != nil {
	// 	log.Println(err)
	// }
	// dataLock.RUnlock()
}
