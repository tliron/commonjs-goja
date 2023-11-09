package api

import (
	"github.com/tliron/commonjs-goja"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateLateBindExtension(jsContext *commonjs.Context) any {
	// commonjs.BindFunc signature
	return func(id string, exportName string) (any, error) {
		if lateBind, err := jsContext.NewLateBind(id, exportName); err == nil {
			return lateBind, nil
		} else {
			return nil, err
		}
	}
}
