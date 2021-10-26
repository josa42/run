package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

var _ Step = &RunStep{}

type RunStep struct {
	// - shell: bash
	//   script: |
	//     echo "Hello World!"
	RunScript string `yaml:"run"`
	Script    string `yaml:"script"`
	Shell     string `yaml:"shell"`
	RunIn     string `yaml:"run-in"`
	dir       string `yaml:"-"`
}

func (c *RunStep) SetDir(dir string) {
	c.dir = dir
}

func (c RunStep) Run(tasks Tasks) {
	file, err := ioutil.TempFile("", "run-script")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	filePath := file.Name()

	shell := c.Shell
	if shell == "" {
		shell = "bash"
	}

	script := c.RunScript
	if script == "" {
		// TODO echo deprecation
		script = c.Script
	}

	ioutil.WriteFile(filePath, []byte(script), 0777)

	c.exec(shell, filePath)
	return
}

func (c RunStep) exec(command string, args ...string) {
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
