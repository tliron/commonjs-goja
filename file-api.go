package commonjs

import (
	contextpkg "context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

type FileAPI struct {
	context *exturl.Context
}

func NewFileAPI(context *exturl.Context) FileAPI {
	return FileAPI{
		context: context,
	}
}

func (self FileAPI) JoinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func (self FileAPI) Exec(name string, arguments ...string) (string, error) {
	cmd := exec.Command(name, arguments...)
	if out, err := cmd.Output(); err == nil {
		return util.BytesToString(out), nil
	} else if err_, ok := err.(*exec.ExitError); ok {
		return "", fmt.Errorf("%s\n%s", err_.Error(), util.BytesToString(err_.Stderr))
	} else {
		return "", err
	}
}

func (self FileAPI) TemporaryFile(pattern string, directory string) (string, error) {
	if file, err := os.CreateTemp(directory, pattern); err == nil {
		name := file.Name()
		os.Remove(name)
		return name, nil
	} else {
		return "", err
	}
}

func (self FileAPI) TemporaryDirectory(pattern string, directory string) (string, error) {
	return os.MkdirTemp(directory, pattern)
}

func (self FileAPI) Download(sourceUrl string, targetPath string) error {
	if sourceUrl_, err := self.context.NewURL(sourceUrl); err == nil {
		return exturl.DownloadTo(contextpkg.TODO(), sourceUrl_, targetPath)
	} else {
		return err
	}
}
