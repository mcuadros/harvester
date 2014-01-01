package intf

type Format interface {
	Parse(line string) Record
}
