module htdvisser.dev/exp/echo

go 1.16

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.2
	github.com/lyft/protoc-gen-star v0.6.0 // indirect
	github.com/mdempsky/unconvert v0.0.0-20200228143138-95ecdbfc0b5f
	github.com/pires/go-proxyproto v0.6.1
	github.com/spf13/afero v1.7.1 // indirect
	github.com/spf13/pflag v1.0.5
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb
	google.golang.org/grpc v1.44.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.27.1
	htdvisser.dev/exp/backbone v0.0.0-20211202070212-5f8f4790ea84
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	mvdan.cc/gofumpt v0.2.1
)
