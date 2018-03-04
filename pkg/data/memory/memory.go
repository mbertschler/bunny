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
	listOrderPrefix = "lo/"
	listItemPrefix  = "li/"
	itemFocusPrefix = "if/"
	userItemPrefix  = "ui/"
)

// DELETE
func itemKey(id int) string {
	return itemPrefix + strconv.Itoa(id)
}

// DELETE
func itemID(id string) int {
	id = strings.TrimPrefix(id, itemPrefix)
	i, err := strconv.Atoi(id)
	if err != nil {
		log.Println("KEY ERROR:", err)
	}
	return i
}

func (d *DB) ItemByID(id int) (stored.Item, error) {
	var item stored.Item
	tx, err := d.View()
	if err != nil {
		return item, err
	}
	item, err = tx.items.Get(id)
	tx.Close()
	return item, err
}

func setupIndices(db *buntdb.DB) {
	db.CreateIndex("userFocus", itemFocusPrefix+"*", func(a, b string) bool {
		var focusA, focusB stored.UserFocus
		err := decode(a, &focusA)
		if err != nil {
			log.Println(err)
		}
		err = decode(b, &focusB)
		if err != nil {
			log.Println(err)
		}
		aStr := "/" + strconv.Itoa(focusA.UserID) + "/" + strconv.Itoa(focusA.ItemID) + "/"
		bStr := "/" + strconv.Itoa(focusB.UserID) + "/" + strconv.Itoa(focusB.ItemID) + "/"
		return aStr < bStr
	})
	db.CreateIndex("listItem", listItemPrefix+"*", func(a, b string) bool {
		var liA, liB stored.ListItem
		err := decode(a, &liA)
		if err != nil {
			log.Println(err)
		}
		err = decode(b, &liB)
		if err != nil {
			log.Println(err)
		}
		aStr := "/" + strconv.Itoa(liA.ListID) + "/" + strconv.Itoa(liA.ItemID) + "/"
		bStr := "/" + strconv.Itoa(liB.ListID) + "/" + strconv.Itoa(liB.ItemID) + "/"
		return aStr < bStr
	})
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

func listItemKey(list, pos int) string {
	return listItemPrefix + strconv.Itoa(list) +
		"/" + strconv.Itoa(pos)
}

func listItemID(id string) (list, pos int) {
	id = strings.TrimPrefix(id, listItemPrefix)
	parts := strings.Split(id, "/")
	list, _ = strconv.Atoi(parts[0])
	pos, _ = strconv.Atoi(parts[1])
	return
}

func focusKey(user, focus, order int) string {
	return itemFocusPrefix + strconv.Itoa(user) +
		"/" + strconv.Itoa(focus) + "/" + strconv.Itoa(order)
}

func focusIDs(id string) (user, focus, order int) {
	id = strings.TrimPrefix(id, itemFocusPrefix)
	parts := strings.Split(id, "/")
	user, _ = strconv.Atoi(parts[0])
	focus, _ = strconv.Atoi(parts[1])
	order, _ = strconv.Atoi(parts[2])
	return
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
		uf := stored.UserFocus{
			UserID: user, ItemID: id}
		start, err := encode(uf)
		if err != nil {
			return err
		}
		uf = stored.UserFocus{
			UserID: user, ItemID: id + 1}
		end, err := encode(uf)
		if err != nil {
			return err
		}
		return tx.AscendRange("userFocus", start, end, func(key, val string) bool {
			_, item.Focus, _ = focusIDs(key)
			return false
		})
	})
	return item, err
}

func (d *DB) DebugItemList(id int) ([]stored.OrderedListItem, error) {
	var items []stored.OrderedListItem
	var list stored.List
	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(listKey(id))
		if err != nil {
			return err
		}
		decode(val, &list)
		var outerErr error
		err = tx.AscendRange("", listItemPrefix+strconv.Itoa(id),
			listItemPrefix+strconv.Itoa(id+1), func(key, val string) bool {
				var listItem stored.ListItem
				outerErr = decode(val, &listItem)
				if outerErr != nil {
					return false
				}
				_, err := d.ItemByID(listItem.ItemID)
				if err != nil {
					outerErr = err
					return false
				}
				_, pos := listItemID(key)
				items = append(items, stored.OrderedListItem{
					ListItem: listItem, Position: pos})
				return true
			})
		if err != nil {
			return err
		}
		return outerErr
	})
	return items, err
}

func (d *DB) ItemList(id int) (stored.List, []stored.Item, error) {
	var items []stored.Item
	var list stored.List
	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(listKey(id))
		if err != nil {
			return err
		}
		decode(val, &list)
		var outerErr error
		err = tx.AscendRange("", listItemPrefix+strconv.Itoa(id),
			listItemPrefix+strconv.Itoa(id+1), func(key, val string) bool {
				var listItem stored.ListItem
				outerErr = decode(val, &listItem)
				if outerErr != nil {
					log.Println(err)
					return true
				}
				item, err := d.ItemByID(listItem.ItemID)
				if err != nil {
					log.Println(err)
					return true
				}
				items = append(items, item)
				return true
			})
		if err != nil {
			return err
		}
		return outerErr
	})
	return list, items, err
}

func (d *DB) UserItemList(user, id int) (stored.List, []stored.Item, error) {
	var items []stored.Item
	var list stored.List
	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(listKey(id))
		if err != nil {
			return err
		}
		err = decode(val, &list)
		if err != nil {
			return err
		}
		_, err = tx.Get(userKey(user))
		if err != nil {
			return err
		}
		var outerErr error
		err = tx.AscendRange("", listItemPrefix+strconv.Itoa(id),
			listItemPrefix+strconv.Itoa(id+1), func(key, val string) bool {
				var listItem stored.ListItem
				err := decode(val, &listItem)
				if err != nil {
					log.Println(err)
					return true
				}
				item, err := d.UserItemByID(user, listItem.ItemID)
				if err != nil {
					log.Println(err)
					return true
				}
				items = append(items, item)
				return true
			})
		if err != nil {
			return err
		}
		return outerErr
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
		return tx.AscendRange("", itemFocusPrefix+strconv.Itoa(user),
			itemFocusPrefix+strconv.Itoa(user+1), func(key, val string) bool {
				_, focus, _ := focusIDs(key)
				var uf stored.UserFocus
				decode(val, &uf)
				itemStr, err := tx.Get(itemKey(uf.ItemID))
				if err != nil {
					log.Println("wtf items are inconsistent", uf.ItemID)
					return true
				}
				var item stored.Item
				decode(itemStr, &item)
				item.Focus = focus
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

func posByItemID(tx *buntdb.Tx, list, item int) (int, bool, error) {
	li := stored.ListItem{
		ListID: list, ItemID: item}
	start, err := encode(li)
	if err != nil {
		return 0, false, err
	}
	li = stored.ListItem{
		ListID: list, ItemID: item + 1}
	end, err := encode(li)
	if err != nil {
		return 0, false, err
	}
	var found bool
	var oldKey string
	err = tx.AscendRange("listItem", start, end, func(key, val string) bool {
		oldKey = key
		found = true
		return false
	})
	var oldPos int
	if found {
		_, oldPos = listItemID(oldKey)
	}
	return oldPos, found, nil
}

func (d *DB) SetListItemPosition(list, item, pos int) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		oldPos, found, err := posByItemID(tx, list, item)
		if err != nil {
			return err
		}
		if found && pos == oldPos {
			return nil
		}
		type change struct {
			key string
			val string
		}

		var posExists bool
		_, err = tx.Get(listItemKey(list, pos))
		if err == nil {
			posExists = true
		}

		if !found {
			err = tx.DescendRange("", listItemKey(list+1, 0),
				listItemKey(list, 0), func(key, val string) bool {
					_, oldPos = listItemID(key)
					oldPos++
					return false
				})
			if err != nil {
				return err
			}
		}

		if posExists {
			if pos < oldPos {
				changed := []change{}
				err = tx.DescendRange("", listItemKey(list, oldPos-1),
					listItemKey(list, pos-1), func(key, val string) bool {
						list, order := listItemID(key)
						key = listItemKey(list, order+1)
						changed = append(changed, change{key, val})
						return true
					})
				for _, e := range changed {
					_, _, err = tx.Set(e.key, e.val, nil)
					if err != nil {
						return err
					}
				}
			} else {
				changed := []change{}
				err = tx.AscendRange("", listItemKey(list, oldPos+1),
					listItemKey(list, pos+1), func(key, val string) bool {
						list, order := listItemID(key)
						key = listItemKey(list, order-1)
						changed = append(changed, change{key, val})
						return true
					})
				for _, e := range changed {
					_, _, err = tx.Set(e.key, e.val, nil)
					if err != nil {
						return err
					}
				}
			}
		}

		key := listItemKey(list, pos)
		li := stored.ListItem{
			ListID: list, ItemID: item}
		val, err := encode(li)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(key, val, nil)
		return err
	})
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

func unsetFocus(tx *buntdb.Tx, user, item int) error {
	val, err := tx.Get(focusKey(user, 1, 1))
	if err == nil {
		var uf stored.UserFocus
		decode(val, &uf)
		_, err = tx.Delete(focusKey(user, 1, 1))
		if err != nil {
			return err
		}
		return setUserFocus(tx, user, uf.ItemID, 2)
	}
	return nil
}
func setUserFocus(tx *buntdb.Tx, user, item, focus int) error {
	var uf = stored.UserFocus{
		UserID: user,
		ItemID: item,
	}
	val, err := encode(uf)
	if err != nil {
		return err
	}
	uf = stored.UserFocus{
		UserID: user, ItemID: item}
	start, err := encode(uf)
	if err != nil {
		return err
	}
	uf = stored.UserFocus{
		UserID: user, ItemID: item + 1}
	end, err := encode(uf)
	if err != nil {
		return err
	}
	var found bool
	var oldKey string
	err = tx.AscendRange("userFocus", start, end, func(key, val string) bool {
		found = true
		oldKey = key
		return false
	})

	if found {
		_, err = tx.Delete(oldKey)
		if err != nil {
			return err
		}
	}
	var order int
	err = tx.DescendRange("", focusKey(user, focus+1, 0),
		focusKey(user, focus, 0), func(key, val string) bool {
			_, _, order = focusIDs(key)
			return false
		})
	if err != nil {
		return err
	}
	order++
	_, _, err = tx.Set(focusKey(user, focus, order), val, nil)

	return err
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

	return d.db.Update(func(tx *buntdb.Tx) error {
		if focus == 1 {
			err := unsetFocus(tx, user, item)
			if err != nil {
				return err
			}
		}
		return setUserFocus(tx, user, item, focus)
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
