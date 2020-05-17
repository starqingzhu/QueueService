package preload

import (
	"path/filepath"
)

var basePath = ".."

var confPath string

func init() {
	confPath, _ = filepath.Abs(basePath + "/conf")
}
