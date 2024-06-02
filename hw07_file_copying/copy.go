package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported input file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFiles             = errors.New("input and output are same files")
)

func getFileSize(filepath string) int64 {
	fileInfo, err := os.Lstat(filepath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func validateInput(fromPath, toPath string, offset int64) error {
	if file, err := os.OpenFile(fromPath, os.O_RDONLY, 0); err != nil {
		return fmt.Errorf("can not read input file: %w", err)
	} else {
		if err := file.Close(); err != nil {
			return fmt.Errorf("check file read: close input file error: %w", err)
		}
	}
	inputFileSize := getFileSize(fromPath)
	if inputFileSize == 0 {
		return ErrUnsupportedFile
	}
	if fromPath == toPath {
		return ErrSameFiles
	}
	if offset >= inputFileSize {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateInput(fromPath, toPath, offset); err != nil {
		return err
	}
	inputFileSize := getFileSize(fromPath)
	if limit > inputFileSize || limit == 0 {
		limit = inputFileSize - offset
	}
	inputFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open input file error: %w", err)
	}
	if offset > 0 {
		inputFile.Seek(offset, 0)
	}
	outputFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("output file error: %w", err)
	}

	defer func() {
		inputFile.Close()
	}()
	defer func() {
		outputFile.Close()
	}()

	readPartSize := int64(1024)
	bar := pb.Start64(limit)
	buffer := make([]byte, readPartSize)

	for limit > 0 {
		readedBytesCount, err := inputFile.Read(buffer[:min(readPartSize, limit)])
		if err != nil && err != io.EOF {
			return fmt.Errorf("error while reading from input file: %w", err)
		}
		if err == io.EOF && readedBytesCount == 0 {
			break
		}
		_, err = outputFile.Write(buffer[:readedBytesCount])
		if err != nil && err != io.EOF {
			return fmt.Errorf("error while writing to output file: %w", err)
		}
		limit = limit - int64(readedBytesCount)
		bar.Add64(int64(readedBytesCount))
	}

	bar.Finish()
	return nil
}
