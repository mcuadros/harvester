package intf

type Writer interface {
	IsAlive() bool
	Setup()
	Boot()
	GetCounters() (int32, int32, int32, int32)
	ResetCounters()
	Teardown()
}
