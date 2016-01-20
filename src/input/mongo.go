package input

import (
	"time"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"

	"gopkg.in/mgo.v2"
)

type MongoConfig struct {
	Url        string `description:"server urls following the format from: http://godoc.org/labix.org/v2/mgo#Dial"`
	Database   string `description:"database name"`
	Collection string `description:"collection name"`
	BatchSize  int    `description:"sets the default batch size used when fetching documents from the database" default:1000`
}

type Mongo struct {
	collection *mgo.Collection
	session    *mgo.Session
	iter       *mgo.Iter
	eof        bool
}

func NewMongo(c *MongoConfig) *Mongo {
	i := new(Mongo)
	i.Connect(c)

	return i
}

func (i *Mongo) Connect(c *MongoConfig) {
	Debug("Connecting to mongo server '%s' ...", c.Url)
	session, err := mgo.DialWithTimeout(c.Url, time.Second)
	if err != nil {
		Critical("Can't connect to mongo: %v\n", err)
	}

	i.session = session
	i.session.SetSocketTimeout(0)
	i.session.SetBatch(c.BatchSize)
	i.session.SetPrefetch(0.5)

	i.collection = i.session.DB(c.Database).C(c.Collection)
	i.iter = i.collection.Find(nil).Iter()
}

func (i *Mongo) GetRecord() intf.Record {
	r := make(intf.Record)
	if ok := i.iter.Next(r); !ok {
		i.eof = true
		if err := i.iter.Err(); err != nil {
			Critical("error reading from input mongo: %v", err)
		}
	}

	return r
}

func (i *Mongo) IsEOF() bool {
	return i.eof
}

func (i *Mongo) Teardown() {
	i.iter.Close()
	i.session.Close()
}
