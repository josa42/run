package run

var _ Step = &TaskStep{}

type TaskStep struct {
	// Task:
	// - task: <task-name>
	Task string
}

func (c *TaskStep) SetDir(dir string) {}

func (c TaskStep) Run(tasks Tasks) (chan struct{}, CancelFunc) {
	if task, ok := tasks[c.Task]; ok {
		return task.Run(tasks)
	}

	return nil, func() {}
}
