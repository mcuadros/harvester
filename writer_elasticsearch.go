package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type WriterElasticSearchConfig struct {
	Host  string
	Port  int
	Index string
	Type  string
}

type WriterElasticSearch struct {
	failed      int
	created     int
	transferred int
	config      WriterElasticSearchConfig
	url         string
}

func NewWriterElasticSearch(config WriterElasticSearchConfig) *WriterElasticSearch {
	reader := new(WriterElasticSearch)
	reader.SetConfig(config)
	reader.ResetCounters()

	return reader
}

func (self *WriterElasticSearch) ResetCounters() {
	self.created = 0
	self.failed = 0
	self.transferred = 0
}

func (self *WriterElasticSearch) GetCounters() (int, int, int) {
	return self.created, self.failed, self.transferred
}

func (self *WriterElasticSearch) SetConfig(config WriterElasticSearchConfig) {
	self.config = config
	self.url = self.getIndexURL()
}

func (self *WriterElasticSearch) WriteFromChannel(channel chan map[string]string, wait sync.WaitGroup) {
	count := 0
	for record := range channel {
		if self.postRecordToIndex(record) {
			self.created++
		} else {
			self.failed++
		}

		count++
	}

	wait.Done()
}

func (self *WriterElasticSearch) postRecordToIndex(record map[string]string) bool {
	buffer := strings.NewReader(self.encodeToJSON(record))
	client := &http.Client{}
	req, err := http.NewRequest("POST", self.url, buffer)

	// NOTE this !!
	req.Close = true

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	if len(body) > 0 && resp.StatusCode == 201 {
		return true
	}

	return false
}

func (self *WriterElasticSearch) encodeToJSON(record map[string]string) string {
	json, err := json.MarshalIndent(record, " ", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}

	self.transferred += len(json)
	return string(json)
}

func (self *WriterElasticSearch) getIndexURL() string {
	return fmt.Sprintf("http://%s:%d/%s/%s",
		self.config.Host,
		self.config.Port,
		self.config.Index,
		self.config.Type)
}
