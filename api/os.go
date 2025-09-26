package api

import (
	contextpkg "context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/exturl"
	"github.com/tliron/go-kutil/util"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateOSExtension(jsContext *commonjs.Context) any {
	return NewOS(jsContext.Environment.URLContext)
}

//
// OS
//

type OS struct {
	Stdout io.Writer
	Stderr io.Writer

	urlContext *exturl.Context
}

func NewOS(urlContext *exturl.Context) OS {
	return OS{
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		urlContext: urlContext,
	}
}

func (self OS) JoinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func (self OS) Exec(name string, arguments ...string) (string, error) {
	cmd := exec.Command(name, arguments...)
	if out, err := cmd.Output(); err == nil {
		return util.BytesToString(out), nil
	} else if err_, ok := err.(*exec.ExitError); ok {
		return "", fmt.Errorf("%s\n%s", err_.Error(), util.BytesToString(err_.Stderr))
	} else {
		return "", err
	}
}

func (self OS) TemporaryFile(pattern string, directory string) (string, error) {
	if file, err := os.CreateTemp(directory, pattern); err == nil {
		name := file.Name()
		os.Remove(name)
		return name, nil
	} else {
		return "", err
	}
}

func (self OS) TemporaryDirectory(pattern string, directory string) (string, error) {
	return os.MkdirTemp(directory, pattern)
}

func (self OS) Download(sourceUrl string, targetPath string, timeoutSeconds float64) error {
	if sourceUrl_, err := self.urlContext.NewURL(sourceUrl); err == nil {
		context := contextpkg.Background()
		if timeoutSeconds > 0.0 {
			var cancelContext contextpkg.CancelFunc
			context, cancelContext = contextpkg.WithTimeout(context, time.Duration(timeoutSeconds*float64(time.Second)))
			defer cancelContext()
		}

		return exturl.DownloadTo(context, sourceUrl_, targetPath)
	} else {
		return err
	}
}
