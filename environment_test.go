package commonjs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
	"github.com/tliron/exturl"
)

func TestEnvironment(t *testing.T) {
	urlContext := exturl.NewContext()
	defer urlContext.Release()

	path := filepath.Join(getRoot(t), "assets")

	environment := NewEnvironment(urlContext, []exturl.URL{urlContext.NewFileURL(path)})
	defer environment.Release()

	environment.Extensions = append(environment.Extensions, Extension{
		Name: "console",
		Create: func(context *Context) goja.Value {
			return context.Environment.Runtime.ToValue(consoleAPI{})
		},
	})

	// Support for "bind" (late binding)
	environment.Extensions = append(environment.Extensions, Extension{
		Name:   "bind",
		Create: CreateLateBindExtension,
	})

	// Start!
	if _, err := environment.RequireID("start"); err != nil {
		t.Errorf("%s", err)
	}
}

type consoleAPI struct{}

func (self consoleAPI) Log(message string) {
	fmt.Println(message)
}

func getRoot(t *testing.T) string {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("COMMONJS_TEST_ROOT"); !ok {
		var err error
		if root, err = os.Getwd(); err != nil {
			t.Errorf("os.Getwd: %s", err.Error())
		}
	}
	return root
}
