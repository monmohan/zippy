package s3

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/monmohan/zippy"
)

func TestS3DownloadFiles(t *testing.T) {
	fetchFn := CreateFetcher()

	getObject(fetchFn, "", "devnull.jpg")
	//try  a file within a folder
	getObject(fetchFn, "TestFolder", "reddot.png")

}

func getObject(fetchFn func(dlEntry zippy.DownloadEntry) zippy.FetchedStream, folder string, key string) {
	stream := fetchFn(zippy.DownloadEntry{"com.github.monmohan.zippy/" + folder + key, zippy.S3, key})
	f, _ := os.Create("/tmp/" + key)
	if stream.Err != nil {
		log.Fatalf("Error in Stream %v", stream.Err.Error())
	}
	defer stream.Stream.Close()
	io.Copy(f, stream.Stream)
}
