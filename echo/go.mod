module htdvisser.dev/exp/echo

go 1.14

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/iancoleman/strcase v0.1.2 // indirect
	github.com/lyft/protoc-gen-star v0.5.2 // indirect
	github.com/mdempsky/unconvert v0.0.0-20200228143138-95ecdbfc0b5f
	github.com/pires/go-proxyproto v0.2.0
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20201020141929-90550980fc82 // indirect
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.32.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	htdvisser.dev/exp/backbone v0.0.0-20201019161321-60b8f60a16df
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	mvdan.cc/gofumpt v0.0.0-20200927160801-5bfeb2e70dd6
)
