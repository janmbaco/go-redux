module github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator

go 1.26.0

require (
	github.com/janmbaco/go-infrastructure/v2 v2.3.2
	github.com/janmbaco/go-redux/v2 v2.1.0
)

require (
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/janmbaco/copier v1.0.0 // indirect
	golang.org/x/net v0.63.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/grpc v1.83.2 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/janmbaco/go-redux/v2 => ../../..
