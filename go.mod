module uggly-server

go 1.15

replace github.com/rendicott/uggly => ../uggly

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/rendicott/uggly v0.0.1
	google.golang.org/grpc v1.34.0
)
