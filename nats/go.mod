module htdvisser.dev/exp/nats

go 1.14

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/pflagenv => ../pflagenv

require (
	github.com/dgryski/go-rendezvous v0.0.0-20200624174652-8d2f3be8b2d9 // indirect
	github.com/go-redis/redis/v8 v8.0.0-beta.5
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0
	github.com/spf13/pflag v1.0.5
	go.opentelemetry.io/otel v0.7.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/protobuf v1.24.0 // indirect
	htdvisser.dev/exp/clicontext v0.0.0-20200615192925-ba29cadbec9f
	htdvisser.dev/exp/pflagenv v0.0.0-20200615192925-ba29cadbec9f
)
