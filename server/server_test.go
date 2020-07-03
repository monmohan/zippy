package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/monmohan/zippy"
)

func TestZipCreation(t *testing.T) {
	urls := []zippy.DownloadEntry{{"http://example.com", zippy.HTTP, "exmaple.com"},
		{"https://github.com", zippy.HTTP, "github.com"},
		{"https://www.reddit.com/r/starterpacks/comments/hidqze/affluent_suburbanite_rambo_starterpack/", zippy.HTTP, "redditimage"}}
	f, err := os.Create("/Users/singhmo/Downloads/test.zip")
	if err != nil {
		log.Fatalf("Error in creating zip writer")

	}
	if err := Zip(urls, f); err != nil {
		t.Fatalf("error in creating zip %s", err.Error())
	}

}

func TestZipCreationViaHttp(t *testing.T) {
	urls := []zippy.DownloadEntry{
		{"http://example.com", zippy.HTTP, "exmaple.com"},
		{"https://github.com", zippy.HTTP, "github.com"},
		{"https://www.reddit.com/r/starterpacks/comments/hidqze/affluent_suburbanite_rambo_starterpack/", zippy.HTTP, "redditimage"},
		{"com.github.monmohan.zippy/TestFolder/reddot.png", zippy.S3, "reddot.png"},
		{"com.github.monmohan.zippy/devnull.jpg", zippy.S3, "devnull.jpg"},
	}
	dlReq := DownloadAsZipRequest{Entries: urls, ZipName: "TestZipCreationMixed.zip"}
	b, _ := json.Marshal(dlReq)
	fmt.Printf("Request sending %v \n", string(b))
	f, err := os.Create("/Users/singhmo/Downloads/testh.zip")
	if err != nil {
		log.Fatalf("Error in creating zip writer")

	}
	resp, err := http.Post("http://127.0.0.1:9999/downloadzip", "application/json", bytes.NewReader(b))

	if err != nil {
		log.Fatalf("Error in download !! %v", err.Error())
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Error in download !! %v", resp.StatusCode)
	}
	io.Copy(f, resp.Body)

}