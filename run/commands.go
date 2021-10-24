package run

import (
	"fmt"
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

func (t Task) Run() {
	for _, c := range t {
		c.Run()
	}
}

type Command struct {
	Shell   string   `yaml:"shell"`
	Script  string   `yaml:"script"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
	RunIn   string   `yaml:"run_in"`
	Dir     string   `yaml:"-"`
}

func (t Command) Run() {
	if t.Command != "" {
		// fmt.Printf("Run: %s [%s]\n", t.Command, t.Dir)
		cmd := exec.Command(t.Command, t.Args...)
		if t.RunIn != "" {
			cmd.Dir = path.Join(t.Dir, t.RunIn)
			// fmt.Printf("  => %s\n", cmd.Dir)
		}

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		cmd.Run()
	}

	if t.Shell != "" {
		file, err := ioutil.TempFile("", "run-script")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())

		filePath := file.Name()
		fmt.Println(filePath)

		ioutil.WriteFile(filePath, []byte(t.Script), 0777)

		cmd := exec.Command(t.Shell, filePath)
		if t.RunIn != "" {
			cmd.Dir = path.Join(t.Dir, t.RunIn)
			fmt.Printf("  => %s\n", cmd.Dir)
		}

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		cmd.Run()
	}

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
