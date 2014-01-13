package output

import (
	"encoding/json"
	"errors"
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

var (
	httpNonCreatedCode = errors.New("http: received != 201 status code")
	httpNetworkError   = errors.New("http: network error")
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
}

func (self *Elasticsearch) createHTTPClient() {
	var dialer = &net.Dialer{Timeout: 1 * time.Second}
	var transport = &http.Transport{Dial: dialer.Dial}

	self.client = &http.Client{Transport: transport}
}

func (self *Elasticsearch) PutRecord(record intf.Record) bool {
	buffer := strings.NewReader(self.encodeToJSON(record))

	retryCount := 0
	retry := true
	for retry {
		retryCount++

		err, ctx := self.makeRequest(buffer)
		switch err {
		case httpNetworkError:
			Debug("%s, retrying", ctx)
			retry = true
			break
		case httpNonCreatedCode:
			Error("%s: received %d", httpNonCreatedCode, ctx)
			return false
		case nil:
			return true
		}

		if retryCount >= 10 {
			Error("retry limit reached, network issues")
			return false
		}
	}

	return false
}

func (self *Elasticsearch) makeRequest(buffer *strings.Reader) (error, interface{}) {
	req, err := http.NewRequest("POST", self.url, buffer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := self.client.Do(req)
	if err != nil {
		return httpNetworkError, err
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return httpNonCreatedCode, nil
	}

	return nil, nil
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
