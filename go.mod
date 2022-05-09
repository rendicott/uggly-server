module uggly-server

replace github.com/rendicott/uggly-server/pageconfig => ./pageconfig

replace github.com/rendicott/uggly => ../uggly

go 1.17

require (
	github.com/fsnotify/fsnotify v1.5.1
	github.com/rendicott/uggly v0.0.5
	github.com/rendicott/uggly-server/pageconfig v0.0.0-20220307205624-28837372344c
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20210917161153-d61c044b1678 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220401170504-314d38edb7de // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
