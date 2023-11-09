package commonjs

import (
	"github.com/tliron/exturl"
)

type BindFunc func(id string, exportName string) (any, error)

// If value is a [Bind] will unbind it, recursively, and return the bound
// value and [Context].
//
// Otherwise will return the provided value and [Context].
func Unbind(value any, jsContext *Context) (any, *Context, error) {
	if bind, ok := value.(Bind); ok {
		if value, jsContext, err := bind.Unbind(); err == nil {
			return Unbind(value, jsContext)
		} else {
			return nil, nil, err
		}
	} else {
		return value, jsContext, nil
	}
}

//
// Bind
//

type Bind interface {
	Unbind() (any, *Context, error)
}

//
// EarlyBind
//

type EarlyBind struct {
	Value   any
	Context *Context
}

func (self *Context) NewEarlyBind(id string, exportName string) (EarlyBind, error) {
	context, cancelContext := self.Environment.NewTimeoutContext()
	defer cancelContext()

	if url, err := self.ResolveAndWatch(context, id, false); err == nil {
		if value, jsContext, err := self.RequireAndExport(context, url, exportName); err == nil {
			return EarlyBind{
				Value:   value,
				Context: jsContext,
			}, nil
		} else {
			return EarlyBind{}, err
		}
	} else {
		return EarlyBind{}, err
	}
}

// ([Bind] interface)
func (self EarlyBind) Unbind() (any, *Context, error) {
	return self.Value, self.Context, nil
}

//
// LateBind
//

type LateBind struct {
	URL        exturl.URL
	ExportName string
	Context    *Context
}

func (self *Context) NewLateBind(id string, exportName string) (LateBind, error) {
	context, cancelContext := self.Environment.NewTimeoutContext()
	defer cancelContext()

	if url, err := self.ResolveAndWatch(context, id, false); err == nil {
		return LateBind{
			URL:        url,
			ExportName: exportName,
			Context:    self,
		}, nil
	} else {
		return LateBind{}, err
	}
}

// ([Bind] interface)
func (self LateBind) Unbind() (any, *Context, error) {
	context, cancelContext := self.Context.Environment.NewTimeoutContext()
	defer cancelContext()

	return self.Context.RequireAndExport(context, self.URL, self.ExportName)
}
