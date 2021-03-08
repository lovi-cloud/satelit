module github.com/lovi-cloud/satelit/examples/client

go 1.15

require (
	github.com/lovi-cloud/satelit v0.0.1
	github.com/spf13/cobra v1.1.3
	google.golang.org/grpc v1.31.1
)

replace (
	github.com/lovi-cloud/satelit => ../../
)