package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// Config .
type Config struct {
	Command           string            // The command to run
	Args              []string          // The args to pass to the command
	WorkingDir        string            // The working directory
	Env               map[string]string // Additional environment variables to set
	OutputMaxLineSize int               // Max output line size
}

type Command struct {
	*Config
}

// New creates command instance.
func New(cfg *Config) *Command {
	return &Command{Config: cfg}
}

// Run excutes the command return stdout, stderr log.
func (c Command) Run() (string, string, error) {
	var stdout []string
	var stderr []string

	cmd := exec.Command(c.Command, c.Args...)

	cmd.Dir = c.WorkingDir
	cmd.Env = envVars(c.Config)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	stdoutScanner := Scanner(stdoutPipe, c.OutputMaxLineSize)
	stderrScanner := Scanner(stderrPipe, c.OutputMaxLineSize)

	if err = cmd.Start(); err != nil {
		return "", "", err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go readScanner(stdoutScanner, wg, &stdout)
	go readScanner(stderrScanner, wg, &stderr)
	wg.Wait()

	if err = stdoutScanner.Err(); err != nil {
		return "", "", err
	}

	if err = stderrScanner.Err(); err != nil {
		return "", "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	return strings.Join(stdout, "\n"), strings.Join(stderr, "\n"), nil
}

// String returns command.
func (c Command) String() string {
	cmd := exec.Command(c.Command, c.Args...)
	cmd.Dir = c.WorkingDir
	cmd.Env = envVars(c.Config)

	return cmd.String()
}

// Scanner returns scanner object. scanner buffer set to given maxLineSize if maxLineSize > 0.
func Scanner(reader io.ReadCloser, maxLineSize int) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	if maxLineSize > 0 {
		scanner.Buffer(make([]byte, maxLineSize), maxLineSize)
	}
	return scanner
}

// readScanner reads the scanner and append the given out.
func readScanner(scanner *bufio.Scanner, wg *sync.WaitGroup, out *[]string) {
	defer wg.Done()
	for scanner.Scan() {
		*out = append(*out, scanner.Text())
	}
}

func envVars(cfg *Config) []string {
	var env []string
	for key, value := range cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

/*
type Commander interface {
	String() string
	Run() (string, string, error)
}

// New creates command instance.
func New(cfg *Config) Commander {
	return impl(cfg)
}

func Fake(stdout, stderr string, err error) {
	impl = func(cfg *Config) Commander {
		return &fakeCommand{Config: cfg, stdout: stdout, stderr: stderr, err: err}
	}
}

type fakeCommand struct {
	*Config
	stdout string
	stderr string
	err    error
}

func (c fakeCommand) String() string {
	return ""
}

func (c fakeCommand) Run() (string, string, error) {
	if c.stdout != "" || c.stderr != "" || c.err != nil {
		return c.stdout, c.stderr, c.err
	}

	// execute command.
	var (
		stdout []string
		stderr []string
	)

	cmd := exec.Command(c.Command, c.Args...)
	cmd.Dir = c.WorkingDir
	cmd.Env = envVars(c.Config)

	execError := ExecError{
		Cmd:        c.Command + " " + strings.Join(c.Args, " "),
		WorkingDir: c.WorkingDir,
		Env:        strings.Join(cmd.Env, " "),
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		execError.Err = err
		return "", "", execError
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		execError.Err = err
		return "", "", execError
	}

	stdoutScanner := Scanner(stdoutPipe, c.OutputMaxLineSize)
	stderrScanner := Scanner(stderrPipe, c.OutputMaxLineSize)

	if err = cmd.Start(); err != nil {
		execError.Err = err
		return "", "", execError
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go readScanner(stdoutScanner, wg, &stdout)
	go readScanner(stderrScanner, wg, &stderr)

	wg.Wait()

	if err = stdoutScanner.Err(); err != nil {
		execError.Err = err
		return "", "", execError
	}

	if err = stderrScanner.Err(); err != nil {
		execError.Err = err
		return "", "", execError
	}

	if err := cmd.Wait(); err != nil {
		execError.Stdout = strings.Join(stdout, "\n")
		execError.Stderr = strings.Join(stderr, "\n")
		execError.Err = err

		return "", "", execError
	}

	return strings.Join(stdout, "\n"), strings.Join(stderr, "\n"), nil
}

var impl = defaultImpl

var defaultImpl = func(cfg *Config) Commander {
	return &Command{Config: cfg}
}

func Reset() {
	impl = defaultImpl
}


// Error implements the error interface and returns a description of the error
func (e ExecError) Error() string {
	return fmt.Sprintf("command and error details are - \n cmd: %s\n working dir: %s\n, env: %s\n stdout: %s\n stderr: %s\n error: %v",
		e.Cmd, e.WorkingDir, e.Env, e.Stdout, e.Stderr, e.Err)
}
*/
