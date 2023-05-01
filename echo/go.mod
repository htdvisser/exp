module htdvisser.dev/exp/echo

go 1.18

replace htdvisser.dev/exp/backbone => ../backbone

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/envoyproxy/protoc-gen-validate v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.3
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2
	github.com/mdempsky/unconvert v0.0.0-20230125054757-2661c2c99a9b
	github.com/pires/go-proxyproto v0.7.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/genproto v0.0.0-20230331144136-dcfb400f0633
	google.golang.org/grpc v1.54.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0
	google.golang.org/protobuf v1.30.0
	htdvisser.dev/exp/backbone v0.0.0-20230211100859-1e66560f78a2
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	mvdan.cc/gofumpt v0.4.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/golang/glog v1.1.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/improbable-eng/grpc-web v0.15.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lyft/protoc-gen-star/v2 v2.0.3 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rs/cors v1.8.3 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
