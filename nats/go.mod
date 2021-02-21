// Deprecated: This pre-1.0 module will be removed. Switch to htdvisser.dev/exp/natsconfig.
module htdvisser.dev/exp/nats

go 1.15

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/pflagenv => ../pflagenv

replace htdvisser.dev/exp/redis => ../redis

replace htdvisser.dev/exp/tls => ../tls

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-redis/redis/v8 v8.6.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redis v0.0.0-20210110145821-20828ad46ee1
	htdvisser.dev/exp/tls v0.0.0-20210110145821-20828ad46ee1
)
