package output

import (
	"harvesterd/intf"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

import . "launchpad.net/gocheck"

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
	url:            "mongodb://localhost",
	dbName:         "test_foo",
	collectionName: "bar",
})

func (self *MongoSuite) TestGetRecord(c *C) {
	config := MongoConfig{Url: self.url, Database: self.dbName, Collection: self.collectionName}

	output := NewMongo(&config)
	record := intf.Record{"foo": "bar", "qux": "baz"}

	c.Assert(output.PutRecord(record), Equals, true)

	self.connect()
	msg := self.getDocument(c)
	c.Assert(msg.Foo, Equals, "bar")
	c.Assert(msg.Qux, Equals, "baz")
	self.deleteDocument()
}

func (self *MongoSuite) connect() {
	self.session, _ = mgo.Dial(self.url)
	self.collection = self.session.DB(self.dbName).C(self.collectionName)
}

func (self *MongoSuite) getDocument(c *C) msg {
	var msg msg
	err := self.collection.Find(bson.M{}).One(&msg)
	if err != nil {
		c.Assert(false, Equals, true)
	}

	return msg
}

func (self *MongoSuite) deleteDocument() {
	session, _ := mgo.Dial("mongodb://localhost")
	collection := session.DB("test_foo").C("bar")

	collection.RemoveAll(bson.M{})
}
