package run

import (
	"fmt"

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

func (c WatchStep) Run(tasks Tasks) {
	m := watcher.NewMatcher(c.Watch)

	notify := func(ev watcher.Event) {
		if ev.Is(fsnotify.Create, fsnotify.Write) && m.Match(ev.Name) {
			fmt.Printf("=> %s: %s\n", ev.Op, ev.Name)
			c.Do.Run(tasks)
		}
	}

	watcher := watcher.New(
		watcher.Notify(notify),
		watcher.Exclude("*/.git/*", "**/*~", "**/node_modules/**"),
		watcher.Exclude(c.Exclude...),
	)
	defer watcher.Stop()
	watcher.Add(".", true)

	utils.WaitForKill()
}
