module htdvisser.dev/exp/echo

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

replace htdvisser.dev/exp/stringslice => ../stringslice

go 1.13

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.2.0-java
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/lyft/protoc-gen-star v0.4.12 // indirect
	github.com/mdempsky/unconvert v0.0.0-20190921185256-3ecd357795af
	github.com/prometheus/procfs v0.0.7 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	go.opencensus.io v0.22.2 // indirect
	golang.org/x/net v0.0.0-20191119073136-fc4aabc6c914 // indirect
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	golang.org/x/tools v0.0.0-20191122232904-2a6ccf25d769 // indirect
	google.golang.org/genproto v0.0.0-20191115221424-83cc0476cb11
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
	htdvisser.dev/exp/backbone v0.0.0-20191108224210-fce53b940d78
	htdvisser.dev/exp/clicontext v0.0.0-20191108224210-fce53b940d78
	htdvisser.dev/exp/flagenv v0.0.0-20191108224210-fce53b940d78
	htdvisser.dev/exp/stringslice v0.0.0-20191108224210-fce53b940d78 // indirect
	mvdan.cc/gofumpt v0.0.0-20191117124704-dc5fc69fd178
)
