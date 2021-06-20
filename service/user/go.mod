module user

go 1.15

require (
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.0.0-20210620082830-8dc9bf49a1d7
	github.com/asim/go-micro/v3 v3.5.1
	github.com/golang/protobuf v1.4.3
	github.com/gomodule/redigo v1.8.5
	github.com/micro/micro/v3 v3.0.0 // indirect
	google.golang.org/protobuf v1.25.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
