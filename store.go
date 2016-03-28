package tinystore

type StoreItem interface {
	Valid() bool
	Validate() error
	GetKey() interface{}
}

type Filter func (a StoreItem) bool ;

func NotFilter(filter Filter) Filter {
	return func(c StoreItem) bool {
		return !filter(c)
	}
}

type Mutator  func(in StoreItem) (StoreItem, error );

var Always Filter = func( a StoreItem) bool {return  true }

// Store  minimal interface
type Store interface {

	// Load  items, replace store.items with provided items
	Load(c ...StoreItem) error

	// All items
	All() []StoreItem

	// Find first item where filter returns true
	Find(c Filter) (StoreItem, error)

	// Add new item if valid and Key Not Exists
	Add(credential StoreItem) error

	// Remove Key matching Item
	Remove(credential StoreItem) error

	// Empty Store
	Clear()

	// ForEach mutate item with provided mutator
	ForEach(transform Mutator) error

	RemoveWhere(f Filter) error

	// ForEachWhere filter return true , mutate item with provided mutator
	ForEachWhere(f Filter,transform Mutator) error

	// GetKey for item return value as key, identifier
	// GetKeyOf(item StoreItem) interface{}
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
		if filter ==nil || filter(item) {
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


//
//func (store *SimpleCredentialStore) FindByUserName(userName string) (StoreItem, error) {
//	found := &Credential{}
//	for _, credential := range store.All()[:] {
//		if userName == credential.Username {
//			found = credential
//			break
//		}
//	}
//	if found.Valid() {
//		return found, nil
//	}
//
//	return found, ErrNotFound
//}
//


