package data

import (
	"log"
	"reflect"
	"testing"

	"github.com/mbertschler/bunny/pkg/data/memory"
	"github.com/mbertschler/bunny/pkg/data/stored"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func resetDB() {
	db = memory.Open()
	setupTestdata()
}

func TestItemByID(t *testing.T) {
	resetDB()
	target := Item{
		ID:    2,
		State: ItemComplete,
		Title: "Look at Bunny",
		Body:  "By reading this text you alredy completed this item.",
	}
	item, err := ItemByID(2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(item, target) {
		t.Error("item and target are not equal", item, target)
	}
	_, err = ItemByID(22)
	if err == nil {
		t.Error("should cause an error")
	}
}

func TestUserItemByID(t *testing.T) {
	resetDB()
	target := Item{
		ID:    2,
		State: ItemComplete,
		Focus: FocusLater,
		Title: "Look at Bunny",
		Body:  "By reading this text you alredy completed this item.",
	}
	item, err := UserItemByID(1, 2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(item, target) {
		t.Error("item and target are not equal", item, target)
	}
	_, err = UserItemByID(17, 22)
	if err == nil {
		t.Error("should cause an error")
	}
}

func TestDeleteItem(t *testing.T) {
	resetDB()
	_, err := ItemByID(2)
	if err != nil {
		t.Error(err)
	}
	err = DeleteItem(2)
	if err != nil {
		t.Error(err)
	}
	_, err = ItemByID(2)
	if err == nil {
		t.Error("should cause an error")
	}
	items, err := ItemList(1)
	if err != nil {
		t.Error(err)
	}
	for _, i := range items {
		if i.ID == 2 {
			t.Error("ID 2 is still referenced")
		}
	}
	focus, err := FocusList(1)
	if err != nil {
		t.Error(err)
	}
	for _, i := range focus.Focus {
		if i.ID == 2 {
			t.Error("ID 2 is still referenced")
		}
	}
	for _, i := range focus.Later {
		if i.ID == 2 {
			t.Error("ID 2 is still referenced")
		}
	}
	for _, i := range focus.Watch {
		if i.ID == 2 {
			t.Error("ID 2 is still referenced")
		}
	}
}

func TestSetItem(t *testing.T) {
	resetDB()
	item, err := ItemByID(2)
	if err != nil {
		t.Error(err)
	}
	item.Title = "just set"
	err = SetItem(item)
	if err != nil {
		t.Error(err)
	}
	item, err = ItemByID(2)
	if err != nil {
		t.Error(err)
	}
	if item.Title != "just set" {
		t.Error("title was not set", item.Title)
	}
	item.ID = 22
	err = SetItem(item)
	if err == nil {
		t.Error("should cause an error")
	}
}

func TestNewItem(t *testing.T) {
	resetDB()
	item1, err := NewItem()
	if err != nil {
		t.Error(err)
	}
	item2, err := NewItem()
	if err != nil {
		t.Error(err)
	}
	if item1.ID == item2.ID {
		t.Error("ids are identical", item1, item2)
	}
	_, err = ItemByID(item1.ID)
	if err != nil {
		t.Error(err)
	}
	_, err = ItemByID(item2.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestSetList(t *testing.T) {
	resetDB()
	list, err := ListByID(1)
	if err != nil {
		t.Error(err)
	}
	list.Title = "just set"
	err = SetList(list)
	if err != nil {
		t.Error(err)
	}
	list, err = ListByID(1)
	if err != nil {
		t.Error(err)
	}
	if list.Title != "just set" {
		t.Error("title was not set", list.Title)
	}
	list.ID = 22
	err = SetList(list)
	if err == nil {
		t.Error("should cause an error")
	}
}

func TestItemList(t *testing.T) {
	resetDB()
	list, err := ItemList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list) != 5 {
		t.Error("expected 5 items")
	}
	var sum int
	for _, item := range list {
		sum += int(item.Focus)
	}
	if sum != 0 {
		t.Error("expected focus sum to be 0")
	}
	_, err = ItemList(12)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestUserItemList(t *testing.T) {
	resetDB()
	list, err := UserItemList(1, 1)
	if err != nil {
		t.Error(err)
	}
	if len(list) != 5 {
		t.Error("expected 5 items")
	}
	var sum int
	for _, item := range list {
		sum += int(item.Focus)
	}
	if sum != 6 {
		t.Error("expected focus sum to be 6", sum)
	}
	_, err = UserItemList(2, 1)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestFocusList(t *testing.T) {
	resetDB()
	_, err := FocusList(1)
	if err != nil {
		t.Error(err)
	}
	_, err = FocusList(12)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestUser(t *testing.T) {
	resetDB()
	_, err := UserByID(1)
	if err != nil {
		t.Error(err)
	}
	_, err = UserByID(12)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestSetFocus(t *testing.T) {
	resetDB()
	item, err := UserItemByID(1, 2)
	if err != nil {
		t.Error(err)
	}
	if item.Focus != FocusLater {
		t.Error("expected focus to be later")
	}
	err = SetFocus(1, 2, FocusWatch)
	if err != nil {
		t.Error(err)
	}
	item, err = UserItemByID(1, 2)
	if err != nil {
		t.Error(err)
	}
	if item.Focus != FocusWatch {
		t.Error("expected focus to be watch")
	}

	err = SetFocus(1, 12, FocusWatch)
	if err == nil {
		t.Error("expected an error")
	}
	err = SetFocus(12, 2, FocusWatch)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestChangeFocusNow(t *testing.T) {
	resetDB()
	list, err := FocusList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list.Focus) != 1 {
		t.Error("expected one item already focussed")
	}
	if list.Focus[0].ID != 1 {
		t.Error("expected item to be 1", list.Focus[0].ID)
	}
	err = SetFocus(1, 3, FocusNow)
	if err != nil {
		t.Error(err)
	}
	list, err = FocusList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list.Focus) != 1 {
		t.Error("expected only one item to be focussed", len(list.Focus))
	}
	if list.Focus[0].ID != 3 {
		t.Error("expected item to be 3", list.Focus[0].ID)
	}
}

var sortSet = []struct {
	Value  int
	Pos    int
	Output []int
}{
	{
		Value:  2,
		Pos:    1,
		Output: []int{2, 1, 3, 4, 5},
	}, {
		Value:  3,
		Pos:    4,
		Output: []int{1, 2, 4, 3, 5},
	},
	{
		Value:  1,
		Pos:    2,
		Output: []int{2, 1, 3, 4, 5},
	},
	{
		Value:  3,
		Pos:    1,
		Output: []int{3, 1, 2, 4, 5},
	},
	{
		Value:  2,
		Pos:    4,
		Output: []int{1, 3, 4, 2, 5},
	},
	{
		Value:  1,
		Pos:    5,
		Output: []int{2, 3, 4, 5, 1},
	},
	{
		Value:  5,
		Pos:    1,
		Output: []int{5, 1, 2, 3, 4},
	},
	{
		Value:  6,
		Pos:    1,
		Output: []int{6, 1, 2, 3, 4, 5},
	},
	{
		Value:  6,
		Pos:    6,
		Output: []int{1, 2, 3, 4, 5, 6},
	},
	{
		Value:  6,
		Pos:    3,
		Output: []int{1, 2, 6, 3, 4, 5},
	},
}

func TestSortItem(t *testing.T) {
	for i, test := range sortSet {
		resetDB()
		err := forceSetItem(Item{
			ID: 6,
		})
		if err != nil {
			t.Error(err)
		}
		list, err := debugItemList(1)
		if err != nil {
			t.Error(err)
		}
		if !isConsistent(list) {
			t.Error("inconsistent", list)
		}
		should := []int{1, 2, 3, 4, 5}
		ids := extractIDs(list)
		if !reflect.DeepEqual(ids, should) {
			t.Error("pre id order is wrong", ids)
		}
		err = SetListItemPosition(1, test.Value, test.Pos)
		if err != nil {
			t.Error(err)
		}
		list, err = debugItemList(1)
		if err != nil {
			t.Error(err)
		}
		if !isConsistent(list) {
			t.Error("inconsistent", list)
		}
		ids = extractIDs(list)
		if !reflect.DeepEqual(ids, test.Output) {
			t.Error("id order is wrong", ids, i, test)
		}
	}
}

func extractIDs(list []stored.OrderedListItem) []int {
	out := make([]int, len(list))
	for i := range list {
		out[i] = list[i].ItemID
	}
	return out
}

func isConsistent(list []stored.OrderedListItem) bool {
	j := 1
	for _, i := range list {
		if i.Position != j {
			return false
		}
		j++
	}
	return true
}
