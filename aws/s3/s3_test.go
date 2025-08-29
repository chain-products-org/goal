package s3_test

import (
	"bytes"
	"fmt"
	"github.com/chain-products-org/goal/aws/s3"
	"os"
	"testing"
)

const (
	accessKey    = ""
	accessSecret = ""
	region       = ""
	bucket       = ""
)

var client *s3.Client

func init() {
	client = s3.NewS3(&s3.Conf{AccessKey: accessKey, AccessSecret: accessSecret, Region: region})
}

func TestUpload(t *testing.T) {
	f := "/Users/home/Pictures/bing/' '.jpg"
	bs, _ := os.ReadFile(f)
	err := client.UploadFile(bucket, "a/00027c7c193111eebf32063de4e620c1.jpg", bytes.NewReader(bs))
	fmt.Println(err)
}

func TestDelete(t *testing.T) {
	err := client.DeleteFile(bucket, "a/00027c7c193111eebf32063de4e620ce.jpg")
	fmt.Println(err)
}
