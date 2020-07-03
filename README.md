# zippy
A server to download data from different sources as zip stream

# Run Server
`$ go run server.go `
# HTTP Request
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
```
Will stream back a zip file containing the 2 pages (example.com and github.com) and the two files in the S3 bucket (com.github.monmohan.zippy)