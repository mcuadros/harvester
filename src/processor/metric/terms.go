package metric

import (
	. "harvesterd/intf"
)

type Terms struct {
	field        string
	count        map[string]int
	readChannel  chan chan map[string]int
	writeChannel chan counterIncrement
}

type counterIncrement struct {
	key   string
	value int
}

func NewTerms(field string) *Terms {
	counter := &Terms{
		field:        field,
		count:        make(map[string]int),
		readChannel:  make(chan chan map[string]int, 0),
		writeChannel: make(chan counterIncrement, 0),
	}

	counter.Boot()

	return counter
}

func (self *Terms) Boot() {
	go self.syncronized()
}

func (self *Terms) syncronized() {
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

func (self *Terms) Process(record Record) {
	switch record[self.field].(type) {
	case string:
		key := record[self.field].(string)
		self.writeChannel <- counterIncrement{key, 1}
	}
}

func (self *Terms) GetField() string {
	return self.field
}

func (self *Terms) GetValue() interface{} {
	reply := make(chan map[string]int)
	self.readChannel <- reply

	return <-reply
}

func (self *Terms) Reset() {
	self.count = make(map[string]int)
}
