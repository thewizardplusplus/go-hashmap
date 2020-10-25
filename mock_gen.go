package hashmap

//go:generate mockery -name=HandlerInterface -inpkg -case=underscore -testonly

// HandlerInterface ...
//
// It's used only for mock generating.
//
type HandlerInterface interface {
	Handle(key Key, value interface{}) bool
}
