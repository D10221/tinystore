package tinystore_test

import (
	"testing"
	"github.com/D10221/tinyStore"
)


// DummyItem implements tinystore.StoreItem
type DumyyItem struct {
	Username string
	Password string
}

func (this *DumyyItem ) Test(test func(item *DumyyItem) bool ) bool {

	if this == nil {
		return false
	}

	return test(this)
}

func IsDumyyItem(item tinystore.StoreItem) bool {
	 _ ,ok := item.(*DumyyItem)
	return ok
}

func (this *DumyyItem ) Valid() bool {
	return this != nil && this.Username != "" && this.Password != ""
}

func (this *DumyyItem) Validate() error {
	if this != nil && this.Valid() {
		return nil
	}
	return tinystore.ErrInvalidStoreItem
}

func (this *DumyyItem) GetKey() interface{} {
	return this.Username
}

func (this *DumyyItem) Equals(other *DumyyItem) bool {
	if this == nil || other == nil {
		return this == other
	}
	return this.Username == other.Username && this.Password == other.Password
}
func NameEquals(name string) func (item *DumyyItem) bool {
	return func (item *DumyyItem) bool {
		return item.Username == name
	}
}
// NameFilter return true if StoreItem is DumyyItem && DumyyItem.Username == name
func NameFilter(name string) tinystore.Filter {
	return func(item tinystore.StoreItem) bool {
		if item == nil {
			panic("No Item")
		}
		value, ok := item.(*DumyyItem)
		if !ok {
			return false
			//panic("Not a credential")
		}
		return value.Test(NameEquals(name))
	}
}

// Test_NameFilter
func Test_NameFilter(t *testing.T){

	item := &DumyyItem{"1", "1"}
	if !NameFilter("1")(item){
		t.Error("NameFilter Failed")
	}
	item = nil
	if NameFilter("1")(item){
		t.Error("NameFilter Failed")
	}
	var noItem *DumyyItem
	if NameFilter("1")(noItem){
		t.Error("NameFilter Failed")
	}
}

// Test_Valid
func Test_Valid(t *testing.T){
	if (&DumyyItem{"", ""}).Valid() {
		t.Error("It Shouldn't be valid")
	}

	if !(&DumyyItem{"me", "1"}).Valid() {
		t.Error("It Should be valid")
	}
	if e:= (&DumyyItem{"", ""}).Validate() ; e == nil {
		t.Error("It Shouldn't validate")
	}
	if e:= (&DumyyItem{"me", "1"}).Validate(); e!= nil {
		t.Error("It Should validate")
	}

	var dummy *DumyyItem // nil
	if dummy.Valid() {
		t.Error("It Shouldn't be valid, is nil ")
	}
	if e:= dummy.Validate() ; e == nil {
		t.Error("It Shouldn't be valid, is nil ")
	}

}
