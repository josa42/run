package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Task []Command

func (t Task) Run(tasks Tasks) {
	for _, c := range t {
		c.Run(tasks)
	}
}

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

type Tasks map[string]Task

func (t Tasks) Append(tasks Tasks, dir string) {
	for name, task := range tasks {
		for idx := range task {
			task[idx].Dir = dir
		}
		t[name] = task
	}
}

func LoadGlobalTasks(loaded_tasks *Tasks, dir string) {
	fpath := filepath.Join(os.Getenv("HOME"), ".config", "run", "tasks.yml")
	content, _ := os.ReadFile(fpath)

	tasks_map := map[string]Tasks{}
	yaml.Unmarshal(content, &tasks_map)

	if tasks, ok := tasks_map["global"]; ok {
		loaded_tasks.Append(tasks, dir)
	}

	for key, tasks := range tasks_map {
		if key == "global" {
			continue
		}

		// fmt.Printf("isSub(%s, abs(%s)) => %v\n", dir, key, isSub(dir, abs(key)))
		if isSub(dir, abs(key)) {
			loaded_tasks.Append(tasks, abs(key))
		}
	}
}

func LoadLocalTasks(loaded_tasks *Tasks, dir string) {
	// TODO find up the dir tree
	fpath := filepath.Join(dir, "tasks.yml")
	content, _ := os.ReadFile(fpath)

	tasks := Tasks{}
	yaml.Unmarshal(content, &tasks)

	loaded_tasks.Append(tasks, dir)
}

func GetTasks() Tasks {
	pwd, _ := os.Getwd()
	loaded_tasks := Tasks{}

	LoadGlobalTasks(&loaded_tasks, pwd)
	LoadLocalTasks(&loaded_tasks, pwd)

	return loaded_tasks
}

func (t Tasks) Run(name string) {
	if task, ok := t[name]; ok {
		task.Run(t)
	}
}

func isSub(dir, parent string) bool {
	if dir == parent {
		return true
	}
	return strings.HasPrefix(dir, parent+"/")
}

func abs(dir string) string {
	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, "~") {
		return strings.Replace(dir, "~", home, 1)
	}

	return dir
}
