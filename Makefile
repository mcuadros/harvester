# Package configuration
PACKAGE = harvesterd
HOMEPAGE = https://github.com/mcuadros/harvesterd
DESCRIPTION = low footprint collector and parser for events and logs
SUBPACKAGES = harvesterd/util \
harvesterd/input \
harvesterd/output \
harvesterd/format \
harvesterd/processor \
harvesterd/processor/metric

COMMANDS =	tool/harvesterd.go
DEPENDENCIES = gopkg.in/check.v1 \
code.google.com/p/go.tools/cmd/cover \
code.google.com/p/gcfg \
github.com/jarod/log4go \
github.com/ActiveState/tail \
github.com/mcuadros/go-syslog \
github.com/mcuadros/go-defaults \
github.com/rcrowley/go-metrics \
github.com/stretchr/objx \
github.com/ajg/form \
labix.org/v2/mgo \
labix.org/v2/mgo/bson

# Environment
BASE_PATH := $(shell pwd)
REPORT_PATH := $(BASE_PATH)/reports
BUILD_PATH := $(BASE_PATH)/build
BIN_PATH := $(BUILD_PATH)/bin
PACKAGE_PATH := $(BUILD_PATH)/src/$(PACKAGE)
PACKAGE_BASE := $(shell dirname $(PACKAGE_PATH))
ALL_PACKAGES := $(PACKAGE) $(SUBPACKAGES)
INSTALL_PATH ?= /opt/harvesterd
VERSION ?= $(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOINSTALL = $(GOCMD) install
GOTEST = $(GOCMD) test
GOPATH = $(BUILD_PATH)
export GOPATH

# FPM
FPMCMD = fpm

# Rules
all: test build 

build: dependencies
	for binary in $(COMMANDS); do \
		$(GOBUILD) -ldflags "-X main.version $(VERSION)" -o $(BIN_PATH)/$${binary%%.*} $(BASE_PATH)/$$binary; \
	done

install:
	$(GOINSTALL) $(BASE_PATH)/$*

test: dependencies
	for package in $(ALL_PACKAGES); do \
		$(GOTEST) $$package; \
	done

coverage: dependencies
	mkdir -p $(REPORT_PATH)
	for package in $(ALL_PACKAGES); do \
		$(GOTEST) -coverprofile=$(REPORT_PATH)/profile.$$(echo $$package | sed 's/\//\./g').out $$package; \
		$(GOCMD) tool cover -html=$(REPORT_PATH)/profile.$$(echo $$package | sed 's/\//\./g').out -o $(REPORT_PATH)/coverage.$$(echo $$package | sed 's/\//\./g').html; \
	done

dependencies:
	mkdir -p $(PACKAGE_BASE)
	ln -s $(BASE_PATH)/src $(PACKAGE_PATH) 2> /dev/null || true
	for i in $(DEPENDENCIES); do $(GOCMD) get $$i; done

rpm: build
	$(FPMCMD) -s dir -t rpm -n $(PACKAGE) -v $(VERSION) \
		--description "$(DESCRIPTION)" \
		--url "$(HOMEPAGE)" \
			$(foreach binary,$(COMMANDS),$(BIN_PATH)/${subst .go,,${binary}}=$(INSTALL_PATH)/bin/) \
			package/rpm/harvesterd-initd=/etc/init.d/harvesterd

clean:
	echo $(VERSION)
	rm -rf $(BUILD_PATH)
	rm -rf $(BIN_PATH)
	rm -rf $(REPORT_PATH)

	$(GOCLEAN) .
