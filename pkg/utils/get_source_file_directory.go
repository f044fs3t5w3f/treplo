package utils

import (
	"path/filepath"
	"runtime"
)

func GetPackageDirecory() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Failed to retrieve caller information")
	}

	return filepath.Dir(filename)
}
