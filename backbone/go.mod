module htdvisser.dev/exp/backbone

go 1.13

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/flagenv => ../flagenv

replace go.opentelemetry.io/otel => go.opentelemetry.io/otel v0.2.0

replace go.opentelemetry.io/otel/exporter/trace/jaeger => go.opentelemetry.io/otel/exporter/trace/jaeger v0.2.0

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/improbable-eng/grpc-web v0.12.0
	github.com/prometheus/client_golang v1.3.0
	github.com/rs/cors v1.7.0 // indirect
	go.opencensus.io v0.22.2
	go.opentelemetry.io/otel v0.2.0
	go.opentelemetry.io/otel/exporter/trace/jaeger v1.0.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191220220014-0732a990476f // indirect
	google.golang.org/api v0.15.0 // indirect
	google.golang.org/genproto v0.0.0-20191220175831-5c49e3ecc1c1 // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
	htdvisser.dev/exp/clicontext v0.0.0-20191208180355-231e02bfe473
	htdvisser.dev/exp/flagenv v0.0.0-20191208180355-231e02bfe473
)
