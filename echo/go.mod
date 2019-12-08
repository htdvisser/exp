module htdvisser.dev/exp/echo

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

replace htdvisser.dev/exp/stringslice => ../stringslice

go 1.13

require (
	github.com/envoyproxy/protoc-gen-validate v0.2.0-java
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/lyft/protoc-gen-star v0.4.12 // indirect
	github.com/mdempsky/unconvert v0.0.0-20190921185256-3ecd357795af
	golang.org/x/tools v0.0.0-20191206204035-259af5ff87bd // indirect
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f
	google.golang.org/grpc v1.25.1
	htdvisser.dev/exp/backbone v0.0.0-20191208171514-69ee5c3600ef
	htdvisser.dev/exp/clicontext v0.0.0-20191208171514-69ee5c3600ef
	htdvisser.dev/exp/flagenv v0.0.0-20191208171514-69ee5c3600ef
	htdvisser.dev/exp/stringslice v0.0.0-00010101000000-000000000000 // indirect
	mvdan.cc/gofumpt v0.0.0-20191129122120-d936fb752cbd
)
