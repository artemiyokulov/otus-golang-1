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
	testDataInputFilePath := "./testdata/input.txt"

	t.Run("success copy all file", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := testDataInputFilePath
		Copy(inputFilePath, outFile.Name(), 0, -1)
		require.True(t, fileAreEquals(inputFilePath, outFile.Name()))
	})

	t.Run("fail if offset > file size", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := testDataInputFilePath
		err = Copy(inputFilePath, outFile.Name(), getFileSize(inputFilePath)*2, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("fail if offset < 0", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := testDataInputFilePath
		err = Copy(inputFilePath, outFile.Name(), -1, 0)
		require.Equal(t, ErrOffsetBelowZero, err)
	})

	t.Run("/dev/urandom as input raise error", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "*")
		require.NoError(t, err)
		inputFilePath := "/dev/urandom"
		err = Copy(inputFilePath, outFile.Name(), 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("from-to same file", func(t *testing.T) {
		err := Copy("testdata/input.txt", "./testdata/input.txt", 0, 0)
		require.Equal(t, ErrSameFiles, err)
	})
}
