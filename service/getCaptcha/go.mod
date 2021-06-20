module getCaptcha

go 1.15

require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210512092938-c05353c2d58c // indirect
	github.com/afocus/captcha v0.0.0-20191010092841-4bd1f21c8868
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.0.0-20210617101138-08216ccf3148
	github.com/asim/go-micro/v3 v3.5.1
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.5
	github.com/google/uuid v1.2.0 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/miekg/dns v1.1.42 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
