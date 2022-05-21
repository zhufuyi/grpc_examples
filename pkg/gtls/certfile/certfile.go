package certfile

import (
	"path/filepath"
	"runtime"
)

var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path 返回绝对路径
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

/*
useage:

improt "grpc_examples/pkg/gtls/certfile"

filepath=certfile.Path("one-way/server.crt")
*/
