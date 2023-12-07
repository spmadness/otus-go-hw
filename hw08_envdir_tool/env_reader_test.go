package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const envDir = "testdata/env"

var defaultEnvironmentData = Environment{
	"BAR": EnvValue{
		Value:      "bar",
		NeedRemove: false,
	},
	"EMPTY": EnvValue{
		Value:      "",
		NeedRemove: false,
	},
	"FOO": EnvValue{
		Value:      "   foo\nwith new line",
		NeedRemove: false,
	},
	"HELLO": EnvValue{
		Value:      "\"hello\"",
		NeedRemove: false,
	},
	"UNSET": EnvValue{
		Value:      "",
		NeedRemove: true,
	},
}

func TestReadDir(t *testing.T) {
	t.Run("ReadDir success: default case", func(t *testing.T) {
		expectedLen := 5

		env, err := ReadDir(envDir)

		require.Nilf(t, err, "expected nil, got err: %s", err)
		require.Equalf(t, expectedLen, len(env), "expected slice len: %d, actual: %d", expectedLen, len(env))

		require.Equal(t, defaultEnvironmentData, env)
	})

	t.Run("ReadDir success: envs with equal sign case", func(t *testing.T) {
		const dir = "testdata/envwithequal"

		err := os.Mkdir(dir, 0o755)
		if err != nil {
			t.Errorf("Mkdir failed: %v", err)
		}
		defer func() {
			err = os.RemoveAll(dir)
			if err != nil {
				t.Errorf("RemoveAll failed: %v", err)
			}
		}()

		fileNames := []string{"FOO", "FOO=", "BAR=", "FOO2=BAR", "BAR2=FOO", "=BAR2FOO"}

		for _, v := range fileNames {
			err = createFile(t, fmt.Sprintf("%s/%s", dir, v))
			if err != nil {
				t.Errorf("createFile failed: %v", err)
			}
		}

		expectedLen := 1

		env, err := ReadDir(dir)
		require.Nilf(t, err, "expected nil, got err: %s", err)
		require.Equalf(t, expectedLen, len(env), "expected slice len: %d, actual: %d", expectedLen, len(env))
		expectedData := Environment{
			"FOO": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		}

		require.Equal(t, expectedData, env)
	})

	t.Run("ReadDir fail: no env dir", func(t *testing.T) {
		env, err := ReadDir("testdata/fakedir")

		require.Truef(t, errors.Is(err, ErrEnvDirNotFound),
			"expected error: %v, actual err: %v", ErrEnvDirNotFound, err)
		require.Nil(t, nil, env)
	})
}

func createFile(t *testing.T, path string) error {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			t.Errorf("tmp file close failed: %v", err)
		}
	}(f)

	return nil
}
