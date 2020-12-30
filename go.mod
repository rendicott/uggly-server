module uggly-server

go 1.15

replace github.com/rendicott/uggly => ../uggly

replace github.com/rendicott/uggly-server/siteconfig => ./siteconfig

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/rendicott/uggly v0.0.1
	github.com/rendicott/uggly-server/siteconfig v0.0.0
	google.golang.org/grpc v1.34.0
	gopkg.in/yaml.v2 v2.4.0
)
