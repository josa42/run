package run

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/josa42/run/pkg/utils"
	"github.com/josa42/run/pkg/watcher"
)

var _ Step = &WatchStep{}

type WatchStep struct {
	Watch   []string `yaml:"watch"`
	Exclude []string `yaml:"exclude"`
	Do      Task     `yaml:"do"`
}

func (c *WatchStep) SetDir(dir string) {}

func (c WatchStep) Run(tasks Tasks) (chan struct{}, CancelFunc) {
	m := watcher.NewMatcher(c.Watch)
	var cancelFn CancelFunc

	notify := func(evt watcher.Event) {
		if evt.Is(fsnotify.Create, fsnotify.Write) && m.Match(evt.Name) {
			fmt.Print(color.BlueString("%s %s\n", evt.Op, evt.Name))

			if cancelFn != nil {
				cancelFn()
			}
			_, cancel := c.Do.Run(tasks)

			cancelFn = cancel
		}
	}

	watcher := watcher.New(
		watcher.Notify(notify),
		watcher.Exclude("*/.git/*", "**/*~", "**/node_modules/**"),
		watcher.Exclude(c.Exclude...),
	)
	watcher.Add(".", true)

	done := make(chan struct{})
	go func() {
		utils.WaitForKill()
		watcher.Stop()
		close(done)
	}()

	return done, func() {}
}
