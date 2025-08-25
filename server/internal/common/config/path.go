package config

import (
	"path/filepath"
	"runtime"
)

func AppRoot() string {
	_, f, _, _ := runtime.Caller(0)
	appRoot := filepath.Clean(filepath.Join(filepath.Dir(f), "../../"))
	return appRoot
}
