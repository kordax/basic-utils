module github.com/kordax/basic-utils-benchmarks/ristretto

go 1.27.0

require (
	github.com/dgraph-io/ristretto/v2 v2.4.2
	github.com/kordax/basic-utils/v4 v4.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20240924180020-3414d57e47da // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	golang.org/x/exp v0.0.0-20260824195058-e88cd73687aa // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/kordax/basic-utils/v4 => ../..
