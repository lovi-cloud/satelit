module github.com/whywaita/satelit

go 1.13

require (
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/whywaita/go-dorado-sdk v0.0.0-20200414010248-c2270c6866aa
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/whywaita/go-dorado-sdk v0.0.0-20200414010248-c2270c6866aa => ./../go-dorado-sdk
