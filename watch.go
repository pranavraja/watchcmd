package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/howeyc/fsnotify"
)

type RecursiveWatcher struct {
	*fsnotify.Watcher
}

func NewRecursiveWatcher(path string) (watcher RecursiveWatcher, err error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	watcher = RecursiveWatcher{w}
	err = watcher.WatchRecursive(path)
	if err != nil {
		return
	}
	return watcher, nil
}

func (watcher RecursiveWatcher) WatchRecursive(path string) (err error) {
	err = watcher.Watch(path)
	filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		if info.IsDir() {
			err = watcher.Watch(path)
		}
		return nil
	})
	return
}

func (watcher RecursiveWatcher) Handle(ev *fsnotify.FileEvent) {
	if ev.IsCreate() {
		info, err := os.Stat(ev.Name)
		if err == nil {
			if info.IsDir() {
				watcher.WatchRecursive(ev.Name)
			}
		}
	} else if ev.IsDelete() {
		watcher.RemoveWatch(ev.Name)
	}
}

func (watcher RecursiveWatcher) NextBatch(d time.Duration) (events []*fsnotify.FileEvent, err error) {
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		select {
		case ev := <-watcher.Event:
			watcher.Handle(ev)
			events = append(events, ev)
		case err = <-watcher.Error:
			return
		}
		// Collect the rest of the available events that occured over duration d
		timeout := time.After(d)
		for {
			select {
			case ev := <-watcher.Event:
				watcher.Handle(ev)
				events = append(events, ev)
			case err = <-watcher.Error:
				return
			case <-timeout:
				return
			}
		}
	}()
	<-done
	return
}
