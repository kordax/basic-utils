module github.com/kordax/basic-utils-benchmarks/ristretto

go 1.26

toolchain go1.26.4

require (
	github.com/dgraph-io/ristretto/v2 v2.4.2
	github.com/kordax/basic-utils/v3 v3.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20240924180020-3414d57e47da // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976 // indirect
	golang.org/x/sys v0.45.0 // indirect
)

replace github.com/kordax/basic-utils/v3 => ../..
