package collector

import (
	"encoding/json"
	"fmt"
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
	failed  int32
	created int32
	config  WriterElasticSearchConfig
	url     string
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
}
func (self *WriterElasticSearch) PrintAndResetCounters() {
	fmt.Println(fmt.Sprintf("Created %d document(s), Failed %d times(s)",
		self.created,
		self.failed))

	self.ResetCounters()
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
		if count%500 == 0 {
			self.PrintAndResetCounters()
		}
	}

	wait.Done()
}

func (self *WriterElasticSearch) postRecordToIndex(record map[string]string) bool {
	buffer := strings.NewReader(self.encodeToJSON(record))
	resp, err := http.Post(self.url, "application/json", buffer)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	if resp.StatusCode == 201 {
		return true
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	if body {
		return true
	}

	return false
}

func (self *WriterElasticSearch) encodeToJSON(record map[string]string) string {
	json, err := json.MarshalIndent(record, " ", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(json)
}

func (self *WriterElasticSearch) getIndexURL() string {
	return fmt.Sprintf("http://%s:%d/%s/%s",
		self.config.Host,
		self.config.Port,
		self.config.Index,
		self.config.Type)
}
