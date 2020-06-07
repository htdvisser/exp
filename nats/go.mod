module htdvisser.dev/exp/nats

go 1.14

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/pflagenv => ../pflagenv

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.3
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0
	github.com/spf13/pflag v1.0.5
	go.opentelemetry.io/otel v0.6.0 // indirect
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	google.golang.org/protobuf v1.24.0 // indirect
	htdvisser.dev/exp/clicontext v0.0.0-20200522135503-6daabdfc50fa
	htdvisser.dev/exp/pflagenv v0.0.0-20200522135503-6daabdfc50fa
)
