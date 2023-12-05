package main

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const dstFilePath = "/tmp/out.txt"

func TestCopy(t *testing.T) {
	t.Run("copy success: simple case", func(t *testing.T) {
		limit = int64(1000)
		defer removeFile(t, dstFilePath)

		err := Copy("testdata/input.txt", dstFilePath, 0, limit)

		require.Nilf(t, err, "expected nil, actual err: %v", err)
		require.FileExistsf(t, dstFilePath, "file not found: %s", dstFilePath)

		f, err := os.Open(dstFilePath)
		if err != nil {
			t.Error(err)
		}

		fi, err := f.Stat()
		if err != nil {
			t.Error(err)
		}
		fsize := fi.Size()

		require.Truef(t, fsize == limit, "actual file size (%d) doesn't match expected (%d)", fsize, limit)
	})

	t.Run("copy success: eof case", func(t *testing.T) {
		// testdata/input.txt size is 6617 bytes
		expectedDstFileSize := int64(17)
		limit = int64(1000)
		defer removeFile(t, dstFilePath)

		err := Copy("testdata/input.txt", dstFilePath, 6600, limit)

		require.Nilf(t, err, "expected nil, actual err: %v", err)
		require.FileExistsf(t, dstFilePath, "file not found: %s", dstFilePath)

		f, err := os.Open(dstFilePath)
		if err != nil {
			t.Error(err)
		}

		fi, err := f.Stat()
		if err != nil {
			t.Error(err)
		}
		fsize := fi.Size()

		require.Truef(t, fsize == expectedDstFileSize,
			"actual file size (%d) doesn't match expected (%d)", fsize, expectedDstFileSize)
	})

	t.Run("copy fail: no src file", func(t *testing.T) {
		err := Copy("testdata/input123.txt", dstFilePath, 0, 10)

		require.Truef(t, errors.Is(err, ErrSrcFileNotFound),
			"expected error: %v, actual err: %v", ErrSrcFileNotFound, err)
	})

	t.Run("copy fail: empty src file", func(t *testing.T) {
		f, err := os.CreateTemp("", "test")
		if err != nil {
			log.Fatal(err)
		}
		defer removeFile(t, f.Name())

		err = Copy(f.Name(), dstFilePath, 0, 10)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile),
			"expected error: %v, actual err: %v", ErrUnsupportedFile, err)
	})

	t.Run("copy fail: offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", dstFilePath, 15000, 0)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize),
			"expected error: %v, actual err: %v", ErrOffsetExceedsFileSize, err)
	})
}

func removeFile(t *testing.T, name string) {
	t.Helper()
	err := os.Remove(name)
	if err != nil {
		t.Errorf("temp file remove failed: %v", err)
	}
}
