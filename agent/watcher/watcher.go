package watcher

import (
	"github.com/hpcloud/tail"
	"log"
	"time"
)

const TickerDuration = 5 * time.Second

type FileWatcherListener interface {
	OnNewLine(line string)
}

type FileWatcherHandler func(string)

type FileWatcher struct {
	path   string
	tail   *tail.Tail
	ticker *time.Ticker

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

func (f *FileWatcher) Tick() {
	f.ticker.Reset(TickerDuration)
}

func (f *FileWatcher) Cleanup() {
	err := f.tail.Stop()
	if err != nil {
		log.Fatalf("Error stoping watcher tail: %v\n", err)
	}
	f.tail.Cleanup()
}
