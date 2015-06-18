package output

import (
	"time"

	"github.com/mcuadros/harvesterd/src/intf"
	. "github.com/mcuadros/harvesterd/src/logger"

	"gopkg.in/mgo.v2"
)

type MongoConfig struct {
	Url        string
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

func (o *Mongo) SetConfig(config *MongoConfig) {
	o.url = config.Url
	o.dbName = config.Database
	o.collectionName = config.Collection
	o.safe = config.Safe
}

func (o *Mongo) Connect() {
	Debug("Connecting to mongo server '%s' ...", o.url)
	session, err := mgo.DialWithTimeout(o.url, time.Second)
	if err != nil {
		Critical("Can't connect to mongo: %v\n", err)
	}

	o.session = session
	if o.safe {
		o.session.SetSafe(&mgo.Safe{})
	}

	o.collection = o.session.DB(o.dbName).C(o.collectionName)
}

func (o *Mongo) PutRecord(record intf.Record) bool {
	err := o.collection.Insert(record)
	if err != nil {
		Error("Can't insert record in mogo: %v\n", err)
		return false
	}

	return true
}
