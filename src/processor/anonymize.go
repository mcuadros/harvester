package processor

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/hex"
	"harvesterd/intf"
	"strings"
)

type AnonymizeConfig struct {
	Fields       string
	Hash         string
	EmailSupport bool
}

type Anonymize struct {
	fields []string
	hash   crypto.Hash
	email  bool
}

func NewAnonymize(config *AnonymizeConfig) *Anonymize {
	processor := new(Anonymize)
	processor.SetConfig(config)

	return processor
}

//Just for match the interface
func (p *Anonymize) SetChannel(channel chan intf.Record) {

}

func (p *Anonymize) SetConfig(config *AnonymizeConfig) {
	for _, field := range strings.Split(config.Fields, ",") {
		p.fields = append(p.fields, field)
	}

	if config.Hash == "" {
		config.Hash = "md5"
	}

	switch config.Hash {
	case "md5":
		p.hash = crypto.MD5
	case "sha1":
		p.hash = crypto.SHA1
	case "sha256":
		p.hash = crypto.SHA256
	case "sha512":
		p.hash = crypto.SHA512
	}

	p.email = config.EmailSupport
}

func (p *Anonymize) Do(record intf.Record) bool {
	for _, field := range p.fields {
		_, ok := record[field]
		if ok {
			p.encodeField(record, field)
		}
	}

	return true
}

func (p *Anonymize) encodeField(record intf.Record, field string) {
	if p.email {
		parts := strings.Split(record[field].(string), "@")
		parts[0] = p.encodeString(parts[0])

		record[field] = strings.Join(parts, "@")
	} else {
		record[field] = p.encodeString(record[field].(string))
	}
}

func (p *Anonymize) encodeString(value string) string {
	h := p.hash.New()
	h.Write([]byte(value))

	return hex.EncodeToString(h.Sum(nil))
}

//Just for match the interface
func (p *Anonymize) Teardown() {

}
