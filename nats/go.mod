module htdvisser.dev/exp/nats

go 1.14

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/pflagenv => ../pflagenv

replace htdvisser.dev/exp/redis => ../redis

replace htdvisser.dev/exp/tls => ../tls

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20200707235045-ab33eee955e0 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/protobuf v1.24.0 // indirect
	htdvisser.dev/exp/clicontext v0.0.0-20200702193537-51825981449e
	htdvisser.dev/exp/pflagenv v0.0.0-20200702193537-51825981449e
	htdvisser.dev/exp/redis v0.0.0-20200708182233-d8e90d0a048b
	htdvisser.dev/exp/tls v0.0.0-20200708182233-d8e90d0a048b
)
