module htdvisser.dev/exp/backbone

go 1.12

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.10.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/procfs v0.0.4 // indirect
	github.com/rs/cors v1.7.0 // indirect
	go.opencensus.io v0.22.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
	google.golang.org/grpc v1.23.0
	htdvisser.dev/exp/clicontext v0.0.0-20190828181845-1947ca297e5a
	htdvisser.dev/exp/flagenv v0.0.0-20190828181845-1947ca297e5a
)
