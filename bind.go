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
		var err error
		if value, jsContext, err = bind.Unbind(); err == nil {
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
	value   any
	context *Context
	err     error
}

// Will attempt to immediately require the id and export the result, storing the
// result in an [EarlyBind].
func (self *Context) NewEarlyBind(id string, exportName string) (EarlyBind, error) {
	context, cancelContext := self.Environment.NewTimeoutContext()
	defer cancelContext()

	var earlyBind EarlyBind

	var url exturl.URL
	if url, earlyBind.err = self.ResolveAndWatch(context, id, false); earlyBind.err == nil {
		earlyBind.value, earlyBind.context, earlyBind.err = self.RequireAndExport(context, url, exportName)
	}

	return earlyBind, earlyBind.err
}

// ([Bind] interface)
func (self EarlyBind) Unbind() (any, *Context, error) {
	if self.err == nil {
		return self.value, self.context, nil
	} else {
		return nil, nil, self.err
	}
}

//
// LateBind
//

type LateBind struct {
	EarlyBind

	url        exturl.URL
	exportName string

	unbound bool
}

// Will resolve the id and store the URL and exportName in a [LateBind]. Only when
// [LateBind.Unbind] is called will require the URL and export the result (and cache
// the return values).
func (self *Context) NewLateBind(id string, exportName string) (LateBind, error) {
	context, cancelContext := self.Environment.NewTimeoutContext()
	defer cancelContext()

	if url, err := self.ResolveAndWatch(context, id, false); err == nil {
		return LateBind{
			url:        url,
			exportName: exportName,
			EarlyBind: EarlyBind{
				context: self,
			},
		}, nil
	} else {
		return LateBind{}, err
	}
}

// ([Bind] interface)
func (self LateBind) Unbind() (any, *Context, error) {
	if self.unbound {
		return self.EarlyBind.Unbind()
	}

	self.unbound = true

	context, cancelContext := self.context.Environment.NewTimeoutContext()
	defer cancelContext()

	self.value, self.context, self.err = self.context.RequireAndExport(context, self.url, self.exportName)
	return self.value, self.context, self.err
}
