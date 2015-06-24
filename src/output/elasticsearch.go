package output

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/mcuadros/harvesterd/src/intf"
)

const ESIdField = "_id"

type ElasticsearchConfig struct {
	Host    string `default:"localhost" description:"elastic host port "`
	Port    int    `default:"9200" description:"elastic search port "`
	Index   string `description:"index name"`
	Type    string `description:"index type"`
	Timeout int    `default:"1" description:"contection timeout"`
	IdField string `description:"the content of the given field will be copied into _id"`
}

type Elasticsearch struct {
	HTTP
	idField string
}

func NewElasticsearch(config *ElasticsearchConfig) *Elasticsearch {
	defaults.SetDefaults(config)

	output := new(Elasticsearch)
	output.idField = config.IdField
	output.SetConfig(output.TransformConfig(config))

	return output
}

func (o *Elasticsearch) TransformConfig(config *ElasticsearchConfig) *HTTPConfig {
	dest := &HTTPConfig{
		Url:         o.getIndexURL(config),
		Method:      "POST",
		Format:      "json",
		ContentType: "application/json",
		Timeout:     config.Timeout,
	}

	return dest
}

func (o *Elasticsearch) PutRecord(record intf.Record) bool {
	if o.idField != "" {
		record[ESIdField] = record[o.idField]
	}

	return o.HTTP.PutRecord(record)
}

func (o *Elasticsearch) getIndexURL(config *ElasticsearchConfig) string {
	return fmt.Sprintf("http://%s:%d/%s/%s",
		config.Host,
		config.Port,
		config.Index,
		config.Type,
	)
}
