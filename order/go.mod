module github.com/Steadypim/rocket-factory/order

go 1.26.1

replace github.com/Steadypim/rocket-factory/shared => ../shared

require (
	github.com/Steadypim/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.72.2
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/oapi-codegen/runtime v1.4.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
