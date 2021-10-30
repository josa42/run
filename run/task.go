package run

import (
	"fmt"
	"sync"

	"gopkg.in/yaml.v2"
)

type Task struct {
	Steps []Step
}

func (t Task) Run(tasks Tasks) (chan struct{}, CancelFunc) {
	fns := []CancelFunc{}
	canceled := false

	done := make(chan struct{})
	go func() {
		for _, c := range t.Steps {
			runDone, cancel := c.Run(tasks)

			fns = append(fns, cancel)
			<-runDone

			if canceled {
				break
			}
		}
		close(done)
	}()

	return done, func() {
		canceled = true
		for _, cancel := range fns {
			cancel()
		}
	}
}

func (t Task) RunParallel(tasks Tasks) (chan struct{}, CancelFunc) {
	fns := []CancelFunc{}

	var wait sync.WaitGroup
	wait.Add(len(t.Steps))

	for _, c := range t.Steps {
		done, cancel := c.Run(tasks)

		fns = append(fns, cancel)
		go func() {
			<-done
			wait.Done()
		}()
	}

	done := make(chan struct{})
	go func() {
		wait.Wait()
		close(done)
	}()

	return done, func() {
		for _, cancel := range fns {
			cancel()
		}
	}
}

type stepRaw struct{}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := []interface{}{}
	unmarshal(&data)

	runSteps := []RunStep{}
	unmarshal(&runSteps)

	taskSteps := []TaskStep{}
	unmarshal(&taskSteps)

	watchSteps := []WatchStep{}
	unmarshal(&watchSteps)

	parallelSteps := []ParallelStep{}
	unmarshal(&parallelSteps)

	t.Steps = []Step{}
	for idx := range data {
		if runSteps[idx].Script != "" {
			t.Steps = append(t.Steps, &runSteps[idx])

		} else if taskSteps[idx].Task != "" {
			t.Steps = append(t.Steps, &taskSteps[idx])

		} else if watchSteps[idx].Watch != nil {
			t.Steps = append(t.Steps, &watchSteps[idx])

		} else if parallelSteps[idx].Parallel.Steps != nil {
			t.Steps = append(t.Steps, &parallelSteps[idx])

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
