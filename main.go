package main

import (
	"fmt"

	"github.com/adasari/go-examples/shell"
)

func main() {
	try()
}

func try() error {
	cfg := &shell.Config{
		Command: "echo",
		Args:    []string{"1234"},
	}

	cmd := shell.New(cfg)

	// there are few scenarios the command execution has to mocked and skipped.
	stdout, stderr, err := cmd.Run()
	fmt.Printf("stdout: %v\n", stdout)
	fmt.Printf("stderr: %v\n", stderr)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		return err
	}

	return nil
}
