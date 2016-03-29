package tinystore

type StoreItemAdapter interface {
	Convert(item map[string]interface{}) StoreItem
	ConvertMany(items []map[string]interface{}) []StoreItem
}

type DefaultStoreItemAdapter struct {
	convert func(item map[string]interface{}) StoreItem
	//convertMany func(items []map[string]interface{}) []StoreItem
}

func NewDefaultStoreItemAdapter(convert  func(item map[string]interface{}) StoreItem) *DefaultStoreItemAdapter {
	return &DefaultStoreItemAdapter{convert}
}

func (adapter *DefaultStoreItemAdapter) Convert(item map[string]interface{}) StoreItem {
	return adapter.convert(item)
}

func (this *DefaultStoreItemAdapter) ConvertMany(items []map[string]interface{}) []StoreItem {
	var result []StoreItem
	for _, item := range items[:] {
		result = append(result, this.Convert(item))
	}
	return result
}

//
//var convertMany = func(adapter StoreItemAdapter,items []map[string]interface{}) []StoreItem {
//	var result []StoreItem
//	for _, item := range items[:] {
//		result = append(result, adapter.Convert(item))
//	}
//	return result
//}
//
