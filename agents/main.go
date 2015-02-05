package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"strings"
)

func upload(fn string) {
	log.Printf("Uploading file %s", fn)
}

func GetFileName(fn string) string {
	return fn[strings.LastIndex(fn, "/")+1:]
}

func handle(w *fsnotify.Watcher) {
	files := map[string]bool{
		".bashrc":       true,
		".bash_profile": true,
		".vimrc":        true,
		".tmux.conf":    true,
		".gitconfig":    true,
	}

	for {
		select {
		case ev := <-w.Event:
			if ev.IsModify() && !ev.IsAttrib() {
				name := GetFileName(ev.Name)
				log.Println(fmt.Sprintf("Changed: %s", name))
				if files[name] {
					upload(name)
				}
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
