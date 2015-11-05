package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watch() {
	watchPath := watchPath()
	filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 {
				if strings.HasPrefix(filepath.Base(path), ".") {
					return filepath.SkipDir
				}

				for _, folder := range excludeFolder() {
					if len(folder) > 1 && strings.Contains(path, folder) {
						watcherLog("Ignoring %s (cause:%s)", path, folder)
						return filepath.SkipDir
					}
				}
			}

			watchFolder(path)
		}

		return err
	})
}
