package intf

type PostProcessor interface {
	SetChannel(channel chan Record)
	Do(record Record) bool
	Finish()
}
