package intf

type PostProcessor interface {
	Do(record map[string]string)
}
