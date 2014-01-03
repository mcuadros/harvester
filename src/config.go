package harvesterd

import (
	"fmt"
	"harvesterd/format"
	"harvesterd/input"
	"harvesterd/logger"
	"harvesterd/output"
	"harvesterd/processor"
	"os"
)

import "code.google.com/p/gcfg"

type Config struct {
	Reader               ReaderConfig
	Writer               WriterConfig
	Logger               logger.LoggerConfig
	Format_JSON          map[string]*format.JSONConfig
	Format_CSV           map[string]*format.CSVConfig
	Format_RegExp        map[string]*format.RegExpConfig
	Format_Apache2       map[string]*format.Apache2Config
	Format_Nginx         map[string]*format.NginxConfig
	Input_File           map[string]*input.FileConfig
	Input_Tail           map[string]*input.TailConfig
	Output_Elasticsearch map[string]*output.ElasticsearchConfig
	Output_Mongo         map[string]*output.MongoConfig
	Output_Dummy         map[string]*output.DummyConfig
	Processor_Anonymize  map[string]*processor.AnonymizeConfig
	Processor_Metrics    map[string]*processor.MetricsConfig
}

var configInstance *Config = new(Config)

func GetConfig() *Config {
	return configInstance
}

func (self *Config) Load(ini string) {
	err := gcfg.ReadStringInto(self, ini)
	if err != nil {
		fmt.Println("error: cannot parse config", err)
		os.Exit(1)
	}
}

func (self *Config) LoadFile(filename string) {
	err := gcfg.ReadFileInto(self, filename)
	if err != nil {
		fmt.Println("erro:", err)
		os.Exit(1)
	}
}
