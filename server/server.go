package main

import (
	"expvar"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	statRequests = expvar.NewMap("requests")
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	statRequests.Add("upload", 1)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Request is not POST")
		return
	}
	up := r.URL.Query()["user"]
	if len(up) == 0 {
		log.Println("User parameter not found")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user parameter not found"))
		return
	}
	user := up[0]
	if user == "" {
		log.Println("User parameter not found")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user parameter not found"))
		return
	}
	path := fmt.Sprintf("./files/%s/", user)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0771)
		if err != nil {
			log.Panic(fmt.Sprintf("Failed to create user directory: %s", path), err)
		}
	}
	//parse the multipart form in the request
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["myfiles"]
	for i, f := range files {
		//for each fileheader, get a handle to the actual file
		file, err := f.Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create(path + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func admin() error {
	sock, err := net.Listen("tcp", "localhost:8123")
	if err != nil {
		return err
	}
	go func() {
		log.Println("HTTP now available at port 8123")
		http.Serve(sock, nil)
	}()
	return nil
}

func main() {
	err := admin()
	if err != nil {
		log.Panic("Failed to start admin module", err)
	}
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err = os.Mkdir("./files", 0771)
		if err != nil {
			log.Panic("Failed to create files directory", err)
		}
	}
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":7644", nil)
}
