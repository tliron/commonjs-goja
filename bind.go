package commonjs

import (
	"github.com/tliron/exturl"
)

type BindFunc func(id string, exportName string) (any, error)

// If value is a [Bind] will unbind it, recursively, and return the bound
// value and [Context].
//
// Otherwise will return the provided value and [Context] as is.
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
	Error   error
}

// Will attempt to immediately require the id and export the result, storing the
// result in an [EarlyBind].
func (self *Context) NewEarlyBind(id string, exportName string) (EarlyBind, error) {
	context, cancelContext := self.Environment.NewTimeoutContext()
	defer cancelContext()

	var earlyBind EarlyBind

	var url exturl.URL
	if url, earlyBind.Error = self.ResolveAndWatch(context, id, false); earlyBind.Error == nil {
		earlyBind.Value, earlyBind.Context, earlyBind.Error = self.RequireAndExport(context, url, exportName)
	}

	return earlyBind, earlyBind.Error
}

// ([Bind] interface)
func (self EarlyBind) Unbind() (any, *Context, error) {
	if self.Error == nil {
		return self.Value, self.Context, nil
	} else {
		return nil, nil, self.Error
	}
}

//
// LateBind
//

type LateBind struct {
	EarlyBind

	URL        exturl.URL
	ExportName string

	Unbound bool
}

// Will resolve the id and store the URL and exportName in a [LateBind]. Only when
// [LateBind.Unbind] is called will require the URL and export the result (and cache
// the return values).
func (self *Context) NewLateBind(id string, exportName string) (LateBind, error) {
	context, cancelContext := self.Environment.NewTimeoutContext()
	defer cancelContext()

	if url, err := self.ResolveAndWatch(context, id, false); err == nil {
		return LateBind{
			URL:        url,
			ExportName: exportName,
			EarlyBind: EarlyBind{
				Context: self,
			},
		}, nil
	} else {
		return LateBind{}, err
	}
}

// ([Bind] interface)
func (self LateBind) Unbind() (any, *Context, error) {
	if self.Unbound {
		return self.EarlyBind.Unbind()
	}

	self.Unbound = true

	context, cancelContext := self.Context.Environment.NewTimeoutContext()
	defer cancelContext()

	self.Value, self.Context, self.Error = self.Context.RequireAndExport(context, self.URL, self.ExportName)
	return self.Value, self.Context, self.Error
}
