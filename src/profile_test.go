package collector

import (
	"testing"
)

func TestBuildFromConfig(t *testing.T) {
	var raw = string(`
		[output "test"]
		index = bar
		[input "test"]
		type = bar
		file = foo
		[format "test"]
		type = bar
		fields = foo
	`)

	config := NewConfig()
	config.Load(raw)

	profile := NewProfile(config)
}
