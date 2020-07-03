package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/monmohan/zippy"
	zhttp "github.com/monmohan/zippy/http"
	"github.com/monmohan/zippy/s3"
)

//DownloadAsZipRequest is request body for downloading a set of URLs
type DownloadAsZipRequest struct {
	Entries []zippy.DownloadEntry
	ZipName string
}

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
	var dlUrls DownloadAsZipRequest
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
	if err := Zip(dlUrls.Entries, w); err != nil {
		fmt.Printf("There were errors in creating zip %s  Error: %s\n", dlUrls.ZipName, err)
	}
}

//Zip zips a set of download URLs and writes to given writer
func Zip(urls []zippy.DownloadEntry, writeTo io.Writer) error {

	zipQ := make(chan zippy.FetchedStream)
	var fetchFn zippy.Fetcher
	for _, url := range urls {
		fetchFn = zhttp.Fetch
		if url.UrlType == zippy.S3 {
			fetchFn = s3.CreateFetcher()
		}
		go zippy.FetchURL(url, zipQ, fetchFn)
	}
	zipWriter := zip.NewWriter(writeTo)
	for range urls {
		zippy.AddToZip(zipWriter, zipQ)
	}
	return zippy.CloseZipWriter(zipWriter)

}

func main() {
	http.HandleFunc("/downloadzip", allowMethods("POST", downloadZip))
	log.Fatal(http.ListenAndServe("127.0.0.1:9999", nil))

}
