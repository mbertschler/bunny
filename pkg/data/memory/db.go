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

	"github.com/tidwall/buntdb"
)

func Open() *DB {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	setupIndices(db)
	return &DB{
		db: db,
	}
}

type DB struct {
	db *buntdb.DB
}

func (d *DB) View() (Tx, error) {
	tx, err := d.db.Begin(false)
	return Tx{
		rawTx:    tx,
		writable: false,
		items:    itemsTx{tx: tx},
	}, err
}

func (d *DB) Update() (Tx, error) {
	tx, err := d.db.Begin(true)
	return Tx{
		rawTx:    tx,
		writable: true,
		items:    itemsTx{tx: tx},
	}, err
}

type Tx struct {
	rawTx    *buntdb.Tx
	writable bool
	items    itemsTx
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
