package file_test

import (
	"os"
	"testing"

	"github.com/valerykalashnikov/streaming-pipeline/file"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	filename := "1234"
	defer func() {
		os.Remove(filename)
	}()
	fileSize := 2 * 1000 * 1000
	err := file.Generate(filename, fileSize)

	require.NoError(t, err)

	f, err := os.Stat(filename)
	require.NoError(t, err)
	size := f.Size()
	assert.Equal(t, int64(20000), size/100) // divide to 100 to remove fully remove bytes fluctuations leading to false negatives
}

func TestIOReadDir(t *testing.T) {
	err := os.MkdirAll("./test", os.ModePerm)
	require.NoError(t, err)
	filename := "./test/1234"
	defer func() {
		os.Remove(filename)
		os.Remove("./test")
	}()

	// 11 bytes is a maximum length of a generated line in a file
	// so let's generate file with only 1 line
	fileSize := 11
	err = file.Generate(filename, fileSize)
	require.NoError(t, err)
	fileList, err := file.IOReadDir("./test")
	require.NoError(t, err)
	require.NotEmpty(t, fileList)
	assert.Equal(t, "1234", fileList[0])
}
