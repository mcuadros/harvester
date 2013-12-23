package collector

import (
	"collector/format"
	"collector/input"
	"collector/intf"
	. "collector/logger"
	"collector/output"
)

type Container struct {
}

var containerInstance *Container = new(Container)

func GetContainer() *Container {
	return containerInstance
}

func (self *Container) GetFormat(key string) intf.Format {
	csvConfig, ok := GetConfig().Format_CSV[key]
	if ok {
		return format.NewCSV(csvConfig)
	}

	regExpConfig, ok := GetConfig().Format_RegExp[key]
	if ok {
		return format.NewRegExp(regExpConfig)
	}

	Critical("Unable to find '%s' format definition", key)
	return nil
}

func (self *Container) GetInput(key string) intf.Input {
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

func (self *Container) GetReader() *Reader {
	config := GetConfig().Reader

	inputs := make([]intf.Input, len(config.Input))
	for i, key := range config.Input {
		inputs[i] = self.GetInput(key)
	}

	return NewReader(inputs)
}

func (self *Container) GetOutput(key string) intf.Output {
	esConfig, ok := GetConfig().Output_Elasticsearch[key]
	if ok {
		return output.NewElasticsearch(esConfig)
	}

	mgoConfig, ok := GetConfig().Output_Mongo[key]
	if ok {
		return output.NewMongo(mgoConfig)
	}

	Critical("Unable to find '%s' output definition", key)
	return nil
}

func (self *Container) GetWriter() *Writer {
	config := GetConfig().Writer

	outputs := make([]intf.Output, len(config.Output))
	for i, key := range config.Output {
		outputs[i] = self.GetOutput(key)
	}

	return NewWriter(outputs, config.Threads)
}
