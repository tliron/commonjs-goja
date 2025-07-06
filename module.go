package commonjs

import (
	"github.com/dop251/goja"
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

func (self *Environment) AddModule(module *Module) {
	self.Modules.Set(module.Id, module)
}
