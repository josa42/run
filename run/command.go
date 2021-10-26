package run

import (
	"os"
	"os/exec"
	"path"
)

var _ Step = &CommandStep{}

type CommandStep struct {
	// - command: echo
	//   args: ["Hello World!"]
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
	RunIn   string   `yaml:"run-in"`
	dir     string   `yaml:"-"`
}

func (c *CommandStep) SetDir(dir string) {
	c.dir = dir
}

func (c CommandStep) Run(tasks Tasks) {
	c.exec(c.Command, c.Args...)
}

func (c CommandStep) exec(command string, args ...string) {
	cmd := exec.Command(command, args...)
	if c.RunIn != "" {
		cmd.Dir = path.Join(c.dir, c.RunIn)
	} else {
		cmd.Dir = c.dir
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmd.Run()
}
