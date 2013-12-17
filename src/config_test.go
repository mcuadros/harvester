package collector

import (
	"testing"
)

func TestBasic(t *testing.T) {
	var raw = string(`
		[basic]
		threads = 10
	`)

	config := NewConfig()
	config.Load(raw)

	if config.Basic.Threads != 10 {
		t.Errorf("FAIL: Wrong loaded data")
	}
}

func TestFormat(t *testing.T) {
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

	if config.Format["test"].Type != "bar" {
		t.Errorf("FAIL: Wrong loaded data")
	}

	if config.Format["test"].Fields != "foo" {
		t.Errorf("FAIL: Wrong loaded data")
	}

	if config.Input["test"].Type != "bar" {
		t.Errorf("FAIL: Wrong loaded data")
	}

	if config.Input["test"].File != "foo" {
		t.Errorf("FAIL: Wrong loaded data")
	}
}

func TestProfiles(t *testing.T) {
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

	if len(config.Profiles) != 1 {
		t.Errorf("FAIL: Invalid count of profiles")
	}
}
