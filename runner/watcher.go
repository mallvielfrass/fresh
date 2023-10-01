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
	absoluteWatchPath, _ := filepath.Abs(watchPath())
	rootPath, _ := filepath.Abs(root())
	watcherLog("watch_path: %s", absoluteWatchPath)
	watcherLog("root: %s", rootPath)

	err := filepath.Walk(absoluteWatchPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}
			//		fmt.Printf("Watching check: %s\n", path)
			if isIgnoredFolder(path) {
				//	watcherLog("Ignoring %s", path)
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
	if err != nil {

		mainLog("ERR: %s\n", err.Error())
	}
}
