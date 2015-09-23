package mutate

import (
	"fmt"
)

type OperationId int

const (
	CAST OperationId = 1 << iota
)

type Operation struct {
	Id     OperationId
	Field  []string
	Params []string
}

// Applies the Operation o to a record, recursively entering the
// maps and slices to get to the element in the path specified in o.Field
func (o Operation) Apply(record interface{}) error {
	var key string
	var isLeaf bool
	if len(o.Field) > 0 {
		key = o.Field[0]
		o.Field = o.Field[1:]
	}
	isLeaf = len(o.Field) == 0

	switch record.(type) {
	case []interface{}:
		r := record.([]interface{})
		if key == "*" {
			var firsterr error
			for i := range r {
				var err error
				if isLeaf {
					r[i], err = o.execute(r[i])
				} else {
					err = o.Apply(r[i])
				}
				if firsterr == nil && err != nil {
					firsterr = err
				}
			}
			return firsterr
		} else {
			return fmt.Errorf("mutate: slice found and key `%s` was in the path", key)
		}
	case map[string]interface{}:
		r := record.(map[string]interface{})
		if key == "*" {
			var firsterr error
			for k := range r {
				var err error
				if isLeaf {
					r[k], err = o.execute(r[k])
				} else {
					err = o.Apply(r[k])
				}
				if firsterr == nil && err != nil {
					firsterr = err
				}
			}
			return firsterr
		} else if _, ok := r[key]; ok {
			var err error
			if isLeaf {
				r[key], err = o.execute(r[key])
			} else {
				err = o.Apply(r[key])
			}
			return err
		} else {
			return fmt.Errorf("mutate: key `%s` was not found in subtree `%s`", key, r)
		}
	default:
		return fmt.Errorf("mutate: not map or slice element found in the path")
	}
	return nil
}

// Executes the action with id o.ID over `value`
func (o *Operation) execute(value interface{}) (interface{}, error) {
	switch o.Id {
	case CAST:
		if len(o.Params) < 1 {
			return value, fmt.Errorf("mutate: wrong number of params for function `cast`, at least 1 expected, %s found", o.Params)
		}
		return cast(value, o.Params[0], o.Params[1:]...)
	}
	return value, fmt.Errorf("mutate: operation %d not implemented", o.Id)
}
