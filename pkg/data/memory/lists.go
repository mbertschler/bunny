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

type listsTx struct {
	parent *Tx
	tx     *buntdb.Tx
}

func (t *listsTx) Key(id int) string {
	return listPrefix + strconv.Itoa(id)
}

func (t *listsTx) ID(key string) int {
	key = strings.TrimPrefix(key, listPrefix)
	i, err := strconv.Atoi(key)
	if err != nil {
		log.Println("KEY ERROR:", err)
	}
	return i
}

func (t *listsTx) Get(id int) (stored.List, error) {
	var list stored.List
	val, err := t.tx.Get(t.Key(id))
	if err != nil {
		return list, err
	}
	err = decode(val, &list)
	return list, err
}

func (t *listsTx) UserItems(user, list int) ([]stored.Item, error) {
	_, err := t.parent.users.Get(user)
	if err != nil {
		return nil, err
	}
	l, err := t.Get(list)
	if err != nil {
		return nil, err
	}
	var out []stored.Item
	for _, id := range l.Items {
		item, err := t.parent.items.UserItem(user, id)
		if err != nil {
			log.Println("oh no, error in a loop :(", user, id, err)
		}
		out = append(out, item)
	}
	return out, err
}

func (t *listsTx) Items(list int) ([]stored.Item, error) {
	l, err := t.Get(list)
	if err != nil {
		return nil, err
	}
	var out []stored.Item
	for _, id := range l.Items {
		item, err := t.parent.items.Get(id)
		if err != nil {
			log.Println("oh no, error in a loop :(", id, err)
		}
		out = append(out, item)
	}
	return out, err
}

func (t *listsTx) Set(l stored.List) error {
	val, err := encode(l)
	if err != nil {
		return err
	}
	_, _, err = t.tx.Set(t.Key(l.ID), val, nil)
	return err
}

func (t *listsTx) SetItemPos(list, item, pos int) error {
	// TODO: remove from other lists
	l, err := t.Get(list)
	if err != nil {
		return err
	}
	i, ok := findInArray(l.Items, item)
	if !ok {
		i = len(l.Items)
		l.Items = append(l.Items, item)
	}
	l.Items, err = sortArray(l.Items, i, pos-1) // 0 indexed not 1
	if err != nil {
		return err
	}
	return t.Set(l)
}

func (t *listsTx) Delete(id int) error {
	return nil
}
