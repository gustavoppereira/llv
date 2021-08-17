package watcher

import (
	"github.com/hpcloud/tail"
	"log"
)

type FileWatcherListener interface {
	OnNewLine(line string)
}

type FileWatcherHandler func(string)

type FileWatcher struct {
	path string
	tail *tail.Tail

	handler FileWatcherHandler
}

func NewFileWatcher(path string, handler FileWatcherHandler) *FileWatcher {
	return &FileWatcher{
		path:    path,
		handler: handler,
	}
}

func (f *FileWatcher) Watch() error {
	t, err := tail.TailFile(f.path, tail.Config{Follow: true})
	if err != nil {
		log.Println(err)
		return err
	}
	f.tail = t

	go func() {
		for line := range t.Lines {
			f.handler(line.Text)
		}
	}()
	return nil
}

func (f *FileWatcher) Cleanup() {
	f.tail.Cleanup()
}
