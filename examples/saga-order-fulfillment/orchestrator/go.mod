module github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator

go 1.24.13

require (
	github.com/janmbaco/go-infrastructure/v2 v2.1.4
	github.com/janmbaco/go-redux/v2 v2.0.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/janmbaco/copier v1.0.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/janmbaco/go-redux/v2 => ../../..
