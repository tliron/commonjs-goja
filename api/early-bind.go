package api

import (
	"github.com/tliron/commonjs-goja"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateEarlyBindExtension(jsContext *commonjs.Context) any {
	// commonjs.BindFunc signature
	return func(id string, exportName string) (any, error) {
		if earlyBind, err := jsContext.NewEarlyBind(id, exportName); err == nil {
			return earlyBind, nil
		} else {
			return nil, err
		}
	}
}
