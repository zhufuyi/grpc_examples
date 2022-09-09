package include

import (
	"path/filepath"
	"runtime"
)

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

// Path 返回绝对路径
func Path(relative string) string {
	if filepath.IsAbs(relative) {
		return relative
	}

	return filepath.Join(basePath, relative)
}
