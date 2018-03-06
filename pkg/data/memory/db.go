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
	"encoding/json"
	"log"

	"github.com/mbertschler/bunny/pkg/data/stored"
	"github.com/tidwall/buntdb"
)

// "tables" or "buckets"
const (
	itemPrefix = "i/"
	listPrefix = "l/"
	userPrefix = "u/"
)

func Open() *DB {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		db: db,
	}
}

type DB struct {
	db *buntdb.DB
}

func (d *DB) Existing(tx *buntdb.Tx, writable bool) Tx {
	return makeTx(tx, writable)
}

func (d *DB) View() (Tx, error) {
	tx, err := d.db.Begin(false)
	return makeTx(tx, false), err
}

func (d *DB) Update() (Tx, error) {
	tx, err := d.db.Begin(true)
	return makeTx(tx, true), err
}

func makeTx(tx *buntdb.Tx, writable bool) Tx {
	t := Tx{
		rawTx:    tx,
		writable: writable,
	}
	t.items = itemsTx{tx: tx, parent: &t}
	t.lists = listsTx{tx: tx, parent: &t}
	t.users = usersTx{tx: tx, parent: &t}
	return t
}

type Tx struct {
	rawTx    *buntdb.Tx
	writable bool
	items    itemsTx
	lists    listsTx
	users    usersTx
}

func (t *Tx) Close() {
	var err error
	if t.writable {
		err = t.rawTx.Commit()
	} else {
		err = t.rawTx.Rollback()
	}
	if err != nil {
		log.Println("TX ERROR:", err)
	}
}

func (t *Tx) Rollback() {
	err := t.rawTx.Rollback()
	if err != nil {
		log.Println("TX ERROR:", err)
	}
	return
}

func encode(in interface{}) (string, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return string(out), stored.WithCause(err, stored.CauseSerialize)
	}
	return string(out), nil
}

func decode(in string, dest interface{}) error {
	err := json.Unmarshal([]byte(in), dest)
	if err != nil {
		return stored.WithCause(err, stored.CauseMalformed)
	}
	return nil
}
