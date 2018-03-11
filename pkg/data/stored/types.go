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

package stored

type Cause int8

const (
	CauseNotFound Cause = iota + 1
	CauseMalformed
	CauseSerialize
)

type CauseError struct {
	Cause Cause
	Err   error
}

func (e CauseError) Error() string {
	return e.Err.Error()
}

func WithCause(err error, cause Cause) CauseError {
	return CauseError{
		Cause: cause,
		Err:   err,
	}
}

type MultiError struct {
	Errors []error
}

func (e MultiError) Error() string {
	out := "MultiError: "
	for _, e := range e.Errors {
		out += e.Error()
	}
	return out
}

func CombineErr(err ...error) error {
	return MultiError{
		Errors: err,
	}
}

const (
	ItemOpen = iota
	ItemComplete
	ItemArchived
)

const (
	FocusNone = iota
	FocusNow
	FocusLater
	FocusWatch
)

type Item struct {
	ID    int
	State int
	Title string
	Body  string

	// foreign fields
	Focus int
}

func (Item) Type() ThingType { return TypeItem }

type List struct {
	ID    int
	State int
	Title string
	Body  string

	// internal stored fields
	Items []int
}

func (List) Type() ThingType { return TypeList }

type ThingType int8

const (
	TypeItem = iota + 1
	TypeList
)

type ThingID struct {
	Type ThingType
	ID   int
}

type Thing interface {
	Type() ThingType
}

type Area struct {
	ID     int
	Title  string
	Body   string
	Things []ThingID
}

type User struct {
	ID   int
	Name string

	// internal stored fields
	Focus map[int][]int
}

type OrderedListItem struct {
	Position int
	Item
}

/*
Linked List Sorting
===================

Key: li/listID/itemID

li/1/3 {
	ListID: 1
	ItemID: 3
	Next: 8
}

li/1/8 {
	ListID: 1
	ItemID: 8
	Next: 6
}

li/1/6 {
	ListID: 1
	ItemID: 6
	Next: 0
}

steps to insert an item
-----------------------
- afterItem = get(afterID)
- nextID = afterItem.Next
- afterItem.Next = insertID
- set(afterItem)
- insertItem.Next = nextID
- set(insertItem)

*/

/*
Ordered List
============

Key: li/listID/index

li/1/1 {
	ItemID: 3
}

li/1/2 {
	ItemID: 8
}

li/1/3 {
	ItemID: 6
}

*/
