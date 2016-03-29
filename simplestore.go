package tinystore

import (
	"sync"
	"errors"
)

// SimpleStore implements  Store
type SimpleStore struct {
	mutex   sync.Mutex

	items   []StoreItem

	// Name instance name , nick name , identifier , etc...
	Name string
}

func (store SimpleStore) GetName() string {
	return store.Name
}

// All implements Store.All
func (s *SimpleStore) All() []StoreItem {
	return s.items
}

// Length implements Store.Length
func(store *SimpleStore) Length() (int, error ) {
	if store ==nil {
		return 0 , errors.New("Nil Items")

	}
	return len(store.items), nil
}

// Find implements Store.Find
func (s *SimpleStore) Find(f Filter) (StoreItem, error) {
	for _, x := range s.items {
		if f(x) {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

// Clear Implements Store.Clear
func (s *SimpleStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items = make([]StoreItem, 0)
}

// Remove implements Store.Remove
func (store *SimpleStore) Remove(item StoreItem) error {

	if ex:= item.Validate() ; ex!= nil { return ex }

	store.mutex.Lock()

	defer store.mutex.Unlock()

	result := make([]StoreItem, 0)

	var e error = ErrNotFound

	for _, current := range store.items {
		if  AreKeysEqual(current, item){
			e = nil
			continue
		} else {
			result = append(result, current)
		}
	}

	if e == nil {
		store.items = result
	}

	return e
}
func AreKeysEqual(a StoreItem, b StoreItem) bool {
	if a == nil {
		return a == b
	}
	return a.GetKey() == b.GetKey()
}
// Add always...
func (store *SimpleStore) Add(item StoreItem) error {

	if ex:= item.Validate(); ex != nil {
		return ex
	}

	found := Exists(store, func(x StoreItem) bool {
		return x.GetKey() == item.GetKey()
	})

	if !found {
		store.mutex.Lock()
		store.items = append(store.items, item)
		store.mutex.Unlock()
		return nil
	}

	return ErrAlreadyExists
}

// RemoveWhere should go , .. should be <tinystore>.RemoveWhere(store, filter ) error
func (s *SimpleStore) RemoveWhere(find Filter) error  {

	result := make([]StoreItem, 0)
	var e error = ErrNotFound

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, x := range s.items[:] {
		if find(x) { // Skip
			e = nil
			continue
		}
		result = append(result, x)
	}
	if e ==nil {
		s.items = result
	}

	return e
}

// ForEach implements Store.ForEach
func (s *SimpleStore) ForEach(f Mutator) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var e error = nil
	for i, x := range s.items {
		r, e := f(x)
		if e == nil {
			s.items[i] = r
		}
	}
	return e
}

// ForEachWhere implements Store.ForEachWhere
func (s *SimpleStore) ForEachWhere(find Filter, transform Mutator) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error = ErrNotFound
	for i, x := range s.items {
		if find(x) {
			if r, e := transform(x) ;e == nil {
				s.items[i] = r
				return e
			}
		}
	}
	return err
}


// GetKey implements Store.GetKey
//func (store *SimpleStore) GetKeyOf(item StoreItem) interface{} {
//	c, ok := item.(*credentials.Credential)
//	if ok {
//		return c.Username
//	}
//	// Todo: return error
//	panic("Not a Credential")
//}

// Load implements Sore.Load, does not do error  checking
func (store *SimpleStore) Load(c ...StoreItem) error {
	store.items = c
	// satisfy Interface
	return nil
}