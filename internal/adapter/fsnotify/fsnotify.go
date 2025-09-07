package fsnotify

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/leonardoTavaresM/watcher/internal/service/watcher"
)

type FsnotifyAdapter struct {
	service *watcher.WatcherService
}

func NewFsnotifyAdapter(s *watcher.WatcherService) *FsnotifyAdapter {
	return &FsnotifyAdapter{service: s}
}

func (a *FsnotifyAdapter) Start(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	if err := AddDirsRecursively(watcher, path); err != nil {
		return err
	}

	log.Printf("init Watcher em: %s\n", path)

	for {
		select {
		case event := <-watcher.Events:
			if ShouldIgnore(event.Name) {
				continue
			}
			ext := filepath.Ext(event.Name)
			evType := ""
			if event.Op&fsnotify.Create == fsnotify.Create {
				evType = "CREATE"
			} else if event.Op&fsnotify.Write == fsnotify.Write {
				evType = "MODIFY"
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				evType = "Remove"
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				evType = "RENAME"
			} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				evType = "CHMOD"
			}

			if evType != "" {
				a.service.HandleFileEvent(event.Name, ext, evType)
			}
		case err := <-watcher.Errors:
			log.Println("Err:", err)
		}
	}
}

func AddDirsRecursively(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if ShouldIgnore(path) {
				log.Printf("Skipping directory: %s\n", path)
				return filepath.SkipDir
			}
			log.Printf("Watching directory: %s\n", path)
			return watcher.Add(path)
		}
		return nil
	})
}
