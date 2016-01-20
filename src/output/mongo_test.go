package output

import (
	"github.com/mcuadros/harvester/src/intf"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	. "gopkg.in/check.v1"
)

type msg struct {
	Id  bson.ObjectId `bson:"_id"`
	Foo string        `bson:"foo"`
	Qux string        `bson:"qux"`
}

type MongoSuite struct {
	url            string
	dbName         string
	collectionName string
	collection     *mgo.Collection
	session        *mgo.Session
}

var _ = Suite(&MongoSuite{
	url:            "localhost",
	dbName:         "test_harvester_output",
	collectionName: "bar",
})

func (s *MongoSuite) TestGetRecord(c *C) {
	config := MongoConfig{Url: s.url, Database: s.dbName, Collection: s.collectionName}

	output := NewMongo(&config)
	record := intf.Record{"foo": "bar", "qux": "baz"}

	c.Assert(output.PutRecord(record), Equals, true)

	s.connect()
	msg := s.getDocument(c)
	c.Assert(msg.Foo, Equals, "bar")
	c.Assert(msg.Qux, Equals, "baz")
	s.deleteDocument()
}

func (s *MongoSuite) connect() {
	s.session, _ = mgo.Dial(s.url)
	s.collection = s.session.DB(s.dbName).C(s.collectionName)
}

func (s *MongoSuite) getDocument(c *C) msg {
	var msg msg
	err := s.collection.Find(bson.M{}).One(&msg)
	if err != nil {
		c.Assert(false, Equals, true)
	}

	return msg
}

func (s *MongoSuite) deleteDocument() {
	session, _ := mgo.Dial("localhost")
	collection := session.DB("test_foo").C("bar")

	collection.RemoveAll(bson.M{})
}
