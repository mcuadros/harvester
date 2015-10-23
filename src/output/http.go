package output

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"
	"github.com/mcuadros/harvester/src/util"

	"github.com/ajg/form"
	"github.com/mcuadros/go-defaults"
)

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

func (o *HTTP) SetConfig(config *HTTPConfig) {
	defaults.SetDefaults(config)

	o.url = util.NewTemplate(config.Url)
	o.format = config.Format
	o.contentType = config.ContentType
	o.method = config.Method
	o.timeout = time.Duration(config.Timeout)
	o.parseHeadersConfig(config.Header)

	o.createHTTPClient()
}

func (o *HTTP) parseHeadersConfig(headers []string) {
	o.headers = make(map[string]*util.Template, len(headers))

	for _, headerRaw := range headers {
		header := strings.Split(headerRaw, ",")
		if len(header) != 2 {
			Critical("Malformed header setting '%s'", header)
		}

		o.headers[header[0]] = util.NewTemplate(header[1])
	}
}

func (o *HTTP) createHTTPClient() {
	var dialer = &net.Dialer{Timeout: o.timeout * time.Second}
	var transport = &http.Transport{Dial: dialer.Dial}

	o.client = &http.Client{Transport: transport}
}

func (o *HTTP) PutRecord(record intf.Record) bool {
	retryCount := 0
	retry := true
	for retry {
		retryCount++

		err, ctx := o.makeRequest(record)

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

func (o *HTTP) makeRequest(record intf.Record) (error, interface{}) {
	url := o.url.Apply(record)
	buffer := strings.NewReader(o.encode(record))
	req, err := http.NewRequest(o.method, url, buffer)

	if o.contentType != "" {
		req.Header.Set("Content-Type", o.contentType)
	}

	for header, value := range o.headers {
		req.Header.Set(header, value.Apply(record))
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return httpNetworkError, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		Debug("Url %s / StatusCode %d", url, resp.StatusCode)
		Debug("Request %s / Response %s", resp, string(body))

		return httpNonOkCode, nil
	} else {
		io.Copy(ioutil.Discard, resp.Body)
	}

	return nil, nil
}

func (o *HTTP) encode(record intf.Record) string {
	switch o.format {
	case "json":
		return o.encodeToJSON(record)
	case "form":
		return o.encodeToForm(record)
	default:
		Critical("Invalid encode format '%s'", o.format)
	}

	return ""
}

func (o *HTTP) encodeToJSON(record intf.Record) string {
	json, err := json.MarshalIndent(record, " ", "    ")
	if err != nil {
		Error("JSON Error %s", err)
	}

	o.transferred += len(json)
	return string(json)
}

func (o *HTTP) encodeToForm(record intf.Record) string {
	values, err := form.EncodeToValues(record)
	if err != nil {
		Error("Form Error %s", err)
	}

	return values.Encode()
}
