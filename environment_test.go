package commonjs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonjs-goja/api"
	"github.com/tliron/exturl"

	_ "github.com/tliron/commonlog/simple"
)

func TestEnvironment(t *testing.T) {
	urlContext := exturl.NewContext()
	defer urlContext.Release()

	path := filepath.Join(getRoot(t), "examples")

	environment := commonjs.NewEnvironment(urlContext, urlContext.NewFileURL(path))
	defer environment.Release()

	environment.Extensions = api.DefaultExtensions{}.Create()

	testEnvironment(t, environment)
}

func testEnvironment(t *testing.T, environment *commonjs.Environment) {
	// Start!
	if _, err := environment.Require("./start", false, nil); err != nil {
		t.Errorf("%s", err)
	}
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
