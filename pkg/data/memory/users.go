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

package memory

import (
	"log"
	"strconv"
	"strings"

	"github.com/mbertschler/bunny/pkg/data/stored"
	"github.com/tidwall/buntdb"
)

type usersTx struct {
	tx     *buntdb.Tx
	parent *Tx
}

func (t *usersTx) Key(user int) string {
	return userPrefix + strconv.Itoa(user)
}

func (t *usersTx) ID(key string) (user int) {
	key = strings.TrimPrefix(key, userPrefix)
	user, err := strconv.Atoi(key)
	if err != nil {
		log.Println("KEY ERROR:", err)
	}
	return user
}

func (t *usersTx) Get(id int) (stored.User, error) {
	var user stored.User
	val, err := t.tx.Get(t.Key(id))
	if err != nil {
		return user, err
	}
	err = decode(val, &user)
	return user, err
}

func (t *usersTx) Set(user stored.User) error {
	val, err := encode(user)
	if err != nil {
		return err
	}
	_, _, err = t.tx.Set(t.Key(user.ID), val, nil)
	return err
}

func (t *usersTx) ItemFocus(user, item int) (int, error) {
	u, err := t.Get(user)
	if err != nil {
		return 0, err
	}
	focus, _ := findItemInFocusmap(u.Focus, item)
	return focus, nil
}

func (t *usersTx) SetFocus(user, item, focus int) error {
	u, err := t.Get(user)
	if err != nil {
		return err
	}
	if u.Focus == nil {
		u.Focus = make(map[int][]int)
	}
	oldFocus, index := findItemInFocusmap(u.Focus, item)
	if oldFocus != 0 {
		u.Focus[oldFocus] = deleteFromArray(u.Focus[oldFocus], index)
	}
	if focus == stored.FocusNow && len(u.Focus[focus]) > 0 {
		err = t.SetFocus(user, u.Focus[focus][0], stored.FocusLater)
		if err != nil {
			return err
		}
		u, err = t.Get(user)
		if err != nil {
			return err
		}
	}
	u.Focus[focus] = append(u.Focus[focus], item)
	return t.Set(u)
}

func findItemInFocusmap(m map[int][]int, id int) (focus, index int) {
	for _, focus := range []int{1, 2, 3} {
		for i, focusID := range m[focus] {
			if id == focusID {
				return focus, i
			}
		}
	}
	return 0, 0
}

func (t *usersTx) AllByUser(user int) ([]stored.Item, error) {
	var items []stored.Item
	u, err := t.Get(user)
	if err != nil {
		return items, err
	}
	for _, focus := range []int{1, 2, 3} {
		for _, focusID := range u.Focus[focus] {
			i, err := t.parent.items.Get(focusID)
			if err != nil {
				log.Println("oh no, inner loop err :(", i)
			}
			i.Focus = focus
			items = append(items, i)
		}
	}
	return items, nil
}
