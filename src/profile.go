package collector

import (
	"./format"
	"./input"
)

type Profile struct {
	Format Format
	Input  Input
	Writer Writer
}

func NewProfile(config *ProfileConfig) *Profile {
	profile := new(Profile)
	profile.BuildFromConfig(config)

	return profile
}

func (self *Profile) BuildFromConfig(config *ProfileConfig) {
	self.buildFormatFromConfig(config.Format)
	self.buildInputFromConfig(config.Input)
}

func (self *Profile) buildFormatFromConfig(config *FormatConfig) {
	switch config.Type {
	case "csv":
		self.Format = format.NewCSV(config)
	case "regexp":
		self.Format = format.NewRegExp(config)
	}
}

func (self *Profile) buildInputFromConfig(config *InputConfig) {
	switch config.Type {
	case "file":
		self.Input = input.NewFile(config)
	case "tail":
		self.Input = input.NewTail(config)
	}
}
