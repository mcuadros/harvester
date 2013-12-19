package intf

type Format interface {
	Parse(line string) map[string]string
}
