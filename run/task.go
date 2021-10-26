package run

var _ Step = &TaskStep{}

type TaskStep struct {
	// Task:
	// - task: <task-name>
	Task string
}

func (c *TaskStep) SetDir(dir string) {}

func (c TaskStep) Run(tasks Tasks) {
	if task, ok := tasks[c.Task]; ok {
		task.Run(tasks)
		return
	}
}
