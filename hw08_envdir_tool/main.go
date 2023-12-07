package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	if len(os.Args) < 3 {
		fmt.Println("go-envdir usage: go-envdir dir cmd")
		os.Exit(1)
	}

	os.Exit(EnvDir(os.Args[1], os.Args[2:]...))
}

func EnvDir(dir string, args ...string) (returnCode int) {
	env, err := ReadDir(dir)
	if err != nil {
		fmt.Printf("%v", err)
		return 1
	}
	return RunCmd(args, env)
}
