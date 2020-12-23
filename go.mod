module github.com/lovi-cloud/satelit

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.7
	github.com/goccy/go-yaml v1.8.1
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.8.0 // indirect
	github.com/ory/dockertest/v3 v3.6.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/lovi-cloud/go-dorado-sdk v0.8.9
	github.com/lovi-cloud/go-os-brick v0.1.0
	github.com/whywaita/satelit-isucon v0.0.0-20200923052945-91d6c7429e55
	github.com/whywaita/satelit-isucon/qualify/team v0.0.0-20200923053817-268b6eacd659
	github.com/whywaita/satelit-isucon/sshkey v0.0.0-20200923053300-cf351b450037
	github.com/whywaita/teleskop v0.0.0-20200908054150-a5abdfc311ea
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20200909081042-eff7692f9009 // indirect
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)
