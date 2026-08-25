[![Coverage](https://raw.githubusercontent.com/kordax/basic-utils/badges/.badges/main/coverage.svg)](https://github.com/kordax/basic-utils/tree/badges)

# Basic Utils

This repository contains a collection of utility libraries implemented in Go, designed to assist in a variety of common
programming tasks. Each module addresses a particular set of functions or data structures.

## Minimum Go Version Requirement

To use or contribute to this project, you need to have at least Go 1.27 installed.
This is due to the usage of features and packages introduced in this version.

## Modules

- [uarray](uarray/README.md): Slice utilities for searching, filtering, mapping, grouping, set-like operations, safe indexing, chunking, and string conversion.
- [uasync](uasync/README.md): Async execution helpers, futures, scheduled tasks, retry helpers, and grouped task execution.
- [ucache](ucache/README.md): Cache implementations, TTL support, composite keys, multi-cache structures, and managed cleanup wrappers.
- [ucast](ucast/README.md): Bi-directional utilities to convert basic types.
- [uconst](uconst/README.md): Shared constraints and common generic contracts.
- [uctx](uctx/README.md): Small context-related helpers.
- [uerror](uerror/README.md): Error handling helpers.
- [uevent](uevent/README.md): Channel watcher helpers.
- [ufile](ufile/README.md): File-system helpers for common read/write/list workflows.
- [umap](umap/README.md): Map utilities and multi-map implementations.
- [umath](umath/README.md): Mathematical utilities and helpers.
- [unumber](unumber/README.md): Versatile numeric representation and denomination helpers.
- [uonce](uonce/README.md): Helpers for one-time execution semantics.
- [uopt](uopt/README.md): Optional type implementation with JSON/SQL support and functional helpers.
- [uos](uos/README.md): Operating-system and environment variable helpers.
- [upair](upair/README.md): Pair type helpers.
- [uqueue](uqueue/README.md): FIFO, lock-free concurrent FIFO, and priority queues.
- [uref](uref/README.md): Reference helpers.
- [uset](uset/README.md): Set implementations and set algebra helpers.
- [usql](usql/README.md): SQL helper types and functions.
- [usrlz](usrlz/README.md): Serialization helpers.
- [usize](usize/README.md): Byte-sized data unit constants.
- [ustr](ustr/README.md): String helpers.
- [ustream](ustream/README.md): Experimental stream-style helpers for slice processing.

## Installation

Make sure you have Go installed on your machine. Then, use `go get` to install the package:

```shell
go get -u github.com/kordax/basic-utils/v4
```

## Usage

You can import each module individually or import the main module and it depends on your needs. For example, to use the
queue library:

```shell
go get -u github.com/kordax/basic-utils/v4/uqueue
```

then...

```go
import "github.com/kordax/basic-utils/v4/uqueue"
```

Then, refer to the individual documentation or code comments of each module for specific usage patterns.

## Static Analysis

The repository also includes a staticcheck.conf file, indicating that it might be set up to use the staticcheck tool for
static code analysis. Run staticcheck in the root directory to perform a code quality check.
