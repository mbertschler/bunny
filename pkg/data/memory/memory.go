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
	"strconv"
	"strings"

	"github.com/tidwall/buntdb"

	"github.com/mbertschler/bunny/pkg/data/stored"
)

const (
	itemPrefix      = "i/"
	listPrefix      = "l/"
	listItemPrefix  = "li/"
	userFocusPrefix = "uf/"
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

func setupIndices(db *buntdb.DB) {

}

type DB struct {
	db *buntdb.DB
}

func itemKey(id int) string {
	return itemPrefix + strconv.Itoa(id)
}

func itemID(id string) int {
	id = strings.TrimPrefix(id, itemPrefix)
	i, _ := strconv.Atoi(id)
	return i
}

func listKey(id int) string {
	return listPrefix + strconv.Itoa(id)
}

func listID(id string) int {
	id = strings.TrimPrefix(id, listPrefix)
	i, _ := strconv.Atoi(id)
	return i
}

func listItemKey(listID, itemID int) string {
	return listItemPrefix + strconv.Itoa(listID) +
		"/" + strconv.Itoa(itemID)
}

func listItemID(id string) (listID, itemID int) {
	id = strings.TrimPrefix(id, listItemPrefix)
	parts := strings.Split(id, "/")
	listID, _ = strconv.Atoi(parts[0])
	itemID, _ = strconv.Atoi(parts[1])
	return
}

func userFocusKey(user, item, focus int) string {
	return userFocusPrefix + strconv.Itoa(user) +
		"/" + strconv.Itoa(item) + "/" + strconv.Itoa(focus)
}

func userFocusIDs(id string) (user, item, focus int) {
	id = strings.TrimPrefix(id, userFocusPrefix)
	parts := strings.Split(id, "/")
	user, _ = strconv.Atoi(parts[0])
	item, _ = strconv.Atoi(parts[1])
	focus, _ = strconv.Atoi(parts[2])
	return
}

func (d *DB) ItemByID(id int) stored.Item {
	var val string
	var err error
	err = d.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(itemKey(id))
		return err
	})
	if err != nil {
		log.Println(err)
	}
	var item stored.Item
	decode(val, &item)
	return item
}

func (d *DB) ItemList(id int) []stored.Item {
	var items []stored.Item
	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(itemPrefix+"*", func(key, val string) bool {
			var item stored.Item
			decode(val, &item)
			items = append(items, item)
			return true
		})
	})
	if err != nil {
		log.Println(err)
	}
	return items
}

func (d *DB) FocusList(user int) []stored.Item {
	var items []stored.Item
	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(userFocusPrefix+"*", func(key, val string) bool {
			var item stored.Item
			decode(val, &item)
			items = append(items, item)
			return true
		})
	})
	if err != nil {
		log.Println(err)
	}
	return items
}

func (d *DB) SetItem(i stored.Item) {
	d.db.Update(func(tx *buntdb.Tx) error {
		tx.Set(itemKey(i.ID), encode(i), nil)
		return nil
	})
}

func (d *DB) SetList(l stored.List) {
	d.db.Update(func(tx *buntdb.Tx) error {
		tx.Set(listKey(l.ID), encode(l), nil)
		return nil
	})
}

func (d *DB) DeleteList(id int) {
	d.db.Update(func(tx *buntdb.Tx) error {
		tx.Delete(listKey(id))
		// TODO: delete from referenced tables
		return nil
	})
}

func (d *DB) DeleteItem(id int) {
	d.db.Update(func(tx *buntdb.Tx) error {
		tx.Delete(itemKey(id))
		// TODO: delete from referenced tables
		return nil
	})
}

func (d *DB) NewItem(i stored.Item) int {
	var id int
	d.db.Update(func(tx *buntdb.Tx) error {
		tx.DescendKeys(itemPrefix, func(key, value string) bool {
			id = itemID(key)
			return false
		})
		id++
		i.ID = id
		tx.Set(itemKey(i.ID), encode(i), nil)
		return nil
	})
	return id
}

func (d *DB) SortListItemAfter(listID, itemID, after int) {
	// var nextID int
	// d.db.Update(func(tx *buntdb.Tx) error {
	// 	afterKey := listItemKey(listID, after)
	// 	val, err := tx.Get(afterKey)
	// 	if err == nil {
	// 		var li stored.ListItem
	// 		decode(val, &li)
	// 		nextID = li.Next
	// 		li.Next = itemID
	// 		tx.Set(afterKey, encode(li), nil)
	// 	} else {
	// 		log.Println(err)
	// 	}

	// 	tx.DescendKeys(itemPrefix, func(key, value string) bool {
	// 		id = itemID(key)
	// 		return false
	// 	})
	// 	id++
	// 	i.ID = id
	// 	tx.Set(itemKey(i.ID), encode(i), nil)
	// 	return nil
	// })
	// return id
}

func (d *DB) SortUserFocusAfter(user, id, after int) {
	// var id int
	// d.db.Update(func(tx *buntdb.Tx) error {
	// 	tx.DescendKeys(itemPrefix, func(key, value string) bool {
	// 		id = itemID(key)
	// 		return false
	// 	})
	// 	id++
	// 	i.ID = id
	// 	tx.Set(itemKey(i.ID), encode(i), nil)
	// 	return nil
	// })
	// return id
}

func (d *DB) SetUserFocus(user, item, focus int) {
	// var uf = stored.UserFocus{
	// 	UserID: user,
	// 	ItemID: item,
	// 	Focus:  focus,
	// }
	// d.db.Update(func(tx *buntdb.Tx) error {
	// 	var found string
	// 	tx.AscendRange("", userFocusKey(user, focus, 0),
	// 		userFocusKey(user, focus+1, 0), func(key, value string) bool {
	// 			found = key
	// 			return true
	// 		})
	// 	if found != "" {
	// 		tx.Delete(found)
	// 	}
	// 	tx.Set(userFocusKey(user, item, focus), encode(uf), nil)
	// 	return nil
	// })
}

func encode(in interface{}) string {
	out, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func decode(in string, dest interface{}) {
	err := json.Unmarshal([]byte(in), dest)
	if err != nil {
		log.Fatal(err)
	}
}
