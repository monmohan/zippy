package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/monmohan/zippy"
)

func allowMethods(method string, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(405)
			return
		}
		handler(w, r)
	}

}

func downloadZip(w http.ResponseWriter, r *http.Request) {
	var dlUrls zippy.DownloadAsZipRequest
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = json.Unmarshal(bytes, &dlUrls)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, dlUrls.ZipName))
	fmt.Printf("request %v\n", dlUrls)
	if err := zippy.Zip(dlUrls.Entries, w); err != nil {
		fmt.Printf("There were errors in creating zip %s  Error: %s\n", dlUrls.ZipName, err)
	}
}

func main() {
	http.HandleFunc("/downloadzip", allowMethods("POST", downloadZip))
	log.Fatal(http.ListenAndServe("127.0.0.1:9999", nil))

}
