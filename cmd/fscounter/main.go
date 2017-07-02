package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

const FSN_MOST = fsnotify.FSN_DELETE | fsnotify.FSN_RENAME | fsnotify.FSN_CREATE

var (
	path = flag.String("path", "", "Path of root directory to watch.")
)

func isDir(path string) (bool, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	if stat.IsDir() {
		return true, nil
	}
	return false, nil
}

func startHandler(watcher *fsnotify.Watcher, ch chan *fsnotify.FileEvent) {
	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsCreate() {
				isdir, err := isDir(ev.Name)
				if err != nil {
					continue
				}
				if isdir {
					// log.Printf("Creating new watcher on %s", ev.Name)
					fmt.Println(watcher.WatchFlags(ev.Name, FSN_MOST))
					fmt.Println(watcher.WatchFlags(ev.Name, FSN_MOST))
					filepath.Walk(ev.Name, func(path string, info os.FileInfo, err error) error {
						if info.IsDir() {
							// log.Printf("w:Creating new watcher on %s", path)
							watcher.WatchFlags(path, FSN_MOST)
						}
						return nil
					})
				}
			} else if ev.IsDelete() {
				// log.Println("Close watcher if delete a dir.")
				watcher.RemoveWatch(ev.Name)
			}
			ch <- ev
		case err := <-watcher.Error:
			log.Println("fserror:", err)
		}
	}
}

func createWatcher(dir string, ch chan *fsnotify.FileEvent) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go startHandler(watcher, ch)

	err = watcher.WatchFlags(dir, FSN_MOST)
	if err != nil {
		return err
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			//log.Printf("w:Creating new watcher on %s", path)
			watcher.WatchFlags(path, FSN_MOST)
		}
		return nil
	})
	return nil
}

func main() {
	flag.Parse()

	listener := make(chan *fsnotify.FileEvent)
	// errors := make(chan error)

	fmt.Printf("watching: %s", *path)

	err := createWatcher(*path, listener)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case ev := <-listener:
			log.Println("event:", ev)
			//case err := <-errors:
			//	log.Println("error:", err)
		}
	}
}
