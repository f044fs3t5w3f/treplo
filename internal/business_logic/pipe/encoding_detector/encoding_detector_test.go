package encoding_detector

import (
	"os"
	"path"
	"testing"

	"github.com/a-kuleshov/treplo/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestDetectEncoding_mp3(t *testing.T) {
	directory := utils.GetPackageDirecory()
	file, err := os.Open(path.Join(directory, "samples/test_mp3.mp3"))
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer file.Close()
	encoding, err := detectEncoding(file)
	assert.NoError(t, err)
	assert.Equal(t, encodingMP3, encoding)
}
func TestDetectEncoding_mp3_fail(t *testing.T) {
	directory := utils.GetPackageDirecory()
	file, err := os.Open(path.Join(directory, "samples/test_nomedia.txt"))
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer file.Close()
	encoding, err := detectEncoding(file)
	assert.Error(t, err)
	assert.Equal(t, "", encoding)
}
