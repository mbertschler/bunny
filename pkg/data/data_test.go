package data

import (
	"log"
	"reflect"
	"testing"

	"github.com/mbertschler/bunny/pkg/data/memory"
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
	item, err = ItemByID(22)
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
	item, err = UserItemByID(17, 22)
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
}
