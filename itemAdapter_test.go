package tinystore_test

import (
	"testing"
	"github.com/D10221/tinystore"
	"io/ioutil"
)

// Test_ItemStoreItemAdapter
func Test_ItemStoreItemAdapter(t *testing.T) {

	var adapter tinystore.StoreItemAdapter = tinystore.NewDefaultStoreItemAdapter(convert)

	items := []map[string]interface{}{
		{
			"Username": "me",
			"Password": "1234",
		},
		{
			"Username": "el",
			"Password": "1234",
		},
	}
	// ...
	{
		converted := adapter.Convert(items[0])
		expected:= &DumyyItem{"me", "1234"}
		if result, ok := converted.(*DumyyItem); ok && result.Equals(expected) {
			// ok
		} else {
			t.Errorf("Expected %v got %v", expected , result)
			return
		}
	}
	// ...
	{
		var expected = []*DumyyItem {
			&DumyyItem {"me", "1234"},
			&DumyyItem {"el", "1234"},
		}

		many:= adapter.ConvertMany(items)

		for i, item := range many {

			if result, ok := item.(*DumyyItem); ok && result.Equals(expected[i]) {
				// ok
			} else {
				t.Errorf("Expected %v got %v", expected[i] , result)
			}
		}
	}


}


var convert = func(item map[string]interface{}) tinystore.StoreItem {

	newItem:= &DumyyItem{}

	newItem.Username = GetStringOrPanic(item, "Username")

	if keyValue, keyExists := item["Password"]; keyExists  {
		if value, ok := keyValue.(string); ok {
			newItem.Password = value
		}
	}
	return newItem
}

// GetString get key from map , returns value , success
func GetString(m map[string]interface{}, key string) (string, bool) {
	if keyValue, exists:= m[key]; exists {
		if value, ok := keyValue.(string) ; ok {
			return value , ok
		}
	}
	return "", false
}

// GetStringOrPanic gets key from map as string or panic
func GetStringOrPanic(m map[string]interface{} , key string) string {
	s, ok := GetString(m , key)
	if !ok {
		panic("Invalid Conversion")
	}
	return s
}


// Test_Store_json_load
func Test_StoreLoadsJsonFile(t *testing.T) {

	store := &tinystore.SimpleStore{ Name: "SimpleStore"}

	tinystore.RegisterStoreAdapter(store, tinystore.NewDefaultStoreItemAdapter(convert))

	e := tinystore.LoadJsonFile(store, "testdata/credentials.json")

	if e != nil {
		t.Error(e)
		return
	}

	item := store.All()[0]
	found, ok := item.(*DumyyItem)

	if !ok || found.Username != "admin" {
		t.Error("LoadJson Failed")
	}

}

// Test_Store_json_load
func Test_StoreLoadsJson(t *testing.T) {

	store := &tinystore.SimpleStore{ Name: "SimpleStore"}

	tinystore.RegisterStoreAdapter(store, tinystore.NewDefaultStoreItemAdapter(convert))

	j, e := ioutil.ReadFile("testdata/credentials.json")
	if e!= nil {
		t.Error(e)
	}

	e = tinystore.LoadJson(store, j)

	if e != nil {
		t.Error(e)
		return
	}

	item := store.All()[0]
	found, ok := item.(*DumyyItem)

	if !ok || found.Username != "admin" {
		t.Error("LoadJson Failed")
	}

}
