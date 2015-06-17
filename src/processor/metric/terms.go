package metric

import . "github.com/mcuadros/harvesterd/src/intf"

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

func (m *Terms) Boot() {
	go m.syncronized()
}

func (m *Terms) syncronized() {
	for {
		select {
		case ci := <-m.writeChannel:
			if len(ci.key) == 0 {
				return
			}

			m.count[ci.key] += ci.value
			break
		case cl := <-m.readChannel:
			nm := make(map[string]int)
			for k, v := range m.count {
				nm[k] = v
			}
			cl <- nm
			break
		}
	}
}

func (m *Terms) Process(record Record) {
	switch record[m.field].(type) {
	case string:
		key := record[m.field].(string)
		m.writeChannel <- counterIncrement{key, 1}
	}
}

func (m *Terms) GetField() string {
	return m.field
}

func (m *Terms) GetValue() interface{} {
	reply := make(chan map[string]int)
	m.readChannel <- reply

	return <-reply
}

func (m *Terms) Reset() {
	m.count = make(map[string]int)
}
