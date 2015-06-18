package input

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v3/testutil"
	. "gopkg.in/check.v1"
)

const MockRegion = "faux-region-1"

var testServer = testutil.NewHTTPServer()

type S3FileSuite struct{}

var _ = Suite(&S3FileSuite{})

func (s *S3FileSuite) SetUpSuite(c *C) {
	testServer.Start()
}

func (s *S3FileSuite) TestGetRecord(c *C) {
	aws.Regions[MockRegion] = aws.Region{
		Name:       MockRegion,
		S3Endpoint: testServer.URL,
	}

	config := S3Config{
		AccessKey: "abc",
		SecretKey: "123",
		Region:    MockRegion,
		Bucket:    "bucket",
		Gzip:      false,
	}

	testServer.Response(200, nil, GetListResultDump)
	s3 := NewS3(&config, new(MockFormat))

	testServer.Response(200, nil, "content")
	c.Assert(s3.GetLine(), Equals, "content")
}

func (s *S3FileSuite) TearDownSuite(c *C) {
	testServer.Stop()
}

func (s *S3FileSuite) TearDownTest(c *C) {
	testServer.Flush()
}

var GetListResultDump = `
<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01">
  <Name>quotes</Name>
  <Prefix>N</Prefix>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>Nelson</Key>
    <LastModified>2006-01-01T12:00:00.000Z</LastModified>
    <ETag>&quot;828ef3fdfa96f00ad9f27c383fc9ac7f&quot;</ETag>
    <Size>5</Size>
    <StorageClass>STANDARD</StorageClass>
    <Owner>
      <ID>bcaf161ca5fb16fd081034f</ID>
      <DisplayName>webfile</DisplayName>
     </Owner>
  </Contents>
  <Contents>
    <Key>Neo</Key>
    <LastModified>2006-01-01T12:00:00.000Z</LastModified>
    <ETag>&quot;828ef3fdfa96f00ad9f27c383fc9ac7f&quot;</ETag>
    <Size>4</Size>
    <StorageClass>STANDARD</StorageClass>
     <Owner>
      <ID>bcaf1ffd86a5fb16fd081034f</ID>
      <DisplayName>webfile</DisplayName>
    </Owner>
 </Contents>
</ListBucketResult>
`
