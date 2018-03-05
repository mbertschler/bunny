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
	"github.com/mbertschler/bunny/pkg/data/stored"
)

func (d *DB) UserByID(id int) (stored.User, error) {
	var user stored.User
	tx, err := d.View()
	if err != nil {
		return user, err
	}
	user, err = tx.users.Get(id)
	tx.Close()
	return user, err
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

func (d *DB) UserItemByID(user, id int) (stored.Item, error) {
	var item stored.Item
	tx, err := d.View()
	if err != nil {
		return item, err
	}
	defer tx.Close()
	item, err = tx.items.Get(id)
	if err != nil {
		return item, err
	}
	focus, err := tx.users.ItemFocus(user, id)
	item.Focus = focus
	return item, err
}

func (d *DB) DebugItemList(id int) ([]stored.OrderedListItem, error) {
	var items []stored.OrderedListItem
	_, raw, err := d.ItemList(id)
	if err != nil {
		return items, err
	}
	for i, el := range raw {
		items = append(items, stored.OrderedListItem{
			Position: i + 1,
			Item:     el,
		})
	}
	return items, err
}

func (d *DB) ItemList(id int) (stored.List, []stored.Item, error) {
	var list stored.List
	var items []stored.Item
	tx, err := d.View()
	if err != nil {
		return list, items, err
	}
	defer tx.Close()
	list, err = tx.lists.Get(id)
	if err != nil {
		return list, items, err
	}
	items, err = tx.lists.Items(id)
	return list, items, err
}

func (d *DB) UserItemList(user, id int) (stored.List, []stored.Item, error) {
	var items []stored.Item
	var list stored.List
	tx, err := d.View()
	if err != nil {
		return list, items, err
	}
	items, err = tx.lists.UserItems(user, id)
	return list, items, err
}

func (d *DB) FocusList(user int) ([]stored.Item, error) {
	tx, err := d.View()
	if err != nil {
		return nil, err
	}
	defer tx.Close()
	return tx.users.AllByUser(user)
}

func (d *DB) SetItem(i stored.Item) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	defer tx.Close()
	_, err = tx.items.Get(i.ID)
	if err != nil {
		return err
	}
	return tx.items.Set(i)
}

func (d *DB) ForceSetItem(i stored.Item) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	defer tx.Close()
	return tx.items.Set(i)
}

func (d *DB) SetList(l stored.List) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	defer tx.Close()
	_, err = tx.lists.Get(l.ID)
	if err != nil {
		return err
	}
	return tx.lists.Set(l)
}

func (d *DB) ForceSetList(l stored.List) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	defer tx.Close()
	return tx.lists.Set(l)
}

func (d *DB) DeleteList(id int) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	defer tx.Close()
	return tx.lists.Delete(id)
}

func (d *DB) DeleteItem(id int) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	defer tx.Close()
	return tx.items.Delete(id)
}

func (d *DB) NewItem(i stored.Item) (int, error) {
	tx, err := d.Update()
	if err != nil {
		return 0, err
	}
	defer tx.Close()
	return tx.items.New(i)
}

func (d *DB) SetListItemPosition(list, item, pos int) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	err = tx.lists.SetItemPos(list, item, pos)
	tx.Close()
	return err
}

func (d *DB) SortUserFocusAfter(user, id, after int) error {
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

	tx, err := d.Update()
	if err != nil {
		return err
	}
	err = tx.users.SetFocus(user, item, focus)
	tx.Close()
	return err
}

func (d *DB) ListByID(id int) (stored.List, error) {
	var list stored.List
	tx, err := d.View()
	if err != nil {
		return list, err
	}
	list, err = tx.lists.Get(id)
	tx.Close()
	return list, err
}

func (d *DB) ForceSetUser(u stored.User) error {
	tx, err := d.Update()
	if err != nil {
		return err
	}
	err = tx.users.Set(u)
	tx.Close()
	return err
}
