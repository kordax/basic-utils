module github.com/kordax/basic-utils

go 1.22.1

replace github.com/kordax/basic-utils/uarr => ./uarr

replace github.com/kordax/basic-utils/uasync => ./uasync

replace github.com/kordax/basic-utils/ufile => ./ufile

replace github.com/kordax/basic-utils/umap => ./umap

replace github.com/kordax/basic-utils/umath => ./umath

replace github.com/kordax/basic-utils/unum => ./unum

replace github.com/kordax/basic-utils/uopt => ./uopt

replace github.com/kordax/basic-utils/uqueue => ./uqueue

replace github.com/kordax/basic-utils/uref => ./uref

replace github.com/kordax/basic-utils/ustr => ./ustr

require (
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
