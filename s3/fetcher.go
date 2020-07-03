package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/monmohan/zippy"
)

func CreateFetcher() func(dlEntry zippy.DownloadEntry) zippy.FetchedStream {
	sess := session.Must(session.NewSession())
	ctx := context.Background()
	svc := s3.New(sess)
	return func(dlEntry zippy.DownloadEntry) zippy.FetchedStream {
		bucketKeySplit := strings.Index(dlEntry.Url, "/")
		if bucketKeySplit == -1 || bucketKeySplit == 0 {
			return zippy.FetchedStream{Stream: nil, Err: fmt.Errorf("Invalid object %v", dlEntry.Url), Name: dlEntry.Name}
		}
		result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(dlEntry.Url[:bucketKeySplit]),
			Key:    aws.String(dlEntry.Url[bucketKeySplit+1:]),
		})
		if err != nil {
			fmt.Printf("Error %v", err.Error())
			// Cast err to awserr.Error to handle specific error codes.
			aerr, ok := err.(awserr.Error)
			if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
				// Specific error code handling
				fmt.Printf("Error %v", aerr)
			}
			return zippy.FetchedStream{Stream: nil, Err: err, Name: dlEntry.Name}
		}
		fmt.Printf("Result %v\n", *result)
		return zippy.FetchedStream{Stream: result.Body, Err: nil, Name: dlEntry.Name}
	}
}
