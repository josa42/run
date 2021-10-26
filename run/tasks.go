package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/josa42/run/pkg/utils"
	"gopkg.in/yaml.v2"
)

type Step interface {
	Run(tasks Tasks)
	SetDir(dir string)
}

type Task struct {
	Steps []Step
}

func (t Task) Run(tasks Tasks) {
	for _, c := range t.Steps {
		c.Run(tasks)
	}
}

type stepRaw struct{}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := []interface{}{}
	unmarshal(&data)

	commands := []CommandStep{}
	unmarshal(&commands)

	shellSteps := []ShellStep{}
	unmarshal(&shellSteps)

	taskSteps := []TaskStep{}
	unmarshal(&taskSteps)

	t.Steps = []Step{}
	for idx := range data {
		if commands[idx].Command != "" {
			t.Steps = append(t.Steps, &commands[idx])

		} else if shellSteps[idx].Shell != "" {
			t.Steps = append(t.Steps, &shellSteps[idx])

		} else if taskSteps[idx].Task != "" {
			t.Steps = append(t.Steps, &taskSteps[idx])

		} else {
			func(v interface{}) {
				s, _ := yaml.Marshal(v)
				fmt.Printf("%s\n", s)
			}(data[idx])
			panic("Unknown step")
		}
	}

	return nil
}

type Tasks map[string]Task

func (t Tasks) Append(tasks Tasks, dir string) {
	for name, task := range tasks {
		for idx := range task.Steps {
			task.Steps[idx].SetDir(dir)
		}
		t[name] = task
	}
}

func GetTasks() Tasks {
	pwd, _ := os.Getwd()
	loaded_tasks := Tasks{}

	dir, _ := utils.FindUp(pwd, "tasks.yml")
	if dir != "" {
		loadProjectTasks(&loaded_tasks, dir)
		loadGlobalTasks(&loaded_tasks, dir)

	} else {
		loadGlobalTasks(&loaded_tasks, pwd)
	}

	return loaded_tasks
}

func loadGlobalTasks(loaded_tasks *Tasks, dir string) {
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

		if utils.IsSubDir(dir, utils.Abs(key)) {
			loaded_tasks.Append(tasks, utils.Abs(key))
		}
	}
}

func loadProjectTasks(loaded_tasks *Tasks, dir string) {
	fpath := filepath.Join(dir, "tasks.yml")
	content, _ := os.ReadFile(fpath)

	tasks := Tasks{}
	yaml.Unmarshal(content, &tasks)

	loaded_tasks.Append(tasks, dir)
}

func (t Tasks) Run(name string) {
	if task, ok := t[name]; ok {
		task.Run(t)
	}
}
