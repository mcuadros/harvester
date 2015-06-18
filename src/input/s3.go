package input

import (
	"bufio"
	"fmt"

	"github.com/mcuadros/harvesterd/src/intf"
	. "github.com/mcuadros/harvesterd/src/logger"

	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

type S3Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	Format    string `description:"A valid format name."`
	Prefix    string `description:"Limits the response to keys that begin with the specified prefix."`
	Delimiter string `description:"A delimiter is a character you use to group keys."`
	Marker    string `description:"Specifies the key to start with when listing objects in a bucket."`
	MaxKeys   int    `description:"Sets the maximum number of keys returned."`
}

type S3 struct {
	*helper
	bucket *s3.Bucket
}

func NewS3(config *S3Config, format intf.Format) *S3 {
	input := &S3{helper: newHelper(format)}
	input.SetConfig(config)

	return input
}

func (i *S3) SetConfig(c *S3Config) {
	client := s3.New(aws.Auth{c.AccessKey, c.SecretKey}, aws.Regions[c.Region])
	i.bucket = client.Bucket(c.Bucket)

	r, err := i.bucket.List(c.Prefix, c.Delimiter, c.Marker, c.MaxKeys)
	if err != nil {
		Critical("list bucket %s: %v", c.Bucket, err)
	}

	for _, key := range r.Contents {
		i.createBufioReader(key)
	}

	if len(i.files) == 0 {
		i.empty = true
		i.eof = true
	}
}

func (i *S3) createBufioReader(key s3.Key) *bufio.Scanner {
	reader, err := i.bucket.GetReader(key.Key)
	fmt.Println(key)
	if err != nil {
		Critical("open %s: %v", key.Key, err)
	}

	return bufio.NewScanner(reader)
}
