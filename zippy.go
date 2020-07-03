package zippy

import (
	"archive/zip"
	"fmt"
	"io"
)

//DownloadURLType S3 or HTTP
type DownloadURLType int

const (
	//HTTP Urls that can downloaded with simple HTTP Get
	HTTP DownloadURLType = iota
	//S3 Object URLs
	S3
)

type DownloadEntry struct {
	Url     string
	UrlType DownloadURLType
	Name    string
}

type FetchedStream struct {
	Stream io.ReadCloser
	Err    error
	Name   string
}

type Fetcher func(entry DownloadEntry) FetchedStream

//FetchURL fetch a given URL using the fetcher func
func FetchURL(entry DownloadEntry, zipQ chan FetchedStream, fetcher Fetcher) {
	zipQ <- fetcher(entry)
}

//AddToZip create entries in zip
func AddToZip(zip *zip.Writer, zipQ chan FetchedStream) error {
	fetched := <-zipQ
	if fetched.Err != nil {
		return fetched.Err
	}
	defer fetched.Stream.Close()
	urlEntry, err := zip.Create(fetched.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Adding stream entry to zip %s\n", fetched.Name)
	io.Copy(urlEntry, fetched.Stream)
	return nil
}

func CloseZipWriter(zipWriter *zip.Writer) error {
	err := zipWriter.Close()
	if err != nil {
		return err
	}
	return nil
}
