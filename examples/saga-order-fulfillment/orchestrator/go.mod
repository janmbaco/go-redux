module github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator

go 1.24.0

require (
	github.com/janmbaco/go-infrastructure/v2 v2.1.1
	github.com/janmbaco/go-redux/v2 v2.0.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/janmbaco/copier v1.0.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251103181224-f26f9409b101 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/janmbaco/go-redux/v2 => ../../..
