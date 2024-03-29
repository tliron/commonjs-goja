package commonjs

import (
	"sync"

	"github.com/dop251/goja"
)

func NewObject(runtime *goja.Runtime, value any) *goja.Object {
	var object_ *goja.Object

	switch value__ := value.(type) {
	case *goja.Object:
		object_ = value__

	case goja.Value:
		object_ = value__.ToObject(runtime)

	default:
		object_ = runtime.ToValue(value).ToObject(runtime)
	}

	object := runtime.NewObject()

	// Copy keys
	for _, key := range object_.Keys() {
		object.Set(key, object_.Get(key))
	}

	return object
}

//
// ThreadSafeObject
//

type ThreadSafeObject struct {
	map_ sync.Map
}

func NewThreadSafeObject() *ThreadSafeObject {
	return &ThreadSafeObject{}
}

func (self *ThreadSafeObject) NewDynamicObject(runtime *goja.Runtime) *goja.Object {
	return runtime.NewDynamicObject(self)
}

// ([goja.DynamicObject] interface)
func (self *ThreadSafeObject) Get(key string) goja.Value {
	if value, ok := self.map_.Load(key); ok {
		return value.(goja.Value)
	} else {
		return nil
	}
}

// ([goja.DynamicObject] interface)
func (self *ThreadSafeObject) Set(key string, value goja.Value) bool {
	self.map_.Store(key, value)
	return true
}

// ([goja.DynamicObject] interface)
func (self *ThreadSafeObject) Has(key string) bool {
	_, ok := self.map_.Load(key)
	return ok
}

// ([goja.DynamicObject] interface)
func (self *ThreadSafeObject) Delete(key string) bool {
	self.map_.Delete(key)
	return true
}

// ([goja.DynamicObject] interface)
func (self *ThreadSafeObject) Keys() []string {
	var keys []string
	self.map_.Range(func(key any, value any) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}
