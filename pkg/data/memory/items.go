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
	tx *buntdb.Tx
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

func (t *itemsTx) Set() error {
	return nil
}

func (t *itemsTx) Delete() error {
	return nil
}
