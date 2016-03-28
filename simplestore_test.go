package tinystore_test

import (
	"testing"
	"fmt"
	"github.com/D10221/tinyStore"
)


func reverse(s string) string {
	var result []byte
	for i := len([]byte(s)) - 1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return string(result)
}


func reversePassword(item tinystore.StoreItem) (tinystore.StoreItem, error) {
	// Credential
	c, ok := item.(*DumyyItem)
	if !ok {
		return nil, tinystore.ErrInvalidStoreItem
	}
	c.Password = reverse(c.Password)
	return c, nil
}

func changePassword(newPassword string) func(item tinystore.StoreItem) (tinystore.StoreItem, error) {
	return func(item tinystore.StoreItem) (tinystore.StoreItem, error) {
		c, ok := item.(*DumyyItem)
		if !ok || c.Password == "" {
			return nil, tinystore.ErrInvalidStoreItem
		}
		c.Password = newPassword
		return item, nil
	}
}

func Test_NewLocalStore(t *testing.T) {

	store := &tinystore.SimpleStore{}

	e := store.Add(&DumyyItem{"me", "1234"})
	if e != nil {
		t.Error(e); return
	}

	if x, e := store.Find(tinystore.KeyEqualsFilter("me")); e != nil || x.GetKey() != "me" {
		t.Error(e)
		return
	}

	e = store.ForEach(reversePassword)
	if e != nil {
		t.Error(e)
		return
	}

	if x, e := tinystore.FindByKey(store, "me"); e != nil || AsCredential(x).Password != "4321" {
		if e != nil {
			t.Error(e)
			return
		}
		t.Error("UpdateAll Fail")
		return
	}

	if ex := store.Add(&DumyyItem{"you", "1234"}); ex != nil {
		t.Error(e); return
	}

	if ex := store.ForEachWhere(NameFilter("you"), reversePassword); ex != nil {
		t.Error(ex)
		return
	}

	if x, e := store.Find(NameFilter("you")); e != nil || AsCredential(x).Password != "4321" {
		if e != nil {
			t.Error(e)
			return
		}
		t.Error("UpdateWhere Fail")
		return
	}

	if ex := store.Remove(&DumyyItem{"me", "1234"}); ex != nil {
		t.Error(ex)
		return
	}

	if x, ex := store.Find(NameFilter("me")); ex != tinystore.ErrNotFound {

		t.Log(x)

		if ex != nil {
			t.Error(ex)
			return
		}
		t.Error("Not FOund?")
	}

	if found, ex := store.Find(NameFilter("you")); ex != nil || AsCredential(found).Username != "you" {
		t.Error(ex)
		return
	}
	if ex := store.RemoveWhere(NameFilter("you")); ex != nil {
		t.Error(ex)
		return
	}
	if found, ex := store.Find(NameFilter("you")); ex != tinystore.ErrNotFound {
		t.Log(found)
		t.Error(ex)
		return
	}

	if l, e := store.Length(); e != nil || l != 0 {
		if e != nil {
			t.Error(e)
			return
		}
		t.Error("Failed count")
		return
	}

	for i := 0; i <= 99; i++ {
		name := fmt.Sprintf("%s", i)
		store.Add(&DumyyItem{name, "1234"})
	}

	if l, e := store.Length(); e != nil || l != 100 {
		if e != nil {
			t.Error(e)
			return
		}
		t.Errorf("Bad Length: %s", l)
	}

	store.Clear()

	if l, e := store.Length(); e != nil || l != 0 {
		if e != nil {
			t.Error(e)
			return
		}
		t.Errorf("Bad Length: %s", l)
	}

}

func Test_Store(t *testing.T) {

	store := &tinystore.SimpleStore{}

	if e := store.Add(&DumyyItem{"admin", "password"}); e != nil {
		t.Error(e)
	}

	item, e := store.Find(NameFilter("admin"))
	if e != nil {
		t.Error(e); return
	}

	credential, ok := item.(*DumyyItem)
	if !ok {
		t.Error("Not A Credential")
	}

	if credential.Username != "admin" || credential.Password != "password" {
		t.Error("Bad store")
	}

	store.Clear()
	store.Add(&DumyyItem{Username: "me", Password:"1234"})

	// ...
	item, e = store.Find(NameFilter("me"))
	user, ok := item.(*DumyyItem)

	if e != nil {
		t.Error(e); return
	}

	if user.Username != "me" {
		t.Error("Not found")
	}

	if tinystore.Length(store) != 1 {
		t.Error("Wtf")
	}

}

// ThisTestDefaultStoreItemAdapter
/*
	ThisTestDefaultStoreItemAdapter
*/


func Test_Add(t *testing.T) {
	store := tinystore.SimpleStore{}
	e := store.Add(&DumyyItem{"me", "1234"})
	if e != nil {
		t.Error(e)
	}
	if len(store.All()) != 1 {
		t.Error("Failed len")
		return
	}
	credential, ok := store.All()[0].(*DumyyItem)
	if !ok {
		t.Error("Not a credential")
	}

	if credential.Username != "me" {
		t.Error("Fail add")
	}
}

func Test_Add_Existing(t *testing.T) {
	store := &tinystore.SimpleStore{}
	e := store.Add(&DumyyItem{"me", "1234"})
	if e != nil {
		t.Error(e)
	}
	e = store.Add(&DumyyItem{"me", "1234"})
	if e != tinystore.ErrAlreadyExists {
		t.Error("Should return AlreadyExists")
	}
	if len(store.All()) != 1 {
		t.Error("Failed len")
		return
	}
	credential, ok := store.All()[0].(*DumyyItem)
	if !ok {
		t.Error("Not a credential")
		return
	}

	if credential.Username != "me" {
		t.Error("Fail add")
	}
}

func Test_Add_Empty(t *testing.T) {
	store := tinystore.SimpleStore{}
	e := store.Add(&DumyyItem{"", ""})
	if e != tinystore.ErrInvalidStoreItem {
		t.Error("Should return InvalidCredential")
	}
	e = nil
	e = store.Add(&DumyyItem{})
	if e != tinystore.ErrInvalidStoreItem {
		t.Error("Should return InvalidCredential")
	}
	if len(store.All()) != 0 {
		t.Error("Failed len")
	}
}

func Test_Remove(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	e := store.Remove(&DumyyItem{})

	if e != tinystore.ErrInvalidStoreItem {
		t.Error("Should be invalid")
		return
	}
	e = nil
	e = store.Remove(&DumyyItem{"me", "1234"})
	if e != tinystore.ErrNotFound {
		t.Error("Should be NotFound")
		return
	}
	e = nil
	e = store.Add(&DumyyItem{"me", "1234"})
	if e != nil {
		t.Error(e); return
	}
	e = store.Remove(&DumyyItem{"me", "1234"})
	if e != nil {
		t.Errorf("Shouldn't error %s", e.Error())
		return
	}
	if tinystore.Length(store) != 0 {
		t.Error("len fail")
	}
}

func Test_Remove_Add(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	e := store.Add(&DumyyItem{"me", "1234"})
	// password is not Checked
	e = store.Remove(&DumyyItem{"me", "1111"})
	if e != nil {
		t.Error(e); return
	}
	e = store.Add(&DumyyItem{"me", "4321"})
	if e != nil {
		t.Error(e); return
	}
	if len(store.All()) != 1 {
		t.Error("Len Fail")
	}
	c, ok := store.All()[0].(*DumyyItem)
	if !ok {
		t.Error("Not a credential")
	}
	if c.Username != "me" && c.Password != "4321" {
		t.Error("Wrong values")
	}

	if u, e := tinystore.FindByKey(store, "me"); e != nil {
		c, ok := u.(*DumyyItem)
		if !ok {
			t.Error("Not a credentia"); return
		}
		if !c.Valid() || c.Username != "me" || c.Password != "4321" {
			t.Error("wrong values")
		}
		t.Error("Wtf")
	}
}

func Test_Update(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	e := store.Add(&DumyyItem{"me", "1234"})

	if e != nil {
		t.Error(e)
	}

	e = store.ForEach(changePassword("abcd"))
	if e != nil {
		t.Error(e)
	}
	if len(store.All()) != 1 || AsCredential(store.All()[0]).Password != "abcd" {
		t.Error("Doesn't work")
	}
}

func Test_UpdateWhere(t *testing.T) {

	store := tinystore.SimpleStore{}

	e := store.Add(&DumyyItem{"me", "1234"})

	if e != nil {
		t.Error(e)
	}

	e = store.ForEachWhere(NameFilter("me"), changePassword("abcd"))
	if e != nil {
		t.Error(e)
	}
	if len(store.All()) != 1 || AsCredential(store.All()[0]).Password != "abcd" {
		t.Error("Doesn't work")
	}
}

func Test_Find(t *testing.T) {

	store := tinystore.SimpleStore{}

	e := store.Add(&DumyyItem{"me", "1234"})

	if e != nil {
		t.Error(e); return
	}

	found, e := store.Find(NameFilter("me"))

	if e != nil {
		t.Error(e)
	}

	if AsCredential(found).Username != "me" {
		t.Error("Wtf")
	}
}

func Test_RemoveWhere(t *testing.T) {

	store := &tinystore.SimpleStore{}

	e := store.Add(&DumyyItem{"me", "1234"})
	if e != nil {
		t.Error(e);
		return
	}

	e = store.RemoveWhere(NameFilter("me"))
	if e != nil {
		t.Error(e)
		return
	}

	found, e := store.Find(NameFilter("me"))

	if  e != tinystore.ErrNotFound {
		if e!=nil {
			t.Error(e)
			return
		}
		t.Error("Wrong Error Type")
		return
	}

	if  IsDumyyItem(found) || ( found!=nil && found.Valid() ) {
		t.Error("Shouldn't be found or valid ")
	}

}

