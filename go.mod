module github.com/whywaita/satelit

go 1.13

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/goccy/go-yaml v1.4.3
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/whywaita/go-dorado-sdk v0.0.0-20200414010248-c2270c6866aa
	go.uber.org/zap v1.14.1
	google.golang.org/grpc v1.29.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	grpc.go4.org v0.0.0-20170609214715-11d0a25b4919 // indirect
)

replace github.com/whywaita/go-dorado-sdk v0.0.0-20200414010248-c2270c6866aa => ./../go-dorado-sdk
