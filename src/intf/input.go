package intf

type Input interface {
	GetLine() string
	GetRecord() Record
	IsEOF() bool
	Finish()
}
