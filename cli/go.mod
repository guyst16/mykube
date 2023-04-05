module github.com/guyst16/mykube/cli

go 1.19

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli v1.22.12 // indirect
	github.com/guyst16/mykube/virtualmachine v0.0.0
)

replace github.com/guyst16/mykube/virtualmachine v0.0.0 => ../virtualmachine
