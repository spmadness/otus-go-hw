package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var args []string

	if len(cmd) == 0 {
		return 1
	}

	if len(cmd) > 1 {
		args = cmd[1:]
	}

	c := exec.Command(cmd[0], args...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	for k, v := range env {
		envV := v.Value

		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return 1
			}
			continue
		}

		err := os.Setenv(k, envV)
		if err != nil {
			return 1
		}
	}

	err := c.Run()
	if err != nil {
		return 1
	}

	return 0
}
