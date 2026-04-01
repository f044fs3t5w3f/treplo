package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func IsDirectoryExistsAndWrible(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}
	if !info.IsDir() {
		return errors.New("file is not directory")
	}
	testFile := filepath.Join(path, ".perm_check")
	f, err := os.OpenFile(testFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("test file was not written: %w", err)
	}
	f.Close()

	return os.Remove(testFile)
}
