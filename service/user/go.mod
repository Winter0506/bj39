module user

go 1.15

require (
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.0.0-20210620082830-8dc9bf49a1d7
	github.com/asim/go-micro/v3 v3.5.1
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/gomodule/redigo v1.8.5
	github.com/google/uuid v1.1.2 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/kr/pretty v0.2.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/protobuf v1.25.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
