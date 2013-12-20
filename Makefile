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
GOFMT=gofmt -w

# Package lists
PACKAGE_NAME := collector
TOPLEVEL_PKG := .
BASE_LIST := src
IMPL_LIST := src/input src/output
CMD_LIST :=	tool
BIN_PATH := bin/collector
DEPENDENCIES_LIST = launchpad.net/gocheck \
code.google.com/p/gcfg \
github.com/jarod/log4go \
github.com/ActiveState/tail


# List building
ALL_LIST = $(INT_LIST) $(BASE_LIST) $(IMPL_LIST) $(CMD_LIST)

BUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(ALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(ALL_LIST), $(int)_install)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_TEST = $(foreach int, $(ALL_LIST), $(int)_fmt)
IREF_LIST = $(foreach int, $(BASE_LIST), $(int)_iref)

GOPATH := $(shell pwd)/build
PACKAGE_PATH := $(GOPATH)/src/$(PACKAGE_NAME)
PACKAGE_BASE := $(shell dirname $(PACKAGE_PATH))
export GOPATH

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: iref build 
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: iref $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_TEST)

$(BUILD_LIST): %_build: %_fmt
	$(GOBUILD) -o $(BIN_PATH) $(TOPLEVEL_PKG)/$*
$(CLEAN_LIST): %_clean:
	rm -rf $(GOPATH)/src
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(IREF_LIST): %_iref:
	mkdir -p $(GOPATH)/src
	mkdir -p $(PACKAGE_BASE)
	ln -s $(shell pwd)/src $(PACKAGE_PATH) 2> /dev/null || true
	for i in $(DEPENDENCIES_LIST); do $(GOCMD) get $$i; done
$(TEST_LIST): %_test:
	$(GOTEST) $(TOPLEVEL_PKG)/$*
$(FMT_TEST): %_fmt:
	$(GOFMT) ./$*

