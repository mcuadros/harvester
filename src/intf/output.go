package intf

type Output interface {
	PutRecord(record Record) bool
}
