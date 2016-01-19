package intf

type Input interface {
	GetRecord() Record
	IsEOF() bool
	Teardown()
}
