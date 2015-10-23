# Harvester [![Build Status](https://travis-ci.org/mcuadros/harvester.png?branch=master)](https://travis-ci.org/mcuadros/harvester)

Harvester is a data collector, parser and processor. It reads, modifies and writes data in a lightweight, highly configurable pipeline. It is written in Go, which ensures both high performance a low footprint to optimize system resources.

Both the input and output sources and the intermediate processors can be easily configured by declaring blocks in the config file. Currently CSV, JSON and apache2/nginx logs format are supported, and there is also a RegExp reader that allows reading any arbitrary source that can be parsed with regular expressions. Outputting to MongoDB, HTTP and ElasticSearch is also possible.

## Running Harvester

From the project root, run:

```sh
go get -t ./...
```

and compile with:

```sh
go build commands/harvester.go
```

so you get the self-contained binary. At this point you might want to move it to `/usr/local/bin` or any other place in your `PATH` you feel it's appropriate. To run harvester, and provided that `harvester` is in your path, type:

```sh
harvester -config /path/to/your/config/file.conf
```

## Configuring Harvester

The config file follows the git-config/`ini` format specified in the [gcfg](https://godoc.org/gopkg.in/gcfg.v1) package.

Read/write and processor modules are declared with `ini` sections, and their properties declared within the block:

```ini
[output-dummy "my-writer-name"]
print = true
```

In this block, a dummy output module was declared with the name "my-writer-name", and its flag `print` was set to true.

## Available Modules

There are five type of modules in Harvester.

* `input`: An `input` is a module that declares a source from which data can be read. It holds information on how to connecto to a certain location or server to stream data (i.e.: read a log in the filesystem, or connect to a certain folder in an Amazon S3 server).
* `output`: `output` models are, similarly to `input` ones, a declaration of a location or server where the data strem will be saved.
* `processor`: A `processor` is a module that is placed in the middle of the stream and updates, changes, removes or collects information about the data.
* `reader`: A `reader` is a module that connects an `input` source with one or more `processors`.
* `writer`: A `writer` is a module that connects a `reader` with one or more `output` modules.

There are several `input`, `output` and `processor` modules available. The documentation for each one of them is located in the `src/input`, `src/output` and `src/processor` packages, respectively.

## License

Copyright 2013-2015 Máximo Cuadros. Licensed under the [MIT License](LICENSE).
