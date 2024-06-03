package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var env = Environment{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := filepath.Base(path)
		if strings.Contains(fileName, "=") {
			return nil
		}

		env[fileName] = EnvValue{
			Value:      "",
			NeedRemove: info.Size() == 0,
		}

		if info.Size() != 0 {
			file, err := os.Open(path)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()

			if err != nil {
				log.Fatal(err)
			}

			scanner := bufio.NewScanner(file)
			var firstLine string
			for scanner.Scan() {
				firstLine = scanner.Text()
				break
			}

			value := strings.ReplaceAll(
				strings.TrimRight(firstLine, " \t\n"),
				"\x00",
				"\n")
			envValue := env[fileName]
			envValue.Value = value
			env[fileName] = envValue
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return env, nil
}
