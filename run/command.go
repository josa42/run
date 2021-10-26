package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

var _ Step = &Command{}

type Command struct {
	RunIn string `yaml:"run_in"`
	Dir   string `yaml:"-"`

	// Shell:
	// - shell: bash
	//   script: |
	//     echo "Hello World!"
	Shell  string `yaml:"shell"`
	Script string `yaml:"script"`

	// Command:
	// - shell: echo
	//   args: ["Hello World!"]
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`

	// Task:
	// - task: <task-name>
	Task string
}

func (c *Command) SetDir(dir string) {
	c.Dir = dir
}

func (c Command) Run(tasks Tasks) {
	if c.Command != "" {
		c.exec(c.Command, c.Args...)
		return
	}

	if c.Shell != "" {
		file, err := ioutil.TempFile("", "run-script")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())

		filePath := file.Name()

		ioutil.WriteFile(filePath, []byte(c.Script), 0777)

		c.exec(c.Shell, filePath)
		return
	}

	if c.Task != "" {
		if task, ok := tasks[c.Task]; ok {
			task.Run(tasks)
			return
		}
	}
}

func (c Command) exec(command string, args ...string) {
	cmd := exec.Command(command, args...)
	if c.RunIn != "" {
		cmd.Dir = path.Join(c.Dir, c.RunIn)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmd.Run()
}
