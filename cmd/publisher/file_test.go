package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valerykalashnikov/streaming-pipeline/file"
)

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
	fileList, err := IOReadDir("./test")
	require.NoError(t, err)
	fmt.Println(fileList)
}
