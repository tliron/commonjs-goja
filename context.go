package commonjs

import (
	contextpkg "context"

	"github.com/dop251/goja"
	"github.com/tliron/exturl"
)

//
// Context
//

type Context struct {
	Environment *Environment
	Parent      *Context
	Module      *Module
	Resolve     ResolveFunc
	Extensions  []goja.Value
}

func (self *Environment) NewContext(url exturl.URL, parent *Context) *Context {
	jsContext := Context{
		Environment: self,
		Parent:      parent,
		Module:      self.NewModule(),
	}

	if url != nil {
		self.AddModule(url, jsContext.Module)
	}

	jsContext.Resolve = self.CreateResolver(url, &jsContext)

	for _, extension := range self.Extensions {
		jsContext.AppendExtension(extension)
	}

	// See: https://nodejs.org/api/modules.html#modules_the_module_object
	// See: https://nodejs.org/api/modules.html#modules_require_id

	jsContext.Module.Require = self.NewRequire(&jsContext)

	if parent != nil {
		parent.Module.Children = append(parent.Module.Children, jsContext.Module)
	}

	return &jsContext
}

func (self *Context) ResolveAndWatch(context contextpkg.Context, id string, bareId bool) (exturl.URL, error) {
	if url, err := self.Resolve(context, id, bareId); err == nil {
		// If it's a file, add to watch
		if fileUrl, ok := url.(*exturl.FileURL); ok {
			if err := self.Environment.Watch(fileUrl.Path); err != nil {
				return nil, err
			}
		}

		return url, nil
	} else {
		return nil, err
	}
}

func (self *Context) RequireAndExport(context contextpkg.Context, url exturl.URL, exportName string) (value any, childJsContext *Context, err error) {
	defer func() {
		if err_ := HandleJavaScriptPanic(recover()); err_ != nil {
			err = err_
		}
	}()

	environment := self.Environment.NewChild()
	jsContext := environment.NewContext(url, self)

	if exports, err := environment.cachedRequireUrl(context, url, jsContext); err == nil {
		if exportName == "" {
			return exports.Export(), jsContext, nil
		} else {
			return exports.Get(exportName).Export(), jsContext, nil // Get can panic
		}
	} else {
		return nil, nil, err
	}
}
