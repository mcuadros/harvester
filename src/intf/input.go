package intf

type Input interface {
	GetLine() string
	GetRecord() map[string]string
	IsEOF() bool
}
