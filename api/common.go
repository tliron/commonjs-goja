package api

import (
	"io"

	"github.com/tliron/commonjs-goja"
)

//
// DefaultExtensions
//

type DefaultExtensions struct {
	LateBind  bool
	Arguments map[string]string
	Stdout    io.Writer
	Stderr    io.Writer
}

func (self DefaultExtensions) Create() []commonjs.Extension {
	var createBind commonjs.CreateExtensionFunc
	if self.LateBind {
		createBind = CreateLateBindExtension
	} else {
		createBind = CreateEarlyBindExtension
	}

	return []commonjs.Extension{{
		Name:   "bind",
		Create: createBind,
	}, {
		Name:   "console",
		Create: CreateConsoleExtension,
	}, {
		Name:   "env",
		Create: CreateEnvExtension(self.Arguments),
	}, {
		Name:   "util",
		Create: CreateUtilExtension,
	}, {
		Name:   "transcribe",
		Create: CreateTranscribeExtension(self.Stdout, self.Stderr),
	}, {
		Name:   "ard",
		Create: CreateARDExtension,
	}, {
		Name:   "os",
		Create: CreateOSExtension,
	}}
}
