package output

import (
	"encoding/json"
	"fmt"
	"harvesterd/intf"
	. "harvesterd/logger"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type ElasticsearchConfig struct {
	Host  string
	Port  int
	Index string
	Type  string
}

type Elasticsearch struct {
	config      *ElasticsearchConfig
	failed      int
	created     int
	transferred int
	url         string
	client      *http.Client
}

func NewElasticsearch(config *ElasticsearchConfig) *Elasticsearch {
	output := new(Elasticsearch)
	output.SetConfig(config)

	return output
}

func (self *Elasticsearch) SetConfig(config *ElasticsearchConfig) {
	self.config = config
	self.url = self.getIndexURL()

	self.createHTTPClient()
	Info("Created HTTP client")
}

func (self *Elasticsearch) createHTTPClient() {
	var dialer = &net.Dialer{Timeout: 1 * time.Second}
	var transport = &http.Transport{Dial: dialer.Dial}

	self.client = &http.Client{Transport: transport}
}

func (self *Elasticsearch) PutRecord(record intf.Record) bool {
	buffer := strings.NewReader(self.encodeToJSON(record))

	req, err := http.NewRequest("POST", self.url, buffer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := self.client.Do(req)
	if err != nil {
		Error("HTTP Error %s", err)
		return false
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return true
	}

	return false
}

func (self *Elasticsearch) encodeToJSON(record intf.Record) string {
	json, err := json.MarshalIndent(record, " ", "    ")
	if err != nil {
		Error("JSON Error %s", err)
	}

	self.transferred += len(json)
	return string(json)
}

func (self *Elasticsearch) getIndexURL() string {
	return fmt.Sprintf("http://%s:%d/%s/%s",
		self.config.Host,
		self.config.Port,
		self.config.Index,
		self.config.Type)
}
