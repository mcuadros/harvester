package metric

import (
	. "harvesterd/intf"
)

type Counter struct {
	field        string
	count        map[string]int
	readChannel  chan chan map[string]int
	writeChannel chan counterIncrement
}

type counterIncrement struct {
	key   string
	value int
}

func NewCounter(field string) *Counter {
	counter := &Counter{
		field:        field,
		count:        make(map[string]int),
		readChannel:  make(chan chan map[string]int, 0),
		writeChannel: make(chan counterIncrement, 0),
	}

	counter.Boot()

	return counter
}

func (self *Counter) Boot() {
	go self.syncronized()
}

func (self *Counter) syncronized() {
	for {
		select {
		case ci := <-self.writeChannel:
			if len(ci.key) == 0 {
				return
			}

			self.count[ci.key] += ci.value
			break
		case cl := <-self.readChannel:
			nm := make(map[string]int)
			for k, v := range self.count {
				nm[k] = v
			}
			cl <- nm
			break
		}
	}
}

func (self *Counter) Process(record Record) {
	key := record[self.field].(string)

	self.writeChannel <- counterIncrement{key, 1}
}

func (self *Counter) GetField() string {
	return self.field
}

func (self *Counter) GetValue() interface{} {
	reply := make(chan map[string]int)
	self.readChannel <- reply

	return <-reply
}

func (self *Counter) Reset() {
	self.count = make(map[string]int)
}
