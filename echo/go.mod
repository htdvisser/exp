module htdvisser.dev/exp/echo

go 1.15

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v0.5.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/iancoleman/strcase v0.1.3 // indirect
	github.com/lyft/protoc-gen-star v0.5.2 // indirect
	github.com/mdempsky/unconvert v0.0.0-20200228143138-95ecdbfc0b5f
	github.com/pires/go-proxyproto v0.5.0
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.4.2 // indirect
	google.golang.org/genproto v0.0.0-20210401141331-865547bb08e2
	google.golang.org/grpc v1.36.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	htdvisser.dev/exp/backbone v0.0.0-20210326073310-b9d6065f847d
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	mvdan.cc/gofumpt v0.1.1
)
