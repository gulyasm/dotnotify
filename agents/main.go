package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func upload(fn string) error {
	log.Printf("Uploading file %s", fn)
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	uri := fmt.Sprintf("http://localhost:7644/upload?user=%s", "gulyasm")
	if err != nil {
		return err
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"myfiles", filepath.Base(fn)))
	h.Set("Content-Type", http.DetectContentType(body.Bytes()))
	formfile, err := writer.CreatePart(h)
	_, err = io.Copy(formfile, file)
	if err != nil {
		return err
	}
	params := map[string]string{}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", uri, body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(r)
	return err
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
					err := upload(ev.Name)
					if err != nil {
						log.Println("Failed to upload", err)
					}
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
	err = watcher.Watch("/home/gulyasm")
	if err != nil {
		log.Panic("Panic error")
	}
	<-done
	watcher.Close()
}
