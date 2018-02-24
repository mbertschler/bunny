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
	userPrefix      = "u/"
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

func userKey(id int) string {
	return userPrefix + strconv.Itoa(id)
}

func userID(id string) int {
	id = strings.TrimPrefix(id, userPrefix)
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

func userFocusKey(user, item int) string {
	return userFocusPrefix + strconv.Itoa(user) +
		"/" + strconv.Itoa(item)
}

func userFocusIDs(id string) (user, item int) {
	id = strings.TrimPrefix(id, userFocusPrefix)
	parts := strings.Split(id, "/")
	user, _ = strconv.Atoi(parts[0])
	item, _ = strconv.Atoi(parts[1])
	return
}

func (d *DB) ItemByID(id int) (stored.Item, error) {
	var val string
	var err error
	err = d.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(itemKey(id))
		return err
	})
	var item stored.Item
	if err != nil {
		return item, err
	}
	err = decode(val, &item)
	return item, err
}

func (d *DB) UserByID(id int) (stored.User, error) {
	var val string
	var err error
	err = d.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(userKey(id))
		return err
	})
	var user stored.User
	if err != nil {
		return user, err
	}
	err = decode(val, &user)
	return user, err
}

func (d *DB) UserItemByID(user, id int) (stored.Item, error) {
	var item stored.Item
	err := d.db.View(func(tx *buntdb.Tx) error {
		itemStr, err := tx.Get(itemKey(id))
		if err != nil {
			return err
		}
		err = decode(itemStr, &item)
		if err != nil {
			return err
		}
		focusStr, err := tx.Get(userFocusKey(user, id))
		if err == nil {
			var focus stored.UserFocus
			err = decode(focusStr, &focus)
			if err != nil {
				return err
			}
			item.Focus = focus.Focus
		}
		return err
	})
	return item, err
}

func (d *DB) ItemList(id int) (stored.List, []stored.Item, error) {
	var items []stored.Item
	var list stored.List
	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(itemKey(id))
		if err != nil {
			return err
		}
		decode(val, &list)
		return tx.AscendKeys(itemPrefix+"*", func(key, val string) bool {
			var item stored.Item
			decode(val, &item)
			items = append(items, item)
			return true
		})
	})
	return list, items, err
}

func (d *DB) UserItemList(user, id int) (stored.List, []stored.Item, error) {
	var items []stored.Item
	var list stored.List
	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(itemKey(id))
		if err != nil {
			return err
		}
		decode(val, &list)
		_, err = tx.Get(userKey(user))
		if err != nil {
			return err
		}
		return tx.AscendKeys(itemPrefix+"*", func(key, val string) bool {
			var item stored.Item
			decode(val, &item)
			uf, err := tx.Get(userFocusKey(user, item.ID))
			if err == nil {
				var focus stored.UserFocus
				decode(uf, &focus)
				item.Focus = focus.Focus
			}
			items = append(items, item)
			return true
		})
	})
	return list, items, err
}

func (d *DB) FocusList(user int) ([]stored.Item, error) {
	var items []stored.Item
	err := d.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(userKey(user))
		if err != nil {
			return err
		}
		return tx.AscendKeys(userFocusPrefix+"*", func(key, val string) bool {
			_, itemID := userFocusIDs(key)
			var focus stored.UserFocus
			decode(val, &focus)
			itemStr, err := tx.Get(itemKey(itemID))
			if err != nil {
				log.Println("wtf items are inconsistent", itemID)
				return true
			}
			var item stored.Item
			decode(itemStr, &item)
			item.Focus = focus.Focus
			items = append(items, item)
			return true
		})
	})
	return items, err
}

func (d *DB) SetItem(i stored.Item) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := encode(i)
		if err != nil {
			return err
		}
		key := itemKey(i.ID)
		_, err = tx.Get(key)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(key, val, nil)
		return err
	})
}

func (d *DB) ForceSetItem(i stored.Item) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := encode(i)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(itemKey(i.ID), val, nil)
		return err
	})
}

func (d *DB) SetList(l stored.List) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := encode(l)
		if err != nil {
			return err
		}
		key := listKey(l.ID)
		_, err = tx.Get(key)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(key, val, nil)
		return err
	})
}

func (d *DB) ForceSetList(l stored.List) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := encode(l)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(listKey(l.ID), val, nil)
		return err
	})
}

func (d *DB) DeleteList(id int) {
	d.db.Update(func(tx *buntdb.Tx) error {
		tx.Delete(listKey(id))
		// TODO: delete from referenced tables
		return nil
	})
}

func (d *DB) DeleteItem(id int) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(itemKey(id))
		// TODO: delete from referenced tables
		return err
	})
}

func (d *DB) NewItem(i stored.Item) (int, error) {
	var id int
	err := d.db.Update(func(tx *buntdb.Tx) error {
		err := tx.DescendKeys(itemPrefix+"*", func(key, value string) bool {
			id = itemID(key)
			return false
		})
		if err != nil {
			return err
		}
		id++
		i.ID = id
		val, err := encode(i)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(itemKey(i.ID), val, nil)
		return err
	})
	return id, err
}

func (d *DB) SortListItemAfter(listID, itemID, after int) error {
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
	return nil
}

func (d *DB) SortUserFocusAfter(user, id, after int) error {
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
	return nil
}

func (d *DB) SetUserFocus(user, item, focus int) error {
	_, err := d.UserByID(user)
	if err != nil {
		return err
	}
	_, err = d.ItemByID(item)
	if err != nil {
		return err
	}
	var uf = stored.UserFocus{
		Focus: focus,
	}
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := encode(uf)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(userFocusKey(user, item), val, nil)
		return err
	})
}

func (d *DB) ListByID(id int) (stored.List, error) {
	var val string
	var err error
	err = d.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(listKey(id))
		return err
	})
	var list stored.List
	if err != nil {
		return list, err
	}
	err = decode(val, &list)
	return list, err
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

func (d *DB) ForceSetUser(u stored.User) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := encode(u)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(userKey(u.ID), val, nil)
		return err
	})
}
