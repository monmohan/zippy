package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
)

type DownloadURLType int

const (
	HTTP DownloadURLType = iota
	S3
)

type DownloadEntry struct {
	Url     string
	urlType DownloadURLType
	Name    string
}
type DownloadAsZipRequest struct {
	Entries []DownloadEntry
	ZipName string
}

type fetchedstream struct {
	stream io.ReadCloser
	err    error
	name   string
}

func Zip(urls []DownloadEntry, writeTo io.Writer) error {

	zipQ := make(chan fetchedstream)

	for _, url := range urls {
		go fetchURL(url, zipQ)
	}
	zipWriter := zip.NewWriter(writeTo)
	for range urls {
		addToZip(zipWriter, zipQ)
	}
	return closeZipWriter(zipWriter)

}

func fetchURL(dlUrl DownloadEntry, zipQ chan fetchedstream) {
	fmt.Printf("fetching url %s\n", dlUrl.Url)
	resp, err := http.Get(dlUrl.Url)
	var fetched fetchedstream
	if err != nil {
		fmt.Printf("Error fetching URL..%s", err.Error())
		fetched.stream, fetched.err, fetched.name = nil, err, ""
		zipQ <- fetched
	}
	fetched.stream, fetched.err, fetched.name = resp.Body, nil, dlUrl.Name
	zipQ <- fetched

}

func addToZip(zip *zip.Writer, zipQ chan fetchedstream) error {
	fetched := <-zipQ
	if fetched.err != nil {
		return fetched.err
	}
	defer fetched.stream.Close()
	urlEntry, err := zip.Create(fetched.name)
	if err != nil {
		return err
	}
	fmt.Printf("Adding stream entry to zip %s\n", fetched.name)
	io.Copy(urlEntry, fetched.stream)
	return nil
}

func closeZipWriter(zipWriter *zip.Writer) error {
	err := zipWriter.Close()
	if err != nil {
		return err
	}
	return nil
}
