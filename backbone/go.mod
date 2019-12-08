module htdvisser.dev/exp/backbone

go 1.12

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

replace go.opentelemetry.io/otel => go.opentelemetry.io/otel v0.1.3-0.20191128072033-f25c84f35fef

replace go.opentelemetry.io/otel/exporter/trace/jaeger => go.opentelemetry.io/otel/exporter/trace/jaeger v0.1.3-0.20191128072033-f25c84f35fef

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.0.0-20191202183732-d1d2010b5bee // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/rs/cors v1.7.0 // indirect
	go.opencensus.io v0.22.2
	go.opentelemetry.io/otel v0.2.0
	go.opentelemetry.io/otel/exporter/trace/jaeger v0.1.3-0.20191128072033-f25c84f35fef
	golang.org/x/net v0.0.0-20191207000613-e7e4b65ae663 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab // indirect
	google.golang.org/api v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
	htdvisser.dev/exp/clicontext v0.0.0-20191201145824-f663fbf8747b
	htdvisser.dev/exp/flagenv v0.0.0-20191201145824-f663fbf8747b
)
