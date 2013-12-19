package collector

import (
	"./format"
	"./input"
	"./intf"
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

	GetLogger().Critical("Unable to find '%s' format definition", key)
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

	GetLogger().Critical("Unable to find '%s' input definition", key)
	return nil
}