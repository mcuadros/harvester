package intf

type Input interface {
	GetLine() string
	IsEOF() bool
}
