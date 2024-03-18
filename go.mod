module github.com/kordax/basic-utils

go 1.22.0

replace github.com/kordax/basic-utils/array-utils => ./array-utils

replace github.com/kordax/basic-utils/async-utils => ./async-utils

replace github.com/kordax/basic-utils/file-utils => ./file-utils

replace github.com/kordax/basic-utils/map-utils => ./map-utils

replace github.com/kordax/basic-utils/math-utils => ./math-utils

replace github.com/kordax/basic-utils/number => ./number

replace github.com/kordax/basic-utils/opt => ./opt

replace github.com/kordax/basic-utils/queue => ./queue

replace github.com/kordax/basic-utils/ref-utils => ./ref-utils

replace github.com/kordax/basic-utils/str-utils => ./str-utils

require (
	github.com/stretchr/testify v1.9.0
	golang.org/x/exp v0.0.0-20240314144324-c7f7c6466f7f
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
