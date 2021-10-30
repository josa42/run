package run

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/josa42/run/pkg/prefixwriter"
)

var _ Step = &RunStep{}

var index = 0

type RunStep struct {
	// - shell: bash
	//   script: |
	//     echo "Hello World!"
	Script string `yaml:"run"`
	Shell  string `yaml:"shell"`
	RunIn  string `yaml:"run-in"`
	dir    string `yaml:"-"`
}

func (c *RunStep) SetDir(dir string) {
	c.dir = dir
}

func (c RunStep) Run(tasks Tasks) (chan struct{}, CancelFunc) {
	file, err := ioutil.TempFile("", "run-script")
	if err != nil {
		log.Fatal(err)
	}

	filePath := file.Name()

	shell := c.Shell
	if shell == "" {
		shell = "bash"
	}

	script := c.Script

	ioutil.WriteFile(filePath, []byte(script), 0777)

	index += 1
	prefix := fmt.Sprintf("[%d]", index)

	lines := strings.Split(script, "\n")
	cmd := lines[0]
	if len(lines) > 1 {
		cmd += " (...)"
	}
	fmt.Print(color.BlueString("%s %s\n", prefix, cmd))

	done := make(chan struct{})
	execDone, cancel := c.exec(prefix, shell, filePath)

	go func() {
		<-execDone
		os.Remove(file.Name())
		close(done)
	}()

	return done, cancel
}

func (c RunStep) exec(prefix string, command string, args ...string) (chan struct{}, CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, command, args...)

	if c.RunIn != "" {
		cmd.Dir = path.Join(c.dir, c.RunIn)
	} else {
		cmd.Dir = c.dir
	}

	stdout := prefixwriter.New(os.Stderr, color.GreenString(prefix))
	stderr := prefixwriter.New(os.Stdout, color.RedString(prefix))

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	done := make(chan struct{})

	go func() {
		cmd.Run()
		stdout.Close()
		stderr.Close()
		close(done)
	}()

	return done, func() { cancel() }
}
