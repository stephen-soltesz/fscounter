package main

import (
	"flag"
	"fmt"
	//"github.com/howeyc/fsnotify"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	//	"path/filepath"
	//"time"
)

// const FSN_MOST = fsnotify.FSN_DELETE | fsnotify.FSN_RENAME | fsnotify.FSN_CREATE

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

func startHandler(ch chan notify.EventInfo) {
	for {
		switch ev := <-ch; ev.Event() {
		case notify.Create:
			fmt.Printf("Create: %s\n", ev.Path())
			/*filepath.Walk(ev.Name, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					fmt.Printf("%s %#v\n", path, info)
					// log.Printf("w:Creating new watcher on %s", path)
					watcher.WatchFlags(path, FSN_MOST)
				}
				return nil
			})*/
		case notify.Remove:
			fmt.Printf("Remove: %s\n", ev.Path())
			//case err := <-watcher.Error:
			//	log.Println("fserror:", err)
		default:
			fmt.Printf("Event: %s\n", ev)
		}
	}
}

/*
// Make the channel buffered to ensure no event is dropped. Notify will drop
// an event if the receiver is not able to keep up the sending pace.
c := make(chan notify.EventInfo, 1)

// Set up a watchpoint listening for events within a directory tree rooted
// at current working directory. Dispatch remove events to c.
if err := notify.Watch("./...", c, notify.Remove); err != nil {
    log.Fatal(err)
}
defer notify.Stop(c)

// Block until an event is received.
ei := <-c
log.Println("Got event:", ei)o
*/

func createWatcher(dir string, ch chan notify.EventInfo) error {
	err := notify.Watch(dir, ch, notify.Create, notify.Remove, notify.Rename, notify.InCloseWrite)
	if err != nil {
		return err
	}

	startHandler(ch)

	return nil
}

func main() {
	flag.Parse()

	listener := make(chan notify.EventInfo, 10)
	// errors := make(chan error)

	fmt.Printf("watching: %s", *path)

	err := createWatcher(*path, listener)
	if err != nil {
		log.Fatal(err)
	}
}
