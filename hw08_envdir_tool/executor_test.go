package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("RunCmd success: default case", func(t *testing.T) {
		cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}
		expectedOutput := strings.ReplaceAll(
			`HELLO is ("hello")
				BAR is (bar)
				FOO is (   foo
				with new line)
				UNSET is ()
				ADDED is (from original env)
				EMPTY is ()
				arguments are arg1=1 arg2=2
				`, "\t", "")

		setEnv(t)
		defer unsetEnv(t)
		// сохраняем дефолтный stdout в переменную
		stdout := os.Stdout
		r, w, _ := os.Pipe()
		// подменяем stdout, чтобы можно было захватить вывод функции RunCmd
		os.Stdout = w

		code := RunCmd(cmd, defaultEnvironmentData)

		err := w.Close()
		if err != nil {
			t.Errorf("%v", err)
		}
		actualOutput, _ := io.ReadAll(r)
		// возвращаем дефолтный stdout обратно
		os.Stdout = stdout

		require.Equal(t, 0, code)
		require.Equal(t, expectedOutput, string(actualOutput))

		_, exists := os.LookupEnv("UNSET")
		require.Equal(t, exists, false)

		require.Equal(t, os.Getenv("HELLO"), "\"hello\"")
		require.Equal(t, os.Getenv("BAR"), "bar")
		require.Equal(t, os.Getenv("FOO"), "   foo\nwith new line")
		require.Equal(t, os.Getenv("ADDED"), "from original env")
		require.Equal(t, os.Getenv("EMPTY"), "")
	})

	t.Run("RunCmd fail: no cmd given", func(t *testing.T) {
		cmd := []string{""}

		code := RunCmd(cmd, defaultEnvironmentData)
		require.Equal(t, 1, code)
	})

	t.Run("RunCmd: no cmd args given", func(t *testing.T) {
		cmd := []string{"/bin/bash"}

		code := RunCmd(cmd, defaultEnvironmentData)
		require.Equal(t, 0, code)
	})

	t.Run("RunCmd: no env given", func(t *testing.T) {
		cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}

		code := RunCmd(cmd, Environment{})
		require.Equal(t, 0, code)
	})
}

func setEnv(t *testing.T) {
	t.Helper()

	env := map[string]string{
		"HELLO": "SHOULD_REPLACE",
		"FOO":   "SHOULD_REPLACE",
		"UNSET": "SHOULD_REMOVE",
		"ADDED": "from original env",
		"EMPTY": "SHOULD_BE_EMPTY",
	}

	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			t.Errorf("Setenv failed: %s", err)
		}
	}
}

func unsetEnv(t *testing.T) {
	t.Helper()

	env := []string{"HELLO", "BAR", "FOO", "UNSET", "ADDED", "EMPTY"}

	for _, v := range env {
		err := os.Unsetenv(v)
		if err != nil {
			t.Errorf("Unsetenv failed: %s", err)
		}
	}
}
