package config

import (
	"path/filepath"
	"runtime"
)

var AppRoot string

func InitAppRoot() {
	_, f, _, _ := runtime.Caller(0)
	AppRoot = filepath.Clean(filepath.Join(filepath.Dir(f), "../../../"))
}
