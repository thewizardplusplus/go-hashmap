package hashmap

// this interface is used only for generating mocks
type handlerInterface interface {
	Handle(key Key, value interface{}) bool
}
