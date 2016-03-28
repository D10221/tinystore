package tinystore_test

import (
	"testing"
	"github.com/D10221/tinyStore"
)

func AsCredentialGetName(item tinystore.StoreItem) string {
	value, ok := item.(*DumyyItem)
	if !ok { panic("Not a Credential")}
	return value.Username
}

func AsCredential(item tinystore.StoreItem) *DumyyItem {
	value, ok := item.(*DumyyItem)
	if !ok { panic("Not a Credential")}
	return value
}

func AsCredentialNameEquals(item tinystore.StoreItem, name string) bool {
	value, ok := item.(*DumyyItem)
	if !ok { panic("Not a Credential")}
	return value.Username == name
}

func Test_TinyStoreFuncs(t *testing.T){
	var store tinystore.Store = &tinystore.SimpleStore{}
	c := &DumyyItem{"me", "1234"}
	e := store.Add(c)
	if e != nil {
		t.Error(e);
		return
	}

	if result,count := tinystore.Where(store, NameFilter("me")) ; count != 1  || AsCredentialGetName(result[0]) != "me"{
		t.Errorf("Error: result: %v ,Count: %v", result, count)
	}
}


func Test_Where(t *testing.T){
	store:=  &tinystore.SimpleStore{}
	if e:= store.Add( &DumyyItem{"xyz", "1234"} ) ; e!= nil {
		t.Error(e)
		return
	}
	if e:= store.Add( &DumyyItem{"uuu", "1234"} ); e!=nil {
		t.Error(e)
		return
	}
	if len(store.All()[:])!= 2 {
		t.Error("Failed Add")
	}
	items, count := tinystore.Where(store, NameFilter("xyz"))
	if count != 1 {
		t.Error("Bad Count")
		return
	}
	if AsCredentialGetName(items[0])!= "xyz" {
		t.Error("Bad Item Returned")
	}

	notIems, count := tinystore.WhereNot(store, NameFilter("xyz"))
	if count != 1 {
		t.Error("Bad Count")
		return
	}
	if  AsCredentialGetName(notIems[0]) != "uuu" {
		t.Error("Bad Item Returned")
		return
	}
}

func Test_FindByKey_Load(t *testing.T){
	var store tinystore.Store = &tinystore.SimpleStore{}
	if e:= store.Load(&DumyyItem{"me", "1234"}, &DumyyItem{"el", "999"});e!=nil{
		t.Error(e)
		return
	}
	if c,e:= tinystore.FindByKey(store, "el"); e!=nil  ||  AsCredentialGetName(c)!= "el"{
		t.Error("Bad result")
	}
}

func Test_Store_ForEach(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	if e:= store.Load(&DumyyItem{"me", "1234"}, &DumyyItem{"el", "999"}); e!=nil {
		t.Error(e)
		return
	}

	changePasswords := func (item tinystore.StoreItem) (tinystore.StoreItem, error) {
		c, ok:= item.(*DumyyItem)
		if!ok {
			return nil, tinystore.ErrAlreadyExists
		}
		c.Password = "xxx" ;
		return c, nil
	}

	if e:= tinystore.ForEach(store, changePasswords, nil); e !=nil {
		t.Error(e);
		return
	}

	if x,e:= tinystore.FindByKey(store, "me"); e!=nil || AsCredential(x).Password != "xxx" {
		if e!=nil { t.Error(e) ; return }
		t.Errorf("Not ok => x: %v", x)
	}


}
