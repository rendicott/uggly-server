module uggly-server

replace github.com/rendicott/uggly-server/pageconfig => ./pageconfig

go 1.17

require (
	github.com/rendicott/uggly v0.0.5
	github.com/rendicott/uggly-server/pageconfig v0.0.0-20220307205624-28837372344c
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
