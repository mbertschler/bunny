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

type itemsTx struct {
	parent *Tx
	tx     *buntdb.Tx
}

func (t *itemsTx) Key(id int) string {
	return itemPrefix + strconv.Itoa(id)
}

func (t *itemsTx) ID(key string) int {
	key = strings.TrimPrefix(key, itemPrefix)
	i, err := strconv.Atoi(key)
	if err != nil {
		log.Println("KEY ERROR:", err)
	}
	return i
}

func (t *itemsTx) Get(id int) (stored.Item, error) {
	var item stored.Item
	val, err := t.tx.Get(t.Key(id))
	if err != nil {
		return item, err
	}
	err = decode(val, &item)
	return item, err
}

func (t *itemsTx) UserItem(user, item int) (stored.Item, error) {
	i, err := t.Get(item)
	if err != nil {
		return i, err
	}
	focus, err := t.parent.users.ItemFocus(user, item)
	i.Focus = focus
	return i, err
}

func (t *itemsTx) Set(i stored.Item) error {
	val, err := encode(i)
	if err != nil {
		return err
	}
	_, _, err = t.tx.Set(t.Key(i.ID), val, nil)
	return err
}

func (t *itemsTx) New(i stored.Item) (int, error) {
	id := 0
	err := t.tx.DescendKeys(itemPrefix+"*",
		func(key, val string) bool {
			id = t.ID(key)
			return false
		})
	id++
	i.ID = id
	err = t.Set(i)
	return id, err
}

func (t *itemsTx) Delete(id int) error {
	_, err := t.tx.Delete(t.Key(id))
	return err
}
