# Ristretto Comparison Benchmarks

This is a separate Go module for optional cache benchmarks against `github.com/dgraph-io/ristretto/v2`.

It intentionally lives outside the root module dependency graph, so the library does not depend on Ristretto and regular
root commands such as `go test ./...` do not run or resolve this module.

Run from this directory:

```shell
go test -run '^$' -bench=. -benchtime=1s -benchmem
```

Run from the repository root:

```shell
(cd benchmarks/ristretto && go test -run '^$' -bench=. -benchtime=1s -benchmem)
```
