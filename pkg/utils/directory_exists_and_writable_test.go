package utils

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDirectoryExistsAndWrible(t *testing.T) {
	currentDir := GetPackageDirecory()
	testCases := []struct {
		name      string
		path      string
		wantError bool
	}{
		{
			"Writable directory",
			"writable_directory",
			false,
		}, {
			"Nonwritable directory",
			"non_writable_directory",
			true,
		}, {
			"File",
			"file",
			true,
		},
	}
	for _, testCases := range testCases {
		t.Run(testCases.name, func(t *testing.T) {
			dir := path.Join(currentDir, "test", testCases.path)
			err := IsDirectoryExistsAndWrible(dir)
			if testCases.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
