package output

import (
	"encoding/json"
	"errors"
	"harvesterd/intf"
	. "harvesterd/logger"
	"harvesterd/util"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

import "github.com/ajg/form"
import "github.com/mcuadros/go-defaults"

var (
	httpNonOkCode    = errors.New("http: received non 2xx status code")
	httpNetworkError = errors.New("http: network error")
)

type HTTPConfig struct {
	Url         string   `description:"url of the request"`
	Format      string   `default:"form" description:"format of the request body (json or form)"`
	ContentType string   `default:"application/x-www-form-urlencoded" description:"Content-Type header"`
	Method      string   `default:"POST" description:"request method"`
	Timeout     int      `default:"1" description:"contection timeout"`
	Header      []string `description:"additional headers, format: header,value"`
}

type HTTP struct {
	url         *util.Template
	headers     map[string]*util.Template
	format      string
	contentType string
	method      string
	timeout     time.Duration
	failed      int
	created     int
	transferred int
	client      *http.Client
}

func NewHTTP(config *HTTPConfig) *HTTP {
	output := new(HTTP)
	output.SetConfig(config)

	return output
}

func (self *HTTP) SetConfig(config *HTTPConfig) {
	defaults.SetDefaults(config)

	self.url = util.NewTemplate(config.Url)
	self.format = config.Format
	self.contentType = config.ContentType
	self.method = config.Method
	self.timeout = time.Duration(config.Timeout)
	self.parseHeadersConfig(config.Header)

	self.createHTTPClient()
}

func (self *HTTP) parseHeadersConfig(headers []string) {
	self.headers = make(map[string]*util.Template, len(headers))

	for _, headerRaw := range headers {
		header := strings.Split(headerRaw, ",")
		if len(header) != 2 {
			Critical("Malformed header setting '%s'", header)
		}

		self.headers[header[0]] = util.NewTemplate(header[1])
	}
}

func (self *HTTP) createHTTPClient() {
	var dialer = &net.Dialer{Timeout: self.timeout * time.Second}
	var transport = &http.Transport{Dial: dialer.Dial}

	self.client = &http.Client{Transport: transport}
}

func (self *HTTP) PutRecord(record intf.Record) bool {
	retryCount := 0
	retry := true
	for retry {
		retryCount++

		err, ctx := self.makeRequest(record)

		switch err {
		case httpNetworkError:
			Debug("%s, retrying", ctx)
			retry = true
			break
		case httpNonOkCode:
			Error("%s: received %d", httpNonOkCode, ctx)
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

func (self *HTTP) makeRequest(record intf.Record) (error, interface{}) {
	url := self.url.Apply(record)
	buffer := strings.NewReader(self.encode(record))
	req, err := http.NewRequest(self.method, url, buffer)

	if self.contentType != "" {
		req.Header.Set("Content-Type", self.contentType)
	}

	for header, value := range self.headers {
		req.Header.Set(header, value.Apply(record))
	}

	resp, err := self.client.Do(req)
	if err != nil {
		return httpNetworkError, err
	}

	io.Copy(ioutil.Discard, resp.Body)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(record, url, resp, string(body))
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		print(resp.StatusCode)
		return httpNonOkCode, nil
	}

	return nil, nil
}

func (self *HTTP) encode(record intf.Record) string {
	switch self.format {
	case "json":
		return self.encodeToJSON(record)
	case "form":
		return self.encodeToForm(record)
	default:
		Critical("Invalid encode format '%s'", self.format)
	}

	return ""
}

func (self *HTTP) encodeToJSON(record intf.Record) string {
	json, err := json.MarshalIndent(record, " ", "    ")
	if err != nil {
		Error("JSON Error %s", err)
	}

	self.transferred += len(json)
	return string(json)
}

func (self *HTTP) encodeToForm(record intf.Record) string {
	values, err := form.EncodeToValues(record)
	if err != nil {
		Error("Form Error %s", err)
	}

	return values.Encode()
}
