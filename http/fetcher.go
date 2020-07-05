package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/monmohan/zippy"
)

//CreateFetcher creates an HTTP Fetcher, to downlload the http URL entry
func CreateFetcher(ctx context.Context) zippy.Fetcher {
	return func(dlEntry zippy.DownloadEntry) zippy.FetchedStream {
		fmt.Printf("fetching url %s\n", dlEntry.Url)
		req, err := http.NewRequest("GET", dlEntry.Url, nil)
		req = req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error fetching URL..%s", err.Error())
			return zippy.FetchedStream{Stream: nil, Err: err, Name: dlEntry.Name}
		}
		return zippy.FetchedStream{Stream: resp.Body, Err: nil, Name: dlEntry.Name}
	}
}
