package run

import (
	"os"
	"os/exec"
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
}

func (t Command) Run() {
	// binary, lookErr := exec.LookPath(t.Command)
	// if lookErr != nil {
	// 	panic(lookErr)
	// }
	// err := syscall.Exec(binary, append([]string{t.Command}, t.Args...), os.Environ())
	// if err != nil {
	// 	log.Println(err)
	// }

	cmd := exec.Command(t.Command, t.Args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmd.Run()

	// log.Println("run1")
}

type Tasks map[string]Task

func (t Tasks) Append(tasks Tasks) {
	for name, task := range tasks {
		t[name] = task
	}
}

func LoadGlobalTasks() Tasks {
	fpath := filepath.Join(os.Getenv("HOME"), ".config", "run", "tasks.yml")
	content, _ := os.ReadFile(fpath)

	tasks_map := map[string]Tasks{}
	yaml.Unmarshal(content, &tasks_map)

	loaded_tasks := Tasks{}

	// keys := []string{"global", pwd}

	pwd, _ := os.Getwd()
	// home := os.Getenv("HOME")
	// if strings.HasPrefix(pwd, home) {
	// 	keys = append(keys, strings.Replace(pwd, home, "~", 1))
	// }

	if tasks, ok := tasks_map["global"]; ok {
		loaded_tasks.Append(tasks)
	}

	for key, tasks := range tasks_map {
		if key == "global" {
			continue
		}

		if isSub(pwd, abs(key)) {
			loaded_tasks.Append(tasks)
		}
	}

	return loaded_tasks
}

func GetTasks() Tasks {
	return LoadGlobalTasks()
	// return Tasks{
	// 	"tree": Task{
	// 		Command: "tree",
	// 	},
	// 	"vim": Task{
	// 		Command: "vim",
	// 	},
	// 	"nvim": Task{
	// 		Command: "nvim",
	// 	},
	// }
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
	return strings.HasPrefix(parent+"/", dir)
}

func abs(dir string) string {
	pwd, _ := os.Getwd()
	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, "~") {
		return strings.Replace(pwd, "~", home, 1)
	}

	return dir
}
