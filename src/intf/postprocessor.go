package intf

type PostProcessor interface {
	Do(record Record)
}
