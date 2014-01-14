package harvesterd

import (
	"harvesterd/format"
	"harvesterd/input"
	"harvesterd/logger"
	"harvesterd/output"
	"harvesterd/processor"
	"reflect"
	"strings"
)

import "code.google.com/p/gcfg"

type Config struct {
	Logger               *logger.LoggerConfig
	Writer               map[string]*WriterConfig
	Reader               map[string]*ReaderConfig
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

type Definition struct {
	Name          string
	AllowMultiple bool
	Fields        []*FieldDefinition
}

type FieldDefinition struct {
	Name        string
	Type        string
	Description string
	Default     string
}

var configInstance *Config = new(Config)

func GetConfig() *Config {
	return configInstance
}

func (self *Config) Load(ini string) {
	err := gcfg.ReadStringInto(self, ini)
	if err != nil {
		logger.Critical("error: cannot parse config", err)
	}
}

func (self *Config) LoadFile(filename string) {
	err := gcfg.ReadFileInto(self, filename)
	if err != nil {
		logger.Critical("erro:", err)
	}
}

func (self *Config) GetDescription() []*Definition {
	typeObject := reflect.TypeOf(configInstance).Elem()

	fields := self.getFieldsFromType(typeObject)
	definitions := make([]*Definition, len(fields))
	for index, field := range fields {
		definitions[index] = self.processField(field)
	}

	return definitions
}

func (self *Config) getFieldsFromType(typeObject reflect.Type) []reflect.StructField {
	typeObject.NumField()

	count := typeObject.NumField()
	results := make([]reflect.StructField, count)
	for i := 0; i < count; i++ {
		results[i] = typeObject.Field(i)
	}

	return results
}

func (self *Config) processField(field reflect.StructField) *Definition {
	def := new(Definition)

	var typeObject reflect.Type

	switch field.Type.Kind() {
	case reflect.Ptr:
		typeObject = field.Type.Elem()
		def.AllowMultiple = false
	case reflect.Map:
		typeObject = field.Type.Elem().Elem()
		def.AllowMultiple = true
	}

	def.Name = strings.Replace(strings.ToLower(field.Name), "_", "-", -1)

	for _, field := range self.getFieldsFromType(typeObject) {
		def.Fields = append(def.Fields, &FieldDefinition{
			Name:        strings.ToLower(field.Name),
			Type:        field.Type.String(),
			Default:     field.Tag.Get("default"),
			Description: field.Tag.Get("description"),
		})
	}

	return def
}
