package tinystore


type StoreItemAdapter interface {
	Convert (item map[string]interface{}) StoreItem
	ConvertMany (items []map[string]interface{}) []StoreItem
}

type DefaultStoreItemAdapter struct {
	convert     func(item map[string]interface{}) StoreItem
	convertMany func(items []map[string]interface{}) []StoreItem
}

func NewDefaultStoreItemAdapter(
convert     func(item map[string]interface{}) StoreItem,
convertMany func(items []map[string]interface{}) []StoreItem) *DefaultStoreItemAdapter {
	return &DefaultStoreItemAdapter{convert,convertMany}
}

func (adapter *DefaultStoreItemAdapter) Convert (item map[string]interface{}) StoreItem {
	return adapter.convert(item)
}

func (adapter *DefaultStoreItemAdapter) ConvertMany (items []map[string]interface{}) []StoreItem {
	return adapter.convertMany(items)
}

