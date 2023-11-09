package api

import (
	"bytes"
	"io"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateConsoleExtension(jsContext *commonjs.Context) any {
	return NewConsole(jsContext)
}

//
// Console
//

// See: https://developer.mozilla.org/en-US/docs/Web/API/console

type Console struct {
	log     commonlog.Logger
	runtime *goja.Runtime
}

func NewConsole(jsContext *commonjs.Context) *Console {
	var keyValues []any
	if jsContext.Module != nil {
		keyValues = []any{"_scope", "console", "module", jsContext.Module.Id}
	} else {
		keyValues = []any{"_scope", "console"}
	}

	return &Console{
		log:     commonlog.NewKeyValueLogger(jsContext.Environment.Log, keyValues...),
		runtime: jsContext.Environment.Runtime,
	}
}

func (self *Console) Log(arguments ...any) {
	self.log_(commonlog.Notice, arguments...)
}

func (self *Console) Error(arguments ...any) {
	self.log_(commonlog.Error, arguments...)
}

func (self *Console) Warn(arguments ...any) {
	self.log_(commonlog.Warning, arguments...)
}

func (self *Console) Info(arguments ...any) {
	self.log_(commonlog.Info, arguments...)
}

func (self *Console) Debug(arguments ...any) {
	self.log_(commonlog.Debug, arguments...)
}

func (self *Console) Trace(arguments ...any) {
	var buffer bytes.Buffer
	for index, frame := range self.runtime.CaptureCallStack(0, nil) {
		if index == 0 {
			// Skip first frame
			continue
		}

		io.WriteString(&buffer, "\n")
		frame.Write(&buffer)
	}

	self.log_(commonlog.Notice, arguments...)
	self.log_(commonlog.Notice, buffer.String())
}

const LOG_DEPTH = 4 // skip through reflection to goja code

func (self *Console) log_(level commonlog.Level, arguments ...any) {
	length := len(arguments)
	if length == 0 {
		return
	} else if message, ok := arguments[0].(string); ok {
		if length == 1 {
			self.log.Log(level, LOG_DEPTH, message)
		} else {
			self.log.Logf(level, LOG_DEPTH, message, arguments[1:]...)
		}
	} else {
		for _, object := range arguments {
			self.log.Logf(level, LOG_DEPTH, "%T: %+v", object, object)
		}
	}
}
