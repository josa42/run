package run

var _ Step = &ParallelStep{}

type ParallelStep struct {
	Parallel Task `yaml:"parallel"`
}

func (c *ParallelStep) SetDir(dir string) {
	for _, s := range c.Parallel.Steps {
		s.SetDir(dir)
	}
}

func (c ParallelStep) Run(tasks Tasks) (chan struct{}, CancelFunc) {
	return c.Parallel.RunParallel(tasks)
}
