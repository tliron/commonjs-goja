package commonjs

import (
	"github.com/dop251/goja"
	"github.com/tliron/exturl"
)

//
// Module
//

type Module struct {
	Id           string
	Children     []*Module
	Filename     string
	Path         string
	Paths        []string
	Exports      *goja.Object
	Require      *goja.Object
	IsPreloading bool
	Loaded       bool
}

func (self *Environment) NewModule() *Module {
	var path []string
	for _, url := range self.BasePaths {
		path = append(path, url.String())
	}

	return &Module{
		Paths:        path,
		Exports:      self.Runtime.NewObject(),
		IsPreloading: true,
	}
}

func (self *Environment) AddModule(url exturl.URL, module *Module) {
	module.Id = url.Key()
	module.IsPreloading = false
	module.Loaded = true
	if fileUrl, ok := url.(*exturl.FileURL); ok {
		module.Filename = fileUrl.Path
		if fileBase, ok := fileUrl.Base().(*exturl.FileURL); ok {
			module.Path = fileBase.Path
		}

		if err := self.Watch(module.Filename); err != nil {
			self.Log.Errorf("%s", err.Error())
		}
	}

	self.Modules.Set(module.Id, module)
}
