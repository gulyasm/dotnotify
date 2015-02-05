package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
)

func upload(fn string) {
	log.Printf("Uploading file %s")
}

func handle(w *fsnotify.Watcher) {
	for {
		select {
		case ev := <-w.Event:
			if ev.IsModify() && !ev.IsAttrib() {
				log.Println(fmt.Sprintf("Changed: %s", ev.Name))
				upload(ev.Name)
			}
		case err := <-w.Error:
			log.Println(err)
		}
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic("Error occured.")
	}

	done := make(chan bool)
	go handle(watcher)
	err = watcher.Watch(".")
	if err != nil {
		log.Panic("Panic error")
	}
	<-done
	watcher.Close()
}
