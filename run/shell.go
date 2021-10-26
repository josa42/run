package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

var _ Step = &ShellStep{}

type ShellStep struct {
	// - shell: bash
	//   script: |
	//     echo "Hello World!"
	Shell  string `yaml:"shell"`
	Script string `yaml:"script"`
	RunIn  string `yaml:"run-in"`
	dir    string `yaml:"-"`
}

func (c *ShellStep) SetDir(dir string) {
	c.dir = dir
}

func (c ShellStep) Run(tasks Tasks) {
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

func (c ShellStep) exec(command string, args ...string) {
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
