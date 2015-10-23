package output

import (
	"strings"
	"time"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"

	"gopkg.in/mgo.v2"
)

type MongoConfig struct {
	Url         string `description:"server urls following the format from: http://godoc.org/labix.org/v2/mgo#Dial"`
	Database    string `description:"database name"`
	Collection  string `description:"collection name"`
	Safe        bool   `description:"sets the session safe mode: http://godoc.org/labix.org/v2/mgo#Session.SetSafe"`
	KillOnError bool   `description:"if true the server will die on any error, excep duplicate key error."`
}

type Mongo struct {
	url            string
	dbName         string
	collectionName string
	collection     *mgo.Collection
	session        *mgo.Session
	safe           bool
	killOnError    bool
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
	o.killOnError = config.KillOnError
}

func (o *Mongo) Connect() {
	Debug("Connecting to mongo server '%s' ...", o.url)
	session, err := mgo.DialWithTimeout(o.url, time.Second)
	if err != nil {
		Critical("Can't connect to mongo: %v\n", err)
	}

	o.session = session
	o.session.SetSocketTimeout(0)

	if o.safe {
		o.session.SetSafe(&mgo.Safe{})
	} else {
		o.session.SetSafe(nil)
	}

	o.collection = o.session.DB(o.dbName).C(o.collectionName)
}

func (o *Mongo) PutRecord(record intf.Record) bool {
	err := o.collection.Insert(record)
	if err != nil {
		if !o.killOnError || strings.Contains(err.Error(), "duplicate key") {
			Error("Can't insert record in mogo: %v\n", err)
		} else {
			Critical("Can't insert record in mogo: %v\n", err)
		}

		return false
	}

	return true
}
