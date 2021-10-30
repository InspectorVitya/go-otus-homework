package main

import (
	"bufio"
	"bytes"
	"io"
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
	filesEnv, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := make(Environment, len(filesEnv))

	for _, file := range filesEnv {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			continue
		}

		openFile, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		r := bufio.NewReader(openFile)
		l, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF { //nolint:errorlint
				res[file.Name()] = EnvValue{Value: "", NeedRemove: true}
				continue
			}
			return nil, err
		}
		openFile.Close()
		str := bytes.TrimRight(l, " \t")
		str = bytes.ReplaceAll(str, []byte("\x00"), []byte("\n"))
		if len(str) == 0 {
			res[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		res[file.Name()] = EnvValue{Value: string(str)}
	}

	return res, nil
}
