package output

import (
	. "collector/logger"
	"labix.org/v2/mgo"
)

type MongoConfig struct {
	Url        string //Format: http://godoc.org/labix.org/v2/mgo#Dial
	Database   string
	Collection string
	Safe       bool
}

type Mongo struct {
	url            string
	dbName         string
	collectionName string
	collection     *mgo.Collection
	session        *mgo.Session
	safe           bool
	failed         int
	created        int
	transferred    int
}

func NewMongo(config *MongoConfig) *Mongo {
	output := new(Mongo)
	output.SetConfig(config)
	output.Connect()

	return output
}

func (self *Mongo) SetConfig(config *MongoConfig) {
	self.url = config.Url
	self.dbName = config.Database
	self.collectionName = config.Collection
	self.safe = config.Safe
}

func (self *Mongo) Connect() {
	Debug("Connecting to mongo server '%s' ...", self.url)
	session, err := mgo.Dial(self.url)
	if err != nil {
		Critical("Can't connect to mongo, go error %v\n", err)
	}

	self.session = session
	if self.safe {
		self.session.SetSafe(&mgo.Safe{})
	}

	self.collection = self.session.DB(self.dbName).C(self.collectionName)
}

func (self *Mongo) PutRecord(record map[string]string) bool {
	err := self.collection.Insert(record)
	if err != nil {
		Error("Can't insert record in mogo: %v\n", err)
		return false
	}

	return true
}
