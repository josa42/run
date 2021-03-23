package run

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Task []Command

func (t Task) Run() {
	for _, c := range t {
		c.Run()
	}
}

type Command struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
	RunIn   string   `yaml:"run_in"`
	Dir     string   `yaml:"-"`
}

func (t Command) Run() {
	// fmt.Printf("Run: %s [%s]\n", t.Command, t.Dir)
	cmd := exec.Command(t.Command, t.Args...)
	if t.RunIn != "" {
		cmd.Dir = path.Join(t.Dir, t.RunIn)
		fmt.Printf("  => %s\n", cmd.Dir)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmd.Run()

	// log.Println("run1")
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

func LoadGlobalTasks(dir string) Tasks {
	fpath := filepath.Join(os.Getenv("HOME"), ".config", "run", "tasks.yml")
	content, _ := os.ReadFile(fpath)

	tasks_map := map[string]Tasks{}
	yaml.Unmarshal(content, &tasks_map)

	loaded_tasks := Tasks{}

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

	return loaded_tasks
}

func GetTasks() Tasks {
	pwd, _ := os.Getwd()
	return LoadGlobalTasks(pwd)
}

func (t Tasks) Run(name string) {
	if task, ok := t[name]; ok {
		task.Run()
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
