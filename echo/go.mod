module htdvisser.dev/exp/echo

go 1.13

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

replace htdvisser.dev/exp/stringslice => ../stringslice

replace go.opentelemetry.io/otel => go.opentelemetry.io/otel v0.2.0

replace go.opentelemetry.io/otel/exporter/trace/jaeger => go.opentelemetry.io/otel/exporter/trace/jaeger v0.2.0

require (
	github.com/envoyproxy/protoc-gen-validate v0.2.0-java
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/lyft/protoc-gen-star v0.4.14 // indirect
	github.com/mdempsky/unconvert v0.0.0-20190921185256-3ecd357795af
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20191220234730-f13409bbebaf // indirect
	google.golang.org/genproto v0.0.0-20191220175831-5c49e3ecc1c1
	google.golang.org/grpc v1.26.0
	htdvisser.dev/exp/backbone v0.0.0-20191221112745-2bfdd273c983
	htdvisser.dev/exp/clicontext v0.0.0-20191221112745-2bfdd273c983
	htdvisser.dev/exp/flagenv v0.0.0-20191221112745-2bfdd273c983
	htdvisser.dev/exp/pflagenv v0.0.0-20200210170633-61b6379ea10f
	mvdan.cc/gofumpt v0.0.0-20191220113447-b896b372089f
)
