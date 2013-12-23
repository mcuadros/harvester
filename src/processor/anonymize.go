package processor

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/hex"
	"hash"
	"strings"
)

type AnonymizeConfig struct {
	Fields       string
	Hash         string
	EmailSupport bool
}

type Anonymize struct {
	fields []string
	hash   hash.Hash
	email  bool
}

func NewAnonymize(config *AnonymizeConfig) *Anonymize {
	processor := new(Anonymize)
	processor.SetConfig(config)

	return processor
}

func (self *Anonymize) SetConfig(config *AnonymizeConfig) {
	for _, field := range strings.Split(config.Fields, ",") {
		self.fields = append(self.fields, field)
	}

	if config.Hash == "" {
		config.Hash = "md5"
	}

	switch config.Hash {
	case "md5":
		self.hash = crypto.MD5.New()
	case "sha1":
		self.hash = crypto.SHA1.New()
	case "sha256":
		self.hash = crypto.SHA256.New()
	case "sha512":
		self.hash = crypto.SHA512.New()
	}

	self.email = config.EmailSupport
}

func (self *Anonymize) Do(record map[string]string) {
	for _, field := range self.fields {
		_, ok := record[field]
		if ok {
			self.encodeField(record, field)
		}
	}
}

func (self *Anonymize) encodeField(record map[string]string, field string) {
	if self.email {
		parts := strings.Split(record[field], "@")
		parts[0] = self.encodeString(parts[0])

		record[field] = strings.Join(parts, "@")
	} else {
		record[field] = self.encodeString(record[field])
	}
}

func (self *Anonymize) encodeString(value string) string {
	self.hash.Reset()
	self.hash.Write([]byte(value))

	return hex.EncodeToString(self.hash.Sum(nil))
}
