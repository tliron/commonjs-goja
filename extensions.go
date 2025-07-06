package commonjs

import (
	"github.com/dop251/goja"
)

// Can return a goja.Value, nil, or other values, which will be converted to a
// goja.Value.
type CreateExtensionFunc func(jsContext *Context) any

//
// Extension
//

type Extension struct {
	Name   string
	Create CreateExtensionFunc
}

func NewExtensions(extensions map[string]CreateExtensionFunc) []Extension {
	if len(extensions) == 0 {
		return nil
	}

	extensions_ := make([]Extension, len(extensions))
	index := 0
	for name, api := range extensions {
		extensions_[index] = Extension{
			Name:   name,
			Create: api,
		}
		index++
	}
	return extensions_
}

func (self *Context) AppendExtensions() {
	for _, extension := range self.Environment.Extensions {
		self.AppendExtension(extension)
	}
}

func (self *Context) AppendExtension(extension Extension) {
	if extension := self.CreateExtension(extension); extension != nil {
		self.Extensions = append(self.Extensions, extension)
	}
}

func (self *Context) CreateExtension(extension Extension) goja.Value {
	if value := extension.Create(self); value != nil {
		if value_, ok := value.(goja.Value); ok {
			return value_
		} else {
			return self.Environment.Runtime.ToValue(value)
		}
	} else {
		return nil
	}
}
