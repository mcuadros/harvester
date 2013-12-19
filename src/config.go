package collector

import (
	"./format"
	"./input"
	"fmt"
)

import "code.google.com/p/gcfg"

type Config struct {
	Basic struct {
		Threads int
	}
	Reader        ReaderConfig
	Logger        LoggerConfig
	Format_CSV    map[string]*format.CSVConfig
	Format_RegExp map[string]*format.RegExpConfig
	Input_File    map[string]*input.FileConfig
	Input_Tail    map[string]*input.TailConfig
	ElasticSearch WriterElasticSearchConfig
}

var configInstance *Config = new(Config)

func GetConfig() *Config {
	return configInstance
}

func (self *Config) Load(ini string) {
	err := gcfg.ReadStringInto(self, ini)
	if err != nil {
		panic(fmt.Sprintf("open config: %v", err))
	}
}

func (self *Config) LoadFile(filename string) {
	err := gcfg.ReadFileInto(self, filename)
	if err != nil {
		panic(fmt.Sprintf("open config: %v", err))
	}
}