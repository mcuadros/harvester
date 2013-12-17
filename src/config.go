package collector

import (
	"./format"
	"./input"
	"fmt"
)

import "code.google.com/p/gcfg"

type ProfileConfig struct {
	Format *FormatConfig
	Input  *InputConfig
	Output *OutputConfig
}

type InputConfig struct {
	Type string
	File string
	input.FileConfig
	input.TailConfig
}

type FormatConfig struct {
	Type string
	format.CSVConfig
	format.RegExpConfig
}

type OutputConfig struct {
	Type string
	WriterElasticSearchConfig
}

type Config struct {
	Profiles []string
	Basic    struct {
		Threads int
	}
	Logger LoggerConfig
	Format map[string]*FormatConfig
	Input  map[string]*InputConfig
	Output map[string]*OutputConfig
}

func NewConfig() *Config {
	config := new(Config)

	return config
}

func (self *Config) Load(ini string) {
	err := gcfg.ReadStringInto(self, ini)
	if err != nil {
		panic(fmt.Sprintf("open config: %v", err))
	}

	self.initialize()
}

func (self *Config) LoadFile(filename string) {
	err := gcfg.ReadFileInto(self, filename)
	if err != nil {
		panic(fmt.Sprintf("open config: %v", err))
	}

	self.initialize()
}

func (self *Config) GetProfile(profile string) ProfileConfig {
	return ProfileConfig{
		Format: self.Format[profile],
		Input:  self.Input[profile],
		Output: self.Output[profile]}
}

func (self *Config) initialize() {
	self.loadProfiles()
	self.validate()
}

func (self *Config) loadProfiles() {
	keys := make(map[string]bool)

	for key, _ := range self.Format {
		keys[key] = true
	}

	for key, _ := range self.Input {
		keys[key] = true
	}

	for key, _ := range self.Output {
		keys[key] = true
	}

	for key, _ := range keys {
		self.Profiles = append(self.Profiles, key)
	}
}

func (self *Config) validate() {
	for _, key := range self.Profiles {
		if _, ok := self.Format[key]; !ok {
			self.throwInvalidPanic("format", key)
		}

		if _, ok := self.Input[key]; !ok {
			self.throwInvalidPanic("input", key)
		}

		if _, ok := self.Output[key]; !ok {
			self.throwInvalidPanic("output", key)
		}
	}
}

func (self *Config) throwInvalidPanic(group, key string) {
	panic(fmt.Sprintf("Missing %s in %s group", group, key))
}
