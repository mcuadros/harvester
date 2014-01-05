package harvesterd

import (
	. "harvesterd/intf"
	. "harvesterd/logger"

	"harvesterd/format"
	"harvesterd/input"
	"harvesterd/output"
	"harvesterd/processor"
)

type Container struct {
}

var containerInstance *Container = new(Container)

func GetContainer() *Container {
	return containerInstance
}

func (self *Container) GetFormat(key string) Format {
	jsonConfig, ok := GetConfig().Format_JSON[key]
	if ok {
		return format.NewJSON(jsonConfig)
	}

	csvConfig, ok := GetConfig().Format_CSV[key]
	if ok {
		return format.NewCSV(csvConfig)
	}

	regExpConfig, ok := GetConfig().Format_RegExp[key]
	if ok {
		return format.NewRegExp(regExpConfig)
	}

	apache2Config, ok := GetConfig().Format_Apache2[key]
	if ok {
		return format.NewApache2(apache2Config)
	}

	nginxConfig, ok := GetConfig().Format_Nginx[key]
	if ok {
		return format.NewNginx(nginxConfig)
	}

	Critical("Unable to find '%s' format definition", key)
	return nil
}

func (self *Container) GetInput(key string) Input {
	fileConfig, ok := GetConfig().Input_File[key]
	if ok {
		format := self.GetFormat(fileConfig.Format)
		return input.NewFile(fileConfig, format)
	}

	tailConfig, ok := GetConfig().Input_Tail[key]
	if ok {
		format := self.GetFormat(tailConfig.Format)
		return input.NewTail(tailConfig, format)
	}

	Critical("Unable to find '%s' input definition", key)
	return nil
}

func (self *Container) GetReader(key string) *Reader {
	config, ok := GetConfig().Reader[key]
	if !ok {
		return nil
	}

	inputs := make([]Input, len(config.Input))
	for i, key := range config.Input {
		inputs[i] = self.GetInput(key)
	}

	processors := make([]PostProcessor, len(config.Processor))
	for i, key := range config.Processor {
		processors[i] = self.GetPostProcessor(key)
	}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.SetProcessors(processors)

	return reader
}

func (self *Container) GetOutput(key string) Output {
	esConfig, ok := GetConfig().Output_Elasticsearch[key]
	if ok {
		return output.NewElasticsearch(esConfig)
	}

	mgoConfig, ok := GetConfig().Output_Mongo[key]
	if ok {
		return output.NewMongo(mgoConfig)
	}

	dummyConfig, ok := GetConfig().Output_Dummy[key]
	if ok {
		return output.NewDummy(dummyConfig)
	}

	Critical("Unable to find '%s' output definition", key)
	return nil
}

func (self *Container) GetPostProcessor(key string) PostProcessor {
	anonConfig, ok := GetConfig().Processor_Anonymize[key]
	if ok {
		return processor.NewAnonymize(anonConfig)
	}

	metricsConfig, ok := GetConfig().Processor_Metrics[key]
	if ok {
		return processor.NewMetrics(metricsConfig)
	}

	Critical("Unable to find '%s' processor definition", key)
	return nil
}

func (self *Container) GetWriter() *Writer {
	config := GetConfig().Writer

	outputs := make([]Output, len(config.Output))
	for i, key := range config.Output {
		outputs[i] = self.GetOutput(key)
	}

	readers := make([]*Reader, len(config.Reader))
	for i, key := range config.Reader {
		readers[i] = self.GetReader(key)
	}

	writer := NewWriter()
	writer.SetOutputs(outputs)
	writer.SetReaders(readers)
	writer.SetThreads(config.Threads)

	return writer
}
