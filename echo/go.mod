module htdvisser.dev/exp/echo

go 1.15

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/iancoleman/strcase v0.1.2 // indirect
	github.com/lyft/protoc-gen-star v0.5.2 // indirect
	github.com/mdempsky/unconvert v0.0.0-20200228143138-95ecdbfc0b5f
	github.com/pires/go-proxyproto v0.3.3
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/tools v0.0.0-20210108195828-e2f9c7f1fc8e // indirect
	google.golang.org/genproto v0.0.0-20210108203827-ffc7fda8c3d7
	google.golang.org/grpc v1.34.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	htdvisser.dev/exp/backbone v0.0.0-20201215110738-6786571a568e
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	mvdan.cc/gofumpt v0.1.0
)
