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
	assert.Equal(t, int64(200000), size/10) // divide to 10 to remove some bytes fluctuations
}
