package output

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
)

type ElasticsearchConfig struct {
	Host    string `default:"localhost" description:"elastic host port "`
	Port    int    `default:"9200" description:"elastic search port "`
	Index   string `description:"index name"`
	Type    string `description:"index type"`
	Timeout int    `default:"1" description:"contection timeout"`
}

type Elasticsearch struct {
	HTTP
}

func NewElasticsearch(config *ElasticsearchConfig) *Elasticsearch {
	defaults.SetDefaults(config)

	output := new(Elasticsearch)
	output.SetConfig(output.TransformConfig(config))

	return output
}

func (self *Elasticsearch) TransformConfig(config *ElasticsearchConfig) *HTTPConfig {
	dest := &HTTPConfig{
		Url:         self.getIndexURL(config),
		Method:      "POST",
		Format:      "json",
		ContentType: "application/json",
		Timeout:     config.Timeout,
	}

	return dest
}

func (self *Elasticsearch) getIndexURL(config *ElasticsearchConfig) string {
	return fmt.Sprintf("http://%s:%d/%s/%s",
		config.Host,
		config.Port,
		config.Index,
		config.Type,
	)
}
