package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSrcFileNotFound       = errors.New("src file not found")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return ErrSrcFileNotFound
	}
	defer closeFile(src)

	fs, err := getFileSize(src)
	if err != nil {
		return err
	}

	if offset > fs {
		return ErrOffsetExceedsFileSize
	}

	bytesCopyTotal := fs - offset
	if limit > 0 {
		bytesCopyTotal = int64(math.Min(float64(bytesCopyTotal), float64(limit)))
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeFile(dst)

	err = copyFile(src, dst, offset, bytesCopyTotal)
	if err != nil {
		return err
	}

	return nil
}

func getFileSize(src *os.File) (int64, error) {
	fi, err := src.Stat()
	if err != nil {
		return 0, err
	}

	fs := fi.Size()

	if fs == 0 {
		return 0, ErrUnsupportedFile
	}

	return fs, nil
}

func copyFile(src *os.File, dst *os.File, offset int64, limit int64) error {
	var sb strings.Builder

	for cur := int64(1); cur <= limit; cur++ {
		_, err := src.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = io.CopyN(dst, src, 1)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		offset++

		progressBar(&sb, cur, limit)
	}

	return nil
}

func progressBar(sb *strings.Builder, cur int64, limit int64) {
	progress := int((float64(cur) / float64(limit)) * 100)

	bLen := 100
	bProgressLen := progress * bLen / 100

	sb.WriteString(strings.Repeat("-", bProgressLen))
	sb.WriteString(">")
	sb.WriteString(strings.Repeat("_", bLen-bProgressLen))

	fmt.Printf("\r %d B / %d B  [%s] %d%%", cur, limit, sb.String(), progress)
	if cur == limit {
		fmt.Println()
	}
	sb.Reset()
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Printf("Failed to close file: %s, error: %v", file.Name(), err)
	}
}
