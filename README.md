# zippy
A server to download data from different sources as zip stream. Currently HTTP and S3 downloads are supported

# Run Server
`$ go run server.go `

# Golang client code sample
```
	//Create download entries, any number can be created and the entries can be mixed
	//example below shows one HTTP and one S3 URL
	urls := []zippy.DownloadEntry{
		{"http://example.com", zippy.HTTP, "exmaple.com"},
		{"com.github.monmohan.zippy/devnull.jpg", zippy.S3, "devnull.jpg"},
	}
	dlReq := DownloadAsZipRequest{Entries: urls, ZipName: "TestZipCreationMixed.zip"}
	b, _ := json.Marshal(dlReq)
	//create the zip file to write the stream to,
	//ignoring error handling since this is just a sample
	f, _ := os.Create("/tmp/TestZipCreationMixed.zip")
	resp, _ := http.Post("http://127.0.0.1:9999/downloadzip", "application/json", bytes.NewReader(b))
	io.Copy(f, resp.Body)
	// /tmp/TestZipCreationMixed.zip file should be available on disk now
  
  ```

# Direct HTTP Post request 
```
POST /downloadzip HTTP/1.1
Host: 127.0.0.1:9999
Content-Type: application/json

{
  "Entries":[
    {
      "Url":"http://example.com",
      "UrlType":0,
      "Name":"exmaple.com"
    },
    {
      "Url":"https://github.com",
      "UrlType":0,
      "Name":"githubhomepage"
    },
    {
      "Url":"com.github.monmohan.zippy/TestFolder/reddot.png",
      "UrlType":1,
      "Name":"reddot.png"
    },
    {
      "Url":"com.github.monmohan.zippy/devnull.jpg",
      "UrlType":1,
      "Name":"devnull.jpg"
    }
  ],
  "ZipName":"Somefilename.zip"
}
Will stream back a zip file containing the 2 pages (example.com and github.com) and the two files in the S3 bucket (com.github.monmohan.zippy)
```
