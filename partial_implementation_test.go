package tinystore_test

import (
	"testing"
	"github.com/D10221/tinystore"
)

// Test_Partial_Implementaion , test handling partial implementaions
func Test_Partial_Implementation(t *testing.T) {
	var store tinystore.Store = &PartialStore{}
	if store ==nil {
		t.Error("Wtf")
	}
	// TODO:
}

type PartialStore struct {
	Functions map[string]interface{}
	Name string
}

// Load  items, replace store.items with provided items
func (store *PartialStore) Load(c ...tinystore.StoreItem) error {
	return tinystore.ErrNotImplemented
}

func (store PartialStore) GetName() string {
	return store.Name
}

// All items
func (store *PartialStore) All() []tinystore.StoreItem {
	panic(tinystore.ErrNotImplemented)
}

// Find first item where filter returns true
func (store *PartialStore) Find(c tinystore.Filter) (tinystore.StoreItem, error){
	return nil, tinystore.ErrNotImplemented
}

// Add new item if valid and Key Not Exists
func (store *PartialStore) Add(credential tinystore.StoreItem) error {
	return tinystore.ErrNotImplemented
}

// Remove Key matching Item
func (store *PartialStore) Remove(credential tinystore.StoreItem) error {
	return tinystore.ErrNotImplemented
}

// Empty Store
func (store *PartialStore) Clear() {
	panic(tinystore.ErrNotImplemented)
}

// ForEach mutate item with provided mutator
func (store *PartialStore) ForEach(transform tinystore.Mutator) error {
	return tinystore.ErrNotImplemented
}

// RemoveWhere filter returns true
// Note: there is a higher order func doing the same
func (store *PartialStore) RemoveWhere(f tinystore.Filter) error {
	return tinystore.ErrNotImplemented
}

// ForEachWhere filter return true , mutate item with provided mutator
// Note: there is a higher order ForEach func doing the same
func (store *PartialStore) ForEachWhere(f tinystore.Filter,transform tinystore.Mutator) error {

	return tinystore.ErrNotImplemented
}

