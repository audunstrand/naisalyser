module github.com/example/test

go 1.20

require (
	github.com/original/module v1.0.0
)

replace github.com/original/module => github.com/fork/module v1.0.1
