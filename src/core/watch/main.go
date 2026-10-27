package watch

import (
	"errors"
	"merodi/src/core/state"
	"merodi/src/parser"
	"merodi/src/utils"
	"os"
	"time"

	. "github.com/pixel2175/errs"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	State *state.State
}

func (self *Watcher) Run(state *state.State, opt *parser.Option) (err error) {
	defer Catch(&err)
	self.State = state

	watcher := TryV(fsnotify.NewWatcher())
	defer watcher.Close()

	if len(state.Watch.Add) == 0 {
		TryE(watcher.Add(state.Config.Tree.Markdown))
	} else {
		for _, file := range self.State.Watch.Add {
			_, statErr := os.Stat(file)
			if os.IsNotExist(statErr) {
				TryE(statErr)
			}
			TryE(watcher.Add(file))
		}
	}

	TryE(self.State.RunHook("on_start_watching"))

	lastEvents := make(map[string]time.Time)
	const debounce = 100 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			now := time.Now()
			if last, ok := lastEvents[event.Name]; ok && now.Sub(last) < debounce {
				continue
			}

			lastEvents[event.Name] = now
			TryE(state.RunHook("on_file_changed", event.Op.String(), event.Name))

		case werr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			if !errors.Is(werr, utils.ErrAborted) {
				TryE(werr)
			}
		}
	}
}
