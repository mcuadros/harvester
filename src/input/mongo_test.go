package input

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

type MongoSuite struct {
	url            string
	dbName         string
	collectionName string
	collection     *mgo.Collection
	session        *mgo.Session
}

var _ = Suite(&MongoSuite{
	url:            "localhost",
	dbName:         "test_harvester_input",
	collectionName: fmt.Sprintf("d%s", time.Now()),
})

func (s *MongoSuite) SetUpSuite(c *C) {
	var err error
	s.session, err = mgo.DialWithTimeout(s.url, time.Second)
	c.Assert(err, IsNil)

	s.collection = s.session.DB(s.dbName).C(s.collectionName)

	for i := 0; i < 10; i++ {
		err := s.collection.Insert(map[string]interface{}{
			"text":   fmt.Sprintf("foo %d", i),
			"number": i,
		})

		c.Assert(err, IsNil)

	}
}

func (s *MongoSuite) TestGetRecord(c *C) {
	config := MongoConfig{Url: s.url, Database: s.dbName, Collection: s.collectionName}
	input := NewMongo(&config)

	for i := 0; i < 11; i++ {
		r := input.GetRecord()

		if i >= 10 {
			c.Assert(len(r), Equals, 0)
		} else {
			c.Assert(r["number"], Equals, i)
		}
	}

	c.Assert(input.IsEOF(), Equals, true)
}
