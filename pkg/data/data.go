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
	"log"

	"github.com/mbertschler/bunny/pkg/data/memory"
	"github.com/mbertschler/bunny/pkg/data/stored"
)

var db *memory.DB

func init() {
	db = memory.Open()
	setupTestdata()
}

func setupTestdata() {
	logErr(forceSetUser(User{
		ID:   1,
		Name: "martin",
	}))
	logErr(forceSetItem(Item{
		ID:    1,
		State: ItemOpen,
		Title: "Hello world!",
		Body:  "Let's have some fun with bunny!",
	}))
	logErr(forceSetItem(Item{
		ID:    2,
		State: ItemComplete,
		Title: "Look at Bunny",
		Body:  "By reading this text you alredy completed this item.",
	}))
	logErr(forceSetItem(Item{
		ID:    3,
		State: ItemOpen,
		Title: "Somebody else does it",
		Body:  "This is something that I am interested in. On the other hand I don't intend to work on it.",
	}))
	logErr(forceSetItem(Item{
		ID:    4,
		State: ItemArchived,
		Title: "Nevermind me, I'm old",
		Body:  "I am done and no longer relevant, so I got archived.",
	}))
	logErr(forceSetItem(Item{
		ID:    5,
		State: ItemOpen,
		Title: "I started it but don't know how to finish",
		Body:  "Somebody please help me so that I can complete this item.",
	}))
	logErr(forceSetItem(Item{
		ID:    6,
		State: ItemComplete,
		Title: "Create areas",
		Body:  "A single list is not enough, we need areas that hold many lists and items.",
	}))
	logErr(forceSetItem(Item{
		ID:    7,
		State: ItemOpen,
		Title: "This item clearly doesn't belong into a list",
		Body:  "What to do with an item if you don't know how to categorize it?",
	}))

	logErr(forceSetList(List{
		ID:    1,
		Title: "Testlist",
		Body:  "just for testing",
	}))

	logErr(SetListItemPosition(1, 1, 1))
	logErr(SetListItemPosition(1, 2, 2))
	logErr(SetListItemPosition(1, 3, 3))
	logErr(SetListItemPosition(1, 4, 4))
	logErr(SetListItemPosition(1, 5, 5))

	logErr(forceSetArea(Area{
		ID:    1,
		Title: "Starting area",
		Body:  "area that keeps inital data",
	}))

	logErr(SetAreaThingPosition(1, TypeItem, 7, 1))
	logErr(SetAreaThingPosition(1, TypeList, 1, 2))
	logErr(SetAreaThingPosition(1, TypeItem, 6, 3))

	logErr(SetUserFocus(1, 1, FocusNow))
	logErr(SetUserFocus(1, 2, FocusLater))
	logErr(SetUserFocus(1, 3, FocusWatch))
	logErr(SetUserFocus(1, 7, FocusLater))
}

func logErr(err error) {
	if err != nil {
		log.Println("setup error:", err)
	}
}

type FocusData struct {
	Focus []Item
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

func (Item) thingType() ThingType { return TypeItem }

func (i Item) Archived() bool {
	return i.State == ItemArchived
}

type List struct {
	ID    int
	State ItemState
	Title string
	Body  string
	Items []Item
}

func (List) thingType() ThingType { return TypeList }

func (l List) Archived() bool {
	return l.State == ItemArchived
}

type Thing interface {
	Archived() bool
	thingType() ThingType
}

type ThingType int8

const (
	TypeItem ThingType = stored.TypeItem
	TypeList ThingType = stored.TypeList
)

type Area struct {
	ID    int
	Title string
	Body  string
	List  []Thing
}

type ItemState int8

const (
	ItemOpen     ItemState = stored.ItemOpen
	ItemComplete ItemState = stored.ItemComplete
	ItemArchived ItemState = stored.ItemArchived
)

type FocusState int8

const (
	FocusNone  FocusState = stored.FocusNone
	FocusNow   FocusState = stored.FocusNow
	FocusLater FocusState = stored.FocusLater
	FocusWatch FocusState = stored.FocusWatch
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

func storedArea(in Area) stored.Area {
	return stored.Area{
		ID:    in.ID,
		Title: in.Title,
		Body:  in.Body,
	}
}

func restoreArea(in stored.Area) Area {
	return Area{
		ID:    in.ID,
		Title: in.Title,
		Body:  in.Body,
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

func storedThing(in Thing) stored.Thing {
	switch in.thingType() {
	case TypeList:
		return storedList(in.(List))
	case TypeItem:
		return storedItem(in.(Item))
	}
	return nil
}

func restoreThing(in stored.Thing) Thing {
	switch in.Type() {
	case stored.TypeList:
		return restoreList(in.(stored.List))
	case stored.TypeItem:
		return restoreItem(in.(stored.Item))
	}
	return nil
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

func UserArea(user, id int) (Area, []Thing, error) {
	var out []Thing
	area, items, err := db.UserArea(user, id)
	if err != nil {
		return restoreArea(area), nil, err
	}
	for _, i := range items {
		out = append(out, restoreThing(i))
	}
	return restoreArea(area), out, nil
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
			out.Focus = append(out.Focus, restoreItem(i))
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

func SetAreaThingPosition(area int, typ ThingType, id, pos int) error {
	return db.SetAreaThingPosition(area, stored.ThingType(typ), id, pos)
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

func forceSetArea(in Area) error {
	return db.ForceSetArea(storedArea(in))
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

func AreaByID(id int) (Area, error) {
	stored, err := db.AreaByID(id)
	a := restoreArea(stored)
	return a, err
}
