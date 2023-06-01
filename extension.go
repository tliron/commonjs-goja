package commonjs

import "github.com/dop251/goja"

//
// Extension
//

type Extension struct {
	Name   string
	Create CreateExtensionFunc
}

type CreateExtensionFunc func(context *Context) goja.Value
