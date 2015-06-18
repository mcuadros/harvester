package input

import . "gopkg.in/check.v1"

type S3FileSuite struct{}

var _ = Suite(&S3FileSuite{})

func (s *S3FileSuite) TestGetRecord(c *C) {
	config := S3Config{
		AccessKey: "AKIAJ4DFHKQW6TIDDSRA",
		SecretKey: "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",
		Region:    "eu-west-1",
		Bucket:    "tyba-linkedin",
	}

	NewS3(&config, new(MockFormat))
}
