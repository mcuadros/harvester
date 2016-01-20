package harvester

import (
	"github.com/mcuadros/harvester/src/format"
	"github.com/mcuadros/harvester/src/input"
	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"
	"github.com/mcuadros/harvester/src/output"
	"github.com/mcuadros/harvester/src/processor"
)

type OutputsFactory func() []intf.Output

type Container struct {
}

var containerInstance *Container = new(Container)

func GetContainer() *Container {
	return containerInstance
}

func (c *Container) GetFormat(key string) intf.Format {
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

func (c *Container) GetInput(key string) intf.Input {
	fileConfig, ok := GetConfig().Input_File[key]
	if ok {
		format := c.GetFormat(fileConfig.Format)
		return input.NewFile(fileConfig, format)
	}

	tailConfig, ok := GetConfig().Input_Tail[key]
	if ok {
		format := c.GetFormat(tailConfig.Format)
		return input.NewTail(tailConfig, format)
	}

	s3Config, ok := GetConfig().Input_S3[key]
	if ok {
		format := c.GetFormat(s3Config.Format)
		return input.NewS3(s3Config, format)
	}

	mongoConfig, ok := GetConfig().Input_Mongo[key]
	if ok {
		return input.NewMongo(mongoConfig)
	}

	Critical("Unable to find '%s' input definition", key)
	return nil
}

func (c *Container) GetReader(key string) *Reader {
	config, ok := GetConfig().Reader[key]
	if !ok {
		return nil
	}

	inputs := make([]intf.Input, len(config.Input))
	for i, key := range config.Input {
		inputs[i] = c.GetInput(key)
	}

	processors := make([]intf.PostProcessor, len(config.Processor))
	for i, key := range config.Processor {
		processors[i] = c.GetPostProcessor(key)
	}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.SetProcessors(processors)

	return reader
}

func (c *Container) GetOutput(key string) intf.Output {
	httpConfig, ok := GetConfig().Output_HTTP[key]
	if ok {
		return output.NewHTTP(httpConfig)
	}

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

func (c *Container) GetPostProcessor(key string) intf.PostProcessor {
	anonConfig, ok := GetConfig().Processor_Anonymize[key]
	if ok {
		return processor.NewAnonymize(anonConfig)
	}

	metricsConfig, ok := GetConfig().Processor_Metrics[key]
	if ok {
		return processor.NewMetrics(metricsConfig)
	}

	mConfig, ok := GetConfig().Processor_Mutate[key]
	if ok {
		return processor.NewMutate(mConfig)
	}

	Critical("Unable to find '%s' processor definition", key)
	return nil
}

func (c *Container) GetWriter(key string) *Writer {
	config, ok := GetConfig().Writer[key]
	if !ok {
		return nil
	}

	outputsFactory := func() []intf.Output {
		outputs := make([]intf.Output, len(config.Output))
		for i, key := range config.Output {
			outputs[i] = c.GetOutput(key)
		}

		return outputs
	}

	readers := make([]*Reader, len(config.Reader))
	for i, key := range config.Reader {
		readers[i] = c.GetReader(key)
	}

	writer := NewWriter()
	writer.SetOutputsFactory(outputsFactory)

	if len(readers) == 0 {
		Critical("Invalid writer config: alteast one reader should be provided.")
	}
	writer.SetReaders(readers)

	if config.Threads == 0 {
		Critical("Invalid writer config: num. of threads should be >0.")
	}
	writer.SetThreads(config.Threads)

	return writer
}

func (c *Container) GetWriterGroup() *WriterGroup {
	writers := make([]intf.Writer, len(GetConfig().Writer))

	i := 0
	for key, _ := range GetConfig().Writer {
		writers[i] = c.GetWriter(key)
		i++
	}

	writerGroup := NewWriterGroup()
	writerGroup.SetWriters(writers)

	return writerGroup
}
