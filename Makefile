# Makefile for a go project
#
# Author: Jon Eisen
# 	site: joneisen.me
# 	
# Targets:
# 	all: Builds the code
# 	build: Builds the code
# 	fmt: Formats the source files
# 	clean: cleans the code
# 	install: Installs the code to the GOPATH
# 	iref: Installs referenced projects
#	test: Runs the tests
#	
#  Blog post on it: http://joneisen.me/post/25503842796
#

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i

# Package lists
PACKAGE_NAME := collector
TOPLEVEL_PKG := .
BASE_LIST := collector
IMPL_LIST := collector/input collector/output collector/format
CMD_LIST :=	tool/collectord.go
BIN_PATH := bin
DEPENDENCIES_LIST = launchpad.net/gocheck \
code.google.com/p/go.tools/cmd/cover \
code.google.com/p/gcfg \
github.com/jarod/log4go \
github.com/ActiveState/tail \
github.com/codegangsta/cli \
labix.org/v2/mgo \
labix.org/v2/mgo/bson


# List building
ALL_LIST = $(INT_LIST) $(BASE_LIST) $(IMPL_LIST) $(CMD_LIST)

BUILD_LIST = $(foreach int, $(CMD_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(ALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(ALL_LIST), $(int)_install)
COVERAGE_LIST = $(foreach int, $(BASE_LIST) $(IMPL_LIST), $(int)_coverage)
TEST_LIST = $(foreach int, $(BASE_LIST) $(IMPL_LIST), $(int)_test)
FMT_TEST = $(foreach int, $(ALL_LIST), $(int)_fmt)
IREF_LIST = $(foreach int, $(BASE_LIST), $(int)_iref)

GOPATH := $(shell pwd)/build
PACKAGE_PATH := $(GOPATH)/src/$(PACKAGE_NAME)
PACKAGE_BASE := $(shell dirname $(PACKAGE_PATH))
REPORT_PATH := $(shell pwd)/reports
export GOPATH

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: iref test build 
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(BUILD_LIST)
test: iref $(TEST_LIST)
coverage: iref $(COVERAGE_LIST)
iref: $(IREF_LIST)

$(BUILD_LIST): %_build:
	$(GOBUILD) -o $(subst .go,,$(BIN_PATH)/$(shell basename $(*))) $(TOPLEVEL_PKG)/$*
$(CLEAN_LIST): %_clean:
	rm -rf $(GOPATH)/src
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(IREF_LIST): %_iref:
	mkdir -p $(PACKAGE_BASE)
	ln -s $(shell pwd)/src $(PACKAGE_PATH) 2> /dev/null || true
	for i in $(DEPENDENCIES_LIST); do $(GOCMD) get $$i; done
$(TEST_LIST): %_test:
	$(GOTEST) $*
$(COVERAGE_LIST): %_coverage:
	mkdir -p $(REPORT_PATH)
	$(GOTEST) -coverprofile=$(REPORT_PATH)/profile.$(subst /,.,$(*)).out $*
	$(GOCMD) tool cover -html=$(REPORT_PATH)/profile.$(subst /,.,$(*)).out -o $(REPORT_PATH)/coverage.$(subst /,.,$(*)).html

