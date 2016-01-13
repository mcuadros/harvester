package output

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/mcuadros/harvester/src/intf"
)

type ElasticsearchConfig struct {
	Host     string `default:"localhost" description:"elastic host port "`
	Port     int    `default:"9200" description:"elastic search port "`
	Index    string `description:"index name"`
	Type     string `description:"index type"`
	Timeout  int    `default:"1" description:"contection timeout"`
	UidField string `description:"copied as id into into _uid, which consists of 'type#id'"`
}

type Elasticsearch struct {
	HTTP
	uidField string
}

func NewElasticsearch(config *ElasticsearchConfig) *Elasticsearch {
	defaults.SetDefaults(config)

	output := new(Elasticsearch)
	output.uidField = config.UidField
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
	return o.HTTP.PutRecord(record)
}

func (o *Elasticsearch) getIndexURL(config *ElasticsearchConfig) string {
	url := fmt.Sprintf("http://%s:%d/%s/%s",
		config.Host,
		config.Port,
		config.Index,
		config.Type,
	)

	if o.uidField != "" {
		url = fmt.Sprintf("%s/%%{%s}", url, o.uidField)
	}

	return url
}
