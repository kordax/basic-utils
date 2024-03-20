[![Tests](https://github.com/kordax/basic-utils/actions/workflows/Tests.yml/badge.svg?branch=main)](https://github.com/kordax/basic-utils/actions/workflows/Tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kordax/basic-utils)](https://goreportcard.com/report/github.com/kordax/basic-utils)
[![Coverage](https://raw.githubusercontent.com/kordax/basic-utils/badges/.badges/main/coverage.svg)](https://github.com/kordax/basic-utils/tree/badges)

# Basic Utils

This repository contains a collection of utility libraries implemented in Go, designed to assist in a variety of common
programming tasks. Each module addresses a particular set of functions or data structures.

## Minimum Go Version Requirement

To use or contribute to this project, you need to have at least Go 1.22 installed.
This is due to the usage of features and packages introduced in this version.

## Modules

- **array-utils**: Utilities related to array manipulations and operations.

- **async-utils**: Utilities that help to organize async operations.

- **file-utils**: Utilities for efficient file handling.

- **map-utils**: Helper functions for working with maps in Go.

- **math-utils**: Mathematical utilities and helpers.

- **number**: Versatile numeric representation.

- **opt**: Optional type implementations, which may hold a value or represent the absence of one.

- **queue**: Implements both a FIFO (First-In-First-Out) queue and a priority queue with thread safety and various
  utility functions.

- **ref-utils**: Utilities related to references.

- **sql-utils**: Utilities related to sql types and methods.

- **str-utils**: Utilities related to string operations.

## Installation

Make sure you have Go installed on your machine. Then, use `go get` to install the package:

```shell
go get -u github.com/kordax/basic-utils
```

## Usage

You can import each module individually or import the main module and it depends on your needs. For example, to use the
queue library:

```shell
go get -u github.com/kordax/basic-utils/queue
```

then...

```go
import "github.com/kordax/basic-utils/queue"
```

Or you can do:

```shell
go get -u github.com/kordax/basic-utils
```

then...

```go
import "github.com/kordax/basic-utils"
```

Then, refer to the individual documentation or code comments of each module for specific usage patterns.

## Static Analysis

The repository also includes a staticcheck.conf file, indicating that it might be set up to use the staticcheck tool for
static code analysis. Run staticcheck in the root directory to perform a code quality check.

## Author

Developed by [@kordax](mailto:dmorozov@valoru-software.com) (Dmitry Morozov)

[Valoru Software](https://valoru-software.com)
