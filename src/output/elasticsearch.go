package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

func NewElasticsearch(config *ElasticsearchConfig) *Elasticsearch {
	output := new(Elasticsearch)
	output.SetConfig(config)

	return output
}

func (self *Elasticsearch) SetConfig(config *ElasticsearchConfig) {
	self.config = config
	self.url = self.getIndexURL()
}

func (self *Elasticsearch) PutRecord(record map[string]string) bool {
	buffer := strings.NewReader(self.encodeToJSON(record))
	transport := &http.Transport{ResponseHeaderTimeout: time.Second * 45}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("POST", self.url, buffer)
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if len(body) > 0 && resp.StatusCode == http.StatusCreated {
		return true
	}

	return false
}

func (self *Elasticsearch) encodeToJSON(record map[string]string) string {
	json, err := json.MarshalIndent(record, " ", "    ")
	if err != nil {
		fmt.Println(err)
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
