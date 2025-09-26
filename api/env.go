package api

import (
	contextpkg "context"
	"io"
	"time"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-kutil/util"
)

func CreateEnvExtension(arguments map[string]string) commonjs.CreateExtensionFunc {
	return func(jsContext *commonjs.Context) any {
		return NewEnv(jsContext, arguments)
	}
}

var Variables = commonjs.NewThreadSafeObject()

//
// Env
//

type Env struct {
	Context   *commonjs.Context
	Variables *goja.Object
	Arguments map[string]string
	Log       commonlog.Logger
}

func NewEnv(jsContext *commonjs.Context, arguments map[string]string) *Env {
	if arguments == nil {
		arguments = make(map[string]string)
	}

	log := jsContext.Environment.Log
	if jsContext.Module != nil {
		log = commonlog.NewKeyValueLogger(log, "module", jsContext.Module.Id)
	}

	return &Env{
		Context:   jsContext,
		Variables: Variables.NewDynamicObject(jsContext.Environment.Runtime),
		Arguments: arguments,
		Log:       log,
	}
}

func (self *Env) LoadString(id string, timeoutSeconds float64) (string, error) {
	if bytes, err := self.LoadBytes(id, timeoutSeconds); err == nil {
		return util.BytesToString(bytes), nil
	} else {
		return "", err
	}
}

func (self *Env) LoadBytes(id string, timeoutSeconds float64) ([]byte, error) {
	context := contextpkg.Background()
	if timeoutSeconds > 0.0 {
		var cancelContext contextpkg.CancelFunc
		context, cancelContext = contextpkg.WithTimeout(context, time.Duration(timeoutSeconds*float64(time.Second)))
		defer cancelContext()
	}

	if url, err := self.Context.ResolveAndWatch(context, id, false); err == nil {
		return exturl.ReadBytes(context, url)
	} else {
		return nil, err
	}
}

func (self *Env) WriteFrom(writer io.Writer, id string, timeoutSeconds float64) error {
	context := contextpkg.Background()
	if timeoutSeconds > 0.0 {
		var cancelContext contextpkg.CancelFunc
		context, cancelContext = contextpkg.WithTimeout(context, time.Duration(timeoutSeconds*float64(time.Second)))
		defer cancelContext()
	}

	if url, err := self.Context.ResolveAndWatch(context, id, false); err == nil {
		if reader, err := url.Open(context); err == nil {
			if _, err := io.Copy(writer, reader); err == nil {
				return reader.Close()
			} else {
				commonlog.CallAndLogWarning(reader.Close, "Env.WriteFrom", self.Context.Environment.Log)
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
