package processor

import (
	"regexp"
	"strings"

	"github.com/mcuadros/harvesterd/src/intf"
	. "github.com/mcuadros/harvesterd/src/logger"
	"github.com/mcuadros/harvesterd/src/processor/mutate"
)

const FIELDSEP = '.'

type MutateConfig struct {
	Verbose bool
	Cast    []string
}

func (mc *MutateConfig) ParseOperations() []*mutate.Operation {
	var operations []*mutate.Operation
	if len(mc.Cast) > 0 {
		for _, rawparams := range mc.Cast {
			op := mc.parseOperation(mutate.CAST, rawparams)
			operations = append(operations, op)
		}
	}
	return operations
}

func (p *MutateConfig) parseOperation(id mutate.OperationId, rawparams string) *mutate.Operation {
	// Parse the raw string into:
	// * single keywords with no spaces
	// * groups of single-quoted strings
	re := regexp.MustCompile("[^' ]+|'[^']+'")
	splitted := re.FindAllString(rawparams, -1)
	for i, s := range splitted {
		splitted[i] = strings.Trim(s, "'")
	}
	return &mutate.Operation{
		Id:     id,
		Field:  strings.Split(splitted[0], string(FIELDSEP)),
		Params: splitted[1:],
	}
}

type Mutate struct {
	operations []*mutate.Operation
	channel    chan intf.Record
	verbose    bool
	isAlive    bool
}

func NewMutate(config *MutateConfig) *Mutate {
	processor := Mutate{
		operations: []*mutate.Operation{},
	}
	processor.SetConfig(config)
	processor.Setup()

	return &processor
}

func (p *Mutate) SetConfig(config *MutateConfig) {
	p.verbose = config.Verbose
	p.operations = config.ParseOperations()
}

func (p *Mutate) SetChannel(channel chan intf.Record) {
	p.channel = channel
}

func (p *Mutate) Do(record intf.Record) bool {
	for _, op := range p.operations {
		err := op.Apply(map[string]interface{}(record))
		if err != nil && p.verbose {
			Warning(err.Error())
		}
	}
	return true
}

func (p *Mutate) Setup() {
	p.isAlive = true
}

func (p *Mutate) Teardown() {
	p.isAlive = false
}
