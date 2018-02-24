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
	"github.com/mbertschler/bunny/pkg/data/memory"
	"github.com/mbertschler/bunny/pkg/data/stored"
)

var db *memory.DB

func init() {
	db = memory.Open()
	setupTestdata()
}

func setupTestdata() {
	forceSetItem(Item{
		ID:    1,
		State: ItemOpen,
		Title: "Hello world!",
		Body:  "Let's have some fun with bunny!",
	})
	forceSetItem(Item{
		ID:    2,
		State: ItemComplete,
		Title: "Look at Bunny",
		Body:  "By reading this text you alredy completed this item.",
	})
	forceSetItem(Item{
		ID:    3,
		State: ItemOpen,
		Title: "Somebody else does it",
		Body:  "This is something that I am interested in. On the other hand I don't intend to work on it.",
	})
	forceSetItem(Item{
		ID:    4,
		State: ItemArchived,
		Title: "Nevermind me, I'm old",
		Body:  "I am done and no longer relevant, so I got archived.",
	})
	forceSetItem(Item{
		ID:    5,
		State: ItemOpen,
		Title: "I started it but don't know how to finish",
		Body:  "Somebody please help me so that I can complete this item.",
	})

	SetList(List{
		ID: 1,
	})

	SortListItemAfter(1, 1, 0)
	SortListItemAfter(1, 2, 1)
	SortListItemAfter(1, 3, 2)
	SortListItemAfter(1, 5, 3)
	SortListItemAfter(1, 4, 5)

	SetUserFocus(1, 1, FocusNow)
	SetUserFocus(1, 5, FocusPause)
	SetUserFocus(1, 2, FocusLater)
	SetUserFocus(1, 3, FocusWatch)
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

func ItemByID(id int) (Item, error) {
	stored, err := db.ItemByID(id)
	i := restoreItem(stored)
	return i, err
}

func UserItemByID(user, id int) (Item, error) {
	stored, err := db.UserItemByID(user, id)
	i := restoreItem(stored)
	return i, err
}

func SetItem(in Item) error {
	return db.SetItem(storedItem(in))
}

func forceSetItem(in Item) error {
	return db.ForceSetItem(storedItem(in))
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
		Focus: FocusState(in.Focus),
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

func NewItem() (Item, error) {
	i := Item{}
	var err error
	i.ID, err = db.NewItem(storedItem(i))
	return i, err
}

func SortItem(listID, itemID, after int) {
	db.SortListItemAfter(listID, itemID, after)
}

func SortFocusItem(user, id, after int) {
	db.SortUserFocusAfter(user, id, after)
}

func SetFocus(user, id int, focus FocusState) {
	db.SetUserFocus(user, id, int(focus))
}

func FocusByUserItem(user, item int) FocusState {
	// db.UserFocus(user, id, int(focus))
	return 0
}

func DeleteItem(id int) error {
	return db.DeleteItem(id)
}

func ItemList(id int) []Item {
	var out []Item
	for _, i := range db.ItemList(id) {
		out = append(out, restoreItem(i))
	}
	return out
}

func UserItemList(user, id int) []Item {
	var out []Item
	for _, i := range db.UserItemList(user, id) {
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

func SortListItemAfter(list, item, after int) error {
	return db.SortListItemAfter(list, item, after)
}

func SetUserFocus(user, item int, focus FocusState) error {
	return db.SetUserFocus(user, item, int(focus))
}

func ListByID(id int) (List, error) {
	stored, err := db.ListByID(id)
	i := restoreList(stored)
	return i, err
}

func SetList(in List) error {
	return db.SetList(storedList(in))
}
