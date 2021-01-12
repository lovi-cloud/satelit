module github.com/lovi-cloud/satelit

go 1.15

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.7
	github.com/goccy/go-yaml v1.8.1
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lovi-cloud/go-dorado-sdk v0.8.9
	github.com/lovi-cloud/go-os-brick v0.1.0
	github.com/lovi-cloud/go-targetd v0.0.0-20210105195521-23c076911343
	github.com/lovi-cloud/teleskop v0.0.1
	github.com/ory/dockertest/v3 v3.6.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
