package run

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/josa42/run/pkg/utils"
	"gopkg.in/yaml.v2"
)

type Tasks map[string]Task

func (t Tasks) Append(tasks Tasks, dir string) {
	for name, task := range tasks {
		for idx := range task.Steps {
			task.Steps[idx].SetDir(dir)
		}
		t[name] = task
	}
}

func (t Tasks) Run(name string) {
	if task, ok := t[name]; ok {
		done, _ := task.Run(t)
		<-done
		return
	}

	// TODO log error
	os.Exit(1)
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

		for _, ikey := range strings.Split(key, "|") {
			if ok, cdir := utils.IsSubDir(dir, utils.Abs(ikey)); ok {
				loaded_tasks.Append(tasks, cdir)
				break
			}
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
