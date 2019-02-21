package plugin

import (
	"fmt"
	"go-api-playground/logging"
	"time"

	"github.com/fsnotify/fsnotify"
)

var logger logging.Logger

func CreatePluginWatcher(dir string, plugUpdateChan chan string, logChan chan string) {
	logger = logging.NewLogger("", logChan)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Log(fmt.Sprintf("error: failed to create a new FS watcher. %s", err.Error()))
		return
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Log("error: an unknown error occurred while attempting to receive an event from the FS watcher.")
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					time.Sleep(1000 * time.Millisecond)
					logger.Log(fmt.Sprintf("info: event received from plugin FS watcher. [Modified File: %s]", event.Name))
					plugUpdateChan <- event.Name
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					logger.Log("error: an unknown error occurred while attempting to receive an event from the FS watcher.")
				}

				logger.Log(fmt.Sprintf("error: an unknown error occurred while attempting to receive an event from the FS watcher. %s", err.Error()))
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		logger.Log(fmt.Sprintf("error: an error occurred while adding the directory to the FS watcher: [Directory: %s]. %s", dir, err.Error()))
	}

	<-done
}
