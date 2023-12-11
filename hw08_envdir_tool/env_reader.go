package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrEnvDirNotFound = errors.New("env dir not found")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var info fs.FileInfo
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrEnvDirNotFound
	}

	for _, file := range files {
		info, err = file.Info()
		if err != nil {
			return nil, err
		}

		if strings.Contains(info.Name(), "=") {
			continue
		}

		var val string

		if !isEmptyFile(info) {
			val, err = ReadFile(fmt.Sprintf("%s/%s", dir, info.Name()))
			if err != nil {
				return nil, err
			}
		}

		env[info.Name()] = EnvValue{
			Value:      val,
			NeedRemove: isEmptyFile(info),
		}
	}

	return env, nil
}

func ReadFile(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			fmt.Printf("file close error: %s", err)
		}
	}(f)

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		b := scanner.Bytes()
		b = bytes.TrimRight(b, " \t")
		b = bytes.ReplaceAll(b, []byte("\x00"), []byte("\n"))

		return string(b), err
	}

	return "", err
}

func isEmptyFile(info fs.FileInfo) bool {
	return info.Size() == 0
}
