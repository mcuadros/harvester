package intf

type Output interface {
	PutRecord(record map[string]string) bool
}
