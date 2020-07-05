package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/monmohan/zippy"
)

func TestZipCreationDirect(t *testing.T) {
	urls := []zippy.DownloadEntry{{"http://example.com", zippy.HTTP, "exmaple.com"},
		{"https://github.com", zippy.HTTP, "github.com"},
		{"https://www.reddit.com/r/starterpacks/comments/hidqze/affluent_suburbanite_rambo_starterpack/", zippy.HTTP, "redditimage"}}
	f, err := os.Create("/tmp/TestZipCreationDirect.zip")
	if err != nil {
		log.Fatalf("Error in creating zip writer")

	}

	if err := Zip(context.Background(), urls, f); err != nil {
		t.Fatalf("error in creating zip %s", err.Error())
	}
	if err := verifyZip("/tmp/TestZipCreationDirect.zip", []string{"exmaple.com", "github.com", "redditimage"}); err != nil {
		t.Fatalf("zip verification failed !!, Error : %s", err.Error())
	}

}

func TestZipCreationMixed(t *testing.T) {
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
	f, err := os.Create("/tmp/TestZipCreationMixed.zip")
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
	if err := verifyZip("/tmp/TestZipCreationMixed.zip", []string{"exmaple.com", "devnull.jpg", "github.com", "reddot.png", "redditimage", "devnull.jpg"}); err != nil {
		t.Fatalf("zip verification failed !!, Error : %s", err.Error())
	}

}

func TestZipCreationContextTimeout(t *testing.T) {
	//create a fake test server which responds slowly
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		fmt.Fprintln(w, "Hello, client. Sorry I am slow")
	}))
	defer ts.Close()

	urls := []zippy.DownloadEntry{
		{"http://example.com", zippy.HTTP, "exmaple.com"},
		{"com.github.monmohan.zippy/devnull.jpg", zippy.S3, "devnull.jpg"},
		{ts.URL, zippy.HTTP, "slow.html"},
	}
	dlReq := DownloadAsZipRequest{Entries: urls, ZipName: "TestZipCreationMixed.zip"}
	b, _ := json.Marshal(dlReq)
	fmt.Printf("Request sending %v \n", string(b))
	f, err := os.Create("/tmp/testslow.zip")
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
	if err := verifyZip("/tmp/testslow.zip", []string{"exmaple.com", "devnull.jpg"}); err != nil {
		t.Fatalf("zip verification failed !!, Error : %s", err.Error())
	}

}
func verifyZip(fileName string, mustHaveFiles []string) error {

	zipReader, err := zip.OpenReader(fileName)
	mustHaves := make(map[string]bool, len(mustHaveFiles))
	for _, f := range mustHaveFiles {
		mustHaves[f] = true
	}

	if err != nil {
		fmt.Printf("Error reading zip file %s , err %s\n", fileName, err.Error())
		return err
	}
	defer zipReader.Close()
	for _, entry := range zipReader.File {
		if !mustHaves[entry.Name] {
			return fmt.Errorf("Zip contains the file %s which is not in must contain list %v", entry.Name, mustHaveFiles)
		}
		if entry.FileInfo().Size() == 0 {
			return fmt.Errorf("Zip contains the file %s which is reported as zero length contain list %v", entry.Name, mustHaveFiles)
		}
		delete(mustHaves, entry.Name)
	}
	var notFoundFiles string
	if len(mustHaves) > 0 {
		for k := range mustHaves {
			notFoundFiles += k
		}
		return fmt.Errorf("Some files were not foind in zip %s", notFoundFiles)
	}
	return nil
}
