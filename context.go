package commonjs

import (
	contextpkg "context"
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/tliron/exturl"
)

//
// Context
//

type Context struct {
	Environment *Environment
	Parent      *Context
	UserContext any
	Module      *Module
	Resolve     ResolveFunc
	Extensions  []goja.Value
}

func (self *Environment) NewContext(url exturl.URL, parent *Context, userContext any) *Context {
	jsContext := Context{
		Environment: self,
		Parent:      parent,
		UserContext: userContext,
		Module:      self.NewModule(),
	}

	jsContext.setUrl(url)

	for _, extension := range self.Extensions {
		jsContext.AppendExtension(extension)
	}

	// See: https://nodejs.org/api/modules.html#modules_the_module_object
	// See: https://nodejs.org/api/modules.html#modules_require_id

	jsContext.Module.Require = jsContext.NewRequire()

	if parent != nil {
		parent.Module.Children = append(parent.Module.Children, jsContext.Module)
		if userContext == nil {
			jsContext.UserContext = parent.UserContext
		}
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

func (self *Context) ResolveAndRequire(context contextpkg.Context, id string, bareId bool, childEnvironment bool, userContext any) (*goja.Object, *Context, error) {
	if url, err := self.Resolve(context, id, bareId); err == nil {
		return self.Require(context, url, childEnvironment, userContext)
	} else {
		return nil, nil, err
	}
}

func (self *Context) Require(context contextpkg.Context, url exturl.URL, childEnvironment bool, userContext any) (*goja.Object, *Context, error) {
	environment := self.Environment
	if childEnvironment {
		environment = environment.NewChild()
	}

	jsContext := environment.NewContext(url, self, userContext)

	if value, err := jsContext.require(context, url); err == nil {
		return value, jsContext, nil
	} else {
		return nil, nil, err
	}
}

func (self *Context) RequireAndExport(context contextpkg.Context, url exturl.URL, childEnvironment bool, userContext any, exportName string) (value any, jsContext *Context, err error) {
	defer func() {
		if err_ := HandleJavaScriptPanic(recover()); err_ != nil {
			err = err_
		}
	}()

	if exports, jsContext, err := self.Require(context, url, childEnvironment, userContext); err == nil {
		if exportName == "" {
			return exports.Export(), jsContext, nil
		} else {
			return exports.Get(exportName).Export(), jsContext, nil // Get can panic
		}
	} else {
		return nil, nil, err
	}
}

func (self *Context) require(context contextpkg.Context, url exturl.URL) (*goja.Object, error) {
	key := url.Key()

	// Try cache
	if exports, loaded := self.Environment.exportsCache.Load(key); loaded {
		// Cache hit
		return exports.(*goja.Object), nil
	} else {
		// Cache miss
		if exports, err := self.runModule(context, url); err == nil {
			if exports_, loaded := self.Environment.exportsCache.LoadOrStore(key, exports); loaded {
				// Cache hit
				return exports_.(*goja.Object), nil
			} else {
				// Cache miss
				return exports, nil
			}
		} else {
			return nil, err
		}
	}
}

func (self *Context) runModule(context contextpkg.Context, url exturl.URL) (*goja.Object, error) {
	if program, err := self.getModule(context, url); err == nil {
		if value, err := self.Environment.Runtime.RunProgram(program); err == nil {
			if call, ok := goja.AssertFunction(value); ok {
				// See: self.compile for arguments
				arguments := []goja.Value{
					self.Module.Exports,
					self.Module.Require,
					self.Environment.Runtime.ToValue(self.Module),
					self.Environment.Runtime.ToValue(self.Module.Filename),
					self.Environment.Runtime.ToValue(self.Module.Path),
				}

				arguments = append(arguments, self.Extensions...)

				if _, err := call(nil, arguments...); err == nil {
					return self.Module.Exports, nil
				} else {
					return nil, UnwrapJavaScriptException(err)
				}
			} else {
				// Should never happen
				return nil, fmt.Errorf("invalid module: %v", value)
			}
		} else {
			return nil, UnwrapJavaScriptException(err)
		}
	} else {
		return nil, UnwrapJavaScriptException(err)
	}
}

func (self *Context) getModule(context contextpkg.Context, url exturl.URL) (*goja.Program, error) {
	key := url.Key()

	// Try cache
	if program, loaded := self.Environment.programCache.Load(key); loaded {
		// Cache hit
		return program.(*goja.Program), nil
	} else {
		// Cache miss
		if program, err := self.compile(context, url); err == nil {
			if program_, loaded := self.Environment.programCache.LoadOrStore(key, program); loaded {
				// Cache hit
				return program_.(*goja.Program), nil
			} else {
				// Cache miss
				return program, nil
			}
		} else {
			return nil, err
		}
	}
}

func (self *Context) compile(context contextpkg.Context, url exturl.URL) (*goja.Program, error) {
	if script, err := exturl.ReadString(context, url); err == nil {
		// Precompile
		if self.Environment.Precompile != nil {
			if script, err = self.Environment.Precompile(url, script, self); err != nil {
				return nil, err
			}
		}

		// See: https://nodejs.org/api/modules.html#modules_the_module_wrapper
		var builder strings.Builder
		builder.WriteString("(function(exports, require, module, __filename, __dirname")
		for _, extension := range self.Environment.Extensions {
			builder.WriteString(", ")
			builder.WriteString(extension.Name)
		}
		builder.WriteString(") {\n")
		builder.WriteString(script)
		builder.WriteString("\n});")
		script = builder.String()
		//log.Infof("%s", script)

		return goja.Compile(url.String(), script, self.Environment.Strict)
	} else {
		return nil, err
	}
}

func (self *Context) setUrl(url exturl.URL) {
	if url != nil {
		self.Environment.AddModule(url, self.Module)
	}
	self.Resolve = self.Environment.CreateResolver(url, self)
}
