package commonjs

import (
	contextpkg "context"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/fswatch"
)

const DEFAULT_TIMEOUT = time.Second * 5

//
// Environment
//

type Environment struct {
	Runtime        *goja.Runtime
	URLContext     *exturl.Context
	BasePaths      []exturl.URL
	Extensions     []Extension
	Modules        *goja.Object
	Precompile     PrecompileFunc
	CreateResolver CreateResolverFunc
	OnFileModified OnFileModifiedFunc
	Timeout        time.Duration
	Strict         bool
	Log            commonlog.Logger
	Lock           sync.Mutex

	watcher      *fswatch.Watcher
	watcherLock  sync.Mutex
	exportsCache sync.Map
	programCache *sync.Map
}

type PrecompileFunc func(url exturl.URL, script string, jsContext *Context) (string, error)

type OnFileModifiedFunc func(id string, module *Module)

func NewEnvironment(urlContext *exturl.Context, basePaths ...exturl.URL) *Environment {
	runtime := goja.New()
	runtime.SetFieldNameMapper(DromedaryCaseMapper)

	return &Environment{
		Runtime:        runtime,
		URLContext:     urlContext,
		BasePaths:      basePaths,
		Modules:        NewThreadSafeObject().NewDynamicObject(runtime),
		CreateResolver: NewDefaultResolverCreator("js", true, urlContext, basePaths...),
		Timeout:        DEFAULT_TIMEOUT,
		Strict:         true,
		Log:            log,
		programCache:   new(sync.Map),
	}
}

func (self *Environment) NewChild() *Environment {
	environment := NewEnvironment(self.URLContext, self.BasePaths...)
	environment.Extensions = self.Extensions
	environment.Precompile = self.Precompile
	environment.CreateResolver = self.CreateResolver
	environment.OnFileModified = self.OnFileModified
	environment.Timeout = self.Timeout
	environment.Strict = self.Strict
	environment.Log = self.Log
	environment.watcher = self.watcher
	environment.programCache = self.programCache
	return environment
}

func (self *Environment) StartWatcher() error {
	self.watcherLock.Lock()
	defer self.watcherLock.Unlock()

	if self.watcher != nil {
		if err := self.watcher.Close(); err == nil {
			self.watcher = nil
		} else {
			return err
		}
	}

	if self.OnFileModified == nil {
		return nil
	}

	var err error
	if self.watcher, err = fswatch.NewWatcher(self.URLContext); err == nil {
		self.watcher.Start(func(fileUrl *exturl.FileURL) {
			self.Lock.Lock()
			id := fileUrl.Key()
			var module *Module
			if module_ := self.Modules.Get(id); module_ != nil {
				module = module_.Export().(*Module)
			}
			self.Lock.Unlock()
			self.OnFileModified(id, module)
		})
		return nil
	} else {
		return err
	}
}

func (self *Environment) StopWatcher() error {
	self.watcherLock.Lock()
	defer self.watcherLock.Unlock()

	if self.watcher != nil {
		if err := self.watcher.Close(); err == nil {
			self.watcher = nil
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

func (self *Environment) Watch(path string) error {
	self.watcherLock.Lock()
	defer self.watcherLock.Unlock()

	if self.watcher != nil {
		return self.watcher.Add(path)
	} else {
		return nil
	}
}

func (self *Environment) Release() error {
	return self.StopWatcher()
}

func (self *Environment) NewTimeoutContext() (contextpkg.Context, contextpkg.CancelFunc) {
	return contextpkg.WithTimeout(contextpkg.Background(), self.Timeout)
}

func (self *Environment) Call(function any, this any, arguments ...any) (any, error) {
	//self.Lock.Lock()
	//defer self.Lock.Unlock()

	return Call(self.Runtime, function, this, arguments...)
}

func (self *Environment) GetAndCall(object *goja.Object, name string, this any, arguments ...any) (any, error) {
	//self.Lock.Lock()
	//defer self.Lock.Unlock()

	return GetAndCall(self.Runtime, object, name, this, arguments...)
}

func (self *Environment) ClearCache() {
	self.exportsCache.Range(func(key any, value any) bool {
		self.exportsCache.Delete(key)
		return true
	})
	self.programCache.Range(func(key any, value any) bool {
		self.programCache.Delete(key)
		return true
	})
	self.Modules = NewThreadSafeObject().NewDynamicObject(self.Runtime)
}

func (self *Environment) Require(id string, bareId bool, userContext any) (*goja.Object, error) {
	context, cancelContext := self.NewTimeoutContext()
	defer cancelContext()

	jsContext := self.NewContext(nil, nil, userContext)
	if url, err := jsContext.Resolve(context, id, bareId); err == nil {
		jsContext.setUrl(url)
		return jsContext.require(context, url)
	} else {
		return nil, err
	}
}

func (self *Environment) RequireURL(url exturl.URL, userContext any) (*goja.Object, error) {
	context, cancelContext := self.NewTimeoutContext()
	defer cancelContext()

	return self.NewContext(url, nil, userContext).require(context, url)
}
