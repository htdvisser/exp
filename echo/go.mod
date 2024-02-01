module htdvisser.dev/exp/echo

go 1.20

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.3
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1
	github.com/mdempsky/unconvert v0.0.0-20230907125504-415706980c06
	github.com/pires/go-proxyproto v0.7.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/genproto/googleapis/api v0.0.0-20231212172506-995d672761c0
	google.golang.org/grpc v1.60.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0
	google.golang.org/protobuf v1.32.0
	htdvisser.dev/exp/backbone v0.0.0-20231206185358-cf15410f4841
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	mvdan.cc/gofumpt v0.5.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/glog v1.2.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/improbable-eng/grpc-web v0.15.0 // indirect
	github.com/lyft/protoc-gen-star/v2 v2.0.3 // indirect
	github.com/rs/cors v1.10.1 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
	google.golang.org/genproto v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nhooyr.io/websocket v1.8.10 // indirect
)
