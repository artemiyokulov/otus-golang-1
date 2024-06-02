package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func fileAreEquals(filePath1, filePath2 string) bool {
	f1, _ := os.ReadFile(filePath1)
	f2, _ := os.ReadFile(filePath2)
	return bytes.Equal(f1, f2)
}

func TestCopy(t *testing.T) {

	t.Run("success copy", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := "./testdata/input.txt"
		Copy(inputFilePath, outFile.Name(), 0, 0)
		require.True(t, fileAreEquals(inputFilePath, outFile.Name()))
	})

	t.Run("fail if offset > file size", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := "./testdata/input.txt"
		err = Copy(inputFilePath, outFile.Name(), getFileSize(inputFilePath)*2, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("/dev/urandom as input raise error", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := "/dev/urandom"
		err = Copy(inputFilePath, outFile.Name(), 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("from-to same file", func(t *testing.T) {
		inputFilePath := "./testdata/input.txt"
		err := Copy(inputFilePath, inputFilePath, 0, 0)
		require.Equal(t, ErrSameFiles, err)
	})
}
