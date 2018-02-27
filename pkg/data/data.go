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

//go:generate stringer -type=ItemState,FocusState -output=states_string.go

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
	forceSetUser(User{
		ID:   1,
		Name: "martin",
	})
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

	forceSetList(List{
		ID: 1,
	})

	SetListItemPosition(1, 1, 1)
	SetListItemPosition(1, 2, 2)
	SetListItemPosition(1, 3, 3)
	SetListItemPosition(1, 4, 4)
	SetListItemPosition(1, 5, 5)

	SetUserFocus(1, 1, FocusNow)
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

type User struct {
	ID   int
	Name string
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
	Items []Item
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

func storedUser(in User) stored.User {
	return stored.User{
		ID:   in.ID,
		Name: in.Name,
	}
}

func restoreUser(in stored.User) User {
	return User{
		ID:   in.ID,
		Name: in.Name,
	}
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
	if err != nil {
		return i, err
	}
	err = db.SetListItemPosition(1, i.ID, 1)
	return i, err
}

func SortFocusItem(user, id, after int) error {
	return db.SortUserFocusAfter(user, id, after)
}

func SetFocus(user, id int, focus FocusState) error {
	return db.SetUserFocus(user, id, int(focus))
}

func DeleteItem(id int) error {
	return db.DeleteItem(id)
}

func ItemList(id int) ([]Item, error) {
	var out []Item
	_, items, err := db.ItemList(id)
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		out = append(out, restoreItem(i))
	}
	return out, nil
}

func UserItemList(user, id int) ([]Item, error) {
	var out []Item
	_, items, err := db.UserItemList(user, id)
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		out = append(out, restoreItem(i))
	}
	return out, nil
}

func FocusList(user int) (FocusData, error) {
	var out FocusData
	list, err := db.FocusList(user)
	if err != nil {
		return out, err
	}
	for _, i := range list {
		switch FocusState(i.Focus) {
		case FocusNow:
			item := restoreItem(i)
			out.Focus = &item
		case FocusLater:
			out.Later = append(out.Later, restoreItem(i))
		case FocusWatch:
			out.Watch = append(out.Watch, restoreItem(i))
		}
	}
	return out, nil
}

func SetListItemPosition(list, item, pos int) error {
	return db.SetListItemPosition(list, item, pos)
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

func forceSetList(in List) error {
	return db.ForceSetList(storedList(in))
}

func forceSetUser(in User) error {
	return db.ForceSetUser(storedUser(in))
}

func debugItemList(list int) ([]stored.OrderedListItem, error) {
	return db.DebugItemList(list)
}

func UserByID(id int) (User, error) {
	stored, err := db.UserByID(id)
	i := restoreUser(stored)
	return i, err
}
