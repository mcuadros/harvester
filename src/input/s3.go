package input

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"

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
	TrackFile string `description:"File for track the read files."`
	Gzip      bool
}

type S3 struct {
	*helper
	bucket    *s3.Bucket
	current   *s3.Key
	trackFile string
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
		i.factories = append(i.factories, i.createReaderFactory(key, c.Gzip))

	}

	if c.TrackFile == "" {
		return
	}

	i.trackFile = c.TrackFile

	i.helper.readerEOF = func() error {
		if i.current == nil {
			return nil
		}

		track := i.getTrackFilename()
		fmt.Println(track)
		if err := os.MkdirAll(filepath.Dir(track), 0644); err != nil {
			return err
		}

		return ioutil.WriteFile(track, []byte(i.current.LastModified), 0644)
	}
}

func (i *S3) getTrackFilename() string {
	return fmt.Sprintf(i.trackFile, i.current.Key)
}

func (i *S3) isCurrentKeyProcessed() bool {
	if i.current == nil || i.trackFile == "" {
		return false
	}

	if _, err := os.Stat(i.getTrackFilename()); err == nil {
		return true
	}

	return false
}

func (i *S3) createReaderFactory(key s3.Key, isGzip bool) ReaderFactory {
	return func() io.Reader {
		i.current = &key
		if i.isCurrentKeyProcessed() {
			Info("Skipping file %s", key.Key)
			return nil
		}

		var reader io.Reader
		var err error
		reader, err = i.bucket.GetReader(key.Key)
		if err != nil {
			Critical("open %s: %v", key.Key, err)
		}

		if isGzip {
			reader, err = gzip.NewReader(reader)
			if err != nil {
				Critical("gzip open %s: %v", key.Key, err)
			}
		}

		return reader
	}
}
