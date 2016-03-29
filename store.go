package tinystore

import (
	"io/ioutil"
	"encoding/json"
)

// StoreItem interface
type StoreItem interface {
	// Valid returns true if the item is valid
	Valid() bool
	// Validate return nil if item is valid return error if not
	//  plan: ValidationError? can be implemented  ?
	//  plan: Register Validators and exec all register validators ?
	Validate() error
	// GetKey returns something? as key , implementing Item choose what to return as key
	GetKey() interface{}
}

// Filter returns true/false takes a StoreItem
type Filter func (a StoreItem) bool ;

// NotFilter negates a Filter
func NotFilter(filter Filter) Filter {
	return func(c StoreItem) bool {
		return !filter(c)
	}
}

// Mutator is used to mutate the input thing
type Mutator  func(in StoreItem) (StoreItem, error );

// Always is shortcut for true
var Always Filter = func( a StoreItem) bool {return  true }

// Store  minimal interface
type Store interface {

	// Load  items, replace store.items with provided items
	Load(items ...StoreItem) error

	// All items
	All() []StoreItem

	// Find first item where filter returns true
	Find(filter Filter) (StoreItem, error)

	// Add new item if valid and Key Not Exists
	Add(item StoreItem) error

	// Remove Key matching Item
	Remove(item StoreItem) error

	// Empty Store
	Clear()

	// ForEach mutate item with provided mutator
	ForEach(f Mutator) error

	// RemoveWhere filter returns true
	// Note: there is a higher order func doing the same
	RemoveWhere(f Filter) error

	// ForEachWhere filter return true , mutate item with provided mutator
	// Note: there is a higher order ForEach func doing the same
	ForEachWhere(f Filter,transform Mutator) error

	// GetName implements Store.GetName
	// instance name, nick name, Identifier , etc ...
	GetName() string

}

// StoreError this package error
type StoreError struct {
	Message string
	Code    int
}

// Error implements error interface
func (e StoreError) Error() string {
	return e.Message
}

// NewError helper
func NewError(message string , code int) *StoreError {
	return &StoreError{message, code}
}

var (
	// ErrNotFound , item was NOT found
	ErrNotFound = NewError("Credential not found", 1)

	// InvalidStoreItem , item is Not valid
	ErrInvalidStoreItem = NewError("Invalid StoreItem", 2)

	// ErrAlreadyExists item's key already in store
	ErrAlreadyExists = NewError("Credential Already Exists", 3)

	// ErrNotImplemented
	ErrNotImplemented = NewError("Not Implemented", 4)
)


// Length returns  len(store.items) or 0(zero) if nil
func Length(store Store) int {
	if store ==nil {
		return 0
	}
	return  len(store.All()[:])
}

// FindByKey returns first item in items who's key == key,
// would be nice to have TKeyType instead of interface{}
func FindByKey(store Store, key interface{}) (StoreItem, error) {
	return store.Find(func(item StoreItem) bool {
		return item.GetKey() == key
	})
}

func KeyEqualsFilter(key interface{}) Filter {
	return func(item StoreItem) bool {
		return item.GetKey() == key
	}
}

// Exists in store an item where find returns true ?
func Exists(store Store, find func(item StoreItem) bool ) bool {
	found := false
	for _, item := range store.All() {
		found = find(item)
		if found {
			break
		}
	}
	return found
}

// Where ... returns slice with items found by find and items count , so we can do if count > 0
// as most of the time we would like to know if returns any
func Where(store Store,find Filter ) (items []StoreItem, count int) {
	for _, item := range store.All()[:] {
		if find(item) {
			items = append(items, item)
			count++
		}
	}
	return items, count
}

// WhereNot items where filter return false , and EXTRA matching item count
func WhereNot(store Store, filter Filter) (items []StoreItem, count int)  {
	return Where(store, NotFilter(filter))
}

// ForEach for each item in Store.items appply mutator if filter is nil or filter returns true
func ForEach(store Store, mutator Mutator, filter Filter) error {

	all := make([]StoreItem, 0)
	var e error = ErrNotFound
	for _, item := range store.All() {
		// Optional Filter
		if filter == nil || filter(item) {
			result, err := mutator(item)
			all = append(all, result)
			e = err
		}
	}
	if e == nil {
		return store.Load(all...)
	}
	return e
}

var StoreAdapters = make(map[string]StoreItemAdapter)

func RegisterStoreAdapter(store Store, adapter StoreItemAdapter)  error {
	name:= store.GetName()
	if name == "" {
		// TODO: specific error
		return ErrNotFound
	}

	StoreAdapters[store.GetName()] = adapter

	return nil
}



// LoadJson
func LoadJsonFile(store Store,path string) error {

	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	return LoadJson(store, bytes)
}


// LoadJson
func LoadJson(store Store,bytes []byte) error {

	items:= make([]map[string]interface{}, 0)

	e := json.Unmarshal(bytes, &items)
	if e != nil {
		return e
	}

	adapter, exists := StoreAdapters[store.GetName()]
	if !exists {
		return ErrNotFound
	}

	e = store.Load(adapter.ConvertMany(items)...)

	return e
}
