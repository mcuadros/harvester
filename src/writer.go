package collector

import (
	"sync"
)

type Writer interface {
	WriteFromChannel(channel chan map[string]string, wait sync.WaitGroup)
	GetCounters() (int, int, int)
	ResetCounters()
}
