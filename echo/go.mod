module htdvisser.dev/exp/echo

go 1.14

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v0.4.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/lyft/protoc-gen-star v0.4.15 // indirect
	github.com/mdempsky/unconvert v0.0.0-20200228143138-95ecdbfc0b5f
	github.com/pires/go-proxyproto v0.1.3
	github.com/spf13/afero v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20200702044944-0cc1aa72b347 // indirect
	google.golang.org/genproto v0.0.0-20200702021140-07506425bd67
	google.golang.org/grpc v1.30.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	htdvisser.dev/exp/backbone v0.0.0-20200615192925-ba29cadbec9f
	htdvisser.dev/exp/clicontext v0.0.0-20200615192925-ba29cadbec9f
	htdvisser.dev/exp/flagenv v0.0.0-20200615192925-ba29cadbec9f
	htdvisser.dev/exp/pflagenv v0.0.0-20200615192925-ba29cadbec9f
	mvdan.cc/gofumpt v0.0.0-20200627213337-90206bd98491
)
