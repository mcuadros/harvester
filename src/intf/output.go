package intf

type Output interface {
	PutRecord(map[string]string) bool
}
