package http

import (
	"fmt"
	"net/http"

	"github.com/monmohan/zippy"
)

//Fetch is used to downlload the http URL entry
func Fetch(dlEntry zippy.DownloadEntry) zippy.FetchedStream {
	fmt.Printf("fetching url %s\n", dlEntry.Url)
	resp, err := http.Get(dlEntry.Url)
	if err != nil {
		fmt.Printf("Error fetching URL..%s", err.Error())
		return zippy.FetchedStream{Stream: nil, Err: err, Name: ""}
	}
	return zippy.FetchedStream{Stream: resp.Body, Err: nil, Name: dlEntry.Name}

}
