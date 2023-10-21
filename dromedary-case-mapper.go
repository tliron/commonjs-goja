package commonjs

import (
	"reflect"

	"github.com/tliron/kutil/util"
)

var DromedaryCaseMapper dromedaryCaseMapper

type dromedaryCaseMapper struct{}

// ([goja.FieldNameMapper] interface)
func (self dromedaryCaseMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return util.ToDromedaryCase(f.Name)
}

// ([goja.FieldNameMapper] interface)
func (self dromedaryCaseMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return util.ToDromedaryCase(m.Name)
}
